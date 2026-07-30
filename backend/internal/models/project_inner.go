package models

import "time"

// ProjectInner 自研项目表
type ProjectInner struct {
	ID                   int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProCode              string     `gorm:"column:pro_code" json:"pro_code"`
	Active               *int       `gorm:"column:active" json:"active"`
	Type                 string     `gorm:"column:type" json:"type"`
	IsBoundary           *int       `gorm:"column:is_boundary" json:"is_boundary"`
	GroupName            string     `gorm:"column:group_name" json:"group_name"`
	DeptName             string     `gorm:"column:dept_name" json:"dept_name"`
	ModuleCode           string     `gorm:"column:module_code" json:"module_code"`
	Branch               string     `gorm:"column:branch" json:"branch"`
	CurrentID            string     `gorm:"column:current_ID" json:"current_ID"`
	CurrentCommitID      string     `gorm:"column:current_commit_id" json:"current_commit_id"`
	ProductionCommitID   string     `gorm:"column:production_commit_id" json:"production_commit_id"`
	CreateBy             string     `gorm:"column:create_by" json:"create_by"`
	UpdateBy             string     `gorm:"column:update_by" json:"update_by"`
	CreateTime           time.Time  `gorm:"column:create_time" json:"create_time"`
	UpdateTime           time.Time  `gorm:"column:update_time" json:"update_time"`
}

func (ProjectInner) TableName() string {
	return "project_inner"
}
