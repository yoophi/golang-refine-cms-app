/** 공개 API 회원(end-user) 표현. */
export interface User {
  id: number
  email: string
  name: string
  avatar?: string
  created_at: string
  updated_at: string
}

export interface AuthResponse {
  accessToken: string
  refreshToken?: string
  user: User
}

export interface RegisterInput {
  email: string
  password: string
  name: string
}

export interface LoginInput {
  email: string
  password: string
}

export interface UpdateProfileInput {
  name?: string
  email?: string
  password?: string
  currentPassword?: string
}
