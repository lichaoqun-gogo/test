# 开源治理平台（OSM）使用文档

## 一、项目概述

开源治理平台（OSM）是一套面向企业内部开源软件全生命周期管理的治理系统，支持资产台账、风险管理、治理任务审批、决策洞察等功能。系统采用前后端分离架构，支持 CAS 单点登录。

### 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Vite + Element Plus + Pinia |
| 后端 | Go + Gin + GORM |
| 数据库 | MySQL 8.0+ |
| 对象存储 | MinIO（源码备份，可选） |
| 认证 | CAS 2.0 协议单点登录 |

### 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| 前端 | 3000 | Vite 开发服务器（固定端口，不可飘移） |
| 后端 | 8080 | Gin HTTP 服务 |
| CAS | 8443 | Mock CAS 测试服务（仅开发环境） |
| MySQL | 3306 | 数据库 |
| MinIO | 9000 | 对象存储（可选） |

---

## 二、后端启动

### 2.1 环境要求

- Go 1.21+
- MySQL 8.0+
- MinIO（可选，不启动则文件上传功能不可用，不影响其他功能）

### 2.2 配置文件

配置文件位于 `backend/config/` 目录。`config.yaml` 是入口配置，其中的 `env` 字段指定当前环境，程序自动加载 `config.{env}.yaml`：

| config.yaml 中的 env | 实际加载文件 | 场景 |
|----------------------|-------------|------|
| `env: dev` | `config.dev.yaml` | 开发环境 |
| `env: prod` | `config.prod.yaml` | 生产环境 |
| `env: prod-bj` | `config.prod-bj.yaml` | 北京机房（示例：可按需扩展） |

切换环境只需修改 `config.yaml` 中的 `env` 值，多机房部署只需新建 `config.{环境名}.yaml`。

**注意：后端必须从 `backend/` 目录启动**，因为配置路径是相对路径 `config/`。

### 2.3 配置项说明

以 `config.dev.yaml` 为例：

```yaml
server:
  port: 8080          # 后端服务端口
  mode: debug          # Gin 模式：debug/release

database:
  host: 127.0.0.1      # MySQL 地址
  port: 3306           # MySQL 端口
  user: root           # 用户名
  password: ""         # 密码
  dbname: osm          # 数据库名

minio:
  endpoint: 127.0.0.1:9000    # MinIO 地址
  accessKeyId: minioadmin     # AccessKey
  secretAccessKey: minioadmin # SecretKey
  bucketName: osm             # 存储桶名称

log:
  level: debug         # 日志级别：debug/info
  format: json         # 日志格式
  filePath: logs/osm.log  # 日志文件路径（为空则只输出到控制台）

auth:
  mode: cas            # 认证模式：local=本地mock/cas=CAS单点登录
  cas:
    serverURL: http://localhost:8443           # CAS 服务器登录地址
    serviceURL: http://localhost:8080/api/v1/auth/cas/callback  # 后端回调地址
    validateURL: http://localhost:8443/serviceValidate          # CAS 票据验证地址
    frontendURL: http://localhost:3000                          # 前端首页地址
  jwt:
    secret: osm-jwt-secret-key-dev  # JWT 签名密钥（生产环境务必修改）
    expiry: 86400                   # JWT 过期时间（秒），86400=24小时
```

### 2.4 启动命令

```bash
cd backend

# 编译（首次或代码更新后执行）
go build -o bin/server cmd/server/main.go

# 启动（无需环境变量，自动读取 config.yaml 的 env 字段决定加载哪个配置）
./bin/server
```

切换环境只需修改 `config.yaml` 中的 `env` 值，无需改启动命令。

### 2.5 日志

日志同时输出到**控制台**和**文件**：

- 日志文件路径：`backend/logs/osm.log`（在 `config.*.yaml` 的 `log.filePath` 中配置）
- 日志目录不存在时会自动创建
- 日志以追加模式写入，不会覆盖

查看日志：
```bash
# 实时查看
tail -f backend/logs/osm.log

# 查看最近 100 行
tail -100 backend/logs/osm.log
```

---

## 三、前端启动

### 3.1 环境要求

- Node.js >= 18.0.0（推荐 20.x LTS）
- npm >= 9.0.0

### 3.2 配置文件

前端配置文件为 `frontend/vite.config.ts`，关键配置：

