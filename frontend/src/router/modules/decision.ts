// ==================== 菜单模块路由：决策与洞察（decision） ====================
// 对应参考项目 router/modules/decision.js：一个菜单模块一个 js，含子菜单。
// 注意：path 必须与权限系统下发菜单 / 本地 menuMap 的 path 完全一致（按可见 path 过滤注册）。

import type { StaticRouteDef, HiddenRouteDef } from '../index'

// 可见路由（决策与洞察模块菜单项）
export const decisionRoutes: StaticRouteDef[] = [
  { path: '/admin/decision/repair-efficiency', name: 'AdminRepairEfficiency', component: () => import('@/views/admin/decision/repair-efficiency.vue'), title: '缺陷修复效能' },
  { path: '/admin/decision/core-dependencies', name: 'AdminCoreDependencies', component: () => import('@/views/admin/decision/core-dependencies.vue'), title: '供应链安全洞察' },
]

// 隐藏路由（无）
export const decisionHiddenRoutes: HiddenRouteDef[] = []
