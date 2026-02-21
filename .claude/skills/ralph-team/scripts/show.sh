#!/usr/bin/env bash
# Ralph Team タスクの内容を表示
set -euo pipefail

NUM="${1:?Usage: show.sh <NUM> (e.g. 01)}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
TASKS_DIR="$REPO_ROOT/ralph-team-tasks"

# 番号にマッチするファイルを検索
FILE=$(find "$TASKS_DIR" -name "${NUM}-*.md" -type f | head -1)

if [ -z "$FILE" ]; then
  echo "Error: タスク ${NUM} が見つかりません"
  echo "利用可能なタスク:"
  ls "$TASKS_DIR"/*.md | xargs -I{} basename {}
  exit 1
fi

echo "=== $(basename "$FILE") ==="
echo ""
cat "$FILE"
