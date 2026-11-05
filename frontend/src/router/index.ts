import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

// 各菜单模块路由（router/modules/*，对应参考项目 router/modules 每菜单一个 js，index 聚合）
// 菜单模块按业务划分（我的任务/我的运维/台账管理/风险管理...），角色不参与文件组织：
// 角色由统一权限系统动态定义（可自由关联菜单/按钮资源），前端按「可见菜单 path 集合」过滤注册。
import { myTaskRoutes, myTaskHiddenRoutes } from './modules/mytask'
import { myOpsRoutes, myOpsHiddenRoutes } from './modules/myops'
import { myUsageRoutes, myUsageHiddenRoutes } from './modules/myusage'
import { myProjectsRoutes, myProjectsHiddenRoutes } from './modules/myprojects'
import { trustedRepoRoutes, trustedRepoHiddenRoutes } from './modules/trustedrepo'
import { personnelRoutes, personnelHiddenRoutes } from './modules/personnel'
import { assetRoutes, assetHiddenRoutes } from './modules/asset'
import { riskRoutes, riskHiddenRoutes } from './modules/risk'
import { governanceRoutes, governanceHiddenRoutes } from './modules/governance'
import { decisionRoutes, decisionHiddenRoutes } from './modules/decision'

const Layout = () => import('@/components/layout/MainLayout.vue')

// ========== 基础路由（静态注册，启动时加载） ==========
const baseRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/login',
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { title: '登录', requiresAuth: false },
  },
  {
    path: '/auth/cas-callback',
    name: 'CasCallback',
    component: () => import('@/views/auth/CasCallback.vue'),
    meta: { title: 'CAS登录', requiresAuth: false },
  },
  {
    path: '/403',
    name: 'Forbidden',
    component: () => import('@/views/auth/Forbidden.vue'),
    meta: { title: '无权限', requiresAuth: false },
  },
  // 动态路由的父级容器（运行时通过 addRoute 填充 children）
  {
    path: '/operator',
    name: 'OperatorLayout',
    component: Layout,
    redirect: '/operator/my-ops',
    meta: { title: '运维工作台', role: 'operator' },
    children: [],
  },
  {
    path: '/specialist',
    name: 'SpecialistLayout',
    component: Layout,
    redirect: '/specialist/my-ops',
    meta: { title: '治理专员工作台', role: 'specialist' },
    children: [],
  },
  {
    path: '/admin',
    name: 'AdminLayout',
    component: Layout,
    redirect: '/admin/asset/software-library',
    meta: { title: '管理控制台', role: 'admin' },
    children: [],
  },
  // 404 必须最后注册（动态路由注册后再添加）
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/auth/NotFound.vue'),
    meta: { title: '页面不存在', requiresAuth: false },
  },
]

// ========== 动态路由映射表（component 路径 -> 实际组件） ==========
// 后端菜单返回的 component 字段对应这里的 key（隐藏路由组件解析用）
const componentMaps: Record<string, () => Promise<any>> = {
  'views/operator/my-assets/MyProjects': () => import('@/views/operator/my-assets/MyProjects.vue'),
  'views/operator/my-assets/ProjectDetail': () => import('@/views/operator/my-assets/ProjectDetail.vue'),
  'views/operator/MyTask': () => import('@/views/operator/MyTask.vue'),
  'views/operator/TaskDetail': () => import('@/views/operator/TaskDetail.vue'),
  'views/operator/MyOps': () => import('@/views/operator/MyOps.vue'),
  'views/operator/MyUsage': () => import('@/views/operator/MyUsage.vue'),
  'views/operator/TrustedRepo': () => import('@/views/operator/TrustedRepo.vue'),
  'views/operator/PersonnelMgmt': () => import('@/views/operator/PersonnelMgmt.vue'),
  'views/operator/software-service/SoftwareService': () => import('@/views/operator/software-service/SoftwareService.vue'),
  // specialist
  'views/specialist/MyTask': () => import('@/views/specialist/MyTask.vue'),
  'views/specialist/MyOps': () => import('@/views/specialist/MyOps.vue'),
  // admin - asset
  'views/admin/asset/SoftwareLibrary': () => import('@/views/admin/asset/SoftwareLibrary.vue'),
  'views/admin/asset/DependencyLibrary': () => import('@/views/admin/asset/DependencyLibrary.vue'),
  'views/admin/asset/ToolLibrary': () => import('@/views/admin/asset/ToolLibrary.vue'),
  'views/admin/asset/SelfDevelopedProjects': () => import('@/views/admin/asset/SelfDevelopedProjects.vue'),
  'views/admin/asset/SourceBackup': () => import('@/views/admin/asset/SourceBackup.vue'),
  'views/admin/asset/TrustedRepository': () => import('@/views/admin/asset/TrustedRepository.vue'),
  // admin - risk
  'views/admin/risk/Blacklist': () => import('@/views/admin/risk/Blacklist.vue'),
  'views/admin/risk/CicdCheck': () => import('@/views/admin/risk/CicdCheck.vue'),
  'views/admin/risk/LicenseLibrary': () => import('@/views/admin/risk/LicenseLibrary.vue'),
  'views/admin/risk/DefectLedger': () => import('@/views/admin/risk/DefectLedger.vue'),
  // admin - governance
  'views/admin/governance/TaskList': () => import('@/views/admin/governance/TaskList.vue'),
  // admin - decision
  'views/admin/decision/repair-efficiency': () => import('@/views/admin/decision/repair-efficiency.vue'),
  'views/admin/decision/core-dependencies': () => import('@/views/admin/decision/core-dependencies.vue'),
}

