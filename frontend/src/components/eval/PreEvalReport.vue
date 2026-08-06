<template>
  <div class="pre-eval-report">
    <div class="report-header">
      <el-icon :size="20" class="header-icon"><Document /></el-icon>
      <span class="header-title">预评估报告</span>
      <el-tag :type="conclusionTagType as any" size="small">{{ conclusionLabel }}</el-tag>
    </div>

    <!-- SCA扫描结果 -->
    <el-collapse v-model="activeCollapse" class="report-collapse">
      <el-collapse-item name="sca">
        <template #title>
          <div class="collapse-title">
            <el-icon><Aim /></el-icon>
            <span>SCA扫描结果</span>
            <el-badge :value="mockData.sca.vulnerabilityCount" :max="99" class="vuln-badge" />
          </div>
        </template>
        <div class="report-section">
          <div class="section-summary">
            <el-descriptions :column="2" border size="small">
              <el-descriptions-item label="组件名称">{{ componentName }}</el-descriptions-item>
              <el-descriptions-item label="扫描版本">{{ componentVersion }}</el-descriptions-item>
              <el-descriptions-item label="已知漏洞数">
                <el-tag :type="vulnCountType" size="small">{{ mockData.sca.vulnerabilityCount }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="依赖项总数">{{ mockData.sca.dependencyCount }}</el-descriptions-item>
            </el-descriptions>
          </div>

          <div v-if="mockData.sca.vulnerabilities.length > 0" class="vulnerability-list">
            <h4 class="subsection-title">漏洞详情</h4>
            <el-table :data="mockData.sca.vulnerabilities" size="small" border>
              <el-table-column prop="cveId" label="CVE编号" width="160" />
              <el-table-column prop="severity" label="严重等级" width="100">
                <template #default="{ row }">
                  <el-tag :type="getRiskTagType(row.severity) as any" size="small">
                    {{ getRiskLabel(row.severity) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="description" label="漏洞描述" />
              <el-table-column prop="fixedVersion" label="修复版本" width="120" />
            </el-table>
          </div>

          <div v-if="mockData.sca.dependencies.length > 0" class="dependency-list">
            <h4 class="subsection-title">依赖树分析（Top 5）</h4>
            <el-table :data="mockData.sca.dependencies" size="small" border>
              <el-table-column prop="name" label="依赖名称" />
              <el-table-column prop="version" label="版本" width="120" />
              <el-table-column prop="license" label="许可证" width="120" />
              <el-table-column prop="status" label="状态" width="80">
                <template #default="{ row }">
                  <el-tag :type="(row.status === 'safe' ? 'success' : 'warning') as any" size="small">
                    {{ row.status === 'safe' ? '安全' : '风险' }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
      </el-collapse-item>

      <!-- 黑白名单匹配 -->
      <el-collapse-item name="blacklist">
        <template #title>
          <div class="collapse-title">
            <el-icon><Switch /></el-icon>
            <span>黑白名单匹配</span>
            <el-tag :type="(mockData.blacklist.hitBlacklist ? 'danger' : 'success') as any" size="small">
              {{ mockData.blacklist.hitBlacklist ? '命中黑名单' : '未命中黑名单' }}
            </el-tag>
          </div>
        </template>
        <div class="report-section">
          <el-alert
            v-if="mockData.blacklist.hitBlacklist"
            :title="mockData.blacklist.blacklistReason"
            type="error"
            :closable="false"
            show-icon
            style="margin-bottom: 16px"
          />
          <el-alert
            v-else
            title="组件未命中任何黑名单规则"
            type="success"
            :closable="false"
            show-icon
            style="margin-bottom: 16px"
          />

          <el-descriptions :column="1" border size="small">
            <el-descriptions-item label="白名单状态">
              <el-tag :type="(mockData.blacklist.inWhitelist ? 'success' : 'info') as any" size="small">
                {{ mockData.blacklist.inWhitelist ? '在白名单中' : '不在白名单中' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="匹配规则数">{{ mockData.blacklist.matchedRules }}</el-descriptions-item>
            <el-descriptions-item label="检查时间">{{ formatTime(mockData.blacklist.checkTime) }}</el-descriptions-item>
          </el-descriptions>

          <div v-if="mockData.blacklist.matchedRulesList.length > 0" class="matched-rules">
            <h4 class="subsection-title">匹配规则详情</h4>
            <ul class="rule-list">
              <li v-for="rule in mockData.blacklist.matchedRulesList" :key="rule.id">
                <el-tag :type="rule.type === 'blacklist' ? 'danger' : 'success'" size="small">
                  {{ rule.type === 'blacklist' ? '黑名单' : '白名单' }}
                </el-tag>
                <span class="rule-name">{{ rule.name }}</span>
                <span class="rule-desc">{{ rule.description }}</span>
              </li>
            </ul>
          </div>
        </div>
      </el-collapse-item>

      <!-- 许可证合规性 -->
      <el-collapse-item name="license">
        <template #title>
          <div class="collapse-title">
            <el-icon><Stamp /></el-icon>
            <span>许可证合规性</span>
            <el-tag :type="licenseRiskTagType as any" size="small">{{ mockData.license.riskLevelLabel }}</el-tag>
          </div>
        </template>
        <div class="report-section">
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="许可证类型">{{ mockData.license.licenseType }}</el-descriptions-item>
            <el-descriptions-item label="风险等级">
              <el-tag :type="licenseRiskTagType as any" size="small">{{ mockData.license.riskLevelLabel }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="OSI认证">
              <el-tag :type="(mockData.license.osiApproved ? 'success' : 'warning') as any" size="small">
                {{ mockData.license.osiApproved ? '是' : '否' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="商用许可">
              <el-tag :type="(mockData.license.commercialUse ? 'success' : 'danger') as any" size="small">
                {{ mockData.license.commercialUse ? '允许' : '禁止' }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>

          <div class="license-compliance">
            <h4 class="subsection-title">合规检查项</h4>
            <ul class="compliance-list">
              <li v-for="item in mockData.license.complianceChecks" :key="item.id" class="compliance-item">
                <el-icon :class="item.passed ? 'check-success' : 'check-fail'">
                  <CircleCheckFilled v-if="item.passed" />
                  <CircleCloseFilled v-else />
                </el-icon>
                <span class="check-label">{{ item.label }}</span>
                <el-tag :type="(item.passed ? 'success' : 'danger') as any" size="small" class="check-result">
                {{ item.passed ? '通过' : '不通过' }}
              </el-tag>
              </li>
            </ul>
          </div>

          <div v-if="mockData.license.suggestions.length > 0" class="license-suggestions">
            <h4 class="subsection-title">合规建议</h4>
            <el-alert
              v-for="suggestion in mockData.license.suggestions"
              :key="suggestion"
              :title="suggestion"
              type="warning"
              :closable="false"
              show-icon
              style="margin-bottom: 8px"
            />
          </div>
        </div>
      </el-collapse-item>

      <!-- 活跃度指标 -->
      <el-collapse-item name="activity">
        <template #title>
          <div class="collapse-title">
            <el-icon><TrendCharts /></el-icon>
            <span>活跃度指标</span>
            <el-tag :type="activityTagType as any" size="small">{{ mockData.activity.healthLevelLabel }}</el-tag>
          </div>
        </template>
        <div class="report-section">
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="社区健康度">
              <div class="health-score">
                <el-progress
                  :percentage="mockData.activity.healthScore"
                  :color="getHealthColor(mockData.activity.healthScore)"
                  :stroke-width="12"
                  :show-text="false"
                  style="width: 120px"
                />
                <span class="score-text">{{ mockData.activity.healthScore }}分</span>
              </div>
            </el-descriptions-item>
            <el-descriptions-item label="健康等级">
              <el-tag :type="activityTagType as any" size="small">{{ mockData.activity.healthLevelLabel }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="最近更新时间">{{ formatTime(mockData.activity.lastUpdated) }}</el-descriptions-item>
            <el-descriptions-item label="维护频率">{{ mockData.activity.maintenanceFreq }}</el-descriptions-item>
            <el-descriptions-item label="Star数">{{ mockData.activity.stars }}</el-descriptions-item>
            <el-descriptions-item label="Fork数">{{ mockData.activity.forks }}</el-descriptions-item>
            <el-descriptions-item label="贡献者数量">{{ mockData.activity.contributors }}</el-descriptions-item>
            <el-descriptions-item label="Open Issues">{{ mockData.activity.openIssues }}</el-descriptions-item>
          </el-descriptions>

          <div class="activity-details">
            <h4 class="subsection-title">详细指标</h4>
            <el-table :data="mockData.activity.metrics" size="small" border>
              <el-table-column prop="name" label="指标名称" />
              <el-table-column prop="value" label="数值" width="120" />
              <el-table-column prop="status" label="状态" width="80">
                <template #default="{ row }">
                  <el-tag :type="(row.status === 'good' ? 'success' : row.status === 'warning' ? 'warning' : 'danger') as any" size="small">
                    {{ row.status === 'good' ? '良好' : row.status === 'warning' ? '警告' : '差' }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  Document,
  Aim,
  Switch,
  Stamp,
  TrendCharts,
  CircleCheckFilled,
  CircleCloseFilled,
} from '@element-plus/icons-vue'
import type { PreEvalConclusion, RiskLevel } from '@/types'

const props = withDefaults(defineProps<{
  componentName: string
  componentVersion: string
  conclusion: PreEvalConclusion
  scaData?: any
  blacklistData?: any
  licenseData?: any
  activityData?: any
}>(), {})

const activeCollapse = ref(['sca', 'blacklist', 'license', 'activity'])

// 假数据
const mockData = computed(() => {
  const conclusion = props.conclusion

  return {
    sca: props.scaData || {
      vulnerabilityCount: conclusion === 'reject' ? 5 : conclusion === 'warning' ? 2 : 0,
      dependencyCount: 42,
      vulnerabilities: conclusion === 'reject'
        ? [
            { cveId: 'CVE-2024-1234', severity: 'critical', description: '远程代码执行漏洞', fixedVersion: '2.1.0' },
            { cveId: 'CVE-2024-5678', severity: 'high', description: 'SQL注入漏洞', fixedVersion: '2.0.5' },
            { cveId: 'CVE-2023-9012', severity: 'medium', description: '跨站脚本攻击', fixedVersion: '2.0.3' },
          ]
        : conclusion === 'warning'
        ? [
            { cveId: 'CVE-2024-3456', severity: 'medium', description: '信息泄露漏洞', fixedVersion: '2.1.0' },
            { cveId: 'CVE-2023-7890', severity: 'low', description: '拒绝服务风险', fixedVersion: '2.0.8' },
          ]
        : [],
      dependencies: [
        { name: 'lodash', version: '4.17.21', license: 'MIT', status: 'safe' },
        { name: 'axios', version: '1.6.0', license: 'MIT', status: 'safe' },
        { name: 'moment', version: '2.29.4', license: 'MIT', status: conclusion === 'warning' ? 'risk' : 'safe' },
        { name: 'express', version: '4.18.2', license: 'MIT', status: 'safe' },
        { name: 'debug', version: '4.3.4', license: 'MIT', status: 'safe' },
      ],
    },
    blacklist: props.blacklistData || {
      hitBlacklist: conclusion === 'reject',
      inWhitelist: conclusion === 'pass',
      blacklistReason: conclusion === 'reject'
        ? '该组件存在已知高危漏洞且长期未维护，已被列入黑名单'
        : '',
      matchedRules: conclusion === 'reject' ? 2 : 1,
      checkTime: '2024-06-13 10:30:00',
      matchedRulesList: conclusion === 'reject'
        ? [
            { id: '1', type: 'blacklist', name: '高危漏洞规则', description: '组件存在未修复的高危漏洞' },
            { id: '2', type: 'blacklist', name: '活跃度规则', description: '组件超过12个月未更新' },
          ]
        : [
            { id: '1', type: 'whitelist', name: '社区认可规则', description: '组件为Apache基金会项目' },
          ],
    },
    license: props.licenseData || {
      licenseType: conclusion === 'reject' ? 'GPL-3.0' : 'Apache-2.0',
      riskLevel: conclusion === 'reject' ? 'high' : conclusion === 'warning' ? 'medium' : 'low',
      riskLevelLabel: conclusion === 'reject' ? '高风险' : conclusion === 'warning' ? '中风险' : '低风险',
      osiApproved: true,
      commercialUse: conclusion !== 'reject',
      complianceChecks: [
        { id: '1', label: '许可证兼容性', passed: conclusion !== 'reject' },
        { id: '2', label: '传染性风险评估', passed: conclusion === 'pass' },
        { id: '3', label: '版权声明要求', passed: true },
        { id: '4', label: '源码公开义务', passed: conclusion !== 'reject' },
      ],
      suggestions: conclusion === 'reject'
        ? [
            'GPL-3.0许可证具有强传染性，建议评估是否适合商用',
            '需要确保所有衍生代码也以GPL-3.0开源',
          ]
        : conclusion === 'warning'
        ? [
            '建议确认许可证与项目当前使用的许可证兼容',
            '注意遵守许可证的版权声明要求',
          ]
        : [],
    },
    activity: props.activityData || {
      healthScore: conclusion === 'reject' ? 35 : conclusion === 'warning' ? 65 : 88,
      healthLevel: conclusion === 'reject' ? 'low' : conclusion === 'warning' ? 'medium' : 'high',
      healthLevelLabel: conclusion === 'reject' ? '不活跃' : conclusion === 'warning' ? '一般' : '活跃',
      lastUpdated: conclusion === 'reject' ? '2022-08-15' : '2024-05-20',
      maintenanceFreq: conclusion === 'reject' ? '12个月/次' : '2周/次',
      stars: 1250,
      forks: 320,
      contributors: 45,
      openIssues: conclusion === 'reject' ? 189 : 23,
      metrics: [
        { name: '代码提交频率', value: conclusion === 'reject' ? '低' : '高', status: conclusion === 'reject' ? 'bad' : 'good' },
        { name: 'Issue响应速度', value: conclusion === 'reject' ? '慢' : '快', status: conclusion === 'reject' ? 'bad' : 'good' },
        { name: 'PR合并率', value: '78%', status: 'good' },
        { name: '核心维护者数量', value: conclusion === 'reject' ? '1' : '5', status: conclusion === 'reject' ? 'bad' : 'good' },
      ],
    },
  }
})

const conclusionLabel = computed(() => {
  const labels: Record<PreEvalConclusion, string> = {
    pass: '通过',
    warning: '风险警告',
    reject: '不通过',
  }
  return labels[props.conclusion]
})

const conclusionTagType = computed(() => {
  const types: Record<PreEvalConclusion, string> = {
    pass: 'success',
    warning: 'warning',
    reject: 'danger',
  }
  return types[props.conclusion]
})

const vulnCountType = computed(() => {
  const count = mockData.value.sca.vulnerabilityCount
  if (count === 0) return 'success'
  if (count <= 2) return 'warning'
  return 'danger'
})

const licenseRiskTagType = computed(() => {
  const types: Record<string, string> = {
    low: 'success',
    medium: 'warning',
    high: 'danger',
  }
  return types[mockData.value.license.riskLevel] || 'info'
})

const activityTagType = computed(() => {
  const types: Record<string, string> = {
    high: 'success',
    medium: 'warning',
    low: 'danger',
  }
  return types[mockData.value.activity.healthLevel] || 'info'
})

function getRiskTagType(severity: string) {
  const types: Record<string, string> = {
    critical: 'danger',
    high: 'danger',
    medium: 'warning',
    low: 'info',
  }
  return types[severity] || 'info'
}

function getRiskLabel(severity: string) {
  const labels: Record<string, string> = {
    critical: '严重',
    high: '高',
    medium: '中',
    low: '低',
  }
  return labels[severity] || severity
}

function getHealthColor(score: number) {
  if (score >= 80) return '#67C23A'
  if (score >= 60) return '#E6A23C'
  return '#F56C6C'
}
</script>

<style scoped>
.pre-eval-report {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #e4e7ed;
}

.report-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e4e7ed;
}

.header-icon {
  color: #409EFF;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.report-collapse {
  border: none;
}

:deep(.el-collapse-item__header) {
  padding: 8px 0;
  font-weight: 500;
}

.collapse-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.vuln-badge {
  margin-left: auto;
}

.report-section {
  padding: 8px 0 16px;
}

.section-summary {
  margin-bottom: 16px;
}

.subsection-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin: 16px 0 12px;
  padding-left: 8px;
  border-left: 3px solid #409EFF;
}

.vulnerability-list,
.dependency-list,
.matched-rules,
.license-compliance,
.license-suggestions,
.activity-details {
  margin-top: 8px;
}

.rule-list {
  list-style: none;
  padding: 0;
}

.rule-list li {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 6px;
  margin-bottom: 8px;
  font-size: 13px;
}

.rule-name {
  font-weight: 500;
  color: #303133;
}

.rule-desc {
  color: #909399;
  margin-left: auto;
}

.compliance-list {
  list-style: none;
  padding: 0;
}

.compliance-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: #fafafa;
  border-radius: 6px;
  margin-bottom: 8px;
}

.check-success {
  color: #67C23A;
  font-size: 18px;
}

.check-fail {
  color: #F56C6C;
  font-size: 18px;
}

.check-label {
  flex: 1;
  font-size: 13px;
  color: #606266;
}

.check-result {
  flex-shrink: 0;
}

.health-score {
  display: flex;
  align-items: center;
  gap: 10px;
}

.score-text {
  font-weight: 600;
  color: #303133;
  font-size: 14px;
}

:deep(.el-descriptions__label) {
  width: 120px;
  font-weight: 500;
}
</style>
