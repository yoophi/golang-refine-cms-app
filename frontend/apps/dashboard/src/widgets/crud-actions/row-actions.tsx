import { CanAccess, useDelete, useNavigation } from '@refinedev/core'
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
      <CanAccess resource={resource} action="show">
        <Button
          variant="ghost"
          size="icon"
          aria-label="상세"
          onClick={() => show(resource, id)}
        >
          <Eye className="size-4" />
        </Button>
      </CanAccess>
      <CanAccess resource={resource} action="edit">
        <Button
          variant="ghost"
          size="icon"
          aria-label="수정"
          onClick={() => edit(resource, id)}
        >
          <Pencil className="size-4" />
        </Button>
      </CanAccess>
      {canDelete ? (
        <CanAccess resource={resource} action="delete">
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
        </CanAccess>
      ) : null}
    </div>
  )
}
