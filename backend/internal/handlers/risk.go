package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"osm-mngt/internal/models"
	"osm-mngt/internal/service/disposal"
	"osm-mngt/pkg/response"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RiskHandler struct {
	db     *gorm.DB
	engine *disposal.Engine
}

func NewRiskHandler(db *gorm.DB) *RiskHandler {
	return &RiskHandler{db: db, engine: disposal.NewEngine(db)}
}

func (h *RiskHandler) getUserInfo(c *gin.Context) (string, string) {
	userID := c.GetHeader("X-User-ID")
	userName := c.GetHeader("X-User-Name")
	return userID, userName
}

// getCurrentUser 返回解码后的用户ID、姓名、角色（缺陷处置流使用）
func (h *RiskHandler) getCurrentUser(c *gin.Context) (string, string, string) {
	userID := c.GetHeader("X-User-ID")
	userName := decodeHeaderValue(c.GetHeader("X-User-Name"))
	role := c.GetHeader("X-User-Role")
	return userID, userName, role
}

// ==================== 黑名单配置（component_core） ====================

func (h *RiskHandler) ListBlacklist(c *gin.Context) {
	page, pageSize := parsePagination(c)

	keyword := c.Query("keyword")
	typeFilter := c.Query("type")
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
	if typeFilter != "" {
		query = query.Where("type = ?", typeFilter)
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
	// level_time 数据库不允许零值（NO_ZERO_DATE），创建时默认当前时间
	req.LevelTime = time.Now()

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
		query = query.Where("process_status = ?", processStatusFilter)
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
		ProcessStatus int    `json:"process_status"`
		Reason        string `json:"reason"`
		JiraNo        string `json:"jira_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	_, userName := h.getUserInfo(c)
	updates := map[string]interface{}{
		"process_status": req.ProcessStatus,
		"reason":         req.Reason,
		"update_by":      userName,
		"update_time":    time.Now(),
	}
	if req.JiraNo != "" {
		updates["jira_no"] = req.JiraNo
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
	licensegroupFilter := c.Query("licensegroup")
	deprecatedFilter := c.Query("IS_DEPRECATED")
	osiFilter := c.Query("IS_OSIAPPROVED")

	query := h.db.Table("license l").
		Select(`l.ID, l.IS_DEPRECATED, l.FSFLIBRE, l.LICENSE, l._NAME, l.IS_OSIAPPROVED,
			g._NAME as LICENSEGROUP, g.RISKWEIGHT`).
		Joins("LEFT JOIN licensegroup_license gl ON l.ID = gl.LICENSE_ID").
		Joins("LEFT JOIN licensegroup g ON gl.LICENSEGROUP_ID = g.ID")

	if keyword != "" {
		query = query.Where("UPPER(l._NAME) LIKE UPPER(?) OR UPPER(l.LICENSE) LIKE UPPER(?)",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if licensegroupFilter != "" {
		query = query.Where("g._NAME = ?", licensegroupFilter)
	}
	if deprecatedFilter != "" {
		if v, err := strconv.Atoi(deprecatedFilter); err == nil {
			query = query.Where("l.IS_DEPRECATED = ?", v)
		}
	}
	if osiFilter != "" {
		if v, err := strconv.Atoi(osiFilter); err == nil {
			query = query.Where("l.IS_OSIAPPROVED = ?", v)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	var licenses []map[string]interface{}
	offset := (page - 1) * pageSize
	if err := query.Order("l._NAME ASC").Offset(offset).Limit(pageSize).Scan(&licenses).Error; err != nil {
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

// ==================== 缺陷台账（外部实体 defect ∪ 内部实体 defect_ledger，LEFT JOIN flow 补流程明细） ====================

// defectLedgerRow 缺陷台账统一行：内部缺陷（风险上报）与外部缺陷（采集台账）共用展示结构
type defectLedgerRow struct {
	ID                    int64      `json:"id"`
	DefectNo              string     `json:"defectNo"`
	Source                string     `json:"source"`
	Title                 string     `json:"title"`
	RiskLevel             string     `json:"riskLevel"`
	AffectedComponent     string     `json:"affectedComponent"`
	AffectedComponentType string     `json:"affectedComponentType"`
	AffectedVersion       string     `json:"affectedVersion"`
	CVEID                 string     `json:"cveId"`
	Description           string     `json:"description"`
	DiscoveredDate        *time.Time `json:"discoveredDate"`
	Status                string     `json:"status"`
	FlowType              string     `json:"flowType"`
	MaintDept             string     `json:"maintDept"`
	Deadline              *time.Time `json:"deadline"`
	DispatchTarget        string     `json:"dispatchTarget"`
	DispatchTargetName    string     `json:"dispatchTargetName"`
	// 分配类型：analyze=已分配研判人；import=导入跟踪人（不算已分配研判人）；空=未分配
	DispatchType          string     `json:"dispatchType"`
	AnalyzeResult         string     `json:"analyzeResult"`
	ReviewMeetingNo       string     `json:"reviewMeetingNo"`
	ReviewConclusion      string     `json:"reviewConclusion"`
	DispatchDept          string     `json:"dispatchDept"`
	DispatchTime          *time.Time `json:"dispatchTime"`
	Dispatcher            string     `json:"dispatcher"`
	AnalyzeBy             string     `json:"analyzeBy"`
	AnalyzeTime           *time.Time `json:"analyzeTime"`
	ReviewBy              string     `json:"reviewBy"`
	ReviewTime            *time.Time `json:"reviewTime"`
	CloseTime             *time.Time `json:"closeTime"`
	// 外部缺陷专有（详情/列表展示）
	Software    string `json:"software"`
	IssueURL    string `json:"issueUrl"`
	TitleZh     string `json:"titleZh"`     // 中文标题（采集工具翻译产物）
	SummaryZh   string `json:"summaryZh"`   // 中文释义
	RepairAdvice string `json:"repairAdvice"` // 中文建议
}

// defectLedgerUnionSQL 缺陷台账合并结果集（外部 defect ∪ 内部 defect_ledger，LEFT JOIN flow 补流程明细）
// 统一列：defect_no, source, title, risk_level, affected_component, discovered_date, status, id, risk_type,
// affected_version, cve_id, description, component_type, flow_type, maint_dept, dispatch_target, deadline,
// analyze_result, review_meeting_no, review_conclusion, dispatch_dept, dispatch_time, dispatcher, analyze_by,
// analyze_time, review_by, review_time, close_time
func defectLedgerUnionSQL() string {
	return `(SELECT dl.defect_no, dl.source, dl.title, dl.risk_level, NULL AS affected_component,
        dl.discovered_date, dl.internal_status AS status, dl.id,
        NULL AS risk_type, NULL AS affected_version, NULL AS cve_id, dl.description, NULL AS component_type,
        df.flow_type, df.maint_dept, df.dispatch_type, df.dispatch_target, df.deadline, df.analyze_result, df.review_meeting_no, df.review_conclusion,
        df.dispatch_dept, df.dispatch_time, df.dispatcher, df.analyze_by, df.analyze_time, df.review_by, df.review_time, df.close_time,
        NULL AS title_zh, NULL AS summary_zh, NULL AS repair_advice
   FROM defect_ledger dl LEFT JOIN defect_flow df ON df.defect_no = dl.defect_no)
  UNION ALL
  (SELECT d.defect_id, d.source, d.title, d.severity AS risk_level, d.software AS affected_component,
        d.created_at AS discovered_date, d.internal_status AS status, NULL AS id,
        NULL AS risk_type, NULL AS affected_version, NULL AS cve_id, d.body AS description, 'external' AS component_type,
        df.flow_type, df.maint_dept, df.dispatch_type, df.dispatch_target, df.deadline, df.analyze_result, df.review_meeting_no, df.review_conclusion,
        df.dispatch_dept, df.dispatch_time, df.dispatcher, df.analyze_by, df.analyze_time, df.review_by, df.review_time, df.close_time,
        d.title_zh AS title_zh, d.summary_zh AS summary_zh, d.repair_advice AS repair_advice
   FROM defect d
   LEFT JOIN defect_flow df ON df.defect_no = d.defect_id)`
}

// defectLedgerPermission 缺陷台账权限过滤（内部/外部/合并视图共用）
// admin / 管理部门专员（全局治理视图）看全部；其他专员：本部门运维的缺陷，或本部门有处置单的缺陷（跨部门裂变后使用部门专员也可看到并改派本部门处置单）
// scope=dept 表示个人工作台视图：管理部门专员也与部门专员一致按部门过滤，仅在全局治理（风险管理-缺陷台账）看全景
// defectNoColumn 为「缺陷编号列」的显式前缀（内部=dl.defect_no，外部=d.defect_id，合并=t.defect_no），避免 JOIN 后列歧义
func (h *RiskHandler) defectLedgerPermission(c *gin.Context, deptColumn, defectNoColumn string) (string, []interface{}) {
	where := []string{"1 = 1"}
	var args []interface{}
	userID, _, role := h.getCurrentUser(c)
	userDept := decodeHeaderValue(c.GetHeader("X-User-Dept"))
	isAdmin := role == "admin"
	isMgmt := h.engine.IsMgmtSpecialist(userID, "")
	scope := c.Query("scope")
	mine := c.Query("mine")
	if mine == "analysis" {
		// 待我研判：真正被分配研判（dispatch_type='analyze'）且待研判的缺陷（个人待办，不按部门过滤）；
		// import 跟踪人只是归属记录，不算研判人，不进入待办
		where = append(where, "dispatch_type = 'analyze' AND dispatch_target = ? AND status = ?")
		args = append(args, userID, "pending_analysis")
	} else if isAdmin || (isMgmt && scope != "dept") {
		// 全局治理视图：admin / 管理部门专员看全部
	} else if userDept == "" {
		// 未获取到部门信息时，不返回任何缺陷数据，避免数据越权
		where = append(where, "1 = 0")
	} else {
		where = append(where, "("+deptColumn+" = ? OR "+defectNoColumn+" IN (SELECT defect_no FROM defect_disposal WHERE dept_name = ?))")
		args = append(args, userDept, userDept)
	}
	return strings.Join(where, " AND "), args
}

// defectListFilters 缺陷台账通用筛选（keyword/status/severity/source）
func defectListFilters(keyword, statusFilter, severityFilter, sourceFilter, keywordCol string) (string, []interface{}) {
	var where []string
	var args []interface{}
	if keyword != "" {
		where = append(where, "(title LIKE ? OR "+keywordCol+" LIKE ?)")
		args = append(args, "%"+keyword+"%", "%"+keyword+"%")
	}
	if statusFilter != "" {
		where = append(where, "status = ?")
		args = append(args, statusFilter)
	}
	if severityFilter != "" {
		where = append(where, "risk_level = ?")
		args = append(args, severityFilter)
	}
	if sourceFilter != "" {
		where = append(where, "source = ?")
		args = append(args, sourceFilter)
	}
	return strings.Join(where, " AND "), args
}

// fillDispatchTargetNames 补充研判人姓名（username -> nickname）
func fillDispatchTargetNames(db *gorm.DB, defects []defectLedgerRow) {
	targets := make([]string, 0, len(defects))
	for _, d := range defects {
		if d.DispatchTarget != "" {
			targets = append(targets, d.DispatchTarget)
		}
	}
	if len(targets) == 0 {
		return
	}
	var users []models.User
	if err := db.Where("username IN ?", targets).Find(&users).Error; err != nil {
		return
	}
	nameMap := map[string]string{}
	for _, u := range users {
		nameMap[u.Username] = u.Nickname
	}
	for i := range defects {
		defects[i].DispatchTargetName = nameMap[defects[i].DispatchTarget]
	}
}

// ListDefects 缺陷台账合并列表（内部 ∪ 外部，供待我研判等跨类型场景）
func (h *RiskHandler) ListDefects(c *gin.Context) {
	page, pageSize := parsePagination(c)
	wherePerm, permArgs := h.defectLedgerPermission(c, "maint_dept", "t.defect_no")
	filterWhere, filterArgs := defectListFilters(c.Query("keyword"), c.Query("status"), c.Query("severity"), c.Query("source"), "affected_component")
	where := []string{wherePerm}
	args := append(permArgs, filterArgs...)
	if filterWhere != "" {
		where = append(where, filterWhere)
	}
	whereStr := strings.Join(where, " AND ")
	union := defectLedgerUnionSQL()

	var total int64
	if err := h.db.Raw("SELECT COUNT(*) FROM ("+union+") t WHERE "+whereStr, args...).Scan(&total).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}
	offset := (page - 1) * pageSize
	var defects []defectLedgerRow
	pageArgs := append(append([]interface{}{}, args...), pageSize, offset)
	if err := h.db.Raw("SELECT t.* FROM ("+union+") t WHERE "+whereStr+" ORDER BY t.discovered_date DESC, t.defect_no DESC LIMIT ? OFFSET ?", pageArgs...).Scan(&defects).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}
	fillDispatchTargetNames(h.db, defects)
	response.PageSuccess(c, defects, total, page, pageSize)
}

// ListInternalDefects 内部缺陷列表（风险上报那套字段，查 defect_ledger + flow）—— 缺陷台账「内部缺陷」tab
func (h *RiskHandler) ListInternalDefects(c *gin.Context) {
	page, pageSize := parsePagination(c)
	wherePerm, permArgs := h.defectLedgerPermission(c, "df.maint_dept", "dl.defect_no")
	filterWhere, filterArgs := defectListFilters(c.Query("keyword"), c.Query("status"), c.Query("severity"), c.Query("source"), "dl.title")
	where := []string{wherePerm}
	args := append(permArgs, filterArgs...)
	if filterWhere != "" {
		where = append(where, filterWhere)
	}
	whereStr := strings.Join(where, " AND ")
	base := "FROM defect_ledger dl LEFT JOIN defect_flow df ON df.defect_no = dl.defect_no"

	var total int64
	if err := h.db.Raw("SELECT COUNT(*) "+base+" WHERE "+whereStr, args...).Scan(&total).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}
	offset := (page - 1) * pageSize
	var defects []defectLedgerRow
	pageArgs := append(append([]interface{}{}, args...), pageSize, offset)
	selectCols := `dl.id, dl.defect_no, dl.source, dl.title, dl.risk_level, NULL AS affected_component, NULL AS affected_component_type, NULL AS affected_version, NULL AS cve_id, dl.description, dl.discovered_date, dl.internal_status AS status, df.flow_type, df.maint_dept, df.deadline, df.dispatch_type, df.dispatch_target, df.analyze_result, df.review_meeting_no, df.review_conclusion, df.dispatch_dept, df.dispatch_time, df.dispatcher, df.analyze_by, df.analyze_time, df.review_by, df.review_time, df.close_time, NULL AS software, NULL AS issue_url`
	if err := h.db.Raw("SELECT "+selectCols+" "+base+" WHERE "+whereStr+" ORDER BY dl.create_time DESC LIMIT ? OFFSET ?", pageArgs...).Scan(&defects).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}
	fillDispatchTargetNames(h.db, defects)
	response.PageSuccess(c, defects, total, page, pageSize)
}

// ListExternalDefects 外部缺陷列表（采集台账，查 defect + flow）—— 缺陷台账「外部缺陷」tab
func (h *RiskHandler) ListExternalDefects(c *gin.Context) {
	page, pageSize := parsePagination(c)
	wherePerm, permArgs := h.defectLedgerPermission(c, "d.maint_dept", "d.defect_id")
	filterWhere, filterArgs := defectListFilters(c.Query("keyword"), c.Query("status"), c.Query("severity"), c.Query("source"), "d.title")
	where := []string{wherePerm}
	args := append(permArgs, filterArgs...)
	if filterWhere != "" {
		where = append(where, filterWhere)
	}
	whereStr := strings.Join(where, " AND ")
	base := "FROM defect d LEFT JOIN defect_flow df ON df.defect_no = d.defect_id"

	var total int64
	if err := h.db.Raw("SELECT COUNT(*) "+base+" WHERE "+whereStr, args...).Scan(&total).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}
	offset := (page - 1) * pageSize
	var defects []defectLedgerRow
	pageArgs := append(append([]interface{}{}, args...), pageSize, offset)
	selectCols := `NULL AS id, d.defect_id AS defect_no, d.source, d.title, d.severity AS risk_level, d.software AS affected_component, 'external' AS affected_component_type, NULL AS affected_version, NULL AS cve_id, d.body AS description, d.created_at AS discovered_date, d.internal_status AS status, df.flow_type, d.maint_dept, df.deadline, df.dispatch_type, df.dispatch_target, df.analyze_result, df.review_meeting_no, df.review_conclusion, df.dispatch_dept, df.dispatch_time, df.dispatcher, df.analyze_by, df.analyze_time, df.review_by, df.review_time, df.close_time, d.software, d.issue_url, d.title_zh, d.summary_zh, d.repair_advice`
	if err := h.db.Raw("SELECT "+selectCols+" "+base+" WHERE "+whereStr+" ORDER BY d.imported_at DESC, d.defect_id DESC LIMIT ? OFFSET ?", pageArgs...).Scan(&defects).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}
	fillDispatchTargetNames(h.db, defects)
	response.PageSuccess(c, defects, total, page, pageSize)
}

// ==================== 缺陷修复效能看板 ====================

// defectSeverityPriority 外部缺陷严重等级展示顺序（看板固定轴，值为 0 也保留）
// 外部缺陷 severity 为导入原始中文值（高危/中危/低危/提示/未知），直接展示
var defectSeverityPriority = []string{"高危", "中危", "低危", "提示", "未知"}

// defectStatusRank 缺陷内部状态等级（就高不就低：多人研判时以更高结论为准，用于导入预判升级判断）
// 等级从低到高：待研判(1) < 忽略(2) < 待上会(3) < 已分配(4) < 已闭环(5)
func defectStatusRank(status string) int {
	switch status {
	case disposal.StatusIgnored:
		return 2
	case disposal.StatusPendingReview:
		return 3
	case disposal.StatusAssigned:
		return 4
	case disposal.StatusClosed:
		return 5
	default: // pending_analysis 及未知
		return 1
	}
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}

// GetDefectStatsOverview 缺陷修复效能看板聚合接口（仅统计外部缺陷 defect 表及其处置工单）
// 参数：start/end（YYYY-MM-DD，可选，作用于"期间"指标）、dept（部门，可选）
// 口径说明：
//   - 存量指标（进行中/逾期未闭环）：只看当前状态，不受时间范围影响，仅受部门筛选影响
//   - 期间指标（新增缺陷/已完成处置/MTTR/闭环率/严重等级分布）：按时间范围统计
//   - 缺陷级：仅外部缺陷（defect），新增缺陷按上游创建时间 created_at，运维部门取实体表 maint_dept
//   - 工单级：仅外部缺陷关联的处置单（defect_no ∈ defect）
//   - MTTR 从处置单派发时间（create_time）到完成时间（completed_at），仅统计已闭环工单
//   - 闭环率（本期）= 本期已完成处置 / 本期新派发处置 × 100
func (h *RiskHandler) GetDefectStatsOverview(c *gin.Context) {
	now := time.Now()

	// 时间范围（闭区间，end 取到当天 23:59:59）
	var start, endExclusive time.Time
	var hasStart, hasEnd bool
	if s := c.Query("start"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			start, hasStart = t, true
		}
	}
	if e := c.Query("end"); e != "" {
		if t, err := time.ParseInLocation("2006-01-02", e, time.Local); err == nil {
			endExclusive, hasEnd = t.AddDate(0, 0, 1), true
		}
	}
	dept := c.Query("dept")

	// ---------- 处置单（工单级，存量口径；仅外部缺陷工单） ----------
	dq := h.db.Model(&models.DefectDisposal{}).
		Where("defect_no IN (SELECT defect_id FROM defect)")
	if dept != "" {
		dq = dq.Where("dept_name = ?", dept)
	}
	var disposals []models.DefectDisposal
	if err := dq.Find(&disposals).Error; err != nil {
		response.Error(c, 500, "统计失败")
		return
	}

	type deptAgg struct {
		total      int
		closed     int
		processing int
		overdue    int
		mttrSum    time.Duration
	}
	agg := map[string]*deptAgg{}
	overallProcessing, overallOverdue, overallClosed := 0, 0, 0
	overallMttrSum := time.Duration(0)
	overallClosedCount := 0

	for i := range disposals {
		d := &disposals[i]
		name := d.DeptName
		if name == "" {
			name = "未指定"
		}
		a := agg[name]
		if a == nil {
			a = &deptAgg{}
			agg[name] = a
		}
		a.total++
		switch d.Status {
		case disposal.DisposalCompleted:
			a.closed++
			overallClosed++
			if d.CompletedAt != nil && !d.CreateTime.IsZero() {
				if diff := d.CompletedAt.Sub(d.CreateTime); diff > 0 {
					a.mttrSum += diff
					overallMttrSum += diff
					overallClosedCount++
				}
			}
		case disposal.DisposalPendingAssign, disposal.DisposalProcessing:
			if d.Deadline != nil && d.Deadline.Before(now) {
				a.overdue++
				overallOverdue++
			} else {
				a.processing++
				overallProcessing++
			}
		}
	}

	// 部门修复效能明细（全量存量口径）
	type deptStat struct {
		DeptName   string  `json:"deptName"`
		Total      int     `json:"total"`
		Closed     int     `json:"closed"`
		Processing int     `json:"processing"`
		Overdue    int     `json:"overdue"`
		Mttr       float64 `json:"mttr"`
		CloseRate  float64 `json:"closeRate"`
	}
	deptList := make([]deptStat, 0, len(agg))
	for name, a := range agg {
		mttr := 0.0
		if a.closed > 0 {
			mttr = round1(a.mttrSum.Hours() / float64(a.closed))
		}
		closeRate := 0.0
		if a.total > 0 {
			closeRate = round1(float64(a.closed) / float64(a.total) * 100)
		}
		deptList = append(deptList, deptStat{
			DeptName: name, Total: a.total, Closed: a.closed,
			Processing: a.processing, Overdue: a.overdue,
			Mttr: mttr, CloseRate: closeRate,
		})
	}
	// 逾期多、盘子大的部门排前面，突出欠账
	sort.Slice(deptList, func(i, j int) bool {
		if deptList[i].Overdue != deptList[j].Overdue {
			return deptList[i].Overdue > deptList[j].Overdue
		}
		return deptList[i].Total > deptList[j].Total
	})

	// 处理情况分布（三态：进行中/逾期未闭环/已完成处置）
	statusDist := []struct {
		Key   string `json:"key"`
		Name  string `json:"name"`
		Count int    `json:"count"`
	}{
		{Key: "processing", Name: "进行中", Count: overallProcessing},
		{Key: "overdue", Name: "逾期未闭环", Count: overallOverdue},
		{Key: "completed", Name: "已完成处置", Count: overallClosed},
	}

	// 部门下拉选项（按处置单数降序）
	deptOptions := []string{}
	{
		var rows []struct {
			DeptName string `gorm:"column:dept_name"`
		}
		h.db.Model(&models.DefectDisposal{}).Select("dept_name").
			Where("dept_name <> '' AND defect_no IN (SELECT defect_id FROM defect)").
			Group("dept_name").
			Order("COUNT(*) DESC").Scan(&rows)
		for _, r := range rows {
			deptOptions = append(deptOptions, r.DeptName)
		}
	}

	// ---------- 缺陷级（期间口径：仅外部缺陷，按上游创建时间，运维部门取实体表 maint_dept） ----------
	baseDefectQ := func() *gorm.DB {
		q := h.db.Table("defect d")
		if hasStart {
			q = q.Where("d.created_at >= ?", start)
		}
		if hasEnd {
			q = q.Where("d.created_at < ?", endExclusive)
		}
		if dept != "" {
			q = q.Where("d.maint_dept = ?", dept)
		}
		return q
	}

	var newDefects int64
	if err := baseDefectQ().Count(&newDefects).Error; err != nil {
		response.Error(c, 500, "统计失败")
		return
	}

	// 严重等级分布
	type sevRow struct {
		Level string `gorm:"column:level"`
		Count int64  `gorm:"column:count"`
	}
	var sevRows []sevRow
	baseDefectQ().Select("d.severity AS level, COUNT(*) AS count").
		Group("d.severity").Scan(&sevRows)
	sevMap := map[string]int64{}
	for _, r := range sevRows {
		sevMap[r.Level] = r.Count
	}
	severity := make([]struct {
		Level string `json:"level"`
		Count int64  `json:"count"`
	}, 0, len(defectSeverityPriority))
	for _, lv := range defectSeverityPriority {
		severity = append(severity, struct {
			Level string `json:"level"`
			Count int64  `json:"count"`
		}{Level: lv, Count: sevMap[lv]})
	}

	// 期间：本期已完成处置 / 本期新派发处置（仅外部缺陷工单）
	baseDisposalQ := func() *gorm.DB {
		q := h.db.Model(&models.DefectDisposal{}).
			Where("defect_no IN (SELECT defect_id FROM defect)")
		if hasStart {
			q = q.Where("completed_at >= ?", start)
		}
		if hasEnd {
			q = q.Where("completed_at < ?", endExclusive)
		}
		if dept != "" {
			q = q.Where("dept_name = ?", dept)
		}
		return q
	}
	var completedInRange int64
	if err := baseDisposalQ().Where("status = ?", disposal.DisposalCompleted).Count(&completedInRange).Error; err != nil {
		response.Error(c, 500, "统计失败")
		return
	}

	// 公司平均 MTTR（期间，从派发 create_time 到 completed_at）
	var completedRows []models.DefectDisposal
	baseDisposalQ().Where("status = ?", disposal.DisposalCompleted).Find(&completedRows)
	mttrSum := time.Duration(0)
	mttrCount := 0
	for i := range completedRows {
		d := &completedRows[i]
		if d.CompletedAt != nil && !d.CreateTime.IsZero() {
			if diff := d.CompletedAt.Sub(d.CreateTime); diff > 0 {
				mttrSum += diff
				mttrCount++
			}
		}
	}
	mttr := 0.0
	if mttrCount > 0 {
		mttr = round1(mttrSum.Hours() / float64(mttrCount))
	}

	// 本期新派发处置单（create_time 在范围内；仅外部缺陷工单）
	newDisposalQ := h.db.Model(&models.DefectDisposal{}).
		Where("defect_no IN (SELECT defect_id FROM defect)")
	if hasStart {
		newDisposalQ = newDisposalQ.Where("create_time >= ?", start)
	}
	if hasEnd {
		newDisposalQ = newDisposalQ.Where("create_time < ?", endExclusive)
	}
	if dept != "" {
		newDisposalQ = newDisposalQ.Where("dept_name = ?", dept)
	}
	var newDisposals int64
	if err := newDisposalQ.Count(&newDisposals).Error; err != nil {
		response.Error(c, 500, "统计失败")
		return
	}
	closeRate := 0.0
	if newDisposals > 0 {
		closeRate = round1(float64(completedInRange) / float64(newDisposals) * 100)
	}

	response.Success(c, gin.H{
		"kpi": gin.H{
			"mttr":       mttr,       // 公司平均 MTTR（天，期间）
			"overdue":    overallOverdue, // 逾期未闭环工单（存量）
			"newDefects": newDefects, // 新增缺陷（期间）
			"completed":  completedInRange, // 已完成处置（期间）
			"closeRate":  closeRate,  // 闭环率 %（期间）
		},
		"severity": severity,
		"status":   statusDist,
		"dept":     deptList,
		"depts":    deptOptions,
	})
}

func (h *RiskHandler) GetDefect(c *gin.Context) {
	defectNo := c.Param("id")
	if defectNo == "" {
		response.Error(c, 400, "无效的缺陷编号")
		return
	}

	var row defectLedgerRow
	if err := h.db.Raw("SELECT t.* FROM ("+defectLedgerUnionSQL()+") t WHERE t.defect_no = ?", defectNo).Scan(&row).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}
	if row.DefectNo == "" {
		response.Error(c, 404, "缺陷不存在")
		return
	}
	response.Success(c, row)
}

func (h *RiskHandler) CreateDefect(c *gin.Context) {
	// 内部缺陷（风险上报）：字段按风险上报那套（标题/描述/等级选填），上报后进入「待审核」，由开源治理管理员审核分配
	var req struct {
		DefectNo    string `json:"defectNo"`
		Source      string `json:"source"`
		Title       string `json:"title"`
		RiskLevel   string `json:"riskLevel"` // 选填
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if req.Title == "" || req.Source == "" {
		response.Error(c, 400, "缺少必填字段")
		return
	}
	if req.RiskLevel == "" {
		req.RiskLevel = "未知"
	}
	if req.DefectNo == "" {
		req.DefectNo = fmt.Sprintf("DF-%d", time.Now().UnixNano())
	}
	now := time.Now()

	// 内部缺陷实体（风险上报：初始状态「待审核」pending_review，写实体表 internal_status）
	entity := models.DefectLedger{
		DefectNo:       req.DefectNo,
		Source:         req.Source,
		Title:          req.Title,
		RiskLevel:      req.RiskLevel,
		Description:    req.Description,
		DiscoveredDate: now,
		InternalStatus: disposal.StatusPendingReview,
		CreateTime:     now,
		UpdateTime:     now,
	}
	if err := h.db.Create(&entity).Error; err != nil {
		response.Error(c, 500, "创建失败")
		return
	}

	// 缺陷级处置流程明细（状态在实体表，flow 只存明细；运维部门默认上报者部门，审核分配时由管理员指定覆盖）
	flow := models.DefectFlow{
		DefectNo:   req.DefectNo,
		FlowType:   "defect",
		MaintDept:  decodeHeaderValue(c.GetHeader("X-User-Dept")),
		CreateTime: now,
		UpdateTime: now,
	}
	if err := h.db.Create(&flow).Error; err != nil {
		response.Error(c, 500, "创建失败")
		return
	}

	response.Success(c, gin.H{"defectNo": entity.DefectNo, "id": entity.ID})
}

// ==================== 外部缺陷导入（采集工具产物，幂等） ====================

// defectImportItem 采集工具输出的单条缺陷（字段与 defect 表一一对应，defect_id 为 upsert 业务键）
type defectImportItem struct {
	DefectID      string          `json:"defect_id"`
	Software      string          `json:"software"`
	IssueURL      string          `json:"issue_url"`
	Title         string          `json:"title"`
	Body          string          `json:"body"`
	Labels        json.RawMessage `json:"labels"`
	Reporter      string          `json:"reporter"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	CommentsCount int             `json:"comments_count"`
	Source        string          `json:"source"`
	Severity      string          `json:"severity"`
	IssueState    string          `json:"issue_state"`
	StateReason   string          `json:"state_reason"`
	Reactions     json.RawMessage `json:"reactions"`
	Milestone     string          `json:"milestone"`
	WatchMatches  json.RawMessage `json:"watch_matches"`
	SummaryZh     string          `json:"summary_zh"`
	RepairAdvice  string          `json:"repair_advice"`
	TitleZh       string          `json:"title_zh"` // 中文标题（采集工具翻译产物）
	ImportedAt    string          `json:"imported_at"`
	DisposalStatus string         `json:"disposal_status"` // 研判状态（导入预判，可选）: pending_review-待上会, ignored-忽略; 缺省=待研判
}

// ImportDefects 导入外部缺陷（采集工具产物，如 redis-all.json 的 defects 数组）
// 幂等：按 defect_id（软件#编号，主键）upsert，同一编号只有一条记录，无论谁导入
// 状态保护：已存在的缺陷只更新采集信息（title/body/severity 等），内部处置状态 internal_status 绝不回退
func (h *RiskHandler) ImportDefects(c *gin.Context) {
	var payload struct {
		Defects []defectImportItem `json:"defects"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		// 兼容直接传数组
		var arr []defectImportItem
		if err2 := c.ShouldBindJSON(&arr); err2 != nil {
			response.Error(c, 400, "参数错误: "+err.Error())
			return
		}
		payload.Defects = arr
	}
	if len(payload.Defects) == 0 {
		response.Error(c, 400, "defects 为空")
		return
	}

	now := time.Now()
	parseTime := func(s string) *time.Time {
		if s == "" {
			return nil
		}
		for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
			if t, e := time.Parse(layout, s); e == nil {
				return &t
			}
		}
		return nil
	}
	toJSON := func(raw json.RawMessage) string {
		if len(raw) == 0 {
			return "[]" // JSON 列不允许空串，缺省给空数组
		}
		return string(raw)
	}

	var imported, updated, skipped int
	var errs []string

	// 导入者信息：跟踪人（研判归属，统一存 username）+ 所属部门（首次导入记录部门，重复导入不覆盖）
	userID, _, _ := h.getCurrentUser(c)
	userDept := decodeHeaderValue(c.GetHeader("X-User-Dept"))

	for _, it := range payload.Defects {
		if strings.TrimSpace(it.DefectID) == "" || strings.TrimSpace(it.Software) == "" || strings.TrimSpace(it.Title) == "" {
			skipped++
			errs = append(errs, "跳过缺少 defect_id/software/title 的记录")
			continue
		}

		// 导入预判状态（disposal_status）：pending_review-待上会 / ignored-忽略 / 缺省=待研判
		targetStatus := disposal.StatusPendingAnalysis
		preJudged := false
		switch it.DisposalStatus {
		case "pending_review", "review", "need_review":
			targetStatus = disposal.StatusPendingReview
			preJudged = true
		case "ignored", "ignore":
			targetStatus = disposal.StatusIgnored
			preJudged = true
		}

		d := models.Defect{
			DefectID:       it.DefectID,
			Software:       it.Software,
			IssueURL:       it.IssueURL,
			Title:          it.Title,
			Body:           it.Body,
			Labels:         toJSON(it.Labels),
			Reporter:       it.Reporter,
			CreatedAt:      parseTime(it.CreatedAt),
			UpdatedAt:      parseTime(it.UpdatedAt),
			CommentsCount:  it.CommentsCount,
			Source:         it.Source,
			Severity:       it.Severity,
			IssueState:     it.IssueState,
			StateReason:    it.StateReason,
			Reactions:      toJSON(it.Reactions),
			Milestone:      it.Milestone,
			WatchMatches:   toJSON(it.WatchMatches),
			SummaryZh:      it.SummaryZh,
			RepairAdvice:   it.RepairAdvice,
			TitleZh:        it.TitleZh,
			MaintDept:      userDept,     // 首次导入记录导入者部门（upsert 更新列不含它 → 不覆盖）
			InternalStatus: targetStatus, // 新记录按导入预判状态（upsert 更新列不含它 → 已存在不回退）
			ImportedAt:     now,
		}
		if d.Source == "" {
			d.Source = "github-issue"
		}
		if d.Severity == "" {
			d.Severity = "未知"
		}
		if d.IssueState == "" {
			d.IssueState = "open"
		}
		if t := parseTime(it.ImportedAt); t != nil {
			d.ImportedAt = *t
		}

		res := h.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "defect_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"software", "issue_url", "title", "body", "labels", "reporter", "created_at", "updated_at", "comments_count", "source", "severity", "issue_state", "state_reason", "reactions", "milestone", "watch_matches", "summary_zh", "repair_advice", "title_zh", "imported_at"}),
		}).Create(&d)
		if res.Error != nil {
			errs = append(errs, it.DefectID+": "+res.Error.Error())
			continue
		}

		// 状态边界 B（就高不就低）：导入预判时，新预判等级高于当前状态等级才升级——
		// 待研判可被忽略/待上会覆盖，忽略可被待上会覆盖（多人研判以更高结论为准）；
		// 绝不降级/回退（已分配/已闭环等流转状态不受导入预判影响）
		if preJudged && targetStatus != disposal.StatusPendingAnalysis {
			var cur models.Defect
			h.db.Select("internal_status").Where("defect_id = ?", d.DefectID).First(&cur)
			if defectStatusRank(targetStatus) > defectStatusRank(cur.InternalStatus) {
				h.db.Model(&models.Defect{}).Where("defect_id = ?", d.DefectID).Update("internal_status", targetStatus)
			}
		}

		// 导入即研判：确保缺陷级处置流程明细存在；归属跟踪人（导入者 username）；预判时记录研判人/时间/结果
		h.db.Exec("INSERT INTO `defect_flow` (`defect_no`, `flow_type`, `maint_dept`, `create_time`, `update_time`) VALUES (?, 'defect', ?, NOW(), NOW()) ON DUPLICATE KEY UPDATE `defect_no` = `defect_no`", d.DefectID, userDept)
		if userID != "" {
			h.db.Exec("UPDATE `defect_flow` SET `dispatch_target` = ?, `dispatch_type` = 'import', `dispatch_time` = NOW(), `dispatcher` = ? WHERE `defect_no` = ? AND (`dispatch_target` IS NULL OR `dispatch_target` = '')", userID, userID, d.DefectID)
			if preJudged {
				h.db.Exec("UPDATE `defect_flow` SET `analyze_by` = ?, `analyze_time` = NOW(), `analyze_result` = ? WHERE `defect_no` = ? AND (`analyze_by` IS NULL OR `analyze_by` = '')", userID, "导入预判："+it.DisposalStatus, d.DefectID)
			}
		}

		switch res.RowsAffected {
		case 1:
			imported++
		case 2:
			updated++
		default:
			skipped++
		}
	}

	response.Success(c, gin.H{"imported": imported, "updated": updated, "skipped": skipped, "errors": errs})
}

// ==================== 漏洞情报（临时安全情报 → 直接裂变处置工单） ====================

var errVulnNoExists = errors.New("漏洞编号已存在")

// canManageVulnIntel 漏洞情报录入权限：开源治理专员 / 开源治理管理员 / 系统管理员
func (h *RiskHandler) canManageVulnIntel(c *gin.Context) bool {
	userID, _, role := h.getCurrentUser(c)
	switch role {
	case "admin", "mgmt", "specialist":
		return true
	}
	return h.engine.IsMgmtSpecialist(userID, "")
}

// vulnAssignOverride 漏洞情报自动分配时对某张处置单的负责人覆盖（前端预览可编辑后回传）
type vulnAssignOverride struct {
	ProjectID        int64  `json:"projectId"`
	ProjectType      string `json:"projectType"`
	AssigneeUsername string `json:"assigneeUsername"`
}

// createVulnIntelInTx 在事务内创建漏洞情报缺陷（外部实体表，软件/组件名支撑裂变）并生成处置工单；
// 漏洞编号已存在时返回 errVulnNoExists（幂等：不重复创建、不重复裂变）
// 处置分配两通道（与缺陷评审一致）：manual-手动指定（多负责人每人一张单）/ auto-自动裂变（默认，按组件类型拆单+匹配，可用 assignOverride 覆盖各单负责人）
func (h *RiskHandler) createVulnIntelInTx(tx *gorm.DB, org, name, vulnNo, title, severity, analyst, dept string, deadline *time.Time, remark, userName, assignMode string, assignees []string, assignOverride []vulnAssignOverride) error {
	var cnt int64
	tx.Model(&models.Defect{}).Where("defect_id = ?", vulnNo).Count(&cnt)
	if cnt > 0 {
		return errVulnNoExists
	}
	now := time.Now()
	entity := models.Defect{
		DefectID:       vulnNo,
		Software:       name,
		Org:            org,
		Title:          title,
		Body:           remark,
		Labels:         "[]",
		Reactions:      "[]",
		WatchMatches:   "[]",
		Severity:       severity,
		Source:         "vuln-intel",
		InternalStatus: disposal.StatusAssigned, // 漏洞情报：直接进入已分配（裂变）
		MaintDept:      dept,
		ImportedAt:     now,
	}
	if err := tx.Create(&entity).Error; err != nil {
		return err
	}
	flow := models.DefectFlow{
		DefectNo:       vulnNo,
		FlowType:       "vuln",
		MaintDept:      dept,
		DispatchTarget: analyst, // 研判人归属
		DispatchType:   "vuln-intel",
		DispatchTime:   &now,
		Dispatcher:     userName,
		Deadline:       deadline,
		CreateTime:     now,
		UpdateTime:     now,
	}
	if err := tx.Create(&flow).Error; err != nil {
		return err
	}
	// 处置分配两通道：manual-手动指定（多负责人每人一张单）/ auto-自动裂变（默认，软件→project / 组件→component_snapshot→project_inner）
	if assignMode == "manual" {
		return h.engine.CreateAssignedDisposalsTx(tx, vulnNo, assignees, userName)
	}
	if err := h.engine.DispatchVulnTx(tx, vulnNo, userName); err != nil {
		return err
	}
	// 自动裂变后按前端编辑的负责人覆盖各处置单（以创建表单人编辑为准）
	return h.applyAssignOverride(tx, vulnNo, assignOverride, userName)
}

// applyAssignOverride 按 projectId+projectType 覆盖裂变出的处置单负责人（漏洞情报 auto 预览编辑回传）
func (h *RiskHandler) applyAssignOverride(tx *gorm.DB, defectNo string, overrides []vulnAssignOverride, operatorName string) error {
	now := time.Now()
	for _, ov := range overrides {
		username := strings.TrimSpace(ov.AssigneeUsername)
		if username == "" {
			continue
		}
		var u models.User
		if err := tx.Where("username = ?", username).First(&u).Error; err != nil {
			return fmt.Errorf("负责人 %s 不是系统用户", username)
		}
		res := tx.Model(&models.DefectDisposal{}).
			Where("defect_no = ? AND project_id = ? AND project_type = ?", defectNo, ov.ProjectID, ov.ProjectType).
			Updates(map[string]interface{}{
				"assignee_username": u.Username,
				"assignee_name":     u.Nickname,
				"assigned_by":       operatorName,
				"assigned_at":       now,
				"update_time":       now,
			})
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

// MatchVulnIntel 漏洞情报自动分配预览：按软件/组件名称返回待裂变处置单清单与候选负责人（不创建、不落库），供前端展示并可编辑
func (h *RiskHandler) MatchVulnIntel(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		response.Error(c, 400, "请填写软件/组件名称")
		return
	}
	preview, err := h.engine.PreviewVulnFission(req.Name)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, preview)
}

// CreateVulnIntel 手动录入漏洞情报：输入软件/组件名称 + 漏洞编号等，创建后生成处置工单（自动裂变或手动指定处置人）
func (h *RiskHandler) CreateVulnIntel(c *gin.Context) {
	if !h.canManageVulnIntel(c) {
		response.Error(c, 403, "无权限：仅开源治理专员/管理员可操作")
		return
	}
	var req struct {
		Org      string `json:"org"`
		Name     string `json:"name"`      // 软件/组件名称（用于裂变匹配）
		VulnNo   string `json:"vulnNo"`    // 漏洞编号（唯一键）
		Title    string `json:"title"`     // 名称（可选，缺省用名称）
		Severity string `json:"severity"`  // 等级（选填）
		Analyst  string `json:"analyst"`   // 研判人（选填，缺省当前用户）
		Dept     string `json:"dept"`      // 研判部门（选填，缺省按软件从 project 解析）
		Deadline string `json:"deadline"`  // 处置截止（选填）
		Remark   string `json:"remark"`    // 详情/备注
		// 处置分配两通道：manual-手动指定（assignees 可多个，每人一张处置单）/ auto-自动裂变（默认）
		AssignMode string   `json:"assignMode"`
		Assignees  []string `json:"assignees"`
		// auto 模式下前端预览编辑后的负责人覆盖（按 projectId+projectType 定位处置单）
		AssignOverride []vulnAssignOverride `json:"assignOverride"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.VulnNo = strings.TrimSpace(req.VulnNo)
	if req.Name == "" || req.VulnNo == "" {
		response.Error(c, 400, "软件/组件名称与漏洞编号必填")
		return
	}
	userID, userName, _ := h.getCurrentUser(c)

	analyst := strings.TrimSpace(req.Analyst)
	if analyst == "" {
		analyst = userID
	}
	dept := strings.TrimSpace(req.Dept)
	if dept == "" {
		var prj models.Project
		if err := h.db.Where("LOWER(_name) = LOWER(?)", req.Name).Order("id ASC").First(&prj).Error; err == nil {
			dept = prj.MaintDept
		}
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = req.Name
	}
	severity := strings.TrimSpace(req.Severity)
	if severity == "" {
		severity = "未知"
	}
	var deadline *time.Time
	if t := strings.TrimSpace(req.Deadline); t != "" {
		if pt, err := time.ParseInLocation("2006-01-02 15:04:05", t, time.Local); err == nil {
			deadline = &pt
		}
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		return h.createVulnIntelInTx(tx, req.Org, req.Name, req.VulnNo, title, severity, analyst, dept, deadline, req.Remark, userName, req.AssignMode, req.Assignees, req.AssignOverride)
	})
	if err != nil {
		if errors.Is(err, errVulnNoExists) {
			response.Error(c, 400, "漏洞编号已存在，请勿重复录入")
			return
		}
		response.Error(c, 400, "创建/裂变失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"defectNo": req.VulnNo})
}

// ImportVulnIntel Excel 导入漏洞情报：必填列（组织/名称/漏洞编号/研判人/研判部门），逐行创建并直接裂变；已存在编号跳过不重复裂变
func (h *RiskHandler) ImportVulnIntel(c *gin.Context) {
	if !h.canManageVulnIntel(c) {
		response.Error(c, 403, "无权限：仅开源治理专员/管理员可操作")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 400, "请上传 Excel 文件")
		return
	}
	src, err := file.Open()
	if err != nil {
		response.Error(c, 400, "无法读取上传文件")
		return
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		response.Error(c, 400, "Excel 解析失败: "+err.Error())
		return
	}
	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil || len(rows) < 2 {
		response.Error(c, 400, "Excel 为空或缺少数据行")
		return
	}

	// 表头映射（允许列顺序不同）
	colIdx := map[string]int{}
	for i, hc := range rows[0] {
		colIdx[strings.TrimSpace(hc)] = i
	}
	required := []string{"组织", "名称", "漏洞编号", "研判人", "研判部门"}
	for _, r := range required {
		if _, ok := colIdx[r]; !ok {
			response.Error(c, 400, "Excel 缺少必填列: "+r)
			return
		}
	}
	get := func(rr []string, key string) string {
		if i, ok := colIdx[key]; ok && i < len(rr) {
			return strings.TrimSpace(rr[i])
		}
		return ""
	}

	userID, userName, _ := h.getCurrentUser(c)
	imported, skipped := 0, 0
	var errs []string
	for ri := 1; ri < len(rows); ri++ {
		rr := rows[ri]
		org := get(rr, "组织")
		name := get(rr, "名称")
		vulnNo := get(rr, "漏洞编号")
		analyst := get(rr, "研判人")
		dept := get(rr, "研判部门")
		if name == "" || vulnNo == "" {
			skipped++
			errs = append(errs, fmt.Sprintf("第%d行：缺少名称/漏洞编号", ri+1))
			continue
		}
		// 研判人：支持账号或姓名（昵称），映射到 users.username
		analystUsername := ""
		if analyst != "" {
			var u models.User
			if err := h.db.Where("username = ? OR nickname = ?", analyst, analyst).First(&u).Error; err == nil {
				analystUsername = u.Username
			} else {
				analystUsername = analyst
			}
		}
		if analystUsername == "" {
			analystUsername = userID
		}
		// 研判部门：优先 Excel，缺省按软件从 project 解析
		if dept == "" {
			var prj models.Project
			if err := h.db.Where("LOWER(_name) = LOWER(?)", name).Order("id ASC").First(&prj).Error; err == nil {
				dept = prj.MaintDept
			}
		}
		err := h.db.Transaction(func(tx *gorm.DB) error {
			// Excel 批量导入：固定自动裂变（auto），不提供手动指定/预览编辑
			return h.createVulnIntelInTx(tx, org, name, vulnNo, name, "未知", analystUsername, dept, nil, "", userName, "", nil, nil)
		})
		if err != nil {
			if errors.Is(err, errVulnNoExists) {
				continue // 已存在：跳过，不重复裂变
			}
			errs = append(errs, fmt.Sprintf("第%d行 %s: %v", ri+1, vulnNo, err))
			continue
		}
		imported++
	}
	response.Success(c, gin.H{"imported": imported, "skipped": skipped, "errors": errs})
}

// GetDefectAnalystCandidates 获取缺陷自动匹配的研判候选负责人（第一个为主岗）
// 基础软件取该软件 project_owner；组件取依赖项目（component_snapshot → project_inner）的 project_owner
func (h *RiskHandler) GetDefectAnalystCandidates(c *gin.Context) {
	defectNo := c.Param("id")
	if defectNo == "" {
		response.Error(c, 400, "无效的缺陷编号")
		return
	}
	response.Success(c, h.engine.MatchCandidatesByNo(defectNo))
}

// AssignDefectAnalyst 第一次分配：把待研判缺陷分配给软件负责人研判
func (h *RiskHandler) AssignDefectAnalyst(c *gin.Context) {
	defectNo := c.Param("id")
	if defectNo == "" {
		response.Error(c, 400, "无效的缺陷编号")
		return
	}

	var req struct {
		AnalystUsername string `json:"analystUsername"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID, userName, role := h.getCurrentUser(c)
	if err := h.engine.AssignAnalyst(defectNo, req.AnalystUsername, userID, userName, role); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, nil)
}

// AutoAssignDefectAnalyst 自动分配研判人：系统取匹配候选的第一个（主岗）
func (h *RiskHandler) AutoAssignDefectAnalyst(c *gin.Context) {
	defectNo := c.Param("id")
	if defectNo == "" {
		response.Error(c, 400, "无效的缺陷编号")
		return
	}

	userID, userName, role := h.getCurrentUser(c)
	analystName, err := h.engine.AutoAssignAnalyst(defectNo, userID, userName, role)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, gin.H{"analystName": analystName, "message": "已自动分配研判人 " + analystName + "（主岗）"})
}

// AutoAssignDefectDisposals 自动分配处置单：对该缺陷全部待分配处置单匹配各部门项目主岗
func (h *RiskHandler) AutoAssignDefectDisposals(c *gin.Context) {
	defectNo := c.Param("id")
	if defectNo == "" {
		response.Error(c, 400, "无效的缺陷编号")
		return
	}

	userID, userName, role := h.getCurrentUser(c)
	assigned, missed, err := h.engine.AutoAssignDisposals(defectNo, userID, userName, role)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, gin.H{"assigned": assigned, "missed": missed})
}

// AnalyzeDefect 研判：pending_analysis -> pending_review / assigned / ignored
func (h *RiskHandler) AnalyzeDefect(c *gin.Context) {
	defectNo := c.Param("id")
	if defectNo == "" {
		response.Error(c, 400, "无效的缺陷编号")
		return
	}

	var req struct {
		Conclusion    string   `json:"conclusion"`
		AnalyzeResult string   `json:"analyzeResult"`
		// 漏洞直接分配（need_fix）处置人：manual-手动指定（assignees 可多个）/ auto-自动分配（默认，裂变+自动匹配）
		AssignMode string   `json:"assignMode"`
		Assignees  []string `json:"assignees"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID, userName, role := h.getCurrentUser(c)
	if err := h.engine.Analyze(defectNo, req.Conclusion, req.AnalyzeResult, req.AssignMode, req.Assignees, userID, userName, role); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, nil)
}

// ReviewDefect 评审（上会）：pending_review -> assigned / ignored
func (h *RiskHandler) ReviewDefect(c *gin.Context) {
	defectNo := c.Param("id")
	if defectNo == "" {
		response.Error(c, 400, "无效的缺陷编号")
		return
	}

	var req struct {
		ReviewConclusion string   `json:"reviewConclusion"`
		ReviewMeetingNo  string   `json:"reviewMeetingNo"`
		Remark           string   `json:"remark"`
		Deadline         string   `json:"deadline"`
		// 处置分配两通道：manual-手动指定（assignees 可多个，每人一张处置单）/ auto-自动分配（默认，裂变+自动匹配）
		AssignMode string   `json:"assignMode"`
		Assignees  []string `json:"assignees"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	var deadline *time.Time
	if req.Deadline != "" {
		if t, parseErr := time.ParseInLocation("2006-01-02 15:04:05", req.Deadline, time.Local); parseErr == nil {
			deadline = &t
		}
	}

	userID, userName, role := h.getCurrentUser(c)
	if err := h.engine.Review(defectNo, req.ReviewConclusion, req.ReviewMeetingNo, req.Remark, deadline, req.AssignMode, req.Assignees, userID, userName, role); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListDefectDisposals 查询某缺陷裂变出的处置工单
func (h *RiskHandler) ListDefectDisposals(c *gin.Context) {
	defectNo := c.Param("id")
	if defectNo == "" {
		response.Error(c, 400, "无效的缺陷编号")
		return
	}

	userID, userName, role := h.getCurrentUser(c)
	userDept := decodeHeaderValue(c.GetHeader("X-User-Dept"))
	// 可见性：能看见该缺陷的用户即可查看其全部处置单（运维部门专员/使用部门专员/管理部门专员均看全景，只读展示），
	// 改派仍受 Assign 权限限制（仅本处置单所属部门专员可改派，管理部门专员不可跨部门改派）
	if !h.engine.CanSeeDefectByNo(userID, userName, role, userDept, defectNo) {
		response.Success(c, []models.DefectDisposal{})
		return
	}

	var disposals []models.DefectDisposal
	if err := h.db.Where("defect_no = ?", defectNo).Order("id ASC").Find(&disposals).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}
	response.Success(c, disposals)
}

// GetDefectHistory 查询缺陷状态变更历史
func (h *RiskHandler) GetDefectHistory(c *gin.Context) {
	defectNo := c.Param("id")
	if defectNo == "" {
		response.Error(c, 400, "无效的缺陷编号")
		return
	}

	var history []models.DefectStatusHistory
	if err := h.db.Where("defect_no = ?", defectNo).Order("id DESC").Find(&history).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}
	response.Success(c, history)
}

// ListDisposals 查询处置工单（待我处置/部门工单），返回父缺陷编号/标题便于处理人识别
func (h *RiskHandler) ListDisposals(c *gin.Context) {
	page, pageSize := parsePagination(c)
	mine := c.Query("mine")
	statusFilter := c.Query("status")

	userID, userName, role := h.getCurrentUser(c)
	isAdmin := role == "admin"
	isMgmt := h.engine.IsMgmtSpecialist(userID, userName)

	buildQuery := func() *gorm.DB {
		q := h.db.Table("defect_disposal dd")
		if mine == "1" {
			// 只看分配给我的工单（按处理人姓名或账号）
			q = q.Where("(dd.assignee_name = ? OR dd.assignee_username = ?)", userName, userID)
		} else if !isAdmin && !isMgmt {
			// 部门专员：只看本部门负责部门的工单
			depts := h.engine.SpecialistDepts(userID, userName)
			if len(depts) > 0 {
				q = q.Where("dd.dept_name IN ?", depts)
			} else {
				q = q.Where("1 = 0")
			}
		}
		if statusFilter != "" {
			q = q.Where("dd.status = ?", statusFilter)
		}
		return q
	}

	var total int64
	if err := buildQuery().Count(&total).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	type disposalRow struct {
		models.DefectDisposal
		DefectTitle string `gorm:"column:defect_title" json:"defectTitle"`
		// 外部缺陷内容详情（处置工单详情弹框展示：中文/英文信息）
		DefectBody   string `gorm:"column:defect_body" json:"defectBody"`
		TitleZh      string `gorm:"column:title_zh" json:"titleZh"`
		SummaryZh    string `gorm:"column:summary_zh" json:"summaryZh"`
		RepairAdvice string `gorm:"column:repair_advice" json:"repairAdvice"`
	}

	var list []disposalRow
	offset := (page - 1) * pageSize
	err := buildQuery().
		Select("dd.*, COALESCE(dl.title, d.title) AS defect_title, d.body AS defect_body, d.title_zh AS title_zh, d.summary_zh AS summary_zh, d.repair_advice AS repair_advice").
		Joins("LEFT JOIN defect_ledger dl ON dl.defect_no = dd.defect_no").
		// 漏洞情报/外部缺陷实体在 defect 表（defect_ledger 无记录），补 join 取标题与内容详情，保证处置工单能展示"处置什么"
		Joins("LEFT JOIN defect d ON d.defect_id = dd.defect_no").
		Order("dd.id DESC").
		Offset(offset).Limit(pageSize).
		Scan(&list).Error
	if err != nil {
		response.Error(c, 500, "查询失败")
		return
	}
	response.PageSuccess(c, list, total, page, pageSize)
}

// AssignDisposal 补分配/改派处理人（工单保持「待处置」）
func (h *RiskHandler) AssignDisposal(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的处置单ID")
		return
	}

	var req struct {
		AssigneeUsername string `json:"assigneeUsername"`
		AssigneeName     string `json:"assigneeName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID, userName, _ := h.getCurrentUser(c)
	if err := h.engine.Assign(id, req.AssigneeUsername, req.AssigneeName, userID, userName); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, nil)
}

// CompleteDisposal 处置完成：待处置/处置中 -> completed
func (h *RiskHandler) CompleteDisposal(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的处置单ID")
		return
	}

	var req struct {
		TicketNo string `json:"ticketNo"`
		Note     string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID, userName, role := h.getCurrentUser(c)
	if err := h.engine.Complete(id, req.TicketNo, req.Note, userID, userName, role); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, nil)
}

// MarkDisposalNA 标记忽略
func (h *RiskHandler) MarkDisposalNA(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的处置单ID")
		return
	}

	var req struct {
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID, userName, role := h.getCurrentUser(c)
	if err := h.engine.MarkNA(id, req.Remark, userID, userName, role); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, nil)
}

