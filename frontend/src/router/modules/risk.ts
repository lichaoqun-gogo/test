// ==================== 菜单模块路由：风险管理（risk） ====================
// 对应参考项目 router/modules/risk.js：一个菜单模块一个 js，含子菜单（名单/CICD/许可证/缺陷台账）。
// 注意：path 必须与权限系统下发菜单 / 本地 menuMap 的 path 完全一致（按可见 path 过滤注册）。

import type { StaticRouteDef, HiddenRouteDef } from '../index'

// 可见路由（风险管理模块下的子菜单项）
export const riskRoutes: StaticRouteDef[] = [
  { path: '/admin/risk/blacklist', name: 'AdminBlacklist', component: () => import('@/views/admin/risk/Blacklist.vue'), title: '名单管理' },
  { path: '/admin/risk/cicd-check', name: 'AdminCicdCheck', component: () => import('@/views/admin/risk/CicdCheck.vue'), title: 'CICD检查' },
  { path: '/admin/risk/license-library', name: 'AdminLicenseLibrary', component: () => import('@/views/admin/risk/LicenseLibrary.vue'), title: '许可证信息库' },
  { path: '/admin/risk/defect-ledger', name: 'AdminDefectLedger', component: () => import('@/views/admin/risk/DefectLedger.vue'), title: '缺陷台账' },
]

// 隐藏路由（无）
export const riskHiddenRoutes: HiddenRouteDef[] = []
