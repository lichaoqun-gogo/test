import type { AssetGrade, SoftwareType, RiskLevel, TicketStatus } from '@/types'

// 全量资产台账
export interface FullAsset {
  id: string
  name: string
  version: string
  type: SoftwareType
  grade: AssetGrade
  language: string
  organization: string
  owner: string
  department: string
  introducedDate: string
  lastUpdated: string
  status: 'active' | 'deprecated' | 'pending_review'
  bugCount: number
  vulnerabilityCount: number
  licenseType: string
  licenseRisk: RiskLevel
}

export const mockFullAssets: FullAsset[] = [
  { id: 'fa001', name: 'Nginx', version: '1.25.4', type: 'middleware', grade: 'A', language: 'C', organization: 'F5 Networks', owner: '张三', department: '技术研发部', introducedDate: '2022-01-15', lastUpdated: '2024-03-15', status: 'active', bugCount: 2, vulnerabilityCount: 1, licenseType: 'BSD-2-Clause', licenseRisk: 'low' },
  { id: 'fa002', name: 'MySQL', version: '8.0.36', type: 'database', grade: 'A', language: 'C++', organization: 'Oracle', owner: '张三', department: '技术研发部', introducedDate: '2021-06-01', lastUpdated: '2024-03-10', status: 'active', bugCount: 5, vulnerabilityCount: 3, licenseType: 'GPL-2.0', licenseRisk: 'high' },
  { id: 'fa003', name: 'Redis', version: '7.2.4', type: 'database', grade: 'A', language: 'C', organization: 'Redis Ltd.', owner: '张三', department: '技术研发部', introducedDate: '2021-08-20', lastUpdated: '2024-03-08', status: 'active', bugCount: 1, vulnerabilityCount: 0, licenseType: 'BSD-3-Clause', licenseRisk: 'low' },
  { id: 'fa004', name: 'Spring Boot', version: '3.2.3', type: 'middleware', grade: 'B', language: 'Java', organization: 'VMware', owner: '李四', department: '技术研发部', introducedDate: '2022-03-10', lastUpdated: '2024-03-12', status: 'active', bugCount: 3, vulnerabilityCount: 1, licenseType: 'Apache-2.0', licenseRisk: 'low' },
  { id: 'fa005', name: 'Vue.js', version: '3.4.21', type: 'library', grade: 'B', language: 'TypeScript', organization: 'Vue Team', owner: '张三', department: '技术研发部', introducedDate: '2023-01-05', lastUpdated: '2024-03-14', status: 'active', bugCount: 2, vulnerabilityCount: 0, licenseType: 'MIT', licenseRisk: 'low' },
  { id: 'fa006', name: 'Kafka', version: '3.6.1', type: 'middleware', grade: 'A', language: 'Java/Scala', organization: 'Apache', owner: '李四', department: '技术研发部', introducedDate: '2021-11-15', lastUpdated: '2024-02-28', status: 'active', bugCount: 4, vulnerabilityCount: 2, licenseType: 'Apache-2.0', licenseRisk: 'low' },
  { id: 'fa007', name: 'Elasticsearch', version: '8.12.2', type: 'database', grade: 'B', language: 'Java', organization: 'Elastic', owner: '张三', department: '技术研发部', introducedDate: '2022-05-20', lastUpdated: '2024-03-05', status: 'active', bugCount: 6, vulnerabilityCount: 2, licenseType: 'SSPL-1.0', licenseRisk: 'medium' },
  { id: 'fa008', name: 'Log4j2', version: '2.22.1', type: 'library', grade: 'C', language: 'Java', organization: 'Apache', owner: '李四', department: '技术研发部', introducedDate: '2022-08-01', lastUpdated: '2024-01-20', status: 'active', bugCount: 1, vulnerabilityCount: 1, licenseType: 'Apache-2.0', licenseRisk: 'low' },
  { id: 'fa009', name: 'Docker', version: '25.0.3', type: 'tool', grade: 'A', language: 'Go', organization: 'Docker Inc.', owner: '张三', department: '技术研发部', introducedDate: '2021-03-01', lastUpdated: '2024-03-01', status: 'active', bugCount: 3, vulnerabilityCount: 1, licenseType: 'Apache-2.0', licenseRisk: 'low' },
  { id: 'fa010', name: 'Prometheus', version: '2.50.1', type: 'tool', grade: 'B', language: 'Go', organization: 'CNCF', owner: '李四', department: '技术研发部', introducedDate: '2022-07-15', lastUpdated: '2024-03-11', status: 'active', bugCount: 2, vulnerabilityCount: 0, licenseType: 'Apache-2.0', licenseRisk: 'low' },
  { id: 'fa011', name: 'lodash', version: '4.17.21', type: 'dependency', grade: 'C', language: 'JavaScript', organization: 'JS Foundation', owner: '张三', department: '技术研发部', introducedDate: '2021-01-01', lastUpdated: '2024-01-15', status: 'active', bugCount: 0, vulnerabilityCount: 0, licenseType: 'MIT', licenseRisk: 'low' },
  { id: 'fa012', name: 'Bower', version: '1.8.14', type: 'tool', grade: 'D', language: 'JavaScript', organization: 'Twitter', owner: '张三', department: '技术研发部', introducedDate: '2019-05-01', lastUpdated: '2022-01-10', status: 'deprecated', bugCount: 2, vulnerabilityCount: 1, licenseType: 'MIT', licenseRisk: 'low' },
]

// 许可证信息库
export interface LicenseInfo {
  id: string
  name: string
  fullName: string
  riskLevel: RiskLevel
  category: string
  description: string
  isApproved: boolean
  usageCount: number
}

export const mockLicenses: LicenseInfo[] = [
  { id: 'lic001', name: 'MIT', fullName: 'MIT License', riskLevel: 'low', category: '宽松型', description: '最宽松的开源许可证之一', isApproved: true, usageCount: 156 },
  { id: 'lic002', name: 'Apache-2.0', fullName: 'Apache License 2.0', riskLevel: 'low', category: '宽松型', description: 'Apache基金会开源许可证', isApproved: true, usageCount: 89 },
  { id: 'lic003', name: 'BSD-2-Clause', fullName: 'BSD 2-Clause License', riskLevel: 'low', category: '宽松型', description: 'BSD简化版许可证', isApproved: true, usageCount: 34 },
  { id: 'lic004', name: 'BSD-3-Clause', fullName: 'BSD 3-Clause License', riskLevel: 'low', category: '宽松型', description: 'BSD标准版许可证', isApproved: true, usageCount: 28 },
  { id: 'lic005', name: 'ISC', fullName: 'ISC License', riskLevel: 'low', category: '宽松型', description: 'ISC互联网系统协会许可证', isApproved: true, usageCount: 15 },
  { id: 'lic006', name: 'MPL-2.0', fullName: 'Mozilla Public License 2.0', riskLevel: 'medium', category: '弱传染性', description: 'Mozilla公共许可证', isApproved: true, usageCount: 12 },
  { id: 'lic007', name: 'LGPL-2.1', fullName: 'GNU Lesser GPL 2.1', riskLevel: 'medium', category: '弱传染性', description: 'GNU宽松通用公共许可证', isApproved: true, usageCount: 8 },
  { id: 'lic008', name: 'GPL-2.0', fullName: 'GNU General Public License 2.0', riskLevel: 'high', category: '强传染性', description: 'GPL标准许可证', isApproved: false, usageCount: 5 },
  { id: 'lic009', name: 'GPL-3.0', fullName: 'GNU General Public License 3.0', riskLevel: 'high', category: '强传染性', description: 'GPL最新版本许可证', isApproved: false, usageCount: 3 },
  { id: 'lic010', name: 'AGPL-3.0', fullName: 'GNU Affero GPL 3.0', riskLevel: 'critical', category: '强传染性', description: '网络使用也需开源', isApproved: false, usageCount: 2 },
  { id: 'lic011', name: 'SSPL-1.0', fullName: 'Server Side Public License', riskLevel: 'medium', category: '特殊型', description: 'MongoDB专用许可证', isApproved: true, usageCount: 1 },
]

// 黑白名单
export interface BlacklistItem {
  id: string
  name: string
  version?: string
  reason: string
  riskLevel: RiskLevel
  addedBy: string
  addedAt: string
  status: 'active' | 'removed'
}

export const mockBlacklist: BlacklistItem[] = [
  { id: 'bl001', name: 'event-stream', version: '3.3.6', reason: '恶意代码注入，窃取加密货币钱包', riskLevel: 'critical', addedBy: '王五', addedAt: '2024-01-15', status: 'active' },
  { id: 'bl002', name: 'log4js', version: '<6.4.0', reason: '已知远程代码执行漏洞未修复', riskLevel: 'high', addedBy: '王五', addedAt: '2024-02-20', status: 'active' },
  { id: 'bl003', name: 'ua-parser-js', version: '0.7.29-0.7.30', reason: '供应链攻击，植入挖矿脚本', riskLevel: 'critical', addedBy: '王五', addedAt: '2024-03-01', status: 'active' },
  { id: 'bl004', name: 'colors', version: '>=1.4.0', reason: '作者故意破坏，无限循环', riskLevel: 'high', addedBy: '王五', addedAt: '2024-01-20', status: 'active' },
  { id: 'bl005', name: 'faker', version: '>=6.0.0', reason: '作者故意破坏，数据损坏', riskLevel: 'high', addedBy: '王五', addedAt: '2024-01-20', status: 'active' },
  { id: 'bl006', name: 'node-ipc', version: '>=11.0.0', reason: '地缘政治抗议，植入恶意代码', riskLevel: 'critical', addedBy: '王五', addedAt: '2024-02-15', status: 'active' },
]

// 统一风险台账
export interface RiskRecord {
  id: string
  source: string
  title: string
  riskLevel: RiskLevel
  riskType: string
  affectedComponent: string
  affectedVersion: string
  description: string
  discoveredDate: string
  status: 'pending' | 'dispatched' | 'processing' | 'resolved'
  dispatchTarget?: string
  dispatchDept?: string
  dispatchTime?: string
  dispatcher?: string
}

