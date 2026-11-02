package models

import "time"

// DefectStatusHistory 缺陷/处置工单状态变更历史
// 覆盖父缺陷与子处置单的状态流转，支持审计与回退追溯
// 按缺陷编号（defect_no）关联缺陷实体
type DefectStatusHistory struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DefectNo     string    `gorm:"column:defect_no;size:100;not null;index:idx_defect_no" json:"defectNo"`
	DisposalID   int64     `gorm:"column:disposal_id;default:0" json:"disposalId"` // 0 表示父缺陷状态变更
	FromStatus   string    `gorm:"column:from_status;size:20" json:"fromStatus"`
	ToStatus     string    `gorm:"column:to_status;size:20" json:"toStatus"`
	Action       string    `gorm:"column:action;size:30" json:"action"` // analyze-研判, review-评审, dispatch-派发, assign-二次分配, complete-完成, rollback-回退, ignore-忽略
	Operator     string    `gorm:"column:operator;size:64" json:"operator"`
	OperatorName string    `gorm:"column:operator_name;size:64" json:"operatorName"`
	Remark       string    `gorm:"column:remark;size:512" json:"remark"`
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`
}

func (DefectStatusHistory) TableName() string {
	return "defect_status_history"
}
