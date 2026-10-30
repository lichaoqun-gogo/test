package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"osm-mngt/internal/config"
	"osm-mngt/internal/models"
	"osm-mngt/internal/service/permission"
	"osm-mngt/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	cfg     *config.Config
	db      *gorm.DB
	permSvc permission.Service
}

func NewAuthHandler(cfg *config.Config, db *gorm.DB) *AuthHandler {
	return &AuthHandler{cfg: cfg, db: db, permSvc: permission.NewLocalService(db)}
}

// CAS XML 响应结构
type casServiceResponse struct {
	XMLName xml.Name `xml:"serviceResponse"`
	Success *casAuthSuccess `xml:"authenticationSuccess"`
	Failure *casAuthFailure `xml:"authenticationFailure"`
}

// casAuthSuccess CAS 认证成功响应的用户属性。
// 属性名与 IAM Account 模型字段对齐（与 mock CAS 测试服务 users.json 返回保持一致），
// 登录后映射到 users 表与 JWT：User→username、AccountName→nickname、Email→email、
// DepName→department、OrgName→team；后续 IAM 接入时同一份数据可直接承接 Account。
type casAuthSuccess struct {
	User        string `xml:"user"`        // 登录账号
	AccountName string `xml:"accountName"` // 姓名
	Email       string `xml:"email"`       // 邮箱
	Phone       string `xml:"phone"`       // 电话
	DepName     string `xml:"depName"`     // 部门
	OrgName     string `xml:"orgName"`     // 组织/团队
}

type casAuthFailure struct {
	Code    string `xml:"code,attr"`
	Message string `xml:",chardata"`
}

// getValidateURL 获取 CAS 票据验证地址，未配置时自动推导为 {serverURL}/serviceValidate
func (h *AuthHandler) getValidateURL() string {
	if h.cfg.Auth.CAS.ValidateURL != "" {
		return h.cfg.Auth.CAS.ValidateURL
	}
	return strings.TrimRight(h.cfg.Auth.CAS.ServerURL, "/") + "/serviceValidate"
}

// getServiceURL 获取后端回调地址，未配置时从请求 Host 自动推导
func (h *AuthHandler) getServiceURL(c *gin.Context) string {
	if h.cfg.Auth.CAS.ServiceURL != "" {
		return h.cfg.Auth.CAS.ServiceURL
	}
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/api/v1/auth/cas/callback", scheme, c.Request.Host)
}

// CasLogin 重定向到 CAS 登录页
func (h *AuthHandler) CasLogin(c *gin.Context) {
	serviceURL := h.getServiceURL(c)
	casURL := fmt.Sprintf("%s/login?service=%s",
		h.cfg.Auth.CAS.ServerURL,
		url.QueryEscape(serviceURL),
	)
	c.Redirect(http.StatusFound, casURL)
}

