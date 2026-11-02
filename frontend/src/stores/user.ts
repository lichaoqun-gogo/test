import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo, RoleType } from '@/types'
import { logoutApi, tokenManager } from '@/api/auth'
import { registerDynamicRoutes, resetRouter, type MenuTreeNode } from '@/router'
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
  const viewMenus = computed<MenuTreeNode[]>(() => {
    if (viewMode.value === 'admin') {
      return roleMenus.value['admin'] || []
    }
    if (role.value === 'admin') {
      return roleMenus.value['specialist'] || []
    }
    return roleMenus.value[role.value || 'operator'] || []
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
  async function fetchUserInfo() {
    if (!token.value) return false
    try {
      const res = await axios.get('/api/v1/auth/me', {
        headers: { Authorization: `Bearer ${token.value}` },
      })
      if (res.data?.code === 0 && res.data.data) {
        const d = res.data.data
        const u = d.user || {}
        const info: UserInfo = {
          id: String(u.id ?? ''),
          username: u.username || '',
          nickname: u.nickname || '',
          role: u.role,
          department: u.department,
          team: u.team || '',
          email: u.email || '',
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
  // 后端按用户全部角色一次返回菜单并集：{ roles: [角色编码], menus: { <role>: [菜单树] } }
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
      const res = await axios.get('/api/v1/auth/menus', { headers })

      const d = res.data?.data
      if (res.data?.code === 0 && d?.menus && typeof d.menus === 'object') {
        const menusByRole: Record<string, MenuTreeNode[]> = d.menus
        roleMenus.value = menusByRole
        // 全部角色的菜单路由均注册
        for (const [menuRole, menuTree] of Object.entries(menusByRole)) {
          if (Array.isArray(menuTree) && menuTree.length > 0) {
            registerDynamicRoutes(menuTree, menuRole)
          }
        }
        routesLoaded.value = true
        return
      }
    } catch (err) {
      console.error('获取菜单失败:', err)
    }

    // 获取菜单失败：清空旧 token，强制重新登录
    console.warn('菜单获取失败，清除认证状态')
    userInfo.value = null
    token.value = ''
    roleMenus.value = {}
    routesLoaded.value = false
    tokenManager.removeToken()
    tokenManager.removeUser()
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
