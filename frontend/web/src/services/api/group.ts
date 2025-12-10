import { http } from '@/lib/http'

export interface GroupInfo {
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

export interface GroupMemberInfo {
  id: number
  groupID: string
  userID: string
  nickname: string
  avatarURL: string
  roleLevel: number
  joinedAt: string
}

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
  groupInfos: GroupInfo[]
}

export interface GetGroupMemberListResp {
  memberList: GroupMemberInfo[]
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

export const createGroup = (data: CreateGroupReq) => {
  return http.post<CreateGroupResp>('/groups', data)
}

export const getGroupsInfo = (params?: GetGroupsInfoParams) => {
  const query = params?.ids?.length ? { ids: params.ids.join(',') } : undefined
  return http.get<GetGroupsInfoResp>('/groups', { params: query })
}

export const getGroupMemberList = (groupID: string) => {
  return http.get<GetGroupMemberListResp>(`/groups/${groupID}/members`)
}

export const joinGroup = (groupID: string, data: JoinGroupReq) => {
  return http.post<JoinGroupResp>(`/groups/${groupID}/join`, data)
}

export const quitGroup = (groupID: string, userID: string) => {
  return http.delete<CommonMsgResp>(`/groups/${groupID}/members/${userID}`, { data: {} })
}

export const inviteUserToGroup = (groupID: string, data: InviteUserToGroupReq) => {
  return http.post<InviteUserToGroupResp>(`/groups/${groupID}/invitations`, data)
}

export const kickGroupMember = (groupID: string, userID: string) => {
  return http.delete<CommonMsgResp>(`/groups/${groupID}/members/${userID}/kick`, { data: {} })
}

export const dismissGroup = (groupID: string) => {
  return http.delete<CommonMsgResp>(`/groups/${groupID}`, { data: {} })
}

export const setGroupInfo = (groupID: string, data: SetGroupInfoReq) => {
  return http.post<CommonMsgResp>(`/groups/${groupID}`, data)
}

export const setGroupMemberInfo = (
  groupID: string,
  userID: string,
  data: SetGroupMemberInfoReq
) => {
  return http.post<CommonMsgResp>(`/groups/${groupID}/members/${userID}`, data)
}
