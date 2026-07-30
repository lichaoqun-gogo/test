package models

import "time"

// LicenseGroupLicense 许可证分组关联表
type LicenseGroupLicense struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	LicenseGroupID int64    `gorm:"column:licensegroup_id" json:"licenseGroupId"`
	LicenseID     int64     `gorm:"column:license_id" json:"licenseId"`
	CreateTime    time.Time `gorm:"column:create_time" json:"createTime"`
	CreateBy      string    `gorm:"column:create_by" json:"createBy"`
	UpdateTime    time.Time `gorm:"column:update_time" json:"updateTime"`
	UpdateBy      string    `gorm:"column:update_by" json:"updateBy"`
}

func (LicenseGroupLicense) TableName() string {
	return "licensegroup_license"
}
