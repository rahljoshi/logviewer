import { format, isSameDay } from 'date-fns'

export function formatBytes(value) {
  if (value < 1024) {
    return `${value} B`
  }

  const units = ['KB', 'MB', 'GB']
  let size = value / 1024
  let unitIndex = 0

  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex += 1
  }

  return `${size.toFixed(1)} ${units[unitIndex]}`
}

export function formatTimestamp(value) {
  const date = new Date(value)
  const now = new Date()

  if (isSameDay(date, now)) {
    return format(date, 'HH:mm:ss.SSS')
  }

  return format(date, 'MMM dd HH:mm:ss')
}
