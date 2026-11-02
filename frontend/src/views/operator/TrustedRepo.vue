<template>
  <div class="trusted-repo">
    <el-tabs v-model="activeTab" class="repo-tabs">
      <!-- Tab1: 核心台账 -->
      <el-tab-pane label="核心台账" name="ledger">
        <div class="tab-content">
          <!-- 类型切换按钮 -->
          <div class="type-switcher">
            <el-button-group>
              <el-button
                :type="ledgerType === 'software' ? 'primary' : ''"
                @click="ledgerType = 'software'"
              >
                基础软件
              </el-button>
              <el-button
                :type="ledgerType === 'dependency' ? 'primary' : ''"
                @click="ledgerType = 'dependency'"
              >
                依赖包
              </el-button>
              <el-button
                :type="ledgerType === 'tool' ? 'primary' : ''"
                @click="ledgerType = 'tool'"
              >
                官方工具包
              </el-button>
            </el-button-group>
          </div>

          <!-- 依赖包表格 -->
          <el-table v-if="ledgerType === 'dependency'" :data="dependencyList" border style="width: 100%">
            <el-table-column prop="name" label="组件名称" />
            <el-table-column prop="version" label="推荐版本" width="120" />
            <el-table-column prop="license" label="许可证" width="120" />
            <el-table-column prop="projectCount" label="关联项目" width="100" align="center">
              <template #default="{ row }">
                <el-button type="primary" link @click="showRelatedProjects(row)">
                  {{ row.projectCount }}
                </el-button>
              </template>
            </el-table-column>
            <el-table-column label="是否核心" width="100" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.isCore" type="danger" size="small">核心</el-tag>
                <el-tag v-else type="info" size="small">普通</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="说明" min-width="200" show-overflow-tooltip />
          </el-table>

          <!-- 基础软件表格 -->
          <el-table v-else-if="ledgerType === 'software'" :data="softwareList" border style="width: 100%">
            <el-table-column prop="name" label="软件名称" />
            <el-table-column prop="version" label="推荐版本" width="120" />
            <el-table-column prop="license" label="许可证" width="120" />
            <el-table-column prop="description" label="说明" min-width="200" show-overflow-tooltip />
          </el-table>

          <!-- 工具包表格 -->
          <el-table v-else :data="toolList" border style="width: 100%">
            <el-table-column prop="name" label="工具名称" />
            <el-table-column prop="version" label="版本" width="120" />
            <el-table-column prop="platform" label="支持平台" width="150" />
            <el-table-column prop="description" label="说明" min-width="200" show-overflow-tooltip />
          </el-table>
        </div>
      </el-tab-pane>

      <!-- Tab2: 软件库 -->
      <el-tab-pane label="软件库" name="library">
        <div class="tab-content">
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

          <!-- 搜索结果 -->
          <div v-if="searched && searchResults.length > 0" class="search-results">
            <el-card
              v-for="item in searchResults"
              :key="item.id"
              class="result-card"
              shadow="hover"
              style="margin-bottom: 16px"
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
                  </div>
                  <el-button
                    v-if="!item.introduced"
                    type="primary"
                    size="small"
                    @click="applySoftwareFromSearch(item)"
                  >
                    申请使用
                  </el-button>
                  <el-button v-else size="small" disabled>已引入</el-button>
                </div>
              </template>

              <p class="component-desc">{{ item.description }}</p>

              <el-row :gutter="20" class="metrics-row">
                <el-col :span="8">
                  <div class="metric-box">
                    <div class="metric-header">
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
                <el-col :span="8">
                  <div class="metric-box">
                    <div class="metric-header">
                      <span class="metric-title">许可证</span>
                    </div>
                    <div class="license-info">
                      <span class="license-name">{{ item.licenseName }}</span>
                      <el-tag :type="item.licenseRisk === 'low' ? 'success' : 'warning'" size="small">
                        {{ item.licenseRisk === 'low' ? '低风险' : '高风险' }}
                      </el-tag>
                    </div>
                  </div>
                </el-col>
                <el-col :span="8">
                  <div class="metric-box">
                    <div class="metric-header">
                      <span class="metric-title">引入状态</span>
                    </div>
                    <span v-if="item.introduced" class="introduced-info">已引入 v{{ item.introducedVersion }}</span>
                    <span v-else class="not-introduced">未引入</span>
                  </div>
                </el-col>
              </el-row>
            </el-card>
          </div>

          <!-- 软件库列表（未搜索时显示） -->
          <div v-if="!searched" class="library-header">
            <el-button type="primary" @click="showApplyPanel">
              <el-icon><Plus /></el-icon>
              申请使用
            </el-button>
          </div>

          <el-table v-if="!searched" :data="softwareLibrary" border style="width: 100%; margin-top: 16px">
            <el-table-column prop="name" label="组件名称" />
            <el-table-column prop="org" label="组织" width="150" />
            <el-table-column prop="version" label="推荐版本" width="120" />
            <el-table-column prop="license" label="许可证" width="120" />
            <el-table-column label="状态" width="120" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.status === 'approved'" type="success" size="small">已通过</el-tag>
                <el-tag v-else-if="row.status === 'pending'" type="warning" size="small">审批中</el-tag>
                <el-tag v-else type="info" size="small">未申请</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" align="center">
              <template #default="{ row }">
                <el-button
                  v-if="row.status === 'approved' && row.downloadUrl"
                  type="primary"
                  link
                  @click="downloadSoftware(row.downloadUrl)"
                >
                  下载
                </el-button>
                <el-button
                  v-else-if="row.status === 'pending'"
                  type="info"
                  link
                  disabled
                >
                  审批中
                </el-button>
                <el-button
                  v-else
                  type="primary"
                  link
                  @click="applySoftware(row)"
                >
                  申请使用
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <!-- 无搜索结果 -->
          <el-empty
            v-else-if="searchResults.length === 0"
            description="未找到匹配的组件，请尝试其他关键词"
            :image-size="120"
          >
            <el-button type="primary" @click="resetSearch">重置搜索</el-button>
          </el-empty>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 关联项目对话框 -->
    <el-dialog
      v-model="relatedProjectsVisible"
      title="关联项目"
      width="700px"
    >
      <div class="related-projects">
        <h4>{{ currentSoftware.name }} - 关联的内部项目</h4>
        <el-table :data="relatedProjects" border>
          <el-table-column prop="name" label="项目名称" />
          <el-table-column prop="team" label="团队" width="150" />
          <el-table-column prop="dept" label="部门" width="150" />
        </el-table>
      </div>
    </el-dialog>

    <!-- 申请使用面板（右侧滑出） -->
    <el-drawer
      v-model="applyDrawerVisible"
      title="申请使用组件"
      direction="rtl"
      size="50%"
    >
      <el-form :model="applyForm" label-width="120px">
        <el-divider content-position="left">组件信息</el-divider>
        <el-form-item label="组织">
          <el-input v-model="applyForm.org" placeholder="如：org.apache" />
        </el-form-item>
        <el-form-item label="组件名称">
          <el-input v-model="applyForm.name" placeholder="如：kafka" />
        </el-form-item>
        <el-form-item label="推荐版本">
          <el-input v-model="applyForm.version" placeholder="如：3.6.1" />
        </el-form-item>

        <el-divider content-position="left">申请信息</el-divider>
        <el-form-item label="申请用途" required>
          <el-input
            v-model="applyForm.purpose"
            type="textarea"
            :rows="4"
            placeholder="请说明申请使用该组件的用途"
          />
        </el-form-item>
        <el-form-item label="期望交付时间" required>
          <el-date-picker
            v-model="applyForm.expectedDate"
            type="date"
            placeholder="选择期望交付时间"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="drawer-footer">
          <el-button @click="applyDrawerVisible = false">取消</el-button>
          <el-button type="primary" @click="submitApplication">提交</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// Tab切换
