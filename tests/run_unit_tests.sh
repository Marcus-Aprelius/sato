#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

PROJECT_DIR="$SCRIPT_DIR"
while [ ! -f "$PROJECT_DIR/go.mod" ]; do
    parent="$(dirname "$PROJECT_DIR")"

    if [ "$parent" = "$PROJECT_DIR" ]; then
        echo "ERROR: project root with go.mod not found" >&2
        exit 1
    fi

    PROJECT_DIR="$parent"
done

REPORT_FILE="$PROJECT_DIR/tests/report_unit_tests.txt"
GO_IMAGE="golang:1.26.5-alpine"

cd "$PROJECT_DIR"

rm -f "$REPORT_FILE"

{
    echo "SATO unit test report"
    echo "Generated: $(date -u +"%Y-%m-%d %H:%M:%S UTC")"
    echo ""

    echo "=== Create playground ==="
    bash tools/playground/playground_create.sh
    echo ""

    cleanup() {
        echo ""
        echo "=== Delete playground ==="
        bash tools/playground/playground_delete.sh
    }

    trap cleanup EXIT

    warn_gofmt() {
        echo "=== Check gofmt ==="

        if command -v gofmt >/dev/null 2>&1; then
            unformatted="$(
                find . \
                    -type f \
                    -name "*.go" \
                    -not -path "./vendor/*" \
                    -not -path "./dist/*" \
                    -not -path "./bin/*" \
                    -exec gofmt -l {} +
            )"
        else
            unformatted="$(
                docker run --rm \
                    -v "$PROJECT_DIR:/src" \
                    -w /src \
                    "$GO_IMAGE" \
                    sh -c 'find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./dist/*" -not -path "./bin/*" -exec gofmt -l {} +'
            )"
        fi

        if [ -n "$unformatted" ]; then
            echo "WARNING: gofmt is required for:"
            echo "$unformatted"
            echo ""
            echo "Run:"
            echo "$unformatted" | sed 's/^/  gofmt -w /'
            echo ""
            echo "Continuing tests..."
            echo ""
            return 0
        fi

        echo "PASS: gofmt"
        echo ""
    }

    run_go_vet() {
        echo "=== Go vet ==="

        if command -v go >/dev/null 2>&1; then
            go vet ./...
        else
            docker run --rm \
                -v "$PROJECT_DIR:/src" \
                -w /src \
                "$GO_IMAGE" \
                go vet ./...
        fi

        echo "PASS: go vet"
        echo ""
    }

    warn_gofmt
    run_go_vet

    if command -v go >/dev/null 2>&1; then
        echo "[Runner] local go"
        go version
        echo ""
        go test -v -coverpkg=./internal/sato ./internal/sato ./tests/unit
    else
        echo "[Runner] docker $GO_IMAGE"
        echo ""
        docker run --rm \
            -v "$PROJECT_DIR:/src" \
            -w /src \
            -v sato-go-mod-cache:/go/pkg/mod \
            -v sato-go-build-cache:/root/.cache/go-build \
            "$GO_IMAGE" \
            go test -v -coverpkg=./internal/sato ./internal/sato ./tests/unit
    fi
} 2>&1 | tee "$REPORT_FILE"

echo ""
echo "Report written to: $REPORT_FILE"
