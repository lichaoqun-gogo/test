<template>
  <div class="main-layout">
    <!-- 顶部导航栏 -->
    <el-header class="layout-header">
      <div class="header-left">
        <div class="logo">
          <div class="logo-icon">
            <el-icon :size="20"><Management /></el-icon>
          </div>
          <span class="logo-text">开源治理平台V3</span>
        </div>

        <!-- 视图切换下拉框（admin角色或有canToggleAdminView权限的用户可见） -->
        <el-dropdown
          v-if="userStore.canToggleView"
          trigger="click"
          @command="handleViewModeChange"
          class="view-mode-dropdown"
        >
          <div class="view-mode-trigger">
            <el-icon :size="16"><component :is="userStore.viewMode === 'admin' ? 'Grid' : 'User'" /></el-icon>
            <span class="view-mode-text">{{ userStore.viewMode === 'admin' ? '全局治理' : '个人工作台' }}</span>
            <el-icon :size="12"><ArrowDown /></el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="admin">
                <el-icon><Grid /></el-icon>
                <span>全局治理</span>
              </el-dropdown-item>
              <el-dropdown-item command="personal">
                <el-icon><User /></el-icon>
                <span>个人工作台</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      <div class="header-right">
        <!-- 用户信息 -->
        <el-dropdown trigger="click" @command="handleCommand">
          <div class="user-info">
            <el-avatar :size="28" class="user-avatar">
              {{ userStore.userInfo?.nickname?.charAt(0) || 'U' }}
            </el-avatar>
            <span class="username">{{ userStore.userInfo?.nickname }}</span>
            <el-tag :type="roleType" size="small" effect="plain" class="role-tag">
              {{ roleLabel }}
            </el-tag>
            <el-icon :size="12" class="arrow-icon"><ArrowDown /></el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item class="user-dropdown-info" disabled>
                <div class="dropdown-user-detail">
                  <p class="dropdown-name">{{ userStore.userInfo?.nickname }}</p>
                  <p class="dropdown-dept">{{ userStore.userInfo?.department }}{{ userStore.userInfo?.team ? ' / ' + userStore.userInfo.team : '' }}</p>
                </div>
              </el-dropdown-item>
              <el-dropdown-item divided command="logout">
                <el-icon><SwitchButton /></el-icon>
                退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>

    <div class="layout-body">
      <!-- 左侧菜单 -->
      <el-aside class="layout-aside" :width="isCollapsed ? '64px' : '240px'">
        <el-menu
          :default-active="activeMenu"
          :collapse="isCollapsed"
          :collapse-transition="false"
          @select="handleMenuSelect"
          class="side-menu"
          background-color="#002140"
          text-color="rgba(255, 255, 255, 0.65)"
          active-text-color="#1890ff"
        >
          <!-- 动态菜单渲染 -->
          <template v-for="menu in displayMenus" :key="menu.id">
            <!-- 有子菜单：渲染为子菜单组 -->
            <el-sub-menu v-if="menu.children && menu.children.length" :index="menu.path || `group-${menu.id}`">
              <template #title>
                <el-icon><component :is="getMenuIcon(menu.icon)" /></el-icon>
                <span>{{ menu.name }}</span>
              </template>
              <!-- 递归渲染子菜单 -->
              <template v-for="child in menu.children" :key="child.id">
                <el-sub-menu
                  v-if="child.children && child.children.length"
                  :index="child.path || `sub-${child.id}`"
                >
                  <template #title>{{ child.name }}</template>
                  <el-menu-item
                    v-for="grandChild in child.children"
                    :key="grandChild.id"
                    :index="grandChild.path"
                  >
                    {{ grandChild.name }}
                  </el-menu-item>
                </el-sub-menu>
                <el-menu-item v-else :index="child.path">
                  {{ child.name }}
                </el-menu-item>
              </template>
            </el-sub-menu>
            <!-- 无子菜单：渲染为独立菜单项 -->
            <el-menu-item v-else :index="menu.path">
              <el-icon><component :is="getMenuIcon(menu.icon)" /></el-icon>
              <template #title>{{ menu.name }}</template>
            </el-menu-item>
          </template>
        </el-menu>
      </el-aside>

      <!-- 菜单收起/展开按钮（浮动在侧边栏右边缘中间） -->
      <div class="collapse-toggle" :style="{ left: isCollapsed ? '64px' : '240px' }" @click="isCollapsed = !isCollapsed">
        <el-icon :size="16">
          <ArrowLeft v-if="!isCollapsed" />
          <ArrowRight v-else />
        </el-icon>
      </div>

      <!-- 主内容区 -->
      <div class="layout-main" style="flex:1;background:#f0f2f5;padding:16px;overflow-y:auto;">
        <router-view />
      </div>
    </div>

    <!-- 底部 -->
    <el-footer class="layout-footer">
      <span>开源治理平台 V3.0.0 &copy; 2024</span>
    </el-footer>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { casLogin } from '@/api/auth'
