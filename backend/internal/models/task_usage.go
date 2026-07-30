package models

import "time"

// TaskUsage 使用申请扩展表
type TaskUsage struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID          int64     `gorm:"column:task_id" json:"taskId"`
	ComponentName   string    `gorm:"column:component_name" json:"componentName"`
	Version         string    `gorm:"column:version" json:"version"`
	IsProduction    int       `gorm:"column:is_production" json:"isProduction"`
	Purpose         string    `gorm:"column:purpose;type:mediumtext" json:"purpose"`
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime      time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (TaskUsage) TableName() string {
	return "task_usage"
}
