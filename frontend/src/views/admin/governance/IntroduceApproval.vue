<template>
  <div class="introduce-approval-page">
    <div class="page-header">
      <h2>引入审批中心</h2>
      <p class="page-desc">集中处理开源组件引入申请，涵盖预评估、业务终评、等级核定等全流程审批</p>
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
        <el-form-item label="安全等级">
          <el-select v-model="filterForm.grade" placeholder="全部" clearable style="width: 120px">
            <el-option label="A - 核心关键" value="A" />
            <el-option label="B - 重要支撑" value="B" />
            <el-option label="C - 一般使用" value="C" />
            <el-option label="D - 探索试用" value="D" />
          </el-select>
        </el-form-item>
        <el-form-item label="申请人">
          <el-input v-model="filterForm.applicant" placeholder="请输入申请人" clearable style="width: 140px" />
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
            <span class="stat-label">总申请数</span>
            <span class="stat-value">{{ introduceList.length }}</span>
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
            <span class="stat-label">本月完成</span>
            <span class="stat-value">{{ completedCount }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 列表 -->
    <el-card shadow="never" class="table-card">
      <el-table :data="filteredData" stripe style="width: 100%">
        <el-table-column prop="applicationNo" label="申请编号" width="160">
          <template #default="{ row }">
            <el-link type="primary" @click="viewDetail(row.id)">{{ row.applicationNo }}</el-link>
          </template>
        </el-table-column>
        <el-table-column label="组件信息" min-width="200">
          <template #default="{ row }">
            <div class="component-cell">
              <span class="component-name">{{ row.componentName }}</span>
              <el-tag size="small" type="info">v{{ row.componentVersion }}</el-tag>
              <el-tag size="small" :type="(getComponentTypeTag(row.componentType) as any)" effect="plain">
                {{ getComponentTypeLabel(row.componentType) }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="applicant" label="申请人" width="100" />
        <el-table-column prop="department" label="部门" width="120" />
        <el-table-column label="核定等级" width="100">
          <template #default="{ row }">
            <GradeTag :grade="row.grade" size="small" />
          </template>
        </el-table-column>
        <el-table-column label="预评估评分" width="110">
          <template #default="{ row }">
            <el-tag
              :type="row.preEvalScore >= 80 ? 'success' : row.preEvalScore >= 60 ? 'warning' : 'danger'"
              size="small"
            >
              {{ row.preEvalScore }}分
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="预评估结论" width="100">
          <template #default="{ row }">
            <el-tag
              :type="row.preEvalConclusion === 'pass' ? 'success' : row.preEvalConclusion === 'warning' ? 'warning' : 'danger'"
              size="small"
            >
              {{ row.preEvalConclusion === 'pass' ? '通过' : row.preEvalConclusion === 'warning' ? '风险警告' : '不通过' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="currentStep" label="当前步骤" width="100" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <TicketStatusTag :status="row.status" />
          </template>
        </el-table-column>
        <el-table-column prop="submittedAt" label="提交时间" width="120">
          <template #default="{ row }">{{ formatTime(row.submittedAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewDetail(row.id)">查看详情</el-button>
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Refresh } from '@element-plus/icons-vue'
import { SOFTWARE_TYPE_LABELS } from '@/types'
import { formatTime } from '@/utils/format'

const router = useRouter()
const currentPage = ref(1)
const pageSize = ref(10)
const introduceList = ref<any[]>([])

const filterForm = ref({
  status: '',
  grade: '',
  applicant: '',
  componentName: '',
})

const filteredData = computed(() => {
  let data = [...introduceList.value]
  if (filterForm.value.status) {
    data = data.filter(i => i.status === filterForm.value.status)
  }
  if (filterForm.value.grade) {
    data = data.filter(i => i.grade === filterForm.value.grade)
  }
  if (filterForm.value.applicant) {
    data = data.filter(i => i.applicant.includes(filterForm.value.applicant))
  }
  if (filterForm.value.componentName) {
    data = data.filter(i => i.componentName.toLowerCase().includes(filterForm.value.componentName.toLowerCase()))
  }
  return data
})

const pendingCount = computed(() => introduceList.value.filter(i => i.status === 'pending').length)
const processingCount = computed(() => introduceList.value.filter(i => i.status === 'processing').length)
const completedCount = computed(() => introduceList.value.filter(i => i.status === 'completed').length)

function loadIntroduceList() {}

onMounted(() => {
  loadIntroduceList()
})

function handleFilter() {
  currentPage.value = 1
}

function handleReset() {
  filterForm.value = { status: '', grade: '', applicant: '', componentName: '' }
  currentPage.value = 1
}

function viewDetail(id: string) {
  router.push(`/admin/governance/introduce/${id}`)
}

function getComponentTypeLabel(type: string) {
  return SOFTWARE_TYPE_LABELS[type as keyof typeof SOFTWARE_TYPE_LABELS] || type
}

function getComponentTypeTag(type: string) {
  const map: Record<string, string> = {
    middleware: '',
    database: 'success',
    tool: 'warning',
    dependency: 'info',
    library: '',
  }
  return map[type] || ''
}
</script>

<style scoped>
.introduce-approval-page {
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
</style>
