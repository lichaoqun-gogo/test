package workflow

import (
	"osm-mngt/internal/models"

	"gorm.io/gorm"
)

// CanUserHandleNode 判断用户是否可处理指定审批节点（审批通过/驳回）
// 支持一个部门配置多名开源治理专员：
//   - 节点指定的审批人（ApproverID）可直接处理；
//   - 该节点对应部门-角色的任一在岗专员（dept_specialist_mapping 中 is_active=1）均可处理。
func CanUserHandleNode(db *gorm.DB, node *models.ApprovalNode, userID string) bool {
	if node == nil || userID == "" {
		return false
	}
	// 节点指定的审批人
	if node.ApproverID == userID {
		return true
	}

	query := db.Table("dept_specialist_mapping").
		Where("specialist_id = ? AND is_active = 1", userID)
	switch node.NodeKey {
	case "node_1":
		// 本部门开源治理专员复审
		query = query.Where("dept = ? AND specialist_role = 'dept_specialist'", node.ApproverDept)
	case "node_2":
		// 使用部门确认
		query = query.Where("dept = ? AND specialist_role = 'usage_dept_specialist'", node.ApproverDept)
	case "node_3":
		// 管理部门评审：任一在岗管理部门专员
		query = query.Where("specialist_role = 'mgmt_dept_specialist'")
	default:
		return false
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
