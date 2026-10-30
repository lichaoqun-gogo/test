// ========== 组织架构类型定义（部门 / 团队 / 成员） ==========
// 对接说明：这些类型对应后端 /org/* 与 /users 接口的返回结构。
// 若外部权限系统的真实返回字段名不同（如 deptName 而非 department），
// 只需改这里的字段名，前端组件引用处同步即可。

// 团队成员
export interface OrgMember {
  username: string   // 账号（CAS 用户名）
  nickname: string   // 姓名
  email?: string
}

// 部门下某个团队（含成员）
export interface OrgTeam {
  team: string           // 团队名，如 前端组
  members: OrgMember[]   // 团队成员
}

// 部门 → 团队 → 成员（GET /org/departments）
export interface OrgDepartment {
  department: string        // 部门名，如 技术部
  teams: OrgTeam[]
}

// 部门 → 所属团队（仅团队名列表，不含成员，GET /org/departments/teams）
export interface DeptTeams {
  department: string
  teams: string[]
}
