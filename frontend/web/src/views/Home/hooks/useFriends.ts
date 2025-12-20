import { listFriends, type Friend } from '@/modules'
import { useCallback, useEffect } from 'react'
import { toast } from 'sonner'
import { useImmer } from 'use-immer'

interface FriendsState {
  friends: Friend[]
  loading: boolean
  error: string | null
}

export function useFriends() {
  const [state, setState] = useImmer<FriendsState>({
    friends: [],
    loading: false,
    error: null,
  })

  const loadFriends = useCallback(async () => {
    setState((draft) => {
      draft.loading = true
      draft.error = null
    })

    try {
      const result = await listFriends({ page: 1, pageSize: 100 })
      setState((draft) => {
        draft.friends = result.data
        draft.loading = false
      })
    } catch (error) {
      setState((draft) => {
        draft.error = '加载失败'
        draft.loading = false
      })
      toast.error('加载好友失败')
      console.log('加载好友失败', error)
    }
  }, [setState])

  useEffect(() => {
    loadFriends()
  }, [loadFriends])

  return {
    friends: state.friends,
    loading: state.loading,
    error: state.error,
    refresh: loadFriends,
  }
}
