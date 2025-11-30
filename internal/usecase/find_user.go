package usecase

import (
	"context"

	"github.com/example/go-react-cqrs-template/internal/domain"
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
	user, err := u.userQuery.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound(id)
	}
	return user, nil
}
