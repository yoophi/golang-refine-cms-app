import { useNavigation, type HttpError } from '@refinedev/core'
import { useForm } from '@refinedev/react-hook-form'

import type { Comment } from '@/entities/comment'
import { Button, Card, CardContent, PageHeader } from '@/shared/ui'

import { CommentFields } from './ui/comment-fields'
import type { CommentFormValues } from './ui/comment-form-values'

export function CommentEdit() {
  const { list } = useNavigation()
  const {
    refineCore: { onFinish, formLoading },
    register,
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<Comment, HttpError, CommentFormValues>({
    refineCoreProps: { resource: 'comments', action: 'edit' },
  })

  return (
    <div>
      <PageHeader title="댓글 수정" />
      <Card>
        <CardContent>
          <form onSubmit={handleSubmit(onFinish)} className="space-y-6">
            <CommentFields
              register={register}
              control={control}
              errors={errors}
            />
            <div className="flex gap-2">
              <Button type="submit" disabled={formLoading}>
                저장
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => list('comments')}
              >
                취소
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
