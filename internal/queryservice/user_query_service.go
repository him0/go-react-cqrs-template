package queryservice

import (
	"context"
	"database/sql"

	"github.com/example/go-react-spec-kit-sample/internal/domain"
	"github.com/example/go-react-spec-kit-sample/internal/infrastructure/dao"
)

// UserQueryService ユーザー読み取り操作を担当
type UserQueryService struct {
	queries *dao.Queries
}

// NewUserQueryService UserQueryServiceのコンストラクタ
func NewUserQueryService(db *sql.DB) *UserQueryService {
	return &UserQueryService{queries: dao.New(db)}
}

// FindByID IDでユーザーを検索
func (q *UserQueryService) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := q.queries.GetUserByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomainUser(user), nil
}

// FindAll すべてのユーザーを取得（ページネーション対応）
func (q *UserQueryService) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users, err := q.queries.ListUsers(ctx, dao.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return toDomainUsers(users), nil
}

// Count ユーザーの総数を取得
func (q *UserQueryService) Count(ctx context.Context) (int, error) {
	count, err := q.queries.CountUsers(ctx)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// FindByEmail メールアドレスでユーザーを検索
func (q *UserQueryService) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := q.queries.GetUserByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomainUser(user), nil
}

// toDomainUser dao.Userをdomain.Userに変換
func toDomainUser(u dao.User) *domain.User {
	return &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// toDomainUsers []dao.Userを[]*domain.Userに変換
func toDomainUsers(users []dao.User) []*domain.User {
	result := make([]*domain.User, len(users))
	for i, u := range users {
		result[i] = toDomainUser(u)
	}
	return result
}
