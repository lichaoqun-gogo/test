<template>
  <div class="blacklist-manage">
    <div class="page-header">
      <el-input v-model="searchText" placeholder="搜索黑名单组件" clearable style="width: 200px" />
      <div>
        <el-button type="success" plain><el-icon><Upload /></el-icon> 批量导入</el-button>
        <el-button type="primary" @click="dialogVisible = true"><el-icon><Plus /></el-icon> 新增</el-button>
      </div>
    </div>

    <el-table :data="blacklist" stripe style="margin-top: 16px">
      <el-table-column prop="name" label="组件名称" min-width="140" />
      <el-table-column prop="version" label="版本限制" width="120" />
      <el-table-column label="风险等级" width="100">
        <template #default="{ row }"><RiskTag :level="row.riskLevel" /></template>
      </el-table-column>
      <el-table-column prop="reason" label="加入原因" min-width="200" />
      <el-table-column prop="addedBy" label="添加人" width="100" />
      <el-table-column prop="addedAt" label="添加时间" width="120" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'danger' : 'info'" size="small">
            {{ row.status === 'active' ? '生效' : '已移除' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="140">
        <template #default><el-button type="primary" link size="small">编辑</el-button><el-button type="danger" link size="small">移除</el-button></template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="新增黑名单" width="480px">
      <el-form label-width="100px">
        <el-form-item label="组件名称" required><el-input placeholder="请输入" /></el-form-item>
        <el-form-item label="版本限制"><el-input placeholder="如 >=1.0.0" /></el-form-item>
        <el-form-item label="风险等级" required>
          <el-select placeholder="请选择" style="width: 100%">
            <el-option label="严重" value="critical" />
            <el-option label="高风险" value="high" />
            <el-option label="中风险" value="medium" />
          </el-select>
        </el-form-item>
        <el-form-item label="加入原因" required><el-input type="textarea" :rows="3" placeholder="请说明原因" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="dialogVisible = false">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { mockBlacklist } from '@/mock/admin'
import { Upload, Plus } from '@element-plus/icons-vue'

const searchText = ref('')
const blacklist = mockBlacklist
const dialogVisible = ref(false)
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; }
</style>
