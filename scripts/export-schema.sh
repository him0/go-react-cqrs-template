#!/bin/bash

# データベーススキーマエクスポートスクリプト
# Usage: ./scripts/export-schema.sh [output_file]
# Example: ./scripts/export-schema.sh db/migrations/001_baseline.sql

set -e

OUTPUT_FILE=${1:-"db/migrations/$(date +%Y%m%d%H%M%S)_baseline.sql"}

# マイグレーションディレクトリを作成
mkdir -p db/migrations

echo "Exporting current database schema..."
echo "Output file: $OUTPUT_FILE"
echo ""

# psqldefでスキーマをエクスポート
go tool psqldef -U postgres -p 5432 -h localhost app_db --password=postgres \
    --export > "$OUTPUT_FILE"

# ファイルが作成されたかチェック
if [ -s "$OUTPUT_FILE" ]; then
    echo "✓ Schema exported successfully!"
    echo ""
    echo "File size: $(wc -l < "$OUTPUT_FILE") lines"
    echo ""
    echo "--- Preview (first 20 lines) ---"
    head -n 20 "$OUTPUT_FILE"
    echo ""
    echo "--- Next Steps ---"
    echo "This file can be used as a baseline migration or reference."
else
    echo "✗ Export failed or database is empty"
    rm "$OUTPUT_FILE"
    exit 1
fi
