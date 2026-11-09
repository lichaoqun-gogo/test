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

// ==================== 权限系统角色编码 → 项目三类角色 映射 ====================
// 后端 /api/user/info 返回的 roles 为统一权限系统编码（贴近内网 COMMON/SPECIALIST/ADMIN），
// 前端据此划分为三类角色形态：operator 普通用户 / specialist 部门治理专员 / admin 管理部门专员。
// 用户约定：权限系统把管理部门专员编码叫 admin 也行，关键是能划分清楚三类角色。
// 对接内网后若出现新编码，在此表补充映射即可（映射不到的按 operator 兜底）。
const ROLE_CODE_MAP: Record<string, RoleType> = {
  COMMON: 'operator',
  SPECIALIST: 'specialist',
  ADMIN: 'admin',
  // 兼容小写/全拼（CAS/JWT 兜底、开发环境直接下发小写编码的场景）
  common: 'operator',
  specialist: 'specialist',
  admin: 'admin',
  operator: 'operator',
}
function mapRoleCode(code?: string): RoleType {
  return ROLE_CODE_MAP[code || ''] || 'operator'
}

// 判断用户信息是否具备视图切换能力（admin 角色，或权限体系授予 global_governance_view）
function canToggleByInfo(info?: UserInfo | null): boolean {
  if (!info) return false
  return info.role === 'admin' || (info.permissions || []).includes(PERM_GLOBAL_GOVERNANCE_VIEW)
}

// ==================== 菜单解析辅助函数（对齐后端 /api/user/info 的 menus 数组结构） ====================
// 内网菜单节点结构：{id, name, appCode, appModule, parentId, orderNum, path, icon, visible, children}
// → 映射为前端 MenuTreeNode（component 由前端 staticRoutes 按 path 匹配，菜单节点无需携带）。

// normalizeMenus 后端 menus 数组 → MenuTreeNode 数组
function normalizeMenus(menus: any[]): MenuTreeNode[] {
  return (menus || []).map(m => ({
    id: m.id ?? 0,
    name: m.name || '',
    path: m.path || '',
    component: m.component || '',
    icon: m.icon || '',
    sort: m.orderNum ?? 0,
    hidden: m.visible === false,
    group: m.appModule || '',
    children: normalizeMenus(m.children),
  }))
}

// 菜单 path 前缀 → 角色（本项目静态路由按 operator/specialist/admin 前缀组织，见 router/modules）
function inferMenuRole(path: string): string {
  if (path.startsWith('/admin')) return 'admin'
  if (path.startsWith('/specialist')) return 'specialist'
  return 'operator'
}

// 节点及子节点归属的角色集合（目录节点 path 可能为空，以子节点 path 为准）
function collectNodeRoles(node: MenuTreeNode): Set<string> {
  const roles = new Set<string>()
  if (node.path) roles.add(inferMenuRole(node.path))
  ;(node.children || []).forEach(ch => {
    collectNodeRoles(ch).forEach(r => roles.add(r))
  })
  return roles
}

// groupMenusByRole 按 path 前缀把 menus 数组分为 {角色: [菜单树]}，用于按角色注册动态路由
function groupMenusByRole(menus: MenuTreeNode[]): Record<string, MenuTreeNode[]> {
  const grouped: Record<string, MenuTreeNode[]> = {}
  for (const node of menus) {
    const role = [...collectNodeRoles(node)][0] || 'operator'
    if (!grouped[role]) grouped[role] = []
    grouped[role].push(node)
  }
  return grouped
}

// menuScope 菜单归属：全局治理（存在 /admin 前缀路径）或个人工作台，用于视图模式过滤
function menuScope(node: MenuTreeNode): 'global' | 'personal' {
  const hasAdmin = (n: MenuTreeNode): boolean => {
    if (n.path && n.path.startsWith('/admin')) return true
    return (n.children || []).some(hasAdmin)
  }
  return hasAdmin(node) ? 'global' : 'personal'
}

