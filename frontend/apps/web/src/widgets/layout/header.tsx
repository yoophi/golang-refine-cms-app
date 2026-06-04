import { Link } from 'react-router-dom'

export function Header() {
  return (
    <header className="border-b">
      <div className="mx-auto flex h-14 max-w-3xl items-center px-4">
        <Link to="/" className="text-lg font-semibold tracking-tight">
          CMS 블로그
        </Link>
      </div>
    </header>
  )
}
