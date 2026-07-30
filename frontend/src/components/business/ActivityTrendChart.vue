<template>
  <div class="activity-trend-chart">
    <div class="chart-header">
      <span class="chart-title">活跃度趋势</span>
      <div class="metric-switcher">
        <el-radio-group v-model="currentMetric" size="small">
          <el-radio-button
            v-for="metric in availableMetrics"
            :key="metric"
            :value="metric"
          >
            {{ metricLabels[metric] }}
          </el-radio-button>
        </el-radio-group>
      </div>
    </div>
    <v-chart
      class="chart-container"
      :option="chartOption"
      autoresize
      :loading="data.length === 0"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
} from 'echarts/components'
import type { EChartsOption } from 'echarts'

use([
  CanvasRenderer,
  LineChart,
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
])

type MetricType = 'star' | 'issue' | 'commit'

interface TrendDataPoint {
  date: string
  star?: number
  issue?: number
  commit?: number
}

const props = withDefaults(defineProps<{
  data: TrendDataPoint[]
  metrics?: MetricType[]
}>(), {
  metrics: () => ['star', 'issue', 'commit'],
})

const metricLabels: Record<MetricType, string> = {
  star: 'Star',
  issue: 'Issue',
  commit: 'Commit',
}

const availableMetrics = computed(() => {
  return props.metrics.filter((m) => metricLabels[m])
})

const currentMetric = ref<MetricType>(availableMetrics.value[0] || 'star')

watch(
  () => props.metrics,
  (newMetrics) => {
    if (newMetrics && newMetrics.length > 0 && !newMetrics.includes(currentMetric.value)) {
      currentMetric.value = newMetrics[0]
    }
  },
)

const chartOption = computed<EChartsOption>(() => {
  const dates = props.data.map((item) => item.date)
  const values = props.data.map((item) => item[currentMetric.value] ?? 0)

  const metricColors: Record<MetricType, string> = {
    star: '#F5A623',
    issue: '#F56C6C',
    commit: '#409EFF',
  }

  return {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        if (!params || params.length === 0) return ''
        const param = params[0]
        return `${param.name}<br/>${metricLabels[currentMetric.value]}: ${param.value}`
      },
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: dates,
      boundaryGap: false,
      axisLine: {
        lineStyle: {
          color: '#DCDFE6',
        },
      },
      axisLabel: {
        color: '#606266',
      },
    },
    yAxis: {
      type: 'value',
      splitLine: {
        lineStyle: {
          type: 'dashed',
          color: '#EBEEF5',
        },
      },
      axisLabel: {
        color: '#606266',
      },
    },
    series: [
      {
        name: metricLabels[currentMetric.value],
        type: 'line',
        data: values,
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: {
          width: 3,
          color: metricColors[currentMetric.value],
        },
        itemStyle: {
          color: metricColors[currentMetric.value],
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              {
                offset: 0,
                color: metricColors[currentMetric.value] + '40',
              },
              {
                offset: 1,
                color: metricColors[currentMetric.value] + '05',
              },
            ],
          },
        },
      },
    ],
  }
})
</script>

<style scoped>
.activity-trend-chart {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.chart-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.metric-switcher {
  display: flex;
  align-items: center;
}

.chart-container {
  height: 300px;
  width: 100%;
}
</style>
