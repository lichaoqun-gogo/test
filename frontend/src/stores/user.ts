import { defineStore } from 'pinia'
import type { UserInfo, RoleType } from '@/types'
import { casLogout, tokenManager } from '@/api/auth'
import { registerDynamicRoutes, resetRouter, type MenuTreeNode } from '@/router'
// ============ 兜底依赖：本地角色-菜单映射（权限系统未对接时使用，对接后可删除 import） ============
import { getRoleMenuMap } from '@/router/menuMap'
import { orgApi } from '@/api/org'
import type { OrgDepartment, DeptTeams } from '@/types/org'
import axios from 'axios'

// ==================== 参考项目接口请求（对应后端 handlers/reference.go） ====================
// 参考项目（osm-front）的 6 个接口，路径不带 /v1，统一 {ok, data} 结构：
//   1. GET /api/env           环境配置
//   2. GET /api/groups        组/部门列表
//   3. GET /api/user/depts    部门树/列表
//   4. GET /api/idcs          IDC 列表
//   5. GET /api/user/info     当前用户完整权限信息
//   6. GET /api/user/refresh  刷新权限缓存
// 参考接口失败不抛错（返回 null），权限系统未对接时后端兜底返回/失败，登录照常。
const refApi = axios.create({ baseURL: '/api', timeout: 10000 })
refApi.interceptors.response.use(
  response => {
    const res = response.data
    if (res && res.ok === true) return res.data
    console.warn('[reference] 接口返回异常:', res)
    return null
  },
  error => {
    console.warn('[reference] 请求失败（权限系统未对接时后端兜底返回，不影响登录）:', error?.message || error)
    return null
  }
)
const getEnv = () => refApi.get('/env') as Promise<{ profiles: { name: string; active: boolean }[]; tips: string[] } | null>
const getGroups = () => refApi.get('/groups') as Promise<{ code: string; name: string }[] | null>
const getUserDepts = () => refApi.get('/user/depts') as Promise<{ deptTree: any[]; deptList: any[] } | null>
const getIdcs = () => refApi.get('/idcs') as Promise<any[] | null>
// headers 用于携带 Authorization Bearer token / X-User-* 等认证头
const getUserInfo = (headers?: Record<string, string>) =>
  refApi.get('/user/info', headers ? { headers } : undefined) as Promise<any | null>

// 全局治理视图权限点编码（与后端 permission 服务约定一致，由权限体系下发）
const PERM_GLOBAL_GOVERNANCE_VIEW = 'global_governance_view'
// 视图模式持久化 key（记住用户上次的 全局治理/个人工作台 选择）
const VIEW_MODE_STORAGE_KEY = 'osm_view_mode'

// 判断用户信息是否具备视图切换能力（admin 角色，或权限体系授予 global_governance_view）
function canToggleByInfo(info?: UserInfo | null): boolean {
  if (!info) return false
  return info.role === 'admin' || (info.permissions || []).includes(PERM_GLOBAL_GOVERNANCE_VIEW)
}

