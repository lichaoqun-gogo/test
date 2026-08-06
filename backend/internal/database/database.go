package database

import (
	"fmt"
	"os"
	"osm-mngt/internal/config"
	"osm-mngt/internal/models"
	"path/filepath"
	"strings"

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
		Logger:                 logger.Default.LogMode(logger.Info),
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
	
	// 自动迁移 users 表（CAS 用户表）
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate users table: %w", err)
	}

	// 自动迁移菜单权限表
	if err := db.AutoMigrate(&models.Menu{}, &models.RoleMenu{}); err != nil {
		return fmt.Errorf("failed to migrate menus table: %w", err)
	}

	// 自动迁移项目-负责人关联表
	if err := db.AutoMigrate(&models.ProjectOwner{}); err != nil {
		return fmt.Errorf("failed to migrate project_owner table: %w", err)
	}

	// 初始化部门专员映射表和Mock数据
	initDeptSpecialistMapping(db)

	// 初始化project表Mock数据（用于测试审批流）
	initProjectMockData(db)

	// 初始化菜单及角色菜单关联数据
	initMenuData(db)

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

// initDeptSpecialistMapping 初始化部门专员映射表和Mock数据
func initDeptSpecialistMapping(db *gorm.DB) {
	// 创建表
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
	  UNIQUE KEY ` + "`uk_dept_role`" + ` (` + "`dept`" + `, ` + "`specialist_role`" + `),
	  KEY ` + "`idx_specialist_id`" + ` (` + "`specialist_id`" + `)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='部门开源治理专员映射表';
	`
	db.Exec(sql)
	
	// 插入Mock数据
	mockData := []struct {
		Dept         string
		SpecialistID string
		SpecialistName string
		SpecialistRole string
	}{
		// 部门开源治理专员
		{"技术研发部", "user001", "李四", "dept_specialist"},
		{"平台研发部", "user002", "王五", "dept_specialist"},
		{"运维部", "user003", "赵六", "dept_specialist"},
		{"安全部", "user004", "钱七", "dept_specialist"},
		// 使用部门专员（测试用例2：使用部门确认 - 王二，运维部专员）
		{"技术研发部", "user001", "李四", "usage_dept_specialist"},
		{"平台研发部", "user002", "王五", "usage_dept_specialist"},
		{"运维部", "user007", "王二", "usage_dept_specialist"},
		{"安全部", "user004", "钱七", "usage_dept_specialist"},
		// 管理部门专员
		{"开源治理办公室", "user010", "孙八", "mgmt_dept_specialist"},
		{"开源治理办公室", "user011", "周九", "mgmt_dept_specialist"},
	}
	
	for _, data := range mockData {
		db.Exec(`
			INSERT INTO ` + "`dept_specialist_mapping`" + ` 
			(` + "`dept`" + `, ` + "`specialist_id`" + `, ` + "`specialist_name`" + `, ` + "`specialist_role`" + `, ` + "`is_active`" + `)
			VALUES (?, ?, ?, ?, 1)
			ON DUPLICATE KEY UPDATE 
				` + "`specialist_name`" + ` = VALUES(` + "`specialist_name`" + `),
				` + "`is_active`" + ` = 1
		`, data.Dept, data.SpecialistID, data.SpecialistName, data.SpecialistRole)
	}
}

// initProjectMockData 初始化project表Mock数据（用于测试审批流）
func initProjectMockData(db *gorm.DB) {
	mockProjects := []struct {
		Name        string
		UseDept     string
		License     string
		LicenseRisk string
		Status      string
		Description string
	}{
		{"kafka", "平台研发部", "Apache-2.0", "low", "active", "Apache Kafka分布式流处理平台"},
		{"redis", "技术研发部", "BSD-3-Clause", "low", "active", "Redis内存数据库"},
		{"spring-boot", "运维部", "Apache-2.0", "low", "active", "Spring Boot应用框架"},
		{"nginx", "安全部", "BSD-2-Clause", "low", "active", "Nginx Web服务器"},
	}
	
	for _, p := range mockProjects {
		db.Exec(`
			INSERT INTO `+"`project`"+` 
			(`+"`_name`"+`, `+"`use_dept`"+`, `+"`license`"+`, `+"`license_risk_level`"+`, `+"`status`"+`, `+"`description`"+`, `+"`create_time`"+`, `+"`update_time`"+`)
			VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
			ON DUPLICATE KEY UPDATE 
				`+"`use_dept`"+` = VALUES(`+"`use_dept`"+`),
				`+"`license`"+` = VALUES(`+"`license`"+`),
				`+"`license_risk_level`"+` = VALUES(`+"`license_risk_level`"+`),
				`+"`status`"+` = VALUES(`+"`status`"+`),
				`+"`description`"+` = VALUES(`+"`description`"+`)
		`, p.Name, p.UseDept, p.License, p.LicenseRisk, p.Status, p.Description)
	}
}

// menuDef 菜单定义（用于初始化数据，支持嵌套子菜单）
type menuDef struct {
	Name      string
	Path      string
	Component string
	Icon      string
	Sort      int
	MenuType  int
	Hidden    bool
	Group     string
	Children  []menuDef
}

// createMenuTree 递归创建菜单及角色-菜单关联（幂等，存在则更新）
func createMenuTree(db *gorm.DB, defs []menuDef, parentID uint, role string) {
	for _, d := range defs {
		m := models.Menu{
			ParentID:  parentID,
			Name:      d.Name,
			Path:      d.Path,
			Component: d.Component,
			Icon:      d.Icon,
			Sort:      d.Sort,
			MenuType:  d.MenuType,
			Hidden:    d.Hidden,
			Group:     d.Group,
		}
		// 先查找，存在则更新，不存在则创建
		var existing models.Menu
		if db.Where("name = ? AND parent_id = ? AND path = ?", d.Name, parentID, d.Path).First(&existing).Error == nil {
			// 更新 icon、component、sort 等字段
			db.Model(&existing).Updates(map[string]interface{}{
				"icon":      d.Icon,
				"component": d.Component,
				"sort":      d.Sort,
				"menu_type": d.MenuType,
				"hidden":    d.Hidden,
				"group":     d.Group,
			})
			m.ID = existing.ID
		} else {
			db.Create(&m)
		}
		// 建立角色-菜单关联
		db.FirstOrCreate(&models.RoleMenu{}, models.RoleMenu{Role: role, MenuID: m.ID})
		// 递归创建子菜单
		createMenuTree(db, d.Children, m.ID, role)
	}
}

// initMenuData 初始化菜单及角色菜单关联数据
func initMenuData(db *gorm.DB) {
	// ===== admin 角色菜单（完整管理菜单树） =====
	adminMenus := []menuDef{
		{
			Name: "台账管理", MenuType: 1, Group: "asset", Sort: 1, Icon: "Coin",
			Children: []menuDef{
				{Name: "开源基础软件", Path: "/admin/asset/software-library", Component: "views/admin/asset/SoftwareLibrary", MenuType: 2, Sort: 1},
				{Name: "开源依赖包", Path: "/admin/asset/dependency-library", Component: "views/admin/asset/DependencyLibrary", MenuType: 2, Sort: 2},
				{Name: "开源工具", Path: "/admin/asset/tool-library", Component: "views/admin/asset/ToolLibrary", MenuType: 2, Sort: 3},
				{Name: "自研项目", Path: "/admin/asset/self-developed", Component: "views/admin/asset/SelfDevelopedProjects", MenuType: 2, Sort: 4},
				{Name: "源码备份", Path: "/admin/asset/source-backup", Component: "views/admin/asset/SourceBackup", MenuType: 2, Sort: 5},
			},
		},
		{Name: "可信仓库", Path: "/admin/trusted-repo", Component: "views/admin/asset/TrustedRepository", Icon: "Files", MenuType: 2, Sort: 2},
		{
			Name: "风险管理", MenuType: 1, Group: "risk", Sort: 3, Icon: "Warning",
			Children: []menuDef{
				{Name: "名单管理", Path: "/admin/risk/blacklist", Component: "views/admin/risk/Blacklist", MenuType: 2, Sort: 1},
				{Name: "CICD检查", Path: "/admin/risk/cicd-check", Component: "views/admin/risk/CicdCheck", MenuType: 2, Sort: 2},
				{Name: "许可证信息库", Path: "/admin/risk/license-library", Component: "views/admin/risk/LicenseLibrary", MenuType: 2, Sort: 3},
				{Name: "缺陷台账", Path: "/admin/risk/defect-ledger", Component: "views/admin/risk/DefectLedger", MenuType: 2, Sort: 4},
				{Name: "活跃度健康度监测", Path: "/admin/risk/activity-monitor", Component: "views/admin/risk/ActivityMonitor", MenuType: 2, Sort: 5},
				{Name: "SCA自动化解析入库", Path: "/admin/risk/sca-parser", Component: "views/admin/risk/ScaParser", MenuType: 2, Sort: 6},
			},
		},
		{
			Name: "治理中心", MenuType: 1, Group: "governance", Sort: 4, Icon: "Setting",
			Children: []menuDef{
				{Name: "任务列表", Path: "/admin/governance/task-list", Component: "views/admin/governance/TaskList", MenuType: 2, Sort: 1},
			},
		},
		{
			Name: "决策与洞察中心", MenuType: 1, Group: "decision", Sort: 5, Icon: "DataAnalysis",
			Children: []menuDef{
				{
					Name: "治理效能看板", MenuType: 1, Sort: 1,
					Children: []menuDef{
						{Name: "违规拦截统计", Path: "/admin/decision/interception", Component: "views/admin/decision/interception", MenuType: 2, Sort: 1},
						{Name: "缺陷修复效能", Path: "/admin/decision/repair-efficiency", Component: "views/admin/decision/repair-efficiency", MenuType: 2, Sort: 2},
						{Name: "资产健康度分布", Path: "/admin/decision/health-distribution", Component: "views/admin/decision/health-distribution", MenuType: 2, Sort: 3},
						{Name: "治理工作量统计", Path: "/admin/decision/workload", Component: "views/admin/decision/workload", MenuType: 2, Sort: 4},
					},
				},
				{Name: "供应链安全洞察", Path: "/admin/decision/core-dependencies", Component: "views/admin/decision/core-dependencies", MenuType: 2, Sort: 2},
			},
		},
	}

	// ===== specialist 角色菜单 =====
	specialistMenus := []menuDef{
		{Name: "我的项目", Path: "/specialist/my-projects", Component: "views/operator/my-assets/MyProjects", Icon: "Monitor", MenuType: 2, Sort: 1},
		{Name: "我的任务", Path: "/specialist/my-task", Component: "views/specialist/MyTask", Icon: "List", MenuType: 2, Sort: 2},
		{
			Name: "我的运维", MenuType: 1, Sort: 3, Icon: "Setting",
			Children: []menuDef{
				{Name: "运维台账", Path: "/specialist/my-ops", Component: "views/specialist/MyOps", MenuType: 2, Sort: 1},
				{Name: "人员管理", Path: "/specialist/personnel-mgmt", Component: "views/operator/PersonnelMgmt", MenuType: 2, Sort: 2},
			},
		},
		{Name: "我的使用", Path: "/specialist/my-usage", Component: "views/operator/MyUsage", Icon: "Files", MenuType: 2, Sort: 4},
		{Name: "可信仓库", Path: "/specialist/trusted-repo", Component: "views/operator/TrustedRepo", Icon: "Files", MenuType: 2, Sort: 5},
	}

	// ===== operator 角色菜单 =====
	operatorMenus := []menuDef{
		{Name: "我的项目", Path: "/operator/my-projects", Component: "views/operator/my-assets/MyProjects", Icon: "Monitor", MenuType: 2, Sort: 1},
		{Name: "我的任务", Path: "/operator/my-task", Component: "views/operator/MyTask", Icon: "List", MenuType: 2, Sort: 2},
		{Name: "我的运维", Path: "/operator/my-ops", Component: "views/operator/MyOps", Icon: "Setting", MenuType: 2, Sort: 3},
		{Name: "我的使用", Path: "/operator/my-usage", Component: "views/operator/MyUsage", Icon: "Files", MenuType: 2, Sort: 4},
		{Name: "可信仓库", Path: "/operator/trusted-repo", Component: "views/operator/TrustedRepo", Icon: "Files", MenuType: 2, Sort: 5},
	}

	createMenuTree(db, adminMenus, 0, "admin")
	createMenuTree(db, specialistMenus, 0, "specialist")
	createMenuTree(db, operatorMenus, 0, "operator")
}
