```text
   _____      _______ ____
  / ____|  /\|__   __| __ \    Secure
 | (___   /  \  | | | |  | |    Access
  \___ \ / /\ \ | | | |  | |     Task
  ____) | ____ \| | | |__| |      Operator
 |_____/_/    \_\_|  \____/    > like sudo, but for secrets
```
# Overview
This release makes some small improvements. See below for more details.

---

# Changes:
**Improved:**
  - added `.devcontainers` folder for easier development
  - added [`SATO Secrets`](https://marketplace.visualstudio.com/items?itemName=MarcusApreliusAntoninus.sato-vscode-ext) VS Code extention to `.devcontainer.json` to view/edit `*.kdbx` files
  - added command `sato get secret <NAME>` with a key `-q|quite`
  - `docker`, `git` and `OS` versions added to `sato version` command output
  - added key `--show-empty-groups` to the command `sato get secrets`
  - `sato help` button  refactored
  - README.md files updated
  
**Fixed:**

---

# Test results:

| Area          | Result | Details                                   |
|---------------|--------|-------------------------------------------|
| Preflight     | PASS   | gofmt, go mod tidy, go vet, go test ./... |
| Unit tests    | PASS   | 19 passed, coverage 14.5%                 |
| E2E tests     | PASS   | 9 passed, 0 failed, 0 skipped             |
| Binary        | PASS   | dist/sato, size 3.0M, commit 20ed140      |
| Final status  | PASS   | All test suites completed successfully    |

---

# Installation

## Prerequisites

* **required**: OS Linux | [docker](https://docs.docker.com/engine/install/) | [docker compose](https://docs.docker.com/compose/install/linux)

* **optional**: [keepassxc](https://keepassxc.org/download/#linux) | [git](https://git-scm.com/install/linux) | [go](https://go.dev/doc/install) |  [manually build bin/rpm/deb files](scripts/build/README.md)

## Install `sato`:
  * from `binary` file:
    ```bash
    wget https://github.com/Marcus-Aprelius/sato/releases/latest/download/sato && chmod +x sato && sudo cp sato /usr/local/bin/
    ```

  * from `.deb` package:
    ```bash
    wget https://github.com/Marcus-Aprelius/sato/releases/download/v0.0.3/sato_0.0.3_amd64.deb && sudo dpkg -i sato_0.0.3_amd64.deb
    ```

  * from `.rpm` package:
    ```bash
    wget https://github.com/Marcus-Aprelius/sato/releases/download/v0.0.3/sato-0.0.3-1.x86_64.rpm && sudo yum install -y sato-0.0.3-1.x86_64.rpm
    ```
