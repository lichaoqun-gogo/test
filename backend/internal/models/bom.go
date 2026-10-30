package models

import "time"

// Bom BOM 清单表（CycloneDX 通用格式）
// 只存 BOM 元数据，文件本体走统一附件表 attachment，store_url 引用 attachment.id。
// 一个项目一条 BOM 记录（按 project_id 判新建/更新）。
type Bom struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BomFormat    string     `gorm:"column:bom_format;size:20" json:"bomFormat"`        // sbom 工具类型
	BomVersion   *int       `gorm:"column:bom_version" json:"bomVersion"`              // bom 生成版本
	SpecVersion  string     `gorm:"column:spec_version;size:20" json:"specVersion"`    // sbom 工具版本
	BomCreate    *time.Time `gorm:"column:bom_create" json:"bomCreate"`                // bom 创建时间
	SerialNumber string     `gorm:"column:serial_number;size:36" json:"serialNumber"`  // sbom 文件 uuid
	StoreURL     *int64     `gorm:"column:store_url" json:"storeUrl"`                  // sbom 文件|附件管理 attachment.id
	ProjectID    int64      `gorm:"column:project_id;not null;index" json:"projectId"` // 关联项目
	CreateTime   time.Time  `gorm:"column:create_time" json:"createTime"`
	CreateBy     string     `gorm:"column:create_by;size:20" json:"createBy"`
	UpdateTime   time.Time  `gorm:"column:update_time" json:"updateTime"`
	UpdateBy     string     `gorm:"column:update_by;size:20" json:"updateBy"`
}

func (Bom) TableName() string {
	return "bom"
}
