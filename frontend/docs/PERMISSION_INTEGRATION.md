# 前端权限对接说明（角色 / 菜单 / 按钮 / 组织架构）

> 用途：把当前 osm-mngt 前端权限实现机制梳理清楚，便于在另一套环境对接统一权限系统时照此实现；
> 并对照参考项目（osm-front，Vue2 + vue-cli）给出结构映射，便于页面/文件对位迁移。
> 核心结论（当前架构）：
> - **登录链路**：CAS 登录只签发基本身份 JWT → 前端登录回调请求权限接口 `/auth/me`（user/roles/permissions）、`/auth/menus`（roles + 各角色菜单树）；菜单完全由权限系统下发，后端只转发。
> - **静态路由 + path 过滤**：前端静态定义全部业务路由（`router/index.ts` 的 `staticRoutes`），登录后按「可见菜单 path 集合」过滤注册（公司参考项目同逻辑）；侧边菜单渲染同一个菜单树。
> - **多角色并集**：`/auth/menus` 返回 `{ roles[], menus: { <role>: [菜单树] } }`，侧边栏 = 各角色菜单按 path 去重合并；admin 菜单归「全局治理视图」，admin 用户个人工作台回退 specialist 菜单。
> - **未对接兜底**：权限系统未对接（baseURL 空）时后端打印日志返回空菜单，前端按角色用本地映射 `src/router/menuMap.ts` 渲染（可整段删除）；`fetchMenus()` 失败不清 token（认证≠权限展示）。
> - **按钮权限**：已落地（`src/directives/auth.ts` + `main.ts` 注册 + `store.hasPerm`），按钮写 `v-auth="'权限编码'"`；权限点由 `/auth/me` 的 `permissions[]` 下发。方法级判断用 `src/utils/auth.ts` 的 `hasPermission(code)`。
> - **组织架构**：后端 `/org/departments`、`/org/departments/teams`、`/users?team=` 优先权限系统、未对接时兜底 users 表；前端 `src/api/org.ts` + `src/types/org.ts` 承接。
> - **参考接口兼容层**：后端 `/api/env`、`/api/groups`、`/api/user/depts`、`/api/idcs`、`/api/user/info`、`/api/user/refresh`（`{ok,data}` 结构，handlers/reference.go），参考前端页面可直接对接。

---

## 附录：当前项目 ↔ 参考项目（osm-front）结构对照

> 目的：参考前端页面/工具文件可对位迁移到本项目，或本项目文件替换到参考项目，便于开发对照。

| 参考项目（Vue2 + vue-cli） | 当前项目（Vue3 + Vite） | 职责 |
|---|---|---|
| `src/main.js`（路由守卫 beforeEach + Vue.prototype.hasPermission） | `src/main.ts`（注册 v-auth 指令 + `utils/auth.ts` installPermission 注入） | 入口 + 全局权限 |
| `src/App.vue`（全局初始化 dispatch initEnv/initGroupTree/initGroupName） | `src/App.vue` | 根组件 |
| `src/router`（modules/每菜单一个 js，index 聚合） | `src/router/modules/*`（mytask/myops/asset/risk...每菜单一个 ts）+ `index.ts` 聚合，按可见 path 过滤注册 | 路由配置 |
| `src/store`（Vuex：用户信息/全局数据） | `src/stores/user.ts`（Pinia：token/userInfo/roles/permissions/roleMenus） | 状态管理 |
| `src/utils/request.js`（Axios 封装 + 拦截器） | `src/api/index.ts`（axios 实例 + 拦截器） | HTTP 请求封装 |
| `src/utils/auth.js`（hasPriv/canRouter/v-auth） | `src/utils/auth.ts`（hasPermission/canRouter）+ `src/directives/auth.ts`（v-auth） | 权限判断 |
| `src/components`（公共组件：顶部栏/左侧菜单） | `src/components/layout/MainLayout.vue` + `src/components/business/*` | 公共组件 |
| `src/pages`（各业务页面） | `src/views/**`（operator/specialist/admin 各工作台页面） | 业务页面 |

**对接要点**：
- 参考前端页面组件（Options API 风格）迁移到 `src/views/**` 后，路由需在 `staticRoutes` 登记 path + 组件，菜单 path 由权限系统下发（或 menuMap 本地映射）
- 权限判断代码统一走 `src/utils/auth.ts`（hasPermission/canRouter），与参考 `utils/auth.js` 职责一一对应
- 请求封装统一走 `src/api/index.ts`（baseURL `/api/v1`）；参考项目请求的 `/api/user/info` 等 6 接口由后端 reference 兼容层提供

