import dayjs from 'dayjs'

/**
 * 格式化时间为 'YYYY-MM-DD HH:mm:ss'
 * @param time 时间字符串、Date 对象或时间戳
 * @returns 格式化后的时间字符串，空值返回 '-'
 */
export function formatTime(time: string | Date | number | undefined | null): string {
  if (!time) return '-'
  const d = dayjs(time)
  if (!d.isValid()) return '-'
  return d.format('YYYY-MM-DD HH:mm:ss')
}
