import type { FriendInfo, FriendRequestInfo, PagedResp } from '@/types/backend'
import { http } from '@/http/http'

export interface AddFriendPayload {
  fromUserID: string
  toUserID: string
  message?: string
}

export interface RespondFriendPayload {
  id: number
  handlerUserID: string
  handleResult: number // 1同意 2拒绝
  handleMsg?: string
}

export interface FriendListQuery {
  userID: string
  page?: number
  pageSize?: number
}

export function applyToAddFriend(payload: AddFriendPayload) {
  return http.post<{ msg?: string }>('/friend/add_friend', payload)
}

export function respondFriendApply(payload: RespondFriendPayload) {
  return http.post<{ msg?: string }>('/friend/add_friend_response', payload)
}

export function deleteFriend(payload: { ownerUserID: string, friendUserID: string }) {
  return http.post<{ msg?: string }>('/friend/delete_friend', payload)
}

export async function getFriendList(query: FriendListQuery) {
  const res = await http.post<{
    data: {
      friends: FriendInfo[]
      total: number
      page: number
      pageSize: number
    }
  }>('/friend/get_friend_list', {
    page: query.page || 1,
    pageSize: query.pageSize || 50,
    userID: query.userID,
  })
  // 后端返回 { data: {...} }
  return (res as any)?.data || (res as any)
}

export async function getIncomingFriendRequests(query: { toUserID: string, page?: number, pageSize?: number }) {
  const res = await http.post<{ data: PagedResp<FriendRequestInfo> }>('/friend/get_friend_apply_list', {
    page: query.page || 1,
    pageSize: query.pageSize || 50,
    toUserID: query.toUserID,
  })
  return (res as any)?.data || (res as any)
}

export async function getOutgoingFriendRequests(query: { fromUserID: string, page?: number, pageSize?: number }) {
  const res = await http.post<{ data: PagedResp<FriendRequestInfo> }>('/friend/get_self_friend_apply_list', {
    page: query.page || 1,
    pageSize: query.pageSize || 50,
    fromUserID: query.fromUserID,
  })
  return (res as any)?.data || (res as any)
}
