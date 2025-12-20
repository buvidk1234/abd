import { http } from '@/lib/http'
import type { PagedParams, PagedR } from '@/modules/types'
import type { UserInfoDTO } from '../user/api'

export interface FriendInfoDTO {
  remark: string
  createdAt: number
  friendUser: UserInfoDTO
  addSource: number
  isPinned: boolean
}

export interface FriendRequestInfoDTO {
  id: number
  fromUser: UserInfoDTO
  toUser: UserInfoDTO
  handleResult: number
  reqMsg: string
  createdAt: string
  updatedAt: string
  handleMsg: string
  handledAt: string
}

export interface BlackInfoDTO {
  ownerUserID: number
  blockUser: UserInfoDTO
  addSource: number
  createdAt: number
}

export interface ApplyToAddFriendReq {
  toUserID: string
  message?: string
}

export interface RespondFriendApplyReq {
  id: string
  handleResult: number
  handleMsg?: string
}

export interface AddBlackReq {
  blockUserID: string
  addSource?: number | string
}

export interface RemoveBlackReq {
  blockUserID: string
}

export interface GetPaginationBlacksResp {
  blacks: BlackInfoDTO[]
  total: number
  page: number
  pageSize: number
}

export interface SearchFriendReq {
  id: string
}

export const getFriendsApi = (params?: PagedParams) =>
  http.get<PagedR<FriendInfoDTO>>('/friend/', { params })

export const searchFriendApi = (params: SearchFriendReq) =>
  http.get<UserInfoDTO>(`/friend/search`, { params })

// TODO: Fix
export const getFriendInfoApi = (friendId: string) => http.get<FriendInfoDTO>(`/friend/${friendId}`)

export const deleteFriendApi = (friendId: string) => http.delete<void>(`/friend/${friendId}`)

export const applyToAddFriendApi = (data: ApplyToAddFriendReq) =>
  http.post<void>('/friend/add', data)

export const respondFriendApplyApi = (data: RespondFriendApplyReq) =>
  http.post<void>('/friend/add-response', data)

export const getFriendApplyListApi = (params: PagedParams) =>
  http.get<PagedR<FriendRequestInfoDTO>>('/friend/apply', { params })

export const addBlackApi = (data: AddBlackReq) => http.post<void>('/friend/add_black', data)

export const removeBlackApi = (data: RemoveBlackReq) =>
  http.post<void>('/friend/remove_black', data)

export const getBlacksApi = (params?: PagedParams) =>
  http.get<GetPaginationBlacksResp>('/friend/black', { params })
