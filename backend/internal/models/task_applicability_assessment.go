package models

import "time"

// TaskApplicabilityAssessment 适用性评估表
type TaskApplicabilityAssessment struct {
	ID                     int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID                 int64     `gorm:"column:task_id" json:"taskId"`
	FunctionalityStatus    string    `gorm:"column:functionality_status" json:"functionalityStatus"`
	FunctionalityDeviation string    `gorm:"column:functionality_deviation;type:mediumtext" json:"functionalityDeviation"`
	SecurityStatus         string    `gorm:"column:security_status" json:"securityStatus"`
	SecurityDeviation      string    `gorm:"column:security_deviation;type:mediumtext" json:"securityDeviation"`
	ReliabilityStatus      string    `gorm:"column:reliability_status" json:"reliabilityStatus"`
	ReliabilityDeviation   string    `gorm:"column:reliability_deviation;type:mediumtext" json:"reliabilityDeviation"`
	PerformanceStatus      string    `gorm:"column:performance_status" json:"performanceStatus"`
	PerformanceDeviation   string    `gorm:"column:performance_deviation;type:mediumtext" json:"performanceDeviation"`
	CompatibilityStatus    string    `gorm:"column:compatibility_status" json:"compatibilityStatus"`
	CompatibilityDeviation string    `gorm:"column:compatibility_deviation;type:mediumtext" json:"compatibilityDeviation"`
	MaintainabilityStatus  string    `gorm:"column:maintainability_status" json:"maintainabilityStatus"`
	MaintainabilityDeviation string  `gorm:"column:maintainability_deviation;type:mediumtext" json:"maintainabilityDeviation"`
	ExtensibilityStatus    string    `gorm:"column:extensibility_status" json:"extensibilityStatus"`
	ExtensibilityDeviation string    `gorm:"column:extensibility_deviation;type:mediumtext" json:"extensibilityDeviation"`
	UsabilityStatus        string    `gorm:"column:usability_status" json:"usabilityStatus"`
	UsabilityDeviation     string    `gorm:"column:usability_deviation;type:mediumtext" json:"usabilityDeviation"`
	CreateTime             time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime             time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (TaskApplicabilityAssessment) TableName() string {
	return "task_applicability_assessment"
}
