import { useMemo, useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'

import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'

const loginSchema = z.object({
  username: z.string().trim().min(1, '请输入账号'),
  password: z.string().min(1, '请输入密码'),
})

const emptyToUndefined = (v: unknown) => (typeof v === 'string' && v.trim() === '' ? undefined : v)

const registerSchema = loginSchema
  .extend({
    confirmPassword: z.string().min(1, '请再次输入密码'),

    nickname: z.preprocess(emptyToUndefined, z.string().trim().optional()),

    phone: z
      .preprocess(emptyToUndefined, z.string().trim().optional())
      .refine((v) => !v || /^1[3-9]\d{9}$/.test(v), '手机号格式不正确'),

    email: z
      .preprocess(emptyToUndefined, z.string().trim().optional())
      .refine((v) => !v || /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v), '邮箱格式不正确'),
  })
  .refine((values) => values.password === values.confirmPassword, {
    path: ['confirmPassword'],
    message: '两次输入的密码不一致',
  })

export type LoginFormValues = z.infer<typeof loginSchema>
export type RegisterFormValues = z.infer<typeof registerSchema>

type AuthFormProps =
  | {
      mode: 'login'
      onSubmit: (values: LoginFormValues) => Promise<void> | void
      submitText?: string
      footerSlot?: React.ReactNode
    }
  | {
      mode: 'register'
      onSubmit: (values: RegisterFormValues) => Promise<void> | void
      submitText?: string
      footerSlot?: React.ReactNode
    }

function AuthForm(props: AuthFormProps) {
  const { mode, onSubmit, submitText = mode === 'login' ? '登录' : '注册', footerSlot } = props
  const [submitting, setSubmitting] = useState(false)

  const schema = useMemo(() => (mode === 'login' ? loginSchema : registerSchema), [mode])

  const form = useForm<LoginFormValues | RegisterFormValues>({
    resolver: zodResolver(schema),
    defaultValues:
      mode === 'login'
        ? {
            username: '',
            password: '',
          }
        : {
            username: '',
            password: '',
            confirmPassword: '',
            nickname: '',
            phone: '',
            email: '',
          },
  })

  const handleSubmit = form.handleSubmit(async (values) => {
    setSubmitting(true)
    try {
      await onSubmit(values as LoginFormValues & RegisterFormValues)
    } finally {
      setSubmitting(false)
    }
  })

  const isRegister = mode === 'register'

  return (
    <Form {...form}>
      <form onSubmit={handleSubmit} className="space-y-4">
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>账号</FormLabel>
              <FormControl>
                <Input placeholder="请输入手机号/用户名" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>密码</FormLabel>
              <FormControl>
                <Input type="password" placeholder="请输入密码" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {isRegister ? (
          <>
            <FormField
              control={form.control}
              name="confirmPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>确认密码</FormLabel>
                  <FormControl>
                    <Input type="password" placeholder="再次输入密码" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="grid gap-4 sm:grid-cols-2">
              <FormField
                control={form.control}
                name="nickname"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>昵称</FormLabel>
                    <FormControl>
                      <Input placeholder="选填" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="phone"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>手机号</FormLabel>
                    <FormControl>
                      <Input placeholder="选填" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <FormLabel>邮箱</FormLabel>
                    <FormControl>
                      <Input type="email" placeholder="选填" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
          </>
        ) : null}

        <Button type="submit" className="w-full" disabled={submitting}>
          {submitting ? '提交中...' : submitText}
        </Button>

        {footerSlot ? (
          <div className="text-center text-sm text-muted-foreground">{footerSlot}</div>
        ) : null}
      </form>
    </Form>
  )
}

export { AuthForm }
