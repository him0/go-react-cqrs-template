package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// pinger DB接続確認用インターフェース（テスト容易性のため）
type pinger interface {
	PingContext(ctx context.Context) error
}

// HealthHandler ヘルスチェック用HTTPハンドラー
type HealthHandler struct {
	db pinger
}

// NewHealthHandler HealthHandlerのコンストラクタ
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// healthResponse ヘルスチェックレスポンス
type healthResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// Liveness サーバーが生きているかを返す (liveness probe)
// 常に 200 OK + {"status": "ok"} を返す
func (h *HealthHandler) Liveness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
}

// Readiness リクエスト受付可能かを返す (readiness probe)
// DB接続を確認し、成功なら 200 + {"status": "ready"}、失敗なら 503 + {"status": "not_ready", "reason": "..."}
func (h *HealthHandler) Readiness(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	if err := h.db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(healthResponse{
			Status: "not_ready",
			Reason: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(healthResponse{Status: "ready"})
}
