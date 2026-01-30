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

Generate queries:

```shell
./fias-exporter --mode copy  ./examples 
./fias-exporter --mode upsert  ./examples
```

Importing to postgresql:

```shell
docker run --name gar -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -d postgres:latest
psql -h localhost -d postgres -U postgres

./fias-exporter --schema=tmp ./examples > examples.sql
psql -h localhost -d postgres -U postgres < ./examples.sql
```

### TODO

- import schemas and use it to create tmp tables
- type casting based on schema