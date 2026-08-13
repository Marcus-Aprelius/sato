<img src="assets/sato_title.png" alt="sato">

> **`sato`** — like `sudo`, but for secrets

`sato` loads secrets from supported encrypted vaults (`.kdbx`, `.psafe3`, `.ibak`) and securely provides them to supported **Docker Compose** and **Git** commands.

`SATO` does not export secrets to the current shell or write secret values to disk.

> Utility `sato` is provided **"as is"** and its usage in a production environment is fully **at your own risk**!

## How It Works
1. `sato` searches for a [supported encrypted vault](#supported-formats) in predefined locations or paths specified by the user. Once a valid vault is found, it is used as the source of secrets.
Additionally, `sato` can safely display secrets' names from the DB (`sato get secrets`).

2. `sato` reads the DB's master password and uses its secrets to run `docker compose ...`, `sato git ...` commands. Password input is not echoed to the terminal and secrets are passed only to the child process, never exported to shell or written to temporary files.
  <img src="assets/compare.jpg" alt="compare">

---

[![Release](https://img.shields.io/github/v/release/Marcus-Aprelius/sato)](https://github.com/Marcus-Aprelius/sato/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go)](https://go.dev/dl)
[![Git](https://img.shields.io/badge/git-2.55-F05032?logo=git&logoColor=orange)](https://git-scm.com/install/linux)
[![KeePassXC](https://img.shields.io/badge/KeePassXC-2.7-8A2BE2?logo=letsencrypt&logoColor=yellow)](https://keepassxc.org/download/#linux)
[![Docker](https://img.shields.io/badge/Docker-29.7-2496ED?logo=docker&logoColor=blue)](https://docs.docker.com/engine/install)
[![Docker Compose](https://img.shields.io/badge/Compose-5.5-1D63ED?logo=docker&logoColor=white)](https://docs.docker.com/compose/install/linux)

[![Binary](https://img.shields.io/badge/Binary-file-success)](https://github.com/Marcus-Aprelius/sato/releases/latest/download/sato)
[![DEB](https://img.shields.io/badge/DEB-package-C71585?logo=debian&logoColor=red)](https://github.com/Marcus-Aprelius/sato/releases/download/v0.0.5/sato_0.0.5_amd64.deb)
[![RPM](https://img.shields.io/badge/RPM-package-EE0000?logo=redhat&logoColor=red)](https://github.com/Marcus-Aprelius/sato/releases/download/v0.0.5/sato-0.0.5-1.x86_64.rpm)

---

## Supported formats

| Encrypted Vaults | Description               |
|------------------|---------------------------|
| `.kdbx`          | KeePass 2 database        | 
| `.psafe3`        | Password Safe V3 database |
| `.ibak`          | Password Safe V3 backup   |

---
## Installation

### Prerequisites

* **required**: OS Linux | [docker](https://docs.docker.com/engine/install/) | [docker compose](https://docs.docker.com/compose/install/linux)

* **optional**: [keepassxc](https://keepassxc.org/download/#linux) | [git](https://git-scm.com/install/linux) | [go](https://go.dev/doc/install) |  [manually build bin/rpm/deb files](tools/build/README.md)

### Install `sato`:
  * from `binary` file:
    ```bash
    wget https://github.com/Marcus-Aprelius/sato/releases/latest/download/sato && chmod +x sato && sudo cp sato /usr/local/bin/
    ```

  * from `.deb` package:
    ```bash
    wget https://github.com/Marcus-Aprelius/sato/releases/download/v0.0.5/sato_0.0.5_amd64.deb && sudo dpkg -i sato_0.0.5_amd64.deb
    ```

  * from `.rpm` package:
    ```bash
    wget https://github.com/Marcus-Aprelius/sato/releases/download/v0.0.5/sato-0.0.5-1.x86_64.rpm && sudo yum install -y sato-0.0.5-1.x86_64.rpm
    ```

---

## Flags and Commands:

* Flags:
  | Flag             | Description                                                          |
  |------------------|----------------------------------------------------------------------|
  | --secret=<NAME>  | Name of the secret for `sato git ...`                                |
  | --db-path=<PATH> | Path to a database ([supported encrypted vault](#supported-formats)) |

* Commands:
  | Command                                                      | Description                                                          |
  |--------------------------------------------------------------|----------------------------------------------------------------------|
  | `sato`                                                       | Show current status                                                  |
  | `sato help`                                                  | Show help                                                            |
  | `sato version`                                               | Show version                                                         |
  | `sato completion bash`                                       | Show Bash completion script                                          |
  | `sato completion bash add\|delete\|update\|status`           | Add/delete update or show status of Bash completion script           |
  |||          
  | `sato get secrets`                                           | List secret names from Secret database (`.kdbx`, `.psafe3`, `.ibak`) |
  | `sato get secrets --tree`                                    | List secret names as a group tree                                    |
  | `sato get secrets --tree --show-empty-groups`                | List secret names as a group tree, including empty groups            |
  | `sato get secret <NAME>`                                     | Show value of a secret                                               |
  | `sato get secret <NAME> -q\|quite`                           | Show value of a secret without any information (for scripts)         |
  |||
  | `sato docker compose <...>`                                  | Run any Docker Compose command with secrets from the configured vault        |
  | `sato docker compose up -d`                                  | Example: start Docker containers in detached mode                            |
  | `sato docker git help`                                       | Show SATO docker help                                                        |
  |||
  | `sato git clone\|push\|pull\|fetch`                          | Run allowed git command with a secret (TOKEN) from from the configured vault |
  | `sato git clone https://github.com/Marcus-Aprelius/sato.git` | Example: clone a repository using a vault secret                             |
  | `sato git help`                                              | Show SATO git help                                                           |

---

## Database Location Priority

| Priority    | Source/Location                                 | Comment                                                           |
|-------------|-------------------------------------------------|-------------------------------------------------------------------|
| 1 (highest) | `--db-path=/path/to/secrets.{kdbx,psafe3,ibak}` | **Specify DB location manually:**<br>if `set` - is used, ignores locations with lower priority<br>if `not set` - finds other locations                 |
| 2           | `~/.sato/*.{kdbx,psafe3,ibak}`                  | **Default location of the DB:**<br>if `present` - is used, ignores location with lower priority<br>if `absent`  - finds other locations                 |
| 3 (lowest)  | `SATO_DB_PATH`                                  | **ENV variable:** (i.e.: `export SATO_DB_PATH=/path/to/secrets.{kdbx,psafe3,ibak}`)<br>if `set`     - is used<br>if `not set` - finds other locations |                                          |

<span style="color:orange">! Pay attention !</span>

1. After specifying DB location, `sato` validates its presence to prevent corruption. Only valid DB locations are used; invalid - ignored as if they were absent.

2. If no valid DB location is set (and DB is absent in the default location), the `sato docker compose` command will not work.

    Specify a valid DB location or place the DB in the default location.

3. If several valid DB locations are available - `sato` will use DB with the highest priority.

---

## Bash Completion
Enable tab completion for `sato` commands:
* for current session: 
  ```bash
  source <(sato completion bash)
  ```
* permanently:
  ```bash
  echo "source <(/usr/local/bin/sato completion bash)" >> ~/.bashrc && source ~/.bashrc
  ```
  or with checks:
  ```bash
  if command -v sato >/dev/null 2>&1; then
    sato completion bash > ~/.sato-completion.bash &&
    (grep -qxF 'source ~/.sato-completion.bash' ~/.bashrc || 
    echo 'source ~/.sato-completion.bash' >> ~/.bashrc) &&
    source ~/.bashrc
  fi
  ```
  or using `sato`:
  ```bash
  sato completion bash add && source ~/.bashrc
  ```

---

## Development

It's not a requirement, rather a `general recommendation` that all development tasks can be done in two ways:

| Development with | Reccomended for                                                                                                               |
|------------------|-------------------------------------------------------------------------------------------------------------------------------|
| `devcontainer`   | - development and fast checks<br>- coding, formatting, tests, `go vet`<br>- quick CLI commands (without build) |
| Docker image     | - release/build process<br>- tools/build/create.sh                                                                          |

Examples of scripts:

| Devcontainers                                                                | Docker Image                                                  |
|------------------------------------------------------------------------------|---------------------------------------------------------------|
| **1. Format Go code:**<br>`gofmt -w internal/sato/*.go tests/unit/*.go`<br>  | **1. Official release build:**<br>`bash create.sh bin`        |
| **2. Quick Go tests:**<br>`go test ./...`<br>`go vet ./...`<br>`go mod tidy` | **2. Package builds:**<br>`bash create.sh deb`<br>`bash create.sh rpm`<br>`bash create.sh all` |
| **3. Work with playground files:**<br>`bash playground_create.sh`<br>`bash playground_delete.sh`    | **3. Final pre-release verification on host:**<br>`bash run_all_tests.sh`<br>`bash create.sh all`                                            |
| **4. `sato` CLI commands:**<br>`go run . help`<br>`go run . version`                              | |
| **5. Work with playground files:**<br>`bash playground_create.sh`<br>`bash playground_delete.sh`  | |
| **6. Tests:**<br>`bash run_unit_tests.sh`<br>`bash run_e2e_tests.sh`<br>`bash run_all_tests.sh`   | |
| **7. Docker / Docker Compose checks**<br>`docker version`<br>`docker compose version`             | |

## Playground and Testing

See [playground](tools/playground/README.md) and [tests](tests/README.md) for the details.

---

[Apache License 2.0](LICENSE)

[VS Code extention SATO secrets](https://marketplace.visualstudio.com/items?itemName=MarcusApreliusAntoninus.sato-vscode-ext) allows work with Vaults and Crypto files.

**Discord:** Marcus.Aprelius.Antoninus

© 2026 **[Marcus-Aprelius](https://github.com/Marcus-Aprelius/sato)**
