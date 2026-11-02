<template>
  <div class="smart-recommendation">
    <div class="recommendation-header">
      <el-icon class="recommend-icon"><MagicStick /></el-icon>
      <span class="recommend-title">智能推荐 - 相似低风险替代组件</span>
      <el-tag type="success" size="small" effect="plain">AI 驱动</el-tag>
    </div>

    <el-row :gutter="16">
      <el-col
        v-for="item in recommendations"
        :key="item.id"
        :xs="24"
        :sm="12"
        :lg="8"
      >
        <el-card class="recommend-card" shadow="hover" @click="handleApply(item)">
          <div class="recommend-card-header">
            <span class="recommend-name">{{ item.name }}</span>
            <el-tag size="small">{{ item.version }}</el-tag>
          </div>

          <p class="recommend-desc">{{ item.description }}</p>

          <div class="recommend-metrics">
            <div class="metric-item">
              <span class="metric-label">活跃度</span>
              <el-progress
                :percentage="item.activityScore"
                :color="getActivityColor(item.activityScore)"
                :stroke-width="6"
                :show-text="false"
              />
              <span class="metric-value">{{ item.activityScore }}</span>
            </div>

            <div class="metric-item">
              <span class="metric-label">已知漏洞</span>
              <span class="metric-value" :class="{ 'vuln-zero': item.vulnerabilityCount === 0 }">
                {{ item.vulnerabilityCount === 0 ? '无' : `${item.vulnerabilityCount}个` }}
              </span>
            </div>

            <div class="metric-item">
              <span class="metric-label">许可证</span>
              <div class="license-info">
                <span>{{ item.licenseName }}</span>
                <RiskTag :level="item.licenseRisk" size="small" />
              </div>
            </div>
          </div>

          <div class="recommend-footer">
            <div class="similarity-info">
              <el-icon><Connection /></el-icon>
              <span>相似度: <strong>{{ item.similarity }}%</strong></span>
            </div>
            <el-button type="primary" size="small" @click.stop="handleApply(item)">
              申请引入
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { MagicStick, Connection } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const props = defineProps<{
  recommendations: any[]
}>()

const router = useRouter()

function getActivityColor(score: number) {
  if (score >= 80) return '#67C23A'
  if (score >= 60) return '#E6A23C'
  return '#F56C6C'
}

function handleApply(item: any) {
  ElMessage.success(`正在跳转到引入申请页面...`)
  router.push({
    path: '/operator/my-task',
  })
}
</script>

<style scoped>
.smart-recommendation {
  margin-top: 24px;
  padding: 16px;
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  border-radius: 8px;
  border: 1px solid #bae6fd;
}

.recommendation-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #bae6fd;
}

.recommend-icon {
  font-size: 20px;
  color: #0ea5e9;
}

.recommend-title {
  font-size: 16px;
  font-weight: 600;
  color: #0c4a6e;
}

.recommend-card {
  margin-bottom: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #e2e8f0;
}

.recommend-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px -8px rgba(14, 165, 233, 0.15);
  border-color: #0ea5e9;
}

.recommend-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.recommend-name {
  font-size: 15px;
  font-weight: 600;
  color: #1e293b;
}

.recommend-desc {
  font-size: 13px;
  color: #64748b;
  margin: 0 0 12px 0;
  line-height: 1.5;
  min-height: 40px;
}

.recommend-metrics {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.metric-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.metric-label {
  font-size: 12px;
  color: #94a3b8;
  min-width: 50px;
}

.metric-value {
  font-size: 13px;
  font-weight: 500;
  color: #334155;
  min-width: 36px;
}

.metric-value.vuln-zero {
  color: #22c55e;
}

.license-info {
  display: flex;
  align-items: center;
  gap: 6px;
}

.recommend-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 10px;
  border-top: 1px solid #f1f5f9;
}

.similarity-info {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #64748b;
}

.similarity-info strong {
  color: #0ea5e9;
}
</style>
