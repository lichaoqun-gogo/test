// ==================== 菜单模块路由：治理中心（governance） ====================
// 对应参考项目 router/modules/governance.js：一个菜单模块一个 js，含子菜单/详情路由。
// 注意：path 必须与权限系统下发菜单 / 本地 menuMap 的 path 完全一致（按可见 path 过滤注册）。

import type { StaticRouteDef, HiddenRouteDef } from '../index'

// 可见路由（治理中心菜单项）
export const governanceRoutes: StaticRouteDef[] = [
  { path: '/admin/governance/task-list', name: 'AdminGovernanceTaskList', component: () => import('@/views/admin/governance/TaskList.vue'), title: '任务列表' },
]

// 隐藏路由（任务详情，不在菜单中，始终注册）
export const governanceHiddenRoutes: HiddenRouteDef[] = [
  { path: 'governance/task-list/:id', name: 'AdminTaskDetail', component: 'views/operator/TaskDetail', layout: 'admin' },
]
