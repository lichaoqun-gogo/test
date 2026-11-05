// ==================== 菜单模块路由：我的运维（myops） ====================
// 对应参考项目 router/modules/myops.js：一个菜单模块一个 js。
// 注意：path 必须与权限系统下发菜单 / 本地 menuMap 的 path 完全一致（按可见 path 过滤注册）。

import type { StaticRouteDef, HiddenRouteDef } from '../index'

// 可见路由（菜单项）
export const myOpsRoutes: StaticRouteDef[] = [
  { path: '/operator/my-ops', name: 'OperatorMyOps', component: () => import('@/views/operator/MyOps.vue'), title: '我的运维' },
  { path: '/specialist/my-ops', name: 'SpecialistMyOps', component: () => import('@/views/specialist/MyOps.vue'), title: '我的运维' },
]

// 隐藏路由（软件服务申请等子页面，不在菜单中，始终注册）
export const myOpsHiddenRoutes: HiddenRouteDef[] = [
  { path: 'software-service', name: 'OperatorSoftwareService', component: 'views/operator/software-service/SoftwareService', layout: 'operator' },
  { path: 'software-service', name: 'SpecialistSoftwareService', component: 'views/operator/software-service/SoftwareService', layout: 'specialist' },
]
