package sqlbuilder

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type SchemaBuilder struct {
	table      string
	fullTable  string
	primaryKey string
}

func NewSchemaBuilder(dbSchema string, tableName string) *SchemaBuilder {
	return &SchemaBuilder{
		table:      tableName,
		fullTable:  buildFullTableName(dbSchema, tableName),
		primaryKey: resolvePrimaryKey(tableName),
	}
}

func (b *SchemaBuilder) Build(data []byte) (string, error) {
	var schema schema
	if err := xml.Unmarshal(data, &schema); err != nil {
		panic(err)
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

	res := fmt.Sprintf("CREATE TABLE %s (\n%s\n);\n%s;\n%s", b.fullTable, b.buildColumns(attrs), b.buildTableComment(descr), b.buildColumnComments(attrs))
	return res, nil
}

func (b *SchemaBuilder) buildColumns(attrs []attribute) string {
	columns := make([]string, len(attrs))
	for i, attr := range attrs {
		columns[i] = fmt.Sprintf("\t%s", b.buildColumn(attr))
	}

	return strings.Join(columns, ",\n")
}

func (b *SchemaBuilder) buildColumn(attr attribute) string {
	var sb strings.Builder
	columnName := resolveColumnName(attr.Name)

	sb.WriteString(escapeColumnName(columnName))
	sb.WriteString(" ")

	sb.WriteString(xsdTypeToSQL(attr.Type))
	if attr.Use == "required" {
		sb.WriteString(" NOT NULL")
	}

	if columnName == b.primaryKey {
		sb.WriteString(" PRIMARY KEY")
	}

	return sb.String()
}

func (b *SchemaBuilder) buildTableComment(descr string) string {
	return fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", b.fullTable, descr)
}

func (b *SchemaBuilder) buildColumnComments(attrs []attribute) string {
	columns := make([]string, len(attrs))
	for i, attr := range attrs {
		columnName := resolveColumnName(attr.Name)
		columns[i] = fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';", b.fullTable, columnName, attr.Annotation.Documentation)
	}

	return strings.Join(columns, "\n")
}

func xsdTypeToSQL(xsdType string) string {
	switch xsdType {
	case "xs:string":
		return "VARCHAR"
	case "xs:int":
		return "INT"
	case "xs:long":
		return "BIGINT"
	case "xs:boolean":
		return "BOOLEAN"
	case "xs:date":
		return "DATE"
	case "xs:dateTime":
		return "TIMESTAMP"
	default:
		return "VARCHAR"
	}
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
