package sqlbuilder

import (
	"fmt"
	"strings"
)

type PostgreSQLSchemaBuilder struct {
	table            string
	fullTable        string
	primaryKey       string
	ignoreRequired   bool
	ignorePrimaryKey bool
}

func NewPostgreSQLSchemaBuilder(dbSchema string, tableName string, ignoreRequired bool, ignorePrimaryKey bool) *PostgreSQLSchemaBuilder {
	return &PostgreSQLSchemaBuilder{
		table:            tableName,
		fullTable:        buildFullTableName(dbSchema, tableName),
		primaryKey:       resolvePrimaryKey(tableName),
		ignoreRequired:   ignoreRequired,
		ignorePrimaryKey: ignorePrimaryKey,
	}
}

func (b *PostgreSQLSchemaBuilder) Build(data []byte) (string, error) {
	attrs, descr, err := ResolveSchemaAttrs(b.table, data)
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf("CREATE TABLE %s (\n%s\n);\n%s;\n%s", b.fullTable, b.buildColumns(attrs), b.buildTableComment(descr), b.buildColumnComments(attrs))
	return res, nil
}

func (b *PostgreSQLSchemaBuilder) BuildPrimaryKey() (string, error) {
	return fmt.Sprintf("ALTER TABLE %s ADD PRIMARY KEY(%s);", b.fullTable, b.primaryKey), nil
}

func (b *PostgreSQLSchemaBuilder) buildColumns(attrs []attribute) string {
	columns := make([]string, len(attrs))
	for i, attr := range attrs {
		columns[i] = fmt.Sprintf("\t%s", b.buildColumn(attr))
	}

	return strings.Join(columns, ",\n")
}

//nolint:dupl
func (b *PostgreSQLSchemaBuilder) buildColumn(attr attribute) string {
	var sb strings.Builder
	columnName := resolveColumnName(attr.Name)

	sb.WriteString(escapeColumnNamePostgreSQL(columnName))
	sb.WriteString(" ")

	sb.WriteString(xsdTypeToSQL(attr.Type))

	if !b.ignoreRequired {
		notNull := resolveNullability(b.table, columnName, attr)
		if notNull != "" {
			sb.WriteString(" " + notNull)
		}
	}

	if !b.ignorePrimaryKey && columnName == b.primaryKey {
		sb.WriteString(" PRIMARY KEY")
	}

	return sb.String()
}

func (b *PostgreSQLSchemaBuilder) buildTableComment(descr string) string {
	return fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", b.fullTable, descr)
}

func (b *PostgreSQLSchemaBuilder) buildColumnComments(attrs []attribute) string {
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
