import React, { useMemo } from 'react'
import { useWebSocketCore } from '@/hooks/useWebSocket'
import { WSContext } from './WSContext'

export function WebSocketProvider(props: {
  token: string
  userId: number | string
  children: React.ReactNode
}) {
  const url = useMemo(() => {
    const base = import.meta.env.VITE_WEBSOCKET_URL
    return `${base}/ws?token=${encodeURIComponent(props.token)}&sendId=${encodeURIComponent(
      String(props.userId)
    )}`
  }, [props.token, props.userId])

  const value = useWebSocketCore({ url })

  return <WSContext.Provider value={value}>{props.children}</WSContext.Provider>
}
