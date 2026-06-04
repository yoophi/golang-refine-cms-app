import type { FieldErrors, UseFormRegister } from 'react-hook-form'

import { Input, Label } from '@/shared/ui'

export interface TagFormValues {
  name: string
  slug: string
}

interface TagFieldsProps {
  register: UseFormRegister<TagFormValues>
  errors: FieldErrors<TagFormValues>
}

export function TagFields({ register, errors }: TagFieldsProps) {
  return (
    <div className="max-w-lg space-y-5">
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
          placeholder="예: react"
          {...register('slug', { required: '슬러그를 입력하세요.' })}
        />
        {errors.slug ? (
          <p className="text-destructive text-sm">{errors.slug.message}</p>
        ) : null}
      </div>
    </div>
  )
}
