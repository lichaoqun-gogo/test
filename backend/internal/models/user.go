package models

import "time"

// User 用户模型（CAS 登录后自动创建/更新）
type User struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Username         string    `gorm:"uniqueIndex;size:64;not null" json:"username"`          // CAS 登录账号
	Nickname         string    `gorm:"size:64" json:"nickname"`                                // 显示名（CAS displayName）
	Email            string    `gorm:"size:128" json:"email"`                                  // 邮箱（CAS mail）
	Phone            string    `gorm:"size:32" json:"phone"`                                   // 电话（CAS phone）
	Department       string    `gorm:"size:128" json:"department"`                             // 部门（CAS department）
	Team             string    `gorm:"size:128" json:"team"`                                   // 团队
	Role             string    `gorm:"size:32;default:operator" json:"role"`                   // 角色：operator/specialist/admin
	CanToggleAdminView bool    `gorm:"default:false" json:"canToggleAdminView"`                // 是否允许切换全局管理视图
	Status           int       `gorm:"default:1" json:"status"`                                // 1=启用 0=禁用
	CreateAt         time.Time `gorm:"autoCreateTime" json:"create_at"`
	UpdateAt         time.Time `gorm:"autoUpdateTime" json:"update_at"`
}

func (User) TableName() string {
	return "users"
}
