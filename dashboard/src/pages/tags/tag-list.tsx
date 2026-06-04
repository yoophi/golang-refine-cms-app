import { useMemo } from 'react'
import { useNavigation } from '@refinedev/core'
import { useTable } from '@refinedev/react-table'
import { type ColumnDef } from '@tanstack/react-table'
import { Plus } from 'lucide-react'

import type { Tag } from '@/entities/tag'
import { formatDateTime } from '@/shared/lib'
import { Button, PageHeader } from '@/shared/ui'
import { DataTable, RowActions } from '@/widgets/crud-actions'

export function TagList() {
  const { create } = useNavigation()

  const columns = useMemo<ColumnDef<Tag>[]>(
    () => [
      { id: 'id', accessorKey: 'id', header: 'ID' },
      { id: 'name', accessorKey: 'name', header: '이름' },
      { id: 'slug', accessorKey: 'slug', header: '슬러그' },
      {
        id: 'createdAt',
        accessorKey: 'createdAt',
        header: '생성일',
        cell: ({ getValue }) => formatDateTime(getValue<string>()),
      },
      {
        id: 'actions',
        header: () => <span className="sr-only">작업</span>,
        cell: ({ row }) => <RowActions resource="tags" id={row.original.id} />,
      },
    ],
    [],
  )

  const { reactTable } = useTable<Tag>({
    columns,
    refineCoreProps: { resource: 'tags' },
  })

  return (
    <div>
      <PageHeader title="태그" description="게시글에 붙일 태그를 관리합니다.">
        <Button onClick={() => create('tags')}>
          <Plus className="size-4" /> 태그 생성
        </Button>
      </PageHeader>
      <DataTable table={reactTable} columnCount={columns.length} />
    </div>
  )
}
