import { useEffect, useCallback, useRef } from 'react'
import { useImmer } from 'use-immer'
import type { Message } from '@/modules'
import { useWS } from './useWS'

export interface GetMaxSeqResp {
  max_seqs: Record<string, number>
  min_seqs: Record<string, number>
}

export interface SeqRange {
  conversation_id: string
  begin: number
  end: number
  num: number
}

export interface PullMessageBySeqsReq {
  seq_ranges: SeqRange[]
  order: number
}

export interface PullMsgs {
  msgs: Message[]
  is_end: boolean
  end_seq: number
}

export interface PullMessageBySeqsResp {
  msgs: Record<string, PullMsgs>
  notification_msgs: Record<string, Message>
}

export interface SendMessageReq {
  sender_id: string
  conv_type: number
  target_id: string
  msg_type: number
  content: string
}

export interface MessageState {
  // 按会话ID存储消息列表
  conversationMessages: Record<string, Message[]>
  // 每个会话的最大 seq
  maxSeqs: Record<string, number>
  // 每个会话的最小 seq
  minSeqs: Record<string, number>
  // 加载状态
  loading: Record<string, boolean>
  // 是否已加载到底
  isEnd: Record<string, boolean>
}

export function useMessage() {
  const { send, subscribe, readyState } = useWS()

  const [state, setState] = useImmer<MessageState>({
    conversationMessages: {},
    maxSeqs: {},
    minSeqs: {},
    loading: {},
    isEnd: {},
  })

  // 用 ref 保存最新的 state，避免 useCallback 依赖 state
  const stateRef = useRef(state)
  stateRef.current = state

  const callbacksRef = useRef<{
    onNewMessage?: (msg: Message) => void
  }>({})

  // 1001: 获取最大 seq
  const getMaxSeq = useCallback(async () => {
    await send({
      req_identifier: 1001,
      data: null,
    })
  }, [send])

  // 1002: 拉取消息
  const pullMessages = useCallback(
    async (seqRanges: SeqRange[], order: number = 1) => {
      // 设置加载状态
      setState((draft) => {
        seqRanges.forEach((range) => {
          draft.loading[range.conversation_id] = true
        })
      })
      await send({
        req_identifier: 1002,
        data: {
          seq_ranges: seqRanges,
          order,
        } as PullMessageBySeqsReq,
      })
    },
    [send, setState]
  )

  // 1003: 发送消息（乐观更新）
  const sendMessage = useCallback(
    async (req: SendMessageReq, conversationId: string) => {
      // 生成临时消息 ID
      const clientMsgId = `client_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`

      // 立即在界面显示消息（乐观更新）
      const optimisticMsg: Message = {
        id: clientMsgId,
        conversation_id: conversationId,
        seq: 0, // 临时 seq
        sender_id: req.sender_id,
        client_msg_id: clientMsgId,
        msg_type: req.msg_type,
        content: req.content,
        status: 1, // 1: 发送中
        send_time: Date.now(),
        create_time: Date.now(),
        conv_type: req.conv_type,
        target_id: req.target_id,
      }

      setState((draft) => {
        if (!draft.conversationMessages[conversationId]) {
          draft.conversationMessages[conversationId] = []
        }
        draft.conversationMessages[conversationId].push(optimisticMsg)
      })

      // 发送到后端
      await send({
        req_identifier: 1003,
        data: {
          ...req,
          client_msg_id: clientMsgId,
        },
      })

      return clientMsgId
    },
    [send, setState]
  )

  // 1001: 获取最大 seq
  useEffect(() => {
    const unsubscribe = subscribe(1001, (msg) => {
      const data = msg.data as GetMaxSeqResp

      if (!data) return
      setState((draft) => {
        draft.maxSeqs = data.max_seqs || {}
        draft.minSeqs = data.min_seqs || {}
      })

      console.log('获取到最大 seq:', data)
    })

    return unsubscribe
  }, [subscribe, setState])

  // 1002: 拉取消息
  useEffect(() => {
    const unsubscribe = subscribe(1002, (msg) => {
      const data = msg.data as PullMessageBySeqsResp
      if (!data) return

      setState((draft) => {
        const msgs = data.msgs || {}

        Object.entries(msgs).forEach(([convId, pullMsgs]) => {
          if (!draft.conversationMessages[convId]) {
            draft.conversationMessages[convId] = []
          }
          const existingIds = new Set(draft.conversationMessages[convId].map((m) => m.id))
          const newMessages = pullMsgs.msgs.filter((m) => !existingIds.has(m.id))
          draft.conversationMessages[convId].push(...newMessages)
          draft.conversationMessages[convId].sort((a, b) => a.seq - b.seq)
          draft.isEnd[convId] = pullMsgs.is_end
          draft.loading[convId] = false
        })
      })
    })

    return unsubscribe
  }, [subscribe, setState])

  // 2001: 接收到消息
  useEffect(() => {
    const unsubscribe = subscribe(2001, (msg) => {
      const newMsg = msg.data as Message
      if (!newMsg) return

      setState((draft) => {
        const convId = newMsg.conversation_id

        if (!draft.conversationMessages[convId]) {
          draft.conversationMessages[convId] = []
        }

        // 如果有 client_msg_id，说明是自己发的消息，需要替换临时消息
        if (newMsg.client_msg_id) {
          const tempMsgIndex = draft.conversationMessages[convId].findIndex(
            (m) => m.client_msg_id === newMsg.client_msg_id
          )

          if (tempMsgIndex !== -1) {
            // 替换临时消息
            draft.conversationMessages[convId][tempMsgIndex] = {
              ...newMsg,
              status: 2, // 2: 已发送
            }
          } else {
            // 没找到临时消息，直接添加
            draft.conversationMessages[convId].push({
              ...newMsg,
              status: 2,
            })
          }
        } else {
          // 普通消息（别人发的），检查是否已存在
          const exists = draft.conversationMessages[convId].some((m) => m.id === newMsg.id)
          if (!exists) {
            draft.conversationMessages[convId].push(newMsg)
          }
        }

        // 按 seq 排序
        draft.conversationMessages[convId].sort((a, b) => a.seq - b.seq)

        // 更新 maxSeq
        if (!draft.maxSeqs[convId] || newMsg.seq > draft.maxSeqs[convId]) {
          draft.maxSeqs[convId] = newMsg.seq
        }
      })

      callbacksRef.current.onNewMessage?.(newMsg)
    })

    return unsubscribe
  }, [subscribe, setState])

  // 连接成功后获取 maxSeq
  useEffect(() => {
    if (readyState === WebSocket.OPEN) {
      getMaxSeq()
    }
  }, [readyState, getMaxSeq])

  // 辅助方法：根据会话ID获取消息
  const getMessagesByConversation = useCallback((conversationId: string): Message[] => {
    return stateRef.current.conversationMessages[conversationId] || []
  }, [])

  // 辅助方法：加载更多历史消息
  const loadMoreMessages = useCallback(
    async (conversationId: string, count: number = 20) => {
      // 从 ref 读取最新状态
      const currentState = stateRef.current
      const messages = currentState.conversationMessages[conversationId] || []

      if (currentState.isEnd[conversationId]) {
        console.log('已经加载到底了')
        return
      }

      // 如果没有消息，从 maxSeq 开始加载
      if (messages.length === 0) {
        const maxSeq = currentState.maxSeqs[conversationId]
        if (!maxSeq) {
          console.warn('没有 maxSeq，无法加载消息')
          return
        }
        await pullMessages(
          [
            {
              conversation_id: conversationId,
              begin: Math.max(1, maxSeq - count + 1),
              end: maxSeq,
              num: count,
            },
          ],
          1 // order: 1 表示从旧到新
        )
      } else {
        // 已有消息，继续向前加载
        const oldestSeq = messages[0]?.seq || 1
        await pullMessages(
          [
            {
              conversation_id: conversationId,
              begin: Math.max(1, oldestSeq - count),
              end: oldestSeq - 1,
              num: count,
            },
          ],
          1
        )
      }
    },
    [pullMessages]
  )

  // 设置新消息回调
  const setOnNewMessage = useCallback((callback: (msg: Message) => void) => {
    callbacksRef.current.onNewMessage = callback
  }, [])

  return {
    // 状态
    conversationMessages: state.conversationMessages,
    maxSeqs: state.maxSeqs,
    minSeqs: state.minSeqs,
    loading: state.loading,
    isEnd: state.isEnd,

    // 方法
    getMaxSeq,
    pullMessages,
    sendMessage,
    getMessagesByConversation,
    loadMoreMessages,
    setOnNewMessage,
  }
}
