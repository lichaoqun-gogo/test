<template>
  <el-dialog
    v-model="visible"
    :title="`${taskData?.taskNo || ''} - ${taskData ? getStatusLabel(taskData.status) : '任务审批'}`"
    width="900px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div v-if="taskData" class="approval-detail">
      <!-- 任务基本信息 -->
      <el-descriptions :column="2" border class="mb-4">
        <el-descriptions-item label="任务标题" :span="2">{{ taskData.title || '-' }}</el-descriptions-item>
        <el-descriptions-item label="组件名称">{{ taskData.componentName || taskData.customComponent || '-' }}</el-descriptions-item>
        <el-descriptions-item label="版本号">{{ taskData.componentVersion || '-' }}</el-descriptions-item>
        <el-descriptions-item label="期望完成时间">{{ taskData.deadline || '-' }}</el-descriptions-item>
        <el-descriptions-item label="任务类型">
          <el-tag :type="getTaskTypeColor(taskData.type) as any">{{ getTaskTypeLabel(taskData.type) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="申请人">{{ taskData.applicant }}</el-descriptions-item>
        <el-descriptions-item label="申请部门">{{ taskData.applicantDept }}</el-descriptions-item>
        <el-descriptions-item v-if="taskData.type === 'introduce'" label="引入原因" :span="2">
          <div style="min-height: 80px; white-space: pre-wrap; line-height: 1.8">
            {{ taskData.introduceReason || '-' }}
          </div>
        </el-descriptions-item>
        <el-descriptions-item v-if="taskData.type === 'withdraw'" label="退出原因" :span="2">
          <div style="min-height: 80px; white-space: pre-wrap; line-height: 1.8">
            {{ taskData.withdrawReason || taskData.exitReason || '-' }}
          </div>
        </el-descriptions-item>
      </el-descriptions>

      <!-- 引入申请表单 -->
      <template v-if="taskData.type === 'introduce'">

        <!-- 适用性评估 -->
        <el-divider content-position="left">适用性评估</el-divider>
        <el-descriptions :column="2" border class="mb-4 assess-descriptions">
          <el-descriptions-item v-for="item in assessItems" :key="item.key" :label="item.label">
            <div>
              <el-tag :type="getAssessStatus(item.key) === 'fully_meet' ? 'success' : 'warning'" size="small">
                {{ getAssessStatus(item.key) === 'fully_meet' ? '完全满足' : getAssessStatus(item.key) === 'deviation' ? '存在偏差' : '未评估' }}
              </el-tag>
              <el-tooltip v-if="getAssessDeviation(item.key)" :content="getAssessDeviation(item.key)" placement="top">
                <div class="deviation-text">
                  <strong>偏差说明：</strong>{{ getAssessDeviation(item.key) }}
                </div>
              </el-tooltip>
            </div>
          </el-descriptions-item>
        </el-descriptions>

        <!-- 备注说明 -->
        <el-descriptions :column="2" border class="mb-4 remark-descriptions">
          <el-descriptions-item label="备注说明" :span="2">
            <div v-if="taskData.remark" style="white-space: pre-wrap">{{ taskData.remark }}</div>
            <div v-else style="color: #909399">-</div>
          </el-descriptions-item>
        </el-descriptions>
      </template>

      <!-- 审批流程（预设流程节点） -->
      <el-divider content-position="left">审批流程</el-divider>
      <div v-if="activeFlowNodes.length > 0" class="approval-flow">
        <div
          v-for="(node, index) in activeFlowNodes"
          :key="node.id || index"
          class="flow-node-wrapper"
        >
          <div class="flow-node" :class="getFlowNodeClass(node)">
            <div class="flow-node-header">
              <span class="flow-node-name">{{ node.nodeName }}</span>
              <el-tag
                :type="getFlowNodeTagType(node)"
                size="small"
                class="flow-node-tag"
              >
                {{ getFlowNodeStatusText(node) }}
              </el-tag>
            </div>
            <div class="flow-node-info" :class="{ 'text-gray': node.status === 'pending' }">
              {{ node.approverDept || '-' }}-{{ node.approver || '-' }}
            </div>
          </div>
          <div v-if="index < activeFlowNodes.length - 1" class="flow-arrow">
            <el-icon :size="18"><ArrowRight /></el-icon>
          </div>
        </div>
      </div>
      <el-empty v-else description="暂无审批流程" :image-size="60" />

      <!-- 审批历史 -->
      <el-divider content-position="left">审批历史</el-divider>
      <el-timeline v-if="deduplicatedHistory.length > 0">
        <el-timeline-item
          v-for="(item, index) in deduplicatedHistory"
          :key="index"
          :timestamp="formatTime(item.createTime || item.approvedAt)"
          placement="top"
          :type="getActionColor(item.action) as any"
        >
          <el-card>
            <p><strong>{{ item.nodeName }}-{{ item.approverDept }}-{{ item.approver }}</strong></p>
            <p>操作：{{ getActionLabel(item.action) }}</p>
            <p v-if="item.opinion">意见：{{ item.opinion }}</p>
          </el-card>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-else description="暂无审批历史" :image-size="60" />

      <!-- 审批操作区（仅当前审批人可见） -->
      <template v-if="canApprove">
        <el-divider content-position="left">审批操作</el-divider>
        <el-form :model="approvalForm" label-width="100px">
          <el-form-item label="审批意见" required>
            <el-input
              v-model="approvalForm.opinion"
              type="textarea"
              :rows="4"
              placeholder="请输入审批意见（驳回时必填）"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>
        </el-form>
      </template>
    </div>

    <template #footer>
      <el-button @click="handleClose">关闭</el-button>
      <template v-if="canApprove">
        <el-button type="danger" @click="handleReject">驳回</el-button>
        <el-button type="primary" @click="handleApprove">通过</el-button>
      </template>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download, ArrowRight } from '@element-plus/icons-vue'
import { TASK_TYPE_LABELS } from '@/types'
import type { TaskType } from '@/types'
import { taskApi } from '@/api'
import { formatTime } from '@/utils/format'

interface Props {
  modelValue: boolean
  taskId: number | null
  currentUser: {
    id: string
    name: string
  }
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'approved'): void
  (e: 'rejected'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

const taskData = ref<any>(null)
const approvalNodes = ref<any[]>([])
const approvalHistory = ref<any[]>([])

const dialogTitle = computed(() => {
  if (!taskData.value) return '任务审批'
  return `${taskData.value.taskNo} - 任务审批`
})

// 审批历史直接使用后端返回数据（后端已过滤为最后一次提交开始）
const deduplicatedHistory = computed(() => {
  return approvalHistory.value
})

// 审批流程节点（过滤掉节点0发起节点，只显示实际审批节点）
const activeFlowNodes = computed(() => {
  return approvalNodes.value.filter(node => node.nodeOrder > 0 && !node.isSkipped)
})

// 判断当前用户是否可以审批
const canApprove = computed(() => {
  if (!taskData.value || !props.currentUser.id) return false
  // 任务必须是审批中状态
  if (taskData.value.status !== 'approving') return false
  // 当前用户必须是处理人（assignedTo 存的是用户ID）
  return taskData.value.assignedTo === props.currentUser.id
})

const approvalForm = ref({
  opinion: '',
})

// 适用性评估项定义（与后端字段匹配）
const assessItems = [
  { key: 'functionality', label: '功能特性' },
  { key: 'security', label: '安全性' },
  { key: 'reliability', label: '可靠性' },
  { key: 'performance', label: '性能效率' },
  { key: 'compatibility', label: '兼容性' },
  { key: 'maintainability', label: '可维护性' },
  { key: 'extensibility', label: '可扩展性' },
  { key: 'usability', label: '易用性' },
]

// 获取评估状态
function getAssessStatus(key: string) {
  const assessment = taskData.value?.applicability_assessment
  if (!assessment) return 'not_assessed'
  // 字段名：functionalityStatus, securityStatus, etc.
  const statusKey = `${key}Status`
  return assessment[statusKey] || 'not_assessed'
}

// 获取偏差说明
function getAssessDeviation(key: string) {
  const assessment = taskData.value?.applicability_assessment
  if (!assessment) return ''
  // 字段名：functionalityDeviation, securityDeviation, etc.
  const deviationKey = `${key}Deviation`
  return assessment[deviationKey] || ''
}

// 加载任务详情
async function loadTaskDetail() {
  if (!props.taskId) return

  console.log('=== [ApprovalDialog] 开始加载任务详情, taskId:', props.taskId, ' ===')

  try {
    // 获取任务完整详情（包含扩展表数据）
    const taskRes = await taskApi.getTaskDetail(props.taskId)
    const data = taskRes.data || taskRes
    
    // 合并任务主表和扩展表数据
    taskData.value = {
      ...data.task,
      ...data.task_detail,  // 包含 component_name, version 等
      applicability_assessment: data.applicability_assessment,
    }

    // 审批节点
    approvalNodes.value = data.approval_nodes || []

    // 审批历史
    approvalHistory.value = data.approval_history || []

    console.log('[ApprovalDialog] 任务详情:', JSON.stringify(taskData.value, null, 2))
    console.log('[ApprovalDialog] componentName:', taskData.value?.componentName)
    console.log('[ApprovalDialog] componentVersion:', taskData.value?.componentVersion)
    console.log('[ApprovalDialog] remarkFiles:', taskData.value?.remarkFiles)
    console.log('[ApprovalDialog] wikiUrl:', taskData.value?.wikiUrl)
    console.log('[ApprovalDialog] remark:', taskData.value?.remark)
    console.log('[ApprovalDialog] 审批节点:', approvalNodes.value)
    console.log('[ApprovalDialog] 审批历史:', approvalHistory.value)
  } catch (e: any) {
    console.error('[ApprovalDialog] 加载任务详情失败:', e)
    ElMessage.error('加载任务详情失败: ' + (e.message || '未知错误'))
  }
}

// 监听对话框打开
watch(() => props.modelValue, (newVal) => {
  if (newVal && props.taskId) {
    loadTaskDetail()
    approvalForm.value.opinion = ''
  }
})

function handleClose() {
  visible.value = false
}

async function handleApprove() {
  if (!props.taskId) return

  // 通过时意见选填，但建议填写
  if (!approvalForm.value.opinion?.trim()) {
    ElMessage.warning('建议填写审批意见')
    // 不阻止提交
  }

  try {
    await ElMessageBox.confirm('确认通过该任务？', '确认通过', {
      type: 'success',
    })

    // 找到当前审批节点的ID
    const currentNode = approvalNodes.value.find(n => n.status === 'pending' && n.approverId === props.currentUser.id)
    if (!currentNode || !currentNode.id) {
      ElMessage.error('未找到当前审批节点')
      return
    }

    await taskApi.approveTask(props.taskId, {
      nodeId: currentNode.id,
      opinion: approvalForm.value.opinion || '同意',
    })

    ElMessage.success('审批通过')
    emit('approved')
    handleClose()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error('[ApprovalDialog] 审批失败:', e)
      ElMessage.error('审批失败: ' + (e.message || '未知错误'))
    }
  }
}

async function handleReject() {
  if (!props.taskId) return

  // 驳回时意见必填
  if (!approvalForm.value.opinion?.trim()) {
    ElMessage.warning('请填写驳回原因')
    return
  }

  try {
    await ElMessageBox.confirm('确认驳回该任务？驳回后将无法恢复。', '确认驳回', {
      type: 'warning',
    })

    // 找到当前审批节点的ID
    const currentNode = approvalNodes.value.find(n => n.status === 'pending' && n.approverId === props.currentUser.id)
    if (!currentNode || !currentNode.id) {
      ElMessage.error('未找到当前审批节点')
      return
    }

    await taskApi.rejectTask(props.taskId, {
      nodeId: currentNode.id,
      opinion: approvalForm.value.opinion,
    })

    ElMessage.success('已驳回')
    emit('rejected')
    handleClose()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error('[ApprovalDialog] 驳回失败:', e)
      ElMessage.error('驳回失败: ' + (e.message || '未知错误'))
    }
  }
}

