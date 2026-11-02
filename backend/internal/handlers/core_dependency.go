package handlers

import (
	"fmt"
	"osm-mngt/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// ==================== 核心依赖统计（CoreDependency） ====================

// CoreDependencySummary 看板汇总指标 - 供应链安全洞察
type CoreDependencySummary struct {
	HighAdoptionCount  int64 `json:"highAdoptionCount"`  // 使用占比大于50%的组件数
	FragmentedCompCount int64 `json:"fragmentedCompCount"` // 碎片化组件数(≥3版本)
	HighRiskProjectCount int64 `json:"highRiskProjectCount"` // 高风险项目数（含黑名单阻断）
	CopyleftCompCount  int64 `json:"copyleftCompCount"`  // 强传染许可证关联组件数
}

// ComponentProjectAssoc 组件-项目关联明细
type ComponentProjectAssoc struct {
	OrgName     string `json:"orgName"`     // 组织（_group）
	CompName    string `json:"compName"`    // 组件名称（_name）
	Version     string `json:"version"`     // 版本号（_version）
	ProjectName string `json:"projectName"` // 项目编码（pro_code）
	ModuleCode  string `json:"moduleCode"`  // 模块编码
	GroupName   string `json:"groupName"`   // 业务组
	DeptName    string `json:"deptName"`    // 部门
}

// ComponentRankingItem 组件排行榜项
type ComponentRankingItem struct {
	CompName      string  `json:"compName"`      // 组件名称
	OrgName       string  `json:"orgName"`       // 组织
	ProjectCount  int64   `json:"projectCount"`  // 关联项目数（按 pro_code 去重）
	UsageRate     float64 `json:"usageRate"`     // 使用占比（%）
	LatestVersion string  `json:"latestVersion"` // 最新版本
}

// GetCoreDependencySummary 获取看板汇总指标 - 供应链安全洞察
func (h *AssetHandler) GetCoreDependencySummary(c *gin.Context) {
	var summary CoreDependencySummary

	// 1. 使用占比大于50%的组件数
	//    = 关联项目数 >= (总项目数 / 2) 的组件数
	//    totalProjects 与排行榜一致：只统计有组件快照的项目
	var totalProjects int64
	h.db.Table("component_snapshot cs").
		Joins("INNER JOIN project_inner pi ON cs.PROJECT_ID = pi.id").
		Where("pi.active = 1 AND cs.PROJECT_ID IS NOT NULL AND cs.PROJECT_ID > 0").
		Distinct("pi.pro_code").
		Count(&totalProjects)

	if totalProjects > 0 {
		// 子查询：每个组件关联的不同项目数
		// 外查询：统计关联项目数 >= totalProjects/2 的组件数
		h.db.Table("(?) as comp_proj",
			h.db.Table("component_snapshot cs").
				Joins("INNER JOIN project_inner pi ON cs.PROJECT_ID = pi.id").
				Where("pi.active = 1 AND cs.PROJECT_ID IS NOT NULL AND cs.PROJECT_ID > 0").
				Select("cs._name, COUNT(DISTINCT pi.pro_code) as proj_cnt").
				Group("cs._name"),
		).Where("proj_cnt >= ?", totalProjects/2).Count(&summary.HighAdoptionCount)
	}

	// 2. 碎片化组件数（≥3版本）
	//    统计 COUNT(DISTINCT _version) >= 3 的组件数
	h.db.Table("(?) as comp_ver",
		h.db.Table("component_snapshot").
			Where("PROJECT_ID IS NOT NULL AND PROJECT_ID > 0").
			Select("_name, COUNT(DISTINCT _version) as ver_cnt").
			Group("_name"),
	).Where("ver_cnt >= 3").Count(&summary.FragmentedCompCount)

	// 3. 高风险项目数（含黑名单阻断记录的项目数）
	//    统计 component_black 中 alert_type=1（阻断）关联的不同 pro_code 数
	h.db.Table("component_black cb").
		Joins("INNER JOIN project_inner pi ON cb.project_id = pi.id").
		Where("cb.alert_type = 1 AND pi.active = 1").
		Distinct("pi.pro_code").
		Count(&summary.HighRiskProjectCount)

	// 4. 强传染许可证关联组件数
	//    统计 LICENSE 包含 GPL（不含LGPL）或 AGPL 的不同组件名数
	h.db.Table("component_snapshot").
		Where("LICENSE IS NOT NULL AND LICENSE != ''").
		Where("LICENSE LIKE ?", "%GPL%").
		Not("LICENSE LIKE ?", "%LGPL%").
		Distinct("_name").
		Count(&summary.CopyleftCompCount)

	response.Success(c, summary)
}

// GetComponentProjectAssoc 获取组件-项目关联明细
func (h *AssetHandler) GetComponentProjectAssoc(c *gin.Context) {
	orgName := c.Query("orgName")
	compName := c.Query("compName")
	projectName := c.Query("projectName")

	page, pageSize := parsePagination(c)

	query := h.db.Table("component_snapshot cs").
		Select("cs._group as org_name, cs._name as comp_name, cs._version as version, "+
			"pi.pro_code as project_name, pi.module_code as module_code, "+
			"pi.group_name as group_name, pi.dept_name as dept_name").
		Joins("INNER JOIN project_inner pi ON cs.PROJECT_ID = pi.id").
		Where("pi.active = 1 AND cs.PROJECT_ID IS NOT NULL AND cs.PROJECT_ID > 0")

	if orgName != "" {
		query = query.Where("cs._group LIKE ?", "%"+orgName+"%")
	}
	if compName != "" {
		query = query.Where("cs._name LIKE ?", "%"+compName+"%")
	}
	if projectName != "" {
		query = query.Where("pi.pro_code LIKE ?", "%"+projectName+"%")
	}

	var total int64
	query.Count(&total)

	var list []ComponentProjectAssoc
	query.Order("cs._name ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&list)

	response.Success(c, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetComponentRanking 获取组件项目数排行榜
func (h *AssetHandler) GetComponentRanking(c *gin.Context) {
	orgName := c.Query("orgName")
	compName := c.Query("compName")

	page, pageSize := parsePagination(c)

	subQuery := h.db.Table("component_snapshot cs").
		Select("cs._name as comp_name, cs._group as org_name, "+
			"COUNT(DISTINCT pi.pro_code) as project_count, "+
			"MAX(cs._version) as latest_version").
		Joins("INNER JOIN project_inner pi ON cs.PROJECT_ID = pi.id").
		Where("pi.active = 1 AND cs.PROJECT_ID IS NOT NULL AND cs.PROJECT_ID > 0")

	if orgName != "" {
		subQuery = subQuery.Where("cs._group LIKE ?", "%"+orgName+"%")
	}
	if compName != "" {
		subQuery = subQuery.Where("cs._name LIKE ?", "%"+compName+"%")
	}

	var total int64
	h.db.Table("(?) as ranked", subQuery.Group("cs._name, cs._group")).Count(&total)

	type rankRow struct {
		CompName      string `json:"compName"`
		OrgName       string `json:"orgName"`
		ProjectCount  int64  `json:"projectCount"`
		LatestVersion string `json:"latestVersion"`
	}

	var rows []rankRow
	subQuery.Group("cs._name, cs._group").
		Order("project_count DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows)

	// 使用占比的分母：按 pro_code 去重的涉及项目总数
	var totalProjects int64
	h.db.Table("component_snapshot cs").
		Joins("INNER JOIN project_inner pi ON cs.PROJECT_ID = pi.id").
		Where("pi.active = 1 AND cs.PROJECT_ID IS NOT NULL AND cs.PROJECT_ID > 0").
		Distinct("pi.pro_code").
		Count(&totalProjects)

	list := make([]ComponentRankingItem, 0, len(rows))
	for _, r := range rows {
		rate := 0.0
		if totalProjects > 0 {
			rate = float64(r.ProjectCount) / float64(totalProjects) * 100
		}
		rate = float64(int(rate*10)) / 10

		list = append(list, ComponentRankingItem{
			CompName:      r.CompName,
			OrgName:       r.OrgName,
			ProjectCount:  r.ProjectCount,
			UsageRate:     rate,
			LatestVersion: r.LatestVersion,
		})
	}

	response.Success(c, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// VersionFragmentationItem 版本碎片化台账项
type VersionFragmentationItem struct {
	CompName        string `json:"compName"`        // 组件名称
	OrgName         string `json:"orgName"`         // 组织
	Classifier      string `json:"classifier"`      // 类别
	VersionCount    int64  `json:"versionCount"`    // 在用版本数
	VersionList     string `json:"versionList"`     // 版本列表（逗号分隔）
	LatestVersion   string `json:"latestVersion"`   // 最新版本
	FragmentLevel   string `json:"fragmentLevel"`   // 碎片化程度: high/medium/low
}

// GetVersionFragmentation 获取版本碎片化台账
// 按 _name 分组统计在用版本数，按版本数降序排列
func (h *AssetHandler) GetVersionFragmentation(c *gin.Context) {
	compName := c.Query("compName")

	page, pageSize := parsePagination(c)

	// 按 _name, _group 分组，统计版本数和版本列表
	// 利用 idx_name_version 索引优化 GROUP BY
	query := h.db.Table("component_snapshot cs").
		Select("cs._name as comp_name, cs._group as org_name, "+
			"MAX(cs.CLASSIFIER) as classifier, "+
			"COUNT(DISTINCT cs._version) as version_count, "+
			"GROUP_CONCAT(DISTINCT cs._version ORDER BY cs._version DESC) as version_list, "+
			"MAX(cs._version) as latest_version").
		Joins("INNER JOIN project_inner pi ON cs.PROJECT_ID = pi.id").
		Where("pi.active = 1 AND cs.PROJECT_ID IS NOT NULL AND cs.PROJECT_ID > 0")

	if compName != "" {
		query = query.Where("cs._name LIKE ?", "%"+compName+"%")
	}

	// 统计总数
	var total int64
	h.db.Table("(?) as vf", query.Group("cs._name, cs._group")).Count(&total)

	// 执行分页查询
	type vfRow struct {
		CompName      string `json:"compName"`
		OrgName       string `json:"orgName"`
		Classifier    string `json:"classifier"`
		VersionCount  int64  `json:"versionCount"`
		VersionList   string `json:"versionList"`
		LatestVersion string `json:"latestVersion"`
	}

	var rows []vfRow
	query.Group("cs._name, cs._group").
		Order("version_count DESC, comp_name ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows)

	// 计算碎片化程度
	list := make([]VersionFragmentationItem, 0, len(rows))
	for _, r := range rows {
		// 碎片化程度：>5 为高，>3 为中，其余为低
		level := "low"
		if r.VersionCount > 5 {
			level = "high"
		} else if r.VersionCount > 3 {
			level = "medium"
		}

		list = append(list, VersionFragmentationItem{
			CompName:      r.CompName,
			OrgName:       r.OrgName,
			Classifier:    r.Classifier,
			VersionCount:  r.VersionCount,
			VersionList:   r.VersionList,
			LatestVersion: r.LatestVersion,
			FragmentLevel: level,
		})
	}

	response.Success(c, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// HighRiskProjectItem 高风险项目台账项
type HighRiskProjectItem struct {
	ProjectName   string `json:"projectName"`   // 项目名称（pro_code）
	BlacklistCount int64 `json:"blacklistCount"` // 黑名单数（警告+提示+阻断都算）
	GroupName     string `json:"groupName"`     // 团队
	DeptName      string `json:"deptName"`      // 部门
}

// GetHighRiskProjects 获取高风险项目台账
// JOIN component_black + project_inner，按 pro_code 分组统计黑名单数
func (h *AssetHandler) GetHighRiskProjects(c *gin.Context) {
	projectName := c.Query("projectName")

	page, pageSize := parsePagination(c)

	// 按 pro_code 分组，统计每个项目的黑名单记录数
	// 利用 idx_project_id 索引加速 JOIN
	query := h.db.Table("component_black cb").
		Select("pi.pro_code as project_name, "+
			"COUNT(cb.id) as blacklist_count, "+
			"MAX(pi.group_name) as group_name, "+
			"MAX(pi.dept_name) as dept_name").
		Joins("INNER JOIN project_inner pi ON cb.project_id = pi.id").
		Where("pi.active = 1 AND cb.project_id IS NOT NULL AND cb.project_id > 0")

	if projectName != "" {
		query = query.Where("pi.pro_code LIKE ?", "%"+projectName+"%")
	}

	// 统计总数
	var total int64
	h.db.Table("(?) as hr", query.Group("pi.pro_code")).Count(&total)

	// 执行分页查询，按黑名单数降序
	var rows []HighRiskProjectItem
	query.Group("pi.pro_code").
		Order("blacklist_count DESC, project_name ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows)

	response.Success(c, gin.H{
		"list":     rows,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// BlacklistDetailItem 黑名单核查记录明细
type BlacklistDetailItem struct {
	ModuleCode    string `json:"moduleCode"`    // 模块名称
	CommitID      string `json:"commitId"`      // 当前版本 commitid
	CreateTime    string `json:"createTime"`    // 核查时间
	AlertType     int    `json:"alertType"`     // 处置：1-阻断, 2-告警
	BlacklistName string `json:"blacklistName"` // 黑名单名称
	Reason        string `json:"reason"`        // 原因
}

// GetProjectBlacklistDetail 获取项目最新提交关联的黑名单核查记录明细
// JOIN component_black + project_inner + component_core
func (h *AssetHandler) GetProjectBlacklistDetail(c *gin.Context) {
	projectName := c.Query("projectName")
	if projectName == "" {
		response.Error(c, 400, "请指定项目名称")
		return
	}

	// 查询该项目所有黑名单记录，关联 project_inner 获取模块和 commitid，关联 component_core 获取黑名单名称和原因
	var rows []BlacklistDetailItem
	h.db.Table("component_black cb").
		Select("pi.module_code as module_code, "+
			"cb.CURRENT_ID as commit_id, "+
			"DATE_FORMAT(cb.create_time, '%Y-%m-%d %H:%i:%s') as create_time, "+
			"cb.alert_type as alert_type, "+
			"cc._name as blacklist_name, "+
			"cc.reason as reason").
		Joins("INNER JOIN project_inner pi ON cb.PROJECT_ID = pi.id").
		Joins("LEFT JOIN component_core cc ON cb.BLACK_ID = cc.id").
		Where("pi.pro_code = ? AND pi.active = 1", projectName).
		Order("cb.create_time DESC").
		Scan(&rows)

	response.Success(c, gin.H{
		"list":  rows,
		"total": len(rows),
	})
}

// SeedHighRiskTestData 插入高风险项目测试数据
// 向 component_black 表插入约 20 条记录，关联已有的 project_inner 和 component_core 记录
func (h *AssetHandler) SeedHighRiskTestData(c *gin.Context) {
	// 清理旧的测试数据
	h.db.Exec("DELETE FROM component_black WHERE create_by = 'test-seed'")

	// 查询已有的 project_inner 记录（取前 10 个不同 pro_code 的项目）
	type projInfo struct {
		ID        int64
		ProCode   string
		CurrentID string
	}
	var projects []projInfo
	h.db.Raw("SELECT MIN(id) as id, pro_code, MAX(current_ID) as current_id FROM project_inner WHERE active = 1 GROUP BY pro_code ORDER BY id LIMIT 10").Scan(&projects)

	if len(projects) == 0 {
		response.Error(c, 400, "没有可用的自研项目数据，请先执行 core-dependencies/seed-test-data")
		return
	}

	// 查询已有的 component_core 记录
	type coreInfo struct {
		ID int64
	}
	var cores []coreInfo
	h.db.Raw("SELECT id FROM component_core WHERE _enable = 1 LIMIT 4").Scan(&cores)

	if len(cores) == 0 {
		response.Error(c, 400, "没有可用的黑名单规则数据")
		return
	}

	// 查询已有的 component_snapshot 记录
	type snapInfo struct {
		ID int64
	}
	var snaps []snapInfo
	h.db.Raw("SELECT id FROM component_snapshot LIMIT 20").Scan(&snaps)

	if len(snaps) == 0 {
		response.Error(c, 400, "没有可用的组件快照数据")
		return
	}

	// 告警类型：1-阻断, 2-警告
	alertTypes := []int{1, 2}
	// 处理状态：0-未处理, 1-处理中, 2-处理完成, 3-误报
	processStatuses := []int{0, 1, 2, 3}

	count := 0
	for i, proj := range projects {
		// 每个项目分配 2-3 条黑名单记录，总共约 20 条
		numRecords := 2
		if i < 5 {
			numRecords = 3
		} else if i >= 7 {
			numRecords = 1
		}

		for j := 0; j < numRecords; j++ {
			core := cores[(i+j)%len(cores)]
			snap := snaps[(i*3+j)%len(snaps)]
			alertType := alertTypes[(i+j)%2]
			processStatus := processStatuses[(i+j)%4]

			result := h.db.Exec(
				"INSERT INTO component_black (BLACK_ID, COMPONENT_ID, PROJECT_ID, PROCESS_STATUS, CURRENT_ID, internal, alert_type, create_by, update_by) VALUES (?, ?, ?, ?, ?, 1, ?, 'test-seed', 'test-seed')",
				core.ID, snap.ID, proj.ID, processStatus, proj.CurrentID, alertType,
			)
			if result.Error != nil {
				continue
			}
			count++
		}
	}

	response.Success(c, gin.H{
		"message":     "高风险项目测试数据插入成功",
		"recordCount": count,
	})
}

// LicenseOverviewItem 许可证台账项
type LicenseOverviewItem struct {
	LicenseName    string `json:"licenseName"`    // 许可证名称
	ComponentCount int64  `json:"componentCount"` // 关联组件数
	RiskLevel      string `json:"riskLevel"`      // 风险等级：high/medium/low/unknown
}

// GetLicenseOverview 获取许可证台账
// 按 LICENSE 字段分组统计 component_snapshot 中的组件数
func (h *AssetHandler) GetLicenseOverview(c *gin.Context) {
	licenseName := c.Query("licenseName")

	page, pageSize := parsePagination(c)

	// 单表 GROUP BY LICENSE，避免 JOIN 提升效率
	query := h.db.Table("component_snapshot").
		Select("LICENSE as license_name, COUNT(*) as component_count").
		Where("LICENSE IS NOT NULL AND LICENSE != ''")

	if licenseName != "" {
		query = query.Where("LICENSE LIKE ?", "%"+licenseName+"%")
	}

	// 统计总数
	var total int64
	h.db.Table("(?) as t", query.Group("LICENSE")).Count(&total)

	// 执行分页查询，按组件数降序
	var rows []LicenseOverviewItem
	query.Group("LICENSE").
		Order("component_count DESC, license_name ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows)

	// 根据许可证名称判定风险等级
	for i := range rows {
		rows[i].RiskLevel = classifyLicenseRisk(rows[i].LicenseName)
	}

	response.Success(c, gin.H{
		"list":     rows,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// classifyLicenseRisk 根据许可证名称判定风险等级
func classifyLicenseRisk(license string) string {
	upper := strings.ToUpper(license)
	switch {
	case strings.Contains(upper, "GPL") && !strings.Contains(upper, "LGPL") && !strings.Contains(upper, "MOZILLA") && !strings.Contains(upper, "MPL"):
		return "high" // 强传染型：GPL
	case strings.Contains(upper, "AGPL"):
		return "high" // 强传染型：AGPL
	case strings.Contains(upper, "LGPL") || strings.Contains(upper, "MPL") || strings.Contains(upper, "MOZILLA") || strings.Contains(upper, "EPL") || strings.Contains(upper, "CDDL"):
		return "medium" // 弱传染型：LGPL/MPL/EPL/CDDL
	case strings.Contains(upper, "MIT") || strings.Contains(upper, "BSD") || strings.Contains(upper, "APACHE") || strings.Contains(upper, "ISC") || strings.Contains(upper, "UNLICENSE"):
		return "low" // 宽松型：MIT/BSD/Apache/ISC
	default:
		return "unknown" // 未知型
	}
}

// SeedLicenseTestData 为已有的 component_snapshot 记录填充 LICENSE 数据
func (h *AssetHandler) SeedLicenseTestData(c *gin.Context) {
	// 查询已有的 component_snapshot 记录中 LICENSE 为空的
	type snapInfo struct {
		ID int64
	}
	var snaps []snapInfo
	h.db.Raw("SELECT id FROM component_snapshot WHERE LICENSE IS NULL OR LICENSE = '' LIMIT 200").Scan(&snaps)

	if len(snaps) == 0 {
		// 如果没有空 LICENSE 的记录，检查是否已经有 LICENSE 数据
		var hasLicense int64
		h.db.Raw("SELECT COUNT(*) FROM component_snapshot WHERE LICENSE IS NOT NULL AND LICENSE != ''").Scan(&hasLicense)
		if hasLicense > 0 {
			response.Success(c, gin.H{
				"message":    "已存在许可证数据，无需重复插入",
				"recordCount": 0,
				"existing":   hasLicense,
			})
			return
		}
		response.Error(c, 400, "没有可用的 component_snapshot 数据")
		return
	}

	// 许可证池及其 license 表 ID
	licenses := []struct {
		Code string
		ID   int64
	}{
		{"MIT", 60002},
		{"Apache-2.0", 60003},
		{"GPL-2.0", 60004},
		{"GPL-3.0", 60005},
		{"BSD-3-Clause", 60006},
		{"BSD-2-Clause", 60007},
		{"LGPL-3.0", 60008},
		{"Unknown", 0}, // 未知许可证
	}

	count := 0
	for i, snap := range snaps {
		lic := licenses[i%len(licenses)]
		licID := lic.ID
		// Unknown 的 LICENSE_ID 设为 NULL
		if licID == 0 {
			h.db.Exec("UPDATE component_snapshot SET LICENSE = ? WHERE id = ?", lic.Code, snap.ID)
		} else {
			h.db.Exec("UPDATE component_snapshot SET LICENSE = ?, LICENSE_ID = ? WHERE id = ?", lic.Code, licID, snap.ID)
		}
		count++
	}

	response.Success(c, gin.H{
		"message":    "许可证测试数据填充成功",
		"recordCount": count,
	})
}

// SeedCoreDependencyTestData 插入测试数据用于演示
// 使用 raw SQL 确保数据正确插入并建立 project_inner -> component_snapshot 的 PROJECT_ID 关联
func (h *AssetHandler) SeedCoreDependencyTestData(c *gin.Context) {
	// 1. 先清理旧的测试数据（先删 component_snapshot 再删 project_inner，避免孤儿引用）
	// 通过子查询找到 TEST-DEP 项目的 ID，删除关联的 component_snapshot
	h.db.Exec("DELETE FROM component_snapshot WHERE PROJECT_ID IN (SELECT id FROM project_inner WHERE pro_code LIKE 'TEST-DEP-%')")
	h.db.Exec("DELETE FROM project_inner WHERE pro_code LIKE 'TEST-DEP-%'")

	// 2. 插入自研项目（使用 raw SQL 确保所有字段正确写入）
	projects := []struct {
		ProCode    string
		Type       string
		GroupName  string
		DeptName   string
		ModuleCode string
		CurrentID  string
	}{
		{"TEST-DEP-001", "java", "平台组", "技术研发部", "platform-core", "test-cid-001"},
		{"TEST-DEP-002", "java", "支付组", "技术研发部", "pay-service", "test-cid-002"},
		{"TEST-DEP-003", "go", "运维组", "基础架构部", "ops-agent", "test-cid-003"},
		{"TEST-DEP-004", "python", "数据组", "数据分析部", "data-pipeline", "test-cid-004"},
		{"TEST-DEP-005", "java", "风控组", "风险控制部", "risk-engine", "test-cid-005"},
		{"TEST-DEP-006", "java", "用户组", "技术研发部", "user-center", "test-cid-006"},
		{"TEST-DEP-007", "go", "网关组", "基础架构部", "api-gateway", "test-cid-007"},
		{"TEST-DEP-008", "python", "推荐组", "数据分析部", "rec-engine", "test-cid-008"},
	}

	projectIDs := make([]int64, 0, len(projects))
	for _, p := range projects {
		result := h.db.Exec(
			"INSERT INTO project_inner (pro_code, active, type, is_boundary, group_name, dept_name, module_code, branch, current_ID, create_by, update_by) VALUES (?, 1, ?, 0, ?, ?, ?, 'main', ?, 'system', 'system')",
			p.ProCode, p.Type, p.GroupName, p.DeptName, p.ModuleCode, p.CurrentID,
		)
		if result.Error != nil {
			response.Error(c, 500, fmt.Sprintf("插入项目 %s 失败: %v", p.ProCode, result.Error))
			return
		}
		// 查询刚插入的自增 ID
		var id int64
		h.db.Raw("SELECT id FROM project_inner WHERE pro_code = ?", p.ProCode).Scan(&id)
		projectIDs = append(projectIDs, id)
	}

	// 3. 组件定义
	components := []struct {
		Group   string
		Name    string
		Version string
	}{
		{"org.springframework", "spring-boot", "3.2.3"},
		{"org.springframework", "spring-boot", "2.7.18"},
		{"com.google.guava", "guava", "33.0.0-jre"},
		{"org.slf4j", "slf4j-api", "2.0.12"},
		{"org.apache.commons", "commons-lang3", "3.14.0"},
		{"com.fasterxml.jackson.core", "jackson-databind", "2.16.1"},
		{"redis.clients", "jedis", "5.1.0"},
		{"com.zaxxer", "HikariCP", "5.1.0"},
		{"org.mybatis", "mybatis", "3.5.15"},
		{"com.alibaba", "druid", "1.2.21"},
		{"io.netty", "netty-all", "4.1.107.Final"},
		{"org.projectlombok", "lombok", "1.18.30"},
		{"cn.hutool", "hutool-all", "5.8.26"},
		{"com.github.pagehelper", "pagehelper", "6.1.0"},
	}

	// 4. 为每个项目分配组件（模拟真实使用分布）
	// 项目索引 -> 组件索引列表
	projectComponents := map[int][]int{
		0: {0, 2, 3, 4, 5, 9},          // TEST-DEP-001: spring-boot3, guava, slf4j, commons, jackson, druid
		1: {0, 3, 5, 6, 7, 8},          // TEST-DEP-002: spring-boot3, slf4j, jackson, jedis, HikariCP, mybatis
		2: {1, 4, 10, 12},              // TEST-DEP-003: spring-boot2, commons, netty, hutool
		3: {3, 4, 5, 11, 12},           // TEST-DEP-004: slf4j, commons, jackson, lombok, hutool
		4: {0, 2, 3, 5, 6, 7, 9, 13},   // TEST-DEP-005: spring-boot3, guava, slf4j, jackson, jedis, HikariCP, druid, pagehelper
		5: {1, 3, 4, 5, 8, 9},          // TEST-DEP-006: spring-boot2, slf4j, commons, jackson, mybatis, druid
		6: {3, 4, 10, 12, 13},          // TEST-DEP-007: slf4j, commons, netty, hutool, pagehelper
		7: {3, 4, 5, 11, 12, 13},       // TEST-DEP-008: slf4j, commons, jackson, lombok, hutool, pagehelper
	}

	count := 0
	for projIdx, compIdxs := range projectComponents {
		if projIdx >= len(projectIDs) {
			continue
		}
		pid := projectIDs[projIdx]
		for _, compIdx := range compIdxs {
			if compIdx >= len(components) {
				continue
			}
			comp := components[compIdx]
			result := h.db.Exec(
				"INSERT INTO component_snapshot (current_id, _group, _name, _version, CLASSIFIER, PROJECT_ID) VALUES (?, ?, ?, ?, 'MAVEN', ?)",
				projects[projIdx].CurrentID, comp.Group, comp.Name, comp.Version, pid,
			)
			if result.Error != nil {
				continue
			}
			count++
		}
	}

	response.Success(c, gin.H{
		"message":        "测试数据插入成功",
		"projectCount":   len(projectIDs),
		"componentCount": count,
	})
}
