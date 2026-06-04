import axios from 'axios'

import { API_URL } from '@/shared/config/api'

import { getAccessToken } from './session'

/**
 * 관리자 API 호출용 axios 인스턴스.
 * - baseURL = 관리자 API 베이스(`/admin/api/v1`)
 * - 모든 요청에 `Authorization: Bearer <accessToken>` 자동 주입
 *
 * simple-rest dataProvider와 authProvider가 공유한다.
 */
export const httpClient = axios.create({ baseURL: API_URL })

httpClient.interceptors.request.use((config) => {
  const token = getAccessToken()
  if (token) {
    config.headers.set('Authorization', `Bearer ${token}`)
  }
  return config
})
