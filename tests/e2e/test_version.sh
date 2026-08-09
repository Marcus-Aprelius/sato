#!/bin/bash
# E2E test: version output contains SATO version, commit SHA,
# and environment versions when available.

set -e

SATO_BIN="${SATO_BIN:-$(cd "$(dirname "$0")/../.." && pwd)/dist/sato}"

if [ ! -x "$SATO_BIN" ]; then
    echo "SKIP: binary not found at $SATO_BIN (run tools/build/create.sh first)" >&2
    exit 77
fi

out=$("$SATO_BIN" version)

if ! echo "$out" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$'; then
    echo "FAIL: version line missing or malformed" >&2
    echo "$out" >&2
    exit 1
fi

if ! echo "$out" | grep -qE '^\[SHA: [a-z0-9]+\]$'; then
    echo "FAIL: commit line missing or malformed" >&2
    echo "$out" >&2
    exit 1
fi

if command -v git >/dev/null 2>&1; then
    if ! echo "$out" | grep -qE '^Git: .+'; then
        echo "FAIL: Git version line missing" >&2
        echo "$out" >&2
        exit 1
    fi
fi

if command -v docker >/dev/null 2>&1; then
    if ! echo "$out" | grep -qE '^Docker: .+'; then
        echo "FAIL: Docker version line missing" >&2
        echo "$out" >&2
        exit 1
    fi

    if docker compose version --short >/dev/null 2>&1; then
        if ! echo "$out" | grep -qE '^Docker Compose: .+'; then
            echo "FAIL: Docker Compose version line missing" >&2
            echo "$out" >&2
            exit 1
        fi
    fi
fi

if [ -f /etc/os-release ]; then
    if ! echo "$out" | grep -qE '^OS: .+'; then
        echo "FAIL: OS version line missing" >&2
        echo "$out" >&2
        exit 1
    fi
fi

echo "PASS: version output correct"
