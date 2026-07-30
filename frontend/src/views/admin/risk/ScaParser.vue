<template>
  <div class="sca-parser">
    <div class="page-header">
      <div class="header-left">
        <h3>SCA 自动化解析入库</h3>
        <p class="desc">手动应急通道：上传 SBOM 文件或源码包进行安全解析与合规校验</p>
      </div>
    </div>

    <el-row :gutter="16" style="margin-top: 16px; margin-bottom: 20px">
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="解析总数" :value="parseRecords.length" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="解析完成" :value="parseRecords.filter(r => r.status === '解析完成').length" value-style="color: #67C23A" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="解析中" :value="parseRecords.filter(r => r.status === '解析中').length" value-style="color: #409EFF" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="解析失败" :value="parseRecords.filter(r => r.status === '解析失败').length" value-style="color: #F56C6C" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 上传区域 -->
    <el-card style="margin-bottom: 20px">
      <template #header>
        <span style="font-weight: 600">手动应急通道入口</span>
      </template>
      <el-row :gutter="24">
        <el-col :span="12">
          <el-upload
            drag
            action="#"
            :auto-upload="false"
            :on-change="handleSbomChange"
            accept=".json,.xml"
          >
            <el-icon :size="48"><UploadFilled /></el-icon>
            <div class="el-upload__text">上传 SBOM 文件 <em>点击或拖拽</em></div>
            <template #tip>
              <div class="el-upload__tip">支持 SPDX (.xml) / CycloneDX (.json) 格式</div>
            </template>
          </el-upload>
          <div style="text-align: center; margin-top: 12px">
            <el-button type="primary" @click="handleParse('sbom')">开始安全解析</el-button>
          </div>
        </el-col>
        <el-col :span="12">
          <el-upload
            drag
            action="#"
            :auto-upload="false"
            :on-change="handleSourceChange"
            accept=".zip,.tar.gz"
          >
            <el-icon :size="48"><FolderOpened /></el-icon>
            <div class="el-upload__text">上传源码包 <em>点击或拖拽</em></div>
            <template #tip>
              <div class="el-upload__tip">支持 .zip / .tar.gz 源码压缩包</div>
            </template>
          </el-upload>
          <div style="text-align: center; margin-top: 12px">
            <el-button type="success" @click="handleParse('source')">开始合规校验</el-button>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 解析记录 -->
    <el-divider content-position="left">解析记录</el-divider>

    <el-table :data="parseRecords" stripe>
      <el-table-column prop="fileName" label="文件名" min-width="200" />
      <el-table-column prop="fileType" label="文件类型" width="150">
        <template #default="{ row }">
          <el-tag :type="fileTypeTag(row.fileType)" size="small">{{ row.fileType }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="解析状态" width="110">
        <template #default="{ row }">
          <el-tag :type="parseStatusType(row.status)" size="small">
            <el-icon v-if="row.status === '解析中'" class="is-loading"><Loading /></el-icon>
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="componentCount" label="组件数" width="90" />
      <el-table-column prop="riskCount" label="安全风险数" width="100">
        <template #default="{ row }">
          <span v-if="row.riskCount > 0" style="color: #F56C6C; font-weight: 600">{{ row.riskCount }}</span>
          <span v-else>{{ row.riskCount }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="licenseRiskCount" label="许可证风险" width="100">
        <template #default="{ row }">
          <span v-if="row.licenseRiskCount > 0" style="color: #E6A23C; font-weight: 600">{{ row.licenseRiskCount }}</span>
          <span v-else>{{ row.licenseRiskCount }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="uploadedBy" label="上传人" width="90" />
      <el-table-column prop="parseTime" label="上传时间" width="160" />
      <el-table-column prop="result" label="解析结果" min-width="200" show-overflow-tooltip />
      <el-table-column label="操作" width="140" fixed="right">
        <template #default="{ row }">
          <el-button v-if="(row as ScaParseRecord).status === '解析完成'" type="primary" link size="small" @click="viewResult(row as ScaParseRecord)">详情</el-button>
          <el-button v-if="(row as ScaParseRecord).status === '解析完成'" type="success" link size="small" @click="handle入库(row as ScaParseRecord)">入库</el-button>
          <el-button v-if="(row as ScaParseRecord).status === '解析失败'" type="danger" link size="small" @click="handleRetry(row as ScaParseRecord)">重试</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 解析结果弹窗 -->
    <el-dialog v-model="resultVisible" title="解析结果详情" width="700px">
      <el-descriptions :column="2" border v-if="currentRecord">
        <el-descriptions-item label="文件名">{{ currentRecord.fileName }}</el-descriptions-item>
        <el-descriptions-item label="文件类型">{{ currentRecord.fileType }}</el-descriptions-item>
        <el-descriptions-item label="组件总数">{{ currentRecord.componentCount }}</el-descriptions-item>
        <el-descriptions-item label="安全风险数" label-class-name="risk-label">{{ currentRecord.riskCount }}</el-descriptions-item>
        <el-descriptions-item label="许可证风险数">{{ currentRecord.licenseRiskCount }}</el-descriptions-item>
        <el-descriptions-item label="上传人">{{ currentRecord.uploadedBy }}</el-descriptions-item>
        <el-descriptions-item label="上传时间">{{ currentRecord.parseTime }}</el-descriptions-item>
        <el-descriptions-item label="入库状态">
          <el-tag type="success" size="small">已入库</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="解析结果" :span="2">{{ currentRecord.result }}</el-descriptions-item>
      </el-descriptions>

      <el-divider content-position="left">发现的风险组件</el-divider>
      <el-table :data="mockRiskComponents" size="small">
        <el-table-column prop="name" label="组件名称" width="140" />
        <el-table-column prop="version" label="版本" width="100" />
        <el-table-column prop="riskType" label="风险类型" width="110" />
        <el-table-column label="风险等级" width="90">
          <template #default="{ row }"><RiskTag :level="row.riskLevel" /></template>
        </el-table-column>
        <el-table-column prop="cveId" label="CVE编号" width="140" />
        <el-table-column prop="description" label="描述" min-width="180" show-overflow-tooltip />
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { mockScaParseRecords } from '@/mock/admin'
import { UploadFilled, FolderOpened, Loading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { ScaParseRecord } from '@/mock/admin'

const parseRecords = ref<ScaParseRecord[]>([...mockScaParseRecords])
const resultVisible = ref(false)
const currentRecord = ref<ScaParseRecord | null>(null)
const sbomFile = ref<File | null>(null)
const sourceFile = ref<File | null>(null)

const mockRiskComponents = [
  { name: 'log4j-core', version: '2.14.1', riskType: '安全漏洞', riskLevel: 'critical', cveId: 'CVE-2021-44228', description: 'JNDI注入导致远程代码执行' },
  { name: 'fastjson', version: '1.2.68', riskType: '安全漏洞', riskLevel: 'high', cveId: 'CVE-2022-2550', description: '反序列化远程代码执行' },
  { name: 'readline-sync', version: '1.4.10', riskType: '许可证风险', riskLevel: 'high', cveId: '-', description: '使用GPL-3.0许可证' },
]

function fileTypeTag(type: string): 'success' | 'warning' | undefined {
  const map: Record<string, 'success' | 'warning' | undefined> = { 'SBOM-SPDX': undefined, 'SBOM-CycloneDX': 'success', '源码包': 'warning' }
  return map[type]
}

function parseStatusType(status: string): 'success' | 'danger' | undefined {
  const map: Record<string, 'success' | 'danger' | undefined> = { '解析中': undefined, '解析完成': 'success', '解析失败': 'danger' }
  return map[status]
}

function handleSbomChange(file: any) {
  sbomFile.value = file.raw
}

function handleSourceChange(file: any) {
  sourceFile.value = file.raw
}

function handleParse(type: 'sbom' | 'source') {
  if (type === 'sbom' && !sbomFile.value) {
    ElMessage.warning('请先选择 SBOM 文件')
    return
  }
  if (type === 'source' && !sourceFile.value) {
    ElMessage.warning('请先选择源码包')
    return
  }

  const fileName = type === 'sbom' ? sbomFile.value?.name : sourceFile.value?.name
  const fileType = type === 'sbom' ? 'SBOM-CycloneDX' : '源码包'

  const newRecord: ScaParseRecord = {
    id: `sca${String(parseRecords.value.length + 1).padStart(3, '0')}`,
    fileName: fileName || 'unknown',
    fileType: fileType as any,
    status: '解析中',
    componentCount: 0,
    riskCount: 0,
    licenseRiskCount: 0,
    parseTime: new Date().toLocaleString('zh-CN'),
    uploadedBy: '王五',
    result: '正在解析中...',
  }
  parseRecords.value.unshift(newRecord)

  // 模拟解析完成
  setTimeout(() => {
    const idx = parseRecords.value.findIndex(r => r.id === newRecord.id)
    if (idx !== -1) {
      parseRecords.value[idx] = {
        ...parseRecords.value[idx],
        status: '解析完成',
        componentCount: Math.floor(Math.random() * 200) + 50,
        riskCount: Math.floor(Math.random() * 5),
        licenseRiskCount: Math.floor(Math.random() * 3),
        result: '解析完成，发现若干安全风险和许可证风险',
      }
      ElMessage.success(`${fileName} 解析完成`)
    }
  }, 3000)

  ElMessage.info(`开始解析 ${fileName}...`)
}

function viewResult(row: ScaParseRecord) {
  currentRecord.value = row
  resultVisible.value = true
}

function handle入库(row: ScaParseRecord) {
  ElMessage.success(`组件数据已入库，共 ${row.componentCount} 个组件`)
}

function handleRetry(row: ScaParseRecord) {
  row.status = '解析中'
  row.result = '重新解析中...'
  setTimeout(() => {
    row.status = '解析完成'
    row.componentCount = Math.floor(Math.random() * 200) + 50
    row.riskCount = Math.floor(Math.random() * 3)
    row.licenseRiskCount = Math.floor(Math.random() * 2)
    row.result = '重新解析完成'
    ElMessage.success(`${row.fileName} 重新解析完成`)
  }, 2000)
  ElMessage.info(`重新解析 ${row.fileName}...`)
}
</script>

<style scoped>
.page-header { margin-bottom: 8px; }
.page-header h3 { margin: 0 0 4px; }
.desc { color: #909399; font-size: 13px; margin: 0; }
:deep(.risk-label) { color: #F56C6C; font-weight: 600; }
</style>
