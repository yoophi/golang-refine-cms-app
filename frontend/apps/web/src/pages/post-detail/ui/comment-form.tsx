import { type FormEvent, useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'

import { createComment } from '@/entities/comment'
import { useReaderStore } from '@/shared/store/reader-store'
import { Button, Input, Label, Textarea } from '@/shared/ui'

export function CommentForm({ postId }: { postId: number }) {
  const queryClient = useQueryClient()
  const storedName = useReaderStore((s) => s.authorName)
  const storedEmail = useReaderStore((s) => s.authorEmail)
  const setAuthor = useReaderStore((s) => s.setAuthor)

  const [authorName, setAuthorName] = useState(storedName)
  const [authorEmail, setAuthorEmail] = useState(storedEmail)
  const [content, setContent] = useState('')

  const mutation = useMutation({
    mutationFn: createComment,
    onSuccess: () => {
      setAuthor({ authorName, authorEmail })
      setContent('')
      queryClient.invalidateQueries({ queryKey: ['comments', postId] })
    },
  })

  const onSubmit = (event: FormEvent) => {
    event.preventDefault()
    mutation.mutate({
      post_id: postId,
      author_name: authorName,
      author_email: authorEmail,
      content,
    })
  }

  return (
    <form onSubmit={onSubmit} className="space-y-4">
      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="authorName">이름</Label>
          <Input
            id="authorName"
            value={authorName}
            onChange={(event) => setAuthorName(event.target.value)}
            required
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="authorEmail">이메일</Label>
          <Input
            id="authorEmail"
            type="email"
            value={authorEmail}
            onChange={(event) => setAuthorEmail(event.target.value)}
            required
          />
        </div>
      </div>
      <div className="space-y-2">
        <Label htmlFor="content">댓글</Label>
        <Textarea
          id="content"
          rows={4}
          value={content}
          onChange={(event) => setContent(event.target.value)}
          required
        />
      </div>
      <div className="flex items-center gap-3">
        <Button type="submit" disabled={mutation.isPending}>
          댓글 등록
        </Button>
        {mutation.isError ? (
          <span className="text-destructive text-sm">
            등록에 실패했습니다. 잠시 후 다시 시도해주세요.
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
