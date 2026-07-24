#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPORT_FILE="$PROJECT_DIR/tests/report_all_tests.txt"

cd "$PROJECT_DIR"

rm -f "$REPORT_FILE"

{
    echo "sato all tests report"
    echo "Generated: $(date -u +"%Y-%m-%d %H:%M:%S UTC")"
    echo ""

    echo "=== Unit tests ==="
    bash tests/run_unit_tests.sh
    echo ""

    echo "=== E2E tests ==="
    bash tests/run_e2e_tests.sh
    echo ""

    echo "All test suites completed successfully"
} 2>&1 | tee "$REPORT_FILE"

echo ""
echo "Report written to: $REPORT_FILE"