// ========== 静态业务路由表（按模块拆分，index 聚合） ==========
// 公司参考项目逻辑：router/modules 每模块一个 js（operator/specialist/admin），index.js 聚合；
// 前端所有业务路由静态定义（含 path），登录后按「可见菜单 path 集合」过滤注册。
// 可见 path 集合来源：统一权限系统菜单（/auth/menus）或前端本地角色映射（menuMap.ts 兜底）。
// 注意：path 必须与菜单树/权限系统下发的菜单 path 完全一致，匹配上的路由才注册可见。
export interface StaticRouteDef {
  path: string
  name: string
  component: () => Promise<any>
  title: string
}

// 隐藏路由定义（不在菜单中、始终注册的详情页等路由）
// layout：注册到哪个工作台布局（operator/specialist/admin），对应基础路由中的父路由
export interface HiddenRouteDef {
  path: string
  name: string
  component: string // componentMaps 的 key
  layout?: string // 缺省按 name 前缀推断（Operator*/Specialist*/Admin*）
}

// 聚合各菜单模块可见路由（所有业务路由静态定义，登录后按可见菜单 path 过滤注册）
const staticRoutes: StaticRouteDef[] = [
  ...myProjectsRoutes,
  ...myTaskRoutes,
  ...myOpsRoutes,
  ...myUsageRoutes,
  ...trustedRepoRoutes,
  ...personnelRoutes,
  ...assetRoutes,
  ...riskRoutes,
  ...governanceRoutes,
  ...decisionRoutes,
]

// ========== 隐藏路由（不在菜单显示但需要注册的路由，各菜单模块定义后按 layout 归组） ==========
const hiddenRoutes: Record<string, HiddenRouteDef[]> = {}
const allHiddenRoutes: HiddenRouteDef[] = [
  ...myProjectsHiddenRoutes,
  ...myTaskHiddenRoutes,
  ...myOpsHiddenRoutes,
  ...myUsageHiddenRoutes,
  ...trustedRepoHiddenRoutes,
  ...personnelHiddenRoutes,
  ...assetHiddenRoutes,
  ...riskHiddenRoutes,
  ...governanceHiddenRoutes,
  ...decisionHiddenRoutes,
]
for (const h of allHiddenRoutes) {
  const layout = h.layout || inferLayout(h.name)
  if (!hiddenRoutes[layout]) hiddenRoutes[layout] = []
  hiddenRoutes[layout].push(h)
}

// 从路由 name 前缀推断布局（Operator*/Specialist*/Admin*）
function inferLayout(name: string): string {
  if (name.startsWith('Admin')) return 'admin'
  if (name.startsWith('Specialist')) return 'specialist'
  return 'operator'
}

// ========== 菜单树节点类型（后端返回） ==========
export interface MenuTreeNode {
  id: number
  name: string
  path: string
  component: string
  icon: string
  sort: number
  hidden: boolean
  group: string
  children: MenuTreeNode[]
}

