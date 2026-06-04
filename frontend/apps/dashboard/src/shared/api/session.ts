import type { AdminUser } from '@/entities/auth'

/**
 * 관리자 세션(액세스/리프레시 토큰 + 신원)을 localStorage에 보관한다.
 * authProvider / accessControlProvider / httpClient가 공유한다.
 */
const ACCESS_TOKEN_KEY = 'cms.admin.accessToken'
const REFRESH_TOKEN_KEY = 'cms.admin.refreshToken'
const USER_KEY = 'cms.admin.user'

export function getAccessToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY)
}

export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN_KEY)
}

export function getStoredUser(): AdminUser | null {
  const raw = localStorage.getItem(USER_KEY)
  if (!raw) return null
  try {
    return JSON.parse(raw) as AdminUser
  } catch {
    return null
  }
}

export function setSession(params: {
  accessToken: string
  refreshToken?: string | null
  user: AdminUser
}): void {
  localStorage.setItem(ACCESS_TOKEN_KEY, params.accessToken)
  if (params.refreshToken) {
    localStorage.setItem(REFRESH_TOKEN_KEY, params.refreshToken)
  }
  localStorage.setItem(USER_KEY, JSON.stringify(params.user))
}

export function clearSession(): void {
  localStorage.removeItem(ACCESS_TOKEN_KEY)
  localStorage.removeItem(REFRESH_TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
}
