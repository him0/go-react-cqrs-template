このプロジェクトのフロントエンドに **ダークモード対応** を追加してください。

## 要件

- システム設定に応じた自動切り替え
- 手動でライト/ダーク/システムを切り替えるトグルボタン
- 選択を localStorage に保存
- shadcn/ui + Tailwind の dark mode クラスを活用

## 実装ガイド

1. `web/src/lib/theme.ts` にテーマ管理ユーティリティを作成
   - `getTheme()`, `setTheme()`, `toggleTheme()`
   - localStorage キー: `"theme"`
   - 値: `"light"` | `"dark"` | `"system"`
   - `<html>` の `class` に `dark` を付与/除去

2. `web/src/components/ThemeToggle.tsx` を作成
   - ライト/ダーク/システムの3つを切り替え
   - Sun/Moon アイコン（lucide-react を使用、既にshadcn/uiの依存にある可能性）
   - shadcn/ui の Button コンポーネントを使用

3. `web/src/routes/__root.tsx` にテーマ初期化処理を追加
   - ページロード時に localStorage からテーマを復元
   - ThemeToggle をヘッダー/ナビゲーションに配置

4. `web/index.html` の `<script>` でフラッシュ防止:
   ```html
   <script>
     (function(){var t=localStorage.getItem('theme');if(t==='dark'||(t!=='light'&&matchMedia('(prefers-color-scheme:dark)').matches))document.documentElement.classList.add('dark')})()
   </script>
   ```

5. Tailwind 設定確認: `web/tailwind.config.js` に `darkMode: 'class'` があるか確認、なければ追加

## 検証

1. `cd web && npx tsc --noEmit` が通ること
2. `cd web && npx vitest run` が通ること
3. `cd web && npm run build` が通ること

## 完了条件

- 上記の検証がすべて通ったら、以下を実行:
  1. `git checkout -b feature/dark-mode` でブランチ作成（既にいれば不要）
  2. 変更をコミット & プッシュ
  3. `gh pr create --title "feat: ダークモード対応" --body "..." --base main` でPR作成
  4. `gh pr merge --auto --squash` で自動マージ設定
  5. 完了を宣言: <promise>DONE</promise>

全検証が通り、PRが作成され自動マージが設定されるまで <promise>DONE</promise> を出力しないでください。
