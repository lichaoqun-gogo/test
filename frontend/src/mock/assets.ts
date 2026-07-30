import type { AssetGrade, SoftwareType, RiskLevel } from '@/types'

// 评级历史记录
export interface GradeHistory {
  grade: AssetGrade
  changedAt: string
  changedBy: string
  reason: string
}

// 软件缺陷信息
export interface SoftwareDefect {
  id: string
  title: string
  severity: RiskLevel
  status: 'open' | 'fixing' | 'fixed' | 'accepted'
  discoveredDate: string
  assignee: string
}

export interface SoftwareAsset {
  id: string
  name: string
  version: string
  type: SoftwareType
  grade: AssetGrade
  language: string
  organization: string
  owner: string
  department: string
  status: 'active' | 'deprecated' | 'pending'
  lastUpdated: string
  bugCount: number
  vulnerabilityCount: number
  description: string
  // 专员分配（仅专员可见）
  assignedSpecialist?: string
  // 评级历史
  gradeHistory?: GradeHistory[]
  // 缺陷信息
  defects?: SoftwareDefect[]
  // 软件描述/维护说明（专员可编辑）
  maintNote?: string
  repoUrl?: string
  officialDocUrl?: string
}

export const mockSoftwareAssets: SoftwareAsset[] = [
  {
    id: 'sw001',
    name: 'Nginx',
    version: '1.25.4',
    type: 'middleware',
    grade: 'A',
    language: 'C',
    organization: 'F5 Networks',
    owner: '张三',
    department: '技术研发部',
    status: 'active',
    lastUpdated: '2024-03-15',
    bugCount: 2,
    vulnerabilityCount: 1,
    description: '高性能HTTP和反向代理服务器',
    assignedSpecialist: '张三',
    repoUrl: 'https://github.com/nginx/nginx',
    officialDocUrl: 'https://nginx.org/en/docs/',
    maintNote: '需要关注HTTP/3支持进展，每季度评估安全补丁',
    gradeHistory: [
      { grade: 'A', changedAt: '2024-03-15', changedBy: '系统', reason: '核心反向代理组件，支撑全部线上业务' },
      { grade: 'B', changedAt: '2023-06-01', changedBy: '张三', reason: '引入初期评估' },
    ],
    defects: [
      { id: 'd001', title: 'HTTP/2 快速重置攻击', severity: 'high', status: 'fixing', discoveredDate: '2024-02-20', assignee: '张三' },
      { id: 'd002', title: '配置解析边界问题', severity: 'medium', status: 'open', discoveredDate: '2024-03-10', assignee: '张三' },
    ],
  },
  {
    id: 'sw002',
    name: 'MySQL',
    version: '8.0.36',
    type: 'database',
    grade: 'A',
    language: 'C++',
    organization: 'Oracle',
    owner: '张三',
    department: '技术研发部',
    status: 'active',
    lastUpdated: '2024-03-10',
    bugCount: 5,
    vulnerabilityCount: 3,
    description: '关系型数据库管理系统',
    assignedSpecialist: '张三',
    repoUrl: 'https://github.com/mysql/mysql-server',
    officialDocUrl: 'https://dev.mysql.com/doc/',
    maintNote: '主业务数据库，需密切关注安全漏洞并及时升级',
    gradeHistory: [
      { grade: 'A', changedAt: '2024-03-10', changedBy: '系统', reason: '核心数据存储，支撑全部交易链路' },
      { grade: 'A', changedAt: '2023-01-15', changedBy: '系统', reason: '初始评级' },
    ],
    defects: [
      { id: 'd003', title: '权限提升漏洞', severity: 'critical', status: 'fixing', discoveredDate: '2024-03-01', assignee: '张三' },
      { id: 'd004', title: '慢查询优化缺陷', severity: 'medium', status: 'open', discoveredDate: '2024-03-05', assignee: '张三' },
      { id: 'd005', title: '备份恢复超时问题', severity: 'low', status: 'accepted', discoveredDate: '2024-02-15', assignee: '张三' },
    ],
  },
  {
    id: 'sw003',
    name: 'Redis',
    version: '7.2.4',
    type: 'database',
    grade: 'A',
    language: 'C',
    organization: 'Redis Ltd.',
    owner: '张三',
    department: '技术研发部',
    status: 'active',
    lastUpdated: '2024-03-08',
    bugCount: 1,
    vulnerabilityCount: 0,
    description: '内存数据结构存储系统',
    assignedSpecialist: '张三',
    repoUrl: 'https://github.com/redis/redis',
    officialDocUrl: 'https://redis.io/documentation',
    maintNote: '缓存核心组件，关注集群模式稳定性',
    gradeHistory: [
      { grade: 'A', changedAt: '2024-03-08', changedBy: '系统', reason: '核心缓存层，影响全站性能' },
    ],
    defects: [
      { id: 'd006', title: '内存碎片率偏高', severity: 'low', status: 'open', discoveredDate: '2024-03-08', assignee: '张三' },
    ],
  },
  {
    id: 'sw004',
    name: 'Spring Boot',
    version: '3.2.3',
    type: 'middleware',
    grade: 'B',
    language: 'Java',
    organization: 'VMware',
    owner: '李四',
    department: '技术研发部',
    status: 'active',
    lastUpdated: '2024-03-12',
    bugCount: 3,
    vulnerabilityCount: 1,
    description: 'Java企业级应用开发框架',
    assignedSpecialist: '李四',
    repoUrl: 'https://github.com/spring-projects/spring-boot',
    officialDocUrl: 'https://docs.spring.io/spring-boot/',
    maintNote: '',
    gradeHistory: [
      { grade: 'B', changedAt: '2024-03-12', changedBy: '系统', reason: 'Java后端核心框架' },
      { grade: 'C', changedAt: '2023-03-01', changedBy: '李四', reason: '初期引入' },
    ],
    defects: [
      { id: 'd007', title: 'Actuator端点暴露风险', severity: 'high', status: 'fixing', discoveredDate: '2024-02-28', assignee: '李四' },
    ],
  },
  {
    id: 'sw005',
    name: 'Vue.js',
    version: '3.4.21',
    type: 'library',
    grade: 'B',
    language: 'TypeScript',
    organization: 'Vue Team',
    owner: '张三',
    department: '技术研发部',
    status: 'active',
    lastUpdated: '2024-03-14',
    bugCount: 2,
    vulnerabilityCount: 0,
    description: '渐进式JavaScript框架',
    assignedSpecialist: '张三',
    repoUrl: 'https://github.com/vuejs/core',
    officialDocUrl: 'https://vuejs.org/',
    maintNote: '前端核心框架，关注v4迁移计划',
    gradeHistory: [
      { grade: 'B', changedAt: '2024-03-14', changedBy: '系统', reason: '前端核心框架' },
    ],
    defects: [],
  },
  {
    id: 'sw006',
    name: 'Kafka',
    version: '3.6.1',
    type: 'middleware',
    grade: 'A',
    language: 'Java/Scala',
    organization: 'Apache',
    owner: '李四',
    department: '技术研发部',
    status: 'active',
    lastUpdated: '2024-02-28',
    bugCount: 4,
    vulnerabilityCount: 2,
    description: '分布式事件流平台',
    assignedSpecialist: '李四',
    repoUrl: 'https://github.com/apache/kafka',
    officialDocUrl: 'https://kafka.apache.org/documentation/',
    maintNote: '消息中间件核心，关注KRaft模式迁移',
    gradeHistory: [
      { grade: 'A', changedAt: '2024-02-28', changedBy: '系统', reason: '全部异步消息通信基础设施' },
      { grade: 'B', changedAt: '2023-06-01', changedBy: '李四', reason: '引入初期' },
    ],
    defects: [
      { id: 'd008', title: '消费者组Rebalance延迟', severity: 'high', status: 'fixing', discoveredDate: '2024-02-15', assignee: '李四' },
      { id: 'd009', title: '磁盘IO瓶颈', severity: 'medium', status: 'open', discoveredDate: '2024-02-28', assignee: '李四' },
    ],
  },
  {
    id: 'sw007',
    name: 'Elasticsearch',
    version: '8.12.2',
    type: 'database',
    grade: 'B',
    language: 'Java',
    organization: 'Elastic',
    owner: '张三',
    department: '技术研发部',
    status: 'active',
    lastUpdated: '2024-03-05',
    bugCount: 6,
    vulnerabilityCount: 2,
    description: '分布式搜索和分析引擎',
    assignedSpecialist: '张三',
    repoUrl: 'https://github.com/elastic/elasticsearch',
    officialDocUrl: 'https://www.elastic.co/guide/en/elasticsearch/reference/',
    maintNote: '日志搜索和全文检索核心',
    gradeHistory: [
      { grade: 'B', changedAt: '2024-03-05', changedBy: '系统', reason: '核心搜索组件' },
    ],
    defects: [
      { id: 'd010', title: 'Lucene索引损坏风险', severity: 'high', status: 'fixing', discoveredDate: '2024-03-01', assignee: '张三' },
      { id: 'd011', title: '聚合查询性能退化', severity: 'medium', status: 'open', discoveredDate: '2024-03-05', assignee: '张三' },
    ],
  },
  {
    id: 'sw008',
    name: 'Log4j2',
    version: '2.22.1',
    type: 'library',
    grade: 'C',
    language: 'Java',
    organization: 'Apache',
    owner: '李四',
    department: '技术研发部',
    status: 'active',
    lastUpdated: '2024-01-20',
    bugCount: 1,
    vulnerabilityCount: 1,
    description: 'Java日志框架',
    assignedSpecialist: '李四',
    repoUrl: 'https://github.com/apache/logging-log4j2',
    officialDocUrl: 'https://logging.apache.org/log4j/2.x/',
    maintNote: '曾受CVE-2021-44228影响，已升级到安全版本',
    gradeHistory: [
      { grade: 'C', changedAt: '2024-01-20', changedBy: '系统', reason: '已修复Log4Shell，降级为一般组件' },
      { grade: 'A', changedAt: '2021-12-10', changedBy: '系统', reason: 'Log4Shell事件，核心关注' },
      { grade: 'B', changedAt: '2023-01-01', changedBy: '张三', reason: '漏洞修复后降级' },
    ],
    defects: [
      { id: 'd012', title: 'JNDI Lookup遗留配置', severity: 'medium', status: 'fixing', discoveredDate: '2024-01-20', assignee: '李四' },
    ],
  },
  {
    id: 'sw009',
    name: 'Docker',
    version: '25.0.3',
    type: 'tool',
    grade: 'A',
    language: 'Go',
    organization: 'Docker Inc.',
    owner: '张三',
    department: '技术研发部',
    status: 'active',
    lastUpdated: '2024-03-01',
    bugCount: 3,
    vulnerabilityCount: 1,
    description: '容器化平台',
    assignedSpecialist: '张三',
    repoUrl: 'https://github.com/moby/moby',
    officialDocUrl: 'https://docs.docker.com/',
    maintNote: '容器化基础设施',
    gradeHistory: [
      { grade: 'A', changedAt: '2024-03-01', changedBy: '系统', reason: '全部容器化运行环境' },
    ],
    defects: [
      { id: 'd013', title: '容器逃逸风险', severity: 'critical', status: 'fixing', discoveredDate: '2024-02-20', assignee: '张三' },
    ],
  },
  {
    id: 'sw010',
    name: 'Prometheus',
    version: '2.50.1',
    type: 'tool',
    grade: 'B',
    language: 'Go',
    organization: 'CNCF',
    owner: '李四',
    department: '技术研发部',
    status: 'active',
    lastUpdated: '2024-03-11',
    bugCount: 2,
    vulnerabilityCount: 0,
    description: '监控系统和时间序列数据库',
    assignedSpecialist: '李四',
    repoUrl: 'https://github.com/prometheus/prometheus',
    officialDocUrl: 'https://prometheus.io/docs/',
    maintNote: '',
    gradeHistory: [
      { grade: 'B', changedAt: '2024-03-11', changedBy: '系统', reason: '核心监控组件' },
    ],
    defects: [],
  },
]

