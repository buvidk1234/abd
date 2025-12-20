import { Outlet } from 'react-router'
import { Toaster } from 'sonner'
import { useInitAuth } from '@/hooks/useInitAuth'

function App() {
  const { isReady } = useInitAuth()

  return (
    <div className="min-h-screen bg-background text-foreground">
      <Toaster theme="system" />
      {isReady ? (
        <Outlet />
      ) : (
        <div className="flex min-h-screen items-center justify-center">
          <div className="text-center">
            <div className="mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-b-2 border-primary"></div>
            <p className="text-muted-foreground">加载中...</p>
          </div>
        </div>
      )}
    </div>
  )
}

export default App
