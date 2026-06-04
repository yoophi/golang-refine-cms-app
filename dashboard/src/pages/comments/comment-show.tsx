import { useNavigation, useOne, useShow } from '@refinedev/core'
import { Pencil } from 'lucide-react'

import {
  COMMENT_STATUS_LABEL,
  COMMENT_STATUS_VARIANT,
  type Comment,
} from '@/entities/comment'
import type { Post } from '@/entities/post'
import { formatDateTime } from '@/shared/lib'
import { Badge, Button, Card, CardContent, PageHeader } from '@/shared/ui'

export function CommentShow() {
  const { list, edit } = useNavigation()
  const { query, result: record } = useShow<Comment>({ resource: 'comments' })

  const { result: post } = useOne<Post>({
    resource: 'posts',
    id: record?.postId ?? 0,
    queryOptions: { enabled: Boolean(record?.postId) },
  })

  return (
    <div>
      <PageHeader title="댓글 상세">
        <Button variant="outline" onClick={() => list('comments')}>
          목록
        </Button>
        {record ? (
          <Button onClick={() => edit('comments', record.id)}>
            <Pencil className="size-4" /> 수정
          </Button>
        ) : null}
      </PageHeader>
      <Card>
        <CardContent>
          {query.isLoading ? (
            <p className="text-muted-foreground text-sm">불러오는 중…</p>
          ) : record ? (
            <dl className="grid grid-cols-[8rem_1fr] gap-y-3 text-sm">
              <dt className="text-muted-foreground">ID</dt>
              <dd>{record.id}</dd>
              <dt className="text-muted-foreground">게시글</dt>
              <dd>{post?.title ?? `#${record.postId}`}</dd>
              <dt className="text-muted-foreground">작성자</dt>
              <dd>{record.authorName}</dd>
              <dt className="text-muted-foreground">이메일</dt>
              <dd>{record.authorEmail || '—'}</dd>
              <dt className="text-muted-foreground">상태</dt>
              <dd>
                <Badge variant={COMMENT_STATUS_VARIANT[record.status]}>
                  {COMMENT_STATUS_LABEL[record.status]}
                </Badge>
              </dd>
              <dt className="text-muted-foreground">내용</dt>
              <dd className="whitespace-pre-wrap">{record.content}</dd>
              <dt className="text-muted-foreground">생성일</dt>
              <dd>{formatDateTime(record.createdAt)}</dd>
              <dt className="text-muted-foreground">수정일</dt>
              <dd>{formatDateTime(record.updatedAt)}</dd>
            </dl>
          ) : (
            <p className="text-muted-foreground text-sm">데이터가 없습니다.</p>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
