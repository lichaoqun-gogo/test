package workflow

import (
	"errors"
	"fmt"
)

// TaskStatus 任务状态定义（五态模型）
const (
	StatusDraft      = "draft"      // 草稿：仅发起人可见、可编辑
	StatusApproving  = "approving"  // 审批中：数据锁定，控制权移交审批流引擎
	StatusWithdrawn  = "withdrawn"  // 已撤回：可重提态
	StatusApproved   = "approved"   // 已通过：终态
	StatusRejected   = "rejected"   // 已驳回：终态，流程彻底终结
)

// TaskAction 任务操作定义
const (
	ActionSubmit   = "submit"   // 提交：草稿 -> 审批中
	ActionWithdraw = "withdraw" // 撤回：审批中 -> 草稿（仅第一节点未处理时）
	ActionApprove  = "approve"  // 同意：流转至下一节点或终态
	ActionReject   = "reject"   // 驳回：直接结束，状态为已驳回
)

// StateMachine 审批状态机
// 严格控制状态转换的触发条件与按钮权限
type StateMachine struct{}

// NewStateMachine 创建状态机
func NewStateMachine() *StateMachine {
	return &StateMachine{}
}

// CanTransition 检查状态转换是否合法
// 参数：
//   - currentStatus: 当前状态
//   - action: 要执行的操作
//   - currentNodeOrder: 当前节点顺序（0-based）
//   - isFirstNode: 是否在第一节点
//   - hasNodeAction: 当前节点是否已有操作记录
func (sm *StateMachine) CanTransition(currentStatus, action string, currentNodeOrder int, isFirstNode, hasNodeAction bool) error {
	switch currentStatus {
	case StatusDraft:
		// 草稿状态只能提交
		if action != ActionSubmit {
			return fmt.Errorf("草稿状态只能执行提交操作，当前操作：%s", action)
		}
		return nil

	case StatusApproving:
		switch action {
		case ActionSubmit:
			return errors.New("审批中状态不能重复提交")
		case ActionWithdraw:
			// 撤回逻辑：仅当 当前处理节点 == 第一节点 且 当前处理人尚未执行任何操作 时，允许发起人撤回
			if !isFirstNode {
				return errors.New("仅第一节点允许撤回")
			}
			if hasNodeAction {
				return errors.New("当前节点已有操作记录，不允许撤回")
			}
			return nil
		case ActionApprove, ActionReject:
			// 审批动作由审批人执行，需要额外校验权限
			return nil
		default:
			return fmt.Errorf("未知操作：%s", action)
		}

	case StatusWithdrawn:
		// 已撤回状态可以重新提交
		if action != ActionSubmit {
			return fmt.Errorf("已撤回状态只能重新提交，当前操作：%s", action)
		}
		return nil

	case StatusApproved:
		return errors.New("已通过状态为终态，不允许任何操作")

	case StatusRejected:
		return errors.New("已驳回状态为终态，不允许任何操作")

	default:
		return fmt.Errorf("未知状态：%s", currentStatus)
	}
}

// GetNextStatus 获取操作后的下一个状态
// 参数：
//   - currentStatus: 当前状态
//   - action: 执行的操作
//   - isLastNode: 是否是最后一个审批节点
//   - isReject: 是否是驳回操作
func (sm *StateMachine) GetNextStatus(currentStatus, action string, isLastNode, isReject bool) (string, error) {
	switch currentStatus {
	case StatusDraft:
		if action == ActionSubmit {
			return StatusApproving, nil
		}
		return "", fmt.Errorf("草稿状态不支持操作：%s", action)

	case StatusApproving:
		switch action {
		case ActionWithdraw:
			return StatusWithdrawn, nil
		case ActionApprove:
			if isLastNode {
				// 最后一个节点同意 -> 已通过
				return StatusApproved, nil
			}
			// 非最后节点同意 -> 仍为审批中
			return StatusApproving, nil
		case ActionReject:
			// 驳回 -> 已驳回（终态）
			return StatusRejected, nil
		default:
			return "", fmt.Errorf("审批中状态不支持操作：%s", action)
		}

	case StatusWithdrawn:
		if action == ActionSubmit {
			return StatusApproving, nil
		}
		return "", fmt.Errorf("已撤回状态不支持操作：%s", action)

	default:
		return "", fmt.Errorf("状态 %s 不支持操作：%s", currentStatus, action)
	}
}

// CanWithdraw 检查是否允许撤回
// 撤回条件：
// 1. 当前状态为审批中
// 2. 当前节点是第一节点（node_order == 0）
// 3. 当前节点尚未有任何操作记录
func (sm *StateMachine) CanWithdraw(currentStatus string, currentNodeOrder int, hasAction bool) bool {
	if currentStatus != StatusApproving {
		return false
	}
	if currentNodeOrder != 0 {
		return false
	}
	if hasAction {
		return false
	}
	return true
}

// CanApprove 检查当前用户是否有审批权限
// 审批条件：
// 1. 当前状态为审批中
// 2. 当前用户是当前节点的处理人
func (sm *StateMachine) CanApprove(currentStatus string, currentApproverID, userID string) bool {
	if currentStatus != StatusApproving {
		return false
	}
	return currentApproverID == userID
}

// CanReject 检查当前用户是否有驳回权限
// 驳回条件同审批
func (sm *StateMachine) CanReject(currentStatus string, currentApproverID, userID string) bool {
	return sm.CanApprove(currentStatus, currentApproverID, userID)
}
