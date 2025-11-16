import { useState } from 'react'
import type { User, UpdateUserRequest } from '../api/generated/models'

export interface UserEditFormProps {
  user: User
  onSubmit: (userId: string, data: UpdateUserRequest) => void
  onCancel: () => void
  isPending: boolean
  isError: boolean
  isSuccess: boolean
  errorMessage?: string
}

export function UserEditForm({
  user,
  onSubmit,
  onCancel,
  isPending,
  isError,
  isSuccess,
  errorMessage,
}: UserEditFormProps) {
  const [formData, setFormData] = useState<UpdateUserRequest>({
    name: user.name,
    email: user.email,
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(user.id, formData)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label className="text-sm font-medium text-muted-foreground">ID</label>
        <p className="font-mono text-sm bg-muted p-2 rounded">{user.id}</p>
      </div>
      <div>
        <label htmlFor="edit-name" className="block text-sm font-medium mb-1">
          Name
        </label>
        <input
          type="text"
          id="edit-name"
          value={formData.name || ''}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          className="w-full px-3 py-2 border rounded-md bg-background"
          minLength={1}
          maxLength={100}
        />
      </div>
      <div>
        <label htmlFor="edit-email" className="block text-sm font-medium mb-1">
          Email
        </label>
        <input
          type="email"
          id="edit-email"
          value={formData.email || ''}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          className="w-full px-3 py-2 border rounded-md bg-background"
        />
      </div>
      <div>
        <label className="text-sm font-medium text-muted-foreground">Created At</label>
        <p className="text-sm">{new Date(user.createdAt).toLocaleString()}</p>
      </div>
      <div>
        <label className="text-sm font-medium text-muted-foreground">Updated At</label>
        <p className="text-sm">{new Date(user.updatedAt).toLocaleString()}</p>
      </div>
      <div className="flex gap-2">
        <button
          type="submit"
          disabled={isPending}
          className="flex-1 px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors disabled:opacity-50"
        >
          {isPending ? 'Updating...' : 'Update'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          disabled={isPending}
          className="flex-1 px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 transition-colors disabled:opacity-50"
        >
          Cancel
        </button>
      </div>
      {isError && (
        <p className="text-red-600 text-sm">Error: {errorMessage || 'Failed to update user'}</p>
      )}
      {isSuccess && <p className="text-green-600 text-sm">User updated successfully!</p>}
    </form>
  )
}
