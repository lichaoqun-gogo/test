// ==================== 参考项目接口封装（对应后端 handlers/reference.go） ====================
// 参考项目（osm-front）请求的 6 个接口，路径不带 /v1，统一 {ok, data} 结构：
//   1. GET /api/env           环境配置（App.vue 挂载时初始化调用）
//   2. GET /api/groups        组/部门列表（App.vue 挂载时初始化调用）
//   3. GET /api/user/depts    部门树/列表（App.vue 挂载时初始化调用）
//   4. GET /api/idcs          IDC 列表（App.vue 挂载时初始化调用）
//   5. GET /api/user/info     当前用户完整权限信息（路由守卫调用，对应参考项目 main.js 守卫访问 /user/info）
//   6. GET /api/user/refresh  刷新权限缓存
// 调用时机对齐参考项目：
//   - App.vue created：initEnv / initGroupTree / initGroupName（env + groups + depts + idcs）
//   - main.js 路由守卫 beforeEach：GET /user/info（本项目等价走 /auth/me，另调本接口对位存储）
// 注意：这些接口失败不阻断流程（后端未对接权限系统时兜底返回/失败），登录照常。

import axios from 'axios'

// 独立 axios 实例：参考接口路径为 /api/*（业务接口是 /api/v1/*），响应结构为 {ok, data}
const refApi = axios.create({ baseURL: '/api', timeout: 10000 })

// 响应拦截器：参考接口统一 {ok: true, data: ...}；失败返回 null 并告警（不抛错，避免阻断初始化/登录）
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

// 环境配置（App 启动时初始化）
export const getEnv = () => refApi.get('/env') as Promise<{ profiles: { name: string; active: boolean }[]; tips: string[] } | null>
// 组/部门列表
export const getGroups = () => refApi.get('/groups') as Promise<{ code: string; name: string }[] | null>
// 部门树/列表
export const getUserDepts = () => refApi.get('/user/depts') as Promise<{ deptTree: any[]; deptList: any[] } | null>
// IDC 列表
export const getIdcs = () => refApi.get('/idcs') as Promise<any[] | null>
// 当前用户完整权限信息（对位参考项目 /user/info）
export const getUserInfo = () => refApi.get('/user/info') as Promise<any | null>
// 刷新权限缓存
export const refreshUserInfo = () => refApi.get('/user/refresh') as Promise<any | null>

export default refApi
