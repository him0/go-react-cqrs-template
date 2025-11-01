import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: Index,
})

function Index() {
  return (
    <div className="max-w-3xl mx-auto space-y-8">
      <div className="text-center space-y-4">
        <h1 className="text-4xl font-bold tracking-tight">
          Welcome to User Management App
        </h1>
        <p className="text-muted-foreground text-lg">
          A sample application using Go (DDD), Vite, React, OpenAPI, and Orval
        </p>
      </div>

      <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
        <h2 className="text-2xl font-semibold mb-4">Tech Stack</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <h3 className="font-medium text-lg">Backend</h3>
            <ul className="space-y-1 text-sm text-muted-foreground">
              <li>• Go with DDD architecture</li>
              <li>• Chi router</li>
              <li>• OpenAPI 3.0 specification</li>
            </ul>
          </div>
          <div className="space-y-2">
            <h3 className="font-medium text-lg">Frontend</h3>
            <ul className="space-y-1 text-sm text-muted-foreground">
              <li>• React 18 + TypeScript</li>
              <li>• Vite build tool</li>
              <li>• TanStack Router + Query</li>
              <li>• Tailwind CSS + shadcn/ui</li>
              <li>• React Query hooks via Orval</li>
            </ul>
          </div>
        </div>
      </div>

      <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
        <h2 className="text-2xl font-semibold mb-4">Getting Started</h2>
        <div className="space-y-3 text-sm">
          <div>
            <code className="bg-muted px-2 py-1 rounded text-xs">
              make install
            </code>
            <p className="text-muted-foreground mt-1">Install dependencies</p>
          </div>
          <div>
            <code className="bg-muted px-2 py-1 rounded text-xs">
              cd web && pnpm run generate:api
            </code>
            <p className="text-muted-foreground mt-1">Generate API code with Orval</p>
          </div>
          <div>
            <code className="bg-muted px-2 py-1 rounded text-xs">
              make run-backend
            </code>
            <p className="text-muted-foreground mt-1">Start backend server (port 8080)</p>
          </div>
          <div>
            <code className="bg-muted px-2 py-1 rounded text-xs">
              make run-frontend
            </code>
            <p className="text-muted-foreground mt-1">Start frontend dev server (port 3000)</p>
          </div>
        </div>
      </div>
    </div>
  )
}
