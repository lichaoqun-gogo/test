package handlers

import (
	"fmt"
	"log"
	"net/url"
	"osm-mngt/internal/models"
	"osm-mngt/internal/service/workflow"
	"osm-mngt/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db           *gorm.DB
	workflowEngine *workflow.WorkflowEngine
}

func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{db: db}
}

func (h *TaskHandler) GetWorkflowEngine() *workflow.WorkflowEngine {
	if h.workflowEngine == nil {
		h.workflowEngine = workflow.NewWorkflowEngine(h.db)
	}
	return h.workflowEngine
}

func decodeHeaderValue(value string) string {
	if value == "" {
		return ""
	}
	decoded, err := url.QueryUnescape(value)
	if err != nil {
		return value
	}
	return decoded
}

func (h *TaskHandler) getUserInfo(c *gin.Context) (string, string, string) {
	userID := c.GetHeader("X-User-ID")
	userName := decodeHeaderValue(c.GetHeader("X-User-Name"))
	userRole := c.GetHeader("X-User-Role")
	return userID, userName, userRole
}

// ListTasks 查询任务列表
// 注意：后端不做权限过滤，所有任务均返回。菜单/数据权限后续由统一权限管理平台控制。
func (h *TaskHandler) ListTasks(c *gin.Context) {
	page, pageSize := parsePagination(c)

	keyword := c.Query("keyword")
	typeFilter := c.Query("type")
	statusFilter := c.Query("status")

	query := h.db.Model(&models.GovernanceTask{})

	if keyword != "" {
		query = query.Where("task_no LIKE ? OR title LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if typeFilter != "" {
		query = query.Where("type = ?", typeFilter)
	}
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	var tasks []models.GovernanceTask
	offset := (page - 1) * pageSize
	if err := query.Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error; err != nil {
		response.Error(c, 500, "查询失败")
		return
	}

	// 为每个任务关联查询task_detail
	userID, _, _ := h.getUserInfo(c)
	tasksWithDetail := make([]gin.H, 0, len(tasks))
	for _, task := range tasks {
		log.Printf("[ListTasks] 任务ID: %d, taskNo: %s, status: %s, applicantID: %s, assignedTo: %s",
			task.ID, task.TaskNo, task.Status, task.ApplicantID, task.AssignedTo)

		taskMap := gin.H{
			"id":           task.ID,
			"taskNo":       task.TaskNo,
			"type":         task.Type,
			"actionType":   task.ActionType,
			"title":        task.Title,
			"applicant":    task.Applicant,
			"applicantId":  task.ApplicantID,
			"applicantDept": task.ApplicantDept,
			"status":       task.Status,
			"priority":     task.Priority,
			"currentStep":  task.CurrentStep,
			"assignedTo":   task.AssignedTo,
			"deadline":     task.Deadline,
			"description":  task.Description,
			"createBy":     task.CreateBy,
			"createTime":   task.CreateTime,
			"updateBy":     task.UpdateBy,
			"updateTime":   task.UpdateTime,
		}

		// 审批中的任务：标记当前用户是否可处理当前节点
		//（一个部门多名专员场景：任一同部门同角色在岗专员均可审批，后端统一下发标记）
		if task.Status == "approving" && userID != "" {
			var node models.ApprovalNode
			if err := h.db.Where("task_id = ? AND status = ? AND is_skipped = ? AND approver_id != '' AND node_order > 0",
				task.ID, "pending", false).
				Order("node_order ASC").First(&node).Error; err == nil {
				taskMap["canApprove"] = workflow.CanUserHandleNode(h.db, &node, userID)
			}
		}

		// 查询task_detail
		var taskDetail interface{}
		switch task.Type {
		case "introduce":
			var introduce models.TaskIntroduce
			if err := h.db.Where("task_id = ?", task.ID).First(&introduce).Error; err == nil {
				taskDetail = introduce
			}
		case "upgrade":
			var upgrade models.TaskUpgrade
			if err := h.db.Where("task_id = ?", task.ID).First(&upgrade).Error; err == nil {
				taskDetail = upgrade
			}
		case "withdraw":
			var withdraw models.TaskWithdraw
			if err := h.db.Where("task_id = ?", task.ID).First(&withdraw).Error; err == nil {
				taskDetail = withdraw
			}
		case "defect":
			var defect models.TaskDefect
			if err := h.db.Where("task_id = ?", task.ID).First(&defect).Error; err == nil {
				taskDetail = defect
			}
		case "usage":
			var usage models.TaskUsage
			if err := h.db.Where("task_id = ?", task.ID).First(&usage).Error; err == nil {
				taskDetail = usage
			}
		}

		if taskDetail != nil {
			taskMap["task_detail"] = taskDetail
		}

		tasksWithDetail = append(tasksWithDetail, taskMap)
	}

	response.PageSuccess(c, tasksWithDetail, total, page, pageSize)
}

// GetTask 获取任务详情
func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的任务ID")
		return
	}

	var task models.GovernanceTask
	if err := h.db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "任务不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, task)
}

