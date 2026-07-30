package models

import "time"

// LicenseGroup 许可证分组表
type LicenseGroup struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:_name" json:"name"`
	RiskWeight int       `gorm:"column:RISKWEIGHT" json:"riskWeight"`
	UUID       string    `gorm:"column:UUID" json:"uuid"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	CreateBy   string    `gorm:"column:create_by" json:"createBy"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
	UpdateBy   string    `gorm:"column:update_by" json:"updateBy"`
}

func (LicenseGroup) TableName() string {
	return "licensegroup"
}
