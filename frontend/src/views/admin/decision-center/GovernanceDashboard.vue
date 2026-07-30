<template>
  <div class="governance-dashboard">
    <h3>治理效能看板</h3>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="12">
        <el-card>
          <template #header>违规拦截统计</template>
          <div class="chart-placeholder">
            <el-empty description="趋势图表区域" :image-size="60" />
            <div v-for="item in dashboardData.blockStats" :key="item.month" class="bar-item">
              <span class="bar-label">{{ item.month }}</span>
              <el-progress :percentage="(item.count / 30) * 100" :stroke-width="16" :show-text="false"
                :color="item.count > 20 ? '#F56C6C' : '#E6A23C'" style="flex: 1; margin: 0 8px" />
              <span class="bar-value">{{ item.count }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>缺陷修复效能（MTTR）</template>
          <el-table :data="dashboardData.repairStats" size="small">
            <el-table-column prop="department" label="部门" />
            <el-table-column label="平均修复时长" width="120">
              <template #default="{ row }">{{ row.mttr }} 天</template>
            </el-table-column>
            <el-table-column label="逾期未修复" width="100">
              <template #default="{ row }">
                <el-tag v-if="row.overdue > 0" type="danger" size="small">{{ row.overdue }}</el-tag>
                <span v-else>0</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="12">
        <el-card>
          <template #header>资产健康度分布</template>
          <div class="health-stats">
            <div v-for="item in dashboardData.healthStats" :key="item.status" class="health-item">
              <el-progress type="dashboard" :percentage="item.percentage"
                :color="healthColor(item.status)" :width="100">
                <template #default="{ percentage }">
                  <span class="percentage-text">{{ percentage }}%</span>
                </template>
              </el-progress>
              <p class="health-label">{{ item.status }} ({{ item.count }})</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>治理工作量统计</template>
          <el-table :data="dashboardData.workloadStats" size="small">
            <el-table-column prop="department" label="部门" />
            <el-table-column prop="introduce" label="引入申请" width="100" />
            <el-table-column prop="upgrade" label="升级审批" width="100" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { mockDashboardData } from '@/mock/admin'

const dashboardData = mockDashboardData

function healthColor(status: string) {
  if (status === '正常') return '#67C23A'
  if (status === '停更') return '#E6A23C'
  return '#F56C6C'
}
</script>

<style scoped>
.chart-placeholder { padding: 12px 0; }
.bar-item { display: flex; align-items: center; margin: 8px 0; }
.bar-label { width: 70px; font-size: 12px; }
.bar-value { width: 30px; text-align: right; font-size: 12px; color: #606266; }
.health-stats { display: flex; justify-content: space-around; padding: 16px 0; }
.health-item { text-align: center; }
.health-label { margin: 8px 0 0; font-size: 13px; color: #606266; }
.percentage-text { font-size: 16px; font-weight: 600; }
</style>
