import { http } from '@/lib/http'

export interface UserInfoDTO {
  username: string
  user_id: string
  nickname: string
  avatar_url: string
  gender: number
  signature: string
}

// ============ Request Params ============

export interface UserLoginReq {
  password: string
  username: string
}

export interface UserRegisterReq {
  username: string
  password: string
  phone?: string
  email?: string
  nickname?: string
  avatar_url?: string
}

export interface UpdateUserInfoReq {
  nickname?: string
  avatar_url?: string
  gender?: number | string
  signature?: string
  birth?: string
  phone?: string
  email?: string
  ex?: string
}

// ============ API Functions ============

export const loginApi = (data: UserLoginReq) => http.post<string>('/user/login', data)

export const registerApi = (data: UserRegisterReq) => http.post<string>('/user/register', data)

export const updateUserInfoApi = (data: UpdateUserInfoReq) =>
  http.post<void>('/user/update-info', data)

export const getUserInfoApi = () => http.get<UserInfoDTO>('/user/info')
