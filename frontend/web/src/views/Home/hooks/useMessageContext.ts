import { useContext } from 'react'
import { MessageContext } from '../providers/MessageContext'

export function useMessageContext() {
  const ctx = useContext(MessageContext)
  if (!ctx) throw new Error('useMessageContext must be used within <MessageProvider />')
  return ctx
}
