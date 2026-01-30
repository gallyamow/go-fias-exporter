## go-fias-exporter

Exports FIAS-dump files to postgresql tables or outputs resulting SQL queries to console.

### Usage

```shell
fias-exporter [flags] <path>
```

| Flag           | Default  | Description                                                                            |
|----------------|----------|----------------------------------------------------------------------------------------|
| `--mode`       | `output` | "copy_from" - generates csv usable with COPY FROM, "upsert" - generates batched upsert |
| `--batch-size` | `1000`   | Minimum size of batch                                                                  |
| `--db`         | `nil`    | Database connection string                                                             |

### Example

```shell
./fias-exporter  ./example
```

### TODO

- import schemas and use it to create tmp tables
- type casting based on schema