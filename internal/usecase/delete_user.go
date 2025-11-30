package usecase

import (
	"context"

	"github.com/example/go-react-cqrs-template/internal/command"
	"github.com/example/go-react-cqrs-template/internal/domain"
	"github.com/example/go-react-cqrs-template/internal/infrastructure"
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
			return domain.ErrUserNotFound(id)
		}

		// ユーザー削除ログを保存
		userLog := domain.NewUserLog(id, domain.UserLogActionDeleted)
		if err := command.SaveUserLog(ctx, tx, userLog); err != nil {
			return err
		}

		// 削除
		return command.Delete(ctx, tx, id)
	})
}