// 自研项目
export interface SelfDevelopedProject {
  id: string
  name: string
  department: string
  team: string
  owner: string
  isBoundary: boolean
  currentCommitId: string
  productionCommitId: string
  licenseCount: number
  dependencyCount: number
  vulnerabilityCount: number
  lastScanned: string
}

export const mockProjects: SelfDevelopedProject[] = [
  {
    id: 'proj001',
    name: '订单管理系统',
    department: '技术研发部',
    team: '业务开发组',
    owner: '张三',
    isBoundary: true,
    currentCommitId: 'a1b2c3d',
    productionCommitId: 'e4f5g6h',
    licenseCount: 12,
    dependencyCount: 45,
    vulnerabilityCount: 3,
    lastScanned: '2024-03-14',
  },
  {
    id: 'proj002',
    name: '用户认证中心',
    department: '技术研发部',
    team: '基础架构组',
    owner: '张三',
    isBoundary: true,
    currentCommitId: 'i7j8k9l',
    productionCommitId: 'i7j8k9l',
    licenseCount: 8,
    dependencyCount: 32,
    vulnerabilityCount: 1,
    lastScanned: '2024-03-13',
  },
  {
    id: 'proj003',
    name: '数据分析平台',
    department: '数据智能部',
    team: '数据分析组',
    owner: '李四',
    isBoundary: false,
    currentCommitId: 'm0n1o2p',
    productionCommitId: 'q3r4s5t',
    licenseCount: 15,
    dependencyCount: 68,
    vulnerabilityCount: 7,
    lastScanned: '2024-03-10',
  },
  {
    id: 'proj004',
    name: 'API网关服务',
    department: '技术研发部',
    team: '基础架构组',
    owner: '张三',
    isBoundary: true,
    currentCommitId: 'u6v7w8x',
    productionCommitId: 'u6v7w8x',
    licenseCount: 6,
    dependencyCount: 28,
    vulnerabilityCount: 0,
    lastScanned: '2024-03-15',
  },
  {
    id: 'proj005',
    name: '消息推送服务',
    department: '技术研发部',
    team: '中间件组',
    owner: '李四',
    isBoundary: true,
    currentCommitId: 'y9z0a1b',
    productionCommitId: 'c2d3e4f',
    licenseCount: 10,
    dependencyCount: 38,
    vulnerabilityCount: 2,
    lastScanned: '2024-03-08',
  },
]

