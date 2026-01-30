CREATE SCHEMA tmp;

CREATE TABLE tmp.addhouse_types (
                                    id INTEGER NOT NULL PRIMARY KEY,
                                    name VARCHAR NOT NULL,
                                    shortname VARCHAR,
                                    "desc" VARCHAR,
                                    updatedate DATE NOT NULL,
                                    startdate DATE NOT NULL,
                                    enddate DATE NOT NULL,
                                    isactive BOOLEAN NOT NULL
);
COMMENT ON TABLE tmp.addhouse_types IS 'Дополнительные сведения по типам домов';
COMMENT ON COLUMN tmp.addhouse_types.id IS 'Идентификатор';
COMMENT ON COLUMN tmp.addhouse_types.name IS 'Наименование';
COMMENT ON COLUMN tmp.addhouse_types.shortname IS 'Краткое наименование';
COMMENT ON COLUMN tmp.addhouse_types.desc IS 'Описание';
COMMENT ON COLUMN tmp.addhouse_types.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.addhouse_types.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.addhouse_types.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.addhouse_types.isactive IS 'Статус активности';

CREATE TABLE tmp.addr_obj (
                              id BIGINT NOT NULL PRIMARY KEY,
                              objectid BIGINT NOT NULL,
                              objectguid VARCHAR NOT NULL,
                              changeid BIGINT NOT NULL,
                              name VARCHAR NOT NULL,
                              typename VARCHAR NOT NULL,
                              level VARCHAR NOT NULL,
                              opertypeid INTEGER NOT NULL,
                              previd BIGINT,
                              nextid BIGINT,
                              updatedate DATE NOT NULL,
                              startdate DATE NOT NULL,
                              enddate DATE NOT NULL,
                              isactual INTEGER NOT NULL,
                              isactive INTEGER NOT NULL
);
COMMENT ON TABLE tmp.addr_obj IS 'Сведения классификатора адресообразующих элементов';
COMMENT ON COLUMN tmp.addr_obj.id IS 'Уникальный идентификатор записи. Ключевое поле';
COMMENT ON COLUMN tmp.addr_obj.objectid IS 'Глобальный уникальный идентификатор адресного объекта типа INTEGER';
COMMENT ON COLUMN tmp.addr_obj.objectguid IS 'Глобальный уникальный идентификатор адресного объекта типа UUID';
COMMENT ON COLUMN tmp.addr_obj.changeid IS 'ID изменившей транзакции';
COMMENT ON COLUMN tmp.addr_obj.name IS 'Наименование';
COMMENT ON COLUMN tmp.addr_obj.typename IS 'Краткое наименование типа объекта';
COMMENT ON COLUMN tmp.addr_obj.level IS 'Уровень адресного объекта';
COMMENT ON COLUMN tmp.addr_obj.opertypeid IS 'Статус действия над записью – причина появления записи';
COMMENT ON COLUMN tmp.addr_obj.previd IS 'Идентификатор записи связывания с предыдущей исторической записью';
COMMENT ON COLUMN tmp.addr_obj.nextid IS 'Идентификатор записи связывания с последующей исторической записью';
COMMENT ON COLUMN tmp.addr_obj.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.addr_obj.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.addr_obj.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.addr_obj.isactual IS 'Статус актуальности адресного объекта ФИАС';
COMMENT ON COLUMN tmp.addr_obj.isactive IS 'Признак действующего адресного объекта';

CREATE TABLE tmp.addr_obj_division (
                                       id BIGINT NOT NULL PRIMARY KEY,
                                       parentid BIGINT NOT NULL,
                                       childid BIGINT NOT NULL,
                                       changeid BIGINT NOT NULL
);
COMMENT ON TABLE tmp.addr_obj_division IS 'Сведения по операциям переподчинения';
COMMENT ON COLUMN tmp.addr_obj_division.id IS 'Уникальный идентификатор записи. Ключевое поле';
COMMENT ON COLUMN tmp.addr_obj_division.parentid IS 'Родительский ID';
COMMENT ON COLUMN tmp.addr_obj_division.childid IS 'Дочерний ID';
COMMENT ON COLUMN tmp.addr_obj_division.changeid IS 'ID изменившей транзакции';

