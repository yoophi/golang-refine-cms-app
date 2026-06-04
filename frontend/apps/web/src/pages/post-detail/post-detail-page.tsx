import { useQuery } from '@tanstack/react-query'
import { Link, useParams } from 'react-router-dom'

import { listPostComments } from '@/entities/comment'
import { getPost } from '@/entities/post'
import { useAuthStore } from '@/entities/user'
import { formatDate } from '@/shared/lib'
import { Badge } from '@/shared/ui'

import { CommentForm } from './ui/comment-form'
import { CommentItem } from './ui/comment-item'

export function PostDetailPage() {
  const { id } = useParams()
  const postId = Number(id)
  const currentUserId = useAuthStore((s) => s.user?.id)

  const {
    data: post,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ['post', postId],
    queryFn: () => getPost(postId),
    enabled: Number.isFinite(postId),
  })

  const { data: comments = [] } = useQuery({
    queryKey: ['comments', postId],
    queryFn: () => listPostComments(postId),
    enabled: Number.isFinite(postId),
  })

  if (isLoading) {
    return <p className="text-muted-foreground text-sm">불러오는 중…</p>
  }
  if (isError || !post) {
    return (
      <div className="space-y-4">
        <p className="text-destructive text-sm">
          게시글을 불러오지 못했습니다.
        </p>
        <Link to="/" className="text-sm underline">
          목록으로
        </Link>
      </div>
    )
  }

  return (
    <article className="space-y-8">
      <Link to="/" className="text-muted-foreground text-sm hover:underline">
        ← 목록으로
      </Link>

      <header className="space-y-3">
        <h1 className="text-3xl font-bold tracking-tight">{post.title}</h1>
        <p className="text-muted-foreground text-sm">
          {formatDate(post.published_at ?? post.created_at)}
        </p>
        {post.tags.length > 0 ? (
          <div className="flex flex-wrap gap-1">
            {post.tags.map((tag) => (
              <Badge key={tag.id} variant="secondary">
                {tag.name}
              </Badge>
            ))}
          </div>
        ) : null}
      </header>

      <div className="whitespace-pre-wrap leading-relaxed">{post.content}</div>

      <section className="space-y-6 border-t pt-8">
        <h2 className="text-xl font-semibold">댓글 {comments.length}</h2>

        {comments.length > 0 ? (
          <ul className="space-y-4">
            {comments.map((comment) => (
              <CommentItem
                key={comment.id}
                comment={comment}
                postId={post.id}
                currentUserId={currentUserId}
              />
            ))}
          </ul>
        ) : (
          <p className="text-muted-foreground text-sm">
            첫 댓글을 남겨보세요.
          </p>
        )}

        <div className="space-y-3">
          <h3 className="font-medium">댓글 작성</h3>
          <CommentForm postId={post.id} />
        </div>
      </section>
    </article>
  )
}
