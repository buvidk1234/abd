/**
 * Infrastructure Layer: Group API calls
 * - Only handles "how to request"
 * - Returns backend DTOs
 * - No business logic conversion, no UI logic
 */

import { http } from '@/lib/http'

// ============ DTOs (Backend Response Types) ============

export interface GroupInfoDTO {
  id: number
  groupName: string
  avatarURL: string
  createdAt: string
  ex: string
  status: number
  creatorUserID: string
  groupType: number
  needVerification: number
  lookMemberInfo: number
  applyMemberFriend: number
}

export interface GroupMemberInfoDTO {
  id: number
  groupID: string
  userID: string
  nickname: string
  avatarURL: string
  roleLevel: number
  joinedAt: string
}

// ============ Request/Response Types ============

export interface CreateGroupReq {
  groupName: string
  avatarURL: string
  ex?: string
  groupType?: number
  needVerification?: number
  lookMemberInfo?: number
  applyMemberFriend?: number
}

export interface CreateGroupResp {
  groupID: string
}

export interface GetGroupsInfoParams {
  ids?: string[]
}

export interface GetGroupsInfoResp {
  groupInfos: GroupInfoDTO[]
}

export interface GetGroupMemberListResp {
  memberList: GroupMemberInfoDTO[]
}

export interface JoinGroupReq {
  reqMsg?: string
  joinSource?: number
  inviterUserID?: string
}

export interface JoinGroupResp {
  status: 'pending' | 'joined' | string
}

export interface InviteUserToGroupReq {
  inviteeUserID: string
}

export interface InviteUserToGroupResp {
  status: 'pending' | 'joined' | string
}

export interface CommonMsgResp {
  msg: string
}

export interface SetGroupInfoReq {
  groupName?: string
  avatarURL?: string
  notification?: string
  introduction?: string
  needVerification?: number
  lookMemberInfo?: number
  applyMemberFriend?: number
}

export interface SetGroupMemberInfoReq {
  nickname?: string
  avatarURL?: string
  roleLevel?: number
}

// ============ API Functions ============

export const createGroupApi = (data: CreateGroupReq) => http.post<CreateGroupResp>('/groups', data)

export const getGroupsInfoApi = (params?: GetGroupsInfoParams) => {
  const query = params?.ids?.length ? { ids: params.ids.join(',') } : undefined
  return http.get<GetGroupsInfoResp>('/groups', { params: query })
}

export const getGroupMemberListApi = (groupID: string) =>
  http.get<GetGroupMemberListResp>(`/groups/${groupID}/members`)

export const joinGroupApi = (groupID: string, data: JoinGroupReq) =>
  http.post<JoinGroupResp>(`/groups/${groupID}/join`, data)

export const quitGroupApi = (groupID: string, userID: string) =>
  http.delete<CommonMsgResp>(`/groups/${groupID}/members/${userID}`, { data: {} })

export const inviteUserToGroupApi = (groupID: string, data: InviteUserToGroupReq) =>
  http.post<InviteUserToGroupResp>(`/groups/${groupID}/invitations`, data)

export const kickGroupMemberApi = (groupID: string, userID: string) =>
  http.delete<CommonMsgResp>(`/groups/${groupID}/members/${userID}/kick`, { data: {} })

export const dismissGroupApi = (groupID: string) =>
  http.delete<CommonMsgResp>(`/groups/${groupID}`, { data: {} })

export const setGroupInfoApi = (groupID: string, data: SetGroupInfoReq) =>
  http.post<CommonMsgResp>(`/groups/${groupID}`, data)

export const setGroupMemberInfoApi = (
  groupID: string,
  userID: string,
  data: SetGroupMemberInfoReq
) => http.post<CommonMsgResp>(`/groups/${groupID}/members/${userID}`, data)
