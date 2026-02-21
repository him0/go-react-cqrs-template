このプロジェクトに **セキュリティヘッダーミドルウェア** を追加してください。

## 要件

以下のHTTPレスポンスヘッダーを全レスポンスに付与:

- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: strict-origin-when-cross-origin`
- `Content-Security-Policy: default-src 'self'`（開発時は緩めにする）
- `Permissions-Policy: camera=(), microphone=(), geolocation=()`

## 実装ガイド

- `internal/handler/middleware/security.go` にミドルウェアを作成
- `chi.Router` の `r.Use()` で適用
- `cmd/server/main.go` に登録（CORS の後、Recoverer の前あたり）
- CSP は環境変数 `CSP_POLICY` で上書き可能にする（デフォルト値あり）
- テストを `internal/handler/middleware/security_test.go` に書く（httptest使用、各ヘッダーの存在を確認）

## 検証

1. `go build ./...` が通ること
2. `go vet ./...` が通ること
3. `go test ./internal/handler/middleware/...` が通ること
4. golangci-lint が通ること（あれば）

## 完了条件

- 上記の検証がすべて通ったら、以下を実行:
  1. `git checkout -b feature/security-headers` でブランチ作成（既にいれば不要）
  2. 変更をコミット & プッシュ
  3. `gh pr create --title "feat: セキュリティヘッダーミドルウェア追加" --body "..." --base main` でPR作成
  4. `gh pr merge --auto --squash` で自動マージ設定
  5. 完了を宣言: <promise>DONE</promise>

全検証が通り、PRが作成され自動マージが設定されるまで <promise>DONE</promise> を出力しないでください。
