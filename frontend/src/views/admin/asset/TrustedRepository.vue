<template>
  <div class="trusted-repository">
    <!-- Webhook 配置卡片 -->
    <el-card class="webhook-card">
      <template #header>
        <div class="card-header">
          <span><el-icon :size="18" style="margin-right: 8px"><Connection /></el-icon> Nexus Webhook 配置</span>
          <el-button type="primary" link @click="editWebhookVisible = true">编辑配置</el-button>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="Webhook URL">
          <el-text type="info" copyable>https://nexus.company.com/webhook/osm-platform</el-text>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag type="success">已启用</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="监听事件" :span="2">
          <el-tag v-for="evt in webhookConfig.events" :key="evt" size="small" style="margin-right: 4px; margin-bottom: 4px">
            {{ evt }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="Secret Token">
          <el-text type="info">****-****-****-{{ webhookConfig.secretToken.slice(-4) }}</el-text>
        </el-descriptions-item>
        <el-descriptions-item label="最后配置时间">{{ formatTime(webhookConfig.lastConfigured) }}</el-descriptions-item>
        <el-descriptions-item label="累计推送次数">{{ webhookConfig.totalPushes }}</el-descriptions-item>
        <el-descriptions-item label="推送成功率">{{ webhookConfig.successRate }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 事件推送日志 -->
    <el-card style="margin-top: 16px">
      <template #header>
        <div class="card-header">
          <span><el-icon :size="18" style="margin-right: 8px"><Document /></el-icon> 事件推送日志</span>
          <div class="header-actions">
            <el-select v-model="eventFilter" placeholder="事件类型" clearable style="width: 160px; margin-right: 12px" @change="handleEventSearch">
              <el-option v-for="t in eventTypes" :key="t" :label="t" :value="t" />
            </el-select>
            <el-select v-model="syncFilter" placeholder="同步状态" clearable style="width: 120px; margin-right: 12px" @change="handleEventSearch">
              <el-option label="成功" value="success" />
              <el-option label="失败" value="failed" />
              <el-option label="待处理" value="pending" />
            </el-select>
            <el-button type="primary" plain @click="handleEventSearch">查询</el-button>
          </div>
        </div>
      </template>

      <el-table :data="paginatedEvents" stripe>
        <el-table-column prop="eventType" label="事件类型" width="180">
          <template #default="{ row }">
            <el-tag :type="getEventTypeColor(row.eventType)" size="small">{{ row.eventType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="componentName" label="组件名称" width="140" />
        <el-table-column prop="version" label="版本" width="100" />
        <el-table-column label="同步状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.syncStatus === 'success' ? 'success' : row.syncStatus === 'failed' ? 'danger' : 'warning'" size="small">
              {{ row.syncStatus === 'success' ? '成功' : row.syncStatus === 'failed' ? '失败' : '待处理' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="receivedAt" label="接收时间" width="160">
          <template #default="{ row }">{{ formatTime(row.receivedAt) }}</template>
        </el-table-column>
        <el-table-column prop="metadata" label="元数据" min-width="220">
          <template #default="{ row }">
            <el-popover trigger="click" placement="left" width="400">
              <template #reference>
                <el-button type="primary" link size="small">查看</el-button>
              </template>
              <pre class="metadata-json">{{ formatJson(row.metadata) }}</pre>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button v-if="row.syncStatus === 'failed'" type="primary" link size="small" @click="handleResync(row)">重推</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="eventPagination.page"
          v-model:page-size="eventPagination.pageSize"
          :total="eventPagination.total"
          :page-sizes="[5, 10, 20]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleEventSizeChange"
          @current-change="handleEventPageChange"
        />
      </div>
    </el-card>

    <!-- 元数据同步结果 -->
    <el-card style="margin-top: 16px">
      <template #header>
        <span><el-icon :size="18" style="margin-right: 8px"><Refresh /></el-icon> 元数据同步结果</span>
      </template>

      <el-table :data="syncResults" stripe>
        <el-table-column prop="source" label="数据源" width="140" />
        <el-table-column prop="componentName" label="组件名称" width="140" />
        <el-table-column prop="version" label="版本" width="100" />
        <el-table-column label="同步状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'synced' ? 'success' : row.status === 'partial' ? 'warning' : 'danger'" size="small">
              {{ row.status === 'synced' ? '已同步' : row.status === 'partial' ? '部分同步' : '同步失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="syncTime" label="同步时间" width="160">
          <template #default="{ row }">{{ formatTime(row.syncTime) }}</template>
        </el-table-column>
        <el-table-column prop="details" label="同步详情" min-width="200">
          <template #default="{ row }">
            <el-text type="info" size="small">{{ row.details }}</el-text>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 编辑 Webhook 配置弹窗 -->
    <el-dialog v-model="editWebhookVisible" title="编辑 Webhook 配置" width="560px">
      <el-form label-width="120px">
        <el-form-item label="Webhook URL">
          <el-input v-model="webhookForm.url" />
        </el-form-item>
        <el-form-item label="启用状态">
          <el-switch v-model="webhookForm.enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
        <el-form-item label="监听事件">
          <el-checkbox-group v-model="webhookForm.events">
            <el-checkbox label="package.published" />
            <el-checkbox label="package.updated" />
            <el-checkbox label="package.vulnerability" />
            <el-checkbox label="package.deleted" />
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="Secret Token">
          <el-input v-model="webhookForm.secretToken" placeholder="留空则不修改" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editWebhookVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveWebhook">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive } from 'vue'
import { mockWebhookEvents } from '@/mock/admin'
import { Connection, Document, Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { formatTime } from '@/utils/format'

// Webhook 配置假数据
const webhookConfig = reactive({
  url: 'https://nexus.company.com/webhook/osm-platform',
  enabled: true,
  events: ['package.published', 'package.updated', 'package.vulnerability'],
  secretToken: 'a1b2c3d4-e5f6',
  lastConfigured: '2024-03-01',
  totalPushes: 12580,
  successRate: '99.7%',
})

const editWebhookVisible = ref(false)
const webhookForm = reactive({
  url: webhookConfig.url,
  enabled: webhookConfig.enabled,
  events: [...webhookConfig.events],
  secretToken: '',
})

// 事件推送日志
const eventFilter = ref('')
const syncFilter = ref('')

const eventTypes = computed(() => {
  const types = new Set(mockWebhookEvents.map(e => e.eventType))
  return Array.from(types)
})

const filteredEvents = computed(() => {
  let list = [...mockWebhookEvents]
  if (eventFilter.value) list = list.filter(e => e.eventType === eventFilter.value)
  if (syncFilter.value) list = list.filter(e => e.syncStatus === syncFilter.value)
  return list
})

const eventPagination = reactive({
  page: 1,
  pageSize: 10,
  total: filteredEvents.value.length,
})

const paginatedEvents = computed(() => {
  const start = (eventPagination.page - 1) * eventPagination.pageSize
  return filteredEvents.value.slice(start, start + eventPagination.pageSize)
})

function handleEventSearch() {
  eventPagination.page = 1
  eventPagination.total = filteredEvents.value.length
}

function handleEventSizeChange() {
  eventPagination.page = 1
}

function handleEventPageChange() {}

function getEventTypeColor(type: string) {
  if (type.includes('vulnerability')) return 'danger'
  if (type.includes('deleted')) return 'info'
  return 'success'
}

function formatJson(jsonStr: string) {
  try {
    return JSON.stringify(JSON.parse(jsonStr), null, 2)
  } catch {
    return jsonStr
  }
}

function handleResync(row: any) {
  ElMessage.info(`正在重新推送事件 "${row.eventType}"`)
}

// 元数据同步结果假数据
const syncResults = ref([
  {
    id: 'sr001',
    source: 'npm Registry',
    componentName: 'lodash',
    version: '4.17.21',
    status: 'synced',
    syncTime: '2024-03-15 10:30:25',
    details: '同步完成：版本信息、依赖树、许可证、安全公告',
  },
  {
    id: 'sr002',
    source: 'Maven Central',
    componentName: 'Spring Boot',
    version: '3.2.3',
    status: 'synced',
    syncTime: '2024-03-14 14:20:18',
    details: '同步完成：POM信息、传递依赖、安全漏洞',
  },
  {
    id: 'sr003',
    source: 'npm Registry',
    componentName: 'Vue.js',
    version: '3.4.21',
    status: 'partial',
    syncTime: '2024-03-14 09:15:42',
    details: '部分同步：缺少安全公告信息',
  },
  {
    id: 'sr004',
    source: 'Maven Central',
    componentName: 'Log4j2',
    version: '2.22.0',
    status: 'synced',
    syncTime: '2024-03-08 16:45:30',
    details: '同步完成：CVE信息已更新',
  },
  {
    id: 'sr005',
    source: 'PyPI',
    componentName: 'requests',
    version: '2.31.0',
    status: 'failed',
    syncTime: '2024-03-07 11:22:08',
    details: '同步失败：网络连接超时',
  },
  {
    id: 'sr006',
    source: 'npm Registry',
    componentName: 'axios',
    version: '1.6.7',
    status: 'synced',
    syncTime: '2024-03-06 08:10:55',
    details: '同步完成：版本信息、依赖树、许可证',
  },
])

function handleSaveWebhook() {
  ElMessage.success('Webhook 配置已保存')
  editWebhookVisible.value = false
}
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.header-actions {
  display: flex;
  align-items: center;
}
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
.metadata-json {
  background: #f5f7fa;
  padding: 12px;
  border-radius: 4px;
  font-size: 12px;
  line-height: 1.6;
  margin: 0;
  max-height: 300px;
  overflow-y: auto;
}
.webhook-card :deep(.el-descriptions__label) {
  font-weight: 500;
}
</style>