export const mockRiskRecords: RiskRecord[] = [
  { id: 'rk001', source: '外部推送', title: 'Log4j2 远程代码执行漏洞', riskLevel: 'critical', riskType: '安全漏洞', affectedComponent: 'Log4j2', affectedVersion: '<2.22.1', description: '发现新的RCE漏洞可利用JNDI注入', discoveredDate: '2024-03-08', status: 'dispatched', dispatchTarget: '李四', dispatchDept: '技术研发部', dispatchTime: '2024-03-09', dispatcher: '王五' },
  { id: 'rk002', source: 'SCA扫描', title: 'event-stream 恶意包检测', riskLevel: 'critical', riskType: '恶意代码', affectedComponent: 'event-stream', affectedVersion: '3.3.6', description: '依赖链中发现恶意代码注入', discoveredDate: '2024-03-15', status: 'pending', dispatchTarget: '', dispatchDept: '', dispatchTime: '', dispatcher: '' },
  { id: 'rk003', source: '外部推送', title: 'Spring4Shell 漏洞预警', riskLevel: 'high', riskType: '安全漏洞', affectedComponent: 'Spring Framework', affectedVersion: '<5.3.18', description: 'Spring框架远程代码执行漏洞', discoveredDate: '2024-02-20', status: 'resolved', dispatchTarget: '张三', dispatchDept: '技术研发部', dispatchTime: '2024-02-21', dispatcher: '王五' },
  { id: 'rk004', source: 'SCA扫描', title: 'Express 路径穿越漏洞', riskLevel: 'medium', riskType: '安全漏洞', affectedComponent: 'express', affectedVersion: '4.18.2', description: '路径穿越可导致文件读取', discoveredDate: '2024-03-05', status: 'dispatched', dispatchTarget: '张三', dispatchDept: '技术研发部', dispatchTime: '2024-03-06', dispatcher: '王五' },
  { id: 'rk005', source: 'SCA扫描', title: 'd3.js 原型链污染', riskLevel: 'high', riskType: '安全漏洞', affectedComponent: 'd3', affectedVersion: '7.8.5', description: '数据绑定过程存在原型链污染', discoveredDate: '2024-03-08', status: 'processing', dispatchTarget: '李四', dispatchDept: '数据智能部', dispatchTime: '2024-03-09', dispatcher: '王五' },
]

// 工单追踪
export interface TicketTracking {
  id: string
  ticketNo: string
  title: string
  type: string
  status: TicketStatus
  assignee: string
  department: string
  riskLevel: RiskLevel
  createdAt: string
  dueDate: string
  progress: number
  isOverdue: boolean
}

export const mockTicketTracking: TicketTracking[] = [
  { id: 'tt001', ticketNo: 'TK-2024-0301', title: '引入 Apache Pulsar', type: '引入申请', status: 'processing', assignee: '张三', department: '技术研发部', riskLevel: 'low', createdAt: '2024-03-01', dueDate: '2024-03-15', progress: 60, isOverdue: false },
  { id: 'tt002', ticketNo: 'TK-2024-0305', title: '升级 Spring Boot', type: '升级申请', status: 'pending', assignee: '李四', department: '技术研发部', riskLevel: 'medium', createdAt: '2024-03-05', dueDate: '2024-03-20', progress: 20, isOverdue: false },
  { id: 'tt003', ticketNo: 'TK-2024-0308', title: 'Log4j2 漏洞修复', type: '缺陷处理', status: 'processing', assignee: '张三', department: '技术研发部', riskLevel: 'critical', createdAt: '2024-03-08', dueDate: '2024-03-12', progress: 80, isOverdue: true },
  { id: 'tt004', ticketNo: 'TK-2024-0310', title: '退出 Bower', type: '退出申请', status: 'completed', assignee: '张三', department: '技术研发部', riskLevel: 'low', createdAt: '2024-03-10', dueDate: '2024-03-20', progress: 100, isOverdue: false },
  { id: 'tt005', ticketNo: 'TK-2024-0315', title: 'event-stream 风险处置', type: '风险工单', status: 'pending', assignee: '待分配', department: '技术研发部', riskLevel: 'critical', createdAt: '2024-03-15', dueDate: '2024-03-18', progress: 0, isOverdue: true },
]

// 引入审批
export interface IntroductionApproval {
  id: string
  applicationNo: string
  componentName: string
  componentVersion: string
  applicant: string
  department: string
  status: TicketStatus
  preEvalScore: number
  preEvalConclusion: string
  submittedAt: string
  currentStep: string
}

export const mockIntroductions: IntroductionApproval[] = [
  { id: 'ia001', applicationNo: 'APP-2024-0301', componentName: 'Apache Pulsar', componentVersion: '3.1.2', applicant: '张三', department: '技术研发部', status: 'processing', preEvalScore: 85, preEvalConclusion: '通过', submittedAt: '2024-03-01', currentStep: '业务终评' },
  { id: 'ia002', applicationNo: 'APP-2024-0312', componentName: 'TailwindCSS', componentVersion: '3.4.1', applicant: '张三', department: '技术研发部', status: 'rejected', preEvalScore: 65, preEvalConclusion: '风险警告', submittedAt: '2024-03-12', currentStep: '已驳回' },
  { id: 'ia003', applicationNo: 'APP-2024-0318', componentName: 'Apache Doris', componentVersion: '2.0.10', applicant: '李四', department: '数据智能部', status: 'pending', preEvalScore: 78, preEvalConclusion: '通过', submittedAt: '2024-03-18', currentStep: '待审批' },
]

// 退出审批
export interface WithdrawalApproval {
  id: string
  applicationNo: string
  componentName: string
  currentGrade: AssetGrade
  disposalMethod: string
  applicant: string
  department: string
  status: TicketStatus
  reason: string
  submittedAt: string
}

export const mockWithdrawals: WithdrawalApproval[] = [
  { id: 'wa001', applicationNo: 'APP-2024-0310', componentName: 'Bower', currentGrade: 'D', disposalMethod: '迁移至 npm', applicant: '张三', department: '技术研发部', status: 'completed', reason: '已停止维护', submittedAt: '2024-03-10' },
  { id: 'wa002', applicationNo: 'APP-2024-0320', componentName: 'Grunt', currentGrade: 'D', disposalMethod: '迁移至 Gulp', applicant: '李四', department: '技术研发部', status: 'pending', reason: '生态衰退，社区活跃度低', submittedAt: '2024-03-20' },
]

// 升级评审
export interface UpgradeReview {
  id: string
  applicationNo: string
  componentName: string
  fromVersion: string
  toVersion: string
  currentGrade: AssetGrade
  targetGrade: AssetGrade
  applicant: string
  department: string
  status: TicketStatus
  reason: string
  submittedAt: string
}

export const mockUpgrades: UpgradeReview[] = [
  { id: 'ur001', applicationNo: 'APP-2024-0305', componentName: 'Spring Boot', fromVersion: '3.1.5', toVersion: '3.2.3', currentGrade: 'A', targetGrade: 'A', applicant: '张三', department: '技术研发部', status: 'pending', reason: '获取安全补丁和新特性', submittedAt: '2024-03-05' },
  { id: 'ur002', applicationNo: 'APP-2024-0322', componentName: 'MySQL', fromVersion: '8.0.35', toVersion: '8.0.36', currentGrade: 'A', targetGrade: 'A', applicant: '李四', department: '技术研发部', status: 'pending', reason: '修复已知漏洞', submittedAt: '2024-03-22' },
]

// 存量治理
export interface InventoryItem {
  id: string
  name: string
  version: string
  project: string
  discoveredBy: string
  isReported: boolean
  grade: AssetGrade
  riskLevel: RiskLevel
  status: 'pending_review' | 'graded' | 'approved' | 'rejected'
  discoveredDate: string
}

export const mockInventory: InventoryItem[] = [
  { id: 'inv001', name: 'moment.js', version: '2.30.1', project: '订单管理系统', discoveredBy: 'SCA扫描', isReported: false, grade: 'C', riskLevel: 'low', status: 'graded', discoveredDate: '2024-03-10' },
  { id: 'inv002', name: 'request', version: '2.88.2', project: '数据分析平台', discoveredBy: 'SCA扫描', isReported: false, grade: 'D', riskLevel: 'medium', status: 'pending_review', discoveredDate: '2024-03-12' },
  { id: 'inv003', name: 'jquery', version: '3.7.1', project: '用户认证中心', discoveredBy: 'SCA扫描', isReported: true, grade: 'C', riskLevel: 'low', status: 'approved', discoveredDate: '2024-03-08' },
  { id: 'inv004', name: 'node-sass', version: '9.0.0', project: '订单管理系统', discoveredBy: 'SCA扫描', isReported: false, grade: 'D', riskLevel: 'medium', status: 'pending_review', discoveredDate: '2024-03-14' },
]

// 违规审计
export interface ViolationRecord {
  id: string
  componentName: string
  version: string
  project: string
  introducedBy: string
  department: string
  violationType: string
  riskLevel: RiskLevel
  status: 'pending' | 'handling' | 'resolved'
  discoveredDate: string
}

export const mockViolations: ViolationRecord[] = [
  { id: 'vio001', componentName: 'event-stream', version: '3.3.6', project: '数据分析平台', introducedBy: '未知', department: '数据智能部', violationType: '未走流程私自引入', riskLevel: 'critical', status: 'handling', discoveredDate: '2024-03-15' },
  { id: 'vio002', componentName: 'moment.js', version: '2.30.1', project: '订单管理系统', introducedBy: '赵六', department: '技术研发部', violationType: '未报备使用', riskLevel: 'low', status: 'resolved', discoveredDate: '2024-03-10' },
]

