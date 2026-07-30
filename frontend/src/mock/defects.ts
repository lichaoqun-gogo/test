import type { RiskLevel } from '@/types'

// 缺陷等级标签映射
export const DEFECT_LEVEL_LABELS: Record<string, string> = {
  critical: '严重',
  high: '高',
  medium: '中',
  low: '低',
}

// 缺陷来源
export type DefectSource = 'external' | 'sca'
export const DEFECT_SOURCE_LABELS: Record<DefectSource, string> = {
  external: '外部推送',
  sca: '内部SCA',
}

// 缺陷处理状态
export type DefectStatus = 'pending' | 'assigned' | 'resolved'
export const DEFECT_STATUS_LABELS: Record<DefectStatus, string> = {
  pending: '待分配',
  assigned: '已分配',
  resolved: '已解决',
}

// 分配方式
export type AssignMethod = 'operator' | 'jira'
export const ASSIGN_METHOD_LABELS: Record<AssignMethod, string> = {
  operator: '分配给专员',
  jira: '生成流程工单（Jira）',
}

// 缺陷记录
export interface DefectRecord {
  id: string
  defectNo: string          // 缺陷编号
  cveId: string             // CVE编号
  title: string             // 缺陷标题
  severityLevel: RiskLevel  // 严重等级
  source: DefectSource      // 来源
  component: string         // 关联组件
  componentVersion: string  // 组件版本
  discoveredAt: string      // 发现时间
  status: DefectStatus      // 处理状态
  description: string       // 缺陷描述
  defectType?: 'dependency' | 'software'  // 缺陷类型：依赖包或开源基础软件
}

// 分配记录
export interface AssignRecord {
  id: string
  taskNo: string            // 任务编号
  defectNo: string          // 缺陷编号
  defectTitle: string       // 缺陷标题
  assignee: string          // 接收人
  assignedBy: string        // 分配人
  assignedAt: string        // 分配时间
  status: DefectStatus      // 任务状态
}

// 运维人员列表
export const operators = [
  { label: '张三', value: '张三' },
  { label: '赵六', value: '赵六' },
  { label: '钱七', value: '钱七' },
  { label: '孙八', value: '孙八' },
]

// 缺陷记录假数据
export const mockDefectRecords: DefectRecord[] = [
  {
    id: 'df001',
    defectNo: 'DEF-2024-0301',
    cveId: 'CVE-2024-21762',
    title: 'FortiOS 越权访问漏洞',
    severityLevel: 'critical',
    source: 'external',
    component: 'FortiOS',
    componentVersion: '< 7.4.4',
    discoveredAt: '2024-03-01',
    status: 'pending',
    defectType: 'software',
    description: 'FortiOS 存在越权访问漏洞，攻击者可绕过认证获取系统权限',
  },
  {
    id: 'df002',
    defectNo: 'DEF-2024-0308',
    cveId: 'CVE-2024-29059',
    title: 'Actuator 信息泄露漏洞',
    severityLevel: 'high',
    source: 'sca',
    component: 'Spring Boot Actuator',
    componentVersion: '< 3.2.4',
    discoveredAt: '2024-03-08',
    status: 'pending',
    defectType: 'dependency',
    description: 'Spring Boot Actuator 端点未正确配置导致敏感信息泄露',
  },
  {
    id: 'df003',
    defectNo: 'DEF-2024-0310',
    cveId: 'CVE-2024-22243',
    title: 'Spring Framework RCE 漏洞',
    severityLevel: 'critical',
    source: 'external',
    component: 'Spring Framework',
    componentVersion: '< 5.3.33',
    discoveredAt: '2024-03-10',
    status: 'assigned',
    defectType: 'software',
    description: 'Spring Framework 存在远程代码执行漏洞',
  },
  {
    id: 'df004',
    defectNo: 'DEF-2024-0312',
    cveId: 'CVE-2024-23456',
    title: 'Log4j2 反序列化漏洞',
    severityLevel: 'high',
    source: 'sca',
    component: 'Log4j2',
    componentVersion: '< 2.23.0',
    discoveredAt: '2024-03-12',
    status: 'assigned',
    defectType: 'dependency',
    description: 'Log4j2 存在反序列化漏洞可导致远程代码执行',
  },
  {
    id: 'df005',
    defectNo: 'DEF-2024-0315',
    cveId: 'CVE-2024-34567',
    title: 'Express.js 路径穿越漏洞',
    severityLevel: 'medium',
    source: 'sca',
    component: 'express',
    componentVersion: '< 4.19.2',
    discoveredAt: '2024-03-15',
    status: 'resolved',
    defectType: 'dependency',
    description: 'Express.js 存在路径穿越漏洞可导致文件读取',
  },
  {
    id: 'df006',
    defectNo: 'DEF-2024-0318',
    cveId: 'CVE-2024-45678',
    title: 'Nginx HTTP 请求走私漏洞',
    severityLevel: 'medium',
    source: 'external',
    component: 'Nginx',
    componentVersion: '< 1.25.5',
    discoveredAt: '2024-03-18',
    status: 'pending',
    defectType: 'software',
    description: 'Nginx 存在 HTTP 请求走私漏洞可导致缓存投毒',
  },
  {
    id: 'df007',
    defectNo: 'DEF-2024-0320',
    cveId: 'CVE-2024-56789',
    title: 'lodash 原型链污染',
    severityLevel: 'low',
    source: 'sca',
    component: 'lodash',
    componentVersion: '< 4.17.22',
    discoveredAt: '2024-03-20',
    status: 'pending',
    defectType: 'dependency',
    description: 'lodash 部分方法存在原型链污染风险',
  },
]

// 分配记录假数据
export const mockAssignRecords: AssignRecord[] = [
  {
    id: 'ar001',
    taskNo: 'TSK-2024-0001',
    defectNo: 'DEF-2024-0310',
    defectTitle: 'Spring Framework RCE 漏洞',
    assignee: '张三',
    assignedBy: '李四',
    assignedAt: '2024-03-11 09:30',
    status: 'assigned',
  },
  {
    id: 'ar002',
    taskNo: 'TSK-2024-0002',
    defectNo: 'DEF-2024-0312',
    defectTitle: 'Log4j2 反序列化漏洞',
    assignee: '赵六',
    assignedBy: '李四',
    assignedAt: '2024-03-13 14:15',
    status: 'assigned',
  },
  {
    id: 'ar003',
    taskNo: 'TSK-2024-0003',
    defectNo: 'DEF-2024-0315',
    defectTitle: 'Express.js 路径穿越漏洞',
    assignee: '钱七',
    assignedBy: '李四',
    assignedAt: '2024-03-16 10:00',
    status: 'resolved',
  },
]
