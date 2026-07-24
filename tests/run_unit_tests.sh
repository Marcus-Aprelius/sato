#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPORT_FILE="$PROJECT_DIR/tests/report_unit_tests.txt"

cd "$PROJECT_DIR"

rm -f "$REPORT_FILE"

{
    echo "SATO unit test report"
    echo "Generated: $(date -u +"%Y-%m-%d %H:%M:%S UTC")"
    echo ""

    echo "=== Create playground ==="
    bash playground/playground_create.sh
    echo ""

    cleanup() {
        echo ""
        echo "=== Delete playground ==="
        bash playground/playground_delete.sh
    }

    trap cleanup EXIT

    if command -v go >/dev/null 2>&1; then
        echo "[Runner] local go"
        go version
        echo ""
        go test -v -coverpkg=./internal/sato ./tests/unit
    else
        echo "[Runner] docker golang:1.26.5-alpine"
        echo ""
        docker run --rm \
            -v "$PROJECT_DIR:/src" \
            -w /src \
            golang:1.26.5-alpine \
            go test -v -coverpkg=./internal/sato ./tests/unit
    fi
} 2>&1 | tee "$REPORT_FILE"

echo ""
echo "Report written to: $REPORT_FILE"
