import axios from 'axios'

// 创建axios实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器：添加用户信息（从localStorage获取）
// CAS 模式下使用 Authorization Bearer token，同时保留 X-User-* 请求头兼容
api.interceptors.request.use(config => {
  const userInfo = localStorage.getItem('osm_user')
  const token = localStorage.getItem('osm_token')
  if (userInfo) {
    const user = JSON.parse(userInfo)
    const userName = user.nickname || user.name || ''
    // X-User-ID 统一使用 CAS 用户名（username）作为用户身份标识，与后端 JWT user_id 保持一致
    config.headers['X-User-ID'] = user.username || user.id || ''
    config.headers['X-User-Name'] = encodeURIComponent(userName)
    config.headers['X-User-Role'] = user.role || ''
    config.headers['X-User-Dept'] = encodeURIComponent(user.department || '')
  }
  // CAS 模式下 token 是 JWT，附加到 Authorization 头
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`
  }
  return config
})

// 响应拦截器
api.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 0) {
      console.error('API Error:', res.message)
      return Promise.reject(new Error(res.message))
    }
    return res.data  // 返回data字段
  },
  error => {
    console.error('Request Error:', error)
    return Promise.reject(error)
  }
)

// 导出类型
export interface SoftwareItem {
  id: number
  _group: string
  _name: string
  _version: string
  PUBLISHER: string
  CLASSIFIER: string
  type: string
  license: string
  applicable_level: string
  maint_dept: string
  use_dept: string
  country: string
  ACTIVE: number
  AUTHOR: string
  PURL: string
  CPE: string
  SWIDTAGID: string
  DESCRIPTION: string
  category: string
  LAST_BOM_IMPORTED: string
  update_by: string
  update_time: string
  create_by: string
  create_time: string
}

// 台账管理API
export const getSoftwareList = async (params?: any): Promise<SoftwareItem[]> => {
  return await api.get('/asset/projects', { params })
}

export const getSoftwareDetail = async (id: number): Promise<SoftwareItem> => {
  return await api.get(`/asset/projects/${id}`)
}

export const createSoftware = async (data: any): Promise<void> => {
  await api.post('/asset/projects', data)
}

export const updateSoftware = async (id: number, data: any): Promise<void> => {
  await api.put(`/asset/projects/${id}`, data)
}

export const deleteSoftware = async (id: number): Promise<void> => {
  await api.delete(`/asset/projects/${id}`)
}

export const getProjectComponents = async (projectId: number, params?: any) => {
  return await api.get(`/asset/projects/${projectId}/components`, { params })
}

export const getProjectLicenseStats = async (projectId: number) => {
  return await api.get(`/asset/projects/${projectId}/license-stats`)
}

export const assetApi = {
  // 开源基础软件
  getProjects: (params: any) => api.get('/asset/projects', { params }),
  getProject: (id: number) => api.get(`/asset/projects/${id}`),
  createProject: (data: any) => api.post('/asset/projects', data),
  updateProject: (id: number, data: any) => api.put(`/asset/projects/${id}`, data),
  deleteProject: (id: number) => api.delete(`/asset/projects/${id}`),

  // 自研项目
  getInnerProjects: (params: any) => api.get('/asset/inner-projects', { params }),
  getInnerProject: (id: number) => api.get(`/asset/inner-projects/${id}`),
  createInnerProject: (data: any) => api.post('/asset/inner-projects', data),
  updateInnerProject: (id: number, data: any) => api.put(`/asset/inner-projects/${id}`, data),
  deleteInnerProject: (id: number) => api.delete(`/asset/inner-projects/${id}`),

  // 开源依赖包
  getComponents: (params: any) => api.get('/asset/components', { params }),
  getComponent: (id: number) => api.get(`/asset/components/${id}`),
  getRelatedProjects: (componentId: number) => api.get(`/asset/components/${componentId}/related-projects`),
}

export const getDependencies = async (params?: any) => {
  return await api.get('/asset/components', { params })
}

export const getRelatedProjects = async (componentId: number) => {
  return await api.get(`/asset/components/${componentId}/related-projects`)
}

// 自研项目API
export const getInnerProjects = async (params?: any) => {
  return await api.get('/asset/inner-projects', { params })
}

export const updateInnerProject = async (id: number, data: any) => {
  return await api.put(`/asset/inner-projects/${id}`, data)
}

// 风险管理API
export const riskApi = {
  // 黑名单
  getBlacklist: (params: any) => api.get('/risk/blacklist', { params }),
  createBlacklist: (data: any) => api.post('/risk/blacklist', data),
  updateBlacklist: (id: number, data: any) => api.put(`/risk/blacklist/${id}`, data),
  deleteBlacklist: (id: number) => api.delete(`/risk/blacklist/${id}`),

  // 黑名单阻断记录
  getBlacklistRecords: (params: any) => api.get('/risk/blacklist/records', { params }),
  processBlacklistRecord: (id: number, data: any) => api.put(`/risk/blacklist/records/${id}/process`, data),

  // 许可证
  getLicenses: (params: any) => api.get('/risk/licenses', { params }),
  getLicenseGroups: () => api.get('/risk/license-groups'),

  // 缺陷台账
  getDefects: (params: any) => api.get('/risk/defects', { params }),
  getInternalDefects: (params: any) => api.get('/risk/defects/internal', { params }),
  getExternalDefects: (params: any) => api.get('/risk/defects/external', { params }),
  getDefect: (defectNo: string) => api.get(`/risk/defects/${encodeURIComponent(defectNo)}`),
  createDefect: (data: any) => api.post('/risk/defects', data),
  importDefects: (data: any) => api.post('/risk/defects/import', data),
  createVulnIntel: (data: any) => api.post('/risk/vuln-intel', data),
  matchVulnIntel: (data: any) => api.post('/risk/vuln-intel/match', data),
  importVulnIntel: (formData: FormData) =>
    api.post('/risk/vuln-intel/import', formData, { headers: { 'Content-Type': 'multipart/form-data' } }),
  getDefectStatsOverview: (params: any) => api.get('/risk/defects/stats/overview', { params }),
  // defectNo 对外部缺陷形如 redis#14611，含 # 需 encodeURIComponent，否则被浏览器当作 URL 片段截断导致 404
  getAnalystCandidates: (defectNo: string) => api.get(`/risk/defects/${encodeURIComponent(defectNo)}/analyst-candidates`),
  assignAnalyst: (defectNo: string, data: any) => api.post(`/risk/defects/${encodeURIComponent(defectNo)}/assign-analyst`, data),
  autoAssignAnalyst: (defectNo: string) => api.post(`/risk/defects/${encodeURIComponent(defectNo)}/auto-assign-analyst`),
  autoAssignDisposals: (defectNo: string) => api.post(`/risk/defects/${encodeURIComponent(defectNo)}/auto-assign-disposals`),
  analyzeDefect: (defectNo: string, data: any) => api.post(`/risk/defects/${encodeURIComponent(defectNo)}/analyze`, data),
  reviewDefect: (defectNo: string, data: any) => api.post(`/risk/defects/${encodeURIComponent(defectNo)}/review`, data),
  getDefectDisposals: (defectNo: string) => api.get(`/risk/defects/${encodeURIComponent(defectNo)}/disposals`),
  getDefectHistory: (defectNo: string) => api.get(`/risk/defects/${encodeURIComponent(defectNo)}/history`),
  getDisposals: (params: any) => api.get('/risk/disposals', { params }),
  assignDisposal: (id: number, data: any) => api.post(`/risk/disposals/${id}/assign`, data),
  completeDisposal: (id: number, data: any) => api.post(`/risk/disposals/${id}/complete`, data),
  markDisposalNA: (id: number, data: any) => api.post(`/risk/disposals/${id}/na`, data),
}

// 任务管理API
export const taskApi = {
  getTasks: (params: any) => api.get('/tasks', { params }),
  getTask: (id: number) => api.get(`/tasks/${id}`),
  getTaskDetail: (id: number) => api.get(`/tasks/${id}/detail`),
  createTask: (data: any) => api.post('/tasks', data),
  updateTask: (id: number, data: any) => api.put(`/tasks/${id}`, data),
  updateTaskStatus: (id: number, data: any) => api.put(`/tasks/${id}/status`, data),
  withdrawTask: (id: number) => api.post(`/workflow/tasks/${id}/withdraw`),
  // 审批流相关
  submitTask: (id: number) => api.post(`/workflow/tasks/${id}/submit`),
  approveTask: (id: number, data: any) => api.post(`/workflow/tasks/${id}/approve`, data),
  rejectTask: (id: number, data: any) => api.post(`/workflow/tasks/${id}/reject`, data),
  getApprovalNodes: (id: number) => api.get(`/workflow/tasks/${id}/nodes`),
  getApprovalHistory: (id: number) => api.get(`/workflow/tasks/${id}/history`),
}

// 审批管理API
export const approvalApi = {
  approve: (taskId: number, data: any) => api.post(`/approvals/${taskId}/approve`, data),
  reject: (taskId: number, data: any) => api.post(`/approvals/${taskId}/reject`, data),
  getApprovalHistory: (taskId: number) => api.get(`/approvals/${taskId}/history`),
}

// 文件管理API
export const fileApi = {
  upload: (formData: FormData) => api.post('/files/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  }),
  download: (id: number) => api.get(`/files/${id}/download`, { responseType: 'blob' }),
  delete: (id: number) => api.delete(`/files/${id}`),
}

// 人员管理API
// 当前数据源：后端 users 表；后期对接统一权限系统时只需修改后端 GetUserList 实现
export const personnelApi = {
  // 获取人员列表（用于负责人选择下拉框）
  // params.dept_name 传入部门名称时只返回该部门人员
  getList: (params?: { dept_name?: string }) => api.get('/users', { params }),
}

// 项目-负责人关联API
// project_type: 'software'(开源基础软件) | 'inner'(自研项目)
export const projectOwnerApi = {
  // 查询某项目的所有负责人
  getOwners: (projectType: string, projectID: number) =>
    api.get('/asset/project-owners', { params: { project_type: projectType, project_id: projectID } }),
  // 批量查询多个项目的负责人映射 { [project_id]: owner_username[] }
  getOwnersMap: (projectType: string, projectIDs: number[]) =>
    api.get('/asset/project-owners/map', { params: { project_type: projectType, project_ids: projectIDs.join(',') } }),
  // 保存项目负责人（全量替换）
  saveOwners: (projectType: string, projectID: number, owners: { owner_username: string; dept_name: string; team: string }[]) =>
    api.post('/asset/project-owners', { project_type: projectType, project_id: projectID, owners }),
}

// 源码备份API
export const sourceBackupApi = {
  // 列出对象（递归列出+按层级分组）
  listObjects: (params: { prefix?: string; query?: string }) =>
    api.get('/source-backups/objects', { params }),
  // 下载单个对象（直接通过浏览器下载，后端流式传输）
  downloadObject: (key: string) => {
    const url = `/api/v1/source-backups/download?key=${encodeURIComponent(key)}`
    const a = document.createElement('a')
    a.href = url
    a.download = ''
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
  },
  // 下载目录下所有文件（打包为 zip）
  downloadFolder: (prefix: string) => {
    const url = `/api/v1/source-backups/download-folder?prefix=${encodeURIComponent(prefix)}`
    const a = document.createElement('a')
    a.href = url
    a.download = ''
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
  },
  // 上传测试数据
  seedTestObjects: () =>
    api.post('/source-backups/seed'),
}

// 核心依赖统计API
export const coreDependencyApi = {
  getSummary: () => api.get('/asset/core-dependencies/summary'),
  getAssoc: (params: { page?: number; pageSize?: number; orgName?: string; compName?: string; projectName?: string }) =>
    api.get('/asset/core-dependencies/assoc', { params }),
  getRanking: (params: { page?: number; pageSize?: number; orgName?: string; compName?: string }) =>
    api.get('/asset/core-dependencies/ranking', { params }),
  getVersionFragmentation: (params: { page?: number; pageSize?: number; compName?: string }) =>
    api.get('/asset/core-dependencies/version-fragmentation', { params }),
  getHighRiskProjects: (params: { page?: number; pageSize?: number; projectName?: string }) =>
    api.get('/asset/core-dependencies/high-risk-projects', { params }),
  getBlacklistDetail: (projectName: string) =>
    api.get('/asset/core-dependencies/blacklist-detail', { params: { projectName } }),
  getLicenseOverview: (params: { page?: number; pageSize?: number; licenseName?: string }) =>
    api.get('/asset/core-dependencies/license-overview', { params }),
}

export default api
