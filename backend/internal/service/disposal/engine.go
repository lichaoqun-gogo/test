package disposal

import (
	"fmt"
	"strings"
	"time"

	"osm-mngt/internal/models"

	"gorm.io/gorm"
)

// 父缺陷状态（缺陷级处置流程状态，存 defect_flow.status）
const (
	StatusPendingAnalysis = "pending_analysis" // 待研判
	StatusPendingReview   = "pending_review"   // 待评审（待上会）
	StatusAssigned        = "assigned"         // 已分配（已裂变处置）
	StatusIgnored         = "ignored"          // 已忽略
	StatusClosed          = "closed"           // 已闭环
)

// 子处置单状态（与 defect_disposal.status 一致）
const (
	DisposalPendingAssign = "pending_assign" // 待分配（待专员二次分配）
	DisposalProcessing    = "processing"     // 处理中
	DisposalCompleted     = "completed"      // 已完成
	DisposalNA            = "na"             // 不适用
)

// 研判结论
const (
	ConclusionNeedReview = "need_review" // 需要上会评审
	ConclusionNeedFix    = "need_fix"    // 需要修复（漏洞可直接分配）
	ConclusionIgnore     = "ignore"      // 忽略
)

// 评审结论
const (
	ReviewNeedFix  = "need_fix"  // 需要修复 -> 裂变处置
	ReviewNoAction = "no_action" // 无需处置 -> 忽略
)

// 处置类型
const (
	DisposalTypeMaint = "maint" // 运维部门
	DisposalTypeUse   = "use"   // 使用部门
)

// Engine 缺陷/漏洞处置引擎（并行处置，与串行审批流独立）
type Engine struct {
	db *gorm.DB
}

func NewEngine(db *gorm.DB) *Engine {
	return &Engine{db: db}
}

// defectCtx 缺陷处置上下文（流程明细 + 实体属性 + 内部处置状态），引擎统一按缺陷编号（defect_no）加载
// 流程明细来自 defect_flow；内部处置状态（internal_status）与影响组件来自缺陷实体表（内部 defect_ledger / 外部 defect）
type defectCtx struct {
	DefectNo          string
	FlowType          string
	MaintDept         string
	Deadline          *time.Time
	Status            string
	DispatchTarget    string
	AffectedComponent string
}

// loadDefectCtx 按缺陷编号加载处置上下文（流程明细 + 实体属性 + 内部处置状态）
func (e *Engine) loadDefectCtx(tx *gorm.DB, defectNo string) (*defectCtx, error) {
	var flow models.DefectFlow
	if err := tx.Where("defect_no = ?", defectNo).First(&flow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("缺陷不存在")
		}
		return nil, err
	}
	ctx := &defectCtx{
		DefectNo:       flow.DefectNo,
		FlowType:       flow.FlowType,
		MaintDept:      flow.MaintDept,
		Deadline:       flow.Deadline,
		DispatchTarget: flow.DispatchTarget,
	}
	// 实体属性 + 内部处置状态：内部缺陷实体（风险上报，漏洞携带影响组件供自动分配裂变）；外部缺陷实体（软件名作为影响组件）
	var ledger models.DefectLedger
	if err := tx.Where("defect_no = ?", defectNo).First(&ledger).Error; err == nil {
		ctx.AffectedComponent = strings.TrimSpace(ledger.AffectedComponent)
		ctx.Status = ledger.InternalStatus
	} else {
		var ext models.Defect
		if err := tx.Where("defect_id = ?", defectNo).First(&ext).Error; err == nil {
			ctx.AffectedComponent = ext.Software
			ctx.Status = ext.InternalStatus
		}
	}
	if ctx.Status == "" {
		return nil, fmt.Errorf("缺陷不存在")
	}
	return ctx, nil
}

// updateStatus 更新缺陷实体表的内部处置状态（内部实体 defect_ledger / 外部实体 defect）
func (e *Engine) updateStatus(tx *gorm.DB, defectNo, status string) error {
	res := tx.Model(&models.DefectLedger{}).Where("defect_no = ?", defectNo).
		Updates(map[string]interface{}{"internal_status": status, "update_time": time.Now()})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return tx.Model(&models.Defect{}).Where("defect_id = ?", defectNo).
			Updates(map[string]interface{}{"internal_status": status}).Error
	}
	return nil
}

type specialist struct {
	SpecialistID   string `gorm:"column:specialist_id"`
	SpecialistName string `gorm:"column:specialist_name"`
	Dept           string `gorm:"column:dept"`
}

// isSpecialistOf 判断用户是否为指定部门的专员（用姓名匹配，兼容 CAS/local 两种模式）
func (e *Engine) isSpecialistOf(userID, userName, dept string, roles ...string) bool {
	if dept == "" {
		return false
	}
	roleList := roles
	if len(roleList) == 0 {
		roleList = []string{"dept_specialist"}
	}
	var count int64
	e.db.Table("dept_specialist_mapping").
		Where("dept = ? AND specialist_role IN ? AND is_active = 1 AND (specialist_id = ? OR specialist_name = ?)",
			dept, roleList, userID, userName).
		Count(&count)
	return count > 0
}

