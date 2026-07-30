# 角色访问路径说明

## 概述

本系统通过**路由路径前缀**区分不同角色的访问权限，不再维护用户权限表。后续将对接公司统一权限管理系统。

## 角色与访问路径

| 角色 | 路径前缀 | 默认首页 | 说明 |
|------|---------|---------|------|
| 运维者 (operator) | `/operator/*` | `/operator/my-ops` | 负责日常运维工作 |
| 开源治理专员 (specialist) | `/specialist/*` | `/specialist/my-ops` | 负责开源治理工作，比运维者多"人员管理"功能 |
| 管理员 (admin) | `/admin/*` | `/admin/asset/software-library` | 负责资产管理、风险管理、治理决策等 |

## 演示账号

登录页面提供三个演示账号（密码均为 `123456`）：

- **普通运维者**：用户名 `operator`
- **开源治理专员**：用户名 `specialist`
- **管理者**：用户名 `admin`

## 各角色可访问页面

### 运维者 (`/operator/*`)

| 页面 | 路径 | 说明 |
|------|------|------|
| 我的项目 | `/operator/my-projects` | 查看关联的项目 |
| 我的任务 | `/operator/my-task` | 查看和处理任务 |
| 我的运维 | `/operator/my-ops` | 运维台账（默认页） |
| 我的使用 | `/operator/my-usage` | 使用记录 |
| 可信仓库 | `/operator/trusted-repo` | 可信仓库列表 |
| 软件服务 | `/operator/software-service` | 软件服务 |
| 项目详情 | `/operator/project-detail/:id` | 项目详情（隐藏路由） |
| 任务详情 | `/operator/my-task/:id` | 任务详情（隐藏路由） |

### 开源治理专员 (`/specialist/*`)

| 页面 | 路径 | 说明 |
|------|------|------|
| 我的项目 | `/specialist/my-projects` | 查看关联的项目 |
| 我的任务 | `/specialist/my-task` | 查看和处理任务 |
| 我的运维 | `/specialist/my-ops` | 运维台账（默认页） |
| 人员管理 | `/specialist/personnel-mgmt` | **专员专属**：管理人员分配 |
| 我的使用 | `/specialist/my-usage` | 使用记录 |
| 可信仓库 | `/specialist/trusted-repo` | 可信仓库列表 |
| 软件服务 | `/specialist/software-service` | 软件服务 |
| 项目详情 | `/specialist/project-detail/:id` | 项目详情（隐藏路由） |
| 任务详情 | `/specialist/my-task/:id` | 任务详情（隐藏路由） |

### 管理员 (`/admin/*`)

#### 台账管理
| 页面 | 路径 |
|------|------|
| 开源基础软件 | `/admin/asset/software-library` |
| 开源依赖包 | `/admin/asset/dependency-library` |
| 开源工具 | `/admin/asset/tool-library` |
| 自研项目 | `/admin/asset/self-developed` |
| 源码备份 | `/admin/asset/source-backup` |
| 可信仓库 | `/admin/trusted-repo` |

#### 风险管理
| 页面 | 路径 |
|------|------|
| 黑白名单管理 | `/admin/risk/blacklist` |
| 许可证信息库 | `/admin/risk/license-library` |
| 缺陷台账 | `/admin/risk/defect-ledger` |
| 活跃度健康度监测 | `/admin/risk/activity-monitor` |
| SCA自动化解析入库 | `/admin/risk/sca-parser` |

#### 治理中心
| 页面 | 路径 |
|------|------|
| 任务列表 | `/admin/governance/task-list` |

#### 决策与洞察中心
| 页面 | 路径 |
|------|------|
| 违规拦截统计 | `/admin/decision/interception` |
| 缺陷修复效能 | `/admin/decision/repair-efficiency` |
| 资产健康度分布 | `/admin/decision/health-distribution` |
| 治理工作量统计 | `/admin/decision/workload` |
| 核心依赖 Top N | `/admin/decision/core-dependencies` |
| 版本碎片化 Top N | `/admin/decision/version-fragmentation` |
| 高风险项目 Top N | `/admin/decision/high-risk-projects` |

## 路由守卫

系统通过路由守卫实现角色验证：

```typescript
// 根据路径前缀判断期望的角色
function getRoleFromPath(path: string): string | null {
  if (path.startsWith('/operator')) return 'operator'
  if (path.startsWith('/specialist')) return 'specialist'
  if (path.startsWith('/admin')) return 'admin'
  return null
}

// 如果用户角色与路径期望的角色不匹配，重定向到 403
router.beforeEach((to, _from, next) => {
  const expectedRole = getRoleFromPath(to.path)
  if (expectedRole && userStore.role !== expectedRole) {
    next({ path: '/403', replace: true })
    return
  }
  next()
})
```

## 按钮显示控制

在组件中使用路径前缀判断角色，示例：

```vue
<template>
  <!-- 运维者和专员可见，管理员不可见 -->
  <el-button v-if="!isAdmin">运维操作按钮</el-button>

  <!-- 仅专员可见 -->
  <el-button v-if="isSpecialist">专员专属按钮</el-button>

  <!-- 仅管理员可见 -->
  <el-button v-if="isAdmin">管理操作按钮</el-button>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()

const isOperator = computed(() => route.path.startsWith('/operator'))
const isSpecialist = computed(() => route.path.startsWith('/specialist'))
const isAdmin = computed(() => route.path.startsWith('/admin'))
</script>
```

## 登录跳转逻辑

登录成功后，根据用户角色跳转到对应首页：

```typescript
if (userStore.isAdmin) {
  router.push('/admin/asset/software-library')
} else if (userStore.isSpecialist) {
  router.push('/specialist/my-ops')
} else {
  router.push('/operator/my-ops')
}
```

## 验证步骤

1. 启动前端项目
2. 使用 `operator` 账号登录，验证：
   - 自动跳转到 `/operator/my-ops`
   - 左侧菜单显示运维者菜单
   - 手动访问 `/admin/asset/software-library` 应跳转到 403 页面
3. 退出登录，使用 `specialist` 账号登录，验证：
   - 自动跳转到 `/specialist/my-ops`
   - 左侧菜单显示专员菜单（包含"人员管理"）
   - "我的运维"页面中显示"缺陷管理"标签页
4. 退出登录，使用 `admin` 账号登录，验证：
   - 自动跳转到 `/admin/asset/software-library`
   - 左侧菜单显示管理员菜单

## 注意事项

- 角色信息目前通过 Mock 数据模拟，后续对接统一权限系统后从外部获取
- 菜单结构硬编码在 `MainLayout.vue` 中，根据路径前缀动态显示
- 运维者和专员使用相同的 Vue 组件，只是路由路径不同
- 跨角色访问（如 operator 访问 /admin/*）会被路由守卫拦截并重定向到 403 页面
