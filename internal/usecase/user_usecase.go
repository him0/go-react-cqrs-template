package usecase

import (
	"context"
	"errors"

	"github.com/example/go-vite-spec-kit-sample/internal/command"
	"github.com/example/go-vite-spec-kit-sample/internal/domain"
	"github.com/example/go-vite-spec-kit-sample/internal/queryservice"
)

// UserUsecase ユーザーのユースケースを管理
type UserUsecase struct {
	userCommand      *command.UserCommand
	userQueryService *queryservice.UserQueryService
}

// NewUserUsecase UserUsecaseのコンストラクタ
func NewUserUsecase(
	userCommand *command.UserCommand,
	userQueryService *queryservice.UserQueryService,
) *UserUsecase {
	return &UserUsecase{
		userCommand:      userCommand,
		userQueryService: userQueryService,
	}
}

// CreateUser ユーザーを作成
func (u *UserUsecase) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
	// メールアドレスの重複チェック
	existingUser, err := u.userQueryService.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// ドメインモデルの作成
	user, err := domain.NewUser(name, email)
	if err != nil {
		return nil, err
	}

	// 永続化
	if err := u.userCommand.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser ユーザーを取得
func (u *UserUsecase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.userQueryService.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// ListUsers ユーザー一覧を取得
func (u *UserUsecase) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, int, error) {
	users, err := u.userQueryService.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.userQueryService.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser ユーザーを更新
func (u *UserUsecase) UpdateUser(ctx context.Context, id, name, email string) (*domain.User, error) {
	// 既存ユーザーの取得
	user, err := u.userQueryService.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// メールアドレスが変更される場合、重複チェック
	if email != "" && email != user.Email {
		existingUser, err := u.userQueryService.FindByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		if existingUser != nil {
			return nil, errors.New("email already exists")
		}
	}

	// ドメインモデルの更新
	if err := user.Update(name, email); err != nil {
		return nil, err
	}

	// 永続化
	if err := u.userCommand.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser ユーザーを削除
func (u *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	// 存在確認
	user, err := u.userQueryService.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// 削除
	return u.userCommand.Delete(ctx, id)
}
