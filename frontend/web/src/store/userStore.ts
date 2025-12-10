import type { UserInfo } from '@/services/api/user'
import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'
import { immer } from 'zustand/middleware/immer'

interface UserState {
  token: string
  setToken: (token: string) => void
  clearToken: () => void
  getToken: () => string
  user: UserInfo
  setUser: (user: UserInfo) => void
  clearUser: () => void
  getUser: () => UserInfo
}

export const useUserStore = create<UserState>()(
  immer(
    persist(
      (set, get) => ({
        user: {} as UserInfo,
        token: '',
        setToken: (token: string) => set({ token }),
        clearToken: () => set({ token: '' }),
        getToken: () => get().token,
        setUser: (user: UserInfo) => set({ user }),
        clearUser: () => set({ user: {} as UserInfo }),
        getUser: () => get().user,
      }),
      {
        name: 'user',
        storage: createJSONStorage(() => localStorage),
        partialize: (state) => ({
          user: state.user,
          token: state.token,
        }),
      }
    )
  )
)
