import type { PropsWithChildren } from 'react'
import { Refine, type ResourceProps } from '@refinedev/core'
import routerProvider from '@refinedev/react-router'
import dataProvider from '@refinedev/simple-rest'
import { FileText, Folder, MessageSquare, Tags } from 'lucide-react'

import { accessControlProvider } from '@/shared/api/access-control-provider'
import { authProvider } from '@/shared/api/auth-provider'
import { httpClient } from '@/shared/api/http-client'
import { notificationProvider } from '@/shared/api/notification-provider'
import { API_URL } from '@/shared/config/api'

export const resources: ResourceProps[] = [
  {
    name: 'posts',
    list: '/posts',
    create: '/posts/create',
    edit: '/posts/edit/:id',
    show: '/posts/show/:id',
    meta: { label: '게시글', icon: <FileText className="size-4" />, canDelete: true },
  },
  {
    name: 'categories',
    list: '/categories',
    create: '/categories/create',
    edit: '/categories/edit/:id',
    show: '/categories/show/:id',
    meta: { label: '카테고리', icon: <Folder className="size-4" />, canDelete: true },
  },
  {
    name: 'tags',
    list: '/tags',
    create: '/tags/create',
    edit: '/tags/edit/:id',
    show: '/tags/show/:id',
    meta: { label: '태그', icon: <Tags className="size-4" />, canDelete: true },
  },
  {
    name: 'comments',
    list: '/comments',
    create: '/comments/create',
    edit: '/comments/edit/:id',
    show: '/comments/show/:id',
    meta: { label: '댓글', icon: <MessageSquare className="size-4" />, canDelete: true },
  },
]

export function AppRefineProvider({ children }: PropsWithChildren) {
  return (
    <Refine
      dataProvider={dataProvider(API_URL, httpClient)}
      routerProvider={routerProvider}
      authProvider={authProvider}
      accessControlProvider={accessControlProvider}
      notificationProvider={notificationProvider}
      resources={resources}
      options={{
        syncWithLocation: true,
        warnWhenUnsavedChanges: true,
        disableTelemetry: true,
        title: { text: 'CMS 대시보드' },
      }}
    >
      {children}
    </Refine>
  )
}
