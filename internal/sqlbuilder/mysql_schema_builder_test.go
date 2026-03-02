package sqlbuilder

import (
	"testing"
)

func TestMySQLSchemaBuilder_Build(t *testing.T) {
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
						<xs:attribute name="NAME" use="required">
							<xs:annotation>
								<xs:documentation>Наименование адресного объекта</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:string">
									<xs:maxLength value="250"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="TYPENAME" use="required">
							<xs:annotation>
								<xs:documentation>Тип (сокращение) адресного объекта</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:string">
									<xs:maxLength value="50"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="LEVEL" use="required">
							<xs:annotation>
								<xs:documentation>Уровень адресного объекта</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:int">
									<xs:totalDigits value="2"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="OPERATIONTYPEID" use="required">
							<xs:annotation>
								<xs:documentation>Тип операции</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:int"/>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="PREVID">
							<xs:annotation>
								<xs:documentation>ID предыдущей записи объекта</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:long">
									<xs:totalDigits value="19"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="NEXTID">
							<xs:annotation>
								<xs:documentation>ID следующей записи объекта</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:long">
									<xs:totalDigits value="19"/>
								</xs:restriction>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="UPDATEDATE" use="required">
							<xs:annotation>
								<xs:documentation>Дата внесения (обновления) записи</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:date"/>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="STARTDATE" use="required">
							<xs:annotation>
								<xs:documentation>Начало действия записи</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:date"/>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="ENDDATE" use="required">
							<xs:annotation>
								<xs:documentation>Окончание действия записи</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:date"/>
							</xs:simpleType>
						</xs:attribute>
						<xs:attribute name="ISACTIVE" use="required">
							<xs:annotation>
								<xs:documentation>Статус актуальности записи</xs:documentation>
							</xs:annotation>
							<xs:simpleType>
								<xs:restriction base="xs:boolean"/>
							</xs:simpleType>
						</xs:attribute>
					</xs:complexType>
				</xs:element>
			</xs:sequence>
		</xs:complexType>
	</xs:element>
</xs:dbSchema>`

		builder := NewMySQLSchemaBuilder("", "addressobjects", false)
		result, err := builder.Build([]byte(xmlData))
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := `CREATE TABLE addressobjects (
	id VARCHAR(500) NOT NULL PRIMARY KEY,
	objectid VARCHAR(500) NOT NULL,
	objectguid VARCHAR(500) NOT NULL,
	name VARCHAR(500) NOT NULL,
	typename VARCHAR(500) NOT NULL,
	level VARCHAR(500) NOT NULL,
	operationtypeid VARCHAR(500) NOT NULL,
	previd VARCHAR(500),
	nextid VARCHAR(500),
	updatedate VARCHAR(500) NOT NULL,
	startdate VARCHAR(500) NOT NULL,
	enddate VARCHAR(500) NOT NULL,
	isactive VARCHAR(500) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
ALTER TABLE addressobjects COMMENT = 'Состав и структура файла со сведениями классификатора адресообразующих элементов БД ФИАС';`

		if result != expected {
			t.Errorf("Build() result =\n%s\n\nexpected =\n%s", result, expected)
		}
	})

	t.Run("with schema", func(t *testing.T) {
		xmlData := `<?xml version="1.0" encoding="utf-8"?>
<xs:dbSchema xmlns:xs="http://www.w3.org/2001/XMLSchema">
	<xs:element name="TEST">
		<xs:annotation>
			<xs:documentation>Test table</xs:documentation>
		</xs:annotation>
		<xs:complexType>
			<xs:sequence>
				<xs:element name="OBJECT">
					<xs:complexType>
						<xs:attribute name="ID" use="required" type="xs:string">
							<xs:annotation>
								<xs:documentation>Test ID</xs:documentation>
							</xs:annotation>
						</xs:attribute>
					</xs:complexType>
				</xs:element>
			</xs:sequence>
		</xs:complexType>
	</xs:element>
</xs:dbSchema>`

		builder := NewMySQLSchemaBuilder("test_schema", "test", false)
		result, err := builder.Build([]byte(xmlData))
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := `CREATE TABLE test_schema.test (
	id VARCHAR(500) NOT NULL PRIMARY KEY
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
ALTER TABLE test_schema.test COMMENT = 'Test table';`

		if result != expected {
			t.Errorf("Build() result =\n%s\n\nexpected =\n%s", result, expected)
		}
	})

	t.Run("ignore not null", func(t *testing.T) {
		xmlData := `<?xml version="1.0" encoding="utf-8"?>
<xs:dbSchema xmlns:xs="http://www.w3.org/2001/XMLSchema">
	<xs:element name="TEST">
		<xs:annotation>
			<xs:documentation>Test table</xs:documentation>
		</xs:annotation>
		<xs:complexType>
			<xs:sequence>
				<xs:element name="OBJECT">
					<xs:complexType>
						<xs:attribute name="ID" use="required" type="xs:string">
							<xs:annotation>
								<xs:documentation>Test ID</xs:documentation>
							</xs:annotation>
						</xs:attribute>
						<xs:attribute name="NAME" use="required" type="xs:string">
							<xs:annotation>
								<xs:documentation>Test Name</xs:documentation>
							</xs:annotation>
						</xs:attribute>
					</xs:complexType>
				</xs:element>
			</xs:sequence>
		</xs:complexType>
	</xs:element>
</xs:dbSchema>`

		builder := NewMySQLSchemaBuilder("", "test", true) // ignoreNotNull = true
		result, err := builder.Build([]byte(xmlData))
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := `CREATE TABLE test (
	id VARCHAR(500) PRIMARY KEY,
	name VARCHAR(500)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
ALTER TABLE test COMMENT = 'Test table';`

		if result != expected {
			t.Errorf("Build() result =\n%s\n\nexpected =\n%s", result, expected)
		}
	})
}

func TestXsdTypeToMySQL(t *testing.T) {
	tests := []struct {
		xsdType  string
		expected string
	}{
		{"xs:string", "VARCHAR(500)"},
		{"xs:int", "INT"},
		{"xs:long", "BIGINT"},
		{"xs:boolean", "BOOLEAN"},
		{"xs:date", "DATE"},
		{"xs:dateTime", "DATETIME"},
		{"unknown", "VARCHAR(500)"},
	}

	for _, tt := range tests {
		t.Run(tt.xsdType, func(t *testing.T) {
			result := xsdTypeToMySQL(tt.xsdType)
			if result != tt.expected {
				t.Errorf("xsdTypeToMySQL(%s) = %s, expected %s", tt.xsdType, result, tt.expected)
			}
		})
	}
}
