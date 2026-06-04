import { useMemo } from 'react'
import { useMany, useNavigation } from '@refinedev/core'
import { useTable } from '@refinedev/react-table'
import { type ColumnDef, flexRender } from '@tanstack/react-table'
import { Plus } from 'lucide-react'

import {
  COMMENT_STATUS_LABEL,
  COMMENT_STATUS_VARIANT,
  type Comment,
  type CommentStatus,
} from '@/entities/comment'
import type { Post } from '@/entities/post'
import { formatDateTime } from '@/shared/lib'
import {
  Badge,
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

interface CommentTableMeta {
  postData?: { data: Post[] }
}

export function CommentList() {
  const { create } = useNavigation()

  const columns = useMemo<ColumnDef<Comment>[]>(
    () => [
      { id: 'id', accessorKey: 'id', header: 'ID' },
      {
        id: 'post',
        accessorKey: 'postId',
        header: '게시글',
        cell: ({ getValue, table }) => {
          const meta = table.options.meta as CommentTableMeta | undefined
          const post = meta?.postData?.data.find(
            (item) => item.id === getValue<number>(),
          )
          return post?.title ?? '—'
        },
      },
      {
        id: 'text',
        accessorKey: 'text',
        header: '내용',
        cell: ({ getValue }) => (
          <span className="line-clamp-1 max-w-md">{getValue<string>()}</span>
        ),
      },
      {
        id: 'status',
        accessorKey: 'status',
        header: '상태',
        cell: ({ getValue }) => {
          const status = getValue<CommentStatus>()
          return (
            <Badge variant={COMMENT_STATUS_VARIANT[status]}>
              {COMMENT_STATUS_LABEL[status]}
            </Badge>
          )
        },
      },
      {
        id: 'createdAt',
        accessorKey: 'createdAt',
        header: '생성일',
        cell: ({ getValue }) => formatDateTime(getValue<string>()),
      },
      {
        id: 'actions',
        header: () => <span className="sr-only">작업</span>,
        cell: ({ row }) => <RowActions resource="comments" id={row.original.id} />,
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
      setOptions,
    },
    refineCore: { tableQuery },
  } = useTable<Comment>({ columns, refineCoreProps: { resource: 'comments' } })

  const comments = tableQuery?.data?.data ?? []
  const postIds = Array.from(new Set(comments.map((comment) => comment.postId)))

  const { result: postData } = useMany<Post>({
    resource: 'posts',
    ids: postIds,
    queryOptions: { enabled: postIds.length > 0 },
  })

  setOptions((prev) => ({
    ...prev,
    meta: { ...prev.meta, postData },
  }))

  const rows = getRowModel().rows

  return (
    <div>
      <PageHeader title="댓글" description="게시글에 달린 댓글을 관리합니다.">
        <Button onClick={() => create('comments')}>
          <Plus className="size-4" /> 댓글 작성
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
