import { useEffect, useCallback, useRef } from 'react'
import { useImmer } from 'use-immer'
import { getFriendById } from '@/modules/friend'
import { useUserStore } from '@/store/userStore'
import { useMessageContext } from './useMessageContext'
import type { ConversationItem } from '../types'

interface ConversationsState {
  conversations: ConversationItem[]
  loading: boolean
}

/**
 * 解析会话 ID，提取对方的用户 ID
 * 格式: single:1_2 (单聊), group:123 (群聊)
 */
function parseConversationId(
  convId: string,
  currentUserId: string | number
): {
  type: 'single' | 'group'
  targetId: string
} | null {
  const [type, ids] = convId.split(':')

  if (type === 'single') {
    const [id1, id2] = ids.split('_')
    const targetId = String(currentUserId) === id1 ? id2 : id1
    return { type: 'single', targetId }
  }

  if (type === 'group') {
    return { type: 'group', targetId: ids }
  }

  return null
}

/**
 * 基于 WebSocket 返回的 maxSeq 构建会话列表
 * 目前支持单聊，通过解析 conversation_id 获取好友信息
 */
export function useConversations() {
  const { user } = useUserStore()
  const { maxSeqs, conversationMessages } = useMessageContext()

  const [state, setState] = useImmer<ConversationsState>({
    conversations: [],
    loading: false,
  })

  // 用 ref 追踪已处理的会话 ID，避免重复请求
  const processedConvIdsRef = useRef<Set<string>>(new Set())
  const creatingConvPromisesRef = useRef<Map<string, Promise<ConversationItem | null>>>(new Map())

  // 当 maxSeqs 变化时，构建会话列表
  useEffect(() => {
    const convIds = Object.keys(maxSeqs)
    if (convIds.length === 0) return

    // 筛选出未处理的会话
    const newConvIds = convIds.filter((id) => !processedConvIdsRef.current.has(id))

    // 如果没有新会话，直接返回
    if (newConvIds.length === 0) {
      return
    }

    console.log('🔄 发现新会话:', newConvIds)

    setState((draft) => {
      draft.loading = true
    })

    // 标记这些会话为已处理（在请求前标记，避免重复请求）
    newConvIds.forEach((id) => processedConvIdsRef.current.add(id))

    // 只获取新增会话的信息
    Promise.all(
      newConvIds.map(async (convId) => {
        const parsed = parseConversationId(convId, user.id)
        if (!parsed) return null

        // 目前只处理单聊
        if (parsed.type === 'single') {
          try {
            const friend = await getFriendById(parsed.targetId)
            const messages = conversationMessages[convId] || []
            const lastMsg = messages[messages.length - 1]

            return {
              id: convId,
              name: friend.remark || friend.friendUser.nickname || friend.friendUser.username,
              title: friend.friendUser.nickname || friend.friendUser.username,
              avatar: friend.friendUser.nickname?.slice(0, 2).toUpperCase() || 'U',
              accent: generateAccentColor(convId),
              lastMessage: lastMsg?.content || '暂无消息',
              time: formatTime(lastMsg?.send_time || lastMsg?.create_time),
              unread: 0, // TODO: 后续接入未读数
              muted: false,
              pinned: friend.isPinned,
              online: false, // TODO: 后续接入在线状态
              typing: false,
              description: friend.friendUser.nickname || friend.friendUser.username,
              messages: [], // 消息由 ChatPage 处理
            } as ConversationItem
          } catch (error) {
            console.error(`Failed to get friend info for ${convId}:`, error)
            return null
          }
        }

        // TODO: 处理群聊
        return null
      })
    ).then((results) => {
      setState((draft) => {
        const newConversations = results.filter((c) => c !== null) as ConversationItem[]
        const existingIds = new Set(draft.conversations.map((conv) => conv.id))
        const uniqueConversations = newConversations.filter((conv) => !existingIds.has(conv.id))
        if (uniqueConversations.length > 0) {
          draft.conversations.push(...uniqueConversations)
        }
        draft.loading = false
      })
    })
  }, [maxSeqs, user.id, setState, conversationMessages])

  // 监听消息变化，更新会话的最后消息
  useEffect(() => {
    setState((draft) => {
      draft.conversations.forEach((conv) => {
        const messages = conversationMessages[conv.id] || []
        const lastMsg = messages[messages.length - 1]
        if (lastMsg) {
          conv.lastMessage = lastMsg.content
          conv.time = formatTime(lastMsg.send_time || lastMsg.create_time)
        }
      })
    })
  }, [conversationMessages, setState])

  const togglePin = useCallback(
    (convId: string) => {
      setState((draft) => {
        const conv = draft.conversations.find((c) => c.id === convId)
        if (conv) {
          conv.pinned = !conv.pinned
        }
      })
    },
    [setState]
  )

  const toggleMute = useCallback(
    (convId: string) => {
      setState((draft) => {
        const conv = draft.conversations.find((c) => c.id === convId)
        if (conv) {
          conv.muted = !conv.muted
        }
      })
    },
    [setState]
  )

  const markAsRead = useCallback(
    (convId: string) => {
      setState((draft) => {
        const conv = draft.conversations.find((c) => c.id === convId)
        if (conv) {
          conv.unread = 0
        }
      })
    },
    [setState]
  )

  // 获取或创建会话（用于点击联系人"发消息"时）
  const getOrCreateConversation = useCallback(
    async (convId: string) => {
      // 检查是否已存在
      const existing = state.conversations.find((c) => c.id === convId)
      if (existing) return existing
      const existingPromise = creatingConvPromisesRef.current.get(convId)
      if (existingPromise) return existingPromise

      const createPromise = (async () => {
        try {
          // 解析会话 ID，创建临时会话
          const parsed = parseConversationId(convId, user.id)
          if (!parsed || parsed.type !== 'single') {
            console.error('暂不支持非单聊会话')
            return null
          }

          try {
            const friend = await getFriendById(parsed.targetId)
            const tempConversation: ConversationItem = {
              id: convId,
              name: friend.remark || friend.friendUser.nickname || friend.friendUser.username,
              title: friend.friendUser.nickname || friend.friendUser.username,
              avatar: friend.friendUser.nickname?.slice(0, 2).toUpperCase() || 'U',
              accent: generateAccentColor(convId),
              lastMessage: '开始聊天吧',
              time: '',
              unread: 0,
              muted: false,
              pinned: friend.isPinned,
              online: false,
              typing: false,
              description: friend.friendUser.nickname || friend.friendUser.username,
              messages: [],
            }

            // 添加到会话列表
            setState((draft) => {
              const alreadyExists = draft.conversations.some((conv) => conv.id === convId)
              if (!alreadyExists) {
                draft.conversations.push(tempConversation)
              }
            })

            return tempConversation
          } catch (error) {
            console.error('创建临时会话失败:', error)
            return null
          }
        } finally {
          creatingConvPromisesRef.current.delete(convId)
        }
      })()

      creatingConvPromisesRef.current.set(convId, createPromise)
      return createPromise
    },
    [state.conversations, user.id, setState]
  )

  return {
    conversations: state.conversations,
    loading: state.loading,
    actions: {
      togglePin,
      toggleMute,
      markAsRead,
    },
    getOrCreateConversation,
  }
}

/**
 * 根据会话 ID 生成一个固定的主题色
 */
function generateAccentColor(convId: string): string {
  const colors = [
    '#3b82f6', // blue
    '#8b5cf6', // purple
    '#ec4899', // pink
    '#f59e0b', // amber
    '#10b981', // emerald
    '#06b6d4', // cyan
    '#6366f1', // indigo
    '#ef4444', // red
  ]

  let hash = 0
  for (let i = 0; i < convId.length; i++) {
    hash = convId.charCodeAt(i) + ((hash << 5) - hash)
  }

  return colors[Math.abs(hash) % colors.length]
}

/**
 * 格式化时间戳为简短显示
 */
function formatTime(timestamp?: number): string {
  if (!timestamp) return ''

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
    return '昨天'
  }

  // 今年
  if (date.getFullYear() === now.getFullYear()) {
    return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
  }

  // 跨年
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}
