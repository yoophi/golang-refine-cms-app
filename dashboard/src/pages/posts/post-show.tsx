import { useMany, useNavigation, useOne, useShow } from '@refinedev/core'
import { Pencil } from 'lucide-react'

import type { Category } from '@/entities/category'
import {
  POST_STATUS_LABEL,
  POST_STATUS_VARIANT,
  type Post,
} from '@/entities/post'
import type { Tag } from '@/entities/tag'
import { formatDateTime } from '@/shared/lib'
import { Badge, Button, Card, CardContent, PageHeader } from '@/shared/ui'

export function PostShow() {
  const { list, edit } = useNavigation()
  const { query, result: record } = useShow<Post>({ resource: 'posts' })

  const { result: category } = useOne<Category>({
    resource: 'categories',
    id: record?.categoryId ?? 0,
    queryOptions: { enabled: Boolean(record?.categoryId) },
  })
  const { result: tagData } = useMany<Tag>({
    resource: 'tags',
    ids: record?.tagIds ?? [],
    queryOptions: { enabled: Boolean(record?.tagIds?.length) },
  })

  return (
    <div>
      <PageHeader title="게시글 상세">
        <Button variant="outline" onClick={() => list('posts')}>
          목록
        </Button>
        {record ? (
          <Button onClick={() => edit('posts', record.id)}>
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
              <dt className="text-muted-foreground">제목</dt>
              <dd>{record.title}</dd>
              <dt className="text-muted-foreground">상태</dt>
              <dd>
                <Badge variant={POST_STATUS_VARIANT[record.status]}>
                  {POST_STATUS_LABEL[record.status]}
                </Badge>
              </dd>
              <dt className="text-muted-foreground">카테고리</dt>
              <dd>{category?.title ?? '—'}</dd>
              <dt className="text-muted-foreground">태그</dt>
              <dd className="flex flex-wrap gap-1">
                {tagData?.data.length ? (
                  tagData.data.map((tag) => (
                    <Badge key={tag.id} variant="outline">
                      {tag.title}
                    </Badge>
                  ))
                ) : (
                  <span>—</span>
                )}
              </dd>
              <dt className="text-muted-foreground">본문</dt>
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
