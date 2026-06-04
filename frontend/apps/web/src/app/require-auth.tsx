import type { PropsWithChildren } from 'react'
import { Navigate } from 'react-router-dom'

import { useAuthStore } from '@/entities/user'

/** 비로그인 상태면 로그인 페이지로 보낸다. */
export function RequireAuth({ children }: PropsWithChildren) {
  const user = useAuthStore((s) => s.user)
  if (!user) {
    return <Navigate to="/login" replace />
  }
  return <>{children}</>
}
