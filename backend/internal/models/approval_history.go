package models

import "time"

// ApprovalHistory 审批历史表
// 记录每个审批节点的操作历史（审批人、时间、意见等）
type ApprovalHistory struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID        int64     `gorm:"column:task_id;index" json:"taskId"`                               // 关联的任务ID
	TaskNo        string    `gorm:"column:task_no;size:50" json:"taskNo"`                             // 任务编号
	NodeID        int64     `gorm:"column:node_id;index" json:"nodeId"`                               // 关联的审批节点ID
	NodeKey       string    `gorm:"column:node_key;size:50" json:"nodeKey"`                           // 节点标识
	NodeName      string    `gorm:"column:node_name;size:100" json:"nodeName"`                        // 节点名称
	ApproverID    string    `gorm:"column:approver_id;size:50" json:"approverId"`                     // 操作人ID
	Approver      string    `gorm:"column:approver;size:100" json:"approver"`                         // 操作人姓名
	ApproverDept  string    `gorm:"column:approver_dept;size:100" json:"approverDept"`                // 操作人部门
	Action        string    `gorm:"column:action;size:20" json:"action"`                              // 操作类型：submit(提交), approve(同意), reject(驳回), withdraw(撤回)
	Opinion       string    `gorm:"column:opinion;size:1000" json:"opinion"`                          // 审批意见
	FromStatus    string    `gorm:"column:from_status;size:20" json:"fromStatus"`                     // 操作前状态
	ToStatus      string    `gorm:"column:to_status;size:20" json:"toStatus"`                         // 操作后状态
	CreateTime    time.Time `gorm:"column:create_time" json:"createTime"`
}

func (ApprovalHistory) TableName() string {
	return "approval_history"
}
