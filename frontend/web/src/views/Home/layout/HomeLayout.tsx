import { Outlet, useLocation, useNavigate } from 'react-router'
import { toast } from 'sonner'
import { useImmer } from 'use-immer'

import { useUserStore } from '@/store/userStore'
import { THEME_COLOR } from '../constants'
import { NavRail } from '../components/shared/NavRail'
import { GlobalSearch } from '../components/shared/GlobalSearch'
import { useMockData } from '../hooks/useMockData'
// import { useWebSocket } from '@/hooks/useWebSocket'
import { WebSocketProvider } from '../providers/WebSocketProvider'
export function HomeLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const userName = useUserStore((state) => state.user?.nickname || state.user?.username || '用户')
  const [showGlobalSearch, setShowGlobalSearch] = useImmer(false)

  const { conversations, friends } = useMockData()

  const { user, token } = useUserStore()

  // Determine active tab from route
  const activeTab = location.pathname.startsWith('/contact') ? 'friends' : 'chat'

  const handleSelectTab = (tab: 'chat' | 'friends') => {
    if (tab === 'chat') {
      // Navigate to first conversation or empty chat
      const firstConv = conversations[0]
      navigate(firstConv ? `/chat/${firstConv.id}` : '/chat')
    } else {
      // Navigate to first friend or special panel
      const firstFriend = friends[0]
      navigate(firstFriend ? `/contact/${firstFriend.id}` : '/contact/special:new-friends')
    }
  }

  const handleGlobalSearchSelect = (payload: { type: 'conversation' | 'friend'; id: string }) => {
    if (payload.type === 'conversation') {
      navigate(`/chat/${payload.id}`)
    } else {
      navigate(`/contact/${payload.id}`)
    }
    setShowGlobalSearch(() => false)
  }

  return (
    <WebSocketProvider token={token} userId={user.id}>
      <div className="flex h-screen w-full overflow-hidden bg-slate-100 text-slate-900">
        <NavRail
          themeColor={THEME_COLOR}
          activeTab={activeTab}
          userName={userName}
          onSelectTab={handleSelectTab}
          onOpenSettings={() => toast.info('设置功能稍后接入')}
        />

        <Outlet
          context={{
            onOpenGlobalSearch: () => setShowGlobalSearch(() => true),
          }}
        />

        <GlobalSearch
          open={showGlobalSearch}
          conversations={conversations}
          friends={friends}
          onClose={() => setShowGlobalSearch(() => false)}
          onSelect={handleGlobalSearchSelect}
        />
      </div>
    </WebSocketProvider>
  )
}
