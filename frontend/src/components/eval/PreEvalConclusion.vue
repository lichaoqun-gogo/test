<template>
  <div class="pre-eval-conclusion" :class="conclusionClass">
    <div class="conclusion-card">
      <div class="conclusion-icon-wrapper">
        <el-icon :size="48" :class="iconClass">
          <component :is="conclusionIcon" />
        </el-icon>
      </div>

      <h3 class="conclusion-title">{{ conclusionTitle }}</h3>
      <p class="conclusion-desc">{{ conclusionDescription }}</p>

      <div class="conclusion-details">
        <div v-for="detail in conclusionDetails" :key="detail.label" class="detail-item">
          <span class="detail-label">{{ detail.label }}</span>
          <el-tag :type="detail.tagType as any" size="small">{{ detail.value }}</el-tag>
        </div>
      </div>

      <div v-if="conclusionRecommendations.length > 0" class="conclusion-recommendations">
        <h4>处理建议</h4>
        <ul>
          <li v-for="(rec, index) in conclusionRecommendations" :key="index">{{ rec }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { SuccessFilled, WarningFilled, CircleCloseFilled } from '@element-plus/icons-vue'
import type { PreEvalConclusion } from '@/types'

const props = withDefaults(defineProps<{
  conclusion: PreEvalConclusion
  componentName?: string
}>(), {
  componentName: '',
})

const conclusionClass = computed(() => {
  const classes: Record<PreEvalConclusion, string> = {
    pass: 'conclusion-pass',
    warning: 'conclusion-warning',
    reject: 'conclusion-reject',
  }
  return classes[props.conclusion]
})

const iconClass = computed(() => {
  const classes: Record<PreEvalConclusion, string> = {
    pass: 'icon-success',
    warning: 'icon-warning',
    reject: 'icon-error',
  }
  return classes[props.conclusion]
})

const conclusionIcon = computed(() => {
  const icons: Record<PreEvalConclusion, any> = {
    pass: SuccessFilled,
    warning: WarningFilled,
    reject: CircleCloseFilled,
  }
  return icons[props.conclusion]
})

const conclusionTitle = computed(() => {
  const titles: Record<PreEvalConclusion, string> = {
    pass: '预评估通过',
    warning: '风险警告',
    reject: '预评估不通过',
  }
  return titles[props.conclusion]
})

const conclusionDescription = computed(() => {
  const descriptions: Record<PreEvalConclusion, string> = {
    pass: '各项指标符合标准，建议审批通过',
    warning: '存在部分风险项，建议人工复核后决定是否通过',
    reject: '存在严重风险，建议驳回申请',
  }
  return descriptions[props.conclusion]
})

const conclusionDetails = computed(() => {
  const componentName = props.componentName || '当前组件'
  switch (props.conclusion) {
    case 'pass':
      return [
        { label: 'SCA扫描', value: '无高危漏洞', tagType: 'success' },
        { label: '黑白名单', value: '未命中黑名单', tagType: 'success' },
        { label: '许可证合规', value: '低风险', tagType: 'success' },
        { label: '活跃度', value: '社区活跃', tagType: 'success' },
      ]
    case 'warning':
      return [
        { label: 'SCA扫描', value: '存在中低风险漏洞', tagType: 'warning' },
        { label: '黑白名单', value: '未命中黑名单', tagType: 'success' },
        { label: '许可证合规', value: '中等风险', tagType: 'warning' },
        { label: '活跃度', value: '社区活跃度一般', tagType: 'warning' },
      ]
    case 'reject':
      return [
        { label: 'SCA扫描', value: '存在高危漏洞', tagType: 'danger' },
        { label: '黑白名单', value: '命中黑名单', tagType: 'danger' },
        { label: '许可证合规', value: '高风险', tagType: 'danger' },
        { label: '活跃度', value: '社区不活跃', tagType: 'danger' },
      ]
  }
})

const conclusionRecommendations = computed(() => {
  switch (props.conclusion) {
    case 'pass':
      return [
        '各项指标均符合要求，可以提交审批',
        '建议在引入后持续关注组件安全动态',
      ]
    case 'warning':
      return [
        '请仔细评估报告中的风险项',
        '建议与开源治理专员沟通确认风险可控性',
        '如确认风险可接受，可提交审批并备注风险处理方案',
      ]
    case 'reject':
      return [
        '存在严重风险，不建议引入该组件',
        '建议寻找替代组件或升级至安全版本后重新申请',
        '如有特殊情况，请联系管理员进行人工审核',
      ]
  }
})
</script>

<style scoped>
.pre-eval-conclusion {
  border-radius: 12px;
  overflow: hidden;
}

.conclusion-card {
  padding: 32px 24px;
  text-align: center;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 12px;
}

.conclusion-icon-wrapper {
  margin-bottom: 16px;
}

.icon-success {
  color: #67C23A;
  animation: scaleIn 0.5s ease;
}

.icon-warning {
  color: #E6A23C;
  animation: scaleIn 0.5s ease;
}

.icon-error {
  color: #F56C6C;
  animation: scaleIn 0.5s ease;
}

@keyframes scaleIn {
  0% {
    transform: scale(0);
    opacity: 0;
  }
  50% {
    transform: scale(1.2);
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

.conclusion-title {
  font-size: 22px;
  font-weight: 700;
  margin: 0 0 8px;
}

.conclusion-pass .conclusion-title {
  color: #67C23A;
}

.conclusion-warning .conclusion-title {
  color: #E6A23C;
}

.conclusion-reject .conclusion-title {
  color: #F56C6C;
}

.conclusion-desc {
  font-size: 14px;
  color: #606266;
  margin: 0 0 24px;
}

.conclusion-details {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
  margin-bottom: 24px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
}

.detail-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: #fff;
  border-radius: 6px;
}

.detail-label {
  font-size: 13px;
  color: #606266;
}

.conclusion-recommendations {
  text-align: left;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
}

.conclusion-recommendations h4 {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 12px;
}

.conclusion-recommendations ul {
  margin: 0;
  padding-left: 20px;
}

.conclusion-recommendations li {
  font-size: 13px;
  color: #606266;
  line-height: 1.8;
}
</style>
