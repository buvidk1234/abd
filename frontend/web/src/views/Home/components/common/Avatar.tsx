import clsx from 'clsx'

interface AvatarProps {
  name: string
  avatar?: string
  accent?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
  status?: 'online' | 'busy' | 'away' | 'offline'
  online?: boolean
  muted?: boolean
  selected?: boolean
  themeColor?: string
  onClick?: () => void
}

const sizeMap = {
  sm: 'size-9',
  md: 'size-11',
  lg: 'size-12',
  xl: 'size-16',
}

const textSizeMap = {
  sm: 'text-xs',
  md: 'text-xs',
  lg: 'text-sm',
  xl: 'text-base',
}

const statusSizeMap = {
  sm: 'size-2.5',
  md: 'size-3',
  lg: 'size-3',
  xl: 'size-3.5',
}

const roundedMap = {
  sm: 'rounded-xl',
  md: 'rounded-2xl',
  lg: 'rounded-2xl',
  xl: 'rounded-3xl',
}

export function Avatar({
  name,
  avatar,
  accent,
  size = 'lg',
  status,
  online,
  muted,
  selected,
  themeColor,
  onClick,
}: AvatarProps) {
  const initials = (avatar || name.slice(0, 2)).toUpperCase()
  const showStatus = status || online

  return (
    <div className="relative">
      <div
        className={clsx(
          'flex items-center justify-center font-bold uppercase text-white shadow-sm transition',
          sizeMap[size],
          textSizeMap[size],
          roundedMap[size],
          onClick && 'cursor-pointer',
          selected ? `ring-4 ring-[${accent}]` : 'ring-4 ring-transparent'
        )}
        style={{
          background: muted
            ? `linear-gradient(135deg, ${accent}, #94a3b8)`
            : themeColor
              ? `linear-gradient(135deg, ${accent}, ${themeColor})`
              : accent,
        }}
        onClick={onClick}
        role={onClick ? 'button' : undefined}
      >
        {initials}
      </div>
      {showStatus && (
        <span
          className={clsx(
            'absolute -bottom-0.5 -right-0.5 flex items-center justify-center rounded-full border-2 border-white',
            statusSizeMap[size],
            status === 'online' || online
              ? 'bg-emerald-500'
              : status === 'busy'
                ? 'bg-amber-500'
                : 'bg-slate-400'
          )}
        />
      )}
    </div>
  )
}
