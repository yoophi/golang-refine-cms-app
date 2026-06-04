import { Link } from 'react-router-dom'

import { Button } from '@/shared/ui'

export function NotFoundPage() {
  return (
    <div className="space-y-4 py-16 text-center">
      <h1 className="text-2xl font-semibold">404</h1>
      <p className="text-muted-foreground">페이지를 찾을 수 없습니다.</p>
      <Button asChild>
        <Link to="/">홈으로</Link>
      </Button>
    </div>
  )
}
