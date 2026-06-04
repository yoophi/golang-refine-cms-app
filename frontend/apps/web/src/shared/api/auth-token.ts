/**
 * 회원 액세스 토큰 보관소(shared 전용).
 * httpClient(요청 인터셉터)와 회원 인증 스토어가 공유한다.
 * 토큰만 다루므로 도메인(entities)에 의존하지 않는다.
 */
const TOKEN_KEY = 'cms.web.accessToken'

let accessToken: string | null = localStorage.getItem(TOKEN_KEY)

export function getAccessToken(): string | null {
  return accessToken
}

export function setAccessToken(token: string | null): void {
  accessToken = token
  if (token) {
    localStorage.setItem(TOKEN_KEY, token)
  } else {
    localStorage.removeItem(TOKEN_KEY)
  }
}
