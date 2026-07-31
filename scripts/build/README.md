# create.sh

Build and package `sato` using Docker (inside Docker containers).

No local Go installation or RPM/DEB tooling is required

You can manually invoke any Go command, for example:
  ```bash
  docker run --rm -v "$PWD:/src" -w /src golang:1.26.5-alpine go mod tidy
  ```

All compiled files will be in folder `dist`, for example:
```text
dist
├── sato
├── sato-0.0.3-1.x86_64.rpm
└── sato_0.0.3_amd64.deb
```

---

### Build `sato`:

* Build `binary` file:
    ```bash
    git clone https://github.com/Marcus-Aprelius/sato.git && cd sato/scripts/build  && bash create.sh bin
    ```

* Build `.deb` package:
    ```bash
    git clone https://github.com/Marcus-Aprelius/sato.git && cd scripts/build && bash create.sh deb
    ```

* Build `.rpm` package:
    ```bash
    git clone https://github.com/Marcus-Aprelius/sato.git && cd scripts/build && bash create.sh rpm
    ```

* Build all packages (`bin`/`.deb`/`.rpm`) at once:
    ```bash
    git clone https://github.com/Marcus-Aprelius/sato.git && cd scripts/build && bash create.sh all
    ```
