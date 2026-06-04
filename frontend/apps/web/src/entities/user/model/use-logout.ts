import { useNavigate } from 'react-router-dom'

import { useAuthStore } from './auth-store'

/** 세션을 비우고 홈으로 이동하는 로그아웃 핸들러. */
export function useLogout() {
  const navigate = useNavigate()
  const clear = useAuthStore((s) => s.clear)
  return () => {
    clear()
    navigate('/')
  }
}
