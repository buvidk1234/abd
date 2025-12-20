/**
 * Domain Layer: Message business models
 * - Page-independent types
 * - Pure functions and rules
 * - No React, no UI dependencies, no HTTP
 */

export enum ConvType {
  SingleChat = 1,
  GroupChat = 2,
}

export enum MessageType {
  Text = 1,
  Image = 2,
  Audio = 3,
  Video = 4,
  File = 5,
}

export enum MessageStatus {
  Sending = 1,
  Sent = 2,
  Failed = 3,
  Read = 4,
}

export interface Message {
  id: string
  conversationId: string
  seq: number
  senderId: string
  clientMsgId?: string
  msgType: MessageType
  content: string
  refMsgId?: string
  status?: MessageStatus
  sendTime?: number
  createTime?: number
  convType?: ConvType
  targetId?: string
}

export interface SendMessageData {
  convType: ConvType
  targetId: string
  msgType: MessageType
  content: string
}

export interface ConversationMessages {
  [conversationId: string]: Message[]
}
