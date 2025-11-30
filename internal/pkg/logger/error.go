package logger

import (
	"log/slog"

	"github.com/example/go-react-cqrs-template/internal/pkg/errors"
)

// LogError はAppErrorを構造化ログとして出力します
func LogError(logger *slog.Logger, appErr *errors.AppError, msg string, additionalAttrs ...slog.Attr) {
	if logger == nil {
		logger = Get()
	}

	// ベース属性を構築
	attrs := []any{
		slog.String("error", appErr.Message()),
		slog.String("user_message", appErr.UserMessage()),
		slog.Int("status_code", appErr.StatusCode()),
		slog.String("error_level", appErr.Level().String()),
	}

	// 元のエラーがある場合
	if appErr.Cause() != nil {
		attrs = append(attrs, slog.String("cause", appErr.Cause().Error()))
	}

	// スタックトレースを追加（ERROR以上のレベルの場合のみ）
	if appErr.Level() >= errors.LevelError && len(appErr.Stack()) > 0 {
		// スタックトレースを文字列の配列として追加
		attrs = append(attrs, slog.Any("stack_trace", appErr.Stack()))
	}

	// 追加の属性をマージ
	for _, attr := range additionalAttrs {
		attrs = append(attrs, attr)
	}

	// エラーレベルに応じてログレベルを決定
	switch appErr.Level() {
	case errors.LevelInfo:
		logger.Info(msg, attrs...)
	case errors.LevelWarning:
		logger.Warn(msg, attrs...)
	case errors.LevelError:
		logger.Error(msg, attrs...)
	case errors.LevelCritical:
		logger.Error(msg, attrs...) // slogにはCriticalがないのでErrorを使用
	default:
		logger.Error(msg, attrs...)
	}
}
