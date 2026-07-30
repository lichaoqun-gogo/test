package models

import "time"

// Project 开源基础软件表（同时存储开源工具，通过 category 字段区分）
type Project struct {
	ID                    int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Active                *int       `gorm:"column:ACTIVE" json:"ACTIVE"`
	Author                string     `gorm:"column:AUTHOR" json:"AUTHOR"`
	Publisher             string     `gorm:"column:PUBLISHER" json:"PUBLISHER"`
	Group                 string     `gorm:"column:_group" json:"_group"`
	Name                  string     `gorm:"column:_name" json:"_name"`
	Version               string     `gorm:"column:_version" json:"_version"`
	Description           string     `gorm:"column:DESCRIPTION" json:"DESCRIPTION"`
	Classifier            string     `gorm:"column:CLASSIFIER" json:"CLASSIFIER"`
	License               string     `gorm:"column:license" json:"license"`
	PURL                  string     `gorm:"column:PURL" json:"PURL"`
	CPE                   string     `gorm:"column:CPE" json:"CPE"`
	SWIDTagID             string     `gorm:"column:SWIDTAGID" json:"SWIDTAGID"`
	DirectDependencies    string     `gorm:"column:DIRECT_DEPENDENCIES;type:mediumtext" json:"DIRECT_DEPENDENCIES"`
	LastBomImported       *time.Time `gorm:"column:LAST_BOM_IMPORTED" json:"LAST_BOM_IMPORTED"`
	LastBomImportedFormat string     `gorm:"column:LAST_BOM_IMPORTED_FORMAT" json:"LAST_BOM_IMPORTED_FORMAT"`
	LastRiskscore         float64    `gorm:"column:LAST_RISKSCORE" json:"LAST_RISKSCORE"`
	UUID                  string     `gorm:"column:UUID" json:"UUID"`
	Type                  string     `gorm:"column:type" json:"type"`
	ApplicableLevel       string     `gorm:"column:applicable_level" json:"applicable_level"`
	MaintDept             string     `gorm:"column:maint_dept" json:"maint_dept"`
	UseDept               string     `gorm:"column:use_dept" json:"use_dept"`
	Country               string     `gorm:"column:country" json:"country"`
	Category              string     `gorm:"column:category" json:"category"` // 类别: software-开源基础软件, tool-开源工具
	CreateTime            time.Time  `gorm:"column:create_time" json:"create_time"`
	CreateBy              string     `gorm:"column:create_by" json:"create_by"`
	UpdateTime            time.Time  `gorm:"column:update_time" json:"update_time"`
	UpdateBy              string     `gorm:"column:update_by" json:"update_by"`
}

func (Project) TableName() string {
	return "project"
}
