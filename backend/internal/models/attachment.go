package models

import "time"

// Attachment 附件表
type Attachment struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	FileName    string    `gorm:"column:file_name" json:"fileName"`
	FileSize    *int64    `gorm:"column:file_size" json:"fileSize"`
	FileType    string    `gorm:"column:file_type" json:"fileType"`
	MinioPath   string    `gorm:"column:minio_path" json:"minioPath"`
	MinioBucket string    `gorm:"column:minio_bucket" json:"minioBucket"`
	RelatedType string    `gorm:"column:related_type" json:"relatedType"`
	RelatedID   *int64    `gorm:"column:related_id" json:"relatedId"`
	Uploader    string    `gorm:"column:uploader" json:"uploader"`
	Description string    `gorm:"column:description" json:"description"`
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (Attachment) TableName() string {
	return "attachment"
}
