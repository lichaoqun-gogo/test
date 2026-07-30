<template>
  <div class="error-page">
    <div class="error-content">
      <div class="error-code">404</div>
      <el-icon :size="48" color="#f5222d" class="error-icon"><WarningFilled /></el-icon>
      <h2 class="error-title">页面不存在</h2>
      <p class="error-desc">抱歉，您访问的页面不存在或已被移除</p>
      <div class="error-actions">
        <el-button type="primary" @click="goBack">
          <el-icon><Back /></el-icon>
          返回首页
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { WarningFilled, Back } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

function goBack() {
  if (userStore.isAdmin) {
    router.push('/admin/asset/software-library')
  } else if (userStore.isSpecialist) {
    router.push('/specialist/my-ops')
  } else if (userStore.isLoggedIn) {
    router.push('/operator/my-ops')
  } else {
    router.push('/login')
  }
}
</script>

<style scoped>
.error-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #001529 0%, #002140 100%);
}

.error-content {
  text-align: center;
  color: #fff;
}

.error-code {
  font-size: 120px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.1);
  line-height: 1;
  margin-bottom: -20px;
}

.error-icon {
  margin-bottom: 16px;
}

.error-title {
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 12px;
  color: #fff;
}

.error-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.65);
  margin: 0 0 32px;
}

.error-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}
</style>
