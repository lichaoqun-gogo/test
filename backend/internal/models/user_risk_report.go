package models

import "time"

// UserRiskReport 用户风险上报表
type UserRiskReport struct {
	ID                  int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ReportNo            string     `gorm:"column:report_no" json:"reportNo"`
	Reporter            string     `gorm:"column:reporter" json:"reporter"`
	ReporterDept        string     `gorm:"column:reporter_dept" json:"reporterDept"`
	Title               string     `gorm:"column:title" json:"title"`
	RiskType            string     `gorm:"column:risk_type" json:"riskType"`
	RiskLevel           string     `gorm:"column:risk_level" json:"riskLevel"`
	AffectedComponent   string     `gorm:"column:affected_component" json:"affectedComponent"`
	AffectedComponentType string   `gorm:"column:affected_component_type" json:"affectedComponentType"`
	AffectedVersion     string     `gorm:"column:affected_version" json:"affectedVersion"`
	Description         string     `gorm:"column:description;type:mediumtext" json:"description"`
	AttachmentURL       string     `gorm:"column:attachment_url" json:"attachmentUrl"`
	RelatedDefectID     *int64     `gorm:"column:related_defect_id" json:"relatedDefectId"`
	Status              string     `gorm:"column:status" json:"status"`
	RejectReason        string     `gorm:"column:reject_reason" json:"rejectReason"`
	ProcessedBy         string     `gorm:"column:processed_by" json:"processedBy"`
	ProcessedTime       *time.Time `gorm:"column:processed_time" json:"processedTime"`
	CreateTime          time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime          time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (UserRiskReport) TableName() string {
	return "user_risk_report"
}
