import { createContext } from 'react'
import { useWebSocketCore } from '@/hooks/useWebSocket'

type WSContextValue = ReturnType<typeof useWebSocketCore>

export const WSContext = createContext<WSContextValue | null>(null)
