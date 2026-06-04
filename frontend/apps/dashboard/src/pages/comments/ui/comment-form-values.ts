import type { CommentStatus } from '@/entities/comment'

export interface CommentFormValues {
  postId: number
  authorName: string
  authorEmail: string
  content: string
  status: CommentStatus
}
