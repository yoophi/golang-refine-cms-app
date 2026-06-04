import type { AuthProvider } from '@refinedev/core'

import type { LoginCredentials } from '@/entities/auth'

import { httpClient } from './http-client'
import {
  clearSession,
  getAccessToken,
  getStoredUser,
  setSession,
} from './session'

interface MaybeHttpError {
  statusCode?: number
  response?: { status?: number }
}

function statusOf(error: unknown): number | undefined {
  const e = error as MaybeHttpError
  return e?.statusCode ?? e?.response?.status
}

/**
 * 관리자 JWT 인증 프로바이더.
 * - `POST /auth/login` 으로 토큰 + 신원을 받아 세션에 저장
 * - 401/403 응답 시 자동 로그아웃 후 `/login` 으로 이동
 */
export const authProvider: AuthProvider = {
  login: async ({ email, password }: LoginCredentials) => {
    try {
      const { data } = await httpClient.post('/auth/login', { email, password })
      setSession({
        accessToken: data.accessToken,
        refreshToken: data.refreshToken,
        user: data.user,
      })
      return { success: true, redirectTo: '/' }
    } catch {
      return {
        success: false,
        error: {
          name: '로그인 실패',
          message: '이메일 또는 비밀번호를 확인하세요.',
        },
      }
    }
  },

  logout: async () => {
    clearSession()
    return { success: true, redirectTo: '/login' }
  },

  check: async () => {
    if (getAccessToken()) {
      return { authenticated: true }
    }
    return { authenticated: false, redirectTo: '/login', logout: true }
  },

  onError: async (error) => {
    const status = statusOf(error)
    if (status === 401 || status === 403) {
      return { logout: true, redirectTo: '/login', error }
    }
    return {}
  },

  getIdentity: async () => getStoredUser(),

  getPermissions: async () => {
    const user = getStoredUser()
    return user ? { role: user.role, permissions: user.permissions } : null
  },
}
