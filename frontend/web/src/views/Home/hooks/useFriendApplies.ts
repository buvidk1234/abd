import { useEffect, useCallback } from 'react'
import { useImmer } from 'use-immer'
import { listFriendApplies, respondFriendApply } from '@/modules'
import type { FriendRequest } from '@/modules'
import { toast } from 'sonner'

interface FriendAppliesState {
  applies: FriendRequest[]
  loading: boolean
  error: string | null
}

export function useFriendApplies() {
  const [state, setState] = useImmer<FriendAppliesState>({
    applies: [],
    loading: false,
    error: null,
  })

  // 加载好友申请列表
  const loadApplies = useCallback(async () => {
    setState((draft) => {
      draft.loading = true
      draft.error = null
    })

    try {
      const result = await listFriendApplies({ page: 1, pageSize: 100 })
      setState((draft) => {
        draft.applies = result.data
        draft.loading = false
      })
    } catch (error) {
      setState((draft) => {
        draft.error = '加载失败'
        draft.loading = false
      })
      toast.error('加载好友申请失败')
      console.log('加载好友申请失败', error)
    }
  }, [setState])

  useEffect(() => {
    loadApplies()
  }, [loadApplies])

  const handleApply = useCallback(
    async (applyId: string, action: 'accept' | 'reject') => {
      try {
        await respondFriendApply({
          id: applyId,
          handleResult: action === 'accept' ? 1 : 2,
        })

        await loadApplies()
      } catch (error) {
        toast.error('响应好友申请失败')
        console.log('响应好友申请失败', error)
        throw error
      }
    },
    [loadApplies]
  )

  return {
    applies: state.applies,
    loading: state.loading,
    error: state.error,
    refresh: loadApplies,
    handleApply,
  }
}
