import type { UserInfo, RoleType } from '@/types'

// 模拟用户列表（供其他 mock 模块引用）
export const mockUsers: Record<string, UserInfo & { password: string }> = {
  // 原有演示账号
  operator: {
    id: 'u001',
    username: 'operator',
    password: '123456',
    nickname: '张三',
    role: 'operator' as RoleType,
    department: '技术研发部',
    team: '交易组',
    email: 'zhangsan@company.com',
  },
  specialist: {
    id: 'u002',
    username: 'specialist',
    password: '123456',
    nickname: '李四',
    role: 'specialist' as RoleType,
    department: '平台研发部',
    team: '支付组',
    email: 'lisi@company.com',
  },
  admin: {
    id: 'u003',
    username: 'admin',
    password: '123456',
    nickname: '王五',
    role: 'admin' as RoleType,
    department: '技术研发部',
    team: '管理部',
    email: 'wangwu@company.com',
  },
  
  // 审批流测试账号
  // 技术研发部专员
  user001: {
    id: 'user001',
    username: 'user001',
    password: '123456',
    nickname: '李四',
    role: 'specialist' as RoleType,
    department: '技术研发部',
    team: '研发组',
    email: 'lisi@tech.com',
  },
  // 平台研发部专员
  user002: {
    id: 'user002',
    username: 'user002',
    password: '123456',
    nickname: '王五',
    role: 'specialist' as RoleType,
    department: '平台研发部',
    team: '平台组',
    email: 'wangwu@platform.com',
  },
  // 运维部专员
  user003: {
    id: 'user003',
    username: 'user003',
    password: '123456',
    nickname: '赵六',
    role: 'specialist' as RoleType,
    department: '运维部',
    team: '运维组',
    email: 'zhaoliu@ops.com',
  },
  // 安全部专员
  user004: {
    id: 'user004',
    username: 'user004',
    password: '123456',
    nickname: '钱七',
    role: 'specialist' as RoleType,
    department: '安全部',
    team: '安全组',
    email: 'qianqi@security.com',
  },
  // 技术研发部使用专员
  user005: {
    id: 'user005',
    username: 'user005',
    password: '123456',
    nickname: '吴十',
    role: 'specialist' as RoleType,
    department: '技术研发部',
    team: '使用组',
    email: 'wushi@tech.com',
  },
  // 平台研发部使用专员
  user006: {
    id: 'user006',
    username: 'user006',
    password: '123456',
    nickname: '郑一',
    role: 'specialist' as RoleType,
    department: '平台研发部',
    team: '使用组',
    email: 'zhengyi@platform.com',
  },
  // 运维部使用专员
  user007: {
    id: 'user007',
    username: 'user007',
    password: '123456',
    nickname: '王二',
    role: 'specialist' as RoleType,
    department: '运维部',
    team: '使用组',
    email: 'wanger@ops.com',
  },
  // 管理部门专员（开源治理办公室）
  user010: {
    id: 'user010',
    username: 'user010',
    password: '123456',
    nickname: '孙八',
    role: 'specialist' as RoleType,
    department: '开源治理办公室',
    team: '治理组',
    email: 'sunba@governance.com',
  },
  user011: {
    id: 'user011',
    username: 'user011',
    password: '123456',
    nickname: '周九',
    role: 'specialist' as RoleType,
    department: '开源治理办公室',
    team: '治理组',
    email: 'zhoujiu@governance.com',
  },
  // 普通用户（技术研发部）
  user100: {
    id: 'user100',
    username: 'user100',
    password: '123456',
    nickname: '张三',
    role: 'operator' as RoleType,
    department: '技术研发部',
    team: '开发组',
    email: 'zhangsan@tech.com',
  },
  // 普通用户（平台研发部）
  user101: {
    id: 'user101',
    username: 'user101',
    password: '123456',
    nickname: '刘洋',
    role: 'operator' as RoleType,
    department: '平台研发部',
    team: '开发组',
    email: 'liuyang@platform.com',
  },
  // 普通用户（运维部）
  user102: {
    id: 'user102',
    username: 'user102',
    password: '123456',
    nickname: '陈明',
    role: 'operator' as RoleType,
    department: '运维部',
    team: '运维组',
    email: 'chenming@ops.com',
  },
}

// 兼容旧版 API（已迁移到 src/api/auth.ts）
export let currentUser: (UserInfo & { password: string }) | null = null

export function login(username: string, password: string): UserInfo | null {
  const user = mockUsers[username]
  if (user && user.password === password) {
    currentUser = user
    const { password: _, ...userInfo } = user
    return userInfo
  }
  return null
}

export function logout() {
  currentUser = null
}

export function getCurrentUser(): UserInfo | null {
  if (!currentUser) return null
  const { password: _, ...userInfo } = currentUser
  return userInfo
}

// 权限菜单配置
export const roleMenus: Record<RoleType, string[]> = {
  operator: [
    '/operator/my-assets',
    '/operator/my-task',
    '/operator/software-service',
  ],
  specialist: [
    '/operator/my-assets',
    '/operator/my-task',
    '/operator/software-service',
    '/operator/defect-records',
  ],
  admin: [
    '/admin/asset-center/software-library',
    '/admin/asset-center/dependency-library',
    '/admin/asset-center/tool-library',
    '/admin/asset-center/backup',
    '/admin/asset-center/trusted-repo',
    '/admin/risk-center/blacklist',
    '/admin/risk-center/license',
    '/admin/risk-center/risk-ledger',
    '/admin/risk-center/ticket-tracking',
    '/admin/governance-center/introduction',
    '/admin/governance-center/withdrawal',
    '/admin/governance-center/upgrade',
    '/admin/governance-center/inventory',
    '/admin/governance-center/violation',
    '/admin/governance-center/policy',
    '/admin/decision-center/dashboard',
    '/admin/decision-center/insights',
  ],
}
