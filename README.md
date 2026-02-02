# go-fias-exporter

Преобразует XML-выгрузки ФИАС (ГАР) в SQL, пригодный для импорта в PostgreSQL.

## Возможности

- Поддержка трех режимов экспорта:
    - **schema** — генерация запросов на `CREATE TABLE`
    - **copy** — генерация запросов на пакетный импорт с использованием `COPY FROM STDIN`
    - **upsert** — генерация запросов добавление и обновление существующих данных через `INSERT … ON CONFLICT`
- Настраиваемый размер батча для оптимальной производительности
- Результат можно:
    - сохранить в файлы
    - передать `psql` для потокового импорта

## Установка

```shell
make build
```

## Использование

```shell
fias-exporter [flags] <путь-к-выгрузке-ФИАС>
```

| Флаг           | По умолчанию | Описание                                                                                                           |
|----------------|--------------|--------------------------------------------------------------------------------------------------------------------|
| `--mode`       | `copy`       | Режим экспорта: `schema` — генерация `CREATE TABLE`, `copy` — `COPY FROM STDIN`, `upsert` — `INSERT … ON CONFLICT` |
| `--db-schema`  | `public`     | Целевая схема базы данных                                                                                          |
| `--batch-size` | `1000000`    | Количество записей в одном batch                                                                                   |

## Пример

```shell
docker pull postgres:latest
docker run --name gar \
  -p 5432:5432 \
  -e POSTGRES_HOST_AUTH_METHOD=trust \
  -d postgres:latest

# 1) Создание схемы (если не хотите в public)
echo 'CREATE SCHEMA tmp;' | docker exec -i gar psql -U postgres

# 2) Импорт таблиц в созданную схему
./fias-exporter --mode schema --db-schema tmp ./example/gar_schemas | docker exec -i gar psql -U postgres -v ON_ERROR_STOP=1

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

# 3) Потоковый импорт данных в созданные таблицы
./fias-exporter --mode copy --db-schema tmp ./example/gar_data | docker exec -i gar psql -U postgres -v ON_ERROR_STOP=1
```
