import { Outlet } from 'react-router-dom'

import { Header } from './header'

export function Layout() {
  return (
    <div className="bg-background text-foreground min-h-screen">
      <Header />
      <main className="mx-auto max-w-3xl px-4 py-8">
        <Outlet />
      </main>
    </div>
  )
}
