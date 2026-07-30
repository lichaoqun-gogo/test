import type { TicketStatus, ApplicationType, RiskLevel, PreEvalStatus, PreEvalConclusion } from '@/types'

export interface Ticket {
  id: string
  ticketNo: string
  type: ApplicationType
  title: string
  applicant: string
  applicantDept: string
  status: TicketStatus
  riskLevel?: RiskLevel
  priority: 'normal' | 'urgent'
  createdAt: string
  updatedAt: string
  assignedTo?: string
  externalLink?: string
  description: string
  currentStep: string
}

export const mockTickets: Ticket[] = [
  {
    id: 't001',
    ticketNo: 'TK-2024-0301',
    type: 'introduce',
    title: '引入 Apache Pulsar 消息中间件',
    applicant: '张三',
    applicantDept: '技术研发部',
    status: 'processing',
    riskLevel: 'low',
    priority: 'normal',
    createdAt: '2024-03-01',
    updatedAt: '2024-03-05',
    assignedTo: '李四',
    externalLink: 'https://ticket.company.com/TK-2024-0301',
    description: '因业务需要引入 Pulsar 作为新一代消息中间件',
    currentStep: '预评估完成',
  },
  {
    id: 't002',
    ticketNo: 'TK-2024-0305',
    type: 'upgrade',
    title: '升级 Spring Boot 至 3.2.3',
    applicant: '张三',
    applicantDept: '技术研发部',
    status: 'pending',
    riskLevel: 'medium',
    priority: 'normal',
    createdAt: '2024-03-05',
    updatedAt: '2024-03-05',
    description: '升级核心框架以获取最新安全补丁',
    currentStep: '待审批',
  },
  {
    id: 't003',
    ticketNo: 'TK-2024-0308',
    type: 'risk_report',
    title: 'Log4j2 高危漏洞上报',
    applicant: '李四',
    applicantDept: '技术研发部',
    status: 'assigned',
    riskLevel: 'critical',
    priority: 'urgent',
    createdAt: '2024-03-08',
    updatedAt: '2024-03-10',
    assignedTo: '张三',
    description: '发现 Log4j2 存在新的远程代码执行漏洞',
    currentStep: '已分配',
  },
  {
    id: 't004',
    ticketNo: 'TK-2024-0310',
    type: 'withdraw',
    title: '退出 Bower 包管理器',
    applicant: '张三',
    applicantDept: '技术研发部',
    status: 'completed',
    riskLevel: 'low',
    priority: 'normal',
    createdAt: '2024-03-10',
    updatedAt: '2024-03-14',
    description: 'Bower 已停止维护，申请退出',
    currentStep: '已完成',
  },
  {
    id: 't005',
    ticketNo: 'TK-2024-0312',
    type: 'introduce',
    title: '引入 TailwindCSS 框架',
    applicant: '张三',
    applicantDept: '技术研发部',
    status: 'rejected',
    riskLevel: 'low',
    priority: 'normal',
    createdAt: '2024-03-12',
    updatedAt: '2024-03-13',
    description: '引入 TailwindCSS 用于前端项目样式管理',
    currentStep: '已驳回',
  },
  {
    id: 't006',
    ticketNo: 'TK-2024-0315',
    type: 'risk_report',
    title: 'event-stream 恶意包风险',
    applicant: '王五',
    applicantDept: '安全管理部',
    status: 'assigned',
    riskLevel: 'critical',
    priority: 'urgent',
    createdAt: '2024-03-15',
    updatedAt: '2024-03-15',
    assignedTo: '李四',
    description: '外部推送 event-stream 存在恶意代码注入',
    currentStep: '待二次分配',
  },
]

// 风险工单二次分配记录
export interface RedistributionRecord {
  id: string
  ticketId: string
  ticketNo: string
  ticketTitle: string
  riskLevel: RiskLevel
  fromDept: string
  assignedTo: string
  assignedBy: string
  assignedAt: string
  status: TicketStatus
  note: string
}

