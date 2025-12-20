/**
 * Domain Layer: Group business models
 * - Page-independent types
 * - Pure functions and rules
 * - No React, no UI dependencies, no HTTP
 */

export interface Group {
  id: string
  name: string
  avatar: string
  createdAt: string
  status: number
  creatorUserId: string
  groupType: number
  needVerification: number
  lookMemberInfo: number
  applyMemberFriend: number
  ex?: string
}

export interface GroupMember {
  id: string
  groupId: string
  userId: string
  nickname: string
  avatar: string
  roleLevel: number
  joinedAt: string
}

export interface CreateGroupData {
  name: string
  avatar: string
  ex?: string
  groupType?: number
  needVerification?: number
  lookMemberInfo?: number
  applyMemberFriend?: number
}

export interface UpdateGroupData {
  name?: string
  avatar?: string
  notification?: string
  introduction?: string
  needVerification?: number
  lookMemberInfo?: number
  applyMemberFriend?: number
}

export interface UpdateGroupMemberData {
  nickname?: string
  avatar?: string
  roleLevel?: number
}

export interface JoinGroupData {
  message?: string
  joinSource?: number
  inviterUserId?: string
}

export interface InviteToGroupData {
  inviteeUserId: string
}

export type JoinStatus = 'pending' | 'joined'
