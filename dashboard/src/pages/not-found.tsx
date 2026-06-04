import { Link } from 'react-router'

import { Button, PageHeader } from '@/shared/ui'

export function NotFound() {
  return (
    <div>
      <PageHeader title="404" description="페이지를 찾을 수 없습니다." />
      <Button asChild>
        <Link to="/">홈으로</Link>
      </Button>
    </div>
  )
}
