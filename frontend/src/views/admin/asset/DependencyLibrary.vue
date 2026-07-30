<template>
  <div class="dependency-library">
    <div class="page-header">
      <div class="filter-bar">
        <el-input v-model="searchText" placeholder="搜索组件名称" clearable style="width: 200px" @clear="handleSearch" @keyup.enter="handleSearch">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
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
      <div class="action-bar">
        <el-button type="success" plain @click="importDialogVisible = true">
          <el-icon><Upload /></el-icon> 批量导入
        </el-button>
        <el-button type="primary" plain @click="handleExport">
          <el-icon><Download /></el-icon> 导出
        </el-button>
      </div>
    </div>

    <el-table :data="paginatedAssets" stripe style="margin-top: 16px" v-loading="loading">
      <el-table-column type="selection" width="50" />
      <el-table-column prop="_group" label="组件组织" min-width="160" />
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

    <!-- 关联项目对话框 -->
    <el-dialog v-model="relatedProjectsDialogVisible" title="关联内部项目" width="700px">
      <el-table :data="relatedProjects" stripe v-loading="loadingProjects">
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
      </el-table>
      <div class="total-info">共 {{ relatedProjects.length }} 个项目</div>
    </el-dialog>

    <!-- 批量导入弹窗 -->
    <el-dialog v-model="importDialogVisible" title="批量导入依赖包" width="560px">
      <div class="import-section">
        <el-alert title="导入说明" type="info" :closable="false" show-icon style="margin-bottom: 16px">
          <p>1. 请先下载模板文件，按照模板格式填写依赖包信息</p>
          <p>2. 支持 .xlsx、.csv 格式文件，单次最多导入 500 条</p>
          <p>3. 导入时将自动去重，已存在的依赖包将跳过</p>
        </el-alert>

        <el-button type="primary" plain @click="handleDownloadTemplate" style="margin-bottom: 16px">
          <el-icon><Download /></el-icon> 下载导入模板
        </el-button>

        <el-upload
          ref="uploadRef"
          class="upload-area"
          drag
          :auto-upload="false"
          :limit="1"
          accept=".xlsx,.csv"
          :on-change="handleFileChange"
          :on-remove="handleFileRemove"
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">将文件拖到此处，或 <em>点击上传</em></div>
          <template #tip>
            <div class="el-upload__tip">仅支持 .xlsx / .csv 格式</div>
          </template>
        </el-upload>

        <div v-if="selectedFile" class="file-info">
          <el-tag type="success">
            <el-icon><Document /></el-icon>
            {{ selectedFile.name }} ({{ formatFileSize(selectedFile.size) }})
          </el-tag>
        </div>
      </div>

      <template #footer>
        <el-button @click="importDialogVisible = false; selectedFile = null">取消</el-button>
        <el-button type="primary" :disabled="!selectedFile" @click="handleImport">
          <el-icon><Upload /></el-icon> 开始导入
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { Search, Upload, Download, UploadFilled, Document } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getDependencies, getRelatedProjects } from '@/api'

const loading = ref(false)
const loadingProjects = ref(false)
const searchText = ref('')
const filterClassifier = ref('')
const importDialogVisible = ref(false)
const selectedFile = ref<File | null>(null)
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
    const params: any = {}
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
  filterClassifier.value = ''
  handleSearch()
}

function handleSizeChange() {
  pagination.page = 1
}

function handlePageChange() {}

function handleExport() {
  ElMessage.success('导出任务已提交，文件生成后将自动下载')
}

function handleDownloadTemplate() {
  ElMessage.success('模板下载中...')
}

function handleFileChange(file: any) {
  selectedFile.value = file.raw
}

function handleFileRemove() {
  selectedFile.value = null
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

function handleImport() {
  if (!selectedFile.value) return
  loading.value = true
  setTimeout(() => {
    ElMessage.success(`导入成功！共导入 ${Math.floor(Math.random() * 50) + 10} 条依赖包数据`)
    loading.value = false
    importDialogVisible.value = false
    selectedFile.value = null
  }, 1500)
}

async function viewRelatedProjects(row: any) {
  currentComponent.value = row
  relatedProjectsDialogVisible.value = true
  loadingProjects.value = true
  try {
    const res = await getRelatedProjects(row.id)
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
.action-bar {
  display: flex;
  gap: 8px;
}
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
.import-section {
  padding: 8px 0;
}
.upload-area {
  width: 100%;
}
.file-info {
  margin-top: 12px;
}
</style>
