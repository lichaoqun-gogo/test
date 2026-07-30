package models

import "time"

// GovernanceTask 治理任务主表
type GovernanceTask struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskNo        string     `gorm:"column:task_no;size:50;uniqueIndex" json:"taskNo"`                 // 任务编号
	Type          string     `gorm:"column:type;size:50" json:"type"`                                  // 任务类型：introduce, upgrade, withdraw, defect, usage
	ActionType    string     `gorm:"column:action_type;size:50" json:"actionType"`                     // 业务类型：introduce(引入), upgrade(升级), withdraw(退出), usage(使用), risk_report(风险上报)
	Title         string     `gorm:"column:title;size:200" json:"title"`                               // 任务标题
	Applicant     string     `gorm:"column:applicant;size:100" json:"applicant"`                       // 申请人姓名
	ApplicantID   string     `gorm:"column:applicant_id;size:50" json:"applicantId"`                   // 申请人ID
	ApplicantDept string     `gorm:"column:applicant_dept;size:100" json:"applicantDept"`              // 申请人部门
	Status        string     `gorm:"column:status;size:20;default:draft" json:"status"`                // 任务状态：draft, approving, withdrawn, approved, rejected
	Priority      string     `gorm:"column:priority;size:20" json:"priority"`                          // 优先级
	CurrentStep   string     `gorm:"column:current_step;size:100" json:"currentStep"`                  // 当前步骤/节点
	CurrentNodeID int64      `gorm:"column:current_node_id" json:"currentNodeId"`                      // 当前审批节点ID
	Description   string     `gorm:"column:description;type:mediumtext" json:"description"`            // 任务描述
	Deadline      *time.Time `gorm:"column:deadline" json:"deadline"`                                  // 截止日期
	AssignedTo    string     `gorm:"column:assigned_to;size:100" json:"assignedTo"`                    // 当前处理人
	ApprovalData  string     `gorm:"column:approval_data;type:mediumtext" json:"approvalData"`         // 审批流数据JSON
	ProjectID     int64      `gorm:"column:project_id" json:"projectId"`                               // 关联的项目ID（开源软件）
	ProjectName   string     `gorm:"column:project_name;size:200" json:"projectName"`                  // 关联项目名称
	ProjectDept   string     `gorm:"column:project_dept;size:100" json:"projectDept"`                  // 项目所属部门
	CreateTime    time.Time  `gorm:"column:create_time" json:"createTime"`
	CreateBy      string     `gorm:"column:create_by;size:50" json:"createBy"`
	UpdateTime    time.Time  `gorm:"column:update_time" json:"updateTime"`
	UpdateBy      string     `gorm:"column:update_by;size:50" json:"updateBy"`
}

func (GovernanceTask) TableName() string {
	return "governance_task"
}
