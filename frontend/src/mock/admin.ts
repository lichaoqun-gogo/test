// Webhook 事件（供活页面 views/admin/asset/TrustedRepository.vue 使用）
// 其余 mock 数据已随废弃页面删除（mock 目录仅保留此文件）。
export interface WebhookEvent {
  id: string
  eventType: string
  componentName: string
  version: string
  syncStatus: 'success' | 'failed' | 'pending'
  receivedAt: string
  metadata: string
}

export const mockWebhookEvents: WebhookEvent[] = [
  { id: 'wh001', eventType: 'package.published', componentName: 'lodash', version: '4.17.21', syncStatus: 'success', receivedAt: '2024-03-15 10:30', metadata: '{"registry":"npm","downloads":15234}' },
  { id: 'wh002', eventType: 'package.updated', componentName: 'Spring Boot', version: '3.2.3', syncStatus: 'success', receivedAt: '2024-03-14 14:20', metadata: '{"registry":"maven","checksum":"abc123"}' },
  { id: 'wh003', eventType: 'package.published', componentName: 'Vue.js', version: '3.4.21', syncStatus: 'success', receivedAt: '2024-03-14 09:15', metadata: '{"registry":"npm","downloads":28900}' },
  { id: 'wh004', eventType: 'package.vulnerability', componentName: 'Log4j2', version: '2.22.0', syncStatus: 'success', receivedAt: '2024-03-08 16:45', metadata: '{"cve":"CVE-2024-1001","severity":"high"}' },
]
