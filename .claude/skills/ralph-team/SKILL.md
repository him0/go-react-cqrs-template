---
name: ralph-team
description: "Ralph Team: Agent Teams による並列タスク自動実行。ralph-team-tasks/*.md に定義されたタスクを worktree 隔離されたワーカーで並列実行し、各タスクを独立した PR として完成させる。"
disable-model-invocation: true
allowed-tools: Bash, Read, Write, Edit, Glob, Grep, Task, TaskCreate, TaskUpdate, TaskList, TaskGet, TeamCreate, TeamDelete, SendMessage
argument-hint: "[--concurrency N] [--tasks 01,02,03|all] [--dry-run]"
---

# Ralph Team リーダー

あなたは Ralph Team のリーダーです。`ralph-team-tasks/*.md` に定義されたタスクを Agent Teams で並列実行し、各タスクを独立した PR として完成させます。

## 引数の解析

ユーザーの入力 `$ARGUMENTS` から以下のオプションを解析してください:
- `--concurrency N` : 同時実行ワーカー数（デフォルト: 3）
- `--tasks "01,02,03"` または `all` : 実行するタスク番号（デフォルト: all）
- `--dry-run` : 実行計画の表示のみ（チーム作成しない）
- `--retry-failed` : 前回失敗したタスクも再実行対象にする

## 実行手順

### Step 1: タスクファイルの読み込み & 完了状態の確認

`ralph-team-tasks/*.md` を全て読み込み、各タスクの番号・名前・概要を把握する。

次に、ステートファイル `ralph-team-tasks/state.json` を読み込んで各タスクの完了状態を判定する。
ファイルが無ければ初回実行として全タスクを未着手とする。

**state.json のフォーマット**:
```json
{
  "tasks": {
    "01": { "status": "done", "pr": "#123", "finished_at": "2026-02-22T10:30:00Z" },
    "04": { "status": "done", "pr": "#124", "finished_at": "2026-02-22T10:45:00Z" },
    "05": { "status": "failed", "reason": "検証失敗（3回リトライ後）", "finished_at": "2026-02-22T11:00:00Z" }
  }
}
```

**判定ルール**:
- `status: "done"` → **完了**: スケジュールから除外
- `status: "failed"` → **失敗**: スケジュールから除外（`--retry-failed` で再実行可能）
- エントリなし → **未着手**: スケジュール対象

**ステート更新タイミング**:
- ワーカーが PR 作成 & auto-merge 設定完了 → リーダーが `"done"` を書き込み
- ワーカーが3回リトライ後も失敗 → リーダーが `"failed"` を書き込み

完了状態をユーザーに表示する:
```
=== タスク状態確認 ===
  [DONE]   04 - バックエンド CI (PR #123)
  [DONE]   08 - マルチステージ Dockerfile (PR #124)
  [FAILED] 05 - Error Boundary (検証失敗)
  [TODO]   01 - ヘルスチェックエンドポイント
  ...

  完了済み: 2 / 12, 失敗: 1, 残り: 9
```

### Step 2: 依存関係の確認（競合グループ）

以下のタスクは同じファイルを変更するため、同時に実行してはいけない:

| 競合グループ | タスク番号 | 理由 |
|---|---|---|
| `cmd/server/main.go` | 01, 02, 03, 09, 11 | サーバー起動コードを変更 |
| フロントエンドフォーム | 06, 12 | フォームコンポーネントを変更 |
| フロントエンドルート | 05, 06, 07 | ルート/レイアウトを変更 |
| 広範囲リファクタ | 10 | 内部全レイヤーに影響 |

### Step 3: スケジュール作成

concurrency に基づいてバッチを組む。競合グループ内のタスクは同じバッチに入れない。

**推奨スケジュール** (concurrency=3):
```
Batch 1: 04 (CI) + 08 (Docker) + 05 (ErrorBoundary)     ← 競合なし
Batch 2: 01 (Health) + 07 (DarkMode)                      ← 独立
Batch 3: 02 (Shutdown) + 12 (FormLib)                      ← 独立
Batch 4: 03 (Security) + 06 (Toast)                        ← 独立
Batch 5: 09 (RateLimit) + 11 (EnvConfig)                   ← main.go 競合だが前が完了済
Batch 6: 10 (RequestID)                                    ← 最後（広範囲リファクタ）
```

### Step 4: dry-run の場合

`--dry-run` が指定されている場合:
1. スケジュール（バッチ割り当て）を表示
2. 各バッチの実行順序と理由を説明
3. チーム作成せずに終了

### Step 5: チーム作成 & タスク登録

```
TeamCreate: team_name="ralph-team"
```

各タスクを TaskCreate で登録し、依存関係（blockedBy）を設定する:
- Batch 2 のタスクは Batch 1 の全タスク完了を待つ
- Batch N のタスクは Batch N-1 の全タスク完了を待つ

### Step 6: ワーカー生成（バッチ単位）

現在のバッチのタスクに対して、Task ツールでワーカーを生成する。
ワーカーは `.claude/agents/ralph-worker.md` で定義されたカスタムエージェントを使用し、`isolation: "worktree"` で git 隔離される:

```
Task(
  subagent_type: "ralph-worker",
  team_name: "ralph-team",
  name: "worker-{task-number}",
  prompt: "以下のタスクを実装してください。

タスク番号: {number}
タスクファイル: ralph-team-tasks/{filename}

タスクファイルの内容を読み込み、指示に従って実装・検証・PR作成・auto-merge まで完了してください。

完了したら TaskUpdate でタスクを completed にし、SendMessage でリーダーに PR URL を報告してください。"
)
```

**重要**: `ralph-worker` エージェントは `isolation: worktree` が設定されているため、Task ツール側で `isolation` を指定する必要はない。エージェント定義の frontmatter で自動的に worktree 隔離される。

### Step 7: ワーカー管理 & ステート更新

- ワーカーからのメッセージを受信して進捗を追跡
- ワーカーが成功を報告 → `ralph-team-tasks/state.json` に `"done"` + PR URL を書き込み
- ワーカーが失敗を報告 → `ralph-team-tasks/state.json` に `"failed"` + 理由を書き込み
- **ステートは各ワーカー完了時に即座に書き込む**（セッション断でも進捗が残る）
- ワーカーが完了したら shutdown_request を送信
- バッチ内の全ワーカーが完了したら、次バッチのワーカーを生成

### Step 8: 最終レポート

全バッチ完了後:
1. 各タスクの結果をサマリー表示（成功/失敗/PR URL）
2. TeamDelete でチームをクリーンアップ

## 出力フォーマット

### 起動時
```
=== Ralph Team 起動 ===
  対象タスク: {task_count} 件
  同時実行数: {concurrency}
  バッチ数: {batch_count}

  Batch 1: 04 (CI), 08 (Docker), 05 (ErrorBoundary)
  Batch 2: 01 (Health), 07 (DarkMode)
  ...
```

### 完了時
```
=== Ralph Team 完了レポート ===
  成功: {N} / {total}
  失敗: {N} / {total}

  [OK] 04 - バックエンド CI → PR #123 (merged)
  [OK] 08 - マルチステージ Dockerfile → PR #124 (merged)
  [NG] 05 - Error Boundary → 検証失敗（3回リトライ後）
  ...
```