// 工具函数
function getTaskTypeColor(type: string): string {
  const colors: Record<string, string> = {
    introduce: 'primary',
    upgrade: 'warning',
    withdraw: 'info',
    usage: 'success',
    risk_report: 'danger',
    defect: 'danger',
  }
  return colors[type] || ''
}

function getTaskTypeLabel(type: string): string {
  return TASK_TYPE_LABELS[type as TaskType] || type
}

function getStatusType(status: string): string {
  const types: Record<string, string> = {
    draft: 'info',
    approving: 'warning',
    withdrawn: 'info',
    approved: 'success',
    rejected: 'danger',
  }
  return types[status] || 'info'
}

function getStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    draft: '草稿',
    approving: '审批中',
    withdrawn: '已撤回',
    approved: '已通过',
    rejected: '已驳回',
  }
  return labels[status] || status
}

function getNodeStatusType(status: string): string {
  const types: Record<string, string> = {
    pending: 'info',
    approved: 'success',
    rejected: 'danger',
    skipped: 'info',
  }
  return types[status] || 'info'
}

function getNodeStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已驳回',
    skipped: '已跳过',
  }
  return labels[status] || status
}

function getActionColor(action: string): string {
  const colors: Record<string, string> = {
    create: 'info',
    submit: 'primary',
    approve: 'success',
    reject: 'danger',
    withdraw: 'warning',
  }
  return colors[action] || 'info'
}

