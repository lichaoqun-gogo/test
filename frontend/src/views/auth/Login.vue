<template>
  <div class="login-redirecting">
    <p>正在跳转至 CAS 统一认证中心...</p>
  </div>
</template>

<script setup lang="ts">
import { casLogin } from '@/api/auth'

// 访问 /login 时直接跳转 CAS 认证（无 UI 闪烁）
casLogin()
</script>

<style scoped>
.login-redirecting {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #001529;
  color: rgba(255, 255, 255, 0.65);
  font-size: 14px;
}
</style>
