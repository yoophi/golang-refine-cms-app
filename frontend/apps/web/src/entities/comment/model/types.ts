export type PublicCommentStatus = 'pending' | 'approved' | 'spam'

/** 공개 API 댓글 표현(snake_case). 공개 목록에는 보통 승인(approved)된 댓글만 노출된다. */
export interface PublicComment {
  id: number
  post_id: number
  parent_id: number | null
  author_name: string
  author_email: string
  content: string
  status: PublicCommentStatus
  created_at: string
}

/** 댓글 작성 요청 본문. */
export interface CreateCommentInput {
  post_id: number
  author_name: string
  author_email: string
  content: string
  parent_id?: number | null
}
