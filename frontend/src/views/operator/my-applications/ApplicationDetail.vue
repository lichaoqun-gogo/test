<template>
  <div class="application-detail">
    <el-page-header @back="goBack" content="申请详情" />

    <el-descriptions :column="2" border style="margin-top: 16px">
      <el-descriptions-item label="申请编号">{{ application?.applicationNo }}</el-descriptions-item>
      <el-descriptions-item label="申请类型">
        {{ APPLICATION_TYPE_LABELS[application?.type as keyof typeof APPLICATION_TYPE_LABELS] }}
      </el-descriptions-item>
      <el-descriptions-item label="组件名称">{{ application?.componentName }}</el-descriptions-item>
      <el-descriptions-item label="组件版本">{{ application?.componentVersion }}</el-descriptions-item>
      <el-descriptions-item label="申请人">{{ application?.applicant }}</el-descriptions-item>
      <el-descriptions-item label="部门">{{ application?.department }}</el-descriptions-item>
      <el-descriptions-item label="预期等级">{{ application?.expectedGrade }}</el-descriptions-item>
      <el-descriptions-item label="紧急程度">
        <el-tag :type="application?.urgency === 'urgent' ? 'danger' : ''" size="small">
          {{ application?.urgency === 'urgent' ? '紧急' : '普通' }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="工单号">{{ application?.ticketNo || '未生成' }}</el-descriptions-item>
      <el-descriptions-item label="创建时间">{{ formatTime(application?.createdAt) }}</el-descriptions-item>
      <el-descriptions-item label="状态" :span="2">
        <TicketStatusTag :status="application?.status || 'pending'" />
      </el-descriptions-item>
      <el-descriptions-item label="使用场景" :span="2">{{ application?.businessScenario }}</el-descriptions-item>
      <el-descriptions-item label="申请说明" :span="2">{{ application?.description }}</el-descriptions-item>
    </el-descriptions>

    <el-divider content-position="left">预评估报告</el-divider>

    <PreEvalProgress
      :status="application?.preEvalStatus || 'pending'"
      :conclusion="application?.preEvalConclusion"
      :score="application?.preEvalScore"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { APPLICATION_TYPE_LABELS } from '@/types'
import { formatTime } from '@/utils/format'

const route = useRoute()
const router = useRouter()

const application = ref<any>(null)

function loadApplication(id: string) {
  // TODO: 调用后端 API 获取申请详情
}

onMounted(() => {
  loadApplication(route.params.id as string)
})

function goBack() {
  router.back()
}
</script>
