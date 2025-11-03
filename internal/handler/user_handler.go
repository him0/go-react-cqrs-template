package handler

import (
	"encoding/json"
	"net/http"

	"github.com/example/go-react-spec-kit-sample/internal/usecase"
	"github.com/example/go-react-spec-kit-sample/pkg/generated/openapi"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// UserHandler HTTPハンドラー（OpenAPI生成のServerInterfaceを実装）
type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

// NewUserHandler UserHandlerのコンストラクタ
func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// CreateUser ユーザーを作成（OpenAPI ServerInterface実装）
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req openapi.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userUsecase.CreateUser(r.Context(), req.Name, string(req.Email))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := openapi.User{
		Id:        openapi_types.UUID(user.ID),
		Name:      user.Name,
		Email:     openapi_types.Email(user.Email),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	respondJSON(w, http.StatusCreated, response)
}

// GetUser ユーザーを取得（OpenAPI ServerInterface実装）
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID) {
	user, err := h.userUsecase.GetUser(r.Context(), userId.String())
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	response := openapi.User{
		Id:        openapi_types.UUID(user.ID),
		Name:      user.Name,
		Email:     openapi_types.Email(user.Email),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

// ListUsers ユーザー一覧を取得（OpenAPI ServerInterface実装）
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request, params openapi.ListUsersParams) {
	// デフォルト値の設定
	limit := 10
	offset := 0

	if params.Limit != nil {
		if *params.Limit > 0 && *params.Limit <= 100 {
			limit = *params.Limit
		}
	}

	if params.Offset != nil && *params.Offset >= 0 {
		offset = *params.Offset
	}

	users, total, err := h.userUsecase.ListUsers(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userResponses := make([]openapi.User, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, openapi.User{
			Id:        openapi_types.UUID(user.ID),
			Name:      user.Name,
			Email:     openapi_types.Email(user.Email),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	response := openapi.UserList{
		Users: userResponses,
		Total: total,
	}

	respondJSON(w, http.StatusOK, response)
}

// UpdateUser ユーザーを更新（OpenAPI ServerInterface実装）
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID) {
	var req openapi.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// 更新値の取得（オプショナルなので既存値を保持）
	name := ""
	email := ""

	if req.Name != nil {
		name = *req.Name
	}
	if req.Email != nil {
		email = string(*req.Email)
	}

	user, err := h.userUsecase.UpdateUser(r.Context(), userId.String(), name, email)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := openapi.User{
		Id:        openapi_types.UUID(user.ID),
		Name:      user.Name,
		Email:     openapi_types.Email(user.Email),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

// DeleteUser ユーザーを削除（OpenAPI ServerInterface実装）
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID) {
	if err := h.userUsecase.DeleteUser(r.Context(), userId.String()); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// respondJSON JSONレスポンスを返す
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// respondError エラーレスポンスを返す
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, openapi.Error{Message: message})
}
