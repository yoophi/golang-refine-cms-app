import { useNavigation, useShow } from '@refinedev/core'
import { Pencil } from 'lucide-react'

import type { Category } from '@/entities/category'
import { formatDateTime } from '@/shared/lib'
import {
  Button,
  Card,
  CardContent,
  DescriptionList,
  PageHeader,
} from '@/shared/ui'

export function CategoryShow() {
  const { list, edit } = useNavigation()
  const { query, result: record } = useShow<Category>({ resource: 'categories' })

  return (
    <div>
      <PageHeader title="카테고리 상세">
        <Button variant="outline" onClick={() => list('categories')}>
          목록
        </Button>
        {record ? (
          <Button onClick={() => edit('categories', record.id)}>
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
                { label: '이름', value: record.name },
                { label: '슬러그', value: record.slug },
                {
                  label: '설명',
                  value: (
                    <span className="whitespace-pre-wrap">
                      {record.description || '—'}
                    </span>
                  ),
                },
                {
                  label: '상위 카테고리',
                  value: record.parentId != null ? `#${record.parentId}` : '—',
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
