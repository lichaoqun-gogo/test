# 开源治理平台 V3 - 前端启动指南

## 环境要求

- **Node.js**: >= 18.0.0（推荐 20.x LTS）
- **npm**: >= 9.0.0

## 检查环境

```bash
node -v   # 应显示 v18.x.x 或更高版本
npm -v    # 应显示 9.x.x 或更高版本
```

## 如果 Node.js 未安装

### macOS 安装方式（使用 Homebrew）：

```bash
# 安装 Homebrew（如果未安装）
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 安装 Node.js
brew install node@20

# 验证安装
node -v
npm -v
```

### 或使用 nvm（推荐，可管理多版本）：

```bash
# 安装 nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

# 重新加载终端配置
source ~/.zshrc

# 安装 Node.js 20 LTS
nvm install 20
nvm use 20

# 验证
node -v
npm -v
```

## 启动项目

```bash
# 进入项目目录
cd /Users/lichaoqun/project/osm-mngt/frontend

# 安装依赖（首次运行或 package.json 变更后）
npm install

# 启动开发服务器
npm run dev
```

启动成功后，浏览器会自动打开 `http://localhost:3000`

## 演示账号

| 账号 | 密码 | 角色 | 权限 |
|------|------|------|------|
| operator | 123456 | 普通运维者 | 我的资产、软件服务、我的申请、我的工单 |
| specialist | 123456 | 开源治理专员 | 运维者权限 + 风险工单二次分配 |
| admin | 123456 | 管理者 | 四大中心完整权限 |

## 常见问题

### 1. `npm: command not found`
- 原因：未安装 Node.js 或环境变量未配置
- 解决：按上述步骤安装 Node.js

### 2. `ERR_OSSL_EVP_UNSUPPORTED`
- 原因：Node.js 版本过新与旧版 OpenSSL 不兼容
- 解决：设置环境变量 `export NODE_OPTIONS=--openssl-legacy-provider`

### 3. 端口 3000 被占用
- 解决：修改 `vite.config.ts` 中的 `server.port` 配置

### 4. TypeScript 类型错误
- 解决：执行 `npm run build` 查看详细错误信息

## 已修复问题

- ✅ 修复 `vite.config.ts` 中 ESM 模块 `__dirname` 不可用问题（改用 `import.meta.url`）
