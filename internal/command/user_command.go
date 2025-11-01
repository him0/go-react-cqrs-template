package command

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/example/go-vite-spec-kit-sample/internal/domain"
)

// UserCommand ユーザー書き込み操作を担当
type UserCommand struct {
	db *sql.DB
}

// NewUserCommand UserCommandのコンストラクタ
func NewUserCommand(db *sql.DB) *UserCommand {
	return &UserCommand{db: db}
}

// Create ユーザーを作成
func (c *UserCommand) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := c.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update ユーザーを更新
func (c *UserCommand) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := c.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %s", user.ID)
	}

	return nil
}

// Delete ユーザーを削除
func (c *UserCommand) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %s", id)
	}

	return nil
}
