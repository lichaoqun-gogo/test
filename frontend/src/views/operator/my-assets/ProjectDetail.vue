<template>
  <div class="project-detail">
    <div class="page-header">
      <el-page-header @back="goBack" :title="'返回'">
        <template #content>
          <span class="page-title">{{ project?.name }}</span>
        </template>
      </el-page-header>
    </div>

    <el-tabs v-model="activeTab" class="detail-tabs">
      <!-- 项目许可证 -->
      <el-tab-pane label="项目许可证" name="licenses">
        <el-table :data="licenses" stripe>
          <el-table-column prop="name" label="许可证名称" width="200" />
          <el-table-column label="风险等级" width="120">
            <template #default="{ row }">
              <el-tag :type="row.riskLevel === 'low' ? 'success' : 'warning'" size="small">
                {{ row.riskLevel === 'low' ? '低风险' : '高风险' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="componentCount" label="涉及组件数" />
        </el-table>
      </el-tab-pane>

      <!-- 项目依赖清单 -->
      <el-tab-pane label="项目依赖清单" name="dependencies">
        <el-table :data="dependencies" stripe>
          <el-table-column prop="name" label="依赖名称" width="150" />
          <el-table-column prop="version" label="版本" width="120" />
          <el-table-column prop="license" label="许可证" width="140" />
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

      <!-- 项目缺陷/漏洞 -->
      <el-tab-pane label="项目缺陷/漏洞" name="vulnerabilities">
        <el-table :data="vulnerabilities" stripe>
          <el-table-column prop="cveId" label="CVE编号" width="160" />
          <el-table-column label="严重程度" width="100">
            <template #default="{ row }">
              <el-tag :type="row.severity === 'critical' || row.severity === 'high' ? 'danger' : 'warning'" size="small">
                {{ row.severity }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="component" label="涉及组件" width="140" />
          <el-table-column prop="version" label="版本" width="100" />
          <el-table-column prop="description" label="描述" min-width="200" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'open' ? 'danger' : 'success'" size="small">
                {{ row.status === 'open' ? '未修复' : '已修复' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="discoveredDate" label="发现日期" width="180">
            <template #default="{ row }">{{ formatTime(row.discoveredDate) }}</template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- 历史阻断记录 -->
      <el-tab-pane label="历史阻断记录" name="blocks">
        <el-table :data="blockRecords" stripe>
          <el-table-column prop="date" label="阻断时间" width="180">
            <template #default="{ row }">{{ formatTime(row.date) }}</template>
          </el-table-column>
          <el-table-column prop="reason" label="阻断原因" min-width="200" />
          <el-table-column prop="rule" label="阻断规则" width="150" />
          <el-table-column prop="handler" label="处理人" width="100" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag v-if="row.status === 'resolved'" type="success" size="small">已解决</el-tag>
              <el-tag v-else type="warning" size="small">待处理</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { formatTime } from '@/utils/format'

const router = useRouter()
const route = useRoute()
const activeTab = ref('licenses')

const projectId = computed(() => route.params.id as string)
const project = ref<any>(null)

const licenses = ref<any[]>([])
const dependencies = ref<any[]>([])
const vulnerabilities = ref<any[]>([])
const blockRecords = ref<any[]>([])

function loadProjectDetail(id: string) {
  // TODO: 调用后端 API 获取项目详情
}

onMounted(() => {
  loadProjectDetail(projectId.value)
})

function goBack() {
  router.back()
}
</script>

<style scoped>
.project-detail {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
}

.detail-tabs {
  background: #fff;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
}
</style>
