import { useDelete, useNavigation } from '@refinedev/core'
import { Eye, Pencil, Trash2 } from 'lucide-react'

import { Button } from '@/shared/ui'

interface RowActionsProps {
  resource: string
  id: number | string
  canDelete?: boolean
}

export function RowActions({ resource, id, canDelete = true }: RowActionsProps) {
  const { show, edit } = useNavigation()
  const { mutate: deleteOne } = useDelete()

  return (
    <div className="flex justify-end gap-1">
      <Button
        variant="ghost"
        size="icon"
        aria-label="상세"
        onClick={() => show(resource, id)}
      >
        <Eye className="size-4" />
      </Button>
      <Button
        variant="ghost"
        size="icon"
        aria-label="수정"
        onClick={() => edit(resource, id)}
      >
        <Pencil className="size-4" />
      </Button>
      {canDelete ? (
        <Button
          variant="ghost"
          size="icon"
          aria-label="삭제"
          onClick={() => {
            if (window.confirm('정말 삭제하시겠습니까?')) {
              deleteOne({ resource, id })
            }
          }}
        >
          <Trash2 className="text-destructive size-4" />
        </Button>
      ) : null}
    </div>
  )
}
