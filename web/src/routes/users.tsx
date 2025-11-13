import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import {
  useUsersListUsers,
  useUsersCreateUser,
  useUsersGetUser,
  useUsersUpdateUser,
  useUsersDeleteUser,
  getUsersListUsersQueryKey,
  getUsersGetUserQueryKey,
} from '../api/generated/users/users'
import type { CreateUserRequest, UpdateUserRequest, User } from '../api/generated/models'

export const Route = createFileRoute('/users')({
  component: Users,
})

function Users() {
  const queryClient = useQueryClient()
  const [selectedUserId, setSelectedUserId] = useState<string | null>(null)
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [isEditMode, setIsEditMode] = useState(false)
  const [formData, setFormData] = useState<CreateUserRequest>({ name: '', email: '' })
  const [editFormData, setEditFormData] = useState<UpdateUserRequest>({ name: '', email: '' })

  // Queries
  const { data: usersList, isLoading: isLoadingList, error: listError } = useUsersListUsers()
  const {
    data: selectedUser,
    isLoading: isLoadingUser,
    error: userError,
  } = useUsersGetUser(selectedUserId || '', {
    query: { enabled: !!selectedUserId },
  })

  // Mutations
  const createUserMutation = useUsersCreateUser({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getUsersListUsersQueryKey() })
        setFormData({ name: '', email: '' })
        setShowCreateForm(false)
      },
    },
  })

  const updateUserMutation = useUsersUpdateUser({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getUsersListUsersQueryKey() })
        if (selectedUserId) {
          queryClient.invalidateQueries({ queryKey: getUsersGetUserQueryKey(selectedUserId) })
        }
        setIsEditMode(false)
      },
    },
  })

  const deleteUserMutation = useUsersDeleteUser({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getUsersListUsersQueryKey() })
        setSelectedUserId(null)
      },
    },
  })

  const handleCreateUser = (e: React.FormEvent) => {
    e.preventDefault()
    createUserMutation.mutate({ data: formData })
  }

  const handleUpdateUser = (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedUserId) return
    updateUserMutation.mutate({ userId: selectedUserId, data: editFormData })
  }

  const handleDeleteUser = (userId: string) => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      deleteUserMutation.mutate({ userId })
    }
  }

  const handleEditClick = (user: User) => {
    setEditFormData({ name: user.name, email: user.email })
    setIsEditMode(true)
  }

  const handleCancelEdit = () => {
    setIsEditMode(false)
    setEditFormData({ name: '', email: '' })
  }

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Users</h1>
          <p className="text-muted-foreground mt-2">Manage your users here.</p>
        </div>
        <button
          onClick={() => setShowCreateForm(!showCreateForm)}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
        >
          {showCreateForm ? 'Cancel' : 'Create User'}
        </button>
      </div>

      {/* Create User Form */}
      {showCreateForm && (
        <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
          <h2 className="text-xl font-semibold mb-4">Create New User</h2>
          <form onSubmit={handleCreateUser} className="space-y-4">
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
                disabled={createUserMutation.isPending}
                className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors disabled:opacity-50"
              >
                {createUserMutation.isPending ? 'Creating...' : 'Create'}
              </button>
              {createUserMutation.isError && (
                <p className="text-red-600 text-sm self-center">
                  Error: {createUserMutation.error?.message || 'Failed to create user'}
                </p>
              )}
              {createUserMutation.isSuccess && (
                <p className="text-green-600 text-sm self-center">User created successfully!</p>
              )}
            </div>
          </form>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* User List */}
        <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
          <h2 className="text-xl font-semibold mb-4">User List</h2>
          {isLoadingList && <p className="text-sm text-muted-foreground">Loading users...</p>}
          {listError && (
            <p className="text-sm text-red-600">
              Error: {listError.message || 'Failed to load users'}
            </p>
          )}
          {usersList?.users && usersList.users.length > 0 ? (
            <div className="space-y-2">
              {usersList.users.map((user: User) => (
                <div
                  key={user.id}
                  className={`p-3 border rounded-md cursor-pointer transition-colors hover:bg-muted ${
                    selectedUserId === user.id ? 'bg-muted border-blue-500' : ''
                  }`}
                  onClick={() => setSelectedUserId(user.id)}
                >
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <p className="font-medium">{user.name}</p>
                      <p className="text-sm text-muted-foreground">{user.email}</p>
                    </div>
                    <button
                      onClick={(e) => {
                        e.stopPropagation()
                        handleDeleteUser(user.id)
                      }}
                      disabled={deleteUserMutation.isPending}
                      className="px-2 py-1 text-xs bg-red-600 text-white rounded hover:bg-red-700 transition-colors disabled:opacity-50"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            !isLoadingList && (
              <p className="text-sm text-muted-foreground">
                No users found. Create one to get started!
              </p>
            )
          )}
          {usersList?.total !== undefined && (
            <p className="text-sm text-muted-foreground mt-4">Total users: {usersList.total}</p>
          )}
        </div>

        {/* User Detail */}
        <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold">User Details</h2>
            {selectedUser && !isEditMode && (
              <button
                onClick={() => handleEditClick(selectedUser)}
                className="px-3 py-1 text-sm bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
              >
                Edit
              </button>
            )}
          </div>
          {!selectedUserId && (
            <p className="text-sm text-muted-foreground">
              Select a user from the list to view details.
            </p>
          )}
          {isLoadingUser && (
            <p className="text-sm text-muted-foreground">Loading user details...</p>
          )}
          {userError && (
            <p className="text-sm text-red-600">
              Error: {userError.message || 'Failed to load user details'}
            </p>
          )}
          {selectedUser && !isEditMode && (
            <div className="space-y-4">
              <div>
                <label className="text-sm font-medium text-muted-foreground">ID</label>
                <p className="font-mono text-sm bg-muted p-2 rounded">{selectedUser.id}</p>
              </div>
              <div>
                <label className="text-sm font-medium text-muted-foreground">Name</label>
                <p className="text-lg font-medium">{selectedUser.name}</p>
              </div>
              <div>
                <label className="text-sm font-medium text-muted-foreground">Email</label>
                <p className="text-lg">{selectedUser.email}</p>
              </div>
              <div>
                <label className="text-sm font-medium text-muted-foreground">Created At</label>
                <p className="text-sm">{new Date(selectedUser.createdAt).toLocaleString()}</p>
              </div>
              <div>
                <label className="text-sm font-medium text-muted-foreground">Updated At</label>
                <p className="text-sm">{new Date(selectedUser.updatedAt).toLocaleString()}</p>
              </div>
              <button
                onClick={() => handleDeleteUser(selectedUser.id)}
                disabled={deleteUserMutation.isPending}
                className="w-full px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors disabled:opacity-50"
              >
                {deleteUserMutation.isPending ? 'Deleting...' : 'Delete User'}
              </button>
            </div>
          )}
          {selectedUser && isEditMode && (
            <form onSubmit={handleUpdateUser} className="space-y-4">
              <div>
                <label className="text-sm font-medium text-muted-foreground">ID</label>
                <p className="font-mono text-sm bg-muted p-2 rounded">{selectedUser.id}</p>
              </div>
              <div>
                <label htmlFor="edit-name" className="block text-sm font-medium mb-1">
                  Name
                </label>
                <input
                  type="text"
                  id="edit-name"
                  value={editFormData.name || ''}
                  onChange={(e) => setEditFormData({ ...editFormData, name: e.target.value })}
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
                  value={editFormData.email || ''}
                  onChange={(e) => setEditFormData({ ...editFormData, email: e.target.value })}
                  className="w-full px-3 py-2 border rounded-md bg-background"
                />
              </div>
              <div>
                <label className="text-sm font-medium text-muted-foreground">Created At</label>
                <p className="text-sm">{new Date(selectedUser.createdAt).toLocaleString()}</p>
              </div>
              <div>
                <label className="text-sm font-medium text-muted-foreground">Updated At</label>
                <p className="text-sm">{new Date(selectedUser.updatedAt).toLocaleString()}</p>
              </div>
              <div className="flex gap-2">
                <button
                  type="submit"
                  disabled={updateUserMutation.isPending}
                  className="flex-1 px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors disabled:opacity-50"
                >
                  {updateUserMutation.isPending ? 'Updating...' : 'Update'}
                </button>
                <button
                  type="button"
                  onClick={handleCancelEdit}
                  disabled={updateUserMutation.isPending}
                  className="flex-1 px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 transition-colors disabled:opacity-50"
                >
                  Cancel
                </button>
              </div>
              {updateUserMutation.isError && (
                <p className="text-red-600 text-sm">
                  Error: {updateUserMutation.error?.message || 'Failed to update user'}
                </p>
              )}
              {updateUserMutation.isSuccess && (
                <p className="text-green-600 text-sm">User updated successfully!</p>
              )}
            </form>
          )}
        </div>
      </div>
    </div>
  )
}
