import type {
  Friend,
  FriendRequest,
  BlackUser,
  SearchFriendParams,
  AddFriendData,
  HandleFriendApplyReq,
} from './model'
import type {
  FriendInfoDTO,
  FriendRequestInfoDTO,
  BlackInfoDTO,
  AddBlackReq,
  RemoveBlackReq,
  SearchFriendReq,
} from './api'
import {
  getFriendsApi,
  getFriendInfoApi,
  deleteFriendApi,
  applyToAddFriendApi,
  respondFriendApplyApi,
  getFriendApplyListApi,
  addBlackApi,
  removeBlackApi,
  getBlacksApi,
  searchFriendApi,
} from './api'
import { toUser } from '../user/repo'
import type { User } from '../user'
import type { PagedParams, PagedR } from '../types'

const toFriend = (dto: FriendInfoDTO): Friend => ({
  friendUser: toUser(dto.friendUser),
  remark: dto.remark,
  isPinned: dto.isPinned,
  addSource: dto.addSource,
  createdAt: dto.createdAt,
})

const toSearchFriendReq = (params: SearchFriendParams): SearchFriendReq => ({
  id: params.id,
})

const toFriendRequest = (dto: FriendRequestInfoDTO): FriendRequest => ({
  id: String(dto.id),
  fromUser: toUser(dto.fromUser),
  toUser: toUser(dto.toUser),
  message: dto.reqMsg || '',
  handleResult: dto.handleResult,
  handleMessage: dto.handleMsg,
  createdAt: dto.createdAt,
  handledAt: dto.handledAt,
})

const toBlackUser = (dto: BlackInfoDTO): BlackUser => ({
  id: dto.blockUser.user_id,
  userId: dto.blockUser.user_id,
  name: dto.blockUser.nickname || dto.blockUser.username || '未命名',
  avatar: dto.blockUser.avatar_url,
  addSource: dto.addSource,
  createdAt: dto.createdAt,
})

export async function listFriends(params?: PagedParams): Promise<PagedR<Friend>> {
  const { data } = await getFriendsApi(params)
  return {
    data: data.data.map(toFriend),
    total: data.total,
    page: data.page,
    pageSize: data.pageSize,
  }
}

// 搜索好友
export async function searchFriend(params: SearchFriendParams): Promise<User> {
  const { data } = await searchFriendApi(toSearchFriendReq(params))
  return toUser(data)
}

export async function getFriendById(friendId: string): Promise<Friend> {
  const { data } = await getFriendInfoApi(friendId)
  return toFriend(data)
}

export async function removeFriend(friendId: string): Promise<void> {
  await deleteFriendApi(friendId)
}

// 发出添加好友请求
export async function applyToAddFriend(req: AddFriendData): Promise<void> {
  await applyToAddFriendApi(req)
}

export async function respondFriendApply(req: HandleFriendApplyReq): Promise<void> {
  await respondFriendApplyApi(req)
}

export async function listFriendApplies(req: PagedParams): Promise<PagedR<FriendRequest>> {
  const { data } = await getFriendApplyListApi(req)
  return {
    data: data.data.map(toFriendRequest),
    total: data.total,
    page: data.page,
    pageSize: data.pageSize,
  }
}

export async function addToBlacklist(req: AddBlackReq): Promise<void> {
  await addBlackApi(req)
}

export async function removeFromBlacklist(req: RemoveBlackReq): Promise<void> {
  await removeBlackApi(req)
}

export async function listBlackUsers(params?: PagedParams): Promise<PagedR<BlackUser>> {
  const { data } = await getBlacksApi(params)
  return {
    data: data.blacks.map(toBlackUser),
    total: data.total,
    page: data.page,
    pageSize: data.pageSize,
  }
}
