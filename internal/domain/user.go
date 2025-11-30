package domain

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
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
		return nil, ErrNameRequired()
	}
	if email == "" {
		return nil, ErrEmailRequired()
	}

	now := time.Now()
	return &User{
		ID:        ulid.MustNew(ulid.Timestamp(now), rand.Reader).String(),
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
