package models

import "time"

// TaskUpgrade 升级申请扩展表
type TaskUpgrade struct {
	ID                int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID            int64      `gorm:"column:task_id" json:"taskId"`
	ComponentID       *int64     `gorm:"column:component_id" json:"componentId"`
	ComponentName     string     `gorm:"column:component_name" json:"componentName"`
	CurrentVersion    string     `gorm:"column:current_version" json:"currentVersion"`
	TargetVersion     string     `gorm:"column:target_version" json:"targetVersion"`
	UpgradeReason     string     `gorm:"column:upgrade_reason;type:mediumtext" json:"upgradeReason"`
	ImpactAnalysis    string     `gorm:"column:impact_analysis;type:mediumtext" json:"impactAnalysis"`
	TestReportLink    string     `gorm:"column:test_report_link" json:"testReportLink"`
	BreakingChanges   string     `gorm:"column:breaking_changes;type:mediumtext" json:"breakingChanges"`
	SecurityFixes     string     `gorm:"column:security_fixes;type:mediumtext" json:"securityFixes"`
	CompatibilityRisk string     `gorm:"column:compatibility_risk" json:"compatibilityRisk"`
	RollbackPlan      string     `gorm:"column:rollback_plan;type:mediumtext" json:"rollbackPlan"`
	PreEvalCacheID    *int64     `gorm:"column:pre_eval_cache_id" json:"preEvalCacheId"`
	PreEvalConclusion string     `gorm:"column:pre_eval_conclusion" json:"preEvalConclusion"`
	CreateTime        time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime        time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (TaskUpgrade) TableName() string {
	return "task_upgrade"
}
