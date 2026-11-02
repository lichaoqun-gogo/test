<template>
  <el-tooltip :content="gradeTooltip" placement="top" :disabled="!showTooltip">
    <el-tag
      :color="gradeColor"
      effect="dark"
      :size="size"
      class="grade-tag"
    >
      {{ grade }}级（{{ gradeLabel }}）
    </el-tag>
  </el-tooltip>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { AssetGrade } from '@/types'
import { GRADE_LABELS, GRADE_COLORS } from '@/types'

const props = withDefaults(defineProps<{
  grade: AssetGrade
  size?: 'large' | 'default' | 'small'
  showTooltip?: boolean
}>(), {
  size: 'default',
  showTooltip: true,
})

const gradeLabel = computed(() => GRADE_LABELS[props.grade] || '')
const gradeColor = computed(() => GRADE_COLORS[props.grade] || '#909399')

const gradeDescriptions: Record<AssetGrade, string> = {
  A: '可安全使用（低风险）：经过充分评估，可以安全使用',
  B: '谨慎使用（中低风险）：使用时需注意相关风险，建议制定缓解措施',
  C: '限期整改（高风险）：存在较高风险，需限期整改或替换',
  D: '禁止使用（高风险）：存在严重风险，禁止使用',
}

const gradeTooltip = computed(() => {
  if (!props.showTooltip) return ''
  return `${props.grade}级 - ${GRADE_LABELS[props.grade]} - ${gradeDescriptions[props.grade]}`
})
</script>

<style scoped>
.grade-tag {
  font-weight: 600;
  border: none;
}
</style>
