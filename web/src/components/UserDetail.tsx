import type { User } from '../api/generated/models'

export interface UserDetailProps {
  user: User
  onEdit: () => void
  onDelete: (userId: string) => void
  isDeleting: boolean
}

export function UserDetail({ user, onEdit, onDelete, isDeleting }: UserDetailProps) {
  return (
    <div className="space-y-4">
      <div>
        <label className="text-sm font-medium text-muted-foreground">ID</label>
        <p className="font-mono text-sm bg-muted p-2 rounded">{user.id}</p>
      </div>
      <div>
        <label className="text-sm font-medium text-muted-foreground">Name</label>
        <p className="text-lg font-medium">{user.name}</p>
      </div>
      <div>
        <label className="text-sm font-medium text-muted-foreground">Email</label>
        <p className="text-lg">{user.email}</p>
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
          onClick={onEdit}
          className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
        >
          Edit
        </button>
        <button
          onClick={() => onDelete(user.id)}
          disabled={isDeleting}
          className="flex-1 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors disabled:opacity-50"
        >
          {isDeleting ? 'Deleting...' : 'Delete User'}
        </button>
      </div>
    </div>
  )
}
