<template>
  <el-dialog v-model="visible" :title="`${projectData?.pro_code || ''} - 项目详情`" width="1000px" @close="handleClose">
    <!-- 头部信息 -->
    <div class="project-header">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="项目名称">{{ projectData?.pro_code }}</el-descriptions-item>
        <el-descriptions-item label="当前Commit ID">{{ projectData?.current_commit_id }}</el-descriptions-item>
      </el-descriptions>
      <div class="query-bar" style="margin-top: 16px">
        <el-input v-model="commitQuery" placeholder="输入Commit ID查询（支持历史版本）" clearable style="width: 220px" @clear="handleQuery" @keyup.enter="handleQuery">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-select v-model="moduleQuery" placeholder="选择子模块" clearable style="width: 180px; margin-left: 12px" @change="handleQuery" @clear="handleQuery">
          <el-option v-for="module in moduleOptions" :key="module" :label="module" :value="module" />
        </el-select>
        <el-button type="primary" @click="handleQuery" style="margin-left: 12px">
          <el-icon><Search /></el-icon> 查询
        </el-button>
        <el-button @click="resetQuery">重置</el-button>
      </div>
    </div>

    <el-tabs v-model="activeTab" style="margin-top: 16px">
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
          <el-table-column prop="module_code" label="模块名称" min-width="120" />
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
const commitQuery = ref('')
const moduleQuery = ref('')
const moduleOptions = ref<string[]>([])

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
  console.log('=== 弹框状态变化 ===')
  console.log('modelValue:', val)
  console.log('projectData:', props.projectData)
  console.log('current_commit_id:', props.projectData?.current_commit_id)
  
  if (val && props.projectData?.current_commit_id) {
    console.log('开始加载数据，commit_id:', props.projectData.current_commit_id)
    activeTab.value = 'license'
    commitQuery.value = ''
    moduleQuery.value = ''
    filterKeyword.value = ''
    filterVersion.value = ''
    await loadModuleOptions()
    handleQuery()
  } else {
    console.warn('不满足加载条件，modelValue:', val, 'current_commit_id:', props.projectData?.current_commit_id)
  }
})

// 加载子模块选项
async function loadModuleOptions() {
  try {
    const res = await api.get('/asset/inner-projects/components', {
      params: { commit_id: props.projectData.current_commit_id }
    })
    const data = Array.isArray(res) ? res : (res?.data || [])
    // 提取唯一的模块名
    const modules = [...new Set(data.map((item: any) => item.module_code).filter(Boolean))]
    moduleOptions.value = modules
    console.log('子模块选项:', modules)
  } catch (error) {
    console.error('加载子模块选项失败:', error)
    moduleOptions.value = []
  }
}

// 统一查询处理
function handleQuery() {
  loadLicenseStats()
  loadComponents()
}

// 重置查询
function resetQuery() {
  commitQuery.value = ''
  moduleQuery.value = ''
  filterKeyword.value = ''
  filterVersion.value = ''
  currentPage.value = 1
  handleQuery()
}

function handlePageChange() {
  // 分页变化时无需重新加载数据
}

async function loadComponents() {
  console.log('=== 加载组件 ===')
  console.log('projectData.current_commit_id:', props.projectData?.current_commit_id)
  console.log('commitQuery:', commitQuery.value)
  console.log('moduleQuery:', moduleQuery.value)
  
  if (!props.projectData?.current_commit_id) {
    console.error('没有current_commit_id，无法加载组件')
    return
  }
  
  loading.value = true
  currentPage.value = 1 // 重置分页
  try {
    const params: any = { commit_id: props.projectData.current_commit_id }
    if (commitQuery.value) params.commit_id = commitQuery.value
    if (moduleQuery.value) params.module_code = moduleQuery.value
    if (filterKeyword.value) params.keyword = filterKeyword.value
    if (filterVersion.value) params.version = filterVersion.value
    
    console.log('请求API，参数:', params)
    const res = await api.get('/asset/inner-projects/components', { params })
    console.log('组件API返回:', res)
    console.log('组件数据类型:', typeof res, Array.isArray(res))
    components.value = Array.isArray(res) ? res : (res?.data || [])
    console.log('最终组件数量:', components.value.length)
  } catch (error) {
    console.error('加载组件失败:', error)
    components.value = []
  } finally {
    loading.value = false
  }
}