// GetTaskDetail 获取任务完整详情（包含扩展表和审批记录）
func (h *TaskHandler) GetTaskDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的任务ID")
		return
	}

	var task models.GovernanceTask
	if err := h.db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "任务不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	// 查询扩展表数据
	var taskDetail interface{}
	switch task.Type {
	case "introduce":
		var introduce models.TaskIntroduce
		if err := h.db.Where("task_id = ?", id).First(&introduce).Error; err != nil && err != gorm.ErrRecordNotFound {
			response.Error(c, 500, "查询扩展数据失败")
			return
		}
		taskDetail = introduce
	case "upgrade":
		var upgrade models.TaskUpgrade
		if err := h.db.Where("task_id = ?", id).First(&upgrade).Error; err != nil && err != gorm.ErrRecordNotFound {
			response.Error(c, 500, "查询扩展数据失败")
			return
		}
		taskDetail = upgrade
	case "withdraw":
		var withdraw models.TaskWithdraw
		if err := h.db.Where("task_id = ?", id).First(&withdraw).Error; err != nil && err != gorm.ErrRecordNotFound {
			response.Error(c, 500, "查询扩展数据失败")
			return
		}
		taskDetail = withdraw
	case "defect":
		var defect models.TaskDefect
		if err := h.db.Where("task_id = ?", id).First(&defect).Error; err != nil && err != gorm.ErrRecordNotFound {
			response.Error(c, 500, "查询扩展数据失败")
			return
		}
		taskDetail = defect
	case "usage":
		var usage models.TaskUsage
		if err := h.db.Where("task_id = ?", id).First(&usage).Error; err != nil && err != gorm.ErrRecordNotFound {
			response.Error(c, 500, "查询扩展数据失败")
			return
		}
		taskDetail = usage
	}

	// 查询审批节点
	var approvalNodes []models.ApprovalNode
	if err := h.db.Where("task_id = ?", id).Order("node_order ASC").Find(&approvalNodes).Error; err != nil {
		response.Error(c, 500, "查询审批节点失败")
		return
	}

	// 查询审批历史（使用 workflow engine 获取最后一次提交开始的历史）
	approvalHistory, err := h.GetWorkflowEngine().GetApprovalHistory(id)
	if err != nil {
		response.Error(c, 500, "查询审批历史失败")
		return
	}

	// 查询适用性评估数据（仅introduce类型）
	var assessment *models.TaskApplicabilityAssessment
	if task.Type == "introduce" {
		var ass models.TaskApplicabilityAssessment
		if err := h.db.Where("task_id = ?", id).First(&ass).Error; err != nil && err != gorm.ErrRecordNotFound {
			response.Error(c, 500, "查询评估数据失败")
			return
		}
		if err == nil {
			assessment = &ass
		}
	}

	response.Success(c, gin.H{
		"task":                       task,
		"task_detail":                taskDetail,
		"applicability_assessment":   assessment,
		"approval_nodes":             approvalNodes,
		"approval_history":           approvalHistory,
	})
}

