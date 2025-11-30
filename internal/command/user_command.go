package command

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/example/go-react-spec-kit-sample/internal/domain"
	"github.com/example/go-react-spec-kit-sample/internal/infrastructure"
	"github.com/example/go-react-spec-kit-sample/internal/infrastructure/dao"
)

// Save ユーザーを保存（トランザクション内で使用）
func Save(ctx context.Context, tx infrastructure.DBTX, user *domain.User) error {
	queries := dao.New(tx)
	err := queries.UpsertUser(ctx, dao.UpsertUserParams{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

// Delete ユーザーを削除（トランザクション内で使用）
func Delete(ctx context.Context, tx infrastructure.DBTX, id string) error {
	queries := dao.New(tx)

	// ユーザーの存在確認（FOR UPDATEでロック取得）
	_, err := queries.GetUserByIDForUpdate(ctx, id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user not found: %s", id)
	}
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// 削除実行
	if err := queries.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// FindByIDForUpdate IDでユーザーを検索しロックを取得（トランザクション内で使用）
func FindByIDForUpdate(ctx context.Context, tx infrastructure.DBTX, id string) (*domain.User, error) {
	queries := dao.New(tx)
	user, err := queries.GetUserByIDForUpdate(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user for update: %w", err)
	}
	return toDomainUser(user), nil
}

// FindByEmailForUpdate メールアドレスでユーザーを検索しロックを取得（トランザクション内で使用）
func FindByEmailForUpdate(ctx context.Context, tx infrastructure.DBTX, email string) (*domain.User, error) {
	queries := dao.New(tx)
	user, err := queries.GetUserByEmailForUpdate(ctx, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email for update: %w", err)
	}
	return toDomainUser(user), nil
}

// toDomainUser dao.Userをdomain.Userに変換
func toDomainUser(u dao.User) *domain.User {
	return &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
