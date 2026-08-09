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

REPORT_FILE="$PROJECT_DIR/tests/report_all_tests.txt"
GO_IMAGE="golang:1.26.5-alpine"

cd "$PROJECT_DIR"

rm -f "$REPORT_FILE"

run_go() {
    if command -v go >/dev/null 2>&1; then
        go "$@" 2> >(grep -v '^go: downloading ' >&2)
    else
        docker run --rm \
            -v "$PROJECT_DIR:/src" \
            -w /src \
            -v sato-go-mod-cache:/go/pkg/mod \
            -v sato-go-build-cache:/root/.cache/go-build \
            "$GO_IMAGE" \
            go "$@" 2> >(grep -v '^go: downloading ' >&2)
    fi
}

run_and_filter() {
    set +e
    "$@" 2>&1 | grep -v '^go: downloading '
    local rc=${PIPESTATUS[0]}
    set -e
    return "$rc"
}

check_gofmt() {
    echo "=== Preflight: gofmt ==="

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
        echo "FAIL: gofmt is required for:"
        echo "$unformatted"
        echo ""
        echo "Run:"
        echo "$unformatted" | sed 's/^/  gofmt -w /'
        exit 1
    fi

    echo "PASS: gofmt"
    echo ""
}

check_go_mod_tidy() {
    echo "=== Preflight: go mod tidy ==="

    before_mod="$(
        {
            sha256sum go.mod
            sha256sum go.sum
        } 2>/dev/null || true
    )"

    run_go mod tidy

    after_mod="$(
        {
            sha256sum go.mod
            sha256sum go.sum
        } 2>/dev/null || true
    )"

    if [ "$before_mod" != "$after_mod" ]; then
        echo "FAIL: go mod tidy changed go.mod or go.sum"
        echo ""
        echo "Review changes:"
        echo "  git diff go.mod go.sum"
        exit 1
    fi

    echo "PASS: go mod tidy"
    echo ""
}

print_summary_banner() {
    local unit_report="$PROJECT_DIR/tests/report_unit_tests.txt"
    local e2e_report="$PROJECT_DIR/tests/report_e2e_tests.txt"

    local unit_passed="0"
    local unit_coverage="n/a"
    local e2e_results="n/a"
    local e2e_passed="0"
    local e2e_failed="0"
    local e2e_skipped="0"
    local binary_size="n/a"
    local git_commit="unknown"

    if [ -f "$unit_report" ]; then
        unit_passed="$(
            grep -c '^--- PASS:' "$unit_report" 2>/dev/null || true
        )"

        unit_coverage="$(
            grep -Eo 'coverage: [0-9.]+% of statements in \./internal/sato' "$unit_report" 2>/dev/null \
                | tail -1 \
                | sed 's/^coverage: //; s/ of statements in \.\/internal\/sato$//' || true
        )"

        if [ -z "$unit_coverage" ]; then
            unit_coverage="n/a"
        fi
    fi

    if [ -f "$e2e_report" ]; then
        e2e_results="$(
            grep -E '^Results:' "$e2e_report" 2>/dev/null | tail -1 || true
        )"

        e2e_passed="$(
            echo "$e2e_results" | sed -n 's/^Results: \([0-9]\+\) passed, \([0-9]\+\) failed, \([0-9]\+\) skipped/\1/p'
        )"

        e2e_failed="$(
            echo "$e2e_results" | sed -n 's/^Results: \([0-9]\+\) passed, \([0-9]\+\) failed, \([0-9]\+\) skipped/\2/p'
        )"

        e2e_skipped="$(
            echo "$e2e_results" | sed -n 's/^Results: \([0-9]\+\) passed, \([0-9]\+\) failed, \([0-9]\+\) skipped/\3/p'
        )"

        if [ -z "$e2e_passed" ]; then
            e2e_passed="0"
        fi

        if [ -z "$e2e_failed" ]; then
            e2e_failed="0"
        fi

        if [ -z "$e2e_skipped" ]; then
            e2e_skipped="0"
        fi
    fi

    if [ -f "$PROJECT_DIR/dist/sato" ]; then
        binary_size="$(du -h "$PROJECT_DIR/dist/sato" | awk '{print $1}')"
    fi

    if command -v git >/dev/null 2>&1; then
        git_commit="$(
            git rev-parse --short HEAD 2>/dev/null || echo "unknown"
        )"
    fi

    echo ""
    echo "============================================================"
    echo " TEST SUMMARY"
    echo "============================================================"
    echo ""
    echo "| Area          | Result | Details |"
    echo "|---------------|--------|---------|"
    echo "| Preflight     | PASS   | gofmt, go mod tidy, go vet, go test ./... |"
    echo "| Unit tests    | PASS   | ${unit_passed} passed, coverage ${unit_coverage} |"
    echo "| E2E tests     | PASS   | ${e2e_passed} passed, ${e2e_failed} failed, ${e2e_skipped} skipped |"
    echo "| Binary        | PASS   | dist/sato, size ${binary_size}, commit ${git_commit} |"
    echo "| Final status  | PASS   | All test suites completed successfully |"
    echo ""
    echo "============================================================"
}

{
    echo "sato all tests report"
    echo "Generated: $(date -u +"%Y-%m-%d %H:%M:%S UTC")"
    echo ""

    check_gofmt

    check_go_mod_tidy

    echo "=== Preflight: go vet ==="
    run_go vet ./...
    echo "PASS: go vet"
    echo ""

    echo "=== Preflight: go test ./... ==="
    run_go test ./...
    echo "PASS: go test ./..."
    echo ""

    echo "=== Unit tests ==="
    run_and_filter bash tests/run_unit_tests.sh
    echo ""

    echo "=== E2E tests ==="
    run_and_filter bash tests/run_e2e_tests.sh
    echo ""

    echo "All test suites completed successfully"

    print_summary_banner
} 2>&1 | tee "$REPORT_FILE"

echo ""
echo "Report written to: $REPORT_FILE"
