import { http } from '@/lib/http'

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

export interface UserInfo {
  username: string
  user_id:  string
  nickname: string
  avatar_url: string
  gender: number | string
  signature: string
}

export const login = (data: UserLoginReq) => {
  return http.post<string>('/user/login', data)
}

export const register = (data: UserRegisterReq) => {
  return http.post<string>('/user/register', data)
}

export const updateUserInfo = (data: UpdateUserInfoReq) => {
  return http.post<void>('/user/update-info', data)
}

export const getUserInfo = () => {
  return http.get<UserInfo>('/user/info')
}
