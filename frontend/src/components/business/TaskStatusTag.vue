<template>
  <el-tag
    :color="statusColor"
    effect="dark"
    :size="size"
    class="task-status-tag"
  >
    {{ statusLabel }}
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type TaskStatus = 'pending' | 'processing' | 'rejected' | 'completed' | 'assigned'

const props = withDefaults(defineProps<{
  status: TaskStatus
  size?: 'large' | 'default' | 'small'
}>(), {
  size: 'default',
})

const statusLabel = computed(() => {
  const labels: Record<TaskStatus, string> = {
    pending: '待审批',
    processing: '处理中',
    rejected: '已驳回',
    completed: '已完成',
    assigned: '已分配',
  }
  return labels[props.status] || ''
})

const statusColor = computed(() => {
  const colors: Record<TaskStatus, string> = {
    pending: '#409EFF',
    processing: '#409EFF',
    rejected: '#F56C6C',
    completed: '#67C23A',
    assigned: '#E6A23C',
  }
  return colors[props.status] || '#909399'
})
</script>

<style scoped>
.task-status-tag {
  font-weight: 600;
  border: none;
}
</style>
