import { render, RenderOptions } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactElement } from 'react'

const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
      mutations: {
        retry: false,
      },
    },
  })

export function renderWithProviders(ui: ReactElement, options?: Omit<RenderOptions, 'wrapper'>) {
  const testQueryClient = createTestQueryClient()

  return render(ui, {
    wrapper: ({ children }) => (
      <QueryClientProvider client={testQueryClient}>{children}</QueryClientProvider>
    ),
    ...options,
  })
}

// eslint-disable-next-line react-refresh/only-export-components
export * from '@testing-library/react'
export { renderWithProviders as render }
