package models

import "time"

// ComponentBlack 组件黑名单处理列表
type ComponentBlack struct {
	ID               int64      `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	ProjectID        int64      `gorm:"column:PROJECT_ID" json:"projectId"`
	ComponentID      *int64     `gorm:"column:COMPONENT_ID" json:"componentId"`
	BlackID          int64      `gorm:"column:BLACK_ID" json:"blackId"`
	ProcessStatus    int        `gorm:"column:PROCESS_STATUS" json:"processStatus"`
	Reason           string     `gorm:"column:reason" json:"reason"`
	JiraNo           string     `gorm:"column:jira_no" json:"jiraNo"`
	ProTime          *time.Time `gorm:"column:pro_time" json:"proTime"`
	TestTime         *time.Time `gorm:"column:test_time" json:"testTime"`
	Description      string     `gorm:"column:DESCRIPTION;type:mediumtext" json:"description"`
	ProcessPlanURL   *int64     `gorm:"column:PROCESS_PLAN_URL" json:"processPlanUrl"`
	ProcessResultURL *int64     `gorm:"column:PROCESS_RESULT_URL" json:"processResultUrl"`
	CreateBy         string     `gorm:"column:create_by" json:"createBy"`
	CreateTime       time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateBy         string     `gorm:"column:update_by" json:"updateBy"`
	UpdateTime       time.Time  `gorm:"column:update_time" json:"updateTime"`
	CurrentID        string     `gorm:"column:CURRENT_ID" json:"currentId"`
	Internal         int        `gorm:"column:internal" json:"internal"`
	AlertType        int        `gorm:"column:alert_type" json:"alertType"`
}

func (ComponentBlack) TableName() string {
	return "component_black"
}
