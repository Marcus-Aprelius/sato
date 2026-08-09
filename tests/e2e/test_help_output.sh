#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
SATO_BIN="${SATO_BIN:-$PROJECT_DIR/dist/sato}"

if [ ! -x "$SATO_BIN" ]; then
    echo "SKIP: binary not found at $SATO_BIN (run tools/build/create.sh bin first)" >&2
    exit 77
fi

out="$("$SATO_BIN" help)"

if ! echo "$out" | grep -q 'git help'; then
    echo "FAIL: main help missing git help" >&2
    echo "$out" >&2
    exit 1
fi

if ! echo "$out" | grep -q 'docker help'; then
    echo "FAIL: main help missing docker help" >&2
    echo "$out" >&2
    exit 1
fi

git_out="$("$SATO_BIN" git)"

if ! echo "$git_out" | grep -q 'git clone <URL>'; then
    echo "FAIL: sato git should show git help" >&2
    echo "$git_out" >&2
    exit 1
fi

for cmd in 'git fetch' 'git pull' 'git push'; do
    if ! echo "$git_out" | grep -q "$cmd"; then
        echo "FAIL: git help missing $cmd" >&2
        echo "$git_out" >&2
        exit 1
    fi
done

docker_out="$("$SATO_BIN" docker)"

if ! echo "$docker_out" | grep -q 'docker compose'; then
    echo "FAIL: sato docker should show docker help" >&2
    echo "$docker_out" >&2
    exit 1
fi

echo "PASS: help output works"
