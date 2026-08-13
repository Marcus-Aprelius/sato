#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
SATO_BIN="${SATO_BIN:-$PROJECT_DIR/dist/sato}"
DB_PATH="$PROJECT_DIR/tools/playground/secrets.kdbx"

if [ ! -x "$SATO_BIN" ]; then
    echo "SKIP: binary not found at $SATO_BIN (run tools/build/create.sh bin first)" >&2
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

cd "$PROJECT_DIR/tools/playground"

out="$(
    DB_PASSWORD="wrong-parent-db-password" \
    API_KEY="wrong-parent-api-key" \
    LEAK_ME="must-not-leak" \
    sh -c 'echo "sato" | "$1" --db-path="$2" docker compose config' sh "$SATO_BIN" "$DB_PATH" 2>/dev/null
)"

if echo "$out" | grep -q 'wrong-parent-db-password'; then
    echo "FAIL: parent DB_PASSWORD leaked into docker compose output" >&2
    echo "$out" >&2
    exit 1
fi

if echo "$out" | grep -q 'wrong-parent-api-key'; then
    echo "FAIL: parent API_KEY leaked into docker compose output" >&2
    echo "$out" >&2
    exit 1
fi

if echo "$out" | grep -q 'LEAK_ME'; then
    echo "FAIL: unrelated parent env variable leaked into docker compose output" >&2
    echo "$out" >&2
    exit 1
fi

if echo "$out" | grep -q 'from_env_file'; then
    echo "FAIL: .env CUSTOM_VAR leaked into docker compose output" >&2
    echo "$out" >&2
    exit 1
fi

if echo "$out" | grep -q 'should_not_appear'; then
    echo "FAIL: .env ANOTHER_VAR leaked into docker compose output" >&2
    echo "$out" >&2
    exit 1
fi

if ! echo "$out" | grep -q 'DB_PASSWORD: super_secret_db_password_123'; then
    echo "FAIL: KeePass DB_PASSWORD was not injected" >&2
    echo "$out" >&2
    exit 1
fi

if ! echo "$out" | grep -q 'API_KEY: sk-test-api-key-xyz789'; then
    echo "FAIL: KeePass API_KEY was not injected" >&2
    echo "$out" >&2
    exit 1
fi

echo "PASS: parent and .env isolation works"
