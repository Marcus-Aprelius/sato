<img src="assets/sato_title.png" alt="sato">

> **`sato`** — like `sudo`, but for secrets

`sato` loads secrets from secure KeePass-compatible `.kdbx` database (DB) and runs **`docker compose`** commands using those variables (secrets). Secrets are not exposed to the shell or written to disk in such case.


## How It Works
1. `sato` searches for a `.kdbx` DB in predefined locations or paths specified by the user. Once a valid DB is found, it is used as the source of secrets. 
Additionally, `sato` can safely display secrets' names from the DB (`sato get secrets`).

2. `sato` reads the DB's master password and uses its secrets to run `docker compose` commands. Password input is not echoed to the terminal and secrets are passed only to child process, never exported to shell or written to temporary files.
  <img src="assets/compare.jpg" alt="compare">

<p><strong><span style="color:red">⚠ Pay Attention!</span></strong></p>

Utility `sato` is provided **"as is"** and its usage in a production environment is fully **at your own risk**!

---

[![Release](https://img.shields.io/github/v/release/Marcus-Aprelius/sato)](https://github.com/Marcus-Aprelius/sato/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)](https://go.dev/dl)
[![Git](https://img.shields.io/badge/git-2.55-F05032?logo=git&logoColor=orange)](https://git-scm.com/install/linux)
[![KeePassXC](https://img.shields.io/badge/KeePassXC-2.7-8A2BE2?logo=letsencrypt&logoColor=yellow)](https://keepassxc.org/download/#linux)
[![Docker](https://img.shields.io/badge/Docker-29.1-2496ED?logo=docker&logoColor=blue)](https://docs.docker.com/engine/install)
[![Docker Compose](https://img.shields.io/badge/Compose-2.4-1D63ED?logo=docker&logoColor=white)](https://docs.docker.com/compose/install/linux)

[![Binary](https://img.shields.io/badge/Binary-file-success)](https://github.com/Marcus-Aprelius/sato/releases/latest/download/sato)
[![DEB](https://img.shields.io/badge/DEB-package-C71585?logo=debian&logoColor=red)](https://github.com/Marcus-Aprelius/sato/releases/download/v0.0.3/sato_0.0.3_amd64.deb)
[![RPM](https://img.shields.io/badge/RPM-package-EE0000?logo=redhat&logoColor=red)](https://github.com/Marcus-Aprelius/sato/releases/download/v0.0.3/sato-0.0.3-1.x86_64.rpm)

---

## Installation

### Prerequisites

* **required**: OS Linux | [docker](https://docs.docker.com/engine/install/) | [docker compose](https://docs.docker.com/compose/install/linux)

* **optional**: [keepassxc](https://keepassxc.org/download/#linux) | [git](https://git-scm.com/install/linux) | [go](https://go.dev/doc/install) |  [manually build bin/rpm/deb files](scripts/build/README.md)

### Install `sato`:
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

---

## Flags and Commands:

* Flags:
  | Flag           | Description                         |
  |----------------|-------------------------------------|  
  | --db-path=PATH | Path to KeePass-compatible .kdbx DB |

* Commands:
  | Command                                       | Description                                                   |
  |-----------------------------------------------|---------------------------------------------------------------|
  | `sato`                                        | Show current status                                           |
  | `sato version`                                | Show version                                                  |
  | `sato help`                                   | Show help                                                     |
  | `sato completion bash`                        | Show Bash completion scripts                                  |
  | `sato get secrets`                            | List secret names from a KeePass-compatible `.kdbx` DB        |
  | `sato get secrets --tree`                     | List secret names as a group tree                             |
  | `sato get secrets --tree --show-empty-groups` | List secret names as a group tree, including empty groups     |
  | `sato get secret <NAME>`                      | Show value of a secret                                        |
  | `sato get secret <NAME> -q\|quite`            | Show value of a secret without any information (for scripts)  |
  | `sato docker compose <...>`                   | Run any Docker Compose command with passwords from `.kdbx` DB |
  | `sato docker compose up -d`                   | Example: start Docker containers in detached mode             |

---

## DB Locations Priority

| Priority    | Source/Location                   | Comment                                                           |
|-------------|-----------------------------------|-------------------------------------------------------------------|
| 1 (highest) | `--db-path=/path/to/secrets.kdbx` | **Specify DB location manually:**<br>if `set` - is used, ignores locations with lower priority<br>if `not set` - finds other locations                 |
| 2           | `~/.sato/secrets.kdbx`            | **Default location of the DB:**<br>if `present` - is used, ignores location with lower priority<br>if `absent`  - finds other locations                 |
| 3 (lowest)  | `SATO_DB_PATH`                    | **ENV variable:** (i.e.: `export SATO_DB_PATH=/path/to/secrets.kdbx`)<br>if `set`     - is used<br>if `not set` - finds other locations |                                          |

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
  sato completion bash >> ~/.bashrc && source ~/.bashrc 
  ```

---

## Development

It's not a requirement, rather a `general recommendation` that all development tasks can be done in two ways:

| Development with | Reccomended for                                                                                                               |
|------------------|-------------------------------------------------------------------------------------------------------------------------------|
| `devcontainer`   | - development and fast checks<br>- coding, formatting, tests, `go vet`<br>- quick CLI commands (without build) |
| Docker image     | - release/build process<br>- scripts/build/create.sh                                                                          |

Examples of scripts:

| Devcontainers                                                         | Docker Images                                                 |
|-----------------------------------------------------------------------|---------------------------------------------------------------|
| **1. Format Go code:**<br>`gofmt -w internal/sato/*.go tests/unit/*.go`<br>  | **1. Official release build:**<br>`bash scripts/build/create.sh bin` |
| **2. Quick Go tests:**<br>`go test ./...`<br>`go vet ./...`<br>`go mod tidy` | **2. Package builds:**<br>`bash scripts/build/create.sh deb`<br>`bash scripts/build/create.sh rpm`<br>`bash scripts/build/create.sh all` |
| **3. Work with playground files:**<br>`bash playground/playground_create.sh`<br>`bash playground/playground_delete.sh`  | **3. Final pre-release verification on host:**<br>`bash tests/run_all_tests.sh`<br>`bash scripts/build/create.sh all`                         |
| **4. `sato` CLI commands:**<br>`go run . help`<br>`go run . version`                                                    | |
| **5. Work with playground files:**<br>`bash playground/playground_create.sh`<br>`bash playground/playground_delete.sh` | |
| **6. Tests:**<br>`bash tests/run_unit_tests.sh`<br>`bash tests/run_e2e_tests.sh`<br>`bash tests/run_all_tests.sh`       | |
| **7. Docker / Docker Compose checks**<br>`docker version`<br>`docker compose version`                                   | |


## Playground and Testing

See [playground](playground/README.md) and [tests](tests/README.md) for the details.

---

[Apache License 2.0](LICENSE)

© 2026 [Marcus-Aprelius](https://github.com/Marcus-Aprelius/sato)
