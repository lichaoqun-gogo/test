package models

import "time"

// DefectLedger 内部缺陷实体表（风险上报那套字段 + 内部处置状态）
// 用户上报风险时创建，编号 DF-xxx；字段与风险上报表单对齐（编号/标题/描述/等级选填），
// 处置流程明细（研判/审核/派发/截止）见 DefectFlow（按 defect_no 关联）
type DefectLedger struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DefectNo       string    `gorm:"column:defect_no;size:100;not null" json:"defectNo"`
	Source         string    `gorm:"column:source;size:30;not null" json:"source"`
	Title          string    `gorm:"column:title;size:255;not null" json:"title"`
	RiskLevel      string    `gorm:"column:risk_level;size:20;not null;default:未知" json:"riskLevel"` // 风险等级（选填）
	Description    string    `gorm:"column:description;type:mediumtext" json:"description"`
	AffectedComponent string `gorm:"column:affected_component;size:128" json:"affectedComponent"` // 影响软件/组件（漏洞自动分配按组件裂变：redis/fastjson 等）
	DiscoveredDate time.Time `gorm:"column:discovered_date;not null" json:"discoveredDate"`
	InternalStatus string    `gorm:"column:internal_status;size:20;not null;default:pending_analysis;index:idx_internal_status" json:"internalStatus"` // 内部处置状态: pending_analysis-待研判, pending_review-待上会/待审核, assigned-已分配, ignored-已忽略, closed-已闭环
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime     time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (DefectLedger) TableName() string {
	return "defect_ledger"
}
