<template>
  <div class="task-detail">
    <!-- 页面头部 -->
    <el-page-header @back="goBack">
      <template #content>
        <span class="page-title">{{ task?.title || '任务详情' }}</span>
      </template>
    </el-page-header>

    <div v-if="task" style="margin-top: 16px">
      <el-row :gutter="16">
        <!-- 左侧：任务信息 -->
        <el-col :span="isReadOnly ? 24 : 16">
          <!-- 草稿状态显示可编辑表单（治理中心只读详情不显示） -->
          <template v-if="task.status === 'draft' && !isReadOnly">
            <el-card shadow="never" class="info-card">
              <template #header>
                <div class="card-header">
                  <span>编辑任务</span>
                </div>
              </template>
              <TaskForm
                ref="taskFormRef"
                :task-data="task"
                :read-only="false"
                @submit="handleSubmit"
                @save-draft="handleSaveDraft"
              />
            </el-card>
          </template>

          <!-- 非草稿状态显示只读详情 -->
          <template v-else>
            <!-- 任务基本信息卡片 -->
            <el-card shadow="never" class="info-card">
            <template #header>
              <div class="card-header">
                <span>任务信息</span>
                <div class="header-actions">
                  <TaskStatusIndicator :status="task.status" />
                </div>
              </div>
            </template>

            <el-descriptions :column="2" border>
              <el-descriptions-item label="任务编号">{{ task.taskNo }}</el-descriptions-item>
              <el-descriptions-item label="任务类型">
                <el-tag
                  v-if="task.type === 'defect'"
                  type="danger"
                  size="small"
                >
                  {{ TASK_TYPE_LABELS[task.type] }}
                </el-tag>
                <span v-else>{{ TASK_TYPE_LABELS[task.type] }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="申请人">{{ task.applicant }}</el-descriptions-item>
              <el-descriptions-item label="所属部门">{{ task.applicantDept }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ task.createdAt }}</el-descriptions-item>
              <el-descriptions-item label="更新时间">{{ task.updatedAt }}</el-descriptions-item>
              <el-descriptions-item label="当前步骤" :span="2">{{ task.currentStep }}</el-descriptions-item>
            </el-descriptions>
          </el-card>

          <!-- 任务专属信息 -->
          <el-card shadow="never" class="info-card" style="margin-top: 16px">
            <template #header>
              <span>{{ TASK_TYPE_LABELS[task.type] }}信息</span>
            </template>

            <!-- 引入申请信息 -->
          <el-descriptions v-if="task.type === 'introduce'" :column="2" border>
            <el-descriptions-item label="选择组件">{{ task.selectedComponent || '-' }}</el-descriptions-item>
            <el-descriptions-item label="组件名称">{{ task.componentName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="版本号">{{ task.componentVersion || '-' }}</el-descriptions-item>
            <el-descriptions-item label="许可证类型">{{ task.licenseType || '-' }}</el-descriptions-item>
            <el-descriptions-item label="引入原因" :span="2">{{ task.introduceReason || '-' }}</el-descriptions-item>
          </el-descriptions>

            <!-- 升级申请信息 -->
            <el-descriptions v-if="task.type === 'upgrade'" :column="2" border>
              <el-descriptions-item label="组件名称">{{ task.componentName || '-' }}</el-descriptions-item>
              <el-descriptions-item label="当前版本">{{ task.currentVersion || '-' }}</el-descriptions-item>
              <el-descriptions-item label="目标版本">{{ task.targetVersion || '-' }}</el-descriptions-item>
              <el-descriptions-item label="测试报告链接" :span="2">
                <el-link v-if="task.testReportLink" :href="task.testReportLink" target="_blank" type="primary">
                  {{ task.testReportLink }}
                  <el-icon><Link /></el-icon>
                </el-link>
                <span v-else>-</span>
              </el-descriptions-item>
              <el-descriptions-item label="升级原因" :span="2">{{ task.description || '-' }}</el-descriptions-item>
            </el-descriptions>

            <!-- 退出申请信息 -->
            <el-descriptions v-if="task.type === 'withdraw'" :column="2" border>
              <el-descriptions-item label="拟退出组件">{{ task.componentName || '-' }}</el-descriptions-item>
              <el-descriptions-item label="当前版本">{{ task.componentVersion || '-' }}</el-descriptions-item>
              <el-descriptions-item label="退出原因" :span="2">{{ task.withdrawReason || '-' }}</el-descriptions-item>
              <el-descriptions-item label="替代组件">{{ task.replacementComponent || '-' }}</el-descriptions-item>
              <el-descriptions-item label="替代方案说明" :span="2">{{ task.alternativePlan || '-' }}</el-descriptions-item>
              <el-descriptions-item label="代码清理计划" :span="2">{{ task.codeCleanupPlan || '-' }}</el-descriptions-item>
              <el-descriptions-item label="受影响项目" :span="2">{{ task.affectedProjects || '-' }}</el-descriptions-item>
            </el-descriptions>

            <!-- 缺陷处置信息 -->
            <template v-if="task.type === 'defect'">
              <el-descriptions :column="2" border>
                <el-descriptions-item label="风险来源">{{ task.riskSource || '-' }}</el-descriptions-item>
                <el-descriptions-item label="CVE编号">{{ task.cveId || '-' }}</el-descriptions-item>
                <el-descriptions-item label="严重等级">
                  <RiskTag v-if="task.severityLevel" :level="task.severityLevel" />
                  <span v-else>-</span>
                </el-descriptions-item>
                <el-descriptions-item label="受影响资产" :span="2">{{ task.affectedAssets || '-' }}</el-descriptions-item>
                <el-descriptions-item label="修复方案" :span="2">{{ task.fixPlan || '-' }}</el-descriptions-item>
              </el-descriptions>

              <!-- 处置状态选择 -->
              <div class="disposal-section" style="margin-top: 20px">
                <el-divider content-position="left">处置信息</el-divider>
                <el-form :model="disposalForm" label-width="120px">
                  <el-form-item label="处置状态">
                    <el-radio-group v-model="disposalForm.status" :disabled="!canEditDisposal">
                      <el-radio value="no-need">无需处置</el-radio>
                      <el-radio value="completed">已处置</el-radio>
                    </el-radio-group>
                  </el-form-item>

                  <el-form-item v-if="disposalForm.status === 'no-need'" label="无需处置原因">
                    <el-input
                      v-model="disposalForm.reason"
                      type="textarea"
                      :rows="3"
                      placeholder="请说明无需处置的原因"
                      :disabled="!canEditDisposal"
                    />
                  </el-form-item>

                  <el-form-item v-if="disposalForm.status === 'completed'" label="生产变更工单">
                    <el-input
                      v-model="disposalForm.changeOrder"
                      placeholder="请输入生产变更工单号"
                      :disabled="!canEditDisposal"
                    />
                  </el-form-item>

                  <el-form-item v-if="canEditDisposal">
                    <el-button type="primary" @click="handleSaveDisposal">保存</el-button>
                    <el-button type="success" @click="handleSubmitDisposal">提交</el-button>
                  </el-form-item>
                </el-form>
              </div>
            </template>
          </el-card>

          <!-- 审批流程（与通用审批弹窗保持一致） -->
          <el-card shadow="never" class="info-card" style="margin-top: 16px">
            <template #header>
              <span>审批流程</span>
            </template>

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
          </el-card>

          <!-- 审批历史（与通用审批弹窗保持一致） -->
          <el-card shadow="never" class="info-card" style="margin-top: 16px">
            <template #header>
              <span>审批历史</span>
            </template>

            <el-timeline v-if="deduplicatedHistory.length > 0">
              <el-timeline-item
                v-for="(item, index) in deduplicatedHistory"
                :key="index"
                :timestamp="formatDate(item.createTime || item.approvedAt)"
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
          </el-card>
          </template>
        </el-col>

        <!-- 右侧：操作面板 -->
        <el-col v-if="!isReadOnly" :span="8">
          <el-card shadow="never" class="action-card">
            <template #header>
              <span>操作</span>
            </template>

            <div class="action-buttons">
              <!-- 草稿状态：提交按钮 -->
              <template v-if="task.status === 'draft'">
                <el-button 
                  type="primary" 
                  class="action-btn" 
                  @click="handleSubmitOrShowReport"
                >
                  <el-icon><Promotion /></el-icon>
                  {{ task.type === 'introduce' ? '提交评估' : '提交' }}
                </el-button>
              </template>

              <!-- 待审批/处理中/已分配：可以审批 -->
              <template v-if="['pending', 'processing', 'assigned'].includes(task.status)">
                <el-button
                  type="success"
                  class="action-btn"
                  @click="showApprovalDialog('approve')"
                >
                  <el-icon><Check /></el-icon>
                  通过
                </el-button>
                <el-button
                  type="danger"
                  class="action-btn"
                  @click="showApprovalDialog('reject')"
                >
                  <el-icon><Close /></el-icon>
                  驳回
                </el-button>
              </template>

              <!-- 已完成/已驳回：只读 -->
              <template v-if="['completed', 'rejected'].includes(task.status)">
                <el-alert
                  :title="task.status === 'completed' ? '任务已完成' : '任务已驳回'"
                  :type="task.status === 'completed' ? 'success' : 'error'"
                  :closable="false"
                  show-icon
                />
              </template>
            </div>
          </el-card>

          <!-- 审批流程（与通用审批弹窗保持一致） -->
          <el-card v-if="task.status !== 'draft'" shadow="never" class="info-card" style="margin-top: 16px">
            <template #header>
              <span>审批流程</span>
            </template>

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
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 审批对话框 -->
    <el-dialog
      v-model="approvalDialogVisible"
      :title="approvalAction === 'approve' ? '审批通过' : '审批驳回'"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="approvalForm" label-width="100px">
        <el-form-item label="审批结果" required>
          <el-radio-group v-model="approvalForm.result">
            <el-radio value="approve">通过</el-radio>
            <el-radio value="reject">驳回</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="审批意见" required>
          <el-input
            v-model="approvalForm.opinion"
            type="textarea"
            :rows="4"
            :placeholder="approvalAction === 'approve' ? '请输入审批通过意见' : '请输入驳回原因'"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="approvalDialogVisible = false">取消</el-button>
        <el-button
          :type="approvalForm.result === 'approve' ? 'success' : 'danger'"
          @click="submitApproval"
        >
          确认
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Link, Check, Close, Promotion, Document, Clock, ArrowRight } from '@element-plus/icons-vue'
import type { Task, TaskStatus } from '@/types'
import { TASK_TYPE_LABELS } from '@/types'
import { useTaskStore } from '@/stores/task'
import TaskStatusIndicator from '@/components/operator/TaskStatusIndicator.vue'
import TaskForm from '@/components/operator/TaskForm.vue'
import RiskTag from '@/components/business/RiskTag.vue'

const router = useRouter()
const route = useRoute()
const taskStore = useTaskStore()
const task = ref<Task | null>(null)
const taskFormRef = ref<InstanceType<typeof TaskForm> | null>(null)

// 治理中心任务列表进入的详情页为只读，不展示操作按钮
const isReadOnly = computed(() => route.path.startsWith('/admin/governance/task-list'))

// 处置表单
const disposalForm = reactive({
  status: 'no-need' as 'no-need' | 'completed',
  reason: '',
  changeOrder: '',
})

// 是否可以编辑处置信息（根据角色和任务状态）
const canEditDisposal = computed(() => {
  if (!task.value) return false
  // 治理中心只读详情不允许编辑处置信息
  if (isReadOnly.value) return false
  // 这里可以根据用户角色和任务状态来判断
  // 简化处理：非草稿状态且任务已分配给当前用户可以编辑
  return task.value.status !== 'draft' && task.value.status !== 'completed' && task.value.status !== 'rejected'
})

// 审批流程表格数据
const approvalHistory = computed(() => {
  if (!task.value) return []
  // 使用后端返回的 approval_history 数据（从最后一次提交开始）
  if (task.value.approvalHistory && task.value.approvalHistory.length > 0) {
    return task.value.approvalHistory.map((item: any, index: number) => ({
      order: index + 1,
      role: item.nodeName || '-',
      approver: item.approver || '-',
      result: item.action === 'approve' ? '通过' : item.action === 'reject' ? '驳回' : item.action === 'submit' ? '提交' : '待审批',
      opinion: item.opinion || '-',
      approvedAt: item.createTime || '-',
    }))
  }
  // 没有 approval_history 时，使用 approval_nodes
  if (task.value.approvalNodes) {
    return task.value.approvalNodes.map((node, index) => ({
      order: index + 1,
      role: node.name,
      approver: node.approver || '-',
      result: node.status === 'approved' ? '通过' : node.status === 'rejected' ? '驳回' : '待审批',
      opinion: node.opinion || '',
      approvedAt: node.approvedAt || '-',
    }))
  }
  return []
})

// 审批流程节点数据（过滤掉节点0发起节点，只显示实际审批节点）
const activeFlowNodes = computed(() => {
  if (!task.value || !task.value.approvalNodes) return []
  return task.value.approvalNodes
    .filter((node: any) => node.nodeOrder > 0 && !node.isSkipped)
    .map((node: any) => ({
      id: node.id,
      nodeName: node.nodeName || node.name,
      approver: node.approver || '',
      approverDept: node.approverDept || '',
      status: node.status,
    }))
})

// 审批历史（直接使用后端返回数据，后端已过滤为最后一次提交开始）
const deduplicatedHistory = computed(() => {
  if (!task.value) return []
  return task.value.approvalHistory || []
})

function getFlowNodeClass(node: any) {
  return {
    'is-completed': node.status === 'approved' || node.status === 'completed',
    'is-active': node.status === 'processing' || node.status === 'pending',
    'is-rejected': node.status === 'rejected',
  }
}

function getFlowNodeTagType(node: any): any {
  if (node.status === 'approved' || node.status === 'completed') return 'success'
  if (node.status === 'rejected') return 'danger'
  if (node.status === 'processing') return 'primary'
  return 'info'
}

function getFlowNodeStatusText(node: any) {
  if (node.status === 'approved' || node.status === 'completed') return '已通过'
  if (node.status === 'rejected') return '已驳回'
  if (node.status === 'processing') return '审批中'
  return '待审批'
}

function getActionColor(action: string): any {
  if (action === 'approve' || action === 'submit') return 'success'
  if (action === 'reject') return 'danger'
  if (action === 'withdraw') return 'warning'
  return 'info'
}

function getActionLabel(action: string) {
  const labels: Record<string, string> = {
    submit: '提交',
    approve: '通过',
    reject: '驳回',
    withdraw: '撤回',
  }
  return labels[action] || action
}

function formatDate(date: string | undefined) {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// 审批对话框
const approvalDialogVisible = ref(false)
const approvalAction = ref<'approve' | 'reject'>('approve')
const approvalForm = reactive({
  result: 'approve' as 'approve' | 'reject',
  opinion: '',
})

async function loadTask() {
  const id = route.params.id as string
  console.log('[TaskDetail] 加载任务详情, id:', id)
  
  // 从后端API获取任务详情
  try {
    const detail: any = await taskStore.getTaskDetail(id)
    console.log('[TaskDetail] 任务详情:', detail)
    console.log('[TaskDetail] detail.task:', detail.task)
    console.log('[TaskDetail] detail.task_detail:', detail.task_detail)
    
    // 合并task和task_detail数据
    const taskData = detail.task || {}
    const taskDetail = detail.task_detail || {}
    
    task.value = {
      ...taskData,
      ...taskDetail,
      // 审批节点
      approvalNodes: detail.approval_nodes || [],
      // 审批历史
      approvalHistory: detail.approval_history || [],
    } as any
    
    console.log('[TaskDetail] 合并后的task.value:', task.value)
  } catch (e: any) {
    console.error('[TaskDetail] 获取任务详情失败:', e)
    ElMessage.error('获取任务详情失败: ' + (e.message || '未知错误'))
    goBack()
  }
}

onMounted(loadTask)

function goBack() {
  router.back()
}

// 审批节点状态类型
function getNodeStatusType(status: string): 'success' | 'warning' | 'info' | 'danger' {
  const types: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    pending: 'info',
    approved: 'success',
    rejected: 'danger',
    not_required: 'info',
  }
  return types[status] || 'info'
}

function getNodeStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    pending: '待审批',
    approved: '已同意',
    rejected: '已驳回',
    not_required: '无需审批',
  }
  return labels[status] || '未知'
}

