package validation

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testOpenAPISpec = []byte(`
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /users:
    post:
      operationId: createUser
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateUserRequest'
      responses:
        '201':
          description: Created
    get:
      operationId: listUsers
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
      responses:
        '200':
          description: OK
  /users/{userId}:
    get:
      operationId: getUser
      parameters:
        - name: userId
          in: path
          required: true
          schema:
            type: string
            pattern: ^[0-9A-HJKMNP-TV-Z]{26}$
      responses:
        '200':
          description: OK
components:
  schemas:
    CreateUserRequest:
      type: object
      required:
        - name
        - email
      properties:
        name:
          type: string
          minLength: 1
          maxLength: 100
        email:
          type: string
          format: email
`)

func TestNewMiddleware(t *testing.T) {
	middleware, err := NewMiddleware(testOpenAPISpec)
	if err != nil {
		t.Fatalf("failed to create middleware: %v", err)
	}
	if middleware == nil {
		t.Fatal("middleware should not be nil")
	}
}

func TestNewMiddleware_InvalidSpec(t *testing.T) {
	_, err := NewMiddleware([]byte("invalid yaml"))
	if err == nil {
		t.Fatal("expected error for invalid spec")
	}
}

func TestMiddleware_ValidRequest(t *testing.T) {
	middleware, err := NewMiddleware(testOpenAPISpec)
	if err != nil {
		t.Fatalf("failed to create middleware: %v", err)
	}

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	body := `{"name": "Test User", "email": "test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestMiddleware_InvalidRequestBody_MissingRequired(t *testing.T) {
	middleware, err := NewMiddleware(testOpenAPISpec)
	if err != nil {
		t.Fatalf("failed to create middleware: %v", err)
	}

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	// email フィールドが欠けている
	body := `{"name": "Test User"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestMiddleware_InvalidRequestBody_InvalidEmail(t *testing.T) {
	middleware, err := NewMiddleware(testOpenAPISpec)
	if err != nil {
		t.Fatalf("failed to create middleware: %v", err)
	}

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	// 無効なメールアドレス
	body := `{"name": "Test User", "email": "invalid-email"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestMiddleware_InvalidRequestBody_EmptyName(t *testing.T) {
	middleware, err := NewMiddleware(testOpenAPISpec)
	if err != nil {
		t.Fatalf("failed to create middleware: %v", err)
	}

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	// 空の名前（minLength: 1 違反）
	body := `{"name": "", "email": "test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestMiddleware_InvalidQueryParameter(t *testing.T) {
	middleware, err := NewMiddleware(testOpenAPISpec)
	if err != nil {
		t.Fatalf("failed to create middleware: %v", err)
	}

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// limit が範囲外（maximum: 100 違反）
	req := httptest.NewRequest(http.MethodGet, "/users?limit=999", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestMiddleware_InvalidPathParameter(t *testing.T) {
	middleware, err := NewMiddleware(testOpenAPISpec)
	if err != nil {
		t.Fatalf("failed to create middleware: %v", err)
	}

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 無効な ULID パターン
	req := httptest.NewRequest(http.MethodGet, "/users/invalid-id", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestMiddleware_ValidPathParameter(t *testing.T) {
	middleware, err := NewMiddleware(testOpenAPISpec)
	if err != nil {
		t.Fatalf("failed to create middleware: %v", err)
	}

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 有効な ULID パターン
	req := httptest.NewRequest(http.MethodGet, "/users/01ARZ3NDEKTSV4RRFFQ69G5FAV", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestMiddleware_UnmatchedRoute(t *testing.T) {
	middleware, err := NewMiddleware(testOpenAPISpec)
	if err != nil {
		t.Fatalf("failed to create middleware: %v", err)
	}

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	// OpenAPI spec に定義されていないルート
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// バリデーションミドルウェアはスキップして次のハンドラーに委譲
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
