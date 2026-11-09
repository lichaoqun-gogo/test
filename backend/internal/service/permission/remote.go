package permission

// ==================== 统一权限系统转发客户端 ====================
// 设计约定：平台的权限类信息（用户档案/角色/权限点/菜单/组织）统一由外部「统一权限系统」下发，
// 后端只做转发、不落库（登录不再写 users 表、菜单不再存表）。
//
// 兜底原则：权限系统未对接（auth.permission.baseURL 为空）或请求失败时，
// 打印警告日志并返回错误，由 handler 层返回兜底数据——保证登录与开发调试不被阻塞。
// 对接后：只需在 config 中填写权限系统地址/接口路径，并在此文件对齐返回字段即可。
//
// 涉及的权限接口（共 5 个，均 GET）：
//   1. 用户档案  {baseURL}{profilePath}?username=xxx → {user, roles, permissions}
//   2. 菜单      {baseURL}{menusPath}?username=xxx  → {roles, menus: {<role>: [菜单树]}}
//   3. 部门-团队-成员 {baseURL}{orgDeptsPath}
//   4. 部门-所属团队   {baseURL}{orgDeptTeamsPath}
//   5. 团队人员   {baseURL}{teamUsersPath}?team=xxx&dept_name=xxx

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"osm-mngt/internal/config"
)

// ProfileUser 用户档案中的用户基本信息（对齐权限系统 user 字段）
type ProfileUser struct {
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	Department string `json:"department"`
	Team       string `json:"team"`
	Role       string `json:"role"` // 基础角色：operator / specialist / admin
}

// ProfileBundle 用户档案响应：user + 角色集 + 权限点集（对齐权限系统返回结构）
type ProfileBundle struct {
	User        ProfileUser `json:"user"`
	Roles       []string    `json:"roles"`       // 角色编码（可多个，如 admin + mgmt_specialist）
	Permissions []string    `json:"permissions"` // 权限点编码（如 global_governance_view）
}

// MenuNode 菜单树节点（与前端 src/router/index.ts 的 MenuTreeNode 结构一致）
type MenuNode struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Component string     `json:"component"` // 前端 componentMaps 的 key，如 views/operator/MyOps
	Icon      string     `json:"icon"`
	Sort      int        `json:"sort"`
	Hidden    bool       `json:"hidden"`
	Group     string     `json:"group"`
	Children  []MenuNode `json:"children"`
}

// MenusBundle 菜单响应：角色集 + 各角色菜单树（后端透传，不构建菜单）
type MenusBundle struct {
	Roles []string              `json:"roles"`
	Menus map[string][]MenuNode `json:"menus"` // key=角色编码，value=该角色菜单树
}

// ==================== 权限系统完整权限信息结构（对齐 /api/user/info 返回） ====================
// 权限系统对接后，GetUserInfo 一次透传 account/roles/menus/orgs/buttons 五个字段，
// 前端按此结构渲染菜单（menus 数组）与控制按钮（buttons 数组），结构见 handlers/reference.go UserInfo。

// UserAccount 用户档案 account（对齐权限系统返回字段）
type UserAccount struct {
	AccountCode    string   `json:"accountCode"`
	AccountName    string   `json:"accountName"`
	Email          string   `json:"email"`
	DepCode        string   `json:"depCode"`
	DepName        string   `json:"depName"`
	OrgCode        string   `json:"orgCode"`
	CustomOrgCode  string   `json:"customOrgCode"`
	OrgName        string   `json:"orgName"`
	OrgDisplayName string   `json:"orgDisplayName"`
	PartOrgCodes   []string `json:"partOrgCodes"`
}

// PermMenuNode 权限系统菜单节点（对齐 /api/user/info 的 menus 数组结构）
type PermMenuNode struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	AppCode   string         `json:"appCode"`
	AppModule string         `json:"appModule"`
	ParentID  uint           `json:"parentId"`
	OrderNum  int            `json:"orderNum"`
	Path      string         `json:"path"`
	Icon      string         `json:"icon"`
	Visible   bool           `json:"visible"`
	Children  []PermMenuNode `json:"children"`
}

// PermButton 权限系统按钮资源（对齐 /api/user/info 的 buttons 数组结构）
type PermButton struct {
	ID           uint     `json:"id"`
	Resource     string   `json:"resource"`
	ResourceType string   `json:"resourceType"`
	Visible      bool     `json:"visible"`
	Actions      []string `json:"actions"`
	Conditions   []any    `json:"conditions"`
}

// PermOrg 权限系统组织资源（对齐 /api/user/info 的 orgs 数组结构）
type PermOrg struct {
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
	Children      []PermOrg `json:"children"`
}

