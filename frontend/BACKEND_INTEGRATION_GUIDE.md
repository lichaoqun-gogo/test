# 前端对接后端指南

## 1. 安装axios

```bash
cd frontend
npm install axios
```

## 2. 配置Vite代理

编辑 `frontend/vite.config.ts`：

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
```

## 3. 删除Mock数据

删除或注释掉以下文件中的Mock数据：

- `frontend/src/mock/admin.ts`
- `frontend/src/mock/operator.ts`
- `frontend/src/mock/specialist.ts`

## 4. 更新页面使用真实API

### 示例：开源基础软件页面

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { assetApi } from '@/api'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await assetApi.getProjects({
      page: page.value,
      pageSize: pageSize.value,
    })
    tableData.value = data.list
    total.value = data.total
  } catch (error) {
    console.error('获取数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>
```

## 5. 用户信息存储

在登录页面保存用户信息到localStorage：

```typescript
// 登录成功后
localStorage.setItem('userInfo', JSON.stringify({
  userId: '123456',
  userName: '张三',
  role: 'admin',  // operator/specialist/admin
  dept: '技术研发部'
}))

// 根据角色跳转
if (role === 'operator') {
  router.push('/operator/my-ops')
} else if (role === 'specialist') {
  router.push('/specialist/my-ops')
} else if (role === 'admin') {
  router.push('/admin/asset/software-library')
}
```

## 6. 启动服务

### 启动后端

```bash
cd backend
make run
```

后端服务运行在：`http://localhost:8080`

### 启动前端

```bash
cd frontend
npm run dev
```

前端服务运行在：`http://localhost:5173`

## 7. 验证

1. 访问 `http://localhost:5173`
2. 登录（任意用户名密码）
3. 查看网络请求，确认请求发送到 `http://localhost:8080/api/v1/...`
4. 验证数据是否正常显示

## 8. 常见问题

### CORS错误

确保后端已配置CORS中间件（已包含在 `middleware/cors.go`）

### 401错误

检查请求头是否正确传递用户信息：
- X-User-ID
- X-User-Name
- X-User-Role
- X-User-Dept

### 数据库连接失败

检查 `config/config.yaml` 中的数据库配置是否正确

### MinIO连接失败

检查 `config/config.yaml` 中的MinIO配置，或暂时不启动MinIO（文件上传功能会失败，但不影响其他功能）

## 9. API调用示例

### 台账管理

```typescript
// 查询开源基础软件
const { data } = await assetApi.getProjects({
  page: 1,
  pageSize: 10,
  keyword: 'Redis',
  type: 'java'
})

// 创建自研项目
await assetApi.createInnerProject({
  pro_code: 'ORDER-SYS',
  module_code: 'order',
  branch: 'main',
  type: 'java',
  group_name: '订单组',
  dept_name: '技术研发部'
})
```

### 任务管理

```typescript
// 创建引入申请任务
await taskApi.createTask({
  type: 'introduce',
  title: '引入Redis 7.0',
  applicant: '张三',
  deadline: '2026-12-31',
  task_detail: {
    componentIds: [1, 2],
    introduceReason: '用于缓存'
  }
})

// 查询任务详情
const { data } = await taskApi.getTaskDetail(taskId)
// data包含: task, task_detail, approval_nodes, approval_history

// 审批通过
await approvalApi.approve(taskId, {
  comment: '同意引入'
})
```

### 文件上传

```typescript
const formData = new FormData()
formData.append('file', file)
formData.append('source', 'task_introduce')
formData.append('description', '评估报告')

const { data } = await fileApi.upload(formData)
console.log('文件ID:', data.id)
```
