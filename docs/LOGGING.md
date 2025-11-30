# ロギング

このプロジェクトでは、Go標準ライブラリの`log/slog`を使用した構造化ロギングを実装しています。

## 機能

- **構造化ログ**: JSON形式でログを出力
- **ログレベル**: DEBUG、INFO、WARN、ERROR
- **リクエストトレーシング**: 自動的にリクエストIDを付与
- **スタックトレース**: エラー時にスタックトレースを記録
- **カスタムエラー**: ユーザー向けと内部向けのメッセージを分離

## 環境変数

### LOG_LEVEL

ログレベルを制御します。以下の値が使用できます：

- `DEBUG`: すべてのログを出力
- `INFO`: INFO以上のログを出力（デフォルト）
- `WARN`: WARN以上のログを出力
- `ERROR`: ERRORログのみ出力

```bash
export LOG_LEVEL=DEBUG
```

### LOG_FORMAT

ログフォーマットを制御します：

- `json`: JSON形式で出力（デフォルト、本番環境推奨）
- `text`: テキスト形式で出力（開発環境向け）

```bash
export LOG_FORMAT=text
```

## 使用例

### 基本的な使い方

```go
import (
    "log/slog"
    "github.com/example/go-react-cqrs-template/internal/pkg/logger"
)

// アプリケーション起動時にロガーをセットアップ
log := logger.Setup()

// ログを出力
log.Info("server starting", slog.String("port", "8080"))
log.Error("failed to connect", slog.String("error", err.Error()))
```

### HTTPハンドラーでの使用

リクエストコンテキストからロガーを取得すると、自動的にリクエストIDが付与されます：

```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // コンテキストからロガーを取得（リクエストID付き）
    log := logger.FromContext(r.Context())

    log.Info("creating user",
        slog.String("email", email),
    )

    // エラー時
    log.Error("failed to create user",
        slog.String("error", err.Error()),
        slog.String("email", email),
    )
}
```

### カスタムエラーの使用

スタックトレースとユーザー向けメッセージを含むカスタムエラー：

```go
import (
    "github.com/example/go-react-cqrs-template/internal/pkg/errors"
    "github.com/example/go-react-cqrs-template/internal/pkg/logger"
)

// エラーを作成
appErr := errors.NotFound("user", "ユーザーが見つかりませんでした")

// エラーをログに記録
logger.LogError(log, appErr, "user operation failed")

// HTTPレスポンスを返す
w.WriteHeader(appErr.StatusCode())
json.NewEncoder(w).Encode(map[string]string{
    "message": appErr.UserMessage(),
})
```

### エラーレベル

カスタムエラーは以下のレベルをサポートします：

- `LevelInfo`: 通常のビジネスロジックエラー（404、400など）
- `LevelWarning`: 注意が必要だが処理は継続可能
- `LevelError`: 予期しないエラー（500など）
- `LevelCritical`: システムに深刻な影響

```go
// 様々なエラーの作成
err1 := errors.NotFound("user", "ユーザーが見つかりません")  // LevelInfo
err2 := errors.BadRequest("invalid email", "メールアドレスが不正です")  // LevelInfo
err3 := errors.Internal(err, "データベースエラー")  // LevelError
err4 := errors.Conflict("email exists", "メールアドレスが既に使用されています")  // LevelWarning
```

## ログ出力例

### JSON形式（本番環境）

```json
{
  "time": "2025-11-03T10:15:30.123456Z",
  "level": "INFO",
  "msg": "request completed",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/v1/users",
  "status": 201,
  "duration": "45ms",
  "duration_ms": 45
}
```

### エラーログ（スタックトレース付き）

```json
{
  "time": "2025-11-03T10:15:30.123456Z",
  "level": "ERROR",
  "msg": "failed to create user",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "error": "database connection failed",
  "user_message": "サーバー内部エラーが発生しました",
  "status_code": 500,
  "error_level": "ERROR",
  "stack_trace": [
    "go-react-cqrs-template/internal/application/user_service.go:45 application.(*UserService).CreateUser",
    "go-react-cqrs-template/internal/interfaces/http/user_handler.go:135 http.(*UserHandler).CreateUser",
    "..."
  ]
}
```

## アーキテクチャ

各層でのロギング：

1. **ミドルウェア層**: すべてのHTTPリクエストを自動的にログ記録、リクエストIDを付与
2. **ハンドラー層**: リクエスト処理の開始・終了、エラーをログ記録
3. **サービス層**: ビジネスロジックの重要な処理をログ記録
4. **リポジトリ層**: データアクセスの操作をログ記録

## 実装ファイル

- `internal/pkg/logger/logger.go`: ロガーのセットアップ
- `internal/pkg/logger/context.go`: コンテキスト関連のユーティリティ
- `internal/pkg/logger/middleware.go`: HTTPミドルウェア
- `internal/pkg/logger/error.go`: エラーロギングヘルパー
- `internal/pkg/errors/errors.go`: カスタムエラー型
