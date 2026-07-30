<template>
  <div class="health-distribution-page">
    <div class="page-header">
      <h3>资产健康度分布</h3>
      <p class="page-desc">图表展示全量资产中正常组件占比、停更/无维护组件占比</p>
    </div>

    <!-- 关键指标 -->
    <el-row :gutter="16" class="metric-row">
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">资产总数</span>
            <span class="metric-value primary">{{ healthData.totalAssets }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">健康资产占比</span>
            <span class="metric-value success">{{ healthData.healthyRate }}%</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">停更组件数</span>
            <span class="metric-value warning">{{ healthData.stoppedCount }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">无维护组件数</span>
            <span class="metric-value danger">{{ healthData.unmaintainedCount }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 主图表区域 -->
    <el-row :gutter="16" class="chart-row">
      <el-col :span="10">
        <el-card>
          <template #header>资产健康度总览</template>
          <v-chart class="chart pie-chart" :option="healthPieOption" autoresize />
        </el-card>
      </el-col>
      <el-col :span="14">
        <el-card>
          <template #header>健康度趋势（近6个月）</template>
          <v-chart class="chart" :option="healthTrendOption" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <!-- 分类分布 -->
    <el-row :gutter="16" class="chart-row">
      <el-col :span="12">
        <el-card>
          <template #header>按资产类型分布</template>
          <v-chart class="chart" :option="typeDistributionOption" autoresize />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>按维护状态分布</template>
          <v-chart class="chart" :option="maintenanceOption" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <!-- 资产健康明细 -->
    <el-card class="detail-table">
      <template #header>资产健康状态明细</template>
      <el-table :data="assetHealthDetails" stripe>
        <el-table-column prop="name" label="资产名称" width="140" />
        <el-table-column prop="version" label="当前版本" width="100" />
        <el-table-column prop="type" label="资产类型" width="100" />
        <el-table-column prop="status" label="健康状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastRelease" label="最后发版" width="120" />
        <el-table-column prop="stars" label="Stars" width="90" />
        <el-table-column prop="commitsLastMonth" label="月提交" width="90" />
        <el-table-column prop="contributors" label="贡献者" width="90" />
        <el-table-column prop="riskDesc" label="风险提示" min-width="160" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { PieChart, LineChart } from 'echarts/charts'
import {
  TooltipComponent,
  GridComponent,
  LegendComponent,
} from 'echarts/components'
import type { EChartsOption } from 'echarts'

use([CanvasRenderer, PieChart, LineChart, TooltipComponent, GridComponent, LegendComponent])

const healthData = {
  totalAssets: 200,
  healthyRate: 78,
  stoppedCount: 28,
  unmaintainedCount: 16,
}

const healthPieOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c}个 ({d}%)' },
  legend: { orient: 'vertical', left: 'left', top: 'center' },
  series: [{
    type: 'pie',
    radius: ['50%', '75%'],
    center: ['60%', '50%'],
    avoidLabelOverlap: false,
    itemStyle: { borderRadius: 8, borderColor: '#fff', borderWidth: 2 },
    label: { show: false },
    emphasis: { label: { show: true, fontWeight: 'bold' } },
    data: [
      { value: 156, name: '正常', itemStyle: { color: '#67C23A' } },
      { value: 28, name: '停更', itemStyle: { color: '#E6A23C' } },
      { value: 16, name: '无维护', itemStyle: { color: '#F56C6C' } },
    ],
  }],
}))

const healthTrendData = [
  { month: '2026-01', healthy: 142, stopped: 30, unmaintained: 20 },
  { month: '2026-02', healthy: 145, stopped: 29, unmaintained: 19 },
  { month: '2026-03', healthy: 148, stopped: 28, unmaintained: 18 },
  { month: '2026-04', healthy: 150, stopped: 28, unmaintained: 17 },
  { month: '2026-05', healthy: 153, stopped: 28, unmaintained: 16 },
  { month: '2026-06', healthy: 156, stopped: 28, unmaintained: 16 },
]

const healthTrendOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'axis' },
  legend: { data: ['正常', '停更', '无维护'], top: 5 },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: healthTrendData.map(d => d.month) },
  yAxis: { type: 'value' },
  series: [
    { name: '正常', type: 'line', data: healthTrendData.map(d => d.healthy), smooth: true, itemStyle: { color: '#67C23A' }, areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: '#67C23A40' }, { offset: 1, color: '#67C23A05' }] } } },
    { name: '停更', type: 'line', data: healthTrendData.map(d => d.stopped), smooth: true, itemStyle: { color: '#E6A23C' }, areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: '#E6A23C40' }, { offset: 1, color: '#E6A23C05' }] } } },
    { name: '无维护', type: 'line', data: healthTrendData.map(d => d.unmaintained), smooth: true, itemStyle: { color: '#F56C6C' }, areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: '#F56C6C40' }, { offset: 1, color: '#F56C6C05' }] } } },
  ],
}))

const typeDistributionOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c}个 ({d}%)' },
  legend: { bottom: '5%', left: 'center' },
  series: [{
    type: 'pie',
    radius: ['40%', '70%'],
    avoidLabelOverlap: false,
    itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
    label: { show: false },
    emphasis: { label: { show: true, fontWeight: 'bold' } },
    data: [
      { value: 56, name: '中间件', itemStyle: { color: '#409EFF' } },
      { value: 42, name: '数据库', itemStyle: { color: '#67C23A' } },
      { value: 68, name: '依赖包', itemStyle: { color: '#E6A23C' } },
      { value: 22, name: '工具库', itemStyle: { color: '#909399' } },
      { value: 12, name: '框架', itemStyle: { color: '#F56C6C' } },
    ],
  }],
}))

const maintenanceOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c}个 ({d}%)' },
  legend: { bottom: '5%', left: 'center' },
  series: [{
    type: 'pie',
    radius: ['40%', '70%'],
    avoidLabelOverlap: false,
    itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
    label: { show: false },
    emphasis: { label: { show: true, fontWeight: 'bold' } },
    data: [
      { value: 98, name: '活跃维护', itemStyle: { color: '#67C23A' } },
      { value: 58, name: '正常维护', itemStyle: { color: '#409EFF' } },
      { value: 28, name: '维护减缓', itemStyle: { color: '#E6A23C' } },
      { value: 16, name: '已停更', itemStyle: { color: '#F56C6C' } },
    ],
  }],
}))

const assetHealthDetails = [
  { name: 'Nginx', version: '1.25.4', type: '中间件', status: '正常', lastRelease: '2026-06-10', stars: '18.5K', commitsLastMonth: 45, contributors: 89, riskDesc: '-' },
  { name: 'Redis', version: '7.2.4', type: '数据库', status: '正常', lastRelease: '2026-06-08', stars: '65.2K', commitsLastMonth: 128, contributors: 456, riskDesc: '-' },
  { name: 'Bower', version: '1.8.14', type: '工具库', status: '停更', lastRelease: '2022-01-10', stars: '2.1K', commitsLastMonth: 0, contributors: 0, riskDesc: '已停止维护，建议迁移' },
  { name: 'Grunt', version: '1.6.1', type: '工具库', status: '无维护', lastRelease: '2021-09-15', stars: '1.8K', commitsLastMonth: 0, contributors: 0, riskDesc: '社区不再维护，需替换' },
  { name: 'Spring Boot', version: '3.2.3', type: '框架', status: '正常', lastRelease: '2026-06-05', stars: '72K', commitsLastMonth: 89, contributors: 312, riskDesc: '-' },
  { name: 'lodash', version: '4.17.21', type: '依赖包', status: '正常', lastRelease: '2026-05-28', stars: '58K', commitsLastMonth: 12, contributors: 89, riskDesc: '长期未大版本更新' },
]

function getStatusType(status: string): 'success' | 'warning' | 'danger' | 'info' {
  const map: Record<string, 'success' | 'warning' | 'danger' | 'info'> = {
    '正常': 'success',
    '停更': 'warning',
    '无维护': 'danger',
  }
  return map[status] || 'info'
}
</script>

<style scoped>
.health-distribution-page { padding: 4px; }
.page-header { margin-bottom: 16px; }
.page-header h3 { margin: 0 0 4px; font-size: 18px; color: #303133; }
.page-desc { margin: 0; font-size: 13px; color: #909399; }
.metric-row { margin-bottom: 16px; }
.metric-card .metric-content { display: flex; flex-direction: column; gap: 8px; }
.metric-label { font-size: 13px; color: #909399; }
.metric-value { font-size: 28px; font-weight: 700; }
.metric-value.primary { color: #409EFF; }
.metric-value.success { color: #67C23A; }
.metric-value.warning { color: #E6A23C; }
.metric-value.danger { color: #F56C6C; }
.chart-row { margin-bottom: 16px; }
.chart { height: 320px; }
.pie-chart { height: 320px; }
.detail-table { margin-top: 8px; }
</style>