// isMgmtSpecialist 判断是否为管理部门专员
func (e *Engine) isMgmtSpecialist(userID, userName string) bool {
	var count int64
	e.db.Table("dept_specialist_mapping").
		Where("specialist_role = 'mgmt_dept_specialist' AND is_active = 1 AND (specialist_id = ? OR specialist_name = ?)",
			userID, userName).
		Count(&count)
	return count > 0
}

// IsMgmtSpecialist 导出：判断是否为管理部门专员
func (e *Engine) IsMgmtSpecialist(userID, userName string) bool {
	return e.isMgmtSpecialist(userID, userName)
}

// SpecialistDepts 导出：返回用户作为部门专员所负责的所有部门
func (e *Engine) SpecialistDepts(userID, userName string) []string {
	return e.specialistDepts(userID, userName)
}

// specialistDepts 返回用户作为部门专员所负责的所有部门
func (e *Engine) specialistDepts(userID, userName string) []string {
	var rows []specialist
	e.db.Table("dept_specialist_mapping").
		Where("specialist_role IN ('dept_specialist','usage_dept_specialist') AND is_active = 1 AND (specialist_id = ? OR specialist_name = ?)",
			userID, userName).
		Find(&rows)
	seen := map[string]bool{}
	var depts []string
	for _, r := range rows {
		if r.Dept == "" || seen[r.Dept] {
			continue
		}
		seen[r.Dept] = true
		depts = append(depts, r.Dept)
	}
	return depts
}

// canOperateDefect 判断用户是否可对缺陷执行研判/分配等操作
func (e *Engine) canOperateDefect(userID, userName, role string, defect *defectCtx) bool {
	if role == "admin" {
		return true
	}
	if e.isMgmtSpecialist(userID, userName) {
		return true
	}
	return e.isSpecialistOf(userID, userName, defect.MaintDept, "dept_specialist")
}

// CanSeeDefect 判断用户能否查看该缺陷（缺陷台账可见性，与 ListDefects 的过滤规则一致）：
//   - admin / 管理部门专员：可见全部（全景）
//   - 部门专员：本部门运维（maint_dept）的缺陷，或本部门有处置单的缺陷（跨部门裂变后使用部门专员也可看到，以便展开改派本部门处置单）
func (e *Engine) CanSeeDefect(userID, userName, role, userDept string, defect *defectCtx) bool {
	if role == "admin" || e.isMgmtSpecialist(userID, userName) {
		return true
	}
	if userDept == "" {
		return false
	}
	if defect.MaintDept == userDept {
		return true
	}
	var n int64
	e.db.Table("defect_disposal").Where("defect_no = ? AND dept_name = ?", defect.DefectNo, userDept).Count(&n)
	return n > 0
}

// canReviewDefect 判断用户是否可评审（上会）该缺陷：仅 admin 或管理部门专员（全局治理）
func (e *Engine) canReviewDefect(userID, userName, role string) bool {
	if role == "admin" {
		return true
	}
	return e.isMgmtSpecialist(userID, userName)
}

// canAnalyzeDefect 判断用户是否可研判该缺陷：仅被分配的研判人本人（admin 除外），专员不再代研判（谁研判就是谁）
func (e *Engine) canAnalyzeDefect(userID, userName, role string, defect *defectCtx) bool {
	if role == "admin" {
		return true
	}
	return defect.DispatchTarget != "" &&
		(defect.DispatchTarget == userID || defect.DispatchTarget == userName)
}

// ==================== 负责人自动匹配（研判分配 / 处置人分配共用） ====================

// OwnerCandidate 负责人候选（自动匹配结果，列表第一个为主岗）
type OwnerCandidate struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	Dept      string `json:"dept"`
	Team      string `json:"team"`
	ProjectID int64  `json:"projectId"`
	IsPrimary bool   `json:"isPrimary"` // 是否主岗（多个负责人取第一个）
}

// queryOwnerCandidates 查询项目（基础软件 software / 自研项目 inner）的负责人列表，按配置顺序返回（第一个为主岗）
func queryOwnerCandidates(db *gorm.DB, projectType string, projectID int64) []OwnerCandidate {
	var rows []struct {
		OwnerUsername string `gorm:"column:owner_username"`
		Nickname      string `gorm:"column:nickname"`
		DeptName      string `gorm:"column:dept_name"`
		Team          string `gorm:"column:team"`
	}
	if err := db.Table("project_owner AS po").
		Select("po.owner_username, IFNULL(u.nickname, po.owner_username) AS nickname, po.dept_name, po.team").
		Joins("LEFT JOIN users u ON u.username = po.owner_username").
		Where("po.project_type = ? AND po.project_id = ?", projectType, projectID).
		Order("po.id ASC").
		Find(&rows).Error; err != nil {
		return nil
	}
	out := make([]OwnerCandidate, 0, len(rows))
	for _, r := range rows {
		out = append(out, OwnerCandidate{
			Username:  r.OwnerUsername,
			Name:      r.Nickname,
			Dept:      r.DeptName,
			Team:      r.Team,
			ProjectID: projectID,
		})
	}
	return out
}

