import { Link, useNavigate } from 'react-router'
import { toast } from 'sonner'

import { AuthForm, type LoginFormValues } from './components/AuthForm'
import { AuthLayout } from './components/AuthLayout'
import { login, getCurrentUser } from '@/modules/user'
import { useUserStore } from '@/store/userStore'

function LoginPage() {
  const navigate = useNavigate()

  const handleLogin = async (values: LoginFormValues) => {
    const token = await login(values)

    useUserStore.getState().setToken(token)
    const profile = await getCurrentUser()
    useUserStore.getState().setUser(profile)
    toast('登录成功')
    navigate('/')
  }

  return (
    <AuthLayout title="登陆" description="更愉快的与朋友交流">
      <AuthForm
        mode="login"
        submitText="登录"
        onSubmit={handleLogin}
        footerSlot={
          <span>
            还没有账号？{' '}
            <Link to="/register" className="text-primary hover:underline">
              立即注册
            </Link>
          </span>
        }
      />
    </AuthLayout>
  )
}

export default LoginPage
