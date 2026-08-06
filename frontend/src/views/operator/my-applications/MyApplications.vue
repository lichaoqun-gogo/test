<template>
  <div class="my-applications">
    <div class="page-header">
      <el-button type="primary" @click="goToNewApplication">
        <el-icon><Plus /></el-icon>
        发起新申请
      </el-button>
    </div>

    <el-table :data="applications" stripe>
      <el-table-column prop="applicationNo" label="申请编号" width="160">
        <template #default="{ row }">
          <el-link type="primary" @click="viewDetail(row.id)">{{ row.applicationNo }}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="申请类型" width="110">
        <template #default="{ row }">
          <el-tag size="small">{{ APPLICATION_TYPE_LABELS[row.type as keyof typeof APPLICATION_TYPE_LABELS] }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="组件名称" width="150">
        <template #default="{ row }">
          {{ row.componentName }}
          <el-tag size="small" type="info" style="margin-left: 4px">{{ row.componentVersion }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="applicant" label="申请人" width="100" />
      <el-table-column label="预评估状态" width="120">
        <template #default="{ row }">
          <el-tag
            :type="preEvalStatusType(row.preEvalStatus)"
            size="small"
          >
            {{ PRE_EVAL_STATUS_LABELS[row.preEvalStatus] }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column v-if="rowHasConclusion" label="预评估结论" width="120">
        <template #default="{ row }">
          <span v-if="row.preEvalConclusion">{{ PRE_EVAL_CONCLUSION_LABELS[row.preEvalConclusion] }}</span>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column v-if="rowHasScore" label="评分" width="80">
        <template #default="{ row }">
          <el-tag v-if="row.preEvalScore !== undefined" :type="scoreTagType(row.preEvalScore)" size="small">
            {{ row.preEvalScore }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="工单号" width="140">
        <template #default="{ row }">
          <span v-if="row.ticketNo">{{ row.ticketNo }}</span>
          <span v-else class="text-muted">未生成</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <TicketStatusTag :status="row.status" />
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="创建时间" width="180">
        <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrapper">
      <el-pagination
        :current-page="1"
        :page-size="10"
        :total="applications.length"
        layout="total, prev, pager, next"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { APPLICATION_TYPE_LABELS, PRE_EVAL_STATUS_LABELS, PRE_EVAL_CONCLUSION_LABELS } from '@/types'
import { Plus } from '@element-plus/icons-vue'
import { formatTime } from '@/utils/format'

const router = useRouter()
const applications = ref<any[]>([])

const rowHasConclusion = computed(() => applications.value.some(a => a.preEvalConclusion))
const rowHasScore = computed(() => applications.value.some(a => a.preEvalScore !== undefined))

function loadApplications() {
  // TODO: 调用后端 API 获取申请列表
}

onMounted(() => {
  loadApplications()
})

function preEvalStatusType(status: string) {
  const types: Record<string, string> = { pending: 'info', scanning: '', completed: 'success', failed: 'danger' }
  return types[status] || ''
}

function scoreTagType(score: number) {
  if (score >= 80) return 'success'
  if (score >= 60) return 'warning'
  return 'danger'
}

function goToNewApplication() {
  router.push('/operator/my-applications/new')
}

function viewDetail(id: string) {
  router.push(`/operator/my-applications/${id}`)
}
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 16px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.text-muted {
  color: #C0C4CC;
}
</style>
