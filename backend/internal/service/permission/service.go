package permission

// ==================== 权限点与角色编码定义 ====================
// 说明：平台权限体系已改为「统一权限系统下发、后端只转发」模型（见 remote.go），
// 此处仅保留权限点/角色编码常量，供业务代码判断权限点时引用，编码与权限系统保持一致。

// 权限点编码（对接外部权限系统时沿用同一套编码）
const (
	// PermGlobalGovernanceView 全局治理视图：控制「全局治理/个人工作台」切换入口及全局治理菜单
	PermGlobalGovernanceView = "global_governance_view"
)

// 角色编码（在基础角色 operator/specialist/admin 之上，表达授权来源身份）
const (
	RoleMgmtSpecialist = "mgmt_specialist" // 管理部门开源治理专员
)
