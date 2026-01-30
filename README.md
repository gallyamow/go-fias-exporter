## go-fias-exporter

Exports FIAS-dump files to postgresql tables or outputs resulting SQL queries to console.

### Usage

```shell
fias-exporter [flags] <path>
```

| Flag           | Default | Description                                                                   |
|----------------|---------|-------------------------------------------------------------------------------|
| `--mode`       | `copy`  | "copy" - generates `COPY FROM csv`, "upsert" - generates `INSERT ON CONFLICT` |
| `--schema`     | ``      | database schema or public will used by default                                |
| `--batch-size` | `1000`  | minimum size of batch                                                         |

### Example

```shell
./fias-exporter --mode copy  ./examples > examples.sql
./fias-exporter --mode upsert  ./examples > examples.sql
```

Importing to postgresql:

```shell
docker run --name gar -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -d postgres:latest
psql -h localhost -d postgres -U postgres
```

### TODO

- import schemas and use it to create tmp tables
- type casting based on schema