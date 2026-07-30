<template>
  <div class="cicd-check">
    <div class="page-header">
      <div class="header-left">
        <h3>CICD检查</h3>
      </div>
    </div>

    <div class="toolbar">
      <el-select v-model="queryParams.level" placeholder="适用等级" clearable style="width: 140px">
        <el-option label="A 可安全使用" value="A" />
        <el-option label="B 谨慎使用" value="B" />
        <el-option label="C 限期整改" value="C" />
        <el-option label="D 禁止使用" value="D" />
      </el-select>
      <el-select v-model="queryParams.deptName" placeholder="部门" clearable style="width: 150px">
        <el-option v-for="dept in departmentOptions" :key="dept" :label="dept" :value="dept" />
      </el-select>
      <el-select v-model="queryParams.groupName" placeholder="团队" clearable style="width: 150px">
        <el-option v-for="group in groupOptions" :key="group" :label="group" :value="group" />
      </el-select>
      <el-select v-model="queryParams.disposalType" placeholder="处置类型" clearable style="width: 130px">
        <el-option label="阻断" value="阻断" />
        <el-option label="警告" value="警告" />
      </el-select>
      <el-input v-model="queryParams.projectName" placeholder="项目名称" clearable style="width: 200px" @keyup.enter="handleSearch" />
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <el-button @click="handleReset">重置</el-button>
    </div>

    <el-table :data="paginatedList" stripe style="margin-top: 16px" v-loading="loading">
      <el-table-column prop="pro_code" label="项目名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="disposalType" label="处置类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.disposalType === '阻断' ? 'danger' : 'warning'" size="small">
            {{ row.disposalType }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="current_commit_id" label="当前版本" width="320">
        <template #default="{ row }">
          <el-link type="primary" @click="openDetail(row as any)">{{ row.current_commit_id }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="group_name" label="团队" width="120" />
      <el-table-column prop="dept_name" label="部门" width="120" />
    </el-table>

    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadData"
        @size-change="loadData"
      />
    </div>

    <!-- 详情弹框 -->
    <el-dialog v-model="detailVisible" :title="detailProjectName" width="900px">
      <div class="commit-query">
        <span class="query-label">当前Commit ID：{{ detailCurrentCommitId }}</span>
        <el-input v-model="detailCommitId" placeholder="输入历史Commit ID查询" style="width: 300px; margin-left: 16px" @keyup.enter="loadDetailData">
          <template #append>
            <el-button @click="loadDetailData">查询</el-button>
          </template>
        </el-input>
      </div>

      <el-table :data="detailList" stripe style="margin-top: 16px" v-loading="detailLoading">
        <el-table-column prop="moduleName" label="模块名称" min-width="150" />
        <el-table-column prop="checkTime" label="核查时间" width="180" />
        <el-table-column prop="disposal" label="处置" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.disposal === '阻断' ? 'danger' : 'warning'" size="small">
              {{ row.disposal }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="blacklistName" label="黑名单" min-width="150" />
        <el-table-column prop="reason" label="原因" min-width="200" show-overflow-tooltip />
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

interface CicdItem {
  id: number
  pro_code: string
  current_commit_id: string
  group_name: string
  dept_name: string
  disposalType: string
  level: string
}

interface CicdDetail {
  moduleName: string
  checkTime: string
  disposal: string
  blacklistName: string
  reason: string
}

const cicdList = ref<CicdItem[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const queryParams = ref({
  projectName: '',
  deptName: '',
  groupName: '',
  disposalType: '',
  level: '',
})

const departmentOptions = computed(() => {
  const depts = cicdList.value.map(item => item.dept_name).filter(Boolean)
  return [...new Set(depts)]
})

const groupOptions = computed(() => {
  const groups = cicdList.value.map(item => item.group_name).filter(Boolean)
  return [...new Set(groups)]
})

const paginatedList = computed(() => cicdList.value)

async function loadData() {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      pageSize: pageSize.value,
      projectName: queryParams.value.projectName,
      deptName: queryParams.value.deptName,
      groupName: queryParams.value.groupName,
      disposalType: queryParams.value.disposalType,
      level: queryParams.value.level,
    }
    const res = await axios.get('/api/v1/risk/cicd-check', { params })
    if (res.data.code === 200 || res.data.code === 0) {
      cicdList.value = res.data.data.list || []
      total.value = res.data.data.total || 0
    } else {
      ElMessage.error(res.data.msg || '请求失败')
    }
  } catch (error) {
    console.error('加载CICD数据失败:', error)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  currentPage.value = 1
  loadData()
}

function handleReset() {
  queryParams.value = {
    projectName: '',
    deptName: '',
    groupName: '',
    disposalType: '',
    level: '',
  }
  currentPage.value = 1
  loadData()
}

// 详情弹框
const detailVisible = ref(false)
const detailProjectName = ref('')
const detailCurrentCommitId = ref('')
const detailCommitId = ref('')
const detailList = ref<CicdDetail[]>([])
const detailLoading = ref(false)

function openDetail(row: CicdItem) {
  detailProjectName.value = row.pro_code
  detailCurrentCommitId.value = row.current_commit_id
  detailCommitId.value = row.current_commit_id
  detailVisible.value = true
  loadDetailData()
}

async function loadDetailData() {
  detailLoading.value = true
  try {
    const params: any = {
      projectName: detailProjectName.value,
    }
    if (detailCommitId.value) {
      params.commitId = detailCommitId.value
    }
    const res = await axios.get('/api/v1/risk/cicd-check/detail', { params })
    if (res.data.code === 200 || res.data.code === 0) {
      detailList.value = res.data.data || []
    }
  } catch (error) {
    console.error('加载详情失败:', error)
  } finally {
    detailLoading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-header {
  margin-bottom: 16px;
}

.page-header h3 {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.commit-query {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.query-label {
  white-space: nowrap;
  font-weight: 500;
}
</style>
