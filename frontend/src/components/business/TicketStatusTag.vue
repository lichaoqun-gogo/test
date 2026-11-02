<template>
  <el-tag
    :type="statusType"
    :effect="effect"
    :size="size"
    class="ticket-status-tag"
  >
    {{ statusLabel }}
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { TicketStatus } from '@/types'
import { TICKET_STATUS_LABELS } from '@/types'

const props = withDefaults(defineProps<{
  status: TicketStatus
  size?: 'large' | 'default' | 'small'
  effect?: 'dark' | 'light' | 'plain'
}>(), {
  size: 'default',
  effect: 'light',
})

const statusLabel = computed(() => TICKET_STATUS_LABELS[props.status] || '')

const statusType = computed(() => {
  const types: Record<TicketStatus, 'success' | 'warning' | 'danger' | 'info' | undefined> = {
    pending: undefined,
    processing: 'warning',
    completed: 'success',
    rejected: 'danger',
    assigned: 'info',
  }
  return types[props.status]
})
</script>
