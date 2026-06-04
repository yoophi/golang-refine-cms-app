import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'

import {
  deleteComment,
  updateComment,
  type PublicComment,
} from '@/entities/comment'
import { formatDate } from '@/shared/lib'
import { Button, Textarea } from '@/shared/ui'

interface CommentItemProps {
  comment: PublicComment
  postId: number
  currentUserId?: number
}

export function CommentItem({
  comment,
  postId,
  currentUserId,
}: CommentItemProps) {
  const queryClient = useQueryClient()
  const isOwner = currentUserId != null && comment.user_id === currentUserId
  const [editing, setEditing] = useState(false)
  const [content, setContent] = useState(comment.content)

  const invalidate = () =>
    queryClient.invalidateQueries({ queryKey: ['comments', postId] })

  const updateMutation = useMutation({
    mutationFn: () => updateComment(comment.id, { content }),
    onSuccess: () => {
      setEditing(false)
      invalidate()
    },
  })
  const deleteMutation = useMutation({
    mutationFn: () => deleteComment(comment.id),
    onSuccess: invalidate,
  })

  return (
    <li className="border-b pb-4 last:border-b-0">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2 text-sm">
          <span className="font-medium">{comment.author_name}</span>
          <span className="text-muted-foreground">
            {formatDate(comment.created_at)}
          </span>
        </div>
        {isOwner && !editing ? (
          <div className="flex gap-1">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                setContent(comment.content)
                setEditing(true)
              }}
            >
              수정
            </Button>
            <Button
              variant="ghost"
              size="sm"
              disabled={deleteMutation.isPending}
              onClick={() => {
                if (window.confirm('댓글을 삭제할까요?')) {
                  deleteMutation.mutate()
                }
              }}
            >
              삭제
            </Button>
          </div>
        ) : null}
      </div>

      {editing ? (
        <div className="mt-2 space-y-2">
          <Textarea
            rows={3}
            value={content}
            onChange={(event) => setContent(event.target.value)}
          />
          <div className="flex gap-2">
            <Button
              size="sm"
              disabled={updateMutation.isPending}
              onClick={() => updateMutation.mutate()}
            >
              저장
            </Button>
            <Button size="sm" variant="outline" onClick={() => setEditing(false)}>
              취소
            </Button>
          </div>
        </div>
      ) : (
        <p className="mt-1 whitespace-pre-wrap text-sm">{comment.content}</p>
      )}
    </li>
  )
}
