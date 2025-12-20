/**
 * Infrastructure Layer: Message API calls
 * - Only handles "how to request"
 * - Returns backend DTOs
 * - No business logic conversion, no UI logic
 */

import { http } from '@/lib/http'

// ============ DTOs (Backend Response Types) ============

export interface MessageDTO {
  ID: number
  ConversationID: string
  Seq: number
  SenderID: string
  ClientMsgID?: string
  MsgType: number
  Content: string
  RefMsgID?: number
  Status?: number
  SendTime?: number
  CreateTime?: number
  ConvType?: number
  TargetID?: string
}

// ============ Request/Response Types ============

export interface SendMessageReq {
  conv_type: number
  target_id: string
  msg_type: number
  content: string
}

export interface PullConvListParams {
  user_seq?: number
}

export interface PullConvListResp {
  pull_msgs: Record<string, MessageDTO[]>
}

export interface PullSpecifiedConvParams {
  conv_seq?: number
}

export interface PullSpecifiedConvResp {
  messages: MessageDTO[]
}

// ============ API Functions ============

export const sendMessageApi = (data: SendMessageReq) => http.post<void>('/msg/send', data)

export const pullConvListApi = (params?: PullConvListParams) =>
  http.get<PullConvListResp>('/msg/pull', { params })

export const pullConversationApi = (convID: string, params?: PullSpecifiedConvParams) =>
  http.get<PullSpecifiedConvResp>(`/msg/pull/${convID}`, {
    params: { conv_id: convID, ...params },
  })
