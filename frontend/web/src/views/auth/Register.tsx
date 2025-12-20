import { Link, useNavigate } from 'react-router'
import { toast } from 'sonner'

import { AuthForm, type RegisterFormValues } from './components/AuthForm'
import { AuthLayout } from './components/AuthLayout'
import { register } from '@/modules/user'

function RegisterPage() {
  const navigate = useNavigate()

  const handleRegister = async (values: RegisterFormValues) => {
    await register(values)
    toast('注册成功，请登录')
    navigate('/login')
  }

  return (
    <AuthLayout title="注册" description="创建你的账号，开始与朋友交流">
      <AuthForm
        mode="register"
        submitText="注册"
        onSubmit={handleRegister}
        footerSlot={
          <span>
            已有账号？{' '}
            <Link to="/login" className="text-primary hover:underline">
              立即登录
            </Link>
          </span>
        }
      />
    </AuthLayout>
  )
}

export default RegisterPage
