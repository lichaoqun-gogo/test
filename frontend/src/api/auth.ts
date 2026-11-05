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
// 参考接口 /api/user/info 返回 { ok, data:{ account, roles, permissions, ... } }，此处映射为前端 UserInfo。
// account 字段：id=账号、displayName=姓名、groupName=部门、orgCode=团队；role 取 roles 首个基础角色。
export async function casCallback(token: string): Promise<ApiResponse<UserInfo>> {
  try {
    const res = await axios.get('/api/user/info', {
      headers: { Authorization: `Bearer ${token}` },
    })
    if (res.data?.ok === true && res.data.data) {
      const d = res.data.data
      const a = d.account || {}
      const info: UserInfo = {
        id: String(a.id ?? ''),
        username: String(a.id ?? ''),
        nickname: a.displayName || '',
        role: (d.roles && d.roles[0]) || 'operator',
        department: a.groupName || '',
        team: a.orgCode || '',
        email: a.email || '',
        roles: d.roles || [],
        permissions: d.permissions || [],
      }
      return { code: 200, message: 'success', data: info }
    }
    // /api/user/info 返回业务错误（如未识别到用户），保留后端消息便于定位
    const bizMsg = res.data?.message || '获取用户信息失败（/api/user/info 返回异常）'
    console.error('[casCallback] /api/user/info 业务失败:', bizMsg, res.data)
    return { code: 401, message: bizMsg, data: null as any }
  } catch (err: any) {
    // 网络/服务异常（接口 404、后端未启动、跨域等），打印完整信息辅助排查
    const status = err?.response?.status
    const netMsg = err?.response?.data?.message || (status ? `认证服务异常（HTTP ${status}）` : '认证服务不可达（/api/user/info 请求失败）')
    console.error('[casCallback] /api/user/info 请求异常:', status, err)
    return { code: 401, message: netMsg, data: null as any }
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
