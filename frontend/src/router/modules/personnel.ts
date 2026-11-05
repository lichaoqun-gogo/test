// ==================== 菜单模块路由：人员管理（personnel） ====================
// 对应参考项目 router/modules/personnel.js：一个菜单模块一个 js。
// 注意：path 必须与权限系统下发菜单 / 本地 menuMap 的 path 完全一致（按可见 path 过滤注册）。

import type { StaticRouteDef, HiddenRouteDef } from '../index'

// 可见路由（菜单项）
export const personnelRoutes: StaticRouteDef[] = [
  { path: '/specialist/personnel-mgmt', name: 'SpecialistPersonnelMgmt', component: () => import('@/views/operator/PersonnelMgmt.vue'), title: '人员管理' },
]

// 隐藏路由（无）
export const personnelHiddenRoutes: HiddenRouteDef[] = []
