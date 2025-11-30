package domain

import (
	"time"

	"github.com/google/uuid"
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
	return &UserLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    action,
		CreatedAt: time.Now(),
	}
}
