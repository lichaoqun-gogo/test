// ========== 组织架构接口模块（部门 / 团队 / 成员） ==========
// 对接说明：
//   1. 默认走全局 axios 实例（src/api/index.ts，baseURL=/api/v1），请求本平台后端。
//   2. 若组织架构数据来自外部权限系统网关，只需把下方 BASE 换成网关地址
//      （如 https://iam.example.com/openapi），其余调用处不用改。
//   3. 返回结构由 src/types/org.ts 承接；真实返回字段不同时改类型定义即可。
import api from './index'
import type { OrgDepartment, DeptTeams, OrgMember } from '@/types/org'

// 组织架构接口统一前缀（对接外部系统时整体替换）
const BASE = '/org'

export const orgApi = {
  // 部门及其所属团队（不含成员）→ [{ department, teams: ['前端组', ...] }]
  // 泛型第二参 R=T：api 实例的响应拦截器已返回 data（原始响应体），无需再取 .data
  departmentsTeams: () => api.get<DeptTeams[], DeptTeams[]>(`${BASE}/departments/teams`),

  // 部门 → 团队 → 成员 → [{ department, teams: [{ team, members: [...] }] }]
  departments: () => api.get<OrgDepartment[], OrgDepartment[]>(`${BASE}/departments`),

  // 按团队取成员（复用 /users，新增 team 过滤参数）→ OrgMember[]
  teamMembers: (team: string) => api.get<OrgMember[], OrgMember[]>('/users', { params: { team } }),
}

// 按部门取成员（复用 /users 的 dept_name 过滤，人员选择下拉等场景）
export function usersByDept(dept: string) {
  return api.get<OrgMember[], OrgMember[]>('/users', { params: { dept_name: dept } })
}
