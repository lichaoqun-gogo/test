import type { UserInfo, ApiResponse } from '@/types'
import axios from 'axios'

// 获取认证模式
export async function getAuthMode(): Promise<'local' | 'cas'> {
  try {
    const res = await axios.get('/api/v1/auth/mode')
    if (res.data?.mode === 'cas') {
      return 'cas'
    }
  } catch {
    // 接口不可用时默认 local
  }
  return 'local'
}

// CAS 登录：重定向到后端 CAS 入口
export function casLogin(): void {
  window.location.href = '/api/v1/auth/cas/login'
}

// CAS 回调处理：用 token 获取用户信息（含权限档案）
// 后端 /auth/me 返回 { user, roles, permissions }，此处映射为前端 UserInfo
export async function casCallback(token: string): Promise<ApiResponse<UserInfo>> {
  try {
    const res = await axios.get('/api/v1/auth/me', {
      headers: { Authorization: `Bearer ${token}` },
    })
    if (res.data?.code === 0 && res.data.data) {
      const d = res.data.data
      const u = d.user || {}
      const info: UserInfo = {
        id: String(u.id ?? ''),
        username: u.username || '',
        nickname: u.nickname || '',
        role: u.role || 'operator',
        department: u.department || '',
        team: u.team || '',
        email: u.email || '',
        roles: d.roles || [],
        permissions: d.permissions || [],
      }
      return { code: 200, message: 'success', data: info }
    }
    return { code: 401, message: res.data?.message || '获取用户信息失败', data: null as any }
  } catch (err: any) {
    return { code: 401, message: err?.response?.data?.message || '认证失败', data: null as any }
  }
}

// 登出接口（CAS 模式下仅清除本地 token，无需调后端）
export async function logoutApi(): Promise<ApiResponse<null>> {
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
