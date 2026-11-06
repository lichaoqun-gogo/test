package handlers

// ==================== 参考项目兼容接口层 ====================
// 对应 osm-front 参考项目（Vue2 + vue-cli webpack）请求的 6 个接口，统一 {ok, data} 结构，
// 便于参考前端页面/文件对位迁移到本项目（无限接近参考项目）。
//
// 数据源约定：优先「统一权限系统」（config auth.permission 或 IAM OpenAPI），
// 未对接/失败时兜底（users 表 / JWT / 默认数据），保证登录与开发调试不被阻塞。
//
// 参考项目接口清单（见 osm-front 技术文档第七节）：
//   1. GET /api/env             → {ok, data:{profiles, tips}}
//   2. GET /api/groups          → {ok, data:[{code, name}]}
//   3. GET /api/user/depts      → {ok, data:{deptTree, deptList}}
//   4. GET /api/idcs            → {ok, data:[...]}
//   5. GET /api/user/info       → {ok, data:{account, roles, menus, orgs, buttons}}
//   6. GET /api/user/refresh    → 刷新权限缓存

import (
	"log"
	"net/http"

	"osm-mngt/internal/config"
	"osm-mngt/internal/service/permission"
	"osm-mngt/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ReferenceHandler 参考项目兼容接口处理器
type ReferenceHandler struct {
	cfg     *config.Config
	db      *gorm.DB
	permSvc *permission.RemoteService // 统一权限系统转发客户端
}

func NewReferenceHandler(cfg *config.Config, db *gorm.DB) *ReferenceHandler {
	return &ReferenceHandler{cfg: cfg, db: db, permSvc: permission.NewRemoteService(cfg.Auth.Permission)}
}

// okData 参考项目统一响应结构 {ok: true, data: ...}
func okData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"ok": true, "data": data})
}

// jwtClaimsFromRequest 从 Authorization Bearer JWT 或 Cookie osm_token 解析身份（参考接口不经过 AuthMiddleware）
func (h *ReferenceHandler) jwtClaimsFromRequest(c *gin.Context) *jwt.Claims {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		if claims, err := jwt.ParseToken(h.cfg.Auth.JWT.Secret, authHeader[7:]); err == nil {
			return claims
		}
	}
	// 无 Authorization 头时回退到 Cookie（HttpOnly，浏览器自动携带）
	if ck, err := c.Cookie("osm_token"); err == nil && ck != "" {
		if claims, err := jwt.ParseToken(h.cfg.Auth.JWT.Secret, ck); err == nil {
			return claims
		}
	}
	return nil
}

// Env GET /api/env 环境配置（参考前端 App.vue 初始化时调用）
// {ok, data:{profiles:[{name, active}], tips:[]}}
func (h *ReferenceHandler) Env(c *gin.Context) {
	profiles := []gin.H{
		{"name": "dev", "active": true},
	}
	okData(c, gin.H{"profiles": profiles, "tips": []string{}})
}

// Groups GET /api/groups 组/部门列表（参考前端 App.vue 初始化时调用）
// {ok, data:[{code, name}]}
// 数据源：优先权限系统（orgDeptTeams 透传），未对接时兜底 users 表按部门聚合。
func (h *ReferenceHandler) Groups(c *gin.Context) {
	// ============ 兜底分支：权限系统未对接/失败时生效（对接后可整段删除）============
	// 未对接时从 users 表按部门聚合（本地开发数据）
	var rows []struct {
		Department string
	}
	if err := h.db.Table("users").Select("department").Where("department <> ''").Group("department").Order("department ASC").Scan(&rows).Error; err != nil {
		okData(c, []gin.H{})
		return
	}
	out := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		out = append(out, gin.H{"code": r.Department, "name": r.Department})
	}
	okData(c, out)
}

// UserDepts GET /api/user/depts 部门树/部门列表（参考前端渲染左侧部门树）
// {ok, data:{deptTree:[], deptList:[{code, name, parentCode}]}}
// 数据源：未对接时兜底 users 表按 部门→团队 构建。
func (h *ReferenceHandler) UserDepts(c *gin.Context) {
	// ============ 兜底分支：权限系统未对接/失败时生效（对接后可整段删除）============
	type row struct {
		Department string
		Team       string
	}
	var rows []row
	h.db.Table("users").Select("department, team").Where("department <> ''").Group("department, team").Order("department ASC, team ASC").Scan(&rows)

	// deptList：部门列表（code=部门名，parentCode 空=顶级）
	deptList := []gin.H{}
	seen := map[string]bool{}
	for _, r := range rows {
		if !seen[r.Department] {
			seen[r.Department] = true
			deptList = append(deptList, gin.H{"code": r.Department, "name": r.Department, "parentCode": ""})
		}
	}
	// deptTree：简化为一层部门树（团队作为子节点）
	deptTree := []gin.H{}
	seenTree := map[string]bool{}
	for _, r := range rows {
		if !seenTree[r.Department] {
			seenTree[r.Department] = true
			deptTree = append(deptTree, gin.H{"code": r.Department, "name": r.Department, "children": []gin.H{}})
		}
	}
	okData(c, gin.H{"deptTree": deptTree, "deptList": deptList})
}

