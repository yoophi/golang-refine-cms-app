import { useNavigation, type HttpError } from '@refinedev/core'
import { useForm } from '@refinedev/react-hook-form'

import type { Tag } from '@/entities/tag'
import { Button, Card, CardContent, PageHeader } from '@/shared/ui'

import { TagFields, type TagFormValues } from './ui/tag-fields'

export function TagEdit() {
  const { list } = useNavigation()
  const {
    refineCore: { onFinish, formLoading },
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Tag, HttpError, TagFormValues>({
    refineCoreProps: { resource: 'tags', action: 'edit' },
  })

  return (
    <div>
      <PageHeader title="태그 수정" />
      <Card>
        <CardContent>
          <form onSubmit={handleSubmit(onFinish)} className="space-y-6">
            <TagFields register={register} errors={errors} />
            <div className="flex gap-2">
              <Button type="submit" disabled={formLoading}>
                저장
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => list('tags')}
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
