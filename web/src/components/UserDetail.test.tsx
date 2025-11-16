import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@/test/utils'
import userEvent from '@testing-library/user-event'
import { UserDetail } from './UserDetail'
import type { User } from '../api/generated/models'

const mockUser: User = {
  id: '123',
  name: 'John Doe',
  email: 'john@example.com',
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-02T00:00:00Z',
}

describe('UserDetail', () => {
  it('should render user information', () => {
    const onEdit = vi.fn()
    const onDelete = vi.fn()

    render(<UserDetail user={mockUser} onEdit={onEdit} onDelete={onDelete} isDeleting={false} />)

    expect(screen.getByText('123')).toBeInTheDocument()
    expect(screen.getByText('John Doe')).toBeInTheDocument()
    expect(screen.getByText('john@example.com')).toBeInTheDocument()
  })

  it('should format and display dates', () => {
    const onEdit = vi.fn()
    const onDelete = vi.fn()

    render(<UserDetail user={mockUser} onEdit={onEdit} onDelete={onDelete} isDeleting={false} />)

    const createdAt = new Date(mockUser.createdAt).toLocaleString()
    const updatedAt = new Date(mockUser.updatedAt).toLocaleString()

    expect(screen.getByText(createdAt)).toBeInTheDocument()
    expect(screen.getByText(updatedAt)).toBeInTheDocument()
  })

  it('should call onEdit when edit button is clicked', async () => {
    const user = userEvent.setup()
    const onEdit = vi.fn()
    const onDelete = vi.fn()

    render(<UserDetail user={mockUser} onEdit={onEdit} onDelete={onDelete} isDeleting={false} />)

    await user.click(screen.getByRole('button', { name: /edit/i }))

    expect(onEdit).toHaveBeenCalledTimes(1)
  })

  it('should call onDelete when delete button is clicked', async () => {
    const user = userEvent.setup()
    const onEdit = vi.fn()
    const onDelete = vi.fn()

    render(<UserDetail user={mockUser} onEdit={onEdit} onDelete={onDelete} isDeleting={false} />)

    await user.click(screen.getByRole('button', { name: /delete user/i }))

    expect(onDelete).toHaveBeenCalledWith('123')
  })

  it('should show deleting state', () => {
    const onEdit = vi.fn()
    const onDelete = vi.fn()

    render(<UserDetail user={mockUser} onEdit={onEdit} onDelete={onDelete} isDeleting={true} />)

    const deleteButton = screen.getByRole('button', { name: /deleting/i })
    expect(deleteButton).toBeInTheDocument()
    expect(deleteButton).toBeDisabled()
  })

  it('should display all user fields', () => {
    const onEdit = vi.fn()
    const onDelete = vi.fn()

    render(<UserDetail user={mockUser} onEdit={onEdit} onDelete={onDelete} isDeleting={false} />)

    expect(screen.getByText('ID')).toBeInTheDocument()
    expect(screen.getByText('Name')).toBeInTheDocument()
    expect(screen.getByText('Email')).toBeInTheDocument()
    expect(screen.getByText('Created At')).toBeInTheDocument()
    expect(screen.getByText('Updated At')).toBeInTheDocument()
  })

  it('should render action buttons', () => {
    const onEdit = vi.fn()
    const onDelete = vi.fn()

    render(<UserDetail user={mockUser} onEdit={onEdit} onDelete={onDelete} isDeleting={false} />)

    expect(screen.getByRole('button', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /delete user/i })).toBeInTheDocument()
  })
})
