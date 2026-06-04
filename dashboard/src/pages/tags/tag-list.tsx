import { useMemo } from 'react'
import { useNavigation } from '@refinedev/core'
import { useTable } from '@refinedev/react-table'
import { type ColumnDef, flexRender } from '@tanstack/react-table'
import { Plus } from 'lucide-react'

import type { Tag } from '@/entities/tag'
import { formatDateTime } from '@/shared/lib'
import {
  Button,
  Card,
  PageHeader,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/ui'
import { Pagination, RowActions } from '@/widgets/crud-actions'

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

  const {
    reactTable: {
      getHeaderGroups,
      getRowModel,
      getState,
      getPageCount,
      getCanPreviousPage,
      getCanNextPage,
      setPageIndex,
    },
  } = useTable<Tag>({ columns, refineCoreProps: { resource: 'tags' } })

  const rows = getRowModel().rows

  return (
    <div>
      <PageHeader title="태그" description="게시글에 붙일 태그를 관리합니다.">
        <Button onClick={() => create('tags')}>
          <Plus className="size-4" /> 태그 생성
        </Button>
      </PageHeader>
      <Card className="py-0">
        <Table>
          <TableHeader>
            {getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {rows.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="text-muted-foreground h-24 text-center"
                >
                  데이터가 없습니다.
                </TableCell>
              </TableRow>
            ) : (
              rows.map((row) => (
                <TableRow key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </Card>
      <Pagination
        pageIndex={getState().pagination.pageIndex}
        pageCount={getPageCount()}
        canPrevious={getCanPreviousPage()}
        canNext={getCanNextPage()}
        onPageChange={setPageIndex}
      />
    </div>
  )
}
