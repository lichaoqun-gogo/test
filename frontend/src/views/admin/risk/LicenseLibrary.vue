<template>
  <div class="license-library">
    <div class="page-header">
      <div class="header-left">
        <h3>许可证信息库</h3>
      </div>
    </div>

    <!-- 看板 -->
    <el-row :gutter="16" style="margin-top: 16px; margin-bottom: 16px">
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="许可证总数" :value="licenses.length" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="OSI已认证" :value="licenses.filter(l => l.IS_OSIAPPROVED === 1).length" value-style="color: #67C23A" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="FSF Libre" :value="licenses.filter(l => l.FSFLIBRE === 1).length" value-style="color: #409EFF" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="已过期" :value="licenses.filter(l => l.IS_DEPRECATED === 1).length" value-style="color: #F56C6C" />
        </el-card>
      </el-col>
    </el-row>

    <!-- Tab页签 -->
    <el-tabs v-model="activeTab" type="border-card" style="margin-top: 16px">
      <!-- Tab 1: 许可证分类台账 -->
      <el-tab-pane label="许可证分类" name="category">
        <div class="filter-bar">
          <el-input v-model="categorySearch" placeholder="搜索许可证分类" clearable style="width: 260px" />
        </div>

        <el-table :data="filteredCategories" stripe style="margin-top: 12px">
          <el-table-column prop="name" label="许可证分类" min-width="180" />
          <el-table-column prop="riskWeight" label="风险权重" width="120">
            <template #default="{ row }">
              <el-tag :type="getRiskLevelTag(getRiskLevel(row.riskWeight)) as any" size="small">
                {{ row.riskWeight }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="licenseCount" label="许可证数量" width="120" />
          <el-table-column prop="createTime" label="创建时间" width="180">
            <template #default="{ row }">{{ formatTime(row.createTime) }}</template>
          </el-table-column>
          <el-table-column prop="updateTime" label="更新时间" width="180">
            <template #default="{ row }">{{ formatTime(row.updateTime) }}</template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- Tab 2: 许可证列表台账 -->
      <el-tab-pane label="许可证列表" name="list">
        <div class="filter-bar">
          <el-input v-model="searchText" placeholder="搜索许可证代号或名称" clearable style="width: 260px" />
          <el-select v-model="filterOsiApproved" placeholder="OSI认证" clearable style="width: 120px">
            <el-option label="是" :value="1" />
            <el-option label="否" :value="0" />
          </el-select>
          <el-select v-model="filterDeprecated" placeholder="是否过期" clearable style="width: 120px">
            <el-option label="是" :value="1" />
            <el-option label="否" :value="0" />
          </el-select>
          <el-select v-model="filterCategory" placeholder="许可证分类" clearable style="width: 180px">
            <el-option v-for="cat in licenseCategories" :key="cat.id" :label="cat.name" :value="cat.name" />
          </el-select>
          <el-button @click="resetFilters">
            <el-icon><Refresh /></el-icon> 重置
          </el-button>
        </div>

        <el-table :data="filteredLicenses" stripe style="margin-top: 12px">
          <el-table-column prop="_name" label="许可证名称" min-width="200">
            <template #default="{ row }">
              <el-link type="primary" @click="openDetail(row as any)">{{ row._name }}</el-link>
            </template>
          </el-table-column>
          <el-table-column prop="LICENSE" label="许可证代号" width="150" />
          <el-table-column label="FSF Libre" width="100">
            <template #default="{ row }">
              <el-tag :type="row.FSFLIBRE === 1 ? 'success' : 'info'" size="small">
                {{ row.FSFLIBRE === 1 ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="OSI审核" width="100">
            <template #default="{ row }">
              <el-tag :type="row.IS_OSIAPPROVED === 1 ? 'success' : 'info'" size="small">
                {{ row.IS_OSIAPPROVED === 1 ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="是否过期" width="100">
            <template #default="{ row }">
              <el-tag :type="row.IS_DEPRECATED === 1 ? 'danger' : 'success'" size="small">
                {{ row.IS_DEPRECATED === 1 ? '已过期' : '有效' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="RISKWEIGHT" label="风险分值" width="100" sortable />
          <el-table-column prop="LICENSEGROUP" label="许可证分类" width="150">
            <template #default="{ row }">
              <el-tag :type="getRiskLevelTag(getRiskLevel(row.RISKWEIGHT)) as any" size="small">
                {{ row.LICENSEGROUP || '-' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <!-- 许可证详情弹框 -->
    <el-dialog v-model="detailVisible" :title="currentLicense?._name || '许可证详情'" width="800px">
      <el-tabs v-model="detailActiveTab" v-loading="detailLoading">
        <!-- Tab1: 概览信息 -->
        <el-tab-pane label="概览信息" name="overview">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="许可证名称">{{ currentLicense?._name }}</el-descriptions-item>
            <el-descriptions-item label="许可证代号">{{ currentLicense?.LICENSE }}</el-descriptions-item>
            <el-descriptions-item label="是否过期">
              <el-tag :type="currentLicense?.IS_DEPRECATED === 1 ? 'danger' : 'success'" size="small">
                {{ currentLicense?.IS_DEPRECATED === 1 ? '已过期' : '有效' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="FSF Libre">
              <el-tag :type="currentLicense?.FSFLIBRE === 1 ? 'success' : 'info'" size="small">
                {{ currentLicense?.FSFLIBRE === 1 ? '是' : '否' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="OSI审核">
              <el-tag :type="currentLicense?.IS_OSIAPPROVED === 1 ? 'success' : 'info'" size="small">
                {{ currentLicense?.IS_OSIAPPROVED === 1 ? '是' : '否' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="许可证分类">
              {{ currentLicense?.license_group || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="说明" :span="2">{{ currentLicense?._COMMENT || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>

        <!-- Tab2: 许可证说明 -->
        <el-tab-pane label="许可证说明" name="licenseText">
          <div class="license-content">
            <pre>{{ currentLicense?._text || '暂无许可证说明' }}</pre>
          </div>
        </el-tab-pane>

        <!-- Tab3: 许可证模板 -->
        <el-tab-pane label="许可证模板" name="licenseTemplate">
          <div class="license-content">
            <pre>{{ currentLicense?.TEMPLATE || '暂无许可证模板' }}</pre>
          </div>
        </el-tab-pane>

        <!-- Tab4: 文件头说明 -->
        <el-tab-pane label="文件头说明" name="fileHeader">
          <div class="license-content">
            <pre>{{ currentLicense?.HEADER || '暂无文件头说明' }}</pre>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { formatTime } from '@/utils/format'
import api from '@/api'

interface LicenseItem {
  id: number
  LICENSE: string
  _name: string
  IS_OSIAPPROVED: number
  IS_DEPRECATED: number
  FSFLIBRE: number
  _COMMENT?: string
  _text?: string
  TEMPLATE?: string
  HEADER?: string
  LICENSEGROUP?: string
  RISKWEIGHT?: number
  license_group?: string
  risk_level?: string
}

const activeTab = ref('category')
const searchText = ref('')
const categorySearch = ref('')
const filterOsiApproved = ref<number | ''>('')
const filterDeprecated = ref<number | ''>('')
const filterCategory = ref('')

const licenses = ref<LicenseItem[]>([])
const loading = ref(false)

// 详情弹框
const detailVisible = ref(false)
const detailActiveTab = ref('overview')
const currentLicense = ref<LicenseItem | null>(null)
const detailLoading = ref(false)

async function loadLicenses() {
  loading.value = true
  try {
    const res: any = await api.get('/risk/licenses', {
      params: { page: 1, pageSize: 1000 }
    })
    const list = Array.isArray(res) ? res : (res.list || [])
    licenses.value = list.map((item: any) => ({
      id: item.ID || item.id,
      LICENSE: item.LICENSE,
      _name: item._NAME || item._name,
      IS_OSIAPPROVED: item.IS_OSIAPPROVED,
      IS_DEPRECATED: item.IS_DEPRECATED,
      FSFLIBRE: item.FSFLIBRE,
      LICENSEGROUP: item.LICENSEGROUP,
      RISKWEIGHT: item.RISKWEIGHT,
      _COMMENT: item._COMMENT,
      _text: item._text,
      TEMPLATE: item.TEMPLATE,
      HEADER: item.HEADER,
    }))
  } catch (error) {
    console.error('加载许可证失败:', error)
  } finally {
    loading.value = false
  }
}

async function openDetail(license: LicenseItem) {
  currentLicense.value = null
  detailActiveTab.value = 'overview'
  detailVisible.value = true
  detailLoading.value = true

  try {
    const res: any = await api.get(`/risk/licenses/${license.id}`)
    currentLicense.value = {
      id: res.id,
      LICENSE: res.LICENSE,
      _name: res._name,
      IS_OSIAPPROVED: res.IS_OSIAPPROVED,
      IS_DEPRECATED: res.IS_DEPRECATED,
      FSFLIBRE: res.FSFLIBRE,
      _COMMENT: res._COMMENT,
      _text: res._text,
      TEMPLATE: res.TEMPLATE,
      HEADER: res.HEADER,
      license_group: res.license_group,
    }
  } catch (error) {
    console.error('获取详情失败:', error)
    ElMessage.error('获取许可证详情失败')
  } finally {
    detailLoading.value = false
  }
}

// 许可证分类数据
const licenseCategories = ref<any[]>([])

async function loadLicenseCategories() {
  try {
    const res: any = await api.get('/risk/license-groups')
    licenseCategories.value = res || []
  } catch (error) {
    console.error('加载许可证分类失败:', error)
    licenseCategories.value = []
  }
}

const filteredCategories = computed(() => {
  if (!categorySearch.value) return licenseCategories.value
  return licenseCategories.value.filter(cat =>
    cat.name.toLowerCase().includes(categorySearch.value.toLowerCase())
  )
})

const filteredLicenses = computed(() => {
  return licenses.value.filter(item => {
    const search = searchText.value?.toLowerCase() || ''
    const matchSearch = !search || (item.LICENSE || '').toLowerCase().includes(search) || (item._name || '').toLowerCase().includes(search)
    const matchOsi = filterOsiApproved.value === '' || filterOsiApproved.value === undefined || filterOsiApproved.value === null || item.IS_OSIAPPROVED === filterOsiApproved.value
    const matchDeprecated = filterDeprecated.value === '' || filterDeprecated.value === undefined || filterDeprecated.value === null || item.IS_DEPRECATED === filterDeprecated.value
    const matchCategory = !filterCategory.value || item.LICENSEGROUP === filterCategory.value
    return matchSearch && matchOsi && matchDeprecated && matchCategory
  })
})

function getRiskLevel(riskWeight?: number): string {
  if (riskWeight == null) return ''
  if (riskWeight <= 20) return '低风险'
  if (riskWeight <= 50) return '中风险'
  return '高风险'
}

function getRiskLevelTag(riskLevel: string): string {
  if (!riskLevel) return ''
  if (riskLevel === '低风险') return 'success'
  if (riskLevel === '中风险') return 'warning'
  if (riskLevel === '高风险') return 'danger'
  return 'info'
}

function resetFilters() {
  searchText.value = ''
  filterOsiApproved.value = ''
  filterDeprecated.value = ''
  filterCategory.value = ''
}

onMounted(() => {
  loadLicenses()
  loadLicenseCategories()
})
</script>

<style scoped>
.license-library {
  padding: 16px;
}

.page-header {
  margin-bottom: 8px;
}

.page-header h3 {
  margin: 0;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.license-content {
  max-height: 500px;
  overflow-y: auto;
  padding: 16px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.license-content pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  margin: 0;
  font-size: 13px;
  line-height: 1.6;
  color: #303133;
}
</style>
