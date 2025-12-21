import type { Message } from '@/modules'
import type { MessageItem } from '../types'

export function getInitials(name: string, length: number = 2): string {
  return name.slice(0, length).toUpperCase()
}

export function formatUnreadCount(count: number): string {
  return count > 99 ? '99+' : count.toString()
}

export function getStatusText(status?: 'online' | 'offline' | 'busy'): string {
  if (status === 'online') return '在线'
  if (status === 'busy') return '忙碌'
  if (status === 'offline') return '离线'
  return ''
}

/**
 * 格式化时间戳为可读字符串
 */
export function formatTimestamp(timestamp?: number): string {
  if (!timestamp) {
    return new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }

  const date = new Date(timestamp)
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const msgDate = new Date(date.getFullYear(), date.getMonth(), date.getDate())

  // 今天
  if (msgDate.getTime() === today.getTime()) {
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }

  // 昨天
  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)
  if (msgDate.getTime() === yesterday.getTime()) {
    return `昨天 ${date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })}`
  }

  // 今年
  if (date.getFullYear() === now.getFullYear()) {
    return date.toLocaleDateString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  // 跨年
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

/**
 * 将后端 Message 适配为前端 MessageItem
 */
export function adaptMessageToItem(msg: Message, currentUserId: string | number): MessageItem {
  const isMine = msg.sender_id === String(currentUserId)

  return {
    id: msg.id,
    author: isMine ? '我' : msg.sender_id,
    content: msg.content,
    timestamp: formatTimestamp(msg.send_time || msg.create_time),
    direction: isMine ? 'out' : 'in',
    status: msg.status, // 1: 发送中, 2: 已发送, 3: 失败
  }
}
