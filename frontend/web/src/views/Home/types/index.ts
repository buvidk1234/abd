export type MessageDirection = 'in' | 'out'

export interface MessageItem {
  id: string
  author: string
  content: string
  timestamp: string
  direction: MessageDirection
  status?: number // 1: 发送中, 2: 已发送, 3: 失败
}

export interface ConversationItem {
  id: string
  name: string
  title?: string
  avatar: string
  accent: string
  lastMessage: string
  time: string
  unread: number
  muted?: boolean
  pinned?: boolean
  online?: boolean
  typing?: boolean
  draft?: string
  reminders?: string[]
  description?: string
  messages: MessageItem[]
}

export interface ChatMenuItem {
  id: string
  title: string
  description: string
  accent: string
}

export interface FriendItem {
  id: string
  name: string
  title: string
  avatar: string
  accent: string
  department: string
  email: string
  phone: string
  tags?: string[]
  location?: string
  status?: 'online' | 'offline' | 'busy'
  note?: string
}

export interface FriendRequestItem {
  id: string
  from: string
  note: string
  time: string
  status: 'pending' | 'accepted' | 'rejected'
}

export interface SavedGroupItem {
  id: string
  name: string
  members: number
  update: string
  accent: string
}

export interface BlacklistItem {
  id: string
  name: string
  reason?: string
  time: string
}
