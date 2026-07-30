<template>
  <div class="inventory-manage">
    <h3>存量软件盘点</h3>
    <p class="desc">对比实际代码库与台账，发现未报备的"影子IT"资产，按A/B/C/D标准重新定级</p>

    <el-table :data="inventory" stripe style="margin-top: 16px">
      <el-table-column prop="name" label="组件名称" width="140" />
      <el-table-column prop="version" label="版本" width="100" />
      <el-table-column prop="project" label="所属项目" width="140" />
      <el-table-column prop="discoveredBy" label="发现方式" width="100" />
      <el-table-column label="等级" width="100">
        <template #default="{ row }"><GradeTag :grade="row.grade" size="small" /></template>
      </el-table-column>
      <el-table-column label="风险等级" width="100">
        <template #default="{ row }"><RiskTag :level="row.riskLevel" /></template>
      </el-table-column>
      <el-table-column label="已报备" width="80">
        <template #default="{ row }">
          <el-tag :type="row.isReported ? 'success' : 'danger'" size="small">
            {{ row.isReported ? '是' : '否' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="inventoryStatusType(row.status)" size="small">{{ inventoryStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="discoveredDate" label="发现日期" width="110" />
      <el-table-column label="操作" width="100">
        <template #default="{ row }">
          <el-button v-if="row.status === 'pending_review'" type="primary" link size="small">定级</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { mockInventory } from '@/mock/admin'

const inventory = mockInventory

function inventoryStatusType(status: string) {
  const types: Record<string, string> = { pending_review: 'warning', graded: '', approved: 'success', rejected: 'danger' }
  return types[status] || ''
}

function inventoryStatusLabel(status: string) {
  const labels: Record<string, string> = { pending_review: '待审核', graded: '已定级', approved: '已通过', rejected: '已驳回' }
  return labels[status] || status
}
</script>

<style scoped>
.desc { color: #909399; font-size: 13px; margin: 8px 0 0; }
</style>
