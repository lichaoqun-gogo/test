<template>
  <div class="approval-detail">
    <el-page-header @back="goBack" content="审批详情" />

    <el-descriptions :column="2" border style="margin-top: 16px">
      <el-descriptions-item label="申请编号">{{ approval?.applicationNo }}</el-descriptions-item>
      <el-descriptions-item label="当前步骤">{{ approval?.currentStep }}</el-descriptions-item>
      <el-descriptions-item label="组件">{{ approval?.componentName }} v{{ approval?.componentVersion }}</el-descriptions-item>
      <el-descriptions-item label="申请人">{{ approval?.applicant }} ({{ approval?.department }})</el-descriptions-item>
      <el-descriptions-item label="预评估评分">
        <el-tag :type="approval && approval.preEvalScore >= 80 ? 'success' : approval && approval.preEvalScore >= 60 ? 'warning' : 'danger'">
          {{ approval?.preEvalScore }}分
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="预评估结论">{{ approval?.preEvalConclusion }}</el-descriptions-item>
      <el-descriptions-item label="提交时间">{{ formatTime(approval?.submittedAt) }}</el-descriptions-item>
      <el-descriptions-item label="状态">
        <TicketStatusTag :status="approval?.status || 'pending'" />
      </el-descriptions-item>
    </el-descriptions>

    <el-divider content-position="left">业务终评</el-divider>

    <el-card>
      <el-form label-width="140px">
        <el-form-item label="业务运行指标匹配度">
          <el-radio-group>
            <el-radio label="match">完全匹配</el-radio>
            <el-radio label="partial">部分匹配</el-radio>
            <el-radio label="unmatch">不匹配</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="部门适用性评估">
          <el-radio-group>
            <el-radio label="apply">适用</el-radio>
            <el-radio label="conditional">有条件适用</el-radio>
            <el-radio label="unapply">不适用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="最终风险确认">
          <el-select placeholder="请选择" style="width: 200px">
            <el-option label="低风险" value="low" />
            <el-option label="中风险" value="medium" />
            <el-option label="高风险" value="high" />
          </el-select>
        </el-form-item>
        <el-form-item label="最终等级核定">
          <el-select placeholder="系统推荐等级" style="width: 200px">
            <el-option label="A - 核心关键" value="A" />
            <el-option label="B - 重要支撑" value="B" />
            <el-option label="C - 一般使用" value="C" />
            <el-option label="D - 探索试用" value="D" />
          </el-select>
        </el-form-item>
        <el-form-item label="补充意见">
          <el-input type="textarea" :rows="3" placeholder="请输入审批意见" />
        </el-form-item>
      </el-form>

      <div style="margin-top: 16px">
        <el-button type="success" @click="handleApprove">审批通过</el-button>
        <el-button type="danger">驳回</el-button>
        <el-button type="warning">要求补充材料</el-button>
        <el-button type="info" plain><el-icon><Download /></el-icon> 生成Word文档</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { mockIntroductions } from '@/mock/admin'
import { Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { formatTime } from '@/utils/format'

const route = useRoute()
const router = useRouter()

const approval = computed(() => mockIntroductions.find(i => i.id === route.params.id))

function goBack() {
  router.back()
}

function handleApprove() {
  if (approval.value) {
    approval.value.status = 'completed'
    approval.value.currentStep = '审批通过'
  }
  ElMessage.success('审批通过')
}
</script>
