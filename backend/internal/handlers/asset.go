package handlers

import (
	"crypto/rand"
	"fmt"
	"osm-mngt/internal/models"
	"osm-mngt/pkg/response"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetHandler struct {
	db *gorm.DB
}

func NewAssetHandler(db *gorm.DB) *AssetHandler {
	return &AssetHandler{db: db}
}

func (h *AssetHandler) getUserInfo(c *gin.Context) (string, string) {
	userID := c.GetHeader("X-User-ID")
	userName := c.GetHeader("X-User-Name")
	return userID, userName
}

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// ==================== 开源基础软件（project） ====================

func (h *AssetHandler) ListProjects(c *gin.Context) {
	keyword := c.Query("keyword")
	typeFilter := c.Query("type")
	levelFilter := c.Query("applicable_level")
	deptFilter := c.Query("maint_dept")
	categoryFilter := c.Query("category")
	ownerUsername := c.Query("owner_username") // 按负责人过滤（operator 查自己负责的）

	query := h.db.Model(&models.Project{}).Where("ACTIVE = ?", 1)

	// 默认只查开源基础软件（兼容未传 category 的旧调用）
	if categoryFilter != "" {
		query = query.Where("category = ?", categoryFilter)
	} else {
		query = query.Where("(category = ? OR category = ?)", "software", "")
	}

	if keyword != "" {
		query = query.Where("_name LIKE ? OR _group LIKE ? OR type LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if typeFilter != "" {
		query = query.Where("type = ?", typeFilter)
	}
	if levelFilter != "" {
		query = query.Where("applicable_level = ?", levelFilter)
	}
	if deptFilter != "" {
		query = query.Where("maint_dept = ?", deptFilter)
	}
	// 按负责人过滤：只返回 project_owner 表中分配给该用户的项目
	if ownerUsername != "" {
		query = query.Where("id IN (SELECT project_id FROM project_owner WHERE project_type = 'software' AND owner_username = ?)", ownerUsername)
	}

	var projects []models.Project
	if err := query.Order("create_time DESC").Find(&projects).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, projects)
}

func (h *AssetHandler) GetProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的项目ID")
		return
	}

	var project models.Project
	if err := h.db.Where("ACTIVE = ?", 1).First(&project, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "项目不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, project)
}

func (h *AssetHandler) CreateProject(c *gin.Context) {
	var req models.Project
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if req.Name == "" || req.Type == "" || req.ApplicableLevel == "" ||
		req.MaintDept == "" || req.UseDept == "" || req.Country == "" {
		response.Error(c, 400, "缺少必填字段")
		return
	}

	active := 1
	req.Active = &active
	req.UUID = generateUUID()

	// 默认类别为开源基础软件
	if req.Category == "" {
		req.Category = "software"
	}

	userID, userName := h.getUserInfo(c)
	req.CreateBy = userID
	req.UpdateBy = userName
	req.CreateTime = time.Now()
	req.UpdateTime = time.Now()

	if err := h.db.Create(&req).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "duplicate") {
			response.Error(c, 400, "项目已存在")
			return
		}
		response.Error(c, 500, "创建失败")
		return
	}

	response.Success(c, req)
}