// CasCallback CAS 回调：验证 ticket，签发 JWT，重定向到前端
func (h *AuthHandler) CasCallback(c *gin.Context) {
	ticket := c.Query("ticket")
	if ticket == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "缺少 ticket 参数"})
		return
	}

	serviceURL := h.getServiceURL(c)

	// 调用 CAS /serviceValidate 验证票据
	validateURL := fmt.Sprintf("%s?ticket=%s&service=%s",
		h.getValidateURL(),
		url.QueryEscape(ticket),
		url.QueryEscape(serviceURL),
	)

	resp, err := http.Get(validateURL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "CAS 服务不可达: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取 CAS 响应失败"})
		return
	}

	var casResp casServiceResponse
	if err := xml.Unmarshal(body, &casResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "解析 CAS 响应失败"})
		return
	}

	if casResp.Failure != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "CAS 验证失败: " + casResp.Failure.Message})
		return
	}

	if casResp.Success == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "CAS 验证失败"})
		return
	}

	// 打印 CAS 返回的用户信息
	logBytes, _ := json.Marshal(casResp.Success)
	fmt.Printf("[CAS] 认证通过，用户信息: %s\n", string(logBytes))

	// 查找或创建用户（CAS 属性 → users 表：accountName→nickname、email、depName→department、orgName→team）
	var user models.User
	result := h.db.Where("username = ?", casResp.Success.User).First(&user)
	if result.Error != nil {
		// 首次登录，自动创建，默认 operator 角色
		user = models.User{
			Username:   casResp.Success.User,
			Nickname:   casResp.Success.AccountName,
			Email:      casResp.Success.Email,
			Phone:      casResp.Success.Phone,
			Department: casResp.Success.DepName,
			Team:       casResp.Success.OrgName,
			Role:       "operator",
		}
		h.db.Create(&user)
	} else {
		// 更新 CAS 返回的最新信息
		user.Nickname = casResp.Success.AccountName
		user.Email = casResp.Success.Email
		user.Phone = casResp.Success.Phone
		user.Department = casResp.Success.DepName
		user.Team = casResp.Success.OrgName
		h.db.Save(&user)
	}

	// 签发 JWT
	// 统一使用 CAS 用户名（username）作为全系统用户身份标识，
	// 与 dept_specialist_mapping.specialist_id、approval_node.approver_id、governance_task.applicant_id 保持一致
	token, err := jwt.GenerateToken(
		h.cfg.Auth.JWT.Secret,
		h.cfg.Auth.JWT.Expiry,
		user.Username,
		user.Username,
		user.Nickname,
		user.Email,
		user.Role,
		user.Department,
		user.Team,
		user.CanToggleAdminView,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "生成 token 失败"})
		return
	}

	// 重定向到前端，携带 token（frontendURL 为空时用相对路径，前后端同域时自动生效）
	frontendURL := h.cfg.Auth.CAS.FrontendURL
	redirectURL := fmt.Sprintf("%s/auth/cas-callback?token=%s", frontendURL, url.QueryEscape(token))
	c.Redirect(http.StatusFound, redirectURL)
}

// getRequestUsername 从 context/JWT/请求头提取当前用户名（CAS username，全系统用户唯一标识）
func (h *AuthHandler) getRequestUsername(c *gin.Context) string {
	if v, ok := c.Get("username"); ok {
		if s, _ := v.(string); s != "" {
			return s
		}
	}
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		if claims, err := jwt.ParseToken(h.cfg.Auth.JWT.Secret, authHeader[7:]); err == nil {
			return claims.Username
		}
	}
	if id := c.GetHeader("X-User-ID"); id != "" {
		return id
	}
	return decodeHeaderValue(c.GetHeader("X-User-Name"))
}

// GetUserInfo 获取当前登录用户信息（含权限档案）
// 返回结构：{ user: 用户基本信息, roles: 角色编码, permissions: 权限点编码 }。
// roles/permissions 由权限服务计算，前端据此渲染可见菜单与按钮；
// 后续对接外部权限系统时仅需替换 permission.Service 实现，此接口契约不变。
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	username := h.getRequestUsername(c)

	if username == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "未识别到用户"})
		return
	}

	// 查询 users 表获取完整用户信息
	var user models.User
	if err := h.db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "用户不存在"})
		return
	}

	profile := h.permSvc.GetProfile(username)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"user":        user,
			"roles":       profile.Roles,
			"permissions": profile.Permissions,
		},
	})
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// CAS 模式下重定向到 CAS 登出页
	if h.cfg.Auth.Mode == "cas" {
		casLogoutURL := fmt.Sprintf("%s/logout?service=%s",
			h.cfg.Auth.CAS.ServerURL,
			url.QueryEscape(h.cfg.Auth.CAS.FrontendURL),
		)
		c.Redirect(http.StatusFound, casLogoutURL)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "退出成功"})
}

// GetAuthMode 返回当前认证模式（local/cas），前端用于决定登录方式
func (h *AuthHandler) GetAuthMode(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"mode": h.cfg.Auth.Mode})
}

// UserSimple 用户精简信息（用于人员选择下拉框等场景）
// 注意：当前从 users 表查询，后期对接统一权限系统时只需修改此方法实现
type UserSimple struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Department string `json:"department"`
	Team       string `json:"team"`
	Email      string `json:"email"`
}

