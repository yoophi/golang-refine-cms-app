import { useQuery } from '@tanstack/react-query'
import { Link, useSearchParams } from 'react-router-dom'

import { listCategories } from '@/entities/category'
import { listPosts } from '@/entities/post'
import { formatDate } from '@/shared/lib'
import {
  Badge,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/shared/ui'

export function PostListPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const categoryParam = searchParams.get('category')
  const categoryId = categoryParam ? Number(categoryParam) : undefined

  const { data: categories = [] } = useQuery({
    queryKey: ['categories'],
    queryFn: listCategories,
  })
  const {
    data: posts = [],
    isLoading,
    isError,
  } = useQuery({
    queryKey: ['posts', { categoryId }],
    queryFn: () => listPosts({ categoryId }),
  })

  const selectCategory = (id?: number) => {
    setSearchParams(id ? { category: String(id) } : {})
  }

  return (
    <div className="space-y-8">
      <div className="flex flex-wrap gap-2">
        <Button
          variant={categoryId ? 'outline' : 'default'}
          size="sm"
          onClick={() => selectCategory(undefined)}
        >
          전체
        </Button>
        {categories.map((category) => (
          <Button
            key={category.id}
            variant={categoryId === category.id ? 'default' : 'outline'}
            size="sm"
            onClick={() => selectCategory(category.id)}
          >
            {category.name}
          </Button>
        ))}
      </div>

      {isLoading ? (
        <p className="text-muted-foreground text-sm">불러오는 중…</p>
      ) : isError ? (
        <p className="text-destructive text-sm">
          게시글을 불러오지 못했습니다.
        </p>
      ) : posts.length === 0 ? (
        <p className="text-muted-foreground text-sm">게시글이 없습니다.</p>
      ) : (
        <div className="space-y-4">
          {posts.map((post) => (
            <Card key={post.id} className="transition-colors hover:border-ring">
              <CardHeader>
                <CardTitle className="text-xl">
                  <Link to={`/posts/${post.id}`} className="hover:underline">
                    {post.title}
                  </Link>
                </CardTitle>
                <p className="text-muted-foreground text-sm">
                  {formatDate(post.published_at ?? post.created_at)}
                </p>
              </CardHeader>
              <CardContent className="space-y-3">
                {post.excerpt ? (
                  <p className="text-muted-foreground line-clamp-3">
                    {post.excerpt}
                  </p>
                ) : null}
                {post.tags.length > 0 ? (
                  <div className="flex flex-wrap gap-1">
                    {post.tags.map((tag) => (
                      <Badge key={tag.id} variant="secondary">
                        {tag.name}
                      </Badge>
                    ))}
                  </div>
                ) : null}
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