// ==================== CICD检查 ====================

// CICDCheckItem CICD检查台账项
type CICDCheckItem struct {
	ID              int64     `json:"id"`
	ProCode         string    `json:"pro_code"`
	CurrentID       string    `json:"current_ID"`
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
			cb.CURRENT_ID as current_id,
			pi.group_name,
			pi.dept_name,
			cc.applicable_level as level,
			CASE 
				WHEN cc.applicable_level IN ('C', 'D') THEN '阻断'
				ELSE '警告'
			END as disposal_type
		`).
		Joins("LEFT JOIN project_inner pi ON cb.PROJECT_ID = pi.id").
		Joins("LEFT JOIN component_core cc ON cb.BLACK_ID = cc.id")

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

	// 在Go代码中去重：按pro_code+current_ID唯一标识
	type ProjectKey struct {
		ProCode       string
		CommitID      string
	}
	seen := make(map[ProjectKey]bool)
	var uniqueItems []CICDCheckItem
	for _, item := range items {
		key := ProjectKey{item.ProCode, item.CurrentID}
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
			pi.module_code as module_name,
			cb.create_time as check_time,
			CASE 
				WHEN cc.applicable_level IN ('C', 'D') THEN '阻断'
				ELSE '警告'
			END as disposal,
			cc._name as blacklist_name,
			cc.reason
		`).
		Joins("LEFT JOIN project_inner pi ON cb.PROJECT_ID = pi.id").
		Joins("LEFT JOIN component_core cc ON cb.BLACK_ID = cc.id").
		Where("pi.pro_code = ?", projectName)

	if commitID != "" {
		query = query.Where("cb.CURRENT_ID = ?", commitID)
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
			cs._name as component_name,
			cs._version as version,
			cc.reason,
			cb.create_time
		`).
		Joins("LEFT JOIN component_core cc ON cb.BLACK_ID = cc.id").
		Joins("LEFT JOIN component_snapshot cs ON cb.COMPONENT_ID = cs.id").
		Where("cb.PROJECT_ID = ?", projectID)

	var records []map[string]interface{}
	if err := query.Order("cb.create_time DESC").Scan(&records).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, records)
}
