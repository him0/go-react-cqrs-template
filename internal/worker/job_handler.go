package worker

import (
	"context"
	"encoding/json"
	"fmt"
)

// JobHandler ジョブハンドラーのインターフェース
type JobHandler interface {
	Handle(ctx context.Context, payload json.RawMessage) error
}

// JobHandlerFunc 関数型のJobHandler
type JobHandlerFunc func(ctx context.Context, payload json.RawMessage) error

// Handle JobHandlerインターフェースを実装
func (f JobHandlerFunc) Handle(ctx context.Context, payload json.RawMessage) error {
	return f(ctx, payload)
}

// Registry ジョブハンドラーの登録と取得
type Registry struct {
	handlers map[string]JobHandler
}

// NewRegistry Registryのコンストラクタ
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]JobHandler),
	}
}

// Register ジョブハンドラーを登録
func (r *Registry) Register(jobType string, handler JobHandler) {
	r.handlers[jobType] = handler
}

// RegisterFunc 関数型のジョブハンドラーを登録
func (r *Registry) RegisterFunc(jobType string, fn func(ctx context.Context, payload json.RawMessage) error) {
	r.handlers[jobType] = JobHandlerFunc(fn)
}

// Get ジョブタイプに対応するハンドラーを取得
func (r *Registry) Get(jobType string) (JobHandler, error) {
	handler, ok := r.handlers[jobType]
	if !ok {
		return nil, fmt.Errorf("no handler registered for job type: %s", jobType)
	}
	return handler, nil
}
