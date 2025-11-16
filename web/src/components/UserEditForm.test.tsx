import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@/test/utils'
import userEvent from '@testing-library/user-event'
import { UserEditForm } from './UserEditForm'
import type { User } from '../api/generated/models'

const mockUser: User = {
  id: '123',
  name: 'John Doe',
  email: 'john@example.com',
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-02T00:00:00Z',
}

describe('UserEditForm', () => {
  it('should render form with user data', () => {
    const onSubmit = vi.fn()
    const onCancel = vi.fn()

    render(
      <UserEditForm
        user={mockUser}
        onSubmit={onSubmit}
        onCancel={onCancel}
        isPending={false}
        isError={false}
        isSuccess={false}
      />
    )

    expect(screen.getByText('123')).toBeInTheDocument()
    expect(screen.getByDisplayValue('John Doe')).toBeInTheDocument()
    expect(screen.getByDisplayValue('john@example.com')).toBeInTheDocument()
  })

  it('should handle form submission', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    const onCancel = vi.fn()

    render(
      <UserEditForm
        user={mockUser}
        onSubmit={onSubmit}
        onCancel={onCancel}
        isPending={false}
        isError={false}
        isSuccess={false}
      />
    )

    const nameInput = screen.getByLabelText(/name/i)
    const emailInput = screen.getByLabelText(/email/i)

    await user.clear(nameInput)
    await user.type(nameInput, 'Jane Smith')
    await user.clear(emailInput)
    await user.type(emailInput, 'jane@example.com')

    await user.click(screen.getByRole('button', { name: /^update$/i }))

    expect(onSubmit).toHaveBeenCalledWith('123', {
      name: 'Jane Smith',
      email: 'jane@example.com',
    })
  })

  it('should call onCancel when cancel button is clicked', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    const onCancel = vi.fn()

    render(
      <UserEditForm
        user={mockUser}
        onSubmit={onSubmit}
        onCancel={onCancel}
        isPending={false}
        isError={false}
        isSuccess={false}
      />
    )

    await user.click(screen.getByRole('button', { name: /cancel/i }))

    expect(onCancel).toHaveBeenCalledTimes(1)
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it('should show pending state', () => {
    const onSubmit = vi.fn()
    const onCancel = vi.fn()

    render(
      <UserEditForm
        user={mockUser}
        onSubmit={onSubmit}
        onCancel={onCancel}
        isPending={true}
        isError={false}
        isSuccess={false}
      />
    )

    const updateButton = screen.getByRole('button', { name: /updating/i })
    const cancelButton = screen.getByRole('button', { name: /cancel/i })

    expect(updateButton).toBeDisabled()
    expect(cancelButton).toBeDisabled()
  })

  it('should show error message', () => {
    const onSubmit = vi.fn()
    const onCancel = vi.fn()

    render(
      <UserEditForm
        user={mockUser}
        onSubmit={onSubmit}
        onCancel={onCancel}
        isPending={false}
        isError={true}
        isSuccess={false}
        errorMessage="Network error"
      />
    )

    expect(screen.getByText(/error: network error/i)).toBeInTheDocument()
  })

  it('should show default error message when no errorMessage provided', () => {
    const onSubmit = vi.fn()
    const onCancel = vi.fn()

    render(
      <UserEditForm
        user={mockUser}
        onSubmit={onSubmit}
        onCancel={onCancel}
        isPending={false}
        isError={true}
        isSuccess={false}
      />
    )

    expect(screen.getByText(/failed to update user/i)).toBeInTheDocument()
  })

  it('should show success message', () => {
    const onSubmit = vi.fn()
    const onCancel = vi.fn()

    render(
      <UserEditForm
        user={mockUser}
        onSubmit={onSubmit}
        onCancel={onCancel}
        isPending={false}
        isError={false}
        isSuccess={true}
      />
    )

    expect(screen.getByText(/user updated successfully/i)).toBeInTheDocument()
  })

  it('should display read-only fields', () => {
    const onSubmit = vi.fn()
    const onCancel = vi.fn()

    render(
      <UserEditForm
        user={mockUser}
        onSubmit={onSubmit}
        onCancel={onCancel}
        isPending={false}
        isError={false}
        isSuccess={false}
      />
    )

    expect(screen.getByText('ID')).toBeInTheDocument()
    expect(screen.getByText('Created At')).toBeInTheDocument()
    expect(screen.getByText('Updated At')).toBeInTheDocument()
  })

  it('should initialize form with user data', () => {
    const onSubmit = vi.fn()
    const onCancel = vi.fn()

    const updatedUser: User = {
      id: '456',
      name: 'Jane Smith',
      email: 'jane@example.com',
      createdAt: '2024-01-03T00:00:00Z',
      updatedAt: '2024-01-04T00:00:00Z',
    }

    render(
      <UserEditForm
        user={updatedUser}
        onSubmit={onSubmit}
        onCancel={onCancel}
        isPending={false}
        isError={false}
        isSuccess={false}
      />
    )

    expect(screen.getByDisplayValue('Jane Smith')).toBeInTheDocument()
    expect(screen.getByDisplayValue('jane@example.com')).toBeInTheDocument()
  })
})
