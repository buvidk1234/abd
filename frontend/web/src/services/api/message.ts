import { http } from '@/lib/http'

export interface SendMessageReq {
  conv_type: number
  target_id: string
  msg_type: number
  content: string
}

export interface Message {
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

export interface PullConvListParams {
  user_seq?: number
}

export interface PullConvListResp {
  pull_msgs: Record<string, Message[]>
}

export interface PullSpecifiedConvParams {
  conv_seq?: number
}

export interface PullSpecifiedConvResp {
  messages: Message[]
}

export const sendMessage = (data: SendMessageReq) => {
  return http.post<void>('/msg/send', data)
}

export const pullConvList = (params?: PullConvListParams) => {
  return http.get<PullConvListResp>('/msg/pull', { params })
}

export const pullConversation = (convID: string, params?: PullSpecifiedConvParams) => {
  return http.get<PullSpecifiedConvResp>(`/msg/pull/${convID}`, {
    params: { conv_id: convID, ...params },
  })
}
