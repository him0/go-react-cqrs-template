このプロジェクトのフロントエンドに **react-hook-form + zod によるフォーム管理** を導入してください。

## 要件

- react-hook-form + @hookform/resolvers + zod を使用
- 既存の UserCreateForm と UserEditForm を react-hook-form に移行
- zodスキーマでバリデーション
- フィールドレベルのエラー表示

## 実装ガイド

1. `cd web && npm install react-hook-form @hookform/resolvers`（zod は既にあるか確認）

2. `web/src/lib/validations/user.ts` にzodスキーマを作成
   ```typescript
   import { z } from "zod";

   export const createUserSchema = z.object({
     name: z.string().min(1, "名前は必須です").max(100, "名前は100文字以内です"),
     email: z.string().email("有効なメールアドレスを入力してください"),
   });

   export const updateUserSchema = z.object({
     name: z.string().min(1, "名前は必須です").max(100, "名前は100文字以内です"),
     email: z.string().email("有効なメールアドレスを入力してください"),
   });
   ```

3. `web/src/components/UserCreateForm.tsx` を react-hook-form に移行
   - `useForm` + `zodResolver` を使用
   - `register` でフィールドを接続
   - フィールドエラーメッセージの表示
   - `handleSubmit` でサブミット処理

4. `web/src/components/UserEditForm.tsx` も同様に移行
   - 既存値を `defaultValues` にセット

5. テストの更新
   - 既存のフォームテストを react-hook-form に対応させる

## 検証

1. `cd web && npx tsc --noEmit` が通ること
2. `cd web && npx vitest run` が通ること
3. `cd web && npm run build` が通ること

## 完了条件

- 上記の検証がすべて通ったら、以下を実行:
  1. `git checkout -b feature/form-library` でブランチ作成（既にいれば不要）
  2. 変更をコミット & プッシュ
  3. `gh pr create --title "feat: react-hook-form + zod によるフォーム管理導入" --body "..." --base main` でPR作成
  4. `gh pr merge --auto --squash` で自動マージ設定
  5. 完了を宣言: <promise>DONE</promise>

全検証が通り、PRが作成され自動マージが設定されるまで <promise>DONE</promise> を出力しないでください。
