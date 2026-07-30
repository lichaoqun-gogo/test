<template>
  <el-tag
    :color="riskColor"
    effect="dark"
    :size="size"
    class="risk-tag"
  >
    {{ riskLabel }}
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { RiskLevel } from '@/types'
import { RISK_LEVEL_LABELS, RISK_LEVEL_COLORS } from '@/types'

const props = withDefaults(defineProps<{
  level?: RiskLevel
  risk?: RiskLevel
  size?: 'large' | 'default' | 'small'
  effect?: 'dark' | 'light' | 'plain'
}>(), {
  size: 'default',
  effect: 'light',
})

const effectiveLevel = computed(() => props.risk ?? props.level ?? 'low')

const riskLabel = computed(() => {
  const labels: Record<RiskLevel, string> = {
    low: '低风险',
    medium: '中风险',
    high: '高风险',
    critical: '严重风险',
  }
  return labels[effectiveLevel.value] || RISK_LEVEL_LABELS[effectiveLevel.value] || ''
})

const riskColor = computed(() => RISK_LEVEL_COLORS[effectiveLevel.value] || '#909399')
</script>

<style scoped>
.risk-tag {
  font-weight: 600;
  border: none;
}
</style>
