# create.sh

Locally Go installed or RPM/DEB tooling is not required since docker image is used.

Created files are located in folder `dist`.

Examples of usage Go commands:

  ```bash
  docker run --rm -v "$PWD:/src" -w /src golang:1.26.7-alpine go mod tidy
  ```

Or create `bash` alias:
  ```bash
  alias go='docker run --rm -v $(pwd):/src -w /src golang:1.26.7-alpine go'
  ```
  ```bash
  alias gofmt='docker run --rm -v $(pwd):/src -w /src golang:1.26.7-alpine gofmt'
  ```

---

### Build `sato`:
* Clone repository:
  ```bash
  git clone https://github.com/Marcus-Aprelius/sato.git && cd tools/build
  ```
* Build `binary` file:
  ```bash
  bash create.sh bin
  ```
* Build `.deb` package:
  ```bash
  bash create.sh deb
  ```
* Build `.rpm` package:
  ```bash
  bash create.sh rpm
  ```
* Build all packages (`bin`, `.deb`, `.rpm`) at once:
  ```bash
  bash create.sh all
  ```