CREATE TABLE tmp.addr_obj_types (
                                    id INTEGER NOT NULL PRIMARY KEY,
                                    level INTEGER NOT NULL,
                                    shortname VARCHAR NOT NULL,
                                    name VARCHAR NOT NULL,
                                    "desc" VARCHAR,
                                    updatedate DATE NOT NULL,
                                    startdate DATE NOT NULL,
                                    enddate DATE NOT NULL,
                                    isactive BOOLEAN NOT NULL
);
COMMENT ON TABLE tmp.addr_obj_types IS 'Сведения по типам адресных объектов';
COMMENT ON COLUMN tmp.addr_obj_types.id IS 'Идентификатор записи';
COMMENT ON COLUMN tmp.addr_obj_types.level IS 'Уровень адресного объекта';
COMMENT ON COLUMN tmp.addr_obj_types.shortname IS 'Краткое наименование типа объекта';
COMMENT ON COLUMN tmp.addr_obj_types.name IS 'Полное наименование типа объекта';
COMMENT ON COLUMN tmp.addr_obj_types.desc IS 'Описание';
COMMENT ON COLUMN tmp.addr_obj_types.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.addr_obj_types.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.addr_obj_types.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.addr_obj_types.isactive IS 'Статус активности';

CREATE TABLE tmp.adm_hierarchy (
                                   id BIGINT NOT NULL PRIMARY KEY,
                                   objectid BIGINT NOT NULL,
                                   parentobjid BIGINT,
                                   changeid BIGINT NOT NULL,
                                   regioncode VARCHAR,
                                   areacode VARCHAR,
                                   citycode VARCHAR,
                                   placecode VARCHAR,
                                   plancode VARCHAR,
                                   streetcode VARCHAR,
                                   previd BIGINT,
                                   nextid BIGINT,
                                   updatedate DATE NOT NULL,
                                   startdate DATE NOT NULL,
                                   enddate DATE NOT NULL,
                                   isactive INTEGER NOT NULL,
                                   path VARCHAR NOT NULL
);
COMMENT ON TABLE tmp.adm_hierarchy IS 'Сведения по иерархии в административном делении';
COMMENT ON COLUMN tmp.adm_hierarchy.id IS 'Уникальный идентификатор записи. Ключевое поле';
COMMENT ON COLUMN tmp.adm_hierarchy.objectid IS 'Глобальный уникальный идентификатор объекта';
COMMENT ON COLUMN tmp.adm_hierarchy.parentobjid IS 'Идентификатор родительского объекта';
COMMENT ON COLUMN tmp.adm_hierarchy.changeid IS 'ID изменившей транзакции';
COMMENT ON COLUMN tmp.adm_hierarchy.regioncode IS 'Код региона';
COMMENT ON COLUMN tmp.adm_hierarchy.areacode IS 'Код района';
COMMENT ON COLUMN tmp.adm_hierarchy.citycode IS 'Код города';
COMMENT ON COLUMN tmp.adm_hierarchy.placecode IS 'Код населенного пункта';
COMMENT ON COLUMN tmp.adm_hierarchy.plancode IS 'Код ЭПС';
COMMENT ON COLUMN tmp.adm_hierarchy.streetcode IS 'Код улицы';
COMMENT ON COLUMN tmp.adm_hierarchy.previd IS 'Идентификатор записи связывания с предыдущей исторической записью';
COMMENT ON COLUMN tmp.adm_hierarchy.nextid IS 'Идентификатор записи связывания с последующей исторической записью';
COMMENT ON COLUMN tmp.adm_hierarchy.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.adm_hierarchy.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.adm_hierarchy.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.adm_hierarchy.isactive IS 'Признак действующего адресного объекта';
COMMENT ON COLUMN tmp.adm_hierarchy.path IS 'Материализованный путь к объекту (полная иерархия)';

