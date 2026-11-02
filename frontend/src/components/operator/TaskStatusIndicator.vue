<template>
  <div class="task-status-indicator">
    <span
      class="status-dot"
      :style="{ backgroundColor: dotColor }"
    />
    <span class="status-text" :style="{ color: textColor }">
      {{ statusLabel }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { TaskStatus } from '@/types'
import { TASK_STATUS_LABELS, TASK_STATUS_COLORS } from '@/types'

const props = defineProps<{
  status: TaskStatus
}>()

const statusLabel = computed(() => TASK_STATUS_LABELS[props.status] || '')
const textColor = computed(() => TASK_STATUS_COLORS[props.status] || '#909399')
const dotColor = computed(() => TASK_STATUS_COLORS[props.status] || '#909399')
</script>

<style scoped>
.task-status-indicator {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-text {
  font-size: 13px;
  font-weight: 500;
}
</style>