const activeTab = ref('ledger')
const ledgerType = ref<'software' | 'dependency' | 'tool'>('software')

// 搜索相关
const searchQuery = ref('')
const searched = ref(false)
const searchResults = ref<any[]>([])

function handleSearch() {
  if (!searchQuery.value.trim()) {
    ElMessage.warning('请输入搜索关键词')
    return
  }
  searched.value = true
  // TODO: 调用后端搜索 API
  searchResults.value = []
}

function resetSearch() {
  searchQuery.value = ''
  searched.value = false
  searchResults.value = []
}

function getActivityColor(score: number) {
  if (score >= 80) return '#67C23A'
  if (score >= 60) return '#E6A23C'
  return '#F56C6C'
}

function applySoftwareFromSearch(item: any) {
  applyForm.org = ''
  applyForm.name = item.name
  applyForm.version = item.version
  applyForm.purpose = ''
  applyForm.expectedDate = ''
  applyDrawerVisible.value = true
}

// 核心台账数据
const softwareList = ref<any[]>([])
const dependencyList = ref<any[]>([])
const toolList = ref<any[]>([])

// 软件库数据
const softwareLibrary = ref<any[]>([])

// 关联项目
const relatedProjectsVisible = ref(false)
const currentSoftware = ref({ name: '' })
const relatedProjects = ref<any[]>([])

