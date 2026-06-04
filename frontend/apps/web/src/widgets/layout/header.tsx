import { Link } from 'react-router-dom'

import { useAuthStore, useLogout } from '@/entities/user'
import { Button } from '@/shared/ui'

export function Header() {
  const user = useAuthStore((s) => s.user)
  const onLogout = useLogout()

  return (
    <header className="border-b">
      <div className="mx-auto flex h-14 max-w-3xl items-center justify-between px-4">
        <Link to="/" className="text-lg font-semibold tracking-tight">
          CMS 블로그
        </Link>
        <nav className="flex items-center gap-3 text-sm">
          {user ? (
            <>
              <Link to="/profile" className="hover:underline">
                {user.name}
              </Link>
              <Button variant="outline" size="sm" onClick={onLogout}>
                로그아웃
              </Button>
            </>
          ) : (
            <>
              <Link to="/login" className="hover:underline">
                로그인
              </Link>
              <Button asChild size="sm">
                <Link to="/register">회원가입</Link>
              </Button>
            </>
          )}
        </nav>
      </div>
    </header>
  )
}
