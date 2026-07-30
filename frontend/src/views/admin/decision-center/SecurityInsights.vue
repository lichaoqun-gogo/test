<template>
  <div class="security-insights">
    <h3>供应链安全洞察</h3>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="8">
        <el-card>
          <template #header>核心依赖 Top 5</template>
          <p class="desc">统计全公司使用频率最高的开源组件</p>
          <el-table :data="dashboardData.coreDeps" size="small">
            <el-table-column type="index" label="排名" width="60" />
            <el-table-column prop="name" label="组件" width="120" />
            <el-table-column prop="usageCount" label="使用次数" width="90" />
            <el-table-column prop="projects" label="项目数" width="80" />
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>版本碎片化 Top 4</template>
          <p class="desc">统计版本最混乱的组件</p>
          <el-table :data="dashboardData.versionFragment" size="small">
            <el-table-column type="index" label="排名" width="60" />
            <el-table-column prop="name" label="组件" width="120" />
            <el-table-column prop="versionCount" label="版本数" width="70" />
            <el-table-column label="版本列表" min-width="120">
              <template #default="{ row }">
                <el-tag v-for="v in row.versions" :key="v" size="small" style="margin: 2px">{{ v }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>高风险项目 Top 4</template>
          <p class="desc">关联漏洞/缺陷最多的自研项目</p>
          <el-table :data="dashboardData.riskProjects" size="small">
            <el-table-column type="index" label="排名" width="60" />
            <el-table-column prop="name" label="项目" width="120" />
            <el-table-column label="风险等级" width="90">
              <template #default="{ row }"><RiskTag :level="row.riskLevel" size="small" /></template>
            </el-table-column>
            <el-table-column prop="vulnerabilityCount" label="漏洞数" width="80" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-card style="margin-top: 16px">
      <template #header>单点故障风险分析</template>
      <el-alert
        v-for="dep in dashboardData.coreDeps.slice(0, 2)"
        :key="dep.name"
        :title="`${dep.name} 在 ${dep.projects} 个项目中使用，存在单点故障风险`"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 8px"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { mockDashboardData } from '@/mock/admin'

const dashboardData = mockDashboardData
</script>

<style scoped>
.desc { color: #909399; font-size: 12px; margin: 0 0 8px; }
</style>
