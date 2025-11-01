package main

import (
	"log"
	"net/http"

	"github.com/example/go-vite-spec-kit-sample/internal/application"
	"github.com/example/go-vite-spec-kit-sample/internal/infrastructure"
	httphandler "github.com/example/go-vite-spec-kit-sample/internal/interfaces/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// リポジトリの初期化
	userRepo := infrastructure.NewInMemoryUserRepository()

	// サービスの初期化
	userService := application.NewUserService(userRepo)

	// ハンドラーの初期化
	userHandler := httphandler.NewUserHandler(userService)

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

	// APIルート
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.ListUsers)
			r.Post("/", userHandler.CreateUser)
			r.Route("/{userId}", func(r chi.Router) {
				r.Get("/", userHandler.GetUser)
				r.Put("/", userHandler.UpdateUser)
				r.Delete("/", userHandler.DeleteUser)
			})
		})
	})

	// サーバー起動
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
