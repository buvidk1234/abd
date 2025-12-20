import clsx from 'clsx'
import { BellOff, Pin, Sparkles } from 'lucide-react'

import type { ConversationItem } from '../../types'
import { formatUnreadCount } from '../../utils'
import { Avatar } from '../common'

interface ConversationListProps {
  themeColor: string
  conversations: ConversationItem[]
  selectedId: string | null
  onSelect: (id: string) => void
  isMobile: boolean
}

export function ConversationList({
  themeColor,
  conversations,
  selectedId,
  onSelect,
  isMobile,
}: ConversationListProps) {
  return (
    <div className="flex flex-col divide-y divide-slate-100">
      {conversations.map((conversation) => (
        <ConversationRow
          key={conversation.id}
          conversation={conversation}
          themeColor={themeColor}
          selected={conversation.id === selectedId}
          onSelect={onSelect}
          isMobile={isMobile}
        />
      ))}
    </div>
  )
}

interface ConversationRowProps {
  conversation: ConversationItem
  themeColor: string
  selected: boolean
  onSelect: (id: string) => void
  isMobile: boolean
}

function ConversationRow({
  conversation,
  themeColor,
  selected,
  onSelect,
  isMobile,
}: ConversationRowProps) {
  const messageWidthClass = isMobile ? 'max-w-[150px]' : 'max-w-[360px]'

  return (
    <button
      onClick={() => onSelect(conversation.id)}
      className={clsx(
        'w-full px-4 py-3 text-left transition hover:bg-slate-50',
        selected
          ? 'bg-[#e46342] text-white shadow-[0_12px_36px_rgba(228,99,66,0.25)] hover:bg-[#e46342]'
          : ''
      )}
    >
      <div className="flex items-center gap-3">
        <Avatar
          name={conversation.name}
          avatar={conversation.avatar}
          accent={conversation.accent}
          online={conversation.online}
          muted={conversation.muted}
          selected={selected}
          themeColor={themeColor}
          size="lg"
        />

        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <div className="flex min-w-0 items-center gap-2">
              <span className="truncate text-sm font-semibold">{conversation.name}</span>
              {conversation.title && (
                <span
                  className={clsx(
                    'truncate rounded-full px-2 py-0.5 text-[10px] font-semibold',
                    selected ? 'bg-white/20 text-white' : 'bg-slate-100 text-slate-600'
                  )}
                >
                  {conversation.title}
                </span>
              )}
              {conversation.pinned && (
                <Pin className={clsx('size-4', selected ? 'text-white/80' : 'text-slate-400')} />
              )}
              {conversation.muted && (
                <BellOff
                  className={clsx('size-4', selected ? 'text-white/80' : 'text-slate-400')}
                />
              )}
            </div>
            <span
              className={clsx('ml-auto text-[11px]', selected ? 'text-white/70' : 'text-slate-400')}
            >
              {conversation.time}
            </span>
          </div>
          <div className="mt-1 flex items-center gap-2">
            {conversation.draft && (
              <span className="rounded-sm bg-white/10 px-1 py-0.5 text-[11px] font-semibold text-white">
                [草稿]
              </span>
            )}
            {conversation.reminders?.map((reminder) => (
              <span
                key={reminder}
                className={clsx(
                  'truncate rounded-full px-2 py-0.5 text-[11px]',
                  selected ? 'bg-white/10 text-white' : 'bg-orange-50 text-orange-700'
                )}
              >
                {reminder}
              </span>
            ))}
            {conversation.typing ? (
              <TypingIndicator selected={selected} themeColor={themeColor} />
            ) : (
              <span
                className={clsx(
                  'truncate text-xs',
                  messageWidthClass,
                  selected ? 'text-white/80' : 'text-slate-500'
                )}
                title={conversation.lastMessage}
              >
                {conversation.lastMessage}
              </span>
            )}
          </div>
        </div>

        {conversation.unread > 0 && (
          <div
            className={clsx(
              'ml-2 min-w-[28px] rounded-full px-2 py-1 text-center text-[11px] font-semibold',
              conversation.muted ? 'bg-slate-200 text-slate-600' : 'bg-[#e46342] text-white'
            )}
          >
            {formatUnreadCount(conversation.unread)}
          </div>
        )}
      </div>
    </button>
  )
}

function TypingIndicator({ selected, themeColor }: { selected: boolean; themeColor: string }) {
  return (
    <span
      className={clsx(
        'flex items-center gap-1 rounded-full px-2 py-1 text-[11px]',
        selected ? 'bg-white/15 text-white' : 'bg-slate-100 text-slate-600'
      )}
    >
      <Sparkles
        className="size-3 animate-pulse"
        style={{ color: selected ? 'white' : themeColor }}
      />
      正在输入...
    </span>
  )
}
