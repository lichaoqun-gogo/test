package workflow

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"osm-mngt/internal/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

// NodeResolver 审批节点解析器
type NodeResolver struct {
	db *gorm.DB // 数据库连接（用于查询部门专员映射）
}

// SpecialistInfo 专员信息
type SpecialistInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Dept string `json:"dept"`
}

// NewNodeResolver 创建节点解析器
func NewNodeResolver(db *gorm.DB) *NodeResolver {
	return &NodeResolver{
		db: db,
	}
}

// getDeptSpecialist 获取部门的开源治理专员（从数据库查询）
func (r *NodeResolver) getDeptSpecialist(dept string) (SpecialistInfo, bool) {
	var specialist struct {
		SpecialistID   string `gorm:"column:specialist_id"`
		SpecialistName string `gorm:"column:specialist_name"`
		Dept           string `gorm:"column:dept"`
	}
	
	err := r.db.Table("dept_specialist_mapping").
		Where("dept = ? AND specialist_role = 'dept_specialist' AND is_active = 1", dept).
		First(&specialist).Error
	
	if err != nil {
		return SpecialistInfo{}, false
	}
	
	return SpecialistInfo{
		ID:   specialist.SpecialistID,
		Name: specialist.SpecialistName,
		Dept: specialist.Dept,
	}, true
}

// getUsageDeptSpecialist 获取使用部门的开源治理专员（从数据库查询）
func (r *NodeResolver) getUsageDeptSpecialist(dept string) (SpecialistInfo, bool) {
	var specialist struct {
		SpecialistID   string `gorm:"column:specialist_id"`
		SpecialistName string `gorm:"column:specialist_name"`
		Dept           string `gorm:"column:dept"`
	}
	
	err := r.db.Table("dept_specialist_mapping").
		Where("dept = ? AND specialist_role = 'usage_dept_specialist' AND is_active = 1", dept).
		First(&specialist).Error
	
	if err != nil {
		return SpecialistInfo{}, false
	}
	
	return SpecialistInfo{
		ID:   specialist.SpecialistID,
		Name: specialist.SpecialistName,
		Dept: specialist.Dept,
	}, true
}

// getMgmtDeptSpecialists 获取管理部门专员列表（从数据库查询）
func (r *NodeResolver) getMgmtDeptSpecialists() []SpecialistInfo {
	var specialists []struct {
		SpecialistID   string `gorm:"column:specialist_id"`
		SpecialistName string `gorm:"column:specialist_name"`
		Dept           string `gorm:"column:dept"`
	}
	
	r.db.Table("dept_specialist_mapping").
		Where("specialist_role = 'mgmt_dept_specialist' AND is_active = 1").
		Find(&specialists)
	
	result := make([]SpecialistInfo, 0, len(specialists))
	for _, s := range specialists {
		result = append(result, SpecialistInfo{
			ID:   s.SpecialistID,
			Name: s.SpecialistName,
			Dept: s.Dept,
		})
	}
	
	return result
}

