import { useEffect, type RefObject } from 'react'

export function useClickOutside(ref: RefObject<HTMLElement | null>, handler: () => void, enabled: boolean = true) {
  useEffect(() => {
    if (!enabled) return

    const handleClickAway = (event: MouseEvent) => {
      const target = event.target as Node
      if (ref.current && !ref.current.contains(target)) {
        handler()
      }
    }

    document.addEventListener('mousedown', handleClickAway)
    return () => document.removeEventListener('mousedown', handleClickAway)
  }, [ref, handler, enabled])
}
