import { useContext } from 'react'
import { WSContext } from '../providers/WSContext'

export function useWS() {
  const ctx = useContext(WSContext)
  if (!ctx) throw new Error('useWS must be used within <WebSocketProvider />')
  return ctx
}
