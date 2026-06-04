export const COMMENT_STATUSES = ['draft', 'published', 'rejected'] as const

export type CommentStatus = (typeof COMMENT_STATUSES)[number]

export interface Comment {
  id: number
  postId: number
  text: string
  status: CommentStatus
  createdAt: string
  updatedAt: string
}

export const COMMENT_STATUS_LABEL: Record<CommentStatus, string> = {
  draft: '대기',
  published: '승인',
  rejected: '거부',
}

export const COMMENT_STATUS_VARIANT: Record<
  CommentStatus,
  'default' | 'secondary' | 'destructive' | 'outline'
> = {
  draft: 'secondary',
  published: 'default',
  rejected: 'destructive',
}
