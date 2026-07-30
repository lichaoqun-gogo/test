package handlers

import (
	"osm-mngt/internal/models"
	"osm-mngt/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskHandler struct {
	db *gorm.DB
}

func NewRiskHandler(db *gorm.DB) *RiskHandler {
	return &RiskHandler{db: db}
}

func (h *RiskHandler) getUserInfo(c *gin.Context) (string, string) {
	userID := c.GetHeader("X-User-ID")
	userName := c.GetHeader("X-User-Name")
	return userID, userName
}

// ==================== 黑名单配置（component_core） ====================

func (h *RiskHandler) ListBlacklist(c *gin.Context) {
	page, pageSize := parsePagination(c)

	keyword := c.Query("keyword")
	formatFilter := c.Query("_format")
	levelFilter := c.Query("applicable_level")
	isBoundaryFilter := c.Query("is_boundary")
	alertTypeFilter := c.Query("alert_type")
	enableFilter := c.Query("_enable")

	query := h.db.Model(&models.ComponentCore{})

	if keyword != "" {
		query = query.Where("_name LIKE ? OR _group LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if formatFilter != "" {
		if v, err := strconv.Atoi(formatFilter); err == nil {
			query = query.Where("_format = ?", v)
		}
	}
	if levelFilter != "" {
		query = query.Where("applicable_level = ?", levelFilter)
	}
	if isBoundaryFilter != "" {
		if v, err := strconv.Atoi(isBoundaryFilter); err == nil {
			query = query.Where("is_boundary = ?", v)
		}
	}
	if alertTypeFilter != "" {
		if v, err := strconv.Atoi(alertTypeFilter); err == nil {
			query = query.Where("alert_type = ?", v)
		}
	}
	if enableFilter != "" {
		if v, err := strconv.Atoi(enableFilter); err == nil {
			query = query.Where("_enable = ?", v)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	var blacklist []models.ComponentCore
	offset := (page - 1) * pageSize
	if err := query.Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&blacklist).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	response.PageSuccess(c, blacklist, total, page, pageSize)
}

func (h *RiskHandler) CreateBlacklist(c *gin.Context) {
	var req models.ComponentCore
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if req.Name == "" || req.Group == "" || req.VersionExp == "" ||
		req.ApplicableLevel == "" || req.SupportTeam == "" {
		response.Error(c, 400, "缺少必填字段")
		return
	}

	userID, userName := h.getUserInfo(c)
	req.CreateBy = userID
	req.UpdateBy = userName
	req.CreateTime = time.Now()
	req.UpdateTime = time.Now()

	if err := h.db.Create(&req).Error; err != nil {
		response.Error(c, 500, "创建失败")
		return
	}

	response.Success(c, req)
}

func (h *RiskHandler) UpdateBlacklist(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的配置ID")
		return
	}

	var blacklist models.ComponentCore
	if err := h.db.First(&blacklist, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "配置不存在")
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

	if err := h.db.Model(&models.ComponentCore{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.Error(c, 500, "更新失败")
		return
	}

	response.Success(c, nil)
}

func (h *RiskHandler) DeleteBlacklist(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的配置ID")
		return
	}

	var blacklist models.ComponentCore
	if err := h.db.First(&blacklist, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "配置不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	if err := h.db.Delete(&blacklist).Error; err != nil {
		response.Error(c, 500, "删除失败")
		return
	}

	response.Success(c, nil)
}

// ==================== 黑名单阻断记录（component_black） ====================

func (h *RiskHandler) ListBlacklistRecords(c *gin.Context) {
	page, pageSize := parsePagination(c)

	projectIDFilter := c.Query("PROJECT_ID")
	processStatusFilter := c.Query("PROCESS_STATUS")
	alertTypeFilter := c.Query("alert_type")
	internalFilter := c.Query("internal")

	query := h.db.Model(&models.ComponentBlack{})

	if projectIDFilter != "" {
		if v, err := strconv.ParseInt(projectIDFilter, 10, 64); err == nil {
			query = query.Where("project_id = ?", v)
		}
	}
	if processStatusFilter != "" {
		query = query.Where("status = ?", processStatusFilter)
	}
	if alertTypeFilter != "" {
		if v, err := strconv.Atoi(alertTypeFilter); err == nil {
			query = query.Where("alert_type = ?", v)
		}
	}
	if internalFilter != "" {
		if v, err := strconv.Atoi(internalFilter); err == nil {
			query = query.Where("internal = ?", v)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	var records []models.ComponentBlack
	offset := (page - 1) * pageSize
	if err := query.Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	response.PageSuccess(c, records, total, page, pageSize)
}

func (h *RiskHandler) ProcessBlacklistRecord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的记录ID")
		return
	}

	var record models.ComponentBlack
	if err := h.db.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "记录不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	var req struct {
		ProcessStatus string `json:"process_status"`
		Reason        string `json:"reason"`
		JiraNo        string `json:"jira_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	_, userName := h.getUserInfo(c)
	updates := map[string]interface{}{
		"status":       req.ProcessStatus,
		"handle_plan":  req.Reason,
		"update_by":    userName,
		"update_time":  time.Now(),
	}
	if req.JiraNo != "" {
		updates["handle_result"] = req.JiraNo
	}

	if err := h.db.Model(&models.ComponentBlack{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.Error(c, 500, "处理失败")
		return
	}

	response.Success(c, nil)
}

// ==================== 许可证（license） ====================

func (h *RiskHandler) ListLicenses(c *gin.Context) {
	page, pageSize := parsePagination(c)

	keyword := c.Query("keyword")
	deprecatedFilter := c.Query("IS_DEPRECATED")
	osiFilter := c.Query("IS_OSIAPPROVED")

	query := h.db.Model(&models.License{})

	if keyword != "" {
		query = query.Where("LICENSE LIKE ? OR _name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if deprecatedFilter != "" {
		if v, err := strconv.Atoi(deprecatedFilter); err == nil {
			query = query.Where("IS_DEPRECATED = ?", v)
		}
	}
	if osiFilter != "" {
		if v, err := strconv.Atoi(osiFilter); err == nil {
			query = query.Where("IS_OSIAPPROVED = ?", v)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	var licenses []models.License
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&licenses).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	response.PageSuccess(c, licenses, total, page, pageSize)
}

func (h *RiskHandler) GetLicenseDetail(c *gin.Context) {
	id := c.Param("id")

	var license models.License
	if err := h.db.First(&license, id).Error; err != nil {
		response.Error(c, 404, "许可证不存在")
		return
	}

	// 查询许可证分类
	var groupLicense models.LicenseGroupLicense
	h.db.Where("license_id = ?", license.ID).First(&groupLicense)

	var groupName string
	if groupLicense.ID > 0 {
		var group models.LicenseGroup
		h.db.First(&group, groupLicense.LicenseGroupID)
		groupName = group.Name
	}

	response.Success(c, gin.H{
		"id":              license.ID,
		"LICENSE":         license.LicenseCode,
		"_name":           license.Name,
		"IS_OSIAPPROVED":  license.IsOSIApproved,
		"IS_DEPRECATED":   license.IsDeprecated,
		"FSFLIBRE":        license.FSFLibre,
		"_COMMENT":        license.Comment,
		"_text":           license.Text,
		"TEMPLATE":        license.Template,
		"HEADER":          license.Header,
		"license_group":   groupName,
		"create_time":     license.CreateTime,
		"update_time":     license.UpdateTime,
	})
}

// ==================== 许可证分组（licensegroup） ====================

func (h *RiskHandler) ListLicenseGroups(c *gin.Context) {
	var groups []models.LicenseGroup
	if err := h.db.Order("id ASC").Find(&groups).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	type GroupWithCount struct {
		models.LicenseGroup
		LicenseCount int64 `json:"licenseCount"`
	}

	result := make([]GroupWithCount, 0, len(groups))
	for _, group := range groups {
		var count int64
		h.db.Model(&models.LicenseGroupLicense{}).
			Where("licensegroup_id = ?", group.ID).
			Count(&count)

		result = append(result, GroupWithCount{
			LicenseGroup: group,
			LicenseCount: count,
		})
	}

	response.Success(c, result)
}

// ==================== 缺陷台账（defect_ledger） ====================

func (h *RiskHandler) ListDefects(c *gin.Context) {
	page, pageSize := parsePagination(c)

	keyword := c.Query("keyword")
	statusFilter := c.Query("status")
	severityFilter := c.Query("severity")
	sourceFilter := c.Query("source")

	query := h.db.Model(&models.DefectLedger{})

	if keyword != "" {
		query = query.Where("title LIKE ? OR affected_component LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}
	if severityFilter != "" {
		query = query.Where("risk_level = ?", severityFilter)
	}
	if sourceFilter != "" {
		query = query.Where("source = ?", sourceFilter)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	var defects []models.DefectLedger
	offset := (page - 1) * pageSize
	if err := query.Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&defects).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	response.PageSuccess(c, defects, total, page, pageSize)
}

func (h *RiskHandler) GetDefect(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的缺陷ID")
		return
	}

	var defect models.DefectLedger
	if err := h.db.First(&defect, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "缺陷不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, defect)
}

func (h *RiskHandler) CreateDefect(c *gin.Context) {
	var req models.DefectLedger
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if req.Title == "" || req.Source == "" {
		response.Error(c, 400, "缺少必填字段")
		return
	}

	if req.Status == "" {
		req.Status = "pending"
	}

	req.CreateTime = time.Now()
	req.UpdateTime = time.Now()

	if err := h.db.Create(&req).Error; err != nil {
		response.Error(c, 500, "创建失败")
		return
	}

	response.Success(c, req)
}

func (h *RiskHandler) UpdateDefectStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的缺陷ID")
		return
	}

	var defect models.DefectLedger
	if err := h.db.First(&defect, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "缺陷不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	var req struct {
		Status    string `json:"status"`
		Assignee  string `json:"assignee"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{
		"status":      req.Status,
		"update_time": time.Now(),
	}
	if req.Assignee != "" {
		updates["dispatcher"] = req.Assignee
		updates["dispatch_time"] = time.Now()
	}

	if err := h.db.Model(&models.DefectLedger{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.Error(c, 500, "更新失败")
		return
	}

	response.Success(c, nil)
}

// ==================== CICD检查 ====================

// CICDCheckItem CICD检查台账项
type CICDCheckItem struct {
	ID              int64     `json:"id"`
	ProCode         string    `json:"pro_code"`
	CurrentCommitID string    `json:"current_commit_id"`
	GroupName       string    `json:"group_name"`
	DeptName        string    `json:"dept_name"`
	DisposalType    string    `json:"disposalType"`
	Level           string    `json:"level"`
}

// CICDDetailItem CICD详情项
type CICDDetailItem struct {
	ModuleName    string `json:"moduleName"`
	CheckTime     string `json:"checkTime"`
	Disposal      string `json:"disposal"`
	BlacklistName string `json:"blacklistName"`
	Reason        string `json:"reason"`
}

func (h *RiskHandler) ListCicdCheck(c *gin.Context) {
	page, pageSize := parsePagination(c)

	projectName := c.Query("projectName")
	deptName := c.Query("deptName")
	groupName := c.Query("groupName")
	disposalType := c.Query("disposalType")
	level := c.Query("level")

	// 查询所有阻断记录
	query := h.db.Table("component_black cb").
		Select(`
			pi.id,
			pi.pro_code as pro_code,
			pi.current_commit_id,
			pi.group_name,
			pi.dept_name,
			cc.applicable_level as level,
			CASE 
				WHEN cc.applicable_level IN ('C', 'D') THEN '阻断'
				ELSE '警告'
			END as disposal_type
		`).
		Joins("LEFT JOIN project_inner pi ON cb.project_id = pi.id").
		Joins("LEFT JOIN component_core cc ON cb.blacklist_id = cc.id").
		Where("cb.status = ?", "待处理")

	if projectName != "" {
		query = query.Where("pi.pro_code LIKE ?", "%"+projectName+"%")
	}
	if deptName != "" {
		query = query.Where("pi.dept_name = ?", deptName)
	}
	if groupName != "" {
		query = query.Where("pi.group_name = ?", groupName)
	}
	if disposalType != "" {
		if disposalType == "阻断" {
			query = query.Where("cc.applicable_level IN ?", []string{"C", "D"})
		} else {
			query = query.Where("cc.applicable_level IN ?", []string{"A", "B"})
		}
	}
	if level != "" {
		query = query.Where("cc.applicable_level = ?", level)
	}

	var items []CICDCheckItem
	query.Order("pi.pro_code ASC").Find(&items)

	// 在Go代码中去重：按pro_code+current_commit_id唯一标识
	type ProjectKey struct {
		ProCode       string
		CommitID      string
	}
	seen := make(map[ProjectKey]bool)
	var uniqueItems []CICDCheckItem
	for _, item := range items {
		key := ProjectKey{item.ProCode, item.CurrentCommitID}
		if !seen[key] {
			seen[key] = true
			uniqueItems = append(uniqueItems, item)
		}
	}

	// 分页
	total := int64(len(uniqueItems))
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(uniqueItems) {
		uniqueItems = []CICDCheckItem{}
	} else if end > len(uniqueItems) {
		uniqueItems = uniqueItems[start:]
	} else {
		uniqueItems = uniqueItems[start:end]
	}

	response.Success(c, gin.H{
		"list":     uniqueItems,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *RiskHandler) GetCicdDetail(c *gin.Context) {
	commitID := c.Query("commitId")
	projectName := c.Query("projectName")

	// 查询该项目在该commit下的所有阻断记录
	// 主表：component_black（阻断记录）
	// 关联：component_core（获取_name黑名单名称、applicable_level等级、reason原因）
	// 关联：project_inner（获取pro_code模块名称）
	query := h.db.Table("component_black cb").
		Select(`
			pi.pro_code as module_name,
			cb.create_time as check_time,
			CASE 
				WHEN cc.applicable_level IN ('C', 'D') THEN '阻断'
				ELSE '警告'
			END as disposal,
			cc._name as blacklist_name,
			cc.reason
		`).
		Joins("LEFT JOIN project_inner pi ON cb.project_id = pi.id").
		Joins("LEFT JOIN component_core cc ON cb.blacklist_id = cc.id").
		Where("cb.project_name = ?", projectName)

	if commitID != "" {
		query = query.Where("pi.current_commit_id = ?", commitID)
	}

	var items []CICDDetailItem
	query.Order("cb.create_time DESC").Find(&items)

	response.Success(c, items)
}

// GetProjectBlockRecords 获取项目的阻断记录
func (h *RiskHandler) GetProjectBlockRecords(c *gin.Context) {
	projectID := c.Query("project_id")
	if projectID == "" {
		response.Error(c, 400, "缺少project_id参数")
		return
	}

	// 查询component_black表
	query := h.db.Table("component_black cb").
		Select(`
			cb.component_name,
			cb.component_version as version,
			cc.reason,
			cb.create_time
		`).
		Joins("LEFT JOIN component_core cc ON cb.blacklist_id = cc.id").
		Where("cb.project_id = ?", projectID)

	var records []map[string]interface{}
	if err := query.Order("cb.create_time DESC").Scan(&records).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, records)
}