---

## 1. 权限四要素总览

| 要素 | 现状实现 | 关键文件 |
|---|---|---|
| 用户信息（角色/部门/团队） | `GET /api/v1/auth/me` → `{ user, roles[], permissions[] }` | [auth.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/api/auth.ts) |
| 动态路由（角色→菜单→path） | `GET /api/v1/auth/menus` → `{ roles[], menus: { <role>: [菜单树] } }` | [router/index.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/router/index.ts) |
| 侧边菜单渲染 | 布局组件读 `userStore.viewMenus` 递归渲染 `menu.name/path` | [MainLayout.vue](file:///Users/lichaoqun/project/osm-mngt/frontend/src/components/layout/MainLayout.vue) |
| 按钮权限（v-auth） | **未实现**（仅 permissions 用于视图切换） | [stores/user.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/stores/user.ts) |

---

## 2. 动态路由原理（角色 → 菜单 → path 渲染可见菜单）

### 2.1 数据流

```
CAS 登录 → /auth/me（user.role / roles[] / permissions[]）
        → /auth/menus（menus: { <role>: [菜单树] }）

前端：
  store.fetchUserInfo()  → 写入 userInfo（role/department/team/roles/permissions）
  store.fetchMenus()     → 写入 roleMenus（全部角色菜单并集）
                        → 对每个角色递归 registerDynamicRoutes(menuTree, role)
  router.addRoute(parentName, { path, component, meta:{title} })
  MainLayout 读 viewMenus（当前视图角色菜单）渲染 el-menu
```

### 2.2 菜单树节点结构（后端下发）

```ts
interface MenuTreeNode {
  id: number
  name: string        // 菜单显示名
  path: string        // 完整路由路径，如 /admin/risk/defect-ledger
  component: string   // 组件相对路径，如 views/admin/risk/DefectLedger（前端映射表 key）
  icon: string
  sort: number
  hidden: boolean
  group: string
  children: MenuTreeNode[]
}
```

### 2.3 组件映射表（`component` 字符串 → 实际 vue 组件）

[router/index.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/router/index.ts#L67-L98) 的 `componentMaps` 是硬编码映射：

```ts
const componentMaps: Record<string, () => Promise<any>> = {
  'views/operator/MyTask': () => import('@/views/operator/MyTask.vue'),
  'views/admin/risk/DefectLedger': () => import('@/views/admin/risk/DefectLedger.vue'),
  // ... 每个接入页面对应一条
}
```

> **对接要点**：后端菜单的 `component` 值必须与前端 `componentMaps` 的 key 完全一致，否则 `registerDynamicRoutes` 会 `console.warn('[registerDynamicRoutes] 组件未找到')` 并跳过该路由（表现为菜单点击无反应）。

### 2.4 动态注册（`registerDynamicRoutes`）

```ts
function registerMenu(menu: MenuTreeNode) {
  if (menu.component && componentMaps[menu.component]) {
    const routePath = menu.path.replace(`${parentPath}/`, '')  // 去掉 /admin 前缀
    router.addRoute(parentName, {
      path: routePath,
      name: `${role}_${menu.path.replace(/\//g, '_')}`,
      component: componentMaps[menu.component],
      meta: { title: menu.name, role },
    })
  }
  menu.children?.forEach(registerMenu)  // 递归
}
```

### 2.5 侧边菜单渲染（`MainLayout.vue`）

- 布局读 `const displayMenus = computed(() => userStore.viewMenus)`（[MainLayout.vue](file:///Users/lichaoqun/project/osm-mngt/frontend/src/components/layout/MainLayout.vue#L169-L170)）
- `viewMenus` 由 `stores/user.ts` 的 `viewMode`（personal / admin）决定取哪个角色的菜单树
- 视图切换：`role==='admin'` 或 `permissions` 含 `global_governance_view` 才允许切换（[stores/user.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/stores/user.ts#L34-L50)）

```ts
const viewMenus = computed(() => {
  if (viewMode.value === 'admin') return roleMenus.value['admin'] || []
  if (role.value === 'admin') return roleMenus.value['specialist'] || []  // admin 个人工作台=specialist 菜单
  return roleMenus.value[role.value || 'operator'] || []
})
```

### 2.6 路由守卫（`router.beforeEach`）

- 未登录 → `/login`
- 已登录但 `routesLoaded=false` → `fetchUserInfo()` + `fetchMenus()` → 动态路由注册后 `next({ path: to.fullPath, replace: true })` 重新导航
- 菜单获取失败 → 清空 token 强制重新登录（[stores/user.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/stores/user.ts#L122-L129)）

---

## 3. 按钮权限（v-auth）——当前未实现，给出方案

### 3.1 现状

- 后端 `/auth/me` 已返回 `permissions: string[]`（权限点编码，目前仅 `global_governance_view`）
- 前端 `UserInfo.permissions` 已定义（[types/index.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/types/index.ts#L11-L22)），store 已存（`userInfo.permissions`）
- **没有自定义指令**：`main.ts` 未注册任何 `v-auth` 指令，模板里也没有 `v-auth` 用法

### 3.2 原系统 v-auth 机制的原理

1. 后端按角色在权限系统配置「按钮权限值」（如 `defect:review`、`vuln:create`），登录后随用户权限档案下发
2. 前端按钮写 `<el-button v-auth="'defect:review'">评审</el-button>`
3. 前端全局注册自定义指令 `v-auth`：
   - `mounted/updated` 时读取指令值（权限编码）
   - 用 store 里的权限集合判断当前用户是否有该权限
   - 没有 → `el.parentNode.removeChild(el)`（或 `el.style.display='none'`）

### 3.3 本工程落地示例（可直接套用）

**指令定义** `src/directives/auth.ts`：

```ts
import type { Directive, DirectiveBinding } from 'vue'
import { useUserStore } from '@/stores/user'

