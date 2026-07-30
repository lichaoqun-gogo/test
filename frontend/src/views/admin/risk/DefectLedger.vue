<template>
  <div class="defect-ledger">
    <div class="page-header">
      <div class="header-left">
        <h3>缺陷台账</h3>
        <p class="desc">集中展示所有外部推送和内部 SCA 发现的风险，支持状态管理和批量操作</p>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="16" style="margin-top: 16px; margin-bottom: 20px">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #fff3e0; color: #ff9800">
              <el-icon :size="28"><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value" style="color: #ff9800">{{ pendingReviewCount }}</div>
              <div class="stat-label">待研判</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #e3f2fd; color: #2196f3">
              <el-icon :size="28"><Tickets /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value" style="color: #2196f3">{{ pendingMeetingCount }}</div>
              <div class="stat-label">待评审</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #f0f9ff; color: #409EFF">
              <el-icon :size="28"><Operation /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value" style="color: #409EFF">{{ processingCount }}</div>
              <div class="stat-label">已分配</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #ffebee; color: #f44336">
              <el-icon :size="28"><WarningFilled /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value" style="color: #f44336">{{ overdueCount }}</div>
              <div class="stat-label">超期未处置</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 工具栏 -->
    <div class="toolbar">
      <el-input v-model="searchText" placeholder="搜索缺陷标题或组件" clearable style="width: 240px" @clear="handleSearch" @keyup.enter="handleSearch">
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>
      <el-select v-model="filterSource" placeholder="风险来源" clearable style="width: 130px; margin-left: 12px" @change="handleSearch">
        <el-option label="外部推送" value="外部推送" />
        <el-option label="SCA扫描" value="SCA扫描" />
        <el-option label="人工上报" value="人工上报" />
      </el-select>
      <el-select v-model="filterStatus" placeholder="全局状态" clearable style="width: 140px; margin-left: 12px" @change="handleSearch">
        <el-option label="待研判" value="待研判" />
        <el-option label="待评审" value="待评审" />
        <el-option label="已分配" value="已分配" />
        <el-option label="已忽略" value="已忽略" />
      </el-select>
      <el-button type="primary" plain @click="handleSearch" style="margin-left: 12px">
        <el-icon><Search /></el-icon> 查询
      </el-button>
      <el-button @click="handleReset">重置</el-button>

      <!-- 批量处置按钮 -->
      <el-divider direction="vertical" />
      <el-button type="primary" @click="showBatchDisposeDialog" :disabled="selectedRows.length === 0">
        <el-icon><EditPen /></el-icon>
        批量处置 ({{ selectedRows.length }})
      </el-button>
    </div>

    <!-- 缺陷列表 -->
    <el-table 
      :data="filteredDefects" 
      stripe 
      style="margin-top: 12px"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="source" label="来源" width="100">
        <template #default="{ row }">
          <el-tag :type="sourceType(row.source)" size="small">{{ row.source }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="风险标题" min-width="200" show-overflow-tooltip />
      <el-table-column label="风险等级" width="100">
        <template #default="{ row }"><RiskTag :level="row.riskLevel" /></template>
      </el-table-column>
      <el-table-column prop="affectedComponent" label="影响组件" width="130" />
      <el-table-column prop="discoveredDate" label="发现日期" width="110" />
      <el-table-column label="全局状态" width="110">
        <template #default="{ row }">
          <el-tag :type="globalStatusType(row.globalStatus)" size="small">{{ row.globalStatus }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="currentHandler" label="当前处理人" width="100" />
      <el-table-column label="关联任务" width="100">
        <template #default="{ row }">
          <el-link v-if="row.tasks && row.tasks.length > 0" type="primary" @click="viewRelatedTasks(row)">
            {{ row.tasks.length }}
          </el-link>
          <span v-else class="text-muted">0</span>
        </template>
      </el-table-column>
      <el-table-column prop="stayDuration" label="滞留时长" width="100">
        <template #default="{ row }">
          <span :class="{ 'stay-warning': row.stayHours > 48, 'stay-danger': row.stayHours > 72 }">{{ row.stayDuration }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="viewDefectDetail(row)">
            详情
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 批量处置对话框 -->
    <el-dialog v-model="batchDisposeDialogVisible" title="批量处置" width="900px">
      <el-alert :title="`已选择 ${selectedRows.length} 个缺陷`" type="info" :closable="false" show-icon style="margin-bottom: 16px" />
      
      <el-form label-width="120px" style="margin-bottom: 16px">
        <el-form-item label="目标状态" required>
          <el-select v-model="batchDisposeForm.targetStatus" placeholder="请选择目标状态" style="width: 100%" @change="onStatusChange">
            <el-option label="待研判" value="待研判" />
            <el-option label="待评审" value="待评审" />
            <el-option label="已分配" value="已分配" />
            <el-option label="已忽略" value="已忽略" />
          </el-select>
        </el-form-item>
      </el-form>

      <!-- 已分配状态：显示每个缺陷的卡片 -->
      <template v-if="batchDisposeForm.targetStatus === '已分配'">
        <el-card 
          v-for="(defect, index) in selectedRows" 
          :key="defect.id"
          shadow="hover"
          style="margin-bottom: 16px"
        >
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-tag :type="globalStatusType(defect.globalStatus)" size="small">{{ defect.globalStatus }}</el-tag>
                {{ defect.title }}
              </span>
              <el-tag type="info" size="small">缺陷 {{ index + 1 }}/{{ selectedRows.length }}</el-tag>
            </div>
          </template>

          <!-- 缺陷基本信息 -->
          <el-descriptions :column="4" size="small" border>
            <el-descriptions-item label="缺陷ID">{{ defect.id }}</el-descriptions-item>
            <el-descriptions-item label="影响组件">{{ defect.affectedComponent }}</el-descriptions-item>
            <el-descriptions-item label="风险等级">
              <RiskTag :level="defect.riskLevel" />
            </el-descriptions-item>
            <el-descriptions-item label="缺陷类型">
              <el-tag type="warning" size="small">{{ defect.defectType === 'dependency' ? '依赖包' : '开源基础软件' }}</el-tag>
            </el-descriptions-item>
          </el-descriptions>

          <!-- 依赖包：关联自研项目和开源基础软件列表 -->
          <template v-if="defect.defectType === 'dependency'">
            <div class="section-title">关联自研项目</div>
            <el-table :data="defectDisposeForms[defect.id].projectList" stripe size="small" style="margin-bottom: 12px">
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
            <el-table :data="defectDisposeForms[defect.id].softwareList" stripe size="small">
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
            <el-table :data="defectDisposeForms[defect.id].assigneeList" stripe size="small">
              <el-table-column prop="name" label="姓名" width="150" />
              <el-table-column prop="dept" label="部门" width="150" />
              <el-table-column prop="role" label="角色" width="150" />
              <el-table-column label="操作" width="100">
                <template #default="{ $index }">
                  <el-button type="danger" link size="small" @click="removeAssigneeFromList(defect.id, $index)">
                    移除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-button type="primary" link size="small" style="margin-top: 8px" @click="showAddAssigneeDialog(defect.id)">
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
    <el-dialog v-model="addAssigneeDialogVisible" title="添加处置人" width="500px" append-to-body>
      <el-form>
        <el-form-item label="选择处置人">
          <el-select v-model="newAssignee" placeholder="请选择处置人" style="width: 100%">
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
        <el-button @click="addAssigneeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAddAssignee" :disabled="!newAssignee">确认添加</el-button>
      </template>
    </el-dialog>

    <!-- 关联任务对话框 -->
    <el-dialog v-model="relatedTasksDialogVisible" title="关联任务" width="800px">
      <el-alert v-if="currentDefect" :title="`缺陷：${currentDefect.title}`" type="info" :closable="false" show-icon style="margin-bottom: 16px" />
      <el-table :data="relatedTasks" stripe>
        <el-table-column prop="taskNo" label="任务编号" width="160">
          <template #default="{ row }">
            <el-link type="primary" @click="viewTaskDetail(row)">{{ row.taskNo }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="任务标题" min-width="200" />
        <el-table-column prop="assignee" label="处理人" width="100" />
        <el-table-column prop="status" label="任务状态" width="100">
          <template #default="{ row }">
            <el-tag :type="taskStatusType(row.status)" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="160" />
      </el-table>
      <el-empty v-if="relatedTasks.length === 0" description="暂无关联任务" />
    </el-dialog>

    <!-- 缺陷详情对话框 -->
    <el-dialog v-model="defectDetailDialogVisible" title="缺陷详情" width="800px">
      <el-descriptions v-if="defectDetail" :column="2" border>
        <el-descriptions-item label="缺陷ID">{{ defectDetail.id }}</el-descriptions-item>
        <el-descriptions-item label="来源">
          <el-tag :type="sourceType(defectDetail.source)" size="small">{{ defectDetail.source }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="风险标题" :span="2">{{ defectDetail.title }}</el-descriptions-item>
        <el-descriptions-item label="风险等级">
          <RiskTag :level="defectDetail.riskLevel" />
        </el-descriptions-item>
        <el-descriptions-item label="影响组件">{{ defectDetail.affectedComponent }}</el-descriptions-item>
        <el-descriptions-item label="发现日期">{{ defectDetail.discoveredDate }}</el-descriptions-item>
        <el-descriptions-item label="全局状态">
          <el-tag :type="globalStatusType(defectDetail.globalStatus)" size="small">{{ defectDetail.globalStatus }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="当前处理人">{{ defectDetail.currentHandler }}</el-descriptions-item>
        <el-descriptions-item label="滞留时长">{{ defectDetail.stayDuration }}</el-descriptions-item>
        <el-descriptions-item label="关联任务" :span="2">
          <el-tag v-for="task in defectDetail.tasks" :key="task.taskNo" style="margin-right: 8px; margin-bottom: 4px">
            {{ task.taskNo }} - {{ task.title }}
          </el-tag>
          <span v-if="!defectDetail.tasks || defectDetail.tasks.length === 0" class="text-muted">无</span>
        </el-descriptions-item>
        <el-descriptions-item label="关联组件" :span="2">
          <div v-if="defectDetail.affectedComponent">
            <el-tag type="info">{{ defectDetail.affectedComponent }}</el-tag>
          </div>
          <span v-else class="text-muted">无</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, Clock, Tickets, Operation, WarningFilled, EditPen, Plus } from '@element-plus/icons-vue'
import RiskTag from '@/components/business/RiskTag.vue'

const router = useRouter()
const searchText = ref('')
const filterSource = ref('')
const filterStatus = ref('')
const selectedRows = ref<any[]>([])

// Mock数据
const defects = ref([
  {
    id: 'DEF-001',
    source: '外部推送',
    title: 'Apache Log4j2 远程代码执行漏洞',
    riskLevel: 'critical',
    defectType: 'dependency',
    affectedComponent: 'log4j-core',
    discoveredDate: '2024-06-15',
    globalStatus: '待研判',
    currentHandler: '系统',
    stayDuration: '2小时',
    stayHours: 2,
    tasks: [],
  },
  {
    id: 'DEF-002',
    source: 'SCA扫描',
    title: 'Spring Framework 权限绕过漏洞',
    riskLevel: 'high',
    defectType: 'software',
    affectedComponent: 'spring-core',
    discoveredDate: '2024-06-14',
    globalStatus: '待评审',
    currentHandler: '王五',
    stayDuration: '26小时',
    stayHours: 26,
    tasks: [],
  },
  {
    id: 'DEF-003',
    source: '外部推送',
    title: 'OpenSSL 缓冲区溢出漏洞',
    riskLevel: 'critical',
    defectType: 'dependency',
    affectedComponent: 'openssl',
    discoveredDate: '2024-06-13',
    globalStatus: '已分配',
    currentHandler: '张三',
    stayDuration: '50小时',
    stayHours: 50,
    tasks: [
      { taskNo: 'TASK-20240613-001', title: '修复OpenSSL漏洞', assignee: '张三', status: '处理中', createdAt: '2024-06-13 10:30' },
      { taskNo: 'TASK-20240613-002', title: '升级OpenSSL版本', assignee: '李四', status: '待处理', createdAt: '2024-06-13 10:30' },
    ],
  },
  {
    id: 'DEF-004',
    source: 'SCA扫描',
    title: 'FastJSON 反序列化漏洞',
    riskLevel: 'high',
    defectType: 'dependency',
    affectedComponent: 'fastjson',
    discoveredDate: '2024-06-12',
    globalStatus: '已忽略',
    currentHandler: '系统',
    stayDuration: '72小时',
    stayHours: 72,
    tasks: [],
  },
  {
    id: 'DEF-005',
    source: '人工上报',
    title: 'MySQL 权限提升漏洞',
    riskLevel: 'medium',
    defectType: 'software',
    affectedComponent: 'mysql',
    discoveredDate: '2024-06-11',
    globalStatus: '已忽略',
    currentHandler: '赵六',
    stayDuration: '96小时',
    stayHours: 96,
    tasks: [
      { taskNo: 'TASK-20240611-001', title: '修复MySQL权限漏洞', assignee: '赵六', status: '已完成', createdAt: '2024-06-11 14:20' },
    ],
  },
])

// 统计数据
const pendingReviewCount = computed(() => defects.value.filter(d => d.globalStatus === '待研判').length)
const pendingMeetingCount = computed(() => defects.value.filter(d => d.globalStatus === '待评审').length)
const processingCount = computed(() => defects.value.filter(d => d.globalStatus === '已分配').length)
const overdueCount = computed(() => defects.value.filter(d => d.stayHours > 72).length)

const filteredDefects = computed(() => {
  return defects.value.filter(item => {
    const matchSearch = !searchText.value || item.title.toLowerCase().includes(searchText.value.toLowerCase()) || item.affectedComponent.toLowerCase().includes(searchText.value.toLowerCase())
    const matchSource = !filterSource.value || item.source === filterSource.value
    const matchStatus = !filterStatus.value || item.globalStatus === filterStatus.value
    return matchSearch && matchSource && matchStatus
  })
})

function sourceType(source: string): 'danger' | 'warning' | 'info' {
  const map: Record<string, 'danger' | 'warning' | 'info'> = { '外部推送': 'danger', 'SCA扫描': 'warning', '人工上报': 'info' }
  return map[source] || 'info'
}

function globalStatusType(status: string) {
  const map: Record<string, string> = {
    '待研判': 'warning',
    '待评审': 'info',
    '已分配': 'primary',
    '已忽略': 'danger',
  }
  return map[status] || ''
}

function taskStatusType(status: string) {
  const map: Record<string, string> = {
    '待处理': 'info',
    '处理中': 'primary',
    '已完成': 'success',
  }
  return map[status] || ''
}

function handleSearch() {}

function handleReset() {
  searchText.value = ''
  filterSource.value = ''
  filterStatus.value = ''
}

function handleSelectionChange(rows: any[]) {
  selectedRows.value = rows
}

// 批量处置
const batchDisposeDialogVisible = ref(false)
const batchDisposeForm = ref({
  targetStatus: '',
  remark: '',
})

// 每个缺陷的处置表单数据
const defectDisposeForms = ref<Record<string, {
  projectList: Array<{ name: string; dept: string; owner: string }>
  softwareList: Array<{ name: string; version: string; owner: string }>
  assigneeList: Array<{ name: string; dept: string; role: string }>
}>>({})

// Mock数据：项目列表
const mockProjectList = [
  { name: '订单管理系统', dept: '技术研发部', owners: ['张三'] },
  { name: '支付系统', dept: '技术研发部', owners: ['李四'] },
  { name: '用户中心', dept: '平台研发部', owners: ['王五'] },
  { name: '消息平台', dept: '平台研发部', owners: ['赵六'] },
]

// Mock数据：软件列表
const mockSoftwareList = [
  { name: 'Redis', version: '7.0', owners: ['张三'] },
  { name: 'Spring Boot', version: '3.2', owners: ['李四'] },
  { name: 'MySQL', version: '8.0', owners: ['王五'] },
  { name: 'Nginx', version: '1.25', owners: ['赵六'] },
]

// Mock数据：处置人列表
const mockAssigneeList = [
  { name: '张三', dept: '技术研发部', role: '软件负责人' },
  { name: '李四', dept: '平台研发部', role: '软件负责人' },
  { name: '王五', dept: '运维部', role: '运维工程师' },
  { name: '赵六', dept: '安全部', role: '软件负责人' },
]

// 添加处置人对话框
const addAssigneeDialogVisible = ref(false)
const newAssignee = ref('')
const currentDefectId = ref('')

// 项目/软件负责人的映射关系
const projectOwners: Record<string, string> = {
  '订单管理系统': '张三',
  '支付系统': '李四',
  '用户中心': '王五',
  '消息平台': '赵六',
}

const softwareOwners: Record<string, string> = {
  'Redis': '张三',
  'Spring Boot': '李四',
  'MySQL': '王五',
  'Nginx': '赵六',
}

function showBatchDisposeDialog() {
  console.log('selectedRows:', selectedRows.value)
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先选择要处置的缺陷')
    return
  }
  
  batchDisposeForm.value = { targetStatus: '', remark: '' }
  defectDisposeForms.value = {}
  
  // 为每个选中的缺陷初始化表单
  selectedRows.value.forEach(defect => {
    defectDisposeForms.value[defect.id] = {
      projectList: defect.defectType === 'dependency' ? [...mockProjectList] : [],
      softwareList: defect.defectType === 'dependency' ? [...mockSoftwareList] : [],
      assigneeList: defect.defectType === 'software' ? [...mockAssigneeList] : [],
    }
  })
  
  console.log('defectDisposeForms:', defectDisposeForms.value)
  batchDisposeDialogVisible.value = true
}

function onStatusChange() {
  // 切换状态时清空所有缺陷的表单数据
  if (batchDisposeForm.value.targetStatus !== '已分配') {
    Object.keys(defectDisposeForms.value).forEach(defectId => {
      const defect = selectedRows.value.find(d => d.id === defectId)
      if (defect) {
        defectDisposeForms.value[defectId] = {
          projectList: defect.defectType === 'dependency' ? [...mockProjectList] : [],
          softwareList: defect.defectType === 'dependency' ? [...mockSoftwareList] : [],
          assigneeList: defect.defectType === 'software' ? [...mockAssigneeList] : [],
        }
      }
    })
  }
}

function removeAssigneeFromList(defectId: string, index: number) {
  defectDisposeForms.value[defectId].assigneeList.splice(index, 1)
}

function showAddAssigneeDialog(defectId: string) {
  currentDefectId.value = defectId
  newAssignee.value = ''
  addAssigneeDialogVisible.value = true
}

function confirmAddAssignee() {
  if (newAssignee.value) {
    const deptMap: Record<string, string> = {
      '张三': '技术研发部',
      '李四': '平台研发部',
      '王五': '运维部',
      '赵六': '安全部',
      '钱七': '质量管理部',
      '孙八': '数据团队',
    }
    defectDisposeForms.value[currentDefectId.value].assigneeList.push({
      name: newAssignee.value,
      dept: deptMap[newAssignee.value] || '',
      role: '处置人',
    })
  }
  addAssigneeDialogVisible.value = false
}

function handleBatchDispose() {
  const targetStatus = batchDisposeForm.value.targetStatus
  
  selectedRows.value.forEach(row => {
    const idx = defects.value.findIndex(d => d.id === row.id)
    if (idx !== -1) {
      defects.value[idx].globalStatus = targetStatus
      
      // 如果状态为已分配，创建任务
      const form = defectDisposeForms.value[row.id]
      if (targetStatus === '已分配' && form) {
        // 收集所有处置人
        const assignees: string[] = []
        if (row.defectType === 'dependency') {
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
        } else if (row.defectType === 'software') {
          // 软件：从处置人列表中收集
          form.assigneeList.forEach(a => assignees.push(a.name))
        }
        
        // 去重
        const uniqueAssignees = [...new Set(assignees)]
        
        if (uniqueAssignees.length > 0) {
          // 任务标题格式：名称-版本-缺陷编号
          const taskTitle = `${row.affectedComponent}-${defects.value[idx].discoveredDate}-${row.id}`
          const newTasks = uniqueAssignees.map((assignee, i) => ({
            taskNo: `TASK-${Date.now()}-${i}`,
            title: taskTitle,
            assignee,
            status: '待处理',
            createdAt: new Date().toLocaleString('zh-CN'),
          }))
          defects.value[idx].tasks = [...(defects.value[idx].tasks || []), ...newTasks]
          defects.value[idx].currentHandler = uniqueAssignees.join(', ')
        }
      }
    }
  })
  
  ElMessage.success(`已将 ${selectedRows.value.length} 个缺陷处置为"${targetStatus}"状态`)
  batchDisposeDialogVisible.value = false
  selectedRows.value = []
}

// 查看关联任务
const relatedTasksDialogVisible = ref(false)
const currentDefect = ref<any>(null)
const relatedTasks = ref<any[]>([])

function viewRelatedTasks(row: any) {
  currentDefect.value = row
  relatedTasks.value = row.tasks || []
  relatedTasksDialogVisible.value = true
}

function viewTaskDetail(task: any) {
  router.push(`/admin/governance/task-list`)
  ElMessage.info(`跳转到任务：${task.taskNo}`)
}

// 查看缺陷详情
const defectDetailDialogVisible = ref(false)
const defectDetail = ref<any>(null)

function viewDefectDetail(row: any) {
  defectDetail.value = row
  defectDetailDialogVisible.value = true
}
</script>

<style scoped>
.page-header { margin-bottom: 8px; }
.page-header h3 { margin: 0 0 4px; }
.desc { color: #909399; font-size: 13px; margin: 0; }

.stat-card {
  border: none;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
}

.toolbar { 
  display: flex; 
  align-items: center; 
  gap: 12px; 
  flex-wrap: wrap;
  margin-top: 16px;
}

.text-muted {
  color: #c0c4cc;
}

.stay-warning {
  color: #E6A23C;
}

.stay-danger {
  color: #F56C6C;
  font-weight: 600;
}

.assignee-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
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
