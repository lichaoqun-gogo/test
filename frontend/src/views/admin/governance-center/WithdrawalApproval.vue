<template>
  <div class="withdrawal-approval">
    <h3>退出审批中心</h3>
    <el-table :data="withdrawals" stripe style="margin-top: 16px">
      <el-table-column prop="applicationNo" label="申请编号" width="160" />
      <el-table-column prop="componentName" label="组件名称" width="140" />
      <el-table-column label="当前等级" width="100">
        <template #default="{ row }"><GradeTag :grade="row.currentGrade" size="small" /></template>
      </el-table-column>
      <el-table-column prop="disposalMethod" label="处置方式" width="140" />
      <el-table-column prop="applicant" label="申请人" width="100" />
      <el-table-column prop="reason" label="退出原因" min-width="160" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }"><TicketStatusTag :status="row.status" /></template>
      </el-table-column>
      <el-table-column prop="submittedAt" label="提交时间" width="180">
        <template #default="{ row }">{{ formatTime(row.submittedAt) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status === 'pending'" type="success" size="small" @click="handleApprove(row)">通过</el-button>
          <el-button v-if="row.status === 'pending'" type="danger" size="small">驳回</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { formatTime } from '@/utils/format'

const withdrawals = ref<any[]>([])

function loadWithdrawals() {}

onMounted(() => {
  loadWithdrawals()
})

function handleApprove(row: any) {
  row.status = 'completed'
  ElMessage.success('审批通过')
}
</script>
