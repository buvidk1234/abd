import type { BlackInfo, FriendInfo, FriendRequestInfo, PagedResp, UserProfile } from '@/types/backend'
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
  page?: number
  pageSize?: number
}

const toUserProfile = (raw?: any): UserProfile => ({
  userID: raw?.user_id ? String(raw.user_id) : (raw?.userID || raw?.userId || ''),
  nickname: raw?.nickname,
  avatarURL: raw?.avatar_url || raw?.avatarURL,
  gender: raw?.gender,
  signature: raw?.signature,
  username: raw?.username,
  phone: raw?.phone,
  email: raw?.email,
})

const toFriendInfo = (raw: any): FriendInfo => ({
  ownerUserID: String(raw?.ownerUserID ?? raw?.owner_user_id ?? ''),
  remark: raw?.remark || '',
  createdAt: Number(raw?.createdAt ?? raw?.created_at ?? 0),
  friendUser: toUserProfile(raw?.friendUser || raw?.friend_user),
  addSource: Number(raw?.addSource ?? raw?.add_source ?? 0),
  isPinned: !!raw?.isPinned,
})

const toFriendRequest = (raw: any): FriendRequestInfo => ({
  id: Number(raw?.id ?? 0),
  fromUser: toUserProfile(raw?.fromUser || raw?.from_user),
  toUser: toUserProfile(raw?.toUser || raw?.to_user),
  handleResult: Number(raw?.handleResult ?? raw?.handle_result ?? 0),
  reqMsg: raw?.reqMsg || raw?.req_msg || '',
  createdAt: raw?.createdAt ? String(raw.createdAt) : '',
  updatedAt: raw?.updatedAt ? String(raw.updatedAt) : '',
  handleMsg: raw?.handleMsg || raw?.handle_msg || '',
  handledAt: raw?.handledAt ? String(raw.handledAt) : '',
})

const toBlackInfo = (raw: any): BlackInfo => ({
  ownerUserID: String(raw?.ownerUserID ?? raw?.owner_user_id ?? ''),
  blockUser: toUserProfile(raw?.blockUser || raw?.block_user),
  addSource: Number(raw?.addSource ?? raw?.add_source ?? 0),
  createdAt: Number(raw?.createdAt ?? raw?.created_at ?? 0),
})

export function applyToAddFriend(payload: AddFriendPayload) {
  return http.post('/friend/add', {
    fromUserID: payload.fromUserID,
    toUserID: payload.toUserID,
    message: payload.message,
  })
}

export function respondFriendApply(payload: RespondFriendPayload) {
  return http.post('/friend/add-response', payload)
}

export function deleteFriend(friendUserID: string) {
  return http({
    url: `/friend/${friendUserID}`,
    method: 'DELETE',
  })
}

export async function getFriendList(query: FriendListQuery = {}) {
  const res = await http.get<{
    friends: any[]
    total: number
    page: number
    pageSize: number
  }>('/friend', {
    page: query.page || 1,
    pageSize: query.pageSize || 10,
  })
  return {
    friends: (res?.friends || []).map(toFriendInfo),
    total: (res as any)?.total ?? 0,
    page: (res as any)?.page ?? query.page ?? 1,
    pageSize: (res as any)?.pageSize ?? query.pageSize ?? 10,
  }
}

export async function getFriendDetail(friendUserID: string) {
  const res = await http.get<any>(`/friend/${friendUserID}`)
  return toFriendInfo(res)
}

export async function getIncomingFriendRequests(query: { toUserID: string, page?: number, pageSize?: number }) {
  const res = await http.post<PagedResp<any>>('/friend/get_friend_apply_list', {
    page: query.page || 1,
    pageSize: query.pageSize || 20,
    toUserID: query.toUserID,
  })
  const payload = (res as any)?.data || res
  return {
    list: (payload?.list || []).map(toFriendRequest),
    total: payload?.total ?? 0,
    page: payload?.page ?? query.page ?? 1,
    pageSize: payload?.pageSize ?? query.pageSize ?? 20,
  } as PagedResp<FriendRequestInfo>
}

export async function getOutgoingFriendRequests(query: { fromUserID: string, page?: number, pageSize?: number }) {
  const res = await http.post<PagedResp<any>>('/friend/get_self_friend_apply_list', {
    page: query.page || 1,
    pageSize: query.pageSize || 20,
    fromUserID: query.fromUserID,
  })
  const payload = (res as any)?.data || res
  return {
    list: (payload?.list || []).map(toFriendRequest),
    total: payload?.total ?? 0,
    page: payload?.page ?? query.page ?? 1,
    pageSize: payload?.pageSize ?? query.pageSize ?? 20,
  } as PagedResp<FriendRequestInfo>
}

export function addBlack(payload: { blockUserID: string, addSource?: number }) {
  return http.post('/friend/add_black', payload)
}

export function removeBlack(blockUserID: string) {
  return http.post('/friend/remove_black', { blockUserID })
}

export async function getBlackList(query: { page?: number, pageSize?: number } = {}) {
  const res = await http.get<{
    blacks: any[]
    total: number
    page: number
    pageSize: number
  }>('/friend/black', {
    page: query.page || 1,
    pageSize: query.pageSize || 10,
  })
  return {
    blacks: (res as any)?.blacks?.map(toBlackInfo) || [],
    total: (res as any)?.total ?? 0,
    page: (res as any)?.page ?? query.page ?? 1,
    pageSize: (res as any)?.pageSize ?? query.pageSize ?? 10,
  }
}
