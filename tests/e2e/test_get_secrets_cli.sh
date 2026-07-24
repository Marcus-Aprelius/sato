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
    echo "SKIP: database not found at $DB_PATH" >&2
    exit 77
fi

stdout_file="$(mktemp)"
stderr_file="$(mktemp)"
trap 'rm -f "$stdout_file" "$stderr_file"' EXIT

echo "sato" | "$SATO_BIN" --db-path="$DB_PATH" get secrets >"$stdout_file" 2>"$stderr_file"

if grep -q '\[DB:' "$stdout_file"; then
    echo "FAIL: stdout must not contain DB marker" >&2
    cat "$stdout_file" >&2
    exit 1
fi

if grep -q 'KeePass password' "$stdout_file"; then
    echo "FAIL: stdout must not contain password prompt" >&2
    cat "$stdout_file" >&2
    exit 1
fi

if ! grep -q '\[DB:' "$stderr_file"; then
    echo "FAIL: stderr must contain DB marker" >&2
    cat "$stderr_file" >&2
    exit 1
fi

if ! grep -q 'KeePass password:' "$stderr_file"; then
    echo "FAIL: stderr must contain password prompt" >&2
    cat "$stderr_file" >&2
    exit 1
fi

for secret in DB_PASSWORD API_KEY; do
    if ! grep -qx "$secret" "$stdout_file"; then
        echo "FAIL: missing secret name $secret in stdout" >&2
        cat "$stdout_file" >&2
        exit 1
    fi
done

echo "PASS: get secrets stdout/stderr separation works"