// MatchCandidates 按缺陷影响对象自动匹配候选负责人（研判分配、处置单自动匹配共用）
//   - 开源基础软件（affected_component 命中 project 表，如 redis）：
//     取该软件的 project_owner（software），多个取第一个为主岗；
//   - 开源组件（project 表未命中，如 fastjson）：
//     查 component_snapshot 中依赖该组件的项目（PROJECT_ID 关联 project_inner 自研项目），
//     取各项目的 project_owner（inner），整体第一个为主岗
func (e *Engine) MatchCandidates(defect *defectCtx) []OwnerCandidate {
	return matchCandidates(e.db, defect)
}

// MatchCandidatesByNo 按缺陷编号加载上下文后返回匹配候选（研判候选接口用）
func (e *Engine) MatchCandidatesByNo(defectNo string) []OwnerCandidate {
	ctx, err := e.loadDefectCtx(e.db, defectNo)
	if err != nil {
		return nil
	}
	return matchCandidates(e.db, ctx)
}

// CanSeeDefectByNo 按缺陷编号判断用户可见性（处置单列表等用）
func (e *Engine) CanSeeDefectByNo(userID, userName, role, userDept, defectNo string) bool {
	ctx, err := e.loadDefectCtx(e.db, defectNo)
	if err != nil {
		return false
	}
	return e.CanSeeDefect(userID, userName, role, userDept, ctx)
}

func matchCandidates(db *gorm.DB, defect *defectCtx) []OwnerCandidate {
	component := strings.TrimSpace(defect.AffectedComponent)
	if component == "" {
		return nil
	}

	seen := map[string]bool{}
	var candidates []OwnerCandidate

	// 1) 开源基础软件：project 表直接命中
	var prj models.Project
	if err := db.Where("_name = ?", component).Order("id ASC").First(&prj).Error; err == nil {
		for _, c := range queryOwnerCandidates(db, "software", prj.ID) {
			if !seen[c.Username] {
				seen[c.Username] = true
				candidates = append(candidates, c)
			}
		}
		if len(candidates) > 0 {
			candidates[0].IsPrimary = true
		}
		return candidates
	}

	// 2) 开源组件：查 component_snapshot 依赖该组件的项目
	var snaps []models.ComponentSnapshot
	if err := db.Where("_name = ?", component).Find(&snaps).Error; err != nil {
		return nil
	}
	seenPID := map[int64]bool{}
	for _, s := range snaps {
		if s.ProjectID == nil || *s.ProjectID == 0 || seenPID[*s.ProjectID] {
			continue
		}
		seenPID[*s.ProjectID] = true
		for _, c := range queryOwnerCandidates(db, "inner", *s.ProjectID) {
			if !seen[c.Username] {
				seen[c.Username] = true
				candidates = append(candidates, c)
			}
		}
	}
	if len(candidates) > 0 {
		candidates[0].IsPrimary = true
	}
	return candidates
}

// AssignAnalyst 第一次分配：把待研判缺陷分配给软件负责人研判（pending_analysis 阶段，可重复改派）
// 分配信息写 defect_flow 的 dispatch_* 字段：dispatch_type='analyze'、dispatch_target=研判人 username
func (e *Engine) AssignAnalyst(defectNo string, analystUsername, userID, userName, role string) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		defect, err := e.loadDefectCtx(tx, defectNo)
		if err != nil {
			return err
		}

		if defect.Status != StatusPendingAnalysis {
			return fmt.Errorf("当前状态 %s 不允许分配研判人", defect.Status)
		}
		// 权限：专员（admin/管理部门专员/维护部门专员）
		if !e.canOperateDefect(userID, userName, role, defect) {
			return fmt.Errorf("您无权分配该缺陷的研判人")
		}
		if analystUsername == "" {
			return fmt.Errorf("请指定研判人")
		}

		// 研判人必须是系统用户
		var analyst models.User
		if err := tx.Where("username = ?", analystUsername).First(&analyst).Error; err != nil {
			return fmt.Errorf("研判人 %s 不是系统用户", analystUsername)
		}

		now := time.Now()
		if err := tx.Model(&models.DefectFlow{}).Where("defect_no = ?", defectNo).Updates(map[string]interface{}{
			"dispatch_type":   "analyze",
			"dispatch_target": analyst.Username,
			"dispatch_dept":   analyst.Department,
			"dispatch_time":   now,
			"dispatcher":      userID,
			"update_time":     now,
		}).Error; err != nil {
			return err
		}

		return e.addHistory(tx, defectNo, 0, StatusPendingAnalysis, StatusPendingAnalysis,
			"assign_analyst", userName, "分配研判人 "+analyst.Nickname)
	})
}

// AutoAssignAnalyst 自动分配研判人：取自动匹配候选的第一个（主岗）
func (e *Engine) AutoAssignAnalyst(defectNo string, userID, userName, role string) (string, error) {
	return e.autoAssignAnalystTx(e.db, defectNo, userID, userName, role)
}

