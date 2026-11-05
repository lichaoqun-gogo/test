package handlers

import (
	"net/http"

	"osm-mngt/internal/config"
	"osm-mngt/internal/models"
	"osm-mngt/internal/service/permission"
	"osm-mngt/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ==================== 组织架构（部门 / 团队 / 成员） ====================
// 数据源：统一权限系统（组织接口透传，前端 src/types/org.ts 与 src/api/org.ts 承接结构不变）；
// 权限系统未对接/失败时兜底查 users 表（保持开发调试可用），对接后此路径不再生效。

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

// deptTeamRow 部门/团队聚合行（按 department + team 去重）
type deptTeamRow struct {
	Department string
	Team       string
}

// ListDepartmentsTeams 部门及其所属团队
// GET /api/v1/org/departments/teams → [{ department, teams: [...] }]
func (h *OrgHandler) ListDepartmentsTeams(c *gin.Context) {
	// 优先从统一权限系统获取（透传）
	if raw, err := h.permSvc.GetOrgDeptTeams(); err == nil && raw != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": raw})
		return
	}

	// 兜底：users 表聚合（权限系统未对接时的本地数据）
	rows := []deptTeamRow{}
	if err := h.db.Table("users").
		Select("department, team").
		Where("department <> ''").
		Group("department, team").
		Order("department ASC, team ASC").
		Scan(&rows).Error; err != nil {
		response.Error(c, 500, "查询部门团队失败")
		return
	}

	// 按部门分组（保持顺序）
	out := []DeptTeams{}
	idx := map[string]int{}
	for _, r := range rows {
		i, ok := idx[r.Department]
		if !ok {
			idx[r.Department] = len(out)
			out = append(out, DeptTeams{Department: r.Department, Teams: []string{}})
			i = len(out) - 1
		}
		if r.Team != "" {
			out[i].Teams = append(out[i].Teams, r.Team)
		}
	}
	response.Success(c, out)
}

// ListDepartments 部门 → 团队 → 成员
// GET /api/v1/org/departments → [{ department, teams: [{ team, members: [...] }] }]
func (h *OrgHandler) ListDepartments(c *gin.Context) {
	// 优先从统一权限系统获取（透传）
	if raw, err := h.permSvc.GetOrgDepartments(); err == nil && raw != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": raw})
		return
	}

	// 兜底：users 表聚合（权限系统未对接时的本地数据）
	rows := []deptTeamRow{}
	if err := h.db.Table("users").
		Select("department, team").
		Where("department <> ''").
		Group("department, team").
		Order("department ASC, team ASC").
		Scan(&rows).Error; err != nil {
		response.Error(c, 500, "查询部门团队失败")
		return
	}

	out := []OrgDepartment{}
	idx := map[string]int{}
	for _, r := range rows {
		i, ok := idx[r.Department]
		if !ok {
			idx[r.Department] = len(out)
			out = append(out, OrgDepartment{Department: r.Department, Teams: []OrgTeam{}})
			i = len(out) - 1
		}
		team := OrgTeam{Team: r.Team, Members: []OrgMember{}}
		var members []models.User
		if err := h.db.Where("department = ? AND team = ?", r.Department, r.Team).
			Order("id ASC").Find(&members).Error; err != nil {
			continue
		}
		for _, m := range members {
			team.Members = append(team.Members, OrgMember{
				Username: m.Username,
				Nickname: m.Nickname,
				Email:    m.Email,
			})
		}
		out[i].Teams = append(out[i].Teams, team)
	}
	response.Success(c, out)
}
