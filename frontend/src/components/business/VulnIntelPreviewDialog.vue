<template>
  <el-dialog
    :model-value="modelValue"
    title="自动分配确认（可编辑负责人）"
    width="760px"
    :close-on-click-modal="false"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <el-alert
      type="info" :closable="false" show-icon style="margin-bottom: 12px"
      title="系统已按组件类型自动匹配处置单与负责人（基础软件=1 张运维单；组件=每个依赖项目 1 张使用单）。可编辑负责人后确认创建，以编辑结果为准。"
    />
    <el-table :data="rows" border size="small">
      <el-table-column label="类型" width="100">
        <template #default="{ row }">{{ row.projectType === 'software' ? '基础软件' : '组件项目' }}</template>
      </el-table-column>
      <el-table-column label="项目/软件" prop="projectName" min-width="140" />
      <el-table-column label="处置部门" prop="deptName" width="120" />
      <el-table-column label="团队" prop="team" width="110" />
      <el-table-column label="负责人" min-width="200">
        <template #default="{ row }">
          <el-input v-model="row.ownerUsername" placeholder="负责人账号，可编辑（需为系统用户）" clearable />
        </template>
      </el-table-column>
    </el-table>
    <div v-if="rows.length === 0" style="color: #909399; text-align: center; padding: 12px 0">未匹配到待裂变处置单</div>
    <template #footer>
      <el-button @click="emit('update:modelValue', false)">取消</el-button>
      <el-button type="primary" @click="confirm">确认创建</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

// 预览单：{ projectId, projectType, projectName, deptName, team, candidates:[{username,name,dept,team}], matchedOwner }
const props = defineProps<{
  modelValue: boolean
  previewList: any[]
}>()
const emit = defineEmits(['update:modelValue', 'confirm'])

const rows = ref<any[]>([])

watch(
  () => props.modelValue,
  v => {
    if (v) initRows()
  }
)

function initRows() {
  rows.value = (props.previewList || []).map(p => ({
    projectId: p.projectId,
    projectType: p.projectType,
    projectName: p.projectName,
    deptName: p.deptName,
    team: p.team || '',
    candidates: p.candidates || [],
    ownerUsername: p.matchedOwner?.username || p.candidates?.[0]?.username || '',
  }))
}

// 确认：以编辑后的负责人为准回传覆盖项（供后端按 projectId+projectType 覆盖各处置单负责人）
function confirm() {
  const overrides = rows.value
    .filter(r => r.ownerUsername)
    .map(r => ({ projectId: r.projectId, projectType: r.projectType, assigneeUsername: r.ownerUsername }))
  emit('confirm', overrides)
  emit('update:modelValue', false)
}
</script>
