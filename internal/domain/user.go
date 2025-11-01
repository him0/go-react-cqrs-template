package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// User ドメインモデル
type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser ユーザーを作成
func NewUser(name, email string) (*User, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}

	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update ユーザー情報を更新
func (u *User) Update(name, email string) error {
	if name != "" {
		u.Name = name
	}
	if email != "" {
		u.Email = email
	}
	u.UpdatedAt = time.Now()
	return nil
}

// UserRepository ユーザーリポジトリのインターフェース
type UserRepository interface {
	FindByID(id string) (*User, error)
	FindAll(limit, offset int) ([]*User, error)
	Count() (int, error)
	Save(user *User) error
	Delete(id string) error
}
