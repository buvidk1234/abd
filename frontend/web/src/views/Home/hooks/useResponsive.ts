import { useEffect } from 'react'
import { useImmer } from 'use-immer'

export function useResponsive(breakpoint: number = 1024) {
  const [isMobile, setIsMobile] = useImmer(false)

  useEffect(() => {
    const updateIsMobile = () => {
      if (typeof window === 'undefined') return
      setIsMobile(() => window.innerWidth < breakpoint)
    }

    updateIsMobile()
    window.addEventListener('resize', updateIsMobile)
    return () => window.removeEventListener('resize', updateIsMobile)
  }, [breakpoint, setIsMobile])

  return isMobile
}
