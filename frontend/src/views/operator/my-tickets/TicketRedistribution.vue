<template>
  <div class="ticket-redistribution">
    <el-alert
      title="风险工单二次分配"
      description="此处展示管理部门分发到本部门的待处理风险工单，可将工单二次分配给部门内的具体运维人员"
      type="warning"
      :closable="false"
      show-icon
      style="margin-bottom: 16px"
    />

    <el-tabs v-model="activeTab">
      <el-tab-pane label="待分配工单" name="pending">
        <el-table :data="pendingTickets" stripe>
          <el-table-column prop="ticketNo" label="工单编号" width="150" />
          <el-table-column prop="ticketTitle" label="工单标题" min-width="200" />
          <el-table-column label="风险等级" width="100">
            <template #default="{ row }">
              <RiskTag :level="row.riskLevel" />
            </template>
          </el-table-column>
          <el-table-column prop="fromDept" label="来源部门" width="120" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <TicketStatusTag :status="row.status" />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button
                v-if="row.status === 'pending'"
                type="primary"
                size="small"
                @click="openAssignDialog(row)"
              >
                分配
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="分配记录" name="records">
        <el-table :data="assignedRecords" stripe>
          <el-table-column prop="ticketNo" label="工单编号" width="150" />
          <el-table-column prop="ticketTitle" label="工单标题" min-width="160" />
          <el-table-column label="风险等级" width="100">
            <template #default="{ row }">
              <RiskTag :level="row.riskLevel" />
            </template>
          </el-table-column>
          <el-table-column prop="assignedTo" label="接收人" width="100" />
          <el-table-column prop="assignedBy" label="分配人" width="100" />
          <el-table-column prop="assignedAt" label="分配时间" width="160" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <TicketStatusTag :status="row.status" />
            </template>
          </el-table-column>
          <el-table-column prop="note" label="备注" min-width="160" />
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <!-- Assign Dialog -->
    <el-dialog v-model="assignDialogVisible" title="分配工单" width="480px">
      <el-form :model="assignForm" label-width="100px">
        <el-form-item label="工单">
          <el-input :model-value="currentTicket?.ticketTitle" disabled />
        </el-form-item>
        <el-form-item label="接收人" required>
          <el-select v-model="assignForm.assignee" placeholder="选择部门内运维人员" style="width: 100%">
            <el-option label="张三" value="张三" />
            <el-option label="赵六" value="赵六" />
            <el-option label="钱七" value="钱七" />
          </el-select>
        </el-form-item>
        <el-form-item label="分配说明">
          <el-input v-model="assignForm.note" type="textarea" :rows="3" placeholder="请输入分配说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="assignDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAssign">确认分配</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { mockRedistributionRecords } from '@/mock/tickets'
import { ElMessage } from 'element-plus'

const activeTab = ref('pending')
const assignDialogVisible = ref(false)
const currentTicket = ref<any>(null)

const pendingTickets = computed(() =>
  mockRedistributionRecords.filter(r => r.status === 'pending')
)

const assignedRecords = computed(() =>
  mockRedistributionRecords.filter(r => r.status !== 'pending')
)

const assignForm = ref({
  assignee: '',
  note: '',
})

function openAssignDialog(row: any) {
  currentTicket.value = row
  assignForm.value = { assignee: '', note: '' }
  assignDialogVisible.value = true
}

function confirmAssign() {
  if (!assignForm.value.assignee) {
    ElMessage.warning('请选择接收人')
    return
  }
  if (currentTicket.value) {
    currentTicket.value.assignedTo = assignForm.value.assignee
    currentTicket.value.assignedBy = '李四'
    currentTicket.value.assignedAt = new Date().toLocaleString('zh-CN')
    currentTicket.value.status = 'assigned'
    currentTicket.value.note = assignForm.value.note
  }
  assignDialogVisible.value = false
  ElMessage.success('工单分配成功')
}
</script>