// ==================== user store（对应参考项目 store/index.js 的 state/mutations/actions） ====================
// 说明：Pinia 选项式写法，state 字段结构对齐参考项目（user/groupName/idcs/groupTree/groupList/env）；
// mutations 在 Pinia 中无需单独定义，actions 内直接改 state 即可；getters 对应 Vuex getters。
export const useUserStore = defineStore('user', {
  state: () => ({
    // ---------- 参考项目对齐的全局共享数据 ----------
    // 当前登录用户（结构对齐参考项目：id=账号、displayName=姓名、groupName=部门、orgCode=团队）
    user: {
      load: false,
      id: '',
      displayName: '',
      groupName: '',
      groupCode: '',
      orgCode: '',
      groupDisplayName: '',
      email: '',
      roles: [] as string[],
      permissions: [] as string[],
      menus: [] as any[],
      buttons: [] as string[],
    },
    // 组/部门（参考项目：map 供快速查名，array 供下拉展示）
    groupName: { map: {} as Record<string, string>, array: [] as any[] },
    // IDC 列表
    idcs: [] as any[],
    // 部门树
    groupTree: [] as any[],
    // 部门列表
    groupList: {} as Record<string, any>,
    // 环境配置
    env: {} as any,

    // ---------- 本项目内部状态 ----------
    // 登录态标识：HttpOnly Cookie 方式下 localStorage 无 token，但浏览器持 Cookie 即视为已登录；
    // token 字段仅用于兼容 local 开发模式（localStorage），实际认证由后端从 Cookie 完成。
    token: tokenManager.hasCookie() ? '__cookie__' : (tokenManager.getToken() || ''),
    viewMode: 'personal' as 'personal' | 'admin',
    // 各角色的菜单树：后端按用户全部角色返回菜单并集（{ <role>: [菜单树] }）
    roleMenus: {} as Record<string, MenuTreeNode[]>,
    routesLoaded: false,
    // 组织架构数据（部门→团队→成员），登录后按需拉取缓存
    orgDepartments: [] as OrgDepartment[],
    deptTeams: [] as DeptTeams[],
  }),

  getters: {
    // 兼容现有组件：userInfo 从参考项目 user 结构映射为 UserInfo（id=账号/username，displayName=nickname，groupName=department，orgCode=team）
    userInfo(state): UserInfo | null {
      if (!state.user.id) return null
      return {
        id: state.user.id,
        username: state.user.id,
        nickname: state.user.displayName,
        role: (state.user.roles && state.user.roles[0]) || 'operator',
        department: state.user.groupName,
        team: state.user.orgCode,
        email: state.user.email,
        roles: state.user.roles,
        permissions: state.user.permissions,
      }
    },
    // 是否已登录：HttpOnly Cookie 存在（__cookie__ 标记）或 local 模式有 token
    isLoggedIn: (state) => !!state.token,
    role(): RoleType | undefined {
      return this.userInfo?.role
    },
    isOperator(): boolean {
      return this.role === 'operator'
    },
    isSpecialist(): boolean {
      return this.role === 'specialist'
    },
    isAdmin(): boolean {
      return this.role === 'admin'
    },
    // 全局治理视图权限：权限体系下发的 global_governance_view 权限点
    hasGlobalGovernanceView(): boolean {
      return (this.userInfo?.permissions || []).includes(PERM_GLOBAL_GOVERNANCE_VIEW)
    },
    // admin 角色或具备全局治理权限的用户，可切换 全局治理/个人工作台
    canToggleView(): boolean {
      return this.role === 'admin' || this.hasGlobalGovernanceView
    },
    showAdminView(): boolean {
      return this.viewMode === 'admin'
    },
    // 当前侧边栏菜单：
    // - 全局治理视图 → 只看 admin 角色菜单
    // - 个人视图 → 全部非 admin 角色菜单并集（admin 菜单归全局治理视图）；
    //   用户仅 admin 角色时回退 specialist 菜单（管理员的个人工作台，原设计约定）
    viewMenus(): MenuTreeNode[] {
      const merge = (roles: string[]): MenuTreeNode[] => {
        const merged: MenuTreeNode[] = []
        const seen = new Set<string>()
        for (const r of roles) {
          const tree = this.roleMenus[r] || []
          for (const node of tree) {
            const key = node.path || `${node.name}-${node.id}`
            if (seen.has(key)) continue
            seen.add(key)
            merged.push(node)
          }
        }
        return merged
      }
      if (this.viewMode === 'admin') {
        return this.roleMenus['admin'] || []
      }
      const roles = (this.userInfo?.roles?.length ? this.userInfo.roles : [this.userInfo?.role || 'operator'])
        .filter(r => r !== 'admin')
      let menus = merge(roles)
      if (menus.length === 0 && this.userInfo?.role === 'admin') {
        menus = merge(['specialist'])
      }
      return menus
    },
    // 按钮权限判断（v-auth 指令底层）：admin 默认全权限，其余看权限点集合
    hasPerm(): (code: string) => boolean {
      return (code: string) => {
        const u = this.userInfo
        if (!u) return false
        if (u.role === 'admin') return true
        return (u.permissions || []).includes(code)
      }
    },
  },

  actions: {
    // ---------- 参考项目 actions 对齐 ----------
    // 环境配置（对应参考项目 initEnv）
    async initEnv() {
      const data = await getEnv()
      if (data) this.env = data
    },
    // 组/部门列表（对应参考项目 initGroupName）
    async initGroupName() {
      const data = await getGroups()
      if (data) {
        this.groupName.array = data
        this.groupName.map = Object.fromEntries(data.map(g => [g.code, g.name]))
      }
    },
    // 部门树/列表（对应参考项目 initGroupTree）
    async initGroupTree() {
      const data = await getUserDepts()
      if (data) {
        this.groupTree = data.deptTree || []
        this.groupList = data.deptList || {}
      }
    },
    // IDC 列表（对应参考项目 initIdc）
    async initIdc() {
      const data = await getIdcs()
      if (data) this.idcs = data
    },
    // 当前用户完整权限信息（对应参考项目 initUserInfo）：/api/user/info → 写入 state.user
    // 返回 true=成功（后续调 fetchMenus 加载菜单），false=失败（调用方负责跳转 /login）
    async initUserInfo(): Promise<boolean> {
      if (!this.token) return false
      // HttpOnly Cookie 模式（__cookie__ 标记）不构造 Authorization 头，后端从 Cookie 认证
      const data = this.token === '__cookie__'
        ? await getUserInfo()
        : await getUserInfo({ Authorization: `Bearer ${this.token}` })
      if (!data) return false
      const a = data.account || {}
      this.user = {
        load: true,
        id: String(a.id ?? ''),
        displayName: a.displayName || '',
        groupName: a.groupName || '',
        groupCode: a.groupCode || '',
        orgCode: a.orgCode || '',
        groupDisplayName: a.groupDisplayName || '',
        email: a.email || '',
        roles: data.roles || [],
        permissions: data.permissions || [],
        menus: data.menus || [],
        buttons: data.buttons || [],
      }
      tokenManager.setUser(this.userInfo as UserInfo)
      // 权限被收回时强制回到个人工作台；仍具备切换能力的用户保留当前视图
      this.viewMode = canToggleByInfo(this.userInfo) ? this.viewMode : 'personal'
      return true
    },
    // 获取菜单并动态注册路由
    // /api/user/info 的 menusByRole 字段：{ <role>: [菜单树] }（承载原 /auth/menus 职责）
    // 前端一次拉取，全部角色的菜单路由均注册；视图切换（个人工作台/全局治理）只是按角色分组展示
    async fetchMenus() {
      if (!this.token) return
      try {
        const headers: Record<string, string> = {}
        if (this.token && this.token !== '__cookie__') {
          // 非 HttpOnly Cookie 模式才构造 Authorization 头；Cookie 模式后端从 Cookie 认证
          headers['Authorization'] = `Bearer ${this.token}`
        }
        if (this.userInfo) {
          // X-User-ID 统一使用 CAS 用户名（username），与后端用户唯一标识保持一致
          headers['X-User-ID'] = this.userInfo.username || ''
          headers['X-User-Role'] = this.userInfo.role || ''
          headers['X-User-Name'] = encodeURIComponent(this.userInfo.nickname || '')
          headers['X-User-Dept'] = encodeURIComponent(this.userInfo.department || '')
        }
        const data = await getUserInfo(headers)

        if (data?.menusByRole && typeof data.menusByRole === 'object') {
          const menusByRole: Record<string, MenuTreeNode[]> = data.menusByRole

          // 权限系统未下发任何菜单（menusByRole 为空对象、或所有角色菜单树均为空）
          // → 用前端本地角色映射兜底（menuMap.ts）
          // ============ 兜底分支：权限系统未下发菜单时生效（对接后可整段删除）============
          const allEmpty = Object.keys(menusByRole).length === 0 || Object.values(menusByRole).every(v => !Array.isArray(v) || v.length === 0)
          if (allEmpty) {
            console.warn('[fetchMenus] 权限系统返回空菜单，使用前端本地角色映射（menuMap.ts 兜底）')
            const fallbackRole = this.userInfo?.role || 'operator'
            const fallbackMenus = getRoleMenuMap(fallbackRole)
            this.roleMenus = { [fallbackRole]: fallbackMenus }
            registerDynamicRoutes(fallbackMenus, fallbackRole)
            // admin 角色的个人工作台使用 specialist 菜单（原设计约定），一并注册路由，保证视图切换可访问
            if (fallbackRole === 'admin') {
              const personalMenus = getRoleMenuMap('specialist')
              this.roleMenus['specialist'] = personalMenus
              registerDynamicRoutes(personalMenus, 'specialist')
            }
          } else {
            this.roleMenus = menusByRole
            // 全部角色的菜单路由均注册；角色菜单为空时由 registerDynamicRoutes 兜底注册默认首页路由
            for (const [menuRole, menuTree] of Object.entries(menusByRole)) {
              if (Array.isArray(menuTree) && menuTree.length > 0) {
                registerDynamicRoutes(menuTree, menuRole)
              } else {
                registerDynamicRoutes([], menuRole)
              }
            }
          }
          this.routesLoaded = true
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
      if (this.roleMenus && Object.keys(this.roleMenus).length === 0) {
        const fallbackRole = this.userInfo?.role || 'operator'
        const fallbackMenus = getRoleMenuMap(fallbackRole)
        this.roleMenus = { [fallbackRole]: fallbackMenus }
        registerDynamicRoutes(fallbackMenus, fallbackRole)
      }
      this.routesLoaded = true
    },

    // ---------- 登录 / 登出 ----------
    // CAS 登录：回调页拿到 token 后调用。
    // 用 token 换取用户档案（/api/user/info）→ 写入 state/localStorage → 加载菜单路由。
    // 成功返回 true，失败返回 false（调用方负责跳转 /login）。
    async loginWithToken(jwtToken: string): Promise<boolean> {
      if (!jwtToken) return false
      this.token = jwtToken
      tokenManager.setToken(jwtToken)
      const ok = await this.initUserInfo()
      if (!ok) {
        tokenManager.removeToken()
        this.token = ''
        return false
      }
      // admin 默认进入全局治理视图；其他用户（含具备全局治理权限的专员）默认个人工作台
      this.viewMode = this.userInfo?.role === 'admin' ? 'admin' : 'personal'
      this.routesLoaded = false
      await this.fetchMenus()
      return true
    },

    // 拉取组织架构（部门→团队→成员）并缓存，供人员选择/部门团队展示使用
    async fetchOrg() {
      if (!this.token) return
      try {
        const [depts, teams] = await Promise.all([orgApi.departments(), orgApi.departmentsTeams()])
        this.orgDepartments = depts || []
        this.deptTeams = teams || []
      } catch (err) {
        console.error('[fetchOrg] 获取组织架构失败:', err)
      }
    },

    async logout() {
      this.user = {
        load: false,
        id: '',
        displayName: '',
        groupName: '',
        groupCode: '',
        orgCode: '',
        groupDisplayName: '',
        email: '',
        roles: [],
        permissions: [],
        menus: [],
        buttons: [],
      }
      this.token = ''
      this.roleMenus = {}
      this.routesLoaded = false
      resetRouter()
      tokenManager.removeToken()
      tokenManager.removeUser()
      this.viewMode = 'personal'
      localStorage.removeItem(VIEW_MODE_STORAGE_KEY)
      // 调用后端 /api/v1/auth/logout：CAS 模式下后端 302 到 CAS 登出页清除 SSO 会话，再跳回前端登录页
      casLogout()
    },

    toggleViewMode() {
      this.viewMode = this.viewMode === 'personal' ? 'admin' : 'personal'
      // 记住上次的视图模式，刷新后恢复（权限仍以 /auth/me 实时返回为准）
      localStorage.setItem(VIEW_MODE_STORAGE_KEY, this.viewMode)
    },

    // 首次进入时从 localStorage 恢复登录状态（token 与用户信息）
    initFromStorage() {
      const savedToken = tokenManager.getToken()
      const savedUser = tokenManager.getUser()
      if (savedToken && savedUser) {
        this.token = savedToken
        // 恢复用户档案（savedUser 为 UserInfo 结构，映射回参考项目 user 结构）
        this.user = {
          load: true,
          id: savedUser.username || savedUser.id || '',
          displayName: savedUser.nickname || '',
          groupName: savedUser.department || '',
          groupCode: savedUser.department || '',
          orgCode: savedUser.team || '',
          groupDisplayName: savedUser.department || '',
          email: savedUser.email || '',
          roles: savedUser.roles || [],
          permissions: savedUser.permissions || [],
          menus: [],
          buttons: [],
        }
        // 恢复上次视图模式（仅 admin 或具备全局治理权限的用户有效）
        const savedMode = localStorage.getItem(VIEW_MODE_STORAGE_KEY)
        this.viewMode = canToggleByInfo(savedUser) && savedMode === 'admin' ? 'admin' : 'personal'
        this.routesLoaded = false
      }
    },
  },
})
