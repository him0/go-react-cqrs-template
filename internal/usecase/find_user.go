package usecase

import (
	"context"
	"log/slog"

	"github.com/example/go-react-cqrs-template/internal/domain"
	"github.com/example/go-react-cqrs-template/internal/pkg/logger"
)

// FindUserUsecase ユーザー取得ユースケース
type FindUserUsecase struct {
	userQuery UserQueryRepository
}

// NewFindUserUsecase FindUserUsecaseのコンストラクタ
func NewFindUserUsecase(userQuery UserQueryRepository) *FindUserUsecase {
	return &FindUserUsecase{
		userQuery: userQuery,
	}
}

// Execute ユーザーを取得
func (u *FindUserUsecase) Execute(ctx context.Context, id string) (*domain.User, error) {
	log := logger.FromContext(ctx)
	log.Info("finding user", slog.String("user_id", id))

	user, err := u.userQuery.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound(id)
	}
	return user, nil
}
