import { CanAccess, useMenu } from '@refinedev/core'
import { Link } from 'react-router'
import { LayoutDashboard } from 'lucide-react'

import { cn } from '@/shared/lib/utils'

export function Sidebar() {
  const { menuItems, selectedKey } = useMenu()

  return (
    <aside className="bg-sidebar text-sidebar-foreground border-sidebar-border flex w-60 shrink-0 flex-col border-r">
      <div className="flex h-14 items-center gap-2 border-b px-5 font-semibold">
        <LayoutDashboard className="size-5" />
        <span>CMS 대시보드</span>
      </div>
      <nav className="flex-1 space-y-1 p-3">
        {menuItems.map((item) => {
          const active = item.key === selectedKey
          return (
            <CanAccess key={item.key} resource={item.name} action="list">
              <Link
                to={item.route ?? '/'}
                className={cn(
                  'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors',
                  active
                    ? 'bg-sidebar-accent text-sidebar-accent-foreground'
                    : 'text-sidebar-foreground/70 hover:bg-sidebar-accent/60 hover:text-sidebar-accent-foreground',
                )}
              >
                {item.icon}
                <span>{item.label ?? item.name}</span>
              </Link>
            </CanAccess>
          )
        })}
      </nav>
      <div className="text-muted-foreground border-t p-4 text-xs">
        API: <code>localhost:9000</code>
      </div>
    </aside>
  )
}
