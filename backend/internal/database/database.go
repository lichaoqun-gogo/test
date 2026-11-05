package database

import (
	"fmt"
	"log"
	"os"
	"osm-mngt/internal/config"
	"osm-mngt/internal/models"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// 执行初始化SQL（如果表不存在）
	if err := initDatabase(db); err != nil {
		return nil, fmt.Errorf("failed to init database schema: %w", err)
	}

	return db, nil
}

// initDatabase 读取并执行SQL初始化脚本
func initDatabase(db *gorm.DB) error {
	// 1. 执行主初始化SQL
	if err := executeSQLFile(db, filepath.Join("sql", "init.sql")); err != nil {
		return fmt.Errorf("failed to init database schema: %w", err)
	}

	// 2. 执行审批流迁移SQL（如果存在）
	executeSQLFile(db, filepath.Join("migrations", "005_approval_workflow.sql"))

	// 3. 执行任务表字段迁移SQL（如果存在）
	executeSQLFile(db, filepath.Join("migrations", "006_add_task_fields.sql"))

	// 4. 执行project表category字段迁移SQL（如果存在）
	executeSQLFile(db, filepath.Join("migrations", "007_add_project_category.sql"))

	// 5. 执行缺陷采集台账表（GitHub/Jira 统一模型，upsert 导入）
	executeSQLFile(db, filepath.Join("migrations", "008_create_defect_collection_table.sql"))

	// 6. 缺陷实体/处置流程分离，处置体系按缺陷编号关联（幂等）
	executeSQLFile(db, filepath.Join("migrations", "009_split_defect_entity_flow.sql"))
	migrateDefectEntityFlow(db)

	// 7. 缺陷实体表增加内部处置状态 internal_status，flow 退化为流程明细表（幂等）
	executeSQLFile(db, filepath.Join("migrations", "010_add_internal_status.sql"))
	migrateDefectInternalStatus(db)

	// 8. 内部缺陷精简为风险上报字段；外部缺陷增加运维部门列（幂等）
	executeSQLFile(db, filepath.Join("migrations", "011_internal_defect_simplify.sql"))
	migrateInternalDefectSimplify(db)

	// 9. 外部缺陷增加组织列 org（漏洞情报，如 com.alibaba）（幂等）
	executeSQLFile(db, filepath.Join("migrations", "012_defect_add_org.sql"))

	// 10. 统一文件管理：重建 attachment（参考结构，兼容生命周期流程附件）+ 新建 bom 表（幂等）
	executeSQLFile(db, filepath.Join("migrations", "013_bom_attachment_rebuild.sql"))

	// 11. 组件快照定时清理：history/commit_id_history/project_inner_pro_info 建表 + 字段迁移（幂等）
	executeSQLFile(db, filepath.Join("migrations", "014_component_snapshot_cleanup.sql"))
	migrateComponentSnapshotCleanup(db)

	// 12. 风险上报任务扩展表增加内部缺陷编号（系统自动生成 DF-xxx）
	executeSQLFile(db, filepath.Join("migrations", "015_task_defect_add_defect_no.sql"))

	// 13. 审批流脏节点修正（幂等）：撤回/驳回重提曾把跳过节点重置为 pending 空审批人，启动时自动修正为 skipped
	fixOrphanApprovalNodes(db)

	// 自动迁移 users 表（CAS 用户表）
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate users table: %w", err)
	}

	// 预置 CAS 测试用户及其角色（幂等）
	initUserMockData(db)

	// 自动迁移项目-负责人关联表
	if err := db.AutoMigrate(&models.ProjectOwner{}); err != nil {
		return fmt.Errorf("failed to migrate project_owner table: %w", err)
	}

	// 自动迁移缺陷处置相关表（外部实体 defect / 内部实体 defect_ledger / 流程明细 defect_flow / 处置工单 / 状态历史）
	// 内部处置状态以实体表 internal_status 为准；处置级状态机明细存 defect_flow（按 defect_no 关联）
	if err := db.AutoMigrate(&models.Defect{}, &models.DefectLedger{}, &models.DefectFlow{}, &models.DefectDisposal{}, &models.DefectStatusHistory{}); err != nil {
		return fmt.Errorf("failed to migrate defect tables: %w", err)
	}

	// 初始化部门专员映射表和Mock数据
	initDeptSpecialistMapping(db)

	// 初始化project表Mock数据（用于测试审批流）
	initProjectMockData(db)

	// 初始化项目-负责人（维护人员）测试数据
	initProjectOwnerMockData(db)

	// 初始化缺陷处置测试数据
	initDefectMockData(db)

	// 初始化自研项目与组件依赖测试数据（组件缺陷自动匹配处置人链路）
	initComponentDependencyMockData(db)

	return nil
}

