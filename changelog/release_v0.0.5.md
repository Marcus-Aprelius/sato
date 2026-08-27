```text
   _____      _______ ____
  / ____|  /\|__   __| __ \    Secure
 | (___   /  \  | | | |  | |    Access
  \___ \ / /\ \ | | | |  | |     Task
  ____) | ____ \| | | |__| |      Operator
 |_____/_/    \_\_|  \____/    > like sudo, but for secrets
```

## Changed:
  * `README.md` updated
  * `.gitignore` updated
  * `.dockerignore` updated
  * tests updated
  * `Dockerfile` updated
  * code rafactored

## Added:
  * Support of `.psafe3` and `.ibak` vaults

---

# Test results:

| Area          | Result | Details                                   |
|---------------|--------|-------------------------------------------|
| Preflight     | PASS   | gofmt, go mod tidy, go vet, go test ./... |
| Unit tests    | PASS   | 31 passed, coverage 20.5%                 |
| E2E tests     | PASS   | 11 passed, 0 failed, 0 skipped            |
| Binary        | PASS   | dist/sato, size 3.2M, commit 9947986      |
| Final status  | PASS   | All test suites completed successfully    |