CREATE TABLE tmp.apartments (
                                id BIGINT NOT NULL PRIMARY KEY,
                                objectid BIGINT NOT NULL,
                                objectguid VARCHAR NOT NULL,
                                changeid BIGINT NOT NULL,
                                number VARCHAR NOT NULL,
                                aparttype INTEGER NOT NULL,
                                opertypeid BIGINT NOT NULL,
                                previd BIGINT,
                                nextid BIGINT,
                                updatedate DATE NOT NULL,
                                startdate DATE NOT NULL,
                                enddate DATE NOT NULL,
                                isactual INTEGER NOT NULL,
                                isactive INTEGER NOT NULL
);
COMMENT ON TABLE tmp.apartments IS 'Сведения по помещениям';
COMMENT ON COLUMN tmp.apartments.id IS 'Уникальный идентификатор записи. Ключевое поле';
COMMENT ON COLUMN tmp.apartments.objectid IS 'Глобальный уникальный идентификатор объекта типа INTEGER';
COMMENT ON COLUMN tmp.apartments.objectguid IS 'Глобальный уникальный идентификатор адресного объекта типа UUID';
COMMENT ON COLUMN tmp.apartments.changeid IS 'ID изменившей транзакции';
COMMENT ON COLUMN tmp.apartments.number IS 'Номер комнаты';
COMMENT ON COLUMN tmp.apartments.aparttype IS 'Тип комнаты';
COMMENT ON COLUMN tmp.apartments.opertypeid IS 'Статус действия над записью – причина появления записи';
COMMENT ON COLUMN tmp.apartments.previd IS 'Идентификатор записи связывания с предыдущей исторической записью';
COMMENT ON COLUMN tmp.apartments.nextid IS 'Идентификатор записи связывания с последующей исторической записью';
COMMENT ON COLUMN tmp.apartments.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.apartments.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.apartments.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.apartments.isactual IS 'Статус актуальности адресного объекта ФИАС';
COMMENT ON COLUMN tmp.apartments.isactive IS 'Признак действующего адресного объекта';

CREATE TABLE tmp.apartment_types (
                                     id INTEGER NOT NULL PRIMARY KEY,
                                     name VARCHAR NOT NULL,
                                     shortname VARCHAR,
                                     "desc" VARCHAR,
                                     updatedate DATE NOT NULL,
                                     startdate DATE NOT NULL,
                                     enddate DATE NOT NULL,
                                     isactive BOOLEAN NOT NULL
);
COMMENT ON TABLE tmp.apartment_types IS 'Сведения по типам помещений';
COMMENT ON COLUMN tmp.apartment_types.id IS 'Идентификатор типа (ключ)';
COMMENT ON COLUMN tmp.apartment_types.name IS 'Наименование';
COMMENT ON COLUMN tmp.apartment_types.shortname IS 'Краткое наименование';
COMMENT ON COLUMN tmp.apartment_types.desc IS 'Описание';
COMMENT ON COLUMN tmp.apartment_types.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.apartment_types.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.apartment_types.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.apartment_types.isactive IS 'Статус активности';

CREATE TABLE tmp.carplaces (
                               id BIGINT NOT NULL PRIMARY KEY,
                               objectid BIGINT NOT NULL,
                               objectguid VARCHAR NOT NULL,
                               changeid BIGINT NOT NULL,
                               number VARCHAR NOT NULL,
                               opertypeid INTEGER NOT NULL,
                               previd BIGINT,
                               nextid BIGINT,
                               updatedate DATE NOT NULL,
                               startdate DATE NOT NULL,
                               enddate DATE NOT NULL,
                               isactual INTEGER NOT NULL,
                               isactive INTEGER NOT NULL
);
COMMENT ON TABLE tmp.carplaces IS 'Сведения по машино-местам';
COMMENT ON COLUMN tmp.carplaces.id IS 'Уникальный идентификатор записи. Ключевое поле';
COMMENT ON COLUMN tmp.carplaces.objectid IS 'Глобальный уникальный идентификатор объекта типа INTEGER';
COMMENT ON COLUMN tmp.carplaces.objectguid IS 'Глобальный уникальный идентификатор адресного объекта типа UUID';
COMMENT ON COLUMN tmp.carplaces.changeid IS 'ID изменившей транзакции';
COMMENT ON COLUMN tmp.carplaces.number IS 'Номер машиноместа';
COMMENT ON COLUMN tmp.carplaces.opertypeid IS 'Статус действия над записью – причина появления записи';
COMMENT ON COLUMN tmp.carplaces.previd IS 'Идентификатор записи связывания с предыдущей исторической записью';
COMMENT ON COLUMN tmp.carplaces.nextid IS 'Идентификатор записи связывания с последующей исторической записью';
COMMENT ON COLUMN tmp.carplaces.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.carplaces.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.carplaces.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.carplaces.isactual IS 'Статус актуальности адресного объекта ФИАС';
COMMENT ON COLUMN tmp.carplaces.isactive IS 'Признак действующего адресного объекта';

