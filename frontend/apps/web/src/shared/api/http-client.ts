import axios from 'axios'

import { API_URL } from '@/shared/config/api'

import { getAccessToken } from './auth-token'

/**
 * 공개 API 호출용 axios 인스턴스.
 * 회원 로그인 시 발급된 토큰이 있으면 `Authorization: Bearer`를 주입한다
 * (비로그인 요청은 헤더 없이 그대로 동작).
 */
export const httpClient = axios.create({ baseURL: API_URL })

httpClient.interceptors.request.use((config) => {
  const token = getAccessToken()
  if (token) {
    config.headers.set('Authorization', `Bearer ${token}`)
  }
  return config
})
