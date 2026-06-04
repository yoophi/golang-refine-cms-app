import type { PropsWithChildren } from 'react'

import { Header } from './header'
import { Sidebar } from './sidebar'

export function Layout({ children }: PropsWithChildren) {
  return (
    <div className="bg-background flex min-h-screen">
      <Sidebar />
      <div className="flex flex-1 flex-col">
        <Header />
        <main className="mx-auto w-full max-w-6xl flex-1 p-6">{children}</main>
      </div>
    </div>
  )
}
