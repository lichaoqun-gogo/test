import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Task } from '@/types'
import { taskApi } from '@/api'

console.log('[taskStore] taskApi:', taskApi)

export const useTaskStore = defineStore('task', () => {
  const tasks = ref<Task[]>([])

  // 从数据库加载任务列表
  async function loadTasks(params?: { status?: string; type?: string; keyword?: string }) {
    console.log('[taskStore.loadTasks] 调用API获取任务列表, params:', params)
    const res: any = await taskApi.getTasks(params || {})
    console.log('[taskStore.loadTasks] API返回结果:', res)
    // 后端返回分页格式：{list: [...], total, page, pageSize}
    tasks.value = Array.isArray(res?.list) ? res.list : []
    console.log('[taskStore.loadTasks] 任务列表已更新:', tasks.value)
  }

  // 创建任务（调用后端API）
  async function createTask(data: any) {
    const res = await taskApi.createTask(data)
    // 重新加载任务列表
    await loadTasks()
    return res
  }

  // 更新任务（调用后端API）
  async function updateTask(id: string, data: any) {
    await taskApi.updateTask(parseInt(id), data)
    // 重新加载任务列表
    await loadTasks()
  }

  // 获取任务详情
  async function getTaskDetail(id: string) {
    return await taskApi.getTaskDetail(parseInt(id))
  }

  // 更新任务状态
  async function updateTaskStatus(id: string, status: string) {
    await taskApi.updateTaskStatus(parseInt(id), { status })
    // 重新加载任务列表
    await loadTasks()
  }

  return {
    tasks,
    loadTasks,
    createTask,
    updateTask,
    updateTaskStatus,
    getTaskDetail,
  }
})
