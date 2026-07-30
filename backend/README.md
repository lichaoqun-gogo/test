# OSM 开源治理平台 - 后端服务

## 项目简介

开源治理平台 V3 后端服务，基于 Go + Gin + GORM + TiDB 构建。

## 技术栈

- **语言**: Go 1.21+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: TiDB 7.x (MySQL兼容)
- **对象存储**: MinIO
- **配置**: YAML

## 项目结构

```
backend/
├── cmd/server/main.go           # 程序入口
├── internal/
│   ├── config/config.go         # 配置管理
│   ├── database/database.go     # 数据库连接
│   ├── minio/minio.go           # MinIO客户端
│   ├── models/                  # GORM模型（21个表）
│   ├── handlers/                # HTTP处理器
│   │   ├── asset.go             # 台账管理API
│   │   ├── risk.go              # 风险管理API
│   │   ├── task.go              # 任务管理API
│   │   ├── approval.go          # 审批管理API
│   │   └── file.go              # 文件管理API
│   ├── routes/routes.go         # 路由注册
│   └── middleware/              # 中间件
│       ├── cors.go
│       └── logger.go
├── pkg/
│   ├── response/response.go     # 统一响应格式
│   └── external/                # 外部工具调用（预留）
├── config/config.yaml           # 配置文件
├── sql/init.sql                 # 数据库初始化脚本
├── go.mod
└── Makefile
```

## 快速开始

### 1. 环境要求

- Go 1.21+
- TiDB 7.x 或 MySQL 8.0+
- MinIO（可选，文件上传功能需要）

### 2. 创建数据库

```bash
# 连接TiDB/MySQL
mysql -h 127.0.0.1 -P 4000 -u root

# 创建数据库
CREATE DATABASE osm DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
```

### 3. 初始化数据库表

```bash
# 执行初始化脚本
mysql -h 127.0.0.1 -P 4000 -u root osm < sql/init.sql
```

### 4. 配置

编辑 `config/config.yaml`：

```yaml
server:
  port: 8080
  mode: debug

database:
  host: 127.0.0.1
  port: 4000
  user: root
  password: ""
  dbname: osm

minio:
  endpoint: 127.0.0.1:9000
  accessKeyId: minioadmin
  secretAccessKey: minioadmin
  useSSL: false
  bucketName: osm-mngt
```

### 5. 安装依赖

```bash
make install-deps
# 或
go mod tidy
```

### 6. 启动服务

```bash
# 开发模式
make run
# 或
go run cmd/server/main.go

# 编译
make build
```

服务启动后访问：`http://localhost:8080`

## API文档

### 台账管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/asset/projects | 查询开源基础软件列表 |
| GET | /api/v1/asset/projects/:id | 查询开源基础软件详情 |
| POST | /api/v1/asset/projects | 创建开源基础软件 |
| PUT | /api/v1/asset/projects/:id | 更新开源基础软件 |
| DELETE | /api/v1/asset/projects/:id | 删除开源基础软件 |
| GET | /api/v1/asset/inner-projects | 查询自研项目列表 |
| GET | /api/v1/asset/inner-projects/:id | 查询自研项目详情 |
| POST | /api/v1/asset/inner-projects | 创建自研项目 |
| PUT | /api/v1/asset/inner-projects/:id | 更新自研项目 |
| DELETE | /api/v1/asset/inner-projects/:id | 删除自研项目 |
| GET | /api/v1/asset/components | 查询开源依赖包列表 |
| GET | /api/v1/asset/components/:id | 查询开源依赖包详情 |

### 风险管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/risk/blacklist | 查询黑名单配置列表 |
| POST | /api/v1/risk/blacklist | 创建黑名单配置 |
| PUT | /api/v1/risk/blacklist/:id | 更新黑名单配置 |
| DELETE | /api/v1/risk/blacklist/:id | 删除黑名单配置 |
| GET | /api/v1/risk/blacklist/records | 查询黑名单阻断记录 |
| PUT | /api/v1/risk/blacklist/records/:id/process | 处理阻断记录 |
| GET | /api/v1/risk/licenses | 查询许可证列表 |
| GET | /api/v1/risk/license-groups | 查询许可证分组 |
| GET | /api/v1/risk/defects | 查询缺陷台账列表 |
| GET | /api/v1/risk/defects/:id | 查询缺陷详情 |
| POST | /api/v1/risk/defects | 创建缺陷 |
| PUT | /api/v1/risk/defects/:id/status | 更新缺陷状态 |

