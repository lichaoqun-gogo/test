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
import { casCallback, tokenManager } from '@/api/auth'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

onMounted(async () => {
  const token = (route.query.token as string) || ''

  if (!token) {
    ElMessage.error('登录失败：未收到认证令牌')
    router.replace('/login')
    return
  }

  // 存储 token
  tokenManager.setToken(token)

  // 用 token 获取用户信息
  const res = await casCallback(token)
  if (res.code === 200 && res.data) {
    userStore.setCasUser(token, res.data)
    // 直接加载菜单和注册路由，避免路由守卫重复请求 /auth/me
    await userStore.fetchMenus()
    // 根据角色跳转
    if (res.data.role === 'admin') {
      router.replace('/admin/asset/software-library')
    } else if (res.data.role === 'specialist') {
      router.replace('/specialist/my-ops')
    } else {
      router.replace('/operator/my-ops')
    }
  } else {
    ElMessage.error(res.message || '登录失败')
    tokenManager.removeToken()
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
