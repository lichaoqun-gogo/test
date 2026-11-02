<template>
  <div class="dependency-library">
    <div class="page-header">
      <el-input v-model="searchText" placeholder="搜索依赖包名称" clearable style="width: 200px" />
      <div>
        <el-button type="success" plain><el-icon><Upload /></el-icon> 批量导入</el-button>
        <el-button type="primary" plain><el-icon><Download /></el-icon> 导出</el-button>
      </div>
    </div>
    <el-table :data="mockFullAssets.filter(a => a.type === 'dependency' || a.type === 'library')" stripe style="margin-top: 16px">
      <el-table-column prop="name" label="依赖名称" min-width="140" />
      <el-table-column prop="version" label="版本" width="100" />
      <el-table-column prop="language" label="语言" width="100" />
      <el-table-column prop="licenseType" label="许可证" width="120" />
      <el-table-column label="许可证风险" width="110">
        <template #default="{ row }"><RiskTag :level="row.licenseRisk" /></template>
      </el-table-column>
      <el-table-column prop="vulnerabilityCount" label="漏洞数" width="80" />
      <el-table-column label="操作" width="140">
        <template #default><el-button type="primary" link size="small">编辑</el-button><el-button type="danger" link size="small">删除</el-button></template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { mockFullAssets } from '@/mock/admin'
import { Upload, Download } from '@element-plus/icons-vue'

const searchText = ref('')
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; }
</style>
