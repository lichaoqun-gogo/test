<template>
  <div class="dependency-library">
    <div class="page-header">
      <div class="filter-bar">
        <el-input v-model="searchText" placeholder="搜索组件名称" clearable style="width: 200px" @clear="handleSearch" @keyup.enter="handleSearch">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-select v-model="filterSource" placeholder="依赖范围" style="width: 120px; margin-left: 12px" @change="handleSearch">
          <el-option label="自研项目" value="inner" />
          <el-option label="开源软件" value="software" />
        </el-select>
        <el-select v-model="filterClassifier" placeholder="组件类型" clearable style="width: 130px; margin-left: 12px" @change="handleSearch">
          <el-option label="library" value="library" />
          <el-option label="framework" value="framework" />
          <el-option label="tool" value="tool" />
        </el-select>
        <el-button type="primary" plain @click="handleSearch" style="margin-left: 12px">
          <el-icon><Search /></el-icon> 查询
        </el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>
    </div>

    <el-table :data="paginatedAssets" stripe style="margin-top: 16px" v-loading="loading">
      <el-table-column prop="_group" label="组件组织" min-width="160" />
      <el-table-column prop="PUBLISHER" label="供应商" min-width="140" />
      <el-table-column prop="_name" label="组件名称" min-width="160" />
      <el-table-column prop="_version" label="版本" width="110" />
      <el-table-column prop="PURL" label="PURL" min-width="200" show-overflow-tooltip />
      <el-table-column prop="LICENSE" label="许可证" width="130" />
      <el-table-column label="操作" width="140">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="viewRelatedProjects(row)">关联项目</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>

    <!-- 关联项目对话框（自研项目/开源软件两套列，按依赖范围切换） -->
    <el-dialog v-model="relatedProjectsDialogVisible" :title="filterSource === 'software' ? '关联开源软件' : '关联自研项目'" width="700px">
      <el-table :data="relatedProjects" stripe v-loading="loadingProjects">
        <template v-if="filterSource === 'software'">
          <el-table-column prop="_name" label="软件名称" min-width="150" />
          <el-table-column prop="_version" label="版本" width="110" />
          <el-table-column prop="maint_dept" label="维护部门" width="140" />
          <el-table-column prop="applicable_level" label="适用等级" width="120" />
        </template>
        <template v-else>
          <el-table-column prop="pro_code" label="项目名称" min-width="150" />
          <el-table-column prop="group_name" label="项目团队" width="120" />
          <el-table-column prop="dept_name" label="项目部门" width="140" />
          <el-table-column label="项目分类" width="120">
            <template #default="{ row }">
              <el-tag :type="row.is_boundary === 1 ? 'warning' : 'info'" size="small">
                {{ row.is_boundary === 1 ? '边界系统' : '非边界系统' }}
              </el-tag>
            </template>
          </el-table-column>
        </template>
      </el-table>
      <div class="total-info">共 {{ relatedProjects.length }} 个项目</div>
    </el-dialog>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { getDependencies, getRelatedProjects } from '@/api'

const loading = ref(false)
const loadingProjects = ref(false)
const searchText = ref('')
// 依赖范围：inner=自研项目组件（component_snapshot，默认）/ software=开源软件组件（component 表）
const filterSource = ref('inner')
const filterClassifier = ref('')
const relatedProjectsDialogVisible = ref(false)
const currentComponent = ref<any>(null)
const relatedProjects = ref<any[]>([])

const dependencies = ref<any[]>([])

const filteredAssets = computed(() => {
  let list = [...(dependencies.value || [])]
  if (searchText.value) {
    const q = searchText.value.toLowerCase()
    list = list.filter(a => a._name.toLowerCase().includes(q) || a._group.toLowerCase().includes(q))
  }
  if (filterClassifier.value) list = list.filter(a => a.CLASSIFIER === filterClassifier.value)
  return list
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

const paginatedAssets = computed(() => {
  const start = (pagination.page - 1) * pagination.pageSize
  return filteredAssets.value.slice(start, start + pagination.pageSize)
})

onMounted(() => {
  loadDependencies()
})

async function loadDependencies() {
  loading.value = true
  try {
    const params: any = { source: filterSource.value }
    if (searchText.value) params.keyword = searchText.value
    if (filterClassifier.value) params.classifier = filterClassifier.value
    
    const res = await getDependencies(params)
    dependencies.value = res || []
    pagination.total = dependencies.value.length
  } catch (error) {
    console.error('加载依赖包失败:', error)
    dependencies.value = []
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  loadDependencies()
}

function handleReset() {
  searchText.value = ''
  filterSource.value = 'inner'
  filterClassifier.value = ''
  handleSearch()
}

function handleSizeChange() {
  pagination.page = 1
}

function handlePageChange() {}

async function viewRelatedProjects(row: any) {
  currentComponent.value = row
  relatedProjectsDialogVisible.value = true
  loadingProjects.value = true
  try {
    const res = await getRelatedProjects(row.id, filterSource.value)
    relatedProjects.value = res || []
  } catch (error) {
    console.error('加载关联项目失败:', error)
    relatedProjects.value = []
  } finally {
    loadingProjects.value = false
  }
}
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}
.filter-bar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
