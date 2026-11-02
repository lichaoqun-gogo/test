<template>
  <div class="my-ops">
    <div class="page-header">
      <h2>我的运维</h2>
      <p class="page-desc">我负责维护的开源基础软件</p>
    </div>

    <el-tabs v-model="activeTab" class="ops-tabs">
      <!-- 我负责的软件 -->
      <el-tab-pane label="我负责的软件" name="software">
        <el-table :data="softwareList" border v-loading="loading">
          <el-table-column prop="grade" label="等级" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="gradeType(row.grade)" size="small" effect="dark">{{ row.grade }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="name" label="软件名称" min-width="160">
            <template #default="{ row }">
              <el-link type="primary" @click="showSoftwareDetail(row)">{{ row.name }}</el-link>
            </template>
          </el-table-column>
          <el-table-column prop="version" label="当前版本" width="120" />
          <el-table-column prop="type" label="类型" width="120">
            <template #default="{ row }"><el-tag size="small">{{ row.type }}</el-tag></template>
          </el-table-column>
          <el-table-column prop="owner" label="负责人" width="100" />
          <el-table-column prop="defectCount" label="缺陷数" width="100" align="center">
            <template #default="{ row }">
              <el-tag v-if="row.defectCount > 0" type="danger" size="small">{{ row.defectCount }}</el-tag>
              <span v-else>0</span>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- 缺陷管理（内部/外部分tab） -->
      <el-tab-pane label="缺陷管理" name="defects">
        <div style="margin-bottom: 16px; display: flex; gap: 12px; align-items: center; flex-wrap: wrap">
          <el-button type="primary" plain @click="openImportDialog">
            <el-icon><Upload /></el-icon> 导入缺陷
          </el-button>
          <el-button type="primary" plain @click="openVulnIntelDialog">
            <el-icon><Files /></el-icon> 导入漏洞情报
          </el-button>
          <span style="color: #909399; font-size: 13px">漏洞情报输入软件/组件名称与详情，创建后直接裂变处置工单</span>
        </div>

        <el-tabs v-model="defectActiveTab" size="small">
          <!-- 内部缺陷 -->
          <el-tab-pane label="内部缺陷" name="internal">
            <div style="margin-bottom: 16px; display: flex; gap: 12px; align-items: center; flex-wrap: wrap">
              <el-input v-model="internalKeyword" placeholder="搜索缺陷标题" clearable style="width: 220px" @keyup.enter="loadInternalDefects">
                <template #prefix><el-icon><Search /></el-icon></template>
              </el-input>
              <el-select v-model="internalStatus" placeholder="全局状态" clearable style="width: 140px" @change="loadInternalDefects">
                <el-option v-for="(label, value) in DEFECT_STATUS_LABELS" :key="value" :label="label" :value="value" />
              </el-select>
              <el-button type="primary" plain @click="internalPage = 1; loadInternalDefects()"><el-icon><Search /></el-icon> 查询</el-button>
              <el-button @click="resetInternalFilter">重置</el-button>
            </div>
            <el-table :data="internalDefectList" border v-loading="internalLoading" @expand-change="handleExpandChange" row-key="defectNo">
              <el-table-column type="expand">
                <template #default="{ row }">
                  <div class="expand-disposals" style="padding: 12px 24px 12px 48px; background: #fafafa">
                    <div class="section-title" style="display: flex; align-items: center; justify-content: space-between">
                      <span>处置工单</span>
                      <el-button v-if="canAutoAssignDisposals(row)" type="success" size="small" plain @click="autoAssignDisposals(row)">一键自动分配</el-button>
                    </div>
                    <el-table v-if="disposalsOf(row).length > 0" :data="disposalsOf(row)" size="small" border style="margin-top: 8px">
                      <el-table-column prop="disposalNo" label="工单编号" width="170" />
                      <el-table-column label="涉及项目" width="150">
                        <template #default="{ row: d }">{{ d.projectName || '-' }}</template>
                      </el-table-column>
                      <el-table-column prop="deptName" label="处置部门" width="110" />
                      <el-table-column label="类型" width="100">
                        <template #default="{ row: d }">{{ DISPOSAL_TYPE_LABELS[d.disposalType] || d.disposalType }}</template>
                      </el-table-column>
                      <el-table-column label="处理人" width="200">
                        <template #default="{ row: d }">
                          <el-select
                            v-if="d.status === 'pending_assign' && isMyDept(d.deptName)"
                            v-model="d.assigneeUsername"
                            :placeholder="d.assigneeName || '请选择处理人'"
                            size="small" filterable style="width: 100%"
                            @change="(val: string) => changeAssignee(d, val)"
                          >
                            <el-option v-for="u in deptUserMap[d.deptName] || []" :key="u.username" :label="`${u.nickname}（${u.username}）`" :value="u.username" />
                          </el-select>
                          <span v-else>{{ d.assigneeName || '-' }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="状态" width="100">
                        <template #default="{ row: d }">
                          <el-tag :type="disposalStatusType(d.status)" size="small">{{ disposalStatusLabel(d.status) }}</el-tag>
                        </template>
                      </el-table-column>
                    </el-table>
                    <el-empty v-else description="暂无处置工单" :image-size="50" style="padding: 12px 0" />
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="defectNo" label="缺陷编号" width="170" />
              <el-table-column prop="source" label="来源" width="110" />
              <el-table-column prop="title" label="风险标题" min-width="200" show-overflow-tooltip />
              <el-table-column label="风险等级" width="100" align="center">
                <template #default="{ row }">
                  <el-tag :type="severityType(row.riskLevel)" size="small">{{ severityLabel(row.riskLevel) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="maintDept" label="处置部门" width="110">
                <template #default="{ row }">{{ row.maintDept || '-' }}</template>
              </el-table-column>
              <el-table-column label="全局状态" width="100" align="center">
                <template #default="{ row }">
                  <el-tag :type="defectStatusType(row.status)" size="small">{{ defectStatusLabel(row.status) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="90" align="center">
                <template #default="{ row }">
                  <el-button type="primary" link size="small" @click="viewDefectDetail(row)">详情</el-button>
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
                @current-change="loadInternalDefects"
                @size-change="onInternalPageSizeChange"
              />
            </div>
          </el-tab-pane>

          <!-- 外部缺陷 -->
          <el-tab-pane label="外部缺陷" name="external">
            <div style="margin-bottom: 16px; display: flex; gap: 12px; align-items: center; flex-wrap: wrap">
              <el-input v-model="externalKeyword" placeholder="搜索缺陷标题或软件" clearable style="width: 220px" @keyup.enter="loadExternalDefects">
                <template #prefix><el-icon><Search /></el-icon></template>
              </el-input>
              <el-select v-model="externalStatus" placeholder="全局状态" clearable style="width: 140px" @change="loadExternalDefects">
                <el-option v-for="(label, value) in DEFECT_STATUS_LABELS" :key="value" :label="label" :value="value" />
              </el-select>
              <el-button type="primary" plain @click="externalPage = 1; loadExternalDefects()"><el-icon><Search /></el-icon> 查询</el-button>
              <el-button @click="resetExternalFilter">重置</el-button>
            </div>
            <el-table :data="externalDefectList" border v-loading="externalLoading" @expand-change="handleExpandChange" row-key="defectNo">
              <el-table-column type="expand">
                <template #default="{ row }">
                  <div class="expand-disposals" style="padding: 12px 24px 12px 48px; background: #fafafa">
                    <div class="section-title" style="display: flex; align-items: center; justify-content: space-between">
                      <span>处置工单</span>
                      <el-button v-if="canAutoAssignDisposals(row)" type="success" size="small" plain @click="autoAssignDisposals(row)">一键自动分配</el-button>
                    </div>
                    <el-table v-if="disposalsOf(row).length > 0" :data="disposalsOf(row)" size="small" border style="margin-top: 8px">
                      <el-table-column prop="disposalNo" label="工单编号" width="170" />
                      <el-table-column label="涉及项目" width="150">
                        <template #default="{ row: d }">{{ d.projectName || '-' }}</template>
                      </el-table-column>
                      <el-table-column prop="deptName" label="处置部门" width="110" />
                      <el-table-column label="类型" width="100">
                        <template #default="{ row: d }">{{ DISPOSAL_TYPE_LABELS[d.disposalType] || d.disposalType }}</template>
                      </el-table-column>
                      <el-table-column label="处理人" width="200">
                        <template #default="{ row: d }">
                          <el-select
                            v-if="d.status === 'pending_assign' && isMyDept(d.deptName)"
                            v-model="d.assigneeUsername"
                            :placeholder="d.assigneeName || '请选择处理人'"
                            size="small" filterable style="width: 100%"
                            @change="(val: string) => changeAssignee(d, val)"
                          >
                            <el-option v-for="u in deptUserMap[d.deptName] || []" :key="u.username" :label="`${u.nickname}（${u.username}）`" :value="u.username" />
                          </el-select>
                          <span v-else>{{ d.assigneeName || '-' }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="状态" width="100">
                        <template #default="{ row: d }">
                          <el-tag :type="disposalStatusType(d.status)" size="small">{{ disposalStatusLabel(d.status) }}</el-tag>
                        </template>
                      </el-table-column>
                    </el-table>
                    <el-empty v-else description="暂无处置工单" :image-size="50" style="padding: 12px 0" />
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="defectNo" label="缺陷编号" width="190" />
              <el-table-column prop="software" label="软件" width="110" />
              <el-table-column prop="title" label="缺陷标题" min-width="200" show-overflow-tooltip />
              <el-table-column label="严重程度" width="100" align="center">
                <template #default="{ row }">
                  <el-tag :type="severityType(row.riskLevel)" size="small">{{ severityLabel(row.riskLevel) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="maintDept" label="运维部门" width="110">
                <template #default="{ row }">{{ row.maintDept || '-' }}</template>
              </el-table-column>
              <el-table-column label="研判人" width="90" align="center">
                <template #default="{ row }">
                  <span v-if="row.dispatchTargetName">{{ row.dispatchTargetName }}</span>
                  <el-tag v-else type="info" size="small">未分配</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="全局状态" width="100" align="center">
                <template #default="{ row }">
                  <el-tag :type="defectStatusType(row.status)" size="small">{{ defectStatusLabel(row.status) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="90" align="center">
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
                @current-change="loadExternalDefects"
                @size-change="onExternalPageSizeChange"
              />
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-tab-pane>
    </el-tabs>

    <!-- 软件详情对话框 -->
    <ProjectDetailDialog v-model="softwareDetailVisible" :project-data="currentSoftware" />

    <!-- 缺陷详情对话框 -->
    <el-dialog v-model="defectDetailDialogVisible" :title="isExternal(currentDefect) ? '外部缺陷详情' : '内部缺陷详情'" width="900px">
      <template v-if="currentDefect">
        <!-- 外部缺陷：内容多，用 tab 分页展示；内部缺陷内容少，保持单页 -->
        <el-tabs v-if="isExternal(currentDefect)" v-model="externalDetailTab" style="margin-bottom: 4px">
          <el-tab-pane label="基本信息" name="basic">
            <el-descriptions :column="2" border size="small">
              <el-descriptions-item label="缺陷编号">{{ currentDefect.defectNo }}</el-descriptions-item>
              <el-descriptions-item label="来源">{{ currentDefect.source }}</el-descriptions-item>
              <el-descriptions-item label="风险标题" :span="2">{{ currentDefect.title }}</el-descriptions-item>
              <el-descriptions-item label="软件"><el-tag size="small">{{ currentDefect.software }}</el-tag></el-descriptions-item>
              <el-descriptions-item label="严重程度"><el-tag :type="severityType(currentDefect.riskLevel)" size="small">{{ severityLabel(currentDefect.riskLevel) }}</el-tag></el-descriptions-item>
              <el-descriptions-item label="原始链接" :span="2">
                <el-link v-if="currentDefect.issueUrl" :href="currentDefect.issueUrl" target="_blank" type="primary">{{ currentDefect.issueUrl }}</el-link>
                <span v-else>-</span>
              </el-descriptions-item>
              <el-descriptions-item label="处置部门">{{ currentDefect.maintDept || '-' }}</el-descriptions-item>
              <el-descriptions-item label="处置截止时间">{{ formatTime(currentDefect.deadline) }}</el-descriptions-item>
              <el-descriptions-item label="全局状态" :span="2">
                <el-tag :type="defectStatusType(currentDefect.status)" size="small">{{ defectStatusLabel(currentDefect.status) }}</el-tag>
              </el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>
          <el-tab-pane label="内容详情" name="content">
            <el-descriptions :column="1" border size="small">
              <el-descriptions-item label="缺陷英文标题">{{ currentDefect.title || '-' }}</el-descriptions-item>
              <el-descriptions-item label="中文标题">{{ currentDefect.titleZh || '-' }}</el-descriptions-item>
              <el-descriptions-item label="缺陷英文正文"><div class="pre-line">{{ currentDefect.description || '-' }}</div></el-descriptions-item>
              <el-descriptions-item label="中文释义"><div class="pre-line">{{ currentDefect.summaryZh || '-' }}</div></el-descriptions-item>
              <el-descriptions-item label="中文建议"><div class="pre-line">{{ currentDefect.repairAdvice || '-' }}</div></el-descriptions-item>
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
          <el-descriptions-item label="处置截止时间">{{ formatTime(currentDefect.deadline) }}</el-descriptions-item>
          <el-descriptions-item label="全局状态">
            <el-tag :type="defectStatusType(currentDefect.status)" size="small">{{ defectStatusLabel(currentDefect.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">{{ currentDefect.description || '-' }}</el-descriptions-item>
        </el-descriptions>

        <div class="action-bar">
          <!-- 分配研判人：仅待研判且未真正分配研判人（dispatch_type='import' 的跟踪人不算，见后端 dispatchType 语义） -->
          <el-button v-if="currentDefect.status === 'pending_analysis' && currentDefect.dispatchType !== 'analyze'" type="warning" plain @click="openAnalystDialog">
            分配研判人
          </el-button>
          <el-button v-if="currentDefect.status === 'pending_analysis' && currentDefect.dispatchType !== 'analyze'" type="success" plain @click="autoAssignAnalyst">
            自动分配研判人
          </el-button>
        </div>

        <!-- 历史记录 -->
        <div class="section-title" style="margin-top: 16px">历史记录</div>
        <el-table :data="defectHistory" size="small" border style="margin-top: 8px">
          <el-table-column prop="createTime" label="时间" width="170" />
          <el-table-column prop="operatorName" label="操作人" width="100" />
          <el-table-column label="动作" width="140">
            <template #default="{ row }">{{ HISTORY_ACTION_LABELS[row.action] || row.action }}</template>
          </el-table-column>
          <el-table-column label="状态流转" width="180">
            <template #default="{ row }">
              <span v-if="row.fromStatus && row.toStatus && row.fromStatus !== row.toStatus">
                {{ row.disposalId > 0 ? disposalStatusLabel(row.fromStatus) : defectStatusLabel(row.fromStatus) }}
                →
                {{ row.disposalId > 0 ? disposalStatusLabel(row.toStatus) : defectStatusLabel(row.toStatus) }}
              </span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="说明" min-width="200" show-overflow-tooltip />
        </el-table>
      </template>
    </el-dialog>

    <!-- 导入缺陷对话框（采集工具产物，幂等：同一缺陷编号只有一条记录） -->
    <el-dialog v-model="importDialogVisible" title="导入缺陷" width="760px">
      <el-alert
        type="info" :closable="false" show-icon style="margin-bottom: 12px"
        title="粘贴采集工具输出的 JSON（如 redis-all.json 全文）。按缺陷编号（redis#xxx）幂等导入：已存在缺陷只更新采集信息，内部处置状态不会回退。"
      />
      <el-upload
        drag :auto-upload="false" :show-file-list="false" accept=".json,application/json"
        style="margin-bottom: 12px" :on-change="onImportFileChange"
      >
        <el-icon size="40" color="#909399"><UploadFilled /></el-icon>
        <div class="el-upload__text">拖拽 JSON 文件到此处，或 <em>点击选择文件</em></div>
      </el-upload>
      <el-input
        v-model="importText" type="textarea" :rows="9" placeholder='{"software":"Redis","defects":[...]}'
        style="font-family: monospace"
      />
      <div v-if="importResult" style="margin-top: 12px">
        <el-alert type="success" :closable="false" show-icon
          :title="`新增 ${importResult.imported ?? 0} 条，更新 ${importResult.updated ?? 0} 条，跳过 ${importResult.skipped ?? 0} 条`" />
        <el-alert v-if="importResult.errors && importResult.errors.length" type="warning" :closable="false" show-icon
          :title="importResult.errors.slice(0, 5).join('；')" style="margin-top: 8px" />
      </div>
      <template #footer>
        <el-button @click="importDialogVisible = false">关闭</el-button>
        <el-button type="primary" :loading="importing" @click="doImport">导入</el-button>
      </template>
    </el-dialog>

    <!-- 漏洞情报对话框（手动录入 / 上传 Excel → 直接裂变处置工单） -->
    <el-dialog v-model="vulnIntelDialogVisible" title="漏洞情报 → 处置工单" width="720px">
      <!-- 选择模式 -->
      <div v-if="!vulnIntelMode" style="display: flex; gap: 16px; justify-content: center; padding: 30px 0">
        <el-button size="large" type="primary" plain @click="vulnIntelMode = 'manual'">
          <el-icon style="margin-right: 6px"><EditPen /></el-icon>手动录入
        </el-button>
        <el-button size="large" type="primary" plain @click="vulnIntelMode = 'upload'">
          <el-icon style="margin-right: 6px"><Upload /></el-icon>上传 Excel
        </el-button>
      </div>

      <!-- 手动录入 -->
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
          <el-button type="primary" :loading="vulnIntelSubmitting || vulnPreviewSubmitting" @click="submitVulnIntelManual">创建并裂变</el-button>
        </div>
      </template>

      <!-- 上传 Excel -->
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
          <el-button type="primary" :loading="vulnIntelSubmitting" @click="submitVulnIntelUpload">导入并裂变</el-button>
        </div>
      </template>

      <!-- 自动分配预览（可编辑负责人后确认创建，脱离 v-if 链） -->
      <VulnIntelPreviewDialog
        v-model="vulnPreviewVisible"
        :preview-list="vulnPreviewList"
        @confirm="onVulnPreviewConfirm"
      />
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

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, UploadFilled, Files, EditPen } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { riskApi, personnelApi } from '@/api'
import request from '@/api'
import { formatTime } from '@/utils/format'
import {
  DEFECT_STATUS_LABELS,
  FLOW_TYPE_LABELS,
  DISPOSAL_TYPE_LABELS,
  defectStatusLabel,
  defectStatusType,
  disposalStatusLabel,
  disposalStatusType,
  severityLabel,
  severityType,
} from '@/utils/defect'
import ProjectDetailDialog from '@/components/business/ProjectDetailDialog.vue'
import DisposalAssignPicker from '@/components/business/DisposalAssignPicker.vue'
import VulnIntelPreviewDialog from '@/components/business/VulnIntelPreviewDialog.vue'

const userStore = useUserStore()

const activeTab = ref('software')
const loading = ref(false)
const softwareList = ref<any[]>([])

async function loadSoftwareList() {
  loading.value = true
  try {
    const dept = userStore.userInfo?.department || ''
    const res: any = await request.get('/asset/projects', { params: { maint_dept: dept } })
    softwareList.value = (res || []).map((item: any) => ({
      ...item,
      name: item._name,
      version: item._version,
      grade: item.applicable_level || 'A',
      type: getTypeLabel(item.type),
      owner: item.create_by || '-',
      department: item.maint_dept || '-',
      description: item.DESCRIPTION || '-',
      defectCount: 0,
    }))
  } catch (error) {
    ElMessage.error('加载软件列表失败')
  } finally {
    loading.value = false
  }
}

function getTypeLabel(type: string) {
  const map: Record<string, string> = {
    java: 'Java组件', go: 'Go组件', python: 'Python组件', js: 'JavaScript组件',
  }
  return map[type] || type || '-'
}

function gradeType(grade: string) {
  const map: Record<string, any> = { A: 'success', B: '', C: 'warning', D: 'danger' }
  return map[grade] || 'info'
}

// 软件详情
const softwareDetailVisible = ref(false)
const currentSoftware = ref<any>(null)
function showSoftwareDetail(software: any) {
  currentSoftware.value = software
  softwareDetailVisible.value = true
}

// ===== 缺陷管理（内部/外部分tab） =====
const defectActiveTab = ref('internal')
const internalLoading = ref(false)
const externalLoading = ref(false)
const internalKeyword = ref('')
const internalStatus = ref('')
const externalKeyword = ref('')
const externalStatus = ref('')
const internalDefectList = ref<any[]>([])
const externalDefectList = ref<any[]>([])
const internalPage = ref(1)
const internalPageSize = ref(10)
const internalTotal = ref(0)
const externalPage = ref(1)
const externalPageSize = ref(10)
const externalTotal = ref(0)

// 内部缺陷 vs 外部缺陷 判定
function isExternal(row: any): boolean {
  if (!row) return false
  return row.affectedComponentType === 'external' || /#/.test(row.defectNo || '') || row.source === 'github-issue' || row.source === 'jira'
}

async function loadInternalDefects() {
  internalLoading.value = true
  try {
    const res: any = await riskApi.getInternalDefects({
      page: internalPage.value,
      pageSize: internalPageSize.value,
      keyword: internalKeyword.value || undefined,
      status: internalStatus.value || undefined,
      // 个人工作台视图：管理部门专员也与部门专员一致按部门过滤，仅在全局治理（风险管理-缺陷台账）看全景
      scope: 'dept',
    })
    // 后端返回 {list,total} 分页结构，list 可能为 null；空结果时保持空数组，避免把对象赋给表格导致渲染崩溃
    internalDefectList.value = Array.isArray(res) ? res : (Array.isArray(res?.list) ? res.list : [])
    internalTotal.value = res?.total ?? internalDefectList.value.length
  } catch (e: any) {
    ElMessage.error('加载内部缺陷失败: ' + (e.message || ''))
  } finally {
    internalLoading.value = false
  }
}

async function loadExternalDefects() {
  externalLoading.value = true
  try {
    const res: any = await riskApi.getExternalDefects({
      page: externalPage.value,
      pageSize: externalPageSize.value,
      keyword: externalKeyword.value || undefined,
      status: externalStatus.value || undefined,
      scope: 'dept',
    })
    externalDefectList.value = Array.isArray(res) ? res : (Array.isArray(res?.list) ? res.list : [])
    externalTotal.value = res?.total ?? externalDefectList.value.length
  } catch (e: any) {
    ElMessage.error('加载外部缺陷失败: ' + (e.message || ''))
  } finally {
    externalLoading.value = false
  }
}

function onInternalPageSizeChange() {
  internalPage.value = 1
  loadInternalDefects()
}
function onExternalPageSizeChange() {
  externalPage.value = 1
  loadExternalDefects()
}

function loadDefects() {
  loadInternalDefects()
  loadExternalDefects()
}

function resetInternalFilter() {
  internalKeyword.value = ''
  internalStatus.value = ''
  internalPage.value = 1
  loadInternalDefects()
}
function resetExternalFilter() {
  externalKeyword.value = ''
  externalStatus.value = ''
  externalPage.value = 1
  loadExternalDefects()
}

// ===== 导入缺陷（采集工具产物，幂等） =====
const importDialogVisible = ref(false)
const importText = ref('')
const importing = ref(false)
const importResult = ref<any>(null)

function openImportDialog() {
  importText.value = ''
  importResult.value = null
  importDialogVisible.value = true
}

function onImportFileChange(file: any) {
  const raw = file?.raw
  if (!raw) return
  const reader = new FileReader()
  reader.onload = () => {
    importText.value = String(reader.result || '')
  }
  reader.readAsText(raw)
}

async function doImport() {
  if (!importText.value.trim()) {
    ElMessage.warning('请粘贴 JSON 内容或选择文件')
    return
  }
  let payload: any
  try {
    payload = JSON.parse(importText.value)
  } catch (e) {
    ElMessage.error('JSON 解析失败，请检查格式')
    return
  }
  importing.value = true
  try {
    const res: any = await riskApi.importDefects(payload)
    importResult.value = res || {}
    ElMessage.success(`导入完成：新增 ${res?.imported ?? 0}，更新 ${res?.updated ?? 0}，跳过 ${res?.skipped ?? 0}`)
    loadDefects()
  } catch (e: any) {
    ElMessage.error('导入失败: ' + (e.message || ''))
  } finally {
    importing.value = false
  }
}

// ===== 漏洞情报（手动录入 / 上传 Excel → 直接裂变处置工单） =====
const RISK_LEVEL_OPTIONS: Record<string, string> = { critical: '严重', high: '高危', medium: '中危', low: '低危', 未知: '未知' }

const vulnIntelDialogVisible = ref(false)
const vulnIntelMode = ref('') // '' | 'manual' | 'upload'
const vulnIntelForm = ref({ org: '', name: '', vulnNo: '', title: '', severity: '', analyst: '', dept: '', deadline: '', remark: '', assign: { mode: 'auto', assignees: [] as string[] } })
const vulnIntelSubmitting = ref(false)
const vulnIntelFile = ref<File | null>(null)
const vulnIntelResult = ref<any>(null)
const vulnIntelResultErrors = ref('')
// 自动分配预览（可编辑负责人后确认创建）
const vulnPreviewVisible = ref(false)
const vulnPreviewList = ref<any[]>([])
const vulnPreviewSubmitting = ref(false)

function openVulnIntelDialog() {
  vulnIntelDialogVisible.value = true
  vulnIntelMode.value = ''
  vulnIntelResult.value = null
  vulnIntelResultErrors.value = ''
  vulnIntelFile.value = null
  vulnIntelForm.value = { org: '', name: '', vulnNo: '', title: '', severity: '', analyst: '', dept: '', deadline: '', remark: '', assign: { mode: 'auto', assignees: [] } }
}

function onVulnIntelFileChange(file: any) {
  const raw = file?.raw
  vulnIntelFile.value = raw || null
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
    vulnPreviewSubmitting.value = true
    try {
      const preview: any = await riskApi.matchVulnIntel({ name: vulnIntelForm.value.name })
      vulnPreviewList.value = preview || []
      vulnPreviewVisible.value = true
    } catch (e: any) {
      ElMessage.error('匹配失败: ' + (e.message || ''))
    } finally {
      vulnPreviewSubmitting.value = false
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
  vulnIntelSubmitting.value = true
  try {
    await riskApi.createVulnIntel({ ...vulnIntelForm.value, ...extra })
    ElMessage.success('创建成功，已生成处置工单')
    vulnIntelDialogVisible.value = false
    loadDefects()
  } catch (e: any) {
    ElMessage.error('创建失败: ' + (e.message || ''))
  } finally {
    vulnIntelSubmitting.value = false
  }
}

async function submitVulnIntelUpload() {
  if (!vulnIntelFile.value) {
    ElMessage.warning('请选择 Excel 文件')
    return
  }
  const fd = new FormData()
  fd.append('file', vulnIntelFile.value)
  vulnIntelSubmitting.value = true
  try {
    const res: any = await riskApi.importVulnIntel(fd)
    vulnIntelResult.value = res || {}
    vulnIntelResultErrors.value = (res?.errors || []).join('\n')
    ElMessage.success(`导入 ${res?.imported ?? 0} 条，跳过 ${res?.skipped ?? 0} 条`)
    loadDefects()
  } catch (e: any) {
    ElMessage.error('导入失败: ' + (e.message || ''))
  } finally {
    vulnIntelSubmitting.value = false
  }
}

// 详情 + 历史
const defectDetailDialogVisible = ref(false)
// 外部缺陷详情弹框 tab：basic-基本信息 / content-内容详情
const externalDetailTab = ref('basic')
const currentDefect = ref<any>(null)
const defectHistory = ref<any[]>([])
const submitting = ref(false)

// 展开行：按缺陷编号缓存的处置单列表
const disposalsMap = ref<Record<string, any[]>>({})

function disposalsOf(defect: any): any[] {
  return disposalsMap.value[defect.defectNo] || []
}

// 懒加载展开行的处置单（只在展开时请求一次；force 时强制刷新）
async function loadExpandDisposals(defect: any, force = false) {
  if (!force && disposalsMap.value[defect.defectNo]) return
  try {
    const list = ((await riskApi.getDefectDisposals(defect.defectNo)) as any) || []
    disposalsMap.value = { ...disposalsMap.value, [defect.defectNo]: list }
    loadDisposalDeptUsers(list)
  } catch (e) {
    disposalsMap.value = { ...disposalsMap.value, [defect.defectNo]: [] }
  }
}

function handleExpandChange(row: any, expandedRows: any[] | boolean) {
  // 兼容 Element Plus 的两种事件签名（expandedRows 数组 / expanded 布尔）
  const expanded = Array.isArray(expandedRows) ? expandedRows.some((r: any) => r.defectNo === row.defectNo) : expandedRows
  if (expanded) {
    loadExpandDisposals(row)
  }
}

// 该处置单是否属于当前用户部门（本部门工单可改派处理人，其他部门只读）
function isMyDept(deptName: string) {
  return !!deptName && deptName === userStore.userInfo?.department
}

// 「一键自动分配」仅对运维部门（该缺陷 maint_dept=本部门）且有未分配工单时显示
function canAutoAssignDisposals(defect: any) {
  const disposals = disposalsOf(defect)
  return (
    defect.maintDept === userStore.userInfo?.department &&
    disposals.some((d: any) => d.status === 'pending_assign' && !d.assigneeUsername && !d.assigneeName)
  )
}

// 历史动作中文名
const HISTORY_ACTION_LABELS: Record<string, string> = {
  create: '创建',
  assign_analyst: '分配研判人',
  auto_assign_analyst: '自动分配研判人',
  analyze: '研判',
  review: '评审（上会）',
  auto_assign: '自动分配处理人',
  assign: '分配处理人',
  complete: '处置完成',
  na: '标记忽略',
  ignore: '忽略',
  close: '闭环',
}

async function viewDefectDetail(row: any) {
  currentDefect.value = row
  defectDetailDialogVisible.value = true
  defectHistory.value = []
  try {
    defectHistory.value = ((await riskApi.getDefectHistory(row.defectNo)) as any) || []
  } catch (e) {
    // 忽略
  }
}

// ===== 分配研判人（第一次分配） =====
const analystDialogVisible = ref(false)
const analystCandidates = ref<any[]>([])
const analystForm = ref({ selected: '', manual: '' })

async function openAnalystDialog() {
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
    defectDetailDialogVisible.value = false
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
    defectDetailDialogVisible.value = false
    loadDefects()
  } catch (e: any) {
    ElMessage.error('自动分配失败: ' + (e.message || ''))
  } finally {
    submitting.value = false
  }
}

// 一键自动分配处置单（各部门匹配项目主岗，仅运维部门在展开行操作）
async function autoAssignDisposals(defect: any) {
  if (!defect) return
  submitting.value = true
  try {
    const res: any = await riskApi.autoAssignDisposals(defect.defectNo)
    const assigned = (res?.assigned || []).join('、')
    const missed = (res?.missed || []).join('、')
    if (assigned) ElMessage.success(`已自动分配：${assigned}`)
    if (missed) ElMessage.warning(`未匹配到主岗（需手动分配）：${missed}`)
    // 刷新该缺陷展开行的处置单与缺陷列表
    await loadExpandDisposals(defect, true)
    loadDefects()
  } catch (e: any) {
    ElMessage.error('自动分配失败: ' + (e.message || ''))
  } finally {
    submitting.value = false
  }
}

// 补分配 / 改派处理人（展开行处置工单表格内「处理人」列为可编辑下拉，仅本部门工单可改）
const deptUserMap = ref<Record<string, any[]>>({})

// 按处置部门加载部门用户（供下拉选项使用）
async function loadDisposalDeptUsers(disposals: any[]) {
  const depts = Array.from(new Set(disposals.map(d => d.deptName).filter(Boolean)))
  for (const dept of depts) {
    if (deptUserMap.value[dept]) continue
    try {
      const users = ((await personnelApi.getList({ dept_name: dept })) as any) || []
      deptUserMap.value = { ...deptUserMap.value, [dept]: users }
    } catch (e) {
      deptUserMap.value = { ...deptUserMap.value, [dept]: [] }
    }
  }
}

// 下拉选中处理人 -> 补分配/改派（仅本部门工单，后端 Assign 二次校验）
async function changeAssignee(disposal: any, username: string) {
  if (!username) return
  const user = (deptUserMap.value[disposal.deptName] || []).find(u => u.username === username)
  submitting.value = true
  try {
    await riskApi.assignDisposal(disposal.id, {
      assigneeUsername: username,
      assigneeName: user?.nickname || username,
    })
    // 同步更新展开行缓存中的处理人姓名
    const nickname = user?.nickname || username
    for (const key of Object.keys(disposalsMap.value)) {
      const list = disposalsMap.value[Number(key)]
      const idx = list.findIndex((d: any) => d.id === disposal.id)
      if (idx >= 0) {
        list[idx] = { ...list[idx], assigneeUsername: username, assigneeName: nickname }
        disposalsMap.value = { ...disposalsMap.value, [Number(key)]: list }
        break
      }
    }
    ElMessage.success('处理人已更新')
    loadDefects()
  } catch (e: any) {
    ElMessage.error('更新失败: ' + (e.message || ''))
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  // 支持从「风险上报」创建成功跳转定位：localStorage 标记直接定位到缺陷管理-内部缺陷
  if (localStorage.getItem('osm_defect_tab') === 'defects') {
    activeTab.value = 'defects'
    defectActiveTab.value = 'internal'
    localStorage.removeItem('osm_defect_tab')
  }
  loadSoftwareList()
  loadDefects()
})
</script>

<style scoped>
.pagination-row { display: flex; justify-content: flex-end; margin-top: 12px; }
.my-ops { padding: 20px; }
.page-header { margin-bottom: 20px; }
.page-header h2 { margin: 0 0 8px; font-size: 18px; font-weight: 600; color: #1f2328; }
.page-desc { margin: 0; font-size: 13px; color: #8c8c8c; }
.ops-tabs { background: #fff; padding: 20px; border-radius: 4px; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08); }
.action-bar { margin: 16px 0; }
.section-title {
  font-size: 14px; font-weight: 600; color: #303133;
  margin: 16px 0 0; padding-left: 8px; border-left: 3px solid #409EFF;
}
.text-muted { color: #c0c4cc; }
/* 长文本（英文正文/中文释义/中文建议）保留换行展示 */
.pre-line {
  white-space: pre-wrap; word-break: break-word; line-height: 1.6; max-height: 320px; overflow-y: auto;
}
</style>
