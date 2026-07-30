<template>
  <div class="activity-monitor">
    <div class="page-header">
      <div class="header-left">
        <h3>活跃度健康度监测</h3>
        <p class="desc">展示采集的 GitHub/GitLab 原始指标（Star/Issue/Commit 趋势），支持导出活跃度监测评估报告</p>
      </div>
    </div>

    <el-row :gutter="16" style="margin-top: 16px; margin-bottom: 20px">
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="监测组件数" :value="activityData.length" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="高活跃度" :value="activityData.filter(d => d.stars > 30000).length" value-style="color: #67C23A" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="中活跃度" :value="activityData.filter(d => d.stars > 10000 && d.stars <= 30000).length" value-style="color: #E6A23C" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="低活跃度" :value="activityData.filter(d => d.stars <= 10000).length" value-style="color: #F56C6C" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 组件指标表格 -->
    <el-table :data="activityData" stripe style="margin-bottom: 20px">
      <el-table-column prop="componentName" label="组件名称" width="140" />
      <el-table-column prop="stars" label="Stars" width="100" sortable />
      <el-table-column prop="forks" label="Forks" width="100" sortable />
      <el-table-column prop="openIssues" label="开放Issue" width="100" />
      <el-table-column prop="commitsLastMonth" label="近30天提交" width="110" sortable />
      <el-table-column prop="contributors" label="贡献者" width="90" />
      <el-table-column prop="lastRelease" label="最新发版" width="120" />
      <el-table-column label="活跃度评分" width="160">
        <template #default="{ row }">
          <el-progress
            :percentage="row.activityScore"
            :color="row.activityScore > 80 ? '#67C23A' : row.activityScore > 50 ? '#E6A23C' : '#F56C6C'"
            :stroke-width="10"
            :text-inside="true"
            :style="{ width: '120px' }"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="showChart(row as ActivityItem)">趋势</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="margin-bottom: 16px">
      <el-button type="primary"><el-icon><Download /></el-icon> 导出评估报告</el-button>
    </div>

    <!-- 趋势图表弹窗 -->
    <el-dialog v-model="chartVisible" :title="`${selectedComponent?.componentName} - 趋势分析`" width="900px">
      <el-row :gutter="16">
        <el-col :span="12">
          <h4 style="text-align: center; margin-bottom: 8px">Star 趋势</h4>
          <div ref="starChartRef" style="width: 100%; height: 320px" />
        </el-col>
        <el-col :span="12">
          <h4 style="text-align: center; margin-bottom: 8px">Commit 趋势</h4>
          <div ref="commitChartRef" style="width: 100%; height: 320px" />
        </el-col>
      </el-row>
      <el-row :gutter="16" style="margin-top: 16px">
        <el-col :span="12">
          <h4 style="text-align: center; margin-bottom: 8px">Issue 趋势</h4>
          <div ref="issueChartRef" style="width: 100%; height: 300px" />
        </el-col>
        <el-col :span="12">
          <h4 style="text-align: center; margin-bottom: 8px">综合健康度</h4>
          <div ref="healthChartRef" style="width: 100%; height: 300px" />
        </el-col>
      </el-row>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { mockActivityData } from '@/mock/admin'
import { Download } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

interface ActivityItem {
  componentName: string
  stars: number
  forks: number
  openIssues: number
  closedIssues: number
  commitsLastMonth: number
  contributors: number
  lastRelease: string
  trendData: { month: string; stars: number; commits: number }[]
  activityScore?: number
}

const activityData = ref<ActivityItem[]>(mockActivityData.map(d => ({
  ...d,
  activityScore: d.stars > 50000 ? 95 : d.stars > 10000 ? 85 : 65,
})))

const chartVisible = ref(false)
const selectedComponent = ref<ActivityItem | null>(null)
const starChartRef = ref<HTMLElement>()
const commitChartRef = ref<HTMLElement>()
const issueChartRef = ref<HTMLElement>()
const healthChartRef = ref<HTMLElement>()

function showChart(row: ActivityItem) {
  selectedComponent.value = row
  chartVisible.value = true
  nextTick(() => {
    renderCharts(row)
  })
}

function renderCharts(data: ActivityItem) {
  const months = data.trendData.map(d => d.month)
  const starValues = data.trendData.map(d => d.stars)
  const commitValues = data.trendData.map(d => d.commits)

  // Star 趋势图
  if (starChartRef.value) {
    const chart = echarts.init(starChartRef.value)
    chart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '8%', right: '5%', bottom: '10%', top: '10%' },
      xAxis: { type: 'category', data: months },
      yAxis: { type: 'value' },
      series: [{
        name: 'Stars',
        type: 'line',
        data: starValues,
        smooth: true,
        areaStyle: { color: 'rgba(64, 158, 255, 0.2)' },
        itemStyle: { color: '#409EFF' },
      }],
    })
  }

  // Commit 趋势图
  if (commitChartRef.value) {
    const chart = echarts.init(commitChartRef.value)
    chart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '8%', right: '5%', bottom: '10%', top: '10%' },
      xAxis: { type: 'category', data: months },
      yAxis: { type: 'value', name: '提交数' },
      series: [{
        name: 'Commits',
        type: 'bar',
        data: commitValues,
        itemStyle: { color: '#67C23A' },
      }],
    })
  }

  // Issue 趋势图（模拟数据）
  if (issueChartRef.value) {
    const chart = echarts.init(issueChartRef.value)
    chart.setOption({
      tooltip: { trigger: 'axis' },
      legend: { data: ['新增Issue', '关闭Issue'] },
      grid: { left: '8%', right: '5%', bottom: '10%', top: '15%' },
      xAxis: { type: 'category', data: months },
      yAxis: { type: 'value' },
      series: [
        { name: '新增Issue', type: 'line', data: [12, 18, 15, 22, 20, 16], smooth: true, itemStyle: { color: '#E6A23C' } },
        { name: '关闭Issue', type: 'line', data: [10, 15, 18, 20, 22, 18], smooth: true, itemStyle: { color: '#67C23A' } },
      ],
    })
  }

  // 综合健康度雷达图
  if (healthChartRef.value) {
    const chart = echarts.init(healthChartRef.value)
    chart.setOption({
      tooltip: {},
      radar: {
        indicator: [
          { name: 'Star增长', max: 100 },
          { name: 'Commit活跃度', max: 100 },
          { name: 'Issue响应', max: 100 },
          { name: '贡献者数量', max: 100 },
          { name: '发版频率', max: 100 },
        ],
        radius: '65%',
      },
      series: [{
        type: 'radar',
        data: [{
          value: [85, 78, 72, data.contributors > 100 ? 90 : 60, 68],
          name: data.componentName,
          areaStyle: { color: 'rgba(64, 158, 255, 0.3)' },
          itemStyle: { color: '#409EFF' },
        }],
      }],
    })
  }
}
</script>

<style scoped>
.page-header { margin-bottom: 8px; }
.page-header h3 { margin: 0 0 4px; }
.page-header h4 { margin: 0 0 4px; font-size: 14px; color: #606266; }
.desc { color: #909399; font-size: 13px; margin: 0; }
</style>
