// Mock 数据统一导出
export * from './users'
export * from './assets'
export * from './tickets'
export * from './software'
export * from './admin'

// 模拟异步延迟
export function mockDelay(ms: number = 300): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms))
}
