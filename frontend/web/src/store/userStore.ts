import type { User } from '@/modules/user'
import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'
import { immer } from 'zustand/middleware/immer'

interface UserState {
  token: string
  setToken: (token: string) => void
  clearToken: () => void
  getToken: () => string
  user: User
  setUser: (user: User) => void
  clearUser: () => void
  getUser: () => User
}

export const useUserStore = create<UserState>()(
  immer(
    persist(
      (set, get) => ({
        user: {} as User,
        token: '',
        setToken: (token: string) => set({ token }),
        clearToken: () => set({ token: '' }),
        getToken: () => get().token,
        setUser: (user: User) => set({ user }),
        clearUser: () => set({ user: {} as User }),
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
