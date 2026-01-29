## go-fias-exporter

Exports FIAS-dump to database tables.

### Usage

```shell
fias-exporter [flags] <database> <path>
```

| Flag           | Default | Description                                                        |
|----------------|---------|--------------------------------------------------------------------|
| `--batch-size` | `1000`  | Minimum size of inserts                                            |
| `--delta`      | `""`    | Number of delta                                                    |
| `--replace`    | `true`  | Replace already existing items, commonly used when importing delta |

### Example

```shell
./fias-exporter --exclude-dir=node_modules,vendor,venv,cache,.gradle --min-size=1MB --workers=32 --format=plain /home/ramil/Projects
```