// CreateTask 创建任务
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req struct {
		Type       string                 `json:"type" binding:"required"`
		Title      string                 `json:"title" binding:"required"`
		Applicant  string                 `json:"applicant"`
		Deadline   string                 `json:"deadline"`
		Status     string                 `json:"status"` // draft 或 approving
		Detail     map[string]interface{} `json:"task_detail"`
		Assessment map[string]interface{} `json:"applicability_assessment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	// 从task_detail中提取project_dept（使用部门）
	var projectDept string
	if req.Detail != nil {
		if pd, ok := req.Detail["project_dept"].(string); ok {
			projectDept = pd
		}
	}

	taskNo := generateTaskNo()
	userID, userName, _ := h.getUserInfo(c)
	userDept := decodeHeaderValue(c.GetHeader("X-User-Dept")) // 从请求头获取用户部门并解码

	var task models.GovernanceTask

	// 使用事务
	err := h.db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建任务主表
		now := time.Now()
		
		// 申请人信息优先使用请求头中的真实用户信息
		applicant := userName
		if applicant == "" {
			applicant = req.Applicant // 兜底
		}
		applicantDept := userDept
		if applicantDept == "" {
			applicantDept = "技术部" // 兜底
		}
		
		task = models.GovernanceTask{
			TaskNo:       taskNo,
			Type:         req.Type,
			ActionType:   req.Type, // 默认为任务类型
			Title:        req.Title,
			Applicant:    applicant,
			ApplicantID:  userID,
			ApplicantDept: applicantDept,
			ProjectDept:  projectDept,  // 使用部门（用于审批流解析）
			Status:       "draft",
			CreateTime:   now,
			CreateBy:     userID,
			UpdateTime:   now,
			UpdateBy:     userName,
		}

		if req.Deadline != "" {
			deadline, err := time.Parse("2006-01-02", req.Deadline)
			if err == nil {
				task.Deadline = &deadline
			}
		}

		if err := tx.Create(&task).Error; err != nil {
			return err
		}

		// 2. 创建扩展表记录
		if err := createTaskDetail(tx, task.ID, req.Type, req.Detail); err != nil {
			return err
		}

		// 3. 创建适用性评估记录（仅introduce类型）
		if req.Type == "introduce" && req.Assessment != nil {
			if err := createApplicabilityAssessment(tx, task.ID, req.Assessment); err != nil {
				return err
			}
		}

		// 4. 如果是提交状态，创建审批节点
		taskStatus := req.Status
		if taskStatus == "" {
			taskStatus = "draft"
		}
		
		if taskStatus == "approving" {
			engine := h.GetWorkflowEngine()
			// 非引入申请：从project表查运维部门（升级/退出 node_1 复审以运维部门为准）
			maintDept := ""
			if req.Type != "introduce" {
				if cn, ok := req.Detail["component_name"].(string); ok && cn != "" {
					tx.Table("project").Select("maint_dept").Where("_name = ?", cn).Row().Scan(&maintDept)
				}
			}
			if err := engine.ResolveNodesAndCreate(tx, &task, req.Type, maintDept); err != nil {
				return fmt.Errorf("创建审批流失败: %w", err)
			}
			task.Status = "approving"
			tx.Save(&task)
		}

		return nil
	})

	if err != nil {
		response.Error(c, 500, "创建失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"task_no": taskNo, "id": task.ID})
}

// UpdateTask 更新任务
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的任务ID")
		return
	}

	var task models.GovernanceTask
	if err := h.db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "任务不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	// 只允许更新 draft / withdrawn / rejected（驳回后可编辑重提）状态的任务
	if task.Status != "draft" && task.Status != "withdrawn" && task.Status != "rejected" {
		response.Error(c, 400, "只能更新草稿、已撤回或已驳回状态的任务")
		return
	}

	var req struct {
		Title      string                 `json:"title"`
		Applicant  string                 `json:"applicant"`
		Deadline   string                 `json:"deadline"`
		Detail     map[string]interface{} `json:"task_detail"`
		Assessment map[string]interface{} `json:"applicability_assessment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	_, userName, _ := h.getUserInfo(c)
	now := time.Now()

	// 更新任务主表
	updates := map[string]interface{}{
		"update_by":   userName,
		"update_time": now,
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Applicant != "" {
		updates["applicant"] = req.Applicant
	}
	if req.Deadline != "" {
		if deadline, err := time.Parse("2006-01-02", req.Deadline); err == nil {
			updates["deadline"] = deadline
		}
	}

	if err := h.db.Model(&models.GovernanceTask{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Printf("更新任务主表失败: %v", err)
		response.Error(c, 500, "更新失败")
		return
	}

	// 更新扩展表数据
	if req.Detail != nil {
		if err := updateTaskDetail(h.db, id, task.Type, req.Detail); err != nil {
			log.Printf("更新扩展数据失败: %v", err)
			response.Error(c, 500, "更新扩展数据失败: "+err.Error())
			return
		}
	}

	// 更新适用性评估数据（仅introduce类型）
	if task.Type == "introduce" && req.Assessment != nil {
		if err := updateApplicabilityAssessment(h.db, id, req.Assessment); err != nil {
			log.Printf("更新评估数据失败: %v", err)
			response.Error(c, 500, "更新评估数据失败: "+err.Error())
			return
		}
	}

	response.Success(c, nil)
}

// UpdateTaskStatus 更新任务状态
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的任务ID")
		return
	}

	var task models.GovernanceTask
	if err := h.db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "任务不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	var req struct {
		Status  string `json:"status" binding:"required"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	// 状态流转验证
	if !isValidStatusTransition(task.Status, req.Status) {
		response.Error(c, 400, fmt.Sprintf("状态流转不合法: %s -> %s", task.Status, req.Status))
		return
	}

	_, userName, _ := h.getUserInfo(c)

	if err := h.db.Model(&models.GovernanceTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     req.Status,
		"update_by":  userName,
		"update_time": time.Now(),
	}).Error; err != nil {
		response.Error(c, 500, "更新失败")
		return
	}

	response.Success(c, nil)
}

// WithdrawTask 撤回任务
func (h *TaskHandler) WithdrawTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的任务ID")
		return
	}

	var task models.GovernanceTask
	if err := h.db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "任务不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	if task.Status != "pending" && task.Status != "processing" {
		response.Error(c, 400, "当前状态不允许撤回")
		return
	}

	_, userName, _ := h.getUserInfo(c)

	if err := h.db.Model(&models.GovernanceTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "draft",
		"update_by":  userName,
		"update_time": time.Now(),
	}).Error; err != nil {
		response.Error(c, 500, "撤回失败")
		return
	}

	response.Success(c, nil)
}

// ==================== 辅助函数 ====================

func generateTaskNo() string {
	return fmt.Sprintf("TASK-%d", time.Now().Unix())
}

func isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		"draft":      {"pending"},
		"pending":    {"processing", "completed", "rejected", "draft"},
		"processing": {"completed", "rejected", "draft"},
		"completed":  {"draft"},
		"rejected":   {"draft"},
	}

	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

func createTaskDetail(tx *gorm.DB, taskID int64, taskType string, detail map[string]interface{}) error {
	now := time.Now()

	switch taskType {
	case "introduce":
		introduce := models.TaskIntroduce{
			TaskID:     taskID,
			CreateTime: now,
			UpdateTime: now,
		}
		if v, ok := detail["component_name"].(string); ok {
			introduce.ComponentName = v
		}
		if v, ok := detail["component_version"].(string); ok {
			introduce.ComponentVersion = v
		}
		if v, ok := detail["license_type"].(string); ok {
			introduce.LicenseType = v
		}
		if v, ok := detail["introduce_reason"].(string); ok {
			introduce.IntroduceReason = v
		}
		if v, ok := detail["use_scenario"].(string); ok {
			introduce.UseScenario = v
		}
		if v, ok := detail["project_names"].(string); ok {
			introduce.ProjectNames = v
		}
		if v, ok := detail["remark"].(string); ok {
			introduce.Remark = v
		}
		if v, ok := detail["pre_eval_conclusion"].(string); ok {
			introduce.PreEvalConclusion = v
		}
		return tx.Create(&introduce).Error

	case "upgrade":
		upgrade := models.TaskUpgrade{
			TaskID:     taskID,
			CreateTime: now,
			UpdateTime: now,
		}
		if v, ok := detail["component_name"].(string); ok {
			upgrade.ComponentName = v
		}
		if v, ok := detail["current_version"].(string); ok {
			upgrade.CurrentVersion = v
		}
		if v, ok := detail["target_version"].(string); ok {
			upgrade.TargetVersion = v
		}
		if v, ok := detail["upgrade_reason"].(string); ok {
			upgrade.UpgradeReason = v
		}
		if v, ok := detail["test_report_link"].(string); ok {
			upgrade.TestReportLink = v
		}
		if v, ok := detail["impact_analysis"].(string); ok {
			upgrade.ImpactAnalysis = v
		}
		if v, ok := detail["compatibility_risk"].(string); ok {
			upgrade.CompatibilityRisk = v
		}
		if v, ok := detail["rollback_plan"].(string); ok {
			upgrade.RollbackPlan = v
		}
		return tx.Create(&upgrade).Error

	case "withdraw":
		withdraw := models.TaskWithdraw{
			TaskID:     taskID,
			CreateTime: now,
			UpdateTime: now,
		}
		if v, ok := detail["component_name"].(string); ok {
			withdraw.ComponentName = v
		}
		if v, ok := detail["component_version"].(string); ok {
			withdraw.ComponentVersion = v
		}
		if v, ok := detail["withdraw_reason"].(string); ok {
			withdraw.WithdrawReason = v
		}
		if v, ok := detail["exit_reason"].(string); ok && withdraw.WithdrawReason == "" {
			withdraw.WithdrawReason = v
		}
		if v, ok := detail["replacement_component"].(string); ok {
			withdraw.ReplacementComponent = v
		}
		if v, ok := detail["alternative_plan"].(string); ok {
			withdraw.AlternativePlan = v
		}
		if v, ok := detail["code_cleanup_plan"].(string); ok {
			withdraw.CodeCleanupPlan = v
		}
		if v, ok := detail["affected_projects"].(string); ok {
			withdraw.AffectedProjects = v
		}
		return tx.Create(&withdraw).Error

	case "defect":
		defect := models.TaskDefect{
			TaskID:     taskID,
			CreateTime: now,
			UpdateTime: now,
			// 内部缺陷编号：系统自动生成（DF-纳秒时间戳，与缺陷台账编号风格一致），不允许用户手填
			InternalDefectNo: fmt.Sprintf("DF-%d", time.Now().UnixNano()),
		}
		if v, ok := detail["risk_description"].(string); ok {
			defect.RiskDescription = v
		}
		return tx.Create(&defect).Error

	case "usage":
		usage := models.TaskUsage{
			TaskID:     taskID,
			CreateTime: now,
			UpdateTime: now,
		}
		if v, ok := detail["component_name"].(string); ok {
			usage.ComponentName = v
		}
		if v, ok := detail["version"].(string); ok {
			usage.Version = v
		}
		if v, ok := detail["is_production"].(float64); ok {
			usage.IsProduction = int(v)
		}
		if v, ok := detail["purpose"].(string); ok {
			usage.Purpose = v
		}
		return tx.Create(&usage).Error

	default:
		return nil
	}
}

func createApprovalNodes(tx *gorm.DB, taskID int64, taskType string, taskNo string) error {
	now := time.Now()

	// 根据任务类型定义不同的审批节点
	nodes := []models.ApprovalNode{
		{
			TaskID:     taskID,
			NodeName:   "开源治理专员审批",
			NodeOrder:  1,
			Status:     "pending",
			CreateTime: now,
		},
		{
			TaskID:     taskID,
			NodeName:   "管理员审批",
			NodeOrder:  2,
			Status:     "pending",
			CreateTime: now,
		},
	}

	for _, node := range nodes {
		if err := tx.Create(&node).Error; err != nil {
			return err
		}
	}

	return nil
}

// createApplicabilityAssessment 创建适用性评估记录
func createApplicabilityAssessment(tx *gorm.DB, taskID int64, assessment map[string]interface{}) error {
	now := time.Now()
	ass := models.TaskApplicabilityAssessment{
		TaskID:     taskID,
		CreateTime: now,
		UpdateTime: now,
	}

	dimensions := []string{
		"functionality", "security", "reliability", "performance",
		"compatibility", "maintainability", "extensibility", "usability",
	}

	for _, dim := range dimensions {
		if item, ok := assessment[dim].(map[string]interface{}); ok {
			if status, ok := item["status"].(string); ok {
				deviationExp := ""
				if exp, ok := item["deviationExplanation"].(string); ok {
					deviationExp = exp
				}

				switch dim {
				case "functionality":
					ass.FunctionalityStatus = status
					ass.FunctionalityDeviation = deviationExp
				case "security":
					ass.SecurityStatus = status
					ass.SecurityDeviation = deviationExp
				case "reliability":
					ass.ReliabilityStatus = status
					ass.ReliabilityDeviation = deviationExp
				case "performance":
					ass.PerformanceStatus = status
					ass.PerformanceDeviation = deviationExp
				case "compatibility":
					ass.CompatibilityStatus = status
					ass.CompatibilityDeviation = deviationExp
				case "maintainability":
					ass.MaintainabilityStatus = status
					ass.MaintainabilityDeviation = deviationExp
				case "extensibility":
					ass.ExtensibilityStatus = status
					ass.ExtensibilityDeviation = deviationExp
				case "usability":
					ass.UsabilityStatus = status
					ass.UsabilityDeviation = deviationExp
				}
			}
		}
	}

	return tx.Create(&ass).Error
}

// updateApplicabilityAssessment 更新适用性评估记录
func updateApplicabilityAssessment(db *gorm.DB, taskID int64, assessment map[string]interface{}) error {
	var ass models.TaskApplicabilityAssessment
	err := db.Where("task_id = ?", taskID).First(&ass).Error
	if err == gorm.ErrRecordNotFound {
		// 不存在则创建
		return createApplicabilityAssessment(db, taskID, assessment)
	}
	if err != nil {
		return err
	}

	// 更新各维度字段
	dimensions := []string{
		"functionality", "security", "reliability", "performance",
		"compatibility", "maintainability", "extensibility", "usability",
	}

	for _, dim := range dimensions {
		if item, ok := assessment[dim].(map[string]interface{}); ok {
			if status, ok := item["status"].(string); ok {
				deviationExp := ""
				if exp, ok := item["deviationExplanation"].(string); ok {
					deviationExp = exp
				}

				switch dim {
				case "functionality":
					ass.FunctionalityStatus = status
					ass.FunctionalityDeviation = deviationExp
				case "security":
					ass.SecurityStatus = status
					ass.SecurityDeviation = deviationExp
				case "reliability":
					ass.ReliabilityStatus = status
					ass.ReliabilityDeviation = deviationExp
				case "performance":
					ass.PerformanceStatus = status
					ass.PerformanceDeviation = deviationExp
				case "compatibility":
					ass.CompatibilityStatus = status
					ass.CompatibilityDeviation = deviationExp
				case "maintainability":
					ass.MaintainabilityStatus = status
					ass.MaintainabilityDeviation = deviationExp
				case "extensibility":
					ass.ExtensibilityStatus = status
					ass.ExtensibilityDeviation = deviationExp
				case "usability":
					ass.UsabilityStatus = status
					ass.UsabilityDeviation = deviationExp
				}
			}
		}
	}

	ass.UpdateTime = time.Now()
	return db.Save(&ass).Error
}

// updateTaskDetail 更新任务扩展表数据
func updateTaskDetail(db *gorm.DB, taskID int64, taskType string, detail map[string]interface{}) error {
	now := time.Now()

	switch taskType {
	case "introduce":
		var introduce models.TaskIntroduce
		err := db.Where("task_id = ?", taskID).First(&introduce).Error
		if err == gorm.ErrRecordNotFound {
			// 不存在则创建
			return createTaskDetail(db, taskID, taskType, detail)
		}
		if err != nil {
			return err
		}
		// 存在则更新
		if v, ok := detail["component_name"].(string); ok {
			introduce.ComponentName = v
		}
		if v, ok := detail["component_version"].(string); ok {
			introduce.ComponentVersion = v
		}
		if v, ok := detail["license_type"].(string); ok {
			introduce.LicenseType = v
		}
		if v, ok := detail["introduce_reason"].(string); ok {
			introduce.IntroduceReason = v
		}
		if v, ok := detail["use_scenario"].(string); ok {
			introduce.UseScenario = v
		}
		if v, ok := detail["project_names"].(string); ok {
			introduce.ProjectNames = v
		}
		if v, ok := detail["remark"].(string); ok {
			introduce.Remark = v
		}
		if v, ok := detail["pre_eval_conclusion"].(string); ok {
			introduce.PreEvalConclusion = v
		}
		introduce.UpdateTime = now
		return db.Save(&introduce).Error

	case "upgrade":
		var upgrade models.TaskUpgrade
		if err := db.Where("task_id = ?", taskID).First(&upgrade).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return createTaskDetail(db, taskID, taskType, detail)
			}
			return err
		}
		if v, ok := detail["component_name"].(string); ok {
			upgrade.ComponentName = v
		}
		if v, ok := detail["current_version"].(string); ok {
			upgrade.CurrentVersion = v
		}
		if v, ok := detail["target_version"].(string); ok {
			upgrade.TargetVersion = v
		}
		if v, ok := detail["upgrade_reason"].(string); ok {
			upgrade.UpgradeReason = v
		}
		if v, ok := detail["impact_analysis"].(string); ok {
			upgrade.ImpactAnalysis = v
		}
		if v, ok := detail["compatibility_risk"].(string); ok {
			upgrade.CompatibilityRisk = v
		}
		if v, ok := detail["rollback_plan"].(string); ok {
			upgrade.RollbackPlan = v
		}
		upgrade.UpdateTime = now
		return db.Save(&upgrade).Error

	case "withdraw":
		var withdraw models.TaskWithdraw
		if err := db.Where("task_id = ?", taskID).First(&withdraw).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return createTaskDetail(db, taskID, taskType, detail)
			}
			return err
		}
		if v, ok := detail["component_name"].(string); ok {
			withdraw.ComponentName = v
		}
		if v, ok := detail["withdraw_reason"].(string); ok {
			withdraw.WithdrawReason = v
		}
		if v, ok := detail["exit_reason"].(string); ok && withdraw.WithdrawReason == "" {
			withdraw.WithdrawReason = v
		}
		if v, ok := detail["replacement_component"].(string); ok {
			withdraw.ReplacementComponent = v
		}
		if v, ok := detail["alternative_plan"].(string); ok {
			withdraw.AlternativePlan = v
		}
		if v, ok := detail["code_cleanup_plan"].(string); ok {
			withdraw.CodeCleanupPlan = v
		}
		if v, ok := detail["affected_projects"].(string); ok {
			withdraw.AffectedProjects = v
		}
		withdraw.UpdateTime = now
		return db.Save(&withdraw).Error

	case "defect":
		var defect models.TaskDefect
		if err := db.Where("task_id = ?", taskID).First(&defect).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return createTaskDetail(db, taskID, taskType, detail)
			}
			return err
		}
		if v, ok := detail["risk_description"].(string); ok {
			defect.RiskDescription = v
		}
		defect.UpdateTime = now
		return db.Save(&defect).Error

	case "usage":
		var usage models.TaskUsage
		if err := db.Where("task_id = ?", taskID).First(&usage).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return createTaskDetail(db, taskID, taskType, detail)
			}
			return err
		}
		if v, ok := detail["component_name"].(string); ok {
			usage.ComponentName = v
		}
		if v, ok := detail["version"].(string); ok {
			usage.Version = v
		}
		if v, ok := detail["is_production"].(float64); ok {
			usage.IsProduction = int(v)
		}
		if v, ok := detail["purpose"].(string); ok {
			usage.Purpose = v
		}
		usage.UpdateTime = now
		return db.Save(&usage).Error

	default:
		return nil
	}
}