// 看板数据
export interface DashboardData {
  blockStats: { month: string; count: number }[]
  repairStats: { department: string; mttr: number; overdue: number }[]
  healthStats: { status: string; count: number; percentage: number }[]
  workloadStats: { department: string; introduce: number; upgrade: number }[]
  coreDeps: { name: string; usageCount: number; projects: number }[]
  versionFragment: { name: string; versionCount: number; versions: string[] }[]
  riskProjects: { name: string; vulnerabilityCount: number; riskLevel: RiskLevel }[]
}

export const mockDashboardData: DashboardData = {
  blockStats: [
    { month: '2024-01', count: 12 },
    { month: '2024-02', count: 18 },
    { month: '2024-03', count: 25 },
    { month: '2024-04', count: 15 },
    { month: '2024-05', count: 20 },
    { month: '2024-06', count: 8 },
  ],
  repairStats: [
    { department: '技术研发部', mttr: 3.5, overdue: 2 },
    { department: '数据智能部', mttr: 5.2, overdue: 4 },
    { department: '产品技术部', mttr: 4.1, overdue: 1 },
    { department: '基础架构部', mttr: 2.8, overdue: 0 },
  ],
  healthStats: [
    { status: '正常', count: 156, percentage: 78 },
    { status: '停更', count: 28, percentage: 14 },
    { status: '无维护', count: 16, percentage: 8 },
  ],
  workloadStats: [
    { department: '技术研发部', introduce: 45, upgrade: 12 },
    { department: '数据智能部', introduce: 18, upgrade: 5 },
    { department: '产品技术部', introduce: 22, upgrade: 8 },
    { department: '基础架构部', introduce: 30, upgrade: 15 },
  ],
  coreDeps: [
    { name: 'lodash', usageCount: 89, projects: 45 },
    { name: 'Spring Boot', usageCount: 76, projects: 38 },
    { name: 'Vue.js', usageCount: 65, projects: 32 },
    { name: 'axios', usageCount: 58, projects: 29 },
    { name: 'MySQL', usageCount: 52, projects: 26 },
  ],
  versionFragment: [
    { name: 'Spring Boot', versionCount: 8, versions: ['2.7.x', '3.0.x', '3.1.x', '3.2.x'] },
    { name: 'Vue.js', versionCount: 5, versions: ['2.6.x', '2.7.x', '3.2.x', '3.3.x', '3.4.x'] },
    { name: 'MySQL', versionCount: 4, versions: ['5.7.x', '8.0.32', '8.0.35', '8.0.36'] },
    { name: 'Redis', versionCount: 3, versions: ['6.2.x', '7.0.x', '7.2.x'] },
  ],
  riskProjects: [
    { name: '数据分析平台', vulnerabilityCount: 7, riskLevel: 'critical' },
    { name: '订单管理系统', vulnerabilityCount: 3, riskLevel: 'high' },
    { name: '消息推送服务', vulnerabilityCount: 2, riskLevel: 'medium' },
    { name: '用户认证中心', vulnerabilityCount: 1, riskLevel: 'low' },
  ],
}

// 活跃度监测数据
export interface ActivityData {
  componentName: string
  stars: number
  forks: number
  openIssues: number
  closedIssues: number
  commitsLastMonth: number
  contributors: number
  lastRelease: string
  trendData: { month: string; stars: number; commits: number }[]
}

export const mockActivityData: ActivityData[] = [
  {
    componentName: 'Nginx',
    stars: 18500,
    forks: 4200,
    openIssues: 125,
    closedIssues: 3200,
    commitsLastMonth: 45,
    contributors: 89,
    lastRelease: '2024-03-15',
    trendData: [
      { month: '2023-10', stars: 17800, commits: 52 },
      { month: '2023-11', stars: 17950, commits: 48 },
      { month: '2023-12', stars: 18100, commits: 35 },
      { month: '2024-01', stars: 18250, commits: 42 },
      { month: '2024-02', stars: 18380, commits: 38 },
      { month: '2024-03', stars: 18500, commits: 45 },
    ],
  },
  {
    componentName: 'Redis',
    stars: 65200,
    forks: 22800,
    openIssues: 890,
    closedIssues: 8900,
    commitsLastMonth: 128,
    contributors: 456,
    lastRelease: '2024-03-08',
    trendData: [
      { month: '2023-10', stars: 62000, commits: 145 },
      { month: '2023-11', stars: 62800, commits: 132 },
      { month: '2023-12', stars: 63500, commits: 98 },
      { month: '2024-01', stars: 64200, commits: 115 },
      { month: '2024-02', stars: 64800, commits: 120 },
      { month: '2024-03', stars: 65200, commits: 128 },
    ],
  },
]

// 备份状态
export interface BackupRecord {
  id: string
  componentName: string
  version: string
  status: 'completed' | 'running' | 'failed' | 'pending'
  storagePath: string
  size: string
  lastUpdated: string
  triggeredBy: string
}

export const mockBackups: BackupRecord[] = [
  { id: 'bk001', componentName: 'Nginx', version: '1.25.4', status: 'completed', storagePath: 'oss://backups/nginx/1.25.4/', size: '1.2MB', lastUpdated: '2024-03-15', triggeredBy: '自动' },
  { id: 'bk002', componentName: 'Redis', version: '7.2.4', status: 'completed', storagePath: 'oss://backups/redis/7.2.4/', size: '3.2MB', lastUpdated: '2024-03-08', triggeredBy: '自动' },
  { id: 'bk003', componentName: 'Spring Boot', version: '3.2.3', status: 'running', storagePath: 'oss://backups/spring-boot/3.2.3/', size: '-', lastUpdated: '2024-03-20', triggeredBy: '手动' },
  { id: 'bk004', componentName: 'Kafka', version: '3.6.1', status: 'failed', storagePath: '-', size: '-', lastUpdated: '2024-03-18', triggeredBy: '自动' },
]

// Webhook事件
export interface WebhookEvent {
  id: string
  eventType: string
  componentName: string
  version: string
  syncStatus: 'success' | 'failed' | 'pending'
  receivedAt: string
  metadata: string
}

export const mockWebhookEvents: WebhookEvent[] = [
  { id: 'wh001', eventType: 'package.published', componentName: 'lodash', version: '4.17.21', syncStatus: 'success', receivedAt: '2024-03-15 10:30', metadata: '{"registry":"npm","downloads":15234}' },
  { id: 'wh002', eventType: 'package.updated', componentName: 'Spring Boot', version: '3.2.3', syncStatus: 'success', receivedAt: '2024-03-14 14:20', metadata: '{"registry":"maven","checksum":"abc123"}' },
  { id: 'wh003', eventType: 'package.published', componentName: 'Vue.js', version: '3.4.21', syncStatus: 'success', receivedAt: '2024-03-14 09:15', metadata: '{"registry":"npm","downloads":28900}' },
  { id: 'wh004', eventType: 'package.vulnerability', componentName: 'Log4j2', version: '2.22.0', syncStatus: 'success', receivedAt: '2024-03-08 16:45', metadata: '{"cve":"CVE-2024-1001","severity":"high"}' },
]

// ==================== 治理中心详细数据 ====================

// 引入审批详情数据
export interface IntroduceApprovalDetail {
  id: string
  applicationNo: string
  componentName: string
  componentVersion: string
  componentType: SoftwareType
  applicant: string
  department: string
  email: string
  phone: string
  status: TicketStatus
  preEvalScore: number
  preEvalConclusion: 'pass' | 'warning' | 'reject'
  submittedAt: string
  currentStep: string
  useScenario: string
  projectNames: string[]
  licenseType: string
  officialUrl: string
  description: string
  grade: AssetGrade
  finalMatch?: 'match' | 'partial' | 'unmatch'
  finalApplicability?: 'apply' | 'conditional' | 'unapply'
  finalRisk?: 'low' | 'medium' | 'high'
  finalGrade?: AssetGrade
  finalOpinion?: string
  approvalHistory?: {
    action: string
    operator: string
    opinion: string
    time: string
  }[]
}

