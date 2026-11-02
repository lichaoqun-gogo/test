package models

import "time"

// TaskWithdraw 退出申请扩展表
type TaskWithdraw struct {
	ID                 int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID             int64      `gorm:"column:task_id" json:"taskId"`
	ComponentID        *int64     `gorm:"column:component_id" json:"componentId"`
	ComponentName      string     `gorm:"column:component_name" json:"componentName"`
	ComponentVersion   string     `gorm:"column:component_version" json:"componentVersion"`
	WithdrawReason     string     `gorm:"column:withdraw_reason;type:mediumtext" json:"withdrawReason"`
	ReplacementComponent string   `gorm:"column:replacement_component" json:"replacementComponent"`
	AlternativePlan    string     `gorm:"column:alternative_plan;type:mediumtext" json:"alternativePlan"`
	CodeCleanupPlan    string     `gorm:"column:code_cleanup_plan;type:mediumtext" json:"codeCleanupPlan"`
	AffectedProjects   string     `gorm:"column:affected_projects;type:mediumtext" json:"affectedProjects"`
	CreateTime         time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime         time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (TaskWithdraw) TableName() string {
	return "task_withdraw"
}
