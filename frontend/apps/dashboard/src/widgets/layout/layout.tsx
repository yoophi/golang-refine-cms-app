import type { PropsWithChildren } from 'react'

import { Sidebar } from './sidebar'

export function Layout({ children }: PropsWithChildren) {
  return (
    <div className="bg-background flex min-h-screen">
      <Sidebar />
      <main className="mx-auto w-full max-w-6xl flex-1 p-6">{children}</main>
    </div>
  )
}
