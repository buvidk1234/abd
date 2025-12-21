import { createPortal } from 'react-dom'
import { useEffect, type ReactNode } from 'react'
import clsx from 'clsx'
import { LogOut, MessageCircle, Settings, Users } from 'lucide-react'
import { useImmer } from 'use-immer'

import { Button } from '@/components/ui/button'
import { getInitials } from '../../utils'

interface NavRailProps {
  themeColor: string
  activeTab: 'chat' | 'friends'
  userName: string
  onSelectTab: (tab: 'chat' | 'friends') => void
  onOpenSettings: () => void
  onLogout?: () => void
}

export function NavRail({
  themeColor,
  activeTab,
  userName,
  onSelectTab,
  onOpenSettings,
  onLogout,
}: NavRailProps) {
  const initials = getInitials(userName)
  const [showSettingsActions, setShowSettingsActions] = useImmer(false)
  const [mounted, setMounted] = useImmer(false)

  useEffect(() => {
    setMounted(true)
  }, [setMounted])

  const handleToggleSettings = () => setShowSettingsActions((open) => !open)
  const handleCloseMenu = () => setShowSettingsActions(false)
  const handleOpenSettings = () => {
    handleCloseMenu()
    onOpenSettings()
  }
  const handleLogout = () => {
    handleCloseMenu()
    onLogout?.()
  }

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

      <div className="relative flex flex-col items-center gap-3">
        <NavButton
          label="设置"
          active={showSettingsActions}
          themeColor={themeColor}
          icon={<Settings className="size-5" />}
          onClick={handleToggleSettings}
        />
        {mounted && showSettingsActions
          ? createPortal(
              <div
                className="fixed inset-0 z-50 flex items-end justify-start bg-transparent"
                onClick={handleCloseMenu}
              >
                <div
                  className="relative mb-20 ml-4 w-60 rounded-2xl border border-slate-200 bg-white/95 p-3 text-sm shadow-2xl backdrop-blur animate-in fade-in zoom-in-95"
                  onClick={(e) => e.stopPropagation()}
                >
                  <div className="mb-2 px-1 text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                    快捷操作
                  </div>
                  <div className="space-y-2">
                    <Button
                      variant="ghost"
                      size="sm"
                      className="w-full justify-start text-slate-700"
                      onClick={handleOpenSettings}
                    >
                      <Settings className="size-4" />
                      打开设置
                    </Button>
                    {onLogout ? (
                      <Button
                        variant="destructive"
                        size="sm"
                        className="w-full justify-start"
                        onClick={handleLogout}
                      >
                        <LogOut className="size-4" />
                        退出登录
                      </Button>
                    ) : null}
                  </div>
                </div>
              </div>,
              document.body
            )
          : null}
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