export const mockRedistributionRecords: RedistributionRecord[] = [
  {
    id: 'rd001',
    ticketId: 't003',
    ticketNo: 'TK-2024-0308',
    ticketTitle: 'Log4j2 高危漏洞上报',
    riskLevel: 'critical',
    fromDept: '安全管理部',
    assignedTo: '张三',
    assignedBy: '李四',
    assignedAt: '2024-03-10 14:30',
    status: 'processing',
    note: '请尽快修复，影响核心业务系统',
  },
  {
    id: 'rd002',
    ticketId: 't006',
    ticketNo: 'TK-2024-0315',
    ticketTitle: 'event-stream 恶意包风险',
    riskLevel: 'critical',
    fromDept: '安全管理部',
    assignedTo: '',
    assignedBy: '',
    assignedAt: '',
    status: 'pending',
    note: '',
  },
]

// 申请记录
export interface Application {
  id: string
  applicationNo: string
  type: ApplicationType
  componentName: string
  componentVersion: string
  applicant: string
  department: string
  status: TicketStatus
  preEvalStatus: PreEvalStatus
  preEvalConclusion?: PreEvalConclusion
  preEvalScore?: number
  ticketNo?: string
  createdAt: string
  updatedAt: string
  description: string
  businessScenario: string
  expectedGrade: string
  urgency: 'normal' | 'urgent'
}

export const mockApplications: Application[] = [
  {
    id: 'app001',
    applicationNo: 'APP-2024-0301',
    type: 'introduce',
    componentName: 'Apache Pulsar',
    componentVersion: '3.1.2',
    applicant: '张三',
    department: '技术研发部',
    status: 'processing',
    preEvalStatus: 'completed',
    preEvalConclusion: 'pass',
    preEvalScore: 85,
    ticketNo: 'TK-2024-0301',
    createdAt: '2024-03-01',
    updatedAt: '2024-03-05',
    description: '引入 Pulsar 作为新一代消息中间件替代 Kafka',
    businessScenario: '订单系统消息队列',
    expectedGrade: 'B',
    urgency: 'normal',
  },
  {
    id: 'app002',
    applicationNo: 'APP-2024-0305',
    type: 'upgrade',
    componentName: 'Spring Boot',
    componentVersion: '3.2.3',
    applicant: '张三',
    department: '技术研发部',
    status: 'pending',
    preEvalStatus: 'completed',
    preEvalConclusion: 'warning',
    preEvalScore: 72,
    ticketNo: 'TK-2024-0305',
    createdAt: '2024-03-05',
    updatedAt: '2024-03-05',
    description: '升级核心框架获取安全补丁和新特性',
    businessScenario: '全部 Java 微服务',
    expectedGrade: 'A',
    urgency: 'normal',
  },
  {
    id: 'app003',
    applicationNo: 'APP-2024-0310',
    type: 'withdraw',
    componentName: 'Bower',
    componentVersion: '1.8.14',
    applicant: '张三',
    department: '技术研发部',
    status: 'completed',
    preEvalStatus: 'completed',
    preEvalConclusion: 'pass',
    preEvalScore: 90,
    ticketNo: 'TK-2024-0310',
    createdAt: '2024-03-10',
    updatedAt: '2024-03-14',
    description: 'Bower 已停止维护，不再使用',
    businessScenario: '前端项目管理',
    expectedGrade: 'D',
    urgency: 'normal',
  },
  {
    id: 'app004',
    applicationNo: 'APP-2024-0312',
    type: 'introduce',
    componentName: 'TailwindCSS',
    componentVersion: '3.4.1',
    applicant: '张三',
    department: '技术研发部',
    status: 'rejected',
    preEvalStatus: 'completed',
    preEvalConclusion: 'warning',
    preEvalScore: 65,
    ticketNo: 'TK-2024-0312',
    createdAt: '2024-03-12',
    updatedAt: '2024-03-13',
    description: '引入 TailwindCSS 用于前端样式管理',
    businessScenario: '前端UI开发',
    expectedGrade: 'C',
    urgency: 'normal',
  },
]
