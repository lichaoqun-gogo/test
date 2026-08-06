<template>
  <div class="self-developed-projects">
    <div class="page-header">
      <h2 class="page-title">自研项目管理</h2>
    </div>

    <!-- 查询栏 -->
    <div class="filter-bar">
      <el-input v-model="searchText" placeholder="搜索项目名称" clearable style="width: 200px" @clear="handleSearch" @keyup.enter="handleSearch">
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>
      <el-select v-model="filterDept" placeholder="项目部门" clearable style="width: 150px; margin-left: 12px" @change="handleSearch">
        <el-option label="技术研发部" value="技术研发部" />
        <el-option label="平台研发部" value="平台研发部" />
        <el-option label="运维部" value="运维部" />
        <el-option label="安全部" value="安全部" />
        <el-option label="质量管理部" value="质量管理部" />
      </el-select>
      <el-select v-model="filterTeam" placeholder="项目团队" clearable style="width: 150px; margin-left: 12px" @change="handleSearch">
        <el-option label="订单团队" value="订单团队" />
        <el-option label="支付团队" value="支付团队" />
        <el-option label="用户团队" value="用户团队" />
        <el-option label="基础架构团队" value="基础架构团队" />
      </el-select>
      <el-button type="primary" plain @click="handleSearch" style="margin-left: 12px">
        <el-icon><Search /></el-icon> 查询
      </el-button>
      <el-button @click="handleReset">重置</el-button>
      <el-button type="success" plain @click="batchSetBoundaryDialogVisible = true" style="margin-left: 12px">
        <el-icon><Setting /></el-icon> 设置项目分类
      </el-button>
    </div>

    <!-- 项目列表 -->
    <el-table :data="paginatedProjects" stripe style="margin-top: 16px" v-loading="loading" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="50" />
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
          <el-link type="primary" @click="viewProjectDetail(row)">{{ row.current_ID || '-' }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="branch" label="分支" width="150" />
      <el-table-column prop="group_name" label="团队" width="120" />
      <el-table-column prop="dept_name" label="部门" width="140" />
      <el-table-column prop="create_time" label="创建时间" width="160">
        <template #default="{ row }">{{ formatTime(row.create_time) }}</template>
      </el-table-column>
      <el-table-column prop="update_time" label="更新时间" width="160">
        <template #default="{ row }">{{ formatTime(row.update_time) }}</template>
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

    <!-- 项目详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="项目详情" width="700px">
      <el-descriptions v-if="currentProject" :column="2" border>
        <el-descriptions-item label="项目编码">{{ currentProject.pro_code }}</el-descriptions-item>
        <el-descriptions-item label="模块名">{{ currentProject.module_code }}</el-descriptions-item>
        <el-descriptions-item label="分支名">{{ currentProject.branch }}</el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag :type="typeTagColor(currentProject.type)" size="small">{{ currentProject.type }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="业务组">{{ currentProject.group_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="部门">{{ currentProject.dept_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="边界系统">
          <el-tag :type="currentProject.is_boundary === 1 ? 'warning' : 'info'" size="small">
            {{ currentProject.is_boundary === 1 ? '边界系统' : '非边界系统' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentProject.active === 1 ? 'success' : 'info'" size="small">
            {{ currentProject.active === 1 ? '正常' : '关闭' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="当前版本">{{ currentProject.current_ID || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatTime(currentProject.create_time) }}</el-descriptions-item>
        <el-descriptions-item label="创建者">{{ currentProject.create_by || '-' }}</el-descriptions-item>
        <el-descriptions-item label="最后修改人">{{ currentProject.update_by || '-' }}</el-descriptions-item>
        <el-descriptions-item label="更新时间" :span="2">{{ formatTime(currentProject.update_time) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!-- Add/Edit Dialog -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑自研项目' : '新增自研项目'" width="600px" @close="resetForm">
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="项目编码" prop="pro_code">
          <el-input v-model="formData.pro_code" placeholder="请输入项目编码" />
        </el-form-item>
        <el-form-item label="模块名" prop="module_code">
          <el-input v-model="formData.module_code" placeholder="请输入模块名" />
        </el-form-item>
        <el-form-item label="分支名" prop="branch">
          <el-input v-model="formData.branch" placeholder="请输入分支名" />
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="类型" prop="type">
              <el-select v-model="formData.type" placeholder="请选择类型" style="width: 100%">
                <el-option label="Java" value="java" />
                <el-option label="Python" value="python" />
                <el-option label="Go" value="go" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="边界系统" prop="is_boundary">
              <el-select v-model="formData.is_boundary" placeholder="请选择" style="width: 100%">
                <el-option label="边界系统" :value="1" />
                <el-option label="非边界系统" :value="0" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="业务组" prop="group_name">
          <el-input v-model="formData.group_name" placeholder="请输入业务组" />
        </el-form-item>
        <el-form-item label="部门" prop="dept_name">
          <el-input v-model="formData.dept_name" placeholder="请输入部门" />
        </el-form-item>
        <el-form-item label="状态" prop="active">
          <el-select v-model="formData.active" placeholder="请选择状态" style="width: 100%">
            <el-option label="正常" :value="1" />
            <el-option label="关闭" :value="0" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 设置项目分类对话框 -->
    <el-dialog v-model="batchSetBoundaryDialogVisible" title="设置项目分类" width="600px">
      <el-alert title="请选择需要设置分类的项目" type="info" :closable="false" show-icon style="margin-bottom: 16px" />
      <el-form label-width="100px">
        <el-form-item label="项目分类">
          <el-select v-model="batchBoundaryValue" placeholder="请选择分类" style="width: 100%">
            <el-option label="边界系统" :value="1" />
            <el-option label="非边界系统" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="已选项目">
          <div v-if="selectedRows.length > 0" class="selected-projects">
            <el-tag v-for="row in selectedRows" :key="row.id" type="success" style="margin: 4px">
              {{ row.pro_code }}
            </el-tag>
          </div>
          <el-tag v-else type="info">未选择项目</el-tag>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="batchSetBoundaryDialogVisible = false">取消</el-button>
        <el-button type="primary" :disabled="selectedRows.length === 0" @click="handleBatchSetBoundary">确定设置</el-button>
      </template>
    </el-dialog>

    <!-- 项目详情弹框（许可证和组件） -->
    <InnerProjectDetailDialog v-model="projectDetailDialogVisible" :project-data="currentProject" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Search, Plus, Setting } from '@element-plus/icons-vue'
import { getInnerProjects, updateInnerProject } from '@/api'
import InnerProjectDetailDialog from '@/components/business/InnerProjectDetailDialog.vue'
import { formatTime } from '@/utils/format'

const loading = ref(false)
const searchText = ref('')
const filterDept = ref('')
const filterTeam = ref('')
const detailDialogVisible = ref(false)
const currentProject = ref<any | null>(null)
const batchSetBoundaryDialogVisible = ref(false)
const batchBoundaryValue = ref<number>(0)
const selectedRows = ref<any[]>([])

// 项目详情弹框
const projectDetailDialogVisible = ref(false)

function viewProjectDetail(row: any) {
  currentProject.value = row
  projectDetailDialogVisible.value = true
}

// CRUD
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

interface SelfProjectFormData {
  pro_code: string
  module_code: string
  branch: string
  type: string
  is_boundary: number
  group_name: string
  dept_name: string
  active: number
}

const defaultFormData: SelfProjectFormData = {
  pro_code: '',
  module_code: '',
  branch: '',
  type: '',
  is_boundary: 0,
  group_name: '',
  dept_name: '',
  active: 1,
}

const formData = reactive<SelfProjectFormData>({ ...defaultFormData })

const formRules: FormRules = {
  pro_code: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  module_code: [{ required: true, message: '请输入模块名', trigger: 'blur' }],
  branch: [{ required: true, message: '请输入分支名', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
}

const projects = ref<any[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

const filteredProjects = computed(() => {
  let list = [...(projects.value || [])]
  if (searchText.value) {
    const q = searchText.value.toLowerCase()
    list = list.filter(p => p.pro_code.toLowerCase().includes(q))
  }
  if (filterDept.value) list = list.filter(p => p.dept_name === filterDept.value)
  if (filterTeam.value) list = list.filter(p => p.group_name === filterTeam.value)
  
  // 去重：按pro_code分组，每个项目只保留一条
  const uniqueMap = new Map()
  list.forEach(p => {
    if (!uniqueMap.has(p.pro_code)) {
      uniqueMap.set(p.pro_code, p)
    }
  })
  return Array.from(uniqueMap.values())
})

const paginatedProjects = computed(() => {
  const start = (pagination.page - 1) * pagination.pageSize
  return filteredProjects.value.slice(start, start + pagination.pageSize)
})

watch(filteredProjects, (list) => {
  pagination.total = list.length
  if (pagination.page > 1 && (pagination.page - 1) * pagination.pageSize >= list.length) {
    pagination.page = 1
  }
})

onMounted(() => {
  loadProjects()
})

async function loadProjects() {
  loading.value = true
  try {
    const params: any = {}
    if (searchText.value) params.keyword = searchText.value
    if (filterDept.value) params.dept_name = filterDept.value
    if (filterTeam.value) params.group_name = filterTeam.value
    
    const res = await getInnerProjects(params)
    projects.value = res || []
  } catch (error) {
    console.error('加载自研项目失败:', error)
    projects.value = []
  } finally {
    loading.value = false
  }
}

function handleSelectionChange(selection: any[]) {
  selectedRows.value = selection
}

function typeTagColor(type: string): 'primary' | 'success' | 'warning' | 'info' {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info'> = { java: 'primary', python: 'success', go: 'warning' }
  return map[type.toLowerCase()] || 'info'
}

function handleSearch() {
  pagination.page = 1
  loadProjects()
}

function handleReset() {
  searchText.value = ''
  filterDept.value = ''
  filterTeam.value = ''
  pagination.page = 1
  loadProjects()
}

function handleSizeChange() {
  pagination.page = 1
}

function handlePageChange() {
  // pagination is reactive, no action needed
}

function viewDetail(row: any) {
  currentProject.value = row
  detailDialogVisible.value = true
}

function handleAdd() {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

function handleEdit(row: any) {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(formData, {
    pro_code: row.pro_code,
    module_code: row.module_code,
    branch: row.branch,
    type: row.type,
    is_boundary: row.is_boundary,
    group_name: row.group_name || '',
    dept_name: row.dept_name || '',
    active: row.active,
  })
  dialogVisible.value = true
}

function resetForm() {
  Object.assign(formData, defaultFormData)
  editingId.value = null
  formRef.value?.clearValidate()
}

function handleSubmit() {
  formRef.value?.validate((valid) => {
    if (valid) {
      loading.value = true
      setTimeout(() => {
        if (isEdit.value) {
          ElMessage.success(`项目 "${formData.pro_code}" 更新成功`)
        } else {
          ElMessage.success(`项目 "${formData.pro_code}" 新增成功`)
        }
        loading.value = false
        dialogVisible.value = false
      }, 500)
    }
  })
}

async function handleBatchSetBoundary() {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先选择项目')
    return
  }
  
  loading.value = true
  try {
    const promises = selectedRows.value.map(row => 
      updateInnerProject(row.id, { is_boundary: batchBoundaryValue.value })
    )
    await Promise.all(promises)
    ElMessage.success(`成功设置 ${selectedRows.value.length} 个项目`)
    batchSetBoundaryDialogVisible.value = false
    loadProjects()
  } catch (error) {
    console.error('批量设置失败:', error)
    ElMessage.error('批量设置失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.self-developed-projects {
  padding: 16px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.filter-bar {
  display: flex;
  align-items: center;
  margin-top: 16px;
}

.selected-projects {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  max-height: 150px;
  overflow-y: auto;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
