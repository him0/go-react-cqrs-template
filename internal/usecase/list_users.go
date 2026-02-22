package usecase

import (
	"context"
	"log/slog"

	"github.com/example/go-react-cqrs-template/internal/domain"
	"github.com/example/go-react-cqrs-template/internal/pkg/logger"
)

// ListUsersUsecase ユーザー一覧取得ユースケース
type ListUsersUsecase struct {
	userQuery UserQueryRepository
}

// NewListUsersUsecase ListUsersUsecaseのコンストラクタ
func NewListUsersUsecase(userQuery UserQueryRepository) *ListUsersUsecase {
	return &ListUsersUsecase{
		userQuery: userQuery,
	}
}

// Execute ユーザー一覧を取得
func (u *ListUsersUsecase) Execute(ctx context.Context, limit, offset int) ([]*domain.User, int, error) {
	log := logger.FromContext(ctx)
	log.Info("listing users", slog.Int("limit", limit), slog.Int("offset", offset))

	users, err := u.userQuery.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.userQuery.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
