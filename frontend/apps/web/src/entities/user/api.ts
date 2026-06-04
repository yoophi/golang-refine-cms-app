import { httpClient } from '@/shared/api/http-client'

import type {
  AuthResponse,
  LoginInput,
  RegisterInput,
  UpdateProfileInput,
  User,
} from './model/types'

export async function register(input: RegisterInput): Promise<AuthResponse> {
  const { data } = await httpClient.post<AuthResponse>('/auth/register', input)
  return data
}

export async function login(input: LoginInput): Promise<AuthResponse> {
  const { data } = await httpClient.post<AuthResponse>('/auth/login', input)
  return data
}

export async function getMe(): Promise<User> {
  const { data } = await httpClient.get<User>('/auth/me')
  return data
}

export async function updateMe(input: UpdateProfileInput): Promise<User> {
  const { data } = await httpClient.patch<User>('/auth/me', input)
  return data
}