function getActionLabel(action: string): string {
  const labels: Record<string, string> = {
    create: '创建',
    submit: '提交',
    approve: '通过',
    reject: '驳回',
    withdraw: '撤回',
  }
  return labels[action] || action
}

function formatDate(date: string | null): string {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 适用性评估标签
function getAssessLabel(key: string): string {
  const labels: Record<string, string> = {
    functionality: '功能特性',
    security: '安全性',
    reliability: '可靠性',
    performance: '性能效率',
    compatibility: '兼容性',
    maintainability: '可维护性',
    community: '社区活跃度',
    license_compliance: '许可证合规性',
  }
  return labels[key] || key
}

function getAssessStatusText(status: string): string {
  const texts: Record<string, string> = {
    fully_meet: '完全满足',
    deviation: '存在偏差',
    not_assessed: '未评估',
  }
  return texts[status] || '未评估'
}

// 格式化审批人信息：角色-部门-人名
function formatApproverInfo(item: any): string {
  const nodeName = item.nodeName || ''
  const approver = item.approver || ''
  const approverDept = item.approverDept || ''
  
  // 节点名称通常包含角色信息，如"开源治理专员复审"
  // 格式：节点名称-部门-人名
  if (approver && approverDept) {
    return `${nodeName}-${approverDept}-${approver}`
  } else if (approver) {
    return `${nodeName}-${approver}`
  }
  return nodeName
}

// 审批流节点显示
function getFlowNodeClass(node: any): string {
  const map: Record<string, string> = {
    pending: 'flow-node-pending',
    approved: 'flow-node-approved',
    rejected: 'flow-node-rejected',
    skipped: 'flow-node-skipped',
  }
  return map[node.status] || ''
}

function getFlowNodeTagType(node: any): string {
  const map: Record<string, string> = {
    pending: 'info',
    approved: 'success',
    rejected: 'danger',
    skipped: 'info',
  }
  return map[node.status] || 'info'
}

function getFlowNodeStatusText(node: any): string {
  const map: Record<string, string> = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已驳回',
    skipped: '已跳过',
  }
  return map[node.status] || node.status
}

