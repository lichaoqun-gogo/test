<template>
  <div class="task-list">
    <div class="page-header">
      <h2 class="page-title">任务列表</h2>
    </div>

    <!-- 统计卡片看板 -->
    <el-row :gutter="16" style="margin-top: 16px; margin-bottom: 20px">
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <el-statistic title="全部任务" :value="taskStore.tasks.length" />
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <el-statistic title="草稿" :value="taskStore.tasks.filter(t => t.status === 'draft').length" value-style="color: #909399" />
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <el-statistic title="待审批" :value="taskStore.tasks.filter(t => t.status === 'pending').length" value-style="color: #E6A23C" />
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <el-statistic title="处理中" :value="taskStore.tasks.filter(t => t.status === 'processing').length" value-style="color: #409EFF" />
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <el-statistic title="已完成" :value="taskStore.tasks.filter(t => t.status === 'completed').length" value-style="color: #67C23A" />
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <el-statistic title="已驳回" :value="taskStore.tasks.filter(t => t.status === 'rejected').length" value-style="color: #F56C6C" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-input v-model="searchText" placeholder="搜索任务标题/编号" clearable style="width: 200px" @clear="handleSearch" @keyup.enter="handleSearch">
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>
      <el-select v-model="filterType" placeholder="任务类型" clearable style="width: 130px; margin-left: 12px" @change="handleSearch">
        <el-option label="引入申请" value="introduce" />
        <el-option label="升级申请" value="upgrade" />
        <el-option label="退出申请" value="exit" />
        <el-option label="缺陷处置" value="defect" />
      </el-select>
      <el-select v-model="filterStatus" placeholder="任务状态" clearable style="width: 130px; margin-left: 12px" @change="handleSearch">
        <el-option label="草稿" value="draft" />
        <el-option label="待审批" value="pending" />
        <el-option label="处理中" value="processing" />
        <el-option label="已完成" value="completed" />
        <el-option label="已驳回" value="rejected" />
      </el-select>
      <el-button type="primary" plain @click="handleSearch" style="margin-left: 12px">
        <el-icon><Search /></el-icon> 查询
      </el-button>
      <el-button @click="handleReset">重置</el-button>
    </div>

    <!-- 任务列表 -->
    <el-table :data="filteredTasks" stripe style="margin-top: 16px">
      <el-table-column prop="taskNo" label="任务编号" width="160">
        <template #default="{ row }">
          <el-link type="primary" @click="viewTask(row.id)">{{ row.taskNo }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="任务标题" min-width="200" />
      <el-table-column label="任务类型" width="120">
        <template #default="{ row }">
          <el-tag :type="taskTypeColor(row.type)" size="small">{{ TASK_TYPE_LABELS[row.type] }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="applicant" label="申请人" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="taskStatusColor(row.status)" size="small">{{ taskStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createTime" label="创建时间" width="160">
        <template #default="{ row }">
          {{ formatTime(row.createTime || row.createdAt) }}
        </template>
      </el-table-column>
    </el-table>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { TASK_TYPE_LABELS } from '@/types'
import type { TaskType } from '@/types'
import { formatTime } from '@/utils/format'
import { useTaskStore } from '@/stores/task'

const router = useRouter()
const taskStore = useTaskStore()

// 页面加载时从后端获取任务数据
onMounted(async () => {
  await taskStore.loadTasks()
  console.log('[TaskList] 任务数据:', taskStore.tasks[0])
})

const searchText = ref('')
const filterType = ref('')
const filterStatus = ref('')

// 使用taskStore中的任务数据
const tasks = computed(() => taskStore.tasks)

const filteredTasks = computed(() => {
  let list = [...tasks.value]
  if (searchText.value) {
    const q = searchText.value.toLowerCase()
    list = list.filter(t => t.title.toLowerCase().includes(q) || t.taskNo.toLowerCase().includes(q))
  }
  if (filterType.value) list = list.filter(t => t.type === filterType.value)
  if (filterStatus.value) list = list.filter(t => t.status === filterStatus.value)
  return list
})

function taskTypeColor(type: TaskType) {
  const colors: Record<TaskType, string> = {
    introduce: 'primary',
    upgrade: 'warning',
    exit: 'info',
    defect: 'danger',
    risk: 'danger',
  }
  return colors[type]
}

function taskStatusColor(status: string) {
  const colors: Record<string, string> = {
    draft: 'info',
    pending: 'warning',
    processing: 'primary',
    assigned: 'warning',
    completed: 'success',
    rejected: 'danger',
    approving: 'warning',
    withdrawn: 'info',
    approved: 'success',
  }
  return colors[status] || 'info'
}

function taskStatusLabel(status: string) {
  const labels: Record<string, string> = {
    draft: '草稿',
    pending: '待审批',
    processing: '处理中',
    assigned: '已分配',
    completed: '已完成',
    rejected: '已驳回',
    approving: '审批中',
    withdrawn: '已撤回',
    approved: '已通过',
  }
  return labels[status] || status
}

function handleSearch() {}

function handleReset() {
  searchText.value = ''
  filterType.value = ''
  filterStatus.value = ''
}

function viewTask(id: string) {
  router.push(`/admin/governance/task-list/${id}`)
}


</script>

<style scoped>
.task-list {
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
</style>