func (h *AssetHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的项目ID")
		return
	}

	var project models.Project
	if err := h.db.Where("ACTIVE = ?", 1).First(&project, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "项目不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	// 移除不允许更新的字段
	delete(updates, "id")
	delete(updates, "UUID")
	delete(updates, "create_by")
	delete(updates, "create_time")

	_, userName := h.getUserInfo(c)
	updates["update_by"] = userName
	updates["update_time"] = time.Now()

	if err := h.db.Model(&models.Project{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.Error(c, 500, "更新失败")
		return
	}

	response.Success(c, nil)
}

func (h *AssetHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的项目ID")
		return
	}

	var project models.Project
	if err := h.db.Where("ACTIVE = ?", 1).First(&project, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "项目不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	_, userName := h.getUserInfo(c)

	if err := h.db.Model(&models.Project{}).Where("id = ?", id).Updates(map[string]interface{}{
		"ACTIVE":      0,
		"update_by":   userName,
		"update_time": time.Now(),
	}).Error; err != nil {
		response.Error(c, 500, "删除失败")
		return
	}

	response.Success(c, nil)
}

// GetProjectComponents 获取开源基础软件的关联依赖
func (h *AssetHandler) GetProjectComponents(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的项目ID")
		return
	}

	keyword := c.Query("keyword")
	groupFilter := c.Query("group")
	versionFilter := c.Query("version")

	query := h.db.Model(&models.Component{}).Where("PROJECT_ID = ?", projectID)

	if keyword != "" {
		query = query.Where("_name LIKE ? OR _group LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if groupFilter != "" {
		query = query.Where("_group = ?", groupFilter)
	}
	if versionFilter != "" {
		query = query.Where("_version = ?", versionFilter)
	}

	var components []models.Component
	if err := query.Order("LAST_RISKSCORE DESC").Find(&components).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, components)
}

// GetProjectLicenseStats 获取开源基础软件的许可证统计
func (h *AssetHandler) GetProjectLicenseStats(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的项目ID")
		return
	}

	type LicenseStat struct {
		License string `json:"license"`
		Count   int    `json:"count"`
	}

	var licenseStats []LicenseStat
	h.db.Raw(`
		SELECT LICENSE as license, COUNT(*) as count 
		FROM component 
		WHERE PROJECT_ID = ? AND LICENSE IS NOT NULL AND LICENSE != ''
		GROUP BY LICENSE
	`, projectID).Scan(&licenseStats)

	type LicenseGroupStat struct {
		GroupName string `json:"group_name"`
		RiskLevel string `json:"risk_level"`
		Count     int    `json:"count"`
	}

	var groupStats []LicenseGroupStat
	h.db.Raw(`
		SELECT
			lg._name as group_name,
			CASE
				WHEN lg.RISKWEIGHT <= 20 THEN '低风险'
				WHEN lg.RISKWEIGHT <= 50 THEN '中风险'
				ELSE '高风险'
			END as risk_level,
			COUNT(*) as count
		FROM component cs
		LEFT JOIN license l ON cs.LICENSE = l.LICENSE
		LEFT JOIN licensegroup_license lgl ON l.id = lgl.license_id
		LEFT JOIN licensegroup lg ON lgl.licensegroup_id = lg.id
		WHERE cs.PROJECT_ID = ?
		GROUP BY lg._name, risk_level
	`, projectID).Scan(&groupStats)

	response.Success(c, gin.H{
		"license_stats": licenseStats,
		"group_stats":   groupStats,
	})
}

// GetProjectDependencies 获取项目的依赖组件
func (h *AssetHandler) GetProjectDependencies(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的项目ID")
		return
	}

	currentID := c.Query("current_ID")

	// 查询component表中该项目的依赖组件
	query := h.db.Table("component").
		Select("component._name as component_name, component._version as version, component.LICENSE as license").
		Where("component.PROJECT_ID = ?", id)

	if currentID != "" {
		query = query.Where("component.current_id = ?", currentID)
	}

	var dependencies []map[string]interface{}
	if err := query.Scan(&dependencies).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, dependencies)
}

// ==================== 自研项目（project_inner） ====================

func (h *AssetHandler) ListInnerProjects(c *gin.Context) {
	keyword := c.Query("keyword")
	typeFilter := c.Query("type")
	deptFilter := c.Query("dept_name")
	groupFilter := c.Query("group_name")
	ownerUsername := c.Query("owner_username") // 按负责人过滤（operator 查自己负责的）

	query := h.db.Model(&models.ProjectInner{})

	if keyword != "" {
		query = query.Where("pro_code LIKE ? OR module_code LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if typeFilter != "" {
		query = query.Where("type = ?", typeFilter)
	}
	if deptFilter != "" {
		query = query.Where("dept_name = ?", deptFilter)
	}
	if groupFilter != "" {
		query = query.Where("group_name = ?", groupFilter)
	}
	// 按负责人过滤：只返回 project_owner 表中分配给该用户的项目
	if ownerUsername != "" {
		query = query.Where("id IN (SELECT project_id FROM project_owner WHERE project_type = 'inner' AND owner_username = ?)", ownerUsername)
	}

	var projects []models.ProjectInner
	if err := query.Order("create_time DESC").Find(&projects).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, projects)
}

