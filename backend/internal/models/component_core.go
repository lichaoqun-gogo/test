package models

import "time"

// ComponentCore 名单配置表
type ComponentCore struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Format          int        `gorm:"column:_format" json:"_format"`               // 仓库类型: 1-Maven, 2-Pypi, 3-Npm, 4-Yum, 5-GO_MODULES
	Type            string     `gorm:"column:type" json:"type"`                     // 分类: 安全算法类、服务治理类等
	Name            string     `gorm:"column:_name" json:"_name"`                   // 名称
	Group           string     `gorm:"column:_group" json:"_group"`                 // 组织
	VersionExp      string     `gorm:"column:version_exp" json:"version_exp"`       // 匹配表达式
	ApplicableLevel string     `gorm:"column:applicable_level" json:"applicable_level"` // 适用等级: A/B/C/D
	SupportTeam     string     `gorm:"column:support_team" json:"support_team"`     // 技术支持团队
	Reason          string     `gorm:"column:reason" json:"reason"`                 // 名单理由
	Advise          string     `gorm:"column:advise" json:"advise"`                 // 建议
	Enable          int        `gorm:"column:_enable" json:"_enable"`               // 是否启用
	AlertType       *int       `gorm:"column:alert_type" json:"alert_type"`         // 告警类型: 1-阻断, 2-警告
	CreateTime      time.Time  `gorm:"column:create_time" json:"create_time"`
	CreateBy        string     `gorm:"column:create_by" json:"create_by"`
	UpdateTime      time.Time  `gorm:"column:update_time" json:"update_time"`
	UpdateBy        string     `gorm:"column:update_by" json:"update_by"`
	LevelTime       time.Time  `gorm:"column:level_time" json:"level_time"`
	IsBoundary      *int       `gorm:"column:is_boundary" json:"is_boundary"`
}

func (ComponentCore) TableName() string {
	return "component_core"
}
