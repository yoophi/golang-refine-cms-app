/**
 * 대시보드가 호출하는 백엔드 API 호스트.
 * BE는 refine `@refinedev/simple-rest` 규약(json-server 스타일)에 맞춰 구현된다.
 * 문서: /docs/swagger.json, /docs/api.md
 */
export const API_URL =
  import.meta.env.VITE_API_URL ?? 'http://localhost:9000'
