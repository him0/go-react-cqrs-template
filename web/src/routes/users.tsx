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
import type { CreateUserRequest, UpdateUserRequest } from '../api/generated/models'
import { UserList } from '../components/UserList'
import { UserCreateForm } from '../components/UserCreateForm'
import { UserDetail } from '../components/UserDetail'
import { UserEditForm } from '../components/UserEditForm'

export const Route = createFileRoute('/users')({
  component: Users,
})

function Users() {
  const queryClient = useQueryClient()
  const [selectedUserId, setSelectedUserId] = useState<string | null>(null)
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [isEditMode, setIsEditMode] = useState(false)

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

  const handleCreateUser = (data: CreateUserRequest) => {
    createUserMutation.mutate({ data })
  }

  const handleUpdateUser = (userId: string, data: UpdateUserRequest) => {
    updateUserMutation.mutate({ userId, data })
  }

  const handleDeleteUser = (userId: string) => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      deleteUserMutation.mutate({ userId })
    }
  }

  const handleEditClick = () => {
    setIsEditMode(true)
  }

  const handleCancelEdit = () => {
    setIsEditMode(false)
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

      {showCreateForm && (
        <UserCreateForm
          onSubmit={handleCreateUser}
          isPending={createUserMutation.isPending}
          isError={createUserMutation.isError}
          isSuccess={createUserMutation.isSuccess}
          errorMessage={createUserMutation.error?.message}
        />
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <UserList
          users={usersList?.users}
          total={usersList?.total}
          isLoading={isLoadingList}
          error={listError}
          selectedUserId={selectedUserId}
          onSelectUser={setSelectedUserId}
          onDeleteUser={handleDeleteUser}
          isDeleting={deleteUserMutation.isPending}
        />

        <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold">User Details</h2>
            {selectedUser && !isEditMode && (
              <button
                onClick={handleEditClick}
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
            <UserDetail
              user={selectedUser}
              onEdit={handleEditClick}
              onDelete={handleDeleteUser}
              isDeleting={deleteUserMutation.isPending}
            />
          )}
          {selectedUser && isEditMode && (
            <UserEditForm
              key={selectedUser.id}
              user={selectedUser}
              onSubmit={handleUpdateUser}
              onCancel={handleCancelEdit}
              isPending={updateUserMutation.isPending}
              isError={updateUserMutation.isError}
              isSuccess={updateUserMutation.isSuccess}
              errorMessage={updateUserMutation.error?.message}
            />
          )}
        </div>
      </div>
    </div>
  )
}
