import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/about')({
  component: About,
})

function About() {
  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">About</h1>
        <p className="text-muted-foreground mt-2">
          Learn more about this project
        </p>
      </div>

      <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6 space-y-4">
        <div>
          <h2 className="text-xl font-semibold mb-2">Project Overview</h2>
          <p className="text-sm text-muted-foreground">
            This is a full-stack sample application demonstrating modern web development practices
            with Go backend and React frontend.
          </p>
        </div>

        <div>
          <h3 className="font-semibold mb-2">Key Features</h3>
          <ul className="list-disc list-inside space-y-1 text-sm text-muted-foreground">
            <li>Domain-Driven Design (DDD) architecture on the backend</li>
            <li>Type-safe API with OpenAPI specification</li>
            <li>Automatic code generation for both frontend and backend</li>
            <li>Modern React with TanStack Router and Query</li>
            <li>Beautiful UI with Tailwind CSS and shadcn/ui</li>
            <li>Full TypeScript support</li>
          </ul>
        </div>

        <div>
          <h3 className="font-semibold mb-2">Technologies Used</h3>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p className="font-medium mb-1">Backend</p>
              <ul className="space-y-1 text-muted-foreground">
                <li>Go 1.21+</li>
                <li>Chi Router</li>
                <li>OpenAPI 3.0</li>
              </ul>
            </div>
            <div>
              <p className="font-medium mb-1">Frontend</p>
              <ul className="space-y-1 text-muted-foreground">
                <li>React 18</li>
                <li>TypeScript</li>
                <li>Vite</li>
                <li>TanStack Router</li>
                <li>TanStack Query</li>
                <li>Tailwind CSS</li>
                <li>shadcn/ui</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
