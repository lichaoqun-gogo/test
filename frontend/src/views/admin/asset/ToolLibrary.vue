<template>
  <div class="tool-library">
    <div class="page-header">
      <div class="filter-bar">
        <el-input v-model="searchText" placeholder="搜索工具名称" clearable style="width: 200px" @clear="handleSearch" @keyup.enter="handleSearch">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-select v-model="filterGrade" placeholder="适用等级" clearable style="width: 130px; margin-left: 12px" @change="handleSearch">
          <el-option v-for="g in ['A', 'B', 'C', 'D']" :key="g" :label="`${g}级`" :value="g" />
        </el-select>
        <el-select v-model="filterType" placeholder="类型" clearable style="width: 130px; margin-left: 12px" @change="handleSearch">
          <el-option label="APPLICATION" value="APPLICATION" />
          <el-option label="FRAMEWORK" value="FRAMEWORK" />
          <el-option label="LIBRARY" value="LIBRARY" />
          <el-option label="CONTAINER" value="CONTAINER" />
          <el-option label="OPERATING_SYSTEM" value="OPERATING_SYSTEM" />
          <el-option label="DEVICE" value="DEVICE" />
          <el-option label="FIRMWARE" value="FIRMWARE" />
          <el-option label="FILE" value="FILE" />
        </el-select>
        <el-button type="primary" plain @click="handleSearch" style="margin-left: 12px">
          <el-icon><Search /></el-icon> 查询
        </el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon> 新增
      </el-button>
    </div>

    <el-table :data="toolList" stripe style="margin-top: 16px" v-loading="loading" row-key="id">
      <!-- 展开行 -->
      <el-table-column type="expand">
        <template #default="{ row }">
          <div class="expand-content">
            <el-row :gutter="20">
              <el-col :span="12">
                <div class="expand-item">
                  <span class="expand-label">BOM最后更新时间：</span>
                  <span class="expand-value">{{ formatTime(row.LAST_BOM_IMPORTED) }}</span>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="expand-item">
                  <span class="expand-label">更新者：</span>
                  <span class="expand-value">{{ row.update_by || '-' }}</span>
                </div>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <div class="expand-item">
                  <span class="expand-label">更新时间：</span>
                  <span class="expand-value">{{ formatTime(row.update_time) }}</span>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="expand-item">
                  <span class="expand-label">创建者：</span>
                  <span class="expand-value">{{ row.create_by || '-' }}</span>
                </div>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <div class="expand-item">
                  <span class="expand-label">创建时间：</span>
                  <span class="expand-value">{{ formatTime(row.create_time) }}</span>
                </div>
              </el-col>
              <el-col :span="12"></el-col>
            </el-row>
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="_group" label="工具组织" min-width="120" />
      <el-table-column label="工具名称" min-width="140">
        <template #default="{ row }">
          <el-link type="primary" @click="showDetail(row)">{{ row._name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="_version" label="版本" width="110" />
      <el-table-column prop="PUBLISHER" label="供应商" min-width="120" />
      <el-table-column prop="CLASSIFIER" label="类型" width="140">
        <template #default="{ row }">
          <el-tag :type="typeTagColor(row.CLASSIFIER)" size="small">{{ row.CLASSIFIER }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="license" label="许可证" width="130" />
      <el-table-column label="适用等级" width="100">
        <template #default="{ row }">
          <GradeTag :grade="row.applicable_level" size="small" />
        </template>
      </el-table-column>
      <el-table-column prop="maint_dept" label="维护部门" width="130" />
      <el-table-column prop="use_dept" label="使用部门" width="130" />
      <el-table-column prop="country" label="国家" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.ACTIVE === 1 ? 'success' : 'info'" size="small">
            {{ row.ACTIVE === 1 ? '正常' : '关闭' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
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

    <!-- Add/Edit Dialog -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑开源工具' : '新增开源工具'" width="800px" @close="resetForm">
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="110px">
        <!-- 通用属性 -->
        <el-divider content-position="left">通用属性</el-divider>
        <el-form-item label="软件分类" prop="category">
          <el-select v-model="formData.category" placeholder="请选择软件分类" style="width: 100%">
            <el-option label="开源基础软件" value="software" />
            <el-option label="开源工具" value="tool" />
          </el-select>
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="工具组织" prop="_group">
              <el-input v-model="formData._group" placeholder="请输入项目组织" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="工具名称" prop="_name">
              <el-input v-model="formData._name" placeholder="请输入工具名称" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="版本" prop="_version">
              <el-input v-model="formData._version" placeholder="请输入版本号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="供应商" prop="PUBLISHER">
              <el-input v-model="formData.PUBLISHER" placeholder="请输入供应商" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="类型" prop="CLASSIFIER">
              <el-select v-model="formData.CLASSIFIER" placeholder="请选择类型" style="width: 100%">
                <el-option label="APPLICATION" value="APPLICATION" />
                <el-option label="FRAMEWORK" value="FRAMEWORK" />
                <el-option label="LIBRARY" value="LIBRARY" />
                <el-option label="CONTAINER" value="CONTAINER" />
                <el-option label="OPERATING_SYSTEM" value="OPERATING_SYSTEM" />
                <el-option label="DEVICE" value="DEVICE" />
                <el-option label="FIRMWARE" value="FIRMWARE" />
                <el-option label="FILE" value="FILE" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="开发语言" prop="type">
              <el-input v-model="formData.type" placeholder="请输入开发语言" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="开源国家" prop="country">
              <el-input v-model="formData.country" placeholder="请输入开源国家" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="许可证" prop="license">
              <el-input v-model="formData.license" placeholder="请输入许可证" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述" prop="DESCRIPTION">
          <el-input v-model="formData.DESCRIPTION" type="textarea" :rows="2" placeholder="请输入组件描述" />
        </el-form-item>
        <el-form-item label="状态" prop="ACTIVE">
          <el-switch v-model="formData.ACTIVE" :active-value="1" :inactive-value="0" active-text="正常" inactive-text="关闭" />
        </el-form-item>

        <!-- 标识属性 -->
        <el-divider content-position="left">标识属性</el-divider>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="作者" prop="AUTHOR">
              <el-input v-model="formData.AUTHOR" placeholder="请输入作者" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="PURL" prop="PURL">
              <el-input v-model="formData.PURL" placeholder="Package-URL" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="CPE" prop="CPE">
              <el-input v-model="formData.CPE" placeholder="Common Platform Enumeration" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="SWID TAG ID" prop="SWIDTAGID">
              <el-input v-model="formData.SWIDTAGID" placeholder="SWID" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="维护部门" prop="maint_dept">
              <el-input v-model="formData.maint_dept" placeholder="请输入维护部门" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="使用部门" prop="use_dept">
              <el-input v-model="formData.use_dept" placeholder="请输入使用部门" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="适用等级" prop="applicable_level">
          <el-select v-model="formData.applicable_level" placeholder="请选择适用等级" style="width: 100%">
            <el-option label="A - 可安全使用（低风险）" value="A" />
            <el-option label="B - 谨慎使用（中低风险）" value="B" />
            <el-option label="C - 限期整改（高风险）" value="C" />
            <el-option label="D - 禁止使用（高风险）" value="D" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 项目详情弹框 -->
    <ProjectDetailDialog v-model="detailDialogVisible" :project-data="currentProject" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { Plus, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { getSoftwareList, createSoftware, updateSoftware, deleteSoftware, type SoftwareItem } from '@/api'
import ProjectDetailDialog from '@/components/business/ProjectDetailDialog.vue'
import { formatTime } from '@/utils/format'

const loading = ref(false)
const toolList = ref<SoftwareItem[]>([])
const searchText = ref('')
const filterGrade = ref('')
const filterType = ref('')

// 详情弹框
const detailDialogVisible = ref(false)
const currentProject = ref<SoftwareItem | null>(null)

function showDetail(row: SoftwareItem) {
  currentProject.value = row
  detailDialogVisible.value = true
}

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

function handleSearch() {
  pagination.page = 1
  loadToolList()
}

function handleReset() {
  searchText.value = ''
  filterGrade.value = ''
  filterType.value = ''
  handleSearch()
}

function handleSizeChange() {
  pagination.page = 1
}

function handlePageChange() {
  // pagination is reactive, no action needed
}

function typeTagColor(type: string): 'primary' | 'success' | 'warning' | 'info' {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info'> = {
    APPLICATION: 'primary',
    FRAMEWORK: 'success',
    LIBRARY: 'warning',
    CONTAINER: 'info',
    OPERATING_SYSTEM: 'primary',
    DEVICE: 'success',
    FIRMWARE: 'warning',
    FILE: 'info'
  }
  return map[type] || 'info'
}

// 加载数据
async function loadToolList() {
  loading.value = true
  try {
    const params: any = { category: 'tool' }
    if (searchText.value) params.keyword = searchText.value
    if (filterGrade.value) params.applicable_level = filterGrade.value
    if (filterType.value) params.type = filterType.value

    const res = await getSoftwareList(params)
    const list = Array.isArray(res) ? res : (res.list || res.data?.list || [])
    toolList.value = list
    pagination.total = list.length
  } catch (error) {
    ElMessage.error('加载数据失败')
    toolList.value = []
  } finally {
    loading.value = false
  }
}

// CRUD operations
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

interface ToolFormData {
  _name: string
  _version: string
  PUBLISHER: string
  CLASSIFIER: string
  type: string
  applicable_level: string
  _group: string
  license: string
  DESCRIPTION: string
  maint_dept: string
  use_dept: string
  country: string
  ACTIVE: number
  AUTHOR: string
  PURL: string
  CPE: string
  SWIDTAGID: string
  category: string
}

const defaultFormData: ToolFormData = {
  _name: '',
  _version: '',
  PUBLISHER: '',
  CLASSIFIER: '',
  type: '',
  applicable_level: '',
  _group: '',
  license: '',
  DESCRIPTION: '',
  maint_dept: '',
  use_dept: '',
  country: '',
  ACTIVE: 1,
  AUTHOR: '',
  PURL: '',
  CPE: '',
  SWIDTAGID: '',
  category: 'tool',
}

const formData = reactive<ToolFormData>({ ...defaultFormData })

const formRules: FormRules = {
  _name: [{ required: true, message: '请输入工具名称', trigger: 'blur' }],
  _version: [{ required: true, message: '请输入版本号', trigger: 'blur' }],
  CLASSIFIER: [{ required: true, message: '请选择类型', trigger: 'change' }],
  type: [{ required: true, message: '请选择开发语言', trigger: 'change' }],
  applicable_level: [{ required: true, message: '请选择适用等级', trigger: 'change' }],
  _group: [{ required: true, message: '请输入工具组织', trigger: 'blur' }],
  maint_dept: [{ required: true, message: '请输入维护部门', trigger: 'blur' }],
  use_dept: [{ required: true, message: '请输入使用部门', trigger: 'blur' }],
  country: [{ required: true, message: '请输入开源国家', trigger: 'blur' }],
  category: [{ required: true, message: '请选择软件分类', trigger: 'change' }],
}

function resetForm() {
  Object.assign(formData, defaultFormData)
  editingId.value = null
  formRef.value?.clearValidate()
}

function handleAdd() {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

function handleEdit(row: SoftwareItem) {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(formData, {
    _name: row._name,
    _version: row._version,
    PUBLISHER: row.PUBLISHER,
    CLASSIFIER: row.CLASSIFIER,
    type: row.type,
    applicable_level: row.applicable_level,
    _group: row._group,
    license: row.license,
    DESCRIPTION: row.DESCRIPTION,
    maint_dept: row.maint_dept,
    use_dept: row.use_dept,
    country: row.country,
    ACTIVE: row.ACTIVE,
    AUTHOR: row.AUTHOR,
    PURL: row.PURL,
    CPE: row.CPE,
    SWIDTAGID: row.SWIDTAGID,
    category: row.category || 'tool',
  })
  dialogVisible.value = true
}

async function handleSubmit() {
  formRef.value?.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        if (isEdit.value && editingId.value) {
          await updateSoftware(editingId.value, formData)
          ElMessage.success(`工具 "${formData._name}" 更新成功`)
        } else {
          await createSoftware(formData)
          ElMessage.success(`工具 "${formData._name}" 新增成功`)
        }
        dialogVisible.value = false
        await loadToolList()
      } catch (error) {
        ElMessage.error('操作失败')
      } finally {
        loading.value = false
      }
    }
  })
}

function handleDelete(row: SoftwareItem) {
  ElMessageBox.confirm(`确定要删除工具 "${row._name}" 吗？`, '确认删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    try {
      await deleteSoftware(row.id)
      ElMessage.success(`工具 "${row._name}" 已删除`)
      await loadToolList()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

onMounted(() => {
  loadToolList()
})
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
.expand-content {
  padding: 16px 20px;
  background-color: #fafafa;
}
.expand-item {
  padding: 6px 0;
  line-height: 1.8;
}
.expand-label {
  color: #909399;
  font-size: 13px;
  margin-right: 8px;
  display: inline-block;
  min-width: 100px;
}
.expand-value {
  color: #606266;
  font-size: 13px;
}
.expand-row {
  margin-bottom: 8px;
}
</style>
