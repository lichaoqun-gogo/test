<template>
  <div class="risk-ledger">
    <h3>统一风险台账</h3>
    <p class="desc">集中展示所有外部推送和内部 SCA 发现的风险，支持手动分发</p>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="6">
        <el-statistic title="风险总数" :value="risks.length" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="待分发" :value="risks.filter(r => r.status === 'pending').length" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="已分发" :value="risks.filter(r => r.status === 'dispatched').length" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="已解决" :value="risks.filter(r => r.status === 'resolved').length" />
      </el-col>
    </el-row>

    <el-table :data="risks" stripe style="margin-top: 16px">
      <el-table-column prop="source" label="来源" width="100" />
      <el-table-column prop="title" label="风险标题" min-width="180" />
      <el-table-column label="风险等级" width="100">
        <template #default="{ row }"><RiskTag :level="row.riskLevel" /></template>
      </el-table-column>
      <el-table-column prop="riskType" label="风险类型" width="100" />
      <el-table-column prop="affectedComponent" label="影响组件" width="130" />
      <el-table-column prop="discoveredDate" label="发现日期" width="180">
        <template #default="{ row }">{{ formatTime(row.discoveredDate) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="riskStatusType(row.status)" size="small">{{ riskStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="dispatchDept" label="目标部门" width="120" />
      <el-table-column prop="dispatchTarget" label="目标人员" width="100" />
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status === 'pending'" type="primary" link size="small" @click="openDispatch(row)">分发</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dispatchVisible" title="分发风险工单" width="480px">
      <el-form :model="dispatchForm" label-width="100px">
        <el-form-item label="风险标题"><el-input :model-value="currentRisk?.title" disabled /></el-form-item>
        <el-form-item label="目标部门" required>
          <el-select v-model="dispatchForm.dept" placeholder="请选择" style="width: 100%">
            <el-option label="技术研发部" value="技术研发部" />
            <el-option label="数据智能部" value="数据智能部" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标人员" required>
          <el-select v-model="dispatchForm.target" placeholder="请选择" style="width: 100%">
            <el-option label="张三" value="张三" />
            <el-option label="李四" value="李四" />
          </el-select>
        </el-form-item>
        <el-form-item label="分发说明"><el-input v-model="dispatchForm.note" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dispatchVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmDispatch">确认分发</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { mockRiskRecords } from '@/mock/admin'
import { ElMessage } from 'element-plus'

const risks = ref(mockRiskRecords)
const dispatchVisible = ref(false)
const currentRisk = ref<any>(null)
const dispatchForm = ref({ dept: '', target: '', note: '' })

function riskStatusType(status: string) {
  const types: Record<string, string> = { pending: 'warning', dispatched: '', processing: 'info', resolved: 'success' }
  return types[status] || ''
}

function riskStatusLabel(status: string) {
  const labels: Record<string, string> = { pending: '待分发', dispatched: '已分发', processing: '处理中', resolved: '已解决' }
  return labels[status] || status
}

function openDispatch(row: any) {
  currentRisk.value = row
  dispatchForm.value = { dept: '', target: '', note: '' }
  dispatchVisible.value = true
}

function confirmDispatch() {
  if (currentRisk.value) {
    currentRisk.value.status = 'dispatched'
    currentRisk.value.dispatchDept = dispatchForm.value.dept
    currentRisk.value.dispatchTarget = dispatchForm.value.target
    currentRisk.value.dispatchTime = new Date().toLocaleDateString('zh-CN')
    currentRisk.value.dispatcher = '王五'
  }
  dispatchVisible.value = false
  ElMessage.success('分发成功')
}
</script>

<style scoped>
.desc { color: #909399; font-size: 13px; margin: 8px 0 0; }
</style>
