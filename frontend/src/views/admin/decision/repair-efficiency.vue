<template>
  <div class="repair-efficiency-page" v-loading="loading">
    <div class="page-header">
      <h3>缺陷修复效能</h3>
      <p class="page-desc">按部门统计外部缺陷平均修复时长（MTTR）和逾期未修复数</p>
    </div>

    <!-- 筛选条 -->
    <el-card class="filter-bar" shadow="never">
      <div class="filter-row">
        <span class="filter-label">时间范围</span>
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          value-format="YYYY-MM-DD"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          :shortcuts="shortcuts"
          clearable
          style="width: 260px"
        />
        <span class="filter-label">部门</span>
        <el-select v-model="deptFilter" placeholder="全部部门" clearable style="width: 200px">
          <el-option v-for="d in deptOptions" :key="d" :label="d" :value="d" />
        </el-select>
        <el-button type="primary" :icon="Refresh" circle @click="loadData" />
      </div>
      <p class="filter-tip">逾期未闭环为当前存量，其余指标按所选时间范围统计；MTTR 从派发到完成。</p>
    </el-card>

    <!-- 关键指标 -->
    <el-row :gutter="16" class="metric-row">
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">公司平均 MTTR</span>
            <span class="metric-value primary">{{ kpi.mttr ?? '-' }} 天</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">逾期未闭环工单</span>
            <span class="metric-value danger">{{ kpi.overdue ?? 0 }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">本期新增缺陷</span>
            <span class="metric-value warning">{{ kpi.newDefects ?? 0 }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">本期闭环率</span>
            <span class="metric-value success">{{ kpi.closeRate ?? 0 }}%</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 分布图：严重等级 / 处理情况 -->
    <el-row :gutter="16" class="chart-row">
      <el-col :span="12">
        <el-card>
          <template #header>缺陷严重等级分布</template>
          <v-chart class="chart" :option="severityOption" autoresize />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>缺陷处理情况分布</template>
          <v-chart class="chart" :option="statusOption" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <!-- 部门效能明细 -->
    <el-card class="detail-table">
      <template #header>部门修复效能明细</template>
      <el-table :data="deptStats" stripe v-if="deptStats.length">
        <el-table-column prop="deptName" label="部门" width="150" />
        <el-table-column prop="total" label="处置工单数" width="110" />
        <el-table-column prop="closed" label="已闭环" width="90" />
        <el-table-column prop="processing" label="进行中" width="90" />
        <el-table-column label="逾期未闭环" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.overdue > 0" type="danger" size="small">{{ row.overdue }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="平均 MTTR" width="120">
          <template #default="{ row }">
            <el-tag :type="getMttrType(row.mttr)" size="small">{{ row.mttr ? row.mttr + ' 天' : '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="闭环率" min-width="160">
          <template #default="{ row }">
            <el-progress :percentage="row.closeRate" :stroke-width="8" :color="getProgressColor(row.closeRate)" />
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-else description="暂无处置工单数据" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { PieChart } from 'echarts/charts'
import {
  TooltipComponent,
  LegendComponent,
} from 'echarts/components'
import type { EChartsOption } from 'echarts'
import { riskApi } from '@/api'

use([CanvasRenderer, PieChart, TooltipComponent, LegendComponent])

interface DefectStatsOverview {
  kpi: { mttr: number; overdue: number; newDefects: number; completed: number; closeRate: number }
  severity: { level: string; count: number }[]
  status: { key: string; name: string; count: number }[]
  dept: {
    deptName: string
    total: number
    closed: number
    processing: number
    overdue: number
    mttr: number
    closeRate: number
  }[]
  depts: string[]
}

const loading = ref(false)
const overview = ref<DefectStatsOverview | null>(null)
const dateRange = ref<[string, string]>([
  dayjs().subtract(29, 'day').format('YYYY-MM-DD'),
  dayjs().format('YYYY-MM-DD'),
])
const deptFilter = ref('')

const shortcuts = [
  {
    text: '近7天',
    value: () => {
      const end = dayjs()
      return [end.subtract(6, 'day').startOf('day').toDate(), end.endOf('day').toDate()]
    },
  },
  {
    text: '近30天',
    value: () => {
      const end = dayjs()
      return [end.subtract(29, 'day').startOf('day').toDate(), end.endOf('day').toDate()]
    },
  },
  {
    text: '本月',
    value: () => {
      return [dayjs().startOf('month').toDate(), dayjs().endOf('month').toDate()]
    },
  },
  {
    text: '本季度',
    value: () => {
      const now = dayjs()
      return [now.startOf('quarter').toDate(), now.endOf('quarter').toDate()]
    },
  },
]

async function loadData() {
  loading.value = true
  try {
    const params: Record<string, string> = {}
    if (dateRange.value && dateRange.value.length === 2) {
      params.start = dateRange.value[0]
      params.end = dateRange.value[1]
    }
    if (deptFilter.value) {
      params.dept = deptFilter.value
    }
    overview.value = await riskApi.getDefectStatsOverview(params)
  } catch (e) {
    console.error('加载缺陷修复效能数据失败:', e)
  } finally {
    loading.value = false
  }
}

const kpi = computed(() => overview.value?.kpi || { mttr: 0, overdue: 0, newDefects: 0, completed: 0, closeRate: 0 })
const deptOptions = computed(() => overview.value?.depts || [])
const deptStats = computed(() => overview.value?.dept || [])

// 严重等级分布（外部缺陷等级：高危/中危/低危/提示/未知）
const severityColors: Record<string, string> = {
  '高危': '#F5222D',
  '中危': '#F56C6C',
  '低危': '#E6A23C',
  '提示': '#67C23A',
  '未知': '#909399',
}
const severityOption = computed<EChartsOption>(() => {
  const data = (overview.value?.severity || []).map(s => ({
    name: s.level,
    value: s.count,
    itemStyle: { color: severityColors[s.level] || '#909399' },
  }))
  return {
    tooltip: { trigger: 'item', formatter: '{b}: {c}件 ({d}%)' },
    legend: { bottom: '0', left: 'center' },
    series: [{
      type: 'pie',
      radius: ['40%', '68%'],
      center: ['50%', '46%'],
      avoidLabelOverlap: false,
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: false },
      emphasis: { label: { show: true, fontWeight: 'bold' } },
      data,
    }],
  }
})

// 处理情况分布
const statusColors: Record<string, string> = {
  '进行中': '#409EFF',
  '逾期未闭环': '#F56C6C',
  '已完成处置': '#67C23A',
}
const statusOption = computed<EChartsOption>(() => {
  const data = (overview.value?.status || []).map(s => ({
    name: s.name,
    value: s.count,
    itemStyle: { color: statusColors[s.name] || '#909399' },
  }))
  return {
    tooltip: { trigger: 'item', formatter: '{b}: {c}件 ({d}%)' },
    legend: { bottom: '0', left: 'center' },
    series: [{
      type: 'pie',
      radius: ['40%', '68%'],
      center: ['50%', '46%'],
      avoidLabelOverlap: false,
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: false },
      emphasis: { label: { show: true, fontWeight: 'bold' } },
      data,
    }],
  }
})

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

watch([dateRange, deptFilter], loadData)

onMounted(loadData)
</script>

<style scoped>
.repair-efficiency-page { padding: 4px; }
.page-header { margin-bottom: 16px; }
.page-header h3 { margin: 0 0 4px; font-size: 18px; color: #303133; }
.page-desc { margin: 0; font-size: 13px; color: #909399; }
.filter-bar { margin-bottom: 16px; }
.filter-row { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.filter-label { font-size: 13px; color: #606266; }
.filter-tip { margin: 10px 0 0; font-size: 12px; color: #c0c4cc; }
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
</style>