func (e *Engine) autoAssignAnalystTx(tx *gorm.DB, defectNo string, userID, userName, role string) (string, error) {
	defect, err := e.loadDefectCtx(tx, defectNo)
	if err != nil {
		return "", err
	}
	if defect.Status != StatusPendingAnalysis {
		return "", fmt.Errorf("当前状态 %s 不允许分配研判人", defect.Status)
	}
	if !e.canOperateDefect(userID, userName, role, defect) {
		return "", fmt.Errorf("您无权分配该缺陷的研判人")
	}

	// 自动匹配：候选第一个为主岗
	candidates := matchCandidates(tx, defect)
	if len(candidates) == 0 {
		return "", fmt.Errorf("未匹配到候选负责人（基础软件查 project_owner，组件查依赖项目 owner），请手动分配")
	}
	primary := candidates[0]

	var analyst models.User
	if err := tx.Where("username = ?", primary.Username).First(&analyst).Error; err != nil {
		return "", fmt.Errorf("研判人 %s 不是系统用户", primary.Username)
	}

	now := time.Now()
	if err := tx.Model(&models.DefectFlow{}).Where("defect_no = ?", defectNo).Updates(map[string]interface{}{
		"dispatch_type":   "analyze",
		"dispatch_target": analyst.Username,
		"dispatch_dept":   analyst.Department,
		"dispatch_time":   now,
		"dispatcher":      userID,
		"update_time":     now,
	}).Error; err != nil {
		return "", err
	}
	if err := e.addHistory(tx, defectNo, 0, StatusPendingAnalysis, StatusPendingAnalysis,
		"auto_assign_analyst", userName, "系统自动分配研判人 "+analyst.Nickname+"（主岗）"); err != nil {
		return "", err
	}
	return analyst.Nickname, nil
}

// AutoAssignDisposals 自动分配处置单：对该缺陷全部「待处置」且未指定处理人的处置单，
// 逐张匹配该涉及项目的主岗（project_owner 第一个），匹配到则记录处理人（工单仍为「待处置」）
// 返回：成功分配的部门列表（去重）、未匹配到主岗的部门列表
func (e *Engine) AutoAssignDisposals(defectNo string, userID, userName, role string) (assigned []string, missed []string, err error) {
	return e.autoAssignDisposalsTx(e.db, defectNo, userID, userName, role)
}

func (e *Engine) autoAssignDisposalsTx(tx *gorm.DB, defectNo string, userID, userName, role string) ([]string, []string, error) {
	defect, err := e.loadDefectCtx(tx, defectNo)
	if err != nil {
		return nil, nil, err
	}
	// 权限：专员
	if !e.canOperateDefect(userID, userName, role, defect) {
		return nil, nil, fmt.Errorf("您无权分配该缺陷的处置单")
	}
	return e.autoAssignDisposalsCore(tx, defectNo, userName)
}

// autoAssignDisposalsCore 自动匹配处理人（不校验操作权限，供研判/评审裂变时系统自动调用）
// 按「涉及项目」匹配：每张处置单取该项目 project_owner 的主岗（第一个），并要求负责人部门与处置部门一致；
// 旧数据（project_id=0 未记录项目）回退为按处置部门匹配主岗
func (e *Engine) autoAssignDisposalsCore(tx *gorm.DB, defectNo string, operatorName string) ([]string, []string, error) {
	defect, err := e.loadDefectCtx(tx, defectNo)
	if err != nil {
		return nil, nil, err
	}

	var disposals []models.DefectDisposal
	if err := tx.Where("defect_no = ? AND status = ?", defectNo, DisposalPendingAssign).
		Find(&disposals).Error; err != nil {
		return nil, nil, err
	}
	if len(disposals) == 0 {
		return nil, nil, nil
	}

	// 旧数据（project_id=0）回退：候选负责人按部门分组，取各部门第一个（主岗）
	deptOwner := map[string]OwnerCandidate{}
	for _, c := range matchCandidates(tx, defect) {
		if _, ok := deptOwner[c.Dept]; !ok {
			deptOwner[c.Dept] = c
		}
	}

	assigned := []string{}
	missed := []string{}
	now := time.Now()
	for i := range disposals {
		d := disposals[i]
		// 已有人处理的不覆盖
		if d.AssigneeUsername != "" || d.AssigneeName != "" {
			continue
		}
		var owner OwnerCandidate
		var hit bool
		if d.ProjectID > 0 {
			// 处理人 = 涉及项目负责人主岗：项目负责人部门即项目所属部门，天然对应处置部门，
			// 不再额外校验「负责人部门==处置单部门」（旧 use 单按使用部门拆单导致的误伤已消除）
			if cands := queryOwnerCandidates(tx, d.ProjectType, d.ProjectID); len(cands) > 0 {
				owner, hit = cands[0], true
			}
		} else {
			owner, hit = deptOwner[d.DeptName]
		}
		if !hit {
			missed = append(missed, d.DeptName)
			continue
		}
		if err := tx.Model(&models.DefectDisposal{}).Where("id = ?", d.ID).Updates(map[string]interface{}{
			"assignee_username": owner.Username,
			"assignee_name":     owner.Name,
			"assigned_by":       operatorName,
			"assigned_at":       now,
			"update_time":       now,
		}).Error; err != nil {
			return assigned, missed, fmt.Errorf("自动分配 %s 处置单失败: %w", d.DeptName, err)
		}
		if err := e.addHistory(tx, defectNo, d.ID, DisposalPendingAssign, DisposalPendingAssign,
			"auto_assign", operatorName, "系统自动匹配处理人 "+owner.Name+"（项目主岗）"); err != nil {
			return assigned, missed, err
		}
		assigned = append(assigned, d.DeptName)
	}
	return uniqueStrings(assigned), uniqueStrings(missed), nil
}