export const auth: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string>) {
    const userStore = useUserStore()
    if (!hasAuth(userStore, binding.value)) el.parentNode?.removeChild(el)
  },
  updated(el: HTMLElement, binding: DirectiveBinding<string>) {
    const userStore = useUserStore()
    if (!hasAuth(userStore, binding.value)) el.parentNode?.removeChild(el)
  },
}

function hasAuth(store: ReturnType<typeof useUserStore>, code: string): boolean {
  const perms = store.userInfo?.permissions || []
  // admin 角色默认全权限；否则看权限点集合
  return store.userInfo?.role === 'admin' || perms.includes(code)
}
```

**注册** `src/main.ts`（追加）：

```ts
import { auth } from '@/directives/auth'
app.directive('auth', auth)
```

**store 增加权限判断**（[stores/user.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/stores/user.ts)）：

```ts
const hasPerm = (code: string) =>
  userInfo.value?.role === 'admin' || (userInfo.value?.permissions || []).includes(code)
```

**后端**：`/auth/me` 的 `permissions` 从权限系统按角色聚合下发全部按钮权限点。

> 注：`v-auth` 用 `removeChild` 是原系统常见做法（简单、元素不占位）；若需要"禁用"而非隐藏，可用 `el.setAttribute('disabled','disabled')` + 事件拦截。

---

## 4. 组织架构接口（部门 / 团队 / 成员）

### 4.1 现状

- `GET /api/v1/users?dept_name=技术部` → 返回该部门人员 `[{id, username, nickname, department, team, email}]`（[auth.go](file:///Users/lichaoqun/project/osm-mngt/backend/internal/handlers/auth.go#L268-L293)）
- 数据源：`users` 表（username/nickname/department/team）

### 4.2 缺口（对照你提的两类接口）

| 诉求 | 现状 | 建议新增接口 |
|---|---|---|
| 部门、团队、团队人员 | 仅按部门取人员 | `GET /org/departments` → `[{ dept, teams: [{ team, members: [...] }] }]` |
| 部门及所属团队 | 无 | `GET /org/departments/teams` → `[{ dept, teams: ['前端组','后端组'] }]` |
| 按团队取人员 | 无 | `GET /users?team=前端组`（复用现有 /users 加 team 参数） |

**后端建议实现**（Gin + GORM，`users` 表 group by）：

```go
// GET /api/v1/org/departments 部门→团队→成员
type TeamMember struct { Username, Nickname string }
type DeptTeam struct {
    Dept   string       `json:"dept"`
    Teams  []TeamGroup  `json:"teams"`
}
type TeamGroup struct {
    Team    string        `json:"team"`
    Members []TeamMember  `json:"members"`
}

