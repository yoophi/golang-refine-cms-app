import { useList, useSelect } from '@refinedev/core'
import {
  Controller,
  type Control,
  type FieldErrors,
  type UseFormRegister,
} from 'react-hook-form'

import type { Category } from '@/entities/category'
import { POST_STATUSES, POST_STATUS_LABEL } from '@/entities/post'
import type { Tag } from '@/entities/tag'
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

import type { PostFormValues } from './post-form-values'

interface PostFieldsProps {
  register: UseFormRegister<PostFormValues>
  control: Control<PostFormValues>
  errors: FieldErrors<PostFormValues>
}

export function PostFields({ register, control, errors }: PostFieldsProps) {
  const { options: categoryOptions } = useSelect<Category>({
    resource: 'categories',
    optionLabel: 'name',
    optionValue: 'id',
  })

  const { result: tagsResult } = useList<Tag>({
    resource: 'tags',
    pagination: { pageSize: 100 },
  })
  const tags = tagsResult?.data ?? []

  return (
    <div className="max-w-2xl space-y-5">
      <div className="space-y-2">
        <Label htmlFor="title">제목</Label>
        <Input
          id="title"
          {...register('title', { required: '제목을 입력하세요.' })}
        />
        {errors.title ? (
          <p className="text-destructive text-sm">{errors.title.message}</p>
        ) : null}
      </div>

      <div className="space-y-2">
        <Label htmlFor="slug">슬러그</Label>
        <Input
          id="slug"
          placeholder="예: hello-world"
          {...register('slug', { required: '슬러그를 입력하세요.' })}
        />
        {errors.slug ? (
          <p className="text-destructive text-sm">{errors.slug.message}</p>
        ) : null}
      </div>

      <div className="space-y-2">
        <Label htmlFor="excerpt">요약</Label>
        <Textarea id="excerpt" rows={2} {...register('excerpt')} />
      </div>

      <div className="space-y-2">
        <Label htmlFor="content">본문</Label>
        <Textarea
          id="content"
          rows={8}
          {...register('content', { required: '본문을 입력하세요.' })}
        />
        {errors.content ? (
          <p className="text-destructive text-sm">{errors.content.message}</p>
        ) : null}
      </div>

      <div className="space-y-2">
        <Label>상태</Label>
        <Controller
          control={control}
          name="status"
          defaultValue="draft"
          render={({ field }) => (
            <Select value={field.value} onValueChange={field.onChange}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="상태 선택" />
              </SelectTrigger>
              <SelectContent>
                {POST_STATUSES.map((status) => (
                  <SelectItem key={status} value={status}>
                    {POST_STATUS_LABEL[status]}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        />
      </div>

      <div className="space-y-2">
        <Label>카테고리</Label>
        <Controller
          control={control}
          name="categoryId"
          defaultValue={null}
          render={({ field }) => (
            <Select
              value={field.value != null ? String(field.value) : 'none'}
              onValueChange={(value) =>
                field.onChange(value === 'none' ? null : Number(value))
              }
            >
              <SelectTrigger className="w-full">
                <SelectValue placeholder="카테고리 선택" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">없음</SelectItem>
                {categoryOptions.map((option) => (
                  <SelectItem key={String(option.value)} value={String(option.value)}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        />
      </div>

      <div className="space-y-2">
        <Label>태그</Label>
        <Controller
          control={control}
          name="tagIds"
          defaultValue={[]}
          render={({ field }) => {
            const selected = field.value ?? []
            return (
              <div className="flex flex-wrap gap-3">
                {tags.length === 0 ? (
                  <span className="text-muted-foreground text-sm">
                    등록된 태그가 없습니다.
                  </span>
                ) : (
                  tags.map((tag) => (
                    <label
                      key={tag.id}
                      className="flex items-center gap-2 text-sm"
                    >
                      <input
                        type="checkbox"
                        className="size-4"
                        checked={selected.includes(tag.id)}
                        onChange={(event) => {
                          const next = new Set<number>(selected)
                          if (event.target.checked) {
                            next.add(tag.id)
                          } else {
                            next.delete(tag.id)
                          }
                          field.onChange(Array.from(next))
                        }}
                      />
                      {tag.name}
                    </label>
                  ))
                )}
              </div>
            )
          }}
        />
      </div>
    </div>
  )
}
