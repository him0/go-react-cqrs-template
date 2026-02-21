このプロジェクトに **マルチステージDockerfile** を追加してください。

## 要件

- フロントエンド（Vite build）+ バックエンド（Go binary）を1つのイメージに
- 最終イメージは軽量（distroless or alpine）
- docker-compose.prod.yml で本番構成を定義

## 実装ガイド

### Dockerfile（プロジェクトルートに作成）

```dockerfile
# Stage 1: Frontend build
FROM node:24-alpine AS frontend
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# Stage 2: Backend build
FROM golang:1.24-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/server/main.go

# Stage 3: Final
FROM alpine:3.21
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=backend /server .
COPY --from=frontend /app/web/dist ./web/dist
EXPOSE 8080
CMD ["./server"]
```

### docker-compose.prod.yml

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=app_db
      - PORT=8080
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: app_db
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  pgdata:
```

### .dockerignore

```
node_modules
web/node_modules
web/dist
bin/
tmp/
.git
*.md
coverage.out
```

## 検証

1. `Dockerfile` の構文確認（`docker build --check .` 的なもの、なければ手動確認）
2. `docker-compose.prod.yml` の構文確認（`docker compose -f docker-compose.prod.yml config` あるいは手動確認）
3. `.dockerignore` が存在すること
4. Go ビルドが通ること: `go build ./...`

## 完了条件

- 上記の検証が通ったら、以下を実行:
  1. `git checkout -b feature/multistage-dockerfile` でブランチ作成（既にいれば不要）
  2. 変更をコミット & プッシュ
  3. `gh pr create --title "feat: マルチステージDockerfile追加" --body "..." --base main` でPR作成
  4. `gh pr merge --auto --squash` で自動マージ設定
  5. 完了を宣言: <promise>DONE</promise>

全検証が通り、PRが作成され自動マージが設定されるまで <promise>DONE</promise> を出力しないでください。
