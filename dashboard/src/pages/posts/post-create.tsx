import { useNavigation, type HttpError } from '@refinedev/core'
import { useForm } from '@refinedev/react-hook-form'

import type { Post } from '@/entities/post'
import { Button, Card, CardContent, PageHeader } from '@/shared/ui'

import { PostFields } from './ui/post-fields'
import type { PostFormValues } from './ui/post-form-values'

export function PostCreate() {
  const { list } = useNavigation()
  const {
    refineCore: { onFinish, formLoading },
    register,
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<Post, HttpError, PostFormValues>({
    refineCoreProps: { resource: 'posts', action: 'create' },
  })

  return (
    <div>
      <PageHeader title="게시글 작성" />
      <Card>
        <CardContent>
          <form onSubmit={handleSubmit(onFinish)} className="space-y-6">
            <PostFields register={register} control={control} errors={errors} />
            <div className="flex gap-2">
              <Button type="submit" disabled={formLoading}>
                저장
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => list('posts')}
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
