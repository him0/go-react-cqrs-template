# DAO (Data Access Object) - sqlc生成コード

このディレクトリには、sqlcによって自動生成されるDAO structとクエリ関数が配置されます。

## 生成方法

```bash
task generate:dao
```

または直接（sqlcコマンドを使用）：

```bash
# go.modのtool directiveで管理されているため、直接実行可能
go tool sqlc generate
```

**Note:** sqlcはgo.modの`tool`ディレクティブで依存関係として管理されています。

## 生成されるファイル

- `db.go` - データベース接続のインターフェース
- `models.go` - テーブル定義に対応するGo struct
- `users.sql.go` - `db/queries/users.sql`から生成されたクエリ関数

## 生成される主なstruct例

```go
// models.go
type User struct {
    ID        string    `db:"id" json:"id"`
    Name      string    `db:"name" json:"name"`
    Email     string    `db:"email" json:"email"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
```

## 生成されるクエリ関数例

```go
// users.sql.go
func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error)
func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error
func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error
func (q *Queries) DeleteUser(ctx context.Context, id string) error
```

## 使い方

```go
package main

import (
    "context"
    "github.com/example/go-react-cqrs-template/internal/infrastructure/dao"
    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    pool, _ := pgxpool.New(context.Background(), "postgres://...")
    queries := dao.New(pool)

    // ユーザーを取得
    user, err := queries.GetUserByID(context.Background(), "user-id")

    // ユーザー一覧を取得
    users, err := queries.ListUsers(context.Background(), dao.ListUsersParams{
        Limit:  10,
        Offset: 0,
    })
}
```

## 注意事項

- このディレクトリ内のファイルは自動生成されるため、直接編集しないでください
- クエリを変更したい場合は `db/queries/*.sql` を編集してください
- スキーマを変更した場合は `db/schema/schema.sql` を編集してください
