<template>
  <div class="core-dependencies-page">
    <!-- 第一部分：看板指标 -->
    <el-row :gutter="16" class="metric-row">
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">使用占比大于50%的组件数</span>
            <span class="metric-value danger">{{ summary.highAdoptionCount }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">碎片化组件数(≥3版本)</span>
            <span class="metric-value warning">{{ summary.fragmentedCompCount }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">高风险项目数</span>
            <span class="metric-value danger">{{ summary.highRiskProjectCount }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-content">
            <span class="metric-label">强传染许可证组件数</span>
            <span class="metric-value primary">{{ summary.copyleftCompCount }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 第二部分：Tab 区分统计类型 -->
    <el-card shadow="never">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <!-- Tab1: 组件项目关联 -->
        <el-tab-pane label="组件项目关联" name="assoc">
          <div class="tab-toolbar">
            <div class="search-area">
              <el-input v-model="searchOrgName" placeholder="组织" clearable style="width: 180px" @keyup.enter="loadAssoc" />
              <el-input v-model="searchCompName" placeholder="组件名称" clearable style="width: 180px; margin-left: 8px" @keyup.enter="loadAssoc" />
              <el-input v-model="searchProjectName" placeholder="项目名称" clearable style="width: 180px; margin-left: 8px" @keyup.enter="loadAssoc" />
              <el-button type="primary" @click="loadAssoc" style="margin-left: 8px">
                <el-icon><Search /></el-icon> 查询
              </el-button>
              <el-button @click="resetSearch" style="margin-left: 4px">重置</el-button>
            </div>
            <el-button type="warning" plain @click="openRankingDialog">
              <el-icon><Trophy /></el-icon> 排行榜
            </el-button>
          </div>

          <el-table :data="assocList" stripe v-loading="assocLoading" style="width: 100%; margin-top: 12px">
            <el-table-column prop="orgName" label="组织" min-width="180" show-overflow-tooltip />
            <el-table-column prop="compName" label="名称" min-width="160" show-overflow-tooltip />
            <el-table-column prop="version" label="版本" width="120" />
            <el-table-column prop="projectName" label="项目名" min-width="140" show-overflow-tooltip />
            <el-table-column prop="moduleCode" label="模块" min-width="140" show-overflow-tooltip />
            <el-table-column prop="groupName" label="团队" width="100" />
            <el-table-column prop="deptName" label="部门" width="120" />
          </el-table>

          <el-pagination
            v-model:current-page="assocPage"
            v-model:page-size="assocPageSize"
            :total="assocTotal"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            style="margin-top: 16px; justify-content: flex-end"
            @size-change="loadAssoc"
            @current-change="loadAssoc"
          />
        </el-tab-pane>

        <!-- Tab2: 版本碎片化 -->
        <el-tab-pane label="版本碎片化" name="fragmentation">
          <div class="tab-toolbar">
            <div class="search-area">
              <el-input v-model="vfSearchCompName" placeholder="组件名称" clearable style="width: 200px" @keyup.enter="loadFragmentation" />
              <el-button type="primary" @click="loadFragmentation" style="margin-left: 8px">
                <el-icon><Search /></el-icon> 查询
              </el-button>
              <el-button @click="resetVfSearch" style="margin-left: 4px">重置</el-button>
            </div>
          </div>

          <el-table :data="vfList" stripe v-loading="vfLoading" style="width: 100%; margin-top: 12px">
            <el-table-column type="index" label="排名" width="60" />
            <el-table-column prop="compName" label="组件名称" min-width="160" show-overflow-tooltip />
            <el-table-column prop="versionCount" label="在用版本数" width="120" sortable>
              <template #default="{ row }">
                <el-tag :type="getVersionCountType(row.versionCount)" size="small">{{ row.versionCount }} 个版本</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="版本列表" min-width="280">
              <template #default="{ row }">
                <div class="version-tags">
                  <el-tag
                    v-for="ver in parseVersionList(row.versionList, row.latestVersion)"
                    :key="ver"
                    :type="ver === row.latestVersion ? 'success' : 'warning'"
                    size="small"
                    class="version-tag"
                  >
                    {{ ver }}
                  </el-tag>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="latestVersion" label="最新版本" width="140" />
            <el-table-column label="碎片化程度" width="120">
              <template #default="{ row }">
                <el-tag :type="getFragmentTagType(row.fragmentLevel)" size="small" effect="dark">
                  {{ getFragmentLabel(row.fragmentLevel) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="vfPage"
            v-model:page-size="vfPageSize"
            :total="vfTotal"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            style="margin-top: 16px; justify-content: flex-end"
            @size-change="loadFragmentation"
            @current-change="loadFragmentation"
          />
        </el-tab-pane>

        <!-- Tab3: 高风险项目 -->
        <el-tab-pane label="高风险项目" name="highRisk">
          <div class="tab-toolbar">
            <div class="search-area">
              <el-input v-model="hrSearchProjectName" placeholder="项目名称" clearable style="width: 200px" @keyup.enter="loadHighRisk" />
              <el-button type="primary" @click="loadHighRisk" style="margin-left: 8px">
                <el-icon><Search /></el-icon> 查询
              </el-button>
              <el-button @click="resetHrSearch" style="margin-left: 4px">重置</el-button>
            </div>
          </div>

          <el-table :data="hrList" stripe v-loading="hrLoading" style="width: 100%; margin-top: 12px" @expand-change="handleHrExpand">
            <el-table-column type="expand">
              <template #default="{ row }">
                <div class="expand-detail">
                  <el-table :data="hrDetailMap[row.projectName] || []" size="small" border style="width: 100%">
                    <el-table-column prop="moduleCode" label="模块名称" min-width="140" />
                    <el-table-column prop="commitId" label="commitID（当前版本）" min-width="160" show-overflow-tooltip />
                    <el-table-column prop="createTime" label="核查时间" width="170" />
                    <el-table-column label="处置" width="90">
                      <template #default="{ row: detail }">
                        <el-tag :type="detail.alertType === 1 ? 'danger' : 'warning'" size="small">
                          {{ detail.alertType === 1 ? '阻断' : '告警' }}
                        </el-tag>
                      </template>
                    </el-table-column>
                    <el-table-column prop="blacklistName" label="黑名单" width="120" />
                    <el-table-column prop="reason" label="原因" min-width="200" show-overflow-tooltip />
                  </el-table>
                  <el-empty v-if="hrDetailMap[row.projectName] && hrDetailMap[row.projectName].length === 0" description="暂无黑名单核查记录" :image-size="60" />
                </div>
              </template>
            </el-table-column>
            <el-table-column type="index" label="排名" width="60" />
            <el-table-column prop="projectName" label="项目名称" min-width="180" show-overflow-tooltip />
            <el-table-column prop="blacklistCount" label="黑名单数" width="120" sortable>
              <template #default="{ row }">
                <el-tag :type="getBlacklistTagType(row.blacklistCount)" size="small">{{ row.blacklistCount }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="groupName" label="团队" width="120" />
            <el-table-column prop="deptName" label="部门" width="150" />
          </el-table>

          <el-pagination
            v-model:current-page="hrPage"
            v-model:page-size="hrPageSize"
            :total="hrTotal"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            style="margin-top: 16px; justify-content: flex-end"
            @size-change="loadHighRisk"
            @current-change="loadHighRisk"
          />
        </el-tab-pane>

        <!-- Tab4: 许可证关联 -->
        <el-tab-pane label="许可证关联" name="license">
          <div class="tab-toolbar">
            <div class="search-area">
              <el-input v-model="licSearchName" placeholder="许可证名称" clearable style="width: 200px" @keyup.enter="loadLicense" />
              <el-button type="primary" @click="loadLicense" style="margin-left: 8px">
                <el-icon><Search /></el-icon> 查询
              </el-button>
              <el-button @click="resetLicSearch" style="margin-left: 4px">重置</el-button>
            </div>
          </div>

          <el-table :data="licList" stripe v-loading="licLoading" style="width: 100%; margin-top: 12px">
            <el-table-column type="index" label="排名" width="60" />
            <el-table-column prop="licenseName" label="许可证名称" min-width="180" show-overflow-tooltip />
            <el-table-column prop="componentCount" label="关联组件数" width="140" sortable>
              <template #default="{ row }">
                <el-tag :type="getLicenseCountTagType(row.componentCount)" size="small">{{ row.componentCount }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="风险等级" width="120">
              <template #default="{ row }">
                <el-tag :type="getLicenseRiskTagType(row.riskLevel)" effect="dark" size="small">
                  {{ getLicenseRiskLabel(row.riskLevel) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="licPage"
            v-model:page-size="licPageSize"
            :total="licTotal"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            style="margin-top: 16px; justify-content: flex-end"
            @size-change="loadLicense"
            @current-change="loadLicense"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 排行榜弹框 -->
    <el-dialog v-model="rankingDialogVisible" title="组件项目数排行榜" width="850px" :close-on-click-modal="false">
      <div class="ranking-toolbar">
        <el-input v-model="rankingSearchOrg" placeholder="组织" clearable style="width: 160px" @keyup.enter="loadRanking" />
        <el-input v-model="rankingSearchComp" placeholder="组件名称" clearable style="width: 160px; margin-left: 8px" @keyup.enter="loadRanking" />
        <el-button type="primary" @click="loadRanking" style="margin-left: 8px">
          <el-icon><Search /></el-icon> 查询
        </el-button>
      </div>

      <el-table :data="rankingList" stripe v-loading="rankingLoading" style="width: 100%; margin-top: 12px">
        <el-table-column type="index" label="排名" width="60" />
        <el-table-column prop="compName" label="组件名称" min-width="160" show-overflow-tooltip />
        <el-table-column prop="projectCount" label="项目数" width="90" sortable>
          <template #default="{ row }">
            <el-tag type="primary" size="small">{{ row.projectCount }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="使用占比" width="180">
          <template #default="{ row }">
            <el-progress :percentage="row.usageRate" :stroke-width="8" :color="getUsageColor(row.usageRate)" />
          </template>
        </el-table-column>
        <el-table-column prop="latestVersion" label="最新版本" width="130" />
      </el-table>

      <el-pagination
        v-model:current-page="rankingPage"
        v-model:page-size="rankingPageSize"
        :total="rankingTotal"
        :page-sizes="[10, 20, 50]"
        layout="total, prev, pager, next"
        style="margin-top: 16px; justify-content: flex-end"
        @current-change="loadRanking"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Trophy } from '@element-plus/icons-vue'
import { coreDependencyApi } from '@/api'

// 看板汇总 - 供应链安全洞察
const summary = reactive({
  highAdoptionCount: 0,
  fragmentedCompCount: 0,
  highRiskProjectCount: 0,
  copyleftCompCount: 0,
})

// Tab 状态
const activeTab = ref('assoc')

// 组件项目关联
const assocLoading = ref(false)
const assocList = ref<any[]>([])
const assocPage = ref(1)
const assocPageSize = ref(10)
const assocTotal = ref(0)
const searchOrgName = ref('')
const searchCompName = ref('')
const searchProjectName = ref('')

// 排行榜弹框
const rankingDialogVisible = ref(false)
const rankingLoading = ref(false)
const rankingList = ref<any[]>([])
const rankingPage = ref(1)
const rankingPageSize = ref(10)
const rankingTotal = ref(0)
const rankingSearchOrg = ref('')
const rankingSearchComp = ref('')

// 版本碎片化
const vfLoading = ref(false)
const vfList = ref<any[]>([])
const vfPage = ref(1)
const vfPageSize = ref(10)
const vfTotal = ref(0)
const vfSearchCompName = ref('')

// 高风险项目
const hrLoading = ref(false)
const hrList = ref<any[]>([])
const hrPage = ref(1)
const hrPageSize = ref(10)
const hrTotal = ref(0)
const hrSearchProjectName = ref('')
const hrDetailMap = reactive<Record<string, any[]>>({})

// 许可证关联
const licLoading = ref(false)
const licList = ref<any[]>([])
const licPage = ref(1)
const licPageSize = ref(10)
const licTotal = ref(0)
const licSearchName = ref('')

onMounted(() => {
  loadSummary()
  loadAssoc()
})

async function loadSummary() {
  try {
    const res = await coreDependencyApi.getSummary()
    Object.assign(summary, res)
  } catch (e: any) {
    ElMessage.error(e.message || '加载看板数据失败')
  }
}

async function loadAssoc() {
  assocLoading.value = true
  try {
    const res = await coreDependencyApi.getAssoc({
      page: assocPage.value,
      pageSize: assocPageSize.value,
      orgName: searchOrgName.value,
      compName: searchCompName.value,
      projectName: searchProjectName.value,
    })
    assocList.value = res.list || []
    assocTotal.value = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载组件项目关联失败')
  } finally {
    assocLoading.value = false
  }
}

function resetSearch() {
  searchOrgName.value = ''
  searchCompName.value = ''
  searchProjectName.value = ''
  assocPage.value = 1
  loadAssoc()
}

function openRankingDialog() {
  rankingDialogVisible.value = true
  rankingPage.value = 1
  loadRanking()
}

async function loadRanking() {
  rankingLoading.value = true
  try {
    const res = await coreDependencyApi.getRanking({
      page: rankingPage.value,
      pageSize: rankingPageSize.value,
      orgName: rankingSearchOrg.value,
      compName: rankingSearchComp.value,
    })
    rankingList.value = res.list || []
    rankingTotal.value = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载排行榜失败')
  } finally {
    rankingLoading.value = false
  }
}

function getUsageColor(rate: number): string {
  if (rate >= 70) return '#F56C6C'
  if (rate >= 40) return '#E6A23C'
  return '#67C23A'
}

// ===== 版本碎片化 =====
function handleTabChange(tabName: string) {
  if (tabName === 'fragmentation' && vfList.value.length === 0) {
    loadFragmentation()
  }
  if (tabName === 'highRisk' && hrList.value.length === 0) {
    loadHighRisk()
  }
  if (tabName === 'license' && licList.value.length === 0) {
    loadLicense()
  }
}

async function loadFragmentation() {
  vfLoading.value = true
  try {
    const res = await coreDependencyApi.getVersionFragmentation({
      page: vfPage.value,
      pageSize: vfPageSize.value,
      compName: vfSearchCompName.value,
    })
    vfList.value = res.list || []
    vfTotal.value = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载版本碎片化数据失败')
  } finally {
    vfLoading.value = false
  }
}

function resetVfSearch() {
  vfSearchCompName.value = ''
  vfPage.value = 1
  loadFragmentation()
}

function getFragmentTagType(level: string): string {
  if (level === 'high') return 'danger'
  if (level === 'medium') return 'warning'
  return 'success'
}

function getFragmentLabel(level: string): string {
  if (level === 'high') return '高'
  if (level === 'medium') return '中'
  return '低'
}

function getVersionCountType(count: number): string {
  if (count > 5) return 'danger'
  if (count > 3) return 'warning'
  return 'success'
}

function parseVersionList(versionList: string, latestVersion: string): string[] {
  if (!versionList) return []
  return versionList.split(',').map(v => v.trim()).filter(v => v)
}

// ===== 高风险项目 =====
async function loadHighRisk() {
  hrLoading.value = true
  try {
    const res = await coreDependencyApi.getHighRiskProjects({
      page: hrPage.value,
      pageSize: hrPageSize.value,
      projectName: hrSearchProjectName.value,
    })
    hrList.value = res.list || []
    hrTotal.value = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载高风险项目数据失败')
  } finally {
    hrLoading.value = false
  }
}

function resetHrSearch() {
  hrSearchProjectName.value = ''
  hrPage.value = 1
  loadHighRisk()
}

function getBlacklistTagType(count: number): string {
  if (count >= 5) return 'danger'
  if (count >= 3) return 'warning'
  return 'info'
}

async function handleHrExpand(row: any, expandedRows: any[]) {
  if (expandedRows.length === 0) return
  if (hrDetailMap[row.projectName]) return // 已加载过
  hrDetailMap[row.projectName] = [] // 先占位，标记为加载中
  try {
    const res = await coreDependencyApi.getBlacklistDetail(row.projectName)
    hrDetailMap[row.projectName] = res.list || []
  } catch (e: any) {
    hrDetailMap[row.projectName] = []
  }
}

// ===== 许可证关联 =====
async function loadLicense() {
  licLoading.value = true
  try {
    const res = await coreDependencyApi.getLicenseOverview({
      page: licPage.value,
      pageSize: licPageSize.value,
      licenseName: licSearchName.value,
    })
    licList.value = res.list || []
    licTotal.value = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载许可证数据失败')
  } finally {
    licLoading.value = false
  }
}

function resetLicSearch() {
  licSearchName.value = ''
  licPage.value = 1
  loadLicense()
}

function getLicenseCountTagType(count: number): string {
  if (count >= 20) return 'danger'
  if (count >= 10) return 'warning'
  return 'success'
}

function getLicenseRiskTagType(level: string): string {
  switch (level) {
    case 'high': return 'danger'
    case 'medium': return 'warning'
    case 'low': return 'success'
    default: return 'info'
  }
}

function getLicenseRiskLabel(level: string): string {
  switch (level) {
    case 'high': return '强传染'
    case 'medium': return '弱传染'
    case 'low': return '宽松'
    default: return '未知'
  }
}
</script>

<style scoped>
.core-dependencies-page { padding: 4px; }
.expand-detail { padding: 12px 20px; }
.metric-row { margin-bottom: 16px; }
.metric-card .metric-content { display: flex; flex-direction: column; gap: 8px; }
.metric-label { font-size: 13px; color: #909399; }
.metric-value { font-size: 28px; font-weight: 700; }
.metric-value.primary { color: #409EFF; }
.metric-value.danger { color: #F56C6C; }
.metric-value.warning { color: #E6A23C; }
.metric-value.success { color: #67C23A; }
.tab-toolbar { display: flex; justify-content: space-between; align-items: center; }
.search-area { display: flex; align-items: center; }
.ranking-toolbar { display: flex; align-items: center; }

/* 版本列表标签 */
.version-tags { display: flex; flex-wrap: wrap; gap: 4px; }
.version-tag { margin: 0; }
.expand-detail { padding: 12px 20px; }
</style>
