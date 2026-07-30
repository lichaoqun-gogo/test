<template>
  <el-dialog
    :model-value="visible"
    title="分配缺陷"
    width="520px"
    @update:model-value="$emit('update:visible', $event)"
    @close="resetForm"
  >
    <el-alert
      v-if="defect"
      :title="`缺陷 ${defect.defectNo} - ${defect.title}`"
      type="info"
      :closable="false"
      style="margin-bottom: 20px"
    />

    <el-form :model="form" label-width="80px" ref="formRef" :rules="rules">
      <el-form-item label="接收人" prop="assignee" required>
        <el-select v-model="form.assignee" placeholder="选择部门内运维人员" style="width: 100%">
          <el-option
            v-for="op in operators"
            :key="op.value"
            :label="op.label"
            :value="op.value"
          />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="$emit('update:visible', false)">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="submitting">确认分配</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { operators, type DefectRecord, type AssignMethod } from '@/mock/defects'

const props = defineProps<{
  visible: boolean
  defect: DefectRecord | null
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  submit: [data: { assignee: string; assignMethod: AssignMethod }]
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = ref({
  assignee: '',
})

const rules: FormRules = {
  assignee: [{ required: true, message: '请选择接收人', trigger: 'change' }],
}

function resetForm() {
  form.value = {
    assignee: '',
  }
  formRef.value?.clearValidate()
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    emit('submit', {
      assignee: form.value.assignee,
      assignMethod: 'operator',
    })
  } finally {
    submitting.value = false
  }
}
</script>
