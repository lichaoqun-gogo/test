// ==================== 菜单模块路由：我的使用（myusage） ====================
// 对应参考项目 router/modules/myusage.js：一个菜单模块一个 js。
// 注意：path 必须与权限系统下发菜单 / 本地 menuMap 的 path 完全一致（按可见 path 过滤注册）。

import type { StaticRouteDef, HiddenRouteDef } from '../index'

// 可见路由（菜单项）
export const myUsageRoutes: StaticRouteDef[] = [
  { path: '/operator/my-usage', name: 'OperatorMyUsage', component: () => import('@/views/operator/MyUsage.vue'), title: '我的使用' },
  { path: '/specialist/my-usage', name: 'SpecialistMyUsage', component: () => import('@/views/operator/MyUsage.vue'), title: '我的使用' },
]

// 隐藏路由（无）
export const myUsageHiddenRoutes: HiddenRouteDef[] = []
