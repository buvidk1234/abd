/**
 * Repository Layer: User data transformation
 * - Calls API functions
 * - Converts DTO -> Domain models
 * - Handles field compatibility, defaults
 * - Key value: "When API changes, only modify this layer"
 */

import type { User, LoginCredentials, RegisterData, UpdateUserProfile } from './model'
import type { UserInfoDTO, UserLoginReq, UserRegisterReq, UpdateUserInfoReq } from './api'
import { loginApi, registerApi, updateUserInfoApi, getUserInfoApi } from './api'

export const toUser = (dto: UserInfoDTO): User => ({
  id: dto.user_id,
  username: dto.username,
  nickname: dto.nickname || dto.username || '未命名',
  avatar: dto.avatar_url,
  gender: dto.gender,
  signature: dto.signature,
})

const toLoginReq = (credentials: LoginCredentials): UserLoginReq => ({
  username: credentials.username,
  password: credentials.password,
})

const toRegisterReq = (data: RegisterData): UserRegisterReq => ({
  username: data.username,
  password: data.password,
  phone: data.phone,
  email: data.email,
  nickname: data.nickname,
  avatar_url: data.avatarUrl,
})

const toUpdateUserInfoReq = (profile: UpdateUserProfile): UpdateUserInfoReq => ({
  nickname: profile.nickname,
  avatar_url: profile.avatarUrl,
  gender: profile.gender,
  signature: profile.signature,
  birth: profile.birth,
  phone: profile.phone,
  email: profile.email,
  ex: profile.ex,
})

export async function login(credentials: LoginCredentials): Promise<string> {
  const { data } = await loginApi(toLoginReq(credentials))
  return data
}

export async function register(registerData: RegisterData): Promise<string> {
  const { data } = await registerApi(toRegisterReq(registerData))
  return data
}

export async function updateProfile(profile: UpdateUserProfile): Promise<void> {
  await updateUserInfoApi(toUpdateUserInfoReq(profile))
}

export async function getCurrentUser(): Promise<User> {
  const { data } = await getUserInfoApi()
  return toUser(data)
}
