<template>
  <div class="repair-efficiency-page">
    <div class="page-header">
      <h3>缺陷修复效能</h3>
      <p class="page-desc">按部门统计缺陷平均修复时长（MTTR）和逾期未修复数</p>
    </div>

    <!-- 关键指标 -->
    <el-row :gutter="16" class="metric-row">
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">公司平均 MTTR</span>
            <span class="metric-value primary">{{ overallMttr }} 天</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">逾期未修复总数</span>
            <span class="metric-value danger">{{ totalOverdue }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">本月新增缺陷</span>
            <span class="metric-value warning">{{ monthlyNewDefects }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">本月修复完成率</span>
            <span class="metric-value success">{{ repairCompletionRate }}%</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- MTTR 柱状图 -->
    <el-row :gutter="16" class="chart-row">
      <el-col :span="12">
        <el-card>
          <template #header>各部门 MTTR 对比（天）</template>
          <v-chart class="chart" :option="mttrBarOption" autoresize />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>各部门逾期未修复数</template>
          <v-chart class="chart" :option="overdueBarOption" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <!-- 部门效能明细 -->
    <el-row :gutter="16" class="chart-row">
      <el-col :span="12">
        <el-card>
          <template #header>缺陷修复时效分布</template>
          <v-chart class="chart" :option="timeDistributionOption" autoresize />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>缺陷严重等级分布</template>
          <v-chart class="chart" :option="severityOption" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <!-- 部门明细表 -->
    <el-card class="detail-table">
      <template #header>部门修复效能明细</template>
      <el-table :data="departmentStats" stripe>
        <el-table-column prop="department" label="部门" width="140" />
        <el-table-column label="平均 MTTR" width="120">
          <template #default="{ row }">
            <el-tag :type="getMttrType(row.mttr)" size="small">{{ row.mttr }} 天</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="totalDefects" label="缺陷总数" width="100" />
        <el-table-column prop="fixedDefects" label="已修复" width="90" />
        <el-table-column prop="overdue" label="逾期未修复" width="110">
          <template #default="{ row }">
            <el-tag v-if="row.overdue > 0" type="danger" size="small">{{ row.overdue }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="修复率" width="100">
          <template #default="{ row }">
            <el-progress :percentage="row.fixRate" :stroke-width="8" :color="getProgressColor(row.fixRate)" />
          </template>
        </el-table-column>
        <el-table-column label="趋势" width="100">
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
import { BarChart, PieChart } from 'echarts/charts'
import {
  TooltipComponent,
  GridComponent,
  LegendComponent,
} from 'echarts/components'
import type { EChartsOption } from 'echarts'

use([CanvasRenderer, BarChart, PieChart, TooltipComponent, GridComponent, LegendComponent])

const overallMttr = 3.9
const totalOverdue = 7
const monthlyNewDefects = 23
const repairCompletionRate = 86

const departmentStats = [
  { department: '技术研发部', mttr: 3.5, totalDefects: 45, fixedDefects: 40, overdue: 2, fixRate: 89, trend: -5.2 },
  { department: '数据智能部', mttr: 5.2, totalDefects: 32, fixedDefects: 25, overdue: 4, fixRate: 78, trend: 3.1 },
  { department: '产品技术部', mttr: 4.1, totalDefects: 28, fixedDefects: 25, overdue: 1, fixRate: 89, trend: -8.4 },
  { department: '基础架构部', mttr: 2.8, totalDefects: 18, fixedDefects: 18, overdue: 0, fixRate: 100, trend: -12.5 },
]

const mttrBarOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'axis', formatter: '{b}: {c} 天' },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: departmentStats.map(d => d.department), axisLabel: { interval: 0, rotate: 15 } },
  yAxis: { type: 'value', name: '天数' },
  series: [{
    type: 'bar',
    data: departmentStats.map(d => ({
      value: d.mttr,
      itemStyle: { color: d.mttr > 4 ? '#F56C6C' : d.mttr > 3 ? '#E6A23C' : '#67C23A' },
    })),
    barWidth: '50%',
    label: { show: true, position: 'top', formatter: '{c} 天' },
  }],
}))

const overdueBarOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'axis', formatter: '{b}: {c} 件' },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: departmentStats.map(d => d.department), axisLabel: { interval: 0, rotate: 15 } },
  yAxis: { type: 'value', name: '件数' },
  series: [{
    type: 'bar',
    data: departmentStats.map(d => ({
      value: d.overdue,
      itemStyle: { color: d.overdue > 0 ? '#F56C6C' : '#E8E8E8' },
    })),
    barWidth: '50%',
    label: { show: true, position: 'top' },
  }],
}))

const timeDistributionOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c}件 ({d}%)' },
  legend: { bottom: '5%', left: 'center' },
  series: [{
    type: 'pie',
    radius: ['40%', '70%'],
    avoidLabelOverlap: false,
    itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
    label: { show: false },
    emphasis: { label: { show: true, fontWeight: 'bold' } },
    data: [
      { value: 42, name: '1天内修复', itemStyle: { color: '#67C23A' } },
      { value: 35, name: '1-3天修复', itemStyle: { color: '#409EFF' } },
      { value: 28, name: '3-7天修复', itemStyle: { color: '#E6A23C' } },
      { value: 12, name: '7-14天修复', itemStyle: { color: '#F56C6C' } },
      { value: 7, name: '超过14天', itemStyle: { color: '#909399' } },
    ],
  }],
}))

const severityOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c}件 ({d}%)' },
  legend: { bottom: '5%', left: 'center' },
  series: [{
    type: 'pie',
    radius: ['40%', '70%'],
    avoidLabelOverlap: false,
    itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
    label: { show: false },
    emphasis: { label: { show: true, fontWeight: 'bold' } },
    data: [
      { value: 8, name: '紧急', itemStyle: { color: '#F5222D' } },
      { value: 25, name: '高', itemStyle: { color: '#F56C6C' } },
      { value: 42, name: '中', itemStyle: { color: '#E6A23C' } },
      { value: 49, name: '低', itemStyle: { color: '#67C23A' } },
    ],
  }],
}))

function getMttrType(mttr: number): 'success' | 'warning' | 'danger' {
  if (mttr <= 3) return 'success'
  if (mttr <= 4) return 'warning'
  return 'danger'
}

function getProgressColor(rate: number): string {
  if (rate >= 90) return '#67C23A'
  if (rate >= 80) return '#E6A23C'
  return '#F56C6C'
}
</script>

<style scoped>
.repair-efficiency-page { padding: 4px; }
.page-header { margin-bottom: 16px; }
.page-header h3 { margin: 0 0 4px; font-size: 18px; color: #303133; }
.page-desc { margin: 0; font-size: 13px; color: #909399; }
.metric-row { margin-bottom: 16px; }
.metric-card .metric-content { display: flex; flex-direction: column; gap: 8px; }
.metric-label { font-size: 13px; color: #909399; }
.metric-value { font-size: 28px; font-weight: 700; }
.metric-value.primary { color: #409EFF; }
.metric-value.danger { color: #F56C6C; }
.metric-value.warning { color: #E6A23C; }
.metric-value.success { color: #67C23A; }
.chart-row { margin-bottom: 16px; }
.chart { height: 320px; }
.detail-table { margin-top: 8px; }
.trend-up { color: #F56C6C; }
.trend-down { color: #67C23A; }
</style>
