import { createFileRoute } from '@tanstack/react-router'
import { useQueryClient } from '@tanstack/react-query'
import { z } from 'zod'
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

const usersSearchSchema = z.object({
  userId: z.string().optional(),
  showCreate: z.boolean().optional(),
  isEdit: z.boolean().optional(),
})

export const Route = createFileRoute('/users')({
  component: Users,
  validateSearch: usersSearchSchema,
})

function Users() {
  const navigate = Route.useNavigate()
  const search = Route.useSearch()
  const queryClient = useQueryClient()

  const selectedUserId = search.userId ?? null
  const showCreateForm = search.showCreate ?? false
  const isEditMode = search.isEdit ?? false

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
  const {
    mutate: createUser,
    isPending: isCreating,
    isError: isCreateError,
    isSuccess: isCreateSuccess,
    error: createError,
  } = useUsersCreateUser({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getUsersListUsersQueryKey() })
        navigate({ search: { ...search, showCreate: undefined } })
      },
    },
  })

  const {
    mutate: updateUser,
    isPending: isUpdating,
    isError: isUpdateError,
    isSuccess: isUpdateSuccess,
    error: updateError,
  } = useUsersUpdateUser({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getUsersListUsersQueryKey() })
        if (selectedUserId) {
          queryClient.invalidateQueries({ queryKey: getUsersGetUserQueryKey(selectedUserId) })
        }
        navigate({ search: { ...search, isEdit: undefined } })
      },
    },
  })

  const { mutate: deleteUser, isPending: isDeleting } = useUsersDeleteUser({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getUsersListUsersQueryKey() })
        navigate({ search: { ...search, userId: undefined, isEdit: undefined } })
      },
    },
  })

  const handleCreateUser = (data: CreateUserRequest) => {
    createUser({ data })
  }

  const handleUpdateUser = (userId: string, data: UpdateUserRequest) => {
    updateUser({ userId, data })
  }

  const handleDeleteUser = (userId: string) => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      deleteUser({ userId })
    }
  }

  const handleSelectUser = (userId: string | null) => {
    navigate({ search: { ...search, userId: userId ?? undefined, isEdit: undefined } })
  }

  const handleToggleCreateForm = () => {
    navigate({ search: { ...search, showCreate: !showCreateForm } })
  }

  const handleEditClick = () => {
    navigate({ search: { ...search, isEdit: true } })
  }

  const handleCancelEdit = () => {
    navigate({ search: { ...search, isEdit: undefined } })
  }

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Users</h1>
          <p className="text-muted-foreground mt-2">Manage your users here.</p>
        </div>
        <button
          onClick={handleToggleCreateForm}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
        >
          {showCreateForm ? 'Cancel' : 'Create User'}
        </button>
      </div>

      {showCreateForm && (
        <UserCreateForm
          onSubmit={handleCreateUser}
          isPending={isCreating}
          isError={isCreateError}
          isSuccess={isCreateSuccess}
          errorMessage={createError?.message}
        />
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <UserList
          users={usersList?.users}
          total={usersList?.total}
          isLoading={isLoadingList}
          error={listError}
          selectedUserId={selectedUserId}
          onSelectUser={handleSelectUser}
          onDeleteUser={handleDeleteUser}
          isDeleting={isDeleting}
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
              isDeleting={isDeleting}
            />
          )}
          {selectedUser && isEditMode && (
            <UserEditForm
              key={selectedUser.id}
              user={selectedUser}
              onSubmit={handleUpdateUser}
              onCancel={handleCancelEdit}
              isPending={isUpdating}
              isError={isUpdateError}
              isSuccess={isUpdateSuccess}
              errorMessage={updateError?.message}
            />
          )}
        </div>
      </div>
    </div>
  )
}
