package usecase

import (
	"context"

	"github.com/example/go-react-spec-kit-sample/internal/command"
	"github.com/example/go-react-spec-kit-sample/internal/domain"
	"github.com/example/go-react-spec-kit-sample/internal/infrastructure"
)

// CreateUserUsecase ユーザー作成ユースケース
type CreateUserUsecase struct {
	userQuery UserQueryRepository
	txManager TransactionManager
}

// NewCreateUserUsecase CreateUserUsecaseのコンストラクタ
func NewCreateUserUsecase(
	userQuery UserQueryRepository,
	txManager TransactionManager,
) *CreateUserUsecase {
	return &CreateUserUsecase{
		userQuery: userQuery,
		txManager: txManager,
	}
}

// Execute ユーザーを作成
func (u *CreateUserUsecase) Execute(ctx context.Context, name, email string) (*domain.User, error) {
	var createdUser *domain.User

	err := u.txManager.RunInTransaction(ctx, func(ctx context.Context, tx infrastructure.DBTX) error {
		// メールアドレスの重複チェック（ロック付き）
		existingUser, err := command.FindByEmailForUpdate(ctx, tx, email)
		if err != nil {
			return err
		}
		if existingUser != nil {
			return domain.ErrEmailAlreadyExists(email)
		}

		// ドメインモデルの作成
		user, err := domain.NewUser(name, email)
		if err != nil {
			return err
		}

		// 永続化
		if err := command.Save(ctx, tx, user); err != nil {
			return err
		}

		createdUser = user
		return nil
	})
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
