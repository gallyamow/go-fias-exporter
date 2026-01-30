## go-fias-exporter

Exports FIAS-dump files to postgresql tables or outputs resulting SQL queries to console.

### Usage

```shell
fias-exporter [flags] <path>
```

| Flag           | Default | Description                                                                   |
|----------------|---------|-------------------------------------------------------------------------------|
| `--mode`       | `copy`  | "copy" - generates `COPY FROM csv`, "upsert" - generates `INSERT ON CONFLICT` |
| `--batch-size` | `1000`  | Minimum size of batch                                                         |

### Example

```shell
./fias-exporter --mode copy  ./examples > examples.sql
./fias-exporter --mode upsert  ./examples > examples.sql
```

### TODO

- import schemas and use it to create tmp tables
- type casting based on schema