// 申请使用面板
const applyDrawerVisible = ref(false)
const applyForm = reactive({
  org: '',
  name: '',
  version: '',
  purpose: '',
  expectedDate: '',
})

function showRelatedProjects(row: any) {
  currentSoftware.value = { name: row.name }
  relatedProjectsVisible.value = true
}

function loadLedgerData() {
  // TODO: 调用后端 API 获取核心台账数据
}

function loadSoftwareLibrary() {
  // TODO: 调用后端 API 获取软件库数据
}

onMounted(() => {
  loadLedgerData()
  loadSoftwareLibrary()
})

function showApplyPanel() {
  applyForm.org = ''
  applyForm.name = ''
  applyForm.version = ''
  applyForm.purpose = ''
  applyForm.expectedDate = ''
  applyDrawerVisible.value = true
}

function applySoftware(row: any) {
  applyForm.org = row.org
  applyForm.name = row.name
  applyForm.version = row.version
  applyForm.purpose = ''
  applyForm.expectedDate = ''
  applyDrawerVisible.value = true
}

function submitApplication() {
  if (!applyForm.purpose || !applyForm.expectedDate) {
    ElMessage.warning('请填写完整的申请信息')
    return
  }

  // 模拟提交
  applyDrawerVisible.value = false
  ElMessage.success('已创建审批任务')
}

function downloadSoftware(url: string) {
  // 模拟下载
  ElMessage.success('开始下载...')
  // window.open(url, '_blank')
}
</script>

<style scoped lang="scss">
.trusted-repo {
  padding: 20px;
  background: #f5f7fa;
  min-height: calc(100vh - 60px);
}

.repo-tabs {
  background: #fff;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
}

.tab-content {
  margin-top: 16px;
}

.type-switcher {
  margin-bottom: 16px;
}

.search-section {
  margin-bottom: 24px;
}

.search-tips {
  margin-top: 8px;
}

.search-results {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.result-card {
  .result-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-left {
      display: flex;
      align-items: center;
      gap: 8px;

      .result-name {
        font-size: 16px;
        font-weight: 600;
        color: #303133;
      }
    }
  }

  .component-desc {
    margin: 12px 0;
    color: #606266;
    font-size: 14px;
  }

  .metrics-row {
    margin-top: 16px;
  }

  .metric-box {
    padding: 12px;
    background: #f5f7fa;
    border-radius: 4px;

    .metric-header {
      display: flex;
      align-items: center;
      margin-bottom: 8px;

      .metric-title {
        font-size: 13px;
        color: #606266;
      }
    }

    .license-info {
      display: flex;
      align-items: center;
      gap: 8px;

      .license-name {
        font-size: 14px;
        color: #303133;
      }
    }

    .introduced-info {
      font-size: 14px;
      color: #67C23A;
    }

    .not-introduced {
      font-size: 14px;
      color: #909399;
    }
  }
}

.library-header {
  margin-bottom: 16px;
  display: flex;
  justify-content: flex-end;
}

.related-projects {
  h4 {
    margin-bottom: 16px;
    color: #303133;
  }
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
