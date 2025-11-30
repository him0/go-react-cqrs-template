package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/example/go-react-spec-kit-sample/internal/usecase"
	"github.com/example/go-react-spec-kit-sample/pkg/generated/openapi"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// UserHandler HTTPハンドラー（OpenAPI生成のServerInterfaceを実装）
type UserHandler struct {
	createUser *usecase.CreateUserUsecase
	findUser   *usecase.FindUserUsecase
	listUsers  *usecase.ListUsersUsecase
	updateUser *usecase.UpdateUserUsecase
	deleteUser *usecase.DeleteUserUsecase
	logger     *slog.Logger
}

// NewUserHandler UserHandlerのコンストラクタ
func NewUserHandler(
	createUser *usecase.CreateUserUsecase,
	findUser *usecase.FindUserUsecase,
	listUsers *usecase.ListUsersUsecase,
	updateUser *usecase.UpdateUserUsecase,
	deleteUser *usecase.DeleteUserUsecase,
	logger *slog.Logger,
) *UserHandler {
	return &UserHandler{
		createUser: createUser,
		findUser:   findUser,
		listUsers:  listUsers,
		updateUser: updateUser,
		deleteUser: deleteUser,
		logger:     logger,
	}
}

// UsersCreateUser ユーザーを作成（OpenAPI ServerInterface実装）
func (h *UserHandler) UsersCreateUser(w http.ResponseWriter, r *http.Request) {
	var req openapi.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストの形式が不正です")
		return
	}

	user, err := h.createUser.Execute(r.Context(), req.Name, string(req.Email))
	if err != nil {
		HandleError(w, err, h.logger)
		return
	}

	response := openapi.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     openapi_types.Email(user.Email),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	respondJSON(w, http.StatusCreated, response)
}

// UsersGetUser ユーザーを取得（OpenAPI ServerInterface実装）
func (h *UserHandler) UsersGetUser(w http.ResponseWriter, r *http.Request, userId string) {
	user, err := h.findUser.Execute(r.Context(), userId)
	if err != nil {
		HandleError(w, err, h.logger)
		return
	}

	response := openapi.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     openapi_types.Email(user.Email),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

// UsersListUsers ユーザー一覧を取得（OpenAPI ServerInterface実装）
func (h *UserHandler) UsersListUsers(w http.ResponseWriter, r *http.Request, params openapi.UsersListUsersParams) {
	// デフォルト値の設定
	limit := 10
	offset := 0

	if params.Limit != nil {
		if *params.Limit > 0 && *params.Limit <= 100 {
			limit = int(*params.Limit)
		}
	}

	if params.Offset != nil && *params.Offset >= 0 {
		offset = int(*params.Offset)
	}

	users, total, err := h.listUsers.Execute(r.Context(), limit, offset)
	if err != nil {
		HandleError(w, err, h.logger)
		return
	}

	userResponses := make([]openapi.User, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, openapi.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     openapi_types.Email(user.Email),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	response := openapi.UserList{
		Users: userResponses,
		Total: int32(total),
	}

	respondJSON(w, http.StatusOK, response)
}

// UsersUpdateUser ユーザーを更新（OpenAPI ServerInterface実装）
func (h *UserHandler) UsersUpdateUser(w http.ResponseWriter, r *http.Request, userId string) {
	var req openapi.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストの形式が不正です")
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

	user, err := h.updateUser.Execute(r.Context(), userId, name, email)
	if err != nil {
		HandleError(w, err, h.logger)
		return
	}

	response := openapi.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     openapi_types.Email(user.Email),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

// UsersDeleteUser ユーザーを削除（OpenAPI ServerInterface実装）
func (h *UserHandler) UsersDeleteUser(w http.ResponseWriter, r *http.Request, userId string) {
	if err := h.deleteUser.Execute(r.Context(), userId); err != nil {
		HandleError(w, err, h.logger)
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
