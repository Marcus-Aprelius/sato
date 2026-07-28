# create.sh

Build and package `sato` using Docker (inside Docker containers).

No local Go installation is required or RPM/DEB tooling is required


### Build `sato`:

* Clone repository:
    ```bash
    git clone https://github.com/Marcus-Aprelius/sato.git
    cd sato
    ```

* Build `bin`ary file
    ```bash
    cd scripts/build && bash create.sh bin
    ```

* Build `.deb` package
    ```bash
    cd scripts/build && bash create.sh deb
    ```

* Build `.rpm` package
    ```bash
    cd scripts/build && bash create.sh rpm
    ```

* Build all packages (`bin`/`.deb`/`.rpm`)
    ```bash
    cd scripts/build && bash create.sh all
    ```


All compiled files will be in folder `dist`, for example:
```text
dist
├── sato
├── sato-0.0.2-1.x86_64.rpm
└── sato_0.0.2_amd64.deb
```
