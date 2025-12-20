import clsx from 'clsx'

interface BadgeProps {
  children: React.ReactNode
  variant?: 'default' | 'primary' | 'success' | 'warning' | 'danger' | 'info'
  size?: 'sm' | 'md' | 'lg'
  customColor?: string
}

const sizeMap = {
  sm: 'px-2 py-0.5 text-[10px]',
  md: 'px-2 py-0.5 text-[11px]',
  lg: 'px-3 py-1 text-xs',
}

const variantMap = {
  default: 'bg-slate-100 text-slate-600',
  primary: 'bg-orange-50 text-[#e46342]',
  success: 'bg-emerald-50 text-emerald-600',
  warning: 'bg-amber-50 text-amber-600',
  danger: 'bg-red-50 text-red-600',
  info: 'bg-blue-50 text-blue-600',
}

export function Badge({ children, variant = 'default', size = 'md', customColor }: BadgeProps) {
  return (
    <span
      className={clsx(
        'inline-flex items-center rounded-full font-semibold',
        sizeMap[size],
        variantMap[variant]
      )}
      style={customColor ? { backgroundColor: `${customColor}15`, color: customColor } : undefined}
    >
      {children}
    </span>
  )
}
