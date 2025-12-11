import { useEffect, useState } from 'react'

export function useResponsive(breakpoint: number = 1024) {
  const [isMobile, setIsMobile] = useState(false)

  useEffect(() => {
    const updateIsMobile = () => {
      if (typeof window === 'undefined') return
      setIsMobile(window.innerWidth < breakpoint)
    }

    updateIsMobile()
    window.addEventListener('resize', updateIsMobile)
    return () => window.removeEventListener('resize', updateIsMobile)
  }, [breakpoint])

  return isMobile
}
