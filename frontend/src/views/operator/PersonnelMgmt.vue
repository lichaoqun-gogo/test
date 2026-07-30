<template>
  <div class="personnel-mgmt">
    <div class="page-header">
      <h2>人员管理</h2>
      <p class="page-desc">维护开源基础软件、自研项目的负责人关联信息</p>
    </div>

    <el-tabs v-model="activeTab" class="personnel-tabs">
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

        <el-table :data="filteredSoftwareList" border style="margin-top: 16px">
          <el-table-column prop="name" label="软件名称" min-width="160" />
          <el-table-column prop="version" label="当前版本" width="120" />
          <el-table-column prop="type" label="类型" width="120">
            <template #default="{ row }">
              <el-tag size="small">{{ row.type }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="owner" label="当前负责人" min-width="180" />
          <el-table-column prop="team" label="所属团队" width="140" />
          <el-table-column label="操作" width="120" align="center">
            <template #default="{ row }">
              <el-button type="primary" link @click="showEditSoftwareOwner(row)">
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

        <el-table :data="filteredProjectList" border style="margin-top: 16px">
          <el-table-column prop="name" label="项目名称" min-width="160" />
          <el-table-column prop="team" label="团队" width="140" />
          <el-table-column prop="owner" label="当前负责人" min-width="180" />
          <el-table-column label="操作" width="120" align="center">
            <template #default="{ row }">
              <el-button type="primary" link @click="showEditProjectOwner(row)">
                编辑
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- 人员字典 -->
      <el-tab-pane label="人员字典" name="dictionary">
        <el-table :data="personDictionary" border style="margin-top: 16px">
          <el-table-column prop="id" label="人员ID" width="120" />
          <el-table-column prop="name" label="姓名" width="120" />
          <el-table-column prop="team" label="团队" width="140" />
          <el-table-column prop="email" label="邮箱" min-width="180" />
          <el-table-column prop="softwareCount" label="负责软件数" width="120" align="center" />
          <el-table-column prop="projectCount" label="负责项目数" width="120" align="center" />
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <!-- 编辑软件负责人对话框 -->
    <el-dialog
      v-model="editSoftwareOwnerVisible"
      title="编辑软件负责人"
      width="600px"
    >
      <el-form :model="editSoftwareForm" label-width="100px">
        <el-form-item label="软件名称">
          <el-input v-model="editSoftwareForm.name" disabled />
        </el-form-item>
        <el-form-item label="当前负责人">
          <el-input v-model="editSoftwareForm.owner" disabled />
        </el-form-item>
        <el-form-item label="新负责人">
          <el-select
            v-model="editSoftwareForm.newOwners"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="请选择或输入负责人"
            style="width: 100%"
          >
            <el-option
              v-for="person in personDictionary"
              :key="person.id"
              :label="`${person.name}（${person.team}）`"
              :value="person.name"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editSoftwareOwnerVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmEditSoftwareOwner">确认</el-button>
      </template>
    </el-dialog>

    <!-- 编辑项目负责人对话框 -->
    <el-dialog
      v-model="editProjectOwnerVisible"
      title="编辑项目负责人"
      width="600px"
    >
      <el-form :model="editProjectForm" label-width="100px">
        <el-form-item label="项目名称">
          <el-input v-model="editProjectForm.name" disabled />
        </el-form-item>
        <el-form-item label="当前负责人">
          <el-input v-model="editProjectForm.owner" disabled />
        </el-form-item>
        <el-form-item label="新负责人">
          <el-select
            v-model="editProjectForm.newOwners"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="请选择或输入负责人"
            style="width: 100%"
          >
            <el-option
              v-for="person in personDictionary"
              :key="person.id"
              :label="`${person.name}（${person.team}）`"
              :value="person.name"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editProjectOwnerVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmEditProjectOwner">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'

const activeTab = ref('software')

// 人员字典
const personDictionary = ref([
  { id: 'user001', name: '张三', team: '交易组', email: 'zhangsan@company.com', softwareCount: 4, projectCount: 2 },
  { id: 'user002', name: '李四', team: '支付组', email: 'lisi@company.com', softwareCount: 3, projectCount: 1 },
  { id: 'user003', name: '王五', team: '支付二组', email: 'wangwu@company.com', softwareCount: 2, projectCount: 3 },
  { id: 'user004', name: '赵六', team: '支付三组', email: 'zhaoliu@company.com', softwareCount: 1, projectCount: 0 },
  { id: 'user005', name: '钱七', team: '交易组', email: 'qianqi@company.com', softwareCount: 0, projectCount: 4 },
])

// 开源基础软件列表
const softwareList = ref([
  { id: 'sw001', name: 'Redis', version: '7.2.4', type: '基础软件', owner: '张三', team: '交易组' },
  { id: 'sw002', name: 'Spring Boot', version: '3.2.3', type: '依赖包', owner: '张三', team: '交易组' },
  { id: 'sw003', name: 'MySQL', version: '8.0.36', type: '基础软件', owner: '李四', team: '支付组' },
  { id: 'sw004', name: 'Nginx', version: '1.25.3', type: '基础软件', owner: '王五', team: '支付二组' },
  { id: 'sw005', name: 'Kafka', version: '3.6.1', type: '基础软件', owner: '赵六', team: '支付三组' },
])

// 自研项目列表
const projectList = ref([
  { id: 'pj001', name: '订单系统', team: '交易组', owner: '张三' },
  { id: 'pj002', name: '支付系统', team: '支付组', owner: '钱七' },
  { id: 'pj003', name: '用户中心', team: '支付二组', owner: '李四' },
  { id: 'pj004', name: '消息平台', team: '支付三组', owner: '王五' },
])

// 搜索和过滤
const softwareSearchText = ref('')
const filteredSoftwareList = ref(softwareList.value)

function filterSoftwareList() {
  if (!softwareSearchText.value) {
    filteredSoftwareList.value = softwareList.value
    return
  }
  const text = softwareSearchText.value.toLowerCase()
  filteredSoftwareList.value = softwareList.value.filter(sw =>
    sw.name.toLowerCase().includes(text)
  )
}

const projectSearchText = ref('')
const filteredProjectList = ref(projectList.value)

function filterProjectList() {
  if (!projectSearchText.value) {
    filteredProjectList.value = projectList.value
    return
  }
  const text = projectSearchText.value.toLowerCase()
  filteredProjectList.value = projectList.value.filter(pj =>
    pj.name.toLowerCase().includes(text)
  )
}

// 编辑软件负责人
const editSoftwareOwnerVisible = ref(false)
const editSoftwareForm = ref({
  id: '',
  name: '',
  owner: '',
  newOwners: [] as string[],
})

function showEditSoftwareOwner(row: any) {
  editSoftwareForm.value = {
    id: row.id,
    name: row.name,
    owner: row.owner || '',
    newOwners: row.owner ? row.owner.split(',').map((s: string) => s.trim()).filter((s: string) => s) : [],
  }
  editSoftwareOwnerVisible.value = true
}

function confirmEditSoftwareOwner() {
  // 更新软件负责人（可以为空）
  const software = softwareList.value.find(sw => sw.id === editSoftwareForm.value.id)
  if (software) {
    software.owner = editSoftwareForm.value.newOwners.join(', ') || ''
    filterSoftwareList()
  }
  editSoftwareOwnerVisible.value = false
  ElMessage.success('负责人更新成功')
}

// 编辑项目负责人
const editProjectOwnerVisible = ref(false)
const editProjectForm = ref({
  id: '',
  name: '',
  owner: '',
  newOwners: [] as string[],
})

function showEditProjectOwner(row: any) {
  editProjectForm.value = {
    id: row.id,
    name: row.name,
    owner: row.owner || '',
    newOwners: row.owner ? row.owner.split(',').map((s: string) => s.trim()).filter((s: string) => s) : [],
  }
  editProjectOwnerVisible.value = true
}

function confirmEditProjectOwner() {
  // 更新项目负责人（可以为空）
  const project = projectList.value.find(pj => pj.id === editProjectForm.value.id)
  if (project) {
    project.owner = editProjectForm.value.newOwners.join(', ') || ''
    filterProjectList()
  }
  editProjectOwnerVisible.value = false
  ElMessage.success('负责人更新成功')
}
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
