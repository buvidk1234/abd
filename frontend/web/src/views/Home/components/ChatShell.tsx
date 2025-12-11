import { useRef } from 'react'
import { toast } from 'sonner'

import { useUserStore } from '@/store/userStore'
import { CHAT_MENUS } from '../constants/menus'
import { SPECIAL_KEYS, THEME_COLOR } from '../constants'
import { useClickOutside, useChatShellState, useResponsive } from '../hooks'
import { ChannelPanel } from './ChannelPanel'
import { ChatPanel } from './ChatPanel'
import { ContactDetail } from './ContactDetail'
import { ContactSpecialPanel } from './ContactSpecialPanel'
import { ContactsSidebar } from './ContactsSidebar'
import { GlobalSearch } from './GlobalSearch'
import { NavRail } from './NavRail'
import { Sidebar } from './Sidebar'

export { SPECIAL_KEYS }

export function ChatShell() {
  const userName = useUserStore((state) => state.user?.nickname || state.user?.username || '用户')
  const isMobile = useResponsive(1024)
  const addMenuRef = useRef<HTMLDivElement | null>(null)

  const { state, selectedConversation, selectedContact, actions } = useChatShellState(isMobile)

  useClickOutside(addMenuRef, actions.closeAddMenu, state.showAddMenu)

  const handleMenuClick = (menuId: string) => {
    const menu = CHAT_MENUS.find((item) => item.id === menuId)
    if (menu) {
      toast.info(`${menu.title} 功能正在对接中`)
    }
    actions.closeAddMenu()
  }

  const handleGlobalSearchSelect = (payload: { type: 'conversation' | 'contact'; id: string }) => {
    if (payload.type === 'conversation') {
      actions.selectConversation(payload.id)
    } else {
      actions.selectContact(payload.id)
    }
    actions.closeGlobalSearch()
  }

  return (
    <div className="flex h-[calc(100vh-1px)] w-full overflow-hidden bg-slate-100 text-slate-900">
      <NavRail
        themeColor={THEME_COLOR}
        activeTab={state.activeTab}
        userName={userName}
        onSelectTab={actions.setActiveTab}
        onOpenSettings={() => toast.info('设置功能稍后接入')}
      />

      {state.activeTab === 'chat' ? (
        <Sidebar
          ref={addMenuRef}
          themeColor={THEME_COLOR}
          userName={userName}
          conversations={state.conversations}
          selectedId={state.selectedConversationId}
          onSelect={actions.selectConversation}
          showAddMenu={state.showAddMenu}
          onToggleAddMenu={actions.toggleAddMenu}
          onMenuClick={handleMenuClick}
          onOpenSearch={actions.openGlobalSearch}
          showSidebar={state.showSidebarOnMobile}
          isMobile={state.isMobile}
        />
      ) : (
        <ContactsSidebar
          themeColor={THEME_COLOR}
          contacts={state.contacts}
          selectedId={state.selectedContactId}
          onSelect={actions.selectContact}
          onSelectSpecial={actions.selectContact}
          specialCounts={{
            newFriends: state.newFriends.filter((i) => i.status === 'pending').length,
            savedGroups: state.savedGroups.length,
            blacklist: state.blacklist.length,
          }}
          showSidebar={state.showSidebarOnMobile}
          isMobile={state.isMobile}
        />
      )}

      <div className="relative flex min-w-0 flex-1">
        {state.activeTab === 'chat' ? (
          <>
            <ChatPanel
              themeColor={THEME_COLOR}
              conversation={selectedConversation}
              isMobile={state.isMobile}
              onBack={actions.showSidebar}
              onToggleDetails={actions.toggleDetails}
              onOpenDetails={actions.openDetails}
            />
            <ChannelPanel
              themeColor={THEME_COLOR}
              open={state.showDetails && !!selectedConversation}
              conversation={selectedConversation}
              onClose={actions.closeDetails}
              onTogglePin={() => selectedConversation && actions.togglePin(selectedConversation.id)}
              onToggleMute={() => selectedConversation && actions.toggleMute(selectedConversation.id)}
              onMarkAsRead={() => selectedConversation && actions.markAsRead(selectedConversation.id)}
            />
          </>
        ) : (
          <>
            {state.selectedContactId === SPECIAL_KEYS.newFriends ? (
              <ContactSpecialPanel
                type="new-friends"
                themeColor={THEME_COLOR}
                newFriends={state.newFriends}
                onBack={actions.showSidebar}
              />
            ) : state.selectedContactId === SPECIAL_KEYS.savedGroups ? (
              <ContactSpecialPanel
                type="saved-groups"
                themeColor={THEME_COLOR}
                savedGroups={state.savedGroups}
                onBack={actions.showSidebar}
              />
            ) : state.selectedContactId === SPECIAL_KEYS.blacklist ? (
              <ContactSpecialPanel
                type="blacklist"
                themeColor={THEME_COLOR}
                blacklist={state.blacklist}
                onBack={actions.showSidebar}
              />
            ) : (
              <ContactDetail
                themeColor={THEME_COLOR}
                contact={selectedContact}
                isMobile={state.isMobile}
                onBack={actions.showSidebar}
              />
            )}
          </>
        )}
      </div>

      <GlobalSearch
        open={state.showGlobalSearch}
        conversations={state.conversations}
        contacts={state.contacts}
        onClose={actions.closeGlobalSearch}
        onSelect={handleGlobalSearchSelect}
      />
    </div>
  )
}

export default ChatShell
