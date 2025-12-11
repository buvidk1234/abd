import clsx from 'clsx'
import { ArrowLeft, Ban, CheckCircle, Clock, Users } from 'lucide-react'

import type { BlacklistItem, FriendRequestItem, SavedGroupItem } from '../types'
import { Badge, IconButton } from './common'

type SpecialType = 'new-friends' | 'saved-groups' | 'blacklist'

interface ContactSpecialPanelProps {
  type: SpecialType
  themeColor: string
  newFriends?: FriendRequestItem[]
  savedGroups?: SavedGroupItem[]
  blacklist?: BlacklistItem[]
  onBack: () => void
}

export function ContactSpecialPanel({
  type,
  themeColor,
  newFriends = [],
  savedGroups = [],
  blacklist = [],
  onBack,
}: ContactSpecialPanelProps) {
  return (
    <div className="flex min-w-0 flex-1 flex-col bg-slate-50">
      <div className="flex h-16 items-center justify-between border-b border-slate-200 bg-white px-4 shadow-sm md:px-6">
        <div className="flex items-center gap-3">
          <IconButton icon={<ArrowLeft className="size-5" />} onClick={onBack} ariaLabel="返回" />
          <div className="text-base font-semibold text-slate-900">{titleMap[type]}</div>
        </div>
        <Badge variant="primary">{badgeText(type, { newFriends, savedGroups, blacklist })}</Badge>
      </div>

      <div className="flex-1 overflow-y-auto px-4 py-6 md:px-8">
        {type === 'new-friends' && (
          <div className="space-y-3">
            {newFriends.map((item) => (
              <div
                key={item.id}
                className="flex items-center justify-between rounded-2xl border border-slate-100 bg-white px-4 py-3 shadow-sm"
              >
                <div className="flex items-center gap-3">
                  <div className="flex size-10 items-center justify-center rounded-xl bg-slate-900 text-xs font-semibold uppercase text-white shadow-sm">
                    {item.from.slice(0, 2)}
                  </div>
                  <div>
                    <div className="text-sm font-semibold text-slate-900">{item.from}</div>
                    <div className="text-xs text-slate-500">{item.note}</div>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-[11px] text-slate-400">{item.time}</span>
                  <StatusPill status={item.status} themeColor={themeColor} />
                </div>
              </div>
            ))}
          </div>
        )}

        {type === 'saved-groups' && (
          <div className="space-y-3">
            {savedGroups.map((item) => (
              <div
                key={item.id}
                className="flex items-center justify-between rounded-2xl border border-slate-100 bg-white px-4 py-3 shadow-sm"
              >
                <div className="flex items-center gap-3">
                  <div
                    className="flex size-10 items-center justify-center rounded-xl text-xs font-semibold uppercase text-white shadow-sm"
                    style={{ background: item.accent }}
                  >
                    {item.name.slice(0, 2)}
                  </div>
                  <div>
                    <div className="text-sm font-semibold text-slate-900">{item.name}</div>
                    <div className="text-xs text-slate-500">{item.members} 人 · 更新于 {item.update}</div>
                  </div>
                </div>
                <Users className="size-4 text-slate-400" />
              </div>
            ))}
          </div>
        )}

        {type === 'blacklist' && (
          <div className="space-y-3">
            {blacklist.map((item) => (
              <div
                key={item.id}
                className="flex items-center justify-between rounded-2xl border border-slate-100 bg-white px-4 py-3 shadow-sm"
              >
                <div>
                  <div className="text-sm font-semibold text-slate-900">{item.name}</div>
                  <div className="text-xs text-slate-500">{item.reason || '无备注原因'}</div>
                </div>
                <div className="flex items-center gap-2 text-xs text-slate-400">
                  <Clock className="size-4" />
                  {item.time}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

const titleMap: Record<SpecialType, string> = {
  'new-friends': '新的朋友',
  'saved-groups': '保存的群',
  blacklist: '黑名单',
}

function badgeText(
  type: SpecialType,
  data: { newFriends: FriendRequestItem[]; savedGroups: SavedGroupItem[]; blacklist: BlacklistItem[] }
) {
  if (type === 'new-friends') {
    return `${data.newFriends.filter((item) => item.status === 'pending').length} 个待处理`
  }
  if (type === 'saved-groups') {
    return `${data.savedGroups.length} 个群`
  }
  if (type === 'blacklist') {
    return `${data.blacklist.length} 个成员`
  }
  return ''
}

function StatusPill({ status, themeColor }: { status: FriendRequestItem['status']; themeColor: string }) {
  const icon =
    status === 'pending' ? <Clock className="size-3" /> : <CheckCircle className="size-3" style={{ color: themeColor }} />
  const text =
    status === 'pending' ? '待确认' : status === 'accepted' ? '已通过' : '已拒绝'
  return (
    <span
      className={clsx(
        'inline-flex items-center gap-1 rounded-full px-2 py-1 text-[11px] font-semibold',
        status === 'pending'
          ? 'bg-slate-100 text-slate-600'
          : status === 'accepted'
            ? 'bg-emerald-50 text-emerald-600'
            : 'bg-slate-100 text-slate-600'
      )}
    >
      {status === 'pending' ? <Clock className="size-3" /> : status === 'accepted' ? icon : <Ban className="size-3" />}
      {text}
    </span>
  )
}
