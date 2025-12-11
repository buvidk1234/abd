import type { ChatMenuItem } from '../types'

export const CHAT_MENUS: ChatMenuItem[] = [
  {
    id: 'new-chat',
    title: '发起会话',
    description: '选择联系人快速开聊',
    accent: '#e46342',
  },
  {
    id: 'create-group',
    title: '创建群聊',
    description: '把讨论集中到同一个地方',
    accent: '#0ea5e9',
  },
  {
    id: 'scan',
    title: '扫一扫',
    description: '加入或分享群聊二维码',
    accent: '#22c55e',
  },
]
