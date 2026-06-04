import { httpClient } from '@/shared/api/http-client'

import type { PublicPost } from './model/types'

export interface ListPostsParams {
  categoryId?: number
  /** 기본값 'published' — 공개 사이트는 발행글만 노출. */
  status?: PublicPost['status']
}

/** 공개 게시글 목록. 응답 엔벨로프: `{ data: [...] }`. */
export async function listPosts(
  params: ListPostsParams = {},
): Promise<PublicPost[]> {
  const { data } = await httpClient.get<{ data: PublicPost[] }>('/posts', {
    params: {
      status: params.status ?? 'published',
      category_id: params.categoryId,
    },
  })
  return data.data
}

/** 공개 게시글 단건. 응답은 게시글 객체(엔벨로프 없음). */
export async function getPost(id: number | string): Promise<PublicPost> {
  const { data } = await httpClient.get<PublicPost>(`/posts/${id}`)
  return data
}
