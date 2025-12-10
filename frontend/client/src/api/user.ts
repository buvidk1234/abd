import type { BackendUserInfo, LoginPayload, RegisterPayload, UpdateUserPayload, UserProfile } from '@/types/backend'
import { http } from '@/http/http'

function toUserProfile(raw?: BackendUserInfo | null): UserProfile {
  if (!raw) {
    return { userID: '' }
  }
  return {
    userID: String(raw.user_id ?? ''),
    nickname: raw.nickname,
    avatarURL: raw.avatar_url,
    gender: raw.gender,
    signature: raw.signature,
  }
}

// 用户注册
export async function registerUser(payload: RegisterPayload) {
  const userID = await http.post<string>('/user/register', {
    username: payload.username,
    password: payload.password,
    phone: payload.phone,
    email: payload.email,
    nickname: payload.nickname,
    avatar_url: payload.avatarURL,
  })
  return { userID: typeof userID === 'string' ? userID : String(userID ?? '') }
}

// 用户登录，返回 token 字符串
export async function loginUser(payload: LoginPayload) {
  const token = await http.post<string>('/user/login', payload)
  return { token: typeof token === 'string' ? token : String((token as any) ?? '') }
}

// 更新用户资料（需登录）
export function updateUserInfo(payload: UpdateUserPayload) {
  return http.post('/user/update-info', {
    nickname: payload.nickname,
    avatar_url: payload.avatarURL,
    gender: payload.gender,
    signature: payload.signature,
    birth: payload.birth,
    phone: payload.phone,
    email: payload.email,
    ex: payload.ex,
  })
}

// 获取当前登录用户的公开信息
export async function fetchCurrentUser() {
  const res = await http.get<BackendUserInfo>('/user/info')
  return toUserProfile(res)
}

// 兼容旧用法，返回数组形式的用户列表
export async function fetchUsers() {
  const profile = await fetchCurrentUser()
  return profile.userID ? [profile] : []
}
