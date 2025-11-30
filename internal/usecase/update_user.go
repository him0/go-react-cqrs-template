package usecase

import (
	"context"

	"github.com/example/go-react-spec-kit-sample/internal/command"
	"github.com/example/go-react-spec-kit-sample/internal/domain"
	"github.com/example/go-react-spec-kit-sample/internal/infrastructure"
)

// UpdateUserUsecase ユーザー更新ユースケース
type UpdateUserUsecase struct {
	userQuery UserQueryRepository
	txManager TransactionManager
}

// NewUpdateUserUsecase UpdateUserUsecaseのコンストラクタ
func NewUpdateUserUsecase(
	userQuery UserQueryRepository,
	txManager TransactionManager,
) *UpdateUserUsecase {
	return &UpdateUserUsecase{
		userQuery: userQuery,
		txManager: txManager,
	}
}

// Execute ユーザーを更新
func (u *UpdateUserUsecase) Execute(ctx context.Context, id, name, email string) error {
	return u.txManager.RunInTransaction(ctx, func(ctx context.Context, tx infrastructure.DBTX) error {
		// 行ロック付きでユーザーを取得
		user, err := command.FindByIDForUpdate(ctx, tx, id)
		if err != nil {
			return err
		}
		if user == nil {
			return domain.ErrUserNotFound(id)
		}

		// メールアドレスが変更される場合、重複チェック（ロック付き）
		if email != "" && email != user.Email {
			existingUser, err := command.FindByEmailForUpdate(ctx, tx, email)
			if err != nil {
				return err
			}
			if existingUser != nil {
				return domain.ErrEmailAlreadyExists(email)
			}
		}

		// ドメインモデルの更新
		if err := user.Update(name, email); err != nil {
			return err
		}

		// 永続化
		return command.Save(ctx, tx, user)
	})
}
