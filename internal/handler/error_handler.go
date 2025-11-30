package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/example/go-react-spec-kit-sample/internal/domain"
	apperrors "github.com/example/go-react-spec-kit-sample/internal/pkg/errors"
	"github.com/example/go-react-spec-kit-sample/internal/pkg/logger"
	"github.com/example/go-react-spec-kit-sample/pkg/generated/openapi"
)

// HandleError はエラーをドメインエラーから AppError に変換し、適切な HTTP レスポンスを返す
func HandleError(w http.ResponseWriter, err error, log *slog.Logger) {
	appErr := ToAppError(err)

	// ログ出力
	if log != nil {
		logger.LogError(log, appErr, "request error")
	}

	// HTTP レスポンス
	respondJSON(w, appErr.StatusCode(), openapi.Error{
		Message: appErr.UserMessage(),
	})
}

// ToAppError はドメインエラーを AppError に変換する
func ToAppError(err error) *apperrors.AppError {
	if err == nil {
		return nil
	}

	// 既に AppError の場合はそのまま返す
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	// ValidationError の場合
	var validationErr *domain.ValidationError
	if errors.As(err, &validationErr) {
		return apperrors.BadRequest(
			validationErr.Message,
			validationErr.UserMessage,
		)
	}

	// NotFoundError の場合
	var notFoundErr *domain.NotFoundError
	if errors.As(err, &notFoundErr) {
		return apperrors.NotFound(
			notFoundErr.Resource,
			notFoundErr.UserMessage,
		)
	}

	// ConflictError の場合
	var conflictErr *domain.ConflictError
	if errors.As(err, &conflictErr) {
		return apperrors.Conflict(
			conflictErr.Message,
			conflictErr.UserMessage,
		)
	}

	// DomainError の場合（基底型）
	var domainErr *domain.DomainError
	if errors.As(err, &domainErr) {
		switch domainErr.Code {
		case domain.ErrCodeValidation:
			return apperrors.BadRequest(domainErr.Message, domainErr.UserMessage)
		case domain.ErrCodeNotFound:
			return apperrors.NotFound("resource", domainErr.UserMessage)
		case domain.ErrCodeConflict:
			return apperrors.Conflict(domainErr.Message, domainErr.UserMessage)
		default:
			return apperrors.Internal(err, domainErr.UserMessage)
		}
	}

	// その他の未知のエラーは内部エラーとして処理
	return apperrors.Internal(err, "")
}
