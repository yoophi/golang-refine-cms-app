export const COMMENT_STATUSES = ['pending', 'approved', 'spam'] as const

export type CommentStatus = (typeof COMMENT_STATUSES)[number]

export interface Comment {
  id: number
  postId: number
  parentId: number | null
  authorName: string
  authorEmail: string
  content: string
  status: CommentStatus
  createdAt: string
  updatedAt: string
}

export const COMMENT_STATUS_LABEL: Record<CommentStatus, string> = {
  pending: '대기',
  approved: '승인',
  spam: '스팸',
}

export const COMMENT_STATUS_VARIANT: Record<
  CommentStatus,
  'default' | 'secondary' | 'destructive' | 'outline'
> = {
  pending: 'secondary',
  approved: 'default',
  spam: 'destructive',
}