// ResolveNodes 根据任务信息解析出实际的审批节点列表
// 参数：
//   - actionType: 业务类型（introduce, upgrade, withdraw, usage, risk_report）
//   - applicantID: 申请人ID
//   - applicantDept: 申请人部门
//   - projectDept: 项目所属部门（运维部门）
//
// 返回：实际需要经过的审批节点列表
func (r *NodeResolver) ResolveNodes(actionType, applicantID, applicantDept, projectDept string) ([]models.ApprovalNode, error) {
	fmt.Printf("[NodeResolver.ResolveNodes] actionType=%s, applicantID=%s, applicantDept=%s, projectDept=%s\n",
		actionType, applicantID, applicantDept, projectDept)

	// 获取模板节点定义
	templateNodes := r.getTemplateNodes()

	var nodes []models.ApprovalNode
	nodeOrder := 0

	for _, tmplNode := range templateNodes {
		// Node 0: 发起节点（固定）
		if tmplNode.NodeKey == "node_0" {
			nodes = append(nodes, models.ApprovalNode{
				NodeKey:      tmplNode.NodeKey,
				NodeName:     tmplNode.NodeName,
				NodeOrder:    nodeOrder,
				NodeType:     tmplNode.NodeType,
				ApproverID:   applicantID,
				Approver:     "发起人",
				ApproverDept: applicantDept,
				Status:       "pending",
				IsSkipped:    false,
			})
			nodeOrder++
			continue
		}

		// Node 1: 本部门开源治理专员复审（条件节点）
		if tmplNode.NodeKey == "node_1" {
			shouldSkip, skipReason := r.shouldSkipNode1(applicantID, applicantDept)
			if shouldSkip {
				// 跳过该节点
				nodes = append(nodes, models.ApprovalNode{
					NodeKey:      tmplNode.NodeKey,
					NodeName:     tmplNode.NodeName,
					NodeOrder:    nodeOrder,
					NodeType:     tmplNode.NodeType,
					Status:       "skipped",
					IsSkipped:    true,
					SkipReason:   skipReason,
				})
				nodeOrder++
				continue
			}

			// 不跳过，查找本部门专员
			specialist, exists := r.getDeptSpecialist(applicantDept)
			if !exists {
				return nil, fmt.Errorf("未找到部门 %s 的开源治理专员", applicantDept)
			}

			nodes = append(nodes, models.ApprovalNode{
				NodeKey:      tmplNode.NodeKey,
				NodeName:     tmplNode.NodeName,
				NodeOrder:    nodeOrder,
				NodeType:     tmplNode.NodeType,
				ApproverID:   specialist.ID,
				Approver:     specialist.Name,
				ApproverDept: specialist.Dept,
				Status:       "pending",
				IsSkipped:    false,
			})
			nodeOrder++
			continue
		}

		// Node 2: 使用部门确认（条件节点，仅升级/退出申请激活）
		// 支持多个使用部门，每个使用部门创建独立审批节点，自动跳过与运维部门相同的
		if tmplNode.NodeKey == "node_2" {
			// 仅升级/退出申请需要Node 2
			if actionType != "upgrade" && actionType != "withdraw" {
				nodes = append(nodes, models.ApprovalNode{
					NodeKey:      tmplNode.NodeKey,
					NodeName:     tmplNode.NodeName,
					NodeOrder:    nodeOrder,
					NodeType:     tmplNode.NodeType,
					Status:       "skipped",
					IsSkipped:    true,
					SkipReason:   fmt.Sprintf("业务类型 %s 不需要使用部门确认，自动跳过", actionType),
				})
				nodeOrder++
				continue
			}

			// 解析使用部门列表（逗号分隔）
			usageDepts := parseDeptList(projectDept)
			if len(usageDepts) == 0 {
				nodes = append(nodes, models.ApprovalNode{
					NodeKey:      tmplNode.NodeKey,
					NodeName:     tmplNode.NodeName,
					NodeOrder:    nodeOrder,
					NodeType:     tmplNode.NodeType,
					Status:       "skipped",
					IsSkipped:    true,
					SkipReason:   "未指定使用部门，自动跳过",
				})
				nodeOrder++
				continue
			}

			// 过滤掉与运维部门（applicantDept）相同的使用部门
			var activeDepts []string
			for _, dept := range usageDepts {
				if dept == applicantDept {
					fmt.Printf("[NodeResolver] 使用部门 %s 与运维部门相同，跳过\n", dept)
					continue
				}
				activeDepts = append(activeDepts, dept)
			}

			if len(activeDepts) == 0 {
				// 所有使用部门都与运维部门相同，跳过整个Node 2
				nodes = append(nodes, models.ApprovalNode{
					NodeKey:      tmplNode.NodeKey,
					NodeName:     tmplNode.NodeName,
					NodeOrder:    nodeOrder,
					NodeType:     tmplNode.NodeType,
					Status:       "skipped",
					IsSkipped:    true,
					SkipReason:   "所有使用部门与运维部门相同，自动跳过（防重复）",
				})
				nodeOrder++
				continue
			}

			// 随机打乱使用部门顺序
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(activeDepts), func(i, j int) {
				activeDepts[i], activeDepts[j] = activeDepts[j], activeDepts[i]
			})

			// 为每个使用部门创建独立的审批节点；没有专员的部门直接跳过
			createdCount := 0
			for _, dept := range activeDepts {
				specialist, exists := r.getUsageDeptSpecialist(dept)
				if !exists {
					fmt.Printf("[NodeResolver] 使用部门 %s 未配置开源治理专员，跳过该部门确认节点\n", dept)
					continue
				}

				nodes = append(nodes, models.ApprovalNode{
					NodeKey:      tmplNode.NodeKey,
					NodeName:     fmt.Sprintf("%s（%s）", tmplNode.NodeName, dept),
					NodeOrder:    nodeOrder,
					NodeType:     tmplNode.NodeType,
					ApproverID:   specialist.ID,
					Approver:     specialist.Name,
					ApproverDept: specialist.Dept,
					Status:       "pending",
					IsSkipped:    false,
				})
				nodeOrder++
				createdCount++
			}

			// 如果所有使用部门都没有配置专员，则整体跳过 Node 2
			if createdCount == 0 {
				nodes = append(nodes, models.ApprovalNode{
					NodeKey:      tmplNode.NodeKey,
					NodeName:     tmplNode.NodeName,
					NodeOrder:    nodeOrder,
					NodeType:     tmplNode.NodeType,
					Status:       "skipped",
					IsSkipped:    true,
					SkipReason:   "所有使用部门均未配置开源治理专员，自动跳过",
				})
				nodeOrder++
			}
			continue
		}

		// Node 3: 管理部门评审（固定节点）
		if tmplNode.NodeKey == "node_3" {
			// 从数据库查询管理部门专员
			specialists := r.getMgmtDeptSpecialists()
			if len(specialists) == 0 {
				return nil, fmt.Errorf("未找到管理部门专员")
			}
			
			// 使用第一个管理部门专员
			specialist := specialists[0]
			nodes = append(nodes, models.ApprovalNode{
				NodeKey:      tmplNode.NodeKey,
				NodeName:     tmplNode.NodeName,
				NodeOrder:    nodeOrder,
				NodeType:     tmplNode.NodeType,
				ApproverID:   specialist.ID,
				Approver:     specialist.Name,
				ApproverDept: specialist.Dept,
				Status:       "pending",
				IsSkipped:    false,
			})
			nodeOrder++
			continue
		}
	}

	return nodes, nil
}

