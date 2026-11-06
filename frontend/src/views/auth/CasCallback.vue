<template>
  <div class="cas-callback">
    <div class="callback-box">
      <el-icon :size="48" class="loading-icon"><Loading /></el-icon>
      <p>正在登录，请稍候...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getDefaultHome } from '@/router'
import { tokenManager } from '@/api/auth'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

onMounted(async () => {
  // 新流程：后端已写 Cookie 并 302 直达角色首页，此页不再进入。
  // 兼容旧链接：带 ?token= 参数时（旧后端/手动拼接）处理登录，否则按 Cookie 恢复登录态。
  const token = (route.query.token as string) || ''

  if (token) {
    // 兼容旧链接：token 来自 URL 参数，写入 Cookie 后走登录流程
    tokenManager.setToken(token)
    const ok = await userStore.loginWithToken(token)
    if (ok) {
      router.replace(getDefaultHome(userStore.userInfo?.role))
    } else {
      ElMessage.error('登录失败')
      router.replace('/login')
    }
    return
  }

  // 无 token 参数：从 Cookie 恢复登录态并加载菜单（后端 302 直达首页时，守卫已处理；此处兜底）
  if (userStore.isLoggedIn) {
    await userStore.fetchMenus()
    router.replace(getDefaultHome(userStore.userInfo?.role))
  } else {
    router.replace('/login')
  }
})
</script>

<style scoped>
.cas-callback {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #001529 0%, #002140 40%, #003a70 100%);
}

.callback-box {
  text-align: center;
  color: #fff;
}

.loading-icon {
  animation: spin 1s linear infinite;
  color: #1890ff;
  margin-bottom: 16px;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.callback-box p {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.75);
}
</style>
