package command

import (
	"context"
	"fmt"

	"github.com/example/go-react-cqrs-template/internal/domain"
	"github.com/example/go-react-cqrs-template/internal/infrastructure"
	"github.com/example/go-react-cqrs-template/internal/infrastructure/dao"
)

// SaveUserLog ユーザーログを保存（トランザクション内で使用）
func SaveUserLog(ctx context.Context, tx infrastructure.DBTX, log *domain.UserLog) error {
	queries := dao.New(tx)
	err := queries.CreateUserLog(ctx, dao.CreateUserLogParams{
		ID:        log.ID,
		UserID:    log.UserID,
		Action:    string(log.Action),
		CreatedAt: log.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to save user log: %w", err)
	}
	return nil
}
