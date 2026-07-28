# Release v0.0.2

## Overview

### SATO - Secure Access Task Operator

Secure secret injection for Docker Compose workflows with KeePass integration.

---

## Changes:
Improved:
  - build of all packages (`bin`/`.deb`/`.rpm`) at once: `create.sh all`;
  - checking of Go syntax (`gofmt`,`go vet`). If files need formatting - warning will appear, otherwise - just information.
  - new tests added and code refactored
  - README.md files updated
  - Database "Mode" displaying - a new column "Mode" added in the table instead of the status at the right of DB name:
    If no DB is found, the `Mode` column is hidden to keep the output simple.
  ```text
    [Now]                                                        [Before]
    ...                                                          ...
    Priority | Path           | Status    | Mode                 Priority | Path           | Status
    ---------|----------------|-----------|------                ---------|----------------|-----------
     1       | --db-path      | not set   |                       1       | --db-path      | not set
     2       | ~/.sato/*.kdbx | not found |                       2       | ~/.sato/*.kdbx | not found
    [3]      | SATO_DB_PATH   | set       | RW                   [3]      | SATO_DB_PATH   | set

    Database: /home/user/.sato/secrets.kdbx                      Database: /home/user/.sato/secrets.kdbx [RW]
    ...                                                          ...
  ```
  - Version of docker compose moved from `sato` to sato `version`

Fixed:
  - download links fixed;
  - `.rpm` package build fixed - deletion of temporary folder fixed;
  - some other small issues;

---

## Test results:

| Area          | Result | Details |
|---------------|--------|---------|
| Preflight     | PASS   | gofmt, go mod tidy, go vet, go test ./... |
| Unit tests    | PASS   | 19 passed, coverage 19.0% |
| E2E tests     | PASS   | 9 passed, 0 failed, 0 skipped |
| Binary        | PASS   | dist/sato, size 3.0M, commit 2b9414d   |
| Final status  | PASS   | All test suites completed successfully |

---

## Quick Usage

* Show status
  ```bash
  sato
  ```

* Run docker compose with secrets
  ```bash
  sato docker compose up -d
  ```

* Use specific DB
  ```bash
  sato --db-path=/path/to/secrets.kdbx docker compose up -d
  ```

---

## Installation

### 1. Prerequisites
* required:
  * OS Linux
  * [docker compose](https://docs.docker.com/compose/install/linux) - to have ability to use `sato`
* optional:
  * [git](https://git-scm.com/install/linux)           - for cloning repository
  * [docker](https://docs.docker.com/engine/install/)  - for manual binary building
  * [keepassxc](https://keepassxc.org/download/#linux) - to edit secrets via UI

### 2. Obtain files from Github Releases (or [build  manually](scripts/build/README.md)):

* `bin`ary file:
  ```bash
  wget https://github.com/Marcus-Aprelius/sato/releases/latest/download/sato
  ```

* `.deb` package:
  ```bash
  wget https://github.com/Marcus-Aprelius/sato/releases/download/v0.0.2/sato_0.0.2_amd64.deb
  ```

* `.rpm` package:
  ```bash
  wget https://github.com/Marcus-Aprelius/sato/releases/download/v0.0.2/sato-0.0.2-1.x86_64.rpm
  ```

### 3. Install `sato`:
  * from `bin`ary file:
    ```bash
    chmod +x sato && sudo mv sato /usr/local/bin/
    ```

  * from `.deb` package:
    ```bash
    sudo dpkg -i sato_0.0.2_amd64.deb
    ```

  * from `.rpm` package:
    ```bash
    sudo yum install -y sato-0.0.2-1.x86_64.rpm
    ```


© 2026 [Marcus-Aprelius](https://github.com/Marcus-Aprelius/sato)

[Apache License 2.0](LICENSE)