// uniqueStrings 去重保序
func uniqueStrings(in []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s != "" && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// Analyze 研判：pending_analysis -> pending_review / assigned / ignored
// 缺陷(flow_type=defect)：研判后必须转待评审（上会）
// 漏洞(flow_type=vuln)：可直接分配(need_fix)或上会(need_review)或忽略(ignore)
func (e *Engine) Analyze(defectNo string, conclusion, analyzeResult string, assignMode string, assignees []string, userID, userName, role string) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		defect, err := e.loadDefectCtx(tx, defectNo)
		if err != nil {
			return err
		}

		if defect.Status != StatusPendingAnalysis {
			return fmt.Errorf("当前状态 %s 不允许研判", defect.Status)
		}
		// 研判权限：被分配的研判人；专员（admin/管理部门专员/维护部门专员）可代研判
		if !e.canAnalyzeDefect(userID, userName, role, defect) {
			return fmt.Errorf("您无权研判该缺陷（研判人或专员可操作）")
		}

		now := time.Now()
		from := defect.Status
		to := ""
		switch conclusion {
		case ConclusionNeedReview:
			to = StatusPendingReview
		case ConclusionIgnore:
			// 研判结论与外部研判口径一致（忽略/待上会）：缺陷也可研判直接忽略；
			// 需要修复分配仍须先上会评审（need_fix 仅漏洞直接分配）
			to = StatusIgnored
		case ConclusionNeedFix:
			if defect.FlowType != "vuln" {
				return fmt.Errorf("缺陷必须上会评审后才能分配")
			}
			to = StatusAssigned
		default:
			return fmt.Errorf("无效的研判结论: %s", conclusion)
		}

		// 状态写实体表 internal_status；研判明细写 defect_flow
		if err := e.updateStatus(tx, defectNo, to); err != nil {
			return err
		}
		detailUpdates := map[string]interface{}{
			"analyze_result": analyzeResult,
			"analyze_by":     userID,
			"analyze_time":   now,
			"update_time":    now,
		}
		if err := tx.Model(&models.DefectFlow{}).Where("defect_no = ?", defectNo).Updates(detailUpdates).Error; err != nil {
			return err
		}

		// 漏洞直接分配处置人：手动指定（多负责人各一张单）或自动分配（裂变+匹配），两通道
		if to == StatusAssigned {
			if assignMode == "manual" {
				if err := e.createAssignedDisposals(tx, defect, assignees, userName); err != nil {
					return err
				}
			} else {
				if err := e.fission(tx, defect); err != nil {
					return err
				}
				if _, _, err := e.autoAssignDisposalsCore(tx, defectNo, userName); err != nil {
					return err
				}
			}
		}

		return e.addHistory(tx, defectNo, 0, from, to, "analyze", userName, analyzeResult)
	})
}

// Review 评审（上会）：pending_review -> assigned / ignored
// 处置分配两通道（与研判分配一致）：
//   - assignMode = "manual"：手动指定一个或多个负责人（每人一张处置单，部门=各自所属部门）
//   - assignMode = "auto"（默认）：按涉及项目裂变 + 自动匹配各项目/软件处理人
func (e *Engine) Review(defectNo string, reviewConclusion, reviewMeetingNo, remark string, deadline *time.Time, assignMode string, assignees []string, userID, userName, role string) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		defect, err := e.loadDefectCtx(tx, defectNo)
		if err != nil {
			return err
		}

		if defect.Status != StatusPendingReview {
			return fmt.Errorf("当前状态 %s 不允许评审", defect.Status)
		}
		if !e.canReviewDefect(userID, userName, role) {
			return fmt.Errorf("您无权评审该缺陷")
		}

		now := time.Now()
		from := defect.Status
		to := ""
		switch reviewConclusion {
		case ReviewNeedFix:
			to = StatusAssigned
		case ReviewNoAction:
			to = StatusIgnored
		default:
			return fmt.Errorf("无效的评审结论: %s", reviewConclusion)
		}

		// 状态写实体表 internal_status；评审明细写 defect_flow
		if err := e.updateStatus(tx, defectNo, to); err != nil {
			return err
		}
		detailUpdates := map[string]interface{}{
			"review_meeting_no": reviewMeetingNo,
			"review_conclusion": reviewConclusion,
			"review_by":         userName,
			"review_time":       now,
			"update_time":       now,
		}
		if deadline != nil {
			detailUpdates["deadline"] = deadline
			defect.Deadline = deadline
		}
		if err := tx.Model(&models.DefectFlow{}).Where("defect_no = ?", defectNo).Updates(detailUpdates).Error; err != nil {
			return err
		}

		if to == StatusAssigned {
			if assignMode == "manual" {
				// 手动指定处置人（内外部缺陷/漏洞统一）：每个负责人各建一张处置单（部门=各自所属部门）
				if err := e.createAssignedDisposals(tx, defect, assignees, userName); err != nil {
					return err
				}
			} else {
				// 自动分配（默认）：按涉及项目裂变并自动匹配各项目/软件处理人
				if err := e.fission(tx, defect); err != nil {
					return err
				}
				if _, _, err := e.autoAssignDisposalsCore(tx, defectNo, userName); err != nil {
					return err
				}
			}
		}

		return e.addHistory(tx, defectNo, 0, from, to, "review", userName, remark)
	})
}

