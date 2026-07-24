#!/bin/bash
# E2E test: version output contains version string and commit SHA in [brackets]

set -e
SATO_BIN="${SATO_BIN:-$(cd "$(dirname "$0")/../.." && pwd)/dist/sato}"

if [ ! -x "$SATO_BIN" ]; then
    echo "SKIP: binary not found at $SATO_BIN (run scripts/build/create.sh first)" >&2
    exit 77
fi

out=$("$SATO_BIN" version)

# Expect two lines: "vX.Y.Z" and "[<commit>]"
if ! echo "$out" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$'; then
    echo "FAIL: version line missing or malformed" >&2
    echo "$out" >&2
    exit 1
fi

if ! echo "$out" | grep -qE '^\[[a-z0-9]+\]$'; then
    echo "FAIL: commit line missing or malformed" >&2
    echo "$out" >&2
    exit 1
fi

echo "PASS: version output correct"
