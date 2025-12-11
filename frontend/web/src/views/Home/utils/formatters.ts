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
