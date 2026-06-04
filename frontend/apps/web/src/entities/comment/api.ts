import { httpClient } from '@/shared/api/http-client'

import type {
  CreateCommentInput,
  PublicComment,
  UpdateCommentInput,
} from './model/types'

/** 게시글의 댓글 목록. `GET /posts/:id/comments`, 엔벨로프 `{ data: [...] }`. */
export async function listPostComments(
  postId: number | string,
): Promise<PublicComment[]> {
  const { data } = await httpClient.get<{ data: PublicComment[] }>(
    `/posts/${postId}/comments`,
  )
  return data.data
}

/** 로그인 회원 댓글 작성. `POST /comments`(Bearer). */
export async function createComment(
  input: CreateCommentInput,
): Promise<PublicComment> {
  const { data } = await httpClient.post<PublicComment>('/comments', input)
  return data
}

/** 본인 댓글 수정. `PATCH /comments/:id`(Bearer, owner). */
export async function updateComment(
  id: number,
  input: UpdateCommentInput,
): Promise<PublicComment> {
  const { data } = await httpClient.patch<PublicComment>(
    `/comments/${id}`,
    input,
  )
  return data
}

/** 본인 댓글 삭제. `DELETE /comments/:id`(Bearer, owner). */
export async function deleteComment(id: number): Promise<void> {
  await httpClient.delete(`/comments/${id}`)
}
