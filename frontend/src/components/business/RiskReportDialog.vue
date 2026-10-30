<template>
  <el-dialog v-model="visible" title="上报风险（内部缺陷）" width="600px" :close-on-click-modal="false" @closed="resetForm">
    <el-form :model="form" label-width="100px">
      <el-form-item label="风险标题" required>
        <el-input v-model="form.title" placeholder="缺陷标题" />
      </el-form-item>
      <el-form-item label="来源" required>
        <el-select v-model="form.source" placeholder="选择来源" style="width: 100%">
          <el-option label="漏洞情报" value="vuln-intel" />
          <el-option label="SCA扫描" value="sca-scan" />
          <el-option label="人工上报" value="user-report" />
        </el-select>
      </el-form-item>
      <el-form-item label="风险等级">
        <el-select v-model="form.riskLevel" placeholder="选填" clearable style="width: 100%">
          <el-option v-for="(label, value) in RISK_LEVEL_OPTIONS" :key="value" :label="label" :value="value" />
        </el-select>
      </el-form-item>
      <el-form-item label="风险描述">
        <el-input v-model="form.description" type="textarea" :rows="4" placeholder="描述风险情况" />
      </el-form-item>
    </el-form>
    <el-alert type="info" :closable="false" show-icon title="上报后将自动生成内部缺陷编号（DF-xxx），进入「待上会」，由开源治理专员评审（忽略或上会分配）。" style="margin-top: 4px" />
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="submit">提交</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { riskApi } from '@/api'

const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
  (e: 'created', defectNo: string): void
}>()

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const RISK_LEVEL_OPTIONS: Record<string, string> = { critical: '严重', high: '高危', medium: '中危', low: '低危', 未知: '未知' }
const form = ref({ title: '', source: '', riskLevel: '', description: '' })
const submitting = ref(false)

function resetForm() {
  form.value = { title: '', source: '', riskLevel: '', description: '' }
}

async function submit() {
  if (!form.value.title || !form.value.source) {
    ElMessage.warning('请填写风险标题和来源')
    return
  }
  submitting.value = true
  try {
    const res: any = await riskApi.createDefect({
      title: form.value.title,
      source: form.value.source,
      riskLevel: form.value.riskLevel || '未知',
      description: form.value.description,
    })
    ElMessage.success('上报成功，已进入待上会')
    visible.value = false
    emit('created', res?.defectNo || '')
  } catch (e: any) {
    ElMessage.error('上报失败: ' + (e.message || ''))
  } finally {
    submitting.value = false
  }
}
</script>
