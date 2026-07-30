package models

import "time"

// DefectTaskMapping 缺陷-任务映射表
type DefectTaskMapping struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DefectID    int64     `gorm:"column:defect_id" json:"defectId"`
	DefectNo    string    `gorm:"column:defect_no" json:"defectNo"`
	TaskID      int64     `gorm:"column:task_id" json:"taskId"`
	TaskNo      string    `gorm:"column:task_no" json:"taskNo"`
	MappingType string    `gorm:"column:mapping_type" json:"mappingType"`
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`
}

func (DefectTaskMapping) TableName() string {
	return "defect_task_mapping"
}
