## go-fias-exporter

Reads FIAS-dump files and outputs resulting SQL queries to console.

## Features

* Supports two export modes: `COPY (fast bulk import)` and `UPSERT` (merge/update existing data)
* Configurable batch processing for optimal performance
* Schema support for organized database structure
* Generates SQL files for later execution or direct pipeline import

## Installation

```shell
make build

# Or install directly
go install github.com/gallywow/go-fias-exporter
```

## Usage

```shell
fias-exporter [flags] <path-to-fias-dump>
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
./fias-exporter --schema=tmp ./examples > examples.sql

docker run --name gar -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -d postgres:latest
docker exec -i gar psql -U postgres < ./create-tmp-tables.sql
docker exec -i gar psql -U postgres < ./examples.sql
```

### TODO

- import schemas and use it to create tmp tables
- type casting based on schema