export const mockIntroduceDetails: IntroduceApprovalDetail[] = [
  {
    id: 'ia001',
    applicationNo: 'APP-2024-0301',
    componentName: 'Apache Pulsar',
    componentVersion: '3.1.2',
    componentType: 'middleware',
    applicant: '张三',
    department: '技术研发部',
    email: 'zhangsan@company.com',
    phone: '13800138001',
    status: 'processing',
    preEvalScore: 85,
    preEvalConclusion: 'pass',
    submittedAt: '2024-03-01',
    currentStep: '业务终评',
    useScenario: '用于订单系统的消息队列，替代原有的RabbitMQ，支持更高吞吐量和更低的延迟',
    projectNames: ['订单管理系统', '支付结算平台'],
    licenseType: 'Apache-2.0',
    officialUrl: 'https://pulsar.apache.org',
    description: 'Apache Pulsar是一个云原生分布式消息和流处理平台',
    grade: 'B',
    approvalHistory: [
      { action: '提交申请', operator: '张三', opinion: '申请引入Apache Pulsar 3.1.2作为消息中间件', time: '2024-03-01 09:00' },
      { action: '预评估完成', operator: '系统', opinion: '预评估评分85分，结论：通过', time: '2024-03-01 09:30' },
      { action: '初审通过', operator: '王五', opinion: '预评估结果良好，进入业务终评', time: '2024-03-02 14:00' },
    ],
  },
  {
    id: 'ia002',
    applicationNo: 'APP-2024-0312',
    componentName: 'TailwindCSS',
    componentVersion: '3.4.1',
    componentType: 'library',
    applicant: '张三',
    department: '技术研发部',
    email: 'zhangsan@company.com',
    phone: '13800138001',
    status: 'rejected',
    preEvalScore: 65,
    preEvalConclusion: 'warning',
    submittedAt: '2024-03-12',
    currentStep: '已驳回',
    useScenario: '用于前端项目的CSS框架，提高开发效率',
    projectNames: ['用户中心前端'],
    licenseType: 'MIT',
    officialUrl: 'https://tailwindcss.com',
    description: 'Tailwind CSS是一个功能类优先的CSS框架',
    grade: 'C',
    finalOpinion: '与现有UI框架冲突，建议统一使用Element Plus',
    approvalHistory: [
      { action: '提交申请', operator: '张三', opinion: '申请引入TailwindCSS 3.4.1', time: '2024-03-12 10:00' },
      { action: '预评估完成', operator: '系统', opinion: '预评估评分65分，结论：风险警告', time: '2024-03-12 10:30' },
      { action: '驳回', operator: '王五', opinion: '与现有UI框架冲突，建议统一使用Element Plus', time: '2024-03-13 09:00' },
    ],
  },
  {
    id: 'ia003',
    applicationNo: 'APP-2024-0318',
    componentName: 'Apache Doris',
    componentVersion: '2.0.10',
    componentType: 'database',
    applicant: '李四',
    department: '数据智能部',
    email: 'lisi@company.com',
    phone: '13800138002',
    status: 'pending',
    preEvalScore: 78,
    preEvalConclusion: 'pass',
    submittedAt: '2024-03-18',
    currentStep: '待审批',
    useScenario: '用于实时数据分析平台，提供亚秒级查询响应',
    projectNames: ['数据分析平台', '运营看板系统'],
    licenseType: 'Apache-2.0',
    officialUrl: 'https://doris.apache.org',
    description: 'Apache Doris是一个基于MPP架构的高性能实时分析数据库',
    grade: 'B',
    approvalHistory: [
      { action: '提交申请', operator: '李四', opinion: '申请引入Apache Doris 2.0.10用于实时分析', time: '2024-03-18 11:00' },
      { action: '预评估完成', operator: '系统', opinion: '预评估评分78分，结论：通过', time: '2024-03-18 11:30' },
    ],
  },
  {
    id: 'ia004',
    applicationNo: 'APP-2024-0325',
    componentName: 'MinIO',
    componentVersion: '2024.3.3',
    componentType: 'middleware',
    applicant: '赵六',
    department: '基础架构部',
    email: 'zhaoliu@company.com',
    phone: '13800138003',
    status: 'completed',
    preEvalScore: 92,
    preEvalConclusion: 'pass',
    submittedAt: '2024-03-25',
    currentStep: '审批完成',
    useScenario: '用于对象存储，存储用户上传的文件和图片',
    projectNames: ['文件管理系统', '用户中心'],
    licenseType: 'AGPL-3.0',
    officialUrl: 'https://min.io',
    description: 'MinIO是一个高性能、兼容S3 API的对象存储',
    grade: 'A',
    finalMatch: 'match',
    finalApplicability: 'apply',
    finalRisk: 'low',
    finalGrade: 'A',
    finalOpinion: '完全满足业务需求，S3兼容性好，性能优异',
    approvalHistory: [
      { action: '提交申请', operator: '赵六', opinion: '申请引入MinIO作为对象存储方案', time: '2024-03-25 09:00' },
      { action: '预评估完成', operator: '系统', opinion: '预评估评分92分，结论：通过', time: '2024-03-25 09:30' },
      { action: '初审通过', operator: '王五', opinion: '预评估结果优秀，进入业务终评', time: '2024-03-26 10:00' },
      { action: '业务终评通过', operator: '管理员', opinion: '完全满足业务需求，S3兼容性好，性能优异', time: '2024-03-27 14:00' },
      { action: '审批完成', operator: '系统', opinion: '流程完成，组件已入库', time: '2024-03-27 14:05' },
    ],
  },
]

// 退出审批详情数据
export interface ExitApprovalDetail {
  id: string
  applicationNo: string
  componentName: string
  componentVersion: string
  currentGrade: AssetGrade
  disposalMethod: string
  applicant: string
  department: string
  status: TicketStatus
  reason: string
  submittedAt: string
  affectedProjects: string[]
  alternativePlan: string
  codeCleanupPlan: string
  riskAssessment: string
  approvalHistory?: {
    action: string
    operator: string
    opinion: string
    time: string
  }[]
}

export const mockExitDetails: ExitApprovalDetail[] = [
  {
    id: 'wa001',
    applicationNo: 'APP-2024-0310',
    componentName: 'Bower',
    componentVersion: '1.8.14',
    currentGrade: 'D',
    disposalMethod: '迁移至npm',
    applicant: '张三',
    department: '技术研发部',
    status: 'completed',
    reason: '已停止维护，社区活跃度极低，存在安全风险',
    submittedAt: '2024-03-10',
    affectedProjects: ['前端基础框架'],
    alternativePlan: '已将所有Bower依赖迁移至npm管理',
    codeCleanupPlan: '删除bower.json配置文件，清理bower_components目录',
    riskAssessment: '低风险，迁移方案已验证',
    approvalHistory: [
      { action: '提交退出申请', operator: '张三', opinion: 'Bower已停止维护，申请退出', time: '2024-03-10 10:00' },
      { action: '审批通过', operator: '王五', opinion: '同意退出，请按计划执行迁移', time: '2024-03-11 09:00' },
    ],
  },
  {
    id: 'wa002',
    applicationNo: 'APP-2024-0320',
    componentName: 'Grunt',
    componentVersion: '1.6.1',
    currentGrade: 'D',
    disposalMethod: '迁移至Gulp/Vite',
    applicant: '李四',
    department: '技术研发部',
    status: 'pending',
    reason: '生态衰退，社区活跃度低，构建速度慢',
    submittedAt: '2024-03-20',
    affectedProjects: ['订单管理系统前端', '用户中心前端'],
    alternativePlan: '计划迁移至Vite构建工具',
    codeCleanupPlan: '删除Gruntfile.js及相关配置文件',
    riskAssessment: '中等风险，需要适配新的构建流程',
    approvalHistory: [
      { action: '提交退出申请', operator: '李四', opinion: 'Grunt构建效率低下，申请退出并迁移至Vite', time: '2024-03-20 14:00' },
    ],
  },
  {
    id: 'wa003',
    applicationNo: 'APP-2024-0328',
    componentName: 'Moment.js',
    componentVersion: '2.30.1',
    currentGrade: 'C',
    disposalMethod: '替换为dayjs',
    applicant: '赵六',
    department: '基础架构部',
    status: 'processing',
    reason: '官方已进入维护模式，推荐使用dayjs替代',
    submittedAt: '2024-03-28',
    affectedProjects: ['数据分析平台', '运营看板系统', '订单管理系统'],
    alternativePlan: '使用dayjs替代，API基本兼容',
    codeCleanupPlan: '全局替换moment引用，删除moment依赖',
    riskAssessment: '低风险，dayjs API与moment高度兼容',
    approvalHistory: [
      { action: '提交退出申请', operator: '赵六', opinion: 'Moment.js已进入维护模式，申请替换为dayjs', time: '2024-03-28 09:00' },
      { action: '初审通过', operator: '王五', opinion: '同意退出，请制定详细迁移计划', time: '2024-03-29 10:00' },
    ],
  },
]

// 升级评审详情数据
export interface UpgradeReviewDetail {
  id: string
  applicationNo: string
  componentName: string
  fromVersion: string
  toVersion: string
  currentGrade: AssetGrade
  targetGrade: AssetGrade
  applicant: string
  department: string
  status: TicketStatus
  reason: string
  submittedAt: string
  breakingChanges: string[]
  securityFixes: string[]
  newFeatures: string[]
  compatibilityRisk: string
  testReport: string
  rollbackPlan: string
  approvalHistory?: {
    action: string
    operator: string
    opinion: string
    time: string
  }[]
  beforeEval?: {
    vulnerabilityCount: number
    healthScore: number
    licenseRisk: string
  }
  afterEval?: {
    vulnerabilityCount: number
    healthScore: number
    licenseRisk: string
  }
}

export const mockUpgradeDetails: UpgradeReviewDetail[] = [
  {
    id: 'ur001',
    applicationNo: 'APP-2024-0305',
    componentName: 'Spring Boot',
    fromVersion: '3.1.5',
    toVersion: '3.2.3',
    currentGrade: 'A',
    targetGrade: 'A',
    applicant: '张三',
    department: '技术研发部',
    status: 'pending',
    reason: '获取安全补丁和新特性支持，修复已知CVE漏洞',
    submittedAt: '2024-03-05',
    breakingChanges: [
      '部分旧版API已废弃，需要调整代码',
      '默认嵌入式Tomcat版本升级至10.1.x',
    ],
    securityFixes: [
      'CVE-2024-22243: URL解码漏洞修复',
      'CVE-2024-22259: HTTP/2请求走私修复',
    ],
    newFeatures: [
      '支持虚拟线程（Java 21+）',
      '增强的Observability支持',
      '改进的AOT编译支持',
    ],
    compatibilityRisk: '低风险，已有完整测试覆盖',
    testReport: '单元测试通过率98.5%，集成测试全部通过',
    rollbackPlan: '保留3.1.5版本依赖，升级失败可快速回滚',
    beforeEval: { vulnerabilityCount: 3, healthScore: 82, licenseRisk: '低风险' },
    afterEval: { vulnerabilityCount: 0, healthScore: 90, licenseRisk: '低风险' },
    approvalHistory: [
      { action: '提交升级申请', operator: '张三', opinion: '申请升级Spring Boot至3.2.3获取安全补丁', time: '2024-03-05 10:00' },
    ],
  },
  {
    id: 'ur002',
    applicationNo: 'APP-2024-0322',
    componentName: 'MySQL',
    fromVersion: '8.0.35',
    toVersion: '8.0.36',
    currentGrade: 'A',
    targetGrade: 'A',
    applicant: '李四',
    department: '技术研发部',
    status: 'pending',
    reason: '修复已知安全漏洞，提升数据库稳定性',
    submittedAt: '2024-03-22',
    breakingChanges: [],
    securityFixes: [
      '修复多个远程代码执行漏洞',
      '修复权限提升漏洞',
    ],
    newFeatures: [
      '优化InnoDB性能',
      '改进JSON函数支持',
    ],
    compatibilityRisk: '极低风险，小版本升级',
    testReport: '回归测试全部通过，性能无退化',
    rollbackPlan: '数据库备份已完成，可快速恢复',
    beforeEval: { vulnerabilityCount: 2, healthScore: 88, licenseRisk: '高风险' },
    afterEval: { vulnerabilityCount: 0, healthScore: 92, licenseRisk: '高风险' },
    approvalHistory: [
      { action: '提交升级申请', operator: '李四', opinion: '申请升级MySQL至8.0.36修复安全漏洞', time: '2024-03-22 09:00' },
    ],
  },
  {
    id: 'ur003',
    applicationNo: 'APP-2024-0330',
    componentName: 'Vue.js',
    fromVersion: '3.3.13',
    toVersion: '3.4.21',
    currentGrade: 'B',
    targetGrade: 'A',
    applicant: '赵六',
    department: '基础架构部',
    status: 'processing',
    reason: '获取更好的TypeScript支持和性能优化',
    submittedAt: '2024-03-30',
    breakingChanges: [
      '部分内部API变更，不影响业务代码',
      '编译器行为优化，可能需要调整部分模板写法',
    ],
    securityFixes: [
      '修复SSRF风险（v-model绑定）',
    ],
    newFeatures: [
      'defineModel宏支持',
      '更好的TypeScript类型推导',
      '响应式性能优化',
    ],
    compatibilityRisk: '低风险，已验证主要业务模块',
    testReport: '核心模块测试通过率100%',
    rollbackPlan: '保留旧版构建产物，可快速切换',
    beforeEval: { vulnerabilityCount: 0, healthScore: 85, licenseRisk: '低风险' },
    afterEval: { vulnerabilityCount: 0, healthScore: 92, licenseRisk: '低风险' },
    approvalHistory: [
      { action: '提交升级申请', operator: '赵六', opinion: '申请升级Vue.js至3.4.21获取新特性', time: '2024-03-30 11:00' },
      { action: '初审通过', operator: '王五', opinion: '升级方案合理，进入详细评审', time: '2024-03-31 09:00' },
    ],
  },
]

