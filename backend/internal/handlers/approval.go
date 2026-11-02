package handlers

import (
	"fmt"
	"osm-mngt/internal/models"
	"osm-mngt/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApprovalHandler struct {
	db *gorm.DB
}

func NewApprovalHandler(db *gorm.DB) *ApprovalHandler {
	return &ApprovalHandler{db: db}
}

func (h *ApprovalHandler) getUserInfo(c *gin.Context) (string, string, string) {
	userID := c.GetHeader("X-User-ID")
	userName := c.GetHeader("X-User-Name")
	userDept := c.GetHeader("X-User-Dept")
	return userID, userName, userDept
}

// Approve 审批通过
func (h *ApprovalHandler) Approve(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("taskId"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的任务ID")
		return
	}

	var req struct {
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID, userName, userDept := h.getUserInfo(c)

	// 使用事务
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询任务
		var task models.GovernanceTask
		if err := tx.First(&task, taskID).Error; err != nil {
			return fmt.Errorf("任务不存在: %w", err)
		}

		// 2. 查找当前待审批的节点
		var currentNode models.ApprovalNode
		if err := tx.Where("task_id = ? AND status = ?", taskID, "pending").
			Order("node_order ASC").First(&currentNode).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return gorm.ErrRecordNotFound
			}
			return err
		}

		now := time.Now()

		// 3. 记录审批历史
		history := models.ApprovalHistory{
			TaskID:       taskID,
			TaskNo:       task.TaskNo,
			NodeID:       currentNode.ID,
			NodeKey:      currentNode.NodeKey,
			NodeName:     currentNode.NodeName,
			ApproverID:   userID,
			Approver:     userName,
			ApproverDept: userDept,
			Action:       "approve",
			Opinion:      req.Comment,
			FromStatus:   task.Status,
			ToStatus:     "approved",
			CreateTime:   now,
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		// 3. 更新节点状态
		if err := tx.Model(&currentNode).Updates(map[string]interface{}{
			"status":     "approved",
			"approver":   userName,
			"opinion":    req.Comment,
			"approved_at": now,
		}).Error; err != nil {
			return err
		}

		// 4. 检查是否所有节点都已审批
		var pendingCount int64
		tx.Model(&models.ApprovalNode{}).
			Where("task_id = ? AND status = ?", taskID, "pending").
			Count(&pendingCount)

		if pendingCount == 0 {
			// 所有节点都已审批，更新任务状态为completed
			tx.Model(&models.GovernanceTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
				"status":     "completed",
				"update_time": now,
			})
		} else {
			// 还有节点待审批，更新任务状态为processing
			tx.Model(&models.GovernanceTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
				"status":     "processing",
				"update_time": now,
			})
		}

		return nil
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 400, "没有待审批的节点")
			return
		}
		response.Error(c, 500, "审批失败")
		return
	}

	response.Success(c, nil)
}

// Reject 审批驳回
func (h *ApprovalHandler) Reject(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("taskId"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的任务ID")
		return
	}

	var req struct {
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID, userName, userDept := h.getUserInfo(c)

	// 使用事务
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询任务
		var task models.GovernanceTask
		if err := tx.First(&task, taskID).Error; err != nil {
			return fmt.Errorf("任务不存在: %w", err)
		}

		// 2. 查找当前待审批的节点
		var currentNode models.ApprovalNode
		if err := tx.Where("task_id = ? AND status = ?", taskID, "pending").
			Order("node_order ASC").First(&currentNode).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return gorm.ErrRecordNotFound
			}
			return err
		}

		now := time.Now()

		// 3. 记录审批历史
		history := models.ApprovalHistory{
			TaskID:       taskID,
			TaskNo:       task.TaskNo,
			NodeID:       currentNode.ID,
			NodeKey:      currentNode.NodeKey,
			NodeName:     currentNode.NodeName,
			ApproverID:   userID,
			Approver:     userName,
			ApproverDept: userDept,
			Action:       "reject",
			Opinion:      req.Comment,
			FromStatus:   task.Status,
			ToStatus:     "rejected",
			CreateTime:   now,
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		// 3. 更新任务状态为rejected
		if err := tx.Model(&models.GovernanceTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
			"status":     "rejected",
			"update_time": now,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 400, "没有待审批的节点")
			return
		}
		response.Error(c, 500, "驳回失败")
		return
	}

	response.Success(c, nil)
}

// GetApprovalHistory 获取审批历史
func (h *ApprovalHandler) GetApprovalHistory(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("taskId"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的任务ID")
		return
	}

	var history []models.ApprovalHistory
	if err := h.db.Where("task_id = ?", taskID).Order("create_time DESC").Find(&history).Error; err != nil {
		response.Error(c, 500, "查询审批历史失败")
		return
	}

	response.Success(c, history)
}
