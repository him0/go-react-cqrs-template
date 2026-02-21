import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import type { User, UpdateUserRequest } from '../api/generated/models'
import { updateUserSchema, type UpdateUserFormValues } from '../lib/validations/user'

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
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<UpdateUserFormValues>({
    resolver: zodResolver(updateUserSchema),
    defaultValues: {
      name: user.name,
      email: user.email,
    },
  })

  const onFormSubmit = (data: UpdateUserFormValues) => {
    onSubmit(user.id, data)
  }

  return (
    <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-4">
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
          {...register('name')}
          className="w-full px-3 py-2 border rounded-md bg-background"
        />
        {errors.name && (
          <p className="text-red-600 text-sm mt-1" role="alert">
            {errors.name.message}
          </p>
        )}
      </div>
      <div>
        <label htmlFor="edit-email" className="block text-sm font-medium mb-1">
          Email
        </label>
        <input
          type="text"
          id="edit-email"
          {...register('email')}
          className="w-full px-3 py-2 border rounded-md bg-background"
        />
        {errors.email && (
          <p className="text-red-600 text-sm mt-1" role="alert">
            {errors.email.message}
          </p>
        )}
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
