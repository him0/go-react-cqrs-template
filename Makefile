.PHONY: help install run-backend run-frontend generate-api build clean test test-backend test-frontend test-coverage docker-up docker-down docker-logs db-migrate db-dry-run setup lint fmt vet check-fmt check-imports modernize modernize-check ci-test

help: ## ヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## 開発環境のセットアップ
	go mod download
	go install github.com/sqldef/sqldef/cmd/psqldef@latest
	go install golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest
	go install golang.org/x/tools/cmd/goimports@latest
	cd web && pnpm install

install: ## 依存関係をインストール
	go mod download
	cd web && pnpm install

docker-up: ## Dockerコンテナを起動
	docker-compose up -d

docker-down: ## Dockerコンテナを停止
	docker-compose down

docker-logs: ## Dockerログを表示
	docker-compose logs -f

db-migrate: ## psqldefを使用してデータベースマイグレーションを実行
	psqldef -U postgres -p 5432 -h localhost app_db --password=postgres --file=db/schema/schema.sql

db-dry-run: ## データベースマイグレーションのドライラン
	psqldef -U postgres -p 5432 -h localhost app_db --password=postgres --file=db/schema/schema.sql --dry-run

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
	@echo "開発環境の起動手順："
	@echo "  1. make docker-up   # PostgreSQLを起動"
	@echo "  2. make db-migrate  # データベースマイグレーション"
	@echo "  3. make run-backend # バックエンドを起動（別ターミナル）"
	@echo "  4. make run-frontend # フロントエンドを起動（別ターミナル）"

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

lint: ## golangci-lintを実行
	@which golangci-lint > /dev/null || (echo "golangci-lintがインストールされていません。インストール: https://golangci-lint.run/welcome/install/" && exit 1)
	golangci-lint run --timeout=5m

fmt: ## gofmtでコードをフォーマット
	gofmt -s -w .
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .

vet: ## go vetを実行
	go vet ./...

check-fmt: ## フォーマットチェック（CI用）
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "以下のファイルがgofmtでフォーマットされていません:"; \
		gofmt -l .; \
		exit 1; \
	fi

check-imports: ## goimportsチェック（CI用）
	@go install golang.org/x/tools/cmd/goimports@latest
	@if [ -n "$$(goimports -l .)" ]; then \
		echo "以下のファイルがgoimportsでフォーマットされていません:"; \
		goimports -l .; \
		exit 1; \
	fi

modernize: ## modernize analyzerでコードを自動修正
	@echo "modernize analyzerを実行してコードを現代化します..."
	@go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test ./...

modernize-check: ## modernize analyzerでチェックのみ実行
	@echo "modernize analyzerでチェック中..."
	@go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest ./...

ci-test: lint vet check-fmt check-imports modernize-check test-backend ## CI用のすべてのチェックを実行
	@echo "すべてのチェックが完了しました"
