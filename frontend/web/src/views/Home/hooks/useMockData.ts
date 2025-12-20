import { useMemo } from 'react'
import { useImmer } from 'use-immer'
import type {
  ConversationItem,
  FriendItem,
  FriendRequestItem,
  SavedGroupItem,
  BlacklistItem,
} from '../types'
import {
  createConversations,
  createFriends,
  createNewFriends,
  createSavedGroups,
  createBlacklist,
} from '../utils/mockData'

interface MockDataState {
  conversations: ConversationItem[]
  friends: FriendItem[]
  newFriends: FriendRequestItem[]
  savedGroups: SavedGroupItem[]
  blacklist: BlacklistItem[]
}

export function useMockData() {
  const [state, setState] = useImmer<MockDataState>(() => ({
    conversations: createConversations(),
    friends: createFriends(),
    newFriends: createNewFriends(),
    savedGroups: createSavedGroups(),
    blacklist: createBlacklist(),
  }))

  const actions = useMemo(
    () => ({
      // Conversation actions
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

      // Friend actions
      addFriend: (friend: FriendItem) => {
        setState((draft) => {
          draft.friends.push(friend)
        })
      },

      removeFriend: (id: string) => {
        setState((draft) => {
          draft.friends = draft.friends.filter((f) => f.id !== id)
        })
      },

      updateFriend: (id: string, updates: Partial<FriendItem>) => {
        setState((draft) => {
          const target = draft.friends.find((f) => f.id === id)
          if (target) {
            Object.assign(target, updates)
          }
        })
      },
    }),
    [setState]
  )

  return {
    ...state,
    actions,
  }
}
