package usecase

import (
	"context"
	"errors"

	"github.com/example/go-react-spec-kit-sample/internal/command"
	"github.com/example/go-react-spec-kit-sample/internal/infrastructure"
)

// DeleteUserUsecase ユーザー削除ユースケース
type DeleteUserUsecase struct {
	userQuery UserQueryRepository
	txManager TransactionManager
}

// NewDeleteUserUsecase DeleteUserUsecaseのコンストラクタ
func NewDeleteUserUsecase(
	userQuery UserQueryRepository,
	txManager TransactionManager,
) *DeleteUserUsecase {
	return &DeleteUserUsecase{
		userQuery: userQuery,
		txManager: txManager,
	}
}

// Execute ユーザーを削除
func (u *DeleteUserUsecase) Execute(ctx context.Context, id string) error {
	return u.txManager.RunInTransaction(ctx, func(ctx context.Context, tx infrastructure.DBTX) error {
		// 行ロック付きで存在確認
		user, err := command.FindByIDForUpdate(ctx, tx, id)
		if err != nil {
			return err
		}
		if user == nil {
			return errors.New("user not found")
		}

		// 削除
		return command.Delete(ctx, tx, id)
	})
}
