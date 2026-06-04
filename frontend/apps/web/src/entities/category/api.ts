import { httpClient } from '@/shared/api/http-client'

import type { PublicCategory } from './model/types'

/** 공개 카테고리 목록. 응답 엔벨로프: `{ data: [...] }`. */
export async function listCategories(): Promise<PublicCategory[]> {
  const { data } = await httpClient.get<{ data: PublicCategory[] }>(
    '/categories',
  )
  return data.data
}