// extractPermissions 从权限系统 buttons 提取按钮权限编码（resource）作为权限点集合，供 v-auth/hasPerm 判断
function extractPermissions(data: any): string[] {
  const perms = [...(data.permissions || [])]
  const buttons = data.buttons || []
  ;(Array.isArray(buttons) ? buttons : []).forEach((b: any) => {
    if (b && b.visible !== false && b.resource) perms.push(b.resource)
  })
  return perms
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
    // 后端 /api/user/info 返回的完整菜单数组（权限系统结构，对齐内网）：登录后按视图过滤渲染侧边栏、按 path 注册路由
    menus: [] as MenuTreeNode[],
    // /api/user/info 完整响应缓存：initUserInfo 一次请求后供 fetchMenus 复用，
    // 避免登录时对同一接口发两次请求（原实现 initUserInfo + fetchMenus 各请求一次）。
    userInfoResp: null as any,
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
        role: mapRoleCode((state.user.roles && state.user.roles[0]) as string),
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
    // - 权限系统菜单（menus 数组，后端 /api/user/info 下发）：按视图模式过滤个人工作台/全局治理
    //   （个人=path 以 /operator|/specialist 开头，全局=path 以 /admin 开头，见 menuScope）
    // - 兜底：menuMap 本地角色映射（权限系统未下发菜单时，见 fetchMenus）
    viewMenus(): MenuTreeNode[] {
      if (this.menus.length > 0) {
        const list = this.viewMode === 'admin'
          ? this.menus.filter(n => menuScope(n) === 'global')
          : this.menus.filter(n => menuScope(n) !== 'global')
        // 过滤结果为空（如兜底菜单未含对应视图）时回退角色分组逻辑，保证侧边栏有内容可渲染
        if (list.length > 0) return list
      }
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
        .map(r => mapRoleCode(r))
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
    // account 对齐权限系统结构：accountCode/accountName/depCode/depName/orgCode/orgName/orgDisplayName
    // buttons 对齐权限系统结构：resource 编码提取到 permissions，供 v-auth/hasPerm 判断按钮显隐
    async initUserInfo(): Promise<boolean> {
      if (!this.token) return false
      // HttpOnly Cookie 模式（__cookie__ 标记）不构造 Authorization 头，后端从 Cookie 认证
      const data = this.token === '__cookie__'
        ? await getUserInfo()
        : await getUserInfo({ Authorization: `Bearer ${this.token}` })
      if (!data) return false
      // 缓存完整响应，fetchMenus 直接复用，避免对同一接口二次请求
      this.userInfoResp = data
      // HttpOnly Cookie 模式下认证走 Cookie，清理 localStorage 中历史遗留的旧 token（避免 axios 误带旧 Authorization 头）
      if (this.token === '__cookie__') {
        tokenManager.removeToken()
      }
      const a = data.account || {}
      this.user = {
        load: true,
        id: String(a.accountCode ?? a.id ?? ''),
        displayName: a.accountName ?? a.displayName ?? '',
        groupName: a.depName ?? a.groupName ?? '',
        groupCode: a.depCode ?? a.groupCode ?? '',
        orgCode: a.orgCode ?? '',
        groupDisplayName: a.orgDisplayName ?? a.groupDisplayName ?? '',
        email: a.email ?? '',
        roles: data.roles || [],
        permissions: extractPermissions(data),
        menus: [],
        buttons: data.buttons || [],
      }
      tokenManager.setUser(this.userInfo as UserInfo)
      // 权限被收回时强制回到个人工作台；仍具备切换能力的用户保留当前视图
      this.viewMode = canToggleByInfo(this.userInfo) ? this.viewMode : 'personal'
      return true
    },
    // 获取菜单并动态注册路由
    // 数据源：后端 /api/user/info 的 menus 数组（权限系统结构，对齐内网），前端按 path 前缀分组注册；
    // 权限系统未下发菜单时走前端本地角色映射兜底（menuMap.ts，对接后可删除）。
    async fetchMenus() {
      if (!this.token) return
      try {
        // 优先复用 initUserInfo 已缓存的完整响应（同一接口只请求一次）；
        // 单独调用本方法（如 CasCallback 兜底、无 token 参数场景）时才自行请求。
        let data = this.userInfoResp
        if (!data) {
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
          data = await getUserInfo(headers)
          this.userInfoResp = data
        }
        if (!data) return

        // ===== 主分支：权限系统下发菜单（menus 数组，对齐内网结构）=====
        if (Array.isArray(data.menus)) {
          const normalized = normalizeMenus(data.menus).filter(n => !n.hidden)
          if (normalized.length > 0) {
            this.menus = normalized
            // 按 path 前缀分组角色，分别注册动态路由（admin 菜单注册到 AdminLayout 等）
            const grouped = groupMenusByRole(normalized)
            for (const [menuRole, tree] of Object.entries(grouped)) {
              registerDynamicRoutes(tree, menuRole)
            }
            this.routesLoaded = true
            return
          }
        }

        // ============ 兜底分支：权限系统未下发菜单时生效（对接后可整段删除）============
        // 菜单可见性属于"权限展示"而非"认证判定"：菜单拉取失败不登出用户；
        // 按用户角色用前端本地菜单映射（menuMap.ts）提供默认菜单（公司模型：有 app 即有默认角色）。
        console.warn('[fetchMenus] 权限系统未返回菜单，使用前端本地角色映射（menuMap.ts 兜底）')
        const fallbackRole = this.userInfo?.role || 'operator'
        const fallbackMenus = getRoleMenuMap(fallbackRole)
        this.menus = [...fallbackMenus]
        this.roleMenus = { [fallbackRole]: fallbackMenus }
        registerDynamicRoutes(fallbackMenus, fallbackRole)
        // admin 角色的个人工作台使用 specialist 菜单（原设计约定），并入 menus（个人视图可渲染）并注册路由
        if (fallbackRole === 'admin') {
          const personalMenus = getRoleMenuMap('specialist')
          this.roleMenus['specialist'] = personalMenus
          this.menus = [...fallbackMenus, ...personalMenus]
          registerDynamicRoutes(personalMenus, 'specialist')
        }
        this.routesLoaded = true
      } catch (err) {
        console.error('获取菜单失败:', err)
      }
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
      // 所有角色默认进入个人工作台（管理部门专员同样默认个人工作台，可手动切换全局治理）
      this.viewMode = 'personal'
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
      this.menus = []
      this.userInfoResp = null
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
        // 用户信息缓存置空：恢复登录态后由路由守卫走 initUserInfo 重新拉取（一次请求，避免用旧缓存）
        this.userInfoResp = null
        this.menus = []
      }
    },
  },
})
