package infrastructure

import (
	"errors"
	"sync"

	"github.com/example/go-vite-spec-kit-sample/internal/domain"
)

// InMemoryUserRepository インメモリユーザーリポジトリ
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

// NewInMemoryUserRepository InMemoryUserRepositoryのコンストラクタ
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

// FindByID IDでユーザーを検索
func (r *InMemoryUserRepository) FindByID(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// FindAll 全ユーザーを取得
func (r *InMemoryUserRepository) FindAll(limit, offset int) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*domain.User, 0)
	i := 0
	for _, user := range r.users {
		if i >= offset && i < offset+limit {
			users = append(users, user)
		}
		i++
		if i >= offset+limit {
			break
		}
	}

	return users, nil
}

// Count ユーザー数を取得
func (r *InMemoryUserRepository) Count() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.users), nil
}

// Save ユーザーを保存
func (r *InMemoryUserRepository) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user
	return nil
}

// Delete ユーザーを削除
func (r *InMemoryUserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}
