import { useParams, useNavigate, useOutletContext } from 'react-router'
import { useEffect, useRef } from 'react'
import { useImmer } from 'use-immer'
import { toast } from 'sonner'

import { THEME_COLOR } from '../constants'
import { CHAT_MENUS } from '../constants/menus'
import { Sidebar } from '../components/chat/Sidebar'
import { ChatPanel } from '../components/chat/ChatPanel'
import { ChannelPanel } from '../components/chat/ChannelPanel'
import { useMockData } from '../hooks/useMockData'
import { useClickOutside } from '../hooks/useClickOutside'
import { useResponsive } from '../hooks/useResponsive'
import { useWS } from '../hooks/useWS'

interface LayoutContext {
  onOpenGlobalSearch: () => void
}

interface UIState {
  showAddMenu: boolean
  showDetails: boolean
  showSidebarOnMobile: boolean
}

export function ChatPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { onOpenGlobalSearch } = useOutletContext<LayoutContext>()
  const isMobile = useResponsive()
  const addMenuRef = useRef<HTMLDivElement | null>(null)

  const { conversations, actions: dataActions } = useMockData()
  const selectedConversation = conversations.find((c) => c.id === id) ?? null

  const [uiState, setUIState] = useImmer<UIState>({
    showAddMenu: false,
    showDetails: false,
    showSidebarOnMobile: true,
  })

  // 移动端
  if (isMobile && !uiState.showSidebarOnMobile && !id) {
    setUIState((draft) => {
      draft.showSidebarOnMobile = true
    })
  }

  // 点击外部区域关闭添加菜单
  useClickOutside(
    addMenuRef,
    () =>
      setUIState((draft) => {
        draft.showAddMenu = false
      }),
    uiState.showAddMenu
  )

  const handleSelectConversation = (convId: string) => {
    navigate(`/chat/${convId}`)
    if (isMobile) {
      setUIState((draft) => {
        draft.showSidebarOnMobile = false
      })
    }
  }

  const handleMenuClick = (menuId: string) => {
    const menu = CHAT_MENUS.find((item) => item.id === menuId)
    if (menu) {
      toast.info(`${menu.title} 功能正在对接中`)
    }
    setUIState((draft) => {
      draft.showAddMenu = false
    })
  }

  const { send } = useWS()
  useEffect(() => {
    setInterval(() => {
      send({
        req_identifier: 4001,
        data: {
          id: id,
        }
      })
    }, 2000)
  }, [])

  return (
    <div className="relative flex min-w-0 flex-1">
      <Sidebar
        ref={addMenuRef}
        themeColor={THEME_COLOR}
        conversations={conversations}
        selectedId={id ?? null}
        onSelect={handleSelectConversation}
        showAddMenu={uiState.showAddMenu}
        onToggleAddMenu={() =>
          setUIState((draft) => {
            draft.showAddMenu = !draft.showAddMenu
          })
        }
        onMenuClick={handleMenuClick}
        onOpenSearch={onOpenGlobalSearch}
        showSidebar={uiState.showSidebarOnMobile}
        isMobile={isMobile}
      />

      <ChatPanel
        themeColor={THEME_COLOR}
        conversation={selectedConversation}
        isMobile={isMobile}
        onBack={() =>
          setUIState((draft) => {
            draft.showSidebarOnMobile = true
          })
        }
        onToggleDetails={() =>
          setUIState((draft) => {
            draft.showDetails = !draft.showDetails
          })
        }
        onOpenDetails={() =>
          setUIState((draft) => {
            draft.showDetails = true
          })
        }
      />

      <ChannelPanel
        themeColor={THEME_COLOR}
        open={uiState.showDetails && !!selectedConversation}
        conversation={selectedConversation}
        onClose={() =>
          setUIState((draft) => {
            draft.showDetails = false
          })
        }
        onTogglePin={() => selectedConversation && dataActions.togglePin(selectedConversation.id)}
        onToggleMute={() => selectedConversation && dataActions.toggleMute(selectedConversation.id)}
        onMarkAsRead={() => selectedConversation && dataActions.markAsRead(selectedConversation.id)}
      />
    </div>
  )
}
