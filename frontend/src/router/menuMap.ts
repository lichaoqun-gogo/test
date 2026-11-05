// ==================== 本地角色-菜单映射（兜底用） ====================
// 背景：登录后菜单来源优先级为「统一权限系统」返回的菜单树（后端 /auth/menus 透传）；
// 权限系统未对接/返回为空时，前端按用户角色用本文件提供默认菜单（公司模型：有 app 即有默认角色）。
// ============ 兜底文件：对接真实权限系统后可整段删除 ============
// 注意：各菜单的 path 必须与 router/index.ts 中 staticRoutes 的 path 完全一致
// （前端静态路由 + 按可见 path 集合过滤注册，公司参考项目同逻辑）。

import type { MenuTreeNode } from './index'

/**
 * 角色菜单映射：{ [role]: [菜单树] }
 * 权限系统未对接时，前端按用户角色取对应菜单树；角色未配置时回退 operator（默认普通运维者）。
 */
const roleMenuMap: Record<string, MenuTreeNode[]> = {
  // ===== operator 普通运维者（默认角色） =====
  operator: [
    { id: 1, name: '我的项目', path: '/operator/my-projects', component: 'views/operator/my-assets/MyProjects', icon: 'Monitor', sort: 1, hidden: false, group: '', children: [] },
    { id: 2, name: '我的任务', path: '/operator/my-task', component: 'views/operator/MyTask', icon: 'List', sort: 2, hidden: false, group: '', children: [] },
    { id: 3, name: '我的运维', path: '/operator/my-ops', component: 'views/operator/MyOps', icon: 'Setting', sort: 3, hidden: false, group: '', children: [] },
    { id: 4, name: '我的使用', path: '/operator/my-usage', component: 'views/operator/MyUsage', icon: 'Files', sort: 4, hidden: false, group: '', children: [] },
    { id: 5, name: '可信仓库', path: '/operator/trusted-repo', component: 'views/operator/TrustedRepo', icon: 'Files', sort: 5, hidden: false, group: '', children: [] },
  ],

  // ===== specialist 治理专员 =====
  specialist: [
    { id: 1, name: '我的项目', path: '/specialist/my-projects', component: 'views/operator/my-assets/MyProjects', icon: 'Monitor', sort: 1, hidden: false, group: '', children: [] },
    { id: 2, name: '我的任务', path: '/specialist/my-task', component: 'views/specialist/MyTask', icon: 'List', sort: 2, hidden: false, group: '', children: [] },
    {
      id: 3, name: '我的运维', path: '', component: '', icon: 'Setting', sort: 3, hidden: false, group: 'ops', children: [
        { id: 31, name: '运维台账', path: '/specialist/my-ops', component: 'views/specialist/MyOps', icon: '', sort: 1, hidden: false, group: '', children: [] },
        { id: 32, name: '人员管理', path: '/specialist/personnel-mgmt', component: 'views/operator/PersonnelMgmt', icon: '', sort: 2, hidden: false, group: '', children: [] },
      ],
    },
    { id: 4, name: '我的使用', path: '/specialist/my-usage', component: 'views/operator/MyUsage', icon: 'Files', sort: 4, hidden: false, group: '', children: [] },
    { id: 5, name: '可信仓库', path: '/specialist/trusted-repo', component: 'views/operator/TrustedRepo', icon: 'Files', sort: 5, hidden: false, group: '', children: [] },
  ],

  // ===== admin 管理专员（全局治理视图菜单） =====
  admin: [
    {
      id: 1, name: '台账管理', path: '', component: '', icon: 'Coin', sort: 1, hidden: false, group: 'asset', children: [
        { id: 11, name: '开源基础软件', path: '/admin/asset/software-library', component: 'views/admin/asset/SoftwareLibrary', icon: '', sort: 1, hidden: false, group: '', children: [] },
        { id: 12, name: '开源依赖包', path: '/admin/asset/dependency-library', component: 'views/admin/asset/DependencyLibrary', icon: '', sort: 2, hidden: false, group: '', children: [] },
        { id: 13, name: '开源工具', path: '/admin/asset/tool-library', component: 'views/admin/asset/ToolLibrary', icon: '', sort: 3, hidden: false, group: '', children: [] },
        { id: 14, name: '自研项目', path: '/admin/asset/self-developed', component: 'views/admin/asset/SelfDevelopedProjects', icon: '', sort: 4, hidden: false, group: '', children: [] },
        { id: 15, name: '源码备份', path: '/admin/asset/source-backup', component: 'views/admin/asset/SourceBackup', icon: '', sort: 5, hidden: false, group: '', children: [] },
      ],
    },
    { id: 2, name: '可信仓库', path: '/admin/trusted-repo', component: 'views/admin/asset/TrustedRepository', icon: 'Files', sort: 2, hidden: false, group: '', children: [] },
    {
      id: 3, name: '风险管理', path: '', component: '', icon: 'Warning', sort: 3, hidden: false, group: 'risk', children: [
        { id: 31, name: '名单管理', path: '/admin/risk/blacklist', component: 'views/admin/risk/Blacklist', icon: '', sort: 1, hidden: false, group: '', children: [] },
        { id: 32, name: 'CICD检查', path: '/admin/risk/cicd-check', component: 'views/admin/risk/CicdCheck', icon: '', sort: 2, hidden: false, group: '', children: [] },
        { id: 33, name: '许可证信息库', path: '/admin/risk/license-library', component: 'views/admin/risk/LicenseLibrary', icon: '', sort: 3, hidden: false, group: '', children: [] },
        { id: 34, name: '缺陷台账', path: '/admin/risk/defect-ledger', component: 'views/admin/risk/DefectLedger', icon: '', sort: 4, hidden: false, group: '', children: [] },
      ],
    },
    {
      id: 4, name: '治理中心', path: '', component: '', icon: 'Setting', sort: 4, hidden: false, group: 'governance', children: [
        { id: 41, name: '任务列表', path: '/admin/governance/task-list', component: 'views/admin/governance/TaskList', icon: '', sort: 1, hidden: false, group: '', children: [] },
      ],
    },
    {
      id: 5, name: '决策与洞察中心', path: '', component: '', icon: 'DataAnalysis', sort: 5, hidden: false, group: 'decision', children: [
        { id: 51, name: '缺陷修复效能', path: '/admin/decision/repair-efficiency', component: 'views/admin/decision/repair-efficiency', icon: '', sort: 1, hidden: false, group: '', children: [] },
        { id: 52, name: '供应链安全洞察', path: '/admin/decision/core-dependencies', component: 'views/admin/decision/core-dependencies', icon: '', sort: 2, hidden: false, group: '', children: [] },
      ],
    },
  ],
}

/**
 * 获取角色对应的默认菜单树（兜底用）。
 * 角色未配置时回退 operator（默认普通运维者）；未知角色同样回退，保证登录后必有菜单可渲染。
 */
export function getRoleMenuMap(role: string | undefined): MenuTreeNode[] {
  return roleMenuMap[role || ''] || roleMenuMap['operator']
}
