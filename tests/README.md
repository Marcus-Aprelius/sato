# Tests
Tests for the `SATO` project.

Tests use [playground](playground/README.md) for testing purposes.

Tests consists of two types:
- unit (Go)
- e2e (Shell)


## Unit tests

They test internal `sato` logic, including:
- command allow-list
- DB path discovery
- KeePass loading
- recursive group reading
- empty group detection
- CLI output separation
- completion generation

### Unit test flow:
- creates required files in the playground
- runs Go unit tests
- deletes created files in the playground
- writes test results to the file `report_unit_tests.txt`


## E2E tests
They test the compiled `sato` binary end-to-end, including:
- version output
- command filtering
- KeePass to Docker Compose integration
- parent environment isolation
- .env isolation
- stdout/stderr separation
- get secrets CLI output
- get secrets tree output
- CLI error paths

### E2E test flow:
- builds `sato` binary
- creates required files in the playground
- runs all E2E shell tests
- deletes created files in the playground
- writes test results to the file `report_e2e_tests.txt`


## All tests
Ability to trigger all tests at once.

### All test flow:
- triggers unit tests
- triggers e2e tests
- writes test results to the file `report_all_tests.txt`


## Run tests

### unit:
```bash
cd tests && bash run_unit_tests.sh
```

### end-to-end:
```bash
cd tests && bash run_e2e_tests.sh
```

### all:
```bash
cd tests && bash run_all_tests.sh
```
