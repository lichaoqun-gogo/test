<template>
  <div class="introduce-detail-page">
    <el-page-header @back="goBack" class="page-header">
      <template #content>
        <span class="page-title">引入审批详情</span>
        <el-tag v-if="detail" :type="(getStatusType(detail.status) as any)" size="small" style="margin-left: 12px">
          {{ getStatusLabel(detail.status) }}
        </el-tag>
      </template>
    </el-page-header>

    <template v-if="detail">
      <!-- 审批流程时间线 -->
      <el-card class="section-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><Clock /></el-icon>
            <span>审批流程</span>
          </div>
        </template>
        <el-timeline>
          <el-timeline-item
            v-for="(item, index) in detail.approvalHistory"
            :key="index"
            :type="getTimelineType(item.action)"
            :size="index === (detail.approvalHistory?.length || 0) - 1 ? 'large' : 'normal'"
          >
            <div class="timeline-content">
              <div class="timeline-header">
                <span class="timeline-action">{{ item.action }}</span>
                <span class="timeline-operator">{{ item.operator }}</span>
                <span class="timeline-time">{{ formatTime(item.time) }}</span>
              </div>
              <p class="timeline-opinion">{{ item.opinion }}</p>
            </div>
          </el-timeline-item>
        </el-timeline>
      </el-card>

      <!-- 申请基本信息 -->
      <el-card class="section-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><Document /></el-icon>
            <span>申请基本信息</span>
          </div>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="申请编号">{{ detail.applicationNo }}</el-descriptions-item>
          <el-descriptions-item label="当前步骤">{{ detail.currentStep }}</el-descriptions-item>
          <el-descriptions-item label="组件名称">{{ detail.componentName }}</el-descriptions-item>
          <el-descriptions-item label="组件版本">v{{ detail.componentVersion }}</el-descriptions-item>
          <el-descriptions-item label="组件类型">{{ getComponentTypeLabel(detail.componentType) }}</el-descriptions-item>
          <el-descriptions-item label="许可证类型">{{ detail.licenseType }}</el-descriptions-item>
          <el-descriptions-item label="官方网站" :span="2">
            <el-link :href="detail.officialUrl" target="_blank" type="primary">{{ detail.officialUrl }}</el-link>
          </el-descriptions-item>
          <el-descriptions-item label="申请人">{{ detail.applicant }}</el-descriptions-item>
          <el-descriptions-item label="所属部门">{{ detail.department }}</el-descriptions-item>
          <el-descriptions-item label="联系邮箱">{{ detail.email }}</el-descriptions-item>
          <el-descriptions-item label="联系电话">{{ detail.phone }}</el-descriptions-item>
          <el-descriptions-item label="使用场景" :span="2">{{ detail.useScenario }}</el-descriptions-item>
          <el-descriptions-item label="关联项目" :span="2">
            <el-tag v-for="proj in detail.projectNames" :key="proj" size="small" style="margin-right: 6px">
              {{ proj }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="组件描述" :span="2">{{ detail.description }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 预评估报告 -->
      <el-card class="section-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><DataAnalysis /></el-icon>
            <span>预评估报告</span>
            <el-tag :type="detail.preEvalConclusion === 'pass' ? 'success' : detail.preEvalConclusion === 'warning' ? 'warning' : 'danger'">
              {{ detail.preEvalConclusion === 'pass' ? '通过' : detail.preEvalConclusion === 'warning' ? '风险警告' : '不通过' }}
            </el-tag>
          </div>
        </template>
        <div class="eval-summary">
          <el-row :gutter="16">
            <el-col :span="8">
              <div class="eval-metric">
                <span class="metric-label">预评估评分</span>
                <span class="metric-value" :class="getScoreClass(detail.preEvalScore)">{{ detail.preEvalScore }}分</span>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="eval-metric">
                <span class="metric-label">核定等级</span>
                <GradeTag :grade="detail.grade" />
              </div>
            </el-col>
            <el-col :span="8">
              <div class="eval-metric">
                <span class="metric-label">提交时间</span>
                <span class="metric-value">{{ formatTime(detail.submittedAt) }}</span>
              </div>
            </el-col>
          </el-row>
        </div>
        <PreEvalReport
          :component-name="detail.componentName"
          :component-version="detail.componentVersion"
          :conclusion="detail.preEvalConclusion"
        />
      </el-card>

      <!-- 业务终评表单 -->
      <el-card v-if="canEvaluate" class="section-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><EditPen /></el-icon>
            <span>业务终评</span>
            <el-tag v-if="hasEvaluated" type="success" size="small">已评估</el-tag>
          </div>
        </template>
        <el-form :model="evalForm" label-width="140px">
          <el-form-item label="业务运行指标匹配度" required>
            <el-radio-group v-model="evalForm.match">
              <el-radio value="match">完全匹配</el-radio>
              <el-radio value="partial">部分匹配</el-radio>
              <el-radio value="unmatch">不匹配</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="部门适用性评估" required>
            <el-radio-group v-model="evalForm.applicability">
              <el-radio value="apply">适用</el-radio>
              <el-radio value="conditional">有条件适用</el-radio>
              <el-radio value="unapply">不适用</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="最终风险确认" required>
            <el-select v-model="evalForm.risk" placeholder="请选择风险等级" style="width: 200px">
              <el-option label="低风险" value="low" />
              <el-option label="中风险" value="medium" />
              <el-option label="高风险" value="high" />
            </el-select>
          </el-form-item>
          <el-form-item label="最终等级核定">
            <el-select v-model="evalForm.grade" placeholder="请选择" style="width: 200px">
              <el-option label="A - 核心关键（系统推荐）" value="A" />
              <el-option label="B - 重要支撑" value="B" />
              <el-option label="C - 一般使用" value="C" />
              <el-option label="D - 探索试用" value="D" />
            </el-select>
            <span class="grade-hint">系统根据预评估（{{ detail.preEvalScore }}分）+ 人工评估自动推荐</span>
          </el-form-item>
          <el-form-item label="补充意见">
            <el-input
              v-model="evalForm.opinion"
              type="textarea"
              :rows="3"
              placeholder="请输入审批意见"
            />
          </el-form-item>
        </el-form>

        <!-- 审批操作 -->
        <div class="action-bar">
          <el-button type="success" size="large" @click="handleApproveWithDownload">
            <el-icon><Check /></el-icon> 审批通过
          </el-button>
          <el-button type="danger" size="large" @click="handleReject">
            <el-icon><Close /></el-icon> 驳回
          </el-button>
          <el-button type="warning" size="large" @click="handleRequestSupplement">
            <el-icon><Document /></el-icon> 要求补充材料
          </el-button>
        </div>
      </el-card>

      <!-- 已完成的终评结果展示 -->
      <el-card v-if="hasEvaluated" class="section-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><CircleCheck /></el-icon>
            <span>终评结果</span>
          </div>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="业务匹配度">
            <el-tag :type="(getMatchType(detail.finalMatch) as any)">
              {{ getMatchLabel(detail.finalMatch) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="部门适用性">
            <el-tag :type="(getApplicabilityType(detail.finalApplicability) as any)">
              {{ getApplicabilityLabel(detail.finalApplicability) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="最终风险确认">
            <el-tag :type="(getRiskType(detail.finalRisk) as any)">
              {{ getRiskLabel(detail.finalRisk) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="最终等级">
            <GradeTag :grade="detail.finalGrade || detail.grade" />
          </el-descriptions-item>
          <el-descriptions-item label="审批意见" :span="2">{{ detail.finalOpinion || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 流程文档 -->
      <el-card class="section-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><FolderOpened /></el-icon>
            <span>流程文档</span>
          </div>
        </template>
        <el-button type="info" plain @click="generateWord">
          <el-icon><Download /></el-icon> 生成Word文档
        </el-button>
        <span class="doc-hint">点击后将生成包含完整审批流程的Word文档并下载</span>
      </el-card>
    </template>

    <!-- 驳回对话框 -->
    <el-dialog v-model="rejectDialogVisible" title="驳回申请" width="500px">
      <el-form :model="rejectForm" label-width="80px">
        <el-form-item label="驳回原因" required>
          <el-input
            v-model="rejectForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请输入驳回原因"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmReject">确认驳回</el-button>
      </template>
    </el-dialog>

    <!-- 审批通过对话框（含下载链接） -->
    <el-dialog v-model="approveDialogVisible" title="审批通过" width="600px">
      <el-form :model="approveForm" label-width="100px">
        <el-form-item label="审批意见" required>
          <el-input
            v-model="approveForm.opinion"
            type="textarea"
            :rows="3"
            placeholder="请输入审批意见"
          />
        </el-form-item>
        <el-form-item label="下载地址" required>
          <el-input
            v-model="approveForm.downloadUrl"
            placeholder="请输入软件/组件的下载地址"
          />
        </el-form-item>
        <el-form-item label="版本说明">
          <el-input
            v-model="approveForm.versionNote"
            type="textarea"
            :rows="2"
            placeholder="可选，填写版本相关说明"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="approveDialogVisible = false">取消</el-button>
        <el-button type="success" @click="confirmApprove">确认通过</el-button>
      </template>
    </el-dialog>

    <!-- 补充材料对话框 -->
    <el-dialog v-model="supplementDialogVisible" title="要求补充材料" width="500px">
      <el-form :model="supplementForm" label-width="100px">
        <el-form-item label="需补充内容" required>
          <el-input
            v-model="supplementForm.content"
            type="textarea"
            :rows="4"
            placeholder="请说明需要申请人补充的材料内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="supplementDialogVisible = false">取消</el-button>
        <el-button type="warning" @click="confirmSupplement">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Clock, Document, DataAnalysis, EditPen, Check, Close,
  CircleCheck, FolderOpened, Download,
} from '@element-plus/icons-vue'
import { mockIntroduceDetails } from '@/mock/admin'
import type { IntroduceApprovalDetail } from '@/mock/admin'
import { SOFTWARE_TYPE_LABELS, TICKET_STATUS_LABELS } from '@/types'
import { formatTime } from '@/utils/format'

const route = useRoute()
const router = useRouter()

const detail = ref<IntroduceApprovalDetail | null>(null)

const evalForm = ref({
  match: 'match' as 'match' | 'partial' | 'unmatch',
  applicability: 'apply' as 'apply' | 'conditional' | 'unapply',
  risk: '' as 'low' | 'medium' | 'high' | '',
  grade: '' as string,
  opinion: '',
})

const rejectDialogVisible = ref(false)
const rejectForm = ref({ reason: '' })
const approveDialogVisible = ref(false)
const approveForm = ref({ 
  opinion: '',
  downloadUrl: '',
  versionNote: '',
})
const supplementDialogVisible = ref(false)
const supplementForm = ref({ content: '' })

const canEvaluate = computed(() => {
  return detail.value && (detail.value.status === 'processing' || detail.value.status === 'pending')
})

const hasEvaluated = computed(() => {
  return detail.value && !!detail.value.finalMatch
})

onMounted(() => {
  const id = route.params.id as string
  detail.value = mockIntroduceDetails.find(i => i.id === id) || null
  if (!detail.value) {
    ElMessage.error('未找到该审批记录')
  }
  // 初始化推荐等级
  if (detail.value) {
    evalForm.value.grade = detail.value.grade
  }
})

function goBack() {
  router.back()
}

function getStatusType(status: string) {
  const map: Record<string, string> = {
    pending: '',
    processing: 'warning',
    completed: 'success',
    rejected: 'danger',
  }
  return map[status] || ''
}

function getStatusLabel(status: string) {
  return TICKET_STATUS_LABELS[status as keyof typeof TICKET_STATUS_LABELS] || status
}

function getTimelineType(action: string) {
  if (action.includes('完成')) return 'success'
  if (action.includes('驳回')) return 'danger'
  if (action.includes('通过')) return 'success'
  return ''
}

function getComponentTypeLabel(type: string) {
  return SOFTWARE_TYPE_LABELS[type as keyof typeof SOFTWARE_TYPE_LABELS] || type
}

function getScoreClass(score: number) {
  if (score >= 80) return 'score-success'
  if (score >= 60) return 'score-warning'
  return 'score-danger'
}

function getMatchType(match?: string) {
  const map: Record<string, string> = { match: 'success', partial: 'warning', unmatch: 'danger' }
  return map[match || ''] || 'info'
}

function getMatchLabel(match?: string) {
  const map: Record<string, string> = { match: '完全匹配', partial: '部分匹配', unmatch: '不匹配' }
  return map[match || ''] || '-'
}

function getApplicabilityType(app?: string) {
  const map: Record<string, string> = { apply: 'success', conditional: 'warning', unapply: 'danger' }
  return map[app || ''] || 'info'
}

function getApplicabilityLabel(app?: string) {
  const map: Record<string, string> = { apply: '适用', conditional: '有条件适用', unapply: '不适用' }
  return map[app || ''] || '-'
}

function getRiskType(risk?: string) {
  const map: Record<string, string> = { low: 'success', medium: 'warning', high: 'danger' }
  return map[risk || ''] || 'info'
}

function getRiskLabel(risk?: string) {
  const map: Record<string, string> = { low: '低风险', medium: '中风险', high: '高风险' }
  return map[risk || ''] || '-'
}

function handleApprove() {
  if (!evalForm.value.risk) {
    ElMessage.warning('请选择最终风险确认')
    return
  }
  if (!evalForm.value.grade) {
    ElMessage.warning('请选择最终等级')
    return
  }
  ElMessageBox.confirm('确认审批通过？任务将流转至下一节点。', '审批确认', {
    confirmButtonText: '确认通过',
    cancelButtonText: '取消',
    type: 'success',
  }).then(() => {
    if (detail.value) {
      detail.value.status = 'completed'
      detail.value.currentStep = '审批完成'
      detail.value.finalMatch = evalForm.value.match
      detail.value.finalApplicability = evalForm.value.applicability
      detail.value.finalRisk = evalForm.value.risk as 'low' | 'medium' | 'high'
      detail.value.finalGrade = evalForm.value.grade as 'A' | 'B' | 'C' | 'D'
      detail.value.finalOpinion = evalForm.value.opinion
      detail.value.approvalHistory = [
        ...(detail.value.approvalHistory || []),
        {
          action: '业务终评通过',
          operator: '管理员',
          opinion: evalForm.value.opinion || '审批通过',
          time: new Date().toLocaleString('zh-CN'),
        },
        {
          action: '审批完成',
          operator: '系统',
          opinion: '流程完成，组件已入库',
          time: new Date().toLocaleString('zh-CN'),
        },
      ]
    }
    ElMessage.success('审批通过，任务已流转至下一节点')
  }).catch(() => {})
}

function handleApproveWithDownload() {
  if (!evalForm.value.risk) {
    ElMessage.warning('请选择最终风险确认')
    return
  }
  if (!evalForm.value.grade) {
    ElMessage.warning('请选择最终等级')
    return
  }
  approveForm.value = {
    opinion: evalForm.value.opinion,
    downloadUrl: '',
    versionNote: '',
  }
  approveDialogVisible.value = true
}

function confirmApprove() {
  if (!approveForm.value.opinion?.trim()) {
    ElMessage.warning('请填写审批意见')
    return
  }
  if (!approveForm.value.downloadUrl?.trim()) {
    ElMessage.warning('请填写下载地址')
    return
  }
  
  if (detail.value) {
    detail.value.status = 'completed'
    detail.value.currentStep = '审批完成'
    detail.value.finalMatch = evalForm.value.match
    detail.value.finalApplicability = evalForm.value.applicability
    detail.value.finalRisk = evalForm.value.risk as 'low' | 'medium' | 'high'
    detail.value.finalGrade = evalForm.value.grade as 'A' | 'B' | 'C' | 'D'
    detail.value.finalOpinion = approveForm.value.opinion
    detail.value.downloadUrl = approveForm.value.downloadUrl
    detail.value.approvalHistory = [
      ...(detail.value.approvalHistory || []),
      {
        action: '业务终评通过',
        operator: '管理员',
        opinion: approveForm.value.opinion,
        time: new Date().toLocaleString('zh-CN'),
      },
      {
        action: '审批完成（已入库）',
        operator: '系统',
        opinion: `流程完成，组件已入库，下载地址：${approveForm.value.downloadUrl}`,
        time: new Date().toLocaleString('zh-CN'),
      },
    ]
  }
  
  approveDialogVisible.value = false
  ElMessage.success('审批通过，组件已入库并生成下载地址')
}

function handleReject() {
  rejectForm.value.reason = ''
  rejectDialogVisible.value = true
}

function confirmReject() {
  if (!rejectForm.value.reason.trim()) {
    ElMessage.warning('请填写驳回原因')
    return
  }
  if (detail.value) {
    detail.value.status = 'rejected'
    detail.value.currentStep = '已驳回'
    detail.value.approvalHistory = [
      ...(detail.value.approvalHistory || []),
      {
        action: '驳回',
        operator: '管理员',
        opinion: rejectForm.value.reason,
        time: new Date().toLocaleString('zh-CN'),
      },
    ]
  }
  rejectDialogVisible.value = false
  ElMessage.success('已驳回申请，任务已退回申请人')
}

function handleRequestSupplement() {
  supplementForm.value.content = ''
  supplementDialogVisible.value = true
}

function confirmSupplement() {
  if (!supplementForm.value.content.trim()) {
    ElMessage.warning('请填写需补充的内容')
    return
  }
  if (detail.value) {
    detail.value.status = 'processing'
    detail.value.currentStep = '待补充材料'
    detail.value.approvalHistory = [
      ...(detail.value.approvalHistory || []),
      {
        action: '要求补充材料',
        operator: '管理员',
        opinion: `需补充：${supplementForm.value.content}`,
        time: new Date().toLocaleString('zh-CN'),
      },
    ]
  }
  supplementDialogVisible.value = false
  ElMessage.success('已要求申请人补充材料，任务已退回')
}

function generateWord() {
  ElMessage.success('正在生成Word文档...')
  setTimeout(() => {
    ElMessage.success(`审批文档「${detail.value?.applicationNo}_${detail.value?.componentName}」已生成并开始下载`)
  }, 1000)
}
</script>

<style scoped>
.introduce-detail-page {
  padding: 4px;
}

.page-header {
  margin-bottom: 20px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
}

.section-card {
  margin-bottom: 16px;
  border: 1px solid #e5e6eb;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
}

.timeline-content {
  padding: 4px 0;
}

.timeline-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 4px;
}

.timeline-action {
  font-weight: 600;
  color: #1d2129;
}

.timeline-operator {
  color: #409EFF;
  font-size: 13px;
}

.timeline-time {
  color: #86909c;
  font-size: 12px;
  margin-left: auto;
}

.timeline-opinion {
  margin: 0;
  font-size: 13px;
  color: #606266;
  background: #f7f8fa;
  padding: 8px 12px;
  border-radius: 6px;
}

.eval-summary {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 8px;
}

.eval-metric {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.metric-label {
  font-size: 13px;
  color: #86909c;
}

.metric-value {
  font-size: 24px;
  font-weight: 600;
}

.score-success {
  color: #67C23A;
}

.score-warning {
  color: #E6A23C;
}

.score-danger {
  color: #F56C6C;
}

.grade-hint {
  margin-left: 12px;
  font-size: 12px;
  color: #86909c;
}

.action-bar {
  display: flex;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid #e5e6eb;
}

.doc-hint {
  margin-left: 12px;
  font-size: 13px;
  color: #86909c;
}

:deep(.el-descriptions__label) {
  font-weight: 500;
}
</style>
