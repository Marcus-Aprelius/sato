#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
SATO_BIN="${SATO_BIN:-$PROJECT_DIR/dist/sato}"
PLAYGROUND_DIR="$PROJECT_DIR/tools/playground"

if [ ! -x "$SATO_BIN" ]; then
    echo "SKIP: binary not found at $SATO_BIN" >&2
    exit 77
fi

if ! command -v docker >/dev/null 2>&1; then
    echo "SKIP: docker not installed" >&2
    exit 77
fi

for extension in psafe3 ibak; do
    DB_PATH="$PLAYGROUND_DIR/secrets.$extension"

    if [ ! -f "$DB_PATH" ]; then
        echo "SKIP: database not found at $DB_PATH" >&2
        exit 77
    fi

    value="$(
        echo "sato" |
            "$SATO_BIN" \
                --db-path="$DB_PATH" \
                get secret DB_PASSWORD \
                -q
    )"

    if [ "$value" != "super_secret_db_password_123" ]; then
        echo "FAIL: unexpected DB_PASSWORD from .$extension" >&2
        exit 1
    fi

    output="$(
        cd "$PLAYGROUND_DIR"

        echo "sato" |
            "$SATO_BIN" \
                --db-path="$DB_PATH" \
                docker compose config \
                2>/dev/null
    )"

    if ! grep -q 'DB_PASSWORD: super_secret_db_password_123' <<<"$output"; then
        echo "FAIL: DB_PASSWORD not injected from .$extension" >&2
        exit 1
    fi

    if ! grep -q 'API_KEY: sk-test-api-key-xyz789' <<<"$output"; then
        echo "FAIL: API_KEY not injected from .$extension" >&2
        exit 1
    fi

    echo "PASS: .$extension secrets loaded and injected"
done

echo "PASS: Password Safe integration works"