// shouldSkipNode1 判断是否跳过Node 1（本部门专员复审）
// 防重复机制：若发起人ID == 本部门开源治理专员ID，则自动跳过
func (r *NodeResolver) shouldSkipNode1(applicantID, applicantDept string) (bool, string) {
	// 从数据库查询部门专员
	var specialist struct {
		SpecialistID   string `gorm:"column:specialist_id"`
		SpecialistName string `gorm:"column:specialist_name"`
	}
	err := r.db.Table("dept_specialist_mapping").
		Select("specialist_id, specialist_name").
		Where("dept = ? AND specialist_role = 'dept_specialist' AND is_active = 1", applicantDept).
		First(&specialist).Error
	
	if err != nil {
		// 如果部门没有专员，也跳过
		return true, fmt.Sprintf("部门 %s 无开源治理专员，自动跳过", applicantDept)
	}

	if applicantID == specialist.SpecialistID {
		return true, fmt.Sprintf("发起人即为部门专员，自动跳过（防重复）")
	}

	return false, ""
}

// parseDeptList 解析部门列表（逗号分隔的字符串转为切片）
func parseDeptList(depts string) []string {
	if depts == "" {
		return nil
	}
	parts := strings.Split(depts, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// getTemplateNodes 获取模板节点定义
func (r *NodeResolver) getTemplateNodes() []models.ApprovalTemplateNode {
	// 这里返回万能模板的节点定义
	// 实际项目中可以从数据库加载，这里先硬编码
	return []models.ApprovalTemplateNode{
		{
			NodeKey:      "node_0",
			NodeName:     "发起节点",
			NodeOrder:    0,
			NodeType:     "fixed",
			ApproverRole: "initiator",
			Description:  "任务发起人提交申请",
		},
		{
			NodeKey:       "node_1",
			NodeName:      "本部门开源治理专员复审",
			NodeOrder:     1,
			NodeType:      "conditional",
			ConditionType: "dept_specialist",
			IsSkipable:    true,
			SkipRule:      marshalSkipRule("same_person", "approver_id"),
			ApproverRole:  "dept_specialist",
			Description:  "发起人所在部门的开源治理专员进行复审",
		},
		{
			NodeKey:       "node_2",
			NodeName:      "使用部门确认",
			NodeOrder:     2,
			NodeType:      "conditional",
			ConditionType: "usage_dept",
			IsSkipable:    true,
			SkipRule:      marshalSkipRule("same_dept", "project_dept"),
			ApproverRole:  "usage_dept_specialist",
			Description:  "使用部门的开源治理专员确认（仅升级/退出申请）",
		},
		{
			NodeKey:      "node_3",
			NodeName:     "管理部门评审",
			NodeOrder:    3,
			NodeType:     "fixed",
			ApproverRole: "mgmt_dept_specialist",
			Description:  "开源治理办公室进行最终评审",
		},
	}
}

// marshalSkipRule 将跳过规则序列化为JSON
func marshalSkipRule(ruleType, field string) string {
	rule := map[string]string{
		"type":  ruleType,
		"field": field,
	}
	data, _ := json.Marshal(rule)
	return string(data)
}
