package workflow

import (
	"fmt"
	"osm-mngt/internal/models"
	"time"

	"gorm.io/gorm"
)

// WorkflowEngine 审批流引擎
type WorkflowEngine struct {
	db           *gorm.DB
	nodeResolver *NodeResolver
	stateMachine *StateMachine
}

// NewWorkflowEngine 创建审批流引擎
func NewWorkflowEngine(db *gorm.DB) *WorkflowEngine {
	return &WorkflowEngine{
		db:           db,
		nodeResolver: NewNodeResolver(db),
		stateMachine: NewStateMachine(),
	}
}

// ResolveNodesAndCreate 解析审批节点并创建到数据库
func (e *WorkflowEngine) ResolveNodesAndCreate(tx *gorm.DB, task *models.GovernanceTask, actionType string) error {
	// 解析节点
	nodes, err := e.nodeResolver.ResolveNodes(actionType, task.ApplicantID, task.ApplicantDept, task.ProjectDept)
	if err != nil {
		return fmt.Errorf("解析审批节点失败: %w", err)
	}

	// 创建节点到数据库
	for i := range nodes {
		nodes[i].TaskID = task.ID
		nodes[i].CreateTime = time.Now()
		nodes[i].UpdateTime = time.Now()
		if err := tx.Create(&nodes[i]).Error; err != nil {
			return fmt.Errorf("创建审批节点失败: %w", err)
		}
	}

	fmt.Printf("[ResolveNodesAndCreate] 任务 %d 创建了 %d 个审批节点\n", task.ID, len(nodes))
	for i, node := range nodes {
		fmt.Printf("  节点%d: key=%s, name=%s, approverId=%s, approver=%s, order=%d, skipped=%t\n",
			i+1, node.NodeKey, node.NodeName, node.ApproverID, node.Approver, node.NodeOrder, node.IsSkipped)
	}

	return nil
}

