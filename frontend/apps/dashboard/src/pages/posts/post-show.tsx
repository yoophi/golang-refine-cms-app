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
import {
  Badge,
  Button,
  Card,
  CardContent,
  DescriptionList,
  PageHeader,
} from '@/shared/ui'

export function PostShow() {
  const { list, edit } = useNavigation()
  const { query, result: record } = useShow<Post>({ resource: 'posts' })

  const { result: category } = useOne<Category>({
    resource: 'categories',
    id: record?.categoryId ?? 0,
    queryOptions: { enabled: record?.categoryId != null },
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
            <DescriptionList
              items={[
                { label: 'ID', value: record.id },
                { label: '제목', value: record.title },
                { label: '슬러그', value: record.slug },
                {
                  label: '상태',
                  value: (
                    <Badge variant={POST_STATUS_VARIANT[record.status]}>
                      {POST_STATUS_LABEL[record.status]}
                    </Badge>
                  ),
                },
                { label: '카테고리', value: category?.name ?? '—' },
                {
                  label: '태그',
                  value: (
                    <div className="flex flex-wrap gap-1">
                      {tagData?.data.length ? (
                        tagData.data.map((tag) => (
                          <Badge key={tag.id} variant="outline">
                            {tag.name}
                          </Badge>
                        ))
                      ) : (
                        <span>—</span>
                      )}
                    </div>
                  ),
                },
                {
                  label: '요약',
                  value: (
                    <span className="whitespace-pre-wrap">
                      {record.excerpt || '—'}
                    </span>
                  ),
                },
                {
                  label: '본문',
                  value: (
                    <span className="whitespace-pre-wrap">{record.content}</span>
                  ),
                },
                {
                  label: '발행일',
                  value: record.publishedAt
                    ? formatDateTime(record.publishedAt)
                    : '—',
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
