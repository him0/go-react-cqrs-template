import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@/test/utils'
import userEvent from '@testing-library/user-event'
import { UserCreateForm } from './UserCreateForm'

describe('UserCreateForm', () => {
  it('should render form fields', () => {
    const onSubmit = vi.fn()

    render(
      <UserCreateForm onSubmit={onSubmit} isPending={false} isError={false} isSuccess={false} />
    )

    expect(screen.getByLabelText(/name/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /create$/i })).toBeInTheDocument()
  })

  it('should handle form submission', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()

    render(
      <UserCreateForm onSubmit={onSubmit} isPending={false} isError={false} isSuccess={false} />
    )

    await user.type(screen.getByLabelText(/name/i), 'John Doe')
    await user.type(screen.getByLabelText(/email/i), 'john@example.com')
    await user.click(screen.getByRole('button', { name: /create$/i }))

    expect(onSubmit).toHaveBeenCalledWith({
      name: 'John Doe',
      email: 'john@example.com',
    })
  })

  it('should show pending state', () => {
    const onSubmit = vi.fn()

    render(
      <UserCreateForm onSubmit={onSubmit} isPending={true} isError={false} isSuccess={false} />
    )

    const button = screen.getByRole('button', { name: /creating/i })
    expect(button).toBeInTheDocument()
    expect(button).toBeDisabled()
  })

  it('should show error message', () => {
    const onSubmit = vi.fn()

    render(
      <UserCreateForm
        onSubmit={onSubmit}
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

    render(
      <UserCreateForm onSubmit={onSubmit} isPending={false} isError={true} isSuccess={false} />
    )

    expect(screen.getByText(/failed to create user/i)).toBeInTheDocument()
  })

  it('should show success message', () => {
    const onSubmit = vi.fn()

    render(
      <UserCreateForm onSubmit={onSubmit} isPending={false} isError={false} isSuccess={true} />
    )

    expect(screen.getByText(/user created successfully/i)).toBeInTheDocument()
  })

  it('should update form fields on input', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()

    render(
      <UserCreateForm onSubmit={onSubmit} isPending={false} isError={false} isSuccess={false} />
    )

    const nameInput = screen.getByLabelText(/name/i)
    const emailInput = screen.getByLabelText(/email/i)

    await user.type(nameInput, 'Test User')
    await user.type(emailInput, 'test@example.com')

    expect(nameInput).toHaveValue('Test User')
    expect(emailInput).toHaveValue('test@example.com')
  })

  it('should have required validation on fields', () => {
    const onSubmit = vi.fn()

    render(
      <UserCreateForm onSubmit={onSubmit} isPending={false} isError={false} isSuccess={false} />
    )

    expect(screen.getByLabelText(/name/i)).toBeRequired()
    expect(screen.getByLabelText(/email/i)).toBeRequired()
  })
})
