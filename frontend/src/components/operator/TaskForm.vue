<template>
  <div class="task-form">
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="120px"
      style="max-width: 700px"
      :disabled="readOnly"
    >
      <!-- 通用基础信息 -->
      <el-divider content-position="left">
        <el-icon><InfoFilled /></el-icon>
        基础信息
      </el-divider>

      <el-form-item label="任务标题" prop="title">
        <el-input v-model="formData.title" placeholder="请输入任务标题" />
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="申请人">
            <el-input v-model="formData.applicant" disabled />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="所属部门">
            <el-input v-model="formData.applicantDept" disabled />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="任务类型">
            <el-input :value="TASK_TYPE_LABELS[formData.type]" disabled />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="期望完成日期" prop="deadline" required>
            <el-date-picker
              v-model="formData.deadline"
              type="date"
              placeholder="请选择期望完成日期"
              style="width: 100%"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <!-- 引入申请表单 -->
      <template v-if="formData.type === 'introduce'">
        <el-divider content-position="left">
          <el-icon><Document /></el-icon>
          引入申请信息
        </el-divider>

        <el-form-item label="选择组件" prop="customComponent" required>
          <div style="display: flex; gap: 12px; align-items: flex-start">
            <el-input
              v-model="formData.customComponent"
              placeholder="请输入组件名称（如：Redis 7.2.4）"
              style="flex: 1"
              clearable
            />
            <el-select
              v-model="selectedComponents"
              placeholder="辅助评估"
              style="width: 280px"
              filterable
              multiple
              @change="onComponentSelect"
            >
              <el-option
                v-for="item in availableComponents"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              >
                <div style="display: flex; justify-content: space-between; align-items: center">
                  <span>{{ item.name }} {{ item.version }}</span>
                  <div>
                    <el-tag v-if="item.alreadyIntroduced" type="danger" size="small" style="margin-right: 4px">已引入</el-tag>
                    <el-tag :type="item.evalScore >= 80 ? 'success' : item.evalScore >= 60 ? 'warning' : 'danger'" size="small">
                      评估{{ item.evalScore }}分
                    </el-tag>
                  </div>
                </div>
              </el-option>
            </el-select>
          </div>
        </el-form-item>

        <!-- 单组件评估结果 -->
        <el-card v-if="selectedComponents.length === 1 && selectedComponentInfo" shadow="never" style="margin-bottom: 16px">
          <template #header>
            <span style="font-weight: 600">组件评估结果</span>
          </template>
          <el-table :data="componentEvalDimensions" border size="small">
            <el-table-column prop="dimension" label="评估维度" width="150" />
            <el-table-column prop="result" label="评估结果" width="120" align="center">
              <template #default="{ row }">
                <el-tag :type="row.result === '通过' ? 'success' : row.result === '警告' ? 'warning' : 'danger'" size="small">
                  {{ row.result }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="说明" />
          </el-table>
        </el-card>

        <!-- 多组件对比评估结果 -->
        <el-card v-if="selectedComponents.length > 1" shadow="never" style="margin-bottom: 16px">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span style="font-weight: 600">多组件对比评估</span>
              <el-tag v-if="isComparing" type="warning" size="small">
                <el-icon><Loading /></el-icon>
                智能评估中...
              </el-tag>
              <el-tag v-else-if="comparisonResult" type="success" size="small">
                评估完成
              </el-tag>
            </div>
          </template>

          <!-- 各组件评估详情 -->
          <el-table :data="comparisonDetails" border style="margin-bottom: 16px">
            <el-table-column prop="name" label="组件名称" width="180" />
            <el-table-column prop="version" label="版本" width="100" />
            <el-table-column prop="license" label="许可证" width="120" />
            <el-table-column prop="score" label="评估分数" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="row.score >= 80 ? 'success' : row.score >= 60 ? 'warning' : 'danger'" size="small">
                  {{ row.score }}分
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="licenseCheck" label="许可证合规" width="120" align="center">
              <template #default="{ row }">
                <el-tag :type="row.licenseCheck === '通过' ? 'success' : 'warning'" size="small">
                  {{ row.licenseCheck }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="blacklistCheck" label="黑名单匹配" width="120" align="center">
              <template #default="{ row }">
                <el-tag type="success" size="small">{{ row.blacklistCheck }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="activityCheck" label="活跃度分析" width="120" align="center">
              <template #default="{ row }">
                <el-tag :type="row.activityCheck === '通过' ? 'success' : row.activityCheck === '警告' ? 'warning' : 'danger'" size="small">
                  {{ row.activityCheck }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>

          <!-- 智能推荐结果 -->
          <el-alert v-if="comparisonResult" :type="comparisonResult.recommended ? 'success' : 'warning'" :closable="false">
            <template #title>
              <div style="font-weight: 600">智能推荐结果</div>
            </template>
            <template #default>
              <div style="margin-bottom: 8px">
                <strong>推荐组件：</strong>
                <el-tag type="success" size="large">{{ comparisonResult.recommendedName }}</el-tag>
              </div>
              <div>
                <strong>推荐理由：</strong>{{ comparisonResult.reason }}
              </div>
            </template>
          </el-alert>
        </el-card>

        <el-form-item label="引入原因" prop="introduceReason">
          <el-input
            v-model="formData.introduceReason"
            type="textarea"
            :rows="4"
            placeholder="请详细说明引入该组件的原因和业务场景"
          />
        </el-form-item>

        <el-divider content-position="left">适用性评估（必填）</el-divider>

        <el-alert type="warning" :closable="false" style="margin-bottom: 16px">
          <template #default>
            请根据实际情况评估以下8个方面，选择"存在偏差"时需填写偏差说明
          </template>
        </el-alert>

        <div v-for="item in assessItems" :key="item.key" class="assess-item">
          <el-form-item :label="item.label" :required="true">
            <el-radio-group v-model="assessResults[item.key].status" @change="onAssessChange(item.key)">
              <el-radio value="fully_meet">完全满足</el-radio>
              <el-radio value="deviation">存在偏差</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item v-if="assessResults[item.key].status === 'deviation'" label="偏差说明" :required="true" style="margin-top: -8px; margin-bottom: 16px">
            <el-input
              v-model="assessResults[item.key].deviationExplanation"
              type="textarea"
              :rows="3"
              placeholder="具体是哪个标准没达到？偏差原因是什么？计划如何规避风险？"
            />
          </el-form-item>
        </div>

        <el-divider content-position="left">备注说明（选填）</el-divider>

        <el-form-item label="附件上传">
          <el-upload
            v-model:file-list="formData.remarkFiles"
            :action="uploadUrl"
            :headers="uploadHeaders"
            :data="{ source: 'task_introduce' }"
            :limit="5"
            accept=".pdf,.doc,.docx,.xls,.xlsx,.png,.jpg,.jpeg"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
          >
            <el-button type="primary" size="small">
              <el-icon><Upload /></el-icon>
              上传附件
            </el-button>
            <template #tip>
              <div class="el-upload__tip">
                支持 PDF、Word、Excel、图片格式，最多5个文件，单个不超过10MB
              </div>
            </template>
          </el-upload>
        </el-form-item>

        <el-form-item label="备注说明">
          <el-input
            v-model="formData.wikiUrl"
            type="textarea"
            :rows="6"
            placeholder="请输入备注说明（选填）"
          />
        </el-form-item>
      </template>

      <!-- 升级申请表单 -->
      <template v-if="formData.type === 'upgrade'">
        <el-divider content-position="left">
          <el-icon><Upload /></el-icon>
          升级申请信息
        </el-divider>

        <el-form-item label="组件名称" prop="componentName">
          <el-input v-model="formData.componentName" placeholder="请输入组件名称" />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="当前版本" prop="currentVersion">
              <el-input v-model="formData.currentVersion" placeholder="请输入当前版本号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="目标版本" prop="targetVersion">
              <el-input v-model="formData.targetVersion" placeholder="请输入目标版本号" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="使用部门" prop="projectDept">
          <el-input v-model="formData.projectDept" placeholder="根据组件名称自动查询" disabled />
        </el-form-item>

        <el-form-item label="测试报告链接" prop="testReportLink">
          <el-input v-model="formData.testReportLink" placeholder="请输入测试报告URL地址" />
        </el-form-item>

        <el-form-item label="升级原因" prop="upgradeReason">
          <el-input
            v-model="formData.upgradeReason"
            type="textarea"
            :rows="4"
            placeholder="请说明升级原因和预期收益"
          />
        </el-form-item>
      </template>

      <!-- 退出申请表单 -->
      <template v-if="formData.type === 'withdraw'">
        <el-divider content-position="left">
          <el-icon><WarningFilled /></el-icon>
          退出申请信息
        </el-divider>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="拟退出组件" prop="exitComponent">
              <el-input v-model="formData.exitComponent" placeholder="请输入拟退出的组件名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="当前版本" prop="exitVersion">
              <el-input v-model="formData.exitVersion" placeholder="请输入当前版本号" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="退出原因" prop="withdrawReason">
          <el-input
            v-model="formData.withdrawReason"
            type="textarea"
            :rows="4"
            placeholder="请说明退出原因（如停止维护、安全漏洞等）"
          />
        </el-form-item>

        <el-form-item label="替代组件" prop="replacementComponent">
          <el-input
            v-model="formData.replacementComponent"
            placeholder="请输入替代组件名称（如有）"
          />
        </el-form-item>

        <el-form-item label="替代方案说明" prop="alternativePlan">
          <el-input
            v-model="formData.alternativePlan"
            type="textarea"
            :rows="3"
            placeholder="请说明替代方案"
          />
        </el-form-item>

        <el-form-item label="代码清理计划" prop="codeCleanupPlan">
          <el-input
            v-model="formData.codeCleanupPlan"
            type="textarea"
            :rows="3"
            placeholder="请说明代码清理计划"
          />
        </el-form-item>

        <el-form-item label="受影响项目" prop="affectedProjects">
          <el-input
            v-model="formData.affectedProjects"
            type="textarea"
            :rows="3"
            placeholder="请输入受影响项目清单"
          />
        </el-form-item>
      </template>

      <!-- 风险上报表单 -->
      <template v-if="formData.type === 'defect'">
        <el-divider content-position="left">
          <el-icon><Bell /></el-icon>
          风险上报信息
        </el-divider>

        <el-form-item label="内部缺陷编号">
          <el-input v-model="formData.internalDefectId" placeholder="如在使用时发现的问题编号" />
        </el-form-item>

        <el-form-item label="风险描述" prop="riskDescription">
          <el-input
            v-model="formData.riskDescription"
            type="textarea"
            :rows="6"
            placeholder="请详细描述发现的风险情况"
          />
        </el-form-item>
      </template>

    </el-form>

    <!-- 适用性评估 -->
    <div class="form-footer" v-if="!readOnly">
      <el-button @click="handleSaveDraftClick">保存草稿</el-button>
      <el-button type="primary" @click="handleSubmitOrShowReport">提交</el-button>
    </div>

    <!-- 引入评估报告对话框 -->
    <el-dialog
      v-model="evalReportDialogVisible"
      title="引入评估信息"
      width="800px"
    >
      <el-table :data="evalReportData" border>
        <el-table-column prop="label" label="评估项" width="180" />
        <el-table-column prop="value" label="评估结果">
          <template #default="{ row }">
            <template v-if="row.isEvalItem">
              <el-tag :type="row.value === '是' ? 'success' : 'danger'" size="small">
                {{ row.value }}
              </el-tag>
            </template>
            <template v-else>
              {{ row.value }}
            </template>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <el-button @click="evalReportDialogVisible = false">取消</el-button>
        <el-button type="success" @click="downloadReport">
          <el-icon><Download /></el-icon>
          下载
        </el-button>
        <el-button type="primary" @click="submitFromReport">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted, nextTick } from 'vue'
import type { FormInstance, FormRules, UploadFile, UploadFiles } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { InfoFilled, Document, Upload, WarningFilled, CircleClose, Check, View, Download, Loading } from '@element-plus/icons-vue'
import { TASK_TYPE_LABELS } from '@/types'
import type { TaskType, TaskStatus } from '@/types'
import { useUserStore } from '@/stores/user'

// 可选组件列表（从使用申请获取）
const availableComponents = ref([
  { id: 'c1', name: 'Redis', version: '7.2.4', license: 'BSD', source: '官方仓库', alreadyIntroduced: false, evalScore: 85, evalSummary: '高并发支持良好，高可用需配置集群，安全性优秀' },
  { id: 'c2', name: 'Spring Boot', version: '3.2.3', license: 'Apache-2.0', source: 'Maven中央仓库', alreadyIntroduced: true, evalScore: 92, evalSummary: '全项评估优秀，已引入组件' },
  { id: 'c3', name: 'MySQL', version: '8.0.36', license: 'GPL-2.0', source: '官方仓库', alreadyIntroduced: true, evalScore: 88, evalSummary: '性能稳定，已引入组件' },
  { id: 'c4', name: 'Nginx', version: '1.25.3', license: 'BSD', source: '官方仓库', alreadyIntroduced: false, evalScore: 90, evalSummary: '反向代理优秀，负载均衡能力强' },
  { id: 'c5', name: 'Kafka', version: '3.6.1', license: 'Apache-2.0', source: 'Apache仓库', alreadyIntroduced: false, evalScore: 82, evalSummary: '消息队列性能优秀' },
  { id: 'c6', name: 'Elasticsearch', version: '8.12.2', license: 'SSPL', source: '官方仓库', alreadyIntroduced: false, evalScore: 75, evalSummary: '搜索能力强，许可证需注意' },
])

// 适用性评估项定义（8个方面）
const assessItems = [
  { key: 'functionality', label: '功能特性' },
  { key: 'security', label: '安全性' },
  { key: 'reliability', label: '可靠性' },
  { key: 'performance', label: '性能效率' },
  { key: 'compatibility', label: '兼容性' },
  { key: 'maintainability', label: '可维护性' },
  { key: 'extensibility', label: '可扩展性' },
  { key: 'usability', label: '易用性' },
]

// 适用性评估结果
interface AssessResult {
  status: 'fully_meet' | 'deviation' | ''
  deviationExplanation?: string
}

const assessResults = ref<Record<string, AssessResult>>({})

// 初始化适用性评估结果
function initAssessResults() {
  assessItems.forEach(item => {
    if (!assessResults.value[item.key]) {
      assessResults.value[item.key] = {
        status: '',
        deviationExplanation: '',
      }
    }
  })
}

// 组件加载时初始化
initAssessResults()

function onAssessChange(key: string) {
  // 如果切换到完全满足，清空偏差说明
  if (assessResults.value[key].status === 'fully_meet') {
    assessResults.value[key].deviationExplanation = ''
  }
}

// Wiki地址
const wikiUrl = ref('')

// 附件列表
const fileList = ref<UploadFile[]>([])

// 已选择的组件（支持多选）
const selectedComponents = ref<string[]>([])
const selectedComponentInfo = ref<any>(null)

// 组件评估维度
const componentEvalDimensions = ref<any[]>([])

// 多组件对比相关
const isComparing = ref(false)
const comparisonResult = ref<any>(null)
const comparisonDetails = ref<any[]>([])

function onComponentSelect(componentIds: string[]) {
  selectedComponents.value = componentIds
  
  if (componentIds.length === 1) {
    // 单组件模式
    selectedComponentInfo.value = availableComponents.value.find(c => c.id === componentIds[0]) || null
    comparisonResult.value = null
    comparisonDetails.value = []
    
    // 生成评估维度
    if (selectedComponentInfo.value) {
      formData.selectedComponent = componentIds[0]
      formData.componentName = selectedComponentInfo.value.name
      formData.componentVersion = selectedComponentInfo.value.version
      formData.licenseType = selectedComponentInfo.value.license
      
      const score = selectedComponentInfo.value.evalScore
      componentEvalDimensions.value = [
        {
          dimension: '许可证合规',
          result: selectedComponentInfo.value.license === 'SSPL' ? '警告' : '通过',
          description: selectedComponentInfo.value.license === 'SSPL' ? 'SSPL许可证可能存在限制' : `${selectedComponentInfo.value.license}许可证合规`,
        },
        {
          dimension: '黑名单匹配',
          result: '通过',
          description: '未匹配到黑名单组件',
        },
        {
          dimension: '活跃度分析',
          result: score >= 80 ? '通过' : score >= 60 ? '警告' : '不通过',
          description: `评估分数${score}分，${score >= 80 ? '活跃度良好' : score >= 60 ? '活跃度一般' : '活跃度较低'}`,
        },
      ]
    }
  } else if (componentIds.length > 1) {
    // 多组件对比模式
    selectedComponentInfo.value = null
    componentEvalDimensions.value = []
    startComparison(componentIds)
  } else {
    // 未选择
    selectedComponentInfo.value = null
    comparisonResult.value = null
    comparisonDetails.value = []
    componentEvalDimensions.value = []
    formData.selectedComponent = ''
  }
}

// 开始多组件对比
async function startComparison(componentIds: string[]) {
  isComparing.value = true
  comparisonResult.value = null
  
  // 生成对比详情
  const components = componentIds.map(id => availableComponents.value.find(c => c.id === id)).filter(Boolean)
  
  comparisonDetails.value = components.map(comp => ({
    name: comp!.name,
    version: comp!.version,
    license: comp!.license,
    score: comp!.evalScore,
    licenseCheck: comp!.license === 'SSPL' ? '警告' : '通过',
    blacklistCheck: '通过',
    activityCheck: comp!.evalScore >= 80 ? '通过' : comp!.evalScore >= 60 ? '警告' : '不通过',
  }))
  
  // 模拟智能评估过程（2秒）
  await new Promise(resolve => setTimeout(resolve, 2000))
  
  // 智能推荐：选择分数最高的
  const best = components.reduce((prev, curr) => (curr!.evalScore > prev!.evalScore ? curr : prev), components[0])
  
  comparisonResult.value = {
    recommended: best!.evalScore >= 60,
    recommendedName: `${best!.name} ${best!.version}`,
    reason: `该组件评估分数${best!.evalScore}分，${best!.evalScore >= 80 ? '综合评估最优，社区活跃度高，建议引入' : best!.evalScore >= 60 ? '综合评估良好，可以满足基本需求' : '分数较低，建议谨慎评估'}`,
  }
  
  isComparing.value = false
  
  // 自动选择推荐的组件
  if (comparisonResult.value.recommended) {
    formData.selectedComponent = best!.id
  }
}

// 评估报告对话框
const evalReportDialogVisible = ref(false)

// 评估报告数据
const evalReportData = computed(() => {
  const component = selectedComponentInfo.value
  if (!component) return []
  
  const data = [
    { label: '组件名称', value: `${component.name} ${component.version}`, isEvalItem: false },
    { label: '许可证', value: component.license, isEvalItem: false },
    { label: '来源', value: component.source, isEvalItem: false },
    { label: '评估分数', value: `${component.evalScore}分`, isEvalItem: false },
    { label: '', value: '', isEvalItem: false }, // 空行
  ]
  
  // 添加评估维度
  componentEvalDimensions.value.forEach(d => {
    data.push({ label: `【${d.dimension}】`, value: `${d.result} - ${d.description}`, isEvalItem: false })
  })
  
  data.push({ label: '', value: '', isEvalItem: false }) // 空行
  
  // 添加适用性评估项
  assessItems.forEach(item => {
    const result = assessResults.value[item.key]
    const statusText = result?.status === 'fully_meet' ? '完全满足' : result?.status === 'deviation' ? '存在偏差' : '未评估'
    let value = statusText
    if (result?.status === 'deviation') {
      value += `\n偏差说明：${result.deviationExplanation || '未填写'}`
    }
    data.push({ 
      label: item.label, 
      value,
      isEvalItem: true 
    })
  })
  
  data.push({ label: '', value: '', isEvalItem: false }) // 空行
  data.push({ label: '引入原因', value: formData.introduceReason || '未填写', isEvalItem: false })
  
  return data
})

function showEvalReport() {
  if (!formData.selectedComponent) {
    ElMessage.warning('请先选择组件')
    return
  }
  evalReportDialogVisible.value = true
}

function submitFromReport() {
  // 如果是多组件模式且正在对比中，提示用户
  if (selectedComponents.value.length > 1 && isComparing.value) {
    ElMessageBox.confirm(
      '智能评估正在进行中，是否放弃对比？\n\n选择"是"将回退到单组件模式，选择"否"请等待评估完成。',
      '放弃对比',
      {
        confirmButtonText: '是',
        cancelButtonText: '否',
        type: 'warning',
      }
    ).then(() => {
      // 用户选择放弃，回退到单组件模式
      if (selectedComponents.value.length > 0) {
        selectedComponents.value = [selectedComponents.value[0]]
        onComponentSelect([selectedComponents.value[0]])
      }
      evalReportDialogVisible.value = false
      ElMessage.info('已回退到单组件模式，请重新选择')
    }).catch(() => {
      // 用户选择等待，不做任何操作
      ElMessage.info('请等待智能评估完成')
    })
    return
  }
  
  evalReportDialogVisible.value = false
  handleSubmit()
}

function downloadReport() {
  const component = selectedComponentInfo.value
  if (!component) return
  
  const reportContent = `
引入评估报告
==================

组件信息：
- 名称：${component.name}
- 版本：${component.version}
- 许可证：${component.license}
- 来源：${component.source}
- 评估分数：${component.evalScore}分

评估维度：
${componentEvalDimensions.value.map(d => `- ${d.dimension}: ${d.result}（${d.description}）`).join('\n')}

引入评估项：
${assessItems.map(item => {
  const result = assessResults.value[item.key]
  const statusText = result?.status === 'fully_meet' ? '完全满足' : result?.status === 'deviation' ? '存在偏差' : '未评估'
  return `- ${item.label}: ${statusText}`
}).join('\n')}

引入原因：
${formData.introduceReason || '未填写'}

生成时间：${new Date().toLocaleString()}
  `.trim()
  
  const blob = new Blob([reportContent], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `引入评估报告_${component.name}_${component.version}.txt`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
  
  ElMessage.success('报告下载成功')
}

interface TaskFormData {
  id: string
  taskNo: string
  type: TaskType
  title: string
  applicant: string
  applicantDept: string
  applicantId: string
  project: string
  projectDept: string  // 使用部门（用于审批流解析）
  status: TaskStatus
  priority: 'normal' | 'urgent'
  description: string
  remark: string
  // 引入申请
  customComponent: string  // 自定义输入的组件
  selectedComponent: string  // 选择的组件ID
  componentName: string
  componentVersion: string
  licenseType: string
  introduceReason: string
  // 备注说明（选填）
  remarkFiles: UploadFile[]
  wikiUrl: string
  // 升级申请
  currentVersion: string
  targetVersion: string
  upgradeReason: string  // 升级原因
  testReportLink: string
  // 退出申请
  exitComponent: string
  exitVersion: string
  exitReason: string
  withdrawReason: string
  replacementComponent: string
  alternativePlan: string
  codeCleanupPlan: string
  affectedProjects: string
  // 风险上报
  internalDefectId: string  // 内部缺陷编号
  riskDescription: string
  // 缺陷处置
  defectRiskSource: string
  defectCveId: string
  defectSeverityLevel: string
  defectAffectedAssets: string
  defectFixPlan: string
  needFix: boolean | undefined
  fixReason: string
  productionChangeNo: string
  productionTicketText: string
  // 风险来源/严重等级（复选框临时绑定）
  riskSource: string
  severityLevel: string
  // 缺陷修复计划
  fixPlan: string
  // 期望完成日期
  deadline: string
}

const userStore = useUserStore()

const props = withDefaults(defineProps<{
  taskData?: Partial<TaskFormData>
  readOnly?: boolean
  taskType?: TaskType
}>(), {
  readOnly: false,
})

const emit = defineEmits<{
  'save-draft': [data: TaskFormData]
  'submit': [data: TaskFormData]
}>()

// 获取当前登录用户信息作为默认值
const getCurrentUserInfo = () => {
  const user = userStore.userInfo
  return {
    applicant: user?.nickname || '',
    applicantDept: user?.department || '',
    applicantId: user?.id || '',
  }
}

function handleSaveDraftClick() {
  const data: any = { ...formData }
  data.assessResults = assessResults.value
  emit('save-draft', data)
}

const formRef = ref<FormInstance>()

// 初始化表单数据
const initFormData = (): TaskFormData => {
  const currentUser = getCurrentUserInfo()
  return {
    id: props.taskData?.id || '',
    taskNo: props.taskData?.taskNo || '',
    type: props.taskType || props.taskData?.type || 'introduce',
    title: props.taskData?.title || '',
    applicant: props.taskData?.applicant || currentUser.applicant,
    applicantDept: props.taskData?.applicantDept || currentUser.applicantDept,
    applicantId: props.taskData?.applicantId || currentUser.applicantId,
    project: props.taskData?.project || '',
    projectDept: props.taskData?.projectDept || '',
    status: props.taskData?.status || 'draft',
    priority: props.taskData?.priority || 'normal',
    deadline: props.taskData?.deadline || '',
    description: props.taskData?.description || '',
    remark: props.taskData?.remark || '',
    selectedComponent: props.taskData?.selectedComponent || '',
    customComponent: props.taskData?.customComponent || '',
    componentName: props.taskData?.componentName || '',
    componentVersion: props.taskData?.componentVersion || '',
    licenseType: props.taskData?.licenseType || '',
    introduceReason: props.taskData?.introduceReason || '',
    remarkFiles: props.taskData?.remarkFiles || [],
    wikiUrl: props.taskData?.wikiUrl || '',
    currentVersion: props.taskData?.currentVersion || '',
    targetVersion: props.taskData?.targetVersion || '',
    upgradeReason: props.taskData?.upgradeReason || '',
    testReportLink: props.taskData?.testReportLink || '',
    exitComponent: props.taskData?.exitComponent || props.taskData?.componentName || '',
    exitVersion: props.taskData?.exitVersion || props.taskData?.componentVersion || '',
    exitReason: props.taskData?.exitReason || props.taskData?.withdrawReason || '',
    withdrawReason: props.taskData?.withdrawReason || props.taskData?.exitReason || '',
    replacementComponent: props.taskData?.replacementComponent || '',
    alternativePlan: props.taskData?.alternativePlan || '',
    codeCleanupPlan: props.taskData?.codeCleanupPlan || '',
    affectedProjects: props.taskData?.affectedProjects || '',
    internalDefectId: props.taskData?.internalDefectId || '',
    riskDescription: props.taskData?.riskDescription || '',
    defectRiskSource: props.taskData?.defectRiskSource || '',
    defectCveId: props.taskData?.defectCveId || '',
    defectSeverityLevel: props.taskData?.defectSeverityLevel || '',
    defectAffectedAssets: props.taskData?.defectAffectedAssets || '',
    defectFixPlan: props.taskData?.defectFixPlan || '',
    needFix: props.taskData?.needFix,
    fixReason: props.taskData?.fixReason || '',
    productionChangeNo: props.taskData?.productionChangeNo || '',
    productionTicketText: props.taskData?.productionTicketText || '',
    riskSource: props.taskData?.riskSource || '',
    severityLevel: props.taskData?.severityLevel || '',
    fixPlan: props.taskData?.fixPlan || '',
  }
}

const formData = reactive<TaskFormData>(initFormData())

// 重置表单为当前用户的新建状态
function resetFormToNew() {
  const currentUser = getCurrentUserInfo()
  Object.assign(formData, {
    id: '',
    taskNo: '',
    type: props.taskType || 'introduce',
    title: '',
    applicant: currentUser.applicant,
    applicantDept: currentUser.applicantDept,
    applicantId: currentUser.applicantId,
    project: '',
    projectDept: '',
    status: 'draft',
    priority: 'normal',
    deadline: '',
    description: '',
    remark: '',
    selectedComponent: '',
    customComponent: '',
    componentName: '',
    componentVersion: '',
    licenseType: '',
    introduceReason: '',
    remarkFiles: [],
    wikiUrl: '',
    currentVersion: '',
    targetVersion: '',
    upgradeReason: '',
    testReportLink: '',
    exitComponent: '',
    exitVersion: '',
    exitReason: '',
    withdrawReason: '',
    replacementComponent: '',
    alternativePlan: '',
    codeCleanupPlan: '',
    affectedProjects: '',
    internalDefectId: '',
    riskDescription: '',
    defectRiskSource: '',
    defectCveId: '',
    defectSeverityLevel: '',
    defectAffectedAssets: '',
    defectFixPlan: '',
    needFix: undefined,
    fixReason: '',
    productionChangeNo: '',
    productionTicketText: '',
    riskSource: '',
    severityLevel: '',
    fixPlan: '',
  })
}

// 填充表单数据
function fillFormData(data: any) {
  console.log('[fillFormData] ========== 开始填充 ==========')
  console.log('[fillFormData] 传入的data:', data)
  console.log('[fillFormData] data.currentVersion:', data.currentVersion)
  console.log('[fillFormData] data.targetVersion:', data.targetVersion)
  console.log('[fillFormData] data.testReportLink:', data.testReportLink)
  console.log('[fillFormData] data.upgradeReason:', data.upgradeReason)
  console.log('[fillFormData] data.customComponent:', data.customComponent)
  console.log('[fillFormData] data.introduceReason:', data.introduceReason)
  console.log('[fillFormData] data.componentName:', data.componentName)
  console.log('[fillFormData] data.deadline:', data.deadline)
  console.log('[fillFormData] data.remark:', data.remark)
  
  if (!data || !data.id) {
    console.log('[fillFormData] data或data.id为空，跳过填充')
    return
  }
  
  console.log('[fillFormData] 开始填充formData...')
  Object.assign(formData, {
    id: data.id || '',
    taskNo: data.taskNo || '',
    type: data.type || 'introduce',
    title: data.title || '',
    applicant: data.applicant || '',
    applicantDept: data.applicantDept || '',
    project: data.project || '',
    projectDept: data.projectDept || '',
    status: data.status || 'draft',
    priority: data.priority || 'normal',
    deadline: data.deadline || '',
    description: data.description || '',
    remark: data.remark || '',
    selectedComponent: data.selectedComponent || '',
    customComponent: data.customComponent || '',
    componentName: data.componentName || '',
    componentVersion: data.componentVersion || '',
    licenseType: data.licenseType || '',
    introduceReason: data.introduceReason || '',
    upgradeReason: data.upgradeReason || '',
    remarkFiles: data.remarkFiles || [],
    wikiUrl: data.wikiUrl || '',
    currentVersion: data.currentVersion || '',
    targetVersion: data.targetVersion || '',
    testReportLink: data.testReportLink || '',
    exitComponent: data.exitComponent || data.componentName || '',
    exitVersion: data.exitVersion || data.componentVersion || '',
    exitReason: data.exitReason || data.withdrawReason || '',
    withdrawReason: data.withdrawReason || data.exitReason || '',
    replacementComponent: data.replacementComponent || '',
    alternativePlan: data.alternativePlan || '',
    codeCleanupPlan: data.codeCleanupPlan || '',
    affectedProjects: data.affectedProjects || '',
    internalDefectId: data.internalDefectId || '',
    riskDescription: data.riskDescription || '',
    defectRiskSource: data.defectRiskSource || '',
    defectCveId: data.defectCveId || '',
    defectSeverityLevel: data.defectSeverityLevel || '',
    defectAffectedAssets: data.defectAffectedAssets || '',
    defectFixPlan: data.defectFixPlan || '',
    needFix: data.needFix,
    fixReason: data.fixReason || '',
    productionChangeNo: data.productionChangeNo || '',
    productionTicketText: data.productionTicketText || '',
    riskSource: data.riskSource || '',
    severityLevel: data.severityLevel || '',
    fixPlan: data.fixPlan || '',
  })
  
  console.log('[fillFormData] 填充后的formData.customComponent:', formData.customComponent)
  console.log('[fillFormData] 填充后的formData.introduceReason:', formData.introduceReason)
  console.log('[fillFormData] 填充后的formData.deadline:', formData.deadline)
  console.log('[fillFormData] 填充后的formData.remark:', formData.remark)
  console.log('[fillFormData] 填充后的formData.currentVersion:', formData.currentVersion)
  console.log('[fillFormData] 填充后的formData.targetVersion:', formData.targetVersion)
  console.log('[fillFormData] 填充后的formData.testReportLink:', formData.testReportLink)
  console.log('[fillFormData] 填充后的formData.upgradeReason:', formData.upgradeReason)
  // 同步适用性评估数据
  console.log('[fillFormData] 准备同步assessResults')
  if (data.assessResults) {
    console.log('[fillFormData] 找到assessResults，开始合并:', data.assessResults)
    Object.assign(assessResults.value, data.assessResults)
    console.log('[fillFormData] 合并后assessResults.value:', JSON.parse(JSON.stringify(assessResults.value)))
  }
  // 从后端返回的applicabilityAssessment转换为前端格式
  if (data.applicabilityAssessment) {
    console.log('[fillFormData] 找到applicabilityAssessment，开始转换')
    const dimensions = [
      'functionality', 'security', 'reliability', 'performance',
      'compatibility', 'maintainability', 'extensibility', 'usability'
    ]
    dimensions.forEach(dim => {
      const backendKey = `${dim}Status`
      const deviationKey = `${dim}Deviation`
      if (data.applicabilityAssessment[backendKey] || data.applicabilityAssessment[deviationKey]) {
        assessResults.value[dim] = {
          status: data.applicabilityAssessment[backendKey] || '',
          deviationExplanation: data.applicabilityAssessment[deviationKey] || ''
        }
      }
    })
    console.log('[fillFormData] 转换后assessResults.value:', JSON.parse(JSON.stringify(assessResults.value)))
  }
  console.log('[fillFormData] 填充完成')
  // 清除表单校验状态，避免旧校验错误影响编辑
  nextTick(() => {
    if (formRef.value) {
      formRef.value.clearValidate()
    }
  })
}

// 监听taskData变化，填充表单
watch(() => props.taskData, (newData) => {
  if (newData && newData.id) {
    fillFormData(newData)
  } else {
    // 新建任务：重置为当前用户信息
    resetFormToNew()
  }
}, { deep: true, immediate: true })

// 组件挂载时填充数据
onMounted(() => {
  if (props.taskData && props.taskData.id) {
    fillFormData(props.taskData)
  } else {
    resetFormToNew()
  }
})

// 文件上传相关
const uploadUrl = '/api/v1/files/upload'  // 后端上传接口
const uploadHeaders = {}

function beforeUpload(file: File) {
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB!')
  }
  return isLt10M
}

function handleUploadSuccess(response: any, file: UploadFile) {
  ElMessage.success('文件上传成功')
}

function handleUploadError(error: any) {
  ElMessage.error('文件上传失败')
}

// 风险来源复选框
const riskSourceCheckboxes = ref<string[]>([])
// 严重等级复选框
const severityCheckboxes = ref<string[]>([])

// 同步复选框到表单数据
watch(riskSourceCheckboxes, (val) => {
  formData.riskSource = val.join('/')
})

watch(severityCheckboxes, (val) => {
  formData.severityLevel = val.join('/') as any
})

// 动态验证规则
const formRules = computed<FormRules>(() => {
  const baseRules: FormRules = {
    title: [{ required: true, message: '请输入任务标题', trigger: 'blur' }],
  }

  if (formData.type === 'introduce') {
    Object.assign(baseRules, {
      customComponent: [{ required: true, message: '请输入组件名称', trigger: 'blur' }],
      introduceReason: [{ required: true, message: '请输入引入原因', trigger: 'blur' }],
    })
  }

  if (formData.type === 'upgrade') {
    Object.assign(baseRules, {
      componentName: [{ required: true, message: '请输入组件名称', trigger: 'blur' }],
      currentVersion: [{ required: true, message: '请输入当前版本', trigger: 'blur' }],
      targetVersion: [{ required: true, message: '请输入目标版本', trigger: 'blur' }],
      testReportLink: [{ required: true, message: '请输入测试报告链接', trigger: 'blur' }],
      upgradeReason: [{ required: true, message: '请输入升级原因', trigger: 'blur' }],
    })
  }

  if (formData.type === 'withdraw') {
    Object.assign(baseRules, {
      exitComponent: [{ required: true, message: '请输入拟退出组件', trigger: 'blur' }],
      exitVersion: [{ required: true, message: '请输入当前版本', trigger: 'blur' }],
      withdrawReason: [{ required: true, message: '请输入退出原因', trigger: 'blur' }],
    })
  }

  if (formData.type === 'defect') {
    Object.assign(baseRules, {
      riskDescription: [{ required: true, message: '请输入风险描述', trigger: 'blur' }],
    })
  }

  return baseRules
})

function onNeedFixChange(val: boolean) {
  if (val) {
    formData.fixReason = ''
  } else {
    formData.fixPlan = ''
    formData.productionChangeNo = ''
    formData.productionTicketText = ''
  }
}

// 暂存
async function handleSaveDraft() {
  const data: any = { ...formData }
  data.assessResults = assessResults.value
  emit('save-draft', data)
}

// 提交
async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    
    // 引入申请需要验证适用性评估
    if (formData.type === 'introduce') {
      // 1. 检查所有评估项是否都已选择状态
      const missingStatus = assessItems.find(item => {
        const result = assessResults.value[item.key]
        return !result?.status
      })
      
      if (missingStatus) {
        ElMessage.warning(`请选择【${missingStatus.label}】的评估状态`)
        return
      }
      
      // 2. 检查存在偏差的项是否填写了偏差说明
      const missingExplanation = assessItems.find(item => {
        const result = assessResults.value[item.key]
        return result?.status === 'deviation' && !result?.deviationExplanation?.trim()
      })
      
      if (missingExplanation) {
        ElMessage.warning(`请填写【${missingExplanation.label}】的偏差说明`)
        return
      }
    }
    
    // 引入申请需要显示安全初评提示
    if (formData.type === 'introduce') {
      // 模拟安全初评过程
      await new Promise(resolve => setTimeout(resolve, 500))
    }
    
    const submitData: any = { ...formData }
    submitData.assessResults = assessResults.value
    emit('submit', submitData)
  } catch (error) {
    // 验证失败，显示错误信息
    console.log('[handleSubmit] 表单校验失败:', error)
    if (error && typeof error === 'object') {
      const firstError = Object.values(error)[0]
      if (firstError && firstError[0]) {
        ElMessage.warning(firstError[0].message)
      }
    }
  }
}

// 提交或显示评估报告
function handleSubmitOrShowReport() {
  // 所有类型都先提交
  handleSubmit()
}

defineExpose({
  formData,
  handleSaveDraft,
  handleSubmit,
  formRef,
  showEvalReport,
})
</script>

<style scoped>
.task-form {
  padding: 16px 0;
}

.form-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 0 0;
  border-top: 1px solid #ebeef5;
  margin-top: 20px;
}

:deep(.el-divider__text) {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  color: #303133;
}

/* 适用性评估项样式 */
.assess-item {
  margin-bottom: 16px;
}

.assess-item :deep(.el-form-item__label) {
  width: 140px !important;
}
</style>
