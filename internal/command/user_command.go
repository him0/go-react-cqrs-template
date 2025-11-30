package command

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/example/go-react-spec-kit-sample/internal/domain"
	"github.com/example/go-react-spec-kit-sample/internal/infrastructure"
)

// Save ユーザーを保存（トランザクション内で使用）
func Save(ctx context.Context, tx infrastructure.DBTX, user *domain.User) error {
	query := `
		INSERT INTO users (id, name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			updated_at = EXCLUDED.updated_at
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// Delete ユーザーを削除（トランザクション内で使用）
func Delete(ctx context.Context, tx infrastructure.DBTX, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := tx.ExecContext(ctx, query, id)
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

// FindByIDForUpdate IDでユーザーを検索しロックを取得（トランザクション内で使用）
func FindByIDForUpdate(ctx context.Context, tx infrastructure.DBTX, id string) (*domain.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
		FOR UPDATE
	`

	var user domain.User
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user for update: %w", err)
	}

	return &user, nil
}

// FindByEmailForUpdate メールアドレスでユーザーを検索しロックを取得（トランザクション内で使用）
func FindByEmailForUpdate(ctx context.Context, tx infrastructure.DBTX, email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE email = $1
		FOR UPDATE
	`

	var user domain.User
	err := tx.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email for update: %w", err)
	}

	return &user, nil
}
