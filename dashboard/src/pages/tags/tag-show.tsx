import { useNavigation, useShow } from '@refinedev/core'
import { Pencil } from 'lucide-react'

import type { Tag } from '@/entities/tag'
import { formatDateTime } from '@/shared/lib'
import { Button, Card, CardContent, PageHeader } from '@/shared/ui'

export function TagShow() {
  const { list, edit } = useNavigation()
  const { query } = useShow<Tag>({ resource: 'tags' })
  const record = query.data?.data

  return (
    <div>
      <PageHeader title="태그 상세">
        <Button variant="outline" onClick={() => list('tags')}>
          목록
        </Button>
        {record ? (
          <Button onClick={() => edit('tags', record.id)}>
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
