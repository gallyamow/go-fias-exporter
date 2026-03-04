# go-fias-exporter

Преобразует XML-выгрузки ФИАС (ГАР) в SQL, пригодный для импорта в **PostgreSQL** и **MySQL**.

<a href="https://pkg.go.dev/github.com/gallyamow/go-fias-exporter"><img src="https://pkg.go.dev/badge/github.com/gallyamow/go-fias-exporter.svg" alt="Go Reference"></a>

## Возможности

- Поддержка трех режимов экспорта:
    - **schema** — генерация запросов на `CREATE TABLE`
    - **keys** — генерация запросов на `ADD PRIMARY KEY`
    - **bulk** — генерация запросов на пакетный импорт с использованием `COPY FROM STDIN` (PostgreSQL) или
      `LOAD DATA LOCAL INFILE` (MySQL)
    - **upsert** — генерация запросов добавление и обновление существующих данных через `INSERT … ON CONFLICT` (
      PostgreSQL) или `INSERT … ON DUPLICATE KEY UPDATE` (MySQL)
- Поддержка баз данных:
    - **PostgreSQL** (по умолчанию)
    - **MySQL**
- Настраиваемый размер batch для оптимальной производительности
- Результат можно:
    - сохранить в файлы
    - передать клиенту БД (psql или mysql) для потокового импорта

## Установка

```shell
# (для linux можно просто скачать из releases)
make build
```

## Использование

```shell
fias-exporter [flags] <путь-к-выгрузке-ФИАС>
```

| Флаг                | По умолчанию | Описание                                                                                                                                                                                                              |
|---------------------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--mode`            | `bulk`       | Режим экспорта: `schema` - генерация `CREATE TABLE`, `keys` - генерация `ADD PRIMARY KEY`, `bulk` - `COPY FROM STDIN`/`LOAD DATA LOCAL INFILE`, `upsert` - `INSERT … ON CONFLICT`/`INSERT … ON DUPLICATE KEY UPDATE`. |
| `--db-type`         | `postgres`   | Тип базы данных: `postgres` или `mysql`.                                                                                                                                                                              |
| `--db-schema`       | `public`     | Целевая схема базы данных.                                                                                                                                                                                            |
| `--batch-size`      | `1000000`    | Количество записей в одном batch.                                                                                                                                                                                     |
| `--ignore-required` | true         | Добавлять ли `NOT NULL` для `use="required"` полей.                                                                                                                                                                   |

## Пример

### PostgreSQL

```shell
docker pull postgres:latest
docker run --name gar-psql \
  -p 5432:5432 \
  -e POSTGRES_HOST_AUTH_METHOD=trust \
  -d postgres:latest

# 1) Создание схемы (если не хотите в public)
echo 'CREATE SCHEMA tmp;' | docker exec -i gar-psql psql -U postgres

# 2) Создание таблиц по схемам
# (по умолчанию без PRIMARY KEY и NOT NULL)
./fias-exporter --db-type postgres --mode schema --db-schema tmp ./example/gar_schemas | docker exec -i gar-psql psql -U postgres -v ON_ERROR_STOP=1

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
);' | docker exec -i gar-psql psql -U postgres

# Быстрый импорт данных в созданные таблицы
./fias-exporter --db-type postgres --mode bulk --db-schema tmp ./example/gar_data | docker exec -i gar-psql psql -U postgres -v ON_ERROR_STOP=1

# Добавление PRIMARY KEY ко всем таблицам (по схемам)
# (после импорта данных)
./fias-exporter --db-type postgres --mode keys --db-schema tmp ./example/gar_schemas | docker exec -i gar-psql psql -U postgres -v ON_ERROR_STOP=1

# Альтернативно: UPSERT
# (будет работать только если созданы PRIMARY KEYS)
./fias-exporter --db-type postgres --mode upsert --db-schema tmp ./example/gar_data | docker exec -i gar-psql psql -U postgres -v ON_ERROR_STOP=1

# Проверка
echo 'SELECT COUNT(*) FROM tmp.addhouse_types;' | docker exec -i gar-psql psql -U postgres -v ON_ERROR_STOP=1
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
# (по умолчанию без PRIMARY KEY и NOT NULL)
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
# ./fias-exporter --db-type mysql --mode bulk ./example/gar_data | docker exec -i gar-mysql mysql -u root gar --local-infile=1

# Добавление PRIMARY KEY ко всем таблицам (по схемам)
# (после импорта данных)
./fias-exporter --db-type mysql --mode keys ./example/gar_schemas | docker exec -i gar-mysql mysql -u root gar

# Альтернативно: UPSERT
# (будет работать только если созданы PRIMARY KEYS)
./fias-exporter --db-type mysql --mode upsert ./example/gar_data | docker exec -i gar-mysql mysql -u root gar

# Проверка
echo 'SELECT COUNT(*) FROM addhouse_types;' | docker exec -i gar-mysql mysql -u root gar
```

## Примечания

* Импортировался полный дамп, процесс подробно [описан](https://github.com/gallyamow/go-fias-exporter/issues/2), спасибо @Djoongaar
* По умолчанию таблицы создаются без `PRIMARY KEY` - чтобы не пересчитывался индекс и `NOT NULL` -
  чтобы процесс импорта не прервался из-за ошибок в данных
* Если используется `upsert`-режим (например для импорта дельты), то `PRIMARY KEY` должны быть заранее созданы.
* `--ignore-requried` нужен так как в `data-части` есть записи с пустыми колонками в полях с `use="required"`.
* В ходе работы в stderr будет выводиться статистика и вывод от клиента БД
* Процесс лучше запускать с остановкой на ошибках (пример для psql ON_ERROR_STOP=1)