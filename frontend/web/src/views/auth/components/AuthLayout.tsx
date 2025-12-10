import type { PropsWithChildren, ReactNode } from 'react'

type AuthLayoutProps = PropsWithChildren<{
  title: string
  description?: string
  footer?: ReactNode
}>

function AuthLayout({ title, description, footer, children }: AuthLayoutProps) {
  return (
    <div className="min-h-screen bg-gradient-to-br from-primary/5 via-background to-accent/20">
      <div className="mx-auto flex max-w-5xl flex-col items-center justify-center px-4 py-12 sm:px-6 lg:px-8">
        <div className="w-full max-w-md space-y-6">
          <div className="text-center">
            <div className="mx-auto mb-4 flex size-14 items-center justify-center rounded-full bg-primary/10 text-lg font-semibold text-primary shadow-sm">
              TS
            </div>
            <h1 className="text-2xl font-semibold tracking-tight text-foreground">{title}</h1>
            {description ? (
              <p className="mt-2 text-sm text-muted-foreground">{description}</p>
            ) : null}
          </div>

          <div className="rounded-xl border bg-card p-6 shadow-lg shadow-primary/5">{children}</div>

          {footer ? (
            <div className="text-center text-sm text-muted-foreground">{footer}</div>
          ) : null}
        </div>
      </div>
    </div>
  )
}

export { AuthLayout }
