<template>
  <div class="software-service">
    <div class="page-header">
      <h2 class="page-title">软件服务</h2>
      <p class="page-desc">支持软件检索、智能推荐、可信包下载等功能</p>
    </div>

    <el-tabs v-model="activeTab" type="border-card" class="service-tabs">
      <!-- ==================== Tab 1: 软件检索与推荐 ==================== -->
      <el-tab-pane label="软件检索与推荐" name="search">
        <!-- 搜索区域 -->
        <div class="search-section">
          <el-input
            v-model="searchQuery"
            placeholder="输入组件名称或版本进行搜索，如：Apache Pulsar、3.1.2"
            size="large"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prepend>
              <el-icon><Search /></el-icon>
            </template>
            <template #append>
              <el-button type="primary" @click="handleSearch">搜索</el-button>
            </template>
          </el-input>

          <div class="search-tips">
            <el-text size="small" type="info">
              提示：搜索将返回组件的活跃度、漏洞、许可证风险等安全信息
            </el-text>
          </div>
        </div>

        <!-- 全局风险警告 -->
        <el-alert
          v-if="hasHighRisk"
          title="搜索结果中存在高风险组件，建议优先考虑低风险替代方案"
          type="warning"
          show-icon
          :closable="false"
          class="global-risk-alert"
        >
          <template #default>
            <p class="alert-detail">以下组件存在高危/严重漏洞或许可证风险，引入前需经过严格安全评估</p>
          </template>
        </el-alert>

        <!-- 搜索结果列表 -->
        <div v-if="searchResults.length > 0" class="search-results">
          <el-card
            v-for="item in searchResults"
            :key="item.id"
            class="result-card"
            shadow="hover"
            :class="{ 'risk-card': item.hasRiskWarning }"
          >
            <template #header>
              <div class="result-header">
                <div class="header-left">
                  <span class="result-name">{{ item.name }}</span>
                  <el-tag size="small" type="info">{{ item.version }}</el-tag>
                  <el-tag
                    :type="item.introduced ? 'success' : 'info'"
                    size="small"
                    effect="plain"
                  >
                    {{ item.introduced ? '已引入' : '未引入' }}
                  </el-tag>
                  <el-tag
                    v-if="item.hasRiskWarning"
                    type="danger"
                    size="small"
                    effect="dark"
                  >
                    ⚠️ 存在风险
                  </el-tag>
                </div>
                <el-button
                  v-if="!item.introduced"
                  type="primary"
                  size="small"
                  @click="applyIntroduce(item)"
                >
                  申请引入
                </el-button>
                <el-button v-else size="small" disabled>已引入</el-button>
              </div>
            </template>

            <!-- 组件描述 -->
            <p class="component-desc">{{ item.description }}</p>

            <!-- 核心指标 -->
            <el-row :gutter="20" class="metrics-row">
              <el-col :span="6">
                <div class="metric-box">
                  <div class="metric-header">
                    <el-icon class="metric-icon"><TrendCharts /></el-icon>
                    <span class="metric-title">活跃度评分</span>
                  </div>
                  <el-progress
                    :percentage="item.activityScore"
                    :color="getActivityColor(item.activityScore)"
                    :stroke-width="10"
                    :format="() => `${item.activityScore}分`"
                  />
                </div>
              </el-col>
              <el-col :span="6">
                <div class="metric-box">
                  <div class="metric-header">
                    <el-icon class="metric-icon"><WarningFilled /></el-icon>
                    <span class="metric-title">已知漏洞</span>
                  </div>
                  <div class="vuln-list">
                    <el-tag
                      v-for="v in item.vulnerabilities"
                      :key="v.severity"
                      :type="vulnTagType(v.severity)"
                      size="small"
                      effect="plain"
                    >
                      {{ RISK_LEVEL_LABELS[v.severity] }}: {{ v.count }}
                    </el-tag>
                    <span v-if="item.vulnerabilityCount === 0" class="no-vuln">
                      <el-icon><CircleCheckFilled /></el-icon>
                      无已知漏洞
                    </span>
                  </div>
                </div>
              </el-col>
              <el-col :span="6">
                <div class="metric-box">
                  <div class="metric-header">
                    <el-icon class="metric-icon"><Document /></el-icon>
                    <span class="metric-title">许可证</span>
                  </div>
                  <div class="license-info">
                    <span class="license-name">{{ item.licenseName }}</span>
                    <RiskTag :level="item.licenseRisk" size="small" />
                  </div>
                </div>
              </el-col>
              <el-col :span="6">
                <div class="metric-box">
                  <div class="metric-header">
                    <el-icon class="metric-icon"><Box /></el-icon>
                    <span class="metric-title">引入状态</span>
                  </div>
                  <div v-if="item.introduced" class="introduced-info">
                    <span>已引入 v{{ item.introducedVersion }}</span>
                    <GradeTag v-if="item.introducedGrade" :grade="item.introducedGrade" size="small" />
                  </div>
                  <span v-else class="not-introduced">未引入</span>
                </div>
              </el-col>
            </el-row>

            <!-- 风险警告详情 -->
            <div v-if="item.riskWarnings && item.riskWarnings.length > 0" class="risk-warnings">
              <el-alert
                v-for="(warning, idx) in item.riskWarnings"
                :key="idx"
                :title="warning"
                type="error"
                :closable="false"
                show-icon
                class="risk-warning-item"
              />
            </div>

            <!-- 智能推荐组件 -->
            <SmartRecommendation
              v-if="item.recommendations && item.recommendations.length > 0"
              :recommendations="item.recommendations"
            />
          </el-card>
        </div>

        <!-- 无结果提示 -->
        <el-empty
          v-else-if="searched"
          description="未找到匹配的组件，请尝试其他关键词"
          :image-size="120"
        >
          <el-button type="primary" @click="resetSearch">重置搜索</el-button>
        </el-empty>

        <!-- 初始状态提示 -->
        <el-empty
          v-else
          description="输入组件名称或版本开始搜索"
          :image-size="100"
        />
      </el-tab-pane>

      <!-- ==================== Tab 2: 可信包下载 ==================== -->
      <el-tab-pane label="可信包下载" name="download">
        <!-- 筛选栏 -->
        <div class="filter-bar">
          <div class="filter-left">
            <el-select
              v-model="packageFilter"
              placeholder="按分类筛选"
              clearable
              style="width: 160px"
              @change="handleFilterChange"
            >
              <el-option label="全部分类" value="" />
              <el-option label="中间件" value="中间件" />
              <el-option label="数据库" value="数据库" />
              <el-option label="依赖包" value="依赖包" />
              <el-option label="工具" value="工具" />
            </el-select>

            <el-select
              v-model="scanStatusFilter"
              placeholder="扫描状态"
              clearable
              style="width: 140px; margin-left: 12px"
              @change="handleFilterChange"
            >
              <el-option label="全部状态" value="" />
              <el-option label="通过" value="passed" />
              <el-option label="警告" value="warning" />
            </el-select>
          </div>

          <div class="filter-right">
            <el-input
              v-model="packageSearch"
              placeholder="搜索包名称..."
              clearable
              style="width: 200px"
              @input="handleFilterChange"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>

        <!-- 统计信息 -->
        <div class="stats-bar">
          <el-tag type="info" effect="plain">
            共 {{ filteredPackages.length }} 个可信包
          </el-tag>
          <el-tag type="success" effect="plain" style="margin-left: 8px">
            扫描通过: {{ scanPassedCount }}
          </el-tag>
          <el-tag type="warning" effect="plain" style="margin-left: 8px">
            存在警告: {{ scanWarningCount }}
          </el-tag>
        </div>

        <!-- 可信包列表 -->
        <el-table :data="paginatedPackages" stripe style="margin-top: 16px">
          <el-table-column prop="name" label="名称" width="140">
            <template #default="{ row }">
              <div class="package-name-cell">
                <el-icon class="package-icon"><Box /></el-icon>
                <span>{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="version" label="版本" width="100" />
          <el-table-column prop="category" label="分类" width="90">
            <template #default="{ row }">
              <el-tag size="small" effect="plain">{{ row.category }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="扫描状态" width="100">
            <template #default="{ row }">
              <el-tag
                :type="row.scanStatus === 'passed' ? 'success' : 'warning'"
                size="small"
              >
                {{ row.scanStatus === 'passed' ? '通过' : '警告' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="scanDate" label="扫描日期" width="120">
            <template #default="{ row }">{{ formatTime(row.scanDate) }}</template>
          </el-table-column>
          <el-table-column prop="size" label="大小" width="80" />
          <el-table-column label="下载次数" width="100" sortable :sort-method="sortByDownload">
            <template #default="{ row }">
              <span>{{ row.downloadCount.toLocaleString() }}</span>
            </template>
          </el-table-column>
          <el-table-column label="发布者" width="120">
            <template #default="{ row }">
              <el-text size="small" truncated>{{ row.publisher }}</el-text>
            </template>
          </el-table-column>
          <el-table-column prop="checksum" label="校验和" min-width="200">
            <template #default="{ row }">
              <el-tooltip :content="row.checksum" placement="top">
                <el-text type="info" size="small" truncated class="checksum-text">
                  {{ row.checksum.substring(0, 20) }}...
                </el-text>
              </el-tooltip>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" size="small" @click="handleDownload(row)">
                <el-icon><Download /></el-icon>
                下载
              </el-button>
              <el-button size="small" @click="showPackageDetail(row)">
                详情
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="currentPage"
            :page-size="pageSize"
            :total="filteredPackages.length"
            layout="total, sizes, prev, pager, next"
            :page-sizes="[5, 10, 20]"
          />
        </div>

        <!-- 包详情弹窗 -->
        <el-dialog
          v-model="detailVisible"
          :title="`包详情 - ${selectedPackage?.name}`"
          width="600px"
        >
          <template v-if="selectedPackage">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="包名称">
                {{ selectedPackage.name }}
              </el-descriptions-item>
              <el-descriptions-item label="版本">
                {{ selectedPackage.version }}
              </el-descriptions-item>
              <el-descriptions-item label="分类">
                <el-tag size="small">{{ selectedPackage.category }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="发布者">
                {{ selectedPackage.publisher }}
              </el-descriptions-item>
              <el-descriptions-item label="扫描状态">
                <el-tag :type="selectedPackage.scanStatus === 'passed' ? 'success' : 'warning'" size="small">
                  {{ selectedPackage.scanStatus === 'passed' ? '通过' : '警告' }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="扫描日期">
                {{ formatTime(selectedPackage.scanDate) }}
              </el-descriptions-item>
              <el-descriptions-item label="文件大小">
                {{ selectedPackage.size }}
              </el-descriptions-item>
              <el-descriptions-item label="下载次数">
                {{ selectedPackage.downloadCount.toLocaleString() }}
              </el-descriptions-item>
              <el-descriptions-item label="描述" :span="2">
                {{ selectedPackage.description }}
              </el-descriptions-item>
              <el-descriptions-item label="SHA256校验和" :span="2">
                <el-text type="info" size="small">{{ selectedPackage.checksum }}</el-text>
              </el-descriptions-item>
            </el-descriptions>
          </template>
          <template #footer>
            <el-button @click="detailVisible = false">关闭</el-button>
            <el-button type="primary" @click="handleDownload(selectedPackage)">
              <el-icon><Download /></el-icon>
              下载此包
            </el-button>
          </template>
        </el-dialog>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { RISK_LEVEL_LABELS } from '@/types'
import { Search, Download, TrendCharts, WarningFilled, Document, Box, CircleCheckFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import SmartRecommendation from './components/SmartRecommendation.vue'
import { formatTime } from '@/utils/format'

const router = useRouter()

// ==================== Tab 1: 搜索相关 ====================
const activeTab = ref('search')
const searchQuery = ref('')
const searchResults = ref<any[]>([])
const searched = ref(false)

const hasHighRisk = computed(() =>
  searchResults.value.some(r =>
    r.vulnerabilities.some((v: any) => v.severity === 'critical' || v.severity === 'high') || r.licenseRisk === 'high'
  )
)

function handleSearch() {
  searched.value = true
  if (!searchQuery.value.trim()) {
    searchResults.value = []
    return
  }
  // TODO: 调用后端搜索 API
  searchResults.value = []
}

function resetSearch() {
  searchQuery.value = ''
  searchResults.value = []
  searched.value = true
}

function getActivityColor(score: number) {
  if (score >= 80) return '#67C23A'
  if (score >= 60) return '#E6A23C'
  return '#F56C6C'
}

function vulnTagType(severity: string) {
  const types: Record<string, '' | 'success' | 'warning' | 'danger'> = {
    critical: 'danger',
    high: 'danger',
    medium: 'warning',
    low: 'success',
  }
  return types[severity] || ''
}

function applyIntroduce(item: any) {
  ElMessage.success(`正在跳转到引入申请页面...`)
  router.push({
    path: '/operator/my-task',
  })
}

// ==================== Tab 2: 可信包相关 ====================
const packageFilter = ref('')
const scanStatusFilter = ref('')
const packageSearch = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const detailVisible = ref(false)
const selectedPackage = ref<any>(null)

const trustedPackages = ref<any[]>([])
const filteredPackages = computed(() => {
  let result = trustedPackages.value

  if (packageFilter.value) {
    result = result.filter(p => p.category === packageFilter.value)
  }
  if (scanStatusFilter.value) {
    result = result.filter(p => p.scanStatus === scanStatusFilter.value)
  }
  if (packageSearch.value.trim()) {
    const q = packageSearch.value.toLowerCase()
    result = result.filter(p => p.name.toLowerCase().includes(q))
  }

  return result
})

const paginatedPackages = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredPackages.value.slice(start, start + pageSize.value)
})

const scanPassedCount = computed(() =>
  filteredPackages.value.filter(p => p.scanStatus === 'passed').length
)

const scanWarningCount = computed(() =>
  filteredPackages.value.filter(p => p.scanStatus === 'warning').length
)

function handleFilterChange() {
  currentPage.value = 1
}

function sortByDownload(a: any, b: any) {
  return a.downloadCount - b.downloadCount
}

function handleDownload(pkg: any) {
  if (!pkg) return
  ElMessage.success(`开始下载 ${pkg.name} v${pkg.version}...`)
}

function showPackageDetail(pkg: any) {
  selectedPackage.value = pkg
  detailVisible.value = true
}

function loadTrustedPackages() {
  // TODO: 调用后端 API 获取可信包列表
}

onMounted(() => {
  loadTrustedPackages()
})
</script>

<style scoped>
.software-service {
  padding: 4px;
}

.page-header {
  margin-bottom: 20px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.page-desc {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
}

.service-tabs {
  min-height: 500px;
}

/* ==================== 搜索区域 ==================== */
.search-section {
  margin-bottom: 20px;
}

.search-tips {
  margin-top: 8px;
}

.global-risk-alert {
  margin-bottom: 16px;
}

.alert-detail {
  margin: 4px 0 0 0;
  font-size: 13px;
  color: #92400e;
}

/* ==================== 搜索结果 ==================== */
.search-results {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.result-card {
  transition: all 0.25s ease;
  border: 1px solid #e5e7eb;
}

.result-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px -4px rgba(0, 0, 0, 0.08);
}

.result-card.risk-card {
  border-color: #fecaca;
}

.result-card.risk-card:hover {
  box-shadow: 0 8px 24px -4px rgba(239, 68, 68, 0.12);
}

.result-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.result-name {
  font-size: 17px;
  font-weight: 600;
  color: #111827;
}

.component-desc {
  font-size: 13px;
  color: #6b7280;
  margin: 0 0 16px 0;
  line-height: 1.6;
}

.metrics-row {
  margin-bottom: 4px;
}

.metric-box {
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
}

.metric-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
}

.metric-icon {
  font-size: 14px;
  color: #6b7280;
}

.metric-title {
  font-size: 12px;
  color: #9ca3af;
  font-weight: 500;
}

.vuln-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  align-items: center;
}

.no-vuln {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #22c55e;
  font-size: 13px;
  font-weight: 500;
}

.no-vuln .el-icon {
  font-size: 16px;
}

.license-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.license-name {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.introduced-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #15803d;
}

.not-introduced {
  font-size: 13px;
  color: #9ca3af;
}

/* ==================== 风险警告 ==================== */
.risk-warnings {
  margin-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.risk-warning-item {
  font-size: 13px;
}

/* ==================== 筛选栏 ==================== */
.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.filter-left,
.filter-right {
  display: flex;
  align-items: center;
}

/* ==================== 统计栏 ==================== */
.stats-bar {
  display: flex;
  align-items: center;
  margin-top: 12px;
}

/* ==================== 分页 ==================== */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

/* ==================== 包名称单元格 ==================== */
.package-name-cell {
  display: flex;
  align-items: center;
  gap: 6px;
}

.package-icon {
  font-size: 16px;
  color: #6b7280;
}

.checksum-text {
  max-width: 100%;
}
</style>
