<template>
  <el-dialog v-model="visible" :title="`${projectData?.pro_code || ''} - 项目详情`" width="1000px" @close="handleClose">
    <!-- 头部信息 -->
    <div class="project-header">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="项目名称">{{ projectData?.pro_code }}</el-descriptions-item>
        <el-descriptions-item label="当前版本">{{ projectData?.current_ID || '-' }}</el-descriptions-item>
      </el-descriptions>
      <div class="query-bar" style="margin-top: 16px">
        <el-input v-model="currentIDQuery" placeholder="输入版本标识查询（支持历史版本）" clearable style="width: 220px" @clear="handleQuery" @keyup.enter="handleQuery">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-select v-model="selectedModuleId" placeholder="选择子模块" clearable style="width: 180px; margin-left: 12px" @change="handleModuleChange" @clear="handleModuleClear">
          <el-option v-for="m in moduleOptions" :key="m.id" :label="m.module_code" :value="m.id" />
        </el-select>
        <el-button type="primary" @click="handleQuery" style="margin-left: 12px">
          <el-icon><Search /></el-icon> 查询
        </el-button>
        <el-button @click="resetQuery">重置</el-button>
      </div>
    </div>

    <el-tabs v-model="activeTab" style="margin-top: 16px" @tab-change="handleTabChange">
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

      <!-- Tab2: 组件信息 -->
      <el-tab-pane label="组件信息" name="components">
        <div class="component-filter">
          <el-input v-model="filterKeyword" placeholder="搜索组件名称/组织" clearable style="width: 200px" @clear="handleQuery" @keyup.enter="handleQuery">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
          <el-input v-model="filterVersion" placeholder="版本号" clearable style="width: 150px; margin-left: 12px" @clear="handleQuery" @keyup.enter="handleQuery" />
          <el-button type="primary" @click="handleQuery" style="margin-left: 12px">
            <el-icon><Search /></el-icon> 查询
          </el-button>
        </div>

        <el-table :data="paginatedComponents" stripe style="margin-top: 16px" v-loading="loading">
          <el-table-column prop="module_code" label="模块名称" min-width="120">
            <template #default="{ row }">
              {{ row.module_code || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="_group" label="组织" min-width="160" />
          <el-table-column prop="_name" label="组件名称" min-width="160" />
          <el-table-column prop="_version" label="版本" width="120" />
          <el-table-column prop="LICENSE" label="许可协议" min-width="140" />
          <el-table-column label="是否内部" width="100">
            <template #default="{ row }">
              <el-tag :type="row.INTERNAL === 1 ? 'success' : 'info'" size="small">
                {{ row.INTERNAL === 1 ? '内部' : '外部' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container" style="margin-top: 16px; display: flex; justify-content: flex-end">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="components.length"
            layout="total, sizes, prev, pager, next"
            @size-change="handlePageChange"
            @current-change="handlePageChange"
          />
        </div>
      </el-tab-pane>
    </el-tabs>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, computed } from 'vue'
import { Search } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import api from '@/api'

const props = defineProps<{
  modelValue: boolean
  projectData: any
}>()

const emit = defineEmits(['update:modelValue'])

const activeTab = ref('license')
const loading = ref(false)
const components = ref<any[]>([])
const filterKeyword = ref('')
const filterVersion = ref('')
const currentIDQuery = ref('')
const selectedModuleId = ref<number | null>(null)  // 选中的 project_inner.id
const moduleOptions = ref<{ id: number; module_code: string }[]>([])

// 分页
const currentPage = ref(1)
const pageSize = ref(10)

// 图表实例
let groupChart: echarts.ECharts | null = null
let licenseChart: echarts.ECharts | null = null

const groupChartRef = ref<HTMLElement>()
const licenseChartRef = ref<HTMLElement>()

// 分页后的组件数据
const paginatedComponents = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return components.value.slice(start, end)
})

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

watch(() => props.modelValue, async (val) => {
  if (val && props.projectData?.current_ID) {
    activeTab.value = 'license'
    currentIDQuery.value = ''
    selectedModuleId.value = null
    filterKeyword.value = ''
    filterVersion.value = ''
    await loadModuleOptions()
    handleQuery()
  }
})

// 加载子模块选项：通过 pro_code 获取项目的所有模块列表（包含 self）
async function loadModuleOptions() {
  try {
    const res = await api.get('/asset/inner-projects/modules', {
      params: { pro_code: props.projectData.pro_code }
    })
    const data = Array.isArray(res) ? res : (res?.data || [])
    moduleOptions.value = data
      .filter((item: any) => item.id && item.module_code)
      .map((item: any) => ({ id: item.id, module_code: item.module_code }))
  } catch (error) {
    console.error('加载子模块选项失败:', error)
    moduleOptions.value = []
  }
}

// 统一查询处理：未选择模块时查询该项目全部模块
function handleQuery() {
  let projectIDs = ''
  if (selectedModuleId.value) {
    projectIDs = selectedModuleId.value.toString()
  } else if (moduleOptions.value.length) {
    projectIDs = moduleOptions.value.map(m => m.id).join(',')
  }

  loadLicenseStats(projectIDs)
  loadComponents(projectIDs)
}

// 模块选择变化
function handleModuleChange() {
  handleQuery()
}

// 模块选择清空
function handleModuleClear() {
  handleQuery()
}

// 重置查询
function resetQuery() {
  currentIDQuery.value = ''
  selectedModuleId.value = null
  filterKeyword.value = ''
  filterVersion.value = ''
  currentPage.value = 1
  handleQuery()
}

function handlePageChange() {
  // 分页变化时无需重新加载数据
}

async function loadComponents(projectIDs: string) {
  if (!props.projectData?.current_ID) {
    return
  }

  loading.value = true
  currentPage.value = 1
  try {
    const params: any = { current_ID: props.projectData.current_ID }
    if (currentIDQuery.value) params.current_ID = currentIDQuery.value
    if (projectIDs) params.project_ids = projectIDs
    if (filterKeyword.value) params.keyword = filterKeyword.value
    if (filterVersion.value) params.version = filterVersion.value

    const res = await api.get('/asset/inner-projects/components', { params })
    components.value = Array.isArray(res) ? res : (res?.data || [])
  } catch (error) {
    console.error('加载组件失败:', error)
    components.value = []
  } finally {
    loading.value = false
  }
}

async function loadLicenseStats(projectIDs: string) {
  if (!props.projectData?.current_ID) {
    return
  }

  try {
    const params: any = { current_ID: props.projectData.current_ID }
    if (currentIDQuery.value) params.current_ID = currentIDQuery.value
    if (projectIDs) params.project_ids = projectIDs
    if (filterKeyword.value) params.keyword = filterKeyword.value
    if (filterVersion.value) params.version = filterVersion.value

    const res = await api.get('/asset/inner-projects/license-stats', { params })
    const statsData = res || {}

    await nextTick()

    renderGroupChart(statsData.group_stats || [])
    renderLicenseChart(statsData.license_stats || [])
  } catch (error) {
    console.error('加载许可证统计失败:', error)
  }
}

// tab 切换时 resize 图表，确保容器从隐藏变为可见后尺寸正确
function handleTabChange() {
  nextTick(() => {
    if (groupChart) groupChart.resize()
    if (licenseChart) licenseChart.resize()
  })
}

// 渲染许可证分级饼图
function renderGroupChart(groupStats: any[]) {
  if (!groupChartRef.value) return

  if (groupChart) {
    groupChart.dispose()
    groupChart = null
  }

  groupChart = echarts.init(groupChartRef.value)

  if (groupStats.length === 0) {
    groupChart.setOption({
      title: { text: '暂无数据', left: 'center', top: 'center', textStyle: { color: '#909399', fontSize: 14 } }
    })
    return
  }

  groupChart.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
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
      itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
      label: { show: false, position: 'center' },
      emphasis: { label: { show: true, fontSize: 16, fontWeight: 'bold' } },
      data: groupStats.map((item: any) => ({
        name: item.group_name || '未分类',
        value: item.count
      }))
    }]
  })
}

// 渲染许可证信息饼图
function renderLicenseChart(licenseStats: any[]) {
  if (!licenseChartRef.value) return

  if (licenseChart) {
    licenseChart.dispose()
    licenseChart = null
  }

  licenseChart = echarts.init(licenseChartRef.value)

  if (licenseStats.length === 0) {
    licenseChart.setOption({
      title: { text: '暂无数据', left: 'center', top: 'center', textStyle: { color: '#909399', fontSize: 14 } }
    })
    return
  }

  licenseChart.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
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
      itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
      label: { show: false, position: 'center' },
      emphasis: { label: { show: true, fontSize: 16, fontWeight: 'bold' } },
      data: licenseStats.map((item: any) => ({
        name: item.license || '未定义',
        value: item.count
      }))
    }]
  })
}

function handleClose() {
  // 销毁图表实例
  if (groupChart) {
    groupChart.dispose()
    groupChart = null
  }
  if (licenseChart) {
    licenseChart.dispose()
    licenseChart = null
  }
  emit('update:modelValue', false)
}
</script>

<style scoped>
.project-header {
  padding-bottom: 16px;
  border-bottom: 1px solid #e4e7ed;
}
.query-bar {
  display: flex;
  align-items: center;
}
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
.component-filter {
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
