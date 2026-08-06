<template>
  <div class="source-backup">
    <div class="page-header">
      <span>源码自动备份状态追踪</span>
      <el-button type="primary" @click="handleTriggerBackup">
        <el-icon><RefreshRight /></el-icon> 手动触发备份
      </el-button>
    </div>

    <el-table :data="backups" stripe style="margin-top: 16px">
      <el-table-column prop="componentName" label="组件名称" width="140" />
      <el-table-column prop="version" label="版本" width="100" />
      <el-table-column label="备份状态" width="100">
        <template #default="{ row }">
          <el-tag :type="backupStatusType(row.status)" size="small">
            {{ backupStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="storagePath" label="存储路径" min-width="200" />
      <el-table-column prop="size" label="大小" width="80" />
      <el-table-column prop="lastUpdated" label="最后更新" width="180">
        <template #default="{ row }">{{ formatTime(row.lastUpdated) }}</template>
      </el-table-column>
      <el-table-column prop="triggeredBy" label="触发方式" width="80" />
      <el-table-column label="操作" width="80">
        <template #default="{ row }">
          <el-button type="primary" link size="small" v-if="row.status === 'failed'">重试</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RefreshRight } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { formatTime } from '@/utils/format'

const backups = ref<any[]>([])

function loadBackups() {}

onMounted(() => {
  loadBackups()
})

function backupStatusType(status: string) {
  const types: Record<string, string> = { completed: 'success', running: '', failed: 'danger', pending: 'info' }
  return types[status] || ''
}

function backupStatusLabel(status: string) {
  const labels: Record<string, string> = { completed: '已完成', running: '进行中', failed: '失败', pending: '待执行' }
  return labels[status] || status
}

function handleTriggerBackup() {
  ElMessage.info('备份任务已触发')
}
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; }
</style>