// 审批历史保留从最后一次提交开始的所有记录，不再额外去重
</script>

<style scoped>
.approval-detail {
  max-height: 70vh;
  overflow-y: auto;
}

.mb-4 {
  margin-bottom: 16px;
}

/* 统一适用性评估和备注说明的标签宽度 */
.assess-descriptions :deep(.el-descriptions__label),
.remark-descriptions :deep(.el-descriptions__label) {
  width: 120px;
  min-width: 120px;
  max-width: 120px;
}

/* 偏差说明文字 */
.deviation-text {
  margin-top: 8px;
  color: #909399;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
}

.deviation-text:hover {
  white-space: normal;
}

/* 审批历史时间轴 */
.el-timeline {
  padding-left: 0;
  margin-left: 4px;
}

.el-timeline-item__timestamp {
  font-size: 12px;
  color: #909399;
}

/* ========== 审批流程可视化 ========== */
.approval-flow {
  display: flex;
  align-items: flex-start;
  padding: 16px 0;
  overflow-x: auto;
  gap: 0;
}

.flow-node-wrapper {
  display: flex;
  align-items: flex-start;
  flex-shrink: 0;
}

.flow-node {
  width: 200px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  padding: 12px 14px;
  background: #fff;
  transition: all 0.3s;
}

.flow-node.flow-node-pending {
  border-color: #dcdfe6;
  background: #fafafa;
}

.flow-node.flow-node-pending .flow-node-name,
.flow-node.flow-node-pending .flow-node-info {
  color: #c0c4cc;
}

.flow-node.flow-node-approved {
  border-color: #67c23a;
  background: #f0f9eb;
}

.flow-node.flow-node-approved .flow-node-name,
.flow-node.flow-node-approved .flow-node-info {
  color: #303133;
}

.flow-node.flow-node-rejected {
  border-color: #f56c6c;
  background: #fef0f0;
}

.flow-node.flow-node-rejected .flow-node-name,
.flow-node.flow-node-rejected .flow-node-info {
  color: #303133;
}

.flow-node.flow-node-skipped {
  border-color: #ebeef5;
  background: #f5f7fa;
  opacity: 0.6;
}

.flow-node-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.flow-node-name {
  font-size: 13px;
  font-weight: 600;
  flex: 1;
}

.flow-node-tag {
  flex-shrink: 0;
}

.flow-node-info {
  font-size: 12px;
  color: #606266;
}

.flow-node-info.text-gray {
  color: #c0c4cc;
}

.flow-arrow {
  display: flex;
  align-items: center;
  padding: 0 8px;
  color: #c0c4cc;
  margin-top: 16px;
  flex-shrink: 0;
}

.flow-arrow .el-icon {
  color: #c0c4cc;
}
</style>
