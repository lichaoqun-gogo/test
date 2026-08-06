<template>
  <div class="upgrade-review-page">
    <div class="page-header">
      <h2>升级评审管理</h2>
      <p class="page-desc">审批核心组件重大升级计划，对比升级前后安全指标，重新核定安全等级</p>
    </div>

    <!-- 筛选栏 -->
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="审批状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable style="width: 140px">
            <el-option label="待审批" value="pending" />
            <el-option label="处理中" value="processing" />
            <el-option label="已完成" value="completed" />
            <el-option label="已驳回" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item label="组件名称">
          <el-input v-model="filterForm.componentName" placeholder="请输入组件名称" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleFilter">
            <el-icon><Search /></el-icon> 查询
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon> 重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stat-row">
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-item">
            <span class="stat-label">总升级申请</span>
            <span class="stat-value">{{ mockUpgradeDetails.length }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-item stat-pending">
            <span class="stat-label">待审批</span>
            <span class="stat-value">{{ pendingCount }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-item stat-processing">
            <span class="stat-label">处理中</span>
            <span class="stat-value">{{ processingCount }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-item stat-completed">
            <span class="stat-label">已升级</span>
            <span class="stat-value">{{ completedCount }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 升级申请列表 -->
    <el-card shadow="never" class="table-card">
      <el-table :data="filteredData" stripe style="width: 100%">
        <el-table-column prop="applicationNo" label="申请编号" width="160">
          <template #default="{ row }">
            <el-link type="primary" @click="viewDetail(row.id)">{{ row.applicationNo }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="componentName" label="组件名称" width="140" />
        <el-table-column label="版本升级" width="180">
          <template #default="{ row }">
            <div class="version-cell">
              <span class="version-from">v{{ row.fromVersion }}</span>
              <el-icon class="arrow-icon"><Right /></el-icon>
              <span class="version-to">v{{ row.toVersion }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="当前等级" width="100">
          <template #default="{ row }">
            <GradeTag :grade="row.currentGrade" size="small" />
          </template>
        </el-table-column>
        <el-table-column label="目标等级" width="100">
          <template #default="{ row }">
            <GradeTag :grade="row.targetGrade" size="small" />
          </template>
        </el-table-column>
        <el-table-column label="安全风险变化" width="130">
          <template #default="{ row }">
            <div v-if="row.beforeEval && row.afterEval" class="risk-change">
              <span class="risk-from">{{ row.beforeEval.vulnerabilityCount }}个漏洞</span>
              <el-icon class="arrow-icon-small"><Right /></el-icon>
              <span class="risk-to" :class="row.afterEval.vulnerabilityCount === 0 ? 'risk-zero' : ''">
                {{ row.afterEval.vulnerabilityCount }}个漏洞
              </span>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="applicant" label="申请人" width="100" />
        <el-table-column prop="department" label="部门" width="120" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <TicketStatusTag :status="row.status" />
          </template>
        </el-table-column>
        <el-table-column prop="submittedAt" label="提交时间" width="120">
          <template #default="{ row }">{{ formatTime(row.submittedAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewDetail(row.id)">查看详情</el-button>
            <el-button
              v-if="row.status === 'pending' || row.status === 'processing'"
              type="success"
              link
              size="small"
              @click="handleApprove(row)"
            >
              通过
            </el-button>
            <el-button
              v-if="row.status === 'pending' || row.status === 'processing'"
              type="danger"
              link
              size="small"
              @click="handleReject(row)"
            >
              驳回
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50]"
          :total="filteredData.length"
          layout="total, sizes, prev, pager, next"
        />
      </div>
    </el-card>

    <!-- 升级详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="升级评审详情" width="800px">
      <template v-if="currentDetail">
        <!-- 基本信息 -->
        <el-descriptions :column="2" border>
          <el-descriptions-item label="申请编号">{{ currentDetail.applicationNo }}</el-descriptions-item>
          <el-descriptions-item label="组件名称">{{ currentDetail.componentName }}</el-descriptions-item>
          <el-descriptions-item label="升级前版本">v{{ currentDetail.fromVersion }}</el-descriptions-item>
          <el-descriptions-item label="升级后版本">v{{ currentDetail.toVersion }}</el-descriptions-item>
          <el-descriptions-item label="当前等级">
            <GradeTag :grade="currentDetail.currentGrade" />
          </el-descriptions-item>
          <el-descriptions-item label="目标等级">
            <GradeTag :grade="currentDetail.targetGrade" />
          </el-descriptions-item>
          <el-descriptions-item label="申请人">{{ currentDetail.applicant }}</el-descriptions-item>
          <el-descriptions-item label="部门">{{ currentDetail.department }}</el-descriptions-item>
          <el-descriptions-item label="升级原因" :span="2">{{ currentDetail.reason }}</el-descriptions-item>
          <el-descriptions-item label="兼容性风险" :span="2">{{ currentDetail.compatibilityRisk }}</el-descriptions-item>
          <el-descriptions-item label="测试报告" :span="2">{{ currentDetail.testReport }}</el-descriptions-item>
          <el-descriptions-item label="回滚方案" :span="2">{{ currentDetail.rollbackPlan }}</el-descriptions-item>
        </el-descriptions>

        <!-- 升级前后对比 -->
        <div v-if="currentDetail.beforeEval && currentDetail.afterEval" class="comparison-section">
          <h4>升级前后对比</h4>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-card shadow="never" class="comparison-card before-card">
                <template #header>
                  <div class="comparison-header">
                    <el-tag type="info">升级前 v{{ currentDetail.fromVersion }}</el-tag>
                  </div>
                </template>
                <div class="comparison-item">
                  <span class="comparison-label">已知漏洞数</span>
                  <span class="comparison-value danger">{{ currentDetail.beforeEval.vulnerabilityCount }}</span>
                </div>
                <div class="comparison-item">
                  <span class="comparison-label">健康度评分</span>
                  <span class="comparison-value">{{ currentDetail.beforeEval.healthScore }}分</span>
                </div>
                <div class="comparison-item">
                  <span class="comparison-label">许可证风险</span>
                  <span class="comparison-value">{{ currentDetail.beforeEval.licenseRisk }}</span>
                </div>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card shadow="never" class="comparison-card after-card">
                <template #header>
                  <div class="comparison-header">
                    <el-tag type="success">升级后 v{{ currentDetail.toVersion }}</el-tag>
                  </div>
                </template>
                <div class="comparison-item">
                  <span class="comparison-label">已知漏洞数</span>
                  <span class="comparison-value" :class="currentDetail.afterEval.vulnerabilityCount === 0 ? 'success' : 'warning'">
                    {{ currentDetail.afterEval.vulnerabilityCount }}
                  </span>
                </div>
                <div class="comparison-item">
                  <span class="comparison-label">健康度评分</span>
                  <span class="comparison-value">{{ currentDetail.afterEval.healthScore }}分</span>
                </div>
                <div class="comparison-item">
                  <span class="comparison-label">许可证风险</span>
                  <span class="comparison-value">{{ currentDetail.afterEval.licenseRisk }}</span>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <!-- 安全修复列表 -->
        <div v-if="currentDetail.securityFixes.length > 0" class="detail-section">
          <h4>安全修复</h4>
          <el-tag
            v-for="fix in currentDetail.securityFixes"
            :key="fix"
            type="danger"
            class="detail-tag"
          >
            {{ fix }}
          </el-tag>
        </div>

        <!-- 新特性列表 -->
        <div v-if="currentDetail.newFeatures.length > 0" class="detail-section">
          <h4>新增特性</h4>
          <el-tag
            v-for="feature in currentDetail.newFeatures"
            :key="feature"
            type="success"
            class="detail-tag"
          >
            {{ feature }}
          </el-tag>
        </div>

        <!-- 破坏性变更 -->
        <div v-if="currentDetail.breakingChanges.length > 0" class="detail-section">
          <h4>破坏性变更</h4>
          <el-tag
            v-for="change in currentDetail.breakingChanges"
            :key="change"
            type="warning"
            class="detail-tag"
          >
            {{ change }}
          </el-tag>
        </div>

        <!-- 审批历史 -->
        <div v-if="currentDetail.approvalHistory" class="approval-history">
          <h4>审批历史</h4>
          <el-timeline>
            <el-timeline-item
              v-for="(item, index) in currentDetail.approvalHistory"
              :key="index"
              :type="item.action.includes('通过') ? 'success' : item.action.includes('驳回') ? 'danger' : ''"
            >
              <div class="timeline-header">
                <span>{{ item.action }}</span>
                <span class="timeline-operator">{{ item.operator }}</span>
                <span class="timeline-time">{{ formatTime(item.time) }}</span>
              </div>
              <p class="timeline-opinion">{{ item.opinion }}</p>
            </el-timeline-item>
          </el-timeline>
        </div>
      </template>
    </el-dialog>

    <!-- 驳回对话框 -->
    <el-dialog v-model="rejectDialogVisible" title="驳回升级申请" width="500px">
      <el-form :model="rejectForm" label-width="80px">
        <el-form-item label="驳回原因" required>
          <el-input v-model="rejectForm.reason" type="textarea" :rows="4" placeholder="请输入驳回原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmReject">确认驳回</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Right } from '@element-plus/icons-vue'
import { mockUpgradeDetails } from '@/mock/admin'
import type { UpgradeReviewDetail } from '@/mock/admin'
import { formatTime } from '@/utils/format'

const router = useRouter()
const currentPage = ref(1)
const pageSize = ref(10)

const filterForm = ref({
  status: '',
  componentName: '',
})

const filteredData = computed(() => {
  let data = [...mockUpgradeDetails]
  if (filterForm.value.status) {
    data = data.filter(i => i.status === filterForm.value.status)
  }
  if (filterForm.value.componentName) {
    data = data.filter(i => i.componentName.toLowerCase().includes(filterForm.value.componentName.toLowerCase()))
  }
  return data
})

const pendingCount = computed(() => mockUpgradeDetails.filter(i => i.status === 'pending').length)
const processingCount = computed(() => mockUpgradeDetails.filter(i => i.status === 'processing').length)
const completedCount = computed(() => mockUpgradeDetails.filter(i => i.status === 'completed').length)

const detailDialogVisible = ref(false)
const currentDetail = ref<UpgradeReviewDetail | null>(null)
const rejectDialogVisible = ref(false)
const rejectForm = ref({ reason: '' })
const currentRejectRow = ref<any>(null)

function handleFilter() {
  currentPage.value = 1
}

function handleReset() {
  filterForm.value = { status: '', componentName: '' }
  currentPage.value = 1
}

function viewDetail(id: string) {
  const item = mockUpgradeDetails.find(i => i.id === id)
  if (item) {
    currentDetail.value = item
    detailDialogVisible.value = true
  }
}

function handleApprove(row: any) {
  ElMessageBox.confirm(`确认通过「${row.componentName}」从 v${row.fromVersion} 升级至 v${row.toVersion}？`, '审批确认', {
    confirmButtonText: '确认通过',
    cancelButtonText: '取消',
    type: 'success',
  }).then(() => {
    row.status = 'completed'
    row.approvalHistory = [
      ...(row.approvalHistory || []),
      {
        action: '审批通过',
        operator: '管理员',
        opinion: `同意升级至 v${row.toVersion}`,
        time: new Date().toLocaleString('zh-CN'),
      },
    ]
    ElMessage.success('审批通过，已允许升级')
  }).catch(() => {})
}

function handleReject(row: any) {
  currentRejectRow.value = row
  rejectForm.value.reason = ''
  rejectDialogVisible.value = true
}

function confirmReject() {
  if (!rejectForm.value.reason.trim()) {
    ElMessage.warning('请填写驳回原因')
    return
  }
  if (currentRejectRow.value) {
    currentRejectRow.value.status = 'rejected'
    currentRejectRow.value.approvalHistory = [
      ...(currentRejectRow.value.approvalHistory || []),
      {
        action: '驳回',
        operator: '管理员',
        opinion: rejectForm.value.reason,
        time: new Date().toLocaleString('zh-CN'),
      },
    ]
  }
  rejectDialogVisible.value = false
  ElMessage.success('已驳回升级申请')
}
</script>

<style scoped>
.upgrade-review-page {
  padding: 4px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0 0 8px;
  font-size: 20px;
  font-weight: 600;
  color: #1d2129;
}

.page-desc {
  margin: 0;
  font-size: 14px;
  color: #86909c;
}

.filter-card {
  margin-bottom: 16px;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
}

.stat-row {
  margin-bottom: 16px;
}

.stat-card {
  border: 1px solid #e5e6eb;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-label {
  font-size: 14px;
  color: #86909c;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #1d2129;
}

.stat-pending .stat-value {
  color: #409EFF;
}

.stat-processing .stat-value {
  color: #E6A23C;
}

.stat-completed .stat-value {
  color: #67C23A;
}

.table-card {
  border: 1px solid #e5e6eb;
}

.version-cell {
  display: flex;
  align-items: center;
  gap: 6px;
}

.version-from {
  color: #86909c;
  text-decoration: line-through;
}

.version-to {
  color: #67C23A;
  font-weight: 600;
}

.arrow-icon {
  color: #409EFF;
  font-size: 14px;
}

.arrow-icon-small {
  color: #409EFF;
  font-size: 12px;
}

.risk-change {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
}

.risk-from {
  color: #F56C6C;
}

.risk-to {
  color: #67C23A;
  font-weight: 600;
}

.risk-zero {
  color: #67C23A;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.comparison-section {
  margin-top: 20px;
}

.comparison-section h4,
.detail-section h4,
.approval-history h4 {
  margin: 0 0 12px;
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
}

.comparison-card {
  border: 1px solid #e5e6eb;
}

.before-card {
  border-left: 3px solid #909399;
}

.after-card {
  border-left: 3px solid #67C23A;
}

.comparison-header {
  display: flex;
  align-items: center;
}

.comparison-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px dashed #e5e6eb;
}

.comparison-item:last-child {
  border-bottom: none;
}

.comparison-label {
  font-size: 13px;
  color: #86909c;
}

.comparison-value {
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
}

.comparison-value.danger {
  color: #F56C6C;
}

.comparison-value.success {
  color: #67C23A;
}

.comparison-value.warning {
  color: #E6A23C;
}

.detail-section {
  margin-top: 16px;
}

.detail-tag {
  margin: 0 8px 8px 0;
}

.approval-history {
  margin-top: 20px;
}

.timeline-header {
  display: flex;
  align-items: center;
  gap: 12px;
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
  margin: 4px 0 0;
  font-size: 13px;
  color: #606266;
  background: #f7f8fa;
  padding: 8px 12px;
  border-radius: 6px;
}
</style>
