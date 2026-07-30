<template>
  <div class="violation-audit">
    <h3>违规使用审计</h3>
    <p class="desc">审计未走流程私自引入的软件</p>

    <el-table :data="violations" stripe style="margin-top: 16px">
      <el-table-column prop="componentName" label="组件名称" width="140" />
      <el-table-column prop="version" label="版本" width="100" />
      <el-table-column prop="project" label="所属项目" width="140" />
      <el-table-column prop="introducedBy" label="引入人" width="100" />
      <el-table-column prop="department" label="部门" width="120" />
      <el-table-column prop="violationType" label="违规类型" min-width="160" />
      <el-table-column label="风险等级" width="100">
        <template #default="{ row }"><RiskTag :level="row.riskLevel" /></template>
      </el-table-column>
      <el-table-column label="处置状态" width="100">
        <template #default="{ row }">
          <el-tag :type="violationStatusType(row.status)" size="small">{{ violationStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="discoveredDate" label="发现日期" width="110" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { mockViolations } from '@/mock/admin'

const violations = mockViolations

function violationStatusType(status: string) {
  const types: Record<string, string> = { pending: 'warning', handling: '', resolved: 'success' }
  return types[status] || ''
}

function violationStatusLabel(status: string) {
  const labels: Record<string, string> = { pending: '待处理', handling: '处理中', resolved: '已处置' }
  return labels[status] || status
}
</script>

<style scoped>
.desc { color: #909399; font-size: 13px; margin: 8px 0 0; }
</style>