// ==================== 风险管理 Mock 数据 ====================

// 缺陷台账
export interface DefectLedgerItem {
  id: string
  source: '外部推送' | 'SCA扫描' | '人工上报'
  title: string
  riskLevel: RiskLevel
  riskType: '安全漏洞' | '恶意代码' | '许可证违规' | '供应链风险' | '其他'
  affectedComponent: string
  affectedVersion: string
  description: string
  discoveredDate: string
  status: '待分发' | '已分发' | '处理中' | '已解决' | '已忽略'
  // 分发相关
  dispatchType?: '专员分配' | 'Jira工单'
  dispatchDept?: string
  dispatchTarget?: string
  dispatchTime?: string
  dispatcher?: string
  jiraKey?: string
}

export const mockDefectLedger: DefectLedgerItem[] = [
  { id: 'df001', source: '外部推送', title: 'Log4j2 JNDI 远程代码执行漏洞 (CVE-2024-1234)', riskLevel: 'critical', riskType: '安全漏洞', affectedComponent: 'Log4j2', affectedVersion: '<2.22.1', description: '通过JNDI注入可实现远程代码执行，影响范围广泛', discoveredDate: '2024-03-08', status: '已分发', dispatchType: 'Jira工单', dispatchDept: '技术研发部', dispatchTarget: '李四', dispatchTime: '2024-03-09', dispatcher: '王五', jiraKey: 'JIRA-2401' },
  { id: 'df002', source: 'SCA扫描', title: 'event-stream 依赖链恶意代码检测', riskLevel: 'critical', riskType: '恶意代码', affectedComponent: 'event-stream', affectedVersion: '3.3.6', description: '依赖链中发现flatmap-stream包包含窃取加密货币钱包的恶意代码', discoveredDate: '2024-03-15', status: '待分发', dispatchTarget: '', dispatchDept: '' },
  { id: 'df003', source: '外部推送', title: 'Spring Framework RCE 漏洞预警', riskLevel: 'high', riskType: '安全漏洞', affectedComponent: 'Spring Framework', affectedVersion: '<5.3.18', description: 'Spring框架数据绑定功能存在远程代码执行漏洞', discoveredDate: '2024-02-20', status: '已解决', dispatchType: '专员分配', dispatchDept: '技术研发部', dispatchTarget: '张三', dispatchTime: '2024-02-21', dispatcher: '王五' },
  { id: 'df004', source: 'SCA扫描', title: 'Express 路径穿越漏洞', riskLevel: 'medium', riskType: '安全漏洞', affectedComponent: 'express', affectedVersion: '4.18.2', description: '路径穿越漏洞可导致敏感文件读取', discoveredDate: '2024-03-05', status: '已分发', dispatchType: '专员分配', dispatchDept: '技术研发部', dispatchTarget: '张三', dispatchTime: '2024-03-06', dispatcher: '王五' },
  { id: 'df005', source: 'SCA扫描', title: 'd3.js 原型链污染风险', riskLevel: 'high', riskType: '安全漏洞', affectedComponent: 'd3', affectedVersion: '7.8.5', description: '数据绑定过程中存在原型链污染风险', discoveredDate: '2024-03-08', status: '处理中', dispatchType: 'Jira工单', dispatchDept: '数据智能部', dispatchTarget: '李四', dispatchTime: '2024-03-09', dispatcher: '王五', jiraKey: 'JIRA-2405' },
  { id: 'df006', source: '人工上报', title: 'GPL-3.0 许可证违规使用', riskLevel: 'high', riskType: '许可证违规', affectedComponent: 'readline-sync', affectedVersion: '1.4.10', description: '产品代码中使用了GPL-3.0许可证组件，存在传染性风险', discoveredDate: '2024-03-12', status: '待分发', dispatchTarget: '', dispatchDept: '' },
  { id: 'df007', source: 'SCA扫描', title: 'node-ipc 地缘政治恶意代码', riskLevel: 'critical', riskType: '恶意代码', affectedComponent: 'node-ipc', affectedVersion: '>=11.0.0', description: '作者植入恶意代码，针对特定IP进行文件擦除', discoveredDate: '2024-03-18', status: '已分发', dispatchType: 'Jira工单', dispatchDept: '基础架构部', dispatchTarget: '赵六', dispatchTime: '2024-03-18', dispatcher: '王五', jiraKey: 'JIRA-2408' },
  { id: 'df008', source: '外部推送', title: 'NPM 包供应链投毒风险', riskLevel: 'high', riskType: '供应链风险', affectedComponent: 'coa', affectedVersion: '2.0.2', description: '知名NPM包被投毒，安装时执行恶意脚本', discoveredDate: '2024-03-20', status: '已分发', dispatchType: '专员分配', dispatchDept: '技术研发部', dispatchTarget: '张三', dispatchTime: '2024-03-20', dispatcher: '王五' },
]

// 分发记录
export interface DistributeRecord {
  id: string
  defectId: string
  defectTitle: string
  riskType: string
  distributeTime: string
  distributor: string
  receiveDept: string
  receivePerson: string
  distributeNote: string
  status: '待处理' | '处理中' | '已完成' | '已驳回'
}

export const mockDistributeRecords: DistributeRecord[] = [
  { id: 'dr001', defectId: 'df001', defectTitle: 'Log4j2 JNDI 远程代码执行漏洞', riskType: '安全漏洞', distributeTime: '2024-03-09 10:30', distributor: '王五', receiveDept: '技术研发部', receivePerson: '李四', distributeNote: '请尽快修复，影响生产环境', status: '处理中' },
  { id: 'dr002', defectId: 'df004', defectTitle: 'Express 路径穿越漏洞', riskType: '安全漏洞', distributeTime: '2024-03-06 14:20', distributor: '王五', receiveDept: '技术研发部', receivePerson: '张三', distributeNote: '中风险，安排在下个迭代修复', status: '已完成' },
  { id: 'dr003', defectId: 'df005', defectTitle: 'd3.js 原型链污染风险', riskType: '安全漏洞', distributeTime: '2024-03-09 09:15', distributor: '王五', receiveDept: '数据智能部', receivePerson: '李四', distributeNote: '影响数据可视化模块', status: '处理中' },
  { id: 'dr004', defectId: 'df007', defectTitle: 'node-ipc 地缘政治恶意代码', riskType: '恶意代码', distributeTime: '2024-03-18 16:45', distributor: '王五', receiveDept: '基础架构部', receivePerson: '赵六', distributeNote: '紧急处理，已确认影响CI/CD环境', status: '待处理' },
  { id: 'dr005', defectId: 'df008', defectTitle: 'NPM 包供应链投毒风险', riskType: '供应链风险', distributeTime: '2024-03-20 11:00', distributor: '王五', receiveDept: '技术研发部', receivePerson: '张三', distributeNote: '检查所有项目是否依赖该包', status: '待处理' },
]

// SCA 解析记录
export interface ScaParseRecord {
  id: string
  fileName: string
  fileType: 'SBOM-SPDX' | 'SBOM-CycloneDX' | '源码包'
  status: '解析中' | '解析完成' | '解析失败'
  componentCount: number
  riskCount: number
  licenseRiskCount: number
  parseTime: string
  uploadedBy: string
  result: string
}

export const mockScaParseRecords: ScaParseRecord[] = [
  { id: 'sca001', fileName: 'order-system-sbom.json', fileType: 'SBOM-CycloneDX', status: '解析完成', componentCount: 356, riskCount: 3, licenseRiskCount: 2, parseTime: '2024-03-20 15:30', uploadedBy: '王五', result: '发现3个安全漏洞，2个许可证风险' },
  { id: 'sca002', fileName: 'data-platform-bom.xml', fileType: 'SBOM-SPDX', status: '解析完成', componentCount: 128, riskCount: 1, licenseRiskCount: 0, parseTime: '2024-03-19 10:15', uploadedBy: '张三', result: '发现1个中风险漏洞' },
  { id: 'sca003', fileName: 'auth-service-source.zip', fileType: '源码包', status: '解析中', componentCount: 0, riskCount: 0, licenseRiskCount: 0, parseTime: '2024-03-21 09:00', uploadedBy: '王五', result: '正在解析中...' },
  { id: 'sca004', fileName: 'legacy-sbom.json', fileType: 'SBOM-CycloneDX', status: '解析失败', componentCount: 0, riskCount: 0, licenseRiskCount: 0, parseTime: '2024-03-18 14:20', uploadedBy: '李四', result: '文件格式错误，无法解析' },
]

