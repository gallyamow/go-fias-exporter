package sqlbuilder

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type MySQLSchemaBuilder struct {
	table         string
	fullTable     string
	primaryKey    string
	ignoreNotNull bool
}

func NewMySQLSchemaBuilder(dbSchema string, tableName string, ignoreNotNull bool) *MySQLSchemaBuilder {
	return &MySQLSchemaBuilder{
		table:         tableName,
		fullTable:     buildFullTableName(dbSchema, tableName),
		primaryKey:    resolvePrimaryKey(tableName),
		ignoreNotNull: ignoreNotNull,
	}
}

func (b *MySQLSchemaBuilder) Build(data []byte) (string, error) {
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

	res := fmt.Sprintf("CREATE TABLE %s (\n%s\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;\n%s", b.fullTable, b.buildColumns(attrs), b.buildTableComment(descr))
	return res, nil
}

func (b *MySQLSchemaBuilder) buildColumns(attrs []attribute) string {
	columns := make([]string, len(attrs))
	for i, attr := range attrs {
		columns[i] = fmt.Sprintf("\t%s", b.buildColumn(attr))
	}

	return strings.Join(columns, ",\n")
}

//nolint:dupl
func (b *MySQLSchemaBuilder) buildColumn(attr attribute) string {
	var sb strings.Builder
	columnName := resolveColumnName(attr.Name)

	sb.WriteString(escapeColumnNameMySQL(columnName))
	sb.WriteString(" ")

	sb.WriteString(xsdTypeToMySQL(attr.Type))

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

func (b *MySQLSchemaBuilder) buildTableComment(descr string) string {
	return fmt.Sprintf("ALTER TABLE %s COMMENT = '%s';", b.fullTable, descr)
}

func xsdTypeToMySQL(xsdType string) string {
	switch xsdType {
	case "xs:string":
		return "VARCHAR(500)"
	case "xs:int":
		return "INT"
	case "xs:long":
		return "BIGINT"
	case "xs:boolean":
		return "BOOLEAN"
	case "xs:date":
		return "DATE"
	case "xs:dateTime":
		return "DATETIME"
	default:
		return "VARCHAR(500)"
	}
}
