<template>
  <div class="tool-library">
    <div class="page-header">
      <el-input v-model="searchText" placeholder="搜索工具名称" clearable style="width: 200px" />
      <el-button type="primary"><el-icon><Plus /></el-icon> 新增</el-button>
    </div>
    <el-table :data="mockFullAssets.filter(a => a.type === 'tool')" stripe style="margin-top: 16px">
      <el-table-column prop="name" label="工具名称" min-width="140" />
      <el-table-column prop="version" label="版本" width="100" />
      <el-table-column prop="organization" label="组织" width="140" />
      <el-table-column prop="owner" label="负责人" width="100" />
      <el-table-column prop="introducedDate" label="引入日期" width="120" />
      <el-table-column label="操作" width="140">
        <template #default><el-button type="primary" link size="small">编辑</el-button><el-button type="danger" link size="small">删除</el-button></template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { mockFullAssets } from '@/mock/admin'
import { Plus } from '@element-plus/icons-vue'
const searchText = ref('')
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; }
</style>