func (h *AssetHandler) GetInnerProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的项目ID")
		return
	}

	var project models.ProjectInner
	if err := h.db.Where("id = ?", id).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "项目不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, project)
}

func (h *AssetHandler) CreateInnerProject(c *gin.Context) {
	var req models.ProjectInner
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if req.ProCode == "" || req.ModuleCode == "" || req.Branch == "" || req.Type == "" {
		response.Error(c, 400, "缺少必填字段")
		return
	}

	userID, userName := h.getUserInfo(c)
	req.CreateBy = userID
	req.UpdateBy = userName
	req.CreateTime = time.Now()
	req.UpdateTime = time.Now()

	active := 1
	req.Active = &active

	if err := h.db.Create(&req).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "duplicate") {
			response.Error(c, 400, "项目已存在")
			return
		}
		response.Error(c, 500, "创建失败")
		return
	}

	response.Success(c, req)
}

func (h *AssetHandler) UpdateInnerProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的项目ID")
		return
	}

	var project models.ProjectInner
	if err := h.db.First(&project, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "项目不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	delete(updates, "id")
	delete(updates, "create_by")
	delete(updates, "create_time")

	_, userName := h.getUserInfo(c)
	updates["update_by"] = userName
	updates["update_time"] = time.Now()

	if err := h.db.Model(&models.ProjectInner{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.Error(c, 500, "更新失败")
		return
	}

	response.Success(c, nil)
}

func (h *AssetHandler) DeleteInnerProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的项目ID")
		return
	}

	var project models.ProjectInner
	if err := h.db.First(&project, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "项目不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	_, userName := h.getUserInfo(c)

	if err := h.db.Model(&models.ProjectInner{}).Where("id = ?", id).Updates(map[string]interface{}{
		"active":      0,
		"update_by":   userName,
		"update_time": time.Now(),
	}).Error; err != nil {
		response.Error(c, 500, "删除失败")
		return
	}

	response.Success(c, nil)
}

// ==================== 开源依赖包（component_snapshot） ====================

func (h *AssetHandler) ListComponents(c *gin.Context) {
	keyword := c.Query("keyword")
	classifierFilter := c.Query("classifier")
	// 依赖来源（开源依赖包页数据范围）：inner=自研项目组件（component_snapshot，默认）/ software=开源软件组件（component）
	source := c.Query("source")
	table := "component_snapshot"
	if source == "software" {
		table = "component"
	}

	// 按 _group + _name + _version 去重，取每组最新记录
	query := h.db.Table(table).
		Where("id IN (?)",
			h.db.Table(table).
				Select("MAX(id)").
				Group("_group, _name, _version"))

	if keyword != "" {
		query = query.Where("_name LIKE ? OR _group LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if classifierFilter != "" {
		query = query.Where("CLASSIFIER = ?", classifierFilter)
	}

	var components []models.ComponentSnapshot
	if err := query.Order("_group ASC, _name ASC").Find(&components).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, components)
}

func (h *AssetHandler) GetComponent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的依赖包ID")
		return
	}

	var component models.ComponentSnapshot
	if err := h.db.First(&component, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "依赖包不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, component)
}

