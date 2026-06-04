import type { PostStatus } from '@/entities/post'

export interface PostFormValues {
  title: string
  slug: string
  excerpt: string
  content: string
  status: PostStatus
  categoryId: number | null
  tagIds: number[]
}
