package sqlbuilder

import "testing"

func TestSchemaBuilder_Build(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		xmlData := `<?xml version="1.0" encoding="utf-8"?>
<!-- edited with XMLSpy v2011 rel. 2 (http://www.altova.com) by TeaM DJiNN (TeaM DJiNN) -->
<xs:dbSchema xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:sch="http://purl.oclc.org/dsdl/schematron" xmlns:usch="http://www.unisoftware.ru/schematron-extensions" xmlns:sql="urn:schemas-microsoft-com:mapping-dbSchema" elementFormDefault="qualified" attributeFormDefault="unqualified">
	<xs:element name="ADDRESSOBJECTS">
		<xs:annotation>
			<xs:documentation>Состав и структура файла со сведениями классификатора адресообразующих элементов БД ФИАС</xs:documentation>
		</xs:annotation>
		<xs:complexType>
			<xs:sequence>
				<xs:element name="OBJECT" maxOccurs="unbounded">
					<xs:annotation>
						<xs:documentation>Сведения классификатора адресообразующих элементов</xs:documentation>
					</xs:annotation>
					<xs:complexType>
						<xs:attribute name="ID" use="required">
							<xs:annotation>
								<xs:documentation>Уникальный идентификатор записи. Ключевое поле</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:long">
									<xs:totalDigits value="19"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="OBJECTID" use="required">
							<xs:annotation>
								<xs:documentation>Глобальный уникальный идентификатор адресного объекта типа INTEGER</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:long">
									<xs:totalDigits value="19"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="OBJECTGUID" use="required">
							<xs:annotation>
								<xs:documentation>Глобальный уникальный идентификатор адресного объекта типа UUID</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:string">
									<xs:length value="36"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="CHANGEID" use="required">
							<xs:annotation>
								<xs:documentation>ID изменившей транзакции</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:long">
									<xs:totalDigits value="19"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="NAME" use="required">
							<xs:annotation>
								<xs:documentation>Наименование</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:string">
									<xs:minLength value="1"/>
									<xs:maxLength value="250"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="TYPENAME" use="required">
							<xs:annotation>
								<xs:documentation>Краткое наименование типа объекта</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:string">
									<xs:minLength value="1"/>
									<xs:maxLength value="50"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="LEVEL" use="required">
							<xs:annotation>
								<xs:documentation>Уровень адресного объекта</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:string">
									<xs:minLength value="1"/>
									<xs:maxLength value="10"/>
									<xs:pattern value="[0-9]{1,10}"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="OPERTYPEID" use="required">
							<xs:annotation>
								<xs:documentation>Статус действия над записью – причина появления записи</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:integer"/>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="PREVID" use="optional">
							<xs:annotation>
								<xs:documentation>Идентификатор записи связывания с предыдущей исторической записью</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:long">
									<xs:totalDigits value="19"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="NEXTID" use="optional">
							<xs:annotation>
								<xs:documentation>Идентификатор записи связывания с последующей исторической записью</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:long">
									<xs:totalDigits value="19"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="UPDATEDATE" type="xs:date" use="required">
							<xs:annotation>
								<xs:documentation>Дата внесения (обновления) записи</xs:documentation>
							</xs:annotation>
						</xs:attribute>
						<xs:attribute name="STARTDATE" type="xs:date" use="required">
							<xs:annotation>
								<xs:documentation>Начало действия записи</xs:documentation>
							</xs:annotation>
						</xs:attribute>
						<xs:attribute name="ENDDATE" type="xs:date" use="required">
							<xs:annotation>
								<xs:documentation>Окончание действия записи</xs:documentation>
							</xs:annotation>
						</xs:attribute>
						<xs:attribute name="ISACTUAL" use="required">
							<xs:annotation>
								<xs:documentation>Статус актуальности адресного объекта ФИАС</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:integer">
									<xs:enumeration value="0"/>
									<xs:enumeration value="1"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="ISACTIVE" use="required">
							<xs:annotation>
								<xs:documentation>Признак действующего адресного объекта</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:integer">
									<xs:enumeration value="0"/>
									<xs:enumeration value="1"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
					</xs:complexType>
				</xs:element>
			</xs:sequence>
		</xs:complexType>
	</xs:element>
</xs:dbSchema>
`

		builder := NewSchemaBuilder("tmp", "addr_obj")

		got, err := builder.Build([]byte(xmlData))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `CREATE TABLE tmp.addr_obj (
	id VARCHAR NOT NULL PRIMARY KEY,
	objectid VARCHAR NOT NULL,
	objectguid VARCHAR NOT NULL,
	changeid VARCHAR NOT NULL,
	name VARCHAR NOT NULL,
	typename VARCHAR NOT NULL,
	level VARCHAR NOT NULL,
	opertypeid VARCHAR NOT NULL,
	previd VARCHAR,
	nextid VARCHAR,
	updatedate DATE NOT NULL,
	startdate DATE NOT NULL,
	enddate DATE NOT NULL,
	isactual VARCHAR NOT NULL,
	isactive VARCHAR NOT NULL
);
COMMENT ON TABLE tmp.addr_obj IS 'Состав и структура файла со сведениями классификатора адресообразующих элементов БД ФИАС';
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
COMMENT ON COLUMN tmp.addr_obj.isactive IS 'Признак действующего адресного объекта';`

		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("normative_docs_kinds", func(t *testing.T) {
		xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<!-- edited with XMLSpy v2011 rel. 2 (http://www.altova.com) by TeaM DJiNN (TeaM DJiNN) -->
<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
	<xs:element name="NDOCKINDS">
		<xs:annotation>
			<xs:documentation>Состав и структура файла со сведениями по видам нормативных документов</xs:documentation>
		</xs:annotation>
		<xs:complexType>
			<xs:sequence>
				<xs:element ref="NDOCKIND" maxOccurs="unbounded">
					<xs:annotation>
						<xs:documentation>Сведения по видам нормативных документов</xs:documentation>
					</xs:annotation>
				</xs:element>
			</xs:sequence>
		</xs:complexType>
	</xs:element>
	<xs:element name="NDOCKIND">
		<xs:complexType>
			<xs:attribute name="ID" type="xs:integer" use="required">
				<xs:annotation>
					<xs:documentation>Идентификатор записи</xs:documentation>
				</xs:annotation>
			</xs:attribute>
			<xs:attribute name="NAME" use="required">
				<xs:annotation>
					<xs:documentation>Наименование</xs:documentation>
				</xs:annotation>
				<xs:simpleType>
					<xs:restriction base="xs:string">
						<xs:maxLength value="500"/>
						<xs:minLength value="1"/>
					</xs:restriction>
				</xs:simpleType>
			</xs:attribute>
		</xs:complexType>
	</xs:element>
</xs:schema>
`

		builder := NewSchemaBuilder("tmp", "normative_docs_kinds")

		got, err := builder.Build([]byte(xmlData))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `CREATE TABLE tmp.normative_docs_kinds (
	id VARCHAR NOT NULL PRIMARY KEY,
	name VARCHAR NOT NULL
);
COMMENT ON TABLE tmp.normative_docs_kinds IS 'Сведения по видам нормативных документов';
COMMENT ON COLUMN tmp.normative_docs_kinds.id IS 'Идентификатор записи';
COMMENT ON COLUMN tmp.normative_docs_kinds.name IS 'Наименование';`

		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("normative_docs_types", func(t *testing.T) {
		xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<!-- edited with XMLSpy v2011 rel. 2 (http://www.altova.com) by TeaM DJiNN (TeaM DJiNN) -->
<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
	<xs:element name="NDOCTYPES">
		<xs:annotation>
			<xs:documentation>Состав и структура файла со сведениями по типам нормативных документов</xs:documentation>
		</xs:annotation>
		<xs:complexType>
			<xs:sequence>
				<xs:element ref="NDOCTYPE" maxOccurs="unbounded">
					<xs:annotation>
						<xs:documentation>Сведения по типам нормативных документов</xs:documentation>
					</xs:annotation>
				</xs:element>
			</xs:sequence>
		</xs:complexType>
	</xs:element>
	<xs:element name="NDOCTYPE">
		<xs:complexType>
			<xs:attribute name="ID" type="xs:integer" use="required">
				<xs:annotation>
					<xs:documentation>Идентификатор записи</xs:documentation>
				</xs:annotation>
			</xs:attribute>
			<xs:attribute name="NAME" use="required">
				<xs:annotation>
					<xs:documentation>Наименование</xs:documentation>
				</xs:annotation>
				<xs:simpleType>
					<xs:restriction base="xs:string">
						<xs:maxLength value="500"/>
						<xs:minLength value="1"/>
					</xs:restriction>
				</xs:simpleType>
			</xs:attribute>
			<xs:attribute name="STARTDATE" use="required">
				<xs:annotation>
					<xs:documentation>Дата начала действия записи</xs:documentation>
				</xs:annotation>
				<xs:simpleType>
					<xs:restriction base="xs:date"/>
				</xs:simpleType>
			</xs:attribute>
			<xs:attribute name="ENDDATE" use="required">
				<xs:annotation>
					<xs:documentation>Дата окончания действия записи</xs:documentation>
				</xs:annotation>
				<xs:simpleType>
					<xs:restriction base="xs:date"/>
				</xs:simpleType>
			</xs:attribute>
		</xs:complexType>
	</xs:element>
</xs:schema>
`

		builder := NewSchemaBuilder("tmp", "normative_docs_types")

		got, err := builder.Build([]byte(xmlData))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `CREATE TABLE tmp.normative_docs_types (
	id VARCHAR NOT NULL PRIMARY KEY,
	name VARCHAR NOT NULL,
	startdate VARCHAR NOT NULL,
	enddate VARCHAR NOT NULL
);
COMMENT ON TABLE tmp.normative_docs_types IS 'Сведения по типам нормативных документов';
COMMENT ON COLUMN tmp.normative_docs_types.id IS 'Идентификатор записи';
COMMENT ON COLUMN tmp.normative_docs_types.name IS 'Наименование';
COMMENT ON COLUMN tmp.normative_docs_types.startdate IS 'Дата начала действия записи';
COMMENT ON COLUMN tmp.normative_docs_types.enddate IS 'Дата окончания действия записи';`

		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})
}
