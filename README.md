# Go Vite Spec Kit Sample

GolangのWebサーバー（DDD構成）+ Vite React + OpenAPI + Orvalを使用したサンプルプロジェクト

## 構成

### バックエンド (Go)
- **アーキテクチャ**: Domain-Driven Design (DDD) + CQRS
  - `domain`: ドメイン層（エンティティ、ビジネスロジック）
  - `command`: コマンド層（書き込み操作：Create, Update, Delete）
  - `queryservice`: クエリサービス層（読み取り操作：Read）
  - `usecase`: ユースケース層（アプリケーションロジック）
  - `handler`: ハンドラー層（HTTPリクエスト処理）
  - `infrastructure`: インフラ層（データベース接続、外部サービス連携）
- **データベース**: PostgreSQL
- **スキーマ管理**: sqldef (psqldef)
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
- **バックエンド DAO**: sqlc (型安全なDAO struct生成)
- **バックエンド API**: openapi-generator (Go server code)
- **フロントエンド**: Orval (TypeScript types + React Query hooks)

## プロジェクト構造

```
.
├── openapi/                # OpenAPI仕様
│   ├── openapi.yaml
│   └── generator-config.yaml
├── cmd/
│   └── server/            # アプリケーションエントリーポイント
│       └── main.go
├── db/
│   ├── schema/            # データベーススキーマ
│   │   └── schema.sql
│   └── queries/           # sqlcクエリ定義
│       └── users.sql
├── internal/
│   ├── domain/            # ドメイン層
│   │   └── user.go
│   ├── command/           # コマンド層（書き込み操作）
│   │   └── user_command.go
│   ├── queryservice/      # クエリサービス層（読み取り操作）
│   │   └── user_query_service.go
│   ├── usecase/           # ユースケース層
│   │   └── user_usecase.go
│   ├── handler/           # ハンドラー層
│   │   └── user_handler.go
│   └── infrastructure/    # インフラ層
│       ├── database.go
│       └── dao/           # sqlc生成DAO (自動生成)
│           ├── db.go
│           ├── models.go
│           └── users.sql.go
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
├── sqlc.yaml             # sqlc設定ファイル
└── go.mod
```

## セットアップ

### 前提条件
- Go 1.25+
- Node.js 24+
- pnpm
- Docker & Docker Compose
- mise (推奨) - バージョン管理ツール

### バージョン管理 (mise)

