import { useNavigation, useOne, useShow } from '@refinedev/core'
import { Pencil } from 'lucide-react'

import {
  COMMENT_STATUS_LABEL,
  COMMENT_STATUS_VARIANT,
  type Comment,
} from '@/entities/comment'
import type { Post } from '@/entities/post'
import { formatDateTime } from '@/shared/lib'
import {
  Badge,
  Button,
  Card,
  CardContent,
  DescriptionList,
  PageHeader,
} from '@/shared/ui'

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
            <DescriptionList
              items={[
                { label: 'ID', value: record.id },
                { label: '게시글', value: post?.title ?? `#${record.postId}` },
                { label: '작성자', value: record.authorName },
                { label: '이메일', value: record.authorEmail || '—' },
                {
                  label: '상태',
                  value: (
                    <Badge variant={COMMENT_STATUS_VARIANT[record.status]}>
                      {COMMENT_STATUS_LABEL[record.status]}
                    </Badge>
                  ),
                },
                {
                  label: '내용',
                  value: (
                    <span className="whitespace-pre-wrap">{record.content}</span>
                  ),
                },
                { label: '생성일', value: formatDateTime(record.createdAt) },
                { label: '수정일', value: formatDateTime(record.updatedAt) },
              ]}
            />
          ) : (
            <p className="text-muted-foreground text-sm">데이터가 없습니다.</p>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
