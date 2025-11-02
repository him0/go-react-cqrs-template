#!/bin/bash

# マイグレーションファイル生成スクリプト
# Usage: ./scripts/generate-migration.sh <migration_name>
# Example: ./scripts/generate-migration.sh add_user_status

set -e

# 引数チェック
if [ $# -eq 0 ]; then
    echo "Error: Migration name is required"
    echo "Usage: $0 <migration_name>"
    echo "Example: $0 add_user_status"
    exit 1
fi

MIGRATION_NAME=$1
TIMESTAMP=$(date +%Y%m%d%H%M%S)
MIGRATION_FILE="db/migrations/${TIMESTAMP}_${MIGRATION_NAME}.sql"

# マイグレーションディレクトリを作成
mkdir -p db/migrations

echo "Generating migration: $MIGRATION_NAME"
echo "Output file: $MIGRATION_FILE"
echo ""

# psqldefでdry-runを実行してマイグレーションSQLを生成
go tool psqldef -U postgres -p 5432 -h localhost app_db --password=postgres \
    --file=db/schema/schema.sql --dry-run > "$MIGRATION_FILE"

# ファイルが空でないかチェック
if [ -s "$MIGRATION_FILE" ]; then
    echo "✓ Migration file created successfully!"
    echo ""
    echo "--- Migration SQL ---"
    cat "$MIGRATION_FILE"
    echo ""
    echo "--- Next Steps ---"
    echo "1. Review the migration file: $MIGRATION_FILE"
    echo "2. Apply the migration: task db:migrate"
else
    echo "✗ No schema changes detected"
    rm "$MIGRATION_FILE"
    exit 1
fi
