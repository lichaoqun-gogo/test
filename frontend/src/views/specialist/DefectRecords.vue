<template>
  <div class="defect-records">
    <h3>缺陷记录</h3>

    <!-- 统计卡片 -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="6">
        <el-statistic title="缺陷总数" :value="defects.length" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="待研判" :value="defects.filter(d => d.status === 'pending').length" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="已分配" :value="defects.filter(d => d.status === 'assigned').length" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="已解决" :value="defects.filter(d => d.status === 'resolved').length" />
      </el-col>
    </el-row>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-form :inline="true" :model="filters">
        <el-form-item label="全局状态">
          <el-select v-model="filters.status" placeholder="全部" clearable style="width: 140px">
            <el-option label="待研判" value="待研判" />
            <el-option label="待评审" value="待评审" />
            <el-option label="已分配" value="已分配" />
            <el-option label="已忽略" value="已忽略" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleFilter">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
          <el-button type="primary" plain @click="openBatchDisposeDialog" :disabled="selectedDefects.length === 0">
            <el-icon><EditPen /></el-icon>
            批量处置 ({{ selectedDefects.length }})
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- Tabs: 缺陷列表 / 分配记录 -->
    <el-tabs v-model="activeTab">
      <!-- 缺陷列表 -->
      <el-tab-pane label="缺陷列表" name="defects">
        <el-table :data="filteredDefects" stripe @selection-change="handleSelectionChange">
          <el-table-column type="selection" width="55" />
          <el-table-column prop="defectNo" label="缺陷编号" width="160" />
          <el-table-column prop="cveId" label="CVE编号" width="160" />
          <el-table-column prop="title" label="缺陷标题" min-width="180" />
          <el-table-column label="严重等级" width="100">
            <template #default="{ row }">
              <RiskTag :level="row.severityLevel" />
            </template>
          </el-table-column>
          <el-table-column prop="component" label="关联组件" width="160" />
          <el-table-column prop="assignee" label="处置人" width="120">
            <template #default="{ row }">
              {{ row.assignee || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="discoveredAt" label="发现时间" width="180">
            <template #default="{ row }">{{ formatTime(row.discoveredAt) }}</template>
          </el-table-column>
          <el-table-column label="全局状态" width="100">
            <template #default="{ row }">
              <el-tag :type="statusTagType(row.status)" size="small">
                {{ statusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <el-button
                v-if="row.status === 'pending'"
                type="primary"
                size="small"
                @click="openAssignDialog(row as DefectRecord)"
              >
                分配
              </el-button>
              <span v-else class="text-muted">-</span>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- 分配记录 -->
      <el-tab-pane label="分配记录" name="records">
        <el-table :data="assignRecords" stripe>
          <el-table-column prop="taskNo" label="任务编号" width="160" />
          <el-table-column prop="defectNo" label="缺陷编号" width="160" />
          <el-table-column prop="defectTitle" label="缺陷标题" min-width="180" />
          <el-table-column prop="assignee" label="接收人" width="100" />
          <el-table-column prop="assignedBy" label="分配人" width="100" />
          <el-table-column prop="assignedAt" label="分配时间" width="180">
            <template #default="{ row }">{{ formatTime(row.assignedAt) }}</template>
          </el-table-column>
          <el-table-column label="任务状态" width="100">
            <template #default="{ row }">
              <el-tag :type="recordStatusTagType(row.status)" size="small">
                {{ statusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <!-- 分配弹窗 -->
    <AssignDialog
      v-model:visible="assignDialogVisible"
      :defect="currentDefect"
      @submit="handleAssign"
    />

    <!-- 批量处置对话框 -->
    <el-dialog v-model="batchDisposeDialogVisible" title="批量处置" width="900px">
      <el-alert :title="`已选择 ${selectedDefects.length} 个缺陷`" type="info" :closable="false" show-icon style="margin-bottom: 16px" />
      
      <el-form label-width="120px" style="margin-bottom: 16px">
        <el-form-item label="目标状态" required>
          <el-select v-model="batchDisposeForm.targetStatus" placeholder="请选择目标状态" style="width: 100%" @change="onBatchStatusChange">
            <el-option label="待研判" value="待研判" />
            <el-option label="待评审" value="待评审" />
            <el-option label="已分配" value="已分配" />
            <el-option label="已忽略" value="已忽略" />
          </el-select>
        </el-form-item>
      </el-form>

      <!-- 已分配状态：显示每个缺陷的卡片 -->
      <template v-if="batchDisposeForm.targetStatus === '已分配'">
        <!-- 调试信息 -->
        <el-alert 
          :title="`[调试] targetStatus=${batchDisposeForm.targetStatus}, selectedDefects=${selectedDefects.length}, forms=${Object.keys(batchDisposeForms).length}`" 
          type="warning" 
          :closable="false" 
          style="margin-bottom: 16px"
        />
        
        <el-card 
          v-for="(defect, index) in selectedDefects" 
          :key="defect.id"
          shadow="hover"
          style="margin-bottom: 16px"
        >
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-tag type="primary" size="small">已分配</el-tag>
                {{ defect.title }}
                <span style="color: #909399; font-size: 12px; margin-left: 8px;">[defectType: {{ defect.defectType }}]</span>
              </span>
              <el-tag type="info" size="small">缺陷 {{ index + 1 }}/{{ selectedDefects.length }}</el-tag>
            </div>
          </template>

          <!-- 缺陷基本信息 -->
          <el-descriptions :column="4" size="small" border>
            <el-descriptions-item label="缺陷编号">{{ defect.defectNo }}</el-descriptions-item>
            <el-descriptions-item label="关联组件">{{ defect.component }}</el-descriptions-item>
            <el-descriptions-item label="严重等级">
              <RiskTag :level="defect.severityLevel" />
            </el-descriptions-item>
            <el-descriptions-item label="缺陷类型">
              <el-tag type="warning" size="small">{{ defect.defectType === 'dependency' ? '依赖包' : '开源基础软件' }}</el-tag>
            </el-descriptions-item>
          </el-descriptions>

          <!-- 依赖包：关联自研项目和开源基础软件列表 -->
          <template v-if="defect.defectType === 'dependency'">
            <!-- 调试信息 -->
            <el-alert 
              :title="`[调试] 依赖包类型 - projectList: ${batchDisposeForms[defect.id]?.projectList?.length || 0}, softwareList: ${batchDisposeForms[defect.id]?.softwareList?.length || 0}`" 
              type="info" 
              :closable="false" 
              style="margin-bottom: 12px"
            />
            
            <div class="section-title">关联自研项目</div>
            <el-table :data="batchDisposeForms[defect.id].projectList" stripe size="small" style="margin-bottom: 12px">
              <el-table-column prop="name" label="项目名称" width="180" />
              <el-table-column prop="dept" label="所属部门" width="130" />
              <el-table-column prop="owners" label="处置人" min-width="250">
                <template #default="{ row }">
                  <el-select v-model="row.owners" placeholder="请选择处置人" size="small" multiple collapse-tags collapse-tags-tooltip :max-collapse-tags="3">
                    <el-option label="张三" value="张三" />
                    <el-option label="李四" value="李四" />
                    <el-option label="王五" value="王五" />
                    <el-option label="赵六" value="赵六" />
                    <el-option label="钱七" value="钱七" />
                  </el-select>
                </template>
              </el-table-column>
            </el-table>

            <div class="section-title">关联开源软件</div>
            <el-table :data="batchDisposeForms[defect.id].softwareList" stripe size="small">
              <el-table-column prop="name" label="软件名称" width="180" />
              <el-table-column prop="version" label="当前版本" width="100" />
              <el-table-column prop="owners" label="处置人" min-width="250">
                <template #default="{ row }">
                  <el-select v-model="row.owners" placeholder="请选择处置人" size="small" multiple collapse-tags collapse-tags-tooltip :max-collapse-tags="3">
                    <el-option label="张三" value="张三" />
                    <el-option label="李四" value="李四" />
                    <el-option label="王五" value="王五" />
                    <el-option label="赵六" value="赵六" />
                    <el-option label="钱七" value="钱七" />
                  </el-select>
                </template>
              </el-table-column>
            </el-table>
          </template>

          <!-- 开源基础软件：处置人列表 -->
          <template v-else-if="defect.defectType === 'software'">
            <div class="section-title">处置人列表</div>
            <el-table :data="batchDisposeForms[defect.id].assigneeList" stripe size="small">
              <el-table-column prop="name" label="姓名" width="150" />
              <el-table-column prop="dept" label="部门" width="150" />
              <el-table-column prop="role" label="角色" width="150" />
              <el-table-column label="操作" width="100">
                <template #default="{ $index }">
                  <el-button type="danger" link size="small" @click="removeBatchAssignee(defect.id, $index)">
                    移除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-button type="primary" link size="small" style="margin-top: 8px" @click="showAddBatchAssigneeDialog(defect.id)">
              <el-icon><Plus /></el-icon>
              添加处置人
            </el-button>
          </template>
        </el-card>
      </template>

      <template #footer>
        <el-button @click="batchDisposeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleBatchDispose" :disabled="!batchDisposeForm.targetStatus">确认处置</el-button>
      </template>
    </el-dialog>

    <!-- 添加处置人对话框 -->
    <el-dialog v-model="addBatchAssigneeDialogVisible" title="添加处置人" width="500px" append-to-body>
      <el-form>
        <el-form-item label="选择处置人">
          <el-select v-model="newBatchAssignee" placeholder="请选择处置人" style="width: 100%">
            <el-option label="张三 - 技术研发部" value="张三" />
            <el-option label="李四 - 平台研发部" value="李四" />
            <el-option label="王五 - 运维部" value="王五" />
            <el-option label="赵六 - 安全部" value="赵六" />
            <el-option label="钱七 - 质量管理部" value="钱七" />
            <el-option label="孙八 - 数据团队" value="孙八" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addBatchAssigneeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAddBatchAssignee" :disabled="!newBatchAssignee">确认添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { EditPen, Plus } from '@element-plus/icons-vue'
import AssignDialog from '@/components/specialist/AssignDialog.vue'

// 类型定义（清理 mock 后本地保留）
type DefectRecord = any
type DefectStatus = any
type DefectSource = any
type AssignMethod = any

// 状态管理
const activeTab = ref('defects')
const assignDialogVisible = ref(false)
const batchDisposeDialogVisible = ref(false)
const currentDefect = ref<DefectRecord | null>(null)
const selectedDefects = ref<DefectRecord[]>([])

// 批量处置表单
const batchDisposeForm = ref({
  targetStatus: '',
})

// 每个缺陷的处置表单数据
const batchDisposeForms = ref<Record<string, {
  projectList: Array<{ name: string; dept: string; owners: string[] }>
  softwareList: Array<{ name: string; version: string; owners: string[] }>
  assigneeList: Array<{ name: string; dept: string; role: string }>
}>>({})

// 添加处置人对话框
const addBatchAssigneeDialogVisible = ref(false)
const newBatchAssignee = ref('')
const currentDefectId = ref('')

// 筛选条件
const filters = ref({
  status: '',
})

// 缺陷数据
const defects = ref<any[]>([])
const assignRecords = ref<any[]>([])

// 筛选后的缺陷列表
const filteredDefects = computed(() => {
  let list = defects.value
  if (filters.value.status) {
    // 中文状态映射到英文
    const statusMap: Record<string, string> = {
      '待研判': 'pending',
      '待评审': 'pending_review',
      '已分配': 'assigned',
      '已忽略': 'ignored',
    }
    const englishStatus = statusMap[filters.value.status] || filters.value.status
    list = list.filter(d => d.status === englishStatus)
  }
  return list
})

// 处理选择变化
function handleSelectionChange(selection: DefectRecord[]) {
  console.log('[缺陷管理] handleSelectionChange 被调用', selection)
  console.log('[缺陷管理] 选中的缺陷数量:', selection.length)
  console.log('[缺陷管理] 选中的缺陷详情:', selection.map(d => ({
    id: d.id,
    defectNo: d.defectNo,
    title: d.title,
    defectType: d.defectType,
  })))
  selectedDefects.value = selection
}

// 打开批量处置对话框
function openBatchDisposeDialog() {
  console.log('[缺陷管理] openBatchDisposeDialog 被调用')
  console.log('[缺陷管理] selectedDefects:', selectedDefects.value)
  
  if (selectedDefects.value.length === 0) {
    ElMessage.warning('请先选择要处置的缺陷')
    return
  }
  batchDisposeForm.value.targetStatus = ''
  batchDisposeForms.value = {}
  
  // 为每个选中的缺陷初始化表单
  selectedDefects.value.forEach(defect => {
    console.log('[缺陷管理] 初始化缺陷表单数据', {
      defectId: defect.id,
      defectType: defect.defectType,
      hasProjectList: defect.defectType === 'dependency',
      hasAssigneeList: defect.defectType === 'software',
    })
    
    batchDisposeForms.value[defect.id] = {
      projectList: [],
      softwareList: [],
      assigneeList: [],
    }
  })
  
  console.log('[缺陷管理] batchDisposeForms 初始化完成:', batchDisposeForms.value)
  batchDisposeDialogVisible.value = true
}

function onBatchStatusChange() {
  console.log('[缺陷管理] onBatchStatusChange 被调用, targetStatus:', batchDisposeForm.value.targetStatus)
  
  // 切换状态时清空所有缺陷的表单数据
  if (batchDisposeForm.value.targetStatus !== '已分配') {
    console.log('[缺陷管理] 状态不是已分配，清空表单数据')
    Object.keys(batchDisposeForms.value).forEach(defectId => {
      const defect = selectedDefects.value.find(d => d.id === defectId)
      if (defect) {
        batchDisposeForms.value[defectId] = {
          projectList: [],
          softwareList: [],
          assigneeList: [],
        }
      }
    })
  } else {
    console.log('[缺陷管理] 状态是已分配，应该显示缺陷卡片')
  }
}

function removeBatchAssignee(defectId: string, index: number) {
  batchDisposeForms.value[defectId].assigneeList.splice(index, 1)
}

function showAddBatchAssigneeDialog(defectId: string) {
  currentDefectId.value = defectId
  newBatchAssignee.value = ''
  addBatchAssigneeDialogVisible.value = true
}

function confirmAddBatchAssignee() {
  if (newBatchAssignee.value) {
    const deptMap: Record<string, string> = {
      '张三': '技术研发部',
      '李四': '平台研发部',
      '王五': '运维部',
      '赵六': '安全部',
      '钱七': '质量管理部',
      '孙八': '数据团队',
    }
    batchDisposeForms.value[currentDefectId.value].assigneeList.push({
      name: newBatchAssignee.value,
      dept: deptMap[newBatchAssignee.value] || '',
      role: '处置人',
    })
  }
  addBatchAssigneeDialogVisible.value = false
}

// 批量处置
function handleBatchDispose() {
  if (!batchDisposeForm.value.targetStatus) {
    ElMessage.warning('请选择目标状态')
    return
  }

  const targetStatus = batchDisposeForm.value.targetStatus

  // 更新选中缺陷的状态
  selectedDefects.value.forEach(defect => {
    const idx = defects.value.findIndex(d => d.id === defect.id)
    if (idx !== -1) {
      defects.value[idx].status = targetStatus as DefectStatus
      
      // 如果状态为已分配，创建任务
      if (targetStatus === '已分配' && batchDisposeForms.value[defect.id]) {
        const form = batchDisposeForms.value[defect.id]
        // 收集所有处置人
        const assignees: string[] = []
        
        if (defect.defectType === 'dependency') {
          // 依赖包：从项目和软件列表中收集处置人
          form.projectList.forEach(p => { 
            if (p.owners && p.owners.length > 0) {
              assignees.push(...p.owners)
            }
          })
          form.softwareList.forEach(s => { 
            if (s.owners && s.owners.length > 0) {
              assignees.push(...s.owners)
            }
          })
        } else if (defect.defectType === 'software') {
          // 软件：从处置人列表中收集
          form.assigneeList.forEach(a => assignees.push(a.name))
        }
        
        const uniqueAssignees = [...new Set(assignees)]
        
        if (uniqueAssignees.length > 0) {
          // 生成任务记录
          uniqueAssignees.forEach((assignee, i) => {
            const taskNo = `TSK-${new Date().getFullYear()}-${String(assignRecords.value.length + 1 + i).padStart(4, '0')}`
            // 任务标题格式：名称-版本-缺陷编号
            const taskTitle = `${defect.component}-${defect.componentVersion}-${defect.defectNo}`
            
            const newRecord = {
              id: `ar${Date.now()}-${i}`,
              taskNo,
              defectNo: defect.defectNo,
              defectTitle: taskTitle,
              assignee,
              assignedBy: '李四',
              assignedAt: new Date().toLocaleString('zh-CN'),
              status: 'assigned' as DefectStatus,
            }
            assignRecords.value.unshift(newRecord)
          })
        }
      }
    }
  })

  ElMessage.success(`成功处置 ${selectedDefects.value.length} 个缺陷`)
  batchDisposeDialogVisible.value = false
  selectedDefects.value = []
}

// 工具函数
function sourceLabel(source: string) {
  return source
}

function sourceTagType(source: string): 'danger' | 'info' | undefined {
  return source === 'external' ? 'danger' : 'info'
}

function statusLabel(status: string) {
  const labels: Record<string, string> = {
    pending: '待研判',
    pending_review: '待评审',
    assigned: '已分配',
    ignored: '已忽略',
    resolved: '已解决',
  }
  return labels[status] || status
}

function statusTagType(status: string): 'warning' | 'info' | 'success' | 'danger' | undefined {
  const types: Record<string, 'warning' | 'info' | 'success' | 'danger' | undefined> = {
    pending: 'warning',
    pending_review: 'info',
    assigned: 'primary',
    ignored: 'danger',
    resolved: 'success',
  }
  return types[status]
}

function recordStatusTagType(status: string): 'warning' | 'info' | 'success' | undefined {
  return statusTagType(status)
}

function handleFilter() {
  // computed 自动处理筛选
}

function resetFilter() {
  filters.value = { status: '' }
}

function openAssignDialog(row: DefectRecord) {
  currentDefect.value = row
  assignDialogVisible.value = true
}

function handleAssign(data: { assignee: string; assignMethod: AssignMethod }) {
  if (!currentDefect.value) return

  // 更新缺陷状态
  currentDefect.value.status = 'assigned'

  // 生成任务编号
  const taskNo = `TSK-${new Date().getFullYear()}-${String(assignRecords.value.length + 1).padStart(4, '0')}`

  // 任务标题格式：名称-版本-缺陷编号
  const taskTitle = `${currentDefect.value.component}-${currentDefect.value.componentVersion}-${currentDefect.value.defectNo}`

  // 添加分配记录
  const newRecord = {
    id: `ar${Date.now()}`,
    taskNo,
    defectNo: currentDefect.value.defectNo,
    defectTitle: taskTitle,
    assignee: data.assignee,
    assignedBy: '李四',
    assignedAt: new Date().toLocaleString('zh-CN'),
    status: 'assigned' as DefectStatus,
  }
  assignRecords.value.unshift(newRecord)

  // 关闭弹窗
  assignDialogVisible.value = false
  currentDefect.value = null

  ElMessage.success('缺陷分配成功')

  // 模拟通知接收人
  ElMessage.info({
    message: `已通知接收人 ${data.assignee}，任务编号：${taskNo}，请等待处理`,
    duration: 3000,
  })
}

// 加载缺陷列表
async function loadDefects() {
  // TODO: 对接后端 API 获取缺陷列表
}

// 加载分配记录
async function loadAssignRecords() {
  // TODO: 对接后端 API 获取分配记录
}

onMounted(() => {
  loadDefects()
  loadAssignRecords()
})
</script>

<style scoped>
.defect-records {
  padding: 4px;
}

.desc {
  color: #909399;
  font-size: 13px;
  margin: 8px 0 0;
}

.filter-bar {
  margin-top: 16px;
  padding: 16px;
  background: #fff;
  border-radius: 4px;
}

.text-muted {
  color: #c0c4cc;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
  margin-top: 12px;
  padding-left: 8px;
  border-left: 3px solid #409EFF;
}
</style>
