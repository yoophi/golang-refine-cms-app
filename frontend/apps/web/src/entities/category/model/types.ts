export interface PublicCategory {
  id: number
  name: string
  slug: string
  description?: string
  parent_id: number | null
}
