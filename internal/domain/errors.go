package domain

import "fmt"

// ErrorCode はドメインエラーの種類を識別するコード
type ErrorCode string

const (
	// ErrCodeValidation はバリデーションエラー
	ErrCodeValidation ErrorCode = "VALIDATION_ERROR"
	// ErrCodeNotFound はリソースが見つからないエラー
	ErrCodeNotFound ErrorCode = "NOT_FOUND"
	// ErrCodeConflict はリソースの競合エラー
	ErrCodeConflict ErrorCode = "CONFLICT"
)

// DomainError はドメイン層のエラーを表す基本構造体
type DomainError struct {
	// Code はエラーの種類を識別するコード
	Code ErrorCode
	// Message は内部用のエラーメッセージ（ログ用）
	Message string
	// UserMessage はユーザー向けのメッセージ（API レスポンス用）
	UserMessage string
	// Field はエラーに関連するフィールド名（バリデーションエラーなどで使用）
	Field string
}

func (e *DomainError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s (field: %s)", e.Code, e.Message, e.Field)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// --- バリデーションエラー ---

// ValidationError はバリデーションエラーを表す
type ValidationError struct {
	DomainError
}

// NewValidationError はバリデーションエラーを作成
func NewValidationError(field, message, userMessage string) *ValidationError {
	return &ValidationError{
		DomainError: DomainError{
			Code:        ErrCodeValidation,
			Message:     message,
			UserMessage: userMessage,
			Field:       field,
		},
	}
}

// --- NotFound エラー ---

// NotFoundError はリソースが見つからないエラーを表す
type NotFoundError struct {
	DomainError
	// Resource はリソースの種類（例: "user", "order"）
	Resource string
}

// NewNotFoundError はリソースが見つからないエラーを作成
func NewNotFoundError(resource, message, userMessage string) *NotFoundError {
	return &NotFoundError{
		DomainError: DomainError{
			Code:        ErrCodeNotFound,
			Message:     message,
			UserMessage: userMessage,
		},
		Resource: resource,
	}
}

// --- 競合エラー ---

// ConflictError はリソースの競合エラーを表す
type ConflictError struct {
	DomainError
	// Resource はリソースの種類
	Resource string
}

// NewConflictError は競合エラーを作成
func NewConflictError(resource, message, userMessage string) *ConflictError {
	return &ConflictError{
		DomainError: DomainError{
			Code:        ErrCodeConflict,
			Message:     message,
			UserMessage: userMessage,
		},
		Resource: resource,
	}
}

// --- User 関連のエラー（よく使うものを定義） ---

// ErrUserNotFound はユーザーが見つからないエラー
func ErrUserNotFound(userID string) *NotFoundError {
	return NewNotFoundError(
		"user",
		fmt.Sprintf("user not found: %s", userID),
		"指定されたユーザーが見つかりません",
	)
}

// ErrEmailAlreadyExists はメールアドレスが既に存在するエラー
func ErrEmailAlreadyExists(email string) *ConflictError {
	return NewConflictError(
		"user",
		fmt.Sprintf("email already exists: %s", email),
		"このメールアドレスは既に使用されています",
	)
}

// ErrNameRequired は名前が必須エラー
func ErrNameRequired() *ValidationError {
	return NewValidationError(
		"name",
		"name is required",
		"名前は必須です",
	)
}

// ErrEmailRequired はメールアドレスが必須エラー
func ErrEmailRequired() *ValidationError {
	return NewValidationError(
		"email",
		"email is required",
		"メールアドレスは必須です",
	)
}