```typescript
server: {
  port: 3000,           // 前端端口（固定）
  strictPort: true,     // 严格端口模式：端口被占用时直接报错，不会飘到其他端口
  open: true,           // 启动后自动打开浏览器
  proxy: {
    '/api': {
      target: 'http://localhost:8080',  // 后端地址
      changeOrigin: true,
    },
  },
},
```

**注意事项：**
- `strictPort: true` 确保端口固定为 3000，被占用时会报错而非自动切换
- `proxy` 配置将所有 `/api` 开头的请求代理到后端 8080 端口
- 如需修改后端地址，修改 `proxy.target` 即可
- 如需修改前端端口，修改 `server.port` 即可

### 3.3 启动命令

```bash
cd frontend

# 安装依赖（首次启动）
npm install

# 启动开发服务器
npm run dev

# 构建生产包
npm run build

# 预览生产包
npm run preview
```

启动后访问：http://localhost:3000

---

## 四、CAS 单点登录配置

### 4.1 认证流程

```
用户访问前端 → 前端重定向到后端 /api/v1/auth/cas/login
→ 后端重定向到 CAS 登录页
→ 用户输入账号密码
→ CAS 验证成功，重定向到后端 /api/v1/auth/cas/callback?ticket=xxx
→ 后端验证 ticket，签发 JWT
→ 重定向到前端 /auth/cas-callback?token=xxx
→ 前端保存 token，加载菜单，跳转主页
```

### 4.2 切换为 CAS 模式

在配置文件中将 `auth.mode` 改为 `cas`：

```yaml
auth:
  mode: cas  # local 改为 cas
```

### 4.3 上线替换地址（关键）

上线时只需修改配置文件 `auth.cas.serverURL` 一个地址，其余项全部自动推导，无需手动填写：

```yaml
auth:
  mode: cas
  cas:
    serverURL: https://cas.company.com  # 替换为公司 CAS 服务器地址
```

其他配置项均**可选**，未配置时自动处理：

| 配置项 | 自动推导规则 | 说明 |
|--------|-------------|------|
| `serviceURL` | 从请求 Host 推导：`{scheme}://{host}/api/v1/auth/cas/callback` | 后端回调地址，Nginx 反代场景自动正确 |
| `validateURL` | 从 serverURL 推导：`{serverURL}/serviceValidate` | CAS 2.0 标准票据验证路径 |
| `frontendURL` | 为空时用相对路径 `/auth/cas-callback?token=xxx` 重定向 | 前后端同域/Nginx 反代时自动生效 |
| `jwt.secret` | 内置默认值 | JWT 签名密钥，无需配置 |
| `jwt.expiry` | 默认 86400（24 小时） | Token 过期时间 |

> **特殊情况**：如果前端和后端不同域（未通过 Nginx 反代到同域），需要额外配置 `frontendURL` 指向前端地址。

### 4.4 CAS 用户属性映射

系统从 CAS 返回的 XML 中读取以下属性：

| CAS 属性 | 系统字段 | 说明 |
|----------|----------|------|
| user | username | 用户名（唯一标识） |
| displayName | nickname | 昵称 |
| mail | email | 邮箱 |
| phone | phone | 手机号 |
| department | department | 部门 |

用户首次登录时自动创建（默认 operator 角色），后续登录自动更新昵称、邮箱、手机、部门信息。

### 4.5 Mock CAS 测试服务（仅开发环境）

项目提供了 Mock CAS 服务用于本地开发测试：

```bash
cd /path/to/mock认证/cas-test
./start.sh
```

测试账号：

| 用户名 | 密码 | 角色 | 说明 |
|--------|------|------|------|
| zhangsan | 123456 | operator | 普通运维者 |
| lisi | 123456 | specialist | 开源治理专员 |
| admin | admin123 | admin | 系统管理员（可切换全局治理/个人工作台） |

---

## 五、数据库初始化

### 5.1 建库

```sql
CREATE DATABASE IF NOT EXISTS osm DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
```

### 5.2 建表

执行 `backend/sql/schema.sql`，包含全部 30 张业务表的 DDL：

```bash
mysql -u root osm < backend/sql/schema.sql
```

### 5.3 初始化数据

执行 `backend/sql/init_data.sql`，包含菜单、角色菜单关联、管理员用户、部门专员映射、审批流模板：

