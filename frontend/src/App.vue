<template>
  <router-view :key="layoutKey" />
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const layoutKey = computed(() => route.matched[0]?.path || route.path)

// 全局初始化（环境/组/部门/IDC）由 user store 的 actions 完成，数据存入 state 全局共享
// （对应参考项目 App.vue created 流程，逐个 dispatch initEnv/initGroupName/initGroupTree/initIdc）。
// 这些初始化失败不阻断登录（权限系统未对接时后端兜底返回/失败，前端静默告警）。
const userStore = useUserStore()

onMounted(() => {
  userStore.initEnv()
  userStore.initGroupName()
  userStore.initGroupTree()
  userStore.initIdc()
})
</script>
