# アーキテクチャガイド

このプロジェクトは、ドメイン駆動設計（DDD）とCQRS（Command Query Responsibility Segregation）パターンを採用したGoアプリケーションです。

## ディレクトリ構成

```
internal/
├── domain/           # ドメイン層 - ビジネスロジックとエンティティ
│   └── user.go
├── command/          # コマンド層 - 書き込み操作（Create, Update, Delete）
│   └── user_command.go
├── queryservice/     # クエリサービス層 - 読み取り操作（Read）
│   └── user_query_service.go
├── usecase/          # ユースケース層 - アプリケーションロジック
│   └── user_usecase.go
├── handler/          # ハンドラー層 - HTTPリクエスト処理
│   └── user_handler.go
└── infrastructure/   # インフラ層 - データベース接続など
    └── database.go
```

## 各層の責務

### 1. Domain層 (domain/)
- ビジネスロジックを含むエンティティとバリューオブジェクト
- ドメインルールのバリデーション
- リポジトリインターフェースの定義（従来の実装、現在は使用していない）

**例**: `user.go` - ユーザーエンティティと生成・更新ロジック

### 2. Command層 (command/)
- **書き込み操作**を担当（Create, Update, Delete）
- データベースへの永続化を実行
- トランザクション管理

**例**: `user_command.go` - ユーザーの作成・更新・削除

### 3. QueryService層 (queryservice/)
- **読み取り操作**を担当（Read）
- データベースからのクエリ実行
- ページネーションやフィルタリング

**例**: `user_query_service.go` - ユーザーの検索・一覧取得

### 4. Usecase層 (usecase/)
- ビジネスユースケースの実装
- CommandとQueryServiceを組み合わせて使用
- ビジネスルールの適用（重複チェックなど）

**例**: `user_usecase.go` - ユーザー登録、更新、削除などのユースケース

### 5. Handler層 (handler/)
- HTTPリクエスト・レスポンスの処理
- リクエストのバリデーション
- JSONのシリアライズ・デシリアライズ

**例**: `user_handler.go` - ユーザーAPI のエンドポイント

### 6. Infrastructure層 (infrastructure/)
- データベース接続
- 外部サービスとの連携
- 技術的な詳細の実装

**例**: `database.go` - PostgreSQL接続設定

## データフロー

```
HTTPリクエスト
    ↓
Handler層
    ↓
Usecase層
    ↓
Command層 / QueryService層
    ↓
Database
```

## CQRS パターン

このアプリケーションではCQRSパターンを採用しています：

- **Command（書き込み）**: `command/` パッケージが担当
  - データの変更操作
  - トランザクション管理

- **Query（読み取り）**: `queryservice/` パッケージが担当
  - データの取得操作
  - パフォーマンス最適化されたクエリ

### メリット

1. **関心の分離**: 読み取りと書き込みの責務が明確
2. **スケーラビリティ**: 読み取りと書き込みを独立してスケール可能
3. **最適化**: それぞれの操作に最適な実装が可能

## データベース管理

### sqldef を使用したスキーマ管理

スキーマファイルは `db/schema/schema.sql` に定義されています。

```bash
# マイグレーション実行
make db-migrate

# ドライラン（変更内容の確認）
make db-dry-run
```

## セットアップと実行

### 1. 依存関係のインストール

```bash
make setup
```

### 2. PostgreSQLの起動

```bash
make docker-up
```

### 3. データベースマイグレーション

```bash
make db-migrate
```

### 4. アプリケーションの起動

```bash
make run-backend
```

## 環境変数

`.env.example` を参考に `.env` ファイルを作成してください：

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=app_db
DB_SSLMODE=disable
PORT=8080
```

## 新機能の追加方法

### 1. ドメインモデルの作成 (domain/)

```go
type Product struct {
    ID    string
    Name  string
    Price int
}
```

### 2. Commandの実装 (command/)

```go
type ProductCommand struct {
    db *sql.DB
}

func (c *ProductCommand) Create(ctx context.Context, product *domain.Product) error {
    // 書き込みロジック
}
```

### 3. QueryServiceの実装 (queryservice/)

```go
type ProductQueryService struct {
    db *sql.DB
}

func (q *ProductQueryService) FindByID(ctx context.Context, id string) (*domain.Product, error) {
    // 読み取りロジック
}
```

### 4. Usecaseの実装 (usecase/)

```go
type ProductUsecase struct {
    productCommand      *command.ProductCommand
    productQueryService *queryservice.ProductQueryService
}
```

### 5. Handlerの実装 (handler/)

```go
type ProductHandler struct {
    productUsecase *usecase.ProductUsecase
}
```

### 6. main.goでの登録

```go
productCommand := command.NewProductCommand(db)
productQueryService := queryservice.NewProductQueryService(db)
productUsecase := usecase.NewProductUsecase(productCommand, productQueryService)
productHandler := handler.NewProductHandler(productUsecase)
```

## テスト

```bash
# すべてのテストを実行
make test

# カバレッジ付き
make test-coverage
```

## ベストプラクティス

1. **単一責任の原則**: 各層は明確な責務を持つ
2. **依存性の方向**: 外側から内側への依存（Handler → Usecase → Command/QueryService）
3. **ドメインロジック**: ビジネスルールはDomain層に集約
4. **エラーハンドリング**: 適切なエラーメッセージとHTTPステータスコード
5. **Context の利用**: タイムアウトとキャンセル処理のため
