<template>
  <div class="my-usage">
    <div class="page-header">
      <h2>我的使用</h2>
      <el-button type="primary" class="apply-btn" @click="showUsageForm">
        <el-icon><Plus /></el-icon>
        申请使用
      </el-button>
    </div>

    <el-table :data="usageList" border>
      <el-table-column prop="taskNo" label="任务编号" width="180">
        <template #default="{ row }">
          <el-link type="primary" @click="viewDetail(row.id)">{{ row.taskNo }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="componentName" label="组件名称" min-width="150" />
      <el-table-column prop="version" label="申请版本" width="150" />
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.status === 'pending'" type="warning" size="small">审批中</el-tag>
          <el-tag v-else-if="row.status === 'approved'" type="success" size="small">已通过</el-tag>
          <el-tag v-else-if="row.status === 'rejected'" type="danger" size="small">已驳回</el-tag>
          <el-tag v-else type="info" size="small">草稿</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="申请时间" width="200" />
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <template v-if="row.status === 'draft'">
            <el-button type="primary" link @click="handleSubmitFromList(row)">提交</el-button>
            <el-button type="primary" link @click="editDraft(row)">编辑</el-button>
          </template>
          <template v-else>
            <el-button type="primary" link @click="viewDetail(row.id)">详情</el-button>
          </template>
        </template>
      </el-table-column>
    </el-table>

    <!-- 申请使用表单对话框 -->
    <el-dialog
      v-model="showFormDialog"
      :title="editingTaskId ? '编辑使用申请' : '申请使用'"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="组件名称" prop="componentName">
          <el-input v-model="formData.componentName" placeholder="请输入组件名称" />
        </el-form-item>

        <el-form-item label="版本" prop="version">
          <el-input v-model="formData.version" placeholder="请输入申请版本" />
        </el-form-item>

        <el-form-item label="是否发布生产" prop="isProduction">
          <el-radio-group v-model="formData.isProduction">
            <el-radio :label="true">是</el-radio>
            <el-radio :label="false">否</el-radio>
          </el-radio-group>
          <el-alert
            v-if="formData.isProduction"
            type="warning"
            :closable="false"
            style="margin-top: 8px"
            show-icon
          >
            <template #default>
              <span style="font-size: 12px">还需提交引入申请</span>
            </template>
          </el-alert>
        </el-form-item>

        <el-form-item label="申请用途" prop="purpose">
          <el-input
            v-model="formData.purpose"
            type="textarea"
            :rows="4"
            placeholder="请说明申请用途"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showFormDialog = false">取消</el-button>
        <el-button @click="handleSaveDraft">保存草稿</el-button>
        <el-button type="primary" @click="handleSubmit">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useTaskStore } from '@/stores/task'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const taskStore = useTaskStore()

const formRef = ref<FormInstance>()
const showFormDialog = ref(false)
const editingTaskId = ref('')

// 表单数据
const formData = ref({
  componentName: '',
  version: '',
  isProduction: false,
  purpose: '',
})

// 表单校验规则
const formRules: FormRules = {
  componentName: [{ required: true, message: '请输入组件名称', trigger: 'blur' }],
  version: [{ required: true, message: '请输入申请版本', trigger: 'blur' }],
  purpose: [{ required: true, message: '请输入申请用途', trigger: 'blur' }],
}

// 使用taskStore中的任务数据，过滤出usage类型并映射字段
const usageList = computed(() => {
  const filtered = taskStore.tasks.filter(t => t.type === 'usage')
  
  return filtered.map(task => {
    const detail = task.task_detail || {}
    return {
      id: task.id,
      taskNo: task.taskNo,
      componentName: detail.componentName || '',
      version: detail.version || '',
      purpose: detail.purpose || '',
      status: task.status,
      createdAt: task.createTime || '',
    }
  })
})

// 页面加载时获取数据
onMounted(async () => {
  await taskStore.loadTasks()
  console.log('[MyUsage] onMounted, tasks count:', taskStore.tasks.length)
  if (taskStore.tasks.length > 0) {
    const firstTask = taskStore.tasks[0]
    console.log('[MyUsage] first task type:', firstTask.type)
    console.log('[MyUsage] first task task_detail:', firstTask.task_detail)
  }
})

// 显示申请表单
function showUsageForm() {
  editingTaskId.value = ''
  formData.value = {
    componentName: '',
    version: '',
    isProduction: false,
    purpose: '',
  }
  showFormDialog.value = true
}

// 编辑草稿
async function editDraft(row: any) {
  editingTaskId.value = row.id
  
  // 从后端获取完整详情
  try {
    const detail: any = await taskStore.getTaskDetail(row.id)
    const taskDetail = detail.task_detail || {}
    
    formData.value = {
      componentName: taskDetail.componentName || '',
      version: taskDetail.version || '',
      isProduction: taskDetail.isProduction === 1,
      purpose: taskDetail.purpose || '',
    }
    showFormDialog.value = true
  } catch (e: any) {
    ElMessage.error('获取详情失败: ' + (e.message || '未知错误'))
  }
}

// 查看详情
function viewDetail(id: string) {
  router.push(`/operator/my-task/${id}`)
}

// 保存草稿
async function handleSaveDraft() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    try {
      const requestData = {
        type: 'usage',
        title: `${formData.value.componentName}使用申请`,
        applicant: '张三',
        task_detail: {
          component_name: formData.value.componentName,
          version: formData.value.version,
          is_production: formData.value.isProduction ? 1 : 0,
          purpose: formData.value.purpose,
        },
        status: 'draft',
      }
      
      if (editingTaskId.value) {
        await taskStore.updateTask(editingTaskId.value, requestData)
        ElMessage.success('草稿已更新')
      } else {
        await taskStore.createTask(requestData)
        ElMessage.success('草稿已保存')
      }
      
      showFormDialog.value = false
    } catch (e: any) {
      ElMessage.error('保存失败: ' + (e.message || '未知错误'))
    }
  })
}

// 从列表直接提交
async function handleSubmitFromList(row: any) {
  try {
    await taskStore.updateTaskStatus(row.id, 'pending')
    ElMessage.success('提交成功')
    await taskStore.loadTasks()
  } catch (e: any) {
    ElMessage.error('提交失败: ' + (e.message || '未知错误'))
  }
}

// 提交表单
async function handleSubmit() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    try {
      const requestData = {
        type: 'usage',
        title: `${formData.value.componentName}使用申请`,
        applicant: '张三',
        task_detail: {
          component_name: formData.value.componentName,
          version: formData.value.version,
          is_production: formData.value.isProduction ? 1 : 0,
          purpose: formData.value.purpose,
        },
        status: 'pending',
      }
      
      if (editingTaskId.value) {
        await taskStore.updateTask(editingTaskId.value, requestData)
        ElMessage.success('使用申请已更新')
      } else {
        await taskStore.createTask(requestData)
        ElMessage.success('使用申请已提交')
      }
      
      showFormDialog.value = false
    } catch (e: any) {
      ElMessage.error('提交失败: ' + (e.message || '未知错误'))
    }
  })
}
</script>

<style scoped>
.my-usage {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 600;
  color: #1f2328;
}

.page-desc {
  margin: 0;
  font-size: 13px;
  color: #8c8c8c;
}
</style>
