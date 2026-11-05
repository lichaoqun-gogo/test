// ==================== 菜单模块路由：台账管理（asset） ====================
// 对应参考项目 router/modules/asset.js：一个菜单模块一个 js，含子菜单（基础软件/依赖包/工具等）。
// 注意：path 必须与权限系统下发菜单 / 本地 menuMap 的 path 完全一致（按可见 path 过滤注册）。

import type { StaticRouteDef, HiddenRouteDef } from '../index'

// 可见路由（台账模块下的子菜单项）
export const assetRoutes: StaticRouteDef[] = [
  { path: '/admin/asset/software-library', name: 'AdminSoftwareLibrary', component: () => import('@/views/admin/asset/SoftwareLibrary.vue'), title: '开源基础软件' },
  { path: '/admin/asset/dependency-library', name: 'AdminDependencyLibrary', component: () => import('@/views/admin/asset/DependencyLibrary.vue'), title: '开源依赖包' },
  { path: '/admin/asset/tool-library', name: 'AdminToolLibrary', component: () => import('@/views/admin/asset/ToolLibrary.vue'), title: '开源工具' },
  { path: '/admin/asset/self-developed', name: 'AdminSelfDevelopedProjects', component: () => import('@/views/admin/asset/SelfDevelopedProjects.vue'), title: '自研项目' },
  { path: '/admin/asset/source-backup', name: 'AdminSourceBackup', component: () => import('@/views/admin/asset/SourceBackup.vue'), title: '源码备份' },
  { path: '/admin/trusted-repo', name: 'AdminTrustedRepository', component: () => import('@/views/admin/asset/TrustedRepository.vue'), title: '可信仓库' },
]

// 隐藏路由（无）
export const assetHiddenRoutes: HiddenRouteDef[] = []
