# OpenAPI Code Generation

このディレクトリには、API仕様とコード生成の設定が含まれています。

## ファイル構成

- `openapi.yaml` - OpenAPI 3.0仕様ファイル
- `generator-config.yaml` - oapi-codegenの設定ファイル（現在は使用していません）

## コード生成

### 必要なツール

プロジェクトでは[oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)を使用してGoのサーバーコードを生成します。

### 生成方法

```bash
# OpenAPI仕様からGoコードを生成
task generate:api

# またはすべてのコード生成を実行（DAO + API）
task generate
```

### 生成されるコード

生成されたコードは `pkg/generated/openapi/server.gen.go` に出力されます。

生成されるコードには以下が含まれます：

1. **型定義** - リクエスト/レスポンスの構造体
   - `User`
   - `CreateUserRequest`
   - `UpdateUserRequest`
   - `UserList`
   - `Error`

2. **ServerInterface** - 実装すべきハンドラーインターフェース
   ```go
   type ServerInterface interface {
       ListUsers(w http.ResponseWriter, r *http.Request, params ListUsersParams)
       CreateUser(w http.ResponseWriter, r *http.Request)
       GetUser(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID)
       UpdateUser(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID)
       DeleteUser(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID)
   }
   ```

3. **Chiルーターハンドラー** - chi-routerとの統合コード

## 使用方法

### 1. インターフェースを実装する

```go
package handler

import (
    "net/http"
    "github.com/example/go-react-spec-kit-sample/pkg/generated/openapi"
    openapi_types "github.com/oapi-codegen/runtime/types"
)

type UserHandler struct {
    // 依存関係
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request, params openapi.ListUsersParams) {
    // 実装
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 実装
}

// ... 他のメソッド
```

### 2. ルーターに登録する

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/example/go-react-spec-kit-sample/pkg/generated/openapi"
)

func SetupRouter() *chi.Mux {
    r := chi.NewRouter()

    handler := &UserHandler{}

    // OpenAPI生成のハンドラーを登録
    openapi.HandlerFromMux(handler, r)

    return r
}
```

## OpenAPI仕様の変更

`openapi.yaml`を編集した後、以下のコマンドでコードを再生成してください：

```bash
task generate:api
```

## 参考リンク

- [oapi-codegen Documentation](https://github.com/oapi-codegen/oapi-codegen)
- [OpenAPI Specification](https://swagger.io/specification/)
