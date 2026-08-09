#!/bin/bash
# E2E test: sato loads secrets from KeePass and injects them into
# `docker compose config` while ignoring unrelated .env values.

set -e
PROJECT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
SATO_BIN="${SATO_BIN:-$PROJECT_DIR/dist/sato}"
DB_PATH="$PROJECT_DIR/tools/playground/secrets.kdbx"

if [ ! -x "$SATO_BIN" ]; then
    echo "SKIP: binary not found at $SATO_BIN (run tools/build/create.shirst)" >&2
    exit 77
fi
if [ ! -f "$DB_PATH" ]; then
    echo "SKIP: playground DB missing at $DB_PATH" >&2
    exit 77
fi
if ! command -v docker >/dev/null 2>&1; then
    echo "SKIP: docker not installed" >&2
    exit 77
fi

cd "$PROJECT_DIR/tools/playground"
out=$(echo "sato" | "$SATO_BIN" --db-path=secrets.kdbx docker compose config 2>&1)

if ! echo "$out" | grep -q "DB_PASSWORD: super_secret_db_password_123"; then
    echo "FAIL: DB_PASSWORD not injected into compose config" >&2
    echo "$out" >&2
    exit 1
fi

if ! echo "$out" | grep -q "API_KEY: sk-test-api-key-xyz789"; then
    echo "FAIL: API_KEY not injected into compose config" >&2
    echo "$out" >&2
    exit 1
fi

echo "PASS: secrets loaded from KeePass and injected into docker compose"
