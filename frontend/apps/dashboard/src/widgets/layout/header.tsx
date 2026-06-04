import { useGetIdentity, useLogout } from '@refinedev/core'
import { LogOut } from 'lucide-react'

import type { AdminUser } from '@/entities/auth'
import { Button } from '@/shared/ui'

export function Header() {
  const { data: user } = useGetIdentity<AdminUser>()
  const { mutate: logout } = useLogout()

  return (
    <header className="flex h-14 items-center justify-end gap-4 border-b px-6">
      {user ? (
        <div className="text-sm">
          <span className="font-medium">{user.name}</span>
          <span className="text-muted-foreground ml-2">{user.role}</span>
        </div>
      ) : null}
      <Button variant="outline" size="sm" onClick={() => logout()}>
        <LogOut className="size-4" /> 로그아웃
      </Button>
    </header>
  )
}
