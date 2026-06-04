import type { CommentStatus } from '@/entities/comment'

export interface CommentFormValues {
  postId: number
  text: string
  status: CommentStatus
}
