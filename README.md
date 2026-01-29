## go-fias-exporter

Exports FIAS-dump files to postgresql tables or outputs resulting SQL queries to console.

### Usage

```shell
fias-exporter [flags] <path>
```

| Flag           | Default  | Description                                                                     |
|----------------|----------|---------------------------------------------------------------------------------|
| `--mode`       | `output` | "output" - to print result in console, "execute" - to execute resulting queries |
| `--batch-size` | `1000`   | Minimum size of batch                                                           |
| `--delta`      | `nil`    | Delta key                                                                       |
| `--db`         | `nil`    | Database connection string                                                      |

### Example

```shell
./fias-exporter --exclude-dir=node_modules,vendor,venv,cache,.gradle --min-size=1MB --workers=32 --format=plain /home/ramil/Projects
```