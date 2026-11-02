package models

import "time"

// DefectFlow 缺陷处置流程明细表（缺陷级状态机的明细字段，按缺陷编号关联实体）
// 内部处置状态（pending_analysis/pending_review/assigned/ignored/closed）以实体表 internal_status 为准，
// 本表只存流程明细：研判/评审/派发信息、截止时间、闭环时间等
//   - 外部缺陷：defect_no = defect.defect_id（如 redis#13401）
//   - 内部上报：defect_no = defect_ledger.defect_no（如 DF-xxx）
type DefectFlow struct {
	ID               int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DefectNo         string     `gorm:"column:defect_no;size:100;not null;uniqueIndex:uk_defect_no" json:"defectNo"`
	FlowType         string     `gorm:"column:flow_type;size:20;default:defect" json:"flowType"` // 流程类型: defect-缺陷, vuln-漏洞
	MaintDept        string     `gorm:"column:maint_dept;size:128" json:"maintDept"`             // 运维部门（权限过滤与裂变依据）
	Deadline         *time.Time `gorm:"column:deadline" json:"deadline"`                         // 上会/研判确定的处置截止时间
	AnalyzeResult    string     `gorm:"column:analyze_result;size:500" json:"analyzeResult"`
	ReviewMeetingNo  string     `gorm:"column:review_meeting_no;size:50" json:"reviewMeetingNo"`
	ReviewConclusion string     `gorm:"column:review_conclusion;size:20" json:"reviewConclusion"`
	DispatchType     string     `gorm:"column:dispatch_type;size:20" json:"dispatchType"`
	DispatchTarget   string     `gorm:"column:dispatch_target;size:50" json:"dispatchTarget"`
	DispatchDept     string     `gorm:"column:dispatch_dept;size:100" json:"dispatchDept"`
	DispatchTime     *time.Time `gorm:"column:dispatch_time" json:"dispatchTime"`
	Dispatcher       string     `gorm:"column:dispatcher;size:20" json:"dispatcher"`
	AnalyzeBy        string     `gorm:"column:analyze_by;size:50" json:"analyzeBy"`
	AnalyzeTime      *time.Time `gorm:"column:analyze_time" json:"analyzeTime"`
	ReviewBy         string     `gorm:"column:review_by;size:50" json:"reviewBy"`
	ReviewTime       *time.Time `gorm:"column:review_time" json:"reviewTime"`
	CloseTime        *time.Time `gorm:"column:close_time" json:"closeTime"`
	CreateTime       time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime       time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (DefectFlow) TableName() string {
	return "defect_flow"
}
