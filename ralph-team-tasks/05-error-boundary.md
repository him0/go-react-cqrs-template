このプロジェクトのフロントエンドに **Error Boundary** を追加してください。

## 要件

- React Error Boundary でアプリ全体の未キャッチエラーをハンドリング
- エラー画面を表示（リトライボタン付き）
- 開発時はエラースタックを表示

## 実装ガイド

1. `web/src/components/ErrorBoundary.tsx` を作成
   - クラスコンポーネント（Error Boundary はクラスのみ）
   - `componentDidCatch` でエラーログ出力
   - フォールバックUI: エラーメッセージ + 「再読み込み」ボタン
   - `import.meta.env.DEV` で開発時のみスタック表示

2. `web/src/components/ErrorFallback.tsx` を作成
   - フォールバック表示用コンポーネント
   - shadcn/ui の Button コンポーネントを使用
   - Tailwind でスタイリング

3. `web/src/main.tsx` または `web/src/routes/__root.tsx` で Error Boundary を適用

4. テスト `web/src/components/__tests__/ErrorBoundary.test.tsx` を作成
   - エラーを投げるコンポーネントでフォールバック表示を確認
   - console.error をモックして不要な出力を抑制

## 検証

1. `cd web && npm run typecheck` が通ること（あれば、なければ `npx tsc --noEmit`）
2. `cd web && npx vitest run` が通ること
3. `cd web && npm run build` が通ること

## 完了条件

- 上記の検証がすべて通ったら、以下を実行:
  1. `git checkout -b feature/error-boundary` でブランチ作成（既にいれば不要）
  2. 変更をコミット & プッシュ
  3. `gh pr create --title "feat: Error Boundary追加" --body "..." --base main` でPR作成
  4. `gh pr merge --auto --squash` で自動マージ設定
  5. 完了を宣言: <promise>DONE</promise>

全検証が通り、PRが作成され自動マージが設定されるまで <promise>DONE</promise> を出力しないでください。
