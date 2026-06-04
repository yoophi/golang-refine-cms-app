/**
 * 사용자용 웹앱이 호출하는 공개(public) API 베이스 URL.
 * 관리자 API(`/admin/api/v1`)와 분리된 `/api/v1` 네임스페이스를 사용한다.
 */
export const API_URL =
  import.meta.env.VITE_API_URL ?? 'http://localhost:9000/api/v1'