```bash
mysql -u root osm < backend/sql/init_data.sql
```

### 5.4 角色说明

| 角色 | 标识 | 菜单范围 | 说明 |
|------|------|----------|------|
| 管理者 | admin | 全局治理菜单（台账、风险、治理、决策） | 可切换"全局治理/个人工作台"视图 |
| 专员 | specialist | 专员工作台菜单（我的项目、任务、运维） | 部门开源治理专员 |
| 运维者 | operator | 运维者菜单（我的项目、任务、运维） | 普通用户 |

### 5.5 表清单

数据库共 30 张表，按功能分组：

| 分组 | 表名 | 说明 |
|------|------|------|
| 用户权限 | users, menus, role_menus, dept_specialist_mapping | 用户、菜单、角色菜单、部门专员 |
| 资产台账 | project, project_inner, component, component_snapshot | 开源软件、自研项目、依赖包 |
| 名单风险 | component_core, component_black | 名单配置、阻断记录 |
| 许可证 | license, licensegroup, licensegroup_license | 许可证、分组 |
| 缺陷 | defect_ledger, defect_task_mapping, user_risk_report | 缺陷台账、映射、用户上报 |
| 治理任务 | governance_task, task_introduce, task_upgrade, task_withdraw, task_usage, task_defect, task_applicability_assessment | 任务主表及各类型扩展表 |
| 审批流 | approval_template, approval_template_node, approval_node, approval_history | 模板定义、节点实例、审批历史 |
| 其他 | pre_eval_cache, attachment, specialist_resource_mapping | 预评估缓存、附件、人员资源映射 |

---

## 六、从源码部署（生产环境）

### 6.1 环境准备

部署服务器需安装以下软件：

| 软件 | 版本要求 | 用途 |
|------|----------|------|
| Go | 1.21+ | 编译后端（仅源码部署方式需要） |
| Node.js | 18.0+（推荐 20.x LTS） | 编译前端 |
| npm | 9.0+ | 前端依赖管理 |
| MySQL | 8.0+ | 数据库 |
| Nginx | 任意稳定版本 | 反向代理 + 静态文件托管 |

> **提示**：Go 是编译型语言，`go build` 生成独立二进制文件，目标机不需要装 Go。支持两种部署方式，按需选择。

### 6.2 获取源码

```bash
# 下载源码到部署目录
git clone <项目仓库地址> /opt/osm-mngt
cd /opt/osm-mngt
```

### 6.3 数据库初始化

```bash
# 1. 创建数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS osm DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;"

# 2. 建表（30 张业务表）
mysql -u root -p osm < backend/sql/schema.sql

# 3. 初始化基础数据（菜单、用户、专员、审批模板）
mysql -u root -p osm < backend/sql/init_data.sql
```

### 6.4 编译后端

Go 编译后生成独立二进制文件，目标机不需要安装 Go。有两种部署方式：

#### 方式A：二进制部署（推荐生产）

在编译机交叉编译 Linux 版本，只拷贝二进制 + 配置到目标机：

```bash
cd /opt/osm-mngt/backend

# 修改生产配置
vi config/config.prod.yaml
# 至少修改以下项：
#   database.host / user / password  -> 实际 MySQL 信息
#   auth.cas.serverURL               -> 公司 CAS 服务器地址

# 交叉编译 Linux 版本（在 macOS 上编译 Linux 可执行文件）
GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/main.go

# 拷贝到目标服务器的文件清单：
# 1. bin/server-linux   -> 目标机 /opt/osm-mngt/backend/bin/server
# 2. config/            -> 目标机 /opt/osm-mngt/backend/config/
# 3. sql/               -> 目标机 /opt/osm-mngt/backend/sql/（首次部署）
```

目标机不需要 Go 环境，直接 `./bin/server` 运行。

#### 方式B：源码部署（开发/测试环境）

拷贝源码到目标机，安装 Go 后编译运行：

```bash
cd /opt/osm-mngt/backend

# 修改生产配置
vi config/config.prod.yaml

# 编译
go build -o bin/server cmd/server/main.go
```

### 6.5 编译前端

```bash
cd /opt/osm-mngt/frontend

# 安装依赖
npm install

# 编译，产物输出到 dist/ 目录
npm run build
```

