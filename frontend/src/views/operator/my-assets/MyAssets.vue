<template>
  <div class="my-assets">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="我负责的软件" name="software">
        <div class="filter-bar">
          <el-select v-model="filterGrade" placeholder="按等级筛选" clearable style="width: 150px">
            <el-option v-for="g in ['A', 'B', 'C', 'D']" :key="g" :label="`${g}级`" :value="g" />
          </el-select>
          <el-select v-model="filterType" placeholder="按类型筛选" clearable style="width: 150px; margin-left: 12px">
            <el-option label="中间件" value="middleware" />
            <el-option label="数据库" value="database" />
            <el-option label="工具" value="tool" />
            <el-option label="依赖包" value="dependency" />
            <el-option label="组件库" value="library" />
          </el-select>
        </div>

        <el-table :data="filteredSoftware" stripe style="margin-top: 16px">
          <el-table-column prop="name" label="软件名称" min-width="140">
            <template #default="{ row }">
              <el-link type="primary" @click="viewDetail(row.id)">{{ row.name }}</el-link>
            </template>
          </el-table-column>
          <el-table-column prop="version" label="版本" width="100" />
          <el-table-column prop="type" label="类型" width="100">
            <template #default="{ row }">
              {{ SOFTWARE_TYPE_LABELS[row.type as keyof typeof SOFTWARE_TYPE_LABELS] || row.type }}
            </template>
          </el-table-column>
          <el-table-column label="等级" width="120">
            <template #default="{ row }">
              <GradeTag :grade="row.grade" size="small" />
            </template>
          </el-table-column>
          <el-table-column prop="language" label="开发语言" width="120" />
          <el-table-column prop="bugCount" label="缺陷数" width="80" />
          <el-table-column prop="vulnerabilityCount" label="漏洞数" width="80">
            <template #default="{ row }">
              <el-tag v-if="row.vulnerabilityCount > 0" type="danger" size="small">
                {{ row.vulnerabilityCount }}
              </el-tag>
              <span v-else>0</span>
            </template>
          </el-table-column>
          <el-table-column prop="lastUpdated" label="最后更新" width="120" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
                {{ row.status === 'active' ? '使用中' : row.status === 'deprecated' ? '已停用' : '待审核' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="我的自研项目" name="project">
        <el-table :data="projects" stripe>
          <el-table-column prop="name" label="项目名称" min-width="160">
            <template #default="{ row }">
              <el-link type="primary" @click="viewProjectDetail(row.id)">{{ row.name }}</el-link>
            </template>
          </el-table-column>
          <el-table-column prop="department" label="所属部门" width="140" />
          <el-table-column prop="team" label="所属团队" width="120" />
          <el-table-column prop="owner" label="负责人" width="100" />
          <el-table-column label="边界" width="80">
            <template #default="{ row }">
              <el-tag :type="row.isBoundary ? 'danger' : 'info'" size="small">
                {{ row.isBoundary ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="licenseCount" label="许可证" width="90" />
          <el-table-column prop="dependencyCount" label="依赖数" width="90" />
          <el-table-column prop="vulnerabilityCount" label="漏洞数" width="90">
            <template #default="{ row }">
              <el-tag v-if="row.vulnerabilityCount > 0" type="danger" size="small">
                {{ row.vulnerabilityCount }}
              </el-tag>
              <span v-else>0</span>
            </template>
          </el-table-column>
          <el-table-column prop="lastScanned" label="最后扫描" width="120" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'active' ? 'success' : 'warning'" size="small">
                {{ row.status === 'active' ? '运行中' : '维护中' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { mockSoftwareAssets, mockProjects } from '@/mock/assets'
import { SOFTWARE_TYPE_LABELS } from '@/types'

const router = useRouter()
const activeTab = ref('software')
const filterGrade = ref('')
const filterType = ref('')

const filteredSoftware = computed(() => {
  let list = mockSoftwareAssets
  if (filterGrade.value) {
    list = list.filter(s => s.grade === filterGrade.value)
  }
  if (filterType.value) {
    list = list.filter(s => s.type === filterType.value)
  }
  return list
})

const projects = mockProjects

function viewDetail(id: string) {
  router.push(`/operator/my-assets/${id}`)
}

function viewProjectDetail(id: string) {
  router.push(`/operator/my-assets/${id}`)
}
</script>

<style scoped>
.filter-bar {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}
</style>
