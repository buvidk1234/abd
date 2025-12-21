/**
 * Domain Layer: Message business models
 * - Page-independent types
 * - Pure functions and rules
 * - No React, no UI dependencies, no HTTP
 */

export interface Message {
  id: string
  conversation_id: string
  seq: number
  sender_id: string
  client_msg_id?: string
  msg_type: number
  content: string
  ref_msg_id?: string
  status?: number
  send_time?: number
  create_time?: number
  conv_type?: number
  target_id?: string
}

export interface SendMessageData {
  conv_type: number
  targetId: string
  msg_type: number
  content: string
}

export interface ConversationMessages {
  [conversationId: string]: Message[]
}
