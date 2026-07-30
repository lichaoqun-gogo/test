<template>
  <div class="interception-page">
    <div class="page-header">
      <h3>违规拦截统计</h3>
      <p class="page-desc">统计被可信仓库/准入策略拦截的黑名单及高危组件次数</p>
    </div>

    <!-- 关键指标 -->
    <el-row :gutter="16" class="metric-row">
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">本月拦截次数</span>
            <span class="metric-value danger">{{ summary.monthlyInterceptions }}</span>
            <span class="metric-trend" :class="summary.interceptionTrend >= 0 ? 'up' : 'down'">
              环比 {{ summary.interceptionTrend >= 0 ? '+' : '' }}{{ summary.interceptionTrend }}%
            </span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">累计拦截次数</span>
            <span class="metric-value primary">{{ summary.totalInterceptions }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">黑名单组件拦截</span>
            <span class="metric-value warning">{{ summary.blacklistInterceptions }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">高危组件拦截</span>
            <span class="metric-value">{{ summary.highRiskInterceptions }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 拦截趋势图 -->
    <el-row :gutter="16" class="chart-row">
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>拦截趋势（近12个月）</span>
              <el-radio-group v-model="trendType" size="small">
                <el-radio-button value="month">按月</el-radio-button>
                <el-radio-button value="quarter">按季度</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <v-chart class="chart" :option="interceptionTrendOption" autoresize />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>拦截类型分布</template>
          <v-chart class="chart pie-chart" :option="interceptionTypeOption" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <!-- 拦截明细 -->
    <el-card class="detail-table">
      <template #header>最近拦截记录</template>
      <el-table :data="interceptionRecords" stripe>
        <el-table-column prop="date" label="拦截时间" width="160" />
        <el-table-column prop="componentName" label="组件名称" width="150" />
        <el-table-column prop="version" label="版本" width="100" />
        <el-table-column prop="type" label="拦截类型" width="120">
          <template #default="{ row }">
            <el-tag :type="row.type === '黑名单' ? 'danger' : 'warning'" size="small">{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="policy" label="触发策略" min-width="150" />
        <el-table-column prop="requester" label="申请人" width="100" />
        <el-table-column prop="department" label="部门" width="120" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, BarChart, PieChart } from 'echarts/charts'
import {
  TooltipComponent,
  GridComponent,
  LegendComponent,
  TitleComponent,
} from 'echarts/components'
import type { EChartsOption } from 'echarts'

use([
  CanvasRenderer,
  LineChart,
  BarChart,
  PieChart,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  TitleComponent,
])

const trendType = ref<'month' | 'quarter'>('month')

const summary = {
  monthlyInterceptions: 28,
  totalInterceptions: 156,
  blacklistInterceptions: 89,
  highRiskInterceptions: 67,
  interceptionTrend: 12.5,
}

const monthlyData = [
  { month: '2025-07', blacklist: 8, highRisk: 5 },
  { month: '2025-08', blacklist: 10, highRisk: 7 },
  { month: '2025-09', blacklist: 6, highRisk: 4 },
  { month: '2025-10', blacklist: 12, highRisk: 8 },
  { month: '2025-11', blacklist: 9, highRisk: 6 },
  { month: '2025-12', blacklist: 15, highRisk: 10 },
  { month: '2026-01', blacklist: 11, highRisk: 7 },
  { month: '2026-02', blacklist: 8, highRisk: 5 },
  { month: '2026-03', blacklist: 14, highRisk: 9 },
  { month: '2026-04', blacklist: 10, highRisk: 8 },
  { month: '2026-05', blacklist: 13, highRisk: 11 },
  { month: '2026-06', blacklist: 16, highRisk: 12 },
]

const interceptionTypeOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c}次 ({d}%)' },
  legend: { bottom: '5%', left: 'center' },
  series: [
    {
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: false },
      emphasis: { label: { show: true, fontWeight: 'bold' } },
      data: [
        { value: summary.blacklistInterceptions, name: '黑名单拦截', itemStyle: { color: '#F56C6C' } },
        { value: summary.highRiskInterceptions, name: '高危组件拦截', itemStyle: { color: '#E6A23C' } },
        { value: 23, name: '许可证违规', itemStyle: { color: '#409EFF' } },
        { value: 21, name: '版本不合规', itemStyle: { color: '#67C23A' } },
      ],
    },
  ],
}))

const interceptionTrendOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'axis' },
  legend: { data: ['黑名单拦截', '高危组件拦截'], top: 5 },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: monthlyData.map(d => d.month),
  },
  yAxis: { type: 'value' },
  series: [
    {
      name: '黑名单拦截',
      type: 'line',
      data: monthlyData.map(d => d.blacklist),
      smooth: true,
      itemStyle: { color: '#F56C6C' },
      areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: '#F56C6C40' }, { offset: 1, color: '#F56C6C05' }] } },
    },
    {
      name: '高危组件拦截',
      type: 'line',
      data: monthlyData.map(d => d.highRisk),
      smooth: true,
      itemStyle: { color: '#E6A23C' },
      areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: '#E6A23C40' }, { offset: 1, color: '#E6A23C05' }] } },
    },
  ],
}))

const interceptionRecords = [
  { date: '2026-06-12 14:23', componentName: 'event-stream', version: '3.3.6', type: '黑名单', policy: '黑名单组件准入拦截', requester: '赵六', department: '数据智能部' },
  { date: '2026-06-11 09:45', componentName: 'log4js', version: '5.3.0', type: '高危', policy: '高危漏洞组件拦截', requester: '钱七', department: '产品研发部' },
  { date: '2026-06-10 16:30', componentName: 'colors', version: '1.4.0', type: '黑名单', policy: '黑名单组件准入拦截', requester: '孙八', department: '基础架构部' },
  { date: '2026-06-09 11:15', componentName: 'ua-parser-js', version: '0.7.30', type: '黑名单', policy: '黑名单组件准入拦截', requester: '周九', department: '技术研发部' },
  { date: '2026-06-08 08:50', componentName: 'node-ipc', version: '11.0.1', type: '高危', policy: '供应链安全风险拦截', requester: '吴十', department: '数据智能部' },
  { date: '2026-06-07 15:20', componentName: 'faker', version: '6.6.0', type: '黑名单', policy: '黑名单组件准入拦截', requester: '郑十一', department: '产品研发部' },
  { date: '2026-06-06 10:05', componentName: 'express', version: '4.17.0', type: '高危', policy: '已知CVE漏洞拦截', requester: '王十二', department: '技术研发部' },
]
</script>

<style scoped>
.interception-page { padding: 4px; }
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
.metric-trend { font-size: 12px; }
.metric-trend.up { color: #F56C6C; }
.metric-trend.down { color: #67C23A; }
.chart-row { margin-bottom: 16px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.chart { height: 350px; }
.pie-chart { height: 350px; }
.detail-table { margin-top: 8px; }
</style>
