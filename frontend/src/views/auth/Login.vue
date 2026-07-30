<template>
  <div class="login-page">
    <!-- 背景装饰 -->
    <div class="login-bg">
      <div class="bg-grid"></div>
      <div class="bg-particles"></div>
    </div>

    <div class="login-card">
      <div class="login-header">
        <div class="logo-icon">
          <el-icon :size="36"><Management /></el-icon>
        </div>
        <h1>开源治理平台 V3</h1>
        <p>Open Source Governance Platform</p>
      </div>

      <!-- CAS 统一认证登录 -->
      <div class="cas-login">
        <el-button type="primary" size="large" class="login-btn" @click="handleCasLogin">
          <el-icon><User /></el-icon>
          通过 CAS 统一认证登录
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { casLogin } from '@/api/auth'

function handleCasLogin() {
  casLogin()
}
</script>

<style scoped>
.login-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #001529 0%, #002140 40%, #003a70 100%);
  position: relative;
  overflow: hidden;
}

.login-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
}

.bg-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(24, 144, 255, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(24, 144, 255, 0.05) 1px, transparent 1px);
  background-size: 50px 50px;
}

.bg-particles {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 20% 30%, rgba(24, 144, 255, 0.15) 0%, transparent 50%),
    radial-gradient(circle at 80% 70%, rgba(24, 144, 255, 0.1) 0%, transparent 50%),
    radial-gradient(circle at 50% 50%, rgba(82, 196, 26, 0.05) 0%, transparent 40%);
}

.login-card {
  width: 440px;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 12px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4), 0 0 40px rgba(24, 144, 255, 0.1);
  position: relative;
  z-index: 1;
  backdrop-filter: blur(10px);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #001529, #1890ff);
  border-radius: 16px;
  color: #fff;
  margin-bottom: 16px;
  box-shadow: 0 8px 24px rgba(24, 144, 255, 0.3);
}

.login-header h1 {
  font-size: 26px;
  font-weight: 600;
  color: #001529;
  margin: 0 0 6px;
  letter-spacing: 1px;
}

.login-header p {
  font-size: 13px;
  color: #8c8c8c;
  letter-spacing: 0.5px;
}

.cas-login {
  margin-top: 28px;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  letter-spacing: 2px;
  background: linear-gradient(135deg, #001529, #1890ff);
  border: none;
  transition: all 0.3s;
}

.login-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(24, 144, 255, 0.4);
}
</style>
