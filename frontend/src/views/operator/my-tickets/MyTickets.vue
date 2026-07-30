<template>
  <div class="my-tickets">
    <div class="filter-bar">
      <el-select v-model="filterType" placeholder="按工单类型筛选" clearable style="width: 150px">
        <el-option label="引入申请" value="introduce" />
        <el-option label="升级申请" value="upgrade" />
        <el-option label="退出申请" value="withdraw" />
        <el-option label="风险上报" value="risk_report" />
        <el-option label="分配的风险工单" value="risk" />
      </el-select>
      <el-select v-model="filterStatus" placeholder="按状态筛选" clearable style="width: 150px; margin-left: 12px">
        <el-option label="待审批" value="pending" />
        <el-option label="处理中" value="processing" />
        <el-option label="已完成" value="completed" />
        <el-option label="已驳回" value="rejected" />
        <el-option label="已分配" value="assigned" />
      </el-select>
    </div>

    <el-table :data="filteredTickets" stripe style="margin-top: 16px">
      <el-table-column prop="ticketNo" label="工单编号" width="150">
        <template #default="{ row }">
          <el-link type="primary" @click="viewDetail(row.id)">{{ row.ticketNo }}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="类型" width="110">
        <template #default="{ row }">
          <el-tag size="small">
            {{ APPLICATION_TYPE_LABELS[row.type as keyof typeof APPLICATION_TYPE_LABELS] }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column v-if="hasRiskLevel" label="风险等级" width="100">
        <template #default="{ row }">
          <RiskTag v-if="row.riskLevel" :level="row.riskLevel" size="small" />
        </template>
      </el-table-column>
      <el-table-column prop="applicant" label="申请人" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <TicketStatusTag :status="row.status" />
        </template>
      </el-table-column>
      <el-table-column prop="currentStep" label="当前步骤" width="120" />
      <el-table-column prop="createdAt" label="创建时间" width="120" />
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.externalLink" type="primary" link size="small" @click="openExternal(row.externalLink!)">
            跳转工单
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { mockTickets } from '@/mock/tickets'
import { APPLICATION_TYPE_LABELS } from '@/types'

const router = useRouter()
const filterType = ref('')
const filterStatus = ref('')

const hasRiskLevel = computed(() => mockTickets.some(t => t.riskLevel))

const filteredTickets = computed(() => {
  let list = mockTickets
  if (filterType.value) {
    list = list.filter(t => t.type === filterType.value)
  }
  if (filterStatus.value) {
    list = list.filter(t => t.status === filterStatus.value)
  }
  return list
})

function viewDetail(id: string) {
  router.push(`/operator/my-tickets/${id}`)
}

function openExternal(url: string) {
  window.open(url, '_blank')
}
</script>

<style scoped>
.filter-bar {
  display: flex;
  align-items: center;
}
</style>
