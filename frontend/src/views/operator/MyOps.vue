<template>
  <div class="my-ops">
    <div class="page-header">
      <h2>我的运维</h2>
      <p class="page-desc">我负责维护的开源基础软件</p>
    </div>

    <el-tabs v-model="activeTab" class="ops-tabs">
      <!-- 我负责的软件 -->
      <el-tab-pane label="我负责的软件" name="software">
        <el-table :data="softwareList" border v-loading="loading">
          <el-table-column prop="grade" label="等级" width="80" align="center">
            <template #default="{ row }">
              <el-tag
                :type="gradeType(row.grade)"
                size="small"
                effect="dark"
              >
                {{ row.grade }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="name" label="软件名称" min-width="160">
            <template #default="{ row }">
              <el-link type="primary" @click="showSoftwareDetail(row)">{{ row.name }}</el-link>
            </template>
          </el-table-column>
          <el-table-column prop="version" label="当前版本" width="120" />
          <el-table-column prop="type" label="类型" width="120">
            <template #default="{ row }">
              <el-tag size="small">{{ row.type }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="owner" label="负责人" width="100" />
          <el-table-column prop="defectCount" label="缺陷数" width="100" align="center">
            <template #default="{ row }">
              <el-tag v-if="row.defectCount > 0" type="danger" size="small">{{ row.defectCount }}</el-tag>
              <span v-else>0</span>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- 缺陷管理 -->
      <el-tab-pane label="缺陷管理" name="defects">
        <!-- 查询和批量操作栏 -->
        <div style="margin-bottom: 16px; display: flex; gap: 12px; align-items: center; flex-wrap: wrap">
          <el-select
            v-model="defectFilter.softwareName"
            placeholder="按软件名称筛选"
            clearable
            style="width: 200px"
            @change="filterDefectList"
          >
            <el-option
              v-for="sw in softwareList"
              :key="sw.id"
              :label="sw.name"
              :value="sw.name"
            />
          </el-select>
          <el-select
            v-model="defectFilter.globalStatus"
            placeholder="全局状态"
            clearable
            style="width: 140px"
            @change="filterDefectList"
          >
            <el-option label="待研判" value="待研判" />
            <el-option label="待评审" value="待评审" />
            <el-option label="已分配" value="已分配" />
            <el-option label="已忽略" value="已忽略" />
          </el-select>
          <el-button type="primary" plain @click="filterDefectList">
            <el-icon><Search /></el-icon> 查询
          </el-button>
          <el-button @click="resetDefectFilter">重置</el-button>
          <el-divider direction="vertical" />
          <el-button
            type="primary"
            :disabled="selectedDefects.length === 0"
            @click="showBatchDisposeDialog"
          >
            <el-icon><EditPen /></el-icon>
            批量处置 ({{ selectedDefects.length }})
          </el-button>
        </div>

        <el-table
          :data="filteredDefectList"
          border
          v-loading="loading"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="id" label="缺陷编号" width="140" />
          <el-table-column prop="softwareName" label="软件名称" width="140" />
          <el-table-column prop="title" label="缺陷标题" min-width="200" show-overflow-tooltip />
          <el-table-column prop="severity" label="严重程度" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="severityType(row.severity)" size="small">{{ row.severity }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="source" label="来源" width="120" />
          <el-table-column prop="globalStatus" label="全局状态" width="110" align="center">
            <template #default="{ row }">
              <el-tag :type="globalStatusType(row.globalStatus)" size="small">{{ row.globalStatus }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="currentHandler" label="当前处理人" width="100" />
          <el-table-column label="操作" width="100" align="center">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="viewDefectDetail(row)">
                详情
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <!-- 批量处置对话框 -->
    <el-dialog
      v-model="batchDisposeDialogVisible"
      title="批量处置"
      width="900px"
    >
      <el-alert :title="`已选择 ${selectedDefects.length} 个缺陷`" type="info" :closable="false" show-icon style="margin-bottom: 16px" />
      
      <el-form :model="batchDisposeForm" label-width="120px" style="margin-bottom: 16px">
        <el-form-item label="目标状态" required>
          <el-select v-model="batchDisposeForm.targetStatus" placeholder="请选择目标状态" style="width: 100%">
            <el-option label="待研判" value="待研判" />
            <el-option label="待评审" value="待评审" />
            <el-option label="已分配" value="已分配" />
            <el-option label="已忽略" value="已忽略" />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="batchDisposeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmBatchDispose" :disabled="!batchDisposeForm.targetStatus">确认处置</el-button>
      </template>
    </el-dialog>

    <!-- 缺陷详情对话框 -->
    <el-dialog
      v-model="defectDetailDialogVisible"
      title="缺陷详情"
      width="800px"
    >
      <el-descriptions v-if="currentDefect" :column="2" border>
        <el-descriptions-item label="缺陷编号">{{ currentDefect.id }}</el-descriptions-item>
        <el-descriptions-item label="软件名称">{{ currentDefect.softwareName }}</el-descriptions-item>
        <el-descriptions-item label="缺陷标题" :span="2">{{ currentDefect.title }}</el-descriptions-item>
        <el-descriptions-item label="严重程度">
          <el-tag :type="severityType(currentDefect.severity)" size="small">{{ currentDefect.severity }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="来源">{{ currentDefect.source }}</el-descriptions-item>
        <el-descriptions-item label="全局状态">
          <el-tag :type="globalStatusType(currentDefect.globalStatus)" size="small">{{ currentDefect.globalStatus }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="当前处理人">{{ currentDefect.currentHandler || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!-- 软件详情对话框（复用 ProjectDetailDialog 组件，展示依赖信息） -->
    <ProjectDetailDialog v-model="softwareDetailVisible" :project-data="currentSoftware" />

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, EditPen } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import request from '@/api'
import { formatTime } from '@/utils/format'
import ProjectDetailDialog from '@/components/business/ProjectDetailDialog.vue'

const route = useRoute()
const userStore = useUserStore()
// 是否看部门全部（specialist 或 admin 在个人工作台视图）
const isSpecialist = computed(() => {
  const role = userStore.role
  return role === 'specialist' || (role === 'admin' && userStore.viewMode === 'personal')
})

const activeTab = ref('software')
const loading = ref(false)

// 我负责的软件
const softwareList = ref<any[]>([])

// 加载软件列表
async function loadSoftwareList() {
  loading.value = true
  try {
    const params: any = {}
    if (isSpecialist.value) {
      // specialist/admin(个人工作台)：看部门所有开源软件
      params.maint_dept = userStore.userInfo?.department || ''
    } else {
      // operator：只看 project_owner 表中分配给自己的软件
      params.owner_username = userStore.userInfo?.username || ''
    }
    console.log('[MyOps] loadSoftwareList:', {
      isSpecialist: isSpecialist.value,
      role: userStore.role,
      viewMode: userStore.viewMode,
      department: userStore.userInfo?.department,
      username: userStore.userInfo?.username,
      params,
    })
    const res: any = await request.get('/asset/projects', { params })
    // 响应拦截器已返回res.data，res直接是数组
    // 保留原始字段（_name, _version, id 等）供 ProjectDetailDialog 使用
    softwareList.value = (res || []).map((item: any) => ({
      ...item,
      name: item._name,
      version: item._version,
      grade: item.applicable_level || 'A',
      type: getTypeLabel(item.type),
      owner: item.create_by || '-',
      department: item.maint_dept || '-',
      description: item.DESCRIPTION || '-',
      defectCount: 0,
    }))
  } catch (error) {
    ElMessage.error('加载软件列表失败')
  } finally {
    loading.value = false
  }
}

function getTypeLabel(type: string) {
  const map: Record<string, string> = {
    java: 'Java组件',
    go: 'Go组件',
    python: 'Python组件',
    js: 'JavaScript组件',
  }
  return map[type] || type || '-'
}

// 缺陷列表
const defectList = ref<any[]>([])
const filteredDefectList = ref<any[]>([])

// 加载缺陷列表
async function loadDefectList() {
  try {
    const res: any = await request.get('/risk/defects', {
      params: { page: 1, pageSize: 100 }
    })
    // 响应拦截器已返回res.data，缺陷API返回分页数据 { list: [...], total: N }
    const list = res?.list || res || []
    defectList.value = list.map((item: any) => ({
      id: item.defect_no || `DEF-${item.id}`,
      softwareName: item.affected_component || '-',
      title: item.title || '-',
      severity: item.risk_level || '中危',
      source: item.source || '内部SCA',
      globalStatus: mapDefectStatus(item.status),
      currentHandler: item.current_handler || '-',
    }))
    filteredDefectList.value = defectList.value
  } catch (error) {
    console.error('加载缺陷列表失败:', error)
  }
}

function mapDefectStatus(status: string) {
  const map: Record<string, string> = {
    pending_analysis: '待研判',
    pending_review: '待评审',
    dispatched: '已分配',
    ignored: '已忽略',
    closed: '已忽略',
  }
  return map[status] || '待研判'
}

// 软件详情
const softwareDetailVisible = ref(false)
const currentSoftware = ref<any>(null)

function showSoftwareDetail(software: any) {
  currentSoftware.value = software
  softwareDetailVisible.value = true
}

// 缺陷查询和筛选
const defectFilter = ref({
  softwareName: '',
  globalStatus: '',
})

function filterDefectList() {
  filteredDefectList.value = defectList.value.filter(d => {
    if (defectFilter.value.softwareName && d.softwareName !== defectFilter.value.softwareName) {
      return false
    }
    if (defectFilter.value.globalStatus && d.globalStatus !== defectFilter.value.globalStatus) {
      return false
    }
    return true
  })
}

function resetDefectFilter() {
  defectFilter.value = {
    softwareName: '',
    globalStatus: '',
  }
  filteredDefectList.value = defectList.value
}

// 批量处置
const batchDisposeDialogVisible = ref(false)
const selectedDefects = ref<any[]>([])
const batchDisposeForm = ref({
  targetStatus: '',
})

function handleSelectionChange(selection: any[]) {
  selectedDefects.value = selection
}

function showBatchDisposeDialog() {
  if (selectedDefects.value.length === 0) {
    ElMessage.warning('请选择要处置的缺陷')
    return
  }
  batchDisposeForm.value = { targetStatus: '' }
  batchDisposeDialogVisible.value = true
}

function confirmBatchDispose() {
  if (!batchDisposeForm.value.targetStatus) {
    ElMessage.warning('请选择目标状态')
    return
  }
  selectedDefects.value.forEach(defect => {
    const d = defectList.value.find(item => item.id === defect.id)
    if (d) {
      d.globalStatus = batchDisposeForm.value.targetStatus
    }
  })
  filterDefectList()
  batchDisposeDialogVisible.value = false
  ElMessage.success(`成功处置 ${selectedDefects.value.length} 个缺陷`)
}

// 查看缺陷详情
const defectDetailDialogVisible = ref(false)
const currentDefect = ref<any>(null)

function viewDefectDetail(defect: any) {
  currentDefect.value = defect
  defectDetailDialogVisible.value = true
}

function globalStatusType(status: string) {
  const map: Record<string, string> = {
    '待研判': 'warning',
    '待评审': 'info',
    '已分配': 'primary',
    '已忽略': 'danger',
  }
  return map[status] || ''
}

function gradeType(grade: string) {
  const map: Record<string, any> = {
    A: 'success',
    B: '',
    C: 'warning',
    D: 'danger',
  }
  return map[grade] || 'info'
}

function severityType(severity: string) {
  const map: Record<string, any> = {
    '严重': 'danger',
    '高危': 'danger',
    '中危': 'warning',
    '低危': '',
  }
  return map[severity] || 'info'
}

onMounted(() => {
  loadSoftwareList()
  loadDefectList()
})
</script>

<style scoped>
.my-ops {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 600;
  color: #1f2328;
}

.page-desc {
  margin: 0;
  font-size: 13px;
  color: #8c8c8c;
}

.ops-tabs {
  background: #fff;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
}
</style>
