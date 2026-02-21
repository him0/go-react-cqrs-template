import { describe, it, expect, vi } from 'vitest'
import { render, screen, waitFor } from '@/test/utils'
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

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        name: 'John Doe',
        email: 'john@example.com',
      })
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

  it('should show validation errors for empty fields', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()

    render(
      <UserCreateForm onSubmit={onSubmit} isPending={false} isError={false} isSuccess={false} />
    )

    await user.click(screen.getByRole('button', { name: /create$/i }))

    await waitFor(() => {
      expect(screen.getByText('名前は必須です')).toBeInTheDocument()
    })
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it('should show validation error for invalid email', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()

    render(
      <UserCreateForm onSubmit={onSubmit} isPending={false} isError={false} isSuccess={false} />
    )

    await user.type(screen.getByLabelText(/name/i), 'John Doe')
    await user.type(screen.getByLabelText(/email/i), 'invalid-email')
    await user.click(screen.getByRole('button', { name: /create$/i }))

    await waitFor(() => {
      expect(screen.getByText('有効なメールアドレスを入力してください')).toBeInTheDocument()
    })
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it('should have aria-required attribute on fields', () => {
    const onSubmit = vi.fn()

    render(
      <UserCreateForm onSubmit={onSubmit} isPending={false} isError={false} isSuccess={false} />
    )

    expect(screen.getByLabelText(/name/i)).toHaveAttribute('aria-required', 'true')
    expect(screen.getByLabelText(/email/i)).toHaveAttribute('aria-required', 'true')
  })
})
