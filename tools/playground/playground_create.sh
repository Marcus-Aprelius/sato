#!/usr/bin/env bash
set -euo pipefail

PLAYGROUND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$PLAYGROUND_DIR/../.." && pwd)"
DB_GENERATOR="$PROJECT_DIR/tools/create_db/create_test_db.go"

cd "$PLAYGROUND_DIR"

DB_PASSWORD="sato"

echo "[SATO playground]"
echo "Directory: $PLAYGROUND_DIR"
echo ""

echo "[1/3] Recreating .env..."

cat > .env <<'EOF'
CUSTOM_VAR=from_env_file
ANOTHER_VAR=should_not_appear
EOF

echo "[2/3] Recreating docker-compose.yml..."

cat > docker-compose.yml <<'EOF'
services:
  secret-test:
    image: alpine:latest
    container_name: sato-test
    environment:
      DB_PASSWORD: ${DB_PASSWORD}
      API_KEY: ${API_KEY}
    command:
      - sh
      - -c
      - |
        echo "=== Secrets loaded from KeePass ==="
        echo "DB_PASSWORD=$DB_PASSWORD"
        echo "API_KEY=$API_KEY"
        echo "=== End of output ==="
EOF

echo "[3/3] Recreating secret databases..."
rm -f secrets.kdbx secrets.psafe3 secrets.ibak

if command -v go >/dev/null 2>&1; then
    go run "$DB_GENERATOR"
else
    if ! command -v docker >/dev/null 2>&1; then
        echo "ERROR: neither go nor docker is available"
        exit 1
    fi

    docker run --rm -v "$PROJECT_DIR:/src" -w /src/tools/playground golang:1.26.7-alpine go run ../create_db/create_test_db.go
fi

echo ""
echo "Playground recreated successfully:"
echo "  .env"
echo "  docker-compose.yml"
echo "  secrets.kdbx"
echo "  secrets.psafe3"
echo "  secrets.ibak"
echo ""

echo "Secret databases:"
echo "  $PLAYGROUND_DIR/secrets.kdbx"
echo "  $PLAYGROUND_DIR/secrets.psafe3"
echo "  $PLAYGROUND_DIR/secrets.ibak"
echo "  Password: $DB_PASSWORD"
echo ""

echo "KeePass structure:"
echo "  Secrets/"
echo "  ├── DB_PASSWORD"
echo "  ├── API_KEY"
echo "  ├── NGINX_PASSWORD"
echo "  ├── test/"
echo "  │   ├── test1"
echo "  │   ├── test2"
echo "  │   └── empty-nested/"
echo "  └── empty-group/"
echo ""

echo ".env file:"
echo "  CUSTOM_VAR=from_env_file"
echo "  ANOTHER_VAR=should_not_appear"
echo ""

echo "Docker Compose demo uses only:"
echo "  DB_PASSWORD"
echo "  API_KEY"
echo ""
