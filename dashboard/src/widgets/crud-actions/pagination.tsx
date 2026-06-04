import { Button } from '@/shared/ui'

interface PaginationProps {
  pageIndex: number
  pageCount: number
  canPrevious: boolean
  canNext: boolean
  onPageChange: (pageIndex: number) => void
}

export function Pagination({
  pageIndex,
  pageCount,
  canPrevious,
  canNext,
  onPageChange,
}: PaginationProps) {
  return (
    <div className="mt-4 flex items-center justify-end gap-3 text-sm">
      <span className="text-muted-foreground">
        {pageIndex + 1} / {Math.max(pageCount, 1)} 페이지
      </span>
      <Button
        variant="outline"
        size="sm"
        disabled={!canPrevious}
        onClick={() => onPageChange(pageIndex - 1)}
      >
        이전
      </Button>
      <Button
        variant="outline"
        size="sm"
        disabled={!canNext}
        onClick={() => onPageChange(pageIndex + 1)}
      >
        다음
      </Button>
    </div>
  )
}
