<template>
  <div class="ticket-tracking">
    <h3>工单与整改追踪</h3>
    <p class="desc">集中查询已生成的缺陷风险工单状态、修复进度及逾期情况</p>

    <div style="margin-top: 16px">
      <el-button type="primary" plain><el-icon><Download /></el-icon> 导出报表</el-button>
    </div>

    <el-table :data="tracking" stripe style="margin-top: 16px">
      <el-table-column prop="ticketNo" label="工单编号" width="150" />
      <el-table-column prop="title" label="标题" min-width="160" />
      <el-table-column prop="type" label="类型" width="100" />
      <el-table-column label="风险等级" width="100">
        <template #default="{ row }"><RiskTag :level="row.riskLevel" /></template>
      </el-table-column>
      <el-table-column prop="assignee" label="负责人" width="100" />
      <el-table-column prop="department" label="部门" width="120" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }"><TicketStatusTag :status="row.status" /></template>
      </el-table-column>
      <el-table-column label="进度" width="160">
        <template #default="{ row }">
          <el-progress :percentage="row.progress" :color="row.progress === 100 ? '#67C23A' : '#409EFF'" />
        </template>
      </el-table-column>
      <el-table-column prop="dueDate" label="截止日期" width="110" />
      <el-table-column label="逾期" width="80">
        <template #default="{ row }">
          <el-tag v-if="row.isOverdue" type="danger" size="small">逾期</el-tag>
          <span v-else class="text-ok">正常</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { mockTicketTracking } from '@/mock/admin'
import { Download } from '@element-plus/icons-vue'

const tracking = mockTicketTracking
</script>

<style scoped>
.desc { color: #909399; font-size: 13px; margin: 8px 0 0; }
.text-ok { color: #67C23A; }
</style>
