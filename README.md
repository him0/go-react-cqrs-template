# Go Vite Spec Kit Sample

GolangのWebサーバー（DDD構成）+ Vite React + OpenAPI + Orvalを使用したサンプルプロジェクト

## 構成

### バックエンド (Go)
- **アーキテクチャ**: Domain-Driven Design (DDD)
  - `domain`: ドメイン層（エンティティ、値オブジェクト、リポジトリインターフェース）
  - `application`: アプリケーション層（ユースケース、サービス）
  - `infrastructure`: インフラ層（リポジトリ実装）
  - `interfaces`: インターフェース層（HTTPハンドラー）
- **ルーター**: chi
- **API仕様**: OpenAPI 3.0

### フロントエンド (React)
- **ビルドツール**: Vite
- **UI**: React 18
- **状態管理**: React Query (@tanstack/react-query)
- **HTTPクライアント**: Axios
- **パッケージマネージャー**: pnpm

### コード生成
- **バックエンド**: openapi-generator (Go server code)
- **フロントエンド**: Orval (TypeScript types + React Query hooks)

## プロジェクト構造

```
.
├── api/                    # OpenAPI仕様
│   ├── openapi.yaml
│   └── generator-config.yaml
├── cmd/
│   └── server/            # アプリケーションエントリーポイント
│       └── main.go
├── internal/
│   ├── domain/            # ドメイン層
│   │   └── user.go
│   ├── application/       # アプリケーション層
│   │   └── user_service.go
│   ├── infrastructure/    # インフラ層
│   │   └── inmemory_user_repository.go
│   └── interfaces/        # インターフェース層
│       └── http/
│           └── user_handler.go
├── web/                   # フロントエンド
│   ├── src/
│   │   ├── api/          # Orvalで生成されたAPIコード
│   │   ├── components/   # Reactコンポーネント
│   │   ├── hooks/        # カスタムフック
│   │   └── pages/        # ページコンポーネント
│   ├── package.json
│   ├── orval.config.ts   # Orval設定
│   └── vite.config.ts    # Vite設定
├── scripts/              # ビルドスクリプト
└── go.mod
```

## セットアップ

### 前提条件
- Go 1.21+
- Node.js 18+
- pnpm

### インストール

1. Goの依存関係をインストール:
```bash
go mod download
```

2. フロントエンドの依存関係をインストール:
```bash
cd web
pnpm install
```

### コード生成

#### フロントエンドのAPIコード生成 (Orval)
```bash
cd web
pnpm run generate:api
```

これにより、`web/src/api/generated`に以下が生成されます:
- TypeScript型定義
- React Query hooks
- Axiosクライアント

#### バックエンドのコード生成 (openapi-generator) ※オプション
```bash
./scripts/generate-api.sh
```

## 開発

### バックエンド起動
```bash
go run cmd/server/main.go
```
サーバーは http://localhost:8080 で起動します。

### フロントエンド起動
```bash
cd web
pnpm run dev
```
開発サーバーは http://localhost:3000 で起動します。

## API エンドポイント

### ユーザー管理
- `GET /api/v1/users` - ユーザー一覧取得
  - クエリパラメータ: `limit`, `offset`
- `POST /api/v1/users` - ユーザー作成
- `GET /api/v1/users/{userId}` - ユーザー詳細取得
- `PUT /api/v1/users/{userId}` - ユーザー更新
- `DELETE /api/v1/users/{userId}` - ユーザー削除

### リクエスト例

ユーザー作成:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

ユーザー一覧取得:
```bash
curl http://localhost:8080/api/v1/users?limit=10&offset=0
```

## DDD レイヤーの説明

### Domain層 (`internal/domain`)
- ビジネスロジックの中核
- エンティティ、値オブジェクト、リポジトリインターフェースを定義
- 他の層に依存しない

### Application層 (`internal/application`)
- ユースケースの実装
- ドメインオブジェクトを組み合わせてビジネスフローを実現
- トランザクション境界を定義

### Infrastructure層 (`internal/infrastructure`)
- 外部システムとのインターフェース実装
- データベース、外部API、ファイルシステムなど
- リポジトリインターフェースの具体実装

### Interface層 (`internal/interfaces`)
- 外部からのリクエストを受け付ける
- HTTPハンドラー、CLIコマンドなど
- リクエストのバリデーション、レスポンスの整形

## Orval の使い方

Orvalは OpenAPI 仕様から TypeScript のコードを自動生成します。

### 設定ファイル (`web/orval.config.ts`)
```typescript
export default defineConfig({
  api: {
    input: {
      target: '../api/openapi.yaml',  // OpenAPI仕様のパス
    },
    output: {
      mode: 'tags-split',             // タグごとにファイル分割
      target: './src/api/generated',  // 出力先
      client: 'react-query',          // React Queryフック生成
      override: {
        mutator: {
          path: './src/api/axios-instance.ts',
          name: 'customInstance',
        },
      },
    },
  },
});
```

### 生成されるコード
- **型定義**: OpenAPIスキーマから TypeScript の型を生成
- **React Query hooks**: `useListUsers`, `useCreateUser` など
- **カスタムインスタンス**: Axios インスタンスの設定

### 使用例
```typescript
import { useListUsers, useCreateUser } from '@/api/generated/users';

function UserList() {
  const { data, isLoading } = useListUsers({ limit: 10, offset: 0 });
  const createUser = useCreateUser();

  const handleCreate = async () => {
    await createUser.mutateAsync({
      data: { name: 'John', email: 'john@example.com' }
    });
  };

  // ...
}
```

## ビルド

### バックエンド
```bash
go build -o bin/server cmd/server/main.go
```

### フロントエンド
```bash
cd web
pnpm run build
```

## ライセンス

MIT
