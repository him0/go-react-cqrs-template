import { useState } from 'react'
import type { CreateUserRequest } from '../api/generated/models'

export interface UserCreateFormProps {
  onSubmit: (data: CreateUserRequest) => void
  isPending: boolean
  isError: boolean
  isSuccess: boolean
  errorMessage?: string
}

export function UserCreateForm({
  onSubmit,
  isPending,
  isError,
  isSuccess,
  errorMessage,
}: UserCreateFormProps) {
  const [formData, setFormData] = useState<CreateUserRequest>({ name: '', email: '' })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(formData)
  }

  return (
    <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
      <h2 className="text-xl font-semibold mb-4">Create New User</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="name" className="block text-sm font-medium mb-1">
            Name
          </label>
          <input
            type="text"
            id="name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            className="w-full px-3 py-2 border rounded-md bg-background"
            required
            minLength={1}
            maxLength={100}
          />
        </div>
        <div>
          <label htmlFor="email" className="block text-sm font-medium mb-1">
            Email
          </label>
          <input
            type="email"
            id="email"
            value={formData.email}
            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
            className="w-full px-3 py-2 border rounded-md bg-background"
            required
          />
        </div>
        <div className="flex gap-2">
          <button
            type="submit"
            disabled={isPending}
            className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors disabled:opacity-50"
          >
            {isPending ? 'Creating...' : 'Create'}
          </button>
          {isError && (
            <p className="text-red-600 text-sm self-center">
              Error: {errorMessage || 'Failed to create user'}
            </p>
          )}
          {isSuccess && (
            <p className="text-green-600 text-sm self-center">User created successfully!</p>
          )}
        </div>
      </form>
    </div>
  )
}
