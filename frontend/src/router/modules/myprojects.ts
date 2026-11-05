// ==================== 菜单模块路由：我的项目（myprojects） ====================
// 对应参考项目 router/modules/myprojects.js：一个菜单模块一个 js，含子菜单/详情路由。
// 注意：path 必须与权限系统下发菜单 / 本地 menuMap 的 path 完全一致（按可见 path 过滤注册）。

import type { StaticRouteDef, HiddenRouteDef } from '../index'

// 可见路由（菜单项）
export const myProjectsRoutes: StaticRouteDef[] = [
  { path: '/operator/my-projects', name: 'OperatorMyProjects', component: () => import('@/views/operator/my-assets/MyProjects.vue'), title: '我的项目' },
  { path: '/specialist/my-projects', name: 'SpecialistMyProjects', component: () => import('@/views/operator/my-assets/MyProjects.vue'), title: '我的项目' },
]

// 隐藏路由（项目详情等带参路由，不在菜单中，始终注册）
export const myProjectsHiddenRoutes: HiddenRouteDef[] = [
  { path: 'project-detail/:id', name: 'OperatorProjectDetail', component: 'views/operator/my-assets/ProjectDetail', layout: 'operator' },
  { path: 'project-detail/:id', name: 'SpecialistProjectDetail', component: 'views/operator/my-assets/ProjectDetail', layout: 'specialist' },
]