rows := []struct{ Department, Team string }{}
h.db.Model(&models.User{}).Where("department <> ''").Group("department, team").Find(&rows)
// 按 department 分组，team 下再查 members
```

**前端**：新增 `src/api/org.ts` 模块（独立文件，与 auth.ts 分模块，符合"接口请求独立文件"诉求）：

```ts
// src/api/org.ts —— 组织架构接口（走 api 实例，直接请求后端）
import api from './index' // 或复用 axios 实例

export const orgApi = {
  // 部门及其所属团队
  departmentsTeams: () => api.get('/org/departments/teams'),
  // 部门→团队→成员
  departments: () => api.get('/org/departments'),
}
```

---

## 5. Store 集中（权限状态）

[stores/user.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/stores/user.ts) 已集中所有权限相关状态，对接时扩展即可：

```ts
export const useUserStore = defineStore('user', () => {
  const token = ref<string>(tokenManager.getToken() || '')
  const userInfo = ref<UserInfo | null>(null)   // role/department/team/roles/permissions
  const roleMenus = ref<Record<string, MenuTreeNode[]>>({}) // { <role>: [菜单树] }
  const viewMode = ref<'personal' | 'admin'>('personal')    // 视图切换
  const routesLoaded = ref(false)

  // 派生
  const role = computed(() => userInfo.value?.role)
  const canToggleView = computed(() => role === 'admin' || hasGlobalGovernanceView)
  const viewMenus = computed<MenuTreeNode[]>(() => { ... })   // 当前视图菜单
  const hasPerm = (code: string) => ...                        // 按钮权限判断（新增）

  // 动作
  fetchUserInfo()  // /auth/me
  fetchMenus()     // /auth/menus → registerDynamicRoutes
  setCasUser() / logout() / toggleViewMode() / initFromStorage()
})
```

- 登录态持久化：`localStorage` 的 `osm_token` / `osm_user`（tokenManager）
- 刷新恢复：路由守卫里 `initFromStorage()` → 再 `fetchUserInfo()` + `fetchMenus()`

---

## 6. 对接清单（后端要提供 / 前端要接）

**后端新增/确认：**
1. `GET /auth/me` 返回 `{ user, roles[], permissions[] }`（permissions 按角色聚合下发按钮权限点）
2. `GET /auth/menus` 返回 `{ roles[], menus: { <role>: [MenuTreeNode] } }`（component 值 = 前端 componentMaps key）
3. `GET /org/departments/teams` 部门→所属团队（新增）
4. `GET /org/departments` 部门→团队→成员（新增）
5. `GET /users?dept_name=&team=` 人员查询（补 team 过滤）

**前端要接：**
1. `src/directives/auth.ts` 新增 v-auth 指令 + `main.ts` 注册（按钮权限）
2. `src/api/org.ts` 新增组织接口模块
3. `stores/user.ts` 增加 `hasPerm` 派生（已有 userInfo/roleMenus/permissions 基础）
4. `router/index.ts` 的 `componentMaps` 与后端菜单 `component` 保持一致（新增页面时两边同步）
5. 按钮模板加 `v-auth="'权限编码'"`

---

## 7. 关键文件索引（现状代码）

| 文件 | 说明 |
|---|---|
| [frontend/src/router/index.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/router/index.ts) | 动态路由注册、componentMaps、路由守卫 |
| [frontend/src/stores/user.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/stores/user.ts) | 权限状态集中（token/userInfo/roleMenus/viewMode） |
| [frontend/src/api/auth.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/api/auth.ts) | 认证与用户信息接口封装 |
| [frontend/src/types/index.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/types/index.ts#L1-L22) | UserInfo / RoleType 定义 |
| [frontend/src/components/layout/MainLayout.vue](file:///Users/lichaoqun/project/osm-mngt/frontend/src/components/layout/MainLayout.vue) | 侧边菜单渲染（viewMenus） |
| [frontend/src/main.ts](file:///Users/lichaoqun/project/osm-mngt/frontend/src/main.ts) | 全局注册入口（指令在此注册） |
| [backend/internal/handlers/auth.go](file:///Users/lichaoqun/project/osm-mngt/backend/internal/handlers/auth.go#L268-L396) | /users、/auth/menus、/auth/me 实现 |
