<template>
  <el-dialog
    v-model="dialogVisible"
    title="发起新申请"
    width="700px"
    :close-on-click-modal="false"
  >
    <div class="new-task-dialog">
      <p class="dialog-desc">请选择要发起的申请类型</p>

      <el-row :gutter="16" class="type-grid">
        <el-col :span="12" v-for="type in taskTypes" :key="type.value">
          <el-card
            shadow="hover"
            class="type-card"
            :class="{ selected: selectedType === type.value }"
            @click="handleSelectType(type.value)"
          >
            <div class="type-card-content">
              <div class="type-icon" :style="{ backgroundColor: type.bgColor }">
                <el-icon :size="28" color="#fff">
                  <component :is="type.icon" />
                </el-icon>
              </div>
              <h4 class="type-name">{{ type.label }}</h4>
              <p class="type-desc">{{ type.desc }}</p>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { TaskType } from '@/types'
import {
  Document,
  Upload,
  WarningFilled,
  Bell,
} from '@element-plus/icons-vue'

const props = defineProps<{
  modelValue: boolean
  defaultType?: TaskType
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'create': [type: TaskType]
}>()

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const selectedType = ref<TaskType | ''>(props.defaultType || '')

// 监听defaultType变化
watch(() => props.defaultType, (val) => {
  selectedType.value = val || ''
}, { immediate: true })

const taskTypes = [
  {
    value: 'introduce' as TaskType,
    label: '引入申请',
    desc: '申请引入新的开源组件到项目中',
    icon: Document,
    bgColor: '#409EFF',
  },
  {
    value: 'upgrade' as TaskType,
    label: '升级申请',
    desc: '申请升级已有开源组件到新版本',
    icon: Upload,
    bgColor: '#67C23A',
  },
  {
    value: 'withdraw' as TaskType,
    label: '退出申请',
    desc: '申请退出不再使用的开源组件',
    icon: WarningFilled,
    bgColor: '#E6A23C',
  },
  {
    value: 'defect' as TaskType,
    label: '缺陷上报',
    desc: '上报发现的开源组件安全缺陷',
    icon: Bell,
    bgColor: '#F56C6C',
  },
]

function handleSelectType(type: TaskType) {
  selectedType.value = type
  emit('create', type)
  selectedType.value = ''
  dialogVisible.value = false
}
</script>

<style scoped>
.new-task-dialog {
  padding: 8px 0;
}

.dialog-desc {
  color: #606266;
  font-size: 14px;
  margin-bottom: 20px;
}

.type-grid {
  margin-top: 0;
}

.type-card {
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
  margin-bottom: 16px;
}

.type-card:hover,
.type-card.selected {
  border-color: #409EFF;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.2);
}

.type-card-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 0;
}

.type-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
}

.type-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px;
}

.type-desc {
  font-size: 13px;
  color: #909399;
  margin: 0;
  text-align: center;
  line-height: 1.4;
}
</style>