// executeSQLFile 执行指定的SQL文件
func executeSQLFile(db *gorm.DB, sqlPath string) error {
	if _, statErr := os.Stat(sqlPath); statErr != nil {
		// SQL文件不存在，跳过
		return nil
	}

	sqlContent, err := os.ReadFile(sqlPath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file %s: %w", sqlPath, err)
	}

	// 分割SQL语句并逐条执行
	statements := splitSQLStatements(string(sqlContent))
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}
		if err := db.Exec(stmt).Error; err != nil {
			// 忽略已存在的表错误（TiDB兼容）
			if strings.Contains(err.Error(), "already exists") ||
				strings.Contains(err.Error(), "Duplicate key name") ||
				strings.Contains(err.Error(), "Duplicate column") {
				continue
			}
			return fmt.Errorf("failed to execute SQL [%s...]: %w", stmt[:min(50, len(stmt))], err)
		}
	}

	return nil
}

// splitSQLStatements 按分号分割SQL语句
func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder

	for _, line := range strings.Split(sql, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		current.WriteString(line)
		current.WriteString("\n")

		if strings.HasSuffix(trimmed, ";") {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	if rest := strings.TrimSpace(current.String()); rest != "" {
		statements = append(statements, rest)
	}

	return statements
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// initUserMockData 预置 CAS 测试用户及其角色（幂等）
// 统一使用 CAS 用户名（username）作为系统内用户唯一标识，
// 与 JWT user_id、dept_specialist_mapping.specialist_id、approval_node.approver_id、governance_task.applicant_id 保持一致
func initUserMockData(db *gorm.DB) {
	mockUsers := []struct {
		Username           string
		Nickname           string
		Department         string
		Team               string
		Role               string
		CanToggleAdminView bool
	}{
		// 普通员工（发起引入/升级/退出申请）
		{"zhangsan", "张三", "技术部", "前端组", "operator", false},
		// 各部门开源治理专员（node_1 本部门专员复审 / 缺陷研判与分配）
		{"lisi", "李四", "技术部", "前端组", "specialist", false},
		{"wangwu", "王五", "平台研发部", "后端组", "specialist", false},
		{"zhaoliu", "赵六", "运维部", "后端组", "specialist", false},
		{"sunqi", "孙七", "安全部", "设计组", "specialist", false},
		// 各部门项目维护人员（缺陷处置二次分配的实际处理人）
		{"liuyang", "刘洋", "平台研发部", "后端组", "operator", false},
		{"chenming", "陈明", "运维部", "后端组", "operator", false},
		{"qianqi", "钱七", "安全部", "设计组", "operator", false},
		// 管理部门专员（node_3 管理部门评审，开源治理办公室；一个部门可配置多名专员）
		{"zhouba", "周八", "开源治理办公室", "运营组", "specialist", false},
		{"wujiu", "吴九", "开源治理办公室", "运营组", "specialist", false},
		// 系统管理员
		{"admin", "系统管理员", "管理层", "管理层", "admin", true},
	}

	for _, mu := range mockUsers {
		var user models.User
		err := db.Where("username = ?", mu.Username).First(&user).Error
		if err != nil {
			// 不存在则创建
			user = models.User{
				Username:           mu.Username,
				Nickname:           mu.Nickname,
				Department:         mu.Department,
				Team:               mu.Team,
				Role:               mu.Role,
				CanToggleAdminView: mu.CanToggleAdminView,
				Status:             1,
			}
			db.Create(&user)
			continue
		}
		// 已存在则补齐角色/部门/团队（保留 CAS 登录写入的显示信息）
		updates := map[string]interface{}{
			"role":                  mu.Role,
			"department":            mu.Department,
			"team":                  mu.Team,
			"can_toggle_admin_view": mu.CanToggleAdminView,
			"status":                1,
		}
		if user.Nickname == "" {
			updates["nickname"] = mu.Nickname
		}
		db.Model(&user).Updates(updates)
	}
}

// initDeptSpecialistMapping 初始化部门专员映射表和Mock数据
// 一个部门可配置 1 名或多名专员（dept_specialist / usage_dept_specialist / mgmt_dept_specialist 均支持多人）
func initDeptSpecialistMapping(db *gorm.DB) {
	// 创建表（唯一键为 dept+role+specialist，同一部门同一角色可有多名专员）
	sql := `
	CREATE TABLE IF NOT EXISTS ` + "`dept_specialist_mapping`" + ` (
	  ` + "`id`" + ` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
	  ` + "`dept`" + ` VARCHAR(100) NOT NULL COMMENT '部门名称',
	  ` + "`specialist_id`" + ` VARCHAR(50) NOT NULL COMMENT '专员ID',
	  ` + "`specialist_name`" + ` VARCHAR(50) NOT NULL COMMENT '专员姓名',
	  ` + "`specialist_role`" + ` VARCHAR(30) NOT NULL COMMENT '专员角色：dept_specialist-部门开源治理专员，usage_dept_specialist-使用部门专员，mgmt_dept_specialist-管理部门专员',
	  ` + "`is_active`" + ` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否有效：1-有效，0-无效',
	  ` + "`create_time`" + ` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	  ` + "`update_time`" + ` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
	  PRIMARY KEY (` + "`id`" + `),
	  UNIQUE KEY ` + "`uk_dept_role_specialist`" + ` (` + "`dept`" + `, ` + "`specialist_role`" + `, ` + "`specialist_id`" + `),
	  KEY ` + "`idx_specialist_id`" + ` (` + "`specialist_id`" + `)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='部门开源治理专员映射表';
	`
	db.Exec(sql)

	// 唯一键迁移：旧版本为 (dept, specialist_role)（每部门每角色仅 1 人），
	// 调整为 (dept, specialist_role, specialist_id) 以支持一个部门多名专员
	migrateDeptSpecialistUniqueKey(db)

	// 清理已废弃的部门映射（部门体系统一后，"技术研发部"不再使用）
	db.Exec("DELETE FROM `dept_specialist_mapping` WHERE `dept` = '技术研发部'")

	// 插入Mock数据
	// specialist_id 统一使用 CAS 用户名（username），与 JWT user_id、审批节点 approver_id 保持一致
	mockData := []struct {
		Dept           string
		SpecialistID   string
		SpecialistName string
		SpecialistRole string
	}{
		// 部门开源治理专员（node_1 本部门专员复审；技术部配置 2 名，验证多专员场景）
		{"技术部", "lisi", "李四", "dept_specialist"},
		{"技术部", "zhengjie", "郑杰", "dept_specialist"},
		{"平台研发部", "wangwu", "王五", "dept_specialist"},
		{"运维部", "zhaoliu", "赵六", "dept_specialist"},
		{"安全部", "sunqi", "孙七", "dept_specialist"},
		// 使用部门专员（node_2 使用部门确认）
		{"技术部", "lisi", "李四", "usage_dept_specialist"},
		{"技术部", "zhengjie", "郑杰", "usage_dept_specialist"},
		{"平台研发部", "wangwu", "王五", "usage_dept_specialist"},
		{"运维部", "zhaoliu", "赵六", "usage_dept_specialist"},
		{"安全部", "sunqi", "孙七", "usage_dept_specialist"},
		// 管理部门专员（node_3 管理部门评审，开源治理办公室；配置 2 名，任一在岗专员均可评审）
		{"开源治理办公室", "zhouba", "周八", "mgmt_dept_specialist"},
		{"开源治理办公室", "wujiu", "吴九", "mgmt_dept_specialist"},
	}

	for _, data := range mockData {
		db.Exec(`
			INSERT INTO `+"`dept_specialist_mapping`"+`
			(`+"`dept`"+`, `+"`specialist_id`"+`, `+"`specialist_name`"+`, `+"`specialist_role`"+`, `+"`is_active`"+`)
			VALUES (?, ?, ?, ?, 1)
			ON DUPLICATE KEY UPDATE
				`+"`specialist_name`"+` = VALUES(`+"`specialist_name`"+`),
				`+"`is_active`"+` = 1
		`, data.Dept, data.SpecialistID, data.SpecialistName, data.SpecialistRole)
	}
}

// migrateDeptSpecialistUniqueKey 迁移 dept_specialist_mapping 唯一键：
// uk_dept_role (dept, specialist_role) → uk_dept_role_specialist (dept, specialist_role, specialist_id)
func migrateDeptSpecialistUniqueKey(db *gorm.DB) {
	const table = "`dept_specialist_mapping`"

	var oldIdx int64
	db.Raw("SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dept_specialist_mapping' AND INDEX_NAME = 'uk_dept_role'").Scan(&oldIdx)
	if oldIdx > 0 {
		if err := db.Exec("ALTER TABLE " + table + " DROP INDEX `uk_dept_role`").Error; err != nil {
			log.Printf("[migrateDeptSpecialistUniqueKey] 删除旧唯一键失败: %v", err)
		}
	}

	var newIdx int64
	db.Raw("SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dept_specialist_mapping' AND INDEX_NAME = 'uk_dept_role_specialist'").Scan(&newIdx)
	if newIdx == 0 {
		if err := db.Exec("ALTER TABLE " + table + " ADD UNIQUE KEY `uk_dept_role_specialist` (`dept`, `specialist_role`, `specialist_id`)").Error; err != nil {
			log.Printf("[migrateDeptSpecialistUniqueKey] 创建新唯一键失败: %v", err)
		}
	}
}

// initProjectMockData 初始化project表Mock数据（用于测试审批流）
// 运维部门 maint_dept / 使用部门 use_dept 与 dept_specialist_mapping 的专员映射保持一致，
// 确保升级/退出申请的 node_2 使用部门确认节点能正确解析到 CAS 专员
func initProjectMockData(db *gorm.DB) {
	// 名称、运维部门、使用部门、许可证、描述
	ensureProject(db, "kafka", "技术部", "平台研发部", "Apache-2.0", "Apache Kafka分布式流处理平台")
	ensureProject(db, "redis", "技术部", "技术部", "BSD-3-Clause", "Redis内存数据库")
	ensureProject(db, "spring-boot", "平台研发部", "技术部", "Apache-2.0", "Spring Boot应用框架")
	ensureProject(db, "nginx", "运维部", "安全部", "BSD-2-Clause", "Nginx Web服务器")
	ensureProject(db, "node.js", "技术部", "技术部", "MIT", "Node.js运行时")
}

// initDefectMockData 初始化缺陷处置测试数据（幂等）
//  1. 确保测试用开源软件存在并具备正确的运维/使用部门
//  2. 插入缺陷/漏洞测试数据，覆盖情报源（vuln-intel/sca-scan/user-report）、
//     缺陷/漏洞两条处置路线、不同运维部门与不同状态
func initDefectMockData(db *gorm.DB) {
	// 确保项目数据（运维部门 + 使用部门），用于权限过滤与裂变
	ensureProject(db, "redis", "技术部", "技术部", "BSD-3-Clause", "Redis内存数据库")
	ensureProject(db, "node.js", "技术部", "技术部", "MIT", "Node.js运行时")
	ensureProject(db, "spring-boot", "平台研发部", "技术部", "Apache-2.0", "Spring Boot应用框架")
	ensureProject(db, "nginx", "运维部", "安全部", "BSD-2-Clause", "Nginx Web服务器")

	defects := []struct {
		DefectNo        string
		FlowType        string
		Source          string
		Title           string
		RiskType        string
		RiskLevel       string
		Component       string
		AffectedVersion string
		CVEID           string
		Description     string
		MaintDept       string
		Status          string
	}{
		// 缺陷：研判后必须上会（defect 路线）
		{"DF-TEST-001", "defect", "vuln-intel", "Redis 远程代码执行漏洞（Lua 脚本）", "vulnerability", "critical", "redis", "7.0.x", "CVE-2024-31449", "Redis Lua 脚本引擎存在远程代码执行漏洞，攻击者可构造恶意脚本绕过沙箱执行任意代码。", "技术部", "pending_analysis"},
		{"DF-TEST-002", "defect", "sca-scan", "Redis 数据持久化配置缺陷", "vulnerability", "high", "redis", "6.2.0", "", "SCA 扫描发现 Redis 持久化配置缺陷，可能导致异常宕机时数据丢失。", "技术部", "pending_review"},
		// 漏洞：可直接研判 need_fix 立即裂变（vuln 路线）
		{"DF-TEST-003", "vuln", "vuln-intel", "Node.js 原型污染漏洞", "vulnerability", "high", "node.js", "20.x", "CVE-2023-30581", "Node.js 运行时存在原型污染漏洞，恶意输入可污染 Object 原型导致拒绝服务或权限绕过。", "技术部", "pending_analysis"},
		{"DF-TEST-004", "vuln", "vuln-intel", "Spring Boot Actuator 配置泄露漏洞", "vulnerability", "medium", "spring-boot", "3.2.x", "CVE-2024-22243", "Spring Boot Actuator 端点配置不当导致敏感配置信息泄露。", "平台研发部", "pending_analysis"},
		// 用户上报：运维部门为运维部，使用部门为安全部（验证跨部门裂变）
		{"DF-TEST-005", "defect", "user-report", "Nginx 日志路径穿越缺陷", "vulnerability", "medium", "nginx", "1.24.x", "", "用户上报 Nginx 日志配置存在路径穿越缺陷，可能覆盖任意文件。", "运维部", "pending_analysis"},
	}

	for _, d := range defects {
		// 内部缺陷实体（风险上报那套字段：编号/标题/等级选填/描述 + 影响组件 + 内部处置状态；upsert 不覆盖状态，避免重启把已闭环等状态打回）
		db.Exec(`
			INSERT INTO `+"`defect_ledger`"+`
			(`+"`defect_no`"+`, `+"`source`"+`, `+"`title`"+`, `+"`risk_level`"+`, `+"`description`"+`, `+"`affected_component`"+`, `+"`discovered_date`"+`, `+"`internal_status`"+`, `+"`create_time`"+`, `+"`update_time`"+`)
			VALUES (?, ?, ?, ?, ?, ?, NOW(), ?, NOW(), NOW())
			ON DUPLICATE KEY UPDATE
				`+"`source`"+` = VALUES(`+"`source`"+`),
				`+"`description`"+` = VALUES(`+"`description`"+`),
				`+"`affected_component`"+` = COALESCE(`+"`affected_component`"+`, VALUES(`+"`affected_component`"+`)),
				`+"`update_time`"+` = NOW()
		`, d.DefectNo, d.Source, d.Title, d.RiskLevel, d.Description, d.Component, d.Status)
		// 缺陷级处置流程明细（研判/评审/派发等字段；状态在实体表，不再存 flow）
		db.Exec(`
			INSERT INTO `+"`defect_flow`"+`
			(`+"`defect_no`"+`, `+"`flow_type`"+`, `+"`maint_dept`"+`, `+"`create_time`"+`, `+"`update_time`"+`)
			VALUES (?, ?, ?, NOW(), NOW())
			ON DUPLICATE KEY UPDATE
				`+"`flow_type`"+` = VALUES(`+"`flow_type`"+`),
				`+"`maint_dept`"+` = VALUES(`+"`maint_dept`"+`),
				`+"`update_time`"+` = NOW()
		`, d.DefectNo, d.FlowType, d.MaintDept)
	}
}

// initComponentDependencyMockData 初始化自研项目与组件依赖测试数据（幂等）
// 支撑"组件缺陷 → 查依赖项目 → 自动匹配项目负责人"的处置链路：
//   - 自研项目（project_inner）：订单中心（技术部）、支付平台（平台研发部）
//   - 组件依赖（component_snapshot）：fastjson 被两个自研项目依赖
//   - 项目负责人（project_owner，inner）：订单中心→张三（主岗）、支付平台→刘洋（主岗）
//   - 组件缺陷：DF-TEST-006 fastjson 反序列化漏洞（affected_component_type='component'）
func initComponentDependencyMockData(db *gorm.DB) {
	type innerDef struct {
		ProCode   string
		GroupName string
		DeptName  string
		Owner     string // 项目主岗（CAS 用户名）
		Team      string
	}
	inners := []innerDef{
		{"ORD-CENTER", "订单组", "技术部", "zhangsan", "前端组"},
		{"PAY-PLATFORM", "支付组", "平台研发部", "liuyang", "后端组"},
	}

	now := time.Now()
	projectIDByCode := map[string]int64{}
	for _, p := range inners {
		var pi models.ProjectInner
		err := db.Where("pro_code = ?", p.ProCode).First(&pi).Error
		if err != nil {
			pi = models.ProjectInner{
				ProCode:    p.ProCode,
				Active:     intPtr(1),
				Type:       "inner",
				GroupName:  p.GroupName,
				DeptName:   p.DeptName,
				CreateTime: now,
				UpdateTime: now,
			}
			if err := db.Create(&pi).Error; err != nil {
				continue
			}
		}
		projectIDByCode[p.ProCode] = pi.ID

		// 项目负责人（inner），第一个为主岗
		var count int64
		db.Model(&models.ProjectOwner{}).
			Where("project_type = ? AND project_id = ? AND owner_username = ?", "inner", pi.ID, p.Owner).
			Count(&count)
		if count == 0 {
			db.Create(&models.ProjectOwner{
				ProjectType:   "inner",
				ProjectID:     pi.ID,
				OwnerUsername: p.Owner,
				DeptName:      p.DeptName,
				Team:          p.Team,
				AssignedBy:    "admin",
				AssignedAt:    now,
			})
		}
	}

	// 组件依赖：fastjson 被两个自研项目依赖
	for _, code := range []string{"ORD-CENTER", "PAY-PLATFORM"} {
		pid, ok := projectIDByCode[code]
		if !ok {
			continue
		}
		var count int64
		db.Model(&models.ComponentSnapshot{}).
			Where("_name = ? AND PROJECT_ID = ?", "fastjson", pid).
			Count(&count)
		if count == 0 {
			db.Create(&models.ComponentSnapshot{
				Name:      "fastjson",
				Version:   "1.2.24",
				CurrentID: "",
				ProjectID: &pid,
			})
		}
	}

	// 组件缺陷：fastjson（内部缺陷按风险上报字段存储 + 影响组件供自动分配裂变）
	db.Exec(`
		INSERT INTO `+"`defect_ledger`"+`
		(`+"`defect_no`"+`, `+"`source`"+`, `+"`title`"+`, `+"`risk_level`"+`, `+"`description`"+`, `+"`affected_component`"+`, `+"`discovered_date`"+`, `+"`internal_status`"+`, `+"`create_time`"+`, `+"`update_time`"+`)
		VALUES ('DF-TEST-006', 'sca-scan', 'fastjson 反序列化远程代码执行漏洞', 'critical', 'fastjson 存在反序列化远程代码执行漏洞，被订单中心、支付平台两个自研项目依赖。', 'fastjson', NOW(), 'pending_analysis', NOW(), NOW())
		ON DUPLICATE KEY UPDATE `+"`affected_component`"+` = COALESCE(`+"`affected_component`"+`, 'fastjson'), `+"`update_time`"+` = NOW()
	`)
	db.Exec(`
		INSERT INTO `+"`defect_flow`"+`
		(`+"`defect_no`"+`, `+"`flow_type`"+`, `+"`maint_dept`"+`, `+"`create_time`"+`, `+"`update_time`"+`)
		VALUES ('DF-TEST-006', 'defect', '技术部', NOW(), NOW())
		ON DUPLICATE KEY UPDATE `+"`update_time`"+` = NOW()
	`)
}

func intPtr(i int) *int { return &i }

// colExists 检查表是否存在指定列
func colExists(db *gorm.DB, table, column string) bool {
	var cnt int64
	db.Raw("SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?",
		table, column).Scan(&cnt)
	return cnt > 0
}

// migrateDefectEntityFlow 缺陷实体/处置流程分离数据迁移（幂等，information_schema 守卫）：
//   - 旧版 defect_ledger 混表（实体+流程） -> 流程字段回填到 defect_flow，实体表纯化
//   - defect_disposal / defect_status_history 由主键关联(defect_id) 迁移为按缺陷编号(defect_no) 关联
// 说明：不用 PREPARE 会话变量做守卫，避免 gorm 连接池跨语句丢变量的坑
func migrateDefectEntityFlow(db *gorm.DB) {
	// 1) 回填 defect_flow（旧表仍有流程列时）
	if colExists(db, "defect_ledger", "status") {
		db.Exec(`
			INSERT INTO ` + "`defect_flow`" + `
			(` + "`defect_no`,`flow_type`,`maint_dept`,`deadline`,`status`,`analyze_result`,`review_meeting_no`,`review_conclusion`,`dispatch_type`,`dispatch_target`,`dispatch_dept`,`dispatch_time`,`dispatcher`,`analyze_by`,`analyze_time`,`review_by`,`review_time`,`close_time`,`create_time`,`update_time`" + `)
			SELECT ` + "`defect_no`,`flow_type`,`maint_dept`,`deadline`,`status`,`analyze_result`,`review_meeting_no`,`review_conclusion`,`dispatch_type`,`dispatch_target`,`dispatch_dept`,`dispatch_time`,`dispatcher`,`analyze_by`,`analyze_time`,`review_by`,`review_time`,`close_time`,`create_time`,`update_time`" + `
			FROM ` + "`defect_ledger`" + `
			ON DUPLICATE KEY UPDATE ` + "`defect_no`" + ` = VALUES(` + "`defect_no`" + `)
		`)
	}
	// 2) 回填处置单/状态历史 defect_no（旧表有 defect_id 时按编号回填）
	if colExists(db, "defect_disposal", "defect_id") {
		db.Exec("UPDATE `defect_disposal` dd JOIN `defect_ledger` dl ON dl.id = dd.defect_id SET dd.defect_no = dl.defect_no")
	}
	if colExists(db, "defect_status_history", "defect_id") {
		db.Exec("UPDATE `defect_status_history` sh JOIN `defect_ledger` dl ON dl.id = sh.defect_id SET sh.defect_no = dl.defect_no")
	}
	// 3) 删除旧主键关联列
	if colExists(db, "defect_disposal", "defect_id") {
		db.Exec("ALTER TABLE `defect_disposal` DROP COLUMN `defect_id`")
	}
	if colExists(db, "defect_status_history", "defect_id") {
		db.Exec("ALTER TABLE `defect_status_history` DROP COLUMN `defect_id`")
	}
	// 4) defect_ledger 删除流程字段，纯化为内部缺陷实体表
	if colExists(db, "defect_ledger", "status") {
		db.Exec("ALTER TABLE `defect_ledger` DROP COLUMN `status`, DROP COLUMN `flow_type`, DROP COLUMN `maint_dept`, DROP COLUMN `deadline`, DROP COLUMN `analyze_result`, DROP COLUMN `review_meeting_no`, DROP COLUMN `review_conclusion`, DROP COLUMN `dispatch_type`, DROP COLUMN `dispatch_target`, DROP COLUMN `dispatch_dept`, DROP COLUMN `dispatch_time`, DROP COLUMN `dispatcher`, DROP COLUMN `analyze_by`, DROP COLUMN `analyze_time`, DROP COLUMN `review_by`, DROP COLUMN `review_time`, DROP COLUMN `close_time`, DROP COLUMN `related_task_id`")
	}
}

// migrateDefectInternalStatus 内部处置状态迁移（幂等）：
// 状态以实体表 internal_status 为准；将旧 defect_flow.status 回填到实体表后删除 flow.status 列
func migrateDefectInternalStatus(db *gorm.DB) {
	// 1) 回填：defect_flow.status -> defect_ledger.internal_status / defect.internal_status（旧表仍有 status 时）
	if colExists(db, "defect_flow", "status") {
		db.Exec("UPDATE `defect_ledger` dl JOIN `defect_flow` df ON df.defect_no = dl.defect_no SET dl.internal_status = df.status")
		db.Exec("UPDATE `defect` d JOIN `defect_flow` df ON df.defect_no = d.defect_id SET d.internal_status = df.status")
	}
	// 2) 删除 defect_flow.status 列（flow 退化为流程明细表）
	if colExists(db, "defect_flow", "status") {
		db.Exec("ALTER TABLE `defect_flow` DROP COLUMN `status`")
	}
}

// migrateInternalDefectSimplify 内部缺陷精简迁移（幂等）：
// 内部缺陷去掉结构字段（risk_type/affected_component_type/affected_version/cve_id），按风险上报那套字段来；
// 但保留 affected_component（漏洞自动分配需按影响组件裂变，如 redis/fastjson），缺列时补回
func migrateInternalDefectSimplify(db *gorm.DB) {
	for _, col := range []string{"risk_type", "affected_component_type", "affected_version", "cve_id"} {
		if colExists(db, "defect_ledger", col) {
			db.Exec("ALTER TABLE `defect_ledger` DROP COLUMN `" + col + "`")
		}
	}
	if !colExists(db, "defect_ledger", "affected_component") {
		db.Exec("ALTER TABLE `defect_ledger` ADD COLUMN `affected_component` varchar(128) DEFAULT NULL COMMENT '影响软件/组件（漏洞自动分配裂变依据）' AFTER `description`")
	}
}

// fixOrphanApprovalNodes 修正审批流脏节点：
// 历史上"撤回/驳回重提"曾把跳过的节点重置为 pending 且无审批人，导致流程卡住。
// 启动时将 status=pending 且无审批人的非发起节点修正为 skipped（幂等）。
func fixOrphanApprovalNodes(db *gorm.DB) {
	res := db.Exec("UPDATE approval_node SET status='skipped', is_skipped=1, skip_reason='历史脏节点修正：无审批人，自动跳过' WHERE node_order > 0 AND status='pending' AND is_skipped=0 AND (approver_id='' OR approver_id IS NULL)")
	if res.Error != nil {
		log.Printf("[fixOrphanApprovalNodes] 修正脏节点失败: %v", res.Error)
		return
	}
	if res.RowsAffected > 0 {
		log.Printf("[fixOrphanApprovalNodes] 已修正 %d 个无审批人的脏节点为 skipped", res.RowsAffected)
	}
}

// migrateComponentSnapshotCleanup 组件快照清理字段迁移（幂等）：
// component_snapshot 增加逻辑删除标识 is_delete（默认 1-未删除，归档后置 0-已删除）与入库时间 create_time
// （1.2 待清理 SQL 依赖 MAX(create_time)）
func migrateComponentSnapshotCleanup(db *gorm.DB) {
	if !colExists(db, "component_snapshot", "is_delete") {
		db.Exec("ALTER TABLE `component_snapshot` ADD COLUMN `is_delete` tinyint NOT NULL DEFAULT 1 COMMENT '逻辑删除标识: 0-已删除, 1-未删除' AFTER `PROJECT_ID`")
	}
	if !colExists(db, "component_snapshot", "create_time") {
		db.Exec("ALTER TABLE `component_snapshot` ADD COLUMN `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间' AFTER `is_delete`")
	}
}

// initProjectOwnerMockData 初始化项目-负责人（维护人员）测试数据（幂等）
// 对应“各部门负责项目对应的维护人员”映射，供人员管理、缺陷处置二次分配等场景使用
func initProjectOwnerMockData(db *gorm.DB) {
	// 软件名 -> 维护人员列表（owner_username 统一使用 CAS 用户名）
	type ownerDef struct {
		Username string
		Dept     string
		Team     string
	}
	projectOwners := map[string][]ownerDef{
		"kafka": {
			{"zhangsan", "技术部", "前端组"},  // 运维部门维护人
			{"liuyang", "平台研发部", "后端组"}, // 使用部门维护人
		},
		"redis": {
			{"zhangsan", "技术部", "前端组"},
		},
		"spring-boot": {
			{"liuyang", "平台研发部", "后端组"}, // 运维部门维护人
			{"zhangsan", "技术部", "前端组"},  // 使用部门维护人
		},
		"nginx": {
			{"chenming", "运维部", "后端组"}, // 运维部门维护人
			{"qianqi", "安全部", "设计组"},   // 使用部门维护人
		},
		"node.js": {
			{"zhangsan", "技术部", "前端组"},
		},
	}

	for name, owners := range projectOwners {
		var p models.Project
		if err := db.Where("_name = ?", name).First(&p).Error; err != nil {
			continue
		}
		for _, o := range owners {
			var count int64
			db.Model(&models.ProjectOwner{}).
				Where("project_type = ? AND project_id = ? AND owner_username = ?", "software", p.ID, o.Username).
				Count(&count)
			if count > 0 {
				continue
			}
			db.Create(&models.ProjectOwner{
				ProjectType:   "software",
				ProjectID:     p.ID,
				OwnerUsername: o.Username,
				DeptName:      o.Dept,
				Team:          o.Team,
				AssignedBy:    "admin",
				AssignedAt:    time.Now(),
			})
		}
	}
}

// ensureProject 确保测试用开源软件存在，并补齐运维/使用部门
func ensureProject(db *gorm.DB, name, maintDept, useDept, license, desc string) {
	var p models.Project
	if err := db.Where("_name = ?", name).First(&p).Error; err == nil {
		db.Model(&models.Project{}).Where("_name = ?", name).Updates(map[string]interface{}{
			"maint_dept": maintDept,
			"use_dept":   useDept,
		})
		return
	}
	db.Exec(`
		INSERT INTO `+"`project`"+`
		(`+"`_name`"+`, `+"`_version`"+`, `+"`UUID`"+`, `+"`type`"+`, `+"`applicable_level`"+`, `+"`maint_dept`"+`, `+"`use_dept`"+`, `+"`country`"+`, `+"`category`"+`, `+"`license`"+`, `+"`DESCRIPTION`"+`, `+"`create_time`"+`, `+"`update_time`"+`)
		VALUES (?, '1.0', UUID(), 'java', 'B', ?, ?, '中国', 'software', ?, ?, NOW(), NOW())
	`, name, maintDept, useDept, license, desc)
}

// ===== 菜单数据初始化已废弃 =====
// 菜单统一由「统一权限系统」按角色下发（见 internal/service/permission/remote.go），
// 不再写入 menus / role_menus 表（废弃菜单存表逻辑）。

	// ===== admin 角色菜单（完整管理菜单树） =====

