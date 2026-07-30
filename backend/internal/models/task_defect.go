package models

import "time"

// TaskDefect 风险上报扩展表
type TaskDefect struct {
	ID                   int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID               int64      `gorm:"column:task_id" json:"taskId"`
	RiskDescription      string     `gorm:"column:risk_description;type:mediumtext" json:"riskDescription"`
	CreateTime           time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime           time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (TaskDefect) TableName() string {
	return "task_defect"
}
