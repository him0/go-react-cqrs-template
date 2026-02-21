import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen } from '@/test/utils'
import { ErrorBoundary } from '@/components/ErrorBoundary'

function ThrowingComponent({ message }: { message: string }): never {
  throw new Error(message)
}

function GoodComponent() {
  return <div>All good</div>
}

describe('ErrorBoundary', () => {
  const originalConsoleError = console.error

  beforeEach(() => {
    console.error = vi.fn()
  })

  afterEach(() => {
    console.error = originalConsoleError
  })

  it('should render children when no error occurs', () => {
    render(
      <ErrorBoundary>
        <GoodComponent />
      </ErrorBoundary>
    )

    expect(screen.getByText('All good')).toBeInTheDocument()
  })

  it('should render fallback UI when a child throws an error', () => {
    render(
      <ErrorBoundary>
        <ThrowingComponent message="Test error" />
      </ErrorBoundary>
    )

    expect(screen.getByText('Something went wrong')).toBeInTheDocument()
    expect(screen.getByText('Test error')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /reload page/i })).toBeInTheDocument()
  })

  it('should display error stack in development mode', () => {
    render(
      <ErrorBoundary>
        <ThrowingComponent message="Stack trace test" />
      </ErrorBoundary>
    )

    // In test environment, import.meta.env.DEV is true
    const preElement = screen.getByText(/Stack trace test/i, { selector: 'pre' })
    expect(preElement).toBeInTheDocument()
  })

  it('should log the error via componentDidCatch', () => {
    render(
      <ErrorBoundary>
        <ThrowingComponent message="Catch test" />
      </ErrorBoundary>
    )

    expect(console.error).toHaveBeenCalledWith(
      'ErrorBoundary caught an error:',
      expect.any(Error),
      expect.objectContaining({ componentStack: expect.any(String) })
    )
  })
})
