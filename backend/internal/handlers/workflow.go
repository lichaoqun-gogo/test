package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"osm-mngt/internal/service/workflow"
	"osm-mngt/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// decodeHeader 对 URI 编码的请求头值进行解码
func decodeHeader(value string) string {
	if value == "" {
		return ""
	}
	decoded, err := url.QueryUnescape(value)
	if err != nil {
		return value
	}
	return decoded
}

// WorkflowHandler 审批流处理器
// 处理审批流相关的HTTP请求
type WorkflowHandler struct {
	db      *gorm.DB
	engine  *workflow.WorkflowEngine
}

// NewWorkflowHandler 创建审批流处理器
func NewWorkflowHandler(db *gorm.DB) *WorkflowHandler {
	return &WorkflowHandler{
		db:     db,
		engine: workflow.NewWorkflowEngine(db),
	}
}

// SubmitTaskRequest 提交任务请求
type SubmitTaskRequest struct {
	TaskID int64 `json:"taskId" binding:"required"`
}

// ApproveNodeRequest 审批节点请求
type ApproveNodeRequest struct {
	NodeID int64  `json:"nodeId" binding:"required"`
	Opinion string `json:"opinion"`
}

// RejectTaskRequest 驳回任务请求
type RejectTaskRequest struct {
	NodeID int64  `json:"nodeId" binding:"required"`
	Opinion string `json:"opinion" binding:"required"`
}

// =====================================================
// HTTP Handler 方法
// =====================================================

// SubmitTask 提交任务（草稿/已撤回 -> 审批中）
// POST /api/v1/workflow/tasks/:id/submit
// 请求头：X-User-ID, X-User-Name
func (h *WorkflowHandler) SubmitTask(c *gin.Context) {
	fmt.Printf("[WorkflowHandler.SubmitTask] 收到提交请求，任务ID: %s\n", c.Param("id"))
	
	// 1. 获取任务ID
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "任务ID无效")
		return
	}

	// 2. 从请求头获取用户信息
	userID := c.GetHeader("X-User-ID")
	userName := decodeHeader(c.GetHeader("X-User-Name"))
	fmt.Printf("[WorkflowHandler.SubmitTask] 用户信息: userID=%s, userName=%s\n", userID, userName)
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "未登录")
		return
	}
	if userName == "" {
		response.Error(c, http.StatusBadRequest, "请求头缺少用户名(X-User-Name)")
		return
	}

	// 3. 执行提交
	fmt.Printf("[WorkflowHandler.SubmitTask] 调用 engine.SubmitTask...\n")
	if err := h.engine.SubmitTask(taskID, userID, userName); err != nil {
		fmt.Printf("[WorkflowHandler.SubmitTask] 提交失败: %v\n", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Printf("[WorkflowHandler.SubmitTask] 提交成功\n")

	response.SuccessWithMessage(c, "提交成功", nil)
}

// WithdrawTask 撤回任务（审批中 -> 已撤回）
// POST /api/v1/workflow/tasks/:id/withdraw
// 请求头：X-User-ID, X-User-Name
func (h *WorkflowHandler) WithdrawTask(c *gin.Context) {
	// 1. 获取任务ID
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "任务ID无效")
		return
	}

	// 2. 从请求头获取用户信息
	userID := c.GetHeader("X-User-ID")
	userName := decodeHeader(c.GetHeader("X-User-Name"))
	fmt.Printf("[WorkflowHandler.WithdrawTask] 请求头: taskID=%d, userID=%s, userName=%s\n", taskID, userID, userName)
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "未登录")
		return
	}
	if userName == "" {
		response.Error(c, http.StatusBadRequest, "请求头缺少用户名(X-User-Name)")
		return
	}

	// 3. 执行撤回
	if err := h.engine.WithdrawTask(taskID, userID, userName); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "撤回成功", nil)
}

// ApproveNode 审批通过节点
// POST /api/v1/workflow/tasks/:id/approve
// 请求头：X-User-ID, X-User-Name
func (h *WorkflowHandler) ApproveNode(c *gin.Context) {
	// 1. 获取任务ID
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "任务ID无效")
		return
	}

	// 2. 解析请求体
	var req ApproveNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误："+err.Error())
		return
	}

	// 3. 从请求头获取用户信息
	userID := c.GetHeader("X-User-ID")
	userName := decodeHeader(c.GetHeader("X-User-Name"))
	fmt.Printf("[WorkflowHandler.ApproveNode] 请求头: taskID=%d, userID=%s, userName=%s\n", taskID, userID, userName)
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "未登录")
		return
	}
	if userName == "" {
		response.Error(c, http.StatusBadRequest, "请求头缺少用户名(X-User-Name)")
		return
	}

	// 4. 执行审批
	if err := h.engine.ApproveNode(taskID, req.NodeID, userID, userName, req.Opinion); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "审批通过", nil)
}

// RejectTask 驳回任务
// POST /api/v1/workflow/tasks/:id/reject
// 请求头：X-User-ID, X-User-Name
func (h *WorkflowHandler) RejectTask(c *gin.Context) {
	// 1. 获取任务ID
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "任务ID无效")
		return
	}

	// 2. 解析请求体
	var req RejectTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误："+err.Error())
		return
	}

	// 3. 从请求头获取用户信息
	userID := c.GetHeader("X-User-ID")
	userName := decodeHeader(c.GetHeader("X-User-Name"))
	fmt.Printf("[WorkflowHandler.RejectTask] 请求头: taskID=%d, userID=%s, userName=%s\n", taskID, userID, userName)
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "未登录")
		return
	}
	if userName == "" {
		response.Error(c, http.StatusBadRequest, "请求头缺少用户名(X-User-Name)")
		return
	}

	// 4. 执行驳回
	if err := h.engine.RejectTask(taskID, req.NodeID, userID, userName, req.Opinion); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "驳回成功", nil)
}

// GetApprovalHistory 获取审批历史
// GET /api/v1/workflow/tasks/:id/history
func (h *WorkflowHandler) GetApprovalHistory(c *gin.Context) {
	// 1. 获取任务ID
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "任务ID无效")
		return
	}

	// 2. 查询审批历史
	histories, err := h.engine.GetApprovalHistory(taskID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "查询审批历史失败")
		return
	}

	response.Success(c, histories)
}

// GetApprovalNodes 获取审批节点列表
// GET /api/v1/workflow/tasks/:id/nodes
func (h *WorkflowHandler) GetApprovalNodes(c *gin.Context) {
	// 1. 获取任务ID
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "任务ID无效")
		return
	}

	// 2. 查询审批节点
	nodes, err := h.engine.GetApprovalNodes(taskID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "查询审批节点失败")
		return
	}

	response.Success(c, nodes)
}

// GetCurrentNode 获取当前待处理节点
// GET /api/v1/workflow/tasks/:id/current-node
func (h *WorkflowHandler) GetCurrentNode(c *gin.Context) {
	// 1. 获取任务ID
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "任务ID无效")
		return
	}

	// 2. 查询当前节点
	node, err := h.engine.GetCurrentNode(taskID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "当前无待处理节点")
		return
	}

	response.Success(c, node)
}