// 项目许可证
export interface ProjectLicense {
  id: string
  projectId: string
  name: string
  riskLevel: RiskLevel
  componentCount: number
}

export const mockLicenses: ProjectLicense[] = [
  { id: 'lic001', projectId: 'proj001', name: 'MIT', riskLevel: 'low', componentCount: 25 },
  { id: 'lic002', projectId: 'proj001', name: 'Apache-2.0', riskLevel: 'low', componentCount: 15 },
  { id: 'lic003', projectId: 'proj001', name: 'GPL-3.0', riskLevel: 'high', componentCount: 3 },
  { id: 'lic004', projectId: 'proj001', name: 'BSD-3-Clause', riskLevel: 'low', componentCount: 2 },
  { id: 'lic005', projectId: 'proj002', name: 'MIT', riskLevel: 'low', componentCount: 20 },
  { id: 'lic006', projectId: 'proj002', name: 'Apache-2.0', riskLevel: 'low', componentCount: 10 },
  { id: 'lic007', projectId: 'proj002', name: 'MPL-2.0', riskLevel: 'medium', componentCount: 2 },
  { id: 'lic008', projectId: 'proj003', name: 'MIT', riskLevel: 'low', componentCount: 30 },
  { id: 'lic009', projectId: 'proj003', name: 'AGPL-3.0', riskLevel: 'critical', componentCount: 2 },
  { id: 'lic010', projectId: 'proj003', name: 'GPL-2.0', riskLevel: 'high', componentCount: 5 },
]