function showApprovalDialog(action: 'approve' | 'reject') {
  approvalAction.value = action
  approvalForm.result = action
  approvalForm.opinion = ''
  approvalDialogVisible.value = true
}

function submitApproval() {
  if (!approvalForm.opinion?.trim()) {
    ElMessage.warning('请填写审批意见')
    return
  }

  const actionLabel = approvalForm.result === 'approve' ? '通过' : '驳回'
  ElMessageBox.confirm(`确认${actionLabel}该任务？`, '确认', {
    type: approvalForm.result === 'approve' ? 'success' : 'warning',
  }).then(() => {
    if (task.value?.approvalNodes) {
      // 找到第一个待审批的节点
      const pendingNode = task.value.approvalNodes.find(n => n.status === 'pending')
      if (pendingNode) {
        pendingNode.status = approvalForm.result === 'approve' ? 'approved' : 'rejected'
        pendingNode.opinion = approvalForm.opinion
        pendingNode.approvedAt = new Date().toLocaleString('zh-CN')
        pendingNode.approver = '当前用户'

        // 如果驳回，更新任务状态
        if (approvalForm.result === 'reject') {
          task.value.status = 'rejected'
          task.value.currentStep = '已驳回'
        } else {
          // 检查是否所有节点都已审批
          const allApproved = task.value.approvalNodes.every(
            n => n.status === 'approved' || n.status === 'not_required'
          )
          if (allApproved) {
            task.value.status = 'completed'
            task.value.currentStep = '已完成'
          }
        }
      }
      ElMessage.success(`审批${actionLabel}`)
    }
    approvalDialogVisible.value = false
  }).catch(() => {})
}

