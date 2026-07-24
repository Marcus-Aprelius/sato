# Release v0.0.1

## Notes

Initial release of `sato` - Secure Access Task Operator. A lightweight utility for injecting KeePass secrets into Docker Compose without `.env` files.

---

## Overview

Secure secret injection for Docker Compose workflows with KeePass integration.

---

## Features

- ✅ KeePass integration (read `.kdbx` files)
- ✅ Docker Compose command execution
- ✅ Strict command allow-list (`docker compose` only)
- ✅ DB auto-discovery (--db-path, ~/.sato/, SATO_DB_PATH)
- ✅ Interactive & piped password input
- ✅ Static binary (~3MB)
- ✅ No `.env` files
- ✅ Secrets isolated to child processes only
- ✅ Bash completion
- ✅ Recursive KeePass group reading
- ✅ Secret tree view (`sato get secrets --tree`)
- ✅ Optional empty group display (`--show-empty-groups`)
- ✅ Unit and E2E test reports

---

## Quick Usage

```bash
# Show status
sato

# Run docker compose with secrets
sato docker compose up -d

# Use specific DB
sato --db-path=/path/to/secrets.kdbx docker compose up -d
```

---

## Install

```bash
bash scripts/build/create.sh bin
```

---

## Test

```bash
cd playground
echo "sato" | ../dist/sato --db-path=secrets.kdbx docker compose config
```

---

© 2026 [Marcus-Aprelius](https://github.com/Marcus-Aprelius/sato)

[Apache License 2.0](LICENSE)