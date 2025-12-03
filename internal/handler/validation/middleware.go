package validation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

// emailPattern は基本的なメールアドレスの正規表現パターン
const emailPattern = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

func init() {
	// email format のカスタムバリデーションを登録
	openapi3.DefineStringFormat("email", emailPattern)
}

// ValidationError はバリデーションエラーの詳細を保持する
type ValidationError struct {
	Message string            `json:"message"`
	Code    string            `json:"code,omitempty"`
	Details []ValidationField `json:"details,omitempty"`
}

// ValidationField は個々のフィールドのバリデーションエラー
type ValidationField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Middleware はOpenAPI定義に基づいてリクエストをバリデートするミドルウェアを作成する
type Middleware struct {
	router routers.Router
}

// NewMiddleware は新しいバリデーションミドルウェアを作成する
func NewMiddleware(openapiSpec []byte) (*Middleware, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(openapiSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}

	// サーバーURLをクリアしてルートパスでマッチングするようにする
	// これにより /api/v1 プレフィックスなしでパスマッチングが可能になる
	doc.Servers = nil

	// スキーマのバリデーション
	if err := doc.Validate(context.Background()); err != nil {
		return nil, fmt.Errorf("invalid OpenAPI spec: %w", err)
	}

	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to create router: %w", err)
	}

	return &Middleware{router: router}, nil
}

// Handler はHTTPミドルウェアとして機能する
func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ルートを検索
		route, pathParams, err := m.router.FindRoute(r)
		if err != nil {
			// ルートが見つからない場合は次のハンドラーに委譲
			// (404はoapi-codegenのハンドラーで処理される)
			next.ServeHTTP(w, r)
			return
		}

		// リクエストボディを読み込んで再利用可能にする
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// リクエストをバリデート
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: pathParams,
			Route:      route,
			Options: &openapi3filter.Options{
				MultiError: true,
			},
		}

		if err := openapi3filter.ValidateRequest(r.Context(), requestValidationInput); err != nil {
			handleValidationError(w, err)
			return
		}

		// リクエストボディを復元して次のハンドラーで使えるようにする
		if bodyBytes != nil {
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		next.ServeHTTP(w, r)
	})
}

// handleValidationError はバリデーションエラーをHTTPレスポンスに変換する
func handleValidationError(w http.ResponseWriter, err error) {
	validationErr := ValidationError{
		Code:    "VALIDATION_ERROR",
		Details: []ValidationField{},
	}

	// MultiErrorの場合は詳細を抽出
	if multiErr, ok := err.(openapi3.MultiError); ok {
		messages := make([]string, 0)
		for _, e := range multiErr {
			detail := extractValidationDetail(e)
			if detail != nil {
				validationErr.Details = append(validationErr.Details, *detail)
				messages = append(messages, detail.Message)
			}
		}
		if len(messages) > 0 {
			validationErr.Message = strings.Join(messages, "; ")
		} else {
			validationErr.Message = "リクエストのバリデーションに失敗しました"
		}
	} else {
		// 単一エラーの場合
		detail := extractValidationDetail(err)
		if detail != nil {
			validationErr.Details = append(validationErr.Details, *detail)
			validationErr.Message = detail.Message
		} else {
			validationErr.Message = err.Error()
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(validationErr)
}

// extractValidationDetail はエラーからフィールド情報を抽出する
func extractValidationDetail(err error) *ValidationField {
	switch e := err.(type) {
	case *openapi3filter.RequestError:
		return extractFromRequestError(e)
	case *openapi3.SchemaError:
		return &ValidationField{
			Field:   getSchemaErrorField(e),
			Message: formatSchemaError(e),
		}
	default:
		return &ValidationField{
			Field:   "",
			Message: err.Error(),
		}
	}
}

// extractFromRequestError はRequestErrorから詳細を抽出する
func extractFromRequestError(err *openapi3filter.RequestError) *ValidationField {
	field := ""
	message := err.Error()

	// パラメータエラーの場合
	if err.Parameter != nil {
		field = err.Parameter.Name
		if schemaErr, ok := err.Err.(*openapi3.SchemaError); ok {
			message = formatSchemaError(schemaErr)
		}
	}

	// ボディエラーの場合
	if err.RequestBody != nil {
		// ネストしたエラーを探索
		if multiErr, ok := err.Err.(openapi3.MultiError); ok {
			for _, innerErr := range multiErr {
				if schemaErr, ok := innerErr.(*openapi3.SchemaError); ok {
					field = getSchemaErrorField(schemaErr)
					message = formatSchemaError(schemaErr)
					break
				}
			}
		} else if schemaErr, ok := err.Err.(*openapi3.SchemaError); ok {
			field = getSchemaErrorField(schemaErr)
			message = formatSchemaError(schemaErr)
		}
	}

	return &ValidationField{
		Field:   field,
		Message: message,
	}
}

// getSchemaErrorField はSchemaErrorからフィールド名を取得する
func getSchemaErrorField(err *openapi3.SchemaError) string {
	if len(err.JSONPointer()) > 0 {
		// JSONポインタから最後の要素を取得
		pointer := err.JSONPointer()
		return pointer[len(pointer)-1]
	}
	return ""
}

// formatSchemaError はSchemaErrorを読みやすいメッセージに変換する
func formatSchemaError(err *openapi3.SchemaError) string {
	reason := err.Reason
	schema := err.Schema

	// よくあるエラーパターンをユーザーフレンドリーなメッセージに変換
	switch {
	case strings.Contains(reason, "minimum string length is"):
		if schema != nil && schema.MinLength > 0 {
			return fmt.Sprintf("%d文字以上で入力してください", schema.MinLength)
		}
	case strings.Contains(reason, "maximum string length is"):
		if schema != nil && schema.MaxLength != nil {
			return fmt.Sprintf("%d文字以下で入力してください", *schema.MaxLength)
		}
	case strings.Contains(reason, "minimum"):
		if schema != nil && schema.Min != nil {
			return fmt.Sprintf("%v以上の値を入力してください", *schema.Min)
		}
	case strings.Contains(reason, "maximum"):
		if schema != nil && schema.Max != nil {
			return fmt.Sprintf("%v以下の値を入力してください", *schema.Max)
		}
	case strings.Contains(reason, "Does not match pattern"):
		return "形式が正しくありません"
	case strings.Contains(reason, "Does not match format"):
		if schema != nil {
			switch schema.Format {
			case "email":
				return "有効なメールアドレスを入力してください"
			case "date-time":
				return "有効な日時形式で入力してください"
			case "uri":
				return "有効なURLを入力してください"
			}
		}
		return "形式が正しくありません"
	case strings.Contains(reason, "property") && strings.Contains(reason, "is missing"):
		return "この項目は必須です"
	}

	return reason
}
