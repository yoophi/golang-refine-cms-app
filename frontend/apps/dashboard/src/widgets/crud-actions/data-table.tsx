import { type Table as TanstackTable, flexRender } from '@tanstack/react-table'

import {
  Card,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/ui'

import { Pagination } from './pagination'

interface DataTableProps<TData> {
  table: TanstackTable<TData>
  columnCount: number
  emptyMessage?: string
}

/**
 * refine `useTable`의 TanStack 테이블 인스턴스를 받아 헤더/바디/페이지네이션을
 * 렌더링하는 공용 테이블. 리소스별 컬럼 정의만 다르고 렌더 구조는 동일하다.
 */
export function DataTable<TData>({
  table,
  columnCount,
  emptyMessage = '데이터가 없습니다.',
}: DataTableProps<TData>) {
  const rows = table.getRowModel().rows

  return (
    <>
      <Card className="py-0">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
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
                  colSpan={columnCount}
                  className="text-muted-foreground h-24 text-center"
                >
                  {emptyMessage}
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
        pageIndex={table.getState().pagination.pageIndex}
        pageCount={table.getPageCount()}
        canPrevious={table.getCanPreviousPage()}
        canNext={table.getCanNextPage()}
        onPageChange={table.setPageIndex}
      />
    </>
  )
}
