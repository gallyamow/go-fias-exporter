-- Эта таблица отсутствует в gar_schemas, но в нее есть данные в корне gar_data
CREATE TABLE tmp.addhouse_types (
	id VARCHAR NOT NULL PRIMARY KEY,
	name VARCHAR NOT NULL,
	shortname VARCHAR,
	"desc" VARCHAR,
	updatedate DATE NOT NULL,
	startdate DATE NOT NULL,
	enddate DATE NOT NULL,
	isactive BOOLEAN NOT NULL
);
COMMENT ON TABLE tmp.house_types IS 'Эта таблица отсутствует в gar_schemas, но в нее есть данные в корне gar_data';
