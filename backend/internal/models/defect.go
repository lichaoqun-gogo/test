package models

import "time"

// Defect 外部缺陷实体表（采集台账，GitHub/Jira 多通道 upsert 导入）
// 主键 defect_id = 软件小写#上游编号（redis#13401 / kafka#KAFKA-20449），导入幂等业务键
// 内部处置状态 internal_status 以本表为准（导入只新增不覆盖状态）
type Defect struct {
	DefectID       string     `gorm:"column:defect_id;size:100;primaryKey" json:"defectId"`
	Software       string     `gorm:"column:software;size:100;not null" json:"software"`
	IssueURL       string     `gorm:"column:issue_url;size:500;not null" json:"issueUrl"`
	Title          string     `gorm:"column:title;size:1000;not null" json:"title"`
	Body           string     `gorm:"column:body;type:mediumtext" json:"body"`
	Labels         string     `gorm:"column:labels;type:json" json:"labels"`         // 原始标签数组 JSON 文本
	Reporter       string     `gorm:"column:reporter;size:200" json:"reporter"`
	CreatedAt      *time.Time `gorm:"column:created_at" json:"createdAt"`            // 缺陷创建时间（UTC）
	UpdatedAt      *time.Time `gorm:"column:updated_at" json:"updatedAt"`            // 缺陷更新时间（UTC），增量游标依据
	CommentsCount  int        `gorm:"column:comments_count;not null;default:0" json:"commentsCount"`
	Source         string     `gorm:"column:source;size:20;not null;default:github-issue" json:"source"`
	Severity       string     `gorm:"column:severity;size:10;not null;default:未知" json:"severity"`
	IssueState     string     `gorm:"column:issue_state;size:10;not null;default:open" json:"issueState"`
	StateReason    string     `gorm:"column:state_reason;size:20" json:"stateReason"`
	Reactions      string     `gorm:"column:reactions;type:json" json:"reactions"`   // 社区关注度 JSON 文本
	Milestone      string     `gorm:"column:milestone;size:100" json:"milestone"`
	WatchMatches   string     `gorm:"column:watch_matches;type:json" json:"watchMatches"` // 命中观测版本数组 JSON 文本
	SummaryZh       string     `gorm:"column:summary_zh;type:text" json:"summaryZh"`
	RepairAdvice    string     `gorm:"column:repair_advice;type:text" json:"repairAdvice"`
	TitleZh         string     `gorm:"column:title_zh;size:1000" json:"titleZh"` // 中文标题（采集工具翻译产物，详情展示）
	Org             string     `gorm:"column:org;size:128" json:"org"`                       // 组织（开源软件/组件的组织命名空间，如 com.alibaba）
	MaintDept       string     `gorm:"column:maint_dept;size:128" json:"maintDept"`             // 运维部门（首次导入者部门，重复导入不覆盖）
	InternalStatus  string     `gorm:"column:internal_status;size:20;not null;default:pending_analysis;index:idx_internal_status" json:"internalStatus"` // 内部处置状态: pending_analysis-待研判, pending_review-待上会, assigned-已分配, ignored-已忽略, closed-已闭环
	ImportedAt      time.Time  `gorm:"column:imported_at;not null" json:"importedAt"`          // 本次导入时间
}

func (Defect) TableName() string {
	return "defect"
}