このプロジェクトでは、開発環境のバージョンを統一するために [mise](https://mise.jdx.dev/) を使用することを推奨します。

mise を使用する場合、プロジェクトディレクトリで以下を実行するだけで必要なツールがインストールされます：

```bash
mise install
```

これにより、`.mise.toml` で定義された以下のバージョンが自動的にインストールされます：
- Go 1.25.3
- Node.js 24.11.0 (Krypton LTS)
- pnpm 10.20.0

### インストール

1. すべての依存関係をインストール:
```bash
make setup
```

これにより以下がインストールされます：
- Go依存関係（`go mod download`）
- 開発ツール（psqldef、sqlc、goimportsなど）
- フロントエンド依存関係（pnpm）

または個別にインストール:

```bash
# Goの依存関係
go mod download

# 開発ツール（tools.goに定義）
make install-tools
# または個別に
go install github.com/sqldef/sqldef/cmd/psqldef@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# フロントエンドの依存関係
pnpm install
```

**ツール管理について:**
- 開発ツールは`tools.go`で依存関係として管理
- `go install`でインストール後、コマンドとして直接実行可能
- バージョンはgo.modで管理

### データベースのセットアップ

1. PostgreSQLをDockerで起動:
```bash
make docker-up
```

2. データベースマイグレーションを実行:
```bash
make db-migrate
```

マイグレーションをドライランで確認する場合:
```bash
make db-dry-run
```

### マイグレーションファイルの生成（差分管理）

sqldefは宣言的アプローチですが、マイグレーション履歴を残すこともできます：

#### 1. 現在のスキーマをエクスポート（ベースライン作成）
```bash
make db-export > db/migrations/001_baseline.sql
# または
./scripts/export-schema.sh db/migrations/001_baseline.sql
```

#### 2. スキーマ変更からマイグレーションファイルを生成
```bash
# 1. db/schema/schema.sql を編集してテーブルやカラムを追加

# 2. 差分マイグレーションを生成
make db-generate-migration NAME=add_user_status
# または
./scripts/generate-migration.sh add_user_status

# 3. 生成されたマイグレーションファイルを確認
# → db/migrations/20250101120000_add_user_status.sql

# 4. マイグレーションを適用
make db-migrate
```

**マイグレーションファイルの利点:**
- ✅ スキーマ変更履歴を追跡できる
- ✅ コードレビューでスキーマ変更を確認できる
- ✅ ロールバックの際の参考になる
- ✅ チーム開発でコンフリクトを減らせる

**ファイル命名規則:**
- `db/migrations/20250101120000_baseline.sql` - 初期スキーマ
- `db/migrations/20250102153000_add_user_status.sql` - ステータスカラム追加
- `db/migrations/20250103094500_create_posts_table.sql` - postsテーブル作成

### コード生成

#### バックエンドのDAO生成 (sqlc)
```bash
make generate-dao
```

これにより、`internal/infrastructure/dao`に以下が生成されます:
- テーブル定義に対応するGo struct（models.go）
- 型安全なクエリ関数（users.sql.go など）
- データベース接続インターフェース（db.go）

**sqlcの特徴:**
- SQLスキーマとクエリから型安全なGoコードを自動生成
- ORMではなく、生のSQLを使用できる
- Command/QueryServiceパターンとの親和性が高い
- ボイラープレートの削減
- `tools.go`で依存関係として管理

**クエリの追加方法:**
1. `db/queries/*.sql` にSQLクエリを追加
2. `make generate-dao` で再生成

**直接実行する場合:**
```bash
# sqlcがインストールされていない場合は先にインストール
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# 生成実行
sqlc generate
```

詳細は[db/queries/users.sql](db/queries/users.sql)と[sqlc.yaml](sqlc.yaml)を参照してください。

#### フロントエンドのAPIコード生成 (Orval)
```bash
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

### 開発環境の起動手順

1. PostgreSQLを起動:
```bash
make docker-up
```

2. データベースマイグレーション:
```bash
make db-migrate
```

3. バックエンド起動（別ターミナル）:
```bash
make run-backend
# または
go run cmd/server/main.go
```
サーバーは http://localhost:8080 で起動します。

4. フロントエンド起動（別ターミナル）:
```bash
make run-frontend
# または
pnpm run dev
```
開発サーバーは http://localhost:3000 で起動します。

### 環境変数

環境変数は`.env`ファイルで設定できます（`.env.example`を参照）:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=app_db
DB_SSLMODE=disable

# Server Configuration
PORT=8080
```

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

## アーキテクチャの詳細

このアプリケーションはDDD（ドメイン駆動設計）とCQRS（コマンドクエリ責務分離）パターンを採用しています。

### データフロー

```
HTTPリクエスト → Handler → Usecase → Command/QueryService → Database
```

### CQRS パターン

読み取りと書き込みの責務を分離：
- **Command** (`command/`): データ変更操作（Create/Update/Delete）とトランザクション管理
- **Query** (`queryservice/`): データ取得操作（Read）と最適化されたクエリ

メリット: 関心の分離、独立したスケーリング、各操作に最適な実装が可能

### 各層の責務

#### Domain層 (`internal/domain`)
- ビジネスロジックの中核
- エンティティとビジネスルールの定義

#### Command層 (`internal/command`)
- **書き込み操作**（Create, Update, Delete）を担当
- データベースへの永続化とトランザクション管理

#### QueryService層 (`internal/queryservice`)
- **読み取り操作**（Read）を担当
- ページネーションやフィルタリング

#### Usecase層 (`internal/usecase`)
- ビジネスユースケースの実装
- CommandとQueryServiceを組み合わせて使用

#### Handler層 (`internal/handler`)
- HTTPリクエスト・レスポンスの処理
- リクエストのバリデーションとJSON処理

#### Infrastructure層 (`internal/infrastructure`)
- データベース接続と外部サービス連携

### 新機能の追加手順

1. **Domain層**: エンティティとビジネスルールを定義
2. **Command層**: 書き込み操作を実装
3. **QueryService層**: 読み取り操作を実装
4. **Usecase層**: CommandとQueryServiceを組み合わせてビジネスロジックを実装
5. **Handler層**: HTTPエンドポイントを実装
6. **main.go**: 依存関係を注入して各層を接続

## Orval の使い方

Orvalは OpenAPI 仕様から TypeScript のコードを自動生成します。

### 設定ファイル (`web/orval.config.ts`)
```typescript
export default defineConfig({
  api: {
    input: {
      target: '../openapi/openapi.yaml',  // OpenAPI仕様のパス
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
pnpm --filter web dlx shadcn@latest add button
pnpm --filter web dlx shadcn@latest add card
pnpm --filter web dlx shadcn@latest add input
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

- **バージョン管理**: mise-actionを使用して `.mise.toml` からツールバージョンを自動読み込み
- **バックエンドテスト**: Goテスト実行、カバレッジレポート生成（PostgreSQL付き）
- **バックエンドLint**: golangci-lintによる静的解析
- **フロントエンドテスト**: Vitestでテスト実行、カバレッジレポート生成
- **フロントエンドLint**: ESLint + Prettierによるコード品質チェック
- **フロントエンドビルド**: 本番ビルドの検証

ワークフローファイル: `.github/workflows/ci.yml`

## ビルド

### バックエンド
```bash
go build -o bin/server cmd/server/main.go
```

### フロントエンド
```bash
pnpm run build
```

## ライセンス

MIT
