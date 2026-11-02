<template>
  <div class="disposal-assign-picker">
    <!-- 自动/手动两通道（showAuto=false 时仅手动指定，如内部缺陷审核分配） -->
    <el-radio-group v-if="showAuto" v-model="mode" style="margin-bottom: 8px">
      <el-radio value="auto">自动分配（系统按涉及项目/软件自动匹配处理人）</el-radio>
      <el-radio value="manual">手动指定</el-radio>
    </el-radio-group>
    <div v-if="mode === 'manual'">
      <el-input v-model="manualInput" placeholder="输入负责人账号，多个用逗号/顿号分隔，如 wangwu、lisi" clearable />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  showAuto?: boolean
  modelValue?: { mode: string; assignees: string[] }
}>()
const emit = defineEmits(['update:modelValue'])

// 无自动入口时固定手动模式；有自动入口时默认自动
const mode = ref(props.showAuto ? props.modelValue?.mode || 'auto' : 'manual')
const manualInput = ref(props.modelValue?.assignees?.join('、') || '')

watch(mode, emitValue)
watch(manualInput, emitValue)

function parseAssignees(input: string): string[] {
  return input
    .split(/[,，、;；\s]+/)
    .map((s) => s.trim())
    .filter(Boolean)
}

function emitValue() {
  emit('update:modelValue', {
    mode: mode.value,
    assignees: mode.value === 'manual' ? parseAssignees(manualInput.value) : [],
  })
}
</script>
