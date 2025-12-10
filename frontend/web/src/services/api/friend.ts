import { http } from '@/lib/http'
import type { UserInfo } from './user'

export interface ApplyToAddFriendReq {
  fromUserID: string
  toUserID: string
  message?: string
}

export interface RespondFriendApplyReq {
  id: string
  handlerUserID: string
  handleResult: number | string
  handleMsg?: string
}

export interface GetFriendListParams {
  page?: number
  pageSize?: number
}

export interface FriendInfo {
  ownerUserID: number
  remark: string
  createdAt: number
  friendUser: UserInfo
  addSource: number
  isPinned: boolean
}

export interface GetFriendListResp {
  friends: FriendInfo[]
  total: number
  page: number
  pageSize: number
}

export interface FriendRequestInfo {
  id: number
  fromUser: UserInfo
  toUser: UserInfo
  handleResult: number
  reqMsg: string
  createdAt: string
  updatedAt: string
  handleMsg: string
  handledAt: string
}

export interface GetFriendApplyListReq {
  page?: number
  pageSize?: number
  toUserID: string
}

export interface GetFriendApplyListResp {
  list: FriendRequestInfo[]
  total: number
  page: number
  pageSize: number
}

export interface GetSelfFriendApplyListReq {
  page?: number
  pageSize?: number
  fromUserID: string
}

export interface GetSelfFriendApplyListResp {
  list: FriendRequestInfo[]
  total: number
  page: number
  pageSize: number
}

export interface AddBlackReq {
  blockUserID: string
  addSource?: number | string
}

export interface RemoveBlackReq {
  blockUserID: string
}

export interface GetPaginationBlacksReq {
  page?: number
  pageSize?: number
}

export interface BlackInfo {
  ownerUserID: number
  blockUser: UserInfo
  addSource: number
  createdAt: number
}

export interface GetPaginationBlacksResp {
  blacks: BlackInfo[]
  total: number
  page: number
  pageSize: number
}

export const applyToAddFriend = (data: ApplyToAddFriendReq) => {
  return http.post<void>('/friend/add', data)
}

export const respondFriendApply = (data: RespondFriendApplyReq) => {
  return http.post<void>('/friend/add-response', data)
}

export const getFriendList = (params?: GetFriendListParams) => {
  return http.get<GetFriendListResp>('/friend', { params })
}

export const getFriendInfo = (friendId: string) => {
  return http.get<FriendInfo>(`/friend/${friendId}`)
}

export const deleteFriend = (friendId: string) => {
  return http.delete<void>(`/friend/${friendId}`)
}

export const getFriendApplyList = (data: GetFriendApplyListReq) => {
  return http.post<GetFriendApplyListResp>('/friend/get_friend_apply_list', data)
}

export const getSelfFriendApplyList = (data: GetSelfFriendApplyListReq) => {
  return http.post<GetSelfFriendApplyListResp>('/friend/get_self_friend_apply_list', data)
}

export const addBlack = (data: AddBlackReq) => {
  return http.post<void>('/friend/add_black', data)
}

export const removeBlack = (data: RemoveBlackReq) => {
  return http.post<void>('/friend/remove_black', data)
}

export const getBlacks = (params?: GetPaginationBlacksReq) => {
  return http.get<GetPaginationBlacksResp>('/friend/black', { params })
}
