import type { User } from '../user'
import type { ApplyToAddFriendReq, RespondFriendApplyReq } from './api'

export type AddFriendData = ApplyToAddFriendReq

export type HandleFriendApplyReq = RespondFriendApplyReq

export interface SearchFriendParams {
  id: string
}

export interface Friend {
  friendUser: User
  remark: string
  isPinned: boolean
  addSource: number
  createdAt: number
}

export interface FriendRequest {
  id: string
  fromUser: User
  toUser: User
  message: string
  handleResult: number
  handleMessage?: string
  createdAt: string
  handledAt?: string
}

export interface BlackUser {
  id: string
  userId: string
  name: string
  avatar?: string
  addSource: number
  createdAt: number
}
