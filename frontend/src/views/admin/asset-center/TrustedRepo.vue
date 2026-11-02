<template>
  <div class="trusted-repo">
    <h3>可信仓库资产</h3>
    <p class="desc">配置 Nexus Webhook 监听规则，查看事件推送日志及元数据同步结果</p>

    <el-card style="margin-top: 16px">
      <template #header>Webhook 配置</template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="Webhook URL">https://nexus.company.com/webhook/osm-platform</el-descriptions-item>
        <el-descriptions-item label="状态"><el-tag type="success">已启用</el-tag></el-descriptions-item>
        <el-descriptions-item label="监听事件">package.published, package.updated, package.vulnerability</el-descriptions-item>
        <el-descriptions-item label="最后配置时间">2024-03-01</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-divider content-position="left">事件推送日志</el-divider>
    <el-table :data="events" stripe>
      <el-table-column prop="eventType" label="事件类型" width="180" />
      <el-table-column prop="componentName" label="组件名称" width="140" />
      <el-table-column prop="version" label="版本" width="100" />
      <el-table-column label="同步状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.syncStatus === 'success' ? 'success' : 'danger'" size="small">
            {{ row.syncStatus === 'success' ? '成功' : '失败' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="receivedAt" label="接收时间" width="160" />
      <el-table-column prop="metadata" label="元数据" min-width="200">
        <template #default="{ row }">
          <el-text type="info" size="small" truncated>{{ row.metadata }}</el-text>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { mockWebhookEvents } from '@/mock/admin'
const events = mockWebhookEvents
</script>

<style scoped>
.desc { color: #909399; font-size: 13px; margin: 8px 0 0; }
</style>
