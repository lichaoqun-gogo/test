<template>
  <div class="activity-monitor">
    <h3>活跃度健康度监测</h3>
    <p class="desc">展示采集的 GitHub/GitLab 原始指标（Star/Issue/Commit 趋势），支持导出活跃度监测评估报告</p>

    <el-table :data="activityData" stripe style="margin-top: 16px">
      <el-table-column prop="componentName" label="组件名称" width="140" />
      <el-table-column prop="stars" label="Stars" width="100" />
      <el-table-column prop="forks" label="Forks" width="100" />
      <el-table-column prop="openIssues" label="开放Issue" width="100" />
      <el-table-column prop="commitsLastMonth" label="近30天提交" width="110" />
      <el-table-column prop="contributors" label="贡献者" width="90" />
      <el-table-column prop="lastRelease" label="最新发版" width="120" />
      <el-table-column label="趋势" min-width="200">
        <template #default="{ row }">
          <el-progress
            :percentage="Math.min(row.activityScore || 85, 100)"
            :color="row.activityScore > 80 ? '#67C23A' : '#E6A23C'"
            :stroke-width="10"
            :text-inside="true"
            :style="{ width: '160px' }"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120">
        <template #default><el-button type="primary" link size="small">查看详情</el-button></template>
      </el-table-column>
    </el-table>

    <div style="margin-top: 16px">
      <el-button type="primary"><el-icon><Download /></el-icon> 导出评估报告</el-button>
    </div>

    <el-divider content-position="left">Star 趋势</el-divider>
    <div class="chart-placeholder">
      <el-empty description="趋势图表区域（实际项目将集成 ECharts 展示 Star/Commit 趋势）" :image-size="80" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { mockActivityData } from '@/mock/admin'
import { Download } from '@element-plus/icons-vue'

const activityData = mockActivityData.map(d => ({
  ...d,
  activityScore: d.stars > 50000 ? 95 : 85,
}))
</script>

<style scoped>
.desc { color: #909399; font-size: 13px; margin: 8px 0 0; }
.chart-placeholder { padding: 40px; background: #fafafa; border-radius: 8px; }
</style>
