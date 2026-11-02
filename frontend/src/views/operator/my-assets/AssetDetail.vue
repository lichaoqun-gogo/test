<template>
  <div class="asset-detail">
    <el-page-header @back="goBack" :title="pageTitle" />

    <el-tabs v-model="activeTab" style="margin-top: 16px">
      <!-- 软件详情 -->
      <el-tab-pane v-if="assetDetail" label="基本信息" name="info">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="软件名称">{{ assetDetail.name }}</el-descriptions-item>
          <el-descriptions-item label="版本">{{ assetDetail.version }}</el-descriptions-item>
          <el-descriptions-item label="类型">{{ SOFTWARE_TYPE_LABELS[assetDetail.type as keyof typeof SOFTWARE_TYPE_LABELS] }}</el-descriptions-item>
          <el-descriptions-item label="等级">
            <GradeTag :grade="assetDetail.grade" />
          </el-descriptions-item>
          <el-descriptions-item label="开发语言">{{ assetDetail.language }}</el-descriptions-item>
          <el-descriptions-item label="组织">{{ assetDetail.organization }}</el-descriptions-item>
          <el-descriptions-item label="负责人">{{ assetDetail.owner }}</el-descriptions-item>
          <el-descriptions-item label="部门">{{ assetDetail.department }}</el-descriptions-item>
          <el-descriptions-item label="最后更新">{{ formatTime(assetDetail.lastUpdated) }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="assetDetail.status === 'active' ? 'success' : 'info'">
              {{ assetDetail.status === 'active' ? '使用中' : '已停用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">{{ assetDetail.description }}</el-descriptions-item>
        </el-descriptions>

        <el-row :gutter="16" style="margin-top: 16px">
          <el-col :span="8">
            <el-statistic title="缺陷数" :value="assetDetail.bugCount" />
          </el-col>
          <el-col :span="8">
            <el-statistic title="漏洞数" :value="assetDetail.vulnerabilityCount">
              <template #suffix>
                <el-tag v-if="assetDetail.vulnerabilityCount > 0" type="danger" size="small" style="margin-left: 8px">
                  需关注
                </el-tag>
              </template>
            </el-statistic>
          </el-col>
        </el-row>
      </el-tab-pane>

      <!-- 项目详情 -->
      <template v-if="projectDetail">
        <el-tab-pane label="项目许可证" name="licenses">
          <el-table :data="licenses" stripe>
            <el-table-column prop="name" label="许可证名称" width="200" />
            <el-table-column label="风险等级" width="120">
              <template #default="{ row }">
                <RiskTag :level="row.riskLevel" />
              </template>
            </el-table-column>
            <el-table-column prop="componentCount" label="涉及组件数" />
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="项目依赖清单" name="dependencies">
          <el-alert
            v-if="blacklistDeps.length > 0"
            :title="`发现 ${blacklistDeps.length} 个黑名单依赖`"
            type="error"
            show-icon
            :closable="false"
            style="margin-bottom: 16px"
          />
          <el-table :data="dependencies" stripe>
            <el-table-column prop="name" label="依赖名称" width="150" />
            <el-table-column prop="version" label="版本" width="120" />
            <el-table-column prop="license" label="许可证" width="140" />
            <el-table-column label="黑名单" width="90">
              <template #default="{ row }">
                <el-tag v-if="row.isInBlacklist" type="danger" size="small">是</el-tag>
                <el-tag v-else type="info" size="small">否</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="漏洞数" width="90">
              <template #default="{ row }">
                <el-tag v-if="row.vulnerabilityCount > 0" type="danger" size="small">
                  {{ row.vulnerabilityCount }}
                </el-tag>
                <span v-else>0</span>
              </template>
            </el-table-column>
            <el-table-column prop="lastUpdated" label="最后更新">
              <template #default="{ row }">{{ formatTime(row.lastUpdated) }}</template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="项目缺陷/漏洞" name="vulnerabilities">
          <el-alert
            v-if="criticalVulns.length > 0"
            :title="`发现 ${criticalVulns.length} 个高危及以上漏洞`"
            type="error"
            show-icon
            :closable="false"
            style="margin-bottom: 16px"
          />
          <el-table :data="vulnerabilities" stripe>
            <el-table-column prop="cveId" label="CVE编号" width="160" />
            <el-table-column label="严重程度" width="100">
              <template #default="{ row }">
                <RiskTag :level="row.severity" />
              </template>
            </el-table-column>
            <el-table-column prop="component" label="涉及组件" width="140" />
            <el-table-column prop="version" label="版本" width="100" />
            <el-table-column prop="description" label="描述" min-width="200" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="vulnStatusType(row.status)" size="small">
                  {{ vulnStatusLabel(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="discoveredDate" label="发现日期" width="180">
              <template #default="{ row }">{{ formatTime(row.discoveredDate) }}</template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </template>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { SOFTWARE_TYPE_LABELS } from '@/types'
import { formatTime } from '@/utils/format'

const route = useRoute()
const router = useRouter()

const activeTab = ref('info')

const assetDetail = ref<any>(null)
const projectDetail = ref<any>(null)

const pageTitle = computed(() =>
  assetDetail.value ? assetDetail.value.name : projectDetail.value ? projectDetail.value.name : '详情'
)

const licenses = ref<any[]>([])
const dependencies = ref<any[]>([])

const blacklistDeps = computed(() =>
  dependencies.value.filter(d => d.isInBlacklist)
)

const vulnerabilities = ref<any[]>([])

const criticalVulns = computed(() =>
  vulnerabilities.value.filter(v => v.severity === 'critical' || v.severity === 'high')
)

function loadAssetDetail(id: string) {
  // TODO: 调用后端 API 获取资产/项目详情
}

onMounted(() => {
  loadAssetDetail(route.params.id as string)
})

function goBack() {
  router.back()
}

function vulnStatusType(status: string) {
  const types: Record<string, string> = { open: 'danger', fixing: 'warning', fixed: 'success', accepted: 'info' }
  return types[status] || ''
}

function vulnStatusLabel(status: string) {
  const labels: Record<string, string> = { open: '待修复', fixing: '修复中', fixed: '已修复', accepted: '已接受' }
  return labels[status] || status
}
</script>