// GetRelatedProjects 获取依赖包关联的项目（按依赖来源）
// source=software：component 表 PROJECT_ID → project（开源软件）；默认：component_snapshot → project_inner（自研项目）
func (h *AssetHandler) GetRelatedProjects(c *gin.Context) {
	componentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的依赖包ID")
		return
	}

	source := c.Query("source")
	table := "component_snapshot"
	if source == "software" {
		table = "component"
	}

	// 查询对应表中该依赖包的 PROJECT_ID
	var component models.ComponentSnapshot
	if err := h.db.Table(table).First(&component, componentID).Error; err != nil {
		response.Error(c, 404, "依赖包不存在")
		return
	}

	// 如果没有关联项目，返回空列表
	if component.ProjectID == nil || *component.ProjectID == 0 {
		response.Success(c, []gin.H{})
		return
	}

	// 开源软件组件：关联 project（开源基础软件）表
	if source == "software" {
		var project models.Project
		if err := h.db.Where("id = ?", *component.ProjectID).First(&project).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.Success(c, []gin.H{})
				return
			}
			response.Error(c, 500, "查询失败")
			return
		}
		response.Success(c, []models.Project{project})
		return
	}

	// 自研项目组件：关联 project_inner 表
	var project models.ProjectInner
	if err := h.db.Where("id = ? AND active = ?", *component.ProjectID, 1).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Success(c, []models.ProjectInner{})
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, []models.ProjectInner{project})
}

// GetInnerProjectModules 获取自研项目的所有模块列表
func (h *AssetHandler) GetInnerProjectModules(c *gin.Context) {
	proCode := c.Query("pro_code")
	if proCode == "" {
		response.Error(c, 400, "项目编码不能为空")
		return
	}

	var modules []struct {
		ID         int64  `json:"id"`
		ModuleCode string `json:"module_code"`
	}
	if err := h.db.Model(&models.ProjectInner{}).
		Where("pro_code = ?", proCode).
		Select("id, module_code").
		Find(&modules).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, modules)
}

// GetInnerProjectComponents 获取自研项目的组件信息
func (h *AssetHandler) GetInnerProjectComponents(c *gin.Context) {
	projectIDsStr := c.Query("project_ids")
	currentID := c.Query("current_ID")
	moduleCode := c.Query("module_code")
	keyword := c.Query("keyword")
	version := c.Query("version")

	// 通过 PROJECT_ID 关联 project_inner 获取模块信息
	query := h.db.Table("component_snapshot").
		Select("component_snapshot.*, pi.module_code, pi.pro_code").
		Joins("LEFT JOIN project_inner pi ON component_snapshot.PROJECT_ID = pi.id")

	// 通过 current_ID 查询（全量或版本范围）
	if currentID != "" {
		query = query.Where("component_snapshot.current_id = ?", currentID)
	}

	// 通过 project_ids 过滤特定模块（PROJECT_ID 关联）
	if projectIDsStr != "" {
		var projectIDs []int64
		for _, idStr := range strings.Split(projectIDsStr, ",") {
			if id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64); err == nil {
				projectIDs = append(projectIDs, id)
			}
		}
		if len(projectIDs) > 0 {
			query = query.Where("component_snapshot.PROJECT_ID IN ?", projectIDs)
		}
	}

	// 按模块名查询
	if moduleCode != "" {
		query = query.Where("pi.module_code = ?", moduleCode)
	}

	// 按关键词查询
	if keyword != "" {
		query = query.Where("component_snapshot._name LIKE ? OR component_snapshot._group LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// 按版本查询
	if version != "" {
		query = query.Where("component_snapshot._version = ?", version)
	}

	var components []map[string]interface{}
	if err := query.Order("component_snapshot.id DESC").Scan(&components).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, components)
}

