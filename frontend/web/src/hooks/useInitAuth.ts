import { useEffect, useRef } from 'react'
import { useLocation, useNavigate } from 'react-router'
import { toast } from 'sonner'
import { getUserInfo } from '@/services/api/user'
import { useUserStore } from '@/store/userStore'
import { useImmer } from 'use-immer'

export function useInitAuth() {
  const [isReady, setIsReady] = useImmer<boolean>(false)
  const location = useLocation()
  const navigate = useNavigate()
  const hasInitialized = useRef(false)

  useEffect(() => {
    // 只在首次加载时执行
    if (hasInitialized.current) {
      return
    }
    hasInitialized.current = true

    const initAuth = async () => {
      const token = useUserStore.getState().getToken()
      const isAuthPage = location.pathname === '/login' || location.pathname === '/register'

      if (isAuthPage) {
        // 在登录/注册页面
        if (token) {
          // 验证 token 是否有效
          try {
            const { data: profile } = await getUserInfo()
            // token 有效，更新用户信息并跳转首页
            useUserStore.getState().setUser(profile)
            toast.warning('请先退出后再登录')
            navigate('/', { replace: true })
          } catch {
            // token 无效，清除 token 和用户信息
            useUserStore.getState().clearToken()
            useUserStore.getState().clearUser()
          }
        }
        setIsReady(true)
      } else {
        // 在其他页面
        if (token) {
          // 验证 token 是否有效
          try {
            const { data: profile } = await getUserInfo()
            // token 有效，更新用户信息
            useUserStore.getState().setUser(profile)
            setIsReady(true)
          } catch {
            // token 无效，清除 token 和用户信息
            useUserStore.getState().clearToken()
            useUserStore.getState().clearUser()
            toast.error('登录过期，请重新登录')
            navigate('/login', { replace: true })
          }
        } else {
          // token 不存在，跳转登录
          toast.error('请先登录')
          navigate('/login', { replace: true })
        }
      }
    }

    initAuth()
  }, [location.pathname, navigate])

  return { isReady }
}

