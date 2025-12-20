/**
 * Group Module Entry Point
 * Export domain models and repository functions for use in UI layer
 */

// Domain models
export type {
  Group,
  GroupMember,
  CreateGroupData,
  UpdateGroupData,
  UpdateGroupMemberData,
  JoinGroupData,
  InviteToGroupData,
  JoinStatus,
} from './model'

// Repository functions
export {
  createGroup,
  getGroupsByIds,
  getGroupMembers,
  joinGroup,
  leaveGroup,
  inviteToGroup,
  kickMember,
  deleteGroup,
  updateGroupInfo,
  updateGroupMember,
} from './repo'
