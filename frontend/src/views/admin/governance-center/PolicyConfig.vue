<template>
  <div class="policy-config">
    <h3>评估指标配置</h3>
    <p class="desc">配置引入/升级时的自动化评估规则、等级划分阈值与卡点标准</p>

    <el-card style="margin-top: 16px">
      <template #header>等级划分阈值</template>
      <el-form label-width="120px">
        <el-form-item label="A级最低评分">
          <el-input-number v-model="thresholds.A" :min="0" :max="100" />
          <span class="form-tip">≥ {{ thresholds.A }} 分</span>
        </el-form-item>
        <el-form-item label="B级最低评分">
          <el-input-number v-model="thresholds.B" :min="0" :max="100" />
          <span class="form-tip">≥ {{ thresholds.B }} 分</span>
        </el-form-item>
        <el-form-item label="C级最低评分">
          <el-input-number v-model="thresholds.C" :min="0" :max="100" />
          <span class="form-tip">≥ {{ thresholds.C }} 分</span>
        </el-form-item>
        <el-form-item label="D级最低评分">
          <el-input-number v-model="thresholds.D" :min="0" :max="100" />
          <span class="form-tip">≥ {{ thresholds.D }} 分</span>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card style="margin-top: 16px">
      <template #header>卡点规则</template>
      <el-table :data="rules" border>
        <el-table-column prop="name" label="规则名称" width="200" />
        <el-table-column prop="condition" label="条件" min-width="200" />
        <el-table-column label="动作" width="120">
          <template #default="{ row }">
            <el-tag :type="row.action === 'reject' ? 'danger' : 'warning'" size="small">
              {{ row.action === 'reject' ? '自动驳回' : '警告' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="启用" width="80">
          <template #default="{ row }">
            <el-switch v-model="row.enabled" />
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <div style="margin-top: 16px">
      <el-button type="primary">保存配置</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const thresholds = ref({ A: 85, B: 70, C: 50, D: 0 })

const rules = ref([
  { name: '黑名单命中', condition: '组件或其依赖命中黑名单', action: 'reject', enabled: true },
  { name: '严重漏洞', condition: '存在 Critical 级别漏洞', action: 'reject', enabled: true },
  { name: '高危许可证', condition: '使用不认可的高风险许可证', action: 'reject', enabled: true },
  { name: '低活跃度', condition: '活跃度评分低于 30 分', action: 'warn', enabled: true },
  { name: '中风险许可证', condition: '使用中等风险许可证', action: 'warn', enabled: true },
])
</script>

<style scoped>
.desc { color: #909399; font-size: 13px; margin: 8px 0 0; }
.form-tip { margin-left: 12px; color: #409EFF; }
</style>
