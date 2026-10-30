// 角色类型
export type RoleType = 'operator' | 'specialist' | 'admin'

export const ROLE_LABELS: Record<RoleType, string> = {
  operator: '普通运维者',
  specialist: '开源治理专员',
  admin: '管理者',
}

// 用户信息
export interface UserInfo {
  id: string
  username: string
  nickname: string
  role: RoleType
  department: string
  team: string
  avatar?: string
  email: string
  roles?: string[] // 角色编码（权限体系下发，如 ['specialist', 'mgmt_specialist']）
  permissions?: string[] // 权限点编码（权限体系下发，如 ['global_governance_view']）
}

// 权限等级 (A/B/C/D)
export type AssetGrade = 'A' | 'B' | 'C' | 'D'

export const GRADE_LABELS: Record<AssetGrade, string> = {
  A: '可安全使用',
  B: '谨慎使用',
  C: '限期整改',
  D: '禁止使用',
}

export const GRADE_COLORS: Record<AssetGrade, string> = {
  A: '#67C23A',
  B: '#409EFF',
  C: '#E6A23C',
  D: '#F56C6C',
}

// 风险等级
export type RiskLevel = 'low' | 'medium' | 'high' | 'critical'

export const RISK_LEVEL_LABELS: Record<RiskLevel, string> = {
  low: '低风险',
  medium: '中风险',
  high: '高风险',
  critical: '严重',
}

export const RISK_LEVEL_COLORS: Record<RiskLevel, string> = {
  low: '#67C23A',
  medium: '#E6A23C',
  high: '#F56C6C',
  critical: '#991B1B',
}

// 工单状态
export type TicketStatus = 'pending' | 'processing' | 'completed' | 'rejected' | 'assigned'

export const TICKET_STATUS_LABELS: Record<TicketStatus, string> = {
  pending: '待审批',
  processing: '处理中',
  completed: '已完成',
  rejected: '已驳回',
  assigned: '已分配',
}

export const TICKET_STATUS_COLORS: Record<TicketStatus, string> = {
  pending: '#409EFF',
  processing: '#E6A23C',
  completed: '#67C23A',
  rejected: '#F56C6C',
  assigned: '#909399',
}

// 申请类型
export type ApplicationType = 'introduce' | 'upgrade' | 'withdraw' | 'risk_report'

export const APPLICATION_TYPE_LABELS: Record<ApplicationType, string> = {
  introduce: '引入申请',
  upgrade: '升级申请',
  withdraw: '退出申请',
  risk_report: '风险上报',
}

// 任务类型（统一：引入申请/升级申请/退出申请/缺陷处置）
export type TaskType = 'introduce' | 'upgrade' | 'withdraw' | 'defect' | 'usage'

export const TASK_TYPE_LABELS: Record<TaskType, string> = {
  introduce: '引入申请',
  upgrade: '升级申请',
  withdraw: '退出申请',
  defect: '风险上报',
  usage: '申请使用',
}

export const TASK_TYPE_COLORS: Record<TaskType, string> = {
  introduce: '#909399',
  upgrade: '#909399',
  withdraw: '#909399',
  defect: '#F56C6C',
  usage: '#909399',
}

// 任务状态
export type TaskStatus = 'draft' | 'approving' | 'withdrawn' | 'approved' | 'rejected' | 'completed' | 'pending' | 'processing'

export const TASK_STATUS_LABELS: Record<TaskStatus, string> = {
  draft: '草稿',
  approving: '审批中',
  withdrawn: '已撤回',
  approved: '已通过',
  rejected: '已驳回',
  completed: '已完成',
  pending: '待审批',
  processing: '处理中',
}

export const TASK_STATUS_COLORS: Record<TaskStatus, string> = {
  draft: '#909399',
  approving: '#E6A23C',
  withdrawn: '#909399',
  approved: '#67C23A',
  rejected: '#F56C6C',
  completed: '#67C23A',
  pending: '#E6A23C',
  processing: '#409EFF',
}

// 审批节点
export interface ApprovalNode {
  id: string
  name: string // 使用部门意见/运维部门意见/管理员审批意见
  approver: string
  status: 'pending' | 'approved' | 'rejected' | 'not_required' | 'processing'
  opinion: string
  approvedAt?: string
}