CREATE TABLE tmp.change_history (
                                    changeid BIGINT NOT NULL,
                                    objectid BIGINT NOT NULL,
                                    adrobjectid VARCHAR NOT NULL,
                                    opertypeid INTEGER NOT NULL,
                                    ndocid BIGINT,
                                    changedate DATE NOT NULL
);
COMMENT ON TABLE tmp.change_history IS 'Сведения по истории изменений';
COMMENT ON COLUMN tmp.change_history.changeid IS 'ID изменившей транзакции';
COMMENT ON COLUMN tmp.change_history.objectid IS 'Уникальный ID объекта';
COMMENT ON COLUMN tmp.change_history.adrobjectid IS 'Уникальный ID изменившей транзакции (GUID)';
COMMENT ON COLUMN tmp.change_history.opertypeid IS 'Тип операции';
COMMENT ON COLUMN tmp.change_history.ndocid IS 'ID документа';
COMMENT ON COLUMN tmp.change_history.changedate IS 'Дата изменения';

CREATE TABLE tmp.houses (
                            id BIGINT NOT NULL PRIMARY KEY,
                            objectid BIGINT NOT NULL,
                            objectguid VARCHAR NOT NULL,
                            changeid BIGINT NOT NULL,
                            housenum VARCHAR,
                            addnum1 VARCHAR,
                            addnum2 VARCHAR,
                            housetype INTEGER,
                            addtype1 INTEGER,
                            addtype2 INTEGER,
                            opertypeid INTEGER NOT NULL,
                            previd BIGINT,
                            nextid BIGINT,
                            updatedate DATE NOT NULL,
                            startdate DATE NOT NULL,
                            enddate DATE NOT NULL,
                            isactual INTEGER NOT NULL,
                            isactive INTEGER NOT NULL
);
COMMENT ON TABLE tmp.houses IS 'Сведения по номерам домов улиц городов и населенных пунктов';
COMMENT ON COLUMN tmp.houses.id IS 'Уникальный идентификатор записи. Ключевое поле';
COMMENT ON COLUMN tmp.houses.objectid IS 'Глобальный уникальный идентификатор объекта типа INTEGER';
COMMENT ON COLUMN tmp.houses.objectguid IS 'Глобальный уникальный идентификатор адресного объекта типа UUID';
COMMENT ON COLUMN tmp.houses.changeid IS 'ID изменившей транзакции';
COMMENT ON COLUMN tmp.houses.housenum IS 'Основной номер дома';
COMMENT ON COLUMN tmp.houses.addnum1 IS 'Дополнительный номер дома 1';
COMMENT ON COLUMN tmp.houses.addnum2 IS 'Дополнительный номер дома 1';
COMMENT ON COLUMN tmp.houses.housetype IS 'Основной тип дома';
COMMENT ON COLUMN tmp.houses.addtype1 IS 'Дополнительный тип дома 1';
COMMENT ON COLUMN tmp.houses.addtype2 IS 'Дополнительный тип дома 2';
COMMENT ON COLUMN tmp.houses.opertypeid IS 'Статус действия над записью – причина появления записи';
COMMENT ON COLUMN tmp.houses.previd IS 'Идентификатор записи связывания с предыдущей исторической записью';
COMMENT ON COLUMN tmp.houses.nextid IS 'Идентификатор записи связывания с последующей исторической записью';
COMMENT ON COLUMN tmp.houses.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.houses.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.houses.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.houses.isactual IS 'Статус актуальности адресного объекта ФИАС';
COMMENT ON COLUMN tmp.houses.isactive IS 'Признак действующего адресного объекта';

