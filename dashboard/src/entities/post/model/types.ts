import type { Tag } from '@/entities/tag'

export const POST_STATUSES = ['draft', 'published', 'archived'] as const

export type PostStatus = (typeof POST_STATUSES)[number]

export interface Post {
  id: number
  title: string
  slug: string
  excerpt: string
  content: string
  status: PostStatus
  categoryId: number | null
  /** 목록/상세 응답에 임베드되는 태그 객체(읽기 전용). */
  tags?: Tag[]
  /** 생성/수정 시 사용하는 태그 ID 배열. */
  tagIds: number[]
  publishedAt: string | null
  createdAt: string
  updatedAt: string
}

export const POST_STATUS_LABEL: Record<PostStatus, string> = {
  draft: '초안',
  published: '게시됨',
  archived: '보관됨',
}

export const POST_STATUS_VARIANT: Record<
  PostStatus,
  'default' | 'secondary' | 'destructive' | 'outline'
> = {
  draft: 'secondary',
  published: 'default',
  archived: 'outline',
}
