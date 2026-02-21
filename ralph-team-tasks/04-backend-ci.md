このプロジェクトに **バックエンドCI用GitHub Actionsワークフロー** を追加してください。

## 要件

- PR と main への push 時に実行
- Go のビルド、lint、vet、テストを実行
- フロントエンドの既存 Playwright CI と並行して動作

## 実装ガイド

`.github/workflows/backend.yml` を作成:

```yaml
name: Backend CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  backend:
    runs-on: ubuntu-latest
    timeout-minutes: 10

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version-file: 'go.mod'
          cache: true

      - name: Install dependencies
        run: go mod download

      - name: Build
        run: go build ./...

      - name: Vet
        run: go vet ./...

      - name: Lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest
          args: --timeout=5m

      - name: Test
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: actions/upload-artifact@v4
        if: always()
        with:
          name: coverage-report
          path: coverage.out
          retention-days: 7
```

## 検証

1. YAML の構文が正しいこと（`python3 -c "import yaml; yaml.safe_load(open('.github/workflows/backend.yml'))"` あるいは手動確認）
2. 既存の `.github/workflows/playwright.yml` と競合しないこと

## 完了条件

- 上記の検証が通ったら、以下を実行:
  1. `git checkout -b feature/backend-ci` でブランチ作成（既にいれば不要）
  2. 変更をコミット & プッシュ
  3. `gh pr create --title "ci: バックエンドCI追加" --body "..." --base main` でPR作成
  4. `gh pr merge --auto --squash` で自動マージ設定
  5. 完了を宣言: <promise>DONE</promise>

全検証が通り、PRが作成され自動マージが設定されるまで <promise>DONE</promise> を出力しないでください。
