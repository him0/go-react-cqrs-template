import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import type { CreateUserRequest } from '../api/generated/models'
import { createUserSchema, type CreateUserFormValues } from '../lib/validations/user'

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
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateUserFormValues>({
    resolver: zodResolver(createUserSchema),
    defaultValues: {
      name: '',
      email: '',
    },
  })

  const onFormSubmit = (data: CreateUserFormValues) => {
    onSubmit(data)
  }

  return (
    <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
      <h2 className="text-xl font-semibold mb-4">Create New User</h2>
      <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-4">
        <div>
          <label htmlFor="name" className="block text-sm font-medium mb-1">
            Name
          </label>
          <input
            type="text"
            id="name"
            {...register('name')}
            className="w-full px-3 py-2 border rounded-md bg-background"
            aria-required="true"
          />
          {errors.name && (
            <p className="text-red-600 text-sm mt-1" role="alert">
              {errors.name.message}
            </p>
          )}
        </div>
        <div>
          <label htmlFor="email" className="block text-sm font-medium mb-1">
            Email
          </label>
          <input
            type="text"
            id="email"
            {...register('email')}
            className="w-full px-3 py-2 border rounded-md bg-background"
            aria-required="true"
          />
          {errors.email && (
            <p className="text-red-600 text-sm mt-1" role="alert">
              {errors.email.message}
            </p>
          )}
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
