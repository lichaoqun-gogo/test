<template>
  <div class="pre-eval-progress">
    <div class="progress-header">
      <el-icon :size="20" class="header-icon"><Monitor /></el-icon>
      <span class="header-title">预评估进度</span>
      <el-tag v-if="overallStatus === 'pending'" type="info" size="small">等待中</el-tag>
      <el-tag v-else-if="overallStatus === 'running'" type="warning" size="small">
        <el-icon class="is-loading"><Loading /></el-icon>
        评估中
      </el-tag>
      <el-tag v-else-if="overallStatus === 'completed'" type="success" size="small">已完成</el-tag>
      <el-tag v-else-if="overallStatus === 'failed'" type="danger" size="small">失败</el-tag>
    </div>

    <div class="progress-steps">
      <div
        v-for="(step, index) in steps"
        :key="step.key"
        class="step-item"
        :class="[step.status, { 'fade-in': step.status !== 'pending' }]"
        :style="{ animationDelay: `${index * 0.1}s` }"
      >
        <div class="step-indicator">
          <el-icon v-if="step.status === 'completed'" class="step-icon success">
            <CircleCheckFilled />
          </el-icon>
          <el-icon v-else-if="step.status === 'running'" class="step-icon running">
            <Loading />
          </el-icon>
          <el-icon v-else-if="step.status === 'failed'" class="step-icon error">
            <CircleCloseFilled />
          </el-icon>
          <div v-else class="step-icon pending">
            <div class="pending-dot"></div>
          </div>
        </div>

        <div class="step-content">
          <span class="step-label">{{ step.label }}</span>
          <span v-if="step.status === 'running'" class="step-status-text">检查中...</span>
          <span v-else-if="step.status === 'completed'" class="step-status-text success">完成</span>
          <span v-else-if="step.status === 'failed'" class="step-status-text error">异常</span>
        </div>

        <!-- 连接线 -->
        <div v-if="index < steps.length - 1" class="step-connector">
          <div class="connector-line" :class="{ active: steps[index + 1].status !== 'pending' }"></div>
        </div>
      </div>
    </div>

    <div v-if="overallStatus === 'failed'" class="error-tip">
      <el-icon><WarningFilled /></el-icon>
      <span>预评估执行失败，请检查网络或联系管理员</span>
      <el-button type="primary" size="small" @click="$emit('retry')">重新评估</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import {
  CircleCheckFilled,
  CircleCloseFilled,
  Loading,
  Monitor,
  WarningFilled,
} from '@element-plus/icons-vue'

export type EvalStepStatus = 'pending' | 'running' | 'completed' | 'failed'

export interface EvalStep {
  key: string
  label: string
  status: EvalStepStatus
}

const props = withDefaults(defineProps<{
  status?: 'idle' | 'running' | 'completed' | 'failed'
}>(), {
  status: 'idle',
})

const emit = defineEmits<{
  'retry': []
  'completed': []
}>()

const defaultSteps: EvalStep[] = [
  { key: 'sca', label: 'SCA扫描', status: 'pending' },
  { key: 'blacklist', label: '黑白名单匹配', status: 'pending' },
  { key: 'license', label: '许可证合规性检查', status: 'pending' },
  { key: 'activity', label: '活跃度指标分析', status: 'pending' },
]

const steps = ref<EvalStep[]>(JSON.parse(JSON.stringify(defaultSteps)))

const overallStatus = computed(() => {
  if (props.status === 'idle') return 'pending'
  if (props.status === 'running') return 'running'
  if (props.status === 'completed') return 'completed'
  if (props.status === 'failed') return 'failed'
  return 'pending'
})

// 模拟逐步执行预评估
function simulateEval() {
  steps.value = JSON.parse(JSON.stringify(defaultSteps))

  const stepDelays = [500, 1200, 2000, 2800]
  const completionDelay = 3500

  // 逐步更新每个步骤的状态
  stepDelays.forEach((delay, index) => {
    // 设置为running
    setTimeout(() => {
      if (props.status === 'running') {
        steps.value[index].status = 'running'
      }
    }, delay - 400)

    // 设置为completed
    setTimeout(() => {
      if (props.status === 'running' || props.status === 'completed') {
        steps.value[index].status = 'completed'
      }
    }, delay)
  })

  // 全部完成
  setTimeout(() => {
    if (props.status === 'completed') {
      emit('completed')
    }
  }, completionDelay)
}

// 监听状态变化，重新开始模拟
watch(() => props.status, (newStatus) => {
  if (newStatus === 'running') {
    simulateEval()
  } else if (newStatus === 'completed') {
    // 将所有步骤标记为完成
    steps.value.forEach(step => {
      step.status = 'completed'
    })
    emit('completed')
  } else if (newStatus === 'failed') {
    // 标记当前步骤为失败
    const runningStep = steps.value.find(s => s.status === 'running')
    if (runningStep) {
      runningStep.status = 'failed'
    }
  } else if (newStatus === 'idle') {
    steps.value = JSON.parse(JSON.stringify(defaultSteps))
  }
}, { immediate: true })

defineExpose({
  steps,
  reset: () => {
    steps.value = JSON.parse(JSON.stringify(defaultSteps))
  },
})
</script>

<style scoped>
.pre-eval-progress {
  background: linear-gradient(135deg, #f5f7fa 0%, #e8ecf1 100%);
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #e4e7ed;
}

.progress-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e4e7ed;
}

.header-icon {
  color: #409EFF;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.progress-steps {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.step-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  position: relative;
  opacity: 0.6;
  transition: all 0.3s ease;
}

.step-item.completed {
  opacity: 1;
}

.step-item.running {
  opacity: 1;
}

.step-item.failed {
  opacity: 1;
}

.fade-in {
  animation: fadeIn 0.4s ease forwards;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateX(-10px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.step-indicator {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.step-icon {
  font-size: 22px;
  transition: all 0.3s ease;
}

.step-icon.success {
  color: #67C23A;
  animation: popIn 0.3s ease;
}

@keyframes popIn {
  0% {
    transform: scale(0);
  }
  60% {
    transform: scale(1.2);
  }
  100% {
    transform: scale(1);
  }
}

.step-icon.running {
  color: #409EFF;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.step-icon.error {
  color: #F56C6C;
  animation: shake 0.5s ease;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-3px); }
  75% { transform: translateX(3px); }
}

.pending-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #C0C4CC;
  border: 2px solid #DCDFE6;
}

.step-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.step-label {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.step-item.completed .step-label {
  color: #303133;
}

.step-item.running .step-label {
  color: #409EFF;
}

.step-status-text {
  font-size: 12px;
  color: #909399;
}

.step-status-text.success {
  color: #67C23A;
}

.step-status-text.error {
  color: #F56C6C;
}

.step-connector {
  position: absolute;
  left: 13px;
  top: 40px;
  height: 20px;
  width: 2px;
}

.connector-line {
  width: 100%;
  height: 100%;
  background: #DCDFE6;
  transition: background 0.3s ease;
}

.connector-line.active {
  background: #67C23A;
}

.error-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 16px;
  padding: 12px;
  background: #FEF0F0;
  border-radius: 8px;
  color: #F56C6C;
  font-size: 13px;
}

.error-tip .el-button {
  margin-left: auto;
}
</style>
