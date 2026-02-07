package sqlbuilder

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type SchemaBuilder struct {
	table         string
	fullTable     string
	primaryKey    string
	ignoreNotNull bool
	driver        Driver
}

func NewSchemaBuilder(dbSchema string, tableName string, ignoreNotNull bool, driver Driver) *SchemaBuilder {
	return &SchemaBuilder{
		table:         tableName,
		fullTable:     buildFullTableName(dbSchema, tableName),
		primaryKey:    resolvePrimaryKey(tableName),
		ignoreNotNull: ignoreNotNull,
		driver:        driver,
	}
}

func (b *SchemaBuilder) Build(data []byte) (string, error) {
	var schema schema
	if err := xml.Unmarshal(data, &schema); err != nil {
		return "", err
	}

	if len(schema.Element) == 0 {
		return "", fmt.Errorf("invalid schema attrs for '%s'", b.table)
	}

	var attrs []attribute
	var descr string

	// (hardcoded)
	if b.table == "normative_docs_kinds" || b.table == "normative_docs_types" {
		attrs = schema.Element[1].ComplexType.Attributes
		descr = schema.Element[0].ComplexType.Sequence.Elements[0].Annotation.Documentation
	} else {
		attrs = schema.Element[0].ComplexType.Sequence.Elements[0].ComplexType.Attributes
		descr = schema.Element[0].Annotation.Documentation
	}

	if len(attrs) == 0 {
		return "", fmt.Errorf("empty attrs for '%s'", b.table)
	}

	res := fmt.Sprintf("CREATE TABLE %s (\n%s\n);\n%s;\n%s",
		b.fullTable,
		b.buildColumns(attrs),
		b.driver.BuildTableComment(b.fullTable, descr),
		b.driver.BuildColumnComments(b.fullTable, attrs),
	)
	return res, nil
}

func (b *SchemaBuilder) buildColumn(attr attribute) string {
	var sb strings.Builder
	columnName := resolveColumnName(attr.Name)

	sb.WriteString(escapeColumnName(columnName))
	sb.WriteString(" ")

	sb.WriteString(b.driver.ResolveColumnType(attr.Type))

	if !b.ignoreNotNull {
		notNull := resolveNullability(b.table, columnName, attr)
		if notNull != "" {
			sb.WriteString(" " + notNull)
		}
	}

	if columnName == b.primaryKey {
		sb.WriteString(" PRIMARY KEY")
	}

	return sb.String()
}

type schema struct {
	Element []element `xml:"element"`
}

type element struct {
	Name        string      `xml:"name,attr"`
	ComplexType complexType `xml:"complexType"`
	Annotation  annotation  `xml:"annotation"`
}

type complexType struct {
	Sequence   sequence    `xml:"sequence"`
	Attributes []attribute `xml:"attribute"`
}

type annotation struct {
	Documentation string `xml:"documentation"`
}

type sequence struct {
	Elements []element `xml:"element"`
}

type attribute struct {
	Name       string     `xml:"name,attr"`
	Type       string     `xml:"type,attr"`
	Use        string     `xml:"use,attr"`
	Annotation annotation `xml:"annotation"`
}
