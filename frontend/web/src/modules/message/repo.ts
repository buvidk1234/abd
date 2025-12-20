/**
 * Repository Layer: Message data transformation
 * - Calls API functions
 * - Converts DTO -> Domain models
 * - Handles field compatibility, defaults
 * - Key value: "When API changes, only modify this layer"
 */

import type { Message, SendMessageData, ConversationMessages } from './model'
import type { MessageDTO, SendMessageReq, PullConvListParams, PullSpecifiedConvParams } from './api'
import { sendMessageApi, pullConvListApi, pullConversationApi } from './api'

// ============ DTO to Domain Converters ============

const toMessage = (dto: MessageDTO): Message => ({
  id: String(dto.ID),
  conversationId: dto.ConversationID,
  seq: dto.Seq,
  senderId: dto.SenderID,
  clientMsgId: dto.ClientMsgID,
  msgType: dto.MsgType,
  content: dto.Content,
  refMsgId: dto.RefMsgID ? String(dto.RefMsgID) : undefined,
  status: dto.Status,
  sendTime: dto.SendTime,
  createTime: dto.CreateTime,
  convType: dto.ConvType,
  targetId: dto.TargetID,
})

// ============ Domain to DTO Converters ============

const toSendMessageReq = (data: SendMessageData): SendMessageReq => ({
  conv_type: data.convType,
  target_id: data.targetId,
  msg_type: data.msgType,
  content: data.content,
})

// ============ Repository Functions ============

export async function sendMessage(data: SendMessageData): Promise<void> {
  await sendMessageApi(toSendMessageReq(data))
}

export async function pullAllConversations(userSeq?: number): Promise<ConversationMessages> {
  const { data } = await pullConvListApi({ user_seq: userSeq })
  const result: ConversationMessages = {}

  for (const [convId, messages] of Object.entries(data.pull_msgs)) {
    result[convId] = messages.map(toMessage)
  }

  return result
}

export async function pullConversationMessages(
  convId: string,
  convSeq?: number
): Promise<Message[]> {
  const { data } = await pullConversationApi(convId, { conv_seq: convSeq })
  return data.messages.map(toMessage)
}
