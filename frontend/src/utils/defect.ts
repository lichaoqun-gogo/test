// 缺陷/漏洞处置流状态映射（与后端 service/disposal 状态常量保持一致）

// 父缺陷状态
export const DEFECT_STATUS_LABELS: Record<string, string> = {
  pending_analysis: '待研判',
  pending_review: '待上会评审',
  assigned: '待处置',
  ignored: '已忽略',
  closed: '已完成',
}

export const DEFECT_STATUS_TYPES: Record<string, DefectTagType> = {
  pending_analysis: 'warning',
  pending_review: 'info',
  assigned: 'warning',
  ignored: 'danger',
  closed: 'success',
}

// 子处置单状态
export const DISPOSAL_STATUS_LABELS: Record<string, string> = {
  pending_assign: '待处置',
  processing: '处理中',
  completed: '已完成',
  na: '忽略',
}

// el-tag type 允许的取值
export type DefectTagType = 'primary' | 'success' | 'info' | 'warning' | 'danger'

export const DISPOSAL_STATUS_TYPES: Record<string, DefectTagType> = {
  pending_assign: 'warning',
  processing: 'primary',
  completed: 'success',
  na: 'info',
}

// 流程类型
export const FLOW_TYPE_LABELS: Record<string, string> = {
  defect: '缺陷',
  vuln: '漏洞',
}

// 处置类型
export const DISPOSAL_TYPE_LABELS: Record<string, string> = {
  maint: '运维部门',
  use: '使用部门',
}

// 研判结论
export const ANALYZE_CONCLUSIONS = [
  { value: 'need_review', label: '需要上会评审' },
  { value: 'need_fix', label: '需要修复（直接分配）' },
  { value: 'ignore', label: '忽略' },
]

// 评审结论
export const REVIEW_CONCLUSIONS = [
  { value: 'need_fix', label: '需要修复（裂变处置）' },
  { value: 'no_action', label: '无需处置（忽略）' },
]

export function defectStatusLabel(status: string): string {
  return DEFECT_STATUS_LABELS[status] || status || '-'
}

export function defectStatusType(status: string): DefectTagType {
  return DEFECT_STATUS_TYPES[status] || 'info'
}

export function disposalStatusLabel(status: string): string {
  return DISPOSAL_STATUS_LABELS[status] || status || '-'
}

export function disposalStatusType(status: string): DefectTagType {
  return DISPOSAL_STATUS_TYPES[status] || 'info'
}

// 严重程度（后端 risk_level 存 critical/high/medium/low）
export function severityLabel(level: string): string {
  const map: Record<string, string> = {
    critical: '严重',
    high: '高危',
    medium: '中危',
    low: '低危',
  }
  return map[level] || level || '-'
}

export function severityType(level: string): DefectTagType {
  const map: Record<string, DefectTagType> = {
    critical: 'danger',
    high: 'danger',
    medium: 'warning',
    low: 'info',
  }
  return map[level] || 'info'
}

// 研判人展示：进入「待上会评审」及后续状态必然已完成研判，研判人必有人选；
// 优先显示平台内分配的研判人（dispatch_target -> 姓名），缺失时回退实际研判人（analyze_by），
// 杜绝「已研判却显示未分配」的矛盾；仅待研判（pending_analysis）未分配时才显示占位。
export function analystLabel(row: { dispatchTargetName?: string; analyzeBy?: string; status?: string }): string {
  if (row.dispatchTargetName) {
    return row.dispatchTargetName
  }
  const analyzed = !!row.status && row.status !== 'pending_analysis'
  if (analyzed && row.analyzeBy) {
    return row.analyzeBy
  }
  return ''
}
