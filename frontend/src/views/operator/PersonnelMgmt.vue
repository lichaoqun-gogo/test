<template>
  <div class="personnel-mgmt">
    <div class="page-header">
      <h2>人员管理</h2>
      <p class="page-desc">维护开源基础软件、自研项目的负责人关联信息（仅本部门）</p>
    </div>

    <el-tabs v-model="activeTab" class="personnel-tabs" @tab-change="handleTabChange">
      <!-- 开源基础软件负责人 -->
      <el-tab-pane label="开源基础软件负责人" name="software">
        <div class="tab-header">
          <el-input
            v-model="softwareSearchText"
            placeholder="搜索软件名称"
            style="width: 300px"
            clearable
            @input="filterSoftwareList"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>

        <el-table :data="filteredSoftwareList" border style="margin-top: 16px" v-loading="softwareLoading">
          <el-table-column prop="_name" label="软件名称" min-width="160" />
          <el-table-column prop="_version" label="当前版本" width="120" />
          <el-table-column prop="type" label="类型" width="120">
            <template #default="{ row }">
              <el-tag size="small">{{ row.type || '-' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="当前负责人" min-width="180">
            <template #default="{ row }">
              {{ formatOwners(row.id) }}
            </template>
          </el-table-column>
          <el-table-column prop="maint_dept" label="维护部门" width="140" />
          <el-table-column label="操作" width="120" align="center">
            <template #default="{ row }">
              <el-button type="primary" link v-auth="'software_owner_update'" @click="showEditSoftwareOwner(row)">
                编辑
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- 自研项目负责人 -->
      <el-tab-pane label="自研项目负责人" name="project">
        <div class="tab-header">
          <el-input
            v-model="projectSearchText"
            placeholder="搜索项目名称"
            style="width: 300px"
            clearable
            @input="filterProjectList"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>

        <el-table :data="filteredProjectList" border style="margin-top: 16px" v-loading="projectLoading">
          <el-table-column prop="pro_code" label="项目名称" min-width="160" />
          <el-table-column prop="group_name" label="团队" width="140" />
          <el-table-column label="当前负责人" min-width="180">
            <template #default="{ row }">
              {{ formatProjectOwners(row.id) }}
            </template>
          </el-table-column>
          <el-table-column prop="dept_name" label="部门" width="140" />
          <el-table-column label="操作" width="120" align="center">
            <template #default="{ row }">
              <el-button type="primary" link v-auth="'project_owner_update'" @click="showEditProjectOwner(row)">
                编辑
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- 人员字典 -->
      <el-tab-pane label="人员字典" name="dictionary">
        <el-table :data="personDictionary" border style="margin-top: 16px" v-loading="personLoading">
          <el-table-column prop="username" label="账号" width="140" />
          <el-table-column prop="nickname" label="姓名" width="120" />
          <el-table-column prop="department" label="部门" width="160" />
          <el-table-column prop="team" label="团队" width="140" />
          <el-table-column prop="email" label="邮箱" min-width="180" />
          <el-table-column label="负责软件数" width="120" align="center">
            <template #default="{ row }">
              {{ getPersonSoftwareCount(row.username) }}
            </template>
          </el-table-column>
          <el-table-column label="负责项目数" width="120" align="center">
            <template #default="{ row }">
              {{ getPersonProjectCount(row.username) }}
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <!-- 编辑软件负责人对话框 -->
    <el-dialog v-model="editSoftwareOwnerVisible" title="编辑软件负责人" width="600px">
      <el-form :model="editSoftwareForm" label-width="100px">
        <el-form-item label="软件名称">
          <el-input v-model="editSoftwareForm.name" disabled />
        </el-form-item>
        <el-form-item label="当前负责人">
          <el-input v-model="editSoftwareForm.ownerDisplay" disabled />
        </el-form-item>
        <el-form-item label="新负责人">
          <el-select
            v-model="editSoftwareForm.newOwners"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="请选择或输入负责人账号"
            style="width: 100%"
            :loading="personLoading"
          >
            <el-option
              v-for="person in personDictionary"
              :key="person.id"
              :label="`${person.nickname}（${person.department}${person.team ? '/' + person.team : ''}）`"
              :value="person.username"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editSoftwareOwnerVisible = false">取消</el-button>
        <el-button type="primary" :loading="saveLoading" @click="confirmEditSoftwareOwner">确认</el-button>
      </template>
    </el-dialog>

    <!-- 编辑项目负责人对话框 -->
    <el-dialog v-model="editProjectOwnerVisible" title="编辑项目负责人" width="600px">
      <el-form :model="editProjectForm" label-width="100px">
        <el-form-item label="项目名称">
          <el-input v-model="editProjectForm.name" disabled />
        </el-form-item>
        <el-form-item label="当前负责人">
          <el-input v-model="editProjectForm.ownerDisplay" disabled />
        </el-form-item>
        <el-form-item label="新负责人">
          <el-select
            v-model="editProjectForm.newOwners"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="请选择或输入负责人账号"
            style="width: 100%"
            :loading="personLoading"
          >
            <el-option
              v-for="person in personDictionary"
              :key="person.id"
              :label="`${person.nickname}（${person.department}${person.team ? '/' + person.team : ''}）`"
              :value="person.username"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editProjectOwnerVisible = false">取消</el-button>
        <el-button type="primary" :loading="saveLoading" @click="confirmEditProjectOwner">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { assetApi, getInnerProjects, personnelApi, projectOwnerApi } from '@/api'

const userStore = useUserStore()
const activeTab = ref('software')

// 人员字典（从后端 users 表查询，只显示本部门）
const personDictionary = ref<any[]>([])
const personLoading = ref(false)

// 全部负责人关联数据（用于统计每个人的负责软件数/项目数）
const allSoftwareOwners = ref<Record<number, string[]>>({})
const allProjectOwners = ref<Record<number, string[]>>({})

// 加载人员列表（只查本部门人员）
async function loadPersonnel() {
  personLoading.value = true
  try {
    const deptName = userStore.userInfo?.department || ''
    const data = await personnelApi.getList(deptName ? { dept_name: deptName } : undefined)
    personDictionary.value = Array.isArray(data) ? data : []

    // 同时加载本部门所有软件/项目的负责人关联（用于统计负责数量）
    await loadAllOwnersForStats()
  } catch (err: any) {
    console.error('[PersonnelMgmt] 加载人员列表失败:', err)
    ElMessage.error('加载人员列表失败: ' + (err.message || '未知错误'))
  } finally {
    personLoading.value = false
  }
}

// 加载本部门所有负责人关联数据（用于人员字典统计）
async function loadAllOwnersForStats() {
  const deptName = userStore.userInfo?.department || ''
  if (!deptName) return

  // 加载本部门的所有软件
  try {
    const swData = await assetApi.getProjects({ category: 'software', maint_dept: deptName })
    const swList = Array.isArray(swData) ? swData : []
    if (swList.length > 0) {
      const swIds = swList.map(s => Number(s.id)).filter(Boolean)
      allSoftwareOwners.value = (await projectOwnerApi.getOwnersMap('software', swIds)) || {}
    }
  } catch (err) {
    console.error('[PersonnelMgmt] 加载软件负责人统计失败:', err)
  }

  // 加载本部门的所有自研项目
  try {
    const pjData = await getInnerProjects({ dept_name: deptName })
    const pjList = Array.isArray(pjData) ? pjData : []
    if (pjList.length > 0) {
      const pjIds = pjList.map(p => Number(p.id)).filter(Boolean)
      allProjectOwners.value = (await projectOwnerApi.getOwnersMap('inner', pjIds)) || {}
    }
  } catch (err) {
    console.error('[PersonnelMgmt] 加载项目负责人统计失败:', err)
  }
}

// 统计某人负责的软件数
function getPersonSoftwareCount(username: string): number {
  let count = 0
  for (const usernames of Object.values(allSoftwareOwners.value)) {
    if (usernames.includes(username)) count++
  }
  return count
}

// 统计某人负责的项目数
function getPersonProjectCount(username: string): number {
  let count = 0
  for (const usernames of Object.values(allProjectOwners.value)) {
    if (usernames.includes(username)) count++
  }
  return count
}

// 根据账号查找人员信息
function findPerson(username: string) {
  return personDictionary.value.find(p => p.username === username)
}

// ========== 开源基础软件 ==========
const softwareList = ref<any[]>([])
const softwareOwnersMap = ref<Record<number, string[]>>({})
const softwareLoading = ref(false)
const softwareSearchText = ref('')
const filteredSoftwareList = ref<any[]>([])

async function loadSoftwareList() {
  const deptName = userStore.userInfo?.department || ''
  if (!deptName) {
    ElMessage.warning('未获取到部门信息')
    return
  }
  softwareLoading.value = true
  try {
    // 查询本部门维护的开源基础软件
    const data = await assetApi.getProjects({ category: 'software', maint_dept: deptName })
    softwareList.value = Array.isArray(data) ? data : []
    filteredSoftwareList.value = softwareList.value

    // 批量加载所有软件的负责人
    if (softwareList.value.length > 0) {
      const ids = softwareList.value.map(s => Number(s.id)).filter(Boolean)
      const map = await projectOwnerApi.getOwnersMap('software', ids)
      softwareOwnersMap.value = (map || {}) as Record<number, string[]>
    } else {
      softwareOwnersMap.value = {}
    }
  } catch (err: any) {
    console.error('[PersonnelMgmt] 加载软件列表失败:', err)
    ElMessage.error('加载软件列表失败: ' + (err.message || '未知错误'))
  } finally {
    softwareLoading.value = false
  }
}

function filterSoftwareList() {
  if (!softwareSearchText.value) {
    filteredSoftwareList.value = softwareList.value
    return
  }
  const text = softwareSearchText.value.toLowerCase()
  filteredSoftwareList.value = softwareList.value.filter(sw =>
    (sw._name || '').toLowerCase().includes(text)
  )
}

function formatOwners(projectId: number): string {
  const usernames = softwareOwnersMap.value[projectId] || []
  if (usernames.length === 0) return '-'
  return usernames.map(u => {
    const p = findPerson(u)
    return p ? `${p.nickname}(${p.username})` : u
  }).join('、')
}

// 编辑软件负责人
const editSoftwareOwnerVisible = ref(false)
const editSoftwareForm = ref({
  id: 0,
  name: '',
  ownerDisplay: '',
  newOwners: [] as string[],
})

function showEditSoftwareOwner(row: any) {
  const existingOwners = softwareOwnersMap.value[row.id] || []
  editSoftwareForm.value = {
    id: row.id,
    name: row._name || '',
    ownerDisplay: formatOwners(row.id),
    newOwners: [...existingOwners],
  }
  editSoftwareOwnerVisible.value = true
}

async function confirmEditSoftwareOwner() {
  saveLoading.value = true
  try {
    const owners = editSoftwareForm.value.newOwners.map(username => {
      const p = findPerson(username)
      return {
        owner_username: username,
        dept_name: p?.department || '',
        team: p?.team || '',
      }
    })
    await projectOwnerApi.saveOwners('software', editSoftwareForm.value.id, owners)
    // 更新本地缓存（台账 + 人员字典统计）
    softwareOwnersMap.value[editSoftwareForm.value.id] = editSoftwareForm.value.newOwners
    allSoftwareOwners.value[editSoftwareForm.value.id] = editSoftwareForm.value.newOwners
    editSoftwareForm.value.ownerDisplay = formatOwners(editSoftwareForm.value.id)
    editSoftwareOwnerVisible.value = false
    ElMessage.success('负责人更新成功')
  } catch (err: any) {
    ElMessage.error('保存失败: ' + (err.message || '未知错误'))
  } finally {
    saveLoading.value = false
  }
}

// ========== 自研项目 ==========
const projectList = ref<any[]>([])
const projectOwnersMap = ref<Record<number, string[]>>({})
const projectLoading = ref(false)
const projectSearchText = ref('')
const filteredProjectList = ref<any[]>([])

async function loadProjects() {
  const deptName = userStore.userInfo?.department || ''
  if (!deptName) {
    ElMessage.warning('未获取到部门信息')
    return
  }
  projectLoading.value = true
  try {
    const innerData = await getInnerProjects({ dept_name: deptName })
    const rawList = Array.isArray(innerData) ? innerData : []
    // 按 pro_code 去重
    const seen = new Set<string>()
    projectList.value = rawList.filter((item: any) => {
      const key = item.pro_code || ''
      if (seen.has(key)) return false
      seen.add(key)
      return true
    })
    filteredProjectList.value = projectList.value

    // 批量加载负责人
    if (projectList.value.length > 0) {
      const ids = projectList.value.map(p => Number(p.id)).filter(Boolean)
      const map = await projectOwnerApi.getOwnersMap('inner', ids)
      projectOwnersMap.value = (map || {}) as Record<number, string[]>
    } else {
      projectOwnersMap.value = {}
    }
  } catch (err: any) {
    console.error('[PersonnelMgmt] 加载项目失败:', err)
    ElMessage.error('加载项目失败: ' + (err.message || '未知错误'))
  } finally {
    projectLoading.value = false
  }
}

function filterProjectList() {
  if (!projectSearchText.value) {
    filteredProjectList.value = projectList.value
    return
  }
  const text = projectSearchText.value.toLowerCase()
  filteredProjectList.value = projectList.value.filter(pj =>
    (pj.pro_code || '').toLowerCase().includes(text) ||
    (pj.module_code || '').toLowerCase().includes(text)
  )
}

// 编辑项目负责人
const editProjectOwnerVisible = ref(false)
const editProjectForm = ref({
  id: 0,
  name: '',
  ownerDisplay: '',
  newOwners: [] as string[],
})
const saveLoading = ref(false)

function showEditProjectOwner(row: any) {
  const existingOwners = projectOwnersMap.value[row.id] || []
  editProjectForm.value = {
    id: row.id,
    name: row.pro_code || row.module_code || '',
    ownerDisplay: formatProjectOwners(row.id),
    newOwners: [...existingOwners],
  }
  editProjectOwnerVisible.value = true
}

function formatProjectOwners(projectId: number): string {
  const usernames = projectOwnersMap.value[projectId] || []
  if (usernames.length === 0) return '-'
  return usernames.map(u => {
    const p = findPerson(u)
    return p ? `${p.nickname}(${p.username})` : u
  }).join('、')
}

async function confirmEditProjectOwner() {
  saveLoading.value = true
  try {
    const owners = editProjectForm.value.newOwners.map(username => {
      const p = findPerson(username)
      return {
        owner_username: username,
        dept_name: p?.department || '',
        team: p?.team || '',
      }
    })
    await projectOwnerApi.saveOwners('inner', editProjectForm.value.id, owners)
    // 更新本地缓存（台账 + 人员字典统计）
    projectOwnersMap.value[editProjectForm.value.id] = editProjectForm.value.newOwners
    allProjectOwners.value[editProjectForm.value.id] = editProjectForm.value.newOwners
    editProjectForm.value.ownerDisplay = formatProjectOwners(editProjectForm.value.id)
    editProjectOwnerVisible.value = false
    ElMessage.success('负责人更新成功')
  } catch (err: any) {
    ElMessage.error('保存失败: ' + (err.message || '未知错误'))
  } finally {
    saveLoading.value = false
  }
}

// tab 切换时按需加载
function handleTabChange(tab: string) {
  if (tab === 'software' && softwareList.value.length === 0) {
    loadSoftwareList()
  } else if (tab === 'project' && projectList.value.length === 0) {
    loadProjects()
  }
}

onMounted(() => {
  loadPersonnel()
  loadSoftwareList() // 默认加载第一个 tab
})
</script>

<style scoped>
.personnel-mgmt {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 600;
  color: #1f2328;
}

.page-desc {
  margin: 0;
  font-size: 13px;
  color: #8c8c8c;
}

.personnel-tabs {
  background: #fff;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
}

.tab-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
