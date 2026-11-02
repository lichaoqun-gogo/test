<template>
  <div class="upgrade-review">
    <h3>升级评审管理</h3>
    <el-table :data="upgrades" stripe style="margin-top: 16px">
      <el-table-column prop="applicationNo" label="申请编号" width="160" />
      <el-table-column prop="componentName" label="组件名称" width="140" />
      <el-table-column label="升级前" width="100">
        <template #default="{ row }">v{{ row.fromVersion }}</template>
      </el-table-column>
      <el-table-column label="升级后" width="100">
        <template #default="{ row }"><el-tag type="success" size="small">v{{ row.toVersion }}</el-tag></template>
      </el-table-column>
      <el-table-column label="当前等级" width="100">
        <template #default="{ row }"><GradeTag :grade="row.currentGrade" size="small" /></template>
      </el-table-column>
      <el-table-column label="目标等级" width="100">
        <template #default="{ row }"><GradeTag :grade="row.targetGrade" size="small" /></template>
      </el-table-column>
      <el-table-column prop="applicant" label="申请人" width="100" />
      <el-table-column prop="reason" label="升级原因" min-width="160" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }"><TicketStatusTag :status="row.status" /></template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status === 'pending'" type="success" size="small">通过</el-button>
          <el-button v-if="row.status === 'pending'" type="danger" size="small">驳回</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const upgrades = ref<any[]>([])

function loadUpgrades() {}

onMounted(() => {
  loadUpgrades()
})
</script>