CREATE TABLE tmp.house_types (
                                 id INTEGER NOT NULL PRIMARY KEY,
                                 name VARCHAR NOT NULL,
                                 shortname VARCHAR,
                                 "desc" VARCHAR,
                                 updatedate DATE NOT NULL,
                                 startdate DATE NOT NULL,
                                 enddate DATE NOT NULL,
                                 isactive BOOLEAN NOT NULL
);
COMMENT ON TABLE tmp.house_types IS 'Сведения по типам домов';
COMMENT ON COLUMN tmp.house_types.id IS 'Идентификатор';
COMMENT ON COLUMN tmp.house_types.name IS 'Наименование';
COMMENT ON COLUMN tmp.house_types.shortname IS 'Краткое наименование';
COMMENT ON COLUMN tmp.house_types.desc IS 'Описание';
COMMENT ON COLUMN tmp.house_types.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.house_types.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.house_types.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.house_types.isactive IS 'Статус активности';

CREATE TABLE tmp.mun_hierarchy (
                                   id BIGINT NOT NULL PRIMARY KEY,
                                   objectid BIGINT NOT NULL,
                                   parentobjid BIGINT,
                                   changeid BIGINT NOT NULL,
                                   oktmo VARCHAR,
                                   previd BIGINT,
                                   nextid BIGINT,
                                   updatedate DATE NOT NULL,
                                   startdate DATE NOT NULL,
                                   enddate DATE NOT NULL,
                                   isactive INTEGER NOT NULL,
                                   path VARCHAR NOT NULL
);
COMMENT ON TABLE tmp.mun_hierarchy IS 'Сведения по иерархии в муниципальном делении';
COMMENT ON COLUMN tmp.mun_hierarchy.id IS 'Уникальный идентификатор записи. Ключевое поле';
COMMENT ON COLUMN tmp.mun_hierarchy.objectid IS 'Глобальный уникальный идентификатор адресного объекта ';
COMMENT ON COLUMN tmp.mun_hierarchy.parentobjid IS 'Идентификатор родительского объекта';
COMMENT ON COLUMN tmp.mun_hierarchy.changeid IS 'ID изменившей транзакции';
COMMENT ON COLUMN tmp.mun_hierarchy.oktmo IS 'Код ОКТМО';
COMMENT ON COLUMN tmp.mun_hierarchy.previd IS 'Идентификатор записи связывания с предыдущей исторической записью';
COMMENT ON COLUMN tmp.mun_hierarchy.nextid IS 'Идентификатор записи связывания с последующей исторической записью';
COMMENT ON COLUMN tmp.mun_hierarchy.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.mun_hierarchy.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.mun_hierarchy.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.mun_hierarchy.isactive IS 'Признак действующего адресного объекта';
COMMENT ON COLUMN tmp.mun_hierarchy.path IS 'Материализованный путь к объекту (полная иерархия)';

CREATE TABLE tmp.normative_docs (
                                    id BIGINT NOT NULL PRIMARY KEY,
                                    name VARCHAR NOT NULL,
                                    date DATE NOT NULL,
                                    number VARCHAR NOT NULL,
                                    type INTEGER NOT NULL,
                                    kind INTEGER NOT NULL,
                                    updatedate DATE NOT NULL,
                                    orgname VARCHAR,
                                    regnum VARCHAR,
                                    regdate DATE,
                                    accdate DATE,
                                    comment VARCHAR
);
COMMENT ON TABLE tmp.normative_docs IS 'Сведения о нормативном документе, являющемся основанием присвоения адресному элементу наименования';
COMMENT ON COLUMN tmp.normative_docs.id IS 'Уникальный идентификатор документа';
COMMENT ON COLUMN tmp.normative_docs.name IS 'Наименование документа';
COMMENT ON COLUMN tmp.normative_docs.date IS 'Дата документа';
COMMENT ON COLUMN tmp.normative_docs.number IS 'Номер документа';
COMMENT ON COLUMN tmp.normative_docs.type IS 'Тип документа';
COMMENT ON COLUMN tmp.normative_docs.kind IS 'Вид документа';
COMMENT ON COLUMN tmp.normative_docs.updatedate IS 'Дата обновления';
COMMENT ON COLUMN tmp.normative_docs.orgname IS 'Наименование органа создвшего нормативный документ';
COMMENT ON COLUMN tmp.normative_docs.regnum IS 'Номер государственной регистрации';
COMMENT ON COLUMN tmp.normative_docs.regdate IS 'Дата государственной регистрации';
COMMENT ON COLUMN tmp.normative_docs.accdate IS 'Дата вступления в силу нормативного документа';
COMMENT ON COLUMN tmp.normative_docs.comment IS 'Комментарий';

