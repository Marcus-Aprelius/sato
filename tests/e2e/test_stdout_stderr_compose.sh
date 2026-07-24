#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
SATO_BIN="${SATO_BIN:-$PROJECT_DIR/dist/sato}"
DB_PATH="$PROJECT_DIR/playground/secrets.kdbx"

if [ ! -x "$SATO_BIN" ]; then
    echo "SKIP: binary not found at $SATO_BIN (run scripts/build/create.sh bin first)" >&2
    exit 77
fi

if [ ! -f "$DB_PATH" ]; then
    echo "SKIP: DB not found at $DB_PATH" >&2
    exit 77
fi

if ! command -v docker >/dev/null 2>&1; then
    echo "SKIP: docker not installed" >&2
    exit 77
fi

stdout_file="$(mktemp)"
stderr_file="$(mktemp)"
trap 'rm -f "$stdout_file" "$stderr_file"' EXIT

cd "$PROJECT_DIR/playground"

echo "sato" | "$SATO_BIN" --db-path="$DB_PATH" docker compose config >"$stdout_file" 2>"$stderr_file"

if grep -q '\[DB:' "$stdout_file"; then
    echo "FAIL: docker compose stdout must not contain SATO DB marker" >&2
    cat "$stdout_file" >&2
    exit 1
fi

if grep -q 'KeePass password' "$stdout_file"; then
    echo "FAIL: docker compose stdout must not contain password prompt" >&2
    cat "$stdout_file" >&2
    exit 1
fi

if ! grep -q '\[DB:' "$stderr_file"; then
    echo "FAIL: stderr must contain SATO DB marker" >&2
    cat "$stderr_file" >&2
    exit 1
fi

if ! grep -q 'KeePass password:' "$stderr_file"; then
    echo "FAIL: stderr must contain password prompt" >&2
    cat "$stderr_file" >&2
    exit 1
fi

if ! grep -q 'DB_PASSWORD: super_secret_db_password_123' "$stdout_file"; then
    echo "FAIL: docker compose config missing DB_PASSWORD from KeePass" >&2
    cat "$stdout_file" >&2
    exit 1
fi

if ! grep -q 'API_KEY: sk-test-api-key-xyz789' "$stdout_file"; then
    echo "FAIL: docker compose config missing API_KEY from KeePass" >&2
    cat "$stdout_file" >&2
    exit 1
fi

echo "PASS: docker compose stdout/stderr separation works"
