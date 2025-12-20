import clsx from 'clsx'
import { ArrowLeft, Info, Search } from 'lucide-react'

import chatBg from '@/assets/chat/chat_bg.svg'
import startChat from '@/assets/chat/start_chat.svg'
import type { ConversationItem, MessageItem } from '../../types'
import { IconButton } from '../common'

interface ChatPanelProps {
  themeColor: string
  conversation: ConversationItem | null
  isMobile: boolean
  onBack: () => void
  onToggleDetails: () => void
  onOpenDetails: () => void
}

export function ChatPanel({
  themeColor,
  conversation,
  isMobile,
  onBack,
  onToggleDetails,
  onOpenDetails,
}: ChatPanelProps) {
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
          <div className="flex-1 space-y-4 overflow-y-auto px-4 py-6 md:px-8">
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
          </div>
          <div className="border-t border-slate-200 bg-white px-4 py-3 md:px-8">
            <div className="flex items-center justify-between rounded-2xl border border-dashed border-slate-300 bg-slate-50 px-4 py-3 text-xs text-slate-600">
              <span>收发消息功能稍后接入,当前仅展示对话布局</span>
              <span className="rounded-full bg-white px-3 py-1 text-[11px] font-semibold text-slate-700">
                敬请期待
              </span>
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
            'rounded-2xl px-4 py-2 text-sm shadow-sm ring-1 ring-black/5',
            isMine ? 'text-white' : 'text-slate-800'
          )}
          style={{
            background: isMine ? `linear-gradient(135deg, ${accent}, ${themeColor})` : 'white',
          }}
        >
          {message.content}
        </div>
        <div className="text-[11px] text-slate-400">
          {isMine ? '我 · ' : `${message.author} · `}
          {message.timestamp}
        </div>
      </div>
    </div>
  )
}
