import type { ReactNode } from 'react'
import clsx from 'clsx'
import { BellOff, CheckCircle, Pin, X } from 'lucide-react'

import type { ConversationItem } from '../types'
import { IconButton } from './common'

interface ChannelPanelProps {
  themeColor: string
  open: boolean
  conversation: ConversationItem | null
  onClose: () => void
  onTogglePin: () => void
  onToggleMute: () => void
  onMarkAsRead: () => void
}

export function ChannelPanel({
  themeColor,
  open,
  conversation,
  onClose,
  onTogglePin,
  onToggleMute,
  onMarkAsRead,
}: ChannelPanelProps) {
  if (!conversation) return null

  return (
    <div
      className={clsx(
        'pointer-events-none absolute inset-y-0 right-0 z-20 w-full max-w-[340px] translate-x-full transition-transform duration-300',
        open ? 'pointer-events-auto translate-x-0' : ''
      )}
    >
      <div className="flex h-full flex-col border-l border-slate-200 bg-white shadow-xl">
        <div className="flex items-center justify-between border-b border-slate-200 px-5 py-4">
          <div className="space-y-1">
            <div className="text-sm font-semibold text-slate-900">会话信息</div>
            <div className="text-xs text-slate-500">右侧栏随时收起</div>
          </div>
          <IconButton icon={<X className="size-4" />} onClick={onClose} ariaLabel="关闭" />
        </div>

        <div className="flex-1 space-y-6 overflow-y-auto px-5 py-6">
          <div className="flex items-center gap-3 rounded-2xl bg-slate-50 px-4 py-3">
            <div
              className="flex size-12 items-center justify-center rounded-2xl text-sm font-bold uppercase text-white shadow-sm"
              style={{ background: `linear-gradient(145deg, ${conversation.accent}, ${themeColor})` }}
            >
              {conversation.avatar}
            </div>
            <div className="min-w-0">
              <div className="flex items-center gap-2">
                <div className="truncate text-sm font-semibold text-slate-900">
                  {conversation.name}
                </div>
                {conversation.online && (
                  <span className="flex items-center gap-1 rounded-full bg-emerald-50 px-2 py-0.5 text-[10px] font-semibold text-emerald-600">
                    在线
                  </span>
                )}
              </div>
              <div className="truncate text-xs text-slate-500">
                {conversation.description || '整理常用信息,方便快速查看'}
              </div>
            </div>
          </div>

          <div className="space-y-3">
            <PanelAction
              icon={<Pin className="size-4" />}
              title={conversation.pinned ? '取消置顶' : '置顶会话'}
              description="固定在列表顶部"
              active={conversation.pinned}
              onClick={onTogglePin}
              activeColor={themeColor}
            />
            <PanelAction
              icon={<BellOff className="size-4" />}
              title={conversation.muted ? '取消免打扰' : '开启免打扰'}
              description="静音但保留未读数"
              active={conversation.muted}
              onClick={onToggleMute}
              activeColor={themeColor}
            />
            <PanelAction
              icon={<CheckCircle className="size-4" />}
              title="标记已读"
              description="清空未读提示"
              active={conversation.unread === 0}
              onClick={onMarkAsRead}
              activeColor={themeColor}
            />
          </div>

          <div className="rounded-2xl border border-slate-100 bg-slate-50 p-4">
            <div className="text-xs font-semibold text-slate-600">小提示</div>
            <p className="mt-2 text-xs leading-6 text-slate-500">
              右侧栏可以放置更多与当前会话相关的卡片信息、常用文件或快捷操作。这里先保留空间,方便后续扩展。
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}

function PanelAction({
  icon,
  title,
  description,
  active,
  activeColor = '#e46342',
  onClick,
}: {
  icon: ReactNode
  title: string
  description: string
  active?: boolean
  activeColor?: string
  onClick: () => void
}) {
  return (
    <button
      onClick={onClick}
      className={clsx(
        'flex w-full items-start gap-3 rounded-2xl border px-4 py-3 text-left transition',
        active
          ? 'border-transparent bg-gradient-to-r from-white to-orange-50 shadow-[0_10px_30px_rgba(228,99,66,0.12)]'
          : 'border-slate-200 hover:border-slate-300 hover:bg-slate-50'
      )}
    >
      <div
        className={clsx(
          'mt-0.5 flex size-9 items-center justify-center rounded-xl text-slate-600',
          active ? 'bg-white shadow-sm' : 'bg-slate-100'
        )}
        style={active ? { color: activeColor } : undefined}
      >
        {icon}
      </div>
      <div className="flex-1">
        <div className="text-sm font-semibold text-slate-900">{title}</div>
        <div className="text-xs text-slate-500">{description}</div>
      </div>
    </button>
  )
}