import {
  Management, User, SwitchButton, ArrowDown,
  Coin, Warning, Setting,
  DataAnalysis, Monitor, List, Files,
  Grid, ArrowLeft, ArrowRight,
} from '@element-plus/icons-vue'
import type { MenuTreeNode } from '@/router'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const isCollapsed = ref(false)

const activeMenu = computed(() => route.path)

// 调试：监听路由变化，输出当前匹配的路由记录
watch(() => route.fullPath, (newPath, oldPath) => {
  console.log(`[MainLayout] route change: ${oldPath} -> ${newPath}`)
  console.log(`[MainLayout] matched:`, route.matched.map(m => ({ path: m.path, name: m.name })))
  console.log(`[MainLayout] all routes:`, router.getRoutes().map(r => ({ path: r.path, name: r.name })))
})

// 显示的菜单（admin 角色切换：全局治理用 admin 菜单，个人工作台用 specialist 菜单）
const displayMenus = computed<MenuTreeNode[]>(() => {
  if (userStore.role === 'admin' && userStore.viewMode === 'personal') {
    return userStore.specialistMenus || []
  }
  return userStore.menus || []
})

// 菜单图标映射
const iconMap: Record<string, any> = {
  Coin, Warning, Setting, DataAnalysis, Monitor, List, Files, Grid, User,
}

function getMenuIcon(icon: string) {
  return iconMap[icon] || Monitor
}

const roleLabel = computed(() => {
  const labels: Record<string, string> = {
    operator: '运维者',
    specialist: '治理专员',
    admin: '管理者',
  }
  return labels[userStore.role || ''] || ''
})

const roleType = computed(() => {
  const types: Record<string, 'success' | 'warning' | 'danger'> = {
    operator: 'success',
    specialist: 'warning',
    admin: 'danger',
  }
  return types[userStore.role || ''] || ''
})

// 切换视图模式
function handleViewModeChange(mode: 'personal' | 'admin') {
  userStore.toggleViewMode()
  if (mode === 'admin') {
    router.push('/admin/asset/software-library')
  } else {
    router.push('/specialist/my-task')
  }
  ElMessage.success(mode === 'admin' ? '已切换到全局治理' : '已切换到个人工作台')
}

// 手动处理菜单点击导航（替代 el-menu 的 router 模式，避免导航冲突）
function handleMenuSelect(index: string) {
  console.log(`[MainLayout] menu select: index=${index}, currentPath=${route.path}`)
  if (!index || index === route.path) {
    console.log(`[MainLayout] 跳过导航: index 为空或与当前路径相同`)
    return
  }
  // 检查目标路由是否存在
  const matched = router.getRoutes().find(r => r.path === index || r.path === `/${index}`)
  console.log(`[MainLayout] 目标路由匹配:`, matched ? { path: matched.path, name: matched.name } : '未找到')
  router.push(index).then(() => {
    console.log(`[MainLayout] push 成功, 新路径: ${route.fullPath}`)
  }).catch(err => {
    console.warn(`[MainLayout] 导航失败: ${index}`, err)
  })
}

async function handleCommand(command: string) {
  if (command === 'logout') {
    try {
      await userStore.logout()
    } catch (err) {
      console.error('[MainLayout] logout error:', err)
    }
    ElMessage.success('已退出登录')
    // 直接跳转 CAS 认证中心，不经过 /login 路由（避免中间页闪烁）
    casLogin()
  }
}
</script>

<style scoped>
/* 供应链安全洞察菜单项与治理效能看板二级菜单对齐 */
:deep(.el-menu--inline .el-menu-item.insight-menu-item) {
  padding-left: 40px !important;
}
.main-layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f0f2f5;
}

/* ========== 顶部导航栏 ========== */
.layout-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #001529;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  z-index: 100;
  padding: 0 20px;
  height: 56px !important;
}

