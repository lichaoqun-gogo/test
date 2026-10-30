package models

import "time"

// DefectDisposal 缺陷处置工单（子记录）
// 一个缺陷（父）可裂变为多张处置工单：一张处置单 = 一个缺陷 × 一个涉及项目（及所属部门）
// 同一部门有多个涉及项目时生成多张处置单（如组件缺陷被某部门多个自研项目依赖）
// 各部门独立处置，全部完成后缺陷闭环
// 按缺陷编号（defect_no）关联缺陷实体，不依赖实体表主键（缺陷实体可反复导入 upsert）
type DefectDisposal struct {
	ID               int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DefectNo         string     `gorm:"column:defect_no;size:100;not null;index:idx_defect_no" json:"defectNo"`
	DisposalNo       string     `gorm:"column:disposal_no;size:50" json:"disposalNo"`
	DeptName         string     `gorm:"column:dept_name;size:128" json:"deptName"`                 // 处置部门
	DisposalType     string     `gorm:"column:disposal_type;size:20" json:"disposalType"`          // maint-运维部门, use-使用部门
	ProjectID        int64      `gorm:"column:project_id;default:0" json:"projectId"`              // 涉及项目ID（软件=project.id，组件=project_inner.id）；0=旧数据未记录项目
	ProjectName      string     `gorm:"column:project_name;size:128" json:"projectName"`           // 涉及项目显示名（软件=_name，组件=pro_code）
	ProjectType      string     `gorm:"column:project_type;size:20" json:"projectType"`            // 涉及项目类型: software-基础软件, inner-自研项目（供按项目匹配负责人）
	AssigneeUsername string     `gorm:"column:assignee_username;size:64;index:idx_assignee" json:"assigneeUsername"` // 处置人账号
	AssigneeName     string     `gorm:"column:assignee_name;size:64" json:"assigneeName"`          // 处置人姓名
	TicketNo         string     `gorm:"column:ticket_no;size:64" json:"ticketNo"`                 // 工单号
	DisposalNote     string     `gorm:"column:disposal_note;type:text" json:"disposalNote"`        // 处置备注
	Deadline         *time.Time `gorm:"column:deadline" json:"deadline"`                           // 截止时间（继承缺陷，不可改）
	Status           string     `gorm:"column:status;size:20;default:pending_assign" json:"status"` // pending_assign-待分配, processing-处理中, completed-已完成, na-不适用
	AssignedBy       string     `gorm:"column:assigned_by;size:64" json:"assignedBy"`              // 二次分配人
	AssignedAt       *time.Time `gorm:"column:assigned_at" json:"assignedAt"`
	CompletedAt      *time.Time `gorm:"column:completed_at" json:"completedAt"`
	CreateTime       time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime       time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (DefectDisposal) TableName() string {
	return "defect_disposal"
}
