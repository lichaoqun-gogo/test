<template>
  <div class="software-library">
    <div class="page-header">
      <div class="filter-bar">
        <el-input v-model="searchText" placeholder="搜索软件名称" clearable style="width: 200px" />
        <el-select v-model="filterGrade" placeholder="等级" clearable style="width: 120px; margin-left: 12px">
          <el-option v-for="g in ['A', 'B', 'C', 'D']" :key="g" :label="`${g}级`" :value="g" />
        </el-select>
        <el-select v-model="filterType" placeholder="类型" clearable style="width: 120px; margin-left: 12px">
          <el-option label="中间件" value="middleware" />
          <el-option label="数据库" value="database" />
          <el-option label="工具" value="tool" />
          <el-option label="依赖包" value="dependency" />
          <el-option label="组件库" value="library" />
        </el-select>
      </div>
      <el-button type="primary" @click="dialogVisible = true">
        <el-icon><Plus /></el-icon> 新增
      </el-button>
    </div>

    <el-table :data="filteredAssets" stripe style="margin-top: 16px">
      <el-table-column prop="name" label="软件名称" min-width="140" />
      <el-table-column prop="version" label="版本" width="100" />
      <el-table-column label="类型" width="100">
        <template #default="{ row }">
          {{ SOFTWARE_TYPE_LABELS[row.type as keyof typeof SOFTWARE_TYPE_LABELS] }}
        </template>
      </el-table-column>
      <el-table-column label="等级" width="120">
        <template #default="{ row }">
          <GradeTag :grade="row.grade" size="small" />
        </template>
      </el-table-column>
      <el-table-column prop="language" label="语言" width="100" />
      <el-table-column prop="organization" label="组织" width="140" />
      <el-table-column prop="owner" label="负责人" width="100" />
      <el-table-column prop="vulnerabilityCount" label="漏洞数" width="80">
        <template #default="{ row }">
          <el-tag v-if="row.vulnerabilityCount > 0" type="danger" size="small">{{ row.vulnerabilityCount }}</el-tag>
          <span v-else>0</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : row.status === 'deprecated' ? 'info' : 'warning'" size="small">
            {{ row.status === 'active' ? '使用中' : row.status === 'deprecated' ? '已停用' : '待审核' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="140" fixed="right">
        <template #default>
          <el-button type="primary" link size="small">编辑</el-button>
          <el-button type="danger" link size="small">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- Add Dialog -->
    <el-dialog v-model="dialogVisible" title="新增开源基础软件" width="560px">
      <el-form label-width="100px">
        <el-form-item label="软件名称"><el-input placeholder="请输入" /></el-form-item>
        <el-form-item label="版本"><el-input placeholder="请输入" /></el-form-item>
        <el-form-item label="类型">
          <el-select placeholder="请选择" style="width: 100%">
            <el-option v-for="(label, key) in SOFTWARE_TYPE_LABELS" :key="key" :label="label" :value="key" />
          </el-select>
        </el-form-item>
        <el-form-item label="等级">
          <el-select placeholder="请选择" style="width: 100%">
            <el-option label="A - 核心关键" value="A" />
            <el-option label="B - 重要支撑" value="B" />
            <el-option label="C - 一般使用" value="C" />
            <el-option label="D - 探索试用" value="D" />
          </el-select>
        </el-form-item>
        <el-form-item label="开发语言"><el-input placeholder="请输入" /></el-form-item>
        <el-form-item label="组织"><el-input placeholder="请输入" /></el-form-item>
        <el-form-item label="负责人"><el-input placeholder="请输入" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="dialogVisible = false">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { mockFullAssets } from '@/mock/admin'
import { SOFTWARE_TYPE_LABELS } from '@/types'
import { Plus } from '@element-plus/icons-vue'

const searchText = ref('')
const filterGrade = ref('')
const filterType = ref('')
const dialogVisible = ref(false)

const filteredAssets = computed(() => {
  let list = mockFullAssets
  if (searchText.value) {
    const q = searchText.value.toLowerCase()
    list = list.filter(a => a.name.toLowerCase().includes(q))
  }
  if (filterGrade.value) list = list.filter(a => a.grade === filterGrade.value)
  if (filterType.value) list = list.filter(a => a.type === filterType.value)
  return list
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.filter-bar {
  display: flex;
  align-items: center;
}
</style>