CREATE TABLE tmp.normative_docs_kinds (
                                          id INTEGER NOT NULL PRIMARY KEY,
                                          name VARCHAR NOT NULL
);
COMMENT ON TABLE tmp.normative_docs_kinds IS '';
COMMENT ON COLUMN tmp.normative_docs_kinds.id IS 'Идентификатор записи';
COMMENT ON COLUMN tmp.normative_docs_kinds.name IS 'Наименование';

CREATE TABLE tmp.normative_docs_types (
                                          id INTEGER NOT NULL PRIMARY KEY,
                                          name VARCHAR NOT NULL,
                                          startdate DATE NOT NULL,
                                          enddate DATE NOT NULL
);
COMMENT ON TABLE tmp.normative_docs_types IS '';
COMMENT ON COLUMN tmp.normative_docs_types.id IS 'Идентификатор записи';
COMMENT ON COLUMN tmp.normative_docs_types.name IS 'Наименование';
COMMENT ON COLUMN tmp.normative_docs_types.startdate IS 'Дата начала действия записи';
COMMENT ON COLUMN tmp.normative_docs_types.enddate IS 'Дата окончания действия записи';

CREATE TABLE tmp.object_levels (
                                   level INTEGER NOT NULL PRIMARY KEY,
                                   name VARCHAR NOT NULL,
                                   shortname VARCHAR,
                                   updatedate DATE NOT NULL,
                                   startdate DATE NOT NULL,
                                   enddate DATE NOT NULL,
                                   isactive BOOLEAN NOT NULL
);
COMMENT ON TABLE tmp.object_levels IS 'Сведения по уровням адресных объектов';
COMMENT ON COLUMN tmp.object_levels.level IS 'Уникальный идентификатор записи. Ключевое поле. Номер уровня объекта';
COMMENT ON COLUMN tmp.object_levels.name IS 'Наименование';
COMMENT ON COLUMN tmp.object_levels.shortname IS 'Краткое наименование';
COMMENT ON COLUMN tmp.object_levels.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.object_levels.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.object_levels.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.object_levels.isactive IS 'Признак действующего адресного объекта';

CREATE TABLE tmp.operation_types (
                                     id INTEGER NOT NULL PRIMARY KEY,
                                     name VARCHAR NOT NULL,
                                     shortname VARCHAR,
                                     "desc" VARCHAR,
                                     updatedate DATE NOT NULL,
                                     startdate DATE NOT NULL,
                                     enddate DATE NOT NULL,
                                     isactive BOOLEAN NOT NULL
);
COMMENT ON TABLE tmp.operation_types IS 'Сведения по статусу действия';
COMMENT ON COLUMN tmp.operation_types.id IS 'Идентификатор статуса (ключ)';
COMMENT ON COLUMN tmp.operation_types.name IS 'Наименование';
COMMENT ON COLUMN tmp.operation_types.shortname IS 'Краткое наименование';
COMMENT ON COLUMN tmp.operation_types.desc IS 'Описание';
COMMENT ON COLUMN tmp.operation_types.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.operation_types.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.operation_types.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.operation_types.isactive IS 'Статус активности';

