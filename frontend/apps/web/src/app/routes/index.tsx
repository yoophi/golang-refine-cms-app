import { Route, Routes } from 'react-router-dom'

import {
  LoginPage,
  NotFoundPage,
  PostDetailPage,
  PostListPage,
  ProfilePage,
  RegisterPage,
} from '@/pages'
import { Layout } from '@/widgets/layout'

import { RequireAuth } from '../require-auth'

export function AppRoutes() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route index element={<PostListPage />} />
        <Route path="posts/:id" element={<PostDetailPage />} />
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
        <Route
          path="profile"
          element={
            <RequireAuth>
              <ProfilePage />
            </RequireAuth>
          }
        />
        <Route path="*" element={<NotFoundPage />} />
      </Route>
    </Routes>
  )
}
