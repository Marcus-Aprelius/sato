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

expect_fail() {
    local name="$1"
    shift

    stdout_file="$(mktemp)"
    stderr_file="$(mktemp)"

    if "$@" >"$stdout_file" 2>"$stderr_file"; then
        echo "FAIL: $name should have failed" >&2
        cat "$stdout_file" >&2
        cat "$stderr_file" >&2
        rm -f "$stdout_file" "$stderr_file"
        exit 1
    fi

    rm -f "$stdout_file" "$stderr_file"
}

expect_fail_stdin() {
    local name="$1"
    local expected="$2"
    shift 2

    stdout_file="$(mktemp)"
    stderr_file="$(mktemp)"

    if echo "sato" | "$@" >"$stdout_file" 2>"$stderr_file"; then
        echo "FAIL: $name should have failed" >&2
        cat "$stdout_file" >&2
        cat "$stderr_file" >&2
        rm -f "$stdout_file" "$stderr_file"
        exit 1
    fi

    if ! grep -q -- "$expected" "$stderr_file"; then
        echo "FAIL: $name did not contain expected error: $expected" >&2
        cat "$stderr_file" >&2
        rm -f "$stdout_file" "$stderr_file"
        exit 1
    fi

    rm -f "$stdout_file" "$stderr_file"
}

expect_success_stdin() {
    local name="$1"
    shift

    stdout_file="$(mktemp)"
    stderr_file="$(mktemp)"

    if ! echo "sato" | "$@" >"$stdout_file" 2>"$stderr_file"; then
        echo "FAIL: $name should have succeeded" >&2
        cat "$stdout_file" >&2
        cat "$stderr_file" >&2
        rm -f "$stdout_file" "$stderr_file"
        exit 1
    fi
}

expect_success_stdin \
    "show empty groups in list mode" \
    "$SATO_BIN" --db-path="$DB_PATH" get secrets --show-empty-groups

if ! grep -q '^empty-group/$' "$stdout_file"; then
    echo "FAIL: list mode should show empty-group/" >&2
    cat "$stdout_file" >&2
    rm -f "$stdout_file" "$stderr_file"
    exit 1
fi

if ! grep -q '^test/empty-nested/$' "$stdout_file"; then
    echo "FAIL: list mode should show test/empty-nested/" >&2
    cat "$stdout_file" >&2
    rm -f "$stdout_file" "$stderr_file"
    exit 1
fi

rm -f "$stdout_file" "$stderr_file"

expect_fail_stdin \
    "unknown get secrets option" \
    "unknown option" \
    "$SATO_BIN" --db-path="$DB_PATH" get secrets --unknown-option

tmp_home="$(mktemp -d)"
trap 'rm -rf "$tmp_home"' EXIT

expect_fail_stdin \
    "missing DB" \
    "DB not found" \
    env -u SATO_DB_PATH HOME="$tmp_home" "$SATO_BIN" --db-path="/tmp/sato-missing-db.kdbx" get secrets

expect_fail \
    "docker compose without subcommand" \
    "$SATO_BIN" docker compose

expect_fail \
    "unsupported command" \
    "$SATO_BIN" kubectl get pods

echo "PASS: CLI error paths work"

