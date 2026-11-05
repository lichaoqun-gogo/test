<template>
  <router-view :key="layoutKey" />
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getEnv, getGroups, getUserDepts, getIdcs } from '@/api/reference'

const route = useRoute()
const layoutKey = computed(() => route.matched[0]?.path || route.path)

// ==================== 全局初始化（对应参考项目 App.vue created 流程） ====================
// 统一认证通过、加载 App.vue 后，请求完成 组织 / 部门 / 环境 / IDC 的初始化：
//   GET /api/env           → 环境配置（initEnv）
//   GET /api/groups        → 组/部门列表（initGroupTree）
//   GET /api/user/depts    → 部门树/列表（initGroupName）
//   GET /api/idcs          → IDC 列表
// 数据存入 user store（envProfiles/groupList/deptTree/deptList/idcList）。
// 注意：这些初始化失败不阻断登录（权限系统未对接时后端兜底返回/失败，前端静默告警）。
const userStore = useUserStore()

onMounted(() => {
  initEnv()
  initGroupTree()
  initGroupName()
  initIdcs()
})

// 环境配置（App 启动时初始化）
async function initEnv() {
  const data = await getEnv()
  if (data) userStore.setEnvProfiles(data.profiles || [])
}

// 组/部门列表
async function initGroupTree() {
  const data = await getGroups()
  if (data) userStore.setGroupList(data)
}

// 部门树/列表（组织名称）
async function initGroupName() {
  const data = await getUserDepts()
  if (data) {
    userStore.setDeptTree(data.deptTree || [])
    userStore.setDeptList(data.deptList || [])
  }
}

// IDC 列表
async function initIdcs() {
  const data = await getIdcs()
  if (data) userStore.setIdcList(data)
}
</script>
