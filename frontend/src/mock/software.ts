import type { RiskLevel, AssetGrade } from '@/types'

// 智能推荐组件信息
export interface RecommendedComponent {
  id: string
  name: string
  version: string
  activityScore: number
  vulnerabilityCount: number
  vulnerabilities: { severity: RiskLevel; count: number }[]
  licenseRisk: RiskLevel
  licenseName: string
  description: string
  similarity: number // 相似度百分比
}

// 软件搜索结果（选型辅助）
export interface SoftwareSearchResult {
  id: string
  name: string
  version: string
  description: string
  introduced: boolean
  introducedVersion?: string
  introducedGrade?: AssetGrade
  activityScore: number
  vulnerabilityCount: number
  vulnerabilities: { severity: RiskLevel; count: number }[]
  licenseRisk: RiskLevel
  licenseName: string
  hasRiskWarning: boolean
  riskWarnings?: string[]
  recommendations?: RecommendedComponent[]
}

export const mockSearchResults: SoftwareSearchResult[] = [
  {
    id: 'sr001',
    name: 'Apache Pulsar',
    version: '3.1.2',
    description: '分布式消息流平台，支持多租户、高性能、持久化存储',
    introduced: true,
    introducedVersion: '3.0.0',
    introducedGrade: 'B',
    activityScore: 92,
    vulnerabilityCount: 2,
    vulnerabilities: [
      { severity: 'medium', count: 1 },
      { severity: 'low', count: 1 },
    ],
    licenseRisk: 'low',
    licenseName: 'Apache-2.0',
    hasRiskWarning: false,
  },
  {
    id: 'sr002',
    name: 'RabbitMQ',
    version: '3.13.0',
    description: '开源消息代理软件，实现高级消息队列协议(AMQP)',
    introduced: true,
    introducedVersion: '3.12.10',
    introducedGrade: 'B',
    activityScore: 88,
    vulnerabilityCount: 1,
    vulnerabilities: [{ severity: 'low', count: 1 }],
    licenseRisk: 'low',
    licenseName: 'MPL-2.0',
    hasRiskWarning: false,
  },
  {
    id: 'sr003',
    name: 'RocketMQ',
    version: '5.3.0',
    description: '阿里巴巴开源的分布式消息中间件，低延迟高可靠',
    introduced: false,
    activityScore: 85,
    vulnerabilityCount: 3,
    vulnerabilities: [
      { severity: 'high', count: 1 },
      { severity: 'medium', count: 2 },
    ],
    licenseRisk: 'low',
    licenseName: 'Apache-2.0',
    hasRiskWarning: true,
    riskWarnings: ['存在1个高危漏洞，建议评估后再引入', '同类组件已引入 Apache Pulsar，可优先考虑'],
    recommendations: [
      {
        id: 'rec001',
        name: 'Apache Pulsar',
        version: '3.1.2',
        description: '功能相似的消息流平台，风险更低，已在项目中引入',
        activityScore: 92,
        vulnerabilityCount: 2,
        vulnerabilities: [
          { severity: 'medium', count: 1 },
          { severity: 'low', count: 1 },
        ],
        licenseRisk: 'low',
        licenseName: 'Apache-2.0',
        similarity: 85,
      },
      {
        id: 'rec002',
        name: 'Kafka',
        version: '3.6.1',
        description: '高吞吐量的分布式发布订阅消息系统',
        activityScore: 95,
        vulnerabilityCount: 1,
        vulnerabilities: [{ severity: 'low', count: 1 }],
        licenseRisk: 'low',
        licenseName: 'Apache-2.0',
        similarity: 78,
      },
    ],
  },
  {
    id: 'sr004',
    name: 'event-stream',
    version: '3.3.6',
    description: 'Node.js事件流处理工具，曾被注入恶意代码',
    introduced: false,
    activityScore: 15,
    vulnerabilityCount: 5,
    vulnerabilities: [
      { severity: 'critical', count: 2 },
      { severity: 'high', count: 3 },
    ],
    licenseRisk: 'high',
    licenseName: 'MIT',
    hasRiskWarning: true,
    riskWarnings: ['⚠️ 曾发现供应链安全事件，强烈不建议引入', '存在2个严重漏洞，可能影响系统安全', '活跃度极低，维护状态不明'],
    recommendations: [
      {
        id: 'rec003',
        name: 'lodash',
        version: '4.17.21',
        description: '高性能JavaScript工具库，安全可靠，社区活跃',
        activityScore: 98,
        vulnerabilityCount: 0,
        vulnerabilities: [],
        licenseRisk: 'low',
        licenseName: 'MIT',
        similarity: 72,
      },
      {
        id: 'rec004',
        name: 'RxJS',
        version: '7.8.1',
        description: '响应式编程库，提供强大的事件流处理能力',
        activityScore: 94,
        vulnerabilityCount: 0,
        vulnerabilities: [],
        licenseRisk: 'low',
        licenseName: 'Apache-2.0',
        similarity: 68,
      },
    ],
  },
  {
    id: 'sr005',
    name: 'FastJSON',
    version: '2.0.45',
    description: '阿里巴巴开源的高性能JSON处理库',
    introduced: true,
    introducedVersion: '2.0.40',
    introducedGrade: 'C',
    activityScore: 70,
    vulnerabilityCount: 4,
    vulnerabilities: [
      { severity: 'high', count: 2 },
      { severity: 'medium', count: 2 },
    ],
    licenseRisk: 'medium',
    licenseName: 'Apache-2.0',
    hasRiskWarning: true,
    riskWarnings: ['历史漏洞较多，建议关注安全更新', '许可证风险等级为中等'],
    recommendations: [
      {
        id: 'rec005',
        name: 'Jackson',
        version: '2.16.1',
        description: 'Java生态最流行的JSON处理库，安全性好',
        activityScore: 96,
        vulnerabilityCount: 1,
        vulnerabilities: [{ severity: 'low', count: 1 }],
        licenseRisk: 'low',
        licenseName: 'Apache-2.0',
        similarity: 90,
      },
      {
        id: 'rec006',
        name: 'Gson',
        version: '2.10.1',
        description: 'Google开源的Java JSON序列化/反序列化库',
        activityScore: 82,
        vulnerabilityCount: 0,
        vulnerabilities: [],
        licenseRisk: 'low',
        licenseName: 'Apache-2.0',
        similarity: 85,
      },
    ],
  },
  {
    id: 'sr006',
    name: 'Log4j',
    version: '2.22.1',
    description: 'Apache开源的Java日志框架，企业级应用广泛',
    introduced: true,
    introducedVersion: '2.20.0',
    introducedGrade: 'A',
    activityScore: 96,
    vulnerabilityCount: 1,
    vulnerabilities: [{ severity: 'medium', count: 1 }],
    licenseRisk: 'low',
    licenseName: 'Apache-2.0',
    hasRiskWarning: false,
  },
  {
    id: 'sr007',
    name: 'Struts2',
    version: '6.3.0',
    description: 'Apache开源的Java Web应用框架',
    introduced: false,
    activityScore: 65,
    vulnerabilityCount: 8,
    vulnerabilities: [
      { severity: 'critical', count: 3 },
      { severity: 'high', count: 3 },
      { severity: 'medium', count: 2 },
    ],
    licenseRisk: 'medium',
    licenseName: 'Apache-2.0',
    hasRiskWarning: true,
    riskWarnings: ['⚠️ 历史严重漏洞频发，安全风险极高', '建议采用 Spring Boot 等更安全的替代方案', '活跃度下降明显，社区维护放缓'],
    recommendations: [
      {
        id: 'rec007',
        name: 'Spring Boot',
        version: '3.2.3',
        description: '现代Java Web框架，安全性高，生态完善',
        activityScore: 98,
        vulnerabilityCount: 0,
        vulnerabilities: [],
        licenseRisk: 'low',
        licenseName: 'Apache-2.0',
        similarity: 82,
      },
    ],
  },
]

