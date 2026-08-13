#!/usr/bin/env bash
set -euo pipefail

PLAYGROUND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PLAYGROUND_DIR"

echo "[SATO playground cleanup]"
echo "Directory: $PLAYGROUND_DIR"
echo ""

echo "Removing generated playground files..."

rm -f Dockerfile
rm -f docker-compose.yml
rm -f .env
rm -f test_script.sh
rm -f run_test.sh
rm -f secrets.kdbx

echo ""
echo "Removed if existed:"
echo "  Dockerfile"
echo "  docker-compose.yml"
echo "  .env"
echo "  test_script.sh"
echo "  run_test.sh"
echo "  secrets.kdbx"
echo ""

echo "Kept:"
echo "  README.md"
echo "  playground_create.sh"
echo "  playground_delete.sh"
echo ""

echo "Generator is kept in:"
echo "  ../create_db/create_test_db.go"
echo ""

echo "Playground cleanup completed."
