package models

import "time"

// SpecialistResourceMapping 人员-资源映射表
type SpecialistResourceMapping struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ResourceType  string    `gorm:"column:resource_type" json:"resourceType"`
	ResourceID    int64     `gorm:"column:resource_id" json:"resourceId"`
	ResourceName  string    `gorm:"column:resource_name" json:"resourceName"`
	Owner         string    `gorm:"column:owner" json:"owner"`
	OwnerDept     string    `gorm:"column:owner_dept" json:"ownerDept"`
	Maintainer    string    `gorm:"column:maintainer" json:"maintainer"`
	MaintainerDept string   `gorm:"column:maintainer_dept" json:"maintainerDept"`
	CreateTime    time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime    time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (SpecialistResourceMapping) TableName() string {
	return "specialist_resource_mapping"
}
