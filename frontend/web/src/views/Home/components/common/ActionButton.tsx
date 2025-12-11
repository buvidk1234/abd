import type { ReactNode } from 'react'
import clsx from 'clsx'

interface ActionButtonProps {
  icon?: ReactNode
  label: string
  onClick: () => void
  variant?: 'solid' | 'ghost' | 'outline'
  size?: 'sm' | 'md' | 'lg'
  themeColor?: string
  disabled?: boolean
}

const sizeMap = {
  sm: 'px-2 py-0.5 text-[10px]',
  md: 'px-3 py-1 text-xs',
  lg: 'px-4 py-2 text-sm',
}

export function ActionButton({
  icon,
  label,
  onClick,
  variant = 'solid',
  size = 'md',
  themeColor = '#e46342',
  disabled,
}: ActionButtonProps) {
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={clsx(
        'flex items-center gap-2 rounded-full font-semibold transition',
        sizeMap[size],
        variant === 'solid' && 'bg-slate-900 text-white shadow-sm',
        variant === 'ghost' && 'bg-white text-slate-600 hover:bg-slate-100',
        variant === 'outline' && 'border border-slate-200 bg-white text-slate-700 hover:border-slate-300 hover:bg-slate-50',
        disabled && 'cursor-not-allowed opacity-50'
      )}
      style={
        variant === 'solid' && !disabled
          ? { background: themeColor, boxShadow: `0 10px 30px ${themeColor}22` }
          : undefined
      }
    >
      {icon && <span className="flex items-center justify-center">{icon}</span>}
      {label}
    </button>
  )
}
