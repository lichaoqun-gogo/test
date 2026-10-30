<template>
  <div class="defect-ledger">
    <div class="page-header">
      <div class="header-left">
        <h3>缺陷台账</h3>
        <p class="desc">内部缺陷（风险上报，管理员审核分配）与外部缺陷（采集导入）分开展示，支持研判、评审与处置</p>
      </div>
      <el-button type="primary" @click="openCreateDialog">
        <el-icon><Plus /></el-icon> 上报风险
      </el-button>
      <el-button type="primary" plain @click="openVulnIntelDialog" style="margin-left: 8px">
        <el-icon><Files /></el-icon> 导入漏洞情报
      </el-button>
    </div>

    <!-- 统计卡片（内部 + 外部合计） -->
    <el-row :gutter="16" style="margin-top: 16px; margin-bottom: 20px">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #fff3e0; color: #ff9800">
              <el-icon :size="28"><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value" style="color: #ff9800">{{ countByStatus('pending_analysis') }}</div>
              <div class="stat-label">待研判</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #e3f2fd; color: #2196f3">
              <el-icon :size="28"><Tickets /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value" style="color: #2196f3">{{ countByStatus('pending_review') }}</div>
              <div class="stat-label">待上会/待审核</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #f0f9ff; color: #409EFF">
              <el-icon :size="28"><Operation /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value" style="color: #409EFF">{{ countByStatus('assigned') }}</div>
              <div class="stat-label">待处置</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #f0f9eb; color: #67c23a">
              <el-icon :size="28"><CircleCheck /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value" style="color: #67c23a">{{ countByStatus('closed') }}</div>
              <div class="stat-label">已闭环</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-tabs v-model="activeTab" style="margin-top: 8px">
      <!-- ==================== 内部缺陷 tab ==================== -->
      <el-tab-pane label="内部缺陷" name="internal">
        <div class="toolbar">
          <el-input v-model="internalFilters.search" placeholder="搜索标题" clearable style="width: 240px" @clear="loadInternal" @keyup.enter="loadInternal">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
          <el-select v-model="internalFilters.severity" placeholder="风险等级" clearable style="width: 130px; margin-left: 12px" @change="loadInternal">
            <el-option v-for="(label, value) in RISK_LEVEL_OPTIONS" :key="value" :label="label" :value="value" />
          </el-select>
          <el-select v-model="internalFilters.status" placeholder="状态" clearable style="width: 140px; margin-left: 12px" @change="loadInternal">
            <el-option v-for="(label, value) in DEFECT_STATUS_LABELS" :key="value" :label="label" :value="value" />
          </el-select>
          <el-button type="primary" plain @click="internalPage = 1; loadInternal()" style="margin-left: 12px"><el-icon><Search /></el-icon> 查询</el-button>
          <el-button @click="resetInternalFilter">重置</el-button>
        </div>

        <el-table :data="internalDefects" stripe style="margin-top: 12px" v-loading="loadingInternal">
          <el-table-column label="缺陷编号" width="170">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="viewDefectDetail(row)">{{ row.defectNo }}</el-button>
            </template>
          </el-table-column>
          <el-table-column prop="source" label="来源" width="110">
            <template #default="{ row }"><el-tag size="small">{{ row.source }}</el-tag></template>
          </el-table-column>
          <el-table-column prop="title" label="风险标题" min-width="200" show-overflow-tooltip />
          <el-table-column label="风险等级" width="100">
            <template #default="{ row }">
              <el-tag :type="severityType(row.riskLevel)" size="small">{{ severityLabel(row.riskLevel) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="maintDept" label="处置部门" width="120">
            <template #default="{ row }">{{ row.maintDept || '-' }}</template>
          </el-table-column>
          <el-table-column label="全局状态" width="110">
            <template #default="{ row }">
              <el-tag :type="defectStatusType(row.status)" size="small">{{ defectStatusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="discoveredDate" label="上报时间" width="150">
            <template #default="{ row }">{{ formatTime(row.discoveredDate) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="viewDefectDetail(row)">详情</el-button>
              <el-button v-if="row.status === 'pending_review'" type="warning" link size="small" @click="viewDefectDetail(row); openReviewDialog()">审核</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="pagination-row">
          <el-pagination
            v-model:current-page="internalPage"
            v-model:page-size="internalPageSize"
            :total="internalTotal"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            @current-change="loadInternal"
            @size-change="onInternalPageSizeChange"
          />
        </div>
      </el-tab-pane>

      <!-- ==================== 外部缺陷 tab ==================== -->
      <el-tab-pane label="外部缺陷" name="external">
        <div class="toolbar">
          <el-input v-model="externalFilters.search" placeholder="搜索标题或软件" clearable style="width: 240px" @clear="loadExternal" @keyup.enter="loadExternal">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
          <el-select v-model="externalFilters.severity" placeholder="严重程度" clearable style="width: 130px; margin-left: 12px" @change="loadExternal">
            <el-option v-for="(label, value) in RISK_LEVEL_OPTIONS" :key="value" :label="label" :value="value" />
          </el-select>
          <el-select v-model="externalFilters.status" placeholder="状态" clearable style="width: 140px; margin-left: 12px" @change="loadExternal">
            <el-option v-for="(label, value) in DEFECT_STATUS_LABELS" :key="value" :label="label" :value="value" />
          </el-select>
          <el-button type="primary" plain @click="externalPage = 1; loadExternal()" style="margin-left: 12px"><el-icon><Search /></el-icon> 查询</el-button>
          <el-button @click="resetExternalFilter">重置</el-button>
        </div>

        <el-table :data="externalDefects" stripe style="margin-top: 12px" v-loading="loadingExternal">
          <el-table-column label="缺陷编号" width="180">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="viewDefectDetail(row)">{{ row.defectNo }}</el-button>
            </template>
          </el-table-column>
          <el-table-column prop="software" label="软件" width="100" />
          <el-table-column prop="title" label="缺陷标题" min-width="200" show-overflow-tooltip />
          <el-table-column label="严重程度" width="100">
            <template #default="{ row }">
              <el-tag :type="severityType(row.riskLevel)" size="small">{{ severityLabel(row.riskLevel) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="maintDept" label="运维部门" width="110">
            <template #default="{ row }">{{ row.maintDept || '-' }}</template>
          </el-table-column>
          <el-table-column label="研判人" width="100">
            <template #default="{ row }">
              <span v-if="row.dispatchTargetName">{{ row.dispatchTargetName }}</span>
              <el-tag v-else type="info" size="small">未分配</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="全局状态" width="110">
            <template #default="{ row }">
              <el-tag :type="defectStatusType(row.status)" size="small">{{ defectStatusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="viewDefectDetail(row)">详情</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="pagination-row">
          <el-pagination
            v-model:current-page="externalPage"
            v-model:page-size="externalPageSize"
            :total="externalTotal"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            @current-change="loadExternal"
            @size-change="onExternalPageSizeChange"
          />
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 缺陷详情对话框 -->
    <el-dialog v-model="detailDialogVisible" :title="isExternal(currentDefect) ? '外部缺陷详情' : '内部缺陷详情'" width="900px">
      <template v-if="currentDefect">
        <!-- 外部缺陷：内容多，用 tab 分页展示；内部缺陷内容少，保持单页 -->
        <el-tabs v-if="isExternal(currentDefect)" v-model="externalDetailTab" style="margin-bottom: 4px">
          <el-tab-pane label="基本信息" name="basic">
            <el-descriptions :column="2" border size="small">
              <el-descriptions-item label="缺陷编号">{{ currentDefect.defectNo }}</el-descriptions-item>
              <el-descriptions-item label="来源">{{ currentDefect.source }}</el-descriptions-item>
              <el-descriptions-item label="风险标题" :span="2">{{ currentDefect.title }}</el-descriptions-item>
              <el-descriptions-item label="软件">
                <el-tag size="small">{{ currentDefect.software }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="严重程度">
                <el-tag :type="severityType(currentDefect.riskLevel)" size="small">{{ severityLabel(currentDefect.riskLevel) }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="原始链接" :span="2">
                <el-link v-if="currentDefect.issueUrl" :href="currentDefect.issueUrl" target="_blank" type="primary">{{ currentDefect.issueUrl }}</el-link>
                <span v-else>-</span>
              </el-descriptions-item>
              <el-descriptions-item label="处置部门">{{ currentDefect.maintDept || '-' }}</el-descriptions-item>
              <el-descriptions-item label="研判人">{{ analystLabel(currentDefect) || '未分配' }}</el-descriptions-item>
              <el-descriptions-item label="处置截止时间">{{ formatTime(currentDefect.deadline) }}</el-descriptions-item>
              <el-descriptions-item label="全局状态">
                <el-tag :type="defectStatusType(currentDefect.status)" size="small">{{ defectStatusLabel(currentDefect.status) }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="研判结论" :span="2">{{ currentDefect.analyzeResult || '-' }}</el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>
          <el-tab-pane label="内容详情" name="content">
            <el-descriptions :column="1" border size="small">
              <el-descriptions-item label="缺陷英文标题">{{ currentDefect.title || '-' }}</el-descriptions-item>
              <el-descriptions-item label="中文标题">{{ currentDefect.titleZh || '-' }}</el-descriptions-item>
              <el-descriptions-item label="缺陷英文正文">
                <div class="pre-line">{{ currentDefect.description || '-' }}</div>
              </el-descriptions-item>
              <el-descriptions-item label="中文释义">
                <div class="pre-line">{{ currentDefect.summaryZh || '-' }}</div>
              </el-descriptions-item>
              <el-descriptions-item label="中文建议">
                <div class="pre-line">{{ currentDefect.repairAdvice || '-' }}</div>
              </el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>
        </el-tabs>
        <el-descriptions v-else :column="2" border size="small">
          <el-descriptions-item label="缺陷编号">{{ currentDefect.defectNo }}</el-descriptions-item>
          <el-descriptions-item label="来源">{{ currentDefect.source }}</el-descriptions-item>
          <el-descriptions-item label="风险标题" :span="2">{{ currentDefect.title }}</el-descriptions-item>
          <el-descriptions-item label="风险等级">
            <el-tag :type="severityType(currentDefect.riskLevel)" size="small">{{ severityLabel(currentDefect.riskLevel) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="上报时间">{{ formatTime(currentDefect.discoveredDate) }}</el-descriptions-item>
          <el-descriptions-item label="处置部门">{{ currentDefect.maintDept || '-' }}</el-descriptions-item>
          <el-descriptions-item label="研判人">{{ analystLabel(currentDefect) || '未分配' }}</el-descriptions-item>
          <el-descriptions-item label="处置截止时间">{{ formatTime(currentDefect.deadline) }}</el-descriptions-item>
          <el-descriptions-item label="全局状态">
            <el-tag :type="defectStatusType(currentDefect.status)" size="small">{{ defectStatusLabel(currentDefect.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">{{ currentDefect.description || '-' }}</el-descriptions-item>
          <el-descriptions-item label="研判结论" :span="2">{{ currentDefect.analyzeResult || '-' }}</el-descriptions-item>
        </el-descriptions>

        <!-- 处置动作 -->
        <div class="action-bar">
          <!-- 分配研判人：仅待研判且未真正分配研判人（dispatch_type='import' 的跟踪人不算，见后端 dispatchType 语义） -->
          <el-button v-if="currentDefect.status === 'pending_analysis' && currentDefect.dispatchType !== 'analyze'" type="warning" plain @click="openAnalystDialog">
            分配研判人
          </el-button>
          <el-button v-if="currentDefect.status === 'pending_analysis' && currentDefect.dispatchType !== 'analyze'" type="success" plain @click="autoAssignAnalyst">
            自动分配研判人
          </el-button>
          <el-button v-if="currentDefect.status === 'pending_analysis'" type="warning" @click="openAnalyzeDialog">研判</el-button>
          <el-button v-if="currentDefect.status === 'pending_review'" type="primary" @click="openReviewDialog">
            {{ isExternal(currentDefect) ? '评审（上会）' : '审核分配' }}
          </el-button>
        </div>

        <!-- 处置工单 -->
        <div v-if="disposals.length > 0" class="section-title">处置工单</div>
        <el-table v-if="disposals.length > 0" :data="disposals" size="small" border style="margin-top: 8px">
          <el-table-column prop="disposalNo" label="工单编号" width="170" />
          <el-table-column label="涉及项目" width="150">
            <template #default="{ row }">{{ row.projectName || '-' }}</template>
          </el-table-column>
          <el-table-column label="处置部门" width="110">
            <template #default="{ row }">{{ row.deptName }}</template>
          </el-table-column>
          <el-table-column label="类型" width="100">
            <template #default="{ row }">{{ DISPOSAL_TYPE_LABELS[row.disposalType] || row.disposalType }}</template>
          </el-table-column>
          <el-table-column label="处理人" width="100">
            <template #default="{ row }">{{ row.assigneeName || '-' }}</template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="disposalStatusType(row.status)" size="small">{{ disposalStatusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="ticketNo" label="工单号" width="120">
            <template #default="{ row }">{{ row.ticketNo || '-' }}</template>
          </el-table-column>
        </el-table>

        <!-- 状态历史 -->
        <div v-if="history.length > 0" class="section-title">状态历史</div>
        <el-table v-if="history.length > 0" :data="history" size="small" border style="margin-top: 8px">
          <el-table-column prop="action" label="动作" width="90" />
          <el-table-column prop="fromStatus" label="原状态" width="110">
            <template #default="{ row }">{{ defectStatusLabel(row.fromStatus) }}</template>
          </el-table-column>
          <el-table-column prop="toStatus" label="新状态" width="110">
            <template #default="{ row }">{{ defectStatusLabel(row.toStatus) }}</template>
          </el-table-column>
          <el-table-column prop="operatorName" label="操作人" width="100" />
          <el-table-column prop="remark" label="备注" min-width="160" />
          <el-table-column prop="createTime" label="时间" width="160">
            <template #default="{ row }">{{ formatTime(row.createTime) }}</template>
          </el-table-column>
        </el-table>
      </template>
    </el-dialog>

    <!-- 研判对话框（待研判：外部缺陷 或 未预判的缺陷） -->
    <el-dialog v-model="analyzeDialogVisible" title="缺陷研判" width="560px">
      <el-alert v-if="currentDefect" :title="`缺陷：${currentDefect.title}`" type="info" :closable="false" show-icon style="margin-bottom: 16px" />
      <el-form label-width="110px">
        <el-form-item label="研判结论" required>
          <el-select v-model="analyzeForm.conclusion" placeholder="请选择研判结论" style="width: 100%">
            <el-option v-for="opt in analyzeConclusionOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="研判说明">
          <el-input v-model="analyzeForm.analyzeResult" type="textarea" :rows="4" placeholder="填写研判分析结论" />
        </el-form-item>
        <!-- 漏洞直接分配（need_fix）时：分配处置人两通道（自动/手动），与缺陷一致 -->
        <template v-if="currentDefect?.flowType === 'vuln' && analyzeForm.conclusion === 'need_fix'">
          <el-form-item label="分配处置人">
            <DisposalAssignPicker v-model="analyzeForm.assign" :show-auto="true" />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="analyzeDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitAnalyze">确认研判</el-button>
      </template>
    </el-dialog>

    <!-- 分配研判人（第一次分配）对话框 -->
    <el-dialog v-model="analystDialogVisible" title="分配研判人（第一次分配）" width="520px">
      <el-alert
        type="info" :closable="false" show-icon style="margin-bottom: 12px"
        title="系统已按软件负责人自动匹配候选（第一个为主岗），也可手动指定"
      />
      <el-radio-group
        v-if="analystCandidates.length > 0"
        v-model="analystForm.selected"
        style="display: flex; flex-direction: column; gap: 8px"
      >
        <el-radio v-for="c in analystCandidates" :key="c.username" :value="c.username">
          {{ c.name }}（{{ c.username }}）｜{{ c.dept }} {{ c.team }}
          <el-tag v-if="c.isPrimary" type="success" size="small" style="margin-left: 4px">主岗</el-tag>
        </el-radio>
      </el-radio-group>
      <el-empty v-else description="未匹配到负责人候选，请手动指定" :image-size="60" />
      <el-form label-width="90px" style="margin-top: 12px">
        <el-form-item label="手动指定">
          <el-input v-model="analystForm.manual" placeholder="输入系统用户名（CAS 账号），优先于候选" clearable />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="analystDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" :disabled="!analystForm.manual && !analystForm.selected" @click="submitAnalyst">
          确认分配
        </el-button>
      </template>
    </el-dialog>

    <!-- 评审/审核分配对话框 -->
    <el-dialog v-model="reviewDialogVisible" :title="isExternal(currentDefect) ? '缺陷评审（上会）' : '内部缺陷审核分配'" width="620px">
      <el-alert v-if="currentDefect" :title="`缺陷：${currentDefect.title}`" type="info" :closable="false" show-icon style="margin-bottom: 16px" />
      <el-form label-width="110px">
        <el-form-item label="结论" required>
          <el-select v-model="reviewForm.reviewConclusion" placeholder="请选择结论" style="width: 100%">
            <el-option v-for="opt in REVIEW_CONCLUSIONS" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <!-- 需要修复时：分配处置人（两通道，与研判分配一致；内部缺陷仅手动指定） -->
        <template v-if="reviewForm.reviewConclusion === 'need_fix'">
          <el-form-item :label="(isExternal(currentDefect) || currentDefect?.flowType === 'vuln') ? '分配处置人' : '负责人'">
            <DisposalAssignPicker
              v-model="reviewForm.assign"
              :show-auto="isExternal(currentDefect) || currentDefect?.flowType === 'vuln'"
            />
          </el-form-item>
        </template>
        <!-- 外部缺陷评审（上会）：会议编号/截止时间/备注 -->
        <template v-if="isExternal(currentDefect)">
          <el-form-item label="会议编号">
            <el-input v-model="reviewForm.reviewMeetingNo" placeholder="评审会议编号" />
          </el-form-item>
          <el-form-item label="处置截止时间">
            <el-date-picker v-model="reviewForm.deadline" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" placeholder="选择处置截止时间" style="width: 100%" />
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="reviewForm.remark" type="textarea" :rows="3" placeholder="评审意见/备注" />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="reviewDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitReview">确认</el-button>
      </template>
    </el-dialog>

    <!-- 上报风险对话框（内部缺陷） -->
    <el-dialog v-model="createDialogVisible" title="上报风险（内部缺陷）" width="600px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="风险标题" required>
          <el-input v-model="createForm.title" placeholder="缺陷标题" />
        </el-form-item>
        <el-form-item label="来源" required>
          <el-select v-model="createForm.source" placeholder="选择来源" style="width: 100%">
            <el-option label="漏洞情报" value="vuln-intel" />
            <el-option label="SCA扫描" value="sca-scan" />
            <el-option label="人工上报" value="user-report" />
          </el-select>
        </el-form-item>
        <el-form-item label="风险等级">
          <el-select v-model="createForm.riskLevel" placeholder="选填" clearable style="width: 100%">
            <el-option v-for="(label, value) in RISK_LEVEL_OPTIONS" :key="value" :label="label" :value="value" />
          </el-select>
        </el-form-item>
        <el-form-item label="风险描述">
          <el-input v-model="createForm.description" type="textarea" :rows="4" placeholder="描述风险情况" />
        </el-form-item>
      </el-form>
      <el-alert type="info" :closable="false" show-icon title="上报后将生成内部缺陷编号（DF-xxx），进入「待审核」，由开源治理管理员审核后分配处置。" style="margin-top: 4px" />
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitCreate">提交</el-button>
      </template>
    </el-dialog>

    <!-- 漏洞情报对话框（手动录入 / 上传 Excel → 直接裂变处置工单） -->
    <el-dialog v-model="vulnIntelDialogVisible" title="漏洞情报 → 处置工单" width="720px">
      <div v-if="!vulnIntelMode" style="display: flex; gap: 16px; justify-content: center; padding: 30px 0">
        <el-button size="large" type="primary" plain @click="vulnIntelMode = 'manual'">
          <el-icon style="margin-right: 6px"><EditPen /></el-icon>手动录入
        </el-button>
        <el-button size="large" type="primary" plain @click="vulnIntelMode = 'upload'">
          <el-icon style="margin-right: 6px"><Upload /></el-icon>上传 Excel
        </el-button>
      </div>

      <template v-else-if="vulnIntelMode === 'manual'">
        <el-form :model="vulnIntelForm" label-width="100px">
          <el-form-item label="名称" required>
            <el-input v-model="vulnIntelForm.name" placeholder="软件或组件名称，如 redis / fastjson（用于匹配涉及项目处置人）" />
          </el-form-item>
          <el-form-item label="漏洞编号" required>
            <el-input v-model="vulnIntelForm.vulnNo" placeholder="唯一编号，如 CVE-2024-xxx（重复录入不重复裂变）" />
          </el-form-item>
          <el-form-item label="组织">
            <el-input v-model="vulnIntelForm.org" placeholder="组织命名空间，如 com.alibaba" />
          </el-form-item>
          <el-form-item label="标题">
            <el-input v-model="vulnIntelForm.title" placeholder="缺省用名称" />
          </el-form-item>
          <el-form-item label="等级">
            <el-select v-model="vulnIntelForm.severity" clearable placeholder="选填" style="width: 100%">
              <el-option v-for="(label, value) in RISK_LEVEL_OPTIONS" :key="value" :label="label" :value="value" />
            </el-select>
          </el-form-item>
          <el-form-item label="研判人">
            <el-input v-model="vulnIntelForm.analyst" placeholder="账号或姓名，缺省为当前用户" />
          </el-form-item>
          <el-form-item label="研判部门">
            <el-input v-model="vulnIntelForm.dept" placeholder="缺省按软件从 project 解析" />
          </el-form-item>
          <el-form-item label="截止时间">
            <el-date-picker v-model="vulnIntelForm.deadline" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" style="width: 100%" />
          </el-form-item>
          <el-form-item label="详情">
            <el-input v-model="vulnIntelForm.remark" type="textarea" :rows="3" placeholder="漏洞详情/备注" />
          </el-form-item>
          <el-form-item label="分配处置人">
            <DisposalAssignPicker v-model="vulnIntelForm.assign" :show-auto="true" />
          </el-form-item>
        </el-form>
        <div style="display: flex; gap: 8px; justify-content: flex-end">
          <el-button @click="vulnIntelMode = ''">返回</el-button>
          <el-button type="primary" :loading="submitting || previewSubmitting" @click="submitVulnIntelManual">创建并裂变</el-button>
        </div>
      </template>

      <template v-else>
        <el-alert
          type="info" :closable="false" show-icon style="margin-bottom: 12px"
          title="Excel 必填列：组织、名称、漏洞编号、研判人、研判部门。逐行创建漏洞情报并直接裂变处置工单；已存在的漏洞编号自动跳过（不重复裂变）。"
        />
        <el-upload
          drag :auto-upload="false" :show-file-list="false" accept=".xlsx,.xls"
          style="margin-bottom: 12px" :on-change="onVulnIntelFileChange"
        >
          <el-icon size="40" color="#909399"><UploadFilled /></el-icon>
          <div class="el-upload__text">拖拽 Excel 文件到此处，或 <em>点击选择文件</em></div>
        </el-upload>
        <div v-if="vulnIntelResult" style="margin-top: 12px">
          <el-alert :type="vulnIntelResultErrors ? 'warning' : 'success'" :closable="false" show-icon
            :title="`导入 ${vulnIntelResult.imported ?? 0} 条，跳过 ${vulnIntelResult.skipped ?? 0} 条`" />
          <pre v-if="vulnIntelResultErrors" style="max-height: 150px; overflow: auto; font-size: 12px; color: #e6a23c">{{ vulnIntelResultErrors }}</pre>
        </div>
        <div style="display: flex; gap: 8px; justify-content: flex-end">
          <el-button @click="vulnIntelMode = ''">返回</el-button>
          <el-button type="primary" :loading="submitting" @click="submitVulnIntelUpload">导入并裂变</el-button>
        </div>
      </template>

      <!-- 自动分配预览（可编辑负责人后确认创建，脱离 v-if 链） -->
      <VulnIntelPreviewDialog
        v-model="vulnPreviewVisible"
        :preview-list="vulnPreviewList"
        @confirm="onVulnPreviewConfirm"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Clock, Tickets, Operation, CircleCheck, Plus, Files, EditPen, UploadFilled } from '@element-plus/icons-vue'
import { riskApi } from '@/api'
import { formatTime } from '@/utils/format'
import {
  DEFECT_STATUS_LABELS,
  FLOW_TYPE_LABELS,
  DISPOSAL_TYPE_LABELS,
  ANALYZE_CONCLUSIONS,
  REVIEW_CONCLUSIONS,
  analystLabel,
  defectStatusLabel,
  defectStatusType,
  disposalStatusLabel,
  disposalStatusType,
  severityLabel,
  severityType,
} from '@/utils/defect'

import DisposalAssignPicker from '@/components/business/DisposalAssignPicker.vue'
import VulnIntelPreviewDialog from '@/components/business/VulnIntelPreviewDialog.vue'

const RISK_LEVEL_OPTIONS: Record<string, string> = { critical: '严重', high: '高危', medium: '中危', low: '低危', 未知: '未知' }

// 内部缺陷 vs 外部缺陷 判定：外部缺陷编号形如 软件#编号 或 source 为采集来源
function isExternal(row: any): boolean {
  if (!row) return false
  return row.affectedComponentType === 'external' || /#/.test(row.defectNo || '') || row.source === 'github-issue' || row.source === 'jira'
}

const loadingInternal = ref(false)
const loadingExternal = ref(false)
const submitting = ref(false)
const activeTab = ref('internal')

const internalDefects = ref<any[]>([])
const externalDefects = ref<any[]>([])
const internalFilters = ref({ search: '', severity: '', status: '' })
const externalFilters = ref({ search: '', severity: '', status: '' })
const internalPage = ref(1)
const internalPageSize = ref(10)
const internalTotal = ref(0)
const externalPage = ref(1)
const externalPageSize = ref(10)
const externalTotal = ref(0)

async function loadInternal() {
  loadingInternal.value = true
  try {
    const res: any = await riskApi.getInternalDefects({
      page: internalPage.value,
      pageSize: internalPageSize.value,
      keyword: internalFilters.value.search || undefined,
      severity: internalFilters.value.severity || undefined,
      status: internalFilters.value.status || undefined,
    })
    internalDefects.value = Array.isArray(res) ? res : (Array.isArray(res?.list) ? res.list : [])
    internalTotal.value = res?.total ?? internalDefects.value.length
  } catch (e: any) {
    ElMessage.error('加载内部缺陷失败: ' + (e.message || ''))
  } finally {
    loadingInternal.value = false
  }
}

async function loadExternal() {
  loadingExternal.value = true
  try {
    const res: any = await riskApi.getExternalDefects({
      page: externalPage.value,
      pageSize: externalPageSize.value,
      keyword: externalFilters.value.search || undefined,
      severity: externalFilters.value.severity || undefined,
      status: externalFilters.value.status || undefined,
    })
    externalDefects.value = Array.isArray(res) ? res : (Array.isArray(res?.list) ? res.list : [])
    externalTotal.value = res?.total ?? externalDefects.value.length
  } catch (e: any) {
    ElMessage.error('加载外部缺陷失败: ' + (e.message || ''))
  } finally {
    loadingExternal.value = false
  }
}

function onInternalPageSizeChange() {
  internalPage.value = 1
  loadInternal()
}
function onExternalPageSizeChange() {
  externalPage.value = 1
  loadExternal()
}

function loadDefects() {
  loadInternal()
  loadExternal()
}

const allDefects = computed(() => [...internalDefects.value, ...externalDefects.value])
function countByStatus(status: string): number {
  return allDefects.value.filter(d => d.status === status).length
}

function resetInternalFilter() {
  internalFilters.value = { search: '', severity: '', status: '' }
  internalPage.value = 1
  loadInternal()
}
function resetExternalFilter() {
  externalFilters.value = { search: '', severity: '', status: '' }
  externalPage.value = 1
  loadExternal()
}

// 详情
const detailDialogVisible = ref(false)
// 外部缺陷详情弹框 tab：basic-基本信息 / content-内容详情
const externalDetailTab = ref('basic')
const currentDefect = ref<any>(null)
const disposals = ref<any[]>([])
const history = ref<any[]>([])

async function viewDefectDetail(row: any) {
  currentDefect.value = row
  detailDialogVisible.value = true
  disposals.value = []
  history.value = []
  try {
    const [d, h]: any[] = await Promise.all([
      riskApi.getDefectDisposals(row.defectNo),
      riskApi.getDefectHistory(row.defectNo),
    ])
    disposals.value = d || []
    history.value = h || []
  } catch (e) {
    // 详情数据加载失败不阻断查看
  }
}

// 分配研判人（第一次分配）
const analystDialogVisible = ref(false)
const analystCandidates = ref<any[]>([])
const analystForm = ref({ selected: '', manual: '' })

async function openAnalystDialog() {
  detailDialogVisible.value = false
  analystForm.value = { selected: '', manual: '' }
  analystCandidates.value = []
  analystDialogVisible.value = true
  try {
    const res: any = await riskApi.getAnalystCandidates(currentDefect.value.defectNo)
    analystCandidates.value = res || []
    if (analystCandidates.value.length > 0) {
      const primary = analystCandidates.value.find((c: any) => c.isPrimary)
      analystForm.value.selected = primary?.username || analystCandidates.value[0].username
    }
  } catch (e: any) {
    ElMessage.error('获取候选负责人失败: ' + (e.message || ''))
  }
}

async function submitAnalyst() {
  if (!currentDefect.value) return
  const username = analystForm.value.manual?.trim() || analystForm.value.selected
  if (!username) return
  submitting.value = true
  try {
    await riskApi.assignAnalyst(currentDefect.value.defectNo, { analystUsername: username })
    ElMessage.success('已分配研判人')
    analystDialogVisible.value = false
    detailDialogVisible.value = false
    loadDefects()
  } catch (e: any) {
    ElMessage.error('分配失败: ' + (e.message || ''))
  } finally {
    submitting.value = false
  }
}

// 自动分配研判人（系统取主岗）
async function autoAssignAnalyst() {
  if (!currentDefect.value) return
  submitting.value = true
  try {
    const res: any = await riskApi.autoAssignAnalyst(currentDefect.value.defectNo)
    ElMessage.success(res?.message || '已自动分配研判人')
    detailDialogVisible.value = false
    loadDefects()
  } catch (e: any) {
    ElMessage.error('自动分配失败: ' + (e.message || ''))
  } finally {
    submitting.value = false
  }
}

// 研判
const analyzeDialogVisible = ref(false)
const analyzeForm = ref({ conclusion: '', analyzeResult: '', assign: { mode: 'auto', assignees: [] as string[] } })

const analyzeConclusionOptions = computed(() => {
  if (currentDefect.value?.flowType === 'vuln') {
    return ANALYZE_CONCLUSIONS
  }
  // 缺陷研判：可上会评审或忽略（与外部研判口径一致）；直接分配仅漏洞支持
  return ANALYZE_CONCLUSIONS.filter(c => c.value !== 'need_fix')
})

function openAnalyzeDialog() {
  detailDialogVisible.value = false
  analyzeForm.value = { conclusion: '', analyzeResult: '', assign: { mode: 'auto', assignees: [] } }
  analyzeDialogVisible.value = true
}

async function submitAnalyze() {
  if (!analyzeForm.value.conclusion) {
    ElMessage.warning('请选择研判结论')
    return
  }
  const assign = analyzeForm.value.assign || { mode: 'auto', assignees: [] }
  if (analyzeForm.value.conclusion === 'need_fix' && assign.mode === 'manual' && (!assign.assignees || assign.assignees.length === 0)) {
    ElMessage.warning('请填写至少一个负责人账号')
    return
  }
  submitting.value = true
  try {
    await riskApi.analyzeDefect(currentDefect.value.defectNo, {
      conclusion: analyzeForm.value.conclusion,
      analyzeResult: analyzeForm.value.analyzeResult,
      assignMode: assign.mode,
      assignees: assign.mode === 'manual' ? assign.assignees : [],
    })
    ElMessage.success('研判完成')
    analyzeDialogVisible.value = false
    detailDialogVisible.value = false
    loadDefects()
  } catch (e: any) {
    ElMessage.error('研判失败: ' + (e.message || ''))
  } finally {
    submitting.value = false
  }
}

// 评审 / 内部缺陷审核分配
const reviewDialogVisible = ref(false)
const reviewForm = ref({
  reviewConclusion: '',
  reviewMeetingNo: '',
  deadline: '',
  remark: '',
  assign: { mode: 'auto', assignees: [] as string[] },
})
function openReviewDialog() {
  // 从详情弹框打开评审/审核弹框：先关闭详情，避免两个弹框叠加（内外缺陷一致）
  detailDialogVisible.value = false
  reviewForm.value = { reviewConclusion: '', reviewMeetingNo: '', deadline: '', remark: '', assign: { mode: 'auto', assignees: [] } }
  reviewDialogVisible.value = true
}

async function submitReview() {
  if (!reviewForm.value.reviewConclusion) {
    ElMessage.warning('请选择结论')
    return
  }
  const external = isExternal(currentDefect.value)
  const needFix = reviewForm.value.reviewConclusion === 'need_fix'
  const assign = reviewForm.value.assign || { mode: 'auto', assignees: [] }
  // 手动指定必须填负责人（内部缺陷固定手动通道）
  if (needFix && assign.mode === 'manual' && (!assign.assignees || assign.assignees.length === 0)) {
    ElMessage.warning('请填写至少一个负责人账号')
    return
  }
  submitting.value = true
  try {
    await riskApi.reviewDefect(currentDefect.value.defectNo, {
      reviewConclusion: reviewForm.value.reviewConclusion,
      // 外部缺陷评审（上会）：会议编号/截止时间/备注
      reviewMeetingNo: external ? reviewForm.value.reviewMeetingNo : undefined,
      deadline: external ? (reviewForm.value.deadline || undefined) : undefined,
      remark: external ? reviewForm.value.remark : undefined,
      // 处置分配两通道：manual-手动指定（每人一张处置单）/ auto-自动分配（裂变+匹配）
      assignMode: needFix ? assign.mode : undefined,
      assignees: needFix && assign.mode === 'manual' ? assign.assignees : [],
    })
    ElMessage.success('评审完成')
    reviewDialogVisible.value = false
    detailDialogVisible.value = false
    loadDefects()
  } catch (e: any) {
    ElMessage.error('操作失败: ' + (e.message || ''))
  } finally {
    submitting.value = false
  }
}

// 上报风险（内部缺陷）
const createDialogVisible = ref(false)
const createForm = ref({ title: '', source: '', riskLevel: '', description: '' })

function openCreateDialog() {
  createForm.value = { title: '', source: '', riskLevel: '', description: '' }
  createDialogVisible.value = true
}

async function submitCreate() {
  if (!createForm.value.title || !createForm.value.source) {
    ElMessage.warning('请填写风险标题和来源')
    return
  }
  submitting.value = true
  try {
    await riskApi.createDefect({
      title: createForm.value.title,
      source: createForm.value.source,
      riskLevel: createForm.value.riskLevel || '未知',
      description: createForm.value.description,
    })
    ElMessage.success('上报成功，进入待审核')
    createDialogVisible.value = false
    loadDefects()
  } catch (e: any) {
    ElMessage.error('上报失败: ' + (e.message || ''))
  } finally {
    submitting.value = false
  }
}

// ===== 漏洞情报（手动录入 / 上传 Excel → 直接裂变处置工单） =====
const vulnIntelDialogVisible = ref(false)
const vulnIntelMode = ref('') // '' | 'manual' | 'upload'
const vulnIntelForm = ref({ org: '', name: '', vulnNo: '', title: '', severity: '', analyst: '', dept: '', deadline: '', remark: '', assign: { mode: 'auto', assignees: [] as string[] } })
const vulnIntelFile = ref<File | null>(null)
const vulnIntelResult = ref<any>(null)
const vulnIntelResultErrors = ref('')
// 自动分配预览（可编辑负责人后确认创建）
const vulnPreviewVisible = ref(false)
const vulnPreviewList = ref<any[]>([])
const previewSubmitting = ref(false)

function openVulnIntelDialog() {
  vulnIntelDialogVisible.value = true
  vulnIntelMode.value = ''
  vulnIntelResult.value = null
  vulnIntelResultErrors.value = ''
  vulnIntelFile.value = null
  vulnIntelForm.value = { org: '', name: '', vulnNo: '', title: '', severity: '', analyst: '', dept: '', deadline: '', remark: '', assign: { mode: 'auto', assignees: [] } }
}

function onVulnIntelFileChange(file: any) {
  vulnIntelFile.value = file?.raw || null
  vulnIntelResult.value = null
  vulnIntelResultErrors.value = ''
}

async function submitVulnIntelManual() {
  if (!vulnIntelForm.value.name || !vulnIntelForm.value.vulnNo) {
    ElMessage.warning('请填写名称和漏洞编号')
    return
  }
  const assign = vulnIntelForm.value.assign || { mode: 'auto', assignees: [] }
  if (assign.mode === 'manual' && (!assign.assignees || assign.assignees.length === 0)) {
    ElMessage.warning('请填写至少一个负责人账号')
    return
  }
  // 自动分配：先预览匹配结果（处置单清单+负责人，可编辑），确认后创建裂变
  if (assign.mode === 'auto') {
    previewSubmitting.value = true
    try {
      const preview: any = await riskApi.matchVulnIntel({ name: vulnIntelForm.value.name })
      vulnPreviewList.value = preview || []
      vulnPreviewVisible.value = true
    } catch (e: any) {
      ElMessage.error('匹配失败: ' + (e.message || ''))
    } finally {
      previewSubmitting.value = false
    }
    return
  }
  await createVulnIntelNow({ assignMode: 'manual', assignees: assign.assignees })
}

// 预览确认：以编辑后的负责人为准回传，批量裂变处置单
async function onVulnPreviewConfirm(overrides: any[]) {
  await createVulnIntelNow({ assignMode: 'auto', assignOverride: overrides })
}

async function createVulnIntelNow(extra: any) {
  submitting.value = true
  try {
    await riskApi.createVulnIntel({ ...vulnIntelForm.value, ...extra })
    ElMessage.success('创建成功，已生成处置工单')
    vulnIntelDialogVisible.value = false
    loadDefects()
  } catch (e: any) {
    ElMessage.error('创建失败: ' + (e.message || ''))
  } finally {
    submitting.value = false
  }
}

async function submitVulnIntelUpload() {
  if (!vulnIntelFile.value) {
    ElMessage.warning('请选择 Excel 文件')
    return
  }
  const fd = new FormData()
  fd.append('file', vulnIntelFile.value)
  submitting.value = true
  try {
    const res: any = await riskApi.importVulnIntel(fd)
    vulnIntelResult.value = res || {}
    vulnIntelResultErrors.value = (res?.errors || []).join('\n')
    ElMessage.success(`导入 ${res?.imported ?? 0} 条，跳过 ${res?.skipped ?? 0} 条`)
    loadDefects()
  } catch (e: any) {
    ElMessage.error('导入失败: ' + (e.message || ''))
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadDefects()
})
</script>

<style scoped>
.pagination-row { display: flex; justify-content: flex-end; margin-top: 12px; }
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}
.page-header h3 { margin: 0 0 4px; }
.desc { color: #909399; font-size: 13px; margin: 0; }

.stat-card { border: none; }
.stat-content { display: flex; align-items: center; gap: 16px; }
.stat-icon {
  width: 56px; height: 56px; border-radius: 12px;
  display: flex; align-items: center; justify-content: center;
}
.stat-info { flex: 1; }
.stat-value { font-size: 28px; font-weight: 600; color: #303133; line-height: 1.2; }
.stat-label { font-size: 13px; color: #909399; margin-top: 4px; }

.toolbar {
  display: flex; align-items: center; gap: 12px; flex-wrap: wrap; margin-top: 16px;
}
.action-bar { margin: 16px 0; }
.section-title {
  font-size: 14px; font-weight: 600; color: #303133;
  margin: 16px 0 0; padding-left: 8px; border-left: 3px solid #409EFF;
}
/* 长文本（英文正文/中文释义/中文建议）保留换行展示 */
.pre-line {
  white-space: pre-wrap; word-break: break-word; line-height: 1.6; max-height: 320px; overflow-y: auto;
}
</style>
