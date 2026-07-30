<template>
  <div class="workload-page">
    <div class="page-header">
      <h3>治理工作量统计</h3>
      <p class="page-desc">统计各业务部门处理的引入申请、升级审批数量</p>
    </div>

    <!-- 关键指标 -->
    <el-row :gutter="16" class="metric-row">
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">本月引入申请</span>
            <span class="metric-value primary">{{ summary.monthlyIntroductions }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">本月升级审批</span>
            <span class="metric-value warning">{{ summary.monthlyUpgrades }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">待处理工单</span>
            <span class="metric-value danger">{{ summary.pendingTickets }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">本月审批完成率</span>
            <span class="metric-value success">{{ summary.completionRate }}%</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 部门工作量对比 -->
    <el-row :gutter="16" class="chart-row">
      <el-col :span="12">
        <el-card>
          <template #header>各部门工作量对比</template>
          <v-chart class="chart" :option="departmentWorkloadOption" autoresize />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>工作量占比分布</template>
          <v-chart class="chart" :option="workloadPieOption" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <!-- 工作量趋势 -->
    <el-card class="chart-card">
      <template #header>工作量趋势（近6个月）</template>
      <v-chart class="chart" :option="workloadTrendOption" autoresize />
    </el-card>

    <!-- 部门明细表 -->
    <el-card class="detail-table">
      <template #header>部门治理工作量明细</template>
      <el-table :data="departmentStats" stripe>
        <el-table-column prop="department" label="部门" width="140" />
        <el-table-column prop="introduce" label="引入申请" width="100" />
        <el-table-column prop="upgrade" label="升级审批" width="100" />
        <el-table-column prop="withdrawal" label="退出申请" width="100" />
        <el-table-column prop="total" label="合计" width="90">
          <template #default="{ row }">
            <el-tag type="primary" size="small">{{ row.total }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="工作量占比" width="120">
          <template #default="{ row }">
            <el-progress :percentage="row.percentage" :stroke-width="8" :color="getProgressColor(row.percentage)" />
          </template>
        </el-table-column>
        <el-table-column label="处理效率" width="100">
          <template #default="{ row }">
            <el-tag :type="row.efficiency >= 90 ? 'success' : row.efficiency >= 80 ? 'warning' : 'danger'" size="small">
              {{ row.efficiency }}%
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="环比趋势" width="100">
          <template #default="{ row }">
            <span :class="row.trend >= 0 ? 'trend-up' : 'trend-down'">
              {{ row.trend >= 0 ? '↑' : '↓' }} {{ Math.abs(row.trend) }}%
            </span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart, PieChart, LineChart } from 'echarts/charts'
import {
  TooltipComponent,
  GridComponent,
  LegendComponent,
} from 'echarts/components'
import type { EChartsOption } from 'echarts'

use([CanvasRenderer, BarChart, PieChart, LineChart, TooltipComponent, GridComponent, LegendComponent])

const summary = {
  monthlyIntroductions: 45,
  monthlyUpgrades: 18,
  pendingTickets: 12,
  completionRate: 88,
}

const departmentStats = [
  { department: '技术研发部', introduce: 45, upgrade: 12, withdrawal: 3, total: 60, percentage: 42, efficiency: 92, trend: 5.2 },
  { department: '数据智能部', introduce: 18, upgrade: 5, withdrawal: 1, total: 24, percentage: 17, efficiency: 85, trend: -3.1 },
  { department: '产品技术部', introduce: 22, upgrade: 8, withdrawal: 2, total: 32, percentage: 22, efficiency: 88, trend: 8.4 },
  { department: '基础架构部', introduce: 30, upgrade: 15, withdrawal: 4, total: 49, percentage: 34, efficiency: 95, trend: 12.5 },
]

const workloadMonths = ['2026-01', '2026-02', '2026-03', '2026-04', '2026-05', '2026-06']
const introduceTrend = [32, 38, 41, 35, 42, 45]
const upgradeTrend = [8, 10, 12, 9, 15, 18]
const withdrawalTrend = [2, 3, 1, 4, 2, 3]

const departmentWorkloadOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'axis' },
  legend: { data: ['引入申请', '升级审批', '退出申请'], top: 5 },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: departmentStats.map(d => d.department), axisLabel: { interval: 0, rotate: 15 } },
  yAxis: { type: 'value', name: '工单数' },
  series: [
    { name: '引入申请', type: 'bar', stack: 'total', data: departmentStats.map(d => d.introduce), itemStyle: { color: '#409EFF' }, barWidth: '50%' },
    { name: '升级审批', type: 'bar', stack: 'total', data: departmentStats.map(d => d.upgrade), itemStyle: { color: '#E6A23C' }, barWidth: '50%' },
    { name: '退出申请', type: 'bar', stack: 'total', data: departmentStats.map(d => d.withdrawal), itemStyle: { color: '#909399' }, barWidth: '50%' },
  ],
}))

const workloadPieOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c}单 ({d}%)' },
  legend: { bottom: '5%', left: 'center' },
  series: [{
    type: 'pie',
    radius: ['40%', '70%'],
    avoidLabelOverlap: false,
    itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
    label: { show: false },
    emphasis: { label: { show: true, fontWeight: 'bold' } },
    data: [
      { value: 60, name: '技术研发部', itemStyle: { color: '#409EFF' } },
      { value: 24, name: '数据智能部', itemStyle: { color: '#67C23A' } },
      { value: 32, name: '产品技术部', itemStyle: { color: '#E6A23C' } },
      { value: 49, name: '基础架构部', itemStyle: { color: '#F56C6C' } },
    ],
  }],
}))

const workloadTrendOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'axis' },
  legend: { data: ['引入申请', '升级审批', '退出申请'], top: 5 },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: workloadMonths },
  yAxis: { type: 'value', name: '工单数' },
  series: [
    { name: '引入申请', type: 'line', data: introduceTrend, smooth: true, itemStyle: { color: '#409EFF' }, areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: '#409EFF40' }, { offset: 1, color: '#409EFF05' }] } } },
    { name: '升级审批', type: 'line', data: upgradeTrend, smooth: true, itemStyle: { color: '#E6A23C' }, areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: '#E6A23C40' }, { offset: 1, color: '#E6A23C05' }] } } },
    { name: '退出申请', type: 'line', data: withdrawalTrend, smooth: true, itemStyle: { color: '#909399' } },
  ],
}))

function getProgressColor(rate: number): string {
  if (rate >= 35) return '#F56C6C'
  if (rate >= 25) return '#E6A23C'
  return '#67C23A'
}
</script>

<style scoped>
.workload-page { padding: 4px; }
.page-header { margin-bottom: 16px; }
.page-header h3 { margin: 0 0 4px; font-size: 18px; color: #303133; }
.page-desc { margin: 0; font-size: 13px; color: #909399; }
.metric-row { margin-bottom: 16px; }
.metric-card .metric-content { display: flex; flex-direction: column; gap: 8px; }
.metric-label { font-size: 13px; color: #909399; }
.metric-value { font-size: 28px; font-weight: 700; }
.metric-value.primary { color: #409EFF; }
.metric-value.warning { color: #E6A23C; }
.metric-value.danger { color: #F56C6C; }
.metric-value.success { color: #67C23A; }
.chart-row { margin-bottom: 16px; }
.chart { height: 320px; }
.chart-card { margin-bottom: 16px; }
.detail-table { margin-top: 8px; }
.trend-up { color: #F56C6C; }
.trend-down { color: #67C23A; }
</style>
