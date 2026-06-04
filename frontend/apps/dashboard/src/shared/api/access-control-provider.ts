import type { AccessControlProvider } from '@refinedev/core'

import { SUPERADMIN_ROLE } from '@/entities/auth'

import { getStoredUser } from './session'

/**
 * refine 액션 → 권한 문자열의 action.
 * `clone`은 생성 권한으로 취급한다.
 */
function normalizeAction(action?: string): string | undefined {
  if (action === 'clone') return 'create'
  return action
}

/**
 * RBAC 접근 제어. 신원의 `permissions: ["resource:action"]` 목록으로 판단하며,
 * `role === 'superadmin'`이면 전체 허용한다.
 */
export const accessControlProvider: AccessControlProvider = {
  can: async ({ resource, action }) => {
    const user = getStoredUser()
    if (!user) {
      return { can: false, reason: '로그인이 필요합니다.' }
    }
    if (user.role === SUPERADMIN_ROLE) {
      return { can: true }
    }

    const act = normalizeAction(action)
    if (!resource || !act) {
      return { can: true }
    }

    const allowed = user.permissions.includes(`${resource}:${act}`)
    return allowed
      ? { can: true }
      : { can: false, reason: '권한이 없습니다.' }
  },
  options: {
    buttons: {
      enableAccessControl: true,
      hideIfUnauthorized: true,
    },
  },
}
