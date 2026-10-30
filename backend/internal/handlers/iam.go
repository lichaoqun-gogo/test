package handlers

import (
	"net/http"

	"osm-mngt/internal/iam"
	"osm-mngt/internal/service/identity"

	"github.com/gin-gonic/gin"
)

// IamHandler 用户身份/权限资源接口（对接 IAM 后暴露给前端的对外能力）
type IamHandler struct {
	identitySvc identity.Service
}

func NewIamHandler(identitySvc identity.Service) *IamHandler {
	return &IamHandler{identitySvc: identitySvc}
}

// GetUserInfo 获取当前登录用户的「全量权限资源包」（IAMUser）。
// 身份（account）来自 JWT/请求头（CAS 登录后写入），组织等权限资源由身份服务拼接；
// handler 不感知数据源是本地 users 表还是 IAM 开放接口。
func (h *IamHandler) GetUserInfo(c *gin.Context) {
	account := currentAccount(c)
	if account.AccountCode == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "未识别到用户"})
		return
	}
	user, err := h.identitySvc.GetUserInfo(c.Request.Context(), account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "获取用户信息失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
}

// currentAccount 从 JWT context（CAS 模式）/ 请求头（local 模式）解析当前用户身份，构造 IAM Account。
// 与 AuthMiddleware 注入的字段保持一致（username/nickname/email/department/team），
// 即 CAS 登录返回的 Account 核心属性：accountCode/accountName/email/depName/orgName。
func currentAccount(c *gin.Context) *iam.Account {
	username := c.GetString("username")
	if username == "" {
		username = decodeHeaderValue(c.GetHeader("X-User-ID"))
	}
	nickname := c.GetString("nickname")
	if nickname == "" {
		nickname = decodeHeaderValue(c.GetHeader("X-User-Name"))
	}
	email := c.GetString("email")
	if email == "" {
		email = decodeHeaderValue(c.GetHeader("X-User-Email"))
	}
	department := c.GetString("department")
	if department == "" {
		department = decodeHeaderValue(c.GetHeader("X-User-Dept"))
	}
	team := c.GetString("team")
	if team == "" {
		team = decodeHeaderValue(c.GetHeader("X-User-Team"))
	}
	return &iam.Account{
		AccountCode: username,
		AccountName: nickname,
		Email:       email,
		DepName:     department,
		OrgName:     team,
	}
}
