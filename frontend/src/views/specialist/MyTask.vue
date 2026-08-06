<template>
  <div class="my-task">
    <!-- 页面头部 -->
    <div class="page-header">
      <h2 class="page-title">我的任务</h2>
    </div>

    <!-- 快捷发起区 -->
    <div class="quick-action-section">
      <div class="action-buttons">
        <div class="action-btn" @click="openTaskForm('introduce')">
          <div class="icon-box introduce-box">
            <el-icon class="btn-icon" size="20"><DocumentAdd /></el-icon>
          </div>
          <span class="btn-text">引入申请</span>
        </div>
        <div class="action-btn" @click="openTaskForm('upgrade')">
          <div class="icon-box upgrade-box">
            <el-icon class="btn-icon" size="20"><Upload /></el-icon>
          </div>
          <span class="btn-text">升级申请</span>
        </div>
        <div class="action-btn" @click="openTaskForm('withdraw')">
          <div class="icon-box withdraw-box">
            <el-icon class="btn-icon" size="20"><CircleClose /></el-icon>
          </div>
          <span class="btn-text">退出申请</span>
        </div>
        <div class="action-btn" @click="openTaskForm('defect')">
          <div class="icon-box defect-box">
            <el-icon class="btn-icon" size="20"><Warning /></el-icon>
          </div>
          <span class="btn-text">缺陷上报</span>
        </div>
      </div>
    </div>

    <!-- 主内容区（查询 + Tab 统一容器） -->
    <div class="main-content-section">
      <!-- 查询区域 -->
      <div class="filter-section">
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <el-form-item label="任务类型">
            <el-select v-model="filterForm.type" placeholder="全部类型" clearable style="width: 140px" @change="applyFilter">
              <el-option v-for="item in taskTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="filterForm.status" placeholder="全部状态" clearable style="width: 120px" @change="applyFilter">
              <el-option v-for="item in taskStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-input v-model="filterForm.keyword" placeholder="搜索任务标题" clearable style="width: 200px" @input="applyFilter">
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button @click="resetFilter"><el-icon><Refresh /></el-icon> 重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- Tab 内容区 -->
      <el-tabs v-model="activeTab" class="task-tabs">
        <!-- Tab 1: 待我审批 -->
      <el-tab-pane label="待我审批" name="pending">
        <el-table :data="filteredPendingTasks" stripe>
            <el-table-column label="任务编号" width="160">
              <template #default="{ row }">
                <el-link v-if="row.status === 'draft'" type="primary" @click="editDraft(row)">
                  {{ row.taskNo }}
                </el-link>
                <el-link v-else type="primary" @click="viewDetail(row.id)">{{ row.taskNo }}</el-link>
              </template>
            </el-table-column>
            <el-table-column label="任务类型" width="120" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.type === 'defect'" type="danger" size="small">
                  {{ getTaskTypeLabel(row) }}
                </el-tag>
                <span v-else class="task-type-text">{{ getTaskTypeLabel(row) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="任务标题" min-width="250">
              <template #default="{ row }">
                <el-link v-if="row.status === 'draft'" type="primary" @click="editDraft(row)">
                  {{ row.title }}
                </el-link>
                <span v-else>{{ row.title }}</span>
              </template>
            </el-table-column>
            <el-table-column label="申请人" width="100" align="center">
              <template #default="{ row }">{{ row.applicant }}</template>
            </el-table-column>
            <el-table-column label="申请时间" width="180">
              <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link @click="openApproveDialog(row)">审批</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="filteredPendingTasks.length === 0" description="暂无待审批任务" />
        </el-tab-pane>

        <!-- Tab 2: 我发起的 -->
        <el-tab-pane label="我发起的" name="initiated">
          <el-table :data="filteredInitiatedTasks" stripe>
            <el-table-column label="任务编号" width="160">
              <template #default="{ row }">
                <el-link v-if="row.status === 'draft'" type="primary" @click="editDraft(row)">
                  {{ row.taskNo }}
                </el-link>
                <el-link v-else type="primary" @click="viewDetail(row.id)">{{ row.taskNo }}</el-link>
              </template>
            </el-table-column>
            <el-table-column label="任务类型" width="120" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.type === 'defect'" type="danger" size="small">
                  {{ getTaskTypeLabel(row) }}
                </el-tag>
                <span v-else class="task-type-text">{{ getTaskTypeLabel(row) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="任务标题" min-width="250">
              <template #default="{ row }">{{ row.title }}</template>
            </el-table-column>
            <el-table-column label="状态" width="100" align="center">
              <template #default="{ row }">
                <TaskStatusIndicator :status="row.status" />
              </template>
            </el-table-column>
            <el-table-column label="创建时间" width="180">
              <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <template v-if="row.status === 'draft' || row.status === 'withdrawn'">
                  <el-button type="primary" link @click="handleQuickSubmit(row.id)">提交</el-button>
                  <el-button type="primary" link @click="handleEdit(row.id)">编辑</el-button>
                </template>
                <!-- 审批中状态：发起人可以尝试撤回（后端会校验是否允许） -->
                <template v-else-if="row.status === 'approving' && row.applicantId === userStore.userInfo?.id">
                  <el-button type="warning" link @click="handleWithdraw(row.id)">撤回</el-button>
                </template>
                <!-- 待我审批：显示审批按钮 -->
                <template v-else-if="row.status === 'approving' && row.assignedTo === userStore.userInfo?.id">
                  <el-button type="primary" link @click="handleApprove(row.id)">审批</el-button>
                </template>
                <!-- 其他状态：显示详情 -->
                <el-button v-else type="primary" link @click="viewDetail(row.id)">详情</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="filteredInitiatedTasks.length === 0" description="暂无发起的任务" />
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 任务表单弹框 -->
    <el-dialog v-model="showTaskForm" :title="formTitle" width="800px" :close-on-click-modal="false">
      <TaskForm
        v-if="showTaskForm"
        :key="editingTaskId || 'new'"
        :task-type="selectedTaskType"
        :task-data="editingTaskData"
        @save-draft="handleSaveDraft"
        @submit="handleSubmit"
      />
    </el-dialog>

    <!-- 审批对话框（通用） -->
    <ApprovalDialog
      v-model="showApprovalDialog"
      :task-id="approvalTaskId"
      :current-user="{
        id: userStore.userInfo?.id || '',
        name: userStore.userInfo?.nickname || '',
      }"
      @approved="handleApprovalSuccess"
      @rejected="handleApprovalSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, ElDialog } from 'element-plus'
import { DocumentAdd, Upload, CircleClose, Warning, Search, Refresh } from '@element-plus/icons-vue'
import type { TaskType, TaskStatus, Task } from '@/types'
import { TASK_TYPE_LABELS } from '@/types'
import { useTaskStore } from '@/stores/task'
import { useUserStore } from '@/stores/user'
import TaskStatusIndicator from '@/components/operator/TaskStatusIndicator.vue'
import TaskForm from '@/components/operator/TaskForm.vue'
import ApprovalDialog from '@/components/business/ApprovalDialog.vue'
import { taskApi } from '@/api'
import { formatTime } from '@/utils/format'

const router = useRouter()
const taskStore = useTaskStore()
const userStore = useUserStore()

// 页面加载时从数据库获取任务
onMounted(async () => {
  console.log('[MyTask] 组件已加载')
  console.log('[MyTask] userStore.userInfo:', userStore.userInfo)
  console.log('[MyTask] activeTab:', activeTab.value)
  await taskStore.loadTasks()
  console.log('[MyTask] 任务列表已加载:', taskStore.tasks.length, '个任务')
  console.log('[MyTask] 第一个任务:', taskStore.tasks[0])
})

const activeTab = ref('pending')
const showTaskForm = ref(false)
const selectedTaskType = ref<TaskType>('introduce')
const editingTaskId = ref('')
const showDebugInfo = ref(false)

const editingTaskData = ref<any>(undefined)

// 审批对话框
const showApprovalDialog = ref(false)
const approvalTaskId = ref<number | null>(null)

// 编辑草稿时获取完整详情
async function editDraft(row: any) {
  editingTaskId.value = row.id
  selectedTaskType.value = row.type
  
  // 获取任务完整详情（包含扩展表和评估数据）
  try {
    const detail: any = await taskStore.getTaskDetail(row.id)

    // 后端返回的task_detail已经是驼峰格式，直接使用
    const taskDetail = detail.task_detail || {}
    const convertedDetail: any = {}
    
    // 引入申请字段映射（后端已返回驼峰格式）
    if (taskDetail.componentName) {
      convertedDetail.customComponent = taskDetail.componentName  // 选择组件
      convertedDetail.componentName = taskDetail.componentName
    }
    if (taskDetail.componentVersion) convertedDetail.version = taskDetail.componentVersion
    if (taskDetail.licenseType) convertedDetail.license = taskDetail.licenseType
    if (taskDetail.introduceReason) convertedDetail.introduceReason = taskDetail.introduceReason
    if (taskDetail.useScenario) convertedDetail.useScenario = taskDetail.useScenario
    if (taskDetail.projectNames) convertedDetail.projectNames = taskDetail.projectNames
    if (taskDetail.remark) convertedDetail.wikiUrl = taskDetail.remark  // 备注说明
    if (taskDetail.preEvalConclusion) convertedDetail.preConclusion = taskDetail.preEvalConclusion
    
    // 获取task主表数据
    const taskData = detail.task || {}
    
    // 合并task和转换后的task_detail
    editingTaskData.value = {
      ...taskData,
      ...convertedDetail,
      applicabilityAssessment: detail.applicability_assessment
    }
  } catch (e) {
    console.error('获取任务详情失败:', e)
    // 失败时使用基本信息
    editingTaskData.value = row
  }
  
  showTaskForm.value = true
}

const formTitle = computed(() => {
  const titles: Record<TaskType, string> = {
    introduce: '引入申请',
    upgrade: '升级申请',
    withdraw: '退出申请',
    defect: '风险上报',
    usage: '申请使用',
  }
  return editingTaskId.value ? `编辑草稿 - ${titles[selectedTaskType.value]}` : `新建${titles[selectedTaskType.value]}`
})

// 筛选表单
const filterForm = ref<{
  type: TaskType | ''
  status: TaskStatus | ''
  keyword: string
}>({
  type: '',
  status: '',
  keyword: '',
})

// 任务类型选项
const taskTypeOptions = [
  { value: 'introduce' as TaskType, label: '引入申请' },
  { value: 'upgrade' as TaskType, label: '升级申请' },
  { value: 'withdraw' as TaskType, label: '退出申请' },
  { value: 'defect' as TaskType, label: '缺陷上报' },
]

// 任务状态选项
const taskStatusOptions = [
  { value: 'draft' as TaskStatus, label: '草稿' },
  { value: 'approving' as TaskStatus, label: '审批中' },
  { value: 'withdrawn' as TaskStatus, label: '已撤回' },
  { value: 'approved' as TaskStatus, label: '已通过' },
  { value: 'rejected' as TaskStatus, label: '已驳回' },
]

// 待我审批的任务（当前用户是审批人）
const pendingTasks = computed(() => {
  const currentUserId = userStore.userInfo?.id || ''
  const currentUserName = userStore.userInfo?.nickname || ''
  
  console.log('[待我审批] 当前用户ID:', currentUserId)
  console.log('[待我审批] 当前用户名:', currentUserName)
  console.log('[待我审批] 所有任务:', taskStore.tasks)
  
  const filtered = taskStore.tasks.filter(t => {
    console.log(`[待我审批] 检查任务: ${t.taskNo}, status=${t.status}, assignedTo=${t.assignedTo}`)
    const statusMatch = t.status === 'approving'
    const assigneeMatch = t.assignedTo === currentUserId
    console.log(`  status匹配: ${statusMatch}, assignee匹配: ${assigneeMatch}`)
    return statusMatch && assigneeMatch
  })
  
  console.log('[待我审批] 过滤后:', filtered)
  return filtered
})

const sortedPendingTasks = computed(() => {
  return [...pendingTasks.value].sort((a, b) => {
    return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
  })
})

// 我发起的任务（当前用户是发起人）
const initiatedTasks = computed(() => {
  const currentUserId = userStore.userInfo?.id || ''
  
  return taskStore.tasks.filter(t => t.applicantId === currentUserId || t.applicant === userStore.userInfo?.nickname)
})

const sortedInitiatedTasks = computed(() => {
  return [...initiatedTasks.value].sort((a, b) => {
    return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
  })
})

// 筛选函数
function getTaskTypeLabel(row: any) {
  return TASK_TYPE_LABELS[(row as Task).type]
}

function filterTasks(tasks: Task[]) {
  let result = tasks

  if (filterForm.value.type) {
    result = result.filter(t => t.type === filterForm.value.type)
  }
  if (filterForm.value.status) {
    result = result.filter(t => t.status === filterForm.value.status)
  }
  if (filterForm.value.keyword) {
    const kw = filterForm.value.keyword.toLowerCase()
    result = result.filter(t => t.title.toLowerCase().includes(kw))
  }

  return result
}

// 筛选后的待我审批任务
const filteredPendingTasks = computed(() => {
  return filterTasks(sortedPendingTasks.value)
})

// 筛选后的我发起的任务
const filteredInitiatedTasks = computed(() => {
  return filterTasks(sortedInitiatedTasks.value)
})

function applyFilter() {}

function resetFilter() {
  filterForm.value = { type: '', status: '', keyword: '' }
}

function openTaskForm(type: TaskType) {
  // 新建任务，清空编辑ID
  editingTaskId.value = ''
  selectedTaskType.value = type
  showTaskForm.value = true
}

// 审批（打开审批对话框）
function handleApprove(id: string) {
  approvalTaskId.value = parseInt(id)
  showApprovalDialog.value = true
}

// 审批成功后的回调
function handleApprovalSuccess() {
  // 重新加载任务列表
  taskStore.loadTasks()
}

// 查看详情（打开审批对话框）
function viewDetail(id: string) {
  approvalTaskId.value = parseInt(id)
  showApprovalDialog.value = true
}

// 打开审批对话框
function openApproveDialog(row: any) {
  approvalTaskId.value = row.id
  showApprovalDialog.value = true
}


async function handleQuickSubmit(id: string) {
  try {
    await taskApi.submitTask(parseInt(id))
    ElMessage.success('任务已提交')
    await taskStore.loadTasks()
  } catch (e: any) {
    ElMessage.error('提交失败: ' + (e.message || '未知错误'))
  }
}

// 编辑（通用）
async function handleEdit(id: string) {
  const task = taskStore.tasks.find(t => t.id === id)
  if (task) {
    await editDraft(task)
  }
}

// 撤回
async function handleWithdraw(id: number) {
  console.log('[撤回任务] 尝试撤回任务ID:', id)
  
  try {
    // 获取审批节点列表
    const nodes: any = await taskApi.getApprovalNodes(id)
    console.log('[撤回任务] 审批节点:', nodes)

    // 找到第一个未跳过的审批节点（Node 1，node_order > 0）
    const firstApprovalNode = nodes?.find((n: any) => n.nodeOrder > 0 && !n.isSkipped)
    
    if (!firstApprovalNode) {
      ElMessage.warning('未找到可撤回的审批节点')
      return
    }
    
    console.log('[撤回任务] 第一个审批节点:', {
      nodeName: firstApprovalNode.nodeName,
      status: firstApprovalNode.status,
      approver: firstApprovalNode.approver,
    })
    
    // 如果第一个审批节点还没有被处理（状态是pending），则可以撤回
    if (firstApprovalNode.status !== 'pending') {
      console.log('[撤回任务] 无法撤回：审批节点状态不是pending')
      ElMessage.warning('审批人已处理，无法撤回')
      return
    }
    
    console.log('[撤回任务] 允许撤回，显示确认框')
    await ElMessageBox.confirm(
      '确认撤回该任务？撤回后可以重新编辑并提交。',
      '确认撤回',
      { type: 'warning' }
    )
    
    console.log('[撤回任务] 用户确认撤回，调用API')
    await taskApi.withdrawTask(id)
    console.log('[撤回任务] 撤回成功')
    ElMessage.success('任务已撤回')
    await taskStore.loadTasks()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error('[撤回任务] 失败:', e)
      ElMessage.error('撤回失败: ' + (e.message || '未知错误'))
    }
  }
}

// 测试：提交最新任务
async function testSubmitLatestTask() {
  try {
    ElMessage.info('正在提交任务...')
    const res = await taskApi.submitTask(1783274627)
    console.log('提交结果:', res)
    
    ElMessage.success('提交成功，重新加载任务列表...')
    await taskStore.loadTasks()
    
    // 查找最新任务
    const task = taskStore.tasks.find(t => t.taskNo === 'TASK-1783274627')
    console.log('提交后的任务:', task)
    console.log('assignedTo:', task?.assignedTo)
    console.log('status:', task?.status)
    
    if (task?.assignedTo === '李四') {
      ElMessage.success('任务已分配给李四！')
    } else {
      ElMessage.warning(`assignedTo=${task?.assignedTo}，不是李四`)
    }
  } catch (e: any) {
    console.error('提交失败:', e)
    ElMessage.error('提交失败: ' + (e.message || '未知错误'))
  }
}

// 构建引入申请详情
function buildIntroduceDetail(data: any) {
  return {
    component_name: data.componentName || data.customComponent,
    component_version: data.version,
    license_type: data.license,
    introduce_reason: data.introduceReason,
    use_scenario: data.useScenario,
    project_names: data.projectNames,
    remark: data.wikiUrl,
    pre_eval_conclusion: data.preConclusion,
  }
}

// 构建升级申请详情
function buildUpgradeDetail(data: any) {
  return {
    component_name: data.componentName,
    current_version: data.currentVersion,
    target_version: data.targetVersion,
    upgrade_reason: data.upgradeReason,
    test_report_link: data.testReportLink,
    project_dept: data.projectDept,
  }
}

// 构建退出申请详情
function buildWithdrawDetail(data: any) {
  return {
    component_name: data.exitComponent,
    component_version: data.exitVersion,
    withdraw_reason: data.withdrawReason,
    exit_reason: data.withdrawReason,
    replacement_component: data.replacementComponent,
    alternative_plan: data.alternativePlan,
    code_cleanup_plan: data.codeCleanupPlan,
    affected_projects: data.affectedProjects,
  }
}

// 构建风险上报详情
function buildDefectDetail(data: any) {
  return {
    risk_description: data.riskDescription,
  }
}

async function handleSaveDraft(data: any) {
  try {
    const requestData: any = {
      title: data.title,
      applicant: data.applicant,
      deadline: data.deadline,
    }

    // 根据任务类型构建 task_detail
    switch (data.type) {
      case 'introduce':
        requestData.task_detail = buildIntroduceDetail(data)
        requestData.applicability_assessment = data.applicabilityAssessment || data.assessResults
        break
      case 'upgrade':
        requestData.task_detail = buildUpgradeDetail(data)
        break
      case 'withdraw':
        requestData.task_detail = buildWithdrawDetail(data)
        break
      case 'defect':
        requestData.task_detail = buildDefectDetail(data)
        break
    }

    if (editingTaskId.value) {
      await taskStore.updateTask(editingTaskId.value, requestData)
      ElMessage.success('草稿已更新')
    } else {
      requestData.type = data.type
      requestData.status = 'draft'
      await taskStore.createTask(requestData)
      ElMessage.success('草稿已保存')
    }
    showTaskForm.value = false
    editingTaskId.value = ''
    editingTaskData.value = undefined
    await taskStore.loadTasks()
  } catch (e: any) {
    ElMessage.error('保存失败: ' + (e.message || '未知错误'))
  }
}

async function handleSubmit(data: any) {
  try {
    const requestData: any = {
      title: data.title,
      applicant: data.applicant,
      deadline: data.deadline,
    }

    // 根据任务类型构建 task_detail
    switch (data.type) {
      case 'introduce':
        requestData.task_detail = buildIntroduceDetail(data)
        requestData.applicability_assessment = data.applicabilityAssessment || data.assessResults
        break
      case 'upgrade':
        requestData.task_detail = buildUpgradeDetail(data)
        break
      case 'withdraw':
        requestData.task_detail = buildWithdrawDetail(data)
        break
      case 'defect':
        requestData.task_detail = buildDefectDetail(data)
        break
    }

    if (editingTaskId.value) {
      // 更新现有任务
      await taskStore.updateTask(editingTaskId.value, requestData)
      // 调用 workflow 提交 API
      await taskApi.submitTask(parseInt(editingTaskId.value))
      ElMessage.success('任务已提交')
    } else {
      // 创建新任务：先创建草稿，再提交
      requestData.type = data.type
      requestData.status = 'draft'
      const newTask: any = await taskStore.createTask(requestData)
      // 使用后端返回的任务 ID（数据库主键）
      const taskId = newTask?.id || newTask?.data?.id
      if (taskId) {
        await taskApi.submitTask(taskId)
        ElMessage.success('任务已提交')
      } else {
        ElMessage.error('创建任务失败：未获取到任务ID')
      }
    }
    showTaskForm.value = false
    editingTaskId.value = ''
    editingTaskData.value = undefined
    // 刷新任务列表
    await taskStore.loadTasks()
  } catch (e: any) {
    ElMessage.error('提交失败: ' + (e.message || '未知错误'))
  }
}
</script>

<style scoped>
.my-task {
  padding: 0;
}

.page-header {
  margin-bottom: 20px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0;
}

/* 快捷发起区样式 */
.quick-action-section {
  background: #f5f7fa;
  border-radius: 10px;
  padding: 16px 20px;
  margin-bottom: 20px;
}

.section-title {
  font-size: 12px;
  font-weight: 600;
  color: #909399;
  margin: 0 20px 0 0;
  white-space: nowrap;
  line-height: 1;
}

.action-buttons {
  display: flex;
  gap: 12px;
  flex: 1;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  background: #ffffff;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  flex: 1;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.action-btn:hover {
  border-color: #409EFF;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.15);
  transform: translateY(-1px);
}

.icon-box {
  width: 36px;
  height: 36px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.introduce-box {
  background: #ecf5ff;
}

.upgrade-box {
  background: #f0f9ff;
}

.withdraw-box {
  background: #fdf6ec;
}

.defect-box {
  background: #fef0f0;
}

.btn-icon {
  flex-shrink: 0;
}

.introduce-box .btn-icon {
  color: #409EFF;
}

.upgrade-box .btn-icon {
  color: #67C23A;
}

.withdraw-box .btn-icon {
  color: #E6A23C;
}

.defect-box .btn-icon {
  color: #F56C6C;
}

.btn-text {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
  white-space: nowrap;
}

/* 主内容区（查询 + Tab 统一容器） */
.main-content-section {
  background: #ffffff;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #e8eaed;
}

/* 查询区域 */
.filter-section {
  padding: 16px 24px 12px;
}

.filter-form {
  margin: 0;
}

/* Tab 样式 */
.task-tabs {
  padding: 0 24px;
}

.task-tabs :deep(.el-tabs__header) {
  margin: 0;
  border-bottom: none;
}

.task-tabs :deep(.el-tabs__item) {
  font-size: 14px;
  font-weight: 500;
  color: #606266;
  padding: 16px 24px;
}

.task-tabs :deep(.el-tabs__item.is-active) {
  color: #409EFF;
  font-weight: 600;
}

.task-tabs :deep(.el-tabs__content) {
  padding: 20px 0;
}

.task-type-text {
  color: #606266;
}
</style>