function submitTask() {
  // 引入申请通过评估报告对话框提交，这里处理其他类型
  if (task.value?.type === 'introduce') {
    ElMessage.info('引入申请请通过"提交"按钮查看评估报告')
    return
  }
  
  ElMessageBox.confirm('确认提交任务？提交后将进入审批流程。', '确认提交', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    if (task.value) {
      task.value.status = 'pending'
      task.value.updatedAt = new Date().toLocaleString('zh-CN')
      ElMessage.success('任务已提交，进入审批流程')
    }
  }).catch(() => {})
}

// 提交或显示评估报告
function handleSubmitOrShowReport() {
  if (task.value?.type === 'introduce') {
    // 引入申请显示评估报告
    if (taskFormRef.value) {
      taskFormRef.value.showEvalReport?.()
    }
  } else {
    // 其他类型直接提交
    submitTask()
  }
}

// 处理表单提交
function handleSubmit(data: any) {
  if (task.value) {
    // 更新任务数据
    Object.assign(task.value, data)
    task.value.updatedAt = new Date().toLocaleString('zh-CN')
    
    // 调用提交逻辑
    submitTask()
  }
}

// 处理暂存
function handleSaveDraft(data: any) {
  if (task.value) {
    // 更新任务数据
    Object.assign(task.value, data)
    task.value.updatedAt = new Date().toLocaleString('zh-CN')
    ElMessage.success('草稿已保存')
  }
}

