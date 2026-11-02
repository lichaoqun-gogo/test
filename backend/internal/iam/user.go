package iam

// IAMUser IAM 登录用户模型：用户登录时一次性获取的自身全部权限资源信息。
// 对应 IAM OpenAPI 的登录用户权限查询接口响应，作为登录后初始化用户身份/角色/菜单/按钮的承接结构。
// 与平台登录初始化链路的映射（由适配层完成）：
//   - Account → 用户身份（对应 GetUserInfo / users 表）
//   - Roles   → 角色编码（对应 users.role 及附加角色如 mgmt_specialist）
//   - Menus   → 菜单树（对应 GetUserMenus，注意 IAM 菜单缺 component，需本地补齐）
//   - Orgs    → 组织信息
//   - Buttons → 按钮级权限（对应前端操作按钮显隐）
type IAMUser struct {
	Account *Account      `json:"account"` // 帐户信息
	Roles   []string      `json:"roles"`   // 角色编码列表
	Menus   []*Menu       `json:"menus"`   // 菜单（树形）
	Orgs    []*Org        `json:"orgs"`    // 组织
	Buttons []*Permission `json:"buttons"` // 按钮权限（资源类型 BUTTON 的权限列表）
}
