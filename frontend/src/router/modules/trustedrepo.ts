// ==================== 菜单模块路由：可信仓库（trustedrepo） ====================
// 对应参考项目 router/modules/trustedrepo.js：一个菜单模块一个 js。
// 注意：path 必须与权限系统下发菜单 / 本地 menuMap 的 path 完全一致（按可见 path 过滤注册）。

import type { StaticRouteDef, HiddenRouteDef } from '../index'

// 可见路由（菜单项）
export const trustedRepoRoutes: StaticRouteDef[] = [
  { path: '/operator/trusted-repo', name: 'OperatorTrustedRepo', component: () => import('@/views/operator/TrustedRepo.vue'), title: '可信仓库' },
  { path: '/specialist/trusted-repo', name: 'SpecialistTrustedRepo', component: () => import('@/views/operator/TrustedRepo.vue'), title: '可信仓库' },
]

// 隐藏路由（无）
export const trustedRepoHiddenRoutes: HiddenRouteDef[] = []
