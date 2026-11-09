package handlers

// ==================== 静态权限资源数据（模拟统一权限系统） ====================
// 统一权限系统未对接（auth.permission.baseURL 为空）时，
// /api/user/info 依据登录用户角色返回以下静态权限资源，返回结构对齐权限系统
// （account/roles/menus/orgs/buttons，见 handlers/reference.go UserInfo），
// 保证前端「菜单按权限渲染、按钮按权限控制」的开发与演示不阻塞，
// 并贴近内网权限系统返回结构（减轻对接真实权限系统后的调试工作量）。
// 对接真实权限系统后，UserInfo 优先透传权限系统返回（remote.go GetUserInfo），本文件不再使用。

import (
	"osm-mngt/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// permRoleCode 基础角色（CAS/JWT claims.Role）→ 权限系统角色编码（贴近内网 COMMON/SPECIALIST/ADMIN）
// 说明：权限系统的角色编码由各业务方定义（用户约定：把管理部门专员编码叫 admin 也行），
// 关键是前端能据此划分三类角色形态；此处仅做未对接时的模拟映射。
var permRoleCode = map[string]string{
	"operator":   "COMMON",
	"specialist": "SPECIALIST",
	"admin":      "ADMIN",
}

// permMenuNode 权限系统菜单节点（对齐 /api/user/info 的 menus 数组结构）
type permMenuNode struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	AppCode   string         `json:"appCode"`
	AppModule string         `json:"appModule"`
	ParentID  uint           `json:"parentId"`
	OrderNum  int            `json:"orderNum"`
	Path      string         `json:"path"`
	Icon      string         `json:"icon"`
	Visible   bool           `json:"visible"`
	Children  []permMenuNode `json:"children"`
}

// permButton 权限系统按钮资源（对齐 /api/user/info 的 buttons 数组结构）
type permButton struct {
	ID           uint     `json:"id"`
	Resource     string   `json:"resource"`
	ResourceType string   `json:"resourceType"`
	Visible      bool     `json:"visible"`
	Actions      []string `json:"actions"`
	Conditions   []any    `json:"conditions"`
}

// permOrg 权限系统组织资源（对齐 /api/user/info 的 orgs 数组结构）
type permOrg struct {
	AppCode       string    `json:"appCode"`
	AppOrg        string    `json:"appOrg"`
	OrgCode       string    `json:"orgCode"`
	OrgName       string    `json:"orgName"`
	OrgPrincipal  string    `json:"orgPrincipal"`
	ParentOrgCode string    `json:"parentOrgCode"`
	ParentOrgName string    `json:"parentOrgName"`
	OrgType       string    `json:"orgType"`
	OrgLevel      string    `json:"orgLevel"`
	DepCode       string    `json:"depCode"`
	DepName       string    `json:"depName"`
	Children      []permOrg `json:"children"`
}

// ========== 菜单资源定义（path 必须与前端 staticRoutes / menuMap.ts 完全一致） ==========

// 个人工作台菜单（SPECIALIST 治理专员；ADMIN 的个人视图复用本套）
var specialistPermMenus = []permMenuNode{
	{ID: 1001, Name: "我的项目", AppCode: "osm", AppModule: "personal", Path: "/specialist/my-projects", Icon: "Monitor", Visible: true},
	{ID: 1002, Name: "我的任务", AppCode: "osm", AppModule: "personal", Path: "/specialist/my-task", Icon: "List", Visible: true},
	{
		ID: 1003, Name: "我的运维", AppCode: "osm", AppModule: "personal", Icon: "Setting", Visible: true, Children: []permMenuNode{
			{ID: 10031, Name: "运维台账", AppCode: "osm", AppModule: "personal", ParentID: 1003, Path: "/specialist/my-ops", Icon: "", Visible: true},
			{ID: 10032, Name: "人员管理", AppCode: "osm", AppModule: "personal", ParentID: 1003, Path: "/specialist/personnel-mgmt", Icon: "", Visible: true},
		},
	},
	{ID: 1004, Name: "我的使用", AppCode: "osm", AppModule: "personal", Path: "/specialist/my-usage", Icon: "Files", Visible: true},
	{ID: 1005, Name: "可信仓库", AppCode: "osm", AppModule: "personal", Path: "/specialist/trusted-repo", Icon: "Files", Visible: true},
}

