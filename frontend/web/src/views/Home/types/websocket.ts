type WSEventType =
  | 1001 // 获取最新的seq null GetMaxSeqResp
  | 1002 // 拉取消息 PullMessageBySeqsReq PullMessageBySeqsResp
  | 1003 // 发送消息 SendMessageReq null
  | 2001 // 接收到消息  null Message
  | 4001 // 测试

// 发
export interface WSRequest<T = unknown> {
  req_identifier: WSEventType
  data: T
}

// 收
export interface WSResponse<T = unknown> {
  req_identifier: WSEventType
  code: number
  msg: string
  data: T
}
