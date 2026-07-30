import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo, RoleType } from '@/types'
import { loginApi, logoutApi, tokenManager } from '@/api/auth'
import { registerDynamicRoutes, resetRouter, type MenuTreeNode } from '@/router'
import axios from 'axios'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<UserInfo | null>(null)
  const token = ref<string>(tokenManager.getToken() || '')
  const viewMode = ref<'personal' | 'admin'>('personal')
  const menus = ref<MenuTreeNode[]>([])
  const specialistMenus = ref<MenuTreeNode[]>([])
  const routesLoaded = ref(false)

  const isLoggedIn = computed(() => !!token.value)
  const role = computed<RoleType | undefined>(() => userInfo.value?.role)
  const isOperator = computed(() => role.value === 'operator')
  const isSpecialist = computed(() => role.value === 'specialist')
  const isAdmin = computed(() => role.value === 'admin')
  // admin 角色可以切换全局治理/个人工作台
  const canToggleView = computed(() => role.value === 'admin')
  const showAdminView = computed(() => viewMode.value === 'admin')

  // 获取菜单并动态注册路由
  async function fetchMenus() {
    if (!token.value) return

    try {
      const headers: Record<string, string> = {}
      if (token.value && !token.value.startsWith('mock_token_')) {
        headers['Authorization'] = `Bearer ${token.value}`
      }
      if (userInfo.value) {
        headers['X-User-ID'] = userInfo.value.id || ''
        headers['X-User-Role'] = userInfo.value.role || ''
        headers['X-User-Name'] = encodeURIComponent(userInfo.value.nickname || '')
        headers['X-User-Dept'] = encodeURIComponent(userInfo.value.department || '')
      }
      const res = await axios.get('/api/v1/auth/menus', { headers })

      if (res.data?.code === 0 && Array.isArray(res.data.data) && res.data.data.length > 0) {
        menus.value = res.data.data
        const currentRole = role.value || 'operator'
        registerDynamicRoutes(res.data.data, currentRole)

        // admin 用户额外加载 specialist 菜单和路由（用于个人工作台视图）
        if (currentRole === 'admin') {
          try {
            const specialistRes = await axios.get('/api/v1/auth/menus?role=specialist', { headers })
            if (specialistRes.data?.code === 0 && Array.isArray(specialistRes.data.data)) {
              specialistMenus.value = specialistRes.data.data
              registerDynamicRoutes(specialistRes.data.data, 'specialist')
            }
          } catch (err) {
            console.error('加载 specialist 菜单失败:', err)
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
    menus.value = []
    specialistMenus.value = []
    routesLoaded.value = false
    tokenManager.removeToken()
    tokenManager.removeUser()
  }

  // local 模式登录
  async function login(username: string, password: string): Promise<boolean> {
    const res = await loginApi(username, password)
    if (res.code === 200) {
      token.value = res.data.token
      userInfo.value = res.data.user
      tokenManager.setToken(res.data.token)
      tokenManager.setUser(res.data.user)
      viewMode.value = res.data.user?.role === 'admin' ? 'admin' : 'personal'
      routesLoaded.value = false
      return true
    }
    return false
  }

  // CAS 模式登录（回调页面调用）
  function setCasUser(jwtToken: string, user: UserInfo) {
    token.value = jwtToken
    userInfo.value = user
    tokenManager.setToken(jwtToken)
    tokenManager.setUser(user)
    viewMode.value = user?.role === 'admin' ? 'admin' : 'personal'
    routesLoaded.value = false
  }

  async function logout() {
    await logoutApi()
    userInfo.value = null
    token.value = ''
    menus.value = []
    specialistMenus.value = []
    routesLoaded.value = false
    resetRouter()
    tokenManager.removeToken()
    tokenManager.removeUser()
    viewMode.value = 'personal'
  }

  function toggleViewMode() {
    viewMode.value = viewMode.value === 'personal' ? 'admin' : 'personal'
  }

  function initFromStorage() {
    const savedToken = tokenManager.getToken()
    const savedUser = tokenManager.getUser()
    if (savedToken && savedUser) {
      // 忽略旧的 mock_token（local 模式残留），只信任 JWT（CAS 模式）
      if (savedToken.startsWith('mock_token_')) {
        tokenManager.removeToken()
        tokenManager.removeUser()
        return
      }
      token.value = savedToken
      userInfo.value = savedUser
      viewMode.value = savedUser?.role === 'admin' ? 'admin' : 'personal'
      routesLoaded.value = false
    }
  }

  return {
    userInfo,
    token,
    viewMode,
    menus,
    specialistMenus,
    routesLoaded,
    isLoggedIn,
    role,
    isOperator,
    isSpecialist,
    isAdmin,
    canToggleView,
    showAdminView,
    fetchMenus,
    login,
    setCasUser,
    logout,
    toggleViewMode,
    initFromStorage,
  }
})
