package application

import (
	"errors"

	"github.com/example/go-vite-spec-kit-sample/internal/domain"
)

// UserService アプリケーションサービス
type UserService struct {
	userRepo domain.UserRepository
}

// NewUserService UserServiceのコンストラクタ
func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser ユーザーを作成
func (s *UserService) CreateUser(name, email string) (*domain.User, error) {
	user, err := domain.NewUser(name, email)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser ユーザーを取得
func (s *UserService) GetUser(id string) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}

// ListUsers ユーザー一覧を取得
func (s *UserService) ListUsers(limit, offset int) ([]*domain.User, int, error) {
	users, err := s.userRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.userRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser ユーザーを更新
func (s *UserService) UpdateUser(id, name, email string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := user.Update(name, email); err != nil {
		return nil, err
	}

	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser ユーザーを削除
func (s *UserService) DeleteUser(id string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}
