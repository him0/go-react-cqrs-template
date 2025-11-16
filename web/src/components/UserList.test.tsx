import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@/test/utils'
import userEvent from '@testing-library/user-event'
import { UserList } from './UserList'
import type { User } from '../api/generated/models'

const mockUsers: User[] = [
  {
    id: '1',
    name: 'John Doe',
    email: 'john@example.com',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  },
  {
    id: '2',
    name: 'Jane Smith',
    email: 'jane@example.com',
    createdAt: '2024-01-02T00:00:00Z',
    updatedAt: '2024-01-02T00:00:00Z',
  },
]

describe('UserList', () => {
  it('should render loading state', () => {
    const onSelectUser = vi.fn()
    const onDeleteUser = vi.fn()

    render(
      <UserList
        users={undefined}
        total={undefined}
        isLoading={true}
        error={null}
        selectedUserId={null}
        onSelectUser={onSelectUser}
        onDeleteUser={onDeleteUser}
        isDeleting={false}
      />
    )

    expect(screen.getByText(/loading users/i)).toBeInTheDocument()
  })

  it('should render error state', () => {
    const error = { message: 'Network error' }
    const onSelectUser = vi.fn()
    const onDeleteUser = vi.fn()

    render(
      <UserList
        users={undefined}
        total={undefined}
        isLoading={false}
        error={error}
        selectedUserId={null}
        onSelectUser={onSelectUser}
        onDeleteUser={onDeleteUser}
        isDeleting={false}
      />
    )

    expect(screen.getByText(/error: network error/i)).toBeInTheDocument()
  })

  it('should render empty state when no users', () => {
    const onSelectUser = vi.fn()
    const onDeleteUser = vi.fn()

    render(
      <UserList
        users={[]}
        total={0}
        isLoading={false}
        error={null}
        selectedUserId={null}
        onSelectUser={onSelectUser}
        onDeleteUser={onDeleteUser}
        isDeleting={false}
      />
    )

    expect(screen.getByText(/no users found/i)).toBeInTheDocument()
  })

  it('should render list of users', () => {
    const onSelectUser = vi.fn()
    const onDeleteUser = vi.fn()

    render(
      <UserList
        users={mockUsers}
        total={2}
        isLoading={false}
        error={null}
        selectedUserId={null}
        onSelectUser={onSelectUser}
        onDeleteUser={onDeleteUser}
        isDeleting={false}
      />
    )

    expect(screen.getByText('John Doe')).toBeInTheDocument()
    expect(screen.getByText('john@example.com')).toBeInTheDocument()
    expect(screen.getByText('Jane Smith')).toBeInTheDocument()
    expect(screen.getByText('jane@example.com')).toBeInTheDocument()
    expect(screen.getByText(/total users: 2/i)).toBeInTheDocument()
  })

  it('should call onSelectUser when user item is clicked', async () => {
    const user = userEvent.setup()
    const onSelectUser = vi.fn()
    const onDeleteUser = vi.fn()

    render(
      <UserList
        users={mockUsers}
        total={2}
        isLoading={false}
        error={null}
        selectedUserId={null}
        onSelectUser={onSelectUser}
        onDeleteUser={onDeleteUser}
        isDeleting={false}
      />
    )

    const firstUser = screen.getByText('John Doe').closest('div')
    await user.click(firstUser!)

    expect(onSelectUser).toHaveBeenCalledWith('1')
  })

  it('should call onDeleteUser when delete button is clicked', async () => {
    const user = userEvent.setup()
    const onSelectUser = vi.fn()
    const onDeleteUser = vi.fn()

    render(
      <UserList
        users={mockUsers}
        total={2}
        isLoading={false}
        error={null}
        selectedUserId={null}
        onSelectUser={onSelectUser}
        onDeleteUser={onDeleteUser}
        isDeleting={false}
      />
    )

    const deleteButtons = screen.getAllByRole('button', { name: /delete/i })
    await user.click(deleteButtons[0])

    expect(onDeleteUser).toHaveBeenCalledWith('1')
  })

  it('should highlight selected user', () => {
    const onSelectUser = vi.fn()
    const onDeleteUser = vi.fn()

    const { container } = render(
      <UserList
        users={mockUsers}
        total={2}
        isLoading={false}
        error={null}
        selectedUserId="1"
        onSelectUser={onSelectUser}
        onDeleteUser={onDeleteUser}
        isDeleting={false}
      />
    )

    const selectedItem = container.querySelector('.border-blue-500')
    expect(selectedItem).toBeInTheDocument()
  })

  it('should not show total when total is undefined', () => {
    const onSelectUser = vi.fn()
    const onDeleteUser = vi.fn()

    render(
      <UserList
        users={mockUsers}
        total={undefined}
        isLoading={false}
        error={null}
        selectedUserId={null}
        onSelectUser={onSelectUser}
        onDeleteUser={onDeleteUser}
        isDeleting={false}
      />
    )

    expect(screen.queryByText(/total users/i)).not.toBeInTheDocument()
  })
})
