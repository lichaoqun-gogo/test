<template>
  <div class="blacklist-manage">
    <div class="page-header">
      <div class="header-left">
        <h3>名单配置管理</h3>
      </div>
    </div>

    <div class="toolbar">
      <el-input v-model="searchText" placeholder="搜索名称" clearable style="width: 240px" @keyup.enter="handleSearch" @clear="handleSearch" />
      <el-select v-model="filterCategory" placeholder="分类" clearable style="width: 160px" @change="handleSearch">
        <el-option v-for="cat in categoryOptions" :key="cat" :label="cat" :value="cat" />
      </el-select>
      <el-select v-model="filterLevel" placeholder="适用等级" clearable style="width: 120px" @change="handleSearch">
        <el-option label="A" value="A" />
        <el-option label="B" value="B" />
        <el-option label="C" value="C" />
        <el-option label="D" value="D" />
      </el-select>
      <el-select v-model="filterEnable" placeholder="状态" clearable style="width: 120px" @change="handleSearch">
        <el-option label="启用" :value="1" />
        <el-option label="禁用" :value="0" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <el-button @click="handleReset">重置</el-button>
      <div class="toolbar-right">
        <el-button type="primary" @click="openAdd"><el-icon><Plus /></el-icon> 新增</el-button>
      </div>
    </div>

    <el-table :data="blacklist" stripe style="margin-top: 16px" row-key="id" v-loading="loading">
      <el-table-column type="expand">
        <template #default="{ row }">
          <div class="expand-content">
            <div class="expand-row">
              <span class="expand-label">仓库类型：</span>
              <span class="expand-value">{{ formatLabel(row._format) }}</span>
            </div>
            <div class="expand-row">
              <span class="expand-label">名单理由：</span>
              <span class="expand-value">{{ row.reason || '-' }}</span>
            </div>
            <div class="expand-row">
              <span class="expand-label">建议：</span>
              <span class="expand-value">{{ row.advise || '-' }}</span>
            </div>
            <el-row :gutter="24" class="expand-row">
              <el-col :span="12">
                <span class="expand-label">创建者：</span>
                <span class="expand-value">{{ row.create_by || '-' }}</span>
              </el-col>
              <el-col :span="12">
                <span class="expand-label">创建时间：</span>
                <span class="expand-value">{{ formatTime(row.create_time) }}</span>
              </el-col>
            </el-row>
            <el-row :gutter="24" class="expand-row">
              <el-col :span="12">
                <span class="expand-label">更新者：</span>
                <span class="expand-value">{{ row.update_by || '-' }}</span>
              </el-col>
              <el-col :span="12">
                <span class="expand-label">更新时间：</span>
                <span class="expand-value">{{ formatTime(row.update_time) }}</span>
              </el-col>
            </el-row>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="type" label="分类" width="140" />
      <el-table-column prop="_group" label="组织" width="160" />
      <el-table-column prop="_name" label="名称" min-width="140" />
      <el-table-column prop="applicable_level" label="适用等级" width="100">
        <template #default="{ row }">
          <el-tag :type="getLevelTag(row.applicable_level) as any" size="small">{{ row.applicable_level || '-' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="alert_type" label="告警类型" width="100">
        <template #default="{ row }">
          <el-tag :type="row.alert_type === 1 ? 'danger' : 'warning'" size="small">
            {{ row.alert_type === 1 ? '阻断' : '警告' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="version_exp" label="表达式" width="200" show-overflow-tooltip />
      <el-table-column prop="support_team" label="技术支持团队" width="130" />
      <el-table-column label="是否启用" width="90">
        <template #default="{ row }">
          <el-tag :type="row._enable === 1 ? 'success' : 'info'" size="small">
            {{ row._enable === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="openEdit(row as any)">编辑</el-button>
          <el-button type="danger" link size="small" @click="handleDelete(row as any)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadData"
        @size-change="loadData"
      />
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑名单' : '新增名单'" width="700px">
      <el-form :model="form" label-width="110px">
        <el-form-item label="名称" :required="!isEdit">
          <el-input v-model="form._name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="组织">
          <el-input v-model="form._group" placeholder="请输入组织" />
        </el-form-item>
        <el-form-item label="仓库类型" :required="!isEdit">
          <el-select v-model="form._format" placeholder="请选择" style="width: 100%">
            <el-option label="Maven" :value="1" />
            <el-option label="Pypi" :value="2" />
            <el-option label="Npm" :value="3" />
            <el-option label="Yum" :value="4" />
            <el-option label="GO_MODULES" :value="5" />
          </el-select>
        </el-form-item>
        <el-form-item label="适用等级" :required="!isEdit">
          <el-select v-model="form.applicable_level" placeholder="请选择" style="width: 100%">
            <el-option label="A - 可安全使用(低风险)" value="A" />
            <el-option label="B - 谨慎使用(中低风险)" value="B" />
            <el-option label="C - 限期整改(高风险)" value="C" />
            <el-option label="D - 禁止使用(高风险)" value="D" />
          </el-select>
        </el-form-item>
        <el-form-item label="分类" :required="!isEdit">
          <el-select v-model="form.type" placeholder="请选择" style="width: 100%">
            <el-option v-for="cat in categoryOptions" :key="cat" :label="cat" :value="cat" />
          </el-select>
        </el-form-item>
        <el-form-item label="匹配表达式">
          <div class="version-exp-container">
            <div v-for="(group, idx) in form.versionGroups" :key="idx" class="version-group">
              <div class="version-range">
                <span class="range-label">不包含</span>
                <el-switch v-model="group.startInclusive" :width="40" />
                <span class="range-label">包含</span>
                <el-input v-model="group.startVersion" placeholder="版本号" style="width: 120px" />
                <span class="range-separator">,</span>
                <el-input v-model="group.endVersion" placeholder="版本号" style="width: 120px" />
                <span class="range-label">不包含</span>
                <el-switch v-model="group.endInclusive" :width="40" />
                <span class="range-label">包含</span>
                <el-button v-if="idx === 0" type="primary" link @click="addVersionGroup" style="margin-left: 8px">
                  <el-icon><Plus /></el-icon>
                </el-button>
                <el-button v-else type="danger" link @click="removeVersionGroup(idx)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="技术支持团队" :required="!isEdit">
          <el-input v-model="form.support_team" placeholder="请输入团队名称" />
        </el-form-item>
        <el-form-item label="名单理由" :required="!isEdit">
          <el-input v-model="form.reason" type="textarea" :rows="2" placeholder="请说明加入名单的理由" />
        </el-form-item>
        <el-form-item label="建议">
          <el-input v-model="form.advise" type="textarea" :rows="2" placeholder="请给出处理建议" />
        </el-form-item>
        <el-form-item label="是否启用">
          <span class="range-label">禁用</span>
          <el-switch v-model="form._enable" :active-value="1" :inactive-value="0" active-color="#13ce66" style="margin: 0 8px" />
          <span class="range-label">启用</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatTime } from '@/utils/format'
import api from '@/api'

interface ComponentCoreItem {
  id: number
  _format: number
  type?: string
  _name: string
  _group: string
  version_exp: string
  applicable_level: string
  support_team: string
  reason?: string
  advise?: string
  _enable: number
  alert_type?: number
  create_time?: string
  create_by?: string
  update_time?: string
  update_by?: string
  is_boundary?: number
}

interface VersionGroup {
  startInclusive: boolean
  startVersion: string
  endVersion: string
  endInclusive: boolean
}

interface BlacklistForm {
  id?: number
  _format: number
  type: string
  _name: string
  _group: string
  applicable_level: string
  support_team: string
  reason: string
  advise: string
  _enable: number
  versionGroups: VersionGroup[]
}

const loading = ref(false)
const saving = ref(false)
const blacklist = ref<ComponentCoreItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchText = ref('')
const filterCategory = ref('')
const filterLevel = ref('')
const filterEnable = ref<number | ''>('')

const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)

const categoryOptions = ['安全算法类', '服务治理类', '通用工具包', '基础框架包', '日志组件', '数据库相关', '序列化工具', '中间件相关']

const defaultForm: BlacklistForm = {
  _format: 1,
  type: '',
  _name: '',
  _group: '',
  applicable_level: '',
  support_team: '',
  reason: '',
  advise: '',
  _enable: 1,
  versionGroups: [{ startInclusive: false, startVersion: '', endVersion: '', endInclusive: false }],
}

const form = ref<BlacklistForm>({ ...defaultForm })

function getLevelTag(level: string) {
  const map: Record<string, string> = { A: 'success', B: 'warning', C: 'danger', D: 'danger' }
  return map[level] || 'info'
}

function formatLabel(format: number): string {
  const map: Record<number, string> = { 1: 'Maven', 2: 'Pypi', 3: 'Npm', 4: 'Yum', 5: 'GO_MODULES' }
  return map[format] || '未知'
}

async function loadData() {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      pageSize: pageSize.value,
    }
    if (searchText.value) params.keyword = searchText.value
    if (filterCategory.value) params.type = filterCategory.value
    if (filterLevel.value) params.applicable_level = filterLevel.value
    if (filterEnable.value !== '' && filterEnable.value !== null) params._enable = filterEnable.value

    const res: any = await api.get('/risk/blacklist', { params })
    blacklist.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('加载名单失败:', error)
    ElMessage.error('加载名单失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  currentPage.value = 1
  loadData()
}

function handleReset() {
  searchText.value = ''
  filterCategory.value = ''
  filterLevel.value = ''
  filterEnable.value = ''
  currentPage.value = 1
  loadData()
}

function parseVersionExpToGroups(versionExp: string): VersionGroup[] {
  if (!versionExp) return [{ startInclusive: false, startVersion: '', endVersion: '', endInclusive: false }]
  const groups: VersionGroup[] = []
  const regex = /([[(])([^,]*),([^)\]]*)([)\]])/g
  let match
  while ((match = regex.exec(versionExp)) !== null) {
    groups.push({
      startInclusive: match[1] === '[',
      startVersion: match[2] === '*' ? '' : match[2],
      endVersion: match[3] === '*' ? '' : match[3],
      endInclusive: match[4] === ']',
    })
  }
  return groups.length > 0 ? groups : [{ startInclusive: false, startVersion: '', endVersion: '', endInclusive: false }]
}

function versionGroupsToExp(groups: VersionGroup[]): string {
  return groups
    .map(g => {
      const start = g.startInclusive ? '[' : '('
      const end = g.endInclusive ? ']' : ')'
      const startVer = g.startVersion || '*'
      const endVer = g.endVersion || '*'
      return `${start}${startVer},${endVer}${end}`
    })
    .join('&&')
}

function openAdd() {
  isEdit.value = false
  form.value = { ...defaultForm, versionGroups: [{ startInclusive: false, startVersion: '', endVersion: '', endInclusive: false }] }
  dialogVisible.value = true
}

function openEdit(row: ComponentCoreItem) {
  isEdit.value = true
  editingId.value = row.id
  form.value = {
    id: row.id,
    _format: row._format,
    type: row.type || '',
    _name: row._name,
    _group: row._group || '',
    applicable_level: row.applicable_level || '',
    support_team: row.support_team || '',
    reason: row.reason || '',
    advise: row.advise || '',
    _enable: row._enable,
    versionGroups: parseVersionExpToGroups(row.version_exp),
  }
  dialogVisible.value = true
}

function addVersionGroup() {
  form.value.versionGroups.push({ startInclusive: false, startVersion: '', endVersion: '', endInclusive: false })
}

function removeVersionGroup(idx: number) {
  form.value.versionGroups.splice(idx, 1)
}

async function handleSave() {
  if (!isEdit.value && (!form.value._name || !form.value.type || !form.value.support_team || !form.value.reason || !form.value.applicable_level)) {
    ElMessage.warning('请填写必填项')
    return
  }
  const versionExp = versionGroupsToExp(form.value.versionGroups)
  const payload: any = {
    _format: form.value._format,
    type: form.value.type,
    _name: form.value._name,
    _group: form.value._group,
    version_exp: versionExp,
    applicable_level: form.value.applicable_level,
    support_team: form.value.support_team,
    reason: form.value.reason,
    advise: form.value.advise,
    _enable: form.value._enable,
  }

  saving.value = true
  try {
    if (isEdit.value) {
      await api.put(`/risk/blacklist/${editingId.value}`, payload)
      ElMessage.success('编辑成功')
      dialogVisible.value = false
      loadData()
    } else {
      await api.post('/risk/blacklist', payload)
      ElMessage.success('新增成功')
      dialogVisible.value = false
      loadData()
    }
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

async function handleDelete(row: ComponentCoreItem) {
  try {
    await ElMessageBox.confirm(`确认删除「${row._name}」吗？`, '提示', { type: 'warning' })
    await api.delete(`/risk/blacklist/${row.id}`)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    // 用户取消删除
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-header {
  margin-bottom: 8px;
}

.page-header h3 {
  margin: 0;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.toolbar-right {
  margin-left: auto;
  display: flex;
  gap: 8px;
}

.expand-content {
  padding: 12px 24px;
}

.expand-row {
  line-height: 28px;
  font-size: 13px;
}

.expand-label {
  color: #909399;
  font-weight: 500;
}

.expand-value {
  color: #303133;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.version-exp-container {
  width: 100%;
}

.version-group {
  margin-bottom: 8px;
}

.version-group:last-child {
  margin-bottom: 0;
}

.version-range {
  display: flex;
  align-items: center;
  gap: 8px;
}

.range-label {
  font-size: 12px;
  color: #909399;
  min-width: 40px;
}

.range-separator {
  color: #909399;
}
</style>
