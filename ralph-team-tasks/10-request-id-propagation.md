このプロジェクトの **リクエストIDのコンテキスト伝播** を強化してください。

## 要件

- 既存の logger ミドルウェアが生成するリクエストIDを全レイヤーで利用可能にする
- usecase, command, queryservice の各メソッドに `context.Context` を渡す
- ログ出力時にリクエストIDを含める
- レスポンスヘッダー `X-Request-ID` にリクエストIDを含める

## 実装ガイド

1. `internal/pkg/logger/` のリクエストID関連を確認・活用
   - `RequestIDFromContext(ctx)` でIDを取得するヘルパーがあるか確認
   - なければ作成

2. レスポンスヘッダーの追加
   - logger ミドルウェアで `w.Header().Set("X-Request-ID", requestID)` を追加

3. 各レイヤーの `context.Context` 対応
   - `internal/queryservice/` の各メソッドに `ctx context.Context` パラメータを追加（なければ）
   - `internal/command/` の各メソッドに `ctx context.Context` パラメータを追加（なければ）
   - `internal/usecase/` の各メソッドに `ctx context.Context` パラメータを追加（なければ）
   - handler から `r.Context()` を渡す

4. ログ出力の強化
   - `slog.With("request_id", requestID)` パターンで各レイヤーのログにリクエストIDを含める
   - 既存の `slog` ロガーのコンテキスト対応を活用

## 注意

- 既存コードの変更が広範囲になる可能性あり。コンパイルエラーを1つずつ解消すること
- テストも `context.Background()` を渡すよう修正が必要

## 検証

1. `go build ./...` が通ること
2. `go vet ./...` が通ること
3. `go test ./...` が通ること
4. golangci-lint が通ること（あれば）

## 完了条件

- 上記の検証がすべて通ったら、以下を実行:
  1. `git checkout -b feature/request-id-propagation` でブランチ作成（既にいれば不要）
  2. 変更をコミット & プッシュ
  3. `gh pr create --title "refactor: リクエストIDのコンテキスト伝播強化" --body "..." --base main` でPR作成
  4. `gh pr merge --auto --squash` で自動マージ設定
  5. 完了を宣言: <promise>DONE</promise>

全検証が通り、PRが作成され自動マージが設定されるまで <promise>DONE</promise> を出力しないでください。
