import axios from 'axios'

import { API_URL } from '@/shared/config/api'

/** 공개 API 호출용 axios 인스턴스(인증 불필요). */
export const httpClient = axios.create({ baseURL: API_URL })
