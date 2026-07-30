<template>
  <div class="ticket-detail">
    <el-page-header @back="goBack" content="工单详情" />

    <el-descriptions :column="2" border style="margin-top: 16px">
      <el-descriptions-item label="工单编号">{{ ticket?.ticketNo }}</el-descriptions-item>
      <el-descriptions-item label="类型">
        {{ APPLICATION_TYPE_LABELS[ticket?.type as keyof typeof APPLICATION_TYPE_LABELS] }}
      </el-descriptions-item>
      <el-descriptions-item label="标题" :span="2">{{ ticket?.title }}</el-descriptions-item>
      <el-descriptions-item label="申请人">{{ ticket?.applicant }}</el-descriptions-item>
      <el-descriptions-item label="部门">{{ ticket?.applicantDept }}</el-descriptions-item>
      <el-descriptions-item label="优先级">
        <el-tag :type="ticket?.priority === 'urgent' ? 'danger' : ''" size="small">
          {{ ticket?.priority === 'urgent' ? '紧急' : '普通' }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="风险等级">
        <RiskTag v-if="ticket?.riskLevel" :level="ticket.riskLevel" />
        <span v-else>-</span>
      </el-descriptions-item>
      <el-descriptions-item label="状态" :span="2">
        <TicketStatusTag :status="ticket?.status || 'pending'" />
      </el-descriptions-item>
      <el-descriptions-item label="当前步骤">{{ ticket?.currentStep }}</el-descriptions-item>
      <el-descriptions-item label="分配给">{{ ticket?.assignedTo || '未分配' }}</el-descriptions-item>
      <el-descriptions-item label="创建时间">{{ ticket?.createdAt }}</el-descriptions-item>
      <el-descriptions-item label="更新时间">{{ ticket?.updatedAt }}</el-descriptions-item>
      <el-descriptions-item label="描述" :span="2">{{ ticket?.description }}</el-descriptions-item>
    </el-descriptions>

    <div v-if="ticket?.externalLink" style="margin-top: 16px">
      <el-button type="primary" @click="openExternal(ticket.externalLink!)">
        <el-icon><Link /></el-icon>
        跳转外部工单系统
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { mockTickets } from '@/mock/tickets'
import { APPLICATION_TYPE_LABELS } from '@/types'
import { Link } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const ticket = computed(() =>
  mockTickets.find(t => t.id === route.params.id)
)

function goBack() {
  router.back()
}

function openExternal(url: string) {
  window.open(url, '_blank')
}
</script>
