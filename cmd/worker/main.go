package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/example/go-react-cqrs-template/internal/infrastructure"
	"github.com/example/go-react-cqrs-template/internal/pkg/logger"
	"github.com/example/go-react-cqrs-template/internal/worker"
)

func main() {
	log := logger.Setup()

	// データベース接続設定
	dbConfig := infrastructure.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "app_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	log.Info("connecting to database",
		slog.String("host", dbConfig.Host),
		slog.Int("port", dbConfig.Port),
		slog.String("database", dbConfig.DBName),
	)

	db, err := infrastructure.NewDB(dbConfig)
	if err != nil {
		log.Error("failed to connect to database",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer db.Close()

	log.Info("successfully connected to database")

	txManager := infrastructure.NewTransactionManager(db)

	// ワーカー設定
	workerConfig := worker.Config{
		PollInterval:    getDurationEnv("WORKER_POLL_INTERVAL", 5*time.Second),
		BatchSize:       getEnvInt("WORKER_BATCH_SIZE", 10),
		MaxConcurrency:  getEnvInt("WORKER_MAX_CONCURRENCY", 5),
		ShutdownTimeout: getDurationEnv("WORKER_SHUTDOWN_TIMEOUT", 30*time.Second),
	}

	// ジョブハンドラーの登録
	registry := worker.NewRegistry()
	registerHandlers(registry, log)

	// ワーカーの作成と起動
	w := worker.NewWorker(txManager, registry, workerConfig, log)

	// グレースフルシャットダウン
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Info("received signal, initiating shutdown",
			slog.String("signal", sig.String()),
		)
		cancel()
	}()

	if err := w.Run(ctx); err != nil {
		log.Error("worker exited with error",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	log.Info("worker exited cleanly")
}

// registerHandlers ジョブハンドラーを登録
func registerHandlers(registry *worker.Registry, log *slog.Logger) {
	// サンプル: ウェルカムメール送信ハンドラー
	registry.RegisterFunc("send_welcome_email", func(ctx context.Context, payload json.RawMessage) error {
		var data struct {
			UserID string `json:"user_id"`
			Email  string `json:"email"`
			Name   string `json:"name"`
		}
		if err := json.Unmarshal(payload, &data); err != nil {
			return err
		}
		log.Info("sending welcome email (stub)",
			slog.String("user_id", data.UserID),
			slog.String("email", data.Email),
			slog.String("name", data.Name),
		)
		// TODO: 実際のメール送信処理を実装
		return nil
	})
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return v
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return d
}
