package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/example/go-vite-spec-kit-sample/internal/application"
	"github.com/go-chi/chi/v5"
)

// UserHandler ユーザーハンドラー
type UserHandler struct {
	userService *application.UserService
}

// NewUserHandler UserHandlerのコンストラクタ
func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUserRequest ユーザー作成リクエスト
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserRequest ユーザー更新リクエスト
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// UserResponse ユーザーレスポンス
type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// UserListResponse ユーザー一覧レスポンス
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
}

// ErrorResponse エラーレスポンス
type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// ListUsers ユーザー一覧を取得
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}

	users, total, err := h.userService.ListUsers(limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	respondJSON(w, http.StatusOK, UserListResponse{
		Users: userResponses,
		Total: total,
	})
}

// CreateUser ユーザーを作成
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userService.CreateUser(req.Name, req.Email)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// GetUser ユーザーを取得
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	user, err := h.userService.GetUser(userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	respondJSON(w, http.StatusOK, UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// UpdateUser ユーザーを更新
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userService.UpdateUser(userID, req.Name, req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// DeleteUser ユーザーを削除
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	if err := h.userService.DeleteUser(userID); err != nil {
		if err.Error() == "user not found" {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{
		Message: message,
	})
}
