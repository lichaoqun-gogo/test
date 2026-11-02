# 开源治理平台前端项目

基于 Vue3 + TypeScript + Vite 构建的开源治理平台前端项目。

## 技术栈

- **框架**: Vue 3.4+
- **语言**: TypeScript 5.4+
- **构建工具**: Vite 5.2+
- **UI 组件库**: Element Plus 2.6+
- **路由**: Vue Router 4.3+
- **状态管理**: Pinia 2.1+
- **HTTP 客户端**: Axios 1.6+
- **图表库**: ECharts 5.5+ / Vue-ECharts 6.7+
- **日期处理**: Day.js 1.11+

## 项目目录结构

```
frontend/
├── public/                 # 静态资源
├── src/
│   ├── api/               # API 接口封装
│   │   └── index.ts       # Axios 实例及请求方法
│   ├── assets/            # 静态资源（图片、字体等）
│   ├── components/        # 公共组件
│   │   ├── business/      # 业务组件
│   │   │   ├── GradeTag.vue
│   │   │   ├── PreEvalProgress.vue
│   │   │   ├── RiskTag.vue
│   │   │   └── TicketStatusTag.vue
│   │   └── layout/        # 布局组件
│   │       └── MainLayout.vue
│   ├── layouts/           # 布局组件（预留）
│   ├── mock/              # Mock 数据
│   ├── router/            # 路由配置
│   │   └── index.ts       # 路由定义及守卫
│   ├── store/             # Pinia 状态管理
│   │   └── user.ts        # 用户状态
│   ├── styles/            # 全局样式
│   │   └── global.scss    # 全局样式变量及工具类
│   ├── types/             # TypeScript 类型定义
│   │   └── index.ts
│   ├── utils/             # 工具函数
│   │   └── index.ts       # 常用工具函数
│   ├── views/             # 页面组件
│   │   ├── admin/         # 管理者控制台
│   │   │   ├── asset-center/      # 资产与供应链中心
│   │   │   ├── decision-center/   # 决策与洞察中心
│   │   │   ├── governance-center/ # 治理与运营中心
│   │   │   └── risk-center/       # 风险与合规中心
│   │   ├── auth/          # 认证页面
│   │   │   ├── Login.vue
│   │   │   ├── Forbidden.vue
│   │   │   └── NotFound.vue
│   │   ├── operator/      # 运维者工作台
│   │   │   ├── my-applications/
│   │   │   ├── my-assets/
│   │   │   ├── my-tickets/
│   │   │   └── software-service/
│   │   └── specialist/    # 开源治理专员（预留）
│   ├── App.vue            # 根组件
│   └── main.ts            # 入口文件
├── index.html
├── package.json
├── tsconfig.json
├── tsconfig.node.json
└── vite.config.ts
```

## 快速开始

### 安装依赖

```bash
npm install
```

### 启动开发服务器

```bash
npm run dev
```

访问 http://localhost:3000

### 构建生产版本

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 配置说明

### Vite 配置

- 路径别名: `@` → `src/`
- 开发服务器端口: 3000
- API 代理: `/api` → `http://localhost:8080`
- Element Plus 自动导入

### 路由配置

- **公开路由**: 登录页、错误页面
- **运维者路由**: `/operator/*`
- **管理者路由**: `/admin/*`
- 路由守卫: 认证检查、权限检查

### 全局组件

已注册的业务组件（可在模板中直接使用）:
- `<GradeTag>` - 等级标签
- `<RiskTag>` - 风险标签
- `<TicketStatusTag>` - 工单状态标签
- `<PreEvalProgress>` - 预评估进度

## 开发规范

- 使用 TypeScript 严格模式
- 组件使用 `<script setup>` 语法
- 使用 Composition API
- 遵循 Vue 3 官方风格指南
