package models

import "time"

// DefectLedger 缺陷台账表
type DefectLedger struct {
	ID                  int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DefectNo            string     `gorm:"column:defect_no" json:"defectNo"`
	Source              string     `gorm:"column:source" json:"source"`
	Title               string     `gorm:"column:title" json:"title"`
	RiskType            string     `gorm:"column:risk_type" json:"riskType"`
	RiskLevel           string     `gorm:"column:risk_level" json:"riskLevel"`
	AffectedComponent   string     `gorm:"column:affected_component" json:"affectedComponent"`
	AffectedComponentType string   `gorm:"column:affected_component_type" json:"affectedComponentType"`
	AffectedVersion     string     `gorm:"column:affected_version" json:"affectedVersion"`
	CVEID               string     `gorm:"column:cve_id" json:"cveId"`
	Description         string     `gorm:"column:description;type:mediumtext" json:"description"`
	DiscoveredDate      time.Time  `gorm:"column:discovered_date" json:"discoveredDate"`
	Status              string     `gorm:"column:status" json:"status"`
	AnalyzeResult       string     `gorm:"column:analyze_result" json:"analyzeResult"`
	ReviewMeetingNo     string     `gorm:"column:review_meeting_no" json:"reviewMeetingNo"`
	ReviewConclusion    string     `gorm:"column:review_conclusion" json:"reviewConclusion"`
	DispatchType        string     `gorm:"column:dispatch_type" json:"dispatchType"`
	DispatchTarget      string     `gorm:"column:dispatch_target" json:"dispatchTarget"`
	DispatchDept        string     `gorm:"column:dispatch_dept" json:"dispatchDept"`
	DispatchTime        *time.Time `gorm:"column:dispatch_time" json:"dispatchTime"`
	Dispatcher          string     `gorm:"column:dispatcher" json:"dispatcher"`
	RelatedTaskID       *int64     `gorm:"column:related_task_id" json:"relatedTaskId"`
	AnalyzeBy           string     `gorm:"column:analyze_by" json:"analyzeBy"`
	AnalyzeTime         *time.Time `gorm:"column:analyze_time" json:"analyzeTime"`
	ReviewBy            string     `gorm:"column:review_by" json:"reviewBy"`
	ReviewTime          *time.Time `gorm:"column:review_time" json:"reviewTime"`
	CloseTime           *time.Time `gorm:"column:close_time" json:"closeTime"`
	CreateTime          time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime          time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (DefectLedger) TableName() string {
	return "defect_ledger"
}
