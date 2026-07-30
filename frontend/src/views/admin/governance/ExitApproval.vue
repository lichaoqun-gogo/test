<template>
  <div class="exit-approval-page">
    <div class="page-header">
      <h2>退出审批中心</h2>
      <p class="page-desc">集中处理软件退出申请，审核原安全等级及处置方式，确保退出流程规范有序</p>
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
        <el-form-item label="当前等级">
          <el-select v-model="filterForm.grade" placeholder="全部" clearable style="width: 120px">
            <el-option label="A" value="A" />
            <el-option label="B" value="B" />
            <el-option label="C" value="C" />
            <el-option label="D" value="D" />
          </el-select>
        </el-form-item>
        <el-form-item label="处置方式">
          <el-select v-model="filterForm.disposal" placeholder="全部" clearable style="width: 140px">
            <el-option label="迁移替代" value="迁移" />
            <el-option label="停止使用" value="停止" />
            <el-option label="替换" value="替换" />
          </el-select>
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
            <span class="stat-label">总退出申请</span>
            <span class="stat-value">{{ mockExitDetails.length }}</span>
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
            <span class="stat-label">已退出</span>
            <span class="stat-value">{{ completedCount }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 退出申请列表 -->
    <el-card shadow="never" class="table-card">
      <el-table :data="filteredData" stripe style="width: 100%">
        <el-table-column prop="applicationNo" label="申请编号" width="160">
          <template #default="{ row }">
            <el-link type="primary" @click="viewDetail(row.id)">{{ row.applicationNo }}</el-link>
          </template>
        </el-table-column>
        <el-table-column label="组件信息" min-width="180">
          <template #default="{ row }">
            <div class="component-cell">
              <span class="component-name">{{ row.componentName }}</span>
              <el-tag size="small" type="info">v{{ row.componentVersion }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="原安全等级" width="100">
          <template #default="{ row }">
            <GradeTag :grade="row.currentGrade" size="small" />
          </template>
        </el-table-column>
        <el-table-column prop="disposalMethod" label="处置方式" width="140" />
        <el-table-column label="影响项目" min-width="160">
          <template #default="{ row }">
            <el-tag v-for="proj in row.affectedProjects.slice(0, 2)" :key="proj" size="small" style="margin-right: 4px">
              {{ proj }}
            </el-tag>
            <el-tag v-if="row.affectedProjects.length > 2" size="small" type="info">
              +{{ row.affectedProjects.length - 2 }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="applicant" label="申请人" width="100" />
        <el-table-column prop="department" label="部门" width="120" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <TicketStatusTag :status="row.status" />
          </template>
        </el-table-column>
        <el-table-column prop="submittedAt" label="提交时间" width="120" />
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

    <!-- 退出详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="退出申请详情" width="700px">
      <template v-if="currentDetail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="申请编号">{{ currentDetail.applicationNo }}</el-descriptions-item>
          <el-descriptions-item label="组件名称">{{ currentDetail.componentName }}</el-descriptions-item>
          <el-descriptions-item label="当前版本">v{{ currentDetail.componentVersion }}</el-descriptions-item>
          <el-descriptions-item label="原安全等级">
            <GradeTag :grade="currentDetail.currentGrade" />
          </el-descriptions-item>
          <el-descriptions-item label="处置方式">{{ currentDetail.disposalMethod }}</el-descriptions-item>
          <el-descriptions-item label="申请人">{{ currentDetail.applicant }}</el-descriptions-item>
          <el-descriptions-item label="退出原因" :span="2">{{ currentDetail.reason }}</el-descriptions-item>
          <el-descriptions-item label="影响项目" :span="2">
            <el-tag v-for="proj in currentDetail.affectedProjects" :key="proj" size="small" style="margin-right: 6px">
              {{ proj }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="替代方案" :span="2">{{ currentDetail.alternativePlan }}</el-descriptions-item>
          <el-descriptions-item label="代码清理计划" :span="2">{{ currentDetail.codeCleanupPlan }}</el-descriptions-item>
          <el-descriptions-item label="风险评估" :span="2">{{ currentDetail.riskAssessment }}</el-descriptions-item>
        </el-descriptions>

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
                <span class="timeline-time">{{ item.time }}</span>
              </div>
              <p class="timeline-opinion">{{ item.opinion }}</p>
            </el-timeline-item>
          </el-timeline>
        </div>
      </template>
    </el-dialog>

    <!-- 驳回对话框 -->
    <el-dialog v-model="rejectDialogVisible" title="驳回退出申请" width="500px">
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
import { Search, Refresh } from '@element-plus/icons-vue'
import { mockExitDetails } from '@/mock/admin'
import type { ExitApprovalDetail } from '@/mock/admin'

const router = useRouter()
const currentPage = ref(1)
const pageSize = ref(10)

const filterForm = ref({
  status: '',
  grade: '',
  disposal: '',
})

const filteredData = computed(() => {
  let data = [...mockExitDetails]
  if (filterForm.value.status) {
    data = data.filter(i => i.status === filterForm.value.status)
  }
  if (filterForm.value.grade) {
    data = data.filter(i => i.currentGrade === filterForm.value.grade)
  }
  if (filterForm.value.disposal) {
    data = data.filter(i => i.disposalMethod.includes(filterForm.value.disposal))
  }
  return data
})

const pendingCount = computed(() => mockExitDetails.filter(i => i.status === 'pending').length)
const processingCount = computed(() => mockExitDetails.filter(i => i.status === 'processing').length)
const completedCount = computed(() => mockExitDetails.filter(i => i.status === 'completed').length)

const detailDialogVisible = ref(false)
const currentDetail = ref<ExitApprovalDetail | null>(null)
const rejectDialogVisible = ref(false)
const rejectForm = ref({ reason: '' })
const currentRejectRow = ref<any>(null)

function handleFilter() {
  currentPage.value = 1
}

function handleReset() {
  filterForm.value = { status: '', grade: '', disposal: '' }
  currentPage.value = 1
}

function viewDetail(id: string) {
  const item = mockExitDetails.find(i => i.id === id)
  if (item) {
    currentDetail.value = item
    detailDialogVisible.value = true
  }
}

function handleApprove(row: any) {
  ElMessageBox.confirm(`确认通过「${row.componentName}」的退出申请？`, '审批确认', {
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
        opinion: '同意退出，请按计划执行',
        time: new Date().toLocaleString('zh-CN'),
      },
    ]
    ElMessage.success('审批通过')
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
  ElMessage.success('已驳回退出申请')
}
</script>

<style scoped>
.exit-approval-page {
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

.component-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.component-name {
  font-weight: 500;
  color: #1d2129;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.approval-history {
  margin-top: 20px;
}

.approval-history h4 {
  margin: 0 0 12px;
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
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
