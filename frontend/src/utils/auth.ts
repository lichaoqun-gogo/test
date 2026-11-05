// ==================== 权限判断工具（对应参考项目 src/utils/auth.js） ====================
// 参考项目（osm-front, Vue2）：
//   - hasPriv(code)        操作权限码判断（方法级）
//   - canRouter(route)     路由可访问性判断（路由级）
//   - v-auth 指令           DOM 级权限控制
//   - Vue.prototype.hasPermission(code)  全局方法注入
// 本项目（Vue3）对应实现：
//   - hasPermission(code)  方法级权限判断（复用 user store 的 hasPerm，底层基于 /auth/me 的 permissions）
//   - canRouter(path)      路由级判断（已注册路由 = 可见菜单过滤后的静态路由，未注册则不可访问）
//   - v-auth 指令          已实现（src/directives/auth.ts，见 main.ts 注册）
//   - installPermission(app) 全局方法注入（main.ts 调用，等价 Vue.prototype.hasPermission）
// 说明：本文件与参考项目 utils/auth.js 职责对齐，页面/工具文件可对位互换。

import { useUserStore } from '@/stores/user'
import router from '@/router'

/**
 * 方法级权限判断：判断当前用户是否拥有某权限点（如 global_governance_view）。
 * 权限点来自 /auth/me 返回的 permissions（统一权限系统下发）；admin 角色默认全权限。
 */
export function hasPermission(code: string): boolean {
  const store = useUserStore()
  return store.hasPerm(code)
}

/**
 * 路由级权限判断：判断目标路径是否可访问。
 * 已注册的路由 = 登录后按「可见菜单 path 集合」过滤注册的静态路由；未注册即视为未授权。
 */
export function canRouter(path: string): boolean {
  return router.getRoutes().some(r => r.path === path)
}

/**
 * 全局注入 hasPermission（等价参考项目 Vue.prototype.hasPermission）。
 * 在 main.ts 调用：installPermission(app)
 */
export function installPermission(app: any): void {
  app.config.globalProperties.hasPermission = (code: string) => hasPermission(code)
}

export default { hasPermission, canRouter, installPermission }
