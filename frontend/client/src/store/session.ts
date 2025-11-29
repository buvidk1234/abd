import type { UserProfile } from '@/types/backend'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useSessionStore = defineStore(
  'session',
  () => {
    const currentUser = ref<UserProfile | null>(null)
    const recentUsers = ref<UserProfile[]>([])

    const isLoggedIn = computed(() => !!currentUser.value?.userID)

    const upsertRecentUser = (user: UserProfile) => {
      if (!user.userID) {
        return
      }
      const existIdx = recentUsers.value.findIndex(item => item.userID === user.userID)
      if (existIdx > -1) {
        recentUsers.value[existIdx] = { ...recentUsers.value[existIdx], ...user }
      }
      else {
        recentUsers.value.unshift(user)
      }
      // 最多保留最近的10个
      recentUsers.value = recentUsers.value.slice(0, 10)
    }

    const setCurrentUser = (user: UserProfile | null) => {
      currentUser.value = user
      if (user) {
        upsertRecentUser(user)
      }
    }

    const updateCurrentProfile = (profile: Partial<UserProfile>) => {
      if (!currentUser.value) {
        return
      }
      currentUser.value = { ...currentUser.value, ...profile }
      upsertRecentUser(currentUser.value)
    }

    const clearSession = () => {
      currentUser.value = null
    }

    return {
      currentUser,
      recentUsers,
      isLoggedIn,
      setCurrentUser,
      updateCurrentProfile,
      clearSession,
    }
  },
  {
    persist: true,
  },
)
