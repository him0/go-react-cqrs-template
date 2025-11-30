package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// Level はエラーのレベルを表します
type Level int

const (
	// LevelInfo は情報レベル（通常のビジネスロジックエラー）
	LevelInfo Level = iota
	// LevelWarning は警告レベル（注意が必要だが処理は継続可能）
	LevelWarning
	// LevelError はエラーレベル（予期しないエラー）
	LevelError
	// LevelCritical は重大エラーレベル（システムに深刻な影響）
	LevelCritical
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelWarning:
		return "WARNING"
	case LevelError:
		return "ERROR"
	case LevelCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// AppError はアプリケーション全体で使用するカスタムエラー型です
type AppError struct {
	// 内部エラーメッセージ（ログ用）
	message string
	// ユーザー向けメッセージ
	userMessage string
	// HTTPステータスコード
	statusCode int
	// エラーレベル
	level Level
	// 元のエラー
	cause error
	// スタックトレース
	stack []string
}

// Error は error インターフェースを実装します
func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

// Message は内部エラーメッセージを返します
func (e *AppError) Message() string {
	return e.message
}

// UserMessage はユーザー向けメッセージを返します
func (e *AppError) UserMessage() string {
	if e.userMessage != "" {
		return e.userMessage
	}
	return "予期しないエラーが発生しました"
}

// StatusCode はHTTPステータスコードを返します
func (e *AppError) StatusCode() int {
	return e.statusCode
}

// Level はエラーレベルを返します
func (e *AppError) Level() Level {
	return e.level
}

// Cause は元のエラーを返します
func (e *AppError) Cause() error {
	return e.cause
}

// Stack はスタックトレースを返します
func (e *AppError) Stack() []string {
	return e.stack
}

// Unwrap は元のエラーを返します（Go 1.13+ のエラーチェーン対応）
func (e *AppError) Unwrap() error {
	return e.cause
}

// New は新しいAppErrorを作成します
func New(message string, userMessage string, statusCode int, level Level) *AppError {
	return &AppError{
		message:     message,
		userMessage: userMessage,
		statusCode:  statusCode,
		level:       level,
		stack:       captureStack(2),
	}
}

// Wrap は既存のエラーをラップして新しいAppErrorを作成します
func Wrap(err error, message string, userMessage string, statusCode int, level Level) *AppError {
	if err == nil {
		return nil
	}

	// 既にAppErrorの場合は、スタックトレースを保持
	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			message:     message,
			userMessage: userMessage,
			statusCode:  statusCode,
			level:       level,
			cause:       appErr,
			stack:       appErr.stack, // 元のスタックトレースを保持
		}
	}

	return &AppError{
		message:     message,
		userMessage: userMessage,
		statusCode:  statusCode,
		level:       level,
		cause:       err,
		stack:       captureStack(2),
	}
}

// よく使うエラーのヘルパー関数

// NotFound はリソースが見つからないエラーを作成します
func NotFound(resource string, userMessage string) *AppError {
	if userMessage == "" {
		userMessage = fmt.Sprintf("%sが見つかりませんでした", resource)
	}
	return New(
		fmt.Sprintf("%s not found", resource),
		userMessage,
		http.StatusNotFound,
		LevelInfo,
	)
}

// BadRequest は不正なリクエストエラーを作成します
func BadRequest(message string, userMessage string) *AppError {
	if userMessage == "" {
		userMessage = "リクエストが不正です"
	}
	return New(
		message,
		userMessage,
		http.StatusBadRequest,
		LevelInfo,
	)
}

// Internal は内部サーバーエラーを作成します
func Internal(err error, userMessage string) *AppError {
	if userMessage == "" {
		userMessage = "サーバー内部エラーが発生しました"
	}
	message := "internal server error"
	if err != nil {
		message = err.Error()
	}
	return Wrap(
		err,
		message,
		userMessage,
		http.StatusInternalServerError,
		LevelError,
	)
}

// Unauthorized は認証エラーを作成します
func Unauthorized(message string, userMessage string) *AppError {
	if userMessage == "" {
		userMessage = "認証が必要です"
	}
	return New(
		message,
		userMessage,
		http.StatusUnauthorized,
		LevelInfo,
	)
}

// Forbidden は権限エラーを作成します
func Forbidden(message string, userMessage string) *AppError {
	if userMessage == "" {
		userMessage = "この操作を実行する権限がありません"
	}
	return New(
		message,
		userMessage,
		http.StatusForbidden,
		LevelInfo,
	)
}

// Conflict はリソース競合エラーを作成します
func Conflict(message string, userMessage string) *AppError {
	if userMessage == "" {
		userMessage = "データが競合しています"
	}
	return New(
		message,
		userMessage,
		http.StatusConflict,
		LevelWarning,
	)
}

// captureStack はスタックトレースをキャプチャします
func captureStack(skip int) []string {
	const maxDepth = 32
	var stack []string

	for i := skip; i < maxDepth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// 関数名を取得
		fn := runtime.FuncForPC(pc)
		var funcName string
		if fn != nil {
			funcName = fn.Name()
			// パッケージパスを短縮
			if idx := strings.LastIndex(funcName, "/"); idx >= 0 {
				funcName = funcName[idx+1:]
			}
		}

		// ファイルパスを短縮（プロジェクトルート以降のみ）
		if idx := strings.Index(file, "go-react-cqrs-template/"); idx >= 0 {
			file = file[idx:]
		}

		stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, funcName))
	}

	return stack
}