// 保存处置信息
function handleSaveDisposal() {
  if (!task.value) return

  // 验证必填字段
  if (disposalForm.status === 'no-need' && !disposalForm.reason?.trim()) {
    ElMessage.warning('请填写无需处置的原因')
    return
  }
  if (disposalForm.status === 'completed' && !disposalForm.changeOrder?.trim()) {
    ElMessage.warning('请填写生产变更工单号')
    return
  }

  // 保存到任务数据
  task.value.disposalStatus = disposalForm.status
  task.value.disposalReason = disposalForm.reason
  task.value.disposalChangeOrder = disposalForm.changeOrder
  task.value.updatedAt = new Date().toLocaleString('zh-CN')
  
  ElMessage.success('处置信息已保存')
}

// 提交处置信息
function handleSubmitDisposal() {
  if (!task.value) return

  // 验证必填字段
  if (disposalForm.status === 'no-need' && !disposalForm.reason?.trim()) {
    ElMessage.warning('请填写无需处置的原因')
    return
  }
  if (disposalForm.status === 'completed' && !disposalForm.changeOrder?.trim()) {
    ElMessage.warning('请填写生产变更工单号')
    return
  }

  ElMessageBox.confirm('确认提交处置信息？提交后将进入审批流程。', '确认提交', {
    type: 'info',
  }).then(() => {
    // 保存到任务数据
    task.value!.disposalStatus = disposalForm.status
    task.value!.disposalReason = disposalForm.reason
    task.value!.disposalChangeOrder = disposalForm.changeOrder
    task.value!.status = 'processing' as TaskStatus
    task.value!.currentStep = '处置中'
    task.value!.updatedAt = new Date().toLocaleString('zh-CN')
    
    ElMessage.success('处置信息已提交，进入审批流程')
  }).catch(() => {})
}
</script>

