import { useSelect } from '@refinedev/core'
import {
  Controller,
  type Control,
  type FieldErrors,
  type UseFormRegister,
} from 'react-hook-form'

import type { Category } from '@/entities/category'
import {
  Input,
  Label,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Textarea,
} from '@/shared/ui'

import type { CategoryFormValues } from './category-form-values'

interface CategoryFieldsProps {
  register: UseFormRegister<CategoryFormValues>
  control: Control<CategoryFormValues>
  errors: FieldErrors<CategoryFormValues>
}

export function CategoryFields({
  register,
  control,
  errors,
}: CategoryFieldsProps) {
  const { options: parentOptions } = useSelect<Category>({
    resource: 'categories',
    optionLabel: 'name',
    optionValue: 'id',
  })

  return (
    <div className="max-w-2xl space-y-5">
      <div className="space-y-2">
        <Label htmlFor="name">이름</Label>
        <Input
          id="name"
          {...register('name', { required: '이름을 입력하세요.' })}
        />
        {errors.name ? (
          <p className="text-destructive text-sm">{errors.name.message}</p>
        ) : null}
      </div>

      <div className="space-y-2">
        <Label htmlFor="slug">슬러그</Label>
        <Input
          id="slug"
          placeholder="예: notice"
          {...register('slug', { required: '슬러그를 입력하세요.' })}
        />
        {errors.slug ? (
          <p className="text-destructive text-sm">{errors.slug.message}</p>
        ) : null}
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">설명</Label>
        <Textarea id="description" rows={3} {...register('description')} />
      </div>

      <div className="space-y-2">
        <Label>상위 카테고리</Label>
        <Controller
          control={control}
          name="parentId"
          defaultValue={null}
          render={({ field }) => (
            <Select
              value={field.value != null ? String(field.value) : 'none'}
              onValueChange={(value) =>
                field.onChange(value === 'none' ? null : Number(value))
              }
            >
              <SelectTrigger className="w-full">
                <SelectValue placeholder="상위 카테고리 선택" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">없음 (최상위)</SelectItem>
                {parentOptions.map((option) => (
                  <SelectItem key={String(option.value)} value={String(option.value)}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        />
      </div>
    </div>
  )
}
