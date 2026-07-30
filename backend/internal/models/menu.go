package models

import "time"

// Menu 菜单模型
type Menu struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ParentID  uint      `gorm:"default:0;index" json:"parent_id"`   // 父菜单ID，0=顶级
	Name      string    `gorm:"size:50;not null" json:"name"`       // 菜单名称
	Path      string    `gorm:"size:200" json:"path"`               // 路由路径（完整路径如 /admin/asset/software-library）
	Component string    `gorm:"size:200" json:"component"`          // 组件路径（如 views/admin/asset/SoftwareLibrary）
	Icon      string    `gorm:"size:50" json:"icon"`                // 菜单图标
	Sort      int       `gorm:"default:0" json:"sort"`              // 排序
	MenuType  int       `gorm:"default:1" json:"menu_type"`         // 1=目录 2=菜单 3=按钮（预留）
	Hidden    bool      `gorm:"default:false" json:"hidden"`        // 是否隐藏
	Group     string    `gorm:"size:50" json:"group"`               // 分组标识（如 asset/risk/governance/decision）
	CreateAt  time.Time `gorm:"autoCreateTime" json:"create_at"`
}

func (Menu) TableName() string {
	return "menus"
}

// RoleMenu 角色-菜单关联
type RoleMenu struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Role   string `gorm:"size:32;index" json:"role"`  // 角色：operator/specialist/admin
	MenuID uint   `gorm:"index" json:"menu_id"`       // 菜单ID
}

func (RoleMenu) TableName() string {
	return "role_menus"
}
