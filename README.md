# go-fias-exporter

Преобразует XML-выгрузки ФИАС (ГАР) в SQL, пригодный для импорта в **PostgreSQL** и **MySQL**.

<a href="https://pkg.go.dev/github.com/gallyamow/go-fias-exporter"><img src="https://pkg.go.dev/badge/github.com/gallyamow/go-fias-exporter.svg" alt="Go Reference"></a>

## Возможности

- Поддержка трех режимов экспорта:
    - **schema** — генерация запросов на `CREATE TABLE`
    - **copy** — генерация запросов на пакетный импорт с использованием `COPY FROM STDIN` (PostgreSQL) или `LOAD DATA LOCAL INFILE` (MySQL)
    - **upsert** — генерация запросов добавление и обновление существующих данных через `INSERT … ON CONFLICT` (PostgreSQL) или `INSERT … ON DUPLICATE KEY UPDATE` (MySQL)
- Поддержка баз данных:
    - **PostgreSQL** (по умолчанию)
    - **MySQL**
- Настраиваемый размер батча для оптимальной производительности
- Результат можно:
    - сохранить в файлы
    - передать в клиент базы данных (psql или mysql) для потокового импорта

## Установка

```shell
make build
```

## Использование

```shell
fias-exporter [flags] <путь-к-выгрузке-ФИАС>
```

| Флаг                | По умолчанию | Описание                                                                                                                                                                        |
|---------------------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--mode`            | `copy`       | Режим экспорта: `schema` — генерация `CREATE TABLE`, `copy` — `COPY FROM STDIN`/`LOAD DATA LOCAL INFILE`, `upsert` — `INSERT … ON CONFLICT`/`INSERT … ON DUPLICATE KEY UPDATE`. |
| `--db-type`         | `postgres`   | Тип базы данных: `postgres` или `mysql`.                                                                                                                                        |
| `--db-schema`       | `public`     | Целевая схема базы данных.                                                                                                                                                      |
| `--batch-size`      | `1000000`    | Количество записей в одном batch.                                                                                                                                               |
| `--ignore-not-null` | false        | Игнорировать ли обработку `NOT NULL` по `use="required"` при определении колонок.                                                                                               |

## Пример

### PostgreSQL

```shell
docker pull postgres:latest
docker run --name gar \
  -p 5432:5432 \
  -e POSTGRES_HOST_AUTH_METHOD=trust \
  -d postgres:latest

# 1) Создание схемы (если не хотите в public)
echo 'CREATE SCHEMA tmp;' | docker exec -i gar psql -U postgres

# 2) Создание таблиц по схемам
./fias-exporter --db-type postgres --mode schema --db-schema tmp ./example/gar_schemas | docker exec -i gar psql -U postgres -v ON_ERROR_STOP=1

# По какой-то причине этой таблицы нет в gar_schemas
echo 'CREATE TABLE tmp.addhouse_types (
	id VARCHAR NOT NULL PRIMARY KEY,
	name VARCHAR NOT NULL,
	shortname VARCHAR,
	"desc" VARCHAR,
	updatedate DATE NOT NULL,
	startdate DATE NOT NULL,
	enddate DATE NOT NULL,
	isactive BOOLEAN NOT NULL
);' | docker exec -i gar psql -U postgres

# 3) Быстрый импорт данных в созданные таблицы
./fias-exporter --db-type postgres --mode copy --db-schema tmp ./example/gar_data | docker exec -i gar psql -U postgres -v ON_ERROR_STOP=1

# Альтернативно: UPSERT
./fias-exporter --db-type postgres --mode upsert --db-schema tmp ./example/gar_data | docker exec -i gar psql -U postgres -v ON_ERROR_STOP=1

# Проверка
echo 'SELECT COUNT(*) FROM tmp.addhouse_types;' | docker exec -i gar psql -U postgres -v ON_ERROR_STOP=1
```

### MySQL

```shell
docker pull mysql:latest
docker run --name gar-mysql \
  -p 3306:3306 \
  -e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
  -d mysql:latest

# 1) Создание базы данных (если не хотите в default)
docker exec -i gar-mysql mysql -u root -e "CREATE DATABASE gar;"

# 2) Создание таблиц по схемам
./fias-exporter --db-type mysql --mode schema ./example/gar_schemas | docker exec -i gar-mysql mysql -u root gar

# По какой-то причине этой таблицы нет в gar_schemas
echo 'CREATE TABLE addhouse_types (
	id VARCHAR(500) NOT NULL PRIMARY KEY,
	name TEXT NOT NULL,
	shortname TEXT,
	`desc` TEXT,
	updatedate DATE NOT NULL,
	startdate DATE NOT NULL,
	enddate DATE NOT NULL,
	isactive BOOLEAN NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;' | docker exec -i gar-mysql mysql -u root gar

# 3) Быстрый импорт данных в созданные таблицы
# TODO: пока не работает
# echo 'SET GLOBAL local_infile = 1;' | docker exec -i gar-mysql mysql -u root
# ./fias-exporter --db-type mysql --mode copy ./example/gar_data | docker exec -i gar-mysql mysql -u root gar --local-infile=1

# Альтернативно: UPSERT
./fias-exporter --db-type mysql --mode upsert ./example/gar_data | docker exec -i gar-mysql mysql -u root gar

# Проверка
echo 'SELECT COUNT(*) FROM addhouse_types;' | docker exec -i gar-mysql mysql -u root gar
```

## Примечания

* импортировался полный дамп, процесс подробно [описан](/issues/2)
* `--ignore-not-null` нужен так как в `data-части` есть записи с пустыми колонками в полях с `use="required"`.
