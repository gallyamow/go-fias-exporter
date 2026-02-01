## go-fias-exporter

Transforms FIAS GAR XML dumps into SQL suitable for PostgreSQL import.

## Features

* Supports two export modes:
  - COPY — fast bulk import using COPY FROM STDIN
  - UPSERT — merge/update existing data using INSERT … ON CONFLICT
* Configurable batch size for optimal performance
* Generates SQL output for:
  - saving to files
  - direct pipelined import into PostgreSQL

## Installation

```shell
make build
```

## Usage

```shell
fias-exporter [flags] <path-to-fias-dump>
```

| Flag           | Default  | Description                                                                                                 |
|----------------|----------|-------------------------------------------------------------------------------------------------------------|
| `--mode`       | `copy`   | Export mode: `schema` - generates CREATE TABLE, `copy` (COPY FROM STDIN) or `upsert` (INSERT … ON CONFLICT) |
| `--db-schema`  | `public` | Target database schema                                                                                      |
| `--batch-size` | `1000`   | Number of records per batch                                                                                 |

### Example

```shell
docker pull postgres:latest
docker run --name gar \
  -p 5432:5432 \
  -e POSTGRES_HOST_AUTH_METHOD=trust \
  -d postgres:latest

# create tables
echo 'CREATE SCHEMA tmp;' | docker exec -i gar psql -U postgres
./fias-exporter --mode schema --db-schema tmp ./example/gar_schemas | docker exec -i gar psql -U postgres

# pipelined data import
./fias-exporter --mode copy --db-schema tmp ./example/gar_data | docker exec -i gar psql -U postgres
```

### TODO

- generate database tables from FIAS schemas with column type casting
