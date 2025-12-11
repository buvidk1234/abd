import { useMemo } from 'react'
import { X } from 'lucide-react'
import { useImmer } from 'use-immer'

import type { ContactItem, ConversationItem } from '../types'
import { filterContacts, filterConversations, formatUnreadCount } from '../utils'
import { Avatar, Badge, IconButton } from './common'

interface GlobalSearchProps {
  open: boolean
  conversations: ConversationItem[]
  contacts: ContactItem[]
  onClose: () => void
  onSelect: (payload: { type: 'conversation' | 'contact'; id: string }) => void
}

export function GlobalSearch({ open, conversations, contacts, onClose, onSelect }: GlobalSearchProps) {
  const [keyword, setKeyword] = useImmer<string>('')

  const filteredConversations = useMemo(
    () => filterConversations(conversations, keyword),
    [conversations, keyword]
  )

  const filteredContacts = useMemo(
    () => filterContacts(contacts, keyword),
    [contacts, keyword]
  )

  const hasResult = filteredConversations.length > 0 || filteredContacts.length > 0

  if (!open) return null

  return (
    <div className="fixed inset-0 z-30 flex items-start justify-center bg-black/30 p-4 backdrop-blur-sm">
      <div className="mt-12 w-full max-w-4xl rounded-3xl bg-white p-6 shadow-2xl">
        <div className="mb-4 flex items-center gap-3">
          <div className="flex flex-1 items-center rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
            <input
              autoFocus
              placeholder="搜索联系人、群聊或消息关键词..."
              className="w-full bg-transparent text-sm text-slate-700 outline-none"
              value={keyword}
              onChange={(e) =>
                setKeyword(() => {
                  return e.target.value
                })
              }
            />
          </div>
          <IconButton icon={<X className="size-5" />} onClick={onClose} ariaLabel="关闭" />
        </div>

        <div className="max-h-[60vh] space-y-4 overflow-y-auto">
          {!hasResult ? (
            <div className="flex h-32 items-center justify-center text-sm text-slate-500">
              没有找到匹配的内容
            </div>
          ) : (
            <>
              {filteredConversations.length > 0 && (
                <div className="space-y-2">
                  <div className="px-1 text-xs font-semibold uppercase text-slate-400">会话</div>
                  {filteredConversations.map((item) => (
                    <button
                      key={item.id}
                      onClick={() => onSelect({ type: 'conversation', id: item.id })}
                      className="flex w-full items-center gap-3 rounded-2xl border border-slate-100 px-4 py-3 text-left transition hover:border-slate-200 hover:bg-slate-50"
                    >
                      <Avatar
                        name={item.name}
                        avatar={item.avatar}
                        accent={item.accent}
                        size="md"
                      />
                      <div className="min-w-0 flex-1">
                        <div className="flex items-center gap-2">
                          <div className="truncate text-sm font-semibold text-slate-900">
                            {item.name}
                          </div>
                          <span className="text-[11px] text-slate-400">{item.time}</span>
                        </div>
                        <div className="truncate text-xs text-slate-500">{item.lastMessage}</div>
                      </div>
                      {item.unread > 0 && (
                        <Badge variant="primary">{formatUnreadCount(item.unread)}</Badge>
                      )}
                    </button>
                  ))}
                </div>
              )}

              {filteredContacts.length > 0 && (
                <div className="space-y-2">
                  <div className="px-1 text-xs font-semibold uppercase text-slate-400">联系人</div>
                  {filteredContacts.map((item) => (
                    <button
                      key={item.id}
                      onClick={() => onSelect({ type: 'contact', id: item.id })}
                      className="flex w-full items-center gap-3 rounded-2xl border border-slate-100 px-4 py-3 text-left transition hover:border-slate-200 hover:bg-slate-50"
                    >
                      <Avatar
                        name={item.name}
                        avatar={item.avatar}
                        accent={item.accent}
                        size="md"
                      />
                      <div className="min-w-0 flex-1">
                        <div className="flex items-center gap-2">
                          <div className="truncate text-sm font-semibold text-slate-900">
                            {item.name}
                          </div>
                          <span className="text-[11px] text-slate-400">{item.title}</span>
                        </div>
                        <div className="truncate text-xs text-slate-500">{item.department}</div>
                      </div>
                      {item.tags && item.tags.length > 0 && (
                        <Badge variant="default" size="md">
                          {item.tags.slice(0, 2).join('/')}
                        </Badge>
                      )}
                    </button>
                  ))}
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}
