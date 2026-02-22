package usecase

import (
	"context"
	"log/slog"

	"github.com/example/go-react-cqrs-template/internal/command"
	"github.com/example/go-react-cqrs-template/internal/domain"
	"github.com/example/go-react-cqrs-template/internal/infrastructure"
	"github.com/example/go-react-cqrs-template/internal/pkg/logger"
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
	log := logger.FromContext(ctx)
	log.Info("updating user", slog.String("user_id", id))

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