// 项目依赖
export interface ProjectDependency {
  id: string
  projectId: string
  name: string
  version: string
  license: string
  isInBlacklist: boolean
  vulnerabilityCount: number
  lastUpdated: string
}

export const mockDependencies: ProjectDependency[] = [
  { id: 'dep001', projectId: 'proj001', name: 'lodash', version: '4.17.21', license: 'MIT', isInBlacklist: false, vulnerabilityCount: 0, lastUpdated: '2024-01-15' },
  { id: 'dep002', projectId: 'proj001', name: 'axios', version: '1.6.8', license: 'MIT', isInBlacklist: false, vulnerabilityCount: 0, lastUpdated: '2024-03-01' },
  { id: 'dep003', projectId: 'proj001', name: 'log4js', version: '6.9.1', license: 'Apache-2.0', isInBlacklist: true, vulnerabilityCount: 2, lastUpdated: '2023-06-15' },
  { id: 'dep004', projectId: 'proj001', name: 'express', version: '4.18.2', license: 'MIT', isInBlacklist: false, vulnerabilityCount: 1, lastUpdated: '2024-02-20' },
  { id: 'dep005', projectId: 'proj002', name: 'jsonwebtoken', version: '9.0.2', license: 'MIT', isInBlacklist: false, vulnerabilityCount: 0, lastUpdated: '2024-01-10' },
  { id: 'dep006', projectId: 'proj002', name: 'bcrypt', version: '5.1.1', license: 'MIT', isInBlacklist: false, vulnerabilityCount: 0, lastUpdated: '2024-02-15' },
  { id: 'dep007', projectId: 'proj003', name: 'd3', version: '7.8.5', license: 'ISC', isInBlacklist: false, vulnerabilityCount: 0, lastUpdated: '2024-01-20' },
  { id: 'dep008', projectId: 'proj003', name: 'event-stream', version: '3.3.6', license: 'MIT', isInBlacklist: true, vulnerabilityCount: 3, lastUpdated: '2018-09-10' },
]

