package sqlbuilder

import (
	"fmt"
	"strings"
)

type MySQLSchemaBuilder struct {
	table            string
	fullTable        string
	primaryKey       string
	ignoreRequired   bool
	ignorePrimaryKey bool
}

func NewMySQLSchemaBuilder(dbSchema string, tableName string, ignoreRequired bool, ignorePrimaryKey bool) *MySQLSchemaBuilder {
	return &MySQLSchemaBuilder{
		table:            tableName,
		fullTable:        buildFullTableName(dbSchema, tableName),
		primaryKey:       resolvePrimaryKey(tableName),
		ignoreRequired:   ignoreRequired,
		ignorePrimaryKey: ignorePrimaryKey,
	}
}

func (b *MySQLSchemaBuilder) Build(data []byte) (string, error) {
	attrs, descr, err := ResolveSchemaAttrs(b.table, data)
	if err != nil {
		return "", err
	}
	res := fmt.Sprintf("CREATE TABLE %s (\n%s\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;\n%s", b.fullTable, b.buildColumns(attrs), b.buildTableComment(descr))
	return res, nil
}

func (b *MySQLSchemaBuilder) BuildPrimaryKey() (string, error) {
	return fmt.Sprintf("ALTER TABLE %s ADD PRIMARY KEY(%s);", b.fullTable, b.primaryKey), nil
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
