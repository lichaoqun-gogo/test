package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
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
	permSvc *permission.RemoteService // 统一权限系统转发客户端（用户档案/菜单/组织均由其下发）
}

func NewAuthHandler(cfg *config.Config, db *gorm.DB) *AuthHandler {
	return &AuthHandler{cfg: cfg, db: db, permSvc: permission.NewRemoteService(cfg.Auth.Permission)}
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
	Role        string `xml:"role"`        // 登录角色（本 app 的默认角色；真实 CAS 可在 attributes 中下发，无则默认 operator）
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

	// 校验 CAS 响应必须包含登录账号（user 字段）。
	// 真实公司 CAS 的 serviceValidate 返回结构可能与 mock 不同（cas: 命名空间、属性层级/名称差异），
	// 若 user 解析为空，此前会创建 username 为空的用户并签发 JWT，导致前端 /auth/me 返回
	// "未识别到用户"、判定认证失败清除 token、重新跳转 CAS 的死循环。
	// 此处强制校验并打印原始响应，便于按真实 CAS 结构适配字段。
	if strings.TrimSpace(casResp.Success.User) == "" {
		bodyPreview := string(body)
		if len(bodyPreview) > 800 {
			bodyPreview = bodyPreview[:800]
		}
		fmt.Printf("[CAS] 响应缺少 user 字段，原始响应: %s\n", bodyPreview)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "CAS 响应缺少用户标识（user 字段为空），请核对 CAS serviceValidate 返回的 XML 结构"})
		return
	}

	// 打印 CAS 返回的用户信息
	logBytes, _ := json.Marshal(casResp.Success)
	fmt.Printf("[CAS] 认证通过，用户信息: %s\n", string(logBytes))

	// 签发 JWT：身份凭证仅使用 CAS 返回的登录账号与角色（CAS 未下发角色时默认 operator）；
	// 角色/权限/菜单等授权信息统一由「统一权限系统」下发（/auth/me、/auth/menus），
	// 登录不再写入 users 表（废弃登录落库逻辑）。
	role := casResp.Success.Role
	if role == "" {
		role = "operator"
	}
	token, err := jwt.GenerateToken(
		h.cfg.Auth.JWT.Secret,
		h.cfg.Auth.JWT.Expiry,
		casResp.Success.User,
		casResp.Success.User,
		casResp.Success.AccountName,
		casResp.Success.Email,
		role,
		casResp.Success.DepName,
		casResp.Success.OrgName,
		false,
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
// 数据源：统一权限系统（user/roles/permissions 均由其下发，后端只转发）；
// 权限系统未对接或调用失败时打印日志并返回兜底数据（JWT 身份 + 默认 operator 角色 + 空权限），
// 保证登录与开发调试不被阻塞。
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	username := h.getRequestUsername(c)

	log.Printf("[auth/me] 请求用户信息: username=%s", username)
	if username == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "未识别到用户"})
		return
	}

	// 优先从统一权限系统获取用户档案
	if bundle, err := h.permSvc.GetUserProfile(username); err == nil && bundle != nil {
		log.Printf("[auth/me] 数据源=统一权限系统, username=%s, roles=%v, permissions=%v", username, bundle.Roles, bundle.Permissions)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"user":        bundle.User,
				"roles":       bundle.Roles,
				"permissions": bundle.Permissions,
			},
		})
		return
	}

	// ============ 兜底分支：权限系统未对接/调用失败时生效 ============
	// 对接真实权限系统后本分支不再进入，可整段删除。
	// 基于 JWT 身份返回基础档案（角色 = CAS 下发角色或默认 operator，对应"有 app 即有默认角色"的本地等价）。
	// 注意：/auth/me 不经过 AuthMiddleware（context 无身份信息），故从 Bearer JWT 解析身份。
	nickname, email, department, team, role := "", "", "", "", ""
	if authHeader := c.GetHeader("Authorization"); len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		if claims, err := jwt.ParseToken(h.cfg.Auth.JWT.Secret, authHeader[7:]); err == nil {
			nickname = claims.Nickname
			email = claims.Email
			department = claims.Department
			team = claims.Team
			role = claims.Role
		}
	}
	if role == "" {
		role = "operator"
	}
	log.Printf("[auth/me] 数据源=JWT兜底(权限系统未对接), username=%s, role=%s", username, role)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"user": gin.H{
				"username":   username,
				"nickname":   nickname,
				"email":      email,
				"department": department,
				"team":       team,
				"role":       role,
			},
			"roles":       []string{role},
			"permissions": []string{},
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
// 数据源：统一权限系统团队人员接口（透传）；未对接/失败时兜底查 users 表（保持开发调试可用）。
// 支持 team / dept_name 参数过滤（人员选择下拉框只显示本部门/团队人员）。
func (h *AuthHandler) GetUserList(c *gin.Context) {
	// 优先从统一权限系统获取团队人员（透传）
	if raw, err := h.permSvc.GetTeamUsers(c.Query("team"), c.Query("dept_name")); err == nil && raw != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": raw})
		return
	}

	// 兜底：users 表查询（权限系统未对接时的本地数据，对接后此路径不再生效）
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

// GetUserMenus 获取当前用户的菜单（按角色返回菜单树并集）
// 返回结构：{ roles: 角色编码, menus: { <role>: [菜单树] } }。
// 数据源：统一权限系统（菜单完全由其按角色下发，后端只做转发、不构建菜单、不再存表）；
// 权限系统未对接或调用失败时打印日志并返回空菜单兜底（前端按默认角色进入空工作台，见前端兜底路由）。
func (h *AuthHandler) GetUserMenus(c *gin.Context) {
	username := h.getRequestUsername(c)

	log.Printf("[auth/menus] 请求用户菜单: username=%s", username)
	if username == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "未识别到用户"})
		return
	}

	// 优先从统一权限系统获取菜单（roles + 各角色菜单树，透传）
	if bundle, err := h.permSvc.GetMenus(username); err == nil && bundle != nil {
		menuCount := 0
		for _, tree := range bundle.Menus {
			menuCount += len(tree)
		}
		log.Printf("[auth/menus] 数据源=统一权限系统, username=%s, roles=%v, 菜单数=%d", username, bundle.Roles, menuCount)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"roles": bundle.Roles,
				"menus": bundle.Menus,
			},
		})
		return
	}

	// ============ 兜底分支：权限系统未对接/调用失败时生效 ============
	// 对接真实权限系统后本分支不再进入，可整段删除。
	// 返回空菜单（角色 = JWT/CAS 角色或默认 operator），前端按角色用本地菜单映射渲染（见前端 menuMap.ts）。
	// 注意：/auth/menus 不经过 AuthMiddleware（context 无身份信息），故从 Bearer JWT 解析角色。
	role := ""
	if authHeader := c.GetHeader("Authorization"); len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		if claims, err := jwt.ParseToken(h.cfg.Auth.JWT.Secret, authHeader[7:]); err == nil {
			role = claims.Role
		}
	}
	if role == "" {
		role = "operator"
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"roles": []string{role},
			"menus": gin.H{role: []permission.MenuNode{}},
		},
	})
}
