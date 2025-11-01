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
- **ルーティング**: TanStack Router (@tanstack/react-router)
- **スタイリング**: Tailwind CSS + shadcn/ui
- **状態管理**: TanStack Query (@tanstack/react-query)
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

## Tailwind CSS & shadcn/ui

プロジェクトにはTailwind CSSとshadcn/uiが設定済みです。

### shadcn/uiコンポーネントの追加

shadcn/uiのコンポーネントを追加するには、pnpmを使用します：

```bash
cd web
pnpm dlx shadcn@latest add button
pnpm dlx shadcn@latest add card
pnpm dlx shadcn@latest add input
# など
```

### 利用可能なコンポーネント

shadcn/uiの全コンポーネントは [https://ui.shadcn.com/docs/components](https://ui.shadcn.com/docs/components) で確認できます。

### 設定ファイル

- `tailwind.config.js`: Tailwind CSS設定
- `components.json`: shadcn/ui設定
- `src/lib/utils.ts`: ユーティリティ関数（`cn`など）

### 使用例

```typescript
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

function MyComponent() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Hello World</CardTitle>
      </CardHeader>
      <CardContent>
        <Button variant="default">Click me</Button>
      </CardContent>
    </Card>
  )
}
```

## TanStack Router

プロジェクトはファイルベースルーティングを使用するTanStack Routerで設定されています。

### ルート構造

ルートは `web/src/routes/` ディレクトリに配置されます：

```
web/src/routes/
├── __root.tsx        # ルートレイアウト（ナビゲーション等）
├── index.tsx         # / (ホームページ)
├── users.tsx         # /users
└── about.tsx         # /about
```

### 新しいルートの追加

新しいページを追加するには、`web/src/routes/` に新しいファイルを作成します：

```typescript
// web/src/routes/new-page.tsx
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/new-page')({
  component: NewPage,
})

function NewPage() {
  return <div>New Page</div>
}
```

### 型安全なナビゲーション

TanStack Routerは完全な型安全性を提供します：

```typescript
import { Link } from '@tanstack/react-router'

// 型安全なリンク
<Link to="/users">Users</Link>
<Link to="/about">About</Link>

// パラメータ付きルートも型安全
<Link to="/users/$userId" params={{ userId: '123' }}>User 123</Link>
```

### ルートツリーの自動生成

Viteプラグインが `routeTree.gen.ts` を自動生成します。このファイルは手動で編集しないでください。

## Orval & React Query

OrvalはOpenAPI仕様からReact Query（TanStack Query）のhooksを自動生成します。

### API コード生成

```bash
cd web
pnpm run generate:api
```

これにより、`web/src/api/generated/` に以下が生成されます：
- TypeScript型定義
- React Query hooks（`useListUsers`, `useCreateUser`, など）
- カスタムAxiosインスタンス

### 使用例

```typescript
import { useListUsers, useCreateUser } from '@/api/generated/users'

function UserList() {
  // GETリクエスト用のhook
  const { data, isLoading, error } = useListUsers({
    limit: 10,
    offset: 0,
  })

  // POSTリクエスト用のhook
  const createUser = useCreateUser()

  const handleCreate = async () => {
    await createUser.mutateAsync({
      data: { name: 'John', email: 'john@example.com' }
    })
  }

  if (isLoading) return <div>Loading...</div>
  if (error) return <div>Error: {error.message}</div>

  return (
    <div>
      <button onClick={handleCreate}>Create User</button>
      <ul>
        {data?.users.map(user => (
          <li key={user.id}>{user.name}</li>
        ))}
      </ul>
    </div>
  )
}
```

### React Query Devtools

開発時はReact Query Devtoolsが自動的に有効になり、クエリの状態を確認できます。

## テスト

このプロジェクトには包括的なテストスイートが含まれています。

### すべてのテストを実行

```bash
make test
```

### バックエンドテスト

Goのテストフレームワークを使用：

```bash
# テストを実行
make test-backend

# または直接
go test -v ./...

# カバレッジ付き
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
```

**テストファイル:**
- `internal/domain/user_test.go` - ドメインロジックのテスト
- `internal/application/user_service_test.go` - アプリケーションサービスのテスト
- `internal/infrastructure/inmemory_user_repository_test.go` - リポジトリのテスト

### フロントエンドテスト

Vitest + React Testing Libraryを使用：

```bash
# テストを実行
make test-frontend

# または直接
cd web
pnpm run test

# カバレッジ付き
pnpm run test:coverage

# UIモード（インタラクティブ）
pnpm run test:ui

# ウォッチモード
make test-watch
```

**テストファイル:**
- `web/src/lib/utils.test.ts` - ユーティリティ関数のテスト
- `web/src/components/Button.test.tsx` - コンポーネントのテスト

### カバレッジレポート

```bash
make test-coverage
```

カバレッジレポートは以下に生成されます：
- バックエンド: `coverage.out`
- フロントエンド: `web/coverage/`

### CI/CD

GitHub Actionsを使用した自動テスト：

- **バックエンドテスト**: Go 1.21でテスト実行、カバレッジレポート生成
- **フロントエンドテスト**: Node.js 20 + pnpmでテスト実行、カバレッジレポート生成
- **Lint**: golangci-lintによる静的解析

ワークフローファイル: `.github/workflows/ci.yml`

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
