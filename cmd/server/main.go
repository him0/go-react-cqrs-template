package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/example/go-react-cqrs-template/internal/config"
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

	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load config",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// データベース接続設定
	dbConfig := infrastructure.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
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
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Error("failed to close database", slog.String("error", closeErr.Error()))
		}
	}()

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
		log,
	)

	// CORSオリジンの解析（カンマ区切りで複数指定可能）
	corsOrigins := strings.Split(cfg.Server.CORSOrigins, ",")

	// ルーターの設定
	r := chi.NewRouter()

	// ミドルウェア
	r.Use(logger.Middleware) // 構造化ログミドルウェア（リクエストID付与）
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	log.Info("middleware configured",
		slog.String("cors_origins", cfg.Server.CORSOrigins),
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
	port := cfg.Server.Port
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

	// Graceful shutdown（設定のタイムアウト値を使用）
	shutdownTimeout := time.Duration(cfg.Server.ShutdownTimeout) * time.Second
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown error", slog.String("error", err.Error()))
	}
	log.Info("server stopped")
}