.header-left {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border-radius: 8px;
  color: #fff;
}

.logo-text {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  letter-spacing: 1px;
}

/* 视图切换下拉框 */
.view-mode-dropdown {
  margin-left: 24px;
}

.view-mode-trigger {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  color: rgba(255, 255, 255, 0.75);
}

.view-mode-trigger:hover {
  color: rgba(255, 255, 255, 0.95);
  background: rgba(255, 255, 255, 0.08);
}

.view-mode-text {
  font-size: 14px;
  font-weight: 500;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* 用户信息 */
.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 6px;
  transition: background 0.2s;
}

.user-info:hover {
  background: rgba(255, 255, 255, 0.1);
}

.user-avatar {
  background: #1890ff;
  color: #fff;
  font-size: 14px;
}

.username {
  font-weight: 500;
  color: rgba(255, 255, 255, 0.85);
  font-size: 14px;
}

.role-tag {
  font-size: 12px;
}

.arrow-icon {
  color: rgba(255, 255, 255, 0.65);
}

.dropdown-user-detail {
  padding: 4px 0;
}

.dropdown-name {
  font-weight: 600;
  color: #303133;
  margin: 0 0 4px;
}

.dropdown-dept {
  font-size: 12px;
  color: #909399;
  margin: 0;
}

/* ========== 主体布局 ========== */
.layout-body {
  flex: 1;
  display: flex;
  overflow: hidden;
  position: relative;
}

/* ========== 左侧菜单 ========== */
.layout-aside {
  background: #002140;
  border-right: 1px solid rgba(255, 255, 255, 0.05);
  transition: width 0.3s;
  overflow-y: auto;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.15);
}

/* 菜单收起/展开浮动按钮 */
.collapse-toggle {
  position: absolute;
  left: 240px;
  top: 50%;
  transform: translate(-50%, -50%);
  z-index: 200;
  width: 24px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  background: #1890ff;
  color: #fff;
  border-radius: 0 6px 6px 0;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.2);
  transition: left 0.3s, background 0.2s;
}

.collapse-toggle:hover {
  background: #40a9ff;
}

.side-menu {
  border-right: none;
  height: 100%;
}

/* 主内容区 */
.layout-main {
  flex: 1;
  background: #f0f2f5;
  padding: 16px;
  overflow-y: auto;
}

/* 底部 */
.layout-footer {
  text-align: center;
  color: rgba(0, 0, 0, 0.45);
  font-size: 12px;
  height: 36px !important;
  line-height: 36px;
  background: #f0f2f5;
  border-top: 1px solid #e8e8e8;
}

/* ========== 深色菜单样式覆盖 ========== */
:deep(.el-menu) {
  border-right: none;
  background: #002140 !important;
}

:deep(.el-sub-menu__title),
:deep(.el-menu-item) {
  color: rgba(255, 255, 255, 0.65) !important;
}

:deep(.el-sub-menu__title:hover),
:deep(.el-menu-item:hover) {
  background: rgba(24, 144, 255, 0.15) !important;
  color: #fff !important;
}

:deep(.el-menu-item.is-active) {
  background: #1890ff !important;
  color: #fff !important;
  border-right: 3px solid #40a9ff;
}

:deep(.el-sub-menu.is-active > .el-sub-menu__title) {
  color: #fff !important;
}

:deep(.el-sub-menu.is-active > .el-sub-menu__title .el-icon) {
  color: #1890ff !important;
}

/* 子菜单样式 */
:deep(.el-menu--inline) {
  background: #001529 !important;
}

:deep(.el-menu--inline > .el-sub-menu > .el-sub-menu__title),
:deep(.el-menu--inline > .el-menu-item) {
  height: 48px;
  line-height: 48px;
}

:deep(.el-menu--inline .el-menu-item) {
  background: #001529 !important;
}

:deep(.el-menu--inline > .el-sub-menu > .el-menu--inline .el-menu-item) {
  padding-left: 48px !important;
}

:deep(.el-menu--inline .el-menu-item:hover) {
  background: rgba(24, 144, 255, 0.2) !important;
}

:deep(.el-menu--inline .el-menu-item.is-active) {
  background: #1890ff !important;
}

/* 滚动条 */
.layout-aside::-webkit-scrollbar {
  width: 4px;
}

.layout-aside::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
}

.layout-aside::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.15);
  border-radius: 2px;
}

.layout-aside::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.25);
}
</style>
