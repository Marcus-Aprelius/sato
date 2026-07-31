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

tmp_home="$(mktemp -d)"
trap 'rm -rf "$tmp_home"' EXIT

out_no_db="$(env -u SATO_DB_PATH HOME="$tmp_home" "$SATO_BIN")"

if echo "$out_no_db" | grep -q 'Mode'; then
    echo "FAIL: Mode column must be hidden when DB is not found" >&2
    echo "$out_no_db" >&2
    exit 1
fi

if ! echo "$out_no_db" | grep -q 'Database: not found'; then
    echo "FAIL: missing 'Database: not found'" >&2
    echo "$out_no_db" >&2
    exit 1
fi

if echo "$out_no_db" | grep -q 'Docker Compose:'; then
    echo "FAIL: sato status must not show Docker Compose version" >&2
    echo "$out_no_db" >&2
    exit 1
fi

out_with_db="$(env SATO_DB_PATH="$DB_PATH" HOME="$tmp_home" "$SATO_BIN")"

if ! echo "$out_with_db" | grep -q 'Mode'; then
    echo "FAIL: Mode column must be shown when DB is found" >&2
    echo "$out_with_db" >&2
    exit 1
fi

if ! echo "$out_with_db" | grep -q 'SATO_DB_PATH'; then
    echo "FAIL: SATO_DB_PATH row missing" >&2
    echo "$out_with_db" >&2
    exit 1
fi

if ! echo "$out_with_db" | grep -qE 'SATO_DB_PATH[[:space:]]+\| set[[:space:]]+\| (RO|RW)'; then
    echo "FAIL: SATO_DB_PATH row must contain RO or RW mode" >&2
    echo "$out_with_db" >&2
    exit 1
fi

if ! echo "$out_with_db" | grep -q "Database: $DB_PATH"; then
    echo "FAIL: Database path line missing" >&2
    echo "$out_with_db" >&2
    exit 1
fi

if echo "$out_with_db" | grep -qE 'Database: .*\[(RO|RW)\]'; then
    echo "FAIL: Database line must not contain [RO]/[RW]" >&2
    echo "$out_with_db" >&2
    exit 1
fi

if echo "$out_with_db" | grep -q 'Docker Compose:'; then
    echo "FAIL: sato status must not show Docker Compose version" >&2
    echo "$out_with_db" >&2
    exit 1
fi

echo "PASS: status output works"
