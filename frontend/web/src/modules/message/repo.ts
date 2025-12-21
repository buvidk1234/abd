/**
 * Repository Layer: Message data transformation
 * - Calls API functions
 * - Converts DTO -> Domain models
 * - Handles field compatibility, defaults
 * - Key value: "When API changes, only modify this layer"
 */

import type { Message, SendMessageData, ConversationMessages } from './model'
import type { MessageDTO, SendMessageReq } from './api'
import { sendMessageApi, pullConvListApi, pullConversationApi } from './api'

// ============ DTO to Domain Converters ============

const toMessage = (dto: MessageDTO): Message => ({
  id: String(dto.ID),
  conversation_id: dto.ConversationID,
  seq: dto.Seq,
  sender_id: dto.SenderID,
  client_msg_id: dto.ClientMsgID,
  msg_type: dto.MsgType,
  content: dto.Content,
  ref_msg_id: dto.RefMsgID ? String(dto.RefMsgID) : undefined,
  status: dto.Status,
  send_time: dto.SendTime,
  create_time: dto.CreateTime,
  conv_type: dto.ConvType,
  target_id: dto.TargetID,
})

// ============ Domain to DTO Converters ============

const toSendMessageReq = (data: SendMessageData): SendMessageReq => ({
  conv_type: data.conv_type,
  target_id: data.targetId,
  msg_type: data.msg_type,
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