// 项目漏洞
export interface ProjectVulnerability {
  id: string
  projectId: string
  cveId: string
  severity: RiskLevel
  component: string
  version: string
  description: string
  status: 'open' | 'fixing' | 'fixed' | 'accepted'
  discoveredDate: string
}

export const mockVulnerabilities: ProjectVulnerability[] = [
  { id: 'vul001', projectId: 'proj001', cveId: 'CVE-2024-1001', severity: 'high', component: 'log4js', version: '6.9.1', description: '远程代码执行漏洞', status: 'fixing', discoveredDate: '2024-03-01' },
  { id: 'vul002', projectId: 'proj001', cveId: 'CVE-2024-1002', severity: 'medium', component: 'express', version: '4.18.2', description: '路径穿越漏洞', status: 'open', discoveredDate: '2024-03-05' },
  { id: 'vul003', projectId: 'proj001', cveId: 'CVE-2024-1003', severity: 'low', component: 'express', version: '4.18.2', description: '信息泄露漏洞', status: 'open', discoveredDate: '2024-03-05' },
  { id: 'vul004', projectId: 'proj003', cveId: 'CVE-2024-2001', severity: 'critical', component: 'event-stream', version: '3.3.6', description: '恶意代码注入', status: 'open', discoveredDate: '2024-02-20' },
  { id: 'vul005', projectId: 'proj003', cveId: 'CVE-2024-2002', severity: 'high', component: 'd3', version: '7.8.5', description: '原型链污染漏洞', status: 'fixing', discoveredDate: '2024-03-08' },
  { id: 'vul006', projectId: 'proj003', cveId: 'CVE-2024-2003', severity: 'medium', component: 'webpack', version: '5.88.0', description: '正则表达式拒绝服务', status: 'open', discoveredDate: '2024-03-10' },
]
