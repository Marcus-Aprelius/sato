```text
   _____      _______ ____
  / ____|  /\|__   __| __ \    Secure
 | (___   /  \  | | | |  | |    Access
  \___ \ / /\ \ | | | |  | |     Task
  ____) | ____ \| | | |__| |      Operator
 |_____/_/    \_\_|  \____/    > like sudo, but for secrets
```

# Changes:

## Changed:
  - docker image migrated in `devcontainer.json` migrated from `golang:1.26.5-alpine` to `golang:1.26.5-bookworm` for better compability with Go
  - verbose mode disabled in file `create.sh` to reduce output
  - completion code adjusted - added `sato git` subcommand
  - codebase refactored -file `run.go` splitted into several files:
    - `commands.go`
    - `db_path.go`
    - `git.go`
    - `version.go`
    - `tree.go`
    - `secrets.go`
    - `help`
  - README.md files updated
  - tests adjusted to test `-q` case
  
## Added:
  - subcommand `sato git` to use git with secrets from DB (`.kdbx` file)
  - subcommands `sato completion bash add|delete|update|status` for management of bash completion in ~/.bashrc
  - flag `--secret` to override default secret name (by default equals repo name) for the command `sato --secret <NAME> git clone <URL>`

---

# Test results:

| Area          | Result | Details                                   |
|---------------|--------|-------------------------------------------|
| Preflight     | PASS   | gofmt, go mod tidy, go vet, go test ./... |
| Unit tests    | PASS   | 25 passed, coverage 21.5%                 |
| E2E tests     | PASS   | 10 passed, 0 failed, 0 skipped            |
| Binary        | PASS   | dist/sato, size 3.1M, commit b203b54      |
| Final status  | PASS   | All test suites completed successfully    |

---

[Apache License 2.0](LICENSE)

© 2026 **[Marcus-Aprelius](https://github.com/Marcus-Aprelius/sato)**

**Discord:** Marcus.Aprelius.Antoninus
