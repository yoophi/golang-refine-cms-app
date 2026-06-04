import { httpClient } from '@/shared/api/http-client'

import type { CreateCommentInput, PublicComment } from './model/types'

/** 게시글의 댓글 목록. `GET /posts/:id/comments`, 엔벨로프 `{ data: [...] }`. */
export async function listPostComments(
  postId: number | string,
): Promise<PublicComment[]> {
  const { data } = await httpClient.get<{ data: PublicComment[] }>(
    `/posts/${postId}/comments`,
  )
  return data.data
}

/** 댓글 작성. `POST /comments`, 생성된 댓글 반환. */
export async function createComment(
  input: CreateCommentInput,
): Promise<PublicComment> {
  const { data } = await httpClient.post<PublicComment>('/comments', input)
  return data
}