// ========== 根据菜单数据动态注册路由 ==========
export function registerDynamicRoutes(menus: MenuTreeNode[], role: string) {
  // 父路由的 name（对应基础路由中定义的 name）
  const parentNameMap: Record<string, string> = {
    operator: 'OperatorLayout',
    specialist: 'SpecialistLayout',
    admin: 'AdminLayout',
  }
  const parentName = parentNameMap[role] || 'OperatorLayout'
  const parentPath = `/${role}`

  console.log(`[registerDynamicRoutes] role=${role}, parentName=${parentName}, menus=${menus.length}`)

  // 兜底默认路由：角色无任何可见路由时保证默认首页可访问，避免无路由落到 404。
  // 设计原则：无角色/菜单为空不代表不能登录，用户应能进入空工作台（默认普通运维者体验）。
  const fallbackRoutes: Record<string, { path: string; component: string; title: string }> = {
    operator: { path: 'my-ops', component: 'views/operator/MyOps', title: '我的运维' },
    specialist: { path: 'my-ops', component: 'views/specialist/MyOps', title: '我的运维' },
    admin: { path: 'asset/software-library', component: 'views/admin/asset/SoftwareLibrary', title: '软件台账' },
  }

  let registeredCount = 0

  // 收集菜单树中所有可见 path（含子节点；目录节点 path 为空则跳过）
  const visiblePaths = new Set<string>()
  const collectPaths = (nodes: MenuTreeNode[]) => {
    for (const n of nodes) {
      if (n.path) visiblePaths.add(n.path)
      if (n.children?.length) collectPaths(n.children)
    }
  }
  collectPaths(menus)
  // 隐藏路由 path（不受菜单过滤，始终可见）加入集合便于统计，但不影响其独立注册
  hiddenRoutes[role]?.forEach(h => visiblePaths.add(`${parentPath}/${h.path}`))

  // 按「可见菜单 path 集合」过滤注册静态路由（公司逻辑：前端静态路由 + 权限系统菜单 path 匹配）
  for (const def of staticRoutes) {
    if (!visiblePaths.has(def.path)) continue
    const routePath = def.path.replace(`${parentPath}/`, '')
    const route: RouteRecordRaw = {
      path: routePath,
      name: def.name,
      component: def.component,
      meta: { title: def.title, role },
    }
    router.addRoute(parentName, route)
    registeredCount++
    console.log(`[registerDynamicRoutes] 注册路由: parent=${parentName}, path=${routePath}, name=${def.name}`)
  }

  // 注册隐藏路由（详情页等带参路由，不在菜单中，始终注册）
  const hidden = hiddenRoutes[role] || []
  hidden.forEach(h => {
    const route: RouteRecordRaw = {
      path: h.path,
      name: h.name,
      component: componentMaps[h.component],
      meta: { title: h.name, role, hidden: true },
    }
    router.addRoute(parentName, route)
    registeredCount++
  })

  // 菜单树为空（主页面一个都没注册）或没有任何路由注册时，兜底注册默认首页，避免落到 404
  // ============ 兜底分支：角色无可见菜单时生效（对接真实权限系统后可整段删除）============
  if ((menus.length === 0 || registeredCount === 0) && fallbackRoutes[role]) {
    const fb = fallbackRoutes[role]
    const route: RouteRecordRaw = {
      path: fb.path,
      name: `${role}_fallback_home`,
      component: componentMaps[fb.component],
      meta: { title: fb.title, role, fallback: true },
    }
    router.addRoute(parentName, route)
    console.log(`[registerDynamicRoutes] 兜底注册默认首页: parent=${parentName}, path=${fb.path}`)
  }
}

// 清除动态路由（登出或切换角色时调用）
export function resetRouter() {
  // 移除所有动态注册的路由，保留基础路由
  const dynamicNames: string[] = []
  router.getRoutes().forEach(r => {
    if (r.name && (
      staticRoutes.some(s => s.name === r.name) ||
      r.name.toString().endsWith('_fallback_home') ||
      hiddenRoutes.operator?.some(h => h.name === r.name) ||
      hiddenRoutes.specialist?.some(h => h.name === r.name) ||
      hiddenRoutes.admin?.some(h => h.name === r.name)
    )) {
      dynamicNames.push(r.name.toString())
    }
  })
  dynamicNames.forEach(name => router.removeRoute(name))
}

const router = createRouter({
  history: createWebHistory(),
  routes: baseRoutes,
})

// 获取角色对应的默认首页
export function getDefaultHome(role: string | undefined): string {
  switch (role) {
    case 'admin':
      return '/admin/asset/software-library'
    case 'specialist':
      return '/specialist/my-ops'
    case 'operator':
    default:
      return '/operator/my-ops'
  }
}

// 路由守卫
router.beforeEach(async (to, _from, next) => {
  document.title = `${to.meta.title || ''} - 开源治理平台 V3`

  const userStore = useUserStore()

  // 首次进入时从 localStorage 恢复登录状态
  // 注意：token 在 store 初始化时已恢复，但 userInfo 需要 initFromStorage 设置
  if (!userStore.routesLoaded) {
    userStore.initFromStorage()
  }

  // 不需要认证的路由直接放行（登录页、CAS回调、403）
  // 注意：404 路由 requiresAuth=false，但已登录用户匹配到 404 时需要先加载动态路由
  if (to.meta.requiresAuth === false && to.name !== 'NotFound') {
    next()
    return
  }

  // 需要认证但未登录
  if (!userStore.isLoggedIn) {
    if (to.name === 'NotFound') {
      next({ path: '/login' })
    } else {
      next({ path: '/login', query: { redirect: to.fullPath } })
    }
    return
  }

  // 已登录但尚未加载动态路由 -> 先加载用户信息和菜单
  if (!userStore.routesLoaded) {
    // 先从后端获取最新用户信息（部门、团队等），再加载菜单
    await userStore.fetchUserInfo()
    await userStore.fetchMenus()
    if (!userStore.isLoggedIn) {
      next({ path: '/login', query: { redirect: to.fullPath } })
      return
    }
    // 动态路由注册后重新导航
    // 注意：to.name 可能是 NotFound，必须去掉 name，让 path 重新匹配新注册的动态路由
    next({ path: to.fullPath, replace: true })
    return
  }

  next()
})

export default router
