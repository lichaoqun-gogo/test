<template>
  <el-dialog v-model="visible" title="分发风险工单" width="520px" @close="handleClose">
    <el-form :model="form" label-width="110px">
      <el-form-item label="风险来源">
        <el-tag :type="sourceType(currentDefect?.source)">{{ currentDefect?.source || '-' }}</el-tag>
      </el-form-item>
      <el-form-item label="风险标题">
        <span>{{ currentDefect?.title }}</span>
      </el-form-item>
      <el-form-item label="风险等级">
        <RiskTag :level="(currentDefect?.riskLevel || 'low') as any" />
      </el-form-item>

      <el-form-item label="分发方式" required>
        <el-radio-group v-model="form.distributeType">
          <el-radio value="assignee">分配开源治理专员</el-radio>
          <el-radio value="jira">生成 Jira 工单</el-radio>
        </el-radio-group>
      </el-form-item>

      <template v-if="form.distributeType === 'assignee'">
        <el-form-item label="接收部门" required>
          <el-select v-model="form.dept" placeholder="请选择部门" style="width: 100%" @change="handleDeptChange">
            <el-option label="技术研发部" value="技术研发部" />
            <el-option label="数据智能部" value="数据智能部" />
            <el-option label="产品技术部" value="产品技术部" />
            <el-option label="基础架构部" value="基础架构部" />
          </el-select>
        </el-form-item>
        <el-form-item label="接收专员" required>
          <el-select v-model="form.target" placeholder="请选择专员" style="width: 100%">
            <el-option v-for="s in filteredSpecialists" :key="s" :label="s" :value="s" />
          </el-select>
        </el-form-item>
      </template>

      <template v-if="form.distributeType === 'jira'">
        <el-form-item label="接收部门" required>
          <el-select v-model="form.dept" placeholder="请选择部门" style="width: 100%">
            <el-option label="技术研发部" value="技术研发部" />
            <el-option label="数据智能部" value="数据智能部" />
            <el-option label="产品技术部" value="产品技术部" />
            <el-option label="基础架构部" value="基础架构部" />
          </el-select>
        </el-form-item>
        <el-form-item label="Jira 项目">
          <el-input v-model="form.jiraProject" placeholder="如 OSM-SEC" />
        </el-form-item>
        <el-form-item label="Jira 工单类型">
          <el-select v-model="form.jiraType" style="width: 100%">
            <el-option label="Bug" value="Bug" />
            <el-option label="Task" value="Task" />
            <el-option label="Story" value="Story" />
          </el-select>
        </el-form-item>
      </template>

      <el-form-item label="分发说明">
        <el-input v-model="form.note" type="textarea" :rows="3" placeholder="请输入分发说明和注意事项" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="handleConfirm">确认分发</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'

interface DefectItem {
  id: string
  source: string
  title: string
  riskLevel: string
}

interface DistributeForm {
  distributeType: 'assignee' | 'jira'
  dept: string
  target: string
  jiraProject: string
  jiraType: string
  note: string
}

const visible = ref(false)
const currentDefect = ref<DefectItem | null>(null)
const form = ref<DistributeForm>({
  distributeType: 'assignee',
  dept: '',
  target: '',
  jiraProject: '',
  jiraType: 'Bug',
  note: '',
})

const emit = defineEmits<{
  confirm: [data: { defect: DefectItem; form: DistributeForm }]
}>()

// 部门-专员映射（模拟）
const specialistMap: Record<string, string[]> = {
  '技术研发部': ['张三', '李四', '王五'],
  '数据智能部': ['赵六', '孙七'],
  '产品技术部': ['周八', '吴九'],
  '基础架构部': ['郑十', '冯十一'],
}

const filteredSpecialists = computed(() => {
  return form.value.dept ? (specialistMap[form.value.dept] || []) : []
})

function sourceType(source?: string): 'danger' | 'warning' | 'info' {
  const map: Record<string, 'danger' | 'warning' | 'info'> = { '外部推送': 'danger', 'SCA扫描': 'warning', '人工上报': 'info' }
  return map[source || ''] || 'info'
}

function handleDeptChange() {
  form.value.target = ''
}

function open(defect: DefectItem) {
  currentDefect.value = defect
  visible.value = true
}

function handleClose() {
  form.value = { distributeType: 'assignee', dept: '', target: '', jiraProject: '', jiraType: 'Bug', note: '' }
  currentDefect.value = null
}

function handleConfirm() {
  if (!form.value.dept) {
    ElMessage.warning('请选择接收部门')
    return
  }
  if (form.value.distributeType === 'assignee' && !form.value.target) {
    ElMessage.warning('请选择接收专员')
    return
  }
  if (!currentDefect.value) return

  emit('confirm', { defect: currentDefect.value, form: { ...form.value } })
  visible.value = false
}

defineExpose({ open })
</script>
