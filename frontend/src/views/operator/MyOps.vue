<template>
  <div class="my-ops">
    <div class="page-header">
      <h2>我的运维</h2>
      <p class="page-desc">我负责维护的开源基础软件</p>
    </div>

    <el-tabs v-model="activeTab" class="ops-tabs">
      <!-- 我负责的软件 -->
      <el-tab-pane label="我负责的软件" name="software">
        <el-table :data="softwareList" border v-loading="loading">
          <el-table-column prop="grade" label="等级" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="gradeType(row.grade)" size="small" effect="dark">{{ row.grade }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="name" label="软件名称" min-width="160">
            <template #default="{ row }">
              <el-link type="primary" @click="showSoftwareDetail(row)">{{ row.name }}</el-link>
            </template>
          </el-table-column>
          <el-table-column prop="version" label="当前版本" width="120" />
          <el-table-column prop="type" label="类型" width="120">
            <template #default="{ row }"><el-tag size="small">{{ row.type }}</el-tag></template>
          </el-table-column>
          <el-table-column prop="owner" label="负责人" width="100" />
          <el-table-column prop="defectCount" label="缺陷数" width="100" align="center">
            <template #default="{ row }">
              <el-tag v-if="row.defectCount > 0" type="danger" size="small">{{ row.defectCount }}</el-tag>
              <span v-else>0</span>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- 缺陷管理（本部门缺陷，只读；内部/外部分tab） -->
      <el-tab-pane label="缺陷管理" name="defects">
        <el-tabs v-model="defectActiveTab" size="small">
          <!-- 内部缺陷 -->
          <el-tab-pane label="内部缺陷" name="internal">
            <div style="margin-bottom: 16px; display: flex; gap: 12px; align-items: center; flex-wrap: wrap">
              <el-input v-model="internalKeyword" placeholder="搜索缺陷标题" clearable style="width: 220px" @keyup.enter="loadInternalDefects">
                <template #prefix><el-icon><Search /></el-icon></template>
              </el-input>
              <el-select v-model="internalStatus" placeholder="全局状态" clearable style="width: 140px" @change="loadInternalDefects">
                <el-option v-for="(label, value) in DEFECT_STATUS_LABELS" :key="value" :label="label" :value="value" />
              </el-select>
              <el-button type="primary" plain @click="internalPage = 1; loadInternalDefects()"><el-icon><Search /></el-icon> 查询</el-button>
              <el-button @click="resetInternalFilter">重置</el-button>
            </div>
            <el-table :data="internalDefectList" border v-loading="internalLoading">
              <el-table-column prop="defectNo" label="缺陷编号" width="170" />
              <el-table-column prop="source" label="来源" width="110" />
              <el-table-column prop="title" label="风险标题" min-width="200" show-overflow-tooltip />
              <el-table-column label="风险等级" width="100" align="center">
                <template #default="{ row }">
                  <el-tag :type="severityType(row.riskLevel)" size="small">{{ severityLabel(row.riskLevel) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="maintDept" label="处置部门" width="120">
                <template #default="{ row }">{{ row.maintDept || '-' }}</template>
              </el-table-column>
              <el-table-column label="全局状态" width="110" align="center">
                <template #default="{ row }">
                  <el-tag :type="defectStatusType(row.status)" size="small">{{ defectStatusLabel(row.status) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="100" align="center">
                <template #default="{ row }">
                  <el-button type="primary" link size="small" @click="viewDefectDetail(row)">详情</el-button>
                </template>
              </el-table-column>
            </el-table>
            <div class="pagination-row">
              <el-pagination
                v-model:current-page="internalPage"
                v-model:page-size="internalPageSize"
                :total="internalTotal"
                :page-sizes="[10, 20, 50, 100]"
                layout="total, sizes, prev, pager, next"
                @current-change="loadInternalDefects"
                @size-change="onInternalPageSizeChange"
              />
            </div>
          </el-tab-pane>

          <!-- 外部缺陷 -->
          <el-tab-pane label="外部缺陷" name="external">
            <div style="margin-bottom: 16px; display: flex; gap: 12px; align-items: center; flex-wrap: wrap">
              <el-input v-model="externalKeyword" placeholder="搜索缺陷标题或软件" clearable style="width: 220px" @keyup.enter="loadExternalDefects">
                <template #prefix><el-icon><Search /></el-icon></template>
              </el-input>
              <el-select v-model="externalStatus" placeholder="全局状态" clearable style="width: 140px" @change="loadExternalDefects">
                <el-option v-for="(label, value) in DEFECT_STATUS_LABELS" :key="value" :label="label" :value="value" />
              </el-select>
              <el-button type="primary" plain @click="externalPage = 1; loadExternalDefects()"><el-icon><Search /></el-icon> 查询</el-button>
              <el-button @click="resetExternalFilter">重置</el-button>
            </div>
            <el-table :data="externalDefectList" border v-loading="externalLoading">
              <el-table-column prop="defectNo" label="缺陷编号" width="180" />
              <el-table-column prop="software" label="软件" width="110" />
              <el-table-column prop="title" label="缺陷标题" min-width="200" show-overflow-tooltip />
              <el-table-column label="严重程度" width="100" align="center">
                <template #default="{ row }">
                  <el-tag :type="severityType(row.riskLevel)" size="small">{{ severityLabel(row.riskLevel) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="maintDept" label="运维部门" width="120">
                <template #default="{ row }">{{ row.maintDept || '-' }}</template>
              </el-table-column>
              <el-table-column label="全局状态" width="110" align="center">
                <template #default="{ row }">
                  <el-tag :type="defectStatusType(row.status)" size="small">{{ defectStatusLabel(row.status) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="100" align="center">
                <template #default="{ row }">
                  <el-button type="primary" link size="small" @click="viewDefectDetail(row)">详情</el-button>
                </template>
              </el-table-column>
            </el-table>
            <div class="pagination-row">
              <el-pagination
                v-model:current-page="externalPage"
                v-model:page-size="externalPageSize"
                :total="externalTotal"
                :page-sizes="[10, 20, 50, 100]"
                layout="total, sizes, prev, pager, next"
                @current-change="loadExternalDefects"
                @size-change="onExternalPageSizeChange"
              />
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-tab-pane>
    </el-tabs>

    <!-- 软件详情对话框 -->
    <ProjectDetailDialog v-model="softwareDetailVisible" :project-data="currentSoftware" />

    <!-- 缺陷详情对话框 -->
    <el-dialog v-model="defectDetailDialogVisible" :title="isExternal(currentDefect) ? '外部缺陷详情' : '内部缺陷详情'" width="720px">
      <el-descriptions v-if="currentDefect" :column="2" border size="small">
        <el-descriptions-item label="缺陷编号">{{ currentDefect.defectNo }}</el-descriptions-item>
        <el-descriptions-item label="来源">{{ currentDefect.source }}</el-descriptions-item>
        <el-descriptions-item label="风险标题" :span="2">{{ currentDefect.title }}</el-descriptions-item>
        <template v-if="isExternal(currentDefect)">
          <el-descriptions-item label="软件">
            <el-tag size="small">{{ currentDefect.software }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="严重程度">
            <el-tag :type="severityType(currentDefect.riskLevel)" size="small">{{ severityLabel(currentDefect.riskLevel) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="原始链接" :span="2">
            <el-link v-if="currentDefect.issueUrl" :href="currentDefect.issueUrl" target="_blank" type="primary">{{ currentDefect.issueUrl }}</el-link>
            <span v-else>-</span>
          </el-descriptions-item>
        </template>
        <template v-else>
          <el-descriptions-item label="风险等级">
            <el-tag :type="severityType(currentDefect.riskLevel)" size="small">{{ severityLabel(currentDefect.riskLevel) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="上报时间">{{ formatTime(currentDefect.discoveredDate) }}</el-descriptions-item>
        </template>
        <el-descriptions-item label="处置部门">{{ currentDefect.maintDept || '-' }}</el-descriptions-item>
        <el-descriptions-item label="处置截止时间">{{ formatTime(currentDefect.deadline) }}</el-descriptions-item>
        <el-descriptions-item label="全局状态">
          <el-tag :type="defectStatusType(currentDefect.status)" size="small">{{ defectStatusLabel(currentDefect.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentDefect.description || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { riskApi } from '@/api'
import request from '@/api'
import { formatTime } from '@/utils/format'
import {
  DEFECT_STATUS_LABELS,
  FLOW_TYPE_LABELS,
  defectStatusLabel,
  defectStatusType,
  severityLabel,
  severityType,
} from '@/utils/defect'
import ProjectDetailDialog from '@/components/business/ProjectDetailDialog.vue'

const userStore = useUserStore()

const activeTab = ref('software')
const loading = ref(false)
const softwareList = ref<any[]>([])

async function loadSoftwareList() {
  loading.value = true
  try {
    const params: any = {}
    const role = userStore.role
    const isSpecialistView = role === 'specialist' || (role === 'admin' && userStore.viewMode === 'personal')
    if (isSpecialistView) {
      params.maint_dept = userStore.userInfo?.department || ''
    } else {
      params.owner_username = userStore.userInfo?.username || ''
    }
    const res: any = await request.get('/asset/projects', { params })
    softwareList.value = (res || []).map((item: any) => ({
      ...item,
      name: item._name,
      version: item._version,
      grade: item.applicable_level || 'A',
      type: getTypeLabel(item.type),
      owner: item.create_by || '-',
      department: item.maint_dept || '-',
      description: item.DESCRIPTION || '-',
      defectCount: 0,
    }))
  } catch (error) {
    ElMessage.error('加载软件列表失败')
  } finally {
    loading.value = false
  }
}

function getTypeLabel(type: string) {
  const map: Record<string, string> = {
    java: 'Java组件', go: 'Go组件', python: 'Python组件', js: 'JavaScript组件',
  }
  return map[type] || type || '-'
}

function gradeType(grade: string) {
  const map: Record<string, any> = { A: 'success', B: '', C: 'warning', D: 'danger' }
  return map[grade] || 'info'
}

// 软件详情
const softwareDetailVisible = ref(false)
const currentSoftware = ref<any>(null)
function showSoftwareDetail(software: any) {
  currentSoftware.value = software
  softwareDetailVisible.value = true
}

// ===== 缺陷管理（内部/外部分tab） =====
const defectActiveTab = ref('internal')
const internalLoading = ref(false)
const externalLoading = ref(false)
const internalKeyword = ref('')
const internalStatus = ref('')
const externalKeyword = ref('')
const externalStatus = ref('')
const internalDefectList = ref<any[]>([])
const externalDefectList = ref<any[]>([])
const internalPage = ref(1)
const internalPageSize = ref(10)
const internalTotal = ref(0)
const externalPage = ref(1)
const externalPageSize = ref(10)
const externalTotal = ref(0)

// 内部缺陷 vs 外部缺陷 判定
function isExternal(row: any): boolean {
  if (!row) return false
  return row.affectedComponentType === 'external' || /#/.test(row.defectNo || '') || row.source === 'github-issue' || row.source === 'jira'
}

async function loadInternalDefects() {
  internalLoading.value = true
  try {
    const res: any = await riskApi.getInternalDefects({
      page: internalPage.value,
      pageSize: internalPageSize.value,
      keyword: internalKeyword.value || undefined,
      status: internalStatus.value || undefined,
      // 个人工作台视图：管理部门专员也与部门专员一致按部门过滤，仅在全局治理（风险管理-缺陷台账）看全景
      scope: 'dept',
    })
    internalDefectList.value = Array.isArray(res) ? res : (Array.isArray(res?.list) ? res.list : [])
    internalTotal.value = res?.total ?? internalDefectList.value.length
  } catch (e: any) {
    ElMessage.error('加载内部缺陷失败: ' + (e.message || ''))
  } finally {
    internalLoading.value = false
  }
}

async function loadExternalDefects() {
  externalLoading.value = true
  try {
    const res: any = await riskApi.getExternalDefects({
      page: externalPage.value,
      pageSize: externalPageSize.value,
      keyword: externalKeyword.value || undefined,
      status: externalStatus.value || undefined,
      scope: 'dept',
    })
    externalDefectList.value = Array.isArray(res) ? res : (Array.isArray(res?.list) ? res.list : [])
    externalTotal.value = res?.total ?? externalDefectList.value.length
  } catch (e: any) {
    ElMessage.error('加载外部缺陷失败: ' + (e.message || ''))
  } finally {
    externalLoading.value = false
  }
}

function onInternalPageSizeChange() {
  internalPage.value = 1
  loadInternalDefects()
}
function onExternalPageSizeChange() {
  externalPage.value = 1
  loadExternalDefects()
}

function loadDefects() {
  loadInternalDefects()
  loadExternalDefects()
}

function resetInternalFilter() {
  internalKeyword.value = ''
  internalStatus.value = ''
  internalPage.value = 1
  loadInternalDefects()
}
function resetExternalFilter() {
  externalKeyword.value = ''
  externalStatus.value = ''
  externalPage.value = 1
  loadExternalDefects()
}

const defectDetailDialogVisible = ref(false)
const currentDefect = ref<any>(null)
function viewDefectDetail(row: any) {
  currentDefect.value = row
  defectDetailDialogVisible.value = true
}

onMounted(() => {
  // 支持从「风险上报」创建成功跳转定位：localStorage 标记直接定位到缺陷管理-内部缺陷
  if (localStorage.getItem('osm_defect_tab') === 'defects') {
    activeTab.value = 'defects'
    defectActiveTab.value = 'internal'
    localStorage.removeItem('osm_defect_tab')
  }
  loadSoftwareList()
  loadDefects()
})
</script>

<style scoped>
.my-ops { padding: 20px; }
.pagination-row { display: flex; justify-content: flex-end; margin-top: 12px; }
.page-header { margin-bottom: 20px; }
.page-header h2 { margin: 0 0 8px; font-size: 18px; font-weight: 600; color: #1f2328; }
.page-desc { margin: 0; font-size: 13px; color: #8c8c8c; }
.ops-tabs { background: #fff; padding: 20px; border-radius: 4px; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08); }
</style>
