import clsx from 'clsx'
import { ArrowLeft, Info, Search, Loader2, Send } from 'lucide-react'
import { useRef, useEffect, useState } from 'react'

import chatBg from '@/assets/chat/chat_bg.svg'
import startChat from '@/assets/chat/start_chat.svg'
import type { ConversationItem, MessageItem } from '../../types'
import { IconButton } from '../common'
import { useMessageContext } from '../../hooks/useMessageContext'
import { useUserStore } from '@/store/userStore'
import { toast } from 'sonner'

interface ChatPanelProps {
  themeColor: string
  conversation: ConversationItem | null
  isMobile: boolean
  onBack: () => void
  onToggleDetails: () => void
  onOpenDetails: () => void
  onLoadMore?: () => void
  isLoadingMore?: boolean
}

export function ChatPanel({
  themeColor,
  conversation,
  isMobile,
  onBack,
  onToggleDetails,
  onOpenDetails,
  onLoadMore,
  isLoadingMore,
}: ChatPanelProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const messagesContainerRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLTextAreaElement>(null)

  const [inputValue, setInputValue] = useState('')
  const [isSending, setIsSending] = useState(false)

  const { sendMessage } = useMessageContext()
  const { user } = useUserStore()

  // 自动调整输入框高度
  const adjustTextareaHeight = () => {
    const textarea = inputRef.current
    if (!textarea) return

    textarea.style.height = 'auto'
    textarea.style.height = `${Math.min(textarea.scrollHeight, 120)}px`
  }

  // 处理输入变化
  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInputValue(e.target.value)
    adjustTextareaHeight()
  }

  // 自动滚动到底部（仅在新消息时）
  useEffect(() => {
    if (conversation && conversation.messages.length > 0) {
      messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [conversation?.id, conversation?.messages.length])

  // 监听滚动，触发加载更多
  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const element = e.currentTarget
    // 当滚动到顶部 50px 内时，触发加载更多
    if (element.scrollTop < 50 && onLoadMore && !isLoadingMore) {
      onLoadMore()
    }
  }

  // 处理发送消息
  const handleSend = async () => {
    if (!inputValue.trim() || !conversation || isSending) return

    const content = inputValue.trim()
    setInputValue('')
    setIsSending(true)

    // 重置输入框高度
    if (inputRef.current) {
      inputRef.current.style.height = 'auto'
    }

    try {
      // 解析会话 ID，提取对方 ID 和会话类型
      const [convType, ids] = conversation.id.split(':')
      let targetId = ''
      let convTypeNum = 1 // 默认单聊

      if (convType === 'single') {
        const [id1, id2] = ids.split('_')
        targetId = String(user.id) === id1 ? id2 : id1
        convTypeNum = 1
      } else if (convType === 'group') {
        targetId = ids
        convTypeNum = 2
      }

      await sendMessage(
        {
          sender_id: String(user.id),
          conv_type: convTypeNum,
          target_id: targetId,
          msg_type: 1, // 文本消息
          content,
        },
        conversation.id // 传入会话 ID
      )

      console.log('✅ 消息已发送（乐观更新）')
    } catch (error) {
      console.error('发送失败:', error)
      toast.error('发送失败，请重试')
      setInputValue(content) // 恢复输入
    } finally {
      setIsSending(false)
      inputRef.current?.focus()
    }
  }

  // 处理回车发送
  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  if (!conversation) {
    return (
      <div className="flex min-w-0 flex-1 flex-col items-center justify-center bg-slate-50">
        <div className="text-center">
          <img src={startChat} alt="start chat" className="mx-auto mb-6 w-64 max-w-[70vw]" />
          <div className="text-lg font-semibold text-slate-800">选择一个会话开始沟通</div>
          <div className="mt-2 text-sm text-slate-500">左侧列表为你保留了最近的消息</div>
        </div>
      </div>
    )
  }

  return (
    <div className="flex min-w-0 flex-1 flex-col bg-slate-50">
      <ChatHeader
        conversation={conversation}
        isMobile={isMobile}
        onBack={onBack}
        onToggleDetails={onToggleDetails}
        onOpenDetails={onOpenDetails}
        themeColor={themeColor}
      />

      <div className="relative flex-1 overflow-hidden">
        <div
          className="absolute inset-0 opacity-60"
          style={{
            backgroundImage: `url(${chatBg})`,
            backgroundRepeat: 'repeat',
            backgroundSize: '360px',
          }}
        />
        <div className="relative flex h-full flex-col">
          <div
            ref={messagesContainerRef}
            className="flex-1 space-y-4 overflow-y-auto px-4 py-6 md:px-8"
            onScroll={handleScroll}
          >
            {isLoadingMore && (
              <div className="flex items-center justify-center py-2">
                <Loader2 className="size-5 animate-spin text-slate-400" />
                <span className="ml-2 text-sm text-slate-500">加载中...</span>
              </div>
            )}
            {conversation.messages.length === 0 ? (
              <div className="mt-10 flex h-full items-center justify-center text-sm text-slate-500">
                暂无历史消息
              </div>
            ) : (
              conversation.messages.map((message) => (
                <MessageBubble
                  key={message.id}
                  message={message}
                  accent={conversation.accent}
                  themeColor={themeColor}
                />
              ))
            )}
            <div ref={messagesEndRef} />
          </div>

          {/* 消息输入框 */}
          <div className="border-t border-slate-200 bg-white px-4 py-3 md:px-8">
            <div className="flex items-end gap-3">
              <textarea
                ref={inputRef}
                value={inputValue}
                onChange={handleInputChange}
                onKeyDown={handleKeyDown}
                placeholder="输入消息... (Enter 发送，Shift+Enter 换行)"
                className="flex-1 resize-none rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none transition-colors focus:border-slate-300 focus:bg-white"
                style={{
                  minHeight: '44px',
                  maxHeight: '120px',
                  overflow: 'auto',
                }}
                disabled={isSending}
              />
              <button
                onClick={handleSend}
                disabled={!inputValue.trim() || isSending}
                className={clsx(
                  'flex h-11 w-11 items-center justify-center rounded-2xl text-white transition-all',
                  inputValue.trim() && !isSending
                    ? 'cursor-pointer shadow-sm hover:shadow-md'
                    : 'cursor-not-allowed opacity-50'
                )}
                style={{
                  background:
                    inputValue.trim() && !isSending
                      ? `linear-gradient(145deg, ${conversation?.accent || themeColor}, ${themeColor})`
                      : '#cbd5e1',
                }}
              >
                {isSending ? (
                  <Loader2 className="size-5 animate-spin" />
                ) : (
                  <Send className="size-5" />
                )}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function ChatHeader({
  conversation,
  isMobile,
  onBack,
  onToggleDetails,
  onOpenDetails,
  themeColor,
}: {
  conversation: ConversationItem
  isMobile: boolean
  onBack: () => void
  onToggleDetails: () => void
  onOpenDetails: () => void
  themeColor: string
}) {
  return (
    <div className="flex h-16 items-center justify-between border-b border-slate-200 bg-white px-4 shadow-sm md:px-6">
      <div className="flex min-w-0 items-center gap-3">
        {isMobile && (
          <IconButton icon={<ArrowLeft className="size-5" />} onClick={onBack} ariaLabel="返回" />
        )}
        <div className="flex items-center gap-3">
          <div
            className="flex size-12 items-center justify-center rounded-2xl text-sm font-bold uppercase text-white shadow-sm"
            style={{ background: `linear-gradient(145deg, ${conversation.accent}, ${themeColor})` }}
            onClick={onOpenDetails}
            role="button"
          >
            {conversation.avatar}
          </div>
          <div className="min-w-0">
            <div className="flex items-center gap-2">
              <div className="truncate text-base font-semibold text-slate-900">
                {conversation.name}
              </div>
              {conversation.online ? (
                <span className="flex items-center gap-1 rounded-full bg-emerald-50 px-2 py-0.5 text-[11px] font-semibold text-emerald-600">
                  <span className="size-1.5 rounded-full bg-emerald-500" />
                  在线
                </span>
              ) : (
                <span className="rounded-full bg-slate-100 px-2 py-0.5 text-[11px] text-slate-500">
                  {conversation.description || '保持沟通'}
                </span>
              )}
            </div>
            <div className="text-xs text-slate-500">
              {conversation.description || '搭建中台 · 让协作像聊天一样顺滑'}
            </div>
          </div>
        </div>
      </div>
      <div className="flex items-center gap-2">
        <IconButton ariaLabel="搜索会话" icon={<Search className="size-4" />} />
        <IconButton
          ariaLabel="会话信息"
          icon={<Info className="size-4" />}
          onClick={onToggleDetails}
          active
          activeColor={themeColor}
        />
      </div>
    </div>
  )
}

function MessageBubble({
  message,
  accent,
  themeColor,
}: {
  message: MessageItem
  accent: string
  themeColor: string
}) {
  const isMine = message.direction === 'out'

  // 获取发送状态文本
  const getStatusText = () => {
    if (!isMine || !message.status) return ''
    if (message.status === 1) return '发送中...'
    if (message.status === 3) return '发送失败'
    return ''
  }

  const statusText = getStatusText()

  return (
    <div className={clsx('flex w-full items-end gap-2', isMine ? 'justify-end' : 'justify-start')}>
      {!isMine && (
        <div className="flex size-9 items-center justify-center rounded-xl bg-slate-200 text-xs font-semibold text-slate-700">
          {message.author.slice(0, 2).toUpperCase()}
        </div>
      )}
      <div className="max-w-[78%] space-y-1">
        <div
          className={clsx(
            'rounded-2xl px-4 py-2 text-sm shadow-sm ring-1 ring-black/5 transition-opacity',
            isMine ? 'text-white' : 'text-slate-800',
            message.status === 1 && 'opacity-70' // 发送中显示半透明
          )}
          style={{
            background: isMine ? `linear-gradient(135deg, ${accent}, ${themeColor})` : 'white',
          }}
        >
          {message.content}
        </div>
        <div className="flex items-center gap-1 text-[11px] text-slate-400">
          <span>
            {isMine ? '我 · ' : `${message.author} · `}
            {message.timestamp}
          </span>
          {statusText && (
            <span className={clsx(message.status === 3 && 'text-red-500')}>· {statusText}</span>
          )}
        </div>
      </div>
    </div>
  )
}
