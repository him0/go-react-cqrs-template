import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/users')({
  component: Users,
})

function Users() {
  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Users</h1>
        <p className="text-muted-foreground mt-2">
          Manage your users here. (API integration coming soon)
        </p>
      </div>

      <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
        <h2 className="text-xl font-semibold mb-4">User List</h2>
        <p className="text-sm text-muted-foreground">
          After running <code className="bg-muted px-2 py-1 rounded text-xs">pnpm run generate:api</code>,
          you can use the generated React Query hooks here.
        </p>
        <div className="mt-4 p-4 bg-muted rounded-md">
          <p className="text-sm font-mono">
            Example: useListUsers(), useCreateUser(), etc.
          </p>
        </div>
      </div>
    </div>
  )
}
