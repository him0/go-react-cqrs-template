# Claude Code Project Guide

## Project Overview
Go + React のサンプルプロジェクト。ドメイン駆動設計とCQRS的パターンを採用。

## Tech Stack

### Backend
- Go 1.x, chi v5, PostgreSQL, sqlc, slog

### Frontend
- React 18, TypeScript, Vite
- TanStack Router, TanStack Query
- Tailwind CSS, shadcn/ui
- Orval (APIクライアント自動生成)
- Vitest, Testing Library

### API定義
- TypeSpec → OpenAPI → Go/TypeScript コード生成

## Commands
`task --list` で利用可能なコマンドを確認。詳細は `Taskfile.yml` を参照。
フロントエンドは `web/package.json` の scripts を参照。

## Directory Structure

### Backend
```
internal/
  domain/         # ドメインモデル (User, UserLog)
  command/        # Write操作 (トランザクション内で使用)
  queryservice/   # Read操作
  usecase/        # ビジネスロジック
  handler/        # HTTPハンドラー
    validation/   # OpenAPIバリデーションミドルウェア
  infrastructure/ # DB接続, DAO (sqlc生成)
  pkg/
    logger/       # slogベースのロガー
    errors/       # AppError
pkg/
  generated/      # oapi-codegen生成コード
db/
  schema/         # テーブル定義 (psqldef)
  queries/        # sqlc用SQL
typespec/         # API定義 (TypeSpec)
openapi/          # 生成されたOpenAPI仕様 + embed.go (バリデーション用)
```

### Frontend
```
web/
  src/
    api/
      generated/  # Orvalで自動生成 (モデル, hooks)
      axios-instance.ts
    components/   # UIコンポーネント
    routes/       # TanStack Routerのページ
    lib/          # ユーティリティ
    test/         # テストユーティリティ
```

## Architecture Patterns

### Backend
- **CQRS分離**: Read(queryservice) / Write(command)
- **トランザクション**: `txManager.RunInTransaction()` でusecase層から制御
- **FOR UPDATE**: Race Condition防止用の行ロック
- **APIバリデーション**: OpenAPI定義から自動バリデーション (kin-openapi)

### Frontend
- **データフェッチ**: TanStack Query + Orval生成のhooks
- **ルーティング**: TanStack Router (ファイルベース)
- **スタイリング**: Tailwind CSS + shadcn/uiコンポーネント

## Troubleshooting

### mise
```bash
# go コマンドが見つからない場合
eval "$(mise activate zsh)" && go build ./...

# Config files are not trusted の場合
mise trust && mise install
```

## Adding New Features

### Backend
1. `db/schema/` にテーブル定義追加
2. `db/queries/` にSQL追加
3. `internal/domain/` にモデル追加
4. DAO生成 (task generate:dao)
5. `internal/command/` にWrite関数追加
6. `internal/usecase/` にユースケース追加

### API追加時
1. `typespec/` にAPI定義追加
2. コード生成 (task generate)
3. `internal/handler/` にハンドラー実装
4. `web/` でOrval生成 (npm run generate:api)

### Frontend
1. `web/src/components/` にコンポーネント追加
2. `web/src/routes/` にページ追加
3. 必要に応じてOrval生成のhooksを使用
