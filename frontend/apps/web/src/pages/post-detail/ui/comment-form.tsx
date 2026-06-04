import { type FormEvent, useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Link } from 'react-router-dom'

import { createComment } from '@/entities/comment'
import { useAuthStore } from '@/entities/user'
import { Button, Textarea } from '@/shared/ui'

export function CommentForm({ postId }: { postId: number }) {
  const user = useAuthStore((s) => s.user)
  const queryClient = useQueryClient()
  const [content, setContent] = useState('')

  const mutation = useMutation({
    mutationFn: createComment,
    onSuccess: () => {
      setContent('')
      queryClient.invalidateQueries({ queryKey: ['comments', postId] })
    },
  })

  if (!user) {
    return (
      <p className="text-muted-foreground text-sm">
        댓글을 작성하려면{' '}
        <Link to="/login" className="text-foreground underline">
          로그인
        </Link>
        하세요.
      </p>
    )
  }

  const onSubmit = (event: FormEvent) => {
    event.preventDefault()
    mutation.mutate({ post_id: postId, content })
  }

  return (
    <form onSubmit={onSubmit} className="space-y-3">
      <Textarea
        rows={4}
        value={content}
        onChange={(event) => setContent(event.target.value)}
        placeholder={`${user.name} 님으로 댓글 작성`}
        required
      />
      <div className="flex items-center gap-3">
        <Button type="submit" disabled={mutation.isPending}>
          댓글 등록
        </Button>
        {mutation.isError ? (
          <span className="text-destructive text-sm">
            등록에 실패했습니다.
          </span>
        ) : null}
        {mutation.isSuccess ? (
          <span className="text-muted-foreground text-sm">
            등록되었습니다. 승인 후 표시될 수 있습니다.
          </span>
        ) : null}
      </div>
    </form>
  )
}
