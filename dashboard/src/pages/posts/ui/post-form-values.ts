import type { PostStatus } from '@/entities/post'

export interface PostFormValues {
  title: string
  content: string
  status: PostStatus
  categoryId: number
  tagIds: number[]
}
