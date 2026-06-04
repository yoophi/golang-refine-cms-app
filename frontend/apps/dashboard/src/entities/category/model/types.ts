export interface Category {
  id: number
  name: string
  slug: string
  description: string
  parentId: number | null
  createdAt: string
  updatedAt: string
}
