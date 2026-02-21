このプロジェクトに **レートリミットミドルウェア** を追加してください。

## 要件

- IPアドレスベースのレートリミット
- Token Bucket アルゴリズム
- 設定可能なリクエスト数/秒とバーストサイズ
- 429 Too Many Requests レスポンス
- 外部ライブラリ使用可（`golang.org/x/time/rate` 推奨）

## 実装ガイド

1. `internal/handler/middleware/ratelimit.go` を作成
   - `golang.org/x/time/rate` パッケージを使用
   - IPアドレスごとに `rate.Limiter` を管理（`sync.Map` で保持）
   - 古いエントリのクリーンアップ（最終アクセスからN分後に削除、goroutineで定期実行）
   - 設定: `RateLimitConfig { RequestsPerSecond float64, BurstSize int }`
   - デフォルト: 10 req/s, burst 20
   - 429 時に `Retry-After` ヘッダーを付与

2. `cmd/server/main.go` に登録
   - 環境変数 `RATE_LIMIT_RPS` と `RATE_LIMIT_BURST` で設定可能
   - `/healthz`, `/readyz` には適用しない（あれば）

3. テスト `internal/handler/middleware/ratelimit_test.go`
   - バースト内のリクエストが通ることを確認
   - バースト超過で 429 が返ることを確認
   - 異なるIPは独立してカウントされることを確認

## 検証

1. `go build ./...` が通ること
2. `go vet ./...` が通ること
3. `go test ./internal/handler/middleware/...` が通ること
4. golangci-lint が通ること（あれば）

## 完了条件

- 上記の検証がすべて通ったら、以下を実行:
  1. `git checkout -b feature/rate-limiter` でブランチ作成（既にいれば不要）
  2. 変更をコミット & プッシュ
  3. `gh pr create --title "feat: レートリミットミドルウェア追加" --body "..." --base main` でPR作成
  4. `gh pr merge --auto --squash` で自動マージ設定
  5. 完了を宣言: <promise>DONE</promise>

全検証が通り、PRが作成され自動マージが設定されるまで <promise>DONE</promise> を出力しないでください。
