#!/usr/bin/env bash
# Ralph Team タスクの PR / worktree 状況を確認
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
TASKS_DIR="$REPO_ROOT/ralph-team-tasks"

echo "=== Ralph Team ステータス ==="
echo ""

# PR一覧を取得
echo "--- Pull Requests ---"
if command -v gh &>/dev/null; then
  gh pr list --state all --limit 50 --json number,title,state,headRefName,url \
    --jq '.[] | select(.headRefName | startswith("feature/")) | "\(.state)\t#\(.number)\t\(.headRefName)\t\(.title)"' \
    | while IFS=$'\t' read -r state number branch title; do
      case "$state" in
        OPEN)    icon="[OPEN]   " ;;
        MERGED)  icon="[MERGED] " ;;
        CLOSED)  icon="[CLOSED] " ;;
        *)       icon="[?]      " ;;
      esac
      echo "  ${icon} ${number}  ${branch}  ${title}"
    done
  echo ""
else
  echo "  gh CLI が見つかりません。インストールしてください。"
  echo ""
fi

# Worktree 一覧
echo "--- Git Worktrees ---"
git worktree list --porcelain 2>/dev/null | while read -r line; do
  case "$line" in
    worktree\ *)
      path="${line#worktree }"
      ;;
    branch\ *)
      branch="${line#branch refs/heads/}"
      echo "  ${branch}  →  ${path}"
      ;;
  esac
done
echo ""

# タスクごとのサマリー
echo "--- タスクサマリー ---"
for f in "$TASKS_DIR"/*.md; do
  filename=$(basename "$f")
  num=$(echo "$filename" | sed 's/-.*//')
  name=$(echo "$filename" | sed 's/^[0-9]*-//; s/\.md$//')
  summary=$(head -1 "$f" | sed 's/.*\*\*\(.*\)\*\*.*/\1/')

  # feature/{name} ブランチの PR を検索
  pr_status="未着手"
  if command -v gh &>/dev/null; then
    pr_info=$(gh pr list --state all --head "feature/${name}" --json number,state --jq 'if length > 0 then .[0] | "\(.state) #\(.number)" else empty end' 2>/dev/null || true)
    if [ -n "$pr_info" ]; then
      pr_status="$pr_info"
    fi
  fi

  printf "  %s  %-30s  %s  (%s)\n" "$num" "$summary" "$pr_status" "$name"
done
echo ""
