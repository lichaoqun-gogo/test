<template>
  <div class="simple-query">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">简单查询示例</h2>
    </div>

    <!-- 查询区域 -->
    <div class="filter-section">
      <el-form :inline="true" :model="queryForm" class="query-form">
        <el-form-item label="查询条件">
          <el-select 
            v-model="queryForm.status" 
            placeholder="请选择状态" 
            clearable 
            style="width: 180px"
            @change="handleQuery"
          >
            <el-option label="全部" value="" />
            <el-option label="启用" value="active" />
            <el-option label="禁用" value="inactive" />
            <el-option label="待审核" value="pending" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格区域 -->
    <div class="table-section">
      <el-table 
        :data="tableData" 
        stripe 
        border
        v-loading="loading"
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="名称" min-width="150" />
        <el-table-column prop="type" label="类型" width="120" align="center" />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag 
              :type="getStatusType(row.status)" 
              size="small"
            >
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="createTime" label="创建时间" width="180" />
      </el-table>

      <!-- 空状态 -->
      <el-empty 
        v-if="!loading && tableData.length === 0" 
        description="暂无数据" 
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'

// 查询表单
const queryForm = ref({
  status: ''
})

// 表格数据
const tableData = ref<any[]>([])
const loading = ref(false)

// 状态类型
function getStatusType(status: string): 'success' | 'info' | 'warning' | 'danger' | undefined {
  switch (status) {
    case 'active': return 'success'
    case 'inactive': return 'info'
    case 'pending': return 'warning'
    default: return undefined
  }
}

// 状态标签
function getStatusLabel(status: string) {
  switch (status) {
    case 'active': return '启用'
    case 'inactive': return '禁用'
    case 'pending': return '待审核'
    default: return status
  }
}

// 查询
async function handleQuery() {
  await fetchData(queryForm.value)
}

// 重置
async function handleReset() {
  queryForm.value.status = ''
  await fetchData()
}

// 获取数据
async function fetchData(params?: any) {
  loading.value = true
  try {
    // TODO: 替换为实际的API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    const mockData = [
      { id: 1, name: '测试项目1', type: '项目', status: 'active', description: '这是一个测试项目', createTime: '2024-01-15 10:30:00' },
      { id: 2, name: '测试项目2', type: '组件', status: 'inactive', description: '这是一个测试组件', createTime: '2024-01-16 14:20:00' },
      { id: 3, name: '测试项目3', type: '项目', status: 'pending', description: '等待审核中', createTime: '2024-01-17 09:15:00' },
      { id: 4, name: '测试项目4', type: '软件', status: 'active', description: '开源软件库', createTime: '2024-01-18 16:45:00' },
      { id: 5, name: '测试项目5', type: '组件', status: 'active', description: '核心依赖组件', createTime: '2024-01-19 11:30:00' },
    ]
    
    // 根据条件过滤
    if (params?.status) {
      tableData.value = mockData.filter(item => item.status === params.status)
    } else {
      tableData.value = mockData
    }
  } catch (error) {
    console.error('查询失败:', error)
  } finally {
    loading.value = false
  }
}

// 初始化
onMounted(() => {
  fetchData()
})
</script>
