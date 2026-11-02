package handlers

import (
	"osm-mngt/internal/models"
	"osm-mngt/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ==================== 组织架构（部门 / 团队 / 成员） ====================
// 数据来源：users 表（username/nickname/department/team）。
// 对接外部权限系统时只需替换数据源（如改为请求 IAM 网关），保持响应结构不变，
// 前端 src/types/org.ts 与 src/api/org.ts 承接结构不变。

// OrgHandler 组织架构查询处理器
type OrgHandler struct {
	db *gorm.DB
}

func NewOrgHandler(db *gorm.DB) *OrgHandler {
	return &OrgHandler{db: db}
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
