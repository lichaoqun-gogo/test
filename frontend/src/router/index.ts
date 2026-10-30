import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

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
// 后端菜单返回的 component 字段对应这里的 key
const componentMaps: Record<string, () => Promise<any>> = {
  // operator
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

// ========== 隐藏路由（不在菜单显示但需要注册的路由） ==========
// 这些路由没有对应的菜单项，但需要动态注册
const hiddenRoutes: Record<string, { path: string; name: string; component: string }[]> = {
  operator: [
    { path: 'project-detail/:id', name: 'OperatorProjectDetail', component: 'views/operator/my-assets/ProjectDetail' },
    { path: 'my-task/:id', name: 'OperatorTaskDetail', component: 'views/operator/TaskDetail' },
    { path: 'software-service', name: 'OperatorSoftwareService', component: 'views/operator/software-service/SoftwareService' },
  ],
  specialist: [
    { path: 'project-detail/:id', name: 'SpecialistProjectDetail', component: 'views/operator/my-assets/ProjectDetail' },
    { path: 'my-task/:id', name: 'SpecialistTaskDetail', component: 'views/operator/TaskDetail' },
    { path: 'software-service', name: 'SpecialistSoftwareService', component: 'views/operator/software-service/SoftwareService' },
  ],
  admin: [
    { path: 'governance/task-list/:id', name: 'AdminTaskDetail', component: 'views/operator/TaskDetail' },
  ],
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

  // 递归注册菜单路由
  function registerMenu(menu: MenuTreeNode) {
    if (menu.component && componentMaps[menu.component]) {
      const routePath = menu.path.replace(`${parentPath}/`, '')
      const route: RouteRecordRaw = {
        path: routePath,
        name: `${role}_${menu.path.replace(/\//g, '_')}`,
        component: componentMaps[menu.component],
        meta: { title: menu.name, role },
      }
      router.addRoute(parentName, route)
      console.log(`[registerDynamicRoutes] 注册路由: parent=${parentName}, path=${routePath}, name=${String(route.name)}, fullPath=${menu.path}`)
    } else if (menu.component) {
      console.warn(`[registerDynamicRoutes] 组件未找到: ${menu.component} (menu: ${menu.name})`)
    }
    if (menu.children?.length) {
      menu.children.forEach(child => registerMenu(child))
    }
  }

  menus.forEach(menu => registerMenu(menu))

  // 注册隐藏路由
  const hidden = hiddenRoutes[role] || []
  hidden.forEach(h => {
    const route: RouteRecordRaw = {
      path: h.path,
      name: h.name,
      component: componentMaps[h.component],
      meta: { title: h.name, role, hidden: true },
    }
    router.addRoute(parentName, route)
  })
}

// 清除动态路由（登出或切换角色时调用）
export function resetRouter() {
  // 移除所有动态路由，保留基础路由
  const dynamicNames: string[] = []
  router.getRoutes().forEach(r => {
    if (r.name && (
      r.name.toString().startsWith('operator_/') ||
      r.name.toString().startsWith('specialist_/') ||
      r.name.toString().startsWith('admin_/') ||
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