// 任务列表
export interface RiskTask {
  id: string
  taskNo: string
  title: string
  riskLevel: RiskLevel
  riskType: string
  source: string
  assignee: string
  department: string
  status: '待处理' | '处理中' | '待复核' | '已完成' | '已逾期' | '已忽略'
  progress: number
  createdAt: string
  dueDate: string
  isOverdue: boolean
  isMyReview: boolean // 是否需要自己复核
  jiraKey?: string
}

export const mockRiskTasks: RiskTask[] = [
  { id: 'rt001', taskNo: 'RT-2024-0301', title: 'Log4j2 JNDI 漏洞修复', riskLevel: 'critical', riskType: '安全漏洞', source: '外部推送', assignee: '李四', department: '技术研发部', status: '处理中', progress: 60, createdAt: '2024-03-09', dueDate: '2024-03-15', isOverdue: true, isMyReview: true, jiraKey: 'JIRA-2401' },
  { id: 'rt002', taskNo: 'RT-2024-0305', title: 'Express 路径穿越漏洞修复', riskLevel: 'medium', riskType: '安全漏洞', source: 'SCA扫描', assignee: '张三', department: '技术研发部', status: '已完成', progress: 100, createdAt: '2024-03-06', dueDate: '2024-03-20', isOverdue: false, isMyReview: false, jiraKey: 'JIRA-2403' },
  { id: 'rt003', taskNo: 'RT-2024-0308', title: 'd3.js 原型链污染修复', riskLevel: 'high', riskType: '安全漏洞', source: 'SCA扫描', assignee: '李四', department: '数据智能部', status: '处理中', progress: 40, createdAt: '2024-03-09', dueDate: '2024-03-22', isOverdue: false, isMyReview: true },
  { id: 'rt004', taskNo: 'RT-2024-0315', title: 'event-stream 恶意包处置', riskLevel: 'critical', riskType: '恶意代码', source: 'SCA扫描', assignee: '待分配', department: '技术研发部', status: '待处理', progress: 0, createdAt: '2024-03-15', dueDate: '2024-03-18', isOverdue: true, isMyReview: false },
  { id: 'rt005', taskNo: 'RT-2024-0318', title: 'node-ipc 恶意代码紧急处置', riskLevel: 'critical', riskType: '恶意代码', source: 'SCA扫描', assignee: '赵六', department: '基础架构部', status: '待处理', progress: 0, createdAt: '2024-03-18', dueDate: '2024-03-20', isOverdue: false, isMyReview: true },
  { id: 'rt006', taskNo: 'RT-2024-0320', title: 'NPM 供应链投毒排查', riskLevel: 'high', riskType: '供应链风险', source: '外部推送', assignee: '张三', department: '技术研发部', status: '处理中', progress: 30, createdAt: '2024-03-20', dueDate: '2024-03-25', isOverdue: false, isMyReview: false },
  { id: 'rt007', taskNo: 'RT-2024-0312', title: 'GPL-3.0 许可证违规整改', riskLevel: 'high', riskType: '许可证违规', source: '人工上报', assignee: '待分配', department: '产品技术部', status: '待处理', progress: 0, createdAt: '2024-03-12', dueDate: '2024-03-26', isOverdue: false, isMyReview: true },
  { id: 'rt008', taskNo: 'RT-2024-0225', title: 'Spring4Shell 漏洞修复', riskLevel: 'high', riskType: '安全漏洞', source: '外部推送', assignee: '张三', department: '技术研发部', status: '待复核', progress: 90, createdAt: '2024-02-21', dueDate: '2024-03-01', isOverdue: false, isMyReview: true },
]

// ==================== 数据库表结构对应的 Mock 数据 ====================

// 1. 开源基础软件表 (project)
export interface ProjectItem {
  id: number
  ACTIVE: number
  AUTHOR?: string
  PUBLISHER?: string
  _group: string
  _name: string
  _version: string
  DESCRIPTION?: string
  CLASSIFIER?: string
  license: string
  PURL?: string
  CPE?: string
  SWIDTAGID?: string
  DIRECT_DEPENDENCIES?: string
  LAST_BOM_IMPORTED?: string
  LAST_BOM_IMPORTED_FORMAT?: string
  LAST_RISKSCORE: number
  UUID?: string
  create_time?: string
  create_by?: string
  update_time?: string
  update_by?: string
  type: string
  applicable_level: string
  maint_dept: string
  use_dept: string
  country: string
}

export const mockOpenSourceProjects: ProjectItem[] = [
  { id: 1, ACTIVE: 1, AUTHOR: 'Igor Sysoev', PUBLISHER: 'F5 Networks', _group: 'nginx', _name: 'Nginx', _version: '1.25.4', DESCRIPTION: '高性能HTTP和反向代理服务器', CLASSIFIER: 'middleware', license: 'BSD-2-Clause', PURL: 'pkg:generic/nginx@1.25.4', LAST_RISKSCORE: 0.125, type: 'go', applicable_level: 'A', maint_dept: '技术研发部', use_dept: '技术研发部', country: '俄罗斯', create_time: '2022-01-15 10:00:00', create_by: 'admin' },
  { id: 2, ACTIVE: 1, PUBLISHER: 'Oracle', _group: 'com.mysql', _name: 'MySQL', _version: '8.0.36', DESCRIPTION: '关系型数据库管理系统', CLASSIFIER: 'database', license: 'GPL-2.0', LAST_RISKSCORE: 0.456, type: 'java', applicable_level: 'A', maint_dept: '技术研发部', use_dept: '技术研发部', country: '美国', create_time: '2021-06-01 10:00:00', create_by: 'admin' },
  { id: 3, ACTIVE: 1, PUBLISHER: 'Redis Ltd.', _group: 'redis', _name: 'Redis', _version: '7.2.4', DESCRIPTION: '内存数据结构存储', CLASSIFIER: 'database', license: 'BSD-3-Clause', LAST_RISKSCORE: 0.0, type: 'go', applicable_level: 'A', maint_dept: '技术研发部', use_dept: '平台研发部', country: '美国', create_time: '2021-08-20 10:00:00', create_by: 'admin' },
  { id: 4, ACTIVE: 1, PUBLISHER: 'VMware', _group: 'org.springframework.boot', _name: 'Spring Boot', _version: '3.2.3', DESCRIPTION: 'Spring应用框架', CLASSIFIER: 'middleware', license: 'Apache-2.0', LAST_RISKSCORE: 0.089, type: 'java', applicable_level: 'B', maint_dept: '技术研发部', use_dept: '技术研发部', country: '美国', create_time: '2022-03-10 10:00:00', create_by: 'admin' },
  { id: 5, ACTIVE: 1, AUTHOR: 'Evan You', PUBLISHER: 'Vue Team', _group: 'vue', _name: 'Vue.js', _version: '3.4.21', DESCRIPTION: '渐进式JavaScript框架', CLASSIFIER: 'library', license: 'MIT', LAST_RISKSCORE: 0.0, type: 'python', applicable_level: 'B', maint_dept: '平台研发部', use_dept: '技术研发部', country: '美国', create_time: '2023-01-05 10:00:00', create_by: 'admin' },
  { id: 6, ACTIVE: 1, PUBLISHER: 'Apache', _group: 'org.apache.kafka', _name: 'Kafka', _version: '3.6.1', DESCRIPTION: '分布式事件流平台', CLASSIFIER: 'middleware', license: 'Apache-2.0', LAST_RISKSCORE: 0.234, type: 'java', applicable_level: 'A', maint_dept: '技术研发部', use_dept: '技术研发部', country: '美国', create_time: '2021-11-15 10:00:00', create_by: 'admin' },
  { id: 7, ACTIVE: 0, PUBLISHER: 'Elastic', _group: 'org.elasticsearch', _name: 'Elasticsearch', _version: '8.12.2', DESCRIPTION: '分布式搜索和分析引擎', CLASSIFIER: 'database', license: 'SSPL-1.0', LAST_RISKSCORE: 0.345, type: 'java', applicable_level: 'B', maint_dept: '技术研发部', use_dept: '数据智能部', country: '美国', create_time: '2022-05-20 10:00:00', create_by: 'admin' },
  { id: 8, ACTIVE: 1, PUBLISHER: 'Apache', _group: 'org.apache.logging.log4j', _name: 'Log4j2', _version: '2.22.1', DESCRIPTION: 'Java日志框架', CLASSIFIER: 'library', license: 'Apache-2.0', LAST_RISKSCORE: 0.567, type: 'java', applicable_level: 'C', maint_dept: '技术研发部', use_dept: '技术研发部', country: '美国', create_time: '2022-08-01 10:00:00', create_by: 'admin' },
]

// 2. 自研项目表 (project_inner)
export interface ProjectInnerItem {
  id: number
  pro_code: string
  active: number
  type: string
  is_boundary: number
  group_name?: string
  dept_name?: string
  create_by?: string
  update_by?: string
  module_code: string
  branch: string
  current_ID: string
  create_time?: string
  update_time?: string
}

