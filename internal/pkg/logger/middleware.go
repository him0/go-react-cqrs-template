package logger

import (
	"log/slog"
	"net/http"
	"time"
)

// Middleware はHTTPリクエストにリクエストIDを付与し、ログを記録するミドルウェアです
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// リクエストIDを生成
		requestID := GenerateRequestID()

		// コンテキストにリクエストIDを追加
		ctx := WithRequestID(r.Context(), requestID)
		r = r.WithContext(ctx)

		// レスポンスヘッダーにリクエストIDを追加
		w.Header().Set("X-Request-ID", requestID)

		// レスポンスライターをラップして、ステータスコードを記録できるようにする
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// リクエスト開始ログ
		logger := FromContext(ctx)
		logger.Info("request started",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)

		// 次のハンドラを実行
		next.ServeHTTP(wrapped, r)

		// リクエスト完了ログ
		duration := time.Since(start)
		logger.Info("request completed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrapped.statusCode),
			slog.Duration("duration", duration),
			slog.Int64("duration_ms", duration.Milliseconds()),
		)
	})
}

// responseWriter はステータスコードを記録するためのラッパー
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
