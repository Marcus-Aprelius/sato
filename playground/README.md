# Playground

Demo environment for `SATO` Project. In this folder, you can try `sato` functionality before deciding whether to use it in a real environment.

Also, this playground is used by tests scripts. See [tests/README.md](tests/README.md) for details.

## Files

* `playground_create.sh`  - creates from scratch test files if they are absent.
If files are present, it will reacreate them.

  | File                 | Description                                                  |
  |----------------------|--------------------------------------------------------------|
  | `.env`               | Non-secret Docker Compose variables for isolation tests      |
  | `docker-compose.yml` | Demo Docker Compose file consuming KeePass secrets           |
  | `secrets.kdbx`       | Test KeePass DB with sample secrets. Master password: `sato` |

* `playground_delete.sh` - deletes all created files.


## Playground lifecycle:

* create or recreate playground files:
  ```bash
  cd playground && bash playground_create.sh
  ```
* delete playground files:
  ```bash
  cd playground && bash playground_delete.sh
  ```

## Usage

Run any `docker compose` command:

```bash
# with secrets loaded from KeePass:
echo "sato" | ../dist/sato --db-path=secrets.kdbx docker compose up -d

# with secrets loaded via env variable:
export SATO_DB_PATH="$(pwd)/secrets.kdbx"
echo "sato" | ../dist/sato docker compose config
```

List secret names:

```bash
# as a list
echo "sato" | ../dist/sato --db-path=secrets.kdbx get secrets

# as a group tree:
echo "sato" | ../dist/sato --db-path=secrets.kdbx get secrets --tree

# as a group tree with an empty groups:
echo "sato" | ../dist/sato --db-path=secrets.kdbx get secrets --tree --show-empty-groups
```