export const mockProjectInners: ProjectInnerItem[] = [
  { id: 1, pro_code: 'order-system', active: 1, type: 'java', is_boundary: 0, group_name: '订单团队', dept_name: '技术研发部', module_code: 'order-core', branch: 'main', current_ID: 'ORD-20240315-001', create_by: '张三', create_time: '2023-01-15 10:00:00', update_time: '2024-03-15 10:00:00' },
  { id: 2, pro_code: 'payment-system', active: 1, type: 'java', is_boundary: 1, group_name: '支付团队', dept_name: '技术研发部', module_code: 'payment-gateway', branch: 'release/v2.0', current_ID: 'PAY-20240310-001', create_by: '李四', create_time: '2023-03-20 10:00:00', update_time: '2024-03-10 10:00:00' },
  { id: 3, pro_code: 'user-center', active: 1, type: 'java', is_boundary: 1, group_name: '基础平台团队', dept_name: '平台研发部', module_code: 'auth-service', branch: 'main', current_ID: 'UC-20240308-001', create_by: '王五', create_time: '2022-11-10 10:00:00', update_time: '2024-03-08 10:00:00' },
  { id: 4, pro_code: 'message-platform', active: 1, type: 'java', is_boundary: 0, group_name: '基础平台团队', dept_name: '平台研发部', module_code: 'msg-push', branch: 'develop', current_ID: 'MSG-20240301-001', create_by: '赵六', create_time: '2024-02-01 10:00:00', update_time: '2024-03-01 10:00:00' },
  { id: 5, pro_code: 'risk-control', active: 1, type: 'python', is_boundary: 1, group_name: '风控团队', dept_name: '安全部', module_code: 'risk-engine', branch: 'main', current_ID: 'RC-20240312-001', create_by: '钱七', create_time: '2024-04-15 10:00:00', update_time: '2024-03-12 10:00:00' },
  { id: 6, pro_code: 'data-platform', active: 1, type: 'go', is_boundary: 0, group_name: '数据团队', dept_name: '技术研发部', module_code: 'data-analysis', branch: 'main', current_ID: 'DP-20240305-001', create_by: '孙八', create_time: '2022-06-01 10:00:00', update_time: '2024-03-05 10:00:00' },
  { id: 7, pro_code: 'monitor-platform', active: 1, type: 'go', is_boundary: 0, group_name: '运维团队', dept_name: '运维部', module_code: 'monitor-core', branch: 'main', current_ID: 'MP-20240318-001', create_by: '周九', create_time: '2023-08-20 10:00:00', update_time: '2024-03-18 10:00:00' },
  { id: 8, pro_code: 'test-mgmt', active: 0, type: 'java', is_boundary: 0, group_name: '测试团队', dept_name: '质量管理部', module_code: 'test-case', branch: 'archive', current_ID: 'TM-20240101-001', create_by: '吴十', create_time: '2021-05-10 10:00:00', update_time: '2024-01-01 10:00:00' },
]

// 3. 开源依赖包表 (component_snapshot)
export interface ComponentSnapshotItem {
  id: number
  current_id: string
  PARENT_COMPONENT_ID?: number
  AUTHOR?: string
  PUBLISHER?: string
  _group: string
  _name: string
  _version: string
  DESCRIPTION?: string
  CLASSIFIER?: string
  FILENAME?: string
  EXTENSION?: string
  INTERNAL?: number
  MD5?: string
  SHA1?: string
  SHA_256?: string
  PURL?: string
  LICENSE: string
  LICENSE_ID?: number
  LAST_RISKSCORE: number
  PROJECT_ID?: number
}

export const mockComponentSnapshots: ComponentSnapshotItem[] = [
  { id: 1, current_id: 'ORD-20240315-001', _group: 'org.springframework', _name: 'spring-core', _version: '6.1.4', DESCRIPTION: 'Spring核心模块', CLASSIFIER: 'library', LICENSE: 'Apache-2.0', LAST_RISKSCORE: 0.045, INTERNAL: 0, PURL: 'pkg:maven/org.springframework/spring-core@6.1.4' },
  { id: 2, current_id: 'ORD-20240315-001', _group: 'com.fasterxml.jackson.core', _name: 'jackson-databind', _version: '2.16.1', DESCRIPTION: 'JSON数据绑定库', CLASSIFIER: 'library', LICENSE: 'Apache-2.0', LAST_RISKSCORE: 0.123, INTERNAL: 0, PURL: 'pkg:maven/com.fasterxml.jackson.core/jackson-databind@2.16.1' },
  { id: 3, current_id: 'PAY-20240310-001', _group: 'org.apache.commons', _name: 'commons-lang3', _version: '3.14.0', DESCRIPTION: 'Apache通用语言工具', CLASSIFIER: 'library', LICENSE: 'Apache-2.0', LAST_RISKSCORE: 0.0, INTERNAL: 0 },
  { id: 4, current_id: 'PAY-20240310-001', _group: 'com.google.guava', _name: 'guava', _version: '33.0.0', DESCRIPTION: 'Google核心Java库', CLASSIFIER: 'library', LICENSE: 'Apache-2.0', LAST_RISKSCORE: 0.012, INTERNAL: 0 },
  { id: 5, current_id: 'UC-20240308-001', _group: 'org.springframework.boot', _name: 'spring-boot-starter-security', _version: '3.2.3', DESCRIPTION: 'Spring安全启动器', CLASSIFIER: 'library', LICENSE: 'Apache-2.0', LAST_RISKSCORE: 0.067, INTERNAL: 0 },
  { id: 6, current_id: 'MSG-20240301-001', _group: 'org.apache.kafka', _name: 'kafka-clients', _version: '3.6.1', DESCRIPTION: 'Kafka客户端库', CLASSIFIER: 'library', LICENSE: 'Apache-2.0', LAST_RISKSCORE: 0.089, INTERNAL: 0 },
  { id: 7, current_id: 'RC-20240312-001', AUTHOR: 'scikit-learn team', _group: 'sklearn', _name: 'scikit-learn', _version: '1.4.0', DESCRIPTION: '机器学习库', CLASSIFIER: 'library', LICENSE: 'BSD-3-Clause', LAST_RISKSCORE: 0.0, INTERNAL: 0 },
  { id: 8, current_id: 'DP-20240305-001', _group: 'github.com/gin-gonic/gin', _name: 'gin', _version: '1.9.1', DESCRIPTION: 'Go Web框架', CLASSIFIER: 'library', LICENSE: 'MIT', LAST_RISKSCORE: 0.0, INTERNAL: 0 },
  { id: 9, current_id: 'DP-20240305-001', _group: 'github.com/go-sql-driver/mysql', _name: 'mysql', _version: '1.7.1', DESCRIPTION: 'MySQL Go驱动', CLASSIFIER: 'library', LICENSE: 'MPL-2.0', LAST_RISKSCORE: 0.156, INTERNAL: 0 },
  { id: 10, current_id: 'ORD-20240315-001', _group: 'internal', _name: 'order-common', _version: '1.0.0', DESCRIPTION: '订单公共模块', CLASSIFIER: 'library', LICENSE: '内部许可证', LAST_RISKSCORE: 0.0, INTERNAL: 1 },
]

// 4. 黑名单配置表 (component_core)
export interface ComponentCoreItem {
  id: number
  _format: number
  type?: string
  _name: string
  _group: string
  version_exp: string
  applicable_level: string
  support_team: string
  reason?: string
  advise?: string
  _enable: number
  alert_type: number
  create_time?: string
  create_by?: string
  update_time?: string
  update_by?: string
  level_time?: string
  is_boundary: number
}

export const mockComponentCores: ComponentCoreItem[] = [
  { id: 1, _format: 1, type: '安全算法类', _name: 'log4j-core', _group: 'org.apache.logging.log4j', version_exp: '[*,2.22.1)', applicable_level: 'C', support_team: '基础架构团队', reason: '存在远程代码执行漏洞(CVE-2024-1234)', advise: '升级至2.22.1或更高版本', _enable: 1, alert_type: 1, is_boundary: 2, create_by: '王五', create_time: '2024-01-15 10:00:00' },
  { id: 2, _format: 3, type: '通用工具包', _name: 'event-stream', _group: '', version_exp: '[3.3.0,3.3.7)', applicable_level: 'D', support_team: '前端安全团队', reason: '依赖链中植入恶意代码flatmap-stream', advise: '立即移除，使用3.3.7以上安全版本', _enable: 1, alert_type: 1, is_boundary: 2, create_by: '王五', create_time: '2024-02-20 10:00:00' },
  { id: 3, _format: 2, type: '基础框架包', _name: 'requests', _group: '', version_exp: '[*,2.31.0)', applicable_level: 'B', support_team: '数据科学团队', reason: '存在证书验证绕过漏洞', advise: '升级至2.31.0或更高版本', _enable: 1, alert_type: 2, is_boundary: 1, create_by: '王五', create_time: '2024-03-01 10:00:00' },
  { id: 4, _format: 1, type: '序列化工具', _name: 'fastjson', _group: 'com.alibaba', version_exp: '[*,2.0.45)', applicable_level: 'C', support_team: '基础架构团队', reason: '多个反序列化漏洞', advise: '升级至最新版本或替换为Jackson', _enable: 1, alert_type: 1, is_boundary: 2, create_by: '王五', create_time: '2024-01-20 10:00:00' },
  { id: 5, _format: 3, type: '通用工具包', _name: 'ua-parser-js', _group: '', version_exp: '[0.7.29,0.7.30]', applicable_level: 'D', support_team: '前端安全团队', reason: '供应链攻击，植入挖矿脚本', advise: '使用0.7.31或更高版本', _enable: 1, alert_type: 1, is_boundary: 2, create_by: '王五', create_time: '2024-03-01 10:00:00' },
  { id: 6, _format: 1, type: '通用工具包', _name: 'commons-text', _group: 'org.apache.commons', version_exp: '[*,1.10.0)', applicable_level: 'A', support_team: '基础架构团队', reason: '代码注入漏洞CVE-2022-42889', advise: '升级至1.10.0或更高版本', _enable: 0, alert_type: 2, is_boundary: 1, create_by: '王五', create_time: '2024-02-15 10:00:00' },
]

