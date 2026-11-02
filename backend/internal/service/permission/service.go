package permission

import (
	"osm-mngt/internal/models"

	"gorm.io/gorm"
)

// 权限点编码（对接外部权限系统时沿用同一套编码）
const (
	// PermGlobalGovernanceView 全局治理视图：控制「全局治理/个人工作台」切换入口及全局治理菜单
	PermGlobalGovernanceView = "global_governance_view"
)

// 角色编码（在 users.role 基础角色之上，表达授权来源身份）
const (
	RoleMgmtSpecialist = "mgmt_specialist" // 管理部门开源治理专员
)

// Profile 用户权限档案：登录后由权限体系计算下发，前端据此渲染可见菜单与按钮。
// 这是与外部权限系统对接的数据契约，后续切换实现时保持结构不变。
type Profile struct {
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

// Service 权限服务接口：对接外部权限系统时替换实现（如改为 HTTP 调用），契约不变
type Service interface {
	GetProfile(username string) *Profile
	HasPermission(username string, perm string) bool
}

// LocalService 本地适配实现（外部权限系统对接前的临时方案）。
// 授权规则等价于外部权限系统中的配置：
//  1. users.role = admin → 授予 global_governance_view
//  2. dept_specialist_mapping 中为有效的 mgmt_dept_specialist → 追加 mgmt_specialist 角色并授予 global_governance_view
//  3. users.can_toggle_admin_view = 1 → 授予 global_governance_view（模拟外部权限系统单独授权，可直接改库开关测试）
type LocalService struct {
	db *gorm.DB
}

func NewLocalService(db *gorm.DB) *LocalService {
	return &LocalService{db: db}
}

// GetProfile 计算用户权限档案
func (s *LocalService) GetProfile(username string) *Profile {
	profile := &Profile{
		Username:    username,
		Roles:       []string{},
		Permissions: []string{},
	}
	if username == "" {
		return profile
	}

	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return profile
	}
	profile.Roles = append(profile.Roles, user.Role)

	hasGlobalView := user.Role == "admin" || user.CanToggleAdminView

	// 管理部门专员 → mgmt_specialist 角色 + 全局治理视图权限
	var count int64
	s.db.Table("dept_specialist_mapping").
		Where("specialist_id = ? AND specialist_role = 'mgmt_dept_specialist' AND is_active = 1", username).
		Count(&count)
	if count > 0 {
		profile.Roles = append(profile.Roles, RoleMgmtSpecialist)
		hasGlobalView = true
	}

	if hasGlobalView {
		profile.Permissions = append(profile.Permissions, PermGlobalGovernanceView)
	}
	return profile
}

// HasPermission 判断用户是否具备某权限点
func (s *LocalService) HasPermission(username string, perm string) bool {
	profile := s.GetProfile(username)
	for _, p := range profile.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}