// SubmitTask 提交任务（从草稿状态提交，创建审批节点并进入审批中状态）
func (e *WorkflowEngine) SubmitTask(taskID int64, userID, userName string) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询任务
		var task models.GovernanceTask
		if err := tx.First(&task, taskID).Error; err != nil {
			return fmt.Errorf("任务不存在: %w", err)
		}

		// 2. 状态校验：只有草稿状态或已撤回状态可以提交
		if task.Status != StatusDraft && task.Status != StatusWithdrawn {
			return fmt.Errorf("只有草稿状态或已撤回状态的任务可以提交，当前状态: %s", task.Status)
		}

		// 2.1 权限校验：只有发起人可以提交
		if task.ApplicantID != userID {
			return fmt.Errorf("只有发起人可以提交任务")
		}

		// 3. 创建审批节点（首次提交时创建，撤回后提交时复用已有节点）
		fmt.Printf("[SubmitTask] 开始处理审批节点...\n")
		
		// 检查是否已有审批节点（撤回后再次提交）
		var existingNodes []models.ApprovalNode
		tx.Where("task_id = ?", taskID).Find(&existingNodes)
		
		if len(existingNodes) > 0 {
			// 已有审批节点，复用（撤回后再次提交）
			fmt.Printf("[SubmitTask] 发现已有 %d 个审批节点，复用并重置状态\n", len(existingNodes))
			
			// 重置所有节点状态为pending
			for i := range existingNodes {
				tx.Model(&models.ApprovalNode{}).Where("id = ?", existingNodes[i].ID).Updates(map[string]interface{}{
					"status":      "pending",
					"opinion":     "",
					"approved_at": nil,
					"is_skipped":  false,
					"skip_reason": "",
				})
			}
		} else {
			// 首次提交，创建审批节点
			fmt.Printf("[SubmitTask] 首次提交，创建审批节点...\n")
			
			// 对于非引入申请，从project表查询使用部门（use_dept）
			if task.ActionType != "introduce" && task.ProjectDept == "" {
				// 从task_detail获取组件名称，再查project表的use_dept
				var componentName string
				switch task.Type {
				case "upgrade":
					var detail models.TaskUpgrade
					if err := tx.Where("task_id = ?", task.ID).First(&detail).Error; err == nil {
						componentName = detail.ComponentName
					}
				case "withdraw":
					var detail models.TaskWithdraw
					if err := tx.Where("task_id = ?", task.ID).First(&detail).Error; err == nil {
						componentName = detail.ComponentName
					}
				}

				if componentName != "" {
					var useDept string
					tx.Table("project").Select("use_dept").Where("_name = ?", componentName).Row().Scan(&useDept)
					if useDept != "" {
						task.ProjectDept = useDept
						tx.Save(&task)
						fmt.Printf("[SubmitTask] 从project表查询到使用部门: %s (组件: %s)\n", task.ProjectDept, componentName)
					}
				}
			}
			
			if err := e.ResolveNodesAndCreate(tx, &task, task.ActionType); err != nil {
				fmt.Printf("[SubmitTask] 创建审批节点失败: %v\n", err)
				return fmt.Errorf("创建审批流失败: %w", err)
			}
		}

		// 查询创建的节点
		var nodes []models.ApprovalNode
		tx.Where("task_id = ?", taskID).Order("node_order ASC").Find(&nodes)
		fmt.Printf("[SubmitTask] 创建审批节点成功，共 %d 个节点\n", len(nodes))
		for i, node := range nodes {
			fmt.Printf("[SubmitTask] 节点%d: key=%s, name=%s, approver=%s, approverId=%s, order=%d\n",
				i+1, node.NodeKey, node.NodeName, node.Approver, node.ApproverID, node.NodeOrder)
		}

		// 4. 查找第一个未跳过且审批人非空的审批节点（跳过Node 0发起节点）
		var firstActiveNode models.ApprovalNode
		err := tx.Where("task_id = ? AND status = ? AND is_skipped = ? AND node_order > 0 AND approver_id != ''", taskID, "pending", false).
			Order("node_order ASC").First(&firstActiveNode).Error
		if err != nil {
			// 查询所有节点用于调试
			var allNodes []models.ApprovalNode
			tx.Where("task_id = ?", taskID).Order("node_order ASC").Find(&allNodes)
			fmt.Printf("[SubmitTask] 未找到可用审批节点，任务 %d 的所有节点:\n", taskID)
			for _, n := range allNodes {
				fmt.Printf("  - Node: key=%s, name=%s, approverId=%s, status=%s, isSkipped=%t, order=%d\n",
					n.NodeKey, n.NodeName, n.ApproverID, n.Status, n.IsSkipped, n.NodeOrder)
			}
			return fmt.Errorf("未找到可用的审批节点")
		}

		// 5. 更新任务状态
		now := time.Now()
		task.Status = StatusApproving
		task.CurrentStep = firstActiveNode.NodeName
		task.AssignedTo = firstActiveNode.ApproverID
		task.UpdateTime = now
		task.UpdateBy = userName
		if err := tx.Save(&task).Error; err != nil {
			return fmt.Errorf("更新任务状态失败: %w", err)
		}

		// 6. 记录提交历史
		history := models.ApprovalHistory{
			TaskID:       task.ID,
			TaskNo:       task.TaskNo,
			NodeID:       firstActiveNode.ID,
			NodeKey:      firstActiveNode.NodeKey,
			NodeName:     firstActiveNode.NodeName,
			ApproverID:   userID,
			Approver:     userName,
			ApproverDept: task.ApplicantDept,
			Action:       "submit",
			Opinion:      "提交审批",
			FromStatus:   StatusDraft,
			ToStatus:     StatusApproving,
			CreateTime:   now,
		}
		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("记录审批历史失败: %w", err)
		}

		fmt.Printf("[SubmitTask] 任务 %d 提交成功，当前节点: %s, 审批人: %s\n",
			taskID, firstActiveNode.NodeName, firstActiveNode.Approver)
		return nil
	})
}

