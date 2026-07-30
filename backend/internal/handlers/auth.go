package handlers

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"osm-mngt/internal/config"
	"osm-mngt/internal/models"
	"osm-mngt/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewAuthHandler(cfg *config.Config, db *gorm.DB) *AuthHandler {
	return &AuthHandler{cfg: cfg, db: db}
}

// CAS XML 响应结构
type casServiceResponse struct {
	XMLName xml.Name `xml:"serviceResponse"`
	Success *casAuthSuccess `xml:"authenticationSuccess"`
	Failure *casAuthFailure `xml:"authenticationFailure"`
}

type casAuthSuccess struct {
	User        string `xml:"user"`
	DisplayName string `xml:"displayName"`
	Mail        string `xml:"mail"`
	Phone       string `xml:"phone"`
	Department  string `xml:"department"`
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

	// 查找或创建用户
	var user models.User
	result := h.db.Where("username = ?", casResp.Success.User).First(&user)
	if result.Error != nil {
		// 首次登录，自动创建，默认 operator 角色
		user = models.User{
			Username:   casResp.Success.User,
			Nickname:   casResp.Success.DisplayName,
			Email:      casResp.Success.Mail,
			Phone:      casResp.Success.Phone,
			Department: casResp.Success.Department,
			Role:       "operator",
		}
		h.db.Create(&user)
	} else {
		// 更新 CAS 返回的最新信息
		user.Nickname = casResp.Success.DisplayName
		user.Email = casResp.Success.Mail
		user.Phone = casResp.Success.Phone
		user.Department = casResp.Success.Department
		h.db.Save(&user)
	}

	// 签发 JWT
	token, err := jwt.GenerateToken(
		h.cfg.Auth.JWT.Secret,
		h.cfg.Auth.JWT.Expiry,
		fmt.Sprintf("%d", user.ID),
		user.Username,
		user.Nickname,
		user.Role,
		user.Department,
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

// GetUserInfo 获取当前登录用户信息
// 优先从中间件注入的 context 读取，其次从 Authorization JWT 解析，最后从 header 回退
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	nickname, _ := c.Get("nickname")
	role, _ := c.Get("role")
	department, _ := c.Get("department")
	canToggle, _ := c.Get("canToggleAdminView")

	// 如果 context 中没有值（未走中间件），尝试从 JWT 解析
	if userID == nil {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenStr := authHeader[7:]
			claims, err := jwt.ParseToken(h.cfg.Auth.JWT.Secret, tokenStr)
			if err == nil {
				userID = claims.UserID
				username = claims.Username
				nickname = claims.Nickname
				role = claims.Role
				department = claims.Department
				canToggle = claims.CanToggleAdminView
			}
		}
	}

	// 仍然为空时从 header 回退
	if role == nil || role == "" {
		role = c.GetHeader("X-User-Role")
	}
	if username == nil || username == "" {
		username = decodeHeaderValue(c.GetHeader("X-User-Name"))
	}
	if userID == nil || userID == "" {
		userID = c.GetHeader("X-User-ID")
	}
	if department == nil || department == "" {
		department = decodeHeaderValue(c.GetHeader("X-User-Dept"))
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"id":                 userID,
			"username":           username,
			"nickname":           nickname,
			"role":               role,
			"department":         department,
			"canToggleAdminView": canToggle,
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

// GetUserMenus 获取当前用户的菜单树
func (h *AuthHandler) GetUserMenus(c *gin.Context) {
	// 优先从 context 获取（CAS 模式下中间件注入）
	roleVal, exists := c.Get("role")
	var roleStr string
	if exists {
		roleStr, _ = roleVal.(string)
	}
	// context 中没有时，尝试从 JWT 解析
	if roleStr == "" {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenStr := authHeader[7:]
			claims, err := jwt.ParseToken(h.cfg.Auth.JWT.Secret, tokenStr)
			if err == nil {
				roleStr = claims.Role
			}
		}
	}
	// 仍然为空时从 header 回退
	if roleStr == "" {
		roleStr = c.GetHeader("X-User-Role")
	}
	if roleStr == "" {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []MenuTreeNode{}})
		return
	}

	// 支持 query 参数指定角色（admin 用户切换个人工作台时获取 specialist 菜单）
	if roleParam := c.Query("role"); roleParam != "" {
		roleStr = roleParam
	}

	// 查询角色关联的菜单ID
	var roleMenus []models.RoleMenu
	h.db.Where("role = ?", roleStr).Find(&roleMenus)
	if len(roleMenus) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []MenuTreeNode{}})
		return
	}

	menuIDs := make([]uint, len(roleMenus))
	for i, rm := range roleMenus {
		menuIDs[i] = rm.MenuID
	}

	// 查询菜单
	var menus []models.Menu
	h.db.Where("id IN ?", menuIDs).Order("sort ASC").Find(&menus)

	// 构建菜单树
	tree := buildMenuTree(menus, 0)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": tree})
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
