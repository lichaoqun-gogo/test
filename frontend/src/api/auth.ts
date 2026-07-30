import type { UserInfo, RoleType, ApiResponse } from '@/types'
import axios from 'axios'

// 是否使用 CAS 认证（由后端配置决定，前端通过接口获取）
// 默认 local 模式（前端 mock），CAS 模式走后端跳转
let authMode: 'local' | 'cas' = 'local'

// 本地 mock 用户（仅 local 模式使用，开发测试）
interface MockUser extends UserInfo {
  password: string
}

const MOCK_USERS: Record<string, MockUser> = {
  operator: {
    id: 'u001',
    username: 'operator',
    nickname: '张三',
    role: 'operator' as RoleType,
    department: '技术研发部',
    team: '交易组',
    email: 'zhangsan@company.com',
    password: '123456',
  },
  specialist: {
    id: 'u002',
    username: 'specialist',
    nickname: '李四',
    role: 'specialist' as RoleType,
    department: '平台研发部',
    team: '支付组',
    email: 'lisi@company.com',
    password: '123456',
  },
  admin: {
    id: 'u003',
    username: 'admin',
    nickname: '王五',
    role: 'admin' as RoleType,
    department: '技术研发部',
    team: '管理部',
    email: 'wangwu@company.com',
    password: '123456',
  },
  // 审批流测试账号（专员，可切换全局视图）
  user010: {
    id: 'user010',
    username: 'user010',
    nickname: '孙八',
    role: 'specialist' as RoleType,
    department: '开源治理办公室',
    team: '治理组',
    email: 'sunba@governance.com',
    password: '123456',
    canToggleAdminView: true,
  },
  user011: {
    id: 'user011',
    username: 'user011',
    nickname: '周九',
    role: 'specialist' as RoleType,
    department: '开源治理办公室',
    team: '治理组',
    email: 'zhoujiu@governance.com',
    password: '123456',
    canToggleAdminView: true,
  },
}

// 模拟延迟
function delay(ms: number = 300): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms))
}

// 获取认证模式
export async function getAuthMode(): Promise<'local' | 'cas'> {
  try {
    const res = await axios.get('/api/v1/auth/mode')
    if (res.data?.mode === 'cas') {
      authMode = 'cas'
      return 'cas'
    }
  } catch {
    // 接口不可用时默认 local
  }
  return 'local'
}

export function getAuthModeSync(): 'local' | 'cas' {
  return authMode
}

// CAS 登录：重定向到后端 CAS 入口
export function casLogin(): void {
  window.location.href = '/api/v1/auth/cas/login'
}

// CAS 回调处理：用 token 获取用户信息
export async function casCallback(token: string): Promise<ApiResponse<UserInfo>> {
  try {
    const res = await axios.get('/api/v1/auth/me', {
      headers: { Authorization: `Bearer ${token}` },
    })
    if (res.data?.code === 0) {
      return { code: 200, message: 'success', data: res.data.data }
    }
    return { code: 401, message: res.data?.message || '获取用户信息失败', data: null as any }
  } catch (err: any) {
    return { code: 401, message: err?.response?.data?.message || '认证失败', data: null as any }
  }
}

// 本地登录接口（local 模式）
export async function loginApi(username: string, password: string): Promise<ApiResponse<{ token: string; user: UserInfo }>> {
  await delay()

  const user = MOCK_USERS[username]
  if (!user || user.password !== password) {
    return {
      code: 401,
      message: '用户名或密码错误',
      data: null as any,
    }
  }

  const { password: _, ...userInfo } = user
  const token = `mock_token_${Date.now()}_${Math.random().toString(36).slice(2, 9)}`

  return {
    code: 200,
    message: '登录成功',
    data: { token, user: userInfo },
  }
}

// 获取当前用户信息
export async function getUserInfoApi(): Promise<ApiResponse<UserInfo>> {
  await delay()
  const savedUser = localStorage.getItem('osm_user')
  if (!savedUser) {
    return { code: 401, message: '未登录', data: null as any }
  }
  return { code: 200, message: 'success', data: JSON.parse(savedUser) }
}

// 登出接口
export async function logoutApi(): Promise<ApiResponse<null>> {
  await delay()
  localStorage.removeItem('osm_token')
  localStorage.removeItem('osm_user')
  return { code: 200, message: '退出成功', data: null }
}

// Token 管理
export const tokenManager = {
  getToken(): string | null {
    return localStorage.getItem('osm_token')
  },
  setToken(token: string): void {
    localStorage.setItem('osm_token', token)
  },
  removeToken(): void {
    localStorage.removeItem('osm_token')
  },
  getUser(): UserInfo | null {
    const saved = localStorage.getItem('osm_user')
    if (!saved) return null
    try {
      return JSON.parse(saved)
    } catch {
      return null
    }
  },
  setUser(user: UserInfo): void {
    localStorage.setItem('osm_user', JSON.stringify(user))
  },
  removeUser(): void {
    localStorage.removeItem('osm_user')
  },
}
