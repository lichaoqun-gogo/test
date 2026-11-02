package models

import "time"

// License 许可证表
type License struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Comment        string     `gorm:"column:_COMMENT;type:mediumtext" json:"comment"`
	IsDeprecated   *int       `gorm:"column:IS_DEPRECATED" json:"isDeprecated"`
	FSFLibre       *int       `gorm:"column:FSFLIBRE" json:"fsfLibre"`
	Header         string     `gorm:"column:HEADER;type:mediumtext" json:"header"`
	LicenseCode    string     `gorm:"column:LICENSE" json:"licenseCode"`
	Name           string     `gorm:"column:_name" json:"name"`
	IsOSIApproved  int        `gorm:"column:IS_OSIAPPROVED" json:"isOsiApproved"`
	SeeAlso        []byte     `gorm:"column:SEEALSO;type:mediumblob" json:"seeAlso"`
	Template       string     `gorm:"column:TEMPLATE;type:mediumtext" json:"template"`
	Text           string     `gorm:"column:_text;type:mediumtext" json:"text"`
	UUID           string     `gorm:"column:UUID" json:"uuid"`
	CreateTime     time.Time  `gorm:"column:create_time" json:"createTime"`
	CreateBy       string     `gorm:"column:create_by" json:"createBy"`
	UpdateTime     time.Time  `gorm:"column:update_time" json:"updateTime"`
	UpdateBy       string     `gorm:"column:update_by" json:"updateBy"`
}

func (License) TableName() string {
	return "license"
}
