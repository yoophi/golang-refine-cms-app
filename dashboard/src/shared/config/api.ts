/**
 * 대시보드가 호출하는 관리자(admin) API 베이스 URL.
 *
 * 관리자 API는 사용자용(public) API(`/api/v1`)와 분리된 `/admin/api/v1`
 * 네임스페이스를 사용하며, refine `@refinedev/simple-rest` 규약
 * (목록 배열 + `X-Total-Count` 헤더, `_start/_end/_sort/_order`, `PATCH` 수정)을 따른다.
 * 문서: /docs/swagger.json, /docs/api.md
 */
export const API_URL =
  import.meta.env.VITE_API_URL ?? 'http://localhost:9000/admin/api/v1'
