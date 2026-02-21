#!/bin/bash
set -euo pipefail

# Only run in remote environments (Claude Code on the web)
if [ "${CLAUDE_CODE_REMOTE:-}" != "true" ]; then
  exit 0
fi

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(cd "$(dirname "$0")/../.." && pwd)}"

##############################################################################
# 1. Install Go tools
##############################################################################

# go-task (task runner)
if ! command -v task &>/dev/null; then
  echo "Installing go-task..."
  curl -sL https://github.com/go-task/task/releases/download/v3.40.1/task_linux_amd64.tar.gz \
    | tar xz -C /usr/local/bin task
fi

# sqlc (SQL code generator)
SQLC_VERSION="1.30.0"
if ! command -v sqlc &>/dev/null || [ "$(sqlc version)" != "v${SQLC_VERSION}" ]; then
  echo "Installing sqlc v${SQLC_VERSION}..."
  curl -sL "https://github.com/sqlc-dev/sqlc/releases/download/v${SQLC_VERSION}/sqlc_${SQLC_VERSION}_linux_amd64.tar.gz" \
    | tar xz -C /usr/local/bin sqlc
fi

# goimports
if ! command -v goimports &>/dev/null; then
  echo "Installing goimports..."
  go install golang.org/x/tools/cmd/goimports@latest
fi

# oapi-codegen
if ! command -v oapi-codegen &>/dev/null; then
  echo "Installing oapi-codegen..."
  GOBIN=/usr/local/bin go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
fi

##############################################################################
# 2. Download Go module dependencies
##############################################################################

echo "Downloading Go modules..."
cd "$PROJECT_DIR"
go mod download

##############################################################################
# 3. Install frontend dependencies
##############################################################################

echo "Installing frontend dependencies..."
cd "$PROJECT_DIR"
pnpm install

##############################################################################
# 4. Start PostgreSQL and set up the database
##############################################################################

echo "Starting PostgreSQL..."
pg_ctlcluster 16 main start 2>/dev/null || true

# Wait for PostgreSQL to be ready
for i in $(seq 1 10); do
  if pg_isready -q 2>/dev/null; then
    break
  fi
  sleep 1
done

# Create database and user if they don't exist
sudo -u postgres psql -tc "SELECT 1 FROM pg_database WHERE datname='app_db'" \
  | grep -q 1 || sudo -u postgres psql -c "CREATE DATABASE app_db;"
sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'postgres';" 2>/dev/null

# Apply schema
cat "$PROJECT_DIR"/db/schema/*.sql \
  | PGPASSWORD=postgres psql -U postgres -h localhost -p 5432 -d app_db 2>/dev/null || true

##############################################################################
# 5. Set environment variables for the session
##############################################################################

if [ -n "${CLAUDE_ENV_FILE:-}" ]; then
  echo 'export DB_PORT=5432' >> "$CLAUDE_ENV_FILE"
fi

echo "Session start hook completed successfully."
