import { useNavigation } from '@refinedev/core'
import { useForm } from '@refinedev/react-hook-form'

import {
  Button,
  Card,
  CardContent,
  Input,
  Label,
  PageHeader,
} from '@/shared/ui'

export function TagCreate() {
  const { list } = useNavigation()
  const {
    refineCore: { onFinish, formLoading },
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({ refineCoreProps: { resource: 'tags', action: 'create' } })

  return (
    <div>
      <PageHeader title="태그 생성" />
      <Card>
        <CardContent>
          <form
            onSubmit={handleSubmit(onFinish)}
            className="max-w-lg space-y-4"
          >
            <div className="space-y-2">
              <Label htmlFor="title">제목</Label>
              <Input
                id="title"
                {...register('title', { required: '제목을 입력하세요.' })}
              />
              {errors.title ? (
                <p className="text-destructive text-sm">
                  {errors.title.message as string}
                </p>
              ) : null}
            </div>
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