CREATE TABLE tmp.param (
                           id BIGINT NOT NULL PRIMARY KEY,
                           objectid BIGINT NOT NULL,
                           changeid BIGINT,
                           changeidend BIGINT NOT NULL,
                           typeid INTEGER NOT NULL,
                           value VARCHAR NOT NULL,
                           updatedate DATE NOT NULL,
                           startdate DATE NOT NULL,
                           enddate DATE NOT NULL
);
COMMENT ON TABLE tmp.param IS 'Сведения о классификаторе параметров адресообразующих элементов и объектов недвижимости ';
COMMENT ON COLUMN tmp.param.id IS 'Идентификатор записи';
COMMENT ON COLUMN tmp.param.objectid IS 'Глобальный уникальный идентификатор адресного объекта ';
COMMENT ON COLUMN tmp.param.changeid IS 'ID изменившей транзакции';
COMMENT ON COLUMN tmp.param.changeidend IS 'ID завершившей транзакции';
COMMENT ON COLUMN tmp.param.typeid IS 'Тип параметра';
COMMENT ON COLUMN tmp.param.value IS 'Значение параметра';
COMMENT ON COLUMN tmp.param.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.param.startdate IS 'Дата начала действия записи';
COMMENT ON COLUMN tmp.param.enddate IS 'Дата окончания действия записи';

CREATE TABLE tmp.param_types (
                                 id INTEGER NOT NULL PRIMARY KEY,
                                 name VARCHAR NOT NULL,
                                 code VARCHAR NOT NULL,
                                 "desc" VARCHAR,
                                 updatedate DATE NOT NULL,
                                 startdate DATE NOT NULL,
                                 enddate DATE NOT NULL,
                                 isactive BOOLEAN NOT NULL
);
COMMENT ON TABLE tmp.param_types IS 'Сведения по типу параметра';
COMMENT ON COLUMN tmp.param_types.id IS 'Идентификатор типа параметра (ключ)';
COMMENT ON COLUMN tmp.param_types.name IS 'Наименование';
COMMENT ON COLUMN tmp.param_types.code IS 'Краткое наименование';
COMMENT ON COLUMN tmp.param_types.desc IS 'Описание';
COMMENT ON COLUMN tmp.param_types.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.param_types.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.param_types.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.param_types.isactive IS 'Статус активности';

CREATE TABLE tmp.reestr_objects (
                                    objectid BIGINT NOT NULL PRIMARY KEY,
                                    createdate DATE NOT NULL,
                                    changeid BIGINT NOT NULL,
                                    levelid INTEGER NOT NULL,
                                    updatedate DATE NOT NULL,
                                    objectguid VARCHAR NOT NULL,
                                    isactive INTEGER NOT NULL
);
COMMENT ON TABLE tmp.reestr_objects IS 'Сведения об адресном элементе в части его идентификаторов';
COMMENT ON COLUMN tmp.reestr_objects.objectid IS 'Уникальный идентификатор объекта';
COMMENT ON COLUMN tmp.reestr_objects.createdate IS 'Дата создания';
COMMENT ON COLUMN tmp.reestr_objects.changeid IS 'ID изменившей транзакции';
COMMENT ON COLUMN tmp.reestr_objects.levelid IS 'Уровень объекта';
COMMENT ON COLUMN tmp.reestr_objects.updatedate IS 'Дата обновления';
COMMENT ON COLUMN tmp.reestr_objects.objectguid IS 'GUID объекта';
COMMENT ON COLUMN tmp.reestr_objects.isactive IS 'Признак действующего объекта (1 - действующий, 0 - не действующий)';

