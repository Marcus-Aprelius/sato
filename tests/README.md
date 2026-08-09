# Tests
Tests for the `SATO` project use [playground](../tools/playground/README.md) for their own purposes.

There is an ability to trigger all tests at once and get results to the file `report_all_tests.txt`

---

## Tests consists of two types:

### 1. Unit tests (test internal logic):
- command allow-list
- DB path discovery
- KeePass loading
- recursive group reading
- empty group detection
- CLI output separation
- completion generation

### Unit test workflow:
- creates required files in the playground
- runs Go unit tests
- deletes created files in the playground
- writes test results to the file `report_unit_tests.txt`

### 2. E2E tests (test compiled `sato` binary):
- version output
- command filtering
- KeePass to Docker Compose integration
- parent environment isolation
- .env isolation
- stdout/stderr separation
- get secrets CLI output
- get secrets tree output
- CLI error paths

### E2E test workflow:
- builds `sato` binary
- creates required files in the playground
- runs all E2E shell tests
- deletes created files in the playground
- writes test results to the file `report_e2e_tests.txt`

---

## Run tests
* unit:
  ```bash
  cd tests && bash run_unit_tests.sh
  ```

* end-to-end:
  ```bash
  cd tests && bash run_e2e_tests.sh
  ```

* all:
  ```bash
  cd tests && bash run_all_tests.sh
  ```
