このプロジェクトに **環境別設定管理** を追加してください。

## 要件

- 環境変数ベースの設定管理を構造化
- dev / staging / prod の設定を明確に分離
- `.env.example` で必要な環境変数を文書化
- `envconfig` や `kelseyhightower/envconfig` ライブラリを使用

## 実装ガイド

1. `internal/config/config.go` を作成
   ```go
   type Config struct {
       Server   ServerConfig
       Database DatabaseConfig
       Log      LogConfig
   }

   type ServerConfig struct {
       Port            string `envconfig:"PORT" default:"8080"`
       ShutdownTimeout int    `envconfig:"SHUTDOWN_TIMEOUT" default:"30"`
       CORSOrigins     string `envconfig:"CORS_ORIGINS" default:"http://localhost:3000"`
   }

   type DatabaseConfig struct {
       Host     string `envconfig:"DB_HOST" default:"localhost"`
       Port     int    `envconfig:"DB_PORT" default:"5432"`
       User     string `envconfig:"DB_USER" default:"postgres"`
       Password string `envconfig:"DB_PASSWORD" default:"postgres"`
       DBName   string `envconfig:"DB_NAME" default:"app_db"`
       SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
   }

   type LogConfig struct {
       Level  string `envconfig:"LOG_LEVEL" default:"info"`
       Format string `envconfig:"LOG_FORMAT" default:"json"`
   }
   ```

2. `cmd/server/main.go` で `config.Load()` を使用して既存の `getEnv()` 呼び出しを置き換え

3. `.env.example` を作成（全環境変数とデフォルト値を記載）

4. テスト `internal/config/config_test.go`
   - デフォルト値のテスト
   - 環境変数の上書きテスト

## 検証

1. `go build ./...` が通ること
2. `go vet ./...` が通ること
3. `go test ./internal/config/...` が通ること
4. golangci-lint が通ること（あれば）

## 完了条件

- 上記の検証がすべて通ったら、以下を実行:
  1. `git checkout -b feature/env-config` でブランチ作成（既にいれば不要）
  2. 変更をコミット & プッシュ
  3. `gh pr create --title "feat: 環境別設定管理追加" --body "..." --base main` でPR作成
  4. `gh pr merge --auto --squash` で自動マージ設定
  5. 完了を宣言: <promise>DONE</promise>

全検証が通り、PRが作成され自動マージが設定されるまで <promise>DONE</promise> を出力しないでください。
