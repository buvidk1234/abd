import { Outlet } from 'react-router'
import { Toaster } from 'sonner'
import { useInitAuth } from '@/hooks/useInitAuth'

function App() {
  const { isReady } = useInitAuth()

  if (!isReady) {
    return (
      <div className="min-h-screen bg-background text-foreground flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">加载中...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background text-foreground">
      <Toaster theme="system" />
      <Outlet />
    </div>
  )
}

export default App
