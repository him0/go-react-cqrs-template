package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/example/go-react-cqrs-template/internal/handler"
	"github.com/example/go-react-cqrs-template/internal/handler/validation"
	"github.com/example/go-react-cqrs-template/internal/infrastructure"
	"github.com/example/go-react-cqrs-template/internal/pkg/logger"
	"github.com/example/go-react-cqrs-template/internal/queryservice"
	"github.com/example/go-react-cqrs-template/internal/usecase"
	openapispec "github.com/example/go-react-cqrs-template/openapi"
	"github.com/example/go-react-cqrs-template/pkg/generated/openapi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// ロガーのセットアップ
	log := logger.Setup()

	// データベース接続設定
	dbConfig := infrastructure.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvInt("DB_PORT", 55432),
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

	// データベース接続
	db, err := infrastructure.NewDB(dbConfig)
	if err != nil {
		log.Error("failed to connect to database",
			slog.String("error", err.Error()),
			slog.String("host", dbConfig.Host),
			slog.String("database", dbConfig.DBName),
		)
		os.Exit(1)
	}
	defer db.Close()

	log.Info("successfully connected to database")

	// 各層の初期化
	txManager := infrastructure.NewTransactionManager(db)
	userQueryService := queryservice.NewUserQueryService(db)

	// Usecases
	createUserUsecase := usecase.NewCreateUserUsecase(userQueryService, txManager)
	findUserUsecase := usecase.NewFindUserUsecase(userQueryService)
	listUsersUsecase := usecase.NewListUsersUsecase(userQueryService)
	updateUserUsecase := usecase.NewUpdateUserUsecase(userQueryService, txManager)
	deleteUserUsecase := usecase.NewDeleteUserUsecase(userQueryService, txManager)

	userHandler := handler.NewUserHandler(
		createUserUsecase,
		findUserUsecase,
		listUsersUsecase,
		updateUserUsecase,
		deleteUserUsecase,
	)

	// ルーターの設定
	r := chi.NewRouter()

	// ミドルウェア
	r.Use(logger.Middleware) // 構造化ログミドルウェア（リクエストID付与）
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	log.Info("middleware configured",
		slog.String("cors_origin", "http://localhost:3000"),
	)

	// ヘルスチェックエンドポイント（バリデーションミドルウェア不要）
	healthHandler := handler.NewHealthHandler(db)
	r.Get("/healthz", healthHandler.Liveness)
	r.Get("/readyz", healthHandler.Readiness)

	// OpenAPIバリデーションミドルウェアの初期化
	validationMiddleware, err := validation.NewMiddleware(openapispec.Spec)
	if err != nil {
		log.Error("failed to create validation middleware",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	log.Info("OpenAPI validation middleware initialized")

	// OpenAPI生成のハンドラーを使用してAPIルートを設定
	r.Route("/api/v1", func(r chi.Router) {
		// OpenAPI仕様に基づくリクエストバリデーション
		r.Use(validationMiddleware.Handler)
		// OpenAPI仕様に従ったルーティングを自動生成
		openapi.HandlerFromMux(userHandler, r)
	})

	// シグナルハンドリングの設定
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// サーバー起動
	port := getEnv("PORT", "8080")
	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		log.Info("server starting",
			slog.String("port", port),
			slog.String("address", ":"+port),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error",
				slog.String("error", err.Error()),
				slog.String("port", port),
			)
			os.Exit(1)
		}
	}()

	// シグナル受信を待機
	<-ctx.Done()
	log.Info("shutting down server...")

	// Graceful shutdown（タイムアウト: 30秒）
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown error", slog.String("error", err.Error()))
	}
	log.Info("server stopped")
}

// getEnv 環境変数を取得、なければデフォルト値を返す
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt 環境変数をint型で取得、なければデフォルト値を返す
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