// createAssignedDisposals 手动指定处置人（多负责人，分隔符解析后传数组）：
// 为每个负责人各建一张处置单，部门/姓名按系统用户自动带出（任意部门人员均可指定）；
// 流程明细的 maint_dept 同步为第一个负责人部门（主运维部门）
func (e *Engine) createAssignedDisposals(tx *gorm.DB, defect *defectCtx, assignees []string, operatorName string) error {
	now := time.Now()
	seen := map[string]bool{}
	firstDept := ""
	for _, username := range assignees {
		username = strings.TrimSpace(username)
		if username == "" || seen[username] {
			continue
		}
		seen[username] = true

		var analyst models.User
		if err := tx.Where("username = ?", username).First(&analyst).Error; err != nil {
			return fmt.Errorf("负责人 %s 不是系统用户", username)
		}
		dept := analyst.Department
		if firstDept == "" {
			firstDept = dept
		}
		disposal := models.DefectDisposal{
			DefectNo:         defect.DefectNo,
			DisposalNo:       genNo("DP"),
			DeptName:         dept,
			DisposalType:     DisposalTypeMaint,
			AssigneeUsername: analyst.Username,
			AssigneeName:     analyst.Nickname,
			AssignedBy:       operatorName,
			Status:           DisposalPendingAssign,
			Deadline:         defect.Deadline,
			CreateTime:       now,
			UpdateTime:       now,
		}
		if err := tx.Create(&disposal).Error; err != nil {
			return fmt.Errorf("创建处置单失败: %w", err)
		}
	}
	if len(seen) == 0 {
		return fmt.Errorf("请指定至少一个负责人")
	}
	return tx.Model(&models.DefectFlow{}).Where("defect_no = ?", defect.DefectNo).
		Update("maint_dept", firstDept).Error
}

// DispatchVulnTx 漏洞情报裂变（在调用方事务内执行）：
// 缺陷实体已置「已分配」，直接裂变生成处置工单并自动匹配各涉及项目/软件的运维主岗第一负责人
func (e *Engine) DispatchVulnTx(tx *gorm.DB, defectNo, operatorName string) error {
	defect, err := e.loadDefectCtx(tx, defectNo)
	if err != nil {
		return err
	}
	if defect.Status != StatusAssigned {
		return fmt.Errorf("缺陷状态为 %s，仅「已分配」状态可直接裂变", defect.Status)
	}
	if err := e.fission(tx, defect); err != nil {
		return err
	}
	if _, _, err := e.autoAssignDisposalsCore(tx, defectNo, operatorName); err != nil {
		return err
	}
	return nil
}

// CreateAssignedDisposalsTx 手动指定处置人（在调用方事务内执行，对称 DispatchVulnTx）：
// 漏洞情报手动录入 manual 通道：为每个负责人各建一张处置单（部门=各自所属部门），不做裂变
func (e *Engine) CreateAssignedDisposalsTx(tx *gorm.DB, defectNo string, assignees []string, operatorName string) error {
	defect, err := e.loadDefectCtx(tx, defectNo)
	if err != nil {
		return err
	}
	return e.createAssignedDisposals(tx, defect, assignees, operatorName)
}

// FissionPreview 裂变预览单（漏洞情报自动分配预览：每张待生成处置单 + 候选负责人）
type FissionPreview struct {
	ProjectID    int64            `json:"projectId"`
	ProjectType  string           `json:"projectType"`  // software-基础软件 / inner-自研项目
	ProjectName  string           `json:"projectName"`  // 基础软件=软件名，组件=依赖项目代码
	DeptName     string           `json:"deptName"`     // 处置部门
	Team         string           `json:"team"`         // 团队（组件依赖项目）
	Candidates   []OwnerCandidate `json:"candidates"`   // 负责人候选（第一个为主岗）
	MatchedOwner *OwnerCandidate  `json:"matchedOwner"` // 自动匹配的主岗
}

// PreviewVulnFission 按影响组件名称预览自动裂变结果（不落库）：
//   - 基础软件（命中 project）：1 行（软件名/运维部门/运维负责人）
//   - 开源组件（命中 component_snapshot）：每个依赖项目 1 行（项目/团队/部门/负责人）
func (e *Engine) PreviewVulnFission(component string) ([]FissionPreview, error) {
	component = strings.TrimSpace(component)
	defect := &defectCtx{AffectedComponent: component}
	items, err := resolveFissionItems(e.db, defect)
	if err != nil {
		return nil, err
	}
	out := make([]FissionPreview, 0, len(items))
	for _, it := range items {
		p := FissionPreview{
			ProjectID:   it.projectID,
			ProjectType: it.projectType,
			ProjectName: it.projectName,
			DeptName:    it.deptName,
		}
		cands := queryOwnerCandidates(e.db, it.projectType, it.projectID)
		p.Candidates = cands
		if len(cands) > 0 {
			m := cands[0]
			p.MatchedOwner = &m
			p.Team = m.Team
		}
		out = append(out, p)
	}
	return out, nil
}

