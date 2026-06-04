export type PublicCommentStatus = 'pending' | 'approved' | 'spam'

/** 공개 API 댓글 표현(snake_case). `user_id`로 작성 회원을 식별한다. */
export interface PublicComment {
  id: number
  post_id: number
  user_id: number | null
  author_name: string
  content: string
  status: PublicCommentStatus
  created_at: string
}

/** 로그인 회원의 댓글 작성 요청. 작성자는 토큰에서 유도되므로 보내지 않는다. */
export interface CreateCommentInput {
  post_id: number
  content: string
}

/** 본인 댓글 수정 요청. */
export interface UpdateCommentInput {
  content: string
}
