/**
 * Repository Layer: Group data transformation
 * - Calls API functions
 * - Converts DTO -> Domain models
 * - Handles field compatibility, defaults
 * - Key value: "When API changes, only modify this layer"
 */

import type {
  Group,
  GroupMember,
  CreateGroupData,
  UpdateGroupData,
  UpdateGroupMemberData,
  JoinGroupData,
  InviteToGroupData,
  JoinStatus,
} from './model'
import type {
  GroupInfoDTO,
  GroupMemberInfoDTO,
  CreateGroupReq,
  GetGroupsInfoParams,
  JoinGroupReq,
  InviteUserToGroupReq,
  SetGroupInfoReq,
  SetGroupMemberInfoReq,
} from './api'
import {
  createGroupApi,
  getGroupsInfoApi,
  getGroupMemberListApi,
  joinGroupApi,
  quitGroupApi,
  inviteUserToGroupApi,
  kickGroupMemberApi,
  dismissGroupApi,
  setGroupInfoApi,
  setGroupMemberInfoApi,
} from './api'

// ============ DTO to Domain Converters ============

const toGroup = (dto: GroupInfoDTO): Group => ({
  id: String(dto.id),
  name: dto.groupName || '未命名群组',
  avatar: dto.avatarURL,
  createdAt: dto.createdAt,
  status: dto.status,
  creatorUserId: dto.creatorUserID,
  groupType: dto.groupType,
  needVerification: dto.needVerification,
  lookMemberInfo: dto.lookMemberInfo,
  applyMemberFriend: dto.applyMemberFriend,
  ex: dto.ex,
})

const toGroupMember = (dto: GroupMemberInfoDTO): GroupMember => ({
  id: String(dto.id),
  groupId: dto.groupID,
  userId: dto.userID,
  nickname: dto.nickname || '未命名成员',
  avatar: dto.avatarURL,
  roleLevel: dto.roleLevel,
  joinedAt: dto.joinedAt,
})

// ============ Domain to DTO Converters ============

const toCreateGroupReq = (data: CreateGroupData): CreateGroupReq => ({
  groupName: data.name,
  avatarURL: data.avatar,
  ex: data.ex,
  groupType: data.groupType,
  needVerification: data.needVerification,
  lookMemberInfo: data.lookMemberInfo,
  applyMemberFriend: data.applyMemberFriend,
})

const toSetGroupInfoReq = (data: UpdateGroupData): SetGroupInfoReq => ({
  groupName: data.name,
  avatarURL: data.avatar,
  notification: data.notification,
  introduction: data.introduction,
  needVerification: data.needVerification,
  lookMemberInfo: data.lookMemberInfo,
  applyMemberFriend: data.applyMemberFriend,
})

const toSetGroupMemberInfoReq = (data: UpdateGroupMemberData): SetGroupMemberInfoReq => ({
  nickname: data.nickname,
  avatarURL: data.avatar,
  roleLevel: data.roleLevel,
})

const toJoinGroupReq = (data: JoinGroupData): JoinGroupReq => ({
  reqMsg: data.message,
  joinSource: data.joinSource,
  inviterUserID: data.inviterUserId,
})

const toInviteUserToGroupReq = (data: InviteToGroupData): InviteUserToGroupReq => ({
  inviteeUserID: data.inviteeUserId,
})

// ============ Repository Functions ============

export async function createGroup(data: CreateGroupData): Promise<string> {
  const { data: result } = await createGroupApi(toCreateGroupReq(data))
  return result.groupID
}

export async function getGroupsByIds(ids?: string[]): Promise<Group[]> {
  const { data } = await getGroupsInfoApi({ ids })
  return data.groupInfos.map(toGroup)
}

export async function getGroupMembers(groupId: string): Promise<GroupMember[]> {
  const { data } = await getGroupMemberListApi(groupId)
  return data.memberList.map(toGroupMember)
}

export async function joinGroup(groupId: string, data: JoinGroupData): Promise<JoinStatus> {
  const { data: result } = await joinGroupApi(groupId, toJoinGroupReq(data))
  return result.status as JoinStatus
}

export async function leaveGroup(groupId: string, userId: string): Promise<string> {
  const { data } = await quitGroupApi(groupId, userId)
  return data.msg
}

export async function inviteToGroup(groupId: string, data: InviteToGroupData): Promise<JoinStatus> {
  const { data: result } = await inviteUserToGroupApi(groupId, toInviteUserToGroupReq(data))
  return result.status as JoinStatus
}

export async function kickMember(groupId: string, userId: string): Promise<string> {
  const { data } = await kickGroupMemberApi(groupId, userId)
  return data.msg
}

export async function deleteGroup(groupId: string): Promise<string> {
  const { data } = await dismissGroupApi(groupId)
  return data.msg
}

export async function updateGroupInfo(groupId: string, data: UpdateGroupData): Promise<string> {
  const { data: result } = await setGroupInfoApi(groupId, toSetGroupInfoReq(data))
  return result.msg
}

export async function updateGroupMember(
  groupId: string,
  userId: string,
  data: UpdateGroupMemberData
): Promise<string> {
  const { data: result } = await setGroupMemberInfoApi(
    groupId,
    userId,
    toSetGroupMemberInfoReq(data)
  )
  return result.msg
}
