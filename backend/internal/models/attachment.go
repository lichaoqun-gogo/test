package models

import "time"

// AttachmentSource 附件来源（兼容生命周期治理流程：不同流程用不同 source，复用同一张附件表）
type AttachmentSource string

const (
	// AttachmentSourceBOMUpload BOM 上传
	AttachmentSourceBOMUpload AttachmentSource = "BOM_UPLOAD"
	// AttachmentSourceVuldbPlan 漏洞处理计划
	AttachmentSourceVuldbPlan AttachmentSource = "VULDB_PLAN"
	// AttachmentSourceVuldbResult 漏洞处理结果
	AttachmentSourceVuldbResult AttachmentSource = "VULDB_RESULT"
	// AttachmentSourceFirstVerify 项目引入初审
	AttachmentSourceFirstVerify AttachmentSource = "FIRST_VERIFY"
	// AttachmentSourceFinalVerify 项目引入终审
	AttachmentSourceFinalVerify AttachmentSource = "FINAL_VERIFY"
)

// Attachment 附件管理表（统一文件管理）
// 文件统一走本表 + MinIO，业务表（bom/升级/退出等）通过本表 ID 引用文件。
// 唯一键 (source, filename, store_url) 防重复上传：同来源同目录同文件命中则更新。
// MinIO 目录规则：{source}/{业务对象目录}/{文件名}（如 bom/项目组织/项目名称/项目版本/xxx.xml）
type Attachment struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Source      string    `gorm:"column:source;size:20;not null" json:"source"`           // 附件来源（AttachmentSource）
	FileName    string    `gorm:"column:filename;size:64;not null" json:"filename"`       // 文件名
	MD5         string    `gorm:"column:md5;size:64;not null" json:"md5"`                 // 文件 md5 值
	Description string    `gorm:"column:description;type:mediumtext" json:"description"`  // 描述和备注
	StoreURL    string    `gorm:"column:store_url;size:150;not null" json:"storeUrl"`     // 对象存储中的相对地址
	CreateBy    string    `gorm:"column:create_by;size:20" json:"createBy"`               // 创建者
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`                   // 创建时间
	UpdateBy    string    `gorm:"column:update_by;size:20" json:"updateBy"`               // 更新者
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`                   // 更新时间
}

func (Attachment) TableName() string {
	return "attachment"
}
