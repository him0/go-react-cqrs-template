package domain

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// UserLogAction ユーザーログのアクション種別
type UserLogAction string

const (
	// UserLogActionCreated ユーザー作成
	UserLogActionCreated UserLogAction = "created"
	// UserLogActionDeleted ユーザー削除
	UserLogActionDeleted UserLogAction = "deleted"
)

// UserLog ユーザーログのドメインモデル
type UserLog struct {
	ID        string
	UserID    string
	Action    UserLogAction
	CreatedAt time.Time
}

// NewUserLog ユーザーログを作成
func NewUserLog(userID string, action UserLogAction) *UserLog {
	now := time.Now()
	return &UserLog{
		ID:        ulid.MustNew(ulid.Timestamp(now), rand.Reader).String(),
		UserID:    userID,
		Action:    action,
		CreatedAt: now,
	}
}
