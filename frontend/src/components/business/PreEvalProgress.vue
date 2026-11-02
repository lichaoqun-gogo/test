<template>
  <div class="pre-eval-progress">
    <!-- 新接口：steps + currentStep -->
    <template v-if="steps && steps.length > 0">
      <div class="progress-steps">
        <div
          v-for="(step, index) in steps"
          :key="index"
          class="progress-step"
          :class="{
            'step-done': index < currentStep,
            'step-active': index === currentStep,
            'step-pending': index > currentStep,
          }"
        >
          <div class="step-icon-wrapper">
            <el-icon v-if="index < currentStep" class="step-icon done">
              <CircleCheckFilled />
            </el-icon>
            <el-icon v-else-if="index === currentStep" class="step-icon active">
              <Loading />
            </el-icon>
            <el-icon v-else class="step-icon pending">
              <Clock />
            </el-icon>
          </div>
          <div class="step-info">
            <span class="step-label">{{ typeof step === 'string' ? step : step.label }}</span>
            <span v-if="index < currentStep" class="step-status success">已完成</span>
            <span v-else-if="index === currentStep" class="step-status active">进行中</span>
            <span v-else class="step-status pending">待执行</span>
          </div>
        </div>
      </div>

      <!-- 进度条 -->
      <div class="progress-bar-wrapper">
        <el-progress
          :percentage="progressPercentage"
          :stroke-width="8"
          :show-text="false"
          :color="progressColors"
        />
        <span class="progress-text">{{ currentStep + 1 }} / {{ steps.length }}</span>
      </div>
    </template>

    <!-- 旧接口：status -->
    <template v-else>
      <div v-if="status === 'pending'" class="eval-step">
        <el-icon class="step-icon pending"><Clock /></el-icon>
        <span class="step-text">等待预评估...</span>
      </div>

      <template v-if="status === 'scanning' || status === 'completed'">
        <div
          v-for="step in legacyProgressSteps"
          :key="step.key"
          class="eval-step"
          :class="{ active: step.isActive, done: step.isDone }"
        >
          <el-icon v-if="step.isDone" class="step-icon done"><CircleCheckFilled /></el-icon>
          <el-icon v-else-if="step.isActive" class="step-icon loading"><Loading /></el-icon>
          <el-icon v-else class="step-icon pending"><Clock /></el-icon>
          <span class="step-text">{{ step.label }}</span>
        </div>
      </template>

      <div v-if="status === 'failed'" class="eval-step failed">
        <el-icon class="step-icon error"><CircleCloseFilled /></el-icon>
        <span class="step-text">预评估失败，请重试</span>
      </div>

      <div v-if="status === 'completed' && conclusion" class="eval-conclusion">
        <el-divider />
        <div class="conclusion-header">
          <el-icon :size="20" :class="conclusionIconClass"><component :is="conclusionIcon" /></el-icon>
          <span class="conclusion-text">{{ conclusionLabel }}</span>
          <el-tag v-if="score !== undefined" :type="scoreType" class="score-tag">
            评分: {{ score }}
          </el-tag>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { CircleCheckFilled, CircleCloseFilled, Clock, Loading, SuccessFilled, WarningFilled, CircleClose } from '@element-plus/icons-vue'
import type { PreEvalStatus, PreEvalConclusion } from '@/types'

interface Step {
  label: string
  description?: string
}

const props = withDefaults(defineProps<{
  steps?: (string | Step)[]
  currentStep?: number
  status?: PreEvalStatus
  conclusion?: PreEvalConclusion
  score?: number
}>(), {
  currentStep: 0,
})

// 新接口计算
const progressPercentage = computed(() => {
  if (!props.steps || props.steps.length === 0) return 0
  return Math.min(100, Math.round(((props.currentStep + 1) / props.steps.length) * 100))
})

const progressColors = [
  { color: '#409EFF', percentage: 30 },
  { color: '#E6A23C', percentage: 60 },
  { color: '#67C23A', percentage: 100 },
]

// 旧接口计算
const legacyProgressSteps = [
  { key: 'sca', label: 'SCA扫描', isActive: props.status === 'scanning', isDone: props.status === 'completed' },
  { key: 'blacklist', label: '黑白名单匹配', isActive: props.status === 'scanning', isDone: props.status === 'completed' },
  { key: 'license', label: '许可证合规性检查', isActive: props.status === 'scanning', isDone: props.status === 'completed' },
  { key: 'activity', label: '活跃度指标分析', isActive: props.status === 'scanning', isDone: props.status === 'completed' },
]

const conclusionLabels: Record<PreEvalConclusion, string> = {
  pass: '预评估通过 - 各项指标符合标准，建议审批通过',
  warning: '风险警告 - 存在部分风险项，建议人工复核',
  reject: '预评估不通过 - 存在严重风险，建议驳回',
}

const conclusionLabel = computed(() => props.conclusion ? conclusionLabels[props.conclusion] : '')

const conclusionIcon = computed(() => {
  if (props.conclusion === 'pass') return SuccessFilled
  if (props.conclusion === 'warning') return WarningFilled
  return CircleClose
})

const conclusionIconClass = computed(() => {
  if (props.conclusion === 'pass') return 'icon-success'
  if (props.conclusion === 'warning') return 'icon-warning'
  return 'icon-error'
})

const scoreType = computed(() => {
  if (!props.score) return 'info'
  if (props.score >= 80) return 'success'
  if (props.score >= 60) return 'warning'
  return 'danger'
})
</script>

<style scoped>
.pre-eval-progress {
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;
}

.progress-steps {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.progress-step {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  transition: all 0.3s ease;
}

.step-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.step-active .step-icon-wrapper {
  background: #409EFF;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.4);
  animation: pulse 1.5s ease-in-out infinite;
}

.step-done .step-icon-wrapper {
  background: #67C23A;
  box-shadow: 0 2px 8px rgba(103, 194, 58, 0.3);
}

.step-icon {
  font-size: 16px;
}

.step-icon.done {
  color: #fff;
}

.step-icon.active {
  color: #fff;
  animation: spin 1s linear infinite;
}

.step-icon.pending {
  color: #C0C4CC;
}

.step-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.step-label {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.step-pending .step-label {
  color: #C0C4CC;
}

.step-status {
  font-size: 12px;
}

.step-status.success {
  color: #67C23A;
}

.step-status.active {
  color: #409EFF;
  font-weight: 500;
}

.step-status.pending {
  color: #C0C4CC;
}

.progress-bar-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #EBEEF5;
}

.progress-bar-wrapper :deep(.el-progress) {
  flex: 1;
}

.progress-text {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
  white-space: nowrap;
}

/* 旧接口样式 */
.eval-step {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 0;
}

.step-icon.loading {
  animation: spin 1s linear infinite;
}

.step-icon.error {
  color: #F56C6C;
  font-size: 20px;
}

.step-text {
  font-size: 14px;
  color: #606266;
}

.eval-conclusion {
  margin-top: 12px;
}

.conclusion-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 0;
}

.conclusion-text {
  font-weight: 500;
  font-size: 15px;
}

.icon-success {
  color: #67C23A;
}

.icon-warning {
  color: #E6A23C;
}

.icon-error {
  color: #F56C6C;
}

.score-tag {
  margin-left: auto;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@keyframes pulse {
  0%, 100% {
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.4);
  }
  50% {
    box-shadow: 0 2px 16px rgba(64, 158, 255, 0.6);
  }
}
</style>
