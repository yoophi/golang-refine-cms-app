import { create } from 'zustand'
import { persist } from 'zustand/middleware'

import { setAccessToken } from '@/shared/api/auth-token'

import type { AuthResponse, User } from './types'

interface AuthState {
  user: User | null
  /** 로그인/가입 응답으로 세션을 설정한다(토큰은 shared 보관소에도 기록). */
  setSession: (auth: AuthResponse) => void
  /** 사용자 정보만 갱신(프로필 수정 후). */
  setUser: (user: User) => void
  /** 로그아웃: 토큰·사용자 제거. */
  clear: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      setSession: (auth) => {
        setAccessToken(auth.accessToken)
        set({ user: auth.user })
      },
      setUser: (user) => set({ user }),
      clear: () => {
        setAccessToken(null)
        set({ user: null })
      },
    }),
    { name: 'cms.web.auth', partialize: (state) => ({ user: state.user }) },
  ),
)
