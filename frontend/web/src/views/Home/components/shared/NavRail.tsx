import type { ReactNode } from 'react'
import clsx from 'clsx'
import { MessageCircle, Settings, Users } from 'lucide-react'

import { getInitials } from '../../utils'

interface NavRailProps {
  themeColor: string
  activeTab: 'chat' | 'friends'
  userName: string
  onSelectTab: (tab: 'chat' | 'friends') => void
  onOpenSettings: () => void
}

export function NavRail({
  themeColor,
  activeTab,
  userName,
  onSelectTab,
  onOpenSettings,
}: NavRailProps) {
  const initials = getInitials(userName)

  return (
    <div className="hidden h-full w-20 flex-shrink-0 flex-col items-center justify-between bg-white/90 py-4 shadow-sm backdrop-blur md:flex">
      <div className="flex flex-col items-center gap-3">
        <div
          className="relative flex size-12 items-center justify-center rounded-2xl text-sm font-bold uppercase text-white shadow-sm"
          style={{ background: themeColor }}
        >
          {initials}
          <span className="absolute -bottom-0.5 -right-0.5 flex size-3 items-center justify-center rounded-full border-2 border-white bg-emerald-500" />
        </div>
        <div className="text-[11px] font-semibold text-slate-600">{userName}</div>
      </div>

      <div className="flex flex-col items-center gap-3">
        <NavButton
          label="消息"
          active={activeTab === 'chat'}
          themeColor={themeColor}
          icon={<MessageCircle className="size-5" />}
          onClick={() => onSelectTab('chat')}
        />
        <NavButton
          label="通讯录"
          active={activeTab === 'friends'}
          themeColor={themeColor}
          icon={<Users className="size-5" />}
          onClick={() => onSelectTab('friends')}
        />
      </div>

      <div className="flex flex-col items-center gap-3">
        <NavButton
          label="设置"
          active={false}
          themeColor={themeColor}
          icon={<Settings className="size-5" />}
          onClick={onOpenSettings}
        />
      </div>
    </div>
  )
}

function NavButton({
  label,
  icon,
  active,
  onClick,
  themeColor,
}: {
  label: string
  icon: ReactNode
  active?: boolean
  themeColor: string
  onClick: () => void
}) {
  return (
    <button
      onClick={onClick}
      className={clsx(
        'flex size-12 flex-col items-center justify-center rounded-2xl text-[11px] font-semibold text-slate-500 transition hover:bg-slate-50 hover:text-slate-700',
        active ? 'bg-slate-900 text-white shadow-[0_10px_30px_rgba(0,0,0,0.15)]' : ''
      )}
      style={
        active ? { background: themeColor, boxShadow: `0 10px 30px ${themeColor}44` } : undefined
      }
    >
      <span className="flex items-center justify-center">{icon}</span>
      <span className="mt-1 leading-none">{label}</span>
    </button>
  )
}
