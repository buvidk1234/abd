import type { LoginPayload, RegisterPayload, UpdateUserPayload, UserProfile } from '@/types/backend'
import { http } from '@/http/http'

// 用户注册
export function registerUser(payload: RegisterPayload) {
  return http.post<{ userID: string }>('/user/user_register', payload)
}

// 用户登录
export function loginUser(payload: LoginPayload) {
  return http.post<{ msg?: string }>('/user/user_login', payload)
}

// 更新用户资料
export function updateUserInfo(payload: UpdateUserPayload) {
  return http.post<{ msg?: string }>('/user/update_user_info', payload)
}

// 获取用户公开信息，如果 userIDs 为空则返回全部
export async function fetchUsers(userIDs: string[] = []) {
  const res = await http.post<{ userInfos: UserProfile[] }>('/user/get_users_info', { userIDs })
  // 后端直接返回 { userInfos: [] }，为了兼容直接取用
  return (res as any)?.userInfos || (res as any || []) as UserProfile[]
}
