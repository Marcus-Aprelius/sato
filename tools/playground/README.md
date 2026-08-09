# Playground

Demo environment for `SATO` Project. In this folder, you can try `sato` functionality before deciding whether to use it in a real environment.

Also, this playground is used by tests scripts. See [tests/README.md](../../tests/README.md) for details.

## Files

* `playground_create.sh` - creates from scratch test files if they are absent or recreates them if they are present
* `playground_delete.sh` - deletes all created files

  | File                 | Description                                                  |
  |----------------------|--------------------------------------------------------------|
  | `.env`               | Non-secret Docker Compose variables for isolation tests      |
  | `docker-compose.yml` | Demo Docker Compose file consuming KeePass secrets           |
  | `secrets.kdbx`       | Test KeePass DB with sample secrets. Master password: `sato` |


## Playground lifecycle:

* create or recreate playground files:
  ```bash
  cd tools/playground && bash playground_create.sh
  ```

* delete playground files:
  ```bash
  cd tools/playground && bash playground_delete.sh
  ```

## Usage examples

* start `docker compose` with secrets loaded from KeePass:
  ```bash
  sato --db-path=secrets.kdbx docker compose up -d
  ```

* start `docker compose` with secrets loaded via env variable:
  ```bash
  export SATO_DB_PATH="$(pwd)/secrets.kdbx"
  sato docker compose config
  ```

* list secret names as a list
  ```bash
  sato --db-path=secrets.kdbx get secrets
  ```

* Get the value of the secret called TOKEN in format <TOKEN_VALUE> without any additional infirmation (to use in script)
  ```bash
  sato --db-path=secrets.kdbx get secret TOKEN -q
  ```

* list secret names as a group tree with an empty groups:
  ```bash
  sato --db-path=secrets.kdbx get secrets --tree --show-empty-groups
  ```
