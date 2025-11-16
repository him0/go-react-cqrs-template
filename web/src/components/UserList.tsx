import type { User } from '../api/generated/models'
import { UserListItem } from './UserListItem'

export interface UserListProps {
  users: User[] | undefined
  total: number | undefined
  isLoading: boolean
  error: { message?: string } | null
  selectedUserId: string | null
  onSelectUser: (userId: string) => void
  onDeleteUser: (userId: string) => void
  isDeleting: boolean
}

export function UserList({
  users,
  total,
  isLoading,
  error,
  selectedUserId,
  onSelectUser,
  onDeleteUser,
  isDeleting,
}: UserListProps) {
  return (
    <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
      <h2 className="text-xl font-semibold mb-4">User List</h2>
      {isLoading && <p className="text-sm text-muted-foreground">Loading users...</p>}
      {error && (
        <p className="text-sm text-red-600">Error: {error.message || 'Failed to load users'}</p>
      )}
      {users && users.length > 0 ? (
        <div className="space-y-2">
          {users.map((user: User) => (
            <UserListItem
              key={user.id}
              user={user}
              isSelected={selectedUserId === user.id}
              onSelect={onSelectUser}
              onDelete={onDeleteUser}
              isDeleting={isDeleting}
            />
          ))}
        </div>
      ) : (
        !isLoading && (
          <p className="text-sm text-muted-foreground">
            No users found. Create one to get started!
          </p>
        )
      )}
      {total !== undefined && (
        <p className="text-sm text-muted-foreground mt-4">Total users: {total}</p>
      )}
    </div>
  )
}