// 全局治理菜单（ADMIN 管理部门专员的全局治理视图）
var adminPermMenus = []permMenuNode{
	{
		ID: 1101, Name: "台账管理", AppCode: "osm", AppModule: "asset", Icon: "Coin", Visible: true, Children: []permMenuNode{
			{ID: 11011, Name: "开源基础软件", AppCode: "osm", AppModule: "asset", ParentID: 1101, Path: "/admin/asset/software-library", Icon: "", Visible: true},
			{ID: 11012, Name: "开源依赖包", AppCode: "osm", AppModule: "asset", ParentID: 1101, Path: "/admin/asset/dependency-library", Icon: "", Visible: true},
			{ID: 11013, Name: "开源工具", AppCode: "osm", AppModule: "asset", ParentID: 1101, Path: "/admin/asset/tool-library", Icon: "", Visible: true},
			{ID: 11014, Name: "自研项目", AppCode: "osm", AppModule: "asset", ParentID: 1101, Path: "/admin/asset/self-developed", Icon: "", Visible: true},
			{ID: 11015, Name: "源码备份", AppCode: "osm", AppModule: "asset", ParentID: 1101, Path: "/admin/asset/source-backup", Icon: "", Visible: true},
		},
	},
	{ID: 1102, Name: "可信仓库", AppCode: "osm", AppModule: "trusted", Path: "/admin/trusted-repo", Icon: "Files", Visible: true},
	{
		ID: 1103, Name: "风险管理", AppCode: "osm", AppModule: "risk", Icon: "Warning", Visible: true, Children: []permMenuNode{
			{ID: 11031, Name: "名单管理", AppCode: "osm", AppModule: "risk", ParentID: 1103, Path: "/admin/risk/blacklist", Icon: "", Visible: true},
			{ID: 11032, Name: "CICD检查", AppCode: "osm", AppModule: "risk", ParentID: 1103, Path: "/admin/risk/cicd-check", Icon: "", Visible: true},
			{ID: 11033, Name: "许可证信息库", AppCode: "osm", AppModule: "risk", ParentID: 1103, Path: "/admin/risk/license-library", Icon: "", Visible: true},
			{ID: 11034, Name: "缺陷台账", AppCode: "osm", AppModule: "risk", ParentID: 1103, Path: "/admin/risk/defect-ledger", Icon: "", Visible: true},
		},
	},
	{
		ID: 1104, Name: "治理中心", AppCode: "osm", AppModule: "governance", Icon: "Setting", Visible: true, Children: []permMenuNode{
			{ID: 11041, Name: "任务列表", AppCode: "osm", AppModule: "governance", ParentID: 1104, Path: "/admin/governance/task-list", Icon: "", Visible: true},
		},
	},
	{
		ID: 1105, Name: "决策与洞察中心", AppCode: "osm", AppModule: "decision", Icon: "DataAnalysis", Visible: true, Children: []permMenuNode{
			{ID: 11051, Name: "缺陷修复效能", AppCode: "osm", AppModule: "decision", ParentID: 1105, Path: "/admin/decision/repair-efficiency", Icon: "", Visible: true},
			{ID: 11052, Name: "供应链安全洞察", AppCode: "osm", AppModule: "decision", ParentID: 1105, Path: "/admin/decision/core-dependencies", Icon: "", Visible: true},
		},
	},
}

// 普通用户菜单（COMMON operator）
var commonPermMenus = []permMenuNode{
	{ID: 1001, Name: "我的项目", AppCode: "osm", AppModule: "personal", Path: "/operator/my-projects", Icon: "Monitor", Visible: true},
	{ID: 1002, Name: "我的任务", AppCode: "osm", AppModule: "personal", Path: "/operator/my-task", Icon: "List", Visible: true},
	{ID: 1003, Name: "我的运维", AppCode: "osm", AppModule: "personal", Path: "/operator/my-ops", Icon: "Setting", Visible: true},
	{ID: 1004, Name: "我的使用", AppCode: "osm", AppModule: "personal", Path: "/operator/my-usage", Icon: "Files", Visible: true},
	{ID: 1005, Name: "可信仓库", AppCode: "osm", AppModule: "personal", Path: "/operator/trusted-repo", Icon: "Files", Visible: true},
}

// menusForPermRole 按权限系统角色编码返回菜单树
func menusForPermRole(permRole string) []permMenuNode {
	switch permRole {
	case "SPECIALIST":
		return specialistPermMenus
	case "ADMIN":
		// 管理部门专员：个人工作台（specialist 套）+ 全局治理两套菜单
		return append(append([]permMenuNode{}, specialistPermMenus...), adminPermMenus...)
	case "COMMON":
		return commonPermMenus
	default:
		return commonPermMenus
	}
}

// ========== 按钮资源定义（resource 与前端 v-auth 编码一致） ==========

