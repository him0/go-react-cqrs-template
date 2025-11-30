# Claude Code Project Guide

## Project Overview
Go + React のサンプルプロジェクト。ドメイン駆動設計とCQRS的パターンを採用。

## Tech Stack
- Go 1.x, chi v5, PostgreSQL, sqlc, slog
- OpenAPI (TypeSpec)

## Directory Structure
```
internal/
  domain/         # ドメインモデル (User, UserLog)
  command/        # Write操作 (トランザクション内で使用)
  queryservice/   # Read操作
  usecase/        # ビジネスロジック
  handler/        # HTTPハンドラー
  infrastructure/ # DB接続, DAO (sqlc生成)
  pkg/
    logger/       # slogベースのロガー
    errors/       # AppError
db/
  schema/         # テーブル定義
  queries/        # sqlc用SQL
```

## Architecture Patterns
- **CQRS分離**: Read(queryservice) / Write(command)
- **トランザクション**: `txManager.RunInTransaction()` でusecase層から制御
- **FOR UPDATE**: Race Condition防止用の行ロック

## Commands
```bash
go build ./...    # ビルド
go test ./...     # テスト
sqlc generate     # DAO生成 (要sqlcインストール)
```

## Troubleshooting

### mise
```bash
# go コマンドが見つからない場合
eval "$(mise activate zsh)" && go build ./...

# Config files are not trusted の場合
mise trust && mise install
```

## Adding New Features
1. `db/schema/` にテーブル定義追加
2. `db/queries/` にSQL追加
3. `internal/domain/` にモデル追加
4. `sqlc generate` でDAO生成
5. `internal/command/` にWrite関数追加
6. `internal/usecase/` にユースケース追加