// GetUserList 查询用户列表（用于人员选择等场景）
// 当前数据源：users 表；后期对接统一权限系统时修改此方法即可
// 支持 dept_name 参数过滤本部门人员
func (h *AuthHandler) GetUserList(c *gin.Context) {
	query := h.db.Model(&models.User{})

	// 按部门过滤（人员选择下拉框只显示本部门人员）
	if deptName := c.Query("dept_name"); deptName != "" {
		query = query.Where("department = ?", deptName)
	}
	// 按团队过滤（组织架构-团队人员查询）
	if team := c.Query("team"); team != "" {
		query = query.Where("team = ?", team)
	}

	var users []models.User
	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "查询用户列表失败"})
		return
	}
	list := make([]UserSimple, 0, len(users))
	for _, u := range users {
		list = append(list, UserSimple{
			ID:         u.ID,
			Username:   u.Username,
			Nickname:   u.Nickname,
			Department: u.Department,
			Team:       u.Team,
			Email:      u.Email,
		})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": list})
}

// MenuTreeNode 菜单树节点
type MenuTreeNode struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Component string         `json:"component"`
	Icon      string         `json:"icon"`
	Sort      int            `json:"sort"`
	Hidden    bool           `json:"hidden"`
	Group     string         `json:"group"`
	Children  []MenuTreeNode `json:"children"`
}

// GetUserMenus 获取当前用户的菜单（按用户全部角色返回并集，每角色一棵菜单树）
// 模型：一个用户可具备多个角色，菜单可见性 = 全部角色菜单的并集，由后端统一计算下发，
// 前端一次拉取全部菜单并注册路由，视图切换（个人工作台/全局治理）只是前端按角色分组展示。
// 菜单角色本地映射规则（菜单-角色映射维护在本地 role_menus 表，对接外部权限系统时
// 替换为接口返回的角色集即可）：
//   - 基础角色：users.role（operator / specialist / admin）
//   - 具备 global_governance_view 权限（如管理部门专员、单独授权）→ 附加 admin（全局治理菜单）
//   - admin 用户 → 附加 specialist（系统管理员同样有个人工作台）
// 响应结构：{ code, data: { roles: [角色编码], menus: { <role>: [菜单树] } } }
func (h *AuthHandler) GetUserMenus(c *gin.Context) {
	username := h.getRequestUsername(c)
	if username == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "未识别到用户"})
		return
	}

	var user models.User
	if err := h.db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "用户不存在"})
		return
	}

	// 计算用户的菜单角色集合（去重、保持顺序）
	roleSet := make([]string, 0, 2)
	seen := map[string]bool{}
	addRole := func(r string) {
		if r == "" || seen[r] {
			return
		}
		seen[r] = true
		roleSet = append(roleSet, r)
	}
	addRole(user.Role)
	// admin 用户附加 specialist 菜单（个人工作台）
	if user.Role == "admin" {
		addRole("specialist")
	}
	// 具备全局治理权限 → 附加 admin 菜单（全局治理）
	if h.permSvc.HasPermission(username, permission.PermGlobalGovernanceView) {
		addRole("admin")
	}

	// 按角色查询菜单并构建树（并集）
	menusByRole := gin.H{}
	for _, role := range roleSet {
		var roleMenus []models.RoleMenu
		h.db.Where("role = ?", role).Find(&roleMenus)
		if len(roleMenus) == 0 {
			menusByRole[role] = []MenuTreeNode{}
			continue
		}
		menuIDs := make([]uint, len(roleMenus))
		for i, rm := range roleMenus {
			menuIDs[i] = rm.MenuID
		}
		var menus []models.Menu
		h.db.Where("id IN ?", menuIDs).Order("sort ASC").Find(&menus)
		menusByRole[role] = buildMenuTree(menus, 0)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"roles": roleSet,
			"menus": menusByRole,
		},
	})
}

func buildMenuTree(menus []models.Menu, parentID uint) []MenuTreeNode {
	var tree []MenuTreeNode
	for _, m := range menus {
		if m.ParentID == parentID {
			node := MenuTreeNode{
				ID:        m.ID,
				Name:      m.Name,
				Path:      m.Path,
				Component: m.Component,
				Icon:      m.Icon,
				Sort:      m.Sort,
				Hidden:    m.Hidden,
				Group:     m.Group,
				Children:  buildMenuTree(menus, m.ID),
			}
			tree = append(tree, node)
		}
	}
	return tree
}
