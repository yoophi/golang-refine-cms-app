/** 관리자 역할. superadmin은 모든 권한을 가진다(ACL 우회). */
export const SUPERADMIN_ROLE = 'superadmin'

export type AdminRole = 'superadmin' | 'editor' | 'viewer'

/** ACL에서 사용하는 동작. `resource:action` 권한 문자열의 action 부분. */
export type PermissionAction = 'list' | 'show' | 'create' | 'edit' | 'delete'

/**
 * 로그인한 관리자 신원.
 * `permissions`는 `"posts:create"` 형식의 `resource:action` 문자열 배열이며,
 * `role === 'superadmin'`이면 권한 목록과 무관하게 전체 허용한다.
 */
export interface AdminUser {
  id: number
  name: string
  email: string
  role: AdminRole | string
  permissions: string[]
  avatar?: string
}

/** 로그인 요청 폼 값. */
export interface LoginCredentials {
  email: string
  password: string
}