// ApproveNode 审批通过
func (e *WorkflowEngine) ApproveNode(taskID int64, nodeID int64, userID, userName, opinion string) error {
	fmt.Printf("[WorkflowEngine.ApproveNode] 入参: taskID=%d, nodeID=%d, userID=%s, userName=%s\n", taskID, nodeID, userID, userName)
	return e.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询任务
		var task models.GovernanceTask
		if err := tx.First(&task, taskID).Error; err != nil {
			return fmt.Errorf("任务不存在: %w", err)
		}

		// 2. 状态校验
		if task.Status != StatusApproving {
			return fmt.Errorf("当前状态 %s 不允许审批", task.Status)
		}

		// 3. 查询当前节点
		var currentNode models.ApprovalNode
		if err := tx.First(&currentNode, nodeID).Error; err != nil {
			return fmt.Errorf("审批节点不存在")
		}

		// 4. 权限校验
		if !e.stateMachine.CanApprove(task.Status, currentNode.ApproverID, userID) {
			return fmt.Errorf("您不是当前节点的审批人，无权审批")
		}

		// 5. 更新当前节点状态
		now := time.Now()
		currentNode.Status = "approved"
		currentNode.Opinion = opinion
		currentNode.ApprovedAt = &now
		currentNode.UpdateTime = now
		if err := tx.Save(&currentNode).Error; err != nil {
			return fmt.Errorf("更新节点状态失败: %w", err)
		}

		// 6. 查找下一个未跳过且审批人非空的节点
		var nextNode models.ApprovalNode
		err := tx.Where("task_id = ? AND status = ? AND is_skipped = ? AND node_order > ? AND approver_id != ''",
			taskID, "pending", false, currentNode.NodeOrder).
			Order("node_order ASC").First(&nextNode).Error

		var newStatus string
		var assignedTo string
		var currentStep string

		if err == nil {
			// 还有下一个节点
			newStatus = StatusApproving
			assignedTo = nextNode.ApproverID
			currentStep = nextNode.NodeName
		} else {
			// 所有节点已完成
			newStatus = StatusApproved
			assignedTo = ""
			currentStep = "已通过"
		}

		// 7. 更新任务状态
		oldStatus := task.Status
		task.Status = newStatus
		task.CurrentStep = currentStep
		task.AssignedTo = assignedTo
		task.UpdateTime = now
		task.UpdateBy = userID
		if err := tx.Save(&task).Error; err != nil {
			return fmt.Errorf("更新任务状态失败: %w", err)
		}

		// 8. 记录审批历史
		fmt.Printf("[WorkflowEngine.ApproveNode] 保存审批历史, userName=%s, currentNode.Approver=%s\n", userName, currentNode.Approver)
		history := models.ApprovalHistory{
			TaskID:       task.ID,
			TaskNo:       task.TaskNo,
			NodeID:       currentNode.ID,
			NodeKey:      currentNode.NodeKey,
			NodeName:     currentNode.NodeName,
			ApproverID:   userID,
			Approver:     userName,
			ApproverDept: currentNode.ApproverDept,
			Action:       ActionApprove,
			Opinion:      opinion,
			FromStatus:   oldStatus,
			ToStatus:     newStatus,
			CreateTime:   now,
		}
		return tx.Create(&history).Error
	})
}

// RejectTask 驳回任务
func (e *WorkflowEngine) RejectTask(taskID int64, nodeID int64, userID, userName, opinion string) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询任务
		var task models.GovernanceTask
		if err := tx.First(&task, taskID).Error; err != nil {
			return fmt.Errorf("任务不存在: %w", err)
		}

		// 2. 状态校验
		if task.Status != StatusApproving {
			return fmt.Errorf("当前状态 %s 不允许驳回", task.Status)
		}

		// 3. 查询当前节点
		var currentNode models.ApprovalNode
		if err := tx.First(&currentNode, nodeID).Error; err != nil {
			return fmt.Errorf("审批节点不存在")
		}

		// 4. 权限校验
		if !e.stateMachine.CanReject(task.Status, currentNode.ApproverID, userID) {
			return fmt.Errorf("您不是当前节点的审批人，无权驳回")
		}

		// 5. 更新当前节点状态
		now := time.Now()
		currentNode.Status = "rejected"
		currentNode.Opinion = opinion
		currentNode.ApprovedAt = &now
		currentNode.UpdateTime = now
		if err := tx.Save(&currentNode).Error; err != nil {
			return fmt.Errorf("更新节点状态失败: %w", err)
		}

		// 6. 更新任务状态为已驳回（终态）
		oldStatus := task.Status
		task.Status = StatusRejected
		task.CurrentStep = "已驳回"
		task.AssignedTo = ""
		task.UpdateTime = now
		task.UpdateBy = userID
		if err := tx.Save(&task).Error; err != nil {
			return fmt.Errorf("更新任务状态失败: %w", err)
		}

		// 7. 记录驳回历史
		history := models.ApprovalHistory{
			TaskID:       task.ID,
			TaskNo:       task.TaskNo,
			NodeID:       currentNode.ID,
			NodeKey:      currentNode.NodeKey,
			NodeName:     currentNode.NodeName,
			ApproverID:   userID,
			Approver:     userName,
			ApproverDept: currentNode.ApproverDept,
			Action:       ActionReject,
			Opinion:      opinion,
			FromStatus:   oldStatus,
			ToStatus:     StatusRejected,
			CreateTime:   now,
		}
		return tx.Create(&history).Error
	})
}