// buttonsForPermRole 按权限系统角色编码返回按钮资源
func buttonsForPermRole(permRole string) []permButton {
	switch permRole {
	case "SPECIALIST":
		// 部门治理专员：人员管理可编辑负责人、缺陷可一键自动分配
		return []permButton{
			{ID: 2001, Resource: "software_owner_update", ResourceType: "BUTTON", Visible: true, Actions: []string{"update"}},
			{ID: 2002, Resource: "project_owner_update", ResourceType: "BUTTON", Visible: true, Actions: []string{"update"}},
			{ID: 2003, Resource: "risk_auto_assign", ResourceType: "BUTTON", Visible: true, Actions: []string{"assign"}},
		}
	case "ADMIN":
		// 管理部门专员：继承专员全部按钮 + 全局治理管理按钮
		return append(buttonsForPermRole("SPECIALIST"), []permButton{
			{ID: 2101, Resource: "asset_review", ResourceType: "BUTTON", Visible: true, Actions: []string{"review"}},
			{ID: 2102, Resource: "defect_dispose", ResourceType: "BUTTON", Visible: true, Actions: []string{"dispose"}},
		}...)
	case "COMMON":
		return []permButton{
			{ID: 2001, Resource: "project_create", ResourceType: "BUTTON", Visible: true, Actions: []string{"create"}},
		}
	default:
		return []permButton{}
	}
}

// ========== 组织资源 ==========

// deptCodeMap 部门名 → 部门编码（模拟权限系统 depCode，未对接时兜底）
var deptCodeMap = map[string]string{
	"技术部":         "tech",
	"平台研发部":      "platform-dev",
	"运维部":         "ops",
	"安全部":         "sec",
	"开源治理办公室":   "osm-gov",
}

// teamCodeMap 团队名 → 团队编码（模拟权限系统 orgCode，未对接时兜底）
var teamCodeMap = map[string]string{
	"前端组": "fe",
	"后端组": "be",
	"运营组": "op",
	"设计组": "design",
}

// orgsForDept 按部门/团队构造组织资源
func orgsForDept(dept, team string) []permOrg {
	if dept == "" {
		return []permOrg{}
	}
	return []permOrg{{
		AppCode:       "osm",
		AppOrg:        "osm-governance",
		OrgCode:       teamCodeMap[team],
		OrgName:       team,
		OrgPrincipal:  "",
		ParentOrgCode: "100",
		ParentOrgName: dept,
		OrgType:       "",
		OrgLevel:      "",
		DepCode:       deptCodeMap[dept],
		DepName:       dept,
		Children:      []permOrg{},
	}}
}

// normPermMenus 规范化菜单树：nil 子节点 → 空数组（对齐权限系统返回 [] 而非 null）
func normPermMenus(menus []permMenuNode) []permMenuNode {
	for i := range menus {
		if menus[i].Children == nil {
			menus[i].Children = []permMenuNode{}
		} else {
			normPermMenus(menus[i].Children)
		}
	}
	return menus
}

// normPermButtons 规范化按钮：nil conditions → 空数组（对齐权限系统返回 [] 而非 null）
func normPermButtons(buttons []permButton) []permButton {
	for i := range buttons {
		if buttons[i].Conditions == nil {
			buttons[i].Conditions = []any{}
		}
	}
	return buttons
}

// staticPermInfo 构造 /api/user/info 的静态权限资源兜底数据（结构对齐权限系统返回）
func staticPermInfo(username string, claims *jwt.Claims) gin.H {
	// 基础身份：优先 JWT claims，其次 staticUsers 兜底
	role := "operator"
	nickname := ""
	dept := ""
	team := ""
	email := ""
	if claims != nil {
		role = claims.Role
		nickname = claims.Nickname
		dept = claims.Department
		team = claims.Team
		email = claims.Email
	}
	if u, ok := staticUserByUsername(username); ok {
		if nickname == "" {
			nickname = u.Nickname
		}
		if dept == "" {
			dept = u.Department
		}
		if team == "" {
			team = u.Team
		}
		if email == "" {
			email = u.Email
		}
	}
	permRole := permRoleCode[role]
	if permRole == "" {
		permRole = "COMMON"
	}

	account := gin.H{
		"accountCode":    username,
		"accountName":    nickname,
		"email":          email,
		"depCode":        deptCodeMap[dept],
		"depName":        dept,
		"orgCode":        teamCodeMap[team],
		"customOrgCode":  "",
		"orgName":        team,
		"orgDisplayName": dept + "/" + team,
		"partOrgCodes":   []string{},
	}

	return gin.H{
		"account": account,
		"roles":   []string{permRole},
		"menus":   normPermMenus(menusForPermRole(permRole)),
		"orgs":    orgsForDept(dept, team),
		"buttons": normPermButtons(buttonsForPermRole(permRole)),
	}
}
