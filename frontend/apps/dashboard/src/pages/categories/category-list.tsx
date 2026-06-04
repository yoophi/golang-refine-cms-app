import { useMemo } from 'react'
import { CanAccess, useNavigation } from '@refinedev/core'
import { useTable } from '@refinedev/react-table'
import { type ColumnDef } from '@tanstack/react-table'
import { Plus } from 'lucide-react'

import type { Category } from '@/entities/category'
import { formatDateTime } from '@/shared/lib'
import { Button, PageHeader } from '@/shared/ui'
import { DataTable, RowActions } from '@/widgets/crud-actions'

export function CategoryList() {
  const { create } = useNavigation()

  const columns = useMemo<ColumnDef<Category>[]>(
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
        cell: ({ row }) => <RowActions resource="categories" id={row.original.id} />,
      },
    ],
    [],
  )

  const { reactTable } = useTable<Category>({
    columns,
    refineCoreProps: { resource: 'categories' },
  })

  return (
    <div>
      <PageHeader title="카테고리" description="게시글 분류 카테고리를 관리합니다.">
        <CanAccess resource="categories" action="create">
          <Button onClick={() => create('categories')}>
            <Plus className="size-4" /> 카테고리 생성
          </Button>
        </CanAccess>
      </PageHeader>
      <DataTable table={reactTable} columnCount={columns.length} />
    </div>
  )
}
