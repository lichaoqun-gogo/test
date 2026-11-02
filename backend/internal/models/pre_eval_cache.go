package models

import "time"

// PreEvalCache 预评估缓存表
type PreEvalCache struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ComponentID    *int64     `gorm:"column:component_id" json:"componentId"`
	ComponentName  string     `gorm:"column:component_name" json:"componentName"`
	ComponentVersion string   `gorm:"column:component_version" json:"componentVersion"`
	EvalResultJSON string     `gorm:"column:eval_result_json;type:mediumtext" json:"evalResultJson"`
	EvalTime       time.Time  `gorm:"column:eval_time" json:"evalTime"`
	ExpireTime     time.Time  `gorm:"column:expire_time" json:"expireTime"`
	IsValid        int        `gorm:"column:is_valid" json:"isValid"`
	CreateTime     time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime     time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (PreEvalCache) TableName() string {
	return "pre_eval_cache"
}
