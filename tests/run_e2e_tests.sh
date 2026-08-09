#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPORT_FILE="$PROJECT_DIR/tests/report_e2e_tests.txt"

cd "$PROJECT_DIR"

rm -f "$REPORT_FILE"

{
    echo "sato E2E test report"
    echo "Generated: $(date -u +"%Y-%m-%d %H:%M:%S UTC")"
    echo ""

    echo "=== Build sato binary ==="
    bash tools/build/create.sh bin
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

    pass=0
    fail=0
    skip=0

    for t in tests/e2e/test_*.sh; do
        name="$(basename "$t")"
        echo "=== $name ==="

        if bash "$t"; then
            pass=$((pass + 1))
        else
            rc=$?
            if [ "$rc" -eq 77 ]; then
                skip=$((skip + 1))
            else
                fail=$((fail + 1))
            fi
        fi

        echo ""
    done

    echo "Results: $pass passed, $fail failed, $skip skipped"

    if [ "$fail" -ne 0 ]; then
        exit 1
    fi
} 2>&1 | tee "$REPORT_FILE"

echo ""
echo "Report written to: $REPORT_FILE"
