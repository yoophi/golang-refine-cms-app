import type { PublicTag } from '@/entities/tag'

export type PublicPostStatus = 'draft' | 'published' | 'archived'

/** 공개 API(`/api/v1`) 게시글 표현(snake_case). */
export interface PublicPost {
  id: number
  title: string
  slug: string
  excerpt: string
  content: string
  status: PublicPostStatus
  category_id: number | null
  tags: PublicTag[]
  published_at: string | null
  created_at: string
  updated_at: string
}
