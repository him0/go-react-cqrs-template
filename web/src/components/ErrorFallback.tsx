import { Button } from '@/components/Button'

interface ErrorFallbackProps {
  error: Error
  onReset: () => void
}

export function ErrorFallback({ error, onReset }: ErrorFallbackProps) {
  return (
    <div className="flex min-h-[50vh] flex-col items-center justify-center gap-6 px-4 text-center">
      <div className="space-y-2">
        <h2 className="text-2xl font-bold text-destructive">Something went wrong</h2>
        <p className="text-muted-foreground">{error.message || 'An unexpected error occurred.'}</p>
      </div>

      {import.meta.env.DEV && error.stack && (
        <pre className="max-w-2xl overflow-auto rounded-md bg-muted p-4 text-left text-xs text-muted-foreground">
          {error.stack}
        </pre>
      )}

      <Button onClick={onReset}>Reload page</Button>
    </div>
  )
}
