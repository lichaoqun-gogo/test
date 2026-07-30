package models

import "time"

// TaskIntroduce 引入申请扩展表
type TaskIntroduce struct {
	ID               int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID           int64      `gorm:"column:task_id" json:"taskId"`
	ComponentID      *int64     `gorm:"column:component_id" json:"componentId"`
	ComponentName    string     `gorm:"column:component_name" json:"componentName"`
	ComponentVersion string     `gorm:"column:component_version" json:"componentVersion"`
	LicenseType      string     `gorm:"column:license_type" json:"licenseType"`
	IntroduceReason  string     `gorm:"column:introduce_reason;type:mediumtext" json:"introduceReason"`
	UseScenario      string     `gorm:"column:use_scenario;type:mediumtext" json:"useScenario"`
	ProjectNames     string     `gorm:"column:project_names;type:mediumtext" json:"projectNames"`
	PreEvalCacheID   *int64     `gorm:"column:pre_eval_cache_id" json:"preEvalCacheId"`
	PreEvalConclusion string    `gorm:"column:pre_eval_conclusion" json:"preEvalConclusion"`
	PreEvalReport    string     `gorm:"column:pre_eval_report;type:mediumtext" json:"preEvalReport"`
	Remark           string     `gorm:"column:remark" json:"remark"`
	CreateTime       time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime       time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (TaskIntroduce) TableName() string {
	return "task_introduce"
}
