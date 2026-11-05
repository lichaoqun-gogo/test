import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo, RoleType } from '@/types'
import { logoutApi, tokenManager } from '@/api/auth'
import { registerDynamicRoutes, resetRouter, type MenuTreeNode } from '@/router'
// ============ 兜底依赖：本地角色-菜单映射（权限系统未对接时使用，对接后可删除 import） ============
import { getRoleMenuMap } from '@/router/menuMap'
import { orgApi } from '@/api/org'
import type { OrgDepartment, DeptTeams } from '@/types/org'
import axios from 'axios'

// 全局治理视图权限点编码（与后端 permission 服务约定一致，由权限体系下发）
const PERM_GLOBAL_GOVERNANCE_VIEW = 'global_governance_view'
// 视图模式持久化 key（记住用户上次的 全局治理/个人工作台 选择）
const VIEW_MODE_STORAGE_KEY = 'osm_view_mode'

// 判断用户信息是否具备视图切换能力（admin 角色，或权限体系授予 global_governance_view）
function canToggleByInfo(info?: UserInfo | null): boolean {
  if (!info) return false
  return info.role === 'admin' || (info.permissions || []).includes(PERM_GLOBAL_GOVERNANCE_VIEW)
}

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<UserInfo | null>(null)
  const token = ref<string>(tokenManager.getToken() || '')
  const viewMode = ref<'personal' | 'admin'>('personal')
  // 各角色的菜单树：后端按用户全部角色返回菜单并集（{ <role>: [菜单树] }），
  // 用户可见菜单 = 全部角色菜单的并集，视图切换只是按角色分组展示
  const roleMenus = ref<Record<string, MenuTreeNode[]>>({})
  const routesLoaded = ref(false)

  // 组织架构数据（部门→团队→成员），登录后按需拉取缓存
  const orgDepartments = ref<OrgDepartment[]>([])
  const deptTeams = ref<DeptTeams[]>([])

  // ==================== 参考项目全局初始化数据 ====================
  // App.vue 挂载时通过参考接口（/api/env、/api/groups、/api/user/depts、/api/idcs）加载；
  // 路由守卫通过参考接口 /api/user/info 对位获取完整权限信息（refUserInfo）。
  // 这些数据主要用于对位参考项目流程，权限系统未对接时后端兜底返回/失败，不影响登录。
  const envProfiles = ref<any[]>([])   // 环境配置（/api/env）
  const groupList = ref<any[]>([])     // 组/部门列表（/api/groups）
  const deptTree = ref<any[]>([])      // 部门树（/api/user/depts）
  const deptList = ref<any[]>([])      // 部门列表（/api/user/depts）
  const idcList = ref<any[]>([])       // IDC 列表（/api/idcs）
  const refUserInfo = ref<any | null>(null) // 参考接口 /api/user/info 返回（对位存储）

  function setEnvProfiles(v: any[]) { envProfiles.value = v }
  function setGroupList(v: any[]) { groupList.value = v }
  function setDeptTree(v: any[]) { deptTree.value = v }
  function setDeptList(v: any[]) { deptList.value = v }
  function setIdcList(v: any[]) { idcList.value = v }
  function setRefUserInfo(v: any | null) { refUserInfo.value = v }

  const isLoggedIn = computed(() => !!token.value)
  const role = computed<RoleType | undefined>(() => userInfo.value?.role)
  const isOperator = computed(() => role.value === 'operator')
  const isSpecialist = computed(() => role.value === 'specialist')
  const isAdmin = computed(() => role.value === 'admin')
  // 全局治理视图权限：权限体系下发的 global_governance_view 权限点
  const hasGlobalGovernanceView = computed(() =>
    (userInfo.value?.permissions || []).includes(PERM_GLOBAL_GOVERNANCE_VIEW)
  )
  // admin 角色或具备全局治理权限的用户，可切换 全局治理/个人工作台
  const canToggleView = computed(() => role.value === 'admin' || hasGlobalGovernanceView.value)
  const showAdminView = computed(() => viewMode.value === 'admin')
  // 当前视图显示的菜单：全局治理视图显示 admin 菜单，个人工作台显示自身角色菜单
  //（admin 用户的个人工作台显示 specialist 菜单）
  // 多角色菜单并集（按 path 去重）：把用户全部角色的菜单树合并为侧边栏菜单
  // （公司模型：权限系统按 roles 聚合返回，如 ["common","admin"]，可见菜单 = 各角色菜单并集）
  function mergeRoleMenus(roles: string[]): MenuTreeNode[] {
    const merged: MenuTreeNode[] = []
    const seen = new Set<string>()
    for (const r of roles) {
      const tree = roleMenus.value[r] || []
      for (const node of tree) {
        const key = node.path || `${node.name}-${node.id}`
        if (seen.has(key)) continue
        seen.add(key)
        merged.push(node)
      }
    }
    return merged
  }

  // 当前侧边栏菜单：
  // - 全局治理视图 → 只看 admin 角色菜单
  // - 个人视图 → 全部非 admin 角色菜单并集（admin 菜单归全局治理视图）；
  //   用户仅 admin 角色时回退 specialist 菜单（管理员的个人工作台，原设计约定）
  const viewMenus = computed<MenuTreeNode[]>(() => {
    if (viewMode.value === 'admin') {
      return roleMenus.value['admin'] || []
    }
    const roles = (userInfo.value?.roles?.length ? userInfo.value.roles : [userInfo.value?.role || 'operator'])
      .filter(r => r !== 'admin')
    let menus = mergeRoleMenus(roles)
    if (menus.length === 0 && userInfo.value?.role === 'admin') {
      menus = mergeRoleMenus(['specialist'])
    }
    return menus
  })

  // 按钮权限判断（v-auth 指令底层）：admin 默认全权限，其余看权限点集合
  function hasPerm(code: string): boolean {
    const u = userInfo.value
    if (!u) return false
    if (u.role === 'admin') return true
    return (u.permissions || []).includes(code)
  }

  // 拉取组织架构（部门→团队→成员）并缓存，供人员选择/部门团队展示使用
  async function fetchOrg() {
    if (!token.value) return
    try {
      const [depts, teams] = await Promise.all([orgApi.departments(), orgApi.departmentsTeams()])
      orgDepartments.value = depts || []
      deptTeams.value = teams || []
    } catch (err) {
      console.error('[fetchOrg] 获取组织架构失败:', err)
    }
  }

  // 从后端获取最新用户信息与权限档案（刷新页面时调用，确保部门/权限等信息最新）
  // 参考接口 /api/user/info 返回 { ok, data:{ account, roles, permissions, ... } }
  // account 字段：id=账号、displayName=姓名、groupName=部门、orgCode=团队；role 取 roles 首个基础角色
  async function fetchUserInfo() {
    if (!token.value) return false
    try {
      const res = await axios.get('/api/user/info', {
        headers: { Authorization: `Bearer ${token.value}` },
      })
      if (res.data?.ok === true && res.data.data) {
        const d = res.data.data
        const a = d.account || {}
        const info: UserInfo = {
          id: String(a.id ?? ''),
          username: String(a.id ?? ''),
          nickname: a.displayName || '',
          role: (d.roles && d.roles[0]) || 'operator',
          department: a.groupName || '',
          team: a.orgCode || '',
          email: a.email || '',
          roles: d.roles || [],
          permissions: d.permissions || [],
        }
        userInfo.value = info
        tokenManager.setUser(info)
        // 权限被收回时强制回到个人工作台；仍具备切换能力的用户保留当前视图
        viewMode.value = canToggleByInfo(info) ? viewMode.value : 'personal'
        return true
      }
    } catch (err) {
      console.error('[fetchUserInfo] 获取用户信息失败:', err)
    }
    return false
  }

  // 获取菜单并动态注册路由
  // 参考接口 /api/user/info 的 menusByRole 字段：{ <role>: [菜单树] }（承载原 /auth/menus 职责）
  // 前端一次拉取，全部角色的菜单路由均注册；视图切换（个人工作台/全局治理）只是按角色分组展示
  async function fetchMenus() {
    if (!token.value) return

    try {
      const headers: Record<string, string> = {}
      if (token.value) {
        headers['Authorization'] = `Bearer ${token.value}`
      }
      if (userInfo.value) {
        // X-User-ID 统一使用 CAS 用户名（username），与后端用户唯一标识保持一致
        headers['X-User-ID'] = userInfo.value.username || ''
        headers['X-User-Role'] = userInfo.value.role || ''
        headers['X-User-Name'] = encodeURIComponent(userInfo.value.nickname || '')
        headers['X-User-Dept'] = encodeURIComponent(userInfo.value.department || '')
      }
      const res = await axios.get('/api/user/info', { headers })

      const d = res.data?.data
      if (res.data?.ok === true && d?.menusByRole && typeof d.menusByRole === 'object') {
        const menusByRole: Record<string, MenuTreeNode[]> = d.menusByRole

        // 权限系统未下发任何菜单（menusByRole 为空对象、或所有角色菜单树均为空）
        // → 用前端本地角色映射兜底（menuMap.ts）
        // ============ 兜底分支：权限系统未下发菜单时生效（对接后可整段删除）============
        const allEmpty = Object.keys(menusByRole).length === 0 || Object.values(menusByRole).every(v => !Array.isArray(v) || v.length === 0)
        if (allEmpty) {
          console.warn('[fetchMenus] 权限系统返回空菜单，使用前端本地角色映射（menuMap.ts 兜底）')
          const fallbackRole = userInfo.value?.role || 'operator'
          const fallbackMenus = getRoleMenuMap(fallbackRole)
          roleMenus.value = { [fallbackRole]: fallbackMenus }
          registerDynamicRoutes(fallbackMenus, fallbackRole)
          // admin 角色的个人工作台使用 specialist 菜单（原设计约定），一并注册路由，保证视图切换可访问
          if (fallbackRole === 'admin') {
            const personalMenus = getRoleMenuMap('specialist')
            roleMenus.value['specialist'] = personalMenus
            registerDynamicRoutes(personalMenus, 'specialist')
          }
        } else {
          roleMenus.value = menusByRole
          // 全部角色的菜单路由均注册；角色菜单为空时由 registerDynamicRoutes 兜底注册默认首页路由
          for (const [menuRole, menuTree] of Object.entries(menusByRole)) {
            if (Array.isArray(menuTree) && menuTree.length > 0) {
              registerDynamicRoutes(menuTree, menuRole)
            } else {
              registerDynamicRoutes([], menuRole)
            }
          }
        }
        routesLoaded.value = true
        return
      }
    } catch (err) {
      console.error('获取菜单失败:', err)
    }

    // ============ 兜底分支：权限系统未对接/返回为空时生效 ============
    // 对接真实权限系统后本分支不再进入，可整段删除。
    // 菜单可见性属于"权限展示"而非"认证判定"：菜单拉取失败不登出用户；
    // 按用户角色用前端本地菜单映射（menuMap.ts）提供默认菜单（公司模型：有 app 即有默认角色）。
    console.warn('[fetchMenus] 权限系统未返回菜单，使用前端本地角色映射（menuMap.ts 兜底）')
    if (roleMenus.value && Object.keys(roleMenus.value).length === 0) {
      const fallbackRole = userInfo.value?.role || 'operator'
      const fallbackMenus = getRoleMenuMap(fallbackRole)
      roleMenus.value = { [fallbackRole]: fallbackMenus }
      registerDynamicRoutes(fallbackMenus, fallbackRole)
    }
    routesLoaded.value = true
  }

  // CAS 模式登录（回调页面调用）
  function setCasUser(jwtToken: string, user: UserInfo) {
    token.value = jwtToken
    userInfo.value = user
    tokenManager.setToken(jwtToken)
    tokenManager.setUser(user)
    // admin 默认进入全局治理视图；其他用户（含具备全局治理权限的专员）默认个人工作台
    viewMode.value = user?.role === 'admin' ? 'admin' : 'personal'
    routesLoaded.value = false
  }

  async function logout() {
    await logoutApi()
    userInfo.value = null
    token.value = ''
    roleMenus.value = {}
    routesLoaded.value = false
    resetRouter()
    tokenManager.removeToken()
    tokenManager.removeUser()
    viewMode.value = 'personal'
    localStorage.removeItem(VIEW_MODE_STORAGE_KEY)
  }

  function toggleViewMode() {
    viewMode.value = viewMode.value === 'personal' ? 'admin' : 'personal'
    // 记住上次的视图模式，刷新后恢复（权限仍以 /auth/me 实时返回为准）
    localStorage.setItem(VIEW_MODE_STORAGE_KEY, viewMode.value)
  }

  function initFromStorage() {
    const savedToken = tokenManager.getToken()
    const savedUser = tokenManager.getUser()
    if (savedToken && savedUser) {
      token.value = savedToken
      userInfo.value = savedUser
      // 恢复上次视图模式（仅 admin 或具备全局治理权限的用户有效）
      const savedMode = localStorage.getItem(VIEW_MODE_STORAGE_KEY)
      viewMode.value = canToggleByInfo(savedUser) && savedMode === 'admin' ? 'admin' : 'personal'
      routesLoaded.value = false
    }
  }

  return {
    userInfo,
    token,
    viewMode,
    roleMenus,
    routesLoaded,
    orgDepartments,
    deptTeams,
    envProfiles,
    groupList,
    deptTree,
    deptList,
    idcList,
    refUserInfo,
    setEnvProfiles,
    setGroupList,
    setDeptTree,
    setDeptList,
    setIdcList,
    setRefUserInfo,
    isLoggedIn,
    role,
    isOperator,
    isSpecialist,
    isAdmin,
    hasGlobalGovernanceView,
    canToggleView,
    showAdminView,
    viewMenus,
    hasPerm,
    fetchOrg,
    fetchUserInfo,
    fetchMenus,
    setCasUser,
    logout,
    toggleViewMode,
    initFromStorage,
  }
})
