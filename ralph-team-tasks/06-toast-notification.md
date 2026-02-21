このプロジェクトのフロントエンドに **Toast通知コンポーネント** を追加してください。

## 要件

- 操作成功/失敗のフィードバックをToastで表示
- sonner ライブラリを使用
- 既存のユーザーCRUD操作にToast通知を追加

## 実装ガイド

1. `cd web && npm install sonner`

2. `web/src/components/ui/sonner.tsx` を作成（shadcn/ui パターンに合わせる）
   - Toaster コンポーネントをエクスポート
   - Tailwind でスタイル調整

3. `web/src/main.tsx` または `web/src/routes/__root.tsx` に `<Toaster />` を追加

4. 既存のユーザー操作にToastを追加:
   - `web/src/components/UserCreateForm.tsx`: 作成成功時 `toast.success("ユーザーを作成しました")`
   - `web/src/components/UserEditForm.tsx`: 更新成功時 `toast.success("ユーザーを更新しました")`
   - `web/src/components/UserListItem.tsx`: 削除成功時 `toast.success("ユーザーを削除しました")`
   - エラー時: `toast.error("エラーが発生しました")`

## 検証

1. `cd web && npx tsc --noEmit` が通ること
2. `cd web && npx vitest run` が通ること
3. `cd web && npm run build` が通ること

## 完了条件

- 上記の検証がすべて通ったら、以下を実行:
  1. `git checkout -b feature/toast-notification` でブランチ作成（既にいれば不要）
  2. 変更をコミット & プッシュ
  3. `gh pr create --title "feat: Toast通知コンポーネント追加" --body "..." --base main` でPR作成
  4. `gh pr merge --auto --squash` で自動マージ設定
  5. 完了を宣言: <promise>DONE</promise>

全検証が通り、PRが作成され自動マージが設定されるまで <promise>DONE</promise> を出力しないでください。
