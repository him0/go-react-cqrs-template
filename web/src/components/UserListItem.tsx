import type { User } from '../api/generated/models'

export interface UserListItemProps {
  user: User
  isSelected: boolean
  onSelect: (userId: string) => void
  onDelete: (userId: string) => void
  isDeleting: boolean
}

export function UserListItem({
  user,
  isSelected,
  onSelect,
  onDelete,
  isDeleting,
}: UserListItemProps) {
  return (
    <div
      className={`p-3 border rounded-md cursor-pointer transition-colors hover:bg-muted ${
        isSelected ? 'bg-muted border-blue-500' : ''
      }`}
      onClick={() => onSelect(user.id)}
    >
      <div className="flex justify-between items-start">
        <div className="flex-1">
          <p className="font-medium">{user.name}</p>
          <p className="text-sm text-muted-foreground">{user.email}</p>
        </div>
        <button
          onClick={(e) => {
            e.stopPropagation()
            onDelete(user.id)
          }}
          disabled={isDeleting}
          className="px-2 py-1 text-xs bg-red-600 text-white rounded hover:bg-red-700 transition-colors disabled:opacity-50"
        >
          Delete
        </button>
      </div>
    </div>
  )
}