编译后 `frontend/dist/` 目录包含所有静态文件（HTML、JS、CSS），将其复制到 Nginx 托管目录：

```bash
# 将编译产物复制到 Nginx 目录
mkdir -p /usr/share/nginx/osm
cp -r dist/* /usr/share/nginx/osm/
```

> **说明**：前端编译后是纯静态文件，多台 Nginx 节点部署时只需把 `dist/` 目录复制到每台服务器的 Nginx 目录下即可，无需在每台服务器安装 Node.js。

### 6.6 启动后端服务

```bash
cd /opt/osm-mngt/backend

# 启动（确保 config.yaml 中 env: prod）
./bin/server

# 或用 nohup 后台运行
nohup ./bin/server > /dev/null 2>&1 &
```

日志输出到 `backend/logs/osm.log`（同时在控制台输出）。

### 6.7 配置 Nginx

编辑 Nginx 配置文件（如 `/etc/nginx/conf.d/osm.conf`）：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /usr/share/nginx/osm;
        index index.html;
        # Vue History 路由模式：刷新页面时回退到 index.html
        try_files $uri $uri/ /index.html;
    }

    # 后端 API 反向代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

重载 Nginx：

```bash
nginx -t          # 检查配置语法
nginx -s reload   # 重载配置
```

### 6.8 验证部署

浏览器访问 `http://your-domain.com`，应自动跳转到 CAS 登录页。

### 6.9 前后端代码是否需要调整？

**不需要。** 原因如下：

| 检查项 | 现状 | 说明 |
|--------|------|------|
| 前端 API 地址 | 相对路径 `/api/v1` | Nginx 统一代理，无需改代码 |
| 前端路由模式 | History 模式 | Nginx 配置 `try_files` 支持 |
| 后端 CORS | 已配置允许所有来源 | 同域时不触发 CORS，无影响 |
| 后端 CAS 回调地址 | 从请求 Host 自动推导 | Nginx 传递 `Host` 和 `X-Forwarded-Proto` 头即可 |
| 后端前端重定向 | `frontendURL` 为空时用相对路径 | 同域时自动生效 |

---

## 七、开发环境启动

开发环境前后端分别启动，通过 Vite 内置代理联调：

```bash
# 1. 启动后端（终端 1，确保 config.yaml 中 env: dev）
cd backend
go build -o bin/server cmd/server/main.go && ./bin/server

# 2. 启动前端（终端 2）
cd frontend
npm install  # 首次需要安装依赖
npm run dev

# 3. 启动 Mock CAS（终端 3，仅开发环境）
cd /path/to/mock认证/cas-test
./start.sh
```

开发环境访问：
- 前端：http://localhost:3000
- 后端 API：http://localhost:8080/api/v1/...
- Mock CAS：http://localhost:8443/login

---

## 八、常见问题

### Q: 端口被占用怎么办？

前端 `strictPort: true`，端口 3000 被占用时会直接报错。解决方式：
- 释放 3000 端口：`lsof -i :3000` 找到进程并 kill
- 或修改 `vite.config.ts` 中的 `server.port`

后端端口 8080 被占用时，修改 `config.*.yaml` 中的 `server.port`。

### Q: 后端启动报 "failed to read config file"？

后端必须从 `backend/` 目录启动，因为配置路径是相对路径。确保 `cd backend` 后再执行启动命令。

### Q: Nginx 部署后页面刷新 404？

前端使用 History 路由模式，Nginx 必须配置 `try_files $uri $uri/ /index.html;`，否则刷新页面会 404。

### Q: Nginx 部署后 CAS 登录回调失败？

确保 Nginx 在 `location /api/` 中传递了以下头：
```nginx
proxy_set_header Host $host;
proxy_set_header X-Forwarded-Proto $scheme;
```
后端依赖这些头自动推导回调地址。如果前端和后端不同域，需在配置文件中显式设置 `frontendURL`。

### Q: MinIO 连接失败？

MinIO 是可选服务，连接失败不影响系统启动，仅文件上传/源码备份功能不可用。如需启用，启动 MinIO 服务并修改 `config.*.yaml` 中的 `minio` 配置。

### Q: 如何切换认证模式？

修改 `config.*.yaml` 中 `auth.mode`：
- `local`：本地 mock 账号（开发用）
- `cas`：CAS 单点登录（生产用）
