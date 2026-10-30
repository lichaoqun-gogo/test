<template>
  <div class="governance-dashboard">
    <h3>治理效能看板</h3>

    <el-card style="margin-top: 16px">
      <template #header>缺陷修复效能（MTTR）</template>
      <el-table :data="dashboardData.repairStats" size="small">
        <el-table-column prop="department" label="部门" />
        <el-table-column label="平均修复时长" width="120">
          <template #default="{ row }">{{ row.mttr }} 天</template>
        </el-table-column>
        <el-table-column label="逾期未修复" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.overdue > 0" type="danger" size="small">{{ row.overdue }}</el-tag>
            <span v-else>0</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { mockDashboardData } from '@/mock/admin'

const dashboardData = mockDashboardData
</script>
