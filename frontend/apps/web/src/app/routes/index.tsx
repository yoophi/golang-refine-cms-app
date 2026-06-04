import { Route, Routes } from 'react-router-dom'

import { NotFoundPage, PostDetailPage, PostListPage } from '@/pages'
import { Layout } from '@/widgets/layout'

export function AppRoutes() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route index element={<PostListPage />} />
        <Route path="posts/:id" element={<PostDetailPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Route>
    </Routes>
  )
}
