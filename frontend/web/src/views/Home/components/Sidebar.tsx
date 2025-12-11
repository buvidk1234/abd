import { forwardRef } from 'react'
import { Plus, Search } from 'lucide-react'
import clsx from 'clsx'

import type { ConversationItem } from '../types'
import { CHAT_MENUS } from '../constants/menus'
import { ConversationList } from './ConversationList'
import { IconButton } from './common'

interface SidebarProps {
  themeColor: string
  userName?: string
  conversations: ConversationItem[]
  selectedId: string | null
  onSelect: (id: string) => void
  showAddMenu: boolean
  onToggleAddMenu: () => void
  onMenuClick: (menuId: string) => void
  onOpenSearch: () => void
  showSidebar: boolean
  isMobile: boolean
}

export const Sidebar = forwardRef<HTMLDivElement, SidebarProps>(
  (
    {
      themeColor,
      userName,
      conversations,
      selectedId,
      onSelect,
      showAddMenu,
      onToggleAddMenu,
      onMenuClick,
      onOpenSearch,
      showSidebar,
      isMobile,
    },
    menuRef
  ) => {
    return (
      <div
        className={clsx(
          'relative z-10 flex h-full flex-col border-r border-slate-200 bg-white shadow-sm transition-all duration-300',
          showSidebar ? 'w-full md:w-[320px]' : 'hidden md:flex md:w-[320px]'
        )}
      >
        <div className="flex h-16 items-center justify-between px-5">
          <div className="flex items-center gap-3 leading-tight">
            <div
              className="flex size-11 items-center justify-center rounded-2xl bg-slate-900 text-xs font-bold uppercase text-white shadow-sm md:hidden"
              style={{ background: themeColor }}
            >
              {(userName || 'U').slice(0, 2).toUpperCase()}
            </div>
            <div>
              <div className="text-lg font-semibold text-slate-900">唐僧叨叨</div>
              <div className="text-xs text-slate-500">让沟通回归轻松</div>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <IconButton icon={<Search className="size-5" />} onClick={onOpenSearch} ariaLabel="搜索" />
            <div className="relative" ref={menuRef}>
              <IconButton icon={<Plus className="size-5" />} onClick={onToggleAddMenu} ariaLabel="添加" />
              {showAddMenu && (
                <div className="absolute right-0 top-full mt-2 w-56 rounded-2xl border border-slate-100 bg-white p-2 shadow-xl">
                  {CHAT_MENUS.map((menu) => (
                    <button
                      key={menu.id}
                      onClick={() => onMenuClick(menu.id)}
                      className="flex w-full items-start gap-2 rounded-xl px-3 py-2 text-left transition hover:bg-slate-50"
                    >
                      <div
                        className="mt-1 size-2 rounded-full"
                        style={{ backgroundColor: menu.accent }}
                      />
                      <div className="flex-1">
                        <div className="text-sm font-semibold text-slate-900">{menu.title}</div>
                        <div className="text-xs text-slate-500">{menu.description}</div>
                      </div>
                    </button>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="flex-1 overflow-y-auto">
          <ConversationList
            themeColor={themeColor}
            conversations={conversations}
            selectedId={selectedId}
            onSelect={onSelect}
            isMobile={isMobile}
          />
        </div>
      </div>
    )
  }
)

Sidebar.displayName = 'Sidebar'