### 任务管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/tasks | 查询任务列表 |
| GET | /api/v1/tasks/:id | 查询任务详情 |
| GET | /api/v1/tasks/:id/detail | 查询任务完整详情（含扩展表和审批记录） |
| POST | /api/v1/tasks | 创建任务 |
| PUT | /api/v1/tasks/:id | 更新任务 |
| PUT | /api/v1/tasks/:id/status | 更新任务状态 |
| PUT | /api/v1/tasks/:id/withdraw | 撤回任务 |

### 审批管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/approvals/:taskId/approve | 审批通过 |
| POST | /api/v1/approvals/:taskId/reject | 审批驳回 |
| GET | /api/v1/approvals/:taskId/history | 查询审批历史 |

### 文件管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/files/upload | 上传文件 |
| GET | /api/v1/files/:id/download | 下载文件 |
| DELETE | /api/v1/files/:id | 删除文件 |

## 请求头

所有API请求需要携带用户信息（从统一权限系统获取）：

```
X-User-ID: 123456          # 用户ID
X-User-Name: 张三          # 用户姓名
X-User-Role: operator      # 用户角色（operator/specialist/admin）
X-User-Dept: 技术研发部    # 用户部门
```

## 响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {...}
}
```

### 分页响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "pageSize": 10
  }
}
```

### 错误响应

```json
{
  "code": 400,
  "message": "参数错误"
}
```

## 状态机

### 任务状态流转

```
draft（草稿） → pending（待审批） → processing（审批中） → completed（已完成）
                                                    ↓
                                               rejected（已驳回）
```

**状态转换规则**：
- draft → pending（提交）
- pending → processing/completed/rejected/draft（审批/撤回）
- processing → completed/rejected/draft（审批完成/驳回/撤回）
- rejected → draft（重新提交）

### 缺陷状态流转

```
待研判 → 待上会/已忽略
待上会 → 已分配/已忽略
已分配 → 已闭环
```

## 外部工具集成

以下功能通过调用外部独立工具服务实现（预留接口）：

- 源码备份工具
- SCA分析工具
- 预评估引擎
- 黑名单匹配引擎
- 活跃度分析引擎
- 漏洞引擎
- 通知服务

所有外部工具调用只打印日志，不写入数据库。

## 数据库表

共21张表：

### 台账管理
- project（开源基础软件）
- project_inner（自研项目）
- component_snapshot（开源依赖包）

### 风险管理
- component_core（黑名单配置）
- component_black（黑名单阻断记录）
- license（许可证）
- licensegroup（许可证分组）
- licensegroup_license（许可证分组关联）
- defect_ledger（缺陷台账）
- defect_task_mapping（缺陷关联任务）
- specialist_resource_mapping（人员-资源映射）

### 治理中心
- governance_task（任务主表）
- task_introduce（引入申请扩展表）
- task_upgrade（升级申请扩展表）
- task_withdraw（退出申请扩展表）
- task_defect（缺陷处置扩展表）
- approval_node（审批节点）
- approval_history（审批历史）
- pre_eval_cache（预评估缓存）
- user_risk_report（用户风险上报）

### 公共
- attachment（附件）

## 前端对接

前端API调用文件：`frontend/src/api/index.ts`

### 配置代理

在 `frontend/vite.config.ts` 中添加：

```typescript
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
```

### 使用示例

```typescript
import { assetApi, taskApi } from '@/api'

// 查询开源基础软件列表
const { data } = await assetApi.getProjects({ page: 1, pageSize: 10 })

// 创建任务
await taskApi.createTask({
  type: 'introduce',
  title: '引入Redis 7.0',
  applicant: '张三',
  deadline: '2026-12-31',
})
```

## 开发计划

- [ ] 可信仓库模块
- [ ] 活跃度检测模块
- [ ] 决策与洞察模块（统计查询）
- [ ] 外部工具对接
- [ ] 统一权限系统对接

## License

MIT