// GetInnerProjectLicenseStats 获取自研项目的许可证统计
func (h *AssetHandler) GetInnerProjectLicenseStats(c *gin.Context) {
	projectIDsStr := c.Query("project_ids")
	currentID := c.Query("current_ID")
	moduleCode := c.Query("module_code")
	keyword := c.Query("keyword")
	version := c.Query("version")

	type LicenseStat struct {
		License string `json:"license"`
		Count   int    `json:"count"`
	}

	type LicenseGroupStat struct {
		GroupName string `json:"group_name"`
		RiskLevel string `json:"risk_level"`
		Count     int    `json:"count"`
	}

	// 通过 PROJECT_ID 关联 project_inner 获取模块信息
	joinClause := "LEFT JOIN project_inner pi ON cs.PROJECT_ID = pi.id"

	// 构建WHERE条件
	whereConditions := []string{"cs.LICENSE IS NOT NULL", "cs.LICENSE != ''"}
	args := []interface{}{}

	// 通过 current_ID 查询
	if currentID != "" {
		whereConditions = append(whereConditions, "cs.current_id = ?")
		args = append(args, currentID)
	}

	// 通过 project_ids 过滤特定模块（PROJECT_ID 关联）
	if projectIDsStr != "" {
		var projectIDs []int64
		for _, idStr := range strings.Split(projectIDsStr, ",") {
			if id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64); err == nil {
				projectIDs = append(projectIDs, id)
			}
		}
		if len(projectIDs) > 0 {
			whereConditions = append(whereConditions, "cs.PROJECT_ID IN ?")
			args = append(args, projectIDs)
		}
	}

	// 按模块名查询
	if moduleCode != "" {
		whereConditions = append(whereConditions, "pi.module_code = ?")
		args = append(args, moduleCode)
	}

	// 按关键词查询
	if keyword != "" {
		whereConditions = append(whereConditions, "(cs._name LIKE ? OR cs._group LIKE ?)")
		args = append(args, "%"+keyword+"%", "%"+keyword+"%")
	}

	// 按版本查询
	if version != "" {
		whereConditions = append(whereConditions, "cs._version = ?")
		args = append(args, version)
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// 查询许可证统计
	licenseQuery := fmt.Sprintf(`
		SELECT cs.LICENSE as license, COUNT(*) as count 
		FROM component_snapshot cs %s
		WHERE %s
		GROUP BY cs.LICENSE
	`, joinClause, whereClause)

	// 查询分级统计
	groupQuery := fmt.Sprintf(`
		SELECT 
			group_name,
			risk_level,
			COUNT(*) as count
		FROM (
			SELECT 
				CASE 
					WHEN cs.LICENSE IN ('MIT', 'Apache-2.0', 'BSD-2-Clause', 'BSD-3-Clause', 'ISC', 'Unlicense', '0BSD', 'CC0-1.0') THEN 'Permissive'
					WHEN cs.LICENSE IN ('LGPL-2.1', 'LGPL-3.0', 'MPL-2.0', 'EPL-1.0', 'EPL-2.0', 'CDDL-1.0') THEN 'WeakCopyleft'
					WHEN cs.LICENSE IN ('GPL-2.0', 'GPL-3.0', 'AGPL-3.0', 'SSPL', 'OSL-3.0') THEN 'StrongCopyleft'
					ELSE '未分类'
				END as group_name,
				CASE 
					WHEN cs.LICENSE IN ('MIT', 'Apache-2.0', 'BSD-2-Clause', 'BSD-3-Clause', 'ISC', 'Unlicense', '0BSD', 'CC0-1.0') THEN '低风险'
					WHEN cs.LICENSE IN ('LGPL-2.1', 'LGPL-3.0', 'MPL-2.0', 'EPL-1.0', 'EPL-2.0', 'CDDL-1.0') THEN '中风险'
					WHEN cs.LICENSE IN ('GPL-2.0', 'GPL-3.0', 'AGPL-3.0', 'SSPL', 'OSL-3.0') THEN '高风险'
					ELSE '未分类'
				END as risk_level
			FROM component_snapshot cs %s
			WHERE %s
		) as license_groups
		GROUP BY group_name, risk_level
	`, joinClause, whereClause)

	var licenseStats []LicenseStat
	h.db.Raw(licenseQuery, args...).Scan(&licenseStats)

	var groupStats []LicenseGroupStat
	h.db.Raw(groupQuery, args...).Scan(&groupStats)

	response.Success(c, gin.H{
		"license_stats": licenseStats,
		"group_stats":   groupStats,
	})
}
