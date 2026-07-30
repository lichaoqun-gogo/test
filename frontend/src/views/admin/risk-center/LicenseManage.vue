<template>
  <div class="license-manage">
    <h3>许可证信息库</h3>
    <p class="desc">集中维护公司认可的许可证清单及风险等级，作为全局字典</p>

    <el-table :data="licenses" stripe style="margin-top: 16px">
      <el-table-column prop="name" label="许可证" width="140" />
      <el-table-column prop="fullName" label="全称" min-width="200" />
      <el-table-column prop="category" label="类型" width="100" />
      <el-table-column label="风险等级" width="100">
        <template #default="{ row }"><RiskTag :level="row.riskLevel" /></template>
      </el-table-column>
      <el-table-column label="是否认可" width="100">
        <template #default="{ row }">
          <el-tag :type="row.isApproved ? 'success' : 'danger'" size="small">
            {{ row.isApproved ? '认可' : '不认可' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="usageCount" label="使用次数" width="100" />
      <el-table-column prop="description" label="描述" min-width="200" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { mockLicenses } from '@/mock/admin'
const licenses = mockLicenses
</script>

<style scoped>
.desc { color: #909399; font-size: 13px; margin: 8px 0 0; }
</style>
