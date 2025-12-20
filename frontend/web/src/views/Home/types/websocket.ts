type WSEventType = 1001 | 1002 | 4001 // 测试

export interface WSEvent<T = unknown> {
  req_identifier: WSEventType
  data: T
}