async function loadLicenseStats() {
  console.log('=== 加载许可证统计 ===')
  console.log('projectData.current_commit_id:', props.projectData?.current_commit_id)
  
  if (!props.projectData?.current_commit_id) {
    console.error('没有current_commit_id，无法加载许可证统计')
    return
  }
  
  try {
    // 许可证统计使用与组件相同的筛选条件
    const params: any = { commit_id: props.projectData.current_commit_id }
    if (commitQuery.value) params.commit_id = commitQuery.value
    if (moduleQuery.value) params.module_code = moduleQuery.value
    if (filterKeyword.value) params.keyword = filterKeyword.value
    if (filterVersion.value) params.version = filterVersion.value
    
    const res = await api.get('/asset/inner-projects/license-stats', { params })
    
    console.log('许可证统计API返回:', res)
    console.log('完整res对象:', JSON.stringify(res, null, 2))
    console.log('res类型:', typeof res, 'Array.isArray:', Array.isArray(res))
    console.log('res?.group_stats:', res?.group_stats)
    console.log('res?.data:', res?.data)
    console.log('res?.data?.group_stats:', res?.data?.group_stats)
    console.log('Object.keys(res):', Object.keys(res))
    
    // 处理响应数据结构 - res直接就是{group_stats, license_stats}
    const statsData = res
    
    console.log('处理后数据:', statsData)
    console.log('group_stats:', statsData?.group_stats)
    console.log('license_stats:', statsData?.license_stats)
    
    await nextTick()
    await nextTick()
    
    // 渲染许可证分级饼图 - 销毁旧图表，创建新图表
    if (groupChartRef.value && statsData?.group_stats && Array.isArray(statsData.group_stats) && statsData.group_stats.length > 0) {
      console.log('渲染许可证分级饼图，数据:', statsData.group_stats)
      
      // 销毁旧图表实例
      if (groupChart) {
        groupChart.dispose()
      }
      
      // 创建新图表实例
      groupChart = echarts.init(groupChartRef.value)
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
          itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
          label: { show: false, position: 'center' },
          emphasis: { label: { show: true, fontSize: 20, fontWeight: 'bold' } },
          data: statsData.group_stats.map((item: any) => ({
            name: item.group_name || '未分类',
            value: item.count
          }))
        }]
      })
      console.log('许可证分级饼图渲染完成')
    } else {
      console.warn('group_stats为空或无数据:', statsData?.group_stats)
      // 清空图表
      if (groupChart) {
        groupChart.clear()
      }
    }
    
    // 渲染许可证信息饼图 - 销毁旧图表，创建新图表
    if (licenseChartRef.value && statsData?.license_stats && Array.isArray(statsData.license_stats) && statsData.license_stats.length > 0) {
      console.log('渲染许可证信息饼图，数据:', statsData.license_stats)
      
      // 销毁旧图表实例
      if (licenseChart) {
        licenseChart.dispose()
      }
      
      // 创建新图表实例
      licenseChart = echarts.init(licenseChartRef.value)
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
          itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
          label: { show: false, position: 'center' },
          emphasis: { label: { show: true, fontSize: 20, fontWeight: 'bold' } },
          data: statsData.license_stats.map((item: any) => ({
            name: item.license || '未定义',
            value: item.count
          }))
        }]
      })
      console.log('许可证信息饼图渲染完成')
    } else {
      console.warn('license_stats为空或无数据:', statsData?.license_stats)
      // 清空图表
      if (licenseChart) {
        licenseChart.clear()
      }
    }
  } catch (error) {
    console.error('加载许可证统计失败:', error)
  }
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
