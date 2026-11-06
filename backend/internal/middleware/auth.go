package middleware

import (
	"net/http"
	"osm-mngt/internal/config"
	"osm-mngt/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
// CAS 模式下验证 JWT：优先取 Authorization Bearer token，无则从 Cookie osm_token 读取（HttpOnly，浏览器自动携带）。
// local 模式下兼容前端 X-User-* 请求头（开发用）。
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// local 模式：兼容前端 mock 请求头，不做强制认证
		if cfg.Auth.Mode == "local" {
			c.Next()
			return
		}

		// CAS 模式：验证 JWT
		token := ""
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// 提取 Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
		// Authorization 头缺失/格式不对时，回退到 Cookie（HttpOnly，浏览器自动携带）
		if token == "" {
			if ck, err := c.Cookie("osm_token"); err == nil {
				token = ck
			}
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(cfg.Auth.JWT.Secret, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "token 无效或已过期"})
			c.Abort()
			return
		}

		// 注入用户信息到 context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("nickname", claims.Nickname)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("department", claims.Department)
		c.Set("team", claims.Team)
		c.Set("canToggleAdminView", claims.CanToggleAdminView)

		// 同时设置请求头，兼容现有 handler 的 getUserInfo 逻辑
		c.Request.Header.Set("X-User-ID", claims.UserID)
		c.Request.Header.Set("X-User-Name", claims.Nickname)
		c.Request.Header.Set("X-User-Email", claims.Email)
		c.Request.Header.Set("X-User-Role", claims.Role)
		c.Request.Header.Set("X-User-Dept", claims.Department)
		c.Request.Header.Set("X-User-Team", claims.Team)
		if claims.CanToggleAdminView {
			c.Request.Header.Set("X-User-Can-Admin-View", "true")
		}

		c.Next()
	}
}
