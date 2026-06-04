import { useNavigation, type HttpError } from '@refinedev/core'
import { useForm } from '@refinedev/react-hook-form'

import type { Category } from '@/entities/category'
import { Button, Card, CardContent, PageHeader } from '@/shared/ui'

import { CategoryFields } from './ui/category-fields'
import type { CategoryFormValues } from './ui/category-form-values'

export function CategoryEdit() {
  const { list } = useNavigation()
  const {
    refineCore: { onFinish, formLoading },
    register,
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<Category, HttpError, CategoryFormValues>({
    refineCoreProps: { resource: 'categories', action: 'edit' },
  })

  return (
    <div>
      <PageHeader title="카테고리 수정" />
      <Card>
        <CardContent>
          <form onSubmit={handleSubmit(onFinish)} className="space-y-6">
            <CategoryFields
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
                onClick={() => list('categories')}
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
