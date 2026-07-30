package models

import "time"

// ApprovalNode 审批节点记录表
// 记录每个任务实例化后的实际审批节点
type ApprovalNode struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID        int64      `gorm:"column:task_id;index" json:"taskId"`                                   // 关联的任务ID
	NodeKey       string     `gorm:"column:node_key;size:50" json:"nodeKey"`                               // 节点标识（来自模板）
	NodeName      string     `gorm:"column:node_name;size:100" json:"nodeName"`                            // 节点名称
	NodeOrder     int        `gorm:"column:node_order" json:"nodeOrder"`                                   // 节点顺序
	NodeType      string     `gorm:"column:node_type;size:20" json:"nodeType"`                             // 节点类型：fixed, conditional
	ApproverID    string     `gorm:"column:approver_id;size:50" json:"approverId"`                         // 审批人ID
	Approver      string     `gorm:"column:approver;size:100" json:"approver"`                             // 审批人姓名
	ApproverDept  string     `gorm:"column:approver_dept;size:100" json:"approverDept"`                    // 审批人部门
	Status        string     `gorm:"column:status;size:20;default:pending" json:"status"`                  // 节点状态：pending, approved, rejected, skipped
	Opinion       string     `gorm:"column:opinion;size:500" json:"opinion"`                               // 审批意见
	ApprovedAt    *time.Time `gorm:"column:approved_at" json:"approvedAt"`                                 // 审批时间
	IsSkipped     bool       `gorm:"column:is_skipped;default:false" json:"isSkipped"`                     // 是否被跳过
	SkipReason    string     `gorm:"column:skip_reason;size:200" json:"skipReason"`                        // 跳过原因
	CreateTime    time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime    time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (ApprovalNode) TableName() string {
	return "approval_node"
}