// 5. 黑名单阻断记录表 (component_black)
export interface ComponentBlackItem {
  ID: number
  PROJECT_ID: number
  COMPONENT_ID?: number
  BLACK_ID: number
  PROCESS_STATUS: number
  reason?: string
  jira_no?: string
  pro_time?: string
  test_time?: string
  DESCRIPTION?: string
  PROCESS_PLAN_URL?: number
  PROCESS_RESULT_URL?: number
  create_by?: string
  create_time?: string
  update_by?: string
  update_time?: string
  CURRENT_ID?: string
  internal: number
  alert_type: number
}

export const mockComponentBlacks: ComponentBlackItem[] = [
  { ID: 1, PROJECT_ID: 1, COMPONENT_ID: 2, BLACK_ID: 1, PROCESS_STATUS: 1, reason: '项目依赖log4j-core 2.17.0版本触发黑名单', jira_no: 'JIRA-2401', DESCRIPTION: '正在制定升级计划', create_by: '张三', create_time: '2024-03-15 10:00:00', CURRENT_ID: 'ORD-20240315-001', internal: 0, alert_type: 1 },
  { ID: 2, PROJECT_ID: 2, COMPONENT_ID: 3, BLACK_ID: 2, PROCESS_STATUS: 0, reason: '检测到event-stream 3.3.6版本', jira_no: 'JIRA-2402', create_by: '李四', create_time: '2024-03-10 10:00:00', CURRENT_ID: 'PAY-20240310-001', internal: 0, alert_type: 1 },
  { ID: 3, PROJECT_ID: 4, COMPONENT_ID: 6, BLACK_ID: 5, PROCESS_STATUS: 2, reason: 'ua-parser-js 0.7.30触发阻断', jira_no: 'JIRA-2403', DESCRIPTION: '已升级至0.7.32', create_by: '赵六', create_time: '2024-03-08 10:00:00', update_time: '2024-03-12 10:00:00', CURRENT_ID: 'MSG-20240301-001', internal: 0, alert_type: 1 },
  { ID: 4, PROJECT_ID: 3, COMPONENT_ID: 5, BLACK_ID: 4, PROCESS_STATUS: 3, reason: 'fastjson 1.2.83触发告警', DESCRIPTION: '经确认为误报，已迁移至fastjson2', create_by: '王五', create_time: '2024-03-05 10:00:00', CURRENT_ID: 'UC-20240308-001', internal: 0, alert_type: 2 },
]

// 6. 许可证表 (license) - 更新现有mockLicenses
export interface LicenseDbItem {
  id: number
  _COMMENT?: string
  IS_DEPRECATED: number
  FSFLIBRE: number
  HEADER?: string
  LICENSE: string
  _name: string
  IS_OSIAPPROVED: number
  SEEALSO?: string
  TEMPLATE?: string
  _text?: string
  UUID: string
  create_time?: string
  create_by?: string
  update_time?: string
  update_by?: string
  riskWeight?: number
  license_group?: string
  risk_level?: string
}

export const mockLicenseDb: LicenseDbItem[] = [
  { id: 1, LICENSE: 'MIT', _name: 'MIT License', IS_OSIAPPROVED: 1, IS_DEPRECATED: 0, FSFLIBRE: 1, _COMMENT: '最宽松的开源许可证之一', SEEALSO: 'https://opensource.org/licenses/MIT', _text: 'Permission is hereby granted, free of charge...', UUID: 'uuid-mit-001', create_time: '2022-01-01 10:00:00', riskWeight: 10, license_group: 'Permissive', risk_level: '低风险' },
  { id: 2, LICENSE: 'Apache-2.0', _name: 'Apache License 2.0', IS_OSIAPPROVED: 1, IS_DEPRECATED: 0, FSFLIBRE: 1, _COMMENT: 'Apache基金会开源许可证', SEEALSO: 'https://opensource.org/licenses/Apache-2.0', UUID: 'uuid-apache-002', create_time: '2022-01-01 10:00:00', riskWeight: 15, license_group: 'Permissive', risk_level: '低风险' },
  { id: 3, LICENSE: 'BSD-2-Clause', _name: 'BSD 2-Clause "Simplified" License', IS_OSIAPPROVED: 1, IS_DEPRECATED: 0, FSFLIBRE: 1, _COMMENT: 'BSD简化版许可证', UUID: 'uuid-bsd2-003', create_time: '2022-01-01 10:00:00', riskWeight: 10, license_group: 'Permissive', risk_level: '低风险' },
  { id: 4, LICENSE: 'BSD-3-Clause', _name: 'BSD 3-Clause "New" or "Revised" License', IS_OSIAPPROVED: 1, IS_DEPRECATED: 0, FSFLIBRE: 1, _COMMENT: 'BSD标准版许可证', UUID: 'uuid-bsd3-004', create_time: '2022-01-01 10:00:00', riskWeight: 12, license_group: 'Permissive', risk_level: '低风险' },
  { id: 5, LICENSE: 'GPL-2.0', _name: 'GNU General Public License 2.0', IS_OSIAPPROVED: 1, IS_DEPRECATED: 0, FSFLIBRE: 1, _COMMENT: 'GPL标准许可证v2', SEEALSO: 'https://opensource.org/licenses/GPL-2.0', UUID: 'uuid-gpl2-005', create_time: '2022-01-01 10:00:00', riskWeight: 75, license_group: 'StrongCopyleft', risk_level: '高风险' },
  { id: 6, LICENSE: 'GPL-3.0', _name: 'GNU General Public License 3.0', IS_OSIAPPROVED: 1, IS_DEPRECATED: 0, FSFLIBRE: 1, _COMMENT: 'GPL最新版本许可证', SEEALSO: 'https://opensource.org/licenses/GPL-3.0', UUID: 'uuid-gpl3-006', create_time: '2022-01-01 10:00:00', riskWeight: 80, license_group: 'StrongCopyleft', risk_level: '高风险' },
  { id: 7, LICENSE: 'LGPL-2.1', _name: 'GNU Lesser General Public License 2.1', IS_OSIAPPROVED: 1, IS_DEPRECATED: 0, FSFLIBRE: 1, _COMMENT: 'GNU宽松通用公共许可证', UUID: 'uuid-lgpl21-007', create_time: '2022-01-01 10:00:00', riskWeight: 45, license_group: 'WeakCopyleft', risk_level: '中风险' },
  { id: 8, LICENSE: 'MPL-2.0', _name: 'Mozilla Public License 2.0', IS_OSIAPPROVED: 1, IS_DEPRECATED: 0, FSFLIBRE: 1, _COMMENT: 'Mozilla公共许可证', UUID: 'uuid-mpl2-008', create_time: '2022-01-01 10:00:00', riskWeight: 35, license_group: 'WeakCopyleft', risk_level: '中风险' },
  { id: 9, LICENSE: 'AGPL-3.0', _name: 'GNU Affero General Public License 3.0', IS_OSIAPPROVED: 1, IS_DEPRECATED: 0, FSFLIBRE: 1, _COMMENT: '网络使用也需开源', UUID: 'uuid-agpl3-009', create_time: '2022-01-01 10:00:00', riskWeight: 90, license_group: 'StrongCopyleft', risk_level: '高风险' },
  { id: 10, LICENSE: 'SSPL-1.0', _name: 'Server Side Public License', IS_OSIAPPROVED: 0, IS_DEPRECATED: 0, FSFLIBRE: 0, _COMMENT: 'MongoDB专用许可证', UUID: 'uuid-sspl-010', create_time: '2022-01-01 10:00:00', riskWeight: 60, license_group: '未分类', risk_level: '中风险' },
  { id: 11, LICENSE: 'ISC', _name: 'ISC License', IS_OSIAPPROVED: 1, IS_DEPRECATED: 0, FSFLIBRE: 1, _COMMENT: 'ISC互联网系统协会许可证', UUID: 'uuid-isc-011', create_time: '2022-01-01 10:00:00', riskWeight: 10, license_group: 'Permissive', risk_level: '低风险' },
]

// 7. 许可证分组表 (licensegroup)
export interface LicenseGroupItem {
  id: number
  _name: string
  RISKWEIGHT: number
  UUID: string
  create_time?: string
  create_by?: string
  update_time?: string
  update_by?: string
}

export const mockLicenseGroups: LicenseGroupItem[] = [
  { id: 1, _name: '宽松型许可证', RISKWEIGHT: 10, UUID: 'uuid-lg-001', create_time: '2022-01-01 10:00:00' },
  { id: 2, _name: '弱传染性许可证', RISKWEIGHT: 40, UUID: 'uuid-lg-002', create_time: '2022-01-01 10:00:00' },
  { id: 3, _name: '强传染性许可证', RISKWEIGHT: 80, UUID: 'uuid-lg-003', create_time: '2022-01-01 10:00:00' },
  { id: 4, _name: '特殊许可证', RISKWEIGHT: 60, UUID: 'uuid-lg-004', create_time: '2022-01-01 10:00:00' },
]

// 8. 附件表 (attachment)
export interface AttachmentItem {
  ID: number
  _source: string
  FILENAME: string
  MD5: string
  DESCRIPTION?: string
  STORE_URL: string
  create_by?: string
  create_time?: string
  update_by?: string
  update_time?: string
}

export const mockAttachments: AttachmentItem[] = [
  { ID: 1, _source: 'SBOM导入', FILENAME: 'order-system-sbom.json', MD5: 'a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4', DESCRIPTION: '订单系统SBOM文件', STORE_URL: '/sbom/order-system-sbom.json', create_by: '张三', create_time: '2024-03-15 10:00:00' },
  { ID: 2, _source: 'SCA扫描', FILENAME: 'payment-scan-result.zip', MD5: 'b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5', DESCRIPTION: '支付系统扫描结果', STORE_URL: '/sca/payment-scan-result.zip', create_by: '李四', create_time: '2024-03-10 10:00:00' },
  { ID: 3, _source: '手动上传', FILENAME: '处理计划文档.pdf', MD5: 'c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6', DESCRIPTION: '漏洞处理计划', STORE_URL: '/docs/处理计划文档.pdf', create_by: '王五', create_time: '2024-03-08 10:00:00' },
]