CREATE TABLE tmp.rooms (
                           id BIGINT NOT NULL PRIMARY KEY,
                           objectid BIGINT NOT NULL,
                           objectguid VARCHAR NOT NULL,
                           changeid BIGINT NOT NULL,
                           number VARCHAR NOT NULL,
                           roomtype INTEGER NOT NULL,
                           opertypeid INTEGER NOT NULL,
                           previd BIGINT,
                           nextid BIGINT,
                           updatedate DATE NOT NULL,
                           startdate DATE NOT NULL,
                           enddate DATE NOT NULL,
                           isactual INTEGER NOT NULL,
                           isactive INTEGER NOT NULL
);
COMMENT ON TABLE tmp.rooms IS 'Сведения по комнатам';
COMMENT ON COLUMN tmp.rooms.id IS 'Уникальный идентификатор записи. Ключевое поле';
COMMENT ON COLUMN tmp.rooms.objectid IS 'Глобальный уникальный идентификатор объекта типа INTEGER';
COMMENT ON COLUMN tmp.rooms.objectguid IS 'Глобальный уникальный идентификатор адресного объекта типа UUID';
COMMENT ON COLUMN tmp.rooms.changeid IS 'ID изменившей транзакции';
COMMENT ON COLUMN tmp.rooms.number IS 'Номер комнаты или офиса';
COMMENT ON COLUMN tmp.rooms.roomtype IS 'Тип комнаты или офиса';
COMMENT ON COLUMN tmp.rooms.opertypeid IS 'Статус действия над записью – причина появления записи';
COMMENT ON COLUMN tmp.rooms.previd IS 'Идентификатор записи связывания с предыдущей исторической записью';
COMMENT ON COLUMN tmp.rooms.nextid IS 'Идентификатор записи связывания с последующей исторической записью';
COMMENT ON COLUMN tmp.rooms.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.rooms.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.rooms.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.rooms.isactual IS 'Статус актуальности адресного объекта ФИАС';
COMMENT ON COLUMN tmp.rooms.isactive IS 'Признак действующего адресного объекта';

CREATE TABLE tmp.room_types (
                                id INTEGER NOT NULL PRIMARY KEY,
                                name VARCHAR NOT NULL,
                                shortname VARCHAR,
                                "desc" VARCHAR,
                                updatedate DATE NOT NULL,
                                startdate DATE NOT NULL,
                                enddate DATE NOT NULL,
                                isactive BOOLEAN NOT NULL
);
COMMENT ON TABLE tmp.room_types IS 'Сведения по типам комнат';
COMMENT ON COLUMN tmp.room_types.id IS 'Идентификатор типа (ключ)';
COMMENT ON COLUMN tmp.room_types.name IS 'Наименование';
COMMENT ON COLUMN tmp.room_types.shortname IS 'Краткое наименование';
COMMENT ON COLUMN tmp.room_types.desc IS 'Описание';
COMMENT ON COLUMN tmp.room_types.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.room_types.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.room_types.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.room_types.isactive IS 'Статус активности';

CREATE TABLE tmp.steads (
                            id INTEGER NOT NULL PRIMARY KEY,
                            objectid INTEGER NOT NULL,
                            objectguid VARCHAR NOT NULL,
                            changeid INTEGER NOT NULL,
                            number VARCHAR NOT NULL,
                            opertypeid VARCHAR NOT NULL,
                            previd INTEGER,
                            nextid INTEGER,
                            updatedate DATE NOT NULL,
                            startdate DATE NOT NULL,
                            enddate DATE NOT NULL,
                            isactual INTEGER NOT NULL,
                            isactive INTEGER NOT NULL
);
COMMENT ON TABLE tmp.steads IS 'Сведения по земельным участкам';
COMMENT ON COLUMN tmp.steads.id IS 'Уникальный идентификатор записи. Ключевое поле';
COMMENT ON COLUMN tmp.steads.objectid IS 'Глобальный уникальный идентификатор объекта типа INTEGER';
COMMENT ON COLUMN tmp.steads.objectguid IS 'Глобальный уникальный идентификатор адресного объекта типа UUID';
COMMENT ON COLUMN tmp.steads.changeid IS 'ID изменившей транзакции';
COMMENT ON COLUMN tmp.steads.number IS 'Номер земельного участка';
COMMENT ON COLUMN tmp.steads.opertypeid IS 'Статус действия над записью – причина появления записи';
COMMENT ON COLUMN tmp.steads.previd IS 'Идентификатор записи связывания с предыдущей исторической записью';
COMMENT ON COLUMN tmp.steads.nextid IS 'Идентификатор записи связывания с последующей исторической записью';
COMMENT ON COLUMN tmp.steads.updatedate IS 'Дата внесения (обновления) записи';
COMMENT ON COLUMN tmp.steads.startdate IS 'Начало действия записи';
COMMENT ON COLUMN tmp.steads.enddate IS 'Окончание действия записи';
COMMENT ON COLUMN tmp.steads.isactual IS 'Статус актуальности адресного объекта ФИАС';
COMMENT ON COLUMN tmp.steads.isactive IS 'Признак действующего адресного объекта';
