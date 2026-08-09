#!/bin/bash
# E2E test: non-docker-compose commands must be rejected.
# "docker compose" without subcommand must also be rejected.

set -e

SATO_BIN="${SATO_BIN:-$(cd "$(dirname "$0")/../.." && pwd)/dist/sato}"

if [ ! -x "$SATO_BIN" ]; then
    echo "SKIP: binary not found at $SATO_BIN (run tools/build/create.sh bin first)" >&2
    exit 77
fi

if "$SATO_BIN" kubectl get pods >/dev/null 2>&1; then
    echo "FAIL: 'sato kubectl' should have been rejected" >&2
    exit 1
fi

if "$SATO_BIN" sh -c 'echo hi' >/dev/null 2>&1; then
    echo "FAIL: 'sato sh -c ...' should have been rejected" >&2
    exit 1
fi

if "$SATO_BIN" docker-compose up >/dev/null 2>&1; then
    echo "FAIL: legacy 'docker-compose' should have been rejected" >&2
    exit 1
fi

if "$SATO_BIN" docker run alpine echo hi >/dev/null 2>&1; then
    echo "FAIL: 'sato docker run ...' should have been rejected" >&2
    exit 1
fi

if "$SATO_BIN" docker compose >/dev/null 2>&1; then
    echo "FAIL: 'sato docker compose' without subcommand should have been rejected" >&2
    exit 1
fi

if "$SATO_BIN" git status >/dev/null 2>&1; then
    echo "FAIL: 'sato git s*atus' should have been rejected" >&2
    exit 1
fi

if "$SATO_BIN" gi* config --global user.name test >/dev/null 2>&1; then
    echo "FAIL: 'sato git config' should have been rejected" >&2
    exit 1
fi

echo "PASS: command filter works"
