package handlers

import (
	"net/http"

	"osm-mngt/internal/config"
	"osm-mngt/internal/service/permission"
	"osm-mngt/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ==================== 组织架构（部门 / 团队 / 成员） ====================
// 数据源：统一权限系统（组织接口透传，前端 src/types/org.ts 与 src/api/org.ts 承接结构不变）；
// 权限系统未对接/失败时返回静态用户数据（users 表已废弃），对接后此路径不再生效。

// OrgHandler 组织架构查询处理器
type OrgHandler struct {
	db      *gorm.DB
	permSvc *permission.RemoteService // 统一权限系统转发客户端
}

func NewOrgHandler(db *gorm.DB, permCfg config.PermissionConfig) *OrgHandler {
	return &OrgHandler{db: db, permSvc: permission.NewRemoteService(permCfg)}
}

// OrgMember 团队成员
type OrgMember struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

// OrgTeam 部门下某个团队（含成员）
type OrgTeam struct {
	Team    string      `json:"team"`
	Members []OrgMember `json:"members"`
}

// OrgDepartment 部门 → 团队 → 成员
type OrgDepartment struct {
	Department string    `json:"department"`
	Teams      []OrgTeam `json:"teams"`
}

// DeptTeams 部门 → 所属团队（仅团队名列表，不含成员）
type DeptTeams struct {
	Department string   `json:"department"`
	Teams      []string `json:"teams"`
}

// ListDepartmentsTeams 部门及其所属团队
// GET /api/v1/org/departments/teams → [{ department, teams: [...] }]
func (h *OrgHandler) ListDepartmentsTeams(c *gin.Context) {
	// 优先从统一权限系统获取（透传）
	if raw, err := h.permSvc.GetOrgDeptTeams(); err == nil && raw != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": raw})
		return
	}

	// 兜底：静态用户数据（users 表已废弃），按部门分组
	out := []DeptTeams{}
	idx := map[string]int{}
	for _, u := range staticUsers {
		if u.Department == "" {
			continue
		}
		i, ok := idx[u.Department]
		if !ok {
			idx[u.Department] = len(out)
			out = append(out, DeptTeams{Department: u.Department, Teams: []string{}})
			i = len(out) - 1
		}
		if u.Team != "" && !contains(out[i].Teams, u.Team) {
			out[i].Teams = append(out[i].Teams, u.Team)
		}
	}
	response.Success(c, out)
}

// contains 判断字符串切片是否包含某元素
func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// ListDepartments 部门 → 团队 → 成员
// GET /api/v1/org/departments → [{ department, teams: [{ team, members: [...] }] }]
func (h *OrgHandler) ListDepartments(c *gin.Context) {
	// 优先从统一权限系统获取（透传）
	if raw, err := h.permSvc.GetOrgDepartments(); err == nil && raw != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": raw})
		return
	}

	// 兜底：静态用户数据（users 表已废弃），按 部门→团队→成员 构建
	out := []OrgDepartment{}
	idx := map[string]int{}
	for _, u := range staticUsers {
		if u.Department == "" {
			continue
		}
		i, ok := idx[u.Department]
		if !ok {
			idx[u.Department] = len(out)
			out = append(out, OrgDepartment{Department: u.Department, Teams: []OrgTeam{}})
			i = len(out) - 1
		}
		// 团队聚合（同一部门下相同团队合并，成员追加）
		var team *OrgTeam
		for j := range out[i].Teams {
			if out[i].Teams[j].Team == u.Team {
				team = &out[i].Teams[j]
				break
			}
		}
		if team == nil {
			team = &OrgTeam{Team: u.Team, Members: []OrgMember{}}
			out[i].Teams = append(out[i].Teams, *team)
			team = &out[i].Teams[len(out[i].Teams)-1]
		}
		team.Members = append(team.Members, OrgMember{
			Username: u.Username,
			Nickname: u.Nickname,
			Email:    u.Email,
		})
	}
	response.Success(c, out)
}
