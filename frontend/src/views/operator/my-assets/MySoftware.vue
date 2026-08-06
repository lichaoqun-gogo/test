<template>
  <div class="my-software">
    <div class="page-header">
      <h2>治理清单</h2>
      <p class="page-desc">分配给您负责治理的基础软件和组件列表</p>
    </div>

    <div class="filter-bar">
      <el-select v-model="filterGrade" placeholder="按等级筛选" clearable style="width: 150px">
        <el-option v-for="g in ['A', 'B', 'C', 'D']" :key="g" :label="`${g}级`" :value="g" />
      </el-select>
      <el-select v-model="filterType" placeholder="按类型筛选" clearable style="width: 150px; margin-left: 12px">
        <el-option label="中间件" value="middleware" />
        <el-option label="数据库" value="database" />
        <el-option label="工具" value="tool" />
        <el-option label="依赖包" value="dependency" />
        <el-option label="组件库" value="library" />
      </el-select>
      <el-input
        v-model="filterKeyword"
        placeholder="搜索软件名称"
        clearable
        style="width: 200px; margin-left: 12px"
      />
    </div>

    <el-table :data="filteredSoftware" stripe style="margin-top: 16px">
      <el-table-column prop="name" label="软件/组件名称" min-width="140">
        <template #default="{ row }">
          <el-link type="primary" @click="viewRiskDetail(row)">{{ row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="version" label="版本" width="100" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          {{ SOFTWARE_TYPE_LABELS[row.type as keyof typeof SOFTWARE_TYPE_LABELS] || row.type }}
        </template>
      </el-table-column>
      <el-table-column label="当前等级" width="120">
        <template #default="{ row }">
          <GradeTag :grade="row.grade" size="small" />
        </template>
      </el-table-column>
      <el-table-column label="风险统计" width="140">
        <template #default="{ row }">
          <span class="risk-stat">
            缺陷 <el-tag :type="row.bugCount > 0 ? 'warning' : 'info'" size="small">{{ row.bugCount }}</el-tag>
            漏洞 <el-tag :type="row.vulnerabilityCount > 0 ? 'danger' : 'info'" size="small">{{ row.vulnerabilityCount }}</el-tag>
          </span>
        </template>
      </el-table-column>
      <el-table-column label="关联缺陷" min-width="180">
        <template #default="{ row }">
          <div v-if="row.defects && row.defects.length > 0" class="defect-list">
            <el-tag
              v-for="d in row.defects.slice(0, 2)"
              :key="d.id"
              :type="defectStatusType(d.status)"
              size="small"
              style="margin-right: 4px; margin-bottom: 4px"
            >
              {{ d.title }}
            </el-tag>
            <span v-if="row.defects.length > 2" class="text-muted">+{{ row.defects.length - 2 }}</span>
          </div>
          <span v-else class="text-muted">暂无缺陷</span>
        </template>
      </el-table-column>
      <el-table-column prop="lastUpdated" label="最后更新" width="180">
        <template #default="{ row }">{{ formatTime(row.lastUpdated) }}</template>
      </el-table-column>
    </el-table>

    <!-- 风险信息弹窗 -->
    <el-dialog v-model="riskDialogVisible" :title="`${selectedSoftware?.name} - 风险信息`" width="800px">
      <div v-if="selectedSoftware" class="risk-detail">
        <el-descriptions :column="3" border size="small" style="margin-bottom: 20px">
          <el-descriptions-item label="软件名称">{{ selectedSoftware.name }}</el-descriptions-item>
          <el-descriptions-item label="版本">{{ selectedSoftware.version }}</el-descriptions-item>
          <el-descriptions-item label="类型">
            {{ SOFTWARE_TYPE_LABELS[selectedSoftware.type as keyof typeof SOFTWARE_TYPE_LABELS] }}
          </el-descriptions-item>
          <el-descriptions-item label="等级">
            <GradeTag :grade="selectedSoftware.grade" />
          </el-descriptions-item>
          <el-descriptions-item label="缺陷数">
            <el-tag :type="selectedSoftware.bugCount > 0 ? 'warning' : 'info'" size="small">
              {{ selectedSoftware.bugCount }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="漏洞数">
            <el-tag :type="selectedSoftware.vulnerabilityCount > 0 ? 'danger' : 'info'" size="small">
              {{ selectedSoftware.vulnerabilityCount }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <h4>缺陷/漏洞列表</h4>
        <el-table :data="allRisks" stripe size="small">
          <el-table-column prop="id" label="编号" width="120" />
          <el-table-column prop="title" label="标题" min-width="200" />
          <el-table-column label="类型" width="80">
            <template #default="{ row }">
              <el-tag :type="row.riskType === '缺陷' ? 'warning' : 'danger'" size="small">
                {{ row.riskType }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="defectStatusType(row.status)" size="small">
                {{ statusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="severity" label="严重程度" width="100">
            <template #default="{ row }">
              <el-tag :type="severityType(row.severity)" size="small">{{ row.severity }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="foundDate" label="发现日期" width="180">
            <template #default="{ row }">{{ formatTime(row.foundDate) }}</template>
          </el-table-column>
        </el-table>

        <el-empty v-if="allRisks.length === 0" description="暂无缺陷和漏洞记录" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { SOFTWARE_TYPE_LABELS } from '@/types'
import { useUserStore } from '@/stores/user'
import { formatTime } from '@/utils/format'

const route = useRoute()
const userStore = useUserStore()

const filterGrade = ref('')
const filterType = ref('')
const filterKeyword = ref('')

const riskDialogVisible = ref(false)
const selectedSoftware = ref<any>(null)

interface RiskItem {
  id: string
  title: string
  riskType: '缺陷' | '漏洞'
  status: string
  severity: string
  foundDate: string
}

// 仅显示分配给当前专员的软件
const softwareList = ref<any[]>([])
const mySoftware = computed(() => {
  const currentUserName = userStore.userInfo?.nickname || userStore.userInfo?.username
  if (!currentUserName) return softwareList.value
  return softwareList.value.filter(s => s.assignedSpecialist === currentUserName)
})

const filteredSoftware = computed(() => {
  let list = mySoftware.value
  if (filterGrade.value) {
    list = list.filter(s => s.grade === filterGrade.value)
  }
  if (filterType.value) {
    list = list.filter(s => s.type === filterType.value)
  }
  if (filterKeyword.value) {
    list = list.filter(s => s.name.toLowerCase().includes(filterKeyword.value.toLowerCase()))
  }
  return list
})

// 获取所有风险（缺陷+漏洞）
const allRisks = computed<RiskItem[]>(() => {
  if (!selectedSoftware.value) return []
  const software = selectedSoftware.value
  const risks: RiskItem[] = []

  // 添加缺陷
  if (software.defects) {
    software.defects.forEach((d: any) => {
      risks.push({
        id: d.id,
        title: d.title,
        riskType: '缺陷',
        status: d.status,
        severity: d.severity || '中',
        foundDate: d.foundDate || '-',
      })
    })
  }

  // 添加漏洞（mock数据中可能没有单独的漏洞列表，这里从vulnerabilityCount推断）
  // 如果有单独的vulnerabilities字段，可以添加
  return risks
})

function defectStatusType(status: string): 'danger' | 'warning' | 'success' | 'info' {
  const types: Record<string, 'danger' | 'warning' | 'success' | 'info'> = { open: 'danger', fixing: 'warning', fixed: 'success', accepted: 'info' }
  return types[status] || 'info'
}

function statusLabel(status: string): string {
  const labels: Record<string, string> = { open: '未处理', fixing: '修复中', fixed: '已修复', accepted: '已接受' }
  return labels[status] || status
}

function severityType(severity: string): 'danger' | 'warning' | 'info' {
  const types: Record<string, 'danger' | 'warning' | 'info'> = { 高: 'danger', 中: 'warning', 低: 'info' }
  return types[severity] || 'info'
}

function viewRiskDetail(row: any) {
  selectedSoftware.value = row
  riskDialogVisible.value = true
}

function loadSoftwareAssets() {
  // TODO: 调用后端 API 获取软件资产列表
}

onMounted(() => {
  // 非专员路径访问时提示
  if (!route.path.startsWith('/specialist')) {
    ElMessage.warning('此页面仅开源治理专员可见')
  }
  loadSoftwareAssets()
})
</script>

<style scoped>
.page-header {
  margin-bottom: 16px;
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

.filter-bar {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.risk-stat {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.defect-list {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}

.text-muted {
  color: #C0C4CC;
}

.risk-detail h4 {
  margin: 16px 0 12px;
  font-size: 15px;
  font-weight: 600;
  color: #1f2328;
}
</style>
