package models

import "time"

// ComponentBlack 黑名单阻断记录表
type ComponentBlack struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BlacklistID     int64      `gorm:"column:blacklist_id" json:"blacklistId"`
	ComponentName   string     `gorm:"column:component_name" json:"componentName"`
	ComponentVersion string    `gorm:"column:component_version" json:"componentVersion"`
	ComponentGroup  string     `gorm:"column:component_group" json:"componentGroup"`
	ProjectID       *int64     `gorm:"column:project_id" json:"projectId"`
	ProjectName     string     `gorm:"column:project_name" json:"projectName"`
	AlertType       int        `gorm:"column:alert_type" json:"alertType"`
	Status          string     `gorm:"column:status" json:"status"`
	Handler         string     `gorm:"column:handler" json:"handler"`
	HandlePlan      string     `gorm:"column:handle_plan" json:"handlePlan"`
	HandleResult    string     `gorm:"column:handle_result;type:mediumtext" json:"handleResult"`
	AttachmentURL   string     `gorm:"column:attachment_url" json:"attachmentUrl"`
	Auditor         string     `gorm:"column:auditor" json:"auditor"`
	AuditOpinion    string     `gorm:"column:audit_opinion" json:"auditOpinion"`
	AuditTime       *time.Time `gorm:"column:audit_time" json:"auditTime"`
	CIPipelineID    string     `gorm:"column:ci_pipeline_id" json:"ciPipelineId"`
	CIBuildID       string     `gorm:"column:ci_build_id" json:"ciBuildId"`
	CreateTime      time.Time  `gorm:"column:create_time" json:"createTime"`
	CreateBy        string     `gorm:"column:create_by" json:"createBy"`
	UpdateTime      time.Time  `gorm:"column:update_time" json:"updateTime"`
	UpdateBy        string     `gorm:"column:update_by" json:"updateBy"`
}

func (ComponentBlack) TableName() string {
	return "component_black"
}
