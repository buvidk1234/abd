import { createContext } from 'react'
import type { useMessage } from '../hooks/useMessage'

type MessageContextValue = ReturnType<typeof useMessage>

export const MessageContext = createContext<MessageContextValue | null>(null)
