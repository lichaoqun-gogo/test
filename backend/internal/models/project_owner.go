package models

import "time"

// ProjectOwner 项目-负责人关联表
// 通过 project_type 区分开源基础软件(software)和自研项目(inner)
// 一个项目可对应多个负责人（支持多对多关系）
type ProjectOwner struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProjectType   string    `gorm:"column:project_type;size:32;not null;index:idx_type_project" json:"project_type"` // software | inner
	ProjectID     int64     `gorm:"column:project_id;not null;index:idx_type_project" json:"project_id"`           // 关联 project.id 或 project_inner.id
	OwnerUsername string    `gorm:"column:owner_username;size:64;not null;index:idx_owner" json:"owner_username"` // 账户名(关联 users.username)
	DeptName      string    `gorm:"column:dept_name;size:128" json:"dept_name"`                                  // 所属部门(冗余)
	Team          string    `gorm:"column:team;size:128" json:"team"`                                          // 所属团队(冗余)
	AssignedBy    string    `gorm:"column:assigned_by;size:64" json:"assigned_by"`                              // 分配人(治理专员账号)
	AssignedAt    time.Time `gorm:"column:assigned_at" json:"assigned_at"`                                       // 分配时间
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ProjectOwner) TableName() string {
	return "project_owner"
}
