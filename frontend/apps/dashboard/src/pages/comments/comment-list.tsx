import { useMemo } from 'react'
import { CanAccess, useMany, useNavigation } from '@refinedev/core'
import { useTable } from '@refinedev/react-table'
import { type ColumnDef } from '@tanstack/react-table'
import { Plus } from 'lucide-react'

import {
  COMMENT_STATUS_LABEL,
  COMMENT_STATUS_VARIANT,
  type Comment,
  type CommentStatus,
} from '@/entities/comment'
import type { Post } from '@/entities/post'
import { formatDateTime } from '@/shared/lib'
import { Badge, Button, PageHeader } from '@/shared/ui'
import { DataTable, RowActions } from '@/widgets/crud-actions'

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
      { id: 'authorName', accessorKey: 'authorName', header: '작성자' },
      {
        id: 'content',
        accessorKey: 'content',
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
    reactTable,
    refineCore: { tableQuery },
  } = useTable<Comment>({ columns, refineCoreProps: { resource: 'comments' } })

  const comments = tableQuery?.data?.data ?? []
  const postIds = Array.from(new Set(comments.map((comment) => comment.postId)))

  const { result: postData } = useMany<Post>({
    resource: 'posts',
    ids: postIds,
    queryOptions: { enabled: postIds.length > 0 },
  })

  reactTable.setOptions((prev) => ({
    ...prev,
    meta: { ...prev.meta, postData },
  }))

  return (
    <div>
      <PageHeader title="댓글" description="게시글에 달린 댓글을 관리합니다.">
        <CanAccess resource="comments" action="create">
          <Button onClick={() => create('comments')}>
            <Plus className="size-4" /> 댓글 작성
          </Button>
        </CanAccess>
      </PageHeader>
      <DataTable table={reactTable} columnCount={columns.length} />
    </div>
  )
}
