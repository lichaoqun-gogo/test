<template>
  <el-dialog v-model="visible" :title="`软件详情 - ${projectData?._name || ''}`" width="900px" @close="handleClose">
    <el-tabs v-model="activeTab">
      <!-- Tab1: 许可证信息 -->
      <el-tab-pane label="许可证信息" name="license">
        <div class="license-charts">
          <div class="chart-container">
            <h4>许可证分级总览</h4>
            <div ref="groupChartRef" class="chart"></div>
          </div>
          <div class="chart-container">
            <h4>许可证信息总览</h4>
            <div ref="licenseChartRef" class="chart"></div>
          </div>
        </div>
      </el-tab-pane>

      <!-- Tab2: 依赖信息 -->
      <el-tab-pane label="依赖信息" name="dependency">
        <div class="dependency-filter">
          <el-input v-model="filterKeyword" placeholder="搜索依赖名称/组织" clearable style="width: 200px" @clear="loadDependencies" @keyup.enter="loadDependencies">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
          <el-input v-model="filterVersion" placeholder="版本号" clearable style="width: 150px; margin-left: 12px" @clear="loadDependencies" @keyup.enter="loadDependencies" />
          <el-button type="primary" @click="loadDependencies" style="margin-left: 12px">
            <el-icon><Search /></el-icon> 查询
          </el-button>
        </div>

        <el-table :data="dependencies" stripe style="margin-top: 16px" v-loading="loadingDeps">
          <el-table-column prop="_group" label="组织" min-width="180" />
          <el-table-column prop="_name" label="依赖名称" min-width="160" />
          <el-table-column prop="_version" label="版本" width="130" />
          <el-table-column prop="CLASSIFIER" label="类型" width="120" />
          <el-table-column prop="LICENSE" label="许可证" width="130" />
          <el-table-column label="风险值" width="100">
            <template #default="{ row }">
              <el-tag :type="row.LAST_RISKSCORE > 0.1 ? 'danger' : 'success'" size="small">
                {{ row.LAST_RISKSCORE }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>

        <div class="total-info">共 {{ dependencies.length }} 个依赖</div>
      </el-tab-pane>
    </el-tabs>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, computed } from 'vue'
import { Search } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { getProjectComponents, getProjectLicenseStats } from '@/api'

const props = defineProps<{
  modelValue: boolean
  projectData: any
}>()

const emit = defineEmits(['update:modelValue'])

const activeTab = ref('license')
const loadingDeps = ref(false)
const dependencies = ref<any[]>([])
const filterKeyword = ref('')
const filterVersion = ref('')

const groupChartRef = ref<HTMLElement>()
const licenseChartRef = ref<HTMLElement>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

watch(() => props.modelValue, (val) => {
  if (val && props.projectData?.id) {
    activeTab.value = 'license'
    loadLicenseStats()
    loadDependencies()
  }
})

async function loadDependencies() {
  if (!props.projectData?.id) return
  
  loadingDeps.value = true
  try {
    const params: any = {}
    if (filterKeyword.value) params.keyword = filterKeyword.value
    if (filterVersion.value) params.version = filterVersion.value
    
    const res = await getProjectComponents(props.projectData.id, params)
    console.log('依赖数据:', res)
    dependencies.value = res || []
  } catch (error) {
    console.error('加载依赖失败:', error)
    dependencies.value = []
  } finally {
    loadingDeps.value = false
  }
}

async function loadLicenseStats() {
  if (!props.projectData?.id) return
  
  try {
    const res = await getProjectLicenseStats(props.projectData.id)
    console.log('许可证统计数据:', res)
    
    await nextTick()
    
    // 许可证分级饼图
    if (groupChartRef.value && res.group_stats) {
      const groupChart = echarts.init(groupChartRef.value)
      groupChart.setOption({
        tooltip: { trigger: 'item' },
        legend: {
          type: 'scroll',
          orient: 'vertical',
          right: '5%',
          top: 'center',
          textStyle: { fontSize: 12 }
        },
        series: [{
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['35%', '50%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: { show: false, position: 'center' },
          emphasis: {
            label: {
              show: true,
              fontSize: 20,
              fontWeight: 'bold'
            }
          },
          data: res.group_stats.map((item: any) => ({
            name: item.group_name || '未分类',
            value: item.count
          }))
        }]
      })
    }
    
    // 许可证信息饼图
    if (licenseChartRef.value && res.license_stats) {
      const licenseChart = echarts.init(licenseChartRef.value)
      licenseChart.setOption({
        tooltip: { trigger: 'item' },
        legend: {
          type: 'scroll',
          orient: 'vertical',
          right: '5%',
          top: 'center',
          textStyle: { fontSize: 12 }
        },
        series: [{
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['35%', '50%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: { show: false, position: 'center' },
          emphasis: {
            label: {
              show: true,
              fontSize: 20,
              fontWeight: 'bold'
            }
          },
          data: res.license_stats.map((item: any) => ({
            name: item.license,
            value: item.count
          }))
        }]
      })
    }
  } catch (error) {
    console.error('加载许可证统计失败:', error)
  }
}

function handleClose() {
  emit('update:modelValue', false)
}
</script>

<style scoped>
.license-charts {
  display: flex;
  gap: 20px;
  padding: 20px 0;
}
.chart-container {
  flex: 1;
  text-align: center;
}
.chart-container h4 {
  margin-bottom: 16px;
  color: #606266;
}
.chart {
  width: 100%;
  height: 300px;
}
.dependency-filter {
  display: flex;
  align-items: center;
}
.total-info {
  margin-top: 12px;
  text-align: right;
  color: #909399;
  font-size: 13px;
}
</style>
