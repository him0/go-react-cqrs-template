.PHONY: help install install-tools run-backend run-frontend generate-api generate-dao build clean test test-backend test-frontend test-coverage podman-up podman-down podman-logs podman-ps podman-restart db-migrate db-dry-run db-export db-generate-migration setup lint fmt vet check-fmt check-imports modernize modernize-check ci-test

help: ## ヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## 開発環境のセットアップ
	go mod download
	go install github.com/sqldef/sqldef/cmd/psqldef@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest
	go install golang.org/x/tools/cmd/goimports@latest
	pnpm install

install: ## 依存関係をインストール
	go mod download
	pnpm install

install-tools: ## tools.goに定義されたツールをインストール
	@echo "Installing tools from tools.go..."
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %@latest

podman-up: ## Podmanコンテナを起動
	podman-compose up -d

podman-down: ## Podmanコンテナを停止
	podman-compose down

podman-logs: ## Podmanログを表示
	podman-compose logs -f

podman-ps: ## Podmanコンテナの状態を表示
	podman-compose ps

podman-restart: ## Podmanコンテナを再起動
	podman-compose restart

db-migrate: ## psqldefを使用してデータベースマイグレーションを実行
	psqldef -U postgres -p 5432 -h localhost app_db --password=postgres --file=db/schema/schema.sql

db-dry-run: ## データベースマイグレーションのドライラン
	psqldef -U postgres -p 5432 -h localhost app_db --password=postgres --file=db/schema/schema.sql --dry-run

db-export: ## 現在のDBスキーマをエクスポート
	@echo "Exporting current database schema..."
	@psqldef -U postgres -p 5432 -h localhost app_db --password=postgres --export

db-generate-migration: ## スキーマ変更からマイグレーションファイルを生成
	@echo "Generating migration file from schema changes..."
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make db-generate-migration NAME=add_user_status"; \
		exit 1; \
	fi
	@TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	FILENAME="db/migrations/$${TIMESTAMP}_$(NAME).sql"; \
	psqldef -U postgres -p 5432 -h localhost app_db --password=postgres --file=db/schema/schema.sql --dry-run > $${FILENAME}; \
	if [ -s $${FILENAME} ]; then \
		echo "Migration file created: $${FILENAME}"; \
		cat $${FILENAME}; \
	else \
		echo "No schema changes detected. Removing empty file."; \
		rm $${FILENAME}; \
	fi

run-backend: ## バックエンドサーバーを起動
	go run cmd/server/main.go

run-frontend: ## フロントエンド開発サーバーを起動
	pnpm run dev

generate-api: ## APIコードを生成（フロントエンド）
	pnpm run generate:api

generate-dao: ## DAOコードをsqlcで生成
	@which sqlc > /dev/null || go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@sqlc generate

build-backend: ## バックエンドをビルド
	go build -o bin/server cmd/server/main.go

build-frontend: ## フロントエンドをビルド
	pnpm run build

build: build-backend build-frontend ## すべてをビルド

clean: ## ビルド成果物を削除
	rm -rf bin/
	rm -rf web/dist/
	rm -rf web/src/api/generated/
	rm -rf internal/infrastructure/dao/
	rm -rf web/coverage/
	rm -f coverage.out

dev: ## 開発環境を起動（バックエンド + フロントエンド）
	@echo "開発環境の起動手順："
	@echo "  1. make podman-up   # PostgreSQLを起動"
	@echo "  2. make db-migrate  # データベースマイグレーション"
	@echo "  3. make run-backend # バックエンドを起動（別ターミナル）"
	@echo "  4. make run-frontend # フロントエンドを起動（別ターミナル）"

test: test-backend test-frontend ## すべてのテストを実行

test-backend: ## バックエンドのテストを実行
	go test -v -race ./...

test-frontend: ## フロントエンドのテストを実行
	pnpm run test -- --run

test-coverage: ## カバレッジ付きでテストを実行
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	pnpm run test:coverage -- --run

test-watch: ## ウォッチモードでフロントエンドテストを実行
	pnpm run test

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
