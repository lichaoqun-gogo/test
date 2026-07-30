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
          <el-table-column prop="lastUpdated" label="最后更新" />
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
          <el-table-column prop="discoveredDate" label="发现日期" width="120" />
        </el-table>
      </el-tab-pane>

      <!-- 历史阻断记录 -->
      <el-tab-pane label="历史阻断记录" name="blocks">
        <el-table :data="blockRecords" stripe>
          <el-table-column prop="date" label="阻断时间" width="160" />
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
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { mockProjects } from '@/mock/assets'

const router = useRouter()
const route = useRoute()
const activeTab = ref('licenses')

const projectId = computed(() => route.params.id as string)
const project = computed(() => mockProjects.find(p => p.id === projectId.value))

// Mock数据
const licenses = ref([
  { name: 'Apache-2.0', riskLevel: 'low', componentCount: 45 },
  { name: 'MIT', riskLevel: 'low', componentCount: 32 },
  { name: 'GPL-3.0', riskLevel: 'high', componentCount: 3 },
])

const dependencies = ref([
  { name: 'Spring Boot', version: '3.2.3', license: 'Apache-2.0', vulnerabilityCount: 0, lastUpdated: '2024-03-15' },
  { name: 'Vue.js', version: '3.4.19', license: 'MIT', vulnerabilityCount: 1, lastUpdated: '2024-03-10' },
  { name: 'Log4j2', version: '2.22.1', license: 'Apache-2.0', vulnerabilityCount: 2, lastUpdated: '2024-03-05' },
])

const vulnerabilities = ref([
  {
    cveId: 'CVE-2024-1234',
    severity: 'high',
    component: 'Log4j2',
    version: '2.22.1',
    description: '远程代码执行漏洞',
    status: 'open',
    discoveredDate: '2024-03-15',
  },
  {
    cveId: 'CVE-2024-1235',
    severity: 'medium',
    component: 'Vue.js',
    version: '3.4.19',
    description: 'XSS跨站脚本漏洞',
    status: 'resolved',
    discoveredDate: '2024-03-10',
  },
])

const blockRecords = ref([
  {
    date: '2024-03-15 10:30',
    reason: '引入了未审批的开源组件 log4j 2.14.0',
    rule: '组件白名单检查',
    handler: '张三',
    status: 'resolved',
  },
  {
    date: '2024-03-10 14:20',
    reason: '生产Commit ID与当前Commit ID不一致',
    rule: '版本一致性检查',
    handler: '李四',
    status: 'resolved',
  },
  {
    date: '2024-03-05 09:15',
    reason: '使用了GPL许可证的组件，不符合许可证要求',
    rule: '许可证检查',
    handler: '王五',
    status: 'pending',
  },
])

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
