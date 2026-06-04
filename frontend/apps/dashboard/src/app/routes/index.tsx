import { Authenticated } from '@refinedev/core'
import { CatchAllNavigate, NavigateToResource } from '@refinedev/react-router'
import { Outlet, Route, Routes } from 'react-router'

import {
  CategoryCreate,
  CategoryEdit,
  CategoryList,
  CategoryShow,
  CommentCreate,
  CommentEdit,
  CommentList,
  CommentShow,
  LoginPage,
  NotFound,
  PostCreate,
  PostEdit,
  PostList,
  PostShow,
  TagCreate,
  TagEdit,
  TagList,
  TagShow,
} from '@/pages'
import { Layout } from '@/widgets/layout'

export function AppRoutes() {
  return (
    <Routes>
      <Route
        element={
          <Authenticated
            key="authenticated-routes"
            fallback={<CatchAllNavigate to="/login" />}
          >
            <Layout>
              <Outlet />
            </Layout>
          </Authenticated>
        }
      >
        <Route index element={<NavigateToResource resource="posts" />} />

        <Route path="posts">
          <Route index element={<PostList />} />
          <Route path="create" element={<PostCreate />} />
          <Route path="edit/:id" element={<PostEdit />} />
          <Route path="show/:id" element={<PostShow />} />
        </Route>

        <Route path="categories">
          <Route index element={<CategoryList />} />
          <Route path="create" element={<CategoryCreate />} />
          <Route path="edit/:id" element={<CategoryEdit />} />
          <Route path="show/:id" element={<CategoryShow />} />
        </Route>

        <Route path="tags">
          <Route index element={<TagList />} />
          <Route path="create" element={<TagCreate />} />
          <Route path="edit/:id" element={<TagEdit />} />
          <Route path="show/:id" element={<TagShow />} />
        </Route>

        <Route path="comments">
          <Route index element={<CommentList />} />
          <Route path="create" element={<CommentCreate />} />
          <Route path="edit/:id" element={<CommentEdit />} />
          <Route path="show/:id" element={<CommentShow />} />
        </Route>

        <Route path="*" element={<NotFound />} />
      </Route>

      <Route
        element={
          <Authenticated key="auth-pages" fallback={<Outlet />}>
            <NavigateToResource resource="posts" />
          </Authenticated>
        }
      >
        <Route path="/login" element={<LoginPage />} />
      </Route>
    </Routes>
  )
}
