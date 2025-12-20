/**
 * Domain Layer: User business models
 * - Page-independent types
 * - Pure functions and rules
 * - No React, no UI dependencies, no HTTP
 */

export interface User {
  id: string
  username: string
  nickname: string
  avatar?: string
  gender?: number
  signature?: string
}

export interface LoginCredentials {
  username: string
  password: string
}

export interface RegisterData {
  username: string
  password: string
  phone?: string
  email?: string
  nickname?: string
  avatarUrl?: string
}

export interface UpdateUserProfile {
  nickname?: string
  avatarUrl?: string
  gender?: number | string
  signature?: string
  birth?: string
  phone?: string
  email?: string
  ex?: string
}
