package main

import (
	"log"
	"net/http"
	"os"

	"github.com/example/go-react-spec-kit-sample/internal/command"
	"github.com/example/go-react-spec-kit-sample/internal/handler"
	"github.com/example/go-react-spec-kit-sample/internal/infrastructure"
	"github.com/example/go-react-spec-kit-sample/internal/queryservice"
	"github.com/example/go-react-spec-kit-sample/internal/usecase"
	"github.com/example/go-react-spec-kit-sample/pkg/generated/openapi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// データベース接続設定
	dbConfig := infrastructure.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "app_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// データベース接続
	db, err := infrastructure.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	// 各層の初期化
	userCommand := command.NewUserCommand(db)
	userQueryService := queryservice.NewUserQueryService(db)
	userUsecase := usecase.NewUserUsecase(userCommand, userQueryService)
	userHandler := handler.NewUserHandler(userUsecase)

	// ルーターの設定
	r := chi.NewRouter()

	// ミドルウェア
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// OpenAPI生成のハンドラーを使用してAPIルートを設定
	r.Route("/api/v1", func(r chi.Router) {
		// OpenAPI仕様に従ったルーティングを自動生成
		openapi.HandlerFromMux(userHandler, r)
	})

	// サーバー起動
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
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
