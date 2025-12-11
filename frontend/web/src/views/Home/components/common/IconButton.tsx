import type { ReactNode } from 'react'
import clsx from 'clsx'

interface IconButtonProps {
  icon: ReactNode
  onClick?: () => void
  ariaLabel?: string
  active?: boolean
  activeColor?: string
  disabled?: boolean
  size?: 'sm' | 'md' | 'lg'
}

const sizeMap = {
  sm: 'size-8',
  md: 'size-10',
  lg: 'size-12',
}

const iconSizeMap = {
  sm: 'size-3.5',
  md: 'size-4',
  lg: 'size-5',
}

export function IconButton({
  icon,
  onClick,
  ariaLabel,
  active,
  activeColor = '#e46342',
  disabled,
  size = 'md',
}: IconButtonProps) {
  return (
    <button
      aria-label={ariaLabel}
      onClick={onClick}
      disabled={disabled}
      className={clsx(
        'flex items-center justify-center rounded-xl text-slate-600 transition',
        sizeMap[size],
        active ? 'ring-2 ring-slate-100' : '',
        disabled ? 'cursor-not-allowed opacity-50' : 'hover:bg-slate-100'
      )}
      style={active ? { color: activeColor } : undefined}
    >
      <span className={iconSizeMap[size]}>{icon}</span>
    </button>
  )
}
