package main

import (
	"log/slog"
	"net/http"
	"os"

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
		Port:     5432,
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
		log,
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

	// サーバー起動
	port := getEnv("PORT", "8080")
	log.Info("server starting",
		slog.String("port", port),
		slog.String("address", ":"+port),
	)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Error("server failed to start",
			slog.String("error", err.Error()),
			slog.String("port", port),
		)
		os.Exit(1)
	}
}

// getEnv 環境変数を取得、なければデフォルト値を返す
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