<style scoped>
.task-detail {
  padding: 0;
}

.page-title {
  font-size: 16px;
  font-weight: 600;
}

.info-card {
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
}

.action-card {
  border-radius: 8px;
  position: sticky;
  top: 16px;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-btn {
  width: 100%;
}

/* 审批流程样式 */
.approval-timeline {
  padding: 8px 0;
}

.approval-node {
  border-left: 2px solid #E4E7ED;
  padding-left: 20px;
  padding-bottom: 24px;
  position: relative;
}

.approval-node:last-child {
  border-left-color: transparent;
  padding-bottom: 0;
}

.approval-node::before {
  content: '';
  position: absolute;
  left: -6px;
  top: 0;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #E4E7ED;
}

.approval-node:first-child::before {
  background: #409EFF;
}

.node-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.node-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #303133;
}

.node-index {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #409EFF;
  color: #fff;
  font-size: 12px;
}

.node-content {
  background: #F5F7FA;
  padding: 12px;
  border-radius: 6px;
}

.node-info {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  color: #909399;
  margin-bottom: 8px;
}

.node-approver {
  font-weight: 500;
}

.node-opinion {
  font-size: 14px;
  color: #606266;
  line-height: 1.5;
}

.opinion-label {
  font-weight: 500;
  color: #303133;
}

