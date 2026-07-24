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

out="$(echo "sato" | "$SATO_BIN" --db-path="$DB_PATH" get secrets --tree 2>/dev/null)"

if echo "$out" | grep -q '^root$'; then
    echo "FAIL: tree output must not contain root" >&2
    echo "$out" >&2
    exit 1
fi

if echo "$out" | grep -q 'Secrets/'; then
    echo "FAIL: common top-level KeePass group should be stripped from tree output" >&2
    echo "$out" >&2
    exit 1
fi

for secret in API_KEY DB_PASSWORD; do
    if ! echo "$out" | grep -q "$secret"; then
        echo "FAIL: tree output missing $secret" >&2
        echo "$out" >&2
        exit 1
    fi
done

if echo "$out" | grep -q 'empty-group/'; then
    echo "FAIL: empty groups must be hidden by default" >&2
    echo "$out" >&2
    exit 1
fi

out_empty="$(echo "sato" | "$SATO_BIN" --db-path="$DB_PATH" get secrets --tree --show-empty-groups 2>/dev/null)"

if echo "$out_empty" | grep -q 'empty-group'; then
    if ! echo "$out_empty" | grep -q 'empty-group/'; then
        echo "FAIL: empty group must be displayed with trailing slash" >&2
        echo "$out_empty" >&2
        exit 1
    fi
fi

echo "PASS: get secrets tree output works"
