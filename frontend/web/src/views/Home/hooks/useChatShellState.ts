import { useMemo } from 'react'
import { useImmer } from 'use-immer'

import type {
  BlacklistItem,
  ContactItem,
  ConversationItem,
  FriendRequestItem,
  SavedGroupItem,
} from '../types'
import {
  createBlacklist,
  createContacts,
  createConversations,
  createNewFriends,
  createSavedGroups,
} from '../utils/mockData'

interface ChatShellState {
  activeTab: 'chat' | 'contacts'
  conversations: ConversationItem[]
  contacts: ContactItem[]
  newFriends: FriendRequestItem[]
  savedGroups: SavedGroupItem[]
  blacklist: BlacklistItem[]
  selectedConversationId: string | null
  selectedContactId: string | null
  showAddMenu: boolean
  showGlobalSearch: boolean
  showDetails: boolean
  isMobile: boolean
  showSidebarOnMobile: boolean
}

export function useChatShellState(isMobile: boolean) {
  const [state, setState] = useImmer<ChatShellState>(() => {
    const conversations = createConversations()
    const contacts = createContacts()
    return {
      activeTab: 'chat',
      conversations,
      contacts,
      newFriends: createNewFriends(),
      savedGroups: createSavedGroups(),
      blacklist: createBlacklist(),
      selectedConversationId: conversations[0]?.id ?? null,
      selectedContactId: contacts[0]?.id ?? null,
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

  const selectedContact = useMemo(
    () => state.contacts.find((item) => item.id === state.selectedContactId) ?? null,
    [state.contacts, state.selectedContactId]
  )

  const actions = {
    setActiveTab: (tab: 'chat' | 'contacts') => {
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

    selectContact: (id: string) => {
      setState((draft) => {
        draft.activeTab = 'contacts'
        draft.selectedContactId = id
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
    selectedContact,
    actions,
  }
}
