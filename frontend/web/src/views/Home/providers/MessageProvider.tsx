import React from 'react'
import { useMessage } from '../hooks/useMessage'
import { MessageContext } from './MessageContext'

export function MessageProvider(props: { children: React.ReactNode }) {
  const messageValue = useMessage()

  return <MessageContext.Provider value={messageValue}>{props.children}</MessageContext.Provider>
}
