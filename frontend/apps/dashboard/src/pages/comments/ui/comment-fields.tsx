import { useSelect } from '@refinedev/core'
import {
  Controller,
  type Control,
  type FieldErrors,
  type UseFormRegister,
} from 'react-hook-form'

import { COMMENT_STATUSES, COMMENT_STATUS_LABEL } from '@/entities/comment'
import type { Post } from '@/entities/post'
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

import type { CommentFormValues } from './comment-form-values'

interface CommentFieldsProps {
  register: UseFormRegister<CommentFormValues>
  control: Control<CommentFormValues>
  errors: FieldErrors<CommentFormValues>
}

export function CommentFields({ register, control, errors }: CommentFieldsProps) {
  const { options: postOptions } = useSelect<Post>({
    resource: 'posts',
    optionLabel: 'title',
    optionValue: 'id',
  })

  return (
    <div className="max-w-2xl space-y-5">
      <div className="space-y-2">
        <Label>게시글</Label>
        <Controller
          control={control}
          name="postId"
          rules={{ required: '게시글을 선택하세요.' }}
          render={({ field }) => (
            <Select
              value={field.value != null ? String(field.value) : undefined}
              onValueChange={(value) => field.onChange(Number(value))}
            >
              <SelectTrigger className="w-full">
                <SelectValue placeholder="게시글 선택" />
              </SelectTrigger>
              <SelectContent>
                {postOptions.map((option) => (
                  <SelectItem key={String(option.value)} value={String(option.value)}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        />
        {errors.postId ? (
          <p className="text-destructive text-sm">{errors.postId.message}</p>
        ) : null}
      </div>

      <div className="grid gap-5 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="authorName">작성자 이름</Label>
          <Input
            id="authorName"
            {...register('authorName', { required: '작성자 이름을 입력하세요.' })}
          />
          {errors.authorName ? (
            <p className="text-destructive text-sm">
              {errors.authorName.message}
            </p>
          ) : null}
        </div>
        <div className="space-y-2">
          <Label htmlFor="authorEmail">작성자 이메일</Label>
          <Input id="authorEmail" type="email" {...register('authorEmail')} />
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="content">내용</Label>
        <Textarea
          id="content"
          rows={5}
          {...register('content', { required: '내용을 입력하세요.' })}
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
          defaultValue="pending"
          render={({ field }) => (
            <Select value={field.value} onValueChange={field.onChange}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="상태 선택" />
              </SelectTrigger>
              <SelectContent>
                {COMMENT_STATUSES.map((status) => (
                  <SelectItem key={status} value={status}>
                    {COMMENT_STATUS_LABEL[status]}
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