// UserInfoResp /api/user/info 完整响应（权限系统对接时一次透传）
type UserInfoResp struct {
	Account UserAccount     `json:"account"`
	Roles   []string        `json:"roles"`
	Menus   []PermMenuNode  `json:"menus"`
	Orgs    []PermOrg       `json:"orgs"`
	Buttons []PermButton    `json:"buttons"`
}

// RemoteService 统一权限系统转发客户端
type RemoteService struct {
	cfg    config.PermissionConfig
	client *http.Client
}

// NewRemoteService 创建权限系统客户端（配置见 config 的 auth.permission 段）
func NewRemoteService(cfg config.PermissionConfig) *RemoteService {
	return &RemoteService{
		cfg:    cfg,
		client: &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second},
	}
}

// enabled 权限系统是否已对接：baseURL 为空 = 未对接（跳过远程调用）
func (s *RemoteService) enabled() bool {
	return s.cfg.BaseURL != ""
}

// getJSON 请求权限系统接口并解析 JSON；未对接/失败返回错误（由调用方兜底）。
// 所有请求统一附带 appId（本系统在权限系统的应用标识）。
func (s *RemoteService) getJSON(path string, query url.Values, out interface{}) error {
	if !s.enabled() {
		return fmt.Errorf("统一权限系统未对接（auth.permission.baseURL 为空），跳过远程调用: %s", path)
	}
	if query == nil {
		query = url.Values{}
	}
	if s.cfg.AppID != "" {
		query.Set("appId", s.cfg.AppID)
	}
	u := strings.TrimRight(s.cfg.BaseURL, "/") + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	start := time.Now()
	resp, err := s.client.Get(u)
	if err != nil {
		return fmt.Errorf("统一权限系统请求失败 %s: %v", u, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("统一权限系统返回异常 %s: HTTP %d", u, resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("统一权限系统响应解析失败 %s: %v", u, err)
	}
	log.Printf("[permission] OK 接口=%s 耗时=%dms", path, time.Since(start).Milliseconds())
	return nil
}

// warn 打印权限系统未对接/失败的警告日志（不阻断，由调用方兜底）
func (s *RemoteService) warn(err error) {
	log.Printf("[permission] 警告（权限系统未对接或调用失败，使用兜底数据）: %v", err)
}

// GetUserProfile 用户档案（user + roles + permissions）；未对接/失败返回 nil
func (s *RemoteService) GetUserProfile(username string) (*ProfileBundle, error) {
	var bundle ProfileBundle
	if err := s.getJSON(s.cfg.ProfilePath, url.Values{"username": {username}}, &bundle); err != nil {
		s.warn(err)
		return nil, err
	}
	return &bundle, nil
}

// GetMenus 菜单（roles + menus 透传）；未对接/失败返回 nil
func (s *RemoteService) GetMenus(username string) (*MenusBundle, error) {
	var bundle MenusBundle
	if err := s.getJSON(s.cfg.MenusPath, url.Values{"username": {username}}, &bundle); err != nil {
		s.warn(err)
		return nil, err
	}
	return &bundle, nil
}

// GetUserInfo 获取用户完整权限信息（account/roles/menus/orgs/buttons 一次透传，结构对齐 /api/user/info）；
// 未对接/失败返回 nil（由 handler 兜底静态权限资源）。
func (s *RemoteService) GetUserInfo(username string) (*UserInfoResp, error) {
	var info UserInfoResp
	if err := s.getJSON(s.cfg.ProfilePath, url.Values{"username": {username}}, &info); err != nil {
		s.warn(err)
		return nil, err
	}
	return &info, nil
}

// GetOrgDeptTeams 部门-所属团队（透传 JSON）；未对接/失败返回 nil
func (s *RemoteService) GetOrgDeptTeams() (json.RawMessage, error) {
	var raw json.RawMessage
	if err := s.getJSON(s.cfg.OrgDeptTeamsPath, nil, &raw); err != nil {
		s.warn(err)
		return nil, err
	}
	return raw, nil
}

// GetOrgDepartments 部门-团队-成员（透传 JSON）；未对接/失败返回 nil
func (s *RemoteService) GetOrgDepartments() (json.RawMessage, error) {
	var raw json.RawMessage
	if err := s.getJSON(s.cfg.OrgDeptsPath, nil, &raw); err != nil {
		s.warn(err)
		return nil, err
	}
	return raw, nil
}

// GetTeamUsers 团队人员（透传 JSON）；未对接/失败返回 nil
func (s *RemoteService) GetTeamUsers(team, deptName string) (json.RawMessage, error) {
	query := url.Values{}
	if team != "" {
		query.Set("team", team)
	}
	if deptName != "" {
		query.Set("dept_name", deptName)
	}
	var raw json.RawMessage
	if err := s.getJSON(s.cfg.TeamUsersPath, query, &raw); err != nil {
		s.warn(err)
		return nil, err
	}
	return raw, nil
}