// 任务详情
export interface Task {
  id: string
  taskNo: string // 任务编号
  type: TaskType
  title: string
  applicant: string
  applicantId?: string
  applicantDept: string
  status: TaskStatus
  riskLevel?: RiskLevel
  priority: 'normal' | 'urgent'
  createdAt: string
  createTime?: string
  updatedAt: string
  deadline?: string // 期望完成日期
  assignedTo?: string
  canApprove?: boolean // 后端标记：当前用户是否可处理当前审批节点（一个部门多名专员，任一在岗专员均可审批）
  externalLink?: string // 外部任务链接
  relatedTaskNo?: string // 关联任务编号
  description: string
  currentStep: string
  approvalNodes?: ApprovalNode[]
  approvalHistory?: any[]
  task_detail?: any
  // 缺陷处置字段
  disposalStatus?: 'no-need' | 'completed'
  disposalReason?: string
  disposalChangeOrder?: string
  // 引入申请专属字段
  selectedComponent?: string // 选择的组件ID
  componentName?: string
  componentVersion?: string
  licenseType?: string
  introduceReason?: string
  // 适用性评估数据
  applicabilityAssessment?: Record<string, {
    status: 'fully_meet' | 'deviation' | ''
    deviationExplanation?: string
  }>
  // 升级申请专属字段
  currentVersion?: string
  targetVersion?: string
  testReportLink?: string
  // 退出申请专属字段
  exitComponent?: string
  exitVersion?: string
  exitReason?: string
  withdrawReason?: string
  replacementComponent?: string
  alternativePlan?: string
  codeCleanupPlan?: string
  affectedProjects?: string
  // 风险上报专属字段
  internalDefectNo?: string // 内部缺陷编号（系统自动生成）
  riskDescription?: string
  // 缺陷处置专属字段
  riskSource?: string
  cveId?: string
  severityLevel?: RiskLevel
  affectedAssets?: string
  fixPlan?: string
  needFix?: boolean // 是否需要处置
  fixReason?: string // 不需要处置原因
  productionChangeNo?: string // 生产变更号
  productionTicketText?: string // 工单文本
  draftData?: any // 暂存数据
}

// 软件类型
export type SoftwareType = 'middleware' | 'database' | 'tool' | 'dependency' | 'library'

export const SOFTWARE_TYPE_LABELS: Record<SoftwareType, string> = {
  middleware: '中间件',
  database: '数据库',
  tool: '工具',
  dependency: '依赖包',
  library: '组件库',
}

// 预评估状态
export type PreEvalStatus = 'pending' | 'scanning' | 'completed' | 'failed'

export const PRE_EVAL_STATUS_LABELS: Record<PreEvalStatus, string> = {
  pending: '待评估',
  scanning: '评估中',
  completed: '评估完成',
  failed: '评估失败',
}

// 预评估结论
export type PreEvalConclusion = 'pass' | 'warning' | 'reject'

export const PRE_EVAL_CONCLUSION_LABELS: Record<PreEvalConclusion, string> = {
  pass: '通过',
  warning: '风险警告',
  reject: '不通过',
}

// 分页参数
export interface PaginationParams {
  page: number
  pageSize: number
  total: number
}

// 通用响应
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 通知类型
export type NotificationType = 'system' | 'risk_distribution' | 'secondary_assign' | 'approval_result'

export const NOTIFICATION_TYPE_LABELS: Record<NotificationType, string> = {
  system: '系统通知',
  risk_distribution: '风险分发',
  secondary_assign: '二次分配',
  approval_result: '审批结果',
}

export const NOTIFICATION_TYPE_COLORS: Record<NotificationType, string> = {
  system: '#409EFF',
  risk_distribution: '#F56C6C',
  secondary_assign: '#E6A23C',
  approval_result: '#67C23A',
}

// 通知级别
export type NotificationLevel = 'info' | 'warning' | 'danger'

export const NOTIFICATION_LEVEL_LABELS: Record<NotificationLevel, string> = {
  info: '普通',
  warning: '重要',
  danger: '紧急',
}

// 通知
export interface Notification {
  id: string
  type: NotificationType
  level: NotificationLevel
  title: string
  content: string
  isRead: boolean
  createdAt: string
  // 关联任务信息
  relatedTaskId?: string
  relatedTaskNo?: string
  // 风险分发字段
  distributor?: string // 分发人
  riskType?: string // 风险类型
  targetDepartment?: string // 目标部门
  // 二次分配字段
  assigner?: string // 分配人
  defectNo?: string // 缺陷编号
  assignee?: string // 被分配人
  // 审批结果字段
  approvalResult?: 'pass' | 'reject'
  approvalTaskType?: string
}
