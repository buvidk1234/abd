import { useParams, useNavigate, useOutletContext } from 'react-router'
import { useEffect, useRef, useMemo } from 'react'
import { useImmer } from 'use-immer'
import { toast } from 'sonner'

import { THEME_COLOR } from '../constants'
import { CHAT_MENUS } from '../constants/menus'
import { Sidebar } from '../components/chat/Sidebar'
import { ChatPanel } from '../components/chat/ChatPanel'
import { ChannelPanel } from '../components/chat/ChannelPanel'
import { useClickOutside } from '../hooks/useClickOutside'
import { useResponsive } from '../hooks/useResponsive'
import { useConversations, useMessageContext } from '../hooks'
import { useUserStore } from '@/store/userStore'
import { adaptMessageToItem } from '../utils/formatters'
import type { ConversationItem } from '../types'

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

  const { conversations, actions: dataActions, getOrCreateConversation } = useConversations()
  const { user } = useUserStore()
  const { conversationMessages, loadMoreMessages, loading, setOnNewMessage, maxSeqs } =
    useMessageContext()

  const [uiState, setUIState] = useImmer<UIState>({
    showAddMenu: false,
    showDetails: false,
    showSidebarOnMobile: true,
  })

  // 追踪已加载过消息的会话
  const loadedConversationsRef = useRef<Set<string>>(new Set())

  // 如果会话不存在，尝试创建（用于从联系人页面跳转）
  useEffect(() => {
    if (!id) return
    console.log(conversations)

    const conv = conversations.find((c) => c.id === id)
    if (!conv && getOrCreateConversation) {
      getOrCreateConversation(id)
    }
  }, [id, conversations, getOrCreateConversation])

  // 获取当前会话，并合并真实消息
  const selectedConversation = useMemo((): ConversationItem | null => {
    const conv = conversations.find((c) => c.id === id)
    if (!conv) return null

    // 获取真实消息
    const realMessages = conversationMessages[id || ''] || []
    console.log('realMessages', realMessages)

    const adaptedMessages = realMessages.map((msg) => adaptMessageToItem(msg, user.id))

    return {
      ...conv,
      messages: adaptedMessages,
    }
  }, [conversations, conversationMessages, id, user.id])

  // 监听新消息
  useEffect(() => {
    setOnNewMessage((msg) => {
      if (msg.sender_id !== user.id) {
        toast.success(`收到新消息: ${msg.content.slice(0, 20)}...`)
      }
    })
  }, [setOnNewMessage])

  // 当会话 ID 和 maxSeq 都准备好时，自动加载消息
  useEffect(() => {
    if (!id) return

    // 如果已经加载过，跳过
    if (loadedConversationsRef.current.has(id)) return

    // 检查该会话是否有 maxSeq
    const hasMaxSeq = maxSeqs[id] !== undefined

    if (hasMaxSeq) {
      loadMoreMessages(id, 20)
      loadedConversationsRef.current.add(id)
    }
  }, [id, maxSeqs, loadMoreMessages])

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

  const handleLoadMore = () => {
    if (id && !loading[id]) {
      loadMoreMessages(id, 20)
    }
  }

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
        onLoadMore={handleLoadMore}
        isLoadingMore={id ? loading[id] : false}
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