.text-muted {
  color: #C0C4CC;
}

/* 审批进度鱼骨条样式 */
.approval-progress-vertical {
  padding: 20px 0;
}

.progress-node {
  position: relative;
  display: flex;
  align-items: flex-start;
  padding-left: 40px;
  padding-bottom: 32px;
}

.progress-node:last-child {
  padding-bottom: 0;
}

.node-icon {
  position: absolute;
  left: 0;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #f5f7fa;
  border: 2px solid #dcdfe6;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  z-index: 1;
}

.progress-node.is-completed .node-icon {
  background: #f0f9ff;
  border-color: #67C23A;
}

.progress-node.is-active .node-icon {
  background: #ecf5ff;
  border-color: #409EFF;
  animation: pulse 2s infinite;
}

.progress-node.is-pending .node-icon {
  background: #f5f7fa;
  border-color: #dcdfe6;
}

.node-line {
  position: absolute;
  left: 15px;
  top: 32px;
  bottom: 0;
  width: 2px;
  background: #dcdfe6;
}

.progress-node.is-completed .node-line {
  background: #67C23A;
}

.node-content {
  flex: 1;
}

.node-role {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
}

.node-approver {
  font-size: 13px;
  color: #606266;
  margin-bottom: 4px;
}

.node-opinion {
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
  padding: 4px 8px;
  background: #f5f7fa;
  border-radius: 4px;
}

.node-time {
  font-size: 12px;
  color: #c0c4cc;
}

/* ========== 审批流程可视化（与通用审批弹窗保持一致） ========== */
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

.flow-node.is-completed {
  border-color: #67c23a;
  background: #f0f9eb;
}

.flow-node.is-completed .flow-node-name,
.flow-node.is-completed .flow-node-info {
  color: #303133;
}

.flow-node.is-active {
  border-color: #dcdfe6;
  background: #fafafa;
}

.flow-node.is-active .flow-node-name,
.flow-node.is-active .flow-node-info {
  color: #c0c4cc;
}

.flow-node.is-rejected {
  border-color: #f56c6c;
  background: #fef0f0;
}

.flow-node.is-rejected .flow-node-name,
.flow-node.is-rejected .flow-node-info {
  color: #303133;
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

@keyframes pulse {
  0%, 100% {
    box-shadow: 0 0 0 0 rgba(64, 158, 255, 0.4);
  }
  50% {
    box-shadow: 0 0 0 8px rgba(64, 158, 255, 0);
  }
}
</style>
