import { useMemo } from 'react'
import clsx from 'clsx'
import { Search } from 'lucide-react'
import { useImmer } from 'use-immer'

import type { ContactItem } from '../types'
import { SPECIAL_KEYS } from '../constants'
import { filterContacts, groupContactsByInitial } from '../utils'
import { Avatar, Badge } from './common'

interface ContactsSidebarProps {
  themeColor: string
  contacts: ContactItem[]
  selectedId: string | null
  onSelect: (id: string) => void
  onSelectSpecial: (key: string) => void
  specialCounts: {
    newFriends: number
    savedGroups: number
    blacklist: number
  }
  showSidebar: boolean
  isMobile: boolean
}

export function ContactsSidebar({
  themeColor,
  contacts,
  selectedId,
  onSelect,
  onSelectSpecial,
  specialCounts,
  showSidebar,
  isMobile,
}: ContactsSidebarProps) {
  const [keywordState, setKeywordState] = useImmer<{ keyword: string }>({ keyword: '' })

  const filteredContacts = useMemo(
    () => filterContacts(contacts, keywordState.keyword),
    [contacts, keywordState.keyword]
  )

  const grouped = useMemo(() => groupContactsByInitial(filteredContacts), [filteredContacts])

  return (
    <div
      className={clsx(
        'relative z-10 flex h-full w-full flex-shrink-0 flex-col border-r border-slate-200 bg-white shadow-sm transition-all duration-300 md:w-[320px]',
        showSidebar ? 'translate-x-0' : '-translate-x-full md:translate-x-0'
      )}
    >
      <div className="flex h-16 items-center justify-between px-5">
        <div className="leading-tight">
          <div className="text-lg font-semibold text-slate-900">通讯录</div>
          <div className="text-xs text-slate-500">常用联系人在这里</div>
        </div>
        <div className="hidden items-center gap-2 md:flex">
          <Badge variant="primary" size="md">{contacts.length} 人</Badge>
        </div>
      </div>

      <div className="px-4 pb-2">
        <div className="flex items-center gap-2 rounded-2xl border border-slate-200 bg-slate-50 px-3 py-2">
          <Search className="size-4 text-slate-400" />
          <input
            className="w-full bg-transparent text-sm text-slate-700 outline-none placeholder:text-slate-400"
            placeholder="搜索姓名、部门、标签"
            value={keywordState.keyword}
            onChange={(e) =>
              setKeywordState((draft) => {
                draft.keyword = e.target.value
              })
            }
          />
        </div>
      </div>

      <div className="space-y-3 px-4">
        <div className="grid grid-cols-2 gap-3">
          <QuickCard
            title="新的朋友"
            badge={specialCounts.newFriends > 0 ? `${specialCounts.newFriends}` : undefined}
            accent={themeColor}
            active={selectedId === SPECIAL_KEYS.newFriends}
            onClick={() => onSelectSpecial(SPECIAL_KEYS.newFriends)}
          />
          <QuickCard
            title="保存的群"
            badge={`${specialCounts.savedGroups}`}
            accent="#0ea5e9"
            active={selectedId === SPECIAL_KEYS.savedGroups}
            onClick={() => onSelectSpecial(SPECIAL_KEYS.savedGroups)}
          />
          <QuickCard
            title="黑名单"
            badge={`${specialCounts.blacklist}`}
            accent="#ef4444"
            active={selectedId === SPECIAL_KEYS.blacklist}
            onClick={() => onSelectSpecial(SPECIAL_KEYS.blacklist)}
          />
          <QuickCard
            title="常用联系人"
            badge={`${contacts.length}`}
            accent="#22c55e"
            active={!selectedId || !selectedId.startsWith('special:')}
            onClick={() => {
              if (contacts.length > 0) {
                onSelect(contacts[0].id)
              }
            }}
          />
        </div>
      </div>

      <div className="flex-1 overflow-y-auto px-2 pb-4">
        {grouped.length === 0 ? (
          <div className="mt-10 flex items-center justify-center text-sm text-slate-500">
            没有匹配的联系人
          </div>
        ) : (
          grouped.map(([key, list]) => (
            <div key={key} className="mb-4">
              <div className="px-4 pb-2 text-xs font-semibold uppercase text-slate-400">{key}</div>
              <div className="space-y-2">
                {list.map((contact) => (
                  <button
                    key={contact.id}
                    onClick={() => onSelect(contact.id)}
                    className={clsx(
                      'flex w-full items-center gap-3 rounded-2xl px-3 py-2 text-left transition',
                      selectedId === contact.id
                        ? 'bg-[#e46342] text-white shadow-[0_10px_30px_rgba(228,99,66,0.18)]'
                        : 'hover:bg-slate-50'
                    )}
                  >
                    <Avatar
                      name={contact.name}
                      avatar={contact.avatar}
                      accent={contact.accent}
                      status={contact.status}
                      selected={selectedId === contact.id}
                      themeColor={themeColor}
                      size="md"
                    />
                    <div className="min-w-0 flex-1">
                      <div className="flex items-center gap-2">
                        <div className="truncate text-sm font-semibold">{contact.name}</div>
                        <span
                          className={clsx(
                            'rounded-full px-2 py-0.5 text-[10px] font-semibold',
                            selectedId === contact.id ? 'bg-white/20 text-white' : 'bg-slate-100 text-slate-600'
                          )}
                        >
                          {contact.title}
                        </span>
                      </div>
                      <div
                        className={clsx(
                          'truncate text-xs',
                          selectedId === contact.id ? 'text-white/80' : 'text-slate-500'
                        )}
                      >
                        {contact.department}
                      </div>
                      {contact.tags && contact.tags.length > 0 && (
                        <div className="mt-1 flex flex-wrap gap-1">
                          {contact.tags.map((tag) => (
                            <span
                              key={tag}
                              className={clsx(
                                'rounded-full px-2 py-0.5 text-[10px]',
                                selectedId === contact.id ? 'bg-white/15 text-white' : 'bg-slate-100 text-slate-600'
                              )}
                            >
                              {tag}
                            </span>
                          ))}
                        </div>
                      )}
                    </div>
                    {!isMobile && (
                      <span
                        className={clsx(
                          'text-[11px] font-semibold',
                          selectedId === contact.id ? 'text-white/80' : 'text-slate-400'
                        )}
                      >
                        {contact.location || '在线协作'}
                      </span>
                    )}
                  </button>
                ))}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  )
}

function QuickCard({
  title,
  badge,
  accent,
  active,
  onClick,
}: {
  title: string
  badge?: string
  accent: string
  active?: boolean
  onClick: () => void
}) {
  return (
    <button
      onClick={onClick}
      className={clsx(
        'flex items-center justify-between rounded-2xl border px-4 py-3 text-left shadow-sm transition',
        active
          ? 'border-transparent text-white'
          : 'border-slate-200 bg-white text-slate-700 hover:border-slate-300 hover:bg-slate-50'
      )}
      style={
        active
          ? {
              background: `linear-gradient(135deg, ${accent}, ${accent}dd)`,
              boxShadow: `0 10px 24px ${accent}33`,
            }
          : undefined
      }
    >
      <div className="text-sm font-semibold">{title}</div>
      {badge && (
        <span
          className={clsx(
            'rounded-full px-2 py-1 text-[11px] font-semibold',
            active ? 'bg-white/20 text-white' : 'bg-slate-100 text-slate-600'
          )}
        >
          {badge}
        </span>
      )}
    </button>
  )
}