// GetApprovalHistory 获取审批历史（从最后一次提交开始显示）
func (e *WorkflowEngine) GetApprovalHistory(taskID int64) ([]models.ApprovalHistory, error) {
	var histories []models.ApprovalHistory
	// 1. 查找所有历史记录
	err := e.db.Where("task_id = ?", taskID).
		Order("create_time ASC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}

	// 2. 如果没有记录，直接返回
	if len(histories) == 0 {
		return histories, nil
	}

	// 3. 找到最后一次 submit 记录的位置
	lastSubmitIndex := -1
	for i := len(histories) - 1; i >= 0; i-- {
		if histories[i].Action == "submit" {
			lastSubmitIndex = i
			break
		}
	}

	// 4. 如果找到 submit 记录，只返回该记录及之后的所有记录
	if lastSubmitIndex != -1 {
		return histories[lastSubmitIndex:], nil
	}

	// 5. 没有 submit 记录，返回所有记录
	return histories, nil
}

// GetApprovalNodes 获取审批节点列表
func (e *WorkflowEngine) GetApprovalNodes(taskID int64) ([]models.ApprovalNode, error) {
	var nodes []models.ApprovalNode
	err := e.db.Where("task_id = ?", taskID).
		Order("node_order ASC").
		Find(&nodes).Error
	return nodes, err
}

// GetCurrentNode 获取当前待处理的审批节点（审批人非空）
func (e *WorkflowEngine) GetCurrentNode(taskID int64) (*models.ApprovalNode, error) {
	var node models.ApprovalNode
	err := e.db.Where("task_id = ? AND status = ? AND is_skipped = ? AND approver_id != ''",
		taskID, "pending", false).
		Order("node_order ASC").
		First(&node).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

// WithdrawTask 撤回任务
func (e *WorkflowEngine) WithdrawTask(taskID int64, userID, userName string) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询任务
		var task models.GovernanceTask
		if err := tx.First(&task, taskID).Error; err != nil {
			return fmt.Errorf("任务不存在: %w", err)
		}

		// 2. 状态校验
		if task.Status != StatusApproving {
			return fmt.Errorf("当前状态 %s 不允许撤回", task.Status)
		}

		// 3. 权限校验
		if task.ApplicantID != userID {
			return fmt.Errorf("只有发起人可以撤回任务")
		}

		// 4. 查询当前节点
		currentNode, err := e.GetCurrentNode(taskID)
		if err != nil {
			return fmt.Errorf("未找到当前审批节点")
		}

		// 5. 检查是否允许撤回（第一个实际审批节点尚未有审批记录）
		// 查找第一个未跳过且审批人非空的实际审批节点（node_order > 0）
		var firstApprovalNode models.ApprovalNode
		err = tx.Where("task_id = ? AND node_order > 0 AND is_skipped = ? AND approver_id != ''", taskID, false).
			Order("node_order ASC").First(&firstApprovalNode).Error
		if err != nil {
			return fmt.Errorf("未找到可撤回的审批节点")
		}
		
		// 检查该节点是否已有审批记录
		if e.nodeHasAction(tx, firstApprovalNode.ID) {
			return fmt.Errorf("审批人已处理，不允许撤回")
		}

		// 6. 更新任务状态（撤回到草稿状态，允许再次提交）
		now := time.Now()
		oldStatus := task.Status
		
		fmt.Printf("[WithdrawTask] 准备更新任务 %d 状态: %s -> draft\n", taskID, oldStatus)
		
		// 使用 Update 方法明确更新字段
		if err := tx.Model(&models.GovernanceTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
			"status":        StatusDraft,
			"current_step":  "",
			"assigned_to":   "",
			"update_time":   now,
			"update_by":     userID,
		}).Error; err != nil {
			fmt.Printf("[WithdrawTask] 更新任务状态失败: %v\n", err)
			return fmt.Errorf("更新任务状态失败: %w", err)
		}
		
		// 验证更新是否成功
		var verifyTask models.GovernanceTask
		tx.First(&verifyTask, taskID)
		fmt.Printf("[WithdrawTask] 数据库验证: 任务 %d 状态已更新为: %s (之前是: %s)\n", taskID, verifyTask.Status, oldStatus)

		// 7. 记录撤回历史
		history := models.ApprovalHistory{
			TaskID:       task.ID,
			TaskNo:       task.TaskNo,
			NodeID:       currentNode.ID,
			NodeKey:      currentNode.NodeKey,
			NodeName:     currentNode.NodeName,
			ApproverID:   userID,
			Approver:     userName,
			ApproverDept: task.ApplicantDept,
			Action:       ActionWithdraw,
			Opinion:       "撤回任务",
			FromStatus:   StatusApproving,
			ToStatus:     StatusWithdrawn,
			CreateTime:   now,
		}
		return tx.Create(&history).Error
	})
}

// nodeHasAction 检查节点是否已有操作记录
func (e *WorkflowEngine) nodeHasAction(tx *gorm.DB, nodeID int64) bool {
	var count int64
	tx.Model(&models.ApprovalHistory{}).
		Where("node_id = ? AND action IN (?, ?)", nodeID, ActionApprove, ActionReject).
		Count(&count)
	return count > 0
}