// Idcs GET /api/idcs IDC 列表（参考前端 IDC 选择器数据源，当前无数据源，返回空）
// {ok, data:[]}
func (h *ReferenceHandler) Idcs(c *gin.Context) {
	okData(c, []gin.H{})
}

// UserInfo GET /api/user/info 当前用户完整权限信息（对应参考项目 main.js 路由守卫访问 /user/info）
// {ok, data:{account, roles, permissions, menus, menusByRole, orgs, buttons}}
// account 字段对齐参考项目：id=账号、displayName=姓名、groupName=组/部门、groupCode、orgCode。
// roles/permissions 对应权限档案；menusByRole 为多角色菜单树（承载原 /auth/menus 职责，前端按角色并集渲染）；
// menus 为参考项目扁平菜单结构（{path, name, children}，参考页面使用）。
// 数据源：优先权限系统（GetUserProfile/GetMenus），未对接时兜底 JWT 身份 + 空权限（前端本地映射兜底菜单）。
func (h *ReferenceHandler) UserInfo(c *gin.Context) {
	claims := h.jwtClaimsFromRequest(c)
	username := ""
	if claims != nil {
		username = claims.Username
	}
	if username == "" {
		username = c.GetHeader("X-User-ID")
	}
	log.Printf("[user-info] 请求用户信息: username=%s", username)

	// 默认兜底数据（权限系统未对接时）
	account := gin.H{
		"id":          username,
		"displayName": "",
		"groupName":   "",
		"groupCode":   "",
		"orgCode":     "",
	}
	roles := []string{}
	permissions := []string{}
	buttons := []string{}
	if claims != nil {
		account = gin.H{
			"id":          claims.Username,
			"displayName": claims.Nickname,
			"groupName":   claims.Department,
			"groupCode":   claims.Department,
			"orgCode":     claims.Team,
		}
		roles = []string{claims.Role}
	}

	// 优先从统一权限系统获取用户档案（user/roles/permissions）
	if profile, err := h.permSvc.GetUserProfile(username); err == nil && profile != nil {
		if profile.User.Nickname != "" {
			account["displayName"] = profile.User.Nickname
		}
		if profile.User.Department != "" {
			account["groupName"] = profile.User.Department
			account["groupCode"] = profile.User.Department
		}
		if profile.User.Team != "" {
			account["orgCode"] = profile.User.Team
		}
		roles = profile.Roles
		permissions = profile.Permissions
		buttons = profile.Permissions
	}

	// 菜单：权限系统按角色下发菜单树（menusByRole）与参考扁平结构（menus）；未对接时为空（前端按角色本地映射兜底）
	menusByRole := gin.H{}
	refMenus := []gin.H{}
	if bundle, err := h.permSvc.GetMenus(username); err == nil && bundle != nil {
		for role, tree := range bundle.Menus {
			menusByRole[role] = tree
		}
		refMenus = flattenMenus(bundle)
	}

	okData(c, gin.H{
		"account":     account,
		"roles":       roles,
		"permissions": permissions,
		"menus":       refMenus,
		"menusByRole": menusByRole,
		"orgs":        []gin.H{},
		"buttons":     buttons,
	})
}

// flattenMenus 把权限系统各角色菜单树拍平为参考项目 menus 结构（{path, name, children}）
func flattenMenus(bundle *permission.MenusBundle) []gin.H {
	out := []gin.H{}
	for _, tree := range bundle.Menus {
		for _, node := range tree {
			out = append(out, menuNodeToRef(node))
		}
	}
	return out
}

func menuNodeToRef(node permission.MenuNode) gin.H {
	children := []gin.H{}
	for _, ch := range node.Children {
		children = append(children, menuNodeToRef(ch))
	}
	return gin.H{"path": node.Path, "name": node.Name, "children": children}
}

// UserRefresh GET /api/user/refresh 刷新权限缓存（参考前端手动刷新权限时调用）
// {ok, data:{}}：重查统一权限系统，未对接时返回空（无缓存可刷）。
func (h *ReferenceHandler) UserRefresh(c *gin.Context) {
	username := c.GetHeader("X-User-ID")
	if claims := h.jwtClaimsFromRequest(c); claims != nil {
		username = claims.Username
	}
	log.Printf("[reference/user-refresh] 刷新权限: username=%s", username)
	// 权限系统未对接（baseURL 空）时 GetUserProfile 内部已打印警告，这里只触发一次以刷新缓存语义
	_, _ = h.permSvc.GetUserProfile(username)
	okData(c, gin.H{})
}
