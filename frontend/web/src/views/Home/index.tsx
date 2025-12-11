import { useNavigate } from 'react-router'

import { Button } from '@/components/ui/button'
import { ChatShell } from './components/ChatShell'
import { useUserStore } from '@/store/userStore'

function HomePage() {
  const navigate = useNavigate()
  const { user, clearUser, clearToken } = useUserStore()

  const handleLogout = () => {
    clearUser()
    clearToken()
    navigate('/login')
  }

  return (
    <div className="relative min-h-screen bg-slate-100">
      <div className="absolute right-4 top-4 z-30 flex items-center gap-2 rounded-full bg-white/80 px-3 py-2 text-xs text-slate-600 shadow-md backdrop-blur">
        <span className="font-semibold text-slate-800">{user?.username || '用户'}</span>
        <Button size="sm" variant="ghost" onClick={handleLogout}>
          退出
        </Button>
      </div>
      <ChatShell />
    </div>
  )
}

export default HomePage
