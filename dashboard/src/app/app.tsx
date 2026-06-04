import { BrowserRouter } from 'react-router'
import {
  DocumentTitleHandler,
  UnsavedChangesNotifier,
} from '@refinedev/react-router'

import { Toaster } from '@/shared/ui'

import { AppRefineProvider } from './providers/refine'
import { AppRoutes } from './routes'

export function App() {
  return (
    <BrowserRouter>
      <AppRefineProvider>
        <AppRoutes />
        <UnsavedChangesNotifier />
        <DocumentTitleHandler />
        <Toaster richColors position="top-right" />
      </AppRefineProvider>
    </BrowserRouter>
  )
}
