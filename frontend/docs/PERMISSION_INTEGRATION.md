# 前端权限对接说明（角色 / 菜单 / 按钮 / 组织架构）

> 用途：把当前 osm-mngt 前端权限实现机制梳理清楚，便于在另一套环境对接统一权限系统时照此实现。
> 核心结论：
> - **动态路由**：后端按「角色→菜单树」下发（`/auth/menus`），前端把菜单树的 `component` 字段映射到 vue 组件并 `addRoute`；侧边菜单渲染同一个菜单树。
> - **按钮权限**：框架代码**已落地**（`src/directives/auth.ts` + `main.ts` 注册 + `store.hasPerm`），在按钮上写 `v-auth="'权限编码'"` 即可；权限点由后端 `/auth/me` 的 `permissions[]` 下发。
> - **组织架构**：后端 `/org/departments`、`/org/departments/teams`、`/users?team=` 已实现，前端 `src/api/org.ts` + `src/types/org.ts` 承接。
> - **状态集中**：已集中在 `stores/user.ts`（token / userInfo / roles / permissions / roleMenus / viewMode / orgDepartments）。

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
