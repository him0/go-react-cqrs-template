#!/usr/bin/env bash
# Ralph Team 用タスク一覧を表示
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
TASKS_DIR="$REPO_ROOT/ralph-team-tasks"

echo "=== Ralph Team タスク一覧 ==="
echo ""
echo "  #   | ファイル                        | 概要"
echo "  ----|--------------------------------|---------------------------"

for f in "$TASKS_DIR"/*.md; do
  filename=$(basename "$f")
  num=$(echo "$filename" | sed 's/-.*//')
  # 1行目から概要を取得（** で囲まれた部分）
  summary=$(head -1 "$f" | sed 's/.*\*\*\(.*\)\*\*.*/\1/')
  printf "  %-3s | %-30s | %s\n" "$num" "$filename" "$summary"
done

echo ""
echo "使い方:"
echo "  task ralph:show NUM=01           # タスク内容を表示"
echo "  task ralph:status                # PR/worktree 状況を確認"
echo "  /ralph-team                      # Claude Code で並列実行"
echo "  /ralph-team --tasks 01,02,03     # 指定タスクのみ実行"
echo "  /ralph-team --dry-run            # 実行計画を確認"
echo ""
