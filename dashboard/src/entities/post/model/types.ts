export const POST_STATUSES = ['draft', 'published', 'rejected'] as const

export type PostStatus = (typeof POST_STATUSES)[number]

export interface Post {
  id: number
  title: string
  content: string
  status: PostStatus
  categoryId: number
  tagIds: number[]
  createdAt: string
  updatedAt: string
}

export const POST_STATUS_LABEL: Record<PostStatus, string> = {
  draft: '초안',
  published: '게시됨',
  rejected: '반려됨',
}

export const POST_STATUS_VARIANT: Record<
  PostStatus,
  'default' | 'secondary' | 'destructive' | 'outline'
> = {
  draft: 'secondary',
  published: 'default',
  rejected: 'destructive',
}
