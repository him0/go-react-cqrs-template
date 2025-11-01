.PHONY: help install run-backend run-frontend generate-api build clean test test-backend test-frontend test-coverage

help: ## ヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## 依存関係をインストール
	go mod download
	cd web && pnpm install

run-backend: ## バックエンドサーバーを起動
	go run cmd/server/main.go

run-frontend: ## フロントエンド開発サーバーを起動
	cd web && pnpm run dev

generate-api: ## APIコードを生成（フロントエンド）
	cd web && pnpm run generate:api

build-backend: ## バックエンドをビルド
	go build -o bin/server cmd/server/main.go

build-frontend: ## フロントエンドをビルド
	cd web && pnpm run build

build: build-backend build-frontend ## すべてをビルド

clean: ## ビルド成果物を削除
	rm -rf bin/
	rm -rf web/dist/
	rm -rf web/src/api/generated/
	rm -rf web/coverage/
	rm -f coverage.out

dev: ## 開発環境を起動（バックエンド + フロントエンド）
	@echo "バックエンドとフロントエンドを別々のターミナルで起動してください："
	@echo "  ターミナル1: make run-backend"
	@echo "  ターミナル2: make run-frontend"

test: test-backend test-frontend ## すべてのテストを実行

test-backend: ## バックエンドのテストを実行
	go test -v -race ./...

test-frontend: ## フロントエンドのテストを実行
	cd web && pnpm run test -- --run

test-coverage: ## カバレッジ付きでテストを実行
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	cd web && pnpm run test:coverage -- --run

test-watch: ## ウォッチモードでフロントエンドテストを実行
	cd web && pnpm run test
