<template>
  <div class="task-list">
    <div class="page-header">
      <div class="header-left">
        <h3>任务列表</h3>
        <p class="desc">全局缺陷任务视图 + 自己需要复核的缺陷任务，集中查询缺陷风险任务状态、修复进度及逾期情况</p>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="16" style="margin-top: 16px; margin-bottom: 20px">
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="任务总数" :value="tasks.length" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="待处理" :value="tasks.filter(t => t.status === '待处理').length" value-style="color: #E6A23C" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="逾期任务" :value="tasks.filter(t => t.isOverdue).length" value-style="color: #F56C6C" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="待我复核" :value="tasks.filter(t => t.isMyReview).length" value-style="color: #409EFF" />
        </el-card>
      </el-col>
    </el-row>

    <!-- Tab 切换 -->
    <el-tabs v-model="activeTab" type="card">
      <el-tab-pane label="全局任务视图" name="global">
        <div class="toolbar">
          <el-input v-model="searchText" placeholder="搜索任务编号或标题" clearable style="width: 240px" />
          <el-select v-model="filterRiskLevel" placeholder="风险等级" clearable style="width: 130px">
            <el-option label="严重" value="critical" />
            <el-option label="高风险" value="high" />
            <el-option label="中风险" value="medium" />
            <el-option label="低风险" value="low" />
          </el-select>
          <el-select v-model="filterStatus" placeholder="状态" clearable style="width: 120px">
            <el-option label="待处理" value="待处理" />
            <el-option label="处理中" value="处理中" />
            <el-option label="待复核" value="待复核" />
            <el-option label="已完成" value="已完成" />
            <el-option label="已逾期" value="已逾期" />
            <el-option label="已忽略" value="已忽略" />
          </el-select>
          <el-select v-model="filterSource" placeholder="来源" clearable style="width: 130px">
            <el-option label="外部推送" value="外部推送" />
            <el-option label="SCA扫描" value="SCA扫描" />
            <el-option label="人工上报" value="人工上报" />
          </el-select>
        </div>

        <el-table :data="filteredTasks" stripe style="margin-top: 12px">
          <el-table-column prop="taskNo" label="任务编号" width="150" fixed />
          <el-table-column prop="title" label="任务标题" min-width="200" show-overflow-tooltip />
          <el-table-column label="风险等级" width="100">
            <template #default="{ row }"><RiskTag :level="row.riskLevel" /></template>
          </el-table-column>
          <el-table-column prop="riskType" label="风险类型" width="110" />
          <el-table-column prop="source" label="来源" width="100">
            <template #default="{ row }">
              <el-tag :type="sourceType(row.source)" size="small">{{ row.source }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="assignee" label="负责人" width="90" />
          <el-table-column prop="department" label="部门" width="110" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="taskStatusType(row.status)" size="small">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="进度" width="140">
            <template #default="{ row }">
              <el-progress :percentage="row.progress" :color="progressColor(row.progress)" />
            </template>
          </el-table-column>
          <el-table-column prop="dueDate" label="截止日期" width="110">
            <template #default="{ row }">{{ formatTime(row.dueDate) }}</template>
          </el-table-column>
          <el-table-column label="逾期" width="70">
            <template #default="{ row }">
              <el-tag v-if="row.isOverdue" type="danger" size="small">逾期</el-tag>
              <span v-else class="text-ok">正常</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="viewDetail(row as RiskTask)">详情</el-button>
              <el-button v-if="(row as RiskTask).jiraKey" type="info" link size="small" @click="openJira((row as RiskTask).jiraKey!)">Jira</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="待我复核" name="my-review">
        <el-alert
          title="以下是需要您复核的缺陷任务，请确认修复结果是否符合安全标准"
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />

        <el-table :data="myReviewTasks" stripe>
          <el-table-column prop="taskNo" label="任务编号" width="150" fixed />
          <el-table-column prop="title" label="任务标题" min-width="200" show-overflow-tooltip />
          <el-table-column label="风险等级" width="100">
            <template #default="{ row }"><RiskTag :level="row.riskLevel" /></template>
          </el-table-column>
          <el-table-column prop="riskType" label="风险类型" width="110" />
          <el-table-column prop="assignee" label="修复人" width="90" />
          <el-table-column prop="department" label="部门" width="110" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="taskStatusType(row.status)" size="small">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="进度" width="140">
            <template #default="{ row }">
              <el-progress :percentage="row.progress" :color="progressColor(row.progress)" />
            </template>
          </el-table-column>
          <el-table-column prop="dueDate" label="截止日期" width="110">
            <template #default="{ row }">{{ formatTime(row.dueDate) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="viewDetail(row as RiskTask)">查看详情</el-button>
              <el-button type="success" link size="small" @click="handlePass(row as RiskTask)">通过</el-button>
              <el-button type="danger" link size="small" @click="handleReject(row as RiskTask)">驳回</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <!-- 任务详情弹窗 -->
    <el-dialog v-model="detailVisible" title="任务详情" width="640px">
      <el-descriptions :column="2" border v-if="currentTask">
        <el-descriptions-item label="任务编号">{{ currentTask.taskNo }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="taskStatusType(currentTask.status)" size="small">{{ currentTask.status }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="任务标题" :span="2">{{ currentTask.title }}</el-descriptions-item>
        <el-descriptions-item label="风险等级">
          <RiskTag :level="currentTask.riskLevel" />
        </el-descriptions-item>
        <el-descriptions-item label="风险类型">{{ currentTask.riskType }}</el-descriptions-item>
        <el-descriptions-item label="来源">{{ currentTask.source }}</el-descriptions-item>
        <el-descriptions-item label="负责人">{{ currentTask.assignee }}</el-descriptions-item>
        <el-descriptions-item label="部门">{{ currentTask.department }}</el-descriptions-item>
        <el-descriptions-item label="截止日期">{{ formatTime(currentTask.dueDate) }}</el-descriptions-item>
        <el-descriptions-item label="逾期情况">
          <el-tag v-if="currentTask.isOverdue" type="danger" size="small">已逾期</el-tag>
          <el-tag v-else type="success" size="small">正常</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="修复进度" :span="2">
          <el-progress :percentage="currentTask.progress" :color="progressColor(currentTask.progress)" style="width: 300px" />
        </el-descriptions-item>
        <el-descriptions-item v-if="currentTask.jiraKey" label="Jira工单" :span="2">
          <el-tag type="info" size="small">{{ currentTask.jiraKey }}</el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { mockRiskTasks } from '@/mock/admin'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { RiskTask } from '@/mock/admin'
import { formatTime } from '@/utils/format'

const activeTab = ref('global')
const searchText = ref('')
const filterRiskLevel = ref('')
const filterStatus = ref('')
const filterSource = ref('')

const tasks = ref<RiskTask[]>([...mockRiskTasks])

const detailVisible = ref(false)
const currentTask = ref<RiskTask | null>(null)

const filteredTasks = computed(() => {
  return tasks.value.filter(item => {
    const matchSearch = !searchText.value || item.taskNo.toLowerCase().includes(searchText.value.toLowerCase()) || item.title.toLowerCase().includes(searchText.value.toLowerCase())
    const matchRisk = !filterRiskLevel.value || item.riskLevel === filterRiskLevel.value
    const matchStatus = !filterStatus.value || item.status === filterStatus.value
    const matchSource = !filterSource.value || item.source === filterSource.value
    return matchSearch && matchRisk && matchStatus && matchSource
  })
})

const myReviewTasks = computed(() => {
  return tasks.value.filter(t => t.isMyReview)
})

function sourceType(source: string): 'danger' | 'warning' | 'info' {
  const map: Record<string, 'danger' | 'warning' | 'info'> = { '外部推送': 'danger', 'SCA扫描': 'warning', '人工上报': 'info' }
  return map[source] || 'info'
}

function taskStatusType(status: string): 'warning' | 'info' | 'success' | 'danger' | undefined {
  const map: Record<string, 'warning' | 'info' | 'success' | 'danger' | undefined> = {
    '待处理': 'warning',
    '处理中': undefined,
    '待复核': 'info',
    '已完成': 'success',
    '已逾期': 'danger',
    '已忽略': 'info',
  }
  return map[status]
}

function progressColor(progress: number) {
  if (progress === 100) return '#67C23A'
  if (progress > 60) return '#409EFF'
  return '#E6A23C'
}

function viewDetail(row: RiskTask) {
  currentTask.value = row
  detailVisible.value = true
}

function openJira(key: string) {
  ElMessage.info(`即将打开 Jira 工单: ${key}`)
}

function handlePass(row: RiskTask) {
  ElMessageBox.confirm(`确认通过「${row.title}」的修复结果？`, '复核确认', { type: 'warning' }).then(() => {
    row.status = '已完成'
    row.progress = 100
    ElMessage.success('已通过')
  }).catch(() => {})
}

function handleReject(row: RiskTask) {
  ElMessageBox.confirm(`确认驳回「${row.title}」的修复结果？将退回给修复人重新处理。`, '复核确认', { type: 'warning' }).then(() => {
    row.status = '处理中'
    row.progress = Math.max(row.progress - 20, 0)
    ElMessage.warning('已驳回，退回重新处理')
  }).catch(() => {})
}
</script>

<style scoped>
.page-header { margin-bottom: 8px; }
.page-header h3 { margin: 0 0 4px; }
.desc { color: #909399; font-size: 13px; margin: 0; }
.toolbar { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.text-ok { color: #67C23A; font-size: 12px; }
</style>
