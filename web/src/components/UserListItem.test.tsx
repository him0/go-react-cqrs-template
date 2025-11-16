import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@/test/utils'
import userEvent from '@testing-library/user-event'
import { UserListItem } from './UserListItem'
import type { User } from '../api/generated/models'

const mockUser: User = {
  id: '1',
  name: 'John Doe',
  email: 'john@example.com',
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
}

describe('UserListItem', () => {
  it('should render user information', () => {
    const onSelect = vi.fn()
    const onDelete = vi.fn()

    render(
      <UserListItem
        user={mockUser}
        isSelected={false}
        onSelect={onSelect}
        onDelete={onDelete}
        isDeleting={false}
      />
    )

    expect(screen.getByText('John Doe')).toBeInTheDocument()
    expect(screen.getByText('john@example.com')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /delete/i })).toBeInTheDocument()
  })

  it('should call onSelect when item is clicked', async () => {
    const user = userEvent.setup()
    const onSelect = vi.fn()
    const onDelete = vi.fn()

    render(
      <UserListItem
        user={mockUser}
        isSelected={false}
        onSelect={onSelect}
        onDelete={onDelete}
        isDeleting={false}
      />
    )

    const item = screen.getByText('John Doe').closest('div')
    await user.click(item!)

    expect(onSelect).toHaveBeenCalledWith('1')
    expect(onSelect).toHaveBeenCalledTimes(1)
  })

  it('should call onDelete when delete button is clicked', async () => {
    const user = userEvent.setup()
    const onSelect = vi.fn()
    const onDelete = vi.fn()

    render(
      <UserListItem
        user={mockUser}
        isSelected={false}
        onSelect={onSelect}
        onDelete={onDelete}
        isDeleting={false}
      />
    )

    await user.click(screen.getByRole('button', { name: /delete/i }))

    expect(onDelete).toHaveBeenCalledWith('1')
    expect(onDelete).toHaveBeenCalledTimes(1)
    expect(onSelect).not.toHaveBeenCalled()
  })

  it('should apply selected styles when isSelected is true', () => {
    const onSelect = vi.fn()
    const onDelete = vi.fn()

    const { container } = render(
      <UserListItem
        user={mockUser}
        isSelected={true}
        onSelect={onSelect}
        onDelete={onDelete}
        isDeleting={false}
      />
    )

    const item = container.querySelector('.border-blue-500')
    expect(item).toBeInTheDocument()
  })

  it('should disable delete button when isDeleting is true', () => {
    const onSelect = vi.fn()
    const onDelete = vi.fn()

    render(
      <UserListItem
        user={mockUser}
        isSelected={false}
        onSelect={onSelect}
        onDelete={onDelete}
        isDeleting={true}
      />
    )

    expect(screen.getByRole('button', { name: /delete/i })).toBeDisabled()
  })

  it('should not call onSelect when delete button is clicked', async () => {
    const user = userEvent.setup()
    const onSelect = vi.fn()
    const onDelete = vi.fn()

    render(
      <UserListItem
        user={mockUser}
        isSelected={false}
        onSelect={onSelect}
        onDelete={onDelete}
        isDeleting={false}
      />
    )

    await user.click(screen.getByRole('button', { name: /delete/i }))

    expect(onDelete).toHaveBeenCalled()
    expect(onSelect).not.toHaveBeenCalled()
  })
})
