package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"osm-mngt/internal/models"
	"osm-mngt/pkg/response"
	"gorm.io/gorm"
)

// ProjectOwnerHandler 项目-负责人关联处理器
type ProjectOwnerHandler struct {
	db *gorm.DB
}

func NewProjectOwnerHandler(db *gorm.DB) *ProjectOwnerHandler {
	return &ProjectOwnerHandler{db: db}
}

// OwnerInfo 负责人信息（含冗余部门/团队）
type OwnerInfo struct {
	OwnerUsername string `json:"owner_username"`
	DeptName      string `json:"dept_name"`
	Team          string `json:"team"`
}

// GetOwnersByProject 查询某项目的所有负责人
// GET /api/v1/asset/project-owners?project_type=software&project_id=1
func (h *ProjectOwnerHandler) GetOwnersByProject(c *gin.Context) {
	projectType := c.Query("project_type")
	projectIDStr := c.Query("project_id")
	if projectType == "" || projectIDStr == "" {
		response.Error(c, http.StatusBadRequest, "缺少 project_type 或 project_id")
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "project_id 参数无效")
		return
	}

	var owners []models.ProjectOwner
	if err := h.db.Where("project_type = ? AND project_id = ?", projectType, projectID).
		Find(&owners).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "查询负责人失败")
		return
	}

	response.Success(c, owners)
}

// SaveOwners 批量保存项目负责人（全量替换）
// POST /api/v1/asset/project-owners
// Body: { project_type, project_id, owners: [{owner_username, dept_name, team}] }
func (h *ProjectOwnerHandler) SaveOwners(c *gin.Context) {
	var req struct {
		ProjectType string      `json:"project_type" binding:"required"`
		ProjectID   int64       `json:"project_id" binding:"required"`
		Owners      []OwnerInfo `json:"owners"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 从 header 获取分配人（X-User-Name 已由前端 encode）
	assignedBy := c.GetHeader("X-User-Name")

	// 事务：先删除旧的，再插入新的
	tx := h.db.Begin()
	if err := tx.Where("project_type = ? AND project_id = ?", req.ProjectType, req.ProjectID).
		Delete(&models.ProjectOwner{}).Error; err != nil {
		tx.Rollback()
		response.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	now := time.Now()
	for _, o := range req.Owners {
		if o.OwnerUsername == "" {
			continue
		}
		rec := models.ProjectOwner{
			ProjectType:   req.ProjectType,
			ProjectID:     req.ProjectID,
			OwnerUsername: o.OwnerUsername,
			DeptName:      o.DeptName,
			Team:          o.Team,
			AssignedBy:    assignedBy,
			AssignedAt:    now,
		}
		if err := tx.Create(&rec).Error; err != nil {
			tx.Rollback()
			response.Error(c, http.StatusInternalServerError, "保存负责人失败: "+err.Error())
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "提交事务失败")
		return
	}

	response.Success(c, gin.H{"message": "保存成功"})
}

// GetOwnersMap 批量查询多个项目的负责人（用于列表页一次性加载）
// GET /api/v1/asset/project-owners/map?project_type=software&project_ids=1,2,3
func (h *ProjectOwnerHandler) GetOwnersMap(c *gin.Context) {
	projectType := c.Query("project_type")
	projectIDsStr := c.Query("project_ids")
	if projectType == "" || projectIDsStr == "" {
		response.Success(c, gin.H{})
		return
	}

	// 解析 project_ids
	idStrs := splitAndTrim(projectIDsStr)
	if len(idStrs) == 0 {
		response.Success(c, gin.H{})
		return
	}

	ids := make([]int64, 0, len(idStrs))
	for _, s := range idStrs {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		response.Success(c, gin.H{})
		return
	}

	var owners []models.ProjectOwner
	if err := h.db.Where("project_type = ? AND project_id IN ?", projectType, ids).
		Find(&owners).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "查询负责人失败")
		return
	}

	// map[project_id] -> []owner_username
	result := make(map[int64][]string)
	for _, o := range owners {
		result[o.ProjectID] = append(result[o.ProjectID], o.OwnerUsername)
	}

	response.Success(c, result)
}

func splitAndTrim(s string) []string {
	result := []string{}
	current := ""
	for _, r := range s {
		if r == ',' {
			if current != "" {
				result = append(result, current)
			}
			current = ""
		} else if r != ' ' {
			current += string(r)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
