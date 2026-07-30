<template>
  <div class="my-projects">
    <div class="page-header">
      <h2>我的项目</h2>
    </div>

    <!-- 项目列表 -->
    <el-table :data="projects" stripe v-loading="loading">
      <el-table-column label="项目分类" width="120">
        <template #default="{ row }">
          <el-tag :type="row.is_boundary === 1 ? 'warning' : 'info'" size="small">
            {{ row.is_boundary === 1 ? '边界' : '非边界' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="pro_code" label="项目名称" min-width="160" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag :type="typeTagColor(row.type)" size="small">{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="当前版本" width="180">
        <template #default="{ row }">
          <el-link type="primary" @click="viewProjectDetail(row)">{{ row.current_commit_id || '-' }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="branch" label="分支" width="120" />
      <el-table-column prop="group_name" label="团队" width="120" />
      <el-table-column prop="dept_name" label="部门" width="120" />
      <el-table-column prop="create_time" label="创建时间" width="160">
        <template #default="{ row }">
          {{ formatTime(row.create_time) }}
        </template>
      </el-table-column>
      <el-table-column prop="update_time" label="更新时间" width="160">
        <template #default="{ row }">
          {{ formatTime(row.update_time) }}
        </template>
      </el-table-column>
    </el-table>

    <!-- 项目详情弹框 -->
    <el-dialog v-model="detailVisible" :title="`${currentProject?.pro_code || ''} - 项目详情`" width="1000px" @close="handleClose">
      <!-- 头部信息 -->
      <div class="project-header">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="项目名称">{{ currentProject?.pro_code }}</el-descriptions-item>
          <el-descriptions-item label="当前Commit ID">{{ currentProject?.current_commit_id }}</el-descriptions-item>
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

        <!-- Tab3: 阻断记录 -->
        <el-tab-pane label="阻断记录" name="blockRecords">
          <el-table :data="blockRecords" stripe v-loading="blockLoading">
            <el-table-column prop="component_name" label="组件名称" min-width="160" />
            <el-table-column prop="version" label="版本" width="120" />
            <el-table-column prop="reason" label="阻断原因" min-width="200" />
            <el-table-column prop="create_time" label="阻断时间" width="160">
              <template #default="{ row }">
                {{ formatTime(row.create_time) }}
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { useUserStore } from '@/stores/user'
import { Search } from '@element-plus/icons-vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'

const userStore = useUserStore()
const projects = ref<any[]>([])
const loading = ref(false)

const userTeam = computed(() => userStore.userInfo?.team || '')

// 详情弹框
const detailVisible = ref(false)
const activeTab = ref('license')
const currentProject = ref<any>(null)
const components = ref<any[]>([])
const blockRecords = ref<any[]>([])
const blockLoading = ref(false)
const filterKeyword = ref('')
const filterVersion = ref('')
const commitQuery = ref('')
const moduleQuery = ref('')
const moduleOptions = ref<string[]>([])
const currentPage = ref(1)
const pageSize = ref(10)

// 图表实例
let groupChart: echarts.ECharts | null = null
let licenseChart: echarts.ECharts | null = null
const groupChartRef = ref<HTMLElement>()
const licenseChartRef = ref<HTMLElement>()

const paginatedComponents = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return components.value.slice(start, end)
})

async function loadProjects() {
  if (!userTeam.value) {
    ElMessage.warning('未获取到团队信息')
    return
  }

  loading.value = true
  try {
    const res = await axios.get('/api/v1/asset/inner-projects', {
      params: { group_name: userTeam.value }
    })
    if (res.data.code === 200 || res.data.code === 0) {
      projects.value = res.data.data || []
    }
  } catch (error) {
    console.error('加载项目失败:', error)
    ElMessage.error('加载项目失败')
  } finally {
    loading.value = false
  }
}

function typeTagColor(type: string) {
  const colors: Record<string, any> = {
    '自研': 'primary',
    '外包': 'warning',
    '开源': 'success',
  }
  return colors[type] || 'info'
}

function formatTime(time: string) {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 19)
}

function viewProjectDetail(row: any) {
  currentProject.value = row
  commitQuery.value = ''
  moduleQuery.value = ''
  filterKeyword.value = ''
  filterVersion.value = ''
  activeTab.value = 'license'
  detailVisible.value = true
}

async function loadModuleOptions() {
  try {
    const res = await axios.get('/api/v1/asset/inner-projects/components', {
      params: { commit_id: currentProject.value.current_commit_id }
    })
    const data = res.data.data || []
    const modules = [...new Set(data.map((item: any) => item.module_code).filter(Boolean))]
    moduleOptions.value = modules
  } catch (error) {
    console.error('加载子模块选项失败:', error)
    moduleOptions.value = []
  }
}

function handleQuery() {
  loadLicenseStats()
  loadComponents()
  loadBlockRecords()
}

function resetQuery() {
  commitQuery.value = ''
  moduleQuery.value = ''
  filterKeyword.value = ''
  filterVersion.value = ''
  currentPage.value = 1
  handleQuery()
}

function handlePageChange() {}

async function loadComponents() {
  if (!currentProject.value?.current_commit_id) return
  
  loading.value = true
  currentPage.value = 1
  try {
    const params: any = { commit_id: currentProject.value.current_commit_id }
    if (commitQuery.value) params.commit_id = commitQuery.value
    if (moduleQuery.value) params.module_code = moduleQuery.value
    if (filterKeyword.value) params.keyword = filterKeyword.value
    if (filterVersion.value) params.version = filterVersion.value
    
    const res = await axios.get('/api/v1/asset/inner-projects/components', { params })
    components.value = res.data.data || []
  } catch (error) {
    console.error('加载组件失败:', error)
    components.value = []
  } finally {
    loading.value = false
  }
}

async function loadLicenseStats() {
  if (!currentProject.value?.current_commit_id) return
  
  try {
    const params: any = { commit_id: currentProject.value.current_commit_id }
    if (commitQuery.value) params.commit_id = commitQuery.value
    if (moduleQuery.value) params.module_code = moduleQuery.value
    if (filterKeyword.value) params.keyword = filterKeyword.value
    if (filterVersion.value) params.version = filterVersion.value
    
    const res = await axios.get('/api/v1/asset/inner-projects/license-stats', { params })
    const statsData = res.data.data || res.data || {}
    
    await nextTick()
    await nextTick()
    
    // 渲染许可证分级饼图
    if (groupChartRef.value && statsData?.group_stats && Array.isArray(statsData.group_stats) && statsData.group_stats.length > 0) {
      if (groupChart) groupChart.dispose()
      groupChart = echarts.init(groupChartRef.value)
      groupChart.setOption({
        tooltip: { trigger: 'item' },
        legend: { type: 'scroll', orient: 'vertical', right: '5%', top: 'center', textStyle: { fontSize: 12 } },
        series: [{
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['35%', '50%'],
          avoidLabelOverlap: false,
          itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
          label: { show: false, position: 'center' },
          emphasis: { label: { show: true, fontSize: 20, fontWeight: 'bold' } },
          data: statsData.group_stats.map((item: any) => ({ name: item.group_name || '未分类', value: item.count }))
        }]
      })
    }
    
    // 渲染许可证信息饼图
    if (licenseChartRef.value && statsData?.license_stats && Array.isArray(statsData.license_stats) && statsData.license_stats.length > 0) {
      if (licenseChart) licenseChart.dispose()
      licenseChart = echarts.init(licenseChartRef.value)
      licenseChart.setOption({
        tooltip: { trigger: 'item' },
        legend: { type: 'scroll', orient: 'vertical', right: '5%', top: 'center', textStyle: { fontSize: 12 } },
        series: [{
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['35%', '50%'],
          avoidLabelOverlap: false,
          itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
          label: { show: false, position: 'center' },
          emphasis: { label: { show: true, fontSize: 20, fontWeight: 'bold' } },
          data: statsData.license_stats.map((item: any) => ({ name: item.license || '未定义', value: item.count }))
        }]
      })
    }
  } catch (error) {
    console.error('加载许可证统计失败:', error)
  }
}

async function loadBlockRecords() {
  if (!currentProject.value?.id) return
  
  blockLoading.value = true
  try {
    const res = await axios.get('/api/v1/risk/cicd-check/block-records', {
      params: { project_id: currentProject.value.id }
    })
    blockRecords.value = res.data.data || []
  } catch (error) {
    console.error('加载阻断记录失败:', error)
    blockRecords.value = []
  } finally {
    blockLoading.value = false
  }
}

function handleClose() {
  if (groupChart) {
    groupChart.dispose()
    groupChart = null
  }
  if (licenseChart) {
    licenseChart.dispose()
    licenseChart = null
  }
  detailVisible.value = false
}

// 监听弹框打开
import { watch } from 'vue'
watch(detailVisible, async (val) => {
  if (val && currentProject.value?.current_commit_id) {
    await loadModuleOptions()
    handleQuery()
  }
})

onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.my-projects {
  padding: 16px;
}

.page-header {
  margin-bottom: 16px;
}

.page-header h2 {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 600;
  color: #1f2328;
}

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
</style>
