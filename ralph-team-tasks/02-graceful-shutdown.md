このプロジェクトに **Graceful Shutdown** を実装してください。

## 要件

- SIGINT / SIGTERM を受信したらサーバーを安全に停止
- 処理中のリクエストは完了を待つ（タイムアウト: 30秒）
- DB接続を安全にクローズ
- 停止プロセスをログに記録

## 実装ガイド

- `cmd/server/main.go` を修正
- `http.ListenAndServe` → `http.Server` + `server.Shutdown(ctx)` に変更
- `os/signal` + `context` パッケージを使用
- `signal.NotifyContext` でシグナルをキャッチ
- シャットダウン時に `log.Info("shutting down server...")` をログ出力
- DB Close は defer で既にあるのでそのまま

## 実装パターン

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

srv := &http.Server{Addr: ":" + port, Handler: r}

go func() {
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Error("server error", slog.String("error", err.Error()))
        os.Exit(1)
    }
}()

<-ctx.Done()
log.Info("shutting down server...")

shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := srv.Shutdown(shutdownCtx); err != nil {
    log.Error("server shutdown error", slog.String("error", err.Error()))
}
log.Info("server stopped")
```

## 検証

1. `go build ./...` が通ること
2. `go vet ./...` が通ること
3. golangci-lint が通ること（あれば）

## 完了条件

- 上記の検証がすべて通ったら、以下を実行:
  1. `git checkout -b feature/graceful-shutdown` でブランチ作成（既にいれば不要）
  2. 変更をコミット & プッシュ
  3. `gh pr create --title "feat: Graceful Shutdown実装" --body "..." --base main` でPR作成
  4. `gh pr merge --auto --squash` で自動マージ設定
  5. 完了を宣言: <promise>DONE</promise>

全検証が通り、PRが作成され自動マージが設定されるまで <promise>DONE</promise> を出力しないでください。
