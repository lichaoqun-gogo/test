package models

import "time"

// ApprovalTemplate 审批流模板表
// 定义不同类型的审批流程模板
type ApprovalTemplate struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TemplateKey string    `gorm:"column:template_key;uniqueIndex;size:50" json:"templateKey"` // 模板唯一标识，如：introduce, upgrade, withdraw, usage, risk_report
	TemplateName string   `gorm:"column:template_name;size:100" json:"templateName"`           // 模板名称，如：引入申请审批流
	Description string    `gorm:"column:description;size:500" json:"description"`              // 模板描述
	IsActive    bool      `gorm:"column:is_active;default:true" json:"isActive"`               // 是否启用
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (ApprovalTemplate) TableName() string {
	return "approval_template"
}

// ApprovalTemplateNode 审批流模板节点定义表
// 定义每个模板包含的审批节点及其条件
type ApprovalTemplateNode struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TemplateID    int64     `gorm:"column:template_id;index" json:"templateId"`                        // 关联的模板ID
	NodeKey       string    `gorm:"column:node_key;size:50" json:"nodeKey"`                            // 节点唯一标识，如：node_0, node_1, node_2, node_3
	NodeName      string    `gorm:"column:node_name;size:100" json:"nodeName"`                         // 节点名称
	NodeOrder     int       `gorm:"column:node_order" json:"nodeOrder"`                                // 节点顺序（从0开始）
	NodeType      string    `gorm:"column:node_type;size:20;default:fixed" json:"nodeType"`            // 节点类型：fixed(固定节点), conditional(条件节点)
	ConditionType string    `gorm:"column:condition_type;size:50" json:"conditionType"`                // 条件类型：dept_specialist(本部门专员), usage_dept(使用部门), mgmt_dept(管理部门)
	IsSkipable    bool      `gorm:"column:is_skipable;default:false" json:"isSkipable"`                // 是否可跳过
	SkipRule      string    `gorm:"column:skip_rule;size:200" json:"skipRule"`                         // 跳过规则JSON，如：{"type":"same_person","field":"approver_id"}
	ApproverRole  string    `gorm:"column:approver_role;size:50" json:"approverRole"`                  // 审批人角色：initiator(发起人), dept_specialist(部门专员), usage_dept_specialist(使用部门专员), mgmt_dept_specialist(管理部门专员)
	Description   string    `gorm:"column:description;size:500" json:"description"`                    // 节点描述
	CreateTime    time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime    time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (ApprovalTemplateNode) TableName() string {
	return "approval_template_node"
}
