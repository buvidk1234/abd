import { useMemo } from 'react'
import { useImmer } from 'use-immer'

import type {
  BlacklistItem,
  FriendItem,
  ConversationItem,
  FriendRequestItem,
  SavedGroupItem,
} from '../types'
import {
  createBlacklist,
  createFriends,
  createConversations,
  createNewFriends,
  createSavedGroups,
} from '../utils/mockData'

interface ChatShellState {
  activeTab: 'chat' | 'friends'
  conversations: ConversationItem[]
  friends: FriendItem[]
  newFriends: FriendRequestItem[]
  savedGroups: SavedGroupItem[]
  blacklist: BlacklistItem[]
  selectedConversationId: string | null
  selectedFriendId: string | null
  showAddMenu: boolean
  showGlobalSearch: boolean
  showDetails: boolean
  isMobile: boolean
  showSidebarOnMobile: boolean
}

export function useChatShellState(isMobile: boolean) {
  const [state, setState] = useImmer<ChatShellState>(() => {
    const conversations = createConversations()
    const friends = createFriends()
    return {
      activeTab: 'chat',
      conversations,
      friends: friends,
      newFriends: createNewFriends(),
      savedGroups: createSavedGroups(),
      blacklist: createBlacklist(),
      selectedConversationId: conversations[0]?.id ?? null,
      selectedFriendId: friends[0]?.id ?? null,
      showAddMenu: false,
      showGlobalSearch: false,
      showDetails: false,
      isMobile,
      showSidebarOnMobile: true,
    }
  })

  // Sync isMobile state
  if (state.isMobile !== isMobile) {
    setState((draft) => {
      draft.isMobile = isMobile
      if (!isMobile) {
        draft.showSidebarOnMobile = true
      }
    })
  }

  const selectedConversation = useMemo(
    () => state.conversations.find((item) => item.id === state.selectedConversationId) ?? null,
    [state.conversations, state.selectedConversationId]
  )

  const selectedFriend = useMemo(
    () => state.friends.find((item) => item.id === state.selectedFriendId) ?? null,
    [state.friends, state.selectedFriendId]
  )

  const actions = {
    setActiveTab: (tab: 'chat' | 'friends') => {
      setState((draft) => {
        draft.activeTab = tab
        draft.showSidebarOnMobile = true
        draft.showAddMenu = false
      })
    },

    selectConversation: (id: string) => {
      setState((draft) => {
        draft.activeTab = 'chat'
        draft.selectedConversationId = id
        draft.showDetails = false
        if (draft.isMobile) {
          draft.showSidebarOnMobile = false
        }
      })
    },

    selectFriend: (id: string) => {
      setState((draft) => {
        draft.activeTab = 'friends'
        draft.selectedFriendId = id
        if (draft.isMobile) {
          draft.showSidebarOnMobile = false
        }
      })
    },

    togglePin: (id: string) => {
      setState((draft) => {
        const target = draft.conversations.find((item) => item.id === id)
        if (target) {
          target.pinned = !target.pinned
        }
      })
    },

    toggleMute: (id: string) => {
      setState((draft) => {
        const target = draft.conversations.find((item) => item.id === id)
        if (target) {
          target.muted = !target.muted
        }
      })
    },

    markAsRead: (id: string) => {
      setState((draft) => {
        const target = draft.conversations.find((item) => item.id === id)
        if (target) {
          target.unread = 0
        }
      })
    },

    toggleAddMenu: () => {
      setState((draft) => {
        draft.showAddMenu = !draft.showAddMenu
      })
    },

    closeAddMenu: () => {
      setState((draft) => {
        draft.showAddMenu = false
      })
    },

    openGlobalSearch: () => {
      setState((draft) => {
        draft.showGlobalSearch = true
      })
    },

    closeGlobalSearch: () => {
      setState((draft) => {
        draft.showGlobalSearch = false
      })
    },

    toggleDetails: () => {
      setState((draft) => {
        draft.showDetails = !draft.showDetails
      })
    },

    openDetails: () => {
      setState((draft) => {
        draft.showDetails = true
      })
    },

    closeDetails: () => {
      setState((draft) => {
        draft.showDetails = false
      })
    },

    showSidebar: () => {
      setState((draft) => {
        draft.showSidebarOnMobile = true
      })
    },

    hideSidebar: () => {
      setState((draft) => {
        draft.showSidebarOnMobile = false
      })
    },
  }

  return {
    state,
    selectedConversation,
    selectedFriend,
    actions,
  }
}