// 可信包
export interface TrustedPackage {
  id: string
  name: string
  version: string
  category: '中间件' | '数据库' | '依赖包' | '工具'
  scanStatus: 'passed' | 'warning' | 'failed'
  scanDate: string
  downloadCount: number
  size: string
  checksum: string
  description: string
  publisher: string
  lastUpdated: string
}

export const mockTrustedPackages: TrustedPackage[] = [
  {
    id: 'pkg001', name: 'Nginx', version: '1.25.4', category: '中间件',
    scanStatus: 'passed', scanDate: '2024-03-15', downloadCount: 1523, size: '1.2MB',
    checksum: 'sha256:a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6',
    description: '高性能HTTP和反向代理服务器',
    publisher: 'Nginx Inc.', lastUpdated: '2024-03-15',
  },
  {
    id: 'pkg002', name: 'MySQL', version: '8.0.36', category: '数据库',
    scanStatus: 'passed', scanDate: '2024-03-14', downloadCount: 2341, size: '285MB',
    checksum: 'sha256:b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1',
    description: '关系型数据库管理系统',
    publisher: 'Oracle', lastUpdated: '2024-03-14',
  },
  {
    id: 'pkg003', name: 'Redis', version: '7.2.4', category: '数据库',
    scanStatus: 'passed', scanDate: '2024-03-13', downloadCount: 1876, size: '3.2MB',
    checksum: 'sha256:c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2',
    description: '内存数据结构存储系统',
    publisher: 'Redis Ltd.', lastUpdated: '2024-03-13',
  },
  {
    id: 'pkg004', name: 'Spring Boot', version: '3.2.3', category: '依赖包',
    scanStatus: 'passed', scanDate: '2024-03-12', downloadCount: 3245, size: '18MB',
    checksum: 'sha256:d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3',
    description: 'Java企业级应用开发框架',
    publisher: 'VMware', lastUpdated: '2024-03-12',
  },
  {
    id: 'pkg005', name: 'Vue.js', version: '3.4.21', category: '依赖包',
    scanStatus: 'passed', scanDate: '2024-03-14', downloadCount: 2890, size: '1.5MB',
    checksum: 'sha256:e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4',
    description: '渐进式JavaScript前端框架',
    publisher: 'Vue.js Team', lastUpdated: '2024-03-14',
  },
  {
    id: 'pkg006', name: 'Kafka', version: '3.6.1', category: '中间件',
    scanStatus: 'warning', scanDate: '2024-03-10', downloadCount: 1234, size: '75MB',
    checksum: 'sha256:f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5',
    description: '分布式流处理平台',
    publisher: 'Apache', lastUpdated: '2024-03-10',
  },
  {
    id: 'pkg007', name: 'Docker', version: '25.0.3', category: '工具',
    scanStatus: 'passed', scanDate: '2024-03-11', downloadCount: 4521, size: '65MB',
    checksum: 'sha256:a1c3e5g7i9k1m3o5q7s9u1w3y5a1c3e5g7i9k1m3o5q7s9u1w3y5a1c3e5g7',
    description: '容器化应用平台',
    publisher: 'Docker Inc.', lastUpdated: '2024-03-11',
  },
  {
    id: 'pkg008', name: 'Prometheus', version: '2.50.1', category: '工具',
    scanStatus: 'passed', scanDate: '2024-03-11', downloadCount: 987, size: '95MB',
    checksum: 'sha256:b2d4f6h8j0l2n4p6r8t0v2x4z6b2d4f6h8j0l2n4p6r8t0v2x4z6b2d4f6h8',
    description: '开源监控系统和时间序列数据库',
    publisher: 'CNCF', lastUpdated: '2024-03-11',
  },
  {
    id: 'pkg009', name: 'Apache Pulsar', version: '3.1.2', category: '中间件',
    scanStatus: 'passed', scanDate: '2024-03-16', downloadCount: 756, size: '120MB',
    checksum: 'sha256:c3e5g7i9k1m3o5q7s9u1w3y5a1c3e5g7i9k1m3o5q7s9u1w3y5a1c3e5g7i9',
    description: '云原生分布式消息流平台',
    publisher: 'Apache', lastUpdated: '2024-03-16',
  },
  {
    id: 'pkg010', name: 'PostgreSQL', version: '16.2', category: '数据库',
    scanStatus: 'passed', scanDate: '2024-03-15', downloadCount: 1987, size: '312MB',
    checksum: 'sha256:d4f6h8j0l2n4p6r8t0v2x4z6b2d4f6h8j0l2n4p6r8t0v2x4z6b2d4f6h8j0',
    description: '高级开源关系型数据库',
    publisher: 'PostgreSQL Global', lastUpdated: '2024-03-15',
  },
  {
    id: 'pkg011', name: 'Lodash', version: '4.17.21', category: '依赖包',
    scanStatus: 'passed', scanDate: '2024-03-13', downloadCount: 5678, size: '240KB',
    checksum: 'sha256:e5g7i9k1m3o5q7s9u1w3y5a1c3e5g7i9k1m3o5q7s9u1w3y5a1c3e5g7i9k1',
    description: 'JavaScript实用工具库',
    publisher: 'JS Foundation', lastUpdated: '2024-03-13',
  },
  {
    id: 'pkg012', name: 'Grafana', version: '10.3.3', category: '工具',
    scanStatus: 'passed', scanDate: '2024-03-14', downloadCount: 1432, size: '185MB',
    checksum: 'sha256:f6h8j0l2n4p6r8t0v2x4z6b2d4f6h8j0l2n4p6r8t0v2x4z6b2d4f6h8j0l2',
    description: '开源可观测和数据可视化平台',
    publisher: 'Grafana Labs', lastUpdated: '2024-03-14',
  },
]
