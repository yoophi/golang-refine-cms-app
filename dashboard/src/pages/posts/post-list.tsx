import { useMemo } from 'react'
import { useMany, useNavigation } from '@refinedev/core'
import { useTable } from '@refinedev/react-table'
import { type ColumnDef, flexRender } from '@tanstack/react-table'
import { Plus } from 'lucide-react'

import type { Category } from '@/entities/category'
import {
  POST_STATUS_LABEL,
  POST_STATUS_VARIANT,
  type Post,
  type PostStatus,
} from '@/entities/post'
import type { Tag } from '@/entities/tag'
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

interface PostTableMeta {
  categoryData?: { data: Category[] }
  tagData?: { data: Tag[] }
}

export function PostList() {
  const { create } = useNavigation()

  const columns = useMemo<ColumnDef<Post>[]>(
    () => [
      { id: 'id', accessorKey: 'id', header: 'ID' },
      { id: 'title', accessorKey: 'title', header: '제목' },
      {
        id: 'status',
        accessorKey: 'status',
        header: '상태',
        cell: ({ getValue }) => {
          const status = getValue<PostStatus>()
          return (
            <Badge variant={POST_STATUS_VARIANT[status]}>
              {POST_STATUS_LABEL[status]}
            </Badge>
          )
        },
      },
      {
        id: 'category',
        accessorKey: 'categoryId',
        header: '카테고리',
        cell: ({ getValue, table }) => {
          const meta = table.options.meta as PostTableMeta | undefined
          const category = meta?.categoryData?.data.find(
            (item) => item.id === getValue<number>(),
          )
          return category?.title ?? '—'
        },
      },
      {
        id: 'tags',
        accessorKey: 'tagIds',
        header: '태그',
        cell: ({ getValue, table }) => {
          const meta = table.options.meta as PostTableMeta | undefined
          const ids = getValue<number[]>() ?? []
          const titles = ids
            .map((id) => meta?.tagData?.data.find((item) => item.id === id)?.title)
            .filter((title): title is string => Boolean(title))
          return (
            <div className="flex flex-wrap gap-1">
              {titles.map((title) => (
                <Badge key={title} variant="outline">
                  {title}
                </Badge>
              ))}
            </div>
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
        cell: ({ row }) => <RowActions resource="posts" id={row.original.id} />,
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
  } = useTable<Post>({ columns, refineCoreProps: { resource: 'posts' } })

  const posts = tableQuery?.data?.data ?? []
  const categoryIds = posts.map((post) => post.categoryId)
  const tagIds = Array.from(new Set(posts.flatMap((post) => post.tagIds ?? [])))

  const { result: categoryData } = useMany<Category>({
    resource: 'categories',
    ids: categoryIds,
    queryOptions: { enabled: categoryIds.length > 0 },
  })
  const { result: tagData } = useMany<Tag>({
    resource: 'tags',
    ids: tagIds,
    queryOptions: { enabled: tagIds.length > 0 },
  })

  setOptions((prev) => ({
    ...prev,
    meta: { ...prev.meta, categoryData, tagData },
  }))

  const rows = getRowModel().rows

  return (
    <div>
      <PageHeader title="게시글" description="블로그 게시글을 관리합니다.">
        <Button onClick={() => create('posts')}>
          <Plus className="size-4" /> 게시글 작성
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
