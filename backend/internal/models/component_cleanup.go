package models

import "time"

// CommitIDHistory 历史 commit 记录（组件快照清理归档时记录，commit_id 唯一键防重复）
type CommitIDHistory struct {
	ID                int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CommitID          string    `gorm:"column:commit_id;size:64;not null;uniqueIndex" json:"commitId"`
	CreateTimeHistory *time.Time `gorm:"column:create_time_history" json:"createTimeHistory"` // 历史创建时间
	CreateTime        time.Time `gorm:"column:create_time" json:"createTime"`                 // 后台归档时间
	CreateBy          string    `gorm:"column:create_by;size:20" json:"createBy"`
	UpdateTime        time.Time `gorm:"column:update_time" json:"updateTime"`
	UpdateBy          string    `gorm:"column:update_by;size:20" json:"updateBy"`
}

func (CommitIDHistory) TableName() string {
	return "commit_id_history"
}

// ProjectInnerProInfo 自研项目线上信息（is_deleted=0 在用的 commit_id 为线上版本，不清理）
type ProjectInnerProInfo struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DailyPeriod     string     `gorm:"column:daily_period;size:20;not null" json:"dailyPeriod"` // 信息同步日期，如 2020-01-10
	ProCode         string     `gorm:"column:pro_code;size:60;not null;uniqueIndex" json:"proCode"`
	CommitID        string     `gorm:"column:commit_id;size:64;not null;index" json:"commitId"` // 线上版本号
	IsDeleted       int        `gorm:"column:is_deleted;not null;default:1" json:"isDeleted"`   // 0-在用 1-已删除
	LastReleaseTime *time.Time `gorm:"column:last_release_time" json:"lastReleaseTime"`
	CreateTime      time.Time  `gorm:"column:create_time" json:"createTime"`
	CreateBy        string     `gorm:"column:create_by;size:20" json:"createBy"`
	UpdateTime      time.Time  `gorm:"column:update_time" json:"updateTime"`
	UpdateBy        string     `gorm:"column:update_by;size:20" json:"updateBy"`
}

func (ProjectInnerProInfo) TableName() string {
	return "project_inner_pro_info"
}
