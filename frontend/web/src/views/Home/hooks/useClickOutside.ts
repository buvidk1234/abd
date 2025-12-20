import { useEffect, type RefObject } from 'react'

/**
 * 点击外部区域关闭组件
 * @param ref 组件引用
 * @param handler 关闭回调
 * @param enabled 是否启用
 */
export function useClickOutside(
  ref: RefObject<HTMLElement | null>,
  handler: () => void,
  enabled: boolean = true
) {
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