// fissionItem 裂变出的单张处置单维度：一张处置单 = 一个缺陷 × 一个涉及项目
type fissionItem struct {
	deptName     string // 处置部门
	disposalType string // maint-运维部门, use-使用部门
	projectID    int64  // 涉及项目ID（软件=project.id，组件=project_inner.id）
	projectName  string // 涉及项目显示名（软件=_name，组件=pro_code）
	projectType  string // software-基础软件, inner-自研项目
}

// resolveFissionItems 解析缺陷的「涉及项目」清单（裂变拆单的粒度 = 涉及项目）：
//   - 基础软件（affected_component 命中 project 表）：公共中间件统一运维，**只拆一张运维单**（不再按使用部门拆单）
//   - 开源组件（命中 component_snapshot）：每个依赖该组件的自研项目一张使用单（部门 = 项目所属部门），
//     同一部门有多个依赖项目即生成多张处置单
func resolveFissionItems(tx *gorm.DB, defect *defectCtx) ([]fissionItem, error) {
	component := strings.TrimSpace(defect.AffectedComponent)
	if component == "" {
		return nil, fmt.Errorf("未解析到涉及项目（缺少影响组件），无法裂变处置单")
	}

	// 1) 基础软件：project 表直接命中 —— 只生成一张运维单，由运维团队统一处置
	// 处置部门以软件在 project 表配置的运维部门为准（公共中间件真实运维归属，如 redis→技术部），
	// flow.maint_dept 仅是导入时的归属部门，不作为基础软件的裂变处置部门
	var prj models.Project
	if err := tx.Where("_name = ?", component).Order("id ASC").First(&prj).Error; err == nil {
		maintDept := prj.MaintDept
		if maintDept == "" {
			maintDept = defect.MaintDept
		}
		if maintDept == "" {
			return nil, fmt.Errorf("软件 %s 未配置运维部门，无法裂变处置单", prj.Name)
		}
		return []fissionItem{{maintDept, DisposalTypeMaint, prj.ID, prj.Name, "software"}}, nil
	}

	// 2) 开源组件：查 component_snapshot 依赖该组件的自研项目（PROJECT_ID -> project_inner）
	var snaps []models.ComponentSnapshot
	if err := tx.Where("_name = ?", component).Find(&snaps).Error; err != nil {
		return nil, err
	}
	seenPID := map[int64]bool{}
	var items []fissionItem
	for _, s := range snaps {
		if s.ProjectID == nil || *s.ProjectID == 0 || seenPID[*s.ProjectID] {
			continue
		}
		seenPID[*s.ProjectID] = true
		var inner models.ProjectInner
		if err := tx.First(&inner, *s.ProjectID).Error; err != nil || inner.DeptName == "" {
			continue // 项目已删除或未归属部门，跳过该依赖项目
		}
		items = append(items, fissionItem{inner.DeptName, DisposalTypeUse, inner.ID, inner.ProCode, "inner"})
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("组件 %s 的依赖项目未解析到所属部门，无法裂变处置单", component)
	}
	return items, nil
}

// fission 裂变：按「涉及项目」创建子处置单（状态「待处置」，处理人留空）
// 裂变后立即自动匹配各涉及项目主岗为处理人；未匹配到的工单由部门专员补分配/改派：
//   - 专员补分配/改派：POST /risk/disposals/:id/assign
//   - 系统自动分配：POST /risk/defects/:defectNo/auto-assign-disposals（匹配各涉及项目的主岗）
func (e *Engine) fission(tx *gorm.DB, defect *defectCtx) error {
	items, err := resolveFissionItems(tx, defect)
	if err != nil {
		return err
	}

	deadline := defect.Deadline
	now := time.Now()
	for _, item := range items {
		disposal := models.DefectDisposal{
			DefectNo:     defect.DefectNo,
			DisposalNo:   genNo("DP"),
			DeptName:     item.deptName,
			DisposalType: item.disposalType,
			ProjectID:    item.projectID,
			ProjectName:  item.projectName,
			ProjectType:  item.projectType,
			// 处理人留空：由裂变后自动匹配各涉及项目主岗，未匹配到的由部门专员补分配/改派
			Deadline:   deadline,
			Status:     DisposalPendingAssign,
			CreateTime: now,
			UpdateTime: now,
		}
		if err := tx.Create(&disposal).Error; err != nil {
			return fmt.Errorf("创建处置单失败: %w", err)
		}
	}
	return nil
}

