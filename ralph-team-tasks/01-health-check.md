このプロジェクトに **ヘルスチェックエンドポイント** を追加してください。

## 要件

1. `/healthz` (liveness probe) - サーバーが生きているか
   - 常に `200 OK` + `{"status": "ok"}` を返す
2. `/readyz` (readiness probe) - リクエスト受付可能か
   - DB接続を `db.PingContext` で確認
   - 成功: `200` + `{"status": "ready"}`
   - 失敗: `503` + `{"status": "not_ready", "reason": "..."}`

## 実装ガイド

- `internal/handler/health_handler.go` にハンドラーを作成
- `cmd/server/main.go` の chi ルーターに `/healthz` と `/readyz` を登録（`/api/v1` の外側、バリデーションミドルウェア不要）
- ハンドラーは `*sql.DB` を受け取る
- テストを `internal/handler/health_handler_test.go` に書く（httptest使用）

## 検証

1. `go build ./...` が通ること
2. `go vet ./...` が通ること
3. `go test ./internal/handler/...` が通ること
4. golangci-lint が通ること（`which golangci-lint > /dev/null 2>&1 && golangci-lint run --timeout=5m` で確認、ない場合はスキップ）

## 完了条件

- 上記の検証がすべて通ったら、以下を実行:
  1. `git checkout -b feature/health-check` でブランチ作成（既にいれば不要）
  2. 変更をコミット & プッシュ
  3. `gh pr create --title "feat: ヘルスチェックエンドポイント追加" --body "..." --base main` でPR作成
  4. `gh pr merge --auto --squash` で自動マージ設定
  5. 完了を宣言: <promise>DONE</promise>

全検証が通り、PRが作成され自動マージが設定されるまで <promise>DONE</promise> を出力しないでください。
