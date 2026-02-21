import { z } from 'zod'

export const createUserSchema = z.object({
  name: z.string().min(1, '名前は必須です').max(100, '名前は100文字以内です'),
  email: z.string().email('有効なメールアドレスを入力してください'),
})

export type CreateUserFormValues = z.infer<typeof createUserSchema>

export const updateUserSchema = z.object({
  name: z.string().min(1, '名前は必須です').max(100, '名前は100文字以内です'),
  email: z.string().email('有効なメールアドレスを入力してください'),
})

export type UpdateUserFormValues = z.infer<typeof updateUserSchema>
