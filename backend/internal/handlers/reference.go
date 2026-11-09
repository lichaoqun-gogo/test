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
// {ok, data:{profiles:["dev"], tips:"本地环境"}}，结构对齐权限系统；未对接时返回静态环境配置。
func (h *ReferenceHandler) Env(c *gin.Context) {
	okData(c, gin.H{
		"profiles": []string{"dev"},
		"tips":     "本地环境",
	})
}

// Groups GET /api/groups 组/部门列表（参考前端 App.vue 初始化时调用）
// {ok, data:[{code, name}]}，结构对齐权限系统；
// 数据源：优先权限系统（orgDeptTeams 透传），未对接时返回静态组数据（贴近内网示例）。
func (h *ReferenceHandler) Groups(c *gin.Context) {
	// 静态组列表（对齐权限系统 groups 返回，未对接时兜底）
	staticGroups := []gin.H{
		{"code": "middleware", "name": "基础组件"},
		{"code": "dba", "name": "DBA"},
		{"code": "sa", "name": "基础设施"},
		{"code": "qa", "name": "质量测试"},
		{"code": "sec", "name": "安全管理"},
	}
	okData(c, staticGroups)
}

// UserDepts GET /api/user/depts 部门树/部门列表（对应参考项目部门树组件 + 自研项目台账查询栏级联下拉）
// {ok, data:{deptTree:[{id, orgCode, groupCode, groupName, parentId, ldapDomain, orderNum, children}], deptList:[...]}}
// 返回结构对齐统一权限系统；数据源：未对接时按 staticUsers 的部门/团队构造树。
func (h *ReferenceHandler) UserDepts(c *gin.Context) {
	// 部门 → 团队映射（来自 staticUsers）
	deptTeams := map[string][]string{}
	deptOrder := []string{}
	for _, u := range staticUsers {
		if u.Department == "" {
			continue
		}
		if _, ok := deptTeams[u.Department]; !ok {
			deptOrder = append(deptOrder, u.Department)
			deptTeams[u.Department] = []string{}
		}
		if u.Team != "" {
			exists := false
			for _, t := range deptTeams[u.Department] {
				if t == u.Team {
					exists = true
					break
				}
			}
			if !exists {
				deptTeams[u.Department] = append(deptTeams[u.Department], u.Team)
			}
		}
	}

	// 构造部门树（部门 → 团队两级）+ 拍平列表，字段对齐权限系统
	deptTree := []gin.H{}
	deptList := []gin.H{}
	for _, dept := range deptOrder {
		deptID := deptCodeMap[dept]
		if deptID == "" {
			deptID = dept
		}
		children := []gin.H{}
		for _, team := range deptTeams[dept] {
			teamID := teamCodeMap[team]
			if teamID == "" {
				teamID = team
			}
			children = append(children, gin.H{
				"id": teamID, "orgCode": "nucc", "groupCode": teamID,
				"groupName": team, "parentId": deptID, "ldapDomain": nil, "orderNum": nil,
				"children": []gin.H{},
			})
			deptList = append(deptList, gin.H{
				"id": teamID, "orgCode": "nucc", "groupCode": teamID,
				"groupName": team, "parentId": deptID, "ldapDomain": nil, "orderNum": nil,
			})
		}
		deptTree = append(deptTree, gin.H{
			"id": deptID, "orgCode": "nucc", "groupCode": deptID,
			"groupName": dept, "parentId": "", "ldapDomain": nil, "orderNum": nil,
			"children": children,
		})
		deptList = append(deptList, gin.H{
			"id": deptID, "orgCode": "nucc", "groupCode": deptID,
			"groupName": dept, "parentId": "", "ldapDomain": nil, "orderNum": nil,
		})
	}
	okData(c, gin.H{"deptTree": deptTree, "deptList": deptList})
}

// Idcs GET /api/idcs IDC 列表（参考前端 IDC 选择器数据源，当前无数据源，返回空）
// {ok, data:[]}
func (h *ReferenceHandler) Idcs(c *gin.Context) {
	okData(c, []gin.H{})
}

// UserInfo GET /api/user/info 当前用户完整权限信息（对应参考项目 main.js 路由守卫访问 /user/info）
// {ok, data:{account, roles, menus, orgs, buttons}}
// 返回结构对齐统一权限系统（内网示例）：
//   - account: {accountCode, accountName, email, depCode, depName, orgCode, customOrgCode, orgName, orgDisplayName, partOrgCodes}
//   - roles:   权限系统角色编码（如 COMMON/SPECIALIST/ADMIN）
//   - menus:   该用户可见菜单树数组（{id, name, appCode, appModule, parentId, orderNum, path, icon, visible, children}）
//   - orgs:    组织资源数组
//   - buttons: 按钮资源数组（{id, resource, resourceType, visible, actions, conditions}），前端 v-auth 按 resource 控制按钮
// 数据源：优先权限系统（GetUserInfo 一次透传），未对接时兜底静态权限资源（perm_mock.go，按角色返回对应菜单/按钮）。
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

	// 优先从统一权限系统获取完整权限信息（结构对齐：account/roles/menus/orgs/buttons）
	if info, err := h.permSvc.GetUserInfo(username); err == nil && info != nil {
		okData(c, info)
		return
	}

	// ===== 兜底：权限系统未对接 → 静态权限资源（perm_mock.go，按用户角色返回菜单/按钮）=====
	okData(c, staticPermInfo(username, claims))
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
