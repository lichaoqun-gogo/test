<template>
  <div class="source-backup">
    <div class="page-header">
      <h2 class="page-title">源码备份</h2>
      <div class="header-actions">
        <el-button type="primary" @click="openDownloadDialog">
          <el-icon><Download /></el-icon> 源码下载
        </el-button>
      </div>
    </div>

    <el-card shadow="never" class="browser-card">
      <!-- 面包屑路径导航 -->
      <el-breadcrumb :separator-icon="ArrowRight" class="path-breadcrumb">
        <el-breadcrumb-item
          v-for="(crumb, idx) in breadcrumbs"
          :key="idx"
          @click="navigateTo(crumb.prefix)"
        >
          <span :class="{ 'crumb-clickable': idx < breadcrumbs.length - 1 }">
            {{ crumb.name }}
          </span>
        </el-breadcrumb-item>
      </el-breadcrumb>

      <!-- 搜索栏 -->
      <div class="toolbar">
        <el-input
          v-model="searchQuery"
          placeholder="请输入对象名称前缀可以过滤"
          clearable
          style="width: 320px"
          @keyup.enter="handleSearch"
          @clear="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="handleSearch">
          <el-icon><Search /></el-icon> 查询
        </el-button>
        <el-button @click="handleRefresh">
          <el-icon><Refresh /></el-icon> 刷新
        </el-button>
      </div>

      <!-- 对象列表 -->
      <el-table
        :data="objectList"
        stripe
        v-loading="loading"
        style="width: 100%"
        empty-text="暂无数据"
      >
        <el-table-column label="名称" min-width="300">
          <template #default="{ row }">
            <div class="name-cell" @click="handleNameClick(row)">
              <el-icon class="name-icon" :color="row.type === 'folder' ? '#e6a23c' : '#409eff'">
                <Folder v-if="row.type === 'folder'" />
                <Document v-else />
              </el-icon>
              <span class="name-text" :class="{ 'name-link': row.type === 'folder' }">
                {{ row.name }}
              </span>
              <el-tag v-if="row.type === 'folder' && row.fileCount > 0" size="small" type="info" class="file-count-tag">
                {{ row.fileCount }} 个文件
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="大小" width="140">
          <template #default="{ row }">
            {{ formatFileSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column label="最后更新时间" width="200">
          <template #default="{ row }">
            {{ formatDate(row.lastModified) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" align="center">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleDownload(row)">
              <el-icon><Download /></el-icon> 下载
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 源码下载弹框 -->
    <el-dialog v-model="downloadDialogVisible" title="源码下载" width="420px" :close-on-click-modal="false">
      <el-form ref="downloadFormRef" :model="downloadForm" :rules="downloadRules" label-width="80px">
        <el-form-item label="组织" prop="orgName">
          <el-input v-model="downloadForm.orgName" placeholder="如 apache" clearable />
        </el-form-item>
        <el-form-item label="名称" prop="projectName">
          <el-input v-model="downloadForm.projectName" placeholder="如 kafka" clearable />
        </el-form-item>
        <el-form-item label="版本" prop="version">
          <el-input v-model="downloadForm.version" placeholder="如 v3.7.0" clearable />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="downloadDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="downloading" @click="handleDownloadSubmit">下载</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  Search, Refresh, Download,
  ArrowRight, Folder, Document,
} from '@element-plus/icons-vue'
import { sourceBackupApi } from '@/api'

interface ObjectItem {
  key: string
  name: string
  type: 'folder' | 'file'
  size: number
  lastModified: string
  fileCount: number
}

const ROOT_PREFIX = ''

const loading = ref(false)
const downloading = ref(false)
const objectList = ref<ObjectItem[]>([])
const currentPrefix = ref(ROOT_PREFIX)
const searchQuery = ref('')
const bucketName = ref('osm')

// 下载弹框
const downloadDialogVisible = ref(false)
const downloadFormRef = ref<FormInstance>()
const downloadForm = reactive({
  orgName: '',
  projectName: '',
  version: '',
})
const downloadRules: FormRules = {
  orgName: [{ required: true, message: '请输入组织名称', trigger: 'blur' }],
  projectName: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  version: [{ required: true, message: '请输入版本号', trigger: 'blur' }],
}

const breadcrumbs = computed(() => {
  const crumbs = [{ name: bucketName.value, prefix: '' }]
  const parts = currentPrefix.value.split('/').filter(p => p !== '')
  let acc = ''
  for (const part of parts) {
    acc = acc ? `${acc}/${part}` : part
    crumbs.push({ name: part, prefix: acc })
  }
  return crumbs
})

onMounted(() => {
  loadObjects()
})

async function loadObjects() {
  loading.value = true
  try {
    const res: any = await sourceBackupApi.listObjects({
      prefix: currentPrefix.value,
      query: searchQuery.value,
    })
    objectList.value = res.list || []
    if (res.bucket) {
      bucketName.value = res.bucket
    }
  } catch (e: any) {
    ElMessage.error(e.message || '加载对象列表失败')
    objectList.value = []
  } finally {
    loading.value = false
  }
}

function navigateTo(prefix: string) {
  currentPrefix.value = prefix || ROOT_PREFIX
  searchQuery.value = ''
  loadObjects()
}

function handleSearch() {
  loadObjects()
}

function handleRefresh() {
  searchQuery.value = ''
  loadObjects()
}

function handleNameClick(row: ObjectItem) {
  if (row.type === 'folder') {
    currentPrefix.value = row.key
    searchQuery.value = ''
    loadObjects()
  }
}

function handleDownload(row: ObjectItem) {
  if (row.type === 'folder') {
    sourceBackupApi.downloadFolder(row.key)
    ElMessage.success(`开始下载 ${row.name} (打包ZIP)`)
  } else {
    sourceBackupApi.downloadObject(row.key)
    ElMessage.success(`开始下载 ${row.name}`)
  }
}

// 打开源码下载弹框
function openDownloadDialog() {
  downloadForm.orgName = ''
  downloadForm.projectName = ''
  downloadForm.version = ''
  downloadDialogVisible.value = true
}

// 弹框提交下载
async function handleDownloadSubmit() {
  if (!downloadFormRef.value) return
  await downloadFormRef.value.validate()

  downloading.value = true
  try {
    // 构造对象路径：src/{组织}/{名称}/{版本}
    const key = `src/${downloadForm.orgName}/${downloadForm.projectName}/${downloadForm.version}`
    sourceBackupApi.downloadObject(key)
    ElMessage.success(`开始下载 ${downloadForm.projectName} ${downloadForm.version}`)
    downloadDialogVisible.value = false
  } catch (e: any) {
    ElMessage.error(e.message || '下载失败')
  } finally {
    downloading.value = false
  }
}

function formatFileSize(size: number): string {
  if (!size || size === 0) return '-'
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  if (size < 1024 * 1024 * 1024) return `${(size / (1024 * 1024)).toFixed(1)} MB`
  return `${(size / (1024 * 1024 * 1024)).toFixed(1)} GB`
}

function formatDate(dateStr: string): string {
  if (!dateStr || dateStr.startsWith('0001-')) return '-'
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

<style scoped>
.source-backup {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
}

.browser-card {
  border-radius: 8px;
}

.path-breadcrumb {
  margin-bottom: 16px;
  font-size: 14px;
}

.crumb-clickable {
  cursor: pointer;
  color: #409eff;
}

.crumb-clickable:hover {
  text-decoration: underline;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.name-cell {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 8px;
}

.name-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.name-link {
  color: #409eff;
}

.name-link:hover {
  text-decoration: underline;
}

.file-count-tag {
  margin-left: 4px;
}
</style>
