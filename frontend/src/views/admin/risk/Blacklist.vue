<template>
  <div class="blacklist-manage">
    <div class="page-header">
      <div class="header-left">
        <h3>名单配置管理</h3>
      </div>
    </div>

    <div class="toolbar">
      <el-input v-model="searchText" placeholder="搜索名称" clearable style="width: 240px" />
      <el-select v-model="filterCategory" placeholder="分类" clearable style="width: 160px">
        <el-option v-for="cat in categoryOptions" :key="cat" :label="cat" :value="cat" />
      </el-select>
      <el-select v-model="filterLevel" placeholder="适用等级" clearable style="width: 120px">
        <el-option label="A" value="A" />
        <el-option label="B" value="B" />
        <el-option label="C" value="C" />
        <el-option label="D" value="D" />
      </el-select>
      <el-select v-model="filterEnable" placeholder="状态" clearable style="width: 120px">
        <el-option label="启用" :value="1" />
        <el-option label="禁用" :value="0" />
      </el-select>
      <div class="toolbar-right">
        <el-button type="primary" @click="openAdd"><el-icon><Plus /></el-icon> 新增</el-button>
      </div>
    </div>

    <el-table :data="filteredBlacklist" stripe style="margin-top: 16px" row-key="id">
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
                <span class="expand-value">{{ row.create_time || '-' }}</span>
              </el-col>
            </el-row>
            <el-row :gutter="24" class="expand-row">
              <el-col :span="12">
                <span class="expand-label">更新者：</span>
                <span class="expand-value">{{ row.update_by || '-' }}</span>
              </el-col>
              <el-col :span="12">
                <span class="expand-label">更新时间：</span>
                <span class="expand-value">{{ row.update_time || '-' }}</span>
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
          <el-tag :type="getLevelTag(row.applicable_level)" size="small">{{ row.applicable_level || '-' }}</el-tag>
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
      <el-table-column label="操作" width="80" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="openEdit(row as ComponentCoreItem)">编辑</el-button>
        </template>
      </el-table-column>
    </el-table>

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
        <el-button type="primary" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { mockComponentCores, type ComponentCoreItem } from '@/mock/admin'
import { Plus, Delete } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const searchText = ref('')
const filterCategory = ref('')
const filterLevel = ref('')
const filterEnable = ref<number | ''>('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)

const blacklist = ref<ComponentCoreItem[]>([...mockComponentCores])

const categoryOptions = ['安全算法类', '服务治理类', '通用工具包', '基础框架包', '日志组件', '数据库相关', '序列化工具', '中间件相关']

const filteredBlacklist = computed(() => {
  return blacklist.value.filter(item => {
    const matchSearch = !searchText.value || item._name.toLowerCase().includes(searchText.value.toLowerCase())
    const matchCategory = filterCategory.value ? item.type === filterCategory.value : true
    const matchLevel = filterLevel.value ? item.applicable_level === filterLevel.value : true
    const matchEnable = (filterEnable.value === '' || filterEnable.value === null || filterEnable.value === undefined) ? true : item._enable === filterEnable.value
    return matchSearch && matchCategory && matchLevel && matchEnable
  })
})

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

function getLevelTag(level: string): string {
  const map: Record<string, string> = { A: 'success', B: 'warning', C: 'danger', D: 'danger' }
  return map[level] || 'info'
}

function formatLabel(format: number): string {
  const map: Record<number, string> = { 1: 'Maven', 2: 'Pypi', 3: 'Npm', 4: 'Yum', 5: 'GO_MODULES' }
  return map[format] || '未知'
}

function handleSearch() {
  // computed会自动响应式更新，无需额外处理
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

function handleSave() {
  if (!isEdit.value && (!form.value._name || !form.value.type || !form.value.support_team || !form.value.reason || !form.value.applicable_level)) {
    ElMessage.warning('请填写必填项')
    return
  }
  const versionExp = versionGroupsToExp(form.value.versionGroups)
  if (isEdit.value) {
    const idx = blacklist.value.findIndex(item => item.id === editingId.value)
    if (idx !== -1) {
      blacklist.value[idx] = {
        ...blacklist.value[idx],
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
    }
    ElMessage.success('编辑成功')
  } else {
    const newItem: ComponentCoreItem = {
      id: blacklist.value.length + 1,
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
      is_boundary: 2,
      create_by: 'admin',
      create_time: new Date().toLocaleString('zh-CN'),
    }
    blacklist.value.unshift(newItem)
    ElMessage.success('新增成功')
  }
  dialogVisible.value = false
}
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