// Assign 补分配/改派处理人：工单状态保持「待处置」（pending_assign），仅更新处理人
func (e *Engine) Assign(disposalID int64, assigneeUsername, assigneeName, userID, userName string) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		var d models.DefectDisposal
		if err := tx.First(&d, disposalID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("处置单不存在")
			}
			return err
		}

		if d.Status != DisposalPendingAssign {
			return fmt.Errorf("当前状态 %s 不允许分配", d.Status)
		}
		// 权限：本部门专员可分配/改派
		if !e.isSpecialistOf(userID, userName, d.DeptName, "dept_specialist", "usage_dept_specialist") {
			return fmt.Errorf("您不是该部门专员，无权分配")
		}
		if assigneeUsername == "" && assigneeName == "" {
			return fmt.Errorf("请指定处理人")
		}

		now := time.Now()
		from := d.Status
		remark := "分配给 " + assigneeName
		if d.AssigneeUsername != "" || d.AssigneeName != "" {
			remark = "改派为 " + assigneeName
		}
		if err := tx.Model(&models.DefectDisposal{}).Where("id = ?", disposalID).Updates(map[string]interface{}{
			"assignee_username": assigneeUsername,
			"assignee_name":     assigneeName,
			"assigned_by":       userName,
			"assigned_at":       now,
			"update_time":       now,
		}).Error; err != nil {
			return err
		}

		return e.addHistory(tx, d.DefectNo, disposalID, from, DisposalPendingAssign, "assign", userName, remark)
	})
}

// Complete 处置完成：processing/pending_assign -> completed，全部完成后闭环父缺陷
func (e *Engine) Complete(disposalID int64, ticketNo, note, userID, userName, role string) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		var d models.DefectDisposal
		if err := tx.First(&d, disposalID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("处置单不存在")
			}
			return err
		}

		if d.Status != DisposalProcessing && d.Status != DisposalPendingAssign {
			return fmt.Errorf("当前状态 %s 不允许完成", d.Status)
		}
		// 权限：仅处理人本人可完成（责任到人），admin 除外，专员不再代处置
		if role != "admin" && d.AssigneeName != userName && d.AssigneeUsername != userID {
			return fmt.Errorf("您不是该处置单的处理人，无权完成")
		}

		now := time.Now()
		from := d.Status
		if err := tx.Model(&models.DefectDisposal{}).Where("id = ?", disposalID).Updates(map[string]interface{}{
			"status":        DisposalCompleted,
			"ticket_no":     ticketNo,
			"disposal_note": note,
			"completed_at":  now,
			"update_time":   now,
		}).Error; err != nil {
			return err
		}

		if err := e.addHistory(tx, d.DefectNo, disposalID, from, DisposalCompleted, "complete", userName, note); err != nil {
			return err
		}

		return e.closeIfAllDone(tx, d.DefectNo)
	})
}

// MarkNA 标记忽略：仅处理人本人可将工单置为忽略（责任到人），admin 除外，专员不再代忽略
func (e *Engine) MarkNA(disposalID int64, remark, userID, userName, role string) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		var d models.DefectDisposal
		if err := tx.First(&d, disposalID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("处置单不存在")
			}
			return err
		}

		if d.Status == DisposalCompleted || d.Status == DisposalNA {
			return fmt.Errorf("当前状态 %s 不允许忽略", d.Status)
		}
		if role != "admin" && d.AssigneeName != userName && d.AssigneeUsername != userID {
			return fmt.Errorf("您不是该处置单的处理人，无权操作")
		}

		now := time.Now()
		from := d.Status
		if err := tx.Model(&models.DefectDisposal{}).Where("id = ?", disposalID).Updates(map[string]interface{}{
			"status":      DisposalNA,
			"update_time": now,
		}).Error; err != nil {
			return err
		}

		if err := e.addHistory(tx, d.DefectNo, disposalID, from, DisposalNA, "na", userName, remark); err != nil {
			return err
		}

		return e.closeIfAllDone(tx, d.DefectNo)
	})
}

// closeIfAllDone 若所有处置单均已完成/不适用，则闭环父缺陷
func (e *Engine) closeIfAllDone(tx *gorm.DB, defectNo string) error {
	var pending int64
	tx.Model(&models.DefectDisposal{}).
		Where("defect_no = ? AND status IN ?", defectNo, []string{DisposalPendingAssign, DisposalProcessing}).
		Count(&pending)
	if pending > 0 {
		return nil
	}

	ctx, err := e.loadDefectCtx(tx, defectNo)
	if err != nil {
		return err
	}
	if ctx.Status != StatusAssigned {
		return nil
	}

	now := time.Now()
	if err := e.updateStatus(tx, defectNo, StatusClosed); err != nil {
		return err
	}
	if err := tx.Model(&models.DefectFlow{}).Where("defect_no = ?", defectNo).Updates(map[string]interface{}{
		"close_time":  now,
		"update_time": now,
	}).Error; err != nil {
		return err
	}
	return e.addHistory(tx, defectNo, 0, StatusAssigned, StatusClosed, "close", "系统", "全部处置完成，自动闭环")
}

func (e *Engine) addHistory(tx *gorm.DB, defectNo string, disposalID int64, from, to, action, operator, remark string) error {
	h := models.DefectStatusHistory{
		DefectNo:     defectNo,
		DisposalID:   disposalID,
		FromStatus:   from,
		ToStatus:     to,
		Action:       action,
		Operator:     operator,
		OperatorName: operator,
		Remark:       remark,
		CreateTime:   time.Now(),
	}
	return tx.Create(&h).Error
}

// genNo 生成业务编号（如处置单号 DP-xxx）
func genNo(prefix string) string {
	// 用纳秒时间戳保证同一秒内批量裂变多张处置单时编号唯一
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
