<template>
  <div class="introduction-approval">
    <h3>引入审批中心</h3>
    <div class="filter-bar" style="margin-top: 12px">
      <el-select v-model="filterStatus" placeholder="按状态筛选" clearable style="width: 150px">
        <el-option label="待审批" value="pending" />
        <el-option label="处理中" value="processing" />
        <el-option label="已完成" value="completed" />
        <el-option label="已驳回" value="rejected" />
      </el-select>
    </div>

    <el-table :data="filteredIntroductions" stripe style="margin-top: 16px">
      <el-table-column prop="applicationNo" label="申请编号" width="160">
        <template #default="{ row }">
          <el-link type="primary" @click="viewDetail(row.id)">{{ row.applicationNo }}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="组件" width="160">
        <template #default="{ row }">
          {{ row.componentName }} <el-tag size="small" type="info">{{ row.componentVersion }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="applicant" label="申请人" width="100" />
      <el-table-column prop="department" label="部门" width="120" />
      <el-table-column label="预评估评分" width="110">
        <template #default="{ row }">
          <el-tag :type="row.preEvalScore >= 80 ? 'success' : row.preEvalScore >= 60 ? 'warning' : 'danger'" size="small">
            {{ row.preEvalScore }}分
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="preEvalConclusion" label="预评估结论" width="100" />
      <el-table-column prop="currentStep" label="当前步骤" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }"><TicketStatusTag :status="row.status" /></template>
      </el-table-column>
      <el-table-column prop="submittedAt" label="提交时间" width="110" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { mockIntroductions } from '@/mock/admin'

const router = useRouter()
const filterStatus = ref('')

const filteredIntroductions = computed(() => {
  if (!filterStatus.value) return mockIntroductions
  return mockIntroductions.filter(i => i.status === filterStatus.value)
})

function viewDetail(id: string) {
  router.push(`/admin/governance-center/introduction/${id}`)
}
</script>
