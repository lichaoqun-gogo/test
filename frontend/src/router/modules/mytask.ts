// ==================== 菜单模块路由：我的任务（mytask） ====================
// 对应参考项目 router/modules/mytask.js：一个菜单模块一个 js，含子菜单/详情路由。
// 同一菜单在多个工作台布局下可存在多个 path（角色动态，由权限系统菜单 path 决定注册哪些）。
// 注意：path 必须与权限系统下发菜单 / 本地 menuMap 的 path 完全一致（按可见 path 过滤注册）。

import type { StaticRouteDef, HiddenRouteDef } from '../index'

// 可见路由（菜单项）
export const myTaskRoutes: StaticRouteDef[] = [
  { path: '/operator/my-task', name: 'OperatorMyTask', component: () => import('@/views/operator/MyTask.vue'), title: '我的任务' },
  { path: '/specialist/my-task', name: 'SpecialistMyTask', component: () => import('@/views/specialist/MyTask.vue'), title: '我的任务' },
]

// 隐藏路由（任务详情等带参路由，不在菜单中，始终注册）
export const myTaskHiddenRoutes: HiddenRouteDef[] = [
  { path: 'my-task/:id', name: 'OperatorTaskDetail', component: 'views/operator/TaskDetail', layout: 'operator' },
  { path: 'my-task/:id', name: 'SpecialistTaskDetail', component: 'views/operator/TaskDetail', layout: 'specialist' },
]
