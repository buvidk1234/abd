import { useNavigate } from 'react-router'
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
    <div className="container mx-auto p-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">欢迎回来，{user?.username || '用户'}</h1>
        <button
          onClick={handleLogout}
          className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
        >
          退出登录
        </button>
      </div>
      <div className="bg-card p-6 rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">主页</h2>
        <p>这里是主页内容</p>
      </div>
    </div>
  )
}

export default HomePage
