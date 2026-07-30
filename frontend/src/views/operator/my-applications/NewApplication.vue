<template>
  <div class="new-application">
    <el-page-header @back="goBack" content="发起新申请" />

    <div style="margin-top: 16px">
      <!-- Step 1: Select application type -->
      <el-steps :active="currentStep" finish-status="success" style="margin-bottom: 24px">
        <el-step title="选择类型" />
        <el-step title="填写信息" />
        <el-step title="预评估" />
        <el-step title="确认提交" />
      </el-steps>

      <!-- Type Selection -->
      <div v-if="currentStep === 0" class="step-content">
        <h3>请选择申请类型</h3>
        <el-row :gutter="16" style="margin-top: 20px">
          <el-col :span="6" v-for="type in applicationTypes" :key="type.value">
            <el-card
              shadow="hover"
              class="type-card"
              :class="{ selected: form.type === type.value }"
              @click="form.type = type.value"
            >
              <el-icon :size="32" :color="type.color"><component :is="type.icon" /></el-icon>
              <h4 style="margin: 12px 0 4px">{{ type.label }}</h4>
              <p class="type-desc">{{ type.desc }}</p>
            </el-card>
          </el-col>
        </el-row>
        <div class="step-actions">
          <el-button type="primary" :disabled="!form.type" @click="currentStep = 1">下一步</el-button>
        </div>
      </div>

      <!-- Fill form -->
      <div v-if="currentStep === 1" class="step-content">
        <h3>填写申请信息</h3>
        <el-form :model="form" label-width="120px" style="max-width: 600px; margin-top: 20px">
          <el-form-item label="组件名称" required>
            <el-input v-model="form.componentName" placeholder="请输入组件名称" />
          </el-form-item>
          <el-form-item label="组件版本" required>
            <el-input v-model="form.componentVersion" placeholder="请输入版本号" />
          </el-form-item>
          <el-form-item label="来源URL">
            <el-input v-model="form.sourceUrl" placeholder="请输入组件来源地址" />
          </el-form-item>
          <el-form-item label="使用场景" required>
            <el-input v-model="form.businessScenario" type="textarea" :rows="3" placeholder="请描述业务系统名称和用途" />
          </el-form-item>
          <el-form-item label="预期等级">
            <el-select v-model="form.expectedGrade" placeholder="请选择">
              <el-option label="A - 核心关键" value="A" />
              <el-option label="B - 重要支撑" value="B" />
              <el-option label="C - 一般使用" value="C" />
              <el-option label="D - 探索试用" value="D" />
            </el-select>
          </el-form-item>
          <el-form-item label="紧急程度">
            <el-radio-group v-model="form.urgency">
              <el-radio label="normal">普通</el-radio>
              <el-radio label="urgent">紧急</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="申请说明" required>
            <el-input v-model="form.description" type="textarea" :rows="4" placeholder="请详细说明申请原因" />
          </el-form-item>
          <el-form-item label="附件">
            <el-upload action="#" :auto-upload="false" :limit="5">
              <el-button type="primary">上传文件</el-button>
              <template #tip>
                <div class="el-upload__tip">支持上传技术方案、安全评估报告等文件</div>
              </template>
            </el-upload>
          </el-form-item>
        </el-form>
        <div class="step-actions">
          <el-button @click="currentStep = 0">上一步</el-button>
          <el-button type="primary" @click="startPreEval">开始预评估</el-button>
        </div>
      </div>

      <!-- Pre-evaluation -->
      <div v-if="currentStep === 2" class="step-content">
        <h3>预评估进行中</h3>
        <PreEvalProgress
          :status="preEvalStatus"
          :conclusion="preEvalConclusion"
          :score="preEvalScore"
        />
        <div class="step-actions" v-if="preEvalStatus === 'completed'">
          <el-button @click="currentStep = 1">返回修改</el-button>
          <el-button type="primary" @click="confirmSubmit">确认提交</el-button>
        </div>
      </div>

      <!-- Confirm submit -->
      <div v-if="currentStep === 3" class="step-content">
        <el-result icon="success" title="申请已提交" :sub-title="`工单号: ${ticketNo}`">
          <template #extra>
            <el-form inline style="margin-top: 12px">
              <el-form-item label="工单链接">
                <el-input v-model="externalLink" placeholder="请输入外部工单链接" style="width: 300px" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="goBack">完成</el-button>
              </el-form-item>
            </el-form>
          </template>
        </el-result>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Document, Upload, WarningFilled, CircleClose } from '@element-plus/icons-vue'
import type { ApplicationType, AssetGrade, PreEvalStatus, PreEvalConclusion } from '@/types'

const route = useRoute()
const router = useRouter()

const currentStep = ref(0)
const preEvalStatus = ref<PreEvalStatus>('pending')
const preEvalConclusion = ref<PreEvalConclusion>()
const preEvalScore = ref<number>()
const ticketNo = ref('')
const externalLink = ref('')

const form = reactive({
  type: '' as ApplicationType | '',
  componentName: (route.query.name as string) || '',
  componentVersion: (route.query.version as string) || '',
  sourceUrl: '',
  businessScenario: '',
  expectedGrade: '' as AssetGrade | '',
  urgency: 'normal' as 'normal' | 'urgent',
  description: '',
})

const applicationTypes = [
  { value: 'introduce' as ApplicationType, label: '引入申请', desc: '申请引入新的开源组件', icon: Document, color: '#409EFF' },
  { value: 'upgrade' as ApplicationType, label: '升级申请', desc: '申请升级已有组件版本', icon: Upload, color: '#67C23A' },
  { value: 'withdraw' as ApplicationType, label: '退出申请', desc: '申请退出不再使用的组件', icon: WarningFilled, color: '#E6A23C' },
  { value: 'risk_report' as ApplicationType, label: '风险上报', desc: '上报发现的安全风险', icon: CircleClose, color: '#F56C6C' },
]

onMounted(() => {
  if (route.query.name) {
    form.type = 'introduce'
  }
})

function startPreEval() {
  if (!form.componentName || !form.description) {
    ElMessage.warning('请填写必填项')
    return
  }
  currentStep.value = 2
  preEvalStatus.value = 'scanning'

  // Simulate pre-evaluation progress
  setTimeout(() => {
    preEvalStatus.value = 'completed'
    preEvalScore.value = Math.floor(Math.random() * 30) + 65
    if (preEvalScore.value >= 80) {
      preEvalConclusion.value = 'pass'
    } else if (preEvalScore.value >= 60) {
      preEvalConclusion.value = 'warning'
    } else {
      preEvalConclusion.value = 'reject'
    }
  }, 2500)
}

function confirmSubmit() {
  ticketNo.value = 'TK-2024-' + String(Math.floor(Math.random() * 9000) + 1000)
  currentStep.value = 3
  ElMessage.success('申请已提交，请填写外部工单链接')
}

function goBack() {
  router.push('/operator/my-applications')
}
</script>

<style scoped>
.step-content {
  padding: 20px 0;
}

.step-content h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 0;
}

.type-card {
  cursor: pointer;
  text-align: center;
  transition: all 0.3s;
  border: 2px solid transparent;
}

.type-card:hover,
.type-card.selected {
  border-color: #409EFF;
}

.type-desc {
  font-size: 12px;
  color: #909399;
  margin: 0;
}

.step-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 24px;
}
</style>
