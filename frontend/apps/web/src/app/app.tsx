import { BrowserRouter } from 'react-router-dom'

import { AppQueryProvider } from './providers/query-provider'
import { AppRoutes } from './routes'

export function App() {
  return (
    <AppQueryProvider>
      <BrowserRouter>
        <AppRoutes />
      </BrowserRouter>
    </AppQueryProvider>
  )
}
