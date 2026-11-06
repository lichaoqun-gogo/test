import type { UserInfo } from '@/types'
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

// 登出：调用后端 /api/v1/auth/logout。
// CAS 模式下后端会 302 重定向到 CAS 登出页（清除 CAS SSO 会话）再跳回前端；
// 因此登出跳转必须用整页导航（location），不能用 axios（axios 会跟随 302 但留在页面内）。
export function casLogout(): void {
  window.location.href = '/api/v1/auth/logout'
}

// Token 管理
// CAS 登录后后端将 JWT 写入 HttpOnly Cookie（osm_token，JS 读不到），浏览器同域自动携带。
// 前端 axios 无需手动加 token——后端 AuthMiddleware 从 Cookie 解析 JWT 完成认证。
// localStorage 仅作为兼容保留（旧版本/开发 local 模式），HttpOnly Cookie 方式下 token 不落前端。
export const tokenManager = {
  getToken(): string | null {
    // 优先读 HttpOnly Cookie 标记（无法读值，仅用于判断是否已登录）；实际认证由后端从 Cookie 完成
    // 兼容保留：localStorage 中的旧 token（local 开发模式 / 旧版本升级）
    return localStorage.getItem('osm_token')
  },
  setToken(token: string): void {
    // HttpOnly Cookie 由后端写入，前端不写 Cookie；此处仅兼容 local 开发模式存 localStorage
    localStorage.setItem('osm_token', token)
  },
  removeToken(): void {
    // HttpOnly Cookie 的清除由后端 logout 完成；此处仅清 localStorage 兼容数据
    localStorage.removeItem('osm_token')
  },
  hasCookie(): boolean {
    // 检测非 HttpOnly 的登录标记 Cookie（osm_logged）。HttpOnly 的 osm_token JS 读不到，
    // 只能通过后端登录时同时写入的标记感知登录态（标记不含凭证，XSS 读到也无害）。
    return document.cookie.split(';').some(item => item.trim().startsWith('osm_logged='))